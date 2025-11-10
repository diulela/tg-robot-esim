package models

import (
	"time"

	"gorm.io/gorm"
)

// UserSession 用户会话模型
type UserSession struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	SessionData string         `gorm:"type:text" json:"session_data"`
	MenuPath    string         `gorm:"type:varchar(500)" json:"menu_path"`
	CurrentMenu string         `gorm:"type:varchar(100)" json:"current_menu"`
	Parameters  string         `gorm:"type:text" json:"parameters"` // JSON 格式存储
	LastActive  time.Time      `gorm:"index" json:"last_active"`
	CreatedAt   time.Time      `gorm:"type:datetime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (UserSession) TableName() string {
	return "user_sessions"
}

// BeforeCreate GORM 钩子：创建前
func (us *UserSession) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	us.CreatedAt = now
	us.UpdatedAt = now
	us.LastActive = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (us *UserSession) BeforeUpdate(tx *gorm.DB) error {
	us.UpdatedAt = time.Now()
	us.LastActive = time.Now()
	return nil
}

// IsExpired 检查会话是否过期
func (us *UserSession) IsExpired(timeout time.Duration) bool {
	return time.Since(us.LastActive) > timeout
}
