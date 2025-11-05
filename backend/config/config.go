package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Duration 自定义时间间隔类型，支持 JSON 字符串解析
type Duration time.Duration

// UnmarshalJSON 实现 JSON 反序列化
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("无效的时间间隔格式: %s", s)
	}

	*d = Duration(duration)
	return nil
}

// MarshalJSON 实现 JSON 序列化
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// ToDuration 转换为标准 time.Duration
func (d Duration) ToDuration() time.Duration {
	return time.Duration(d)
}

// String 返回时间间隔的字符串表示
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Seconds 返回秒数
func (d Duration) Seconds() float64 {
	return time.Duration(d).Seconds()
}

// Milliseconds 返回毫秒数
func (d Duration) Milliseconds() int64 {
	return time.Duration(d).Milliseconds()
}

// NewDuration 从字符串创建 Duration
func NewDuration(s string) (Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return Duration(0), fmt.Errorf("无效的时间间隔格式: %s", s)
	}
	return Duration(d), nil
}

// NewDurationFromTime 从 time.Duration 创建 Duration
func NewDurationFromTime(d time.Duration) Duration {
	return Duration(d)
}

// Config 主配置结构
type Config struct {
	Telegram   TelegramConfig   `json:"telegram"`
	Database   DatabaseConfig   `json:"database"`
	Blockchain BlockchainConfig `json:"blockchain"`
	Logging    LoggingConfig    `json:"logging"`
	Server     ServerConfig     `json:"server"`
	EsimSDK    EsimSDKConfig    `json:"esim_sdk"`
	Recharge   RechargeConfig   `json:"recharge"`
}

// TelegramConfig Telegram 相关配置
type TelegramConfig struct {
	BotToken   string   `json:"bot_token"`
	WebhookURL string   `json:"webhook_url"`
	MiniAppURL string   `json:"miniapp_url"`
	Timeout    Duration `json:"timeout"`
	Debug      bool     `json:"debug"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type           string `json:"type"` // sqlite, mysql
	DSN            string `json:"dsn"`  // 数据源名称
	MaxConnections int    `json:"max_connections"`
	MaxIdleConns   int    `json:"max_idle_conns"`
	ConnMaxLife    string `json:"conn_max_life"`
}

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	TronAPIKey            string   `json:"tron_api_key"`
	TronEndpoint          string   `json:"tron_endpoint"`
	MonitorInterval       Duration `json:"monitor_interval"`
	RequiredConfirmations int      `json:"required_confirmations"`
	MaxBlockDelay         int      `json:"max_block_delay"`
	WalletAddress         string   `json:"wallet_address"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level    string `json:"level"`    // debug, info, warn, error
	File     string `json:"file"`     // 日志文件路径
	MaxSize  int    `json:"max_size"` // MB
	MaxAge   int    `json:"max_age"`  // days
	Compress bool   `json:"compress"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int      `json:"port"`
	ReadTimeout  Duration `json:"read_timeout"`
	WriteTimeout Duration `json:"write_timeout"`
	IdleTimeout  Duration `json:"idle_timeout"`
}

// EsimSDKConfig eSIM SDK 配置
type EsimSDKConfig struct {
	APIKey         string `json:"api_key"`
	APISecret      string `json:"api_secret"`
	BaseURL        string `json:"base_url"`
	TimezoneOffset int    `json:"timezone_offset"` // 时区偏移（小时），例如：8 表示 UTC+8
}

// RechargeConfig 充值相关配置
type RechargeConfig struct {
	MinAmount              float64 `json:"min_amount"`               // 最小充值金额
	MaxAmount              float64 `json:"max_amount"`               // 最大充值金额
	OrderExpireMinutes     int     `json:"order_expire_minutes"`     // 订单过期时间（分钟）
	RequiredConfirmations  int     `json:"required_confirmations"`   // 所需确认数
	MonitorIntervalSeconds int     `json:"monitor_interval_seconds"` // 监控间隔（秒）
	DepositAddress         string  `json:"deposit_address"`          // 系统收款地址
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果不存在，创建默认配置文件
		if err := CreateDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 应用环境变量覆盖
	applyEnvironmentOverrides(&config)

	return &config, nil
}

// CreateDefaultConfig 创建默认配置文件
func CreateDefaultConfig(configPath string) error {
	defaultConfig := &Config{
		Telegram: TelegramConfig{
			BotToken:   "${TELEGRAM_BOT_TOKEN}",
			WebhookURL: "",
			MiniAppURL: "${MINIAPP_URL}",
			Timeout:    Duration(60 * time.Second),
			Debug:      false,
		},
		Database: DatabaseConfig{
			Type:           "sqlite",
			DSN:            "bot.db",
			MaxConnections: 10,
			MaxIdleConns:   5,
			ConnMaxLife:    "1h",
		},
		Blockchain: BlockchainConfig{
			TronAPIKey:            "${TRON_API_KEY}",
			TronEndpoint:          "https://api.trongrid.io",
			MonitorInterval:       Duration(30 * time.Second),
			RequiredConfirmations: 12,
			MaxBlockDelay:         100,
			WalletAddress:         "",
		},
		Logging: LoggingConfig{
			Level:    "info",
			File:     "bot.log",
			MaxSize:  100,
			MaxAge:   30,
			Compress: true,
		},
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  Duration(30 * time.Second),
			WriteTimeout: Duration(30 * time.Second),
			IdleTimeout:  Duration(120 * time.Second),
		},
		EsimSDK: EsimSDKConfig{
			APIKey:         "${ESIM_API_KEY}",
			APISecret:      "${ESIM_API_SECRET}",
			BaseURL:        "https://api.your-domain.com",
			TimezoneOffset: 0, // 默认使用 UTC 时间
		},
		Recharge: RechargeConfig{
			MinAmount:              10.0,
			MaxAmount:              10000.0,
			OrderExpireMinutes:     30,
			RequiredConfirmations:  19,
			MonitorIntervalSeconds: 30,
			DepositAddress:         "${DEPOSIT_WALLET_ADDRESS}",
		},
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config file: %w", err)
	}

	return nil
}

// applyEnvironmentOverrides 应用环境变量覆盖
func applyEnvironmentOverrides(config *Config) {
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		config.Telegram.BotToken = token
	}

	if miniAppURL := os.Getenv("MINIAPP_URL"); miniAppURL != "" {
		config.Telegram.MiniAppURL = miniAppURL
	}

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		config.Database.DSN = dsn
	}

	if apiKey := os.Getenv("TRON_API_KEY"); apiKey != "" {
		config.Blockchain.TronAPIKey = apiKey
	}

	if debug := os.Getenv("DEBUG"); debug == "true" {
		config.Telegram.Debug = true
		config.Logging.Level = "debug"
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}

	if esimAPIKey := os.Getenv("ESIM_API_KEY"); esimAPIKey != "" {
		config.EsimSDK.APIKey = esimAPIKey
	}

	if esimAPISecret := os.Getenv("ESIM_API_SECRET"); esimAPISecret != "" {
		config.EsimSDK.APISecret = esimAPISecret
	}

	if esimBaseURL := os.Getenv("ESIM_BASE_URL"); esimBaseURL != "" {
		config.EsimSDK.BaseURL = esimBaseURL
	}

	if depositAddress := os.Getenv("DEPOSIT_WALLET_ADDRESS"); depositAddress != "" {
		config.Recharge.DepositAddress = depositAddress
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Telegram.BotToken == "" || c.Telegram.BotToken == "${TELEGRAM_BOT_TOKEN}" {
		return fmt.Errorf("telegram bot token is required")
	}

	// 验证 Mini App URL
	// if c.Telegram.MiniAppURL != "" && c.Telegram.MiniAppURL != "${MINIAPP_URL}" {
	// 	if !strings.HasPrefix(c.Telegram.MiniAppURL, "https://") {
	// 		return fmt.Errorf("mini app URL must use HTTPS protocol")
	// 	}
	// }

	if c.Database.Type != "sqlite" && c.Database.Type != "mysql" {
		return fmt.Errorf("unsupported database type: %s", c.Database.Type)
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	if c.Blockchain.RequiredConfirmations < 1 {
		return fmt.Errorf("required confirmations must be at least 1")
	}

	// 验证充值配置
	if c.Recharge.MinAmount <= 0 {
		return fmt.Errorf("recharge min amount must be greater than 0")
	}

	if c.Recharge.MaxAmount <= c.Recharge.MinAmount {
		return fmt.Errorf("recharge max amount must be greater than min amount")
	}

	if c.Recharge.OrderExpireMinutes <= 0 {
		return fmt.Errorf("recharge order expire minutes must be greater than 0")
	}

	if c.Recharge.RequiredConfirmations < 1 {
		return fmt.Errorf("recharge required confirmations must be at least 1")
	}

	return nil
}
