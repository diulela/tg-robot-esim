package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"tg-robot-sim/pkg/telegram"
	"tg-robot-sim/services"
	"tg-robot-sim/storage/models"
)

// MiniAppHandler Mini App API 处理器
type MiniAppApiService struct {
	productService       services.ProductService
	walletService        services.WalletService
	orderService         services.OrderService
	walletHistoryService services.WalletHistoryService
	rechargeService      services.RechargeService
}

// NewMiniAppHandler 创建 Mini App 处理器实例
func NewMiniAppApiService(
	productService services.ProductService,
	walletService services.WalletService,
	orderService services.OrderService,
	walletHistoryService services.WalletHistoryService,
	rechargeService services.RechargeService,
) *MiniAppApiService {
	return &MiniAppApiService{
		productService:       productService,
		walletService:        walletService,
		orderService:         orderService,
		walletHistoryService: walletHistoryService,
		rechargeService:      rechargeService,
	}
}

// APIResponse 统一的 API 响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// 错误码定义
const (
	ErrCodeInvalidAmount   = 40001 // 充值金额低于最小限额
	ErrCodeInvalidFormat   = 40002 // 充值金额格式错误
	ErrCodeOrderNotFound   = 40003 // 订单不存在
	ErrCodeOrderExpired    = 40004 // 订单已过期
	ErrCodeOrderCompleted  = 40005 // 订单已完成，不能重复处理
	ErrCodeWalletNotFound  = 40006 // 用户钱包不存在
	ErrCodeGenerateAmount  = 50001 // 生成精确金额失败
	ErrCodeBlockchainQuery = 50002 // 区块链查询失败
	ErrCodeDatabaseError   = 50003 // 数据库操作失败
	ErrCodeBalanceUpdate   = 50004 // 余额更新失败
)

// sendJSON 发送 JSON 响应
func (h *MiniAppApiService) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendSuccess 发送成功响应
func (h *MiniAppApiService) sendSuccess(w http.ResponseWriter, data interface{}) {
	response := APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	h.sendJSON(w, http.StatusOK, response)
}

// sendError 发送错误响应
func (h *MiniAppApiService) sendError(w http.ResponseWriter, statusCode int, message string, details string) {
	response := ErrorResponse{
		Code:    statusCode,
		Message: message,
		Details: details,
	}
	h.sendJSON(w, statusCode, response)
}

// sendErrorWithCode 发送带自定义错误码的错误响应
func (h *MiniAppApiService) sendErrorWithCode(w http.ResponseWriter, statusCode int, errorCode int, message string, details string) {
	response := ErrorResponse{
		Code:    errorCode,
		Message: message,
		Details: details,
	}
	h.sendJSON(w, statusCode, response)
}

// getUserIDFromContext 从上下文获取用户ID
func (h *MiniAppApiService) getUserIDFromContext(r *http.Request) (int64, error) {
	// 从 Telegram Web App 初始化数据中提取用户ID
	initData := r.Header.Get("X-Telegram-Init-Data")
	if initData == "" {
		// 开发模式：从查询参数获取
		initData = r.URL.Query().Get("init_data")
	}

	// 如果有初始化数据，解析用户ID
	if initData != "" {
		userID, err := telegram.GetUserID(initData)
		if err == nil && userID > 0 {
			return userID, nil
		}
	}

	// 开发模式：从查询参数获取用户ID
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return 0, nil
}

// parseIntParam 解析整数参数
func (h *MiniAppApiService) parseIntParam(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// RegisterRoutes 注册路由
func (h *MiniAppApiService) RegisterRoutes(mux *http.ServeMux) {
	// 产品相关
	mux.HandleFunc("/api/miniapp/products", h.handleProducts)
	mux.HandleFunc("/api/miniapp/products/", h.handleProductDetail)

	// 钱包相关
	mux.HandleFunc("/api/miniapp/wallet/balance", h.handleWalletBalance)

	// 充值相关
	mux.HandleFunc("/api/miniapp/wallet/recharge", h.handleCreateRecharge)
	mux.HandleFunc("/api/miniapp/wallet/recharge/", h.handleRechargeDetail)
	mux.HandleFunc("/api/miniapp/wallet/recharge/history", h.handleRechargeHistory)

	// 订单相关
	mux.HandleFunc("/api/miniapp/purchase", h.handlePurchase)
	mux.HandleFunc("/api/miniapp/orders", h.handleOrders)

	// 钱包历史相关
	mux.HandleFunc("/api/miniapp/wallet/history", h.handleWalletHistory)
	mux.HandleFunc("/api/miniapp/wallet/history/stats", h.handleWalletHistoryStats)
	mux.HandleFunc("/api/miniapp/wallet/history/", h.handleHistoryRecord)
}

// handleProducts 处理产品列表请求
func (h *MiniAppApiService) handleProducts(w http.ResponseWriter, r *http.Request) {

	fmt.Println("====获取产品列表 ==========")

	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := r.Context()

	// 获取查询参数
	productType := r.URL.Query().Get("type")
	country := r.URL.Query().Get("country")
	search := r.URL.Query().Get("search")
	limit := h.parseIntParam(r, "limit", 20)
	offset := h.parseIntParam(r, "offset", 0)

	// 构建筛选条件
	filters := services.ProductFilters{
		Type:    productType,
		Country: country,
		Search:  search,
		Limit:   limit,
		Offset:  offset,
	}

	// 获取产品列表
	products, err := h.productService.GetProducts(ctx, filters)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	// 获取总数
	total, _ := h.productService.CountProducts(ctx, filters)

	// 返回响应
	h.sendSuccess(w, map[string]interface{}{
		"products": products,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// handleProductDetail 处理产品详情请求
func (h *MiniAppApiService) handleProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := r.Context()

	// 从 URL 路径提取产品 ID
	// /api/miniapp/products/123
	path := r.URL.Path
	idStr := path[len("/api/miniapp/products/"):]
	if idStr == "" {
		h.sendError(w, http.StatusBadRequest, "Product ID is required", "")
		return
	}

	productID, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	// 获取产品详情
	product, err := h.productService.GetProductByID(ctx, productID)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Product not found", err.Error())
		return
	}

	h.sendSuccess(w, product)
}

// handleWalletBalance 处理钱包余额请求
func (h *MiniAppApiService) handleWalletBalance(w http.ResponseWriter, r *http.Request) {
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

	// 获取钱包信息
	wallet, err := h.walletService.GetWallet(ctx, userID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get balance", err.Error())
		return
	}

	// 转换为前端期望的格式
	response := map[string]interface{}{
		"id":            fmt.Sprintf("%d", wallet.ID),
		"userId":        fmt.Sprintf("%d", wallet.UserID),
		"balance":       parseFloat(wallet.Balance),
		"currency":      "USDT",
		"frozenAmount":  parseFloat(wallet.FrozenBalance),
		"totalRecharge": parseFloat(wallet.TotalIncome),
		"totalSpent":    parseFloat(wallet.TotalExpense),
		"createdAt":     wallet.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updatedAt":     wallet.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	h.sendSuccess(w, response)
}

// parseFloat 将字符串转换为浮点数
func parseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return value
}

// handleCreateRecharge 处理创建充值订单请求
func (h *MiniAppApiService) handleCreateRecharge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// 解析请求体
	var req struct {
		Amount string `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// 验证金额格式
	if req.Amount == "" {
		h.sendError(w, http.StatusBadRequest, "Amount is required", "")
		return
	}

	// 创建充值订单
	order, err := h.rechargeService.CreateRechargeOrder(ctx, userID, req.Amount)
	if err != nil {
		// 根据错误类型返回不同的错误码
		errMsg := err.Error()
		if errMsg == "充值金额格式错误" {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeInvalidFormat, errMsg, "")
		} else if strings.Contains(errMsg, "充值金额不能低于") || strings.Contains(errMsg, "充值金额不能超过") {
			h.sendErrorWithCode(w, http.StatusBadRequest, ErrCodeInvalidAmount, errMsg, "")
		} else if strings.Contains(errMsg, "生成精确金额失败") {
			h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeGenerateAmount, "生成充值订单失败", errMsg)
		} else {
			h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "创建充值订单失败", errMsg)
		}
		return
	}

	// 返回订单信息
	response := map[string]interface{}{
		"order_no":       order.OrderNo,
		"amount":         order.Amount,
		"exact_amount":   order.ExactAmount,
		"wallet_address": order.WalletAddress,
		"status":         order.Status,
		"expires_at":     order.ExpiresAt,
		"created_at":     order.CreatedAt,
	}

	h.sendSuccess(w, response)
}

// handlePurchase 处理购买请求
func (h *MiniAppApiService) handlePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// 解析请求体
	var req struct {
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// 创建订单
	order, err := h.orderService.CreateOrder(ctx, userID, req.ProductID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to create order", err.Error())
		return
	}

	// 处理支付
	if err := h.orderService.PayOrder(ctx, order.OrderNo); err != nil {
		h.sendError(w, http.StatusInternalServerError, "Payment failed", err.Error())
		return
	}

	// 重新获取订单以获取最新状态
	order, _ = h.orderService.GetOrderByOrderNo(ctx, order.OrderNo)

	h.sendSuccess(w, order)
}

// handleOrders 处理订单列表请求
func (h *MiniAppApiService) handleOrders(w http.ResponseWriter, r *http.Request) {
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

	// 获取查询参数
	limit := h.parseIntParam(r, "limit", 20)
	offset := h.parseIntParam(r, "offset", 0)

	// 获取订单列表
	orders, err := h.orderService.GetUserOrders(ctx, userID, limit, offset)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get orders", err.Error())
		return
	}

	// 获取订单统计
	stats, _ := h.orderService.GetOrderStats(ctx, userID)

	h.sendSuccess(w, map[string]interface{}{
		"orders": orders,
		"stats":  stats,
		"limit":  limit,
		"offset": offset,
	})
}

// handleWalletHistory 处理钱包历史记录请求
func (h *MiniAppApiService) handleWalletHistory(w http.ResponseWriter, r *http.Request) {
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

	// 构建筛选条件
	filters := services.WalletHistoryFilters{
		UserID:    userID,
		Type:      models.WalletHistoryType(r.URL.Query().Get("type")),
		Status:    models.WalletHistoryStatus(r.URL.Query().Get("status")),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
		Limit:     h.parseIntParam(r, "limit", 20),
		Offset:    h.parseIntParam(r, "offset", 0),
	}

	// 获取钱包历史记录
	histories, total, err := h.walletHistoryService.GetWalletHistory(ctx, userID, filters)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get wallet history", err.Error())
		return
	}

	// 返回响应
	response := map[string]interface{}{
		"records": histories,
		"total":   total,
		"limit":   filters.Limit,
		"offset":  filters.Offset,
	}

	h.sendSuccess(w, response)
}

// handleWalletHistoryStats 处理钱包历史统计请求
func (h *MiniAppApiService) handleWalletHistoryStats(w http.ResponseWriter, r *http.Request) {
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

	// 获取钱包历史统计
	stats, err := h.walletHistoryService.GetWalletHistoryStats(ctx, userID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get wallet stats", err.Error())
		return
	}

	h.sendSuccess(w, stats)
}

// handleHistoryRecord 处理单条历史记录详情请求
func (h *MiniAppApiService) handleHistoryRecord(w http.ResponseWriter, r *http.Request) {
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

	// 从 URL 路径提取记录 ID
	// /api/miniapp/wallet/history/123
	path := r.URL.Path
	recordIDStr := strings.TrimPrefix(path, "/api/miniapp/wallet/history/")
	if recordIDStr == "" || recordIDStr == path {
		h.sendError(w, http.StatusBadRequest, "Record ID is required", "")
		return
	}

	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid record ID", err.Error())
		return
	}

	// 获取历史记录详情
	record, err := h.walletHistoryService.GetHistoryRecord(ctx, uint(recordID), userID)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Record not found", err.Error())
		return
	}

	h.sendSuccess(w, record)
}

// handleRechargeDetail 处理充值订单详情和状态检查请求
func (h *MiniAppApiService) handleRechargeDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 获取用户 ID
	userID, err := h.getUserIDFromContext(r)
	if err != nil || userID == 0 {
		h.sendError(w, http.StatusUnauthorized, "Unauthorized", "Invalid user ID")
		return
	}

	// 从 URL 路径提取订单号
	// /api/miniapp/wallet/recharge/RCH1730800000001 或 /api/miniapp/wallet/recharge/RCH1730800000001/check
	path := r.URL.Path
	orderNo := ""
	isCheckRequest := false

	if path == "/api/miniapp/wallet/recharge/" {
		h.sendError(w, http.StatusBadRequest, "Order number is required", "")
		return
	}

	// 解析路径
	pathParts := strings.Split(strings.TrimPrefix(path, "/api/miniapp/wallet/recharge/"), "/")
	if len(pathParts) > 0 && pathParts[0] != "" {
		orderNo = pathParts[0]
		if len(pathParts) > 1 && pathParts[1] == "check" {
			isCheckRequest = true
		}
	}

	if orderNo == "" {
		h.sendError(w, http.StatusBadRequest, "Order number is required", "")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// 获取订单详情
		order, err := h.rechargeService.GetRechargeOrder(ctx, orderNo)
		if err != nil {
			if err.Error() == "订单不存在" {
				h.sendErrorWithCode(w, http.StatusNotFound, ErrCodeOrderNotFound, err.Error(), "")
			} else {
				h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "获取充值订单失败", err.Error())
			}
			return
		}

		// 检查订单是否属于当前用户
		if order.UserID != userID {
			h.sendError(w, http.StatusForbidden, "Access denied", "")
			return
		}

		// 返回订单详情
		response := map[string]interface{}{
			"order_no":       order.OrderNo,
			"amount":         order.Amount,
			"exact_amount":   order.ExactAmount,
			"wallet_address": order.WalletAddress,
			"status":         order.Status,
			"tx_hash":        order.TxHash,
			"confirmations":  order.Confirmations,
			"expires_at":     order.ExpiresAt,
			"confirmed_at":   order.ConfirmedAt,
			"created_at":     order.CreatedAt,
		}

		h.sendSuccess(w, response)

	case http.MethodPost:
		// 手动检查充值状态
		if !isCheckRequest {
			h.sendError(w, http.StatusBadRequest, "Invalid request path", "")
			return
		}

		order, err := h.rechargeService.CheckRechargeStatus(ctx, orderNo)
		if err != nil {
			if err.Error() == "订单不存在" {
				h.sendErrorWithCode(w, http.StatusNotFound, ErrCodeOrderNotFound, err.Error(), "")
			} else if strings.Contains(err.Error(), "区块链查询失败") {
				h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeBlockchainQuery, "查询充值状态失败", err.Error())
			} else {
				h.sendErrorWithCode(w, http.StatusInternalServerError, ErrCodeDatabaseError, "检查充值状态失败", err.Error())
			}
			return
		}

		// 检查订单是否属于当前用户
		if order.UserID != userID {
			h.sendError(w, http.StatusForbidden, "Access denied", "")
			return
		}

		// 返回更新后的订单状态
		response := map[string]interface{}{
			"order_no":      order.OrderNo,
			"status":        order.Status,
			"tx_hash":       order.TxHash,
			"confirmations": order.Confirmations,
			"confirmed_at":  order.ConfirmedAt,
		}

		h.sendSuccess(w, response)

	default:
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
	}
}

// handleRechargeHistory 处理充值历史请求
func (h *MiniAppApiService) handleRechargeHistory(w http.ResponseWriter, r *http.Request) {
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

	// 获取查询参数
	limit := h.parseIntParam(r, "limit", 20)
	offset := h.parseIntParam(r, "offset", 0)

	// 获取充值历史
	orders, total, err := h.rechargeService.GetUserRechargeHistory(ctx, userID, limit, offset)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get recharge history", err.Error())
		return
	}

	// 转换订单数据格式
	var orderList []map[string]interface{}
	for _, order := range orders {
		orderData := map[string]interface{}{
			"order_no":     order.OrderNo,
			"amount":       order.Amount,
			"status":       order.Status,
			"tx_hash":      order.TxHash,
			"created_at":   order.CreatedAt,
			"confirmed_at": order.ConfirmedAt,
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
