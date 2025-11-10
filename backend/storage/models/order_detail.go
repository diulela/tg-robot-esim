package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderDetail 订单详情模型
type OrderDetail struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uint      `gorm:"uniqueIndex;not null" json:"order_id"` // 关联订单ID
	ProviderData string    `gorm:"type:longtext" json:"provider_data"`   // JSON 格式存储第三方完整数据
	OrderItems   string    `gorm:"type:longtext" json:"order_items"`     // JSON 格式存储 OrderItems
	Esims        string    `gorm:"type:longtext" json:"esims"`           // JSON 格式存储 Esims
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (OrderDetail) TableName() string {
	return "order_details"
}

// BeforeCreate GORM 钩子：创建前
func (od *OrderDetail) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	od.CreatedAt = now
	od.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (od *OrderDetail) BeforeUpdate(tx *gorm.DB) error {
	od.UpdatedAt = time.Now()
	return nil
}
