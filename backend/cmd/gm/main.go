package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"tg-robot-sim/config"
	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/storage/data"
	"tg-robot-sim/storage/models"
	"tg-robot-sim/storage/repository"
)

const (
	cmdSyncProducts       = "sync-products"
	cmdSyncProductDetails = "sync-product-details"
	cmdListProducts       = "list-products"
	cmdHelp               = "help"
)

func main() {
	// 定义命令行参数
	command := flag.String("cmd", "", "命令: sync-products, list-products, sync-product-details, help")
	configPath := flag.String("config", "config/config.json", "配置文件路径")
	productType := flag.String("type", "", "产品类型: local, regional, global (可选)")
	limit := flag.Int("limit", 0, "限制数量 (0 表示全部)")
	flag.Parse()

	if *command == "" || *command == cmdHelp {
		printHelp()
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := data.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer db.Close()

	// 自动迁移
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 执行命令
	ctx := context.Background()
	switch *command {
	case cmdSyncProducts:
		if err := syncProducts(ctx, cfg, db, *productType, *limit); err != nil {
			log.Fatalf("同步产品失败: %v", err)
		}
	case cmdSyncProductDetails:
		if err := syncProductDetails(ctx, cfg, db, *limit); err != nil {
			log.Fatalf("同步产品详情失败: %v", err)
		}
	case cmdListProducts:
		if err := listProducts(ctx, db, *productType); err != nil {
			log.Fatalf("列出产品失败: %v", err)
		}
	default:
		fmt.Printf("未知命令: %s\n", *command)
		printHelp()
		os.Exit(1)
	}
}

// syncProducts 同步产品数据
func syncProducts(ctx context.Context, cfg *config.Config, db *data.Database, productType string, limit int) error {
	fmt.Println("开始同步产品数据...")
	fmt.Printf("配置: API=%s, 类型=%s, 限制=%d\n", cfg.EsimSDK.BaseURL, productType, limit)

	// 创建 eSIM 客户端
	client := esim.NewClient(esim.Config{
		APIKey:         cfg.EsimSDK.APIKey,
		APISecret:      cfg.EsimSDK.APISecret,
		BaseURL:        cfg.EsimSDK.BaseURL,
		TimezoneOffset: cfg.EsimSDK.TimezoneOffset,
	})

	productRepo := db.GetProductRepository()

	// 同步不同类型的产品
	types := []esim.ProductType{esim.ProductTypeLocal, esim.ProductTypeRegional, esim.ProductTypeGlobal}
	if productType != "" {
		types = []esim.ProductType{esim.ProductType(productType)}
	}

	totalSynced := 0
	totalFailed := 0

	for _, pType := range types {
		fmt.Printf("\n正在同步 %s 产品...\n", pType)

		page := 1
		for {
			// 获取产品列表
			resp, err := client.GetProducts(&esim.ProductParams{
				Type:  pType,
				Page:  page,
				Limit: 20,
			})

			if err != nil {
				fmt.Printf("获取产品列表失败 (页 %d): %v\n", page, err)
				totalFailed++
				break
			}

			if !resp.Success || len(resp.Message.Products) == 0 {
				break
			}

			fmt.Printf("  页 %d/%d: %d 个产品\n", page, resp.Message.Pagination.TotalPages, len(resp.Message.Products))

			// 转换并保存产品
			for _, apiProduct := range resp.Message.Products {
				product, err := convertToModel(&apiProduct)
				if err != nil {
					fmt.Printf("    ✗ 转换产品失败 [%s]: %v\n", apiProduct.Name, err)
					totalFailed++
					continue
				}

				if err := productRepo.Upsert(ctx, product); err != nil {
					fmt.Printf("    ✗ 保存产品失败 [%s]: %v\n", apiProduct.Name, err)
					totalFailed++
					continue
				}

				fmt.Printf("    ✓ %s\n", apiProduct.Name)
				totalSynced++

				// 检查限制
				if limit > 0 && totalSynced >= limit {
					goto done
				}
			}

			// 检查是否还有更多页
			if page >= resp.Message.Pagination.TotalPages {
				break
			}

			page++
			time.Sleep(500 * time.Millisecond) // 避免请求过快
		}
	}

done:
	fmt.Printf("\n同步完成!\n")
	fmt.Printf("  成功: %d\n", totalSynced)
	fmt.Printf("  失败: %d\n", totalFailed)

	return nil
}

// listProducts 列出本地产品
func listProducts(ctx context.Context, db *data.Database, productType string) error {
	productRepo := db.GetProductRepository()

	params := repository.ListParams{
		Type:    productType,
		Status:  "active",
		OrderBy: "sort_order",
		Limit:   100,
	}

	products, total, err := productRepo.List(ctx, params)
	if err != nil {
		return fmt.Errorf("查询产品失败: %w", err)
	}

	fmt.Printf("本地产品列表 (共 %d 个)\n", total)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, product := range products {
		fmt.Printf("%d. [%s] %s\n", i+1, product.Type, product.Name)
		fmt.Printf("   ID: %d | 第三方ID: %s\n", product.ID, product.ThirdPartyID)
		fmt.Printf("   流量: %dMB | 有效期: %d天\n", product.DataSize, product.ValidDays)
		fmt.Printf("   价格: $%.2f | 成本: $%.2f | 利润: $%.2f\n",
			product.Price, product.CostPrice, product.Price-product.CostPrice)
		fmt.Printf("   状态: %s | 同步时间: %s\n",
			product.Status, product.SyncedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

// convertToModel 将 API 产品转换为数据库模型
func convertToModel(apiProduct *esim.Product) (*models.Product, error) {
	// 序列化国家列表
	countriesJSON, err := json.Marshal(apiProduct.Countries)
	if err != nil {
		return nil, fmt.Errorf("序列化国家列表失败: %w", err)
	}

	// 序列化特性列表
	featuresJSON, err := json.Marshal(apiProduct.Features)
	if err != nil {
		return nil, fmt.Errorf("序列化特性列表失败: %w", err)
	}

	// 使用第三方ID作为唯一标识
	thirdPartyID := apiProduct.ThirdPartyID
	if thirdPartyID == "" {
		thirdPartyID = fmt.Sprintf("product-%d", apiProduct.ID)
	}

	// 计算价格
	price := apiProduct.Price
	if price == 0 {
		price = apiProduct.RetailPrice
	}

	costPrice := apiProduct.CostPrice
	if costPrice == 0 {
		costPrice = apiProduct.AgentPrice
	}

	return &models.Product{
		ThirdPartyID:   thirdPartyID,
		Name:           apiProduct.Name,
		NameEn:         apiProduct.NameEn,
		Description:    apiProduct.Description,
		DescriptionEn:  apiProduct.DescriptionEn,
		Type:           string(apiProduct.Type),
		Countries:      string(countriesJSON),
		DataSize:       apiProduct.DataSize,
		ValidDays:      apiProduct.ValidDays,
		Features:       string(featuresJSON),
		Image:          apiProduct.Image,
		Price:          price,
		CostPrice:      costPrice,
		RetailPrice:    apiProduct.RetailPrice,
		AgentPrice:     apiProduct.AgentPrice,
		PlatformProfit: price - costPrice,
		IsHot:          apiProduct.IsHot,
		IsRecommend:    apiProduct.IsRecommend,
		SortOrder:      apiProduct.SortOrder,
		Status:         apiProduct.Status,
		SyncedAt:       time.Now(),
	}, nil
}

// syncProductDetails 同步产品详情
func syncProductDetails(ctx context.Context, cfg *config.Config, db *data.Database, limit int) error {
	fmt.Println("开始同步产品详情...")

	// 创建 eSIM 客户端
	client := esim.NewClient(esim.Config{
		APIKey:         cfg.EsimSDK.APIKey,
		APISecret:      cfg.EsimSDK.APISecret,
		BaseURL:        cfg.EsimSDK.BaseURL,
		TimezoneOffset: cfg.EsimSDK.TimezoneOffset,
	})

	productRepo := db.GetProductRepository()
	detailRepo := db.GetProductDetailRepository()

	// 获取所有产品
	params := repository.ListParams{
		Status:  "active",
		OrderBy: "id",
		Limit:   0, // 获取所有
	}

	products, total, err := productRepo.List(ctx, params)
	if err != nil {
		return fmt.Errorf("获取产品列表失败: %w", err)
	}

	fmt.Printf("找到 %d 个产品需要同步详情\n\n", total)

	totalSynced := 0
	totalFailed := 0

	for i, product := range products {
		// 检查限制
		if limit > 0 && totalSynced >= limit {
			break
		}

		fmt.Printf("[%d/%d] 同步产品: %s (ID: %d)\n", i+1, len(products), product.Name, product.ID)

		// 从 API 获取产品详情
		// 注意：这里需要使用第三方ID，假设存储在 ThirdPartyID 字段中
		thirdPartyID := extractThirdPartyID(product.ThirdPartyID)
		if thirdPartyID == 0 {
			fmt.Printf("  ✗ 跳过: 无效的第三方ID\n")
			totalFailed++
			continue
		}

		resp, err := client.GetProduct(thirdPartyID)
		if err != nil {
			fmt.Printf("  ✗ 获取详情失败: %v\n", err)
			totalFailed++
			continue
		}

		if !resp.Success {
			fmt.Printf("  ✗ API返回失败: %s\n", resp.Data)
			totalFailed++
			continue
		}

		// 转换为详情模型
		if resp.ProductDetail == nil {
			fmt.Printf("  ✗ 产品详情为空\n")
			totalFailed++
			continue
		}

		detail, err := convertToDetailModel(product.ID, resp.ProductDetail)
		if err != nil {
			fmt.Printf("  ✗ 转换失败: %v\n", err)
			totalFailed++
			continue
		}

		// 保存详情
		if err := detailRepo.Upsert(ctx, detail); err != nil {
			fmt.Printf("  ✗ 保存失败: %v\n", err)
			totalFailed++
			continue
		}

		fmt.Printf("  ✓ 同步成功\n")
		totalSynced++

		// 避免请求过快
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Printf("\n同步完成!\n")
	fmt.Printf("  成功: %d\n", totalSynced)
	fmt.Printf("  失败: %d\n", totalFailed)

	return nil
}

// extractThirdPartyID 从字符串中提取第三方ID
func extractThirdPartyID(thirdPartyID string) int {
	// 如果是 "product-123" 格式，提取数字
	if strings.HasPrefix(thirdPartyID, "product-") {
		idStr := strings.TrimPrefix(thirdPartyID, "product-")
		if id, err := strconv.Atoi(idStr); err == nil {
			return id
		}
	}
	// 尝试直接转换
	if id, err := strconv.Atoi(thirdPartyID); err == nil {
		return id
	}
	return 0
}

// convertToDetailModel 将 API 产品详情转换为详情模型
func convertToDetailModel(productID int, apiDetail *esim.ProductDetail) (*models.ProductDetail, error) {
	// 序列化国家列表
	countriesJSON, err := json.Marshal(apiDetail.Countries)
	if err != nil {
		return nil, fmt.Errorf("序列化国家列表失败: %w", err)
	}

	// 序列化特性列表
	featuresJSON, err := json.Marshal(apiDetail.Features)
	if err != nil {
		return nil, fmt.Errorf("序列化特性列表失败: %w", err)
	}

	// 计算数据大小字符串
	dataSize := "无限流量"
	if apiDetail.DataSize > 0 {
		if apiDetail.DataSize >= 1024 {
			dataSize = fmt.Sprintf("%.1fGB", float64(apiDetail.DataSize)/1024)
		} else {
			dataSize = fmt.Sprintf("%dMB", apiDetail.DataSize)
		}
	}

	return &models.ProductDetail{
		ProductID:    productID,
		ThirdPartyID: apiDetail.ID,
		Name:         apiDetail.Name,
		Type:         apiDetail.Type,
		Countries:    string(countriesJSON),
		DataSize:     dataSize,
		ValidDays:    apiDetail.ValidDays,
		Price:        apiDetail.Price,
		CostPrice:    apiDetail.CostPrice,
		Description:  apiDetail.Description,
		Features:     string(featuresJSON),
		Status:       apiDetail.Status,
		ApiCreatedAt: "", // API 响应中没有 createdAt 字段
		SyncedAt:     time.Now(),
	}, nil
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("eSIM 管理工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  gm -cmd <命令> [选项]")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  sync-products         从 API 同步产品数据到本地数据库")
	fmt.Println("  sync-product-details  从 API 同步产品详情到详情表")
	fmt.Println("  list-products         列出本地数据库中的产品")
	fmt.Println("  help                  显示帮助信息")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -config <path>   配置文件路径 (默认: config/config.json)")
	fmt.Println("  -type <type>     产品类型: local, regional, global")
	fmt.Println("  -limit <n>       限制同步数量 (0 表示全部)")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # 同步所有产品")
	fmt.Println("  gm -cmd sync-products")
	fmt.Println()
	fmt.Println("  # 同步产品详情")
	fmt.Println("  gm -cmd sync-product-details")
	fmt.Println()
	fmt.Println("  # 只同步前 10 个产品的详情")
	fmt.Println("  gm -cmd sync-product-details -limit 10")
	fmt.Println()
	fmt.Println("  # 列出本地产品")
	fmt.Println("  gm -cmd list-products")
	fmt.Println()
	fmt.Println("  # 列出指定类型的产品")
	fmt.Println("  gm -cmd list-products -type global")
}
