package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// OrderDetailRepository 订单详情仓储接口
type OrderDetailRepository interface {
	Create(ctx context.Context, orderDetail *models.OrderDetail) error
	GetByOrderID(ctx context.Context, orderID uint) (*models.OrderDetail, error)
	Update(ctx context.Context, orderDetail *models.OrderDetail) error
	Delete(ctx context.Context, id uint) error

	// JSON 数据处理方法
	SaveProviderData(ctx context.Context, orderID uint, providerData interface{}) error
	SaveOrderItems(ctx context.Context, orderID uint, orderItems interface{}) error
	SaveEsims(ctx context.Context, orderID uint, esims interface{}) error

	// 批量操作
	CreateOrUpdate(ctx context.Context, orderDetail *models.OrderDetail) error
}

// orderDetailRepository 订单详情仓储实现
type orderDetailRepository struct {
	db *gorm.DB
}

// NewOrderDetailRepository 创建订单详情仓储实例
func NewOrderDetailRepository(db *gorm.DB) OrderDetailRepository {
	return &orderDetailRepository{db: db}
}

// Create 创建订单详情
func (r *orderDetailRepository) Create(ctx context.Context, orderDetail *models.OrderDetail) error {
	return r.db.WithContext(ctx).Create(orderDetail).Error
}

// GetByOrderID 根据订单ID获取订单详情
func (r *orderDetailRepository) GetByOrderID(ctx context.Context, orderID uint) (*models.OrderDetail, error) {
	var orderDetail models.OrderDetail
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&orderDetail).Error
	if err != nil {
		return nil, err
	}
	return &orderDetail, nil
}

// Update 更新订单详情
func (r *orderDetailRepository) Update(ctx context.Context, orderDetail *models.OrderDetail) error {
	return r.db.WithContext(ctx).Save(orderDetail).Error
}

// Delete 删除订单详情
func (r *orderDetailRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OrderDetail{}, id).Error
}

// SaveProviderData 保存第三方完整数据
func (r *orderDetailRepository) SaveProviderData(ctx context.Context, orderID uint, providerData interface{}) error {
	jsonData, err := json.Marshal(providerData)
	if err != nil {
		return fmt.Errorf("failed to marshal provider data: %w", err)
	}

	return r.db.WithContext(ctx).Model(&models.OrderDetail{}).
		Where("order_id = ?", orderID).
		Update("provider_data", string(jsonData)).Error
}

// SaveOrderItems 保存订单项数据
func (r *orderDetailRepository) SaveOrderItems(ctx context.Context, orderID uint, orderItems interface{}) error {
	jsonData, err := json.Marshal(orderItems)
	if err != nil {
		return fmt.Errorf("failed to marshal order items: %w", err)
	}

	return r.db.WithContext(ctx).Model(&models.OrderDetail{}).
		Where("order_id = ?", orderID).
		Update("order_items", string(jsonData)).Error
}

// SaveEsims 保存 eSIM 数据
func (r *orderDetailRepository) SaveEsims(ctx context.Context, orderID uint, esims interface{}) error {
	jsonData, err := json.Marshal(esims)
	if err != nil {
		return fmt.Errorf("failed to marshal esims: %w", err)
	}

	return r.db.WithContext(ctx).Model(&models.OrderDetail{}).
		Where("order_id = ?", orderID).
		Update("esims", string(jsonData)).Error
}

// CreateOrUpdate 创建或更新订单详情
func (r *orderDetailRepository) CreateOrUpdate(ctx context.Context, orderDetail *models.OrderDetail) error {
	// 先尝试获取现有记录
	existing, err := r.GetByOrderID(ctx, orderDetail.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			return r.Create(ctx, orderDetail)
		}
		return fmt.Errorf("failed to check existing order detail: %w", err)
	}

	// 记录存在，更新现有记录
	orderDetail.ID = existing.ID
	return r.Update(ctx, orderDetail)
}

// ParseProviderData 解析第三方数据（辅助方法）
func (r *orderDetailRepository) ParseProviderData(orderDetail *models.OrderDetail, target interface{}) error {
	if orderDetail.ProviderData == "" {
		return fmt.Errorf("provider data is empty")
	}

	return json.Unmarshal([]byte(orderDetail.ProviderData), target)
}

// ParseOrderItems 解析订单项数据（辅助方法）
func (r *orderDetailRepository) ParseOrderItems(orderDetail *models.OrderDetail, target interface{}) error {
	if orderDetail.OrderItems == "" {
		return fmt.Errorf("order items data is empty")
	}

	return json.Unmarshal([]byte(orderDetail.OrderItems), target)
}

// ParseEsims 解析 eSIM 数据（辅助方法）
func (r *orderDetailRepository) ParseEsims(orderDetail *models.OrderDetail, target interface{}) error {
	if orderDetail.Esims == "" {
		return fmt.Errorf("esims data is empty")
	}

	return json.Unmarshal([]byte(orderDetail.Esims), target)
}
