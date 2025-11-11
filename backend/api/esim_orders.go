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

	// 转换订单数据格式
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
	orderDetail, err := h.orderService.GetOrderWithDetail(ctx, uint(orderID), userID)
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
		"order_id":          orderDetail.OrderID,
		"order_no":          orderDetail.OrderNo,
		"user_id":           orderDetail.UserID,
		"product_id":        orderDetail.ProductID,
		"product_name":      orderDetail.ProductName,
		"quantity":          orderDetail.Quantity,
		"unit_price":        orderDetail.UnitPrice,
		"total_amount":      orderDetail.TotalAmount,
		"status":            orderDetail.Status,
		"provider_order_id": orderDetail.ProviderOrderID,
		"order_items":       orderDetail.OrderItems,
		"esims":             orderDetail.Esims,
		"created_at":        orderDetail.CreatedAt,
		"updated_at":        orderDetail.UpdatedAt,
		"completed_at":      orderDetail.CompletedAt,
	}

	h.sendSuccess(w, response)
}
