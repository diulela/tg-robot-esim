package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// RechargeStatus 充值订单状态枚举
type RechargeStatus string

const (
	RechargeStatusPending   RechargeStatus = "pending"   // 待支付
	RechargeStatusConfirmed RechargeStatus = "confirmed" // 已确认
	RechargeStatusExpired   RechargeStatus = "expired"   // 已过期
	RechargeStatusFailed    RechargeStatus = "failed"    // 失败
)

// RechargeOrder 充值订单模型
type RechargeOrder struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo       string         `gorm:"uniqueIndex;size:32;not null" json:"order_no"`  // 订单号
	UserID        int64          `gorm:"index;not null" json:"user_id"`                 // 用户ID
	Amount        string         `gorm:"type:decimal(10,2);not null" json:"amount"`     // 充值金额
	WalletAddress string         `gorm:"size:100;not null" json:"wallet_address"`       // 充值地址
	Status        RechargeStatus `gorm:"size:20;default:'pending';index" json:"status"` // 订单状态
	TxHash        string         `gorm:"size:100;index" json:"tx_hash"`                 // 交易哈希
	Confirmations int            `gorm:"default:0" json:"confirmations"`                // 确认数
	Remark        string         `gorm:"type:text" json:"remark"`                       // 备注
	ExpiresAt     time.Time      `gorm:"index" json:"expires_at"`                       // 过期时间
	ConfirmedAt   *time.Time     `json:"confirmed_at,omitempty"`                        // 确认时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (RechargeOrder) TableName() string {
	return "recharge_orders"
}

// BeforeCreate GORM 钩子：创建前
func (r *RechargeOrder) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now

	// 生成订单号
	if r.OrderNo == "" {
		r.OrderNo = generateRechargeOrderNo()
	}

	// 设置默认过期时间（30分钟）
	if r.ExpiresAt.IsZero() {
		r.ExpiresAt = now.Add(30 * time.Minute)
	}

	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (r *RechargeOrder) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}

// IsConfirmed 检查订单是否已确认
func (r *RechargeOrder) IsConfirmed() bool {
	return r.Status == RechargeStatusConfirmed
}

// IsExpired 检查订单是否已过期
func (r *RechargeOrder) IsExpired() bool {
	if r.Status == RechargeStatusExpired {
		return true
	}
	return time.Now().After(r.ExpiresAt) && r.Status == RechargeStatusPending
}

// IsPending 检查订单是否待支付
func (r *RechargeOrder) IsPending() bool {
	return r.Status == RechargeStatusPending && !r.IsExpired()
}

// generateRechargeOrderNo 生成充值订单号
func generateRechargeOrderNo() string {
	// 格式: RCH + 时间戳 + 随机数
	return fmt.Sprintf("RCH%d%04d", time.Now().Unix(), time.Now().Nanosecond()%10000)
}
