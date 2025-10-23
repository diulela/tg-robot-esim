package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionStatus 交易状态枚举
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction 区块链交易模型
type Transaction struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	TxHash        string            `gorm:"uniqueIndex;not null" json:"tx_hash"`
	FromAddress   string            `gorm:"index" json:"from_address"`
	ToAddress     string            `gorm:"index" json:"to_address"`
	Amount        string            `gorm:"not null" json:"amount"` // 使用字符串存储精确金额
	TokenSymbol   string            `gorm:"default:'USDT'" json:"token_symbol"`
	Status        TransactionStatus `gorm:"default:'pending'" json:"status"`
	Confirmations int               `gorm:"default:0" json:"confirmations"`
	BlockNumber   int64             `json:"block_number"`
	BlockHash     string            `json:"block_hash"`
	GasUsed       int64             `json:"gas_used"`
	GasPrice      string            `json:"gas_price"`
	Timestamp     time.Time         `json:"timestamp"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate GORM 钩子：创建前
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

// IsConfirmed 检查交易是否已确认
func (t *Transaction) IsConfirmed(requiredConfirmations int) bool {
	return t.Status == TransactionStatusConfirmed && t.Confirmations >= requiredConfirmations
}
