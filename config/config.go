package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config 主配置结构
type Config struct {
	Telegram   TelegramConfig   `json:"telegram"`
	Database   DatabaseConfig   `json:"database"`
	Blockchain BlockchainConfig `json:"blockchain"`
	Logging    LoggingConfig    `json:"logging"`
	Server     ServerConfig     `json:"server"`
}

// TelegramConfig Telegram 相关配置
type TelegramConfig struct {
	BotToken   string        `json:"bot_token"`
	WebhookURL string        `json:"webhook_url"`
	Timeout    time.Duration `json:"timeout"`
	Debug      bool          `json:"debug"`
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
	TronAPIKey            string        `json:"tron_api_key"`
	TronEndpoint          string        `json:"tron_endpoint"`
	MonitorInterval       time.Duration `json:"monitor_interval"`
	RequiredConfirmations int           `json:"required_confirmations"`
	MaxBlockDelay         int           `json:"max_block_delay"`
	WalletAddress         string        `json:"wallet_address"`
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
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
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
			Timeout:    60 * time.Second,
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
			MonitorInterval:       30 * time.Second,
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
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
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
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Telegram.BotToken == "" || c.Telegram.BotToken == "${TELEGRAM_BOT_TOKEN}" {
		return fmt.Errorf("telegram bot token is required")
	}

	if c.Database.Type != "sqlite" && c.Database.Type != "mysql" {
		return fmt.Errorf("unsupported database type: %s", c.Database.Type)
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("database DSN is required")
	}

	if c.Blockchain.RequiredConfirmations < 1 {
		return fmt.Errorf("required confirmations must be at least 1")
	}

	return nil
}
