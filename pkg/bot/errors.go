package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/pkg/retry"
)

// TelegramError Telegram API 错误类型
type TelegramError struct {
	Code        int
	Description string
	Parameters  map[string]interface{}
}

// Error 实现 error 接口
func (e *TelegramError) Error() string {
	return fmt.Sprintf("telegram error %d: %s", e.Code, e.Description)
}

// IsRetryable 判断错误是否可重试
func (e *TelegramError) IsRetryable() bool {
	switch e.Code {
	case 429: // Too Many Requests
		return true
	case 500, 502, 503, 504: // Server errors
		return true
	case 400: // Bad Request - 通常不可重试
		return false
	case 401: // Unauthorized - 不可重试
		return false
	case 403: // Forbidden - 不可重试
		return false
	case 404: // Not Found - 不可重试
		return false
	default:
		// 其他错误默认可重试
		return true
	}
}

// GetRetryAfter 获取重试延迟时间
func (e *TelegramError) GetRetryAfter() time.Duration {
	if retryAfter, ok := e.Parameters["retry_after"]; ok {
		if seconds, ok := retryAfter.(int); ok {
			return time.Duration(seconds) * time.Second
		}
	}
	return 0
}

// ErrorHandler Telegram API 错误处理器
type ErrorHandler struct {
	logger         Logger
	retryConfig    *retry.RetryConfig
	circuitBreaker *retry.CircuitBreaker
}

// Logger 日志接口
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(logger Logger) *ErrorHandler {
	retryConfig := retry.DefaultRetryConfig()
	retryConfig.MaxRetries = 3
	retryConfig.InitialDelay = 1 * time.Second
	retryConfig.MaxDelay = 30 * time.Second

	circuitBreaker := retry.NewCircuitBreaker(5, 1*time.Minute)

	return &ErrorHandler{
		logger:         logger,
		retryConfig:    retryConfig,
		circuitBreaker: circuitBreaker,
	}
}

// HandleAPICall 处理 Telegram API 调用（带重试和断路器）
func (h *ErrorHandler) HandleAPICall(ctx context.Context, fn func() (tgbotapi.Message, error)) (tgbotapi.Message, error) {
	var result tgbotapi.Message

	retryFunc := func() error {
		msg, err := fn()
		if err != nil {
			return err
		}
		result = msg
		return nil
	}

	// 使用断路器保护
	err := h.circuitBreaker.Execute(func() error {
		return retry.Retry(ctx, h.retryConfig, retryFunc, h.isRetryableError)
	})

	if err != nil {
		h.logger.Error("API call failed: %v", err)
		return tgbotapi.Message{}, err
	}

	return result, nil
}

// HandleAPIRequest 处理 Telegram API 请求（带重试和断路器）
func (h *ErrorHandler) HandleAPIRequest(ctx context.Context, fn func() (tgbotapi.APIResponse, error)) (tgbotapi.APIResponse, error) {
	var result tgbotapi.APIResponse

	retryFunc := func() error {
		resp, err := fn()
		if err != nil {
			return err
		}
		result = resp
		return nil
	}

	// 使用断路器保护
	err := h.circuitBreaker.Execute(func() error {
		return retry.Retry(ctx, h.retryConfig, retryFunc, h.isRetryableError)
	})

	if err != nil {
		h.logger.Error("API request failed: %v", err)
		return tgbotapi.APIResponse{}, err
	}

	return result, nil
}

// isRetryableError 判断错误是否可重试
func (h *ErrorHandler) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否是 Telegram API 错误
	if tgErr := h.parseTelegramError(err); tgErr != nil {
		// 如果是 429 错误，需要等待指定时间
		if tgErr.Code == 429 {
			retryAfter := tgErr.GetRetryAfter()
			if retryAfter > 0 {
				h.logger.Warn("Rate limited, waiting %v before retry", retryAfter)
				time.Sleep(retryAfter)
			}
		}
		return tgErr.IsRetryable()
	}

	// 检查网络错误
	errStr := err.Error()
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "no such host") {
		return true
	}

	// 默认不重试
	return false
}

// parseTelegramError 解析 Telegram API 错误
func (h *ErrorHandler) parseTelegramError(err error) *TelegramError {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// 解析常见的 Telegram API 错误
	if strings.Contains(errStr, "Too Many Requests") {
		return &TelegramError{
			Code:        429,
			Description: "Too Many Requests",
		}
	}

	if strings.Contains(errStr, "Bad Request") {
		return &TelegramError{
			Code:        400,
			Description: "Bad Request",
		}
	}

	if strings.Contains(errStr, "Unauthorized") {
		return &TelegramError{
			Code:        401,
			Description: "Unauthorized",
		}
	}

	if strings.Contains(errStr, "Forbidden") {
		return &TelegramError{
			Code:        403,
			Description: "Forbidden",
		}
	}

	// 检查 HTTP 状态码
	if strings.Contains(errStr, "500") {
		return &TelegramError{
			Code:        500,
			Description: "Internal Server Error",
		}
	}

	if strings.Contains(errStr, "502") {
		return &TelegramError{
			Code:        502,
			Description: "Bad Gateway",
		}
	}

	if strings.Contains(errStr, "503") {
		return &TelegramError{
			Code:        503,
			Description: "Service Unavailable",
		}
	}

	return nil
}

// GetCircuitBreakerState 获取断路器状态
func (h *ErrorHandler) GetCircuitBreakerState() retry.CircuitState {
	return h.circuitBreaker.GetState()
}

// GetFailureCount 获取失败次数
func (h *ErrorHandler) GetFailureCount() int {
	return h.circuitBreaker.GetFailureCount()
}
