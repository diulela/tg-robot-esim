package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tg-robot-sim/pkg/telegram"
	"tg-robot-sim/services"
)

// MiniAppApiService Mini App API 处理器
type MiniAppApiService struct {
	productService       services.ProductService
	walletService        services.WalletService
	orderService         services.OrderService
	walletHistoryService services.WalletHistoryService
	rechargeService      services.RechargeService
	esimCardService      services.EsimCardService
}

// NewMiniAppApiService 创建 Mini App 处理器实例
func NewMiniAppApiService(
	productService services.ProductService,
	walletService services.WalletService,
	orderService services.OrderService,
	walletHistoryService services.WalletHistoryService,
	rechargeService services.RechargeService,
	esimCardService services.EsimCardService,
) *MiniAppApiService {
	return &MiniAppApiService{
		productService:       productService,
		walletService:        walletService,
		orderService:         orderService,
		walletHistoryService: walletHistoryService,
		rechargeService:      rechargeService,
		esimCardService:      esimCardService,
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
	// 客户端错误 (40xxx)
	ErrCodeInvalidRequest      = 40000 // 无效的请求
	ErrCodeInvalidAmount       = 40001 // 充值金额低于最小限额或余额不足
	ErrCodeInvalidFormat       = 40002 // 充值金额格式错误
	ErrCodeOrderNotFound       = 40003 // 订单不存在
	ErrCodeOrderExpired        = 40004 // 订单已过期
	ErrCodeOrderCompleted      = 40005 // 订单已完成，不能重复处理
	ErrCodeWalletNotFound      = 40006 // 用户钱包不存在
	ErrCodeProductNotFound     = 40007 // 产品不存在
	ErrCodeProductUnavailable  = 40008 // 产品暂不可用
	ErrCodeUnauthorized        = 40100 // 未授权访问
	ErrCodeInsufficientBalance = 40009 // 余额不足（用于订单创建）
	ErrCodeNotFound            = 40400 // 资源未找到

	// 服务器错误 (50xxx)
	ErrCodeGenerateAmount  = 50001 // 生成精确金额失败
	ErrCodeBlockchainQuery = 50002 // 区块链查询失败
	ErrCodeDatabaseError   = 50003 // 数据库操作失败
	ErrCodeBalanceUpdate   = 50004 // 余额更新失败
	ErrCodeInternalError   = 50000 // 内部服务器错误
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
