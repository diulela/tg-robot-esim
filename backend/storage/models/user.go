package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	TelegramID    int64          `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username      string         `gorm:"index" json:"username"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	Language      string         `gorm:"default:'zh'" json:"language"`
	WalletAddress string         `gorm:"size:100;index" json:"wallet_address"` // 钱包地址
	IsVIP         bool           `gorm:"default:false" json:"is_vip"`          // VIP 状态
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM 钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
