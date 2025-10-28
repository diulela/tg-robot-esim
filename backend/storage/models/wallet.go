package models

import (
	"time"

	"gorm.io/gorm"
)

// Wallet 钱包模型
type Wallet struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"uniqueIndex;not null" json:"user_id"`
	Balance       string    `gorm:"type:decimal(20,8);default:'0'" json:"balance"`        // 可用余额
	FrozenBalance string    `gorm:"type:decimal(20,8);default:'0'" json:"frozen_balance"` // 冻结余额
	TotalIncome   string    `gorm:"type:decimal(20,8);default:'0'" json:"total_income"`   // 总收入
	TotalExpense  string    `gorm:"type:decimal(20,8);default:'0'" json:"total_expense"`  // 总支出
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Wallet) TableName() string {
	return "wallets"
}

// BeforeCreate GORM 钩子：创建前
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	w.CreatedAt = now
	w.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (w *Wallet) BeforeUpdate(tx *gorm.DB) error {
	w.UpdatedAt = time.Now()
	return nil
}
