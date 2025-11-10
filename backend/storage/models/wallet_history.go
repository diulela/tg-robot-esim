package models

import (
	"time"

	"gorm.io/gorm"
)

// WalletHistoryType 钱包历史记录类型
type WalletHistoryType string

const (
	WalletHistoryTypeRecharge WalletHistoryType = "recharge" // 充值
	WalletHistoryTypePayment  WalletHistoryType = "payment"  // 支付（包含冻结金额的最终扣费）
	WalletHistoryTypeRefund   WalletHistoryType = "refund"   // 退款（包含冻结金额的退还）
)

// WalletHistoryStatus 钱包历史记录状态
type WalletHistoryStatus string

const (
	WalletHistoryStatusPending   WalletHistoryStatus = "pending"   // 处理中
	WalletHistoryStatusCompleted WalletHistoryStatus = "completed" // 已完成
	WalletHistoryStatusFailed    WalletHistoryStatus = "failed"    // 失败
	WalletHistoryStatusCancelled WalletHistoryStatus = "cancelled" // 已取消
)

// WalletHistory 钱包历史记录模型
type WalletHistory struct {
	ID            uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64               `gorm:"index;not null" json:"user_id"`                     // 用户ID
	Type          WalletHistoryType   `gorm:"size:20;not null;index" json:"type"`                // 记录类型
	Amount        string              `gorm:"type:decimal(10,2);not null" json:"amount"`         // 金额（正数为收入，负数为支出）
	BalanceBefore string              `gorm:"type:decimal(10,2);not null" json:"balance_before"` // 操作前余额
	BalanceAfter  string              `gorm:"type:decimal(10,2);not null" json:"balance_after"`  // 操作后余额
	Status        WalletHistoryStatus `gorm:"size:20;default:'completed';index" json:"status"`   // 记录状态
	Description   string              `gorm:"size:500;not null" json:"description"`              // 描述
	RelatedType   string              `gorm:"size:50;index" json:"related_type"`                 // 关联类型（recharge_order, order）
	RelatedID     string              `gorm:"size:100;index" json:"related_id"`                  // 关联ID（订单号等）
	TxHash        string              `gorm:"size:100;index" json:"tx_hash"`                     // 区块链交易哈希（如果有）
	Metadata      string              `gorm:"type:text" json:"metadata"`                         // 额外元数据（JSON格式）
	CreatedAt     time.Time           `gorm:"type:datetime" json:"created_at"`
	UpdatedAt     time.Time           `gorm:"type:datetime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (WalletHistory) TableName() string {
	return "wallet_histories"
}

// BeforeCreate GORM 钩子：创建前
func (w *WalletHistory) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	w.CreatedAt = now
	w.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (w *WalletHistory) BeforeUpdate(tx *gorm.DB) error {
	w.UpdatedAt = time.Now()
	return nil
}

// IsIncome 检查是否为收入记录
func (w *WalletHistory) IsIncome() bool {
	return w.Type == WalletHistoryTypeRecharge || w.Type == WalletHistoryTypeRefund
}

// IsExpense 检查是否为支出记录
func (w *WalletHistory) IsExpense() bool {
	return w.Type == WalletHistoryTypePayment
}

// IsCompleted 检查记录是否已完成
func (w *WalletHistory) IsCompleted() bool {
	return w.Status == WalletHistoryStatusCompleted
}
