package repository

import (
	"context"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID int64) (int64, error)
}

// orderRepository 订单仓储实现
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create 创建订单
func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID 根据ID获取订单
func (r *orderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单
func (r *orderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Where("order_no = ?", orderNo).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID 根据用户ID获取订单列表
func (r *orderRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Order, error) {
	var orders []*models.Order
	query := r.db.WithContext(ctx).
		Preload("Product").
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

// Update 更新订单
func (r *orderRepository) Update(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus 更新订单状态
func (r *orderRepository) UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除订单
func (r *orderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Order{}, id).Error
}

// CountByUserID 统计用户订单数量
func (r *orderRepository) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
