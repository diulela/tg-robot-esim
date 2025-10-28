package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态枚举
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusCompleted OrderStatus = "completed" // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

// Order 订单模型
type Order struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo     string         `gorm:"uniqueIndex;size:32;not null" json:"order_no"`  // 订单号
	UserID      int64          `gorm:"index;not null" json:"user_id"`                 // 用户ID
	ProductID   int            `gorm:"index;not null" json:"product_id"`              // 产品ID
	ProductName string         `gorm:"size:200" json:"product_name"`                  // 产品名称（冗余）
	Amount      string         `gorm:"type:decimal(10,2);not null" json:"amount"`     // 订单金额
	Status      OrderStatus    `gorm:"size:20;default:'pending';index" json:"status"` // 订单状态
	TxHash      string         `gorm:"size:100;index" json:"tx_hash"`                 // 交易哈希
	Remark      string         `gorm:"type:text" json:"remark"`                       // 备注
	PaidAt      *time.Time     `json:"paid_at,omitempty"`                             // 支付时间
	CompletedAt *time.Time     `json:"completed_at,omitempty"`                        // 完成时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// BeforeCreate GORM 钩子：创建前
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now

	// 生成订单号
	if o.OrderNo == "" {
		o.OrderNo = generateOrderNo()
	}

	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid || o.Status == OrderStatusCompleted
}

// IsCompleted 检查订单是否已完成
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusCompleted
}

// IsCancelled 检查订单是否已取消
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	// 格式: ORD + 时间戳 + 随机数
	return fmt.Sprintf("ORD%d%04d", time.Now().Unix(), time.Now().Nanosecond()%10000)
}
