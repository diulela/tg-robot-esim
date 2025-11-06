package services

import (
	"context"
	"tg-robot-sim/storage/models"
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

	// GetAddressIncomingTransactions 获取地址的入账交易
	GetAddressIncomingTransactions(ctx context.Context, address string, minAmount string) ([]*TransactionInfo, error)

	// GetTransactionByHash 根据哈希获取交易详情
	GetTransactionByHash(ctx context.Context, txHash string) (*TransactionInfo, error)

	// MatchTransactionAmount 匹配交易金额
	MatchTransactionAmount(txAmount string, targetAmount string) bool
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

	// SendRechargeSuccessNotification 发送充值成功通知
	SendRechargeSuccessNotification(ctx context.Context, userID int64, amount string, orderNo string) error
}

// RechargeService 定义充值服务接口
// 负责处理用户充值相关业务逻辑
type RechargeService interface {
	// CreateRechargeOrder 创建充值订单
	CreateRechargeOrder(ctx context.Context, userID int64, amount string) (*models.RechargeOrder, error)

	// GetRechargeOrder 获取充值订单详情
	GetRechargeOrder(ctx context.Context, orderNo string) (*models.RechargeOrder, error)

	// GetUserRechargeHistory 获取用户充值历史
	GetUserRechargeHistory(ctx context.Context, userID int64, limit, offset int) ([]*models.RechargeOrder, int64, error)

	// CheckRechargeStatus 手动检查充值状态
	CheckRechargeStatus(ctx context.Context, orderNo string) (*models.RechargeOrder, error)

	// ProcessPendingRecharges 处理待支付的充值订单（定时任务调用）
	ProcessPendingRecharges(ctx context.Context) error

	// ConfirmRecharge 确认充值并更新余额，并发送 Telegram 通知
	ConfirmRecharge(ctx context.Context, order *models.RechargeOrder, txHash string) error

	// ExpireOldOrders 将过期订单标记为已过期
	ExpireOldOrders(ctx context.Context) error

	// GenerateExactAmount 生成唯一的精确金额
	GenerateExactAmount(ctx context.Context, baseAmount string) (string, error)
}

// WalletHistoryFilters 钱包历史筛选条件
type WalletHistoryFilters struct {
	UserID    int64                      `json:"user_id"`
	Type      models.WalletHistoryType   `json:"type"`
	Status    models.WalletHistoryStatus `json:"status"`
	StartDate string                     `json:"start_date"`
	EndDate   string                     `json:"end_date"`
	Limit     int                        `json:"limit"`
	Offset    int                        `json:"offset"`
}

// WalletHistoryStats 钱包历史统计
type WalletHistoryStats struct {
	TotalRecords    int64  `json:"total_records"`
	TotalIncome     string `json:"total_income"`     // 总收入（充值+退款）
	TotalExpense    string `json:"total_expense"`    // 总支出（支付）
	PendingAmount   string `json:"pending_amount"`   // 处理中金额
	CompletedAmount string `json:"completed_amount"` // 已完成金额
}

// WalletHistoryService 定义钱包历史服务接口
// 负责管理用户钱包交易历史记录的业务逻辑
type WalletHistoryService interface {
	// GetWalletHistory 获取钱包历史记录
	GetWalletHistory(ctx context.Context, userID int64, filters WalletHistoryFilters) ([]*models.WalletHistory, int64, error)

	// GetWalletHistoryStats 获取钱包历史统计
	GetWalletHistoryStats(ctx context.Context, userID int64) (*WalletHistoryStats, error)

	// GetHistoryRecord 获取单条历史记录详情
	GetHistoryRecord(ctx context.Context, recordID uint, userID int64) (*models.WalletHistory, error)

	// CreateRechargeRecord 创建充值记录
	CreateRechargeRecord(ctx context.Context, userID int64, amount, balanceBefore, balanceAfter string, relatedID, txHash string) error

	// CreatePaymentRecord 创建支付记录
	CreatePaymentRecord(ctx context.Context, userID int64, amount, balanceBefore, balanceAfter string, relatedID, description string) error

	// UpdateRecordStatus 更新记录状态
	UpdateRecordStatus(ctx context.Context, recordID uint, status models.WalletHistoryStatus) error

	// CreateRefundRecord 创建退款记录（用于订单失败时的退款）
	CreateRefundRecord(ctx context.Context, userID int64, amount, balanceBefore, balanceAfter string, relatedID, description string) error
}

// CreateEsimOrderRequest 创建 eSIM 订单请求
type CreateEsimOrderRequest struct {
	UserID      int64  `json:"user_id" validate:"required"`
	ProductID   int    `json:"product_id" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,min=1"`
	TotalAmount string `json:"total_amount" validate:"required"`
	Remark      string `json:"remark,omitempty"`
}

// EsimOrderResponse eSIM 订单响应
type EsimOrderResponse struct {
	OrderID     uint               `json:"order_id"`
	OrderNo     string             `json:"order_no"`
	Status      models.OrderStatus `json:"status"`
	TotalAmount string             `json:"total_amount"`
	CreatedAt   time.Time          `json:"created_at"`
}

// OrderWithDetail 包含详情的订单信息
type OrderWithDetail struct {
	OrderID         uint               `json:"order_id"`
	OrderNo         string             `json:"order_no"`
	UserID          int64              `json:"user_id"`
	ProductID       int                `json:"product_id"`
	ProductName     string             `json:"product_name"`
	Quantity        int                `json:"quantity"`
	UnitPrice       string             `json:"unit_price"`
	TotalAmount     string             `json:"total_amount"`
	Status          models.OrderStatus `json:"status"`
	ProviderOrderID string             `json:"provider_order_id,omitempty"`
	OrderItems      []OrderItemDetail  `json:"order_items,omitempty"`
	Esims           []EsimDetail       `json:"esims,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	CompletedAt     *time.Time         `json:"completed_at,omitempty"`
}

// OrderItemDetail 订单项详情
type OrderItemDetail struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	DataSize    int     `json:"data_size"`
	ValidDays   int     `json:"valid_days"`
}

// EsimDetail eSIM 详情
type EsimDetail struct {
	ID                int    `json:"id"`
	ICCID             string `json:"iccid"`
	Status            string `json:"status"`
	HasActivationCode bool   `json:"has_activation_code"`
	HasQrCode         bool   `json:"has_qr_code"`
}

// ProviderOrderData 第三方订单数据
type ProviderOrderData struct {
	OrderID     int               `json:"order_id"`
	OrderNumber string            `json:"order_number"`
	Status      string            `json:"status"`
	OrderItems  []OrderItemDetail `json:"order_items"`
	Esims       []EsimDetail      `json:"esims"`
}

// OrderFilters 订单筛选条件
type OrderFilters struct {
	Status    models.OrderStatus `json:"status,omitempty"`
	StartDate string             `json:"start_date,omitempty"`
	EndDate   string             `json:"end_date,omitempty"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
}
