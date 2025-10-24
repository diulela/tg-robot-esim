package services

import (
	"context"
	"fmt"

	"tg-robot-sim/pkg/sdk/esim"
)

// EsimService eSIM产品服务接口
type EsimService interface {
	// GetProducts 获取产品列表
	GetProducts(ctx context.Context, params *esim.ProductParams) (*esim.ProductListResponse, error)

	// GetProduct 获取产品详情
	GetProduct(ctx context.Context, productID int) (*esim.ProductDetailResponse, error)

	// CreateOrder 创建订单
	CreateOrder(ctx context.Context, req esim.CreateOrderRequest) (*esim.CreateOrderResponse, error)

	// GetOrder 获取订单详情
	GetOrder(ctx context.Context, orderID int) (*esim.OrderDetailResponse, error)

	// GetEsimUsage 获取eSIM使用情况
	GetEsimUsage(ctx context.Context, orderID int) (*esim.EsimUsageResponse, error)

	// GetTopupPackages 获取充值套餐
	GetTopupPackages(ctx context.Context, orderID int) (*esim.TopupPackagesResponse, error)

	// TopupEsim eSIM充值
	TopupEsim(ctx context.Context, orderID int, req esim.TopupRequest) (*esim.TopupResponse, error)
}

// esimServiceImpl eSIM服务实现
type esimServiceImpl struct {
	client *esim.Client
}

// NewEsimService 创建eSIM服务
func NewEsimService(apiKey, apiSecret, baseURL string) EsimService {
	client := esim.NewClient(esim.Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
	})

	return &esimServiceImpl{
		client: client,
	}
}

// GetProducts 获取产品列表
func (s *esimServiceImpl) GetProducts(ctx context.Context, params *esim.ProductParams) (*esim.ProductListResponse, error) {
	return s.client.GetProducts(params)
}

// GetProduct 获取产品详情
func (s *esimServiceImpl) GetProduct(ctx context.Context, productID int) (*esim.ProductDetailResponse, error) {
	return s.client.GetProduct(productID)
}

// CreateOrder 创建订单
func (s *esimServiceImpl) CreateOrder(ctx context.Context, req esim.CreateOrderRequest) (*esim.CreateOrderResponse, error) {
	return s.client.CreateOrder(req)
}

// GetOrder 获取订单详情
func (s *esimServiceImpl) GetOrder(ctx context.Context, orderID int) (*esim.OrderDetailResponse, error) {
	return s.client.GetOrder(orderID)
}

// GetEsimUsage 获取eSIM使用情况
func (s *esimServiceImpl) GetEsimUsage(ctx context.Context, orderID int) (*esim.EsimUsageResponse, error) {
	return s.client.GetEsimUsage(orderID)
}

// GetTopupPackages 获取充值套餐
func (s *esimServiceImpl) GetTopupPackages(ctx context.Context, orderID int) (*esim.TopupPackagesResponse, error) {
	return s.client.GetTopupPackages(orderID)
}

// TopupEsim eSIM充值
func (s *esimServiceImpl) TopupEsim(ctx context.Context, orderID int, req esim.TopupRequest) (*esim.TopupResponse, error) {
	return s.client.TopupEsim(orderID, req)
}

// FormatProductMessage 格式化产品信息为消息文本
func FormatProductMessage(product *esim.Product) string {
	msg := fmt.Sprintf("📱 *%s*\n\n", product.Name)
	msg += fmt.Sprintf("🌍 类型: %s\n", getProductTypeText(product.Type))

	// 国家列表
	if len(product.Countries) > 0 {
		msg += "🗺️ 支持国家: "
		for i, country := range product.Countries {
			if i > 0 {
				msg += ", "
			}
			msg += country.CN
			if i >= 2 {
				msg += fmt.Sprintf(" 等%d个国家", len(product.Countries))
				break
			}
		}
		msg += "\n"
	}

	// 流量和有效期
	msg += fmt.Sprintf("📊 流量: %s\n", formatDataSize(product.DataSize))
	msg += fmt.Sprintf("⏰ 有效期: %d天\n", product.ValidDays)

	// 价格
	msg += fmt.Sprintf("💰 零售价: $%.2f\n", product.RetailPrice)
	msg += fmt.Sprintf("💵 代理价: $%.2f\n", product.AgentPrice)

	// 特性
	if len(product.Features) > 0 {
		msg += "\n✨ 特性:\n"
		for _, feature := range product.Features {
			msg += fmt.Sprintf("  • %s\n", feature)
		}
	}

	return msg
}

// getProductTypeText 获取产品类型文本
func getProductTypeText(productType esim.ProductType) string {
	switch productType {
	case esim.ProductTypeLocal:
		return "本地"
	case esim.ProductTypeRegional:
		return "区域"
	case esim.ProductTypeGlobal:
		return "全球"
	default:
		return string(productType)
	}
}

// formatDataSize 格式化流量大小
func formatDataSize(sizeMB int) string {
	if sizeMB >= 1024 {
		return fmt.Sprintf("%.1fGB", float64(sizeMB)/1024)
	}
	return fmt.Sprintf("%dMB", sizeMB)
}
