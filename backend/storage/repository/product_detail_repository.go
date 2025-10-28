package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"tg-robot-sim/storage/models"
)

// ProductDetailRepository 产品详情仓储接口
type ProductDetailRepository interface {
	// Create 创建产品详情
	Create(ctx context.Context, detail *models.ProductDetail) error

	// Update 更新产品详情
	Update(ctx context.Context, detail *models.ProductDetail) error

	// GetByProductID 根据产品ID获取详情
	GetByProductID(ctx context.Context, productID int) (*models.ProductDetail, error)

	// GetByThirdPartyID 根据第三方ID获取详情
	GetByThirdPartyID(ctx context.Context, thirdPartyID int) (*models.ProductDetail, error)

	// Upsert 创建或更新产品详情
	Upsert(ctx context.Context, detail *models.ProductDetail) error

	// Delete 删除产品详情
	Delete(ctx context.Context, productID int) error

	// UpdateSyncTime 更新同步时间
	UpdateSyncTime(ctx context.Context, productID int) error
}

// productDetailRepository 产品详情仓储实现
type productDetailRepository struct {
	db *gorm.DB
}

// NewProductDetailRepository 创建产品详情仓储
func NewProductDetailRepository(db *gorm.DB) ProductDetailRepository {
	return &productDetailRepository{db: db}
}

// Create 创建产品详情
func (r *productDetailRepository) Create(ctx context.Context, detail *models.ProductDetail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

// Update 更新产品详情
func (r *productDetailRepository) Update(ctx context.Context, detail *models.ProductDetail) error {
	return r.db.WithContext(ctx).Save(detail).Error
}

// GetByProductID 根据产品ID获取详情
func (r *productDetailRepository) GetByProductID(ctx context.Context, productID int) (*models.ProductDetail, error) {
	var detail models.ProductDetail
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// GetByThirdPartyID 根据第三方ID获取详情
func (r *productDetailRepository) GetByThirdPartyID(ctx context.Context, thirdPartyID int) (*models.ProductDetail, error) {
	var detail models.ProductDetail
	err := r.db.WithContext(ctx).Where("third_party_id = ?", thirdPartyID).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// Upsert 创建或更新产品详情
func (r *productDetailRepository) Upsert(ctx context.Context, detail *models.ProductDetail) error {
	// 尝试查找现有详情
	existing, err := r.GetByProductID(ctx, detail.ProductID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		// 更新现有详情
		detail.ID = existing.ID
		detail.CreatedAt = existing.CreatedAt
		return r.Update(ctx, detail)
	}

	// 创建新详情
	return r.Create(ctx, detail)
}

// Delete 删除产品详情
func (r *productDetailRepository) Delete(ctx context.Context, productID int) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&models.ProductDetail{}).Error
}

// UpdateSyncTime 更新同步时间
func (r *productDetailRepository) UpdateSyncTime(ctx context.Context, productID int) error {
	return r.db.WithContext(ctx).Model(&models.ProductDetail{}).
		Where("product_id = ?", productID).
		Update("synced_at", time.Now()).Error
}
