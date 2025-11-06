package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"tg-robot-sim/services"
	"tg-robot-sim/storage/models"
)

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
