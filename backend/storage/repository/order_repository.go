package repository

import (
	"context"
	"time"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByIDs(ctx context.Context, id []uint) ([]*models.Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*models.Order, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID int64) (int64, error)
	GetUserOrderByID(ctx context.Context, userID int64, orderID uint) (*models.Order, error)

	// eSIM 订单处理相关方法
	// GetByIDWithDetail 根据ID获取订单（包含详情）
	GetByIDWithDetail(ctx context.Context, id uint) (*models.Order, error)

	// GetPendingSyncOrders 获取需要同步的待处理订单
	GetPendingSyncOrders(ctx context.Context, limit int) ([]*models.Order, error)

	// UpdateSyncInfo 更新订单同步信息
	UpdateSyncInfo(ctx context.Context, id uint, syncAttempts int, nextSyncAt *time.Time) error

	// GetByProviderOrderID 根据第三方订单ID获取订单
	GetByProviderOrderID(ctx context.Context, providerOrderID string) (*models.Order, error)

	// BatchUpdateStatus 批量更新订单状态
	BatchUpdateStatus(ctx context.Context, orderIDs []uint, status models.OrderStatus) error

	// GetByUserIDWithFilters 根据用户ID和筛选条件获取订单列表
	GetByUserIDWithFilters(ctx context.Context, userID int64, status models.OrderStatus, limit, offset int) ([]*models.Order, int64, error)
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

// GetByIDWithDetail 根据ID获取订单（包含详情）
func (r *orderRepository) GetByIDWithDetail(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetPendingSyncOrders 获取需要同步的待处理订单
func (r *orderRepository) GetPendingSyncOrders(ctx context.Context, limit int) ([]*models.Order, error) {
	var orders []*models.Order
	query := r.db.WithContext(ctx).
		Where("status = ? AND provider_order_id != ''", models.OrderStatusProcessing).
		Where("next_sync_at IS NULL OR next_sync_at <= ?", time.Now()).
		Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&orders).Error
	return orders, err
}

// UpdateSyncInfo 更新订单同步信息
func (r *orderRepository) UpdateSyncInfo(ctx context.Context, id uint, syncAttempts int, nextSyncAt *time.Time) error {
	updates := map[string]interface{}{
		"sync_attempts": syncAttempts,
		"last_sync_at":  time.Now(),
	}

	if nextSyncAt != nil {
		updates["next_sync_at"] = *nextSyncAt
	}

	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetByProviderOrderID 根据第三方订单ID获取订单
func (r *orderRepository) GetByProviderOrderID(ctx context.Context, providerOrderID string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Where("provider_order_id = ?", providerOrderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// BatchUpdateStatus 批量更新订单状态
func (r *orderRepository) BatchUpdateStatus(ctx context.Context, orderIDs []uint, status models.OrderStatus) error {
	if len(orderIDs) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"status": status,
	}

	// 根据状态设置相应的时间戳
	switch status {
	case models.OrderStatusCompleted:
		updates["completed_at"] = time.Now()
	case models.OrderStatusPaid:
		updates["paid_at"] = time.Now()
	}

	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id IN ?", orderIDs).
		Updates(updates).Error
}

// GetByUserIDWithFilters 根据用户ID和筛选条件获取订单列表
func (r *orderRepository) GetByUserIDWithFilters(ctx context.Context, userID int64, status models.OrderStatus, limit, offset int) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)

	// 如果指定了状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取订单列表
	query = query.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&orders).Error
	return orders, total, err
}

func (r *orderRepository) GetUserOrderByID(ctx context.Context, userID int64, orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, orderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByIDs(ctx context.Context, id []uint) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).
		Where("id IN ?", id).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
