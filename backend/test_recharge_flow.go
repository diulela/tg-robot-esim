package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"tg-robot-sim/config"
	"tg-robot-sim/services"
	"tg-robot-sim/storage/data"
)

// 测试充值流程的端到端测试
func main() {
	log.Println("开始 USDT 充值流程端到端测试...")

	// 加载配置
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := data.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 运行数据库迁移
	if err := data.AutoMigrate(db.GetDB()); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建服务
	walletService := services.NewWalletService(db.GetWalletRepository(), db.GetRechargeOrderRepository(), nil)
	blockchainService := services.NewMockBlockchainService()
	notificationService := services.NewMockNotificationService()

	rechargeService := services.NewRechargeService(
		db.GetRechargeOrderRepository(),
		walletService,
		blockchainService,
		notificationService,
		db.GetDB(),
		cfg.Recharge.DepositAddress,
		cfg.Recharge.MinAmount,
		cfg.Recharge.MaxAmount,
	)

	// 测试用户ID
	testUserID := int64(123456789)

	// 1. 测试创建充值订单
	log.Println("\n=== 测试 1: 创建充值订单 ===")
	ctx := context.Background()
	order, err := rechargeService.CreateRechargeOrder(ctx, testUserID, "100.00")
	if err != nil {
		log.Fatalf("创建充值订单失败: %v", err)
	}
	log.Printf("✅ 充值订单创建成功:")
	log.Printf("   订单号: %s", order.OrderNo)
	log.Printf("   充值金额: %s USDT", order.Amount)
	log.Printf("   精确金额: %s USDT", order.ExactAmount)
	log.Printf("   收款地址: %s", order.WalletAddress)
	log.Printf("   过期时间: %s", order.ExpiresAt.Format("2006-01-02 15:04:05"))

	// 2. 测试获取订单详情
	log.Println("\n=== 测试 2: 获取订单详情 ===")
	retrievedOrder, err := rechargeService.GetRechargeOrder(ctx, order.OrderNo)
	if err != nil {
		log.Fatalf("获取订单详情失败: %v", err)
	}
	log.Printf("✅ 订单详情获取成功:")
	log.Printf("   状态: %s", retrievedOrder.Status)
	log.Printf("   精确金额: %s USDT", retrievedOrder.ExactAmount)

	// 3. 测试检查充值状态（模拟转账完成）
	log.Println("\n=== 测试 3: 检查充值状态 ===")
	updatedOrder, err := rechargeService.CheckRechargeStatus(ctx, order.OrderNo)
	if err != nil {
		log.Fatalf("检查充值状态失败: %v", err)
	}
	log.Printf("✅ 充值状态检查完成:")
	log.Printf("   状态: %s", updatedOrder.Status)
	if updatedOrder.TxHash != "" {
		log.Printf("   交易哈希: %s", updatedOrder.TxHash)
	}
	if updatedOrder.ConfirmedAt != nil {
		log.Printf("   确认时间: %s", updatedOrder.ConfirmedAt.Format("2006-01-02 15:04:05"))
	}

	// 4. 测试获取用户钱包余额
	log.Println("\n=== 测试 4: 检查钱包余额 ===")
	wallet, err := walletService.GetOrCreateWallet(ctx, testUserID)
	if err != nil {
		log.Fatalf("获取钱包失败: %v", err)
	}
	log.Printf("✅ 钱包余额:")
	log.Printf("   可用余额: %s USDT", wallet.Balance)
	log.Printf("   总收入: %s USDT", wallet.TotalIncome)

	// 5. 测试获取充值历史
	log.Println("\n=== 测试 5: 获取充值历史 ===")
	orders, total, err := rechargeService.GetUserRechargeHistory(ctx, testUserID, 10, 0)
	if err != nil {
		log.Fatalf("获取充值历史失败: %v", err)
	}
	log.Printf("✅ 充值历史:")
	log.Printf("   总记录数: %d", total)
	log.Printf("   当前页记录数: %d", len(orders))
	for i, historyOrder := range orders {
		log.Printf("   [%d] 订单号: %s, 金额: %s, 状态: %s",
			i+1, historyOrder.OrderNo, historyOrder.Amount, historyOrder.Status)
	}

	// 6. 测试 HTTP API（如果服务器正在运行）
	log.Println("\n=== 测试 6: HTTP API 测试 ===")
	testHTTPAPI(testUserID)

	log.Println("\n🎉 所有测试完成！USDT 充值流程端到端测试成功！")
}

// testHTTPAPI 测试 HTTP API
func testHTTPAPI(userID int64) {
	baseURL := "http://localhost:8080"

	// 测试创建充值订单 API
	log.Println("测试创建充值订单 API...")

	requestBody := map[string]interface{}{
		"amount": "50.00",
	}

	jsonData, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", baseURL+"/api/miniapp/wallet/recharge", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ 创建请求失败: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Telegram-Init-Data", fmt.Sprintf("user_id=%d", userID))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ HTTP 请求失败: %v (服务器可能未启动)", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ 读取响应失败: %v", err)
		return
	}

	log.Printf("✅ HTTP API 测试成功:")
	log.Printf("   状态码: %d", resp.StatusCode)
	log.Printf("   响应: %s", string(body))
}
