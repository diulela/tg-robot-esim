package models

import (
	"time"

	"gorm.io/gorm"
)

// ProductDetail 产品详情模型（基于API返回结构）
type ProductDetail struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID    int            `gorm:"uniqueIndex;not null" json:"product_id"` // 关联产品表ID
	ThirdPartyID int            `gorm:"index" json:"third_party_id"`            // API返回的产品ID
	Name         string         `gorm:"size:200;not null" json:"name"`
	Type         string         `gorm:"size:20;index" json:"type"`  // local, regional, global
	Countries    string         `gorm:"type:text" json:"countries"` // JSON 数组格式，如 ["US", "CN"]
	DataSize     string         `gorm:"size:50" json:"data_size"`   // 如 "5GB", "10GB"
	ValidDays    int            `gorm:"not null" json:"valid_days"`
	Price        float64        `gorm:"type:decimal(10,2)" json:"price"`
	CostPrice    float64        `gorm:"type:decimal(10,2)" json:"cost_price"`
	Description  string         `gorm:"type:text" json:"description"`
	Features     string         `gorm:"type:text" json:"features"` // JSON 数组格式
	Status       string         `gorm:"size:20;default:'active'" json:"status"`
	ApiCreatedAt string         `gorm:"size:50" json:"api_created_at"` // API返回的创建时间
	SyncedAt     *time.Time     `gorm:"index" json:"synced_at"`        // 最后同步时间
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ProductDetail) TableName() string {
	return "product_details"
}
