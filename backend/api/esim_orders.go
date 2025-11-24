package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"tg-robot-sim/services"
	"tg-robot-sim/storage/models"
)

// handleEsimOrders 处理 eSIM 订单相关请求
func (h *MiniAppApiService) handleEsimOrders(w http.ResponseWriter, r *http.Request) {
	// 获取用户 ID
	userID, err := h.getUserIDFromContext(r)
	if err != nil || userID == 0 {
		h.sendError(w, http.StatusUnauthorized, "Unauthorized", "Invalid user ID")
		return
	}

	switch r.Method {
	case http.MethodPost:
		// 创建 eSIM 订单
		h.handleCreateEsimOrder(w, r, userID)
	case http.MethodGet:
		// 获取用户 eSIM 订单列表
		h.handleGetEsimOrders(w, r, userID)
	default:
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
	}
}

// handleCreateEsimOrder 处理创建 eSIM 订单请求
func (h *MiniAppApiService) handleCreateEsimOrder(w http.ResponseWriter, r *http.Request, userID int64) {
	ctx := r.Context()

	// 解析请求体
	var req services.CreateEsimOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// 设置用户ID
	req.UserID = userID

	// 验证必填字段
	if req.ProductID == 0 {
		h.sendError(w, http.StatusBadRequest, "Product ID is required", "")
		return
	}
	if req.Quantity <= 0 {
		h.sendError(w, http.StatusBadRequest, "Quantity must be greater than 0", "")
		return
	}
	if req.TotalAmount == "" {
		h.sendError(w, http.StatusBadRequest, "Total amount is required", "")
		return
	}
	if req.CustomerEmail == "" {
		h.sendError(w, http.StatusBadRequest, "Customer email is required", "")
		return
	}

	// 创建 eSIM 订单
	order, err := h.orderService.CreateEsimOrder(ctx, &req)
	if err != nil {
		// 根据错误类型返回不同的错误码
		errMsg := err.Error()
		if strings.Contains(errMsg, "余额不足") {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeInsufficientBalance, errMsg, "")
		} else if strings.Contains(errMsg, "产品不存在") {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeProductNotFound, errMsg, "")
		} else if strings.Contains(errMsg, "产品暂不可用") {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeProductUnavailable, errMsg, "")
		} else if strings.Contains(errMsg, "订单金额不匹配") {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeInvalidAmount, errMsg, "")
		} else if strings.Contains(errMsg, "邮箱") {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeInvalidRequest, errMsg, "")
		} else {
			h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "创建订单失败", errMsg)
		}
		return
	}

	// 返回订单信息
	response := map[string]interface{}{
		"order_id":     order.OrderID,
		"order_no":     order.OrderNo,
		"status":       order.Status,
		"total_amount": order.TotalAmount,
		"created_at":   order.CreatedAt,
	}

	h.sendSuccess(w, response)
}

// handleGetEsimOrders 处理获取用户 eSIM 订单列表请求
func (h *MiniAppApiService) handleGetEsimOrders(w http.ResponseWriter, r *http.Request, userID int64) {
	ctx := r.Context()

	// 获取查询参数
	limit := h.parseIntParam(r, "limit", 20)
	offset := h.parseIntParam(r, "offset", 0)
	statusParam := r.URL.Query().Get("status")

	// 构建筛选条件
	var statusFilter models.OrderStatus
	if statusParam != "" {
		statusFilter = models.OrderStatus(statusParam)
	}

	// 获取订单列表（使用带筛选的方法）
	orders, total, err := h.orderService.GetUserOrdersWithFilters(ctx, userID, statusFilter, limit, offset)
	if err != nil {
		h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取订单列表失败", err.Error())
		return
	}

	// 收集订单 IDs
	orderIDs := make([]uint, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	// ordersMap := make(map[uint]*models.Order)

	// 批量获取订单详情
	orders, err = h.orderService.GetOrderByIDs(ctx, orderIDs)
	if err != nil {
		orders = make([]*models.Order, 0)
	}

	// 转换订单数据格式，合并详情
	var orderList []map[string]interface{}
	for _, order := range orders {
		orderData := map[string]interface{}{
			"order_id":          order.ID,
			"order_no":          order.OrderNo,
			"product_id":        order.ProductID,
			"product_name":      order.ProductName,
			"quantity":          order.Quantity,
			"unit_price":        order.UnitPrice,
			"total_amount":      order.Amount,
			"status":            order.Status,
			"provider_order_id": order.ProviderOrderID,
			"provider_order_no": order.ProviderOrderNo,
			"created_at":        order.CreatedAt,
			"updated_at":        order.UpdatedAt,
			"completed_at":      order.CompletedAt,
		}

		orderList = append(orderList, orderData)
	}

	// 返回响应
	h.sendSuccess(w, map[string]interface{}{
		"orders": orderList,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// handleEsimOrderDetail 处理 eSIM 订单详情请求
func (h *MiniAppApiService) handleEsimOrderDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := r.Context()

	// 获取用户 ID
	userID, err := h.getUserIDFromContext(r)
	if err != nil || userID == 0 {
		h.sendError(w, http.StatusUnauthorized, "Unauthorized", "Invalid user ID")
		return
	}

	// 从 URL 路径提取订单 ID
	// /api/miniapp/esim/orders/123
	path := r.URL.Path
	orderIDStr := strings.TrimPrefix(path, "/api/miniapp/esim/orders/")
	if orderIDStr == "" || orderIDStr == path {
		h.sendError(w, http.StatusBadRequest, "Order ID is required", "")
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	// 获取订单详情
	orderDetail, err := h.orderService.GetUserOrderByID(ctx, uint(orderID), userID)
	if err != nil {
		if strings.Contains(err.Error(), "订单不存在") {
			h.sendErrorWithCode(w, http.StatusNotFound, ErrCodeOrderNotFound, err.Error(), "")
		} else if strings.Contains(err.Error(), "无权访问") {
			h.sendError(w, http.StatusForbidden, "Access denied", "")
		} else {
			h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取订单详情失败", err.Error())
		}
		return
	}

	// 转换为前端期望的格式
	response := map[string]interface{}{
		"order_id":          orderDetail.ID,
		"order_no":          orderDetail.OrderNo,
		"user_id":           orderDetail.UserID,
		"product_id":        orderDetail.ProductID,
		"product_name":      orderDetail.ProductName,
		"quantity":          orderDetail.Quantity,
		"unit_price":        orderDetail.UnitPrice,
		"total_amount":      orderDetail.Amount,
		"status":            orderDetail.Status,
		"provider_order_id": orderDetail.ProviderOrderID,
		"provider_order_no": orderDetail.ProviderOrderNo,
		"created_at":        orderDetail.CreatedAt,
		"updated_at":        orderDetail.UpdatedAt,
		"completed_at":      orderDetail.CompletedAt,
	}

	h.sendSuccess(w, response)
}

// handleEsimCards 处理 eSIM 卡相关请求
func (h *MiniAppApiService) handleEsimCards(w http.ResponseWriter, r *http.Request) {
	// 获取用户 ID
	userID, err := h.getUserIDFromContext(r)
	if err != nil || userID == 0 {
		h.sendError(w, http.StatusUnauthorized, "Unauthorized", "Invalid user ID")
		return
	}

	// 根据 URL 路径分发请求
	path := r.URL.Path
	if strings.Contains(path, "/sync") {
		// 同步 eSIM 卡状态
		h.handleSyncEsimCard(w, r, userID)
	} else if strings.HasSuffix(path, "/") || !strings.Contains(strings.TrimPrefix(path, "/api/miniapp/esim/cards/"), "/") {
		// 列表或详情
		switch r.Method {
		case http.MethodGet:
			if strings.Contains(path, "/api/miniapp/esim/cards/") && !strings.HasSuffix(path, "/") {
				// 获取单个 eSIM 卡详情
				h.handleGetEsimCardDetail(w, r, userID)
			} else {
				// 获取 eSIM 卡列表
				h.handleGetEsimCardList(w, r, userID)
			}
		default:
			h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		}
	} else {
		h.sendError(w, http.StatusNotFound, "Not found", "")
	}
}

// handleGetEsimCardList 处理获取 eSIM 卡列表请求
func (h *MiniAppApiService) handleGetEsimCardList(w http.ResponseWriter, r *http.Request, userID int64) {
	ctx := r.Context()

	// 获取查询参数
	limit := h.parseIntParam(r, "limit", 20)
	offset := h.parseIntParam(r, "offset", 0)
	statusParam := r.URL.Query().Get("status")

	// 验证分页参数
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 构建筛选条件
	filters := services.EsimCardFilters{
		Status: models.EsimStatus(statusParam),
		Limit:  limit,
		Offset: offset,
	}

	// 获取 eSIM 卡列表
	esimCards, total, err := h.esimCardService.GetUserEsimCards(ctx, userID, filters)
	if err != nil {
		h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取 eSIM 卡列表失败", err.Error())
		return
	}

	// 转换数据格式
	var cardList []map[string]interface{}
	for _, card := range esimCards {
		cardData := map[string]interface{}{
			"id":             card.ID,
			"iccid":          card.ICCID,
			"status":         card.Status,
			"data_used":      card.DataUsed,
			"data_remaining": card.DataRemaining,
			"usage_percent":  card.UsagePercent,
			"activated_at":   card.ActivatedAt,
			"expires_at":     card.ExpiresAt,
			"last_sync_at":   card.LastSyncAt,
			"created_at":     card.CreatedAt,
		}
		cardList = append(cardList, cardData)
	}

	// 返回响应
	h.sendSuccess(w, map[string]interface{}{
		"esim_cards": cardList,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	})
}

// handleGetEsimCardDetail 处理获取 eSIM 卡详情请求
func (h *MiniAppApiService) handleGetEsimCardDetail(w http.ResponseWriter, r *http.Request, userID int64) {
	ctx := r.Context()

	// 从 URL 路径提取 eSIM 卡 ID
	path := r.URL.Path
	esimIDStr := strings.TrimPrefix(path, "/api/miniapp/esim/cards/")
	esimIDStr = strings.TrimSuffix(esimIDStr, "/")
	if esimIDStr == "" {
		h.sendError(w, http.StatusBadRequest, "eSIM card ID is required", "")
		return
	}

	esimID, err := strconv.ParseUint(esimIDStr, 10, 32)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid eSIM card ID", err.Error())
		return
	}

	// 获取 eSIM 卡详情及关联订单
	esimCardWithOrder, err := h.esimCardService.GetEsimCardWithOrder(ctx, uint(esimID), userID)
	if err != nil {
		if strings.Contains(err.Error(), "eSIM 卡不存在") {
			h.sendErrorWithCode(w, http.StatusNotFound, ErrCodeNotFound, err.Error(), "")
		} else if strings.Contains(err.Error(), "无权访问") {
			h.sendError(w, http.StatusForbidden, "Access denied", "")
		} else {
			h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取 eSIM 卡详情失败", err.Error())
		}
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"id":               esimCardWithOrder.EsimCard.ID,
		"iccid":            esimCardWithOrder.EsimCard.ICCID,
		"activation_code":  esimCardWithOrder.EsimCard.ActivationCode,
		"qr_code":          esimCardWithOrder.EsimCard.QrCode,
		"lpa":              esimCardWithOrder.EsimCard.Lpa,
		"direct_apple_url": esimCardWithOrder.EsimCard.DirectAppleUrl,
		"apn_type":         esimCardWithOrder.EsimCard.ApnType,
		"is_roaming":       esimCardWithOrder.EsimCard.IsRoaming,
		"status":           esimCardWithOrder.EsimCard.Status,
		"data_used":        esimCardWithOrder.EsimCard.DataUsed,
		"data_remaining":   esimCardWithOrder.EsimCard.DataRemaining,
		"usage_percent":    esimCardWithOrder.EsimCard.UsagePercent,
		"activated_at":     esimCardWithOrder.EsimCard.ActivatedAt,
		"expires_at":       esimCardWithOrder.EsimCard.ExpiresAt,
		"last_sync_at":     esimCardWithOrder.EsimCard.LastSyncAt,
		"created_at":       esimCardWithOrder.EsimCard.CreatedAt,
	}

	// 添加关联的购买订单信息
	if esimCardWithOrder.PurchaseOrder != nil {
		response["purchase_order"] = map[string]interface{}{
			"order_id":   esimCardWithOrder.PurchaseOrder.ID,
			"order_no":   esimCardWithOrder.PurchaseOrder.OrderNo,
			"product_id": esimCardWithOrder.PurchaseOrder.ProductID,
			"amount":     esimCardWithOrder.PurchaseOrder.Amount,
			"status":     esimCardWithOrder.PurchaseOrder.Status,
			"created_at": esimCardWithOrder.PurchaseOrder.CreatedAt,
		}
	}

	h.sendSuccess(w, response)
}

// handleSyncEsimCard 处理同步 eSIM 卡状态请求
func (h *MiniAppApiService) handleSyncEsimCard(w http.ResponseWriter, r *http.Request, userID int64) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := r.Context()

	// 从 URL 路径提取 eSIM 卡 ID
	path := r.URL.Path
	esimIDStr := strings.TrimPrefix(path, "/api/miniapp/esim/cards/")
	esimIDStr = strings.TrimSuffix(esimIDStr, "/sync")
	if esimIDStr == "" {
		h.sendError(w, http.StatusBadRequest, "eSIM card ID is required", "")
		return
	}

	esimID, err := strconv.ParseUint(esimIDStr, 10, 32)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid eSIM card ID", err.Error())
		return
	}

	// 验证用户权限（获取 eSIM 卡检查所有权）
	_, err = h.esimCardService.GetEsimCard(ctx, uint(esimID), userID)
	if err != nil {
		if strings.Contains(err.Error(), "eSIM 卡不存在") {
			h.sendErrorWithCode(w, http.StatusNotFound, ErrCodeNotFound, err.Error(), "")
		} else if strings.Contains(err.Error(), "无权访问") {
			h.sendError(w, http.StatusForbidden, "Access denied", "")
		} else {
			h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取 eSIM 卡失败", err.Error())
		}
		return
	}

	// 同步 eSIM 卡状态
	err = h.esimCardService.SyncEsimCardStatus(ctx, uint(esimID))
	if err != nil {
		h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "同步 eSIM 卡状态失败", err.Error())
		return
	}

	// 获取更新后的 eSIM 卡信息
	updatedCard, err := h.esimCardService.GetEsimCard(ctx, uint(esimID), userID)
	if err != nil {
		h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取更新后的 eSIM 卡信息失败", err.Error())
		return
	}

	// 返回更新后的信息
	response := map[string]interface{}{
		"id":             updatedCard.ID,
		"iccid":          updatedCard.ICCID,
		"status":         updatedCard.Status,
		"data_used":      updatedCard.DataUsed,
		"data_remaining": updatedCard.DataRemaining,
		"usage_percent":  updatedCard.UsagePercent,
		"activated_at":   updatedCard.ActivatedAt,
		"expires_at":     updatedCard.ExpiresAt,
		"last_sync_at":   updatedCard.LastSyncAt,
		"updated_at":     updatedCard.UpdatedAt,
	}

	h.sendSuccess(w, response)
}
