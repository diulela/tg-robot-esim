package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// handlePurchase 处理购买请求
// @deprecated 此接口已废弃，请使用 POST /api/miniapp/esim/orders
// 废弃日期: 2025-01-15
// 计划下线日期: 2025-06-01
func (h *MiniAppApiService) handlePurchase(w http.ResponseWriter, r *http.Request) {
	// 添加废弃标记响应头
	w.Header().Set("X-API-Deprecated", "true")
	w.Header().Set("X-API-Deprecated-Message", "此接口已废弃，请使用 POST /api/miniapp/esim/orders")
	w.Header().Set("X-API-Deprecated-Since", "2025-01-15")
	w.Header().Set("X-API-Sunset-Date", "2025-06-01")
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

	// 记录旧接口调用日志
	log.Printf("[DEPRECATED API] /api/miniapp/purchase called by user_id=%d at %s", userID, time.Now().Format(time.RFC3339))

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
// @deprecated 此接口已废弃，请使用 GET /api/miniapp/esim/orders
// 废弃日期: 2025-01-15
// 计划下线日期: 2025-06-01
func (h *MiniAppApiService) handleOrders(w http.ResponseWriter, r *http.Request) {
	// 添加废弃标记响应头
	w.Header().Set("X-API-Deprecated", "true")
	w.Header().Set("X-API-Deprecated-Message", "此接口已废弃，请使用 GET /api/miniapp/esim/orders")
	w.Header().Set("X-API-Deprecated-Since", "2025-01-15")
	w.Header().Set("X-API-Sunset-Date", "2025-06-01")

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

	// 记录旧接口调用日志
	log.Printf("[DEPRECATED API] /api/miniapp/orders called by user_id=%d at %s", userID, time.Now().Format(time.RFC3339))

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
