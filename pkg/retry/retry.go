package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int           // 最大重试次数
	InitialDelay    time.Duration // 初始延迟
	MaxDelay        time.Duration // 最大延迟
	BackoffFactor   float64       // 退避因子
	RetryableErrors []error       // 可重试的错误类型
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:    3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
	}
}

// RetryFunc 重试函数类型
type RetryFunc func() error

// IsRetryableFunc 判断错误是否可重试的函数类型
type IsRetryableFunc func(error) bool

// Retry 执行重试逻辑
func Retry(ctx context.Context, config *RetryConfig, fn RetryFunc, isRetryable IsRetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// 执行函数
		err := fn()
		if err == nil {
			return nil // 成功
		}

		lastErr = err

		// 检查是否可重试
		if isRetryable != nil && !isRetryable(err) {
			return err // 不可重试的错误
		}

		// 最后一次尝试失败
		if attempt == config.MaxRetries {
			break
		}

		// 计算延迟时间
		delay := calculateDelay(config, attempt)

		// 等待重试
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// 继续重试
		}
	}

	return fmt.Errorf("retry failed after %d attempts, last error: %w", config.MaxRetries+1, lastErr)
}

// calculateDelay 计算延迟时间（指数退避）
func calculateDelay(config *RetryConfig, attempt int) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt))

	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// CircuitBreaker 断路器
type CircuitBreaker struct {
	maxFailures     int
	resetTimeout    time.Duration
	failureCount    int
	lastFailureTime time.Time
	state           CircuitState
}

// CircuitState 断路器状态
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// String 返回状态字符串
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// NewCircuitBreaker 创建断路器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Execute 执行函数（带断路器保护）
func (cb *CircuitBreaker) Execute(fn RetryFunc) error {
	// 检查断路器状态
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.failureCount = 0
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	// 执行函数
	err := fn()

	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

// onSuccess 成功时的处理
func (cb *CircuitBreaker) onSuccess() {
	cb.failureCount = 0
	cb.state = StateClosed
}

// onFailure 失败时的处理
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = StateOpen
	}
}

// GetState 获取断路器状态
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// GetFailureCount 获取失败次数
func (cb *CircuitBreaker) GetFailureCount() int {
	return cb.failureCount
}
