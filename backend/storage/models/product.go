package models

import (
	"time"

	"gorm.io/gorm"
)

// Product 产品模型
type Product struct {
	ID             int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ThirdPartyID   string         `gorm:"uniqueIndex;size:100" json:"third_party_id"` // 第三方产品ID
	Name           string         `gorm:"size:200;not null" json:"name"`
	NameEn         string         `gorm:"size:200" json:"name_en"`
	Description    string         `gorm:"type:text" json:"description"`
	DescriptionEn  string         `gorm:"type:text" json:"description_en"`
	Type           string         `gorm:"size:20;index" json:"type"`  // local, regional, global
	Countries      string         `gorm:"type:text" json:"countries"` // JSON 格式存储国家列表
	DataSize       int            `gorm:"not null" json:"data_size"`  // MB
	ValidDays      int            `gorm:"not null" json:"valid_days"`
	Features       string         `gorm:"type:text" json:"features"` // JSON 格式存储特性列表
	Image          string         `gorm:"size:500" json:"image"`
	Price          float64        `gorm:"type:decimal(10,2)" json:"price"`
	CostPrice      float64        `gorm:"type:decimal(10,2)" json:"cost_price"`
	RetailPrice    float64        `gorm:"type:decimal(10,2)" json:"retail_price"`
	AgentPrice     float64        `gorm:"type:decimal(10,2)" json:"agent_price"`
	PlatformProfit float64        `gorm:"type:decimal(10,2)" json:"platform_profit"`
	IsHot          bool           `gorm:"default:false" json:"is_hot"`
	IsRecommend    bool           `gorm:"default:false" json:"is_recommend"`
	SortOrder      int            `gorm:"default:0" json:"sort_order"`
	Status         string         `gorm:"size:20;default:'active'" json:"status"` // active, inactive
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}
