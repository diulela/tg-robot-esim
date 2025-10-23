package handlers

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Middleware 中间件接口
type Middleware interface {
	// ProcessMessage 处理消息中间件
	ProcessMessage(ctx context.Context, message *tgbotapi.Message, next MessageHandlerFunc) error

	// ProcessCallback 处理回调中间件
	ProcessCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, next CallbackHandlerFunc) error
}

// MessageHandlerFunc 消息处理函数类型
type MessageHandlerFunc func(ctx context.Context, message *tgbotapi.Message) error

// CallbackHandlerFunc 回调处理函数类型
type CallbackHandlerFunc func(ctx context.Context, callback *tgbotapi.CallbackQuery) error

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	logger Logger
}

// Logger 日志接口
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(logger Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// ProcessMessage 处理消息日志
func (m *LoggingMiddleware) ProcessMessage(ctx context.Context, message *tgbotapi.Message, next MessageHandlerFunc) error {
	start := time.Now()

	m.logger.Info("Processing message from user %d (@%s): %s",
		message.From.ID, message.From.UserName, message.Text)

	err := next(ctx, message)

	duration := time.Since(start)
	if err != nil {
		m.logger.Error("Message processing failed in %v: %v", duration, err)
	} else {
		m.logger.Debug("Message processed successfully in %v", duration)
	}

	return err
}

// ProcessCallback 处理回调日志
func (m *LoggingMiddleware) ProcessCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, next CallbackHandlerFunc) error {
	start := time.Now()

	m.logger.Info("Processing callback from user %d (@%s): %s",
		callback.From.ID, callback.From.UserName, callback.Data)

	err := next(ctx, callback)

	duration := time.Since(start)
	if err != nil {
		m.logger.Error("Callback processing failed in %v: %v", duration, err)
	} else {
		m.logger.Debug("Callback processed successfully in %v", duration)
	}

	return err
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	userLimits map[int64]time.Time
	interval   time.Duration
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(interval time.Duration) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		userLimits: make(map[int64]time.Time),
		interval:   interval,
	}
}

// ProcessMessage 处理消息限流
func (m *RateLimitMiddleware) ProcessMessage(ctx context.Context, message *tgbotapi.Message, next MessageHandlerFunc) error {
	userID := message.From.ID
	now := time.Now()

	if lastTime, exists := m.userLimits[userID]; exists {
		if now.Sub(lastTime) < m.interval {
			return nil // 忽略过于频繁的请求
		}
	}

	m.userLimits[userID] = now
	return next(ctx, message)
}

// ProcessCallback 处理回调限流
func (m *RateLimitMiddleware) ProcessCallback(ctx context.Context, callback *tgbotapi.CallbackQuery, next CallbackHandlerFunc) error {
	userID := callback.From.ID
	now := time.Now()

	if lastTime, exists := m.userLimits[userID]; exists {
		if now.Sub(lastTime) < m.interval {
			return nil // 忽略过于频繁的请求
		}
	}

	m.userLimits[userID] = now
	return next(ctx, callback)
}
