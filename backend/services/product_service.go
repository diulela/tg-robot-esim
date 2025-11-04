package services

import (
	"context"
	"fmt"
	"strings"

	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

// ProductFilters 产品筛选条件
type ProductFilters struct {
	Type        string // local, regional, global
	Country     string // 国家代码筛选
	Search      string // 搜索关键词
	IsHot       *bool  // 是否热门
	IsRecommend *bool  // 是否推荐
	MinPrice    float64
	MaxPrice    float64
	Limit       int
	Offset      int
}

// ProductService 产品服务接口
type ProductService interface {
	GetProducts(ctx context.Context, filters ProductFilters) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	SearchProducts(ctx context.Context, query string, limit int) ([]*models.Product, error)
	GetHotProducts(ctx context.Context, limit int) ([]*models.Product, error)
	GetRecommendedProducts(ctx context.Context, limit int) ([]*models.Product, error)
	GetProductsByType(ctx context.Context, productType string, limit, offset int) ([]*models.Product, error)
	CountProducts(ctx context.Context, filters ProductFilters) (int64, error)
}

// productService 产品服务实现
type productService struct {
	productRepo repository.ProductRepository
}

// NewProductService 创建产品服务实例
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

// GetProducts 获取产品列表（带筛选）
func (s *productService) GetProducts(ctx context.Context, filters ProductFilters) ([]*models.Product, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})

	if filters.Type != "" && filters.Type != "all" {
		conditions["type"] = filters.Type
	}

	if filters.Country != "" {
		conditions["country"] = filters.Country
	}

	if filters.IsHot != nil {
		conditions["is_hot"] = *filters.IsHot
	}

	if filters.IsRecommend != nil {
		conditions["is_recommend"] = *filters.IsRecommend
	}

	conditions["status"] = "active"

	// 获取产品列表
	products, err := s.productRepo.FindByConditions(ctx, conditions, filters.Limit, filters.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// 如果有搜索关键词，进行过滤
	if filters.Search != "" {
		products = s.filterBySearch(products, filters.Search)
	}

	// 如果有价格范围，进行过滤
	if filters.MinPrice > 0 || filters.MaxPrice > 0 {
		products = s.filterByPrice(products, filters.MinPrice, filters.MaxPrice)
	}

	return products, nil
}

// GetProductByID 根据ID获取产品
func (s *productService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return product, nil
}

// SearchProducts 搜索产品
func (s *productService) SearchProducts(ctx context.Context, query string, limit int) ([]*models.Product, error) {
	if query == "" {
		return s.GetProducts(ctx, ProductFilters{Limit: limit})
	}

	// 获取所有活跃产品
	products, err := s.productRepo.FindByConditions(ctx, map[string]interface{}{
		"status": "active",
	}, 0, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	// 过滤搜索结果
	results := s.filterBySearch(products, query)

	// 限制结果数量
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetHotProducts 获取热门产品
func (s *productService) GetHotProducts(ctx context.Context, limit int) ([]*models.Product, error) {
	products, err := s.productRepo.FindByConditions(ctx, map[string]interface{}{
		"is_hot": true,
		"status": "active",
	}, limit, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to get hot products: %w", err)
	}

	return products, nil
}

// GetRecommendedProducts 获取推荐产品
func (s *productService) GetRecommendedProducts(ctx context.Context, limit int) ([]*models.Product, error) {
	products, err := s.productRepo.FindByConditions(ctx, map[string]interface{}{
		"is_recommend": true,
		"status":       "active",
	}, limit, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to get recommended products: %w", err)
	}

	return products, nil
}

// GetProductsByType 根据类型获取产品
func (s *productService) GetProductsByType(ctx context.Context, productType string, limit, offset int) ([]*models.Product, error) {
	conditions := map[string]interface{}{
		"status": "active",
	}

	if productType != "" && productType != "all" {
		conditions["type"] = productType
	}

	products, err := s.productRepo.FindByConditions(ctx, conditions, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by type: %w", err)
	}

	return products, nil
}

// CountProducts 统计产品数量
func (s *productService) CountProducts(ctx context.Context, filters ProductFilters) (int64, error) {
	conditions := make(map[string]interface{})

	if filters.Type != "" && filters.Type != "all" {
		conditions["type"] = filters.Type
	}

	if filters.Country != "" {
		conditions["country"] = filters.Country
	}

	if filters.IsHot != nil {
		conditions["is_hot"] = *filters.IsHot
	}

	if filters.IsRecommend != nil {
		conditions["is_recommend"] = *filters.IsRecommend
	}

	conditions["status"] = "active"

	count, err := s.productRepo.Count(ctx, conditions)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// filterBySearch 根据搜索关键词过滤产品
func (s *productService) filterBySearch(products []*models.Product, query string) []*models.Product {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return products
	}

	var results []*models.Product
	for _, product := range products {
		// 搜索产品名称和描述
		if strings.Contains(strings.ToLower(product.Name), query) ||
			strings.Contains(strings.ToLower(product.NameEn), query) ||
			strings.Contains(strings.ToLower(product.Description), query) ||
			strings.Contains(strings.ToLower(product.DescriptionEn), query) {
			results = append(results, product)
		}
	}

	return results
}

// filterByPrice 根据价格范围过滤产品
func (s *productService) filterByPrice(products []*models.Product, minPrice, maxPrice float64) []*models.Product {
	var results []*models.Product

	for _, product := range products {
		if minPrice > 0 && product.Price < minPrice {
			continue
		}
		if maxPrice > 0 && product.Price > maxPrice {
			continue
		}
		results = append(results, product)
	}

	return results
}
