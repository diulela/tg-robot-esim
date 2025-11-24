package repository

import (
	"context"

	"tg-robot-sim/storage/models"

	"gorm.io/gorm"
)

// EsimCardRepository eSIM 卡仓储接口
type EsimCardRepository interface {
	// Create 创建 eSIM 卡记录
	Create(ctx context.Context, esimCard *models.EsimCard) error

	// GetByID 根据ID获取 eSIM 卡
	GetByID(ctx context.Context, id uint) (*models.EsimCard, error)

	// GetByICCID 根据ICCID获取 eSIM 卡
	GetByICCID(ctx context.Context, iccid string) (*models.EsimCard, error)

	// GetByUserID 获取用户的所有 eSIM 卡
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.EsimCard, error)

	// GetByUserIDWithFilters 根据筛选条件获取用户的 eSIM 卡
	GetByUserIDWithFilters(ctx context.Context, userID int64, status models.EsimStatus, limit, offset int) ([]*models.EsimCard, int64, error)

	// GetByOrderID 根据订单ID获取 eSIM 卡列表
	GetByOrderID(ctx context.Context, orderID uint) ([]*models.EsimCard, error)

	// Update 更新 eSIM 卡信息
	Update(ctx context.Context, esimCard *models.EsimCard) error

	// UpdateStatus 更新 eSIM 卡状态
	UpdateStatus(ctx context.Context, id uint, status models.EsimStatus) error

	// UpdateUsage 更新 eSIM 卡使用情况
	UpdateUsage(ctx context.Context, id uint, dataUsed, dataRemaining int, usagePercent string) error

	// Delete 删除 eSIM 卡（软删除）
	Delete(ctx context.Context, id uint) error

	// Count 统计 eSIM 卡数量
	Count(ctx context.Context, userID int64, status models.EsimStatus) (int64, error)
}

// esimCardRepository eSIM 卡仓储实现
type esimCardRepository struct {
	db *gorm.DB
}

// NewEsimCardRepository 创建 eSIM 卡仓储实例
func NewEsimCardRepository(db *gorm.DB) EsimCardRepository {
	return &esimCardRepository{db: db}
}

// Create 创建 eSIM 卡记录
func (r *esimCardRepository) Create(ctx context.Context, esimCard *models.EsimCard) error {
	return r.db.WithContext(ctx).Create(esimCard).Error
}

// GetByID 根据ID获取 eSIM 卡
func (r *esimCardRepository) GetByID(ctx context.Context, id uint) (*models.EsimCard, error) {
	var esimCard models.EsimCard
	err := r.db.WithContext(ctx).
		First(&esimCard, id).Error
	if err != nil {
		return nil, err
	}
	return &esimCard, nil
}

// GetByICCID 根据ICCID获取 eSIM 卡
func (r *esimCardRepository) GetByICCID(ctx context.Context, iccid string) (*models.EsimCard, error) {
	var esimCard models.EsimCard
	err := r.db.WithContext(ctx).
		Where("iccid = ?", iccid).
		First(&esimCard).Error
	if err != nil {
		return nil, err
	}
	return &esimCard, nil
}

// GetByUserID 获取用户的所有 eSIM 卡
func (r *esimCardRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.EsimCard, error) {
	var esimCards []*models.EsimCard
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&esimCards).Error
	return esimCards, err
}

// GetByUserIDWithFilters 根据筛选条件获取用户的 eSIM 卡
func (r *esimCardRepository) GetByUserIDWithFilters(ctx context.Context, userID int64, status models.EsimStatus, limit, offset int) ([]*models.EsimCard, int64, error) {
	var esimCards []*models.EsimCard
	var total int64

	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID)

	// 如果指定了状态，添加状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Model(&models.EsimCard{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query = query.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&esimCards).Error
	return esimCards, total, err
}

// GetByOrderID 根据订单ID获取 eSIM 卡列表
func (r *esimCardRepository) GetByOrderID(ctx context.Context, orderID uint) ([]*models.EsimCard, error) {
	var esimCards []*models.EsimCard
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&esimCards).Error
	return esimCards, err
}

// Update 更新 eSIM 卡信息
func (r *esimCardRepository) Update(ctx context.Context, esimCard *models.EsimCard) error {
	return r.db.WithContext(ctx).Save(esimCard).Error
}

// UpdateStatus 更新 eSIM 卡状态
func (r *esimCardRepository) UpdateStatus(ctx context.Context, id uint, status models.EsimStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.EsimCard{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateUsage 更新 eSIM 卡使用情况
func (r *esimCardRepository) UpdateUsage(ctx context.Context, id uint, dataUsed, dataRemaining int, usagePercent string) error {
	return r.db.WithContext(ctx).
		Model(&models.EsimCard{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"data_used":      dataUsed,
			"data_remaining": dataRemaining,
			"usage_percent":  usagePercent,
		}).Error
}

// Delete 删除 eSIM 卡（软删除）
func (r *esimCardRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.EsimCard{}, id).Error
}

// Count 统计 eSIM 卡数量
func (r *esimCardRepository) Count(ctx context.Context, userID int64, status models.EsimStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Model(&models.EsimCard{}).Count(&count).Error
	return count, err
}
