package repository

import (
	"context"
	"time"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// RechargeOrderRepository 充值订单仓储接口
type RechargeOrderRepository interface {
	Create(ctx context.Context, order *models.RechargeOrder) error
	GetByID(ctx context.Context, id uint) (*models.RechargeOrder, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*models.RechargeOrder, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.RechargeOrder, error)
	GetByTxHash(ctx context.Context, txHash string) (*models.RechargeOrder, error)
	GetPendingOrders(ctx context.Context) ([]*models.RechargeOrder, error)
	Update(ctx context.Context, order *models.RechargeOrder) error
	UpdateStatus(ctx context.Context, id uint, status models.RechargeStatus) error
	Delete(ctx context.Context, id uint) error
	ExpireOldOrders(ctx context.Context) error
}

// rechargeOrderRepository 充值订单仓储实现
type rechargeOrderRepository struct {
	db *gorm.DB
}

// NewRechargeOrderRepository 创建充值订单仓储实例
func NewRechargeOrderRepository(db *gorm.DB) RechargeOrderRepository {
	return &rechargeOrderRepository{db: db}
}

// Create 创建充值订单
func (r *rechargeOrderRepository) Create(ctx context.Context, order *models.RechargeOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID 根据ID获取充值订单
func (r *rechargeOrderRepository) GetByID(ctx context.Context, id uint) (*models.RechargeOrder, error) {
	var order models.RechargeOrder
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取充值订单
func (r *rechargeOrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*models.RechargeOrder, error) {
	var order models.RechargeOrder
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("order_no = ?", orderNo).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID 根据用户ID获取充值订单列表
func (r *rechargeOrderRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.RechargeOrder, error) {
	var orders []*models.RechargeOrder
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&orders).Error
	return orders, err
}

// GetByTxHash 根据交易哈希获取充值订单
func (r *rechargeOrderRepository) GetByTxHash(ctx context.Context, txHash string) (*models.RechargeOrder, error) {
	var order models.RechargeOrder
	err := r.db.WithContext(ctx).
		Where("tx_hash = ?", txHash).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetPendingOrders 获取待处理的充值订单
func (r *rechargeOrderRepository) GetPendingOrders(ctx context.Context) ([]*models.RechargeOrder, error) {
	var orders []*models.RechargeOrder
	err := r.db.WithContext(ctx).
		Where("status = ?", models.RechargeStatusPending).
		Where("expires_at > ?", time.Now()).
		Order("created_at ASC").
		Find(&orders).Error
	return orders, err
}

// Update 更新充值订单
func (r *rechargeOrderRepository) Update(ctx context.Context, order *models.RechargeOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus 更新充值订单状态
func (r *rechargeOrderRepository) UpdateStatus(ctx context.Context, id uint, status models.RechargeStatus) error {
	return r.db.WithContext(ctx).Model(&models.RechargeOrder{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除充值订单
func (r *rechargeOrderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.RechargeOrder{}, id).Error
}

// ExpireOldOrders 将过期的订单标记为已过期
func (r *rechargeOrderRepository) ExpireOldOrders(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&models.RechargeOrder{}).
		Where("status = ?", models.RechargeStatusPending).
		Where("expires_at < ?", time.Now()).
		Update("status", models.RechargeStatusExpired).Error
}
