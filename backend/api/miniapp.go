package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tg-robot-sim/pkg/telegram"
	"tg-robot-sim/services"
)

// MiniAppHandler Mini App API 处理器
type MiniAppApiService struct {
	productService     services.ProductService
	walletService      services.WalletService
	orderService       services.OrderService
	transactionService services.TransactionService
}

// NewMiniAppHandler 创建 Mini App 处理器实例
func NewMiniAppApiService(
	productService services.ProductService,
	walletService services.WalletService,
	orderService services.OrderService,
	transactionService services.TransactionService,
) *MiniAppApiService {
	return &MiniAppApiService{
		productService:     productService,
		walletService:      walletService,
		orderService:       orderService,
		transactionService: transactionService,
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
	mux.HandleFunc("/api/miniapp/wallet/recharge", h.handleWalletRecharge)

	// 订单相关
	mux.HandleFunc("/api/miniapp/purchase", h.handlePurchase)
	mux.HandleFunc("/api/miniapp/orders", h.handleOrders)

	// 交易相关
	mux.HandleFunc("/api/miniapp/transactions", h.handleTransactions)
}

// handleProducts 处理产品列表请求
func (h *MiniAppApiService) handleProducts(w http.ResponseWriter, r *http.Request) {
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

	// 获取钱包余额
	balance, err := h.walletService.GetBalance(ctx, userID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get balance", err.Error())
		return
	}

	h.sendSuccess(w, balance)
}

// handleWalletRecharge 处理钱包充值请求
func (h *MiniAppApiService) handleWalletRecharge(w http.ResponseWriter, r *http.Request) {
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

	// 创建充值订单
	order, err := h.walletService.CreateRechargeOrder(ctx, userID, req.Amount)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to create recharge order", err.Error())
		return
	}

	h.sendSuccess(w, order)
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

// handleTransactions 处理交易历史请求
func (h *MiniAppApiService) handleTransactions(w http.ResponseWriter, r *http.Request) {
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

	// 获取交易历史
	transactions, err := h.transactionService.GetTransactionHistory(ctx, userID, limit, offset)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get transactions", err.Error())
		return
	}

	h.sendSuccess(w, map[string]interface{}{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}
