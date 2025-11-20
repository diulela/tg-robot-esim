package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OrderStatus 订单状态枚举
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // 待支付
	OrderStatusPaid       OrderStatus = "paid"       // 已支付
	OrderStatusProcessing OrderStatus = "processing" // 处理中（包含创建、支付、第三方处理等所有中间状态，余额已冻结）
	OrderStatusCompleted  OrderStatus = "completed"  // 已完成（冻结金额已扣除，eSIM 信息已获取）
	OrderStatusCancelled  OrderStatus = "cancelled"  // 已取消
	OrderStatusRefunded   OrderStatus = "refunded"   // 已退款
	OrderStatusFailed     OrderStatus = "failed"     // 失败（冻结金额已退还）
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
	Remark      string         `gorm:"type:text" json:"remark"`                       // 备注
	PaidAt      *time.Time     `json:"paid_at,omitempty"`                             // 支付时间
	CompletedAt *time.Time     `json:"completed_at,omitempty"`                        // 完成时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// eSIM 订单处理新增字段
	Quantity        int        `gorm:"default:1" json:"quantity"`               // 购买数量
	UnitPrice       string     `gorm:"type:decimal(10,4)" json:"unit_price"`    // 单价
	ProviderOrderID string     `gorm:"size:100;index" json:"provider_order_id"` // 第三方订单ID
	ProviderOrderNo string     `gorm:"size:100;index" json:"provider_order_no"` // 第三方订单号
	SyncAttempts    int        `gorm:"default:0" json:"sync_attempts"`          // 同步尝试次数
	LastSyncAt      *time.Time `gorm:"index;type:datetime" json:"last_sync_at"` // 最后同步时间
	NextSyncAt      *time.Time `gorm:"index;type:datetime" json:"next_sync_at"` // 下次同步时间
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
