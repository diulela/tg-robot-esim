package main

import (
	"context"
	"fmt"
	"log"

	"tg-robot-sim/config"
	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/services"
	"tg-robot-sim/storage/data"
)

// 这是一个示例程序，展示如何使用 eSIM 订单处理系统

func main() {
	fmt.Println("eSIM 订单处理系统示例")

	// 1. 初始化配置
	dbConfig := &config.DatabaseConfig{
		Type: "sqlite",
		DSN:  "test.db",
	}

	// 2. 初始化数据库
	database, err := data.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.Close()

	// 3. 自动迁移数据库
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 4. 初始化服务
	walletHistoryService := services.NewWalletHistoryService(
		database.GetWalletHistoryRepository(),
	)

	walletService := services.NewWalletService(
		database.GetWalletRepository(),
		database.GetRechargeOrderRepository(),
		nil, // blockchain service
		walletHistoryService,
	)

	orderService := services.NewOrderService(
		database.GetOrderRepository(),
		database.GetOrderDetailRepository(),
		database.GetProductRepository(),
		walletService,
	)

	// 5. 创建 eSIM 客户端
	esimClient := &esim.Client{
		// 配置 eSIM 客户端
	}

	orderSyncService := services.NewOrderSyncService(
		database.GetOrderRepository(),
		orderService,
		esimClient,
	)

	// 6. 示例：创建 eSIM 订单
	ctx := context.Background()

	// 创建订单请求
	createOrderReq := &services.CreateEsimOrderRequest{
		UserID:      123456,
		ProductID:   1,
		Quantity:    1,
		TotalAmount: "27.7200",
		Remark:      "测试订单",
	}

	fmt.Println("创建 eSIM 订单...")
	orderResp, err := orderService.CreateEsimOrder(ctx, createOrderReq)
	if err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}

	fmt.Printf("订单创建成功: ID=%d, 订单号=%s, 状态=%s\n",
		orderResp.OrderID, orderResp.OrderNo, orderResp.Status)

	// 7. 示例：查询订单详情
	fmt.Println("查询订单详情...")
	orderDetail, err := orderService.GetOrderWithDetail(ctx, orderResp.OrderID, 123456)
	if err != nil {
		log.Printf("查询订单详情失败: %v", err)
		return
	}

	fmt.Printf("订单详情: 产品=%s, 数量=%d, 金额=%s\n",
		orderDetail.ProductName, orderDetail.Quantity, orderDetail.TotalAmount)

	// 8. 示例：手动同步订单状态
	fmt.Println("同步订单状态...")
	syncResult, err := orderSyncService.SyncOrderStatus(ctx, orderResp.OrderID)
	if err != nil {
		log.Printf("同步订单状态失败: %v", err)
		return
	}

	fmt.Printf("同步结果: 成功=%t, 消息=%s\n", syncResult.Success, syncResult.Message)

	fmt.Println("示例程序执行完成")
}
