package services

import (
	"context"
	"time"
)

// DialogResponse 对话响应结构
type DialogResponse struct {
	Message    string                 `json:"message"`
	Keyboard   interface{}            `json:"keyboard,omitempty"`
	ParseMode  string                 `json:"parse_mode,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// UserContext 用户上下文结构
type UserContext struct {
	UserID      int64                  `json:"user_id"`
	CurrentMenu string                 `json:"current_menu"`
	MenuPath    []string               `json:"menu_path"`
	Parameters  map[string]interface{} `json:"parameters"`
	LastActive  time.Time              `json:"last_active"`
}

// MenuResponse 菜单响应结构
type MenuResponse struct {
	Text      string      `json:"text"`
	Keyboard  interface{} `json:"keyboard"`
	ParseMode string      `json:"parse_mode,omitempty"`
	EditMode  bool        `json:"edit_mode"`
}

// TransactionInfo 交易信息结构
type TransactionInfo struct {
	TxHash        string    `json:"tx_hash"`
	FromAddress   string    `json:"from_address"`
	ToAddress     string    `json:"to_address"`
	Amount        string    `json:"amount"`
	Confirmations int       `json:"confirmations"`
	BlockNumber   int64     `json:"block_number"`
	Timestamp     time.Time `json:"timestamp"`
	Status        string    `json:"status"`
}

// TransactionStatus 交易状态枚举
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// DialogService 定义对话服务接口
// 负责处理用户对话逻辑和会话管理
type DialogService interface {
	// ProcessMessage 处理用户消息
	ProcessMessage(ctx context.Context, userID int64, message string) (*DialogResponse, error)

	// GetUserContext 获取用户上下文
	GetUserContext(userID int64) (*UserContext, error)

	// SetUserContext 设置用户上下文
	SetUserContext(userID int64, context *UserContext) error

	// ClearUserContext 清除用户上下文
	ClearUserContext(userID int64) error

	// IsUserActive 检查用户是否活跃
	IsUserActive(userID int64) bool
}

// MenuService 定义菜单服务接口
// 负责管理交互式菜单系统
type MenuService interface {
	// GetMainMenu 获取主菜单
	GetMainMenu(userID int64) (*MenuResponse, error)

	// HandleMenuAction 处理菜单操作
	HandleMenuAction(ctx context.Context, userID int64, action string) (*MenuResponse, error)

	// GetMenuHistory 获取菜单历史
	GetMenuHistory(userID int64) ([]string, error)

	// NavigateBack 返回上级菜单
	NavigateBack(userID int64) (*MenuResponse, error)

	// ResetMenu 重置菜单到主菜单
	ResetMenu(userID int64) (*MenuResponse, error)
}

// BlockchainService 定义区块链服务接口
// 负责监控和验证区块链交易
type BlockchainService interface {
	// StartMonitoring 开始监控区块链交易
	StartMonitoring(ctx context.Context) error

	// StopMonitoring 停止监控
	StopMonitoring() error

	// ValidateTransaction 验证交易
	ValidateTransaction(txHash string) (*TransactionInfo, error)

	// GetTransactionStatus 获取交易状态
	GetTransactionStatus(txHash string) (TransactionStatus, error)

	// MonitorAddress 监控指定地址的交易
	MonitorAddress(address string) error

	// GetAddressTransactions 获取地址的交易记录
	GetAddressTransactions(address string, limit int) ([]*TransactionInfo, error)

	// IsTransactionConfirmed 检查交易是否已确认
	IsTransactionConfirmed(txHash string, requiredConfirmations int) (bool, error)
}

// NotificationService 定义通知服务接口
// 负责发送各种类型的通知
type NotificationService interface {
	// SendMessage 发送消息
	SendMessage(ctx context.Context, userID int64, message string) error

	// SendMenuMessage 发送菜单消息
	SendMenuMessage(ctx context.Context, userID int64, response *MenuResponse) error

	// EditMessage 编辑消息
	EditMessage(ctx context.Context, userID int64, messageID int, newText string) error

	// SendTransactionNotification 发送交易通知
	SendTransactionNotification(ctx context.Context, userID int64, txInfo *TransactionInfo) error
}
