package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
