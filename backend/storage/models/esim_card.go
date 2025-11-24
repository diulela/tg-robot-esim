package models

import (
	"time"

	"gorm.io/gorm"
)

// EsimStatus eSIM 卡状态
type EsimStatus string

const (
	EsimStatusPending    EsimStatus = "pending"    // 待激活
	EsimStatusActive     EsimStatus = "active"     // 已激活
	EsimStatusExpired    EsimStatus = "expired"    // 已过期
	EsimStatusSuspended  EsimStatus = "suspended"  // 已暂停
	EsimStatusTerminated EsimStatus = "terminated" // 已终止
)

// EsimCard eSIM 卡模型
type EsimCard struct {
	ID      uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID  int64 `gorm:"index;not null" json:"user_id"`  // 用户ID
	OrderID uint  `gorm:"index;not null" json:"order_id"` // 购买订单ID

	// eSIM 基本信息
	ICCID          string `gorm:"uniqueIndex;size:50;not null" json:"iccid"` // ICCID号码
	ActivationCode string `gorm:"type:text" json:"activation_code"`          // 激活码
	QrCode         string `gorm:"type:text" json:"qr_code"`                  // 二维码
	Lpa            string `gorm:"type:text" json:"lpa"`                      // LPA
	DirectAppleUrl string `gorm:"type:text" json:"direct_apple_url"`         // Apple直接安装链接
	ApnType        string `gorm:"size:20" json:"apn_type"`                   // APN类型 (manual/automatic)
	IsRoaming      bool   `gorm:"default:false" json:"is_roaming"`           // 是否漫游

	// 使用情况
	Status        EsimStatus `gorm:"size:20;not null;index" json:"status"` // eSIM状态
	DataUsed      int        `gorm:"default:0" json:"data_used"`           // 已使用流量（MB）
	DataRemaining int        `gorm:"default:0" json:"data_remaining"`      // 剩余流量（MB）
	UsagePercent  string     `gorm:"size:10" json:"usage_percent"`         // 使用百分比

	// 时间信息
	ActivatedAt *time.Time     `gorm:"type:datetime" json:"activated_at"` // 激活时间
	ExpiresAt   *time.Time     `gorm:"type:datetime" json:"expires_at"`   // 过期时间
	LastSyncAt  *time.Time     `gorm:"type:datetime" json:"last_sync_at"` // 最后同步时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 第三方数据
	ProviderOrderID int    `gorm:"index" json:"provider_order_id"`     // 第三方订单ID
	ProviderData    string `gorm:"type:longtext" json:"provider_data"` // 第三方完整数据（JSON）
}

// TableName 指定表名
func (EsimCard) TableName() string {
	return "esim_cards"
}

// BeforeCreate GORM 钩子：创建前
func (e *EsimCard) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (e *EsimCard) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}

// IsActive 检查 eSIM 卡是否激活
func (e *EsimCard) IsActive() bool {
	return e.Status == EsimStatusActive
}

// IsExpired 检查 eSIM 卡是否过期
func (e *EsimCard) IsExpired() bool {
	return e.Status == EsimStatusExpired
}

// CanSync 检查 eSIM 卡是否可以同步
func (e *EsimCard) CanSync() bool {
	return e.Status != EsimStatusTerminated && e.Status != EsimStatusExpired
}
