package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"tg-robot-sim/storage/models"
)

// ProductRepository 产品仓储接口
type ProductRepository interface {
	// Create 创建产品
	Create(ctx context.Context, product *models.Product) error

	// Update 更新产品
	Update(ctx context.Context, product *models.Product) error

	// GetByID 根据ID获取产品
	GetByID(ctx context.Context, id int) (*models.Product, error)

	// GetByThirdPartyID 根据第三方ID获取产品
	GetByThirdPartyID(ctx context.Context, thirdPartyID string) (*models.Product, error)

	// List 获取产品列表
	List(ctx context.Context, params ListParams) ([]*models.Product, int64, error)

	// Upsert 创建或更新产品
	Upsert(ctx context.Context, product *models.Product) error

	// BatchUpsert 批量创建或更新产品
	BatchUpsert(ctx context.Context, products []*models.Product) error

	// Delete 删除产品
	Delete(ctx context.Context, id int) error

	// UpdateSyncTime 更新同步时间
	UpdateSyncTime(ctx context.Context, id int) error

	// FindByConditions 根据条件查询产品
	FindByConditions(ctx context.Context, conditions map[string]interface{}, limit, offset int) ([]*models.Product, error)

	// Count 统计产品数量
	Count(ctx context.Context, conditions map[string]interface{}) (int64, error)
}

// ListParams 列表查询参数
type ListParams struct {
	Type      string
	Status    string
	IsHot     *bool
	NameLike  string // 名称模糊搜索
	Page      int
	Limit     int
	OrderBy   string
	OrderDesc bool
}

// ProductModel 产品模型别名（用于避免循环导入）
type ProductModel = models.Product

// productRepository 产品仓储实现
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓储
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create 创建产品
func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// Update 更新产品
func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// GetByID 根据ID获取产品
func (r *productRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetByThirdPartyID 根据第三方ID获取产品
func (r *productRepository) GetByThirdPartyID(ctx context.Context, thirdPartyID string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("third_party_id = ?", thirdPartyID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List 获取产品列表
func (r *productRepository) List(ctx context.Context, params ListParams) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Product{})

	// 过滤条件
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.IsHot != nil {
		query = query.Where("is_hot = ?", *params.IsHot)
	}
	if params.NameLike != "" {
		query = query.Where("name LIKE ?", "%"+params.NameLike+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "sort_order"
	if params.OrderBy != "" {
		orderBy = params.OrderBy
	}
	if params.OrderDesc {
		orderBy += " DESC"
	}
	query = query.Order(orderBy)

	// 分页
	if params.Limit > 0 {
		offset := 0
		if params.Page > 1 {
			offset = (params.Page - 1) * params.Limit
		}
		query = query.Offset(offset).Limit(params.Limit)
	}

	err := query.Find(&products).Error
	return products, total, err
}

// Upsert 创建或更新产品
func (r *productRepository) Upsert(ctx context.Context, product *models.Product) error {
	// 尝试查找现有产品
	existing, err := r.GetByThirdPartyID(ctx, product.ThirdPartyID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existing != nil {
		// 更新现有产品
		product.ID = existing.ID
		product.CreatedAt = existing.CreatedAt
		return r.Update(ctx, product)
	}

	// 创建新产品
	return r.Create(ctx, product)
}

// BatchUpsert 批量创建或更新产品
func (r *productRepository) BatchUpsert(ctx context.Context, products []*models.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			repo := &productRepository{db: tx}
			if err := repo.Upsert(ctx, product); err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete 删除产品
func (r *productRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

// UpdateSyncTime 更新同步时间
func (r *productRepository) UpdateSyncTime(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", id).
		Update("synced_at", time.Now()).Error
}

// FindByConditions 根据条件查询产品
func (r *productRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}, limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	query := r.db.WithContext(ctx).Model(&models.Product{})

	// 应用条件
	for key, value := range conditions {
		switch key {
		case "country":
			// 处理国家筛选，使用 LIKE 查询 JSON 字段
			if countryCode, ok := value.(string); ok && countryCode != "" {
				query = query.Where("countries LIKE ?", "%"+countryCode+"%")
			}
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	// 排序
	query = query.Order("sort_order ASC, created_at DESC")

	// 分页
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&products).Error
	return products, err
}

// Count 统计产品数量
func (r *productRepository) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Product{})

	// 应用条件
	for key, value := range conditions {
		switch key {
		case "country":
			// 处理国家筛选，使用 LIKE 查询 JSON 字段
			if countryCode, ok := value.(string); ok && countryCode != "" {
				query = query.Where("countries LIKE ?", "%"+countryCode+"%")
			}
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	err := query.Count(&count).Error
	return count, err
}
