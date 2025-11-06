package config

import "time"

// OrderConfig 订单处理相关配置
type OrderConfig struct {
	// 同步配置
	SyncInterval    time.Duration `json:"sync_interval"`     // 同步间隔
	MaxSyncAttempts int           `json:"max_sync_attempts"` // 最大同步尝试次数
	OrderTimeout    time.Duration `json:"order_timeout"`     // 订单超时时间

	// 通知配置
	NotificationEnabled bool `json:"notification_enabled"` // 是否启用通知

	// 重试配置
	RetryConfig RetryConfig `json:"retry_config"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts   int           `json:"max_attempts"`   // 最大重试次数
	InitialDelay  time.Duration `json:"initial_delay"`  // 初始延迟
	MaxDelay      time.Duration `json:"max_delay"`      // 最大延迟
	BackoffFactor float64       `json:"backoff_factor"` // 退避因子
}

// DefaultOrderConfig 默认订单配置
func DefaultOrderConfig() *OrderConfig {
	return &OrderConfig{
		SyncInterval:        10 * time.Second,
		MaxSyncAttempts:     100,
		OrderTimeout:        30 * time.Minute,
		NotificationEnabled: true,
		RetryConfig: RetryConfig{
			MaxAttempts:   5,
			InitialDelay:  2 * time.Second,
			MaxDelay:      5 * time.Minute,
			BackoffFactor: 2.0,
		},
	}
}
