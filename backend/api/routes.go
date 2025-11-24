package api

import "net/http"

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

	// eSIM 订单相关
	mux.HandleFunc("/api/miniapp/esim/orders", h.handleEsimOrders)
	mux.HandleFunc("/api/miniapp/esim/orders/", h.handleEsimOrderDetail)

	// eSIM 卡相关
	mux.HandleFunc("/api/miniapp/esim/cards", h.handleEsimCards)

	// 钱包历史相关
	mux.HandleFunc("/api/miniapp/wallet/history", h.handleWalletHistory)
	mux.HandleFunc("/api/miniapp/wallet/history/stats", h.handleWalletHistoryStats)
	mux.HandleFunc("/api/miniapp/wallet/history/", h.handleHistoryRecord)
}
