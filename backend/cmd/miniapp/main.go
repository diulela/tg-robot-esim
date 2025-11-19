package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tg-robot-sim/config"
	"tg-robot-sim/pkg/bot"
	"tg-robot-sim/pkg/logger"
	"tg-robot-sim/pkg/sdk/esim"
	"tg-robot-sim/pkg/tron"
	"tg-robot-sim/server"
	"tg-robot-sim/services"
	service_common "tg-robot-sim/services/common"
	"tg-robot-sim/storage/data"
)

const (
	defaultConfigPath = "config/config.json"
)

func main() {

	// 加载配置
	cfg, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// 初始化数据库
	db, err := data.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 运行数据库迁移
	log.Println("Running database migrations...")
	if err := data.AutoMigrate(db.GetDB()); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 为现有用户创建钱包
	log.Println("Migrating wallets for existing users...")
	if err := data.MigrateWallets(db.GetDB()); err != nil {
		log.Printf("Warning: Failed to migrate wallets: %v", err)
	}

	// 创建模拟服务用于测试
	// 初始化 TRON 客户端
	tronClient := tron.NewClient(cfg.Blockchain.TronEndpoint, cfg.Blockchain.TronAPIKey, appLogger)

	// 初始化区块链服务
	blockchainService := services.NewBlockchainService(tronClient, &cfg.Blockchain, appLogger)

	walletHistoryService := services.NewWalletHistoryService(db.GetWalletHistoryRepository())

	// 创建服务
	// 注意：这里需要实现 BlockchainService，暂时传 nil
	walletService := services.NewWalletService(
		db.GetWalletRepository(),
		db.GetRechargeOrderRepository(),
		blockchainService,
		walletHistoryService,
	)
	// 初始化 eSIM 服务
	var esimService service_common.EsimClientService
	if cfg.EsimSDK.APIKey != "" && cfg.EsimSDK.APIKey != "${ESIM_API_KEY}" {
		esimService = service_common.NewEsimClientService(
			cfg.EsimSDK.APIKey,
			cfg.EsimSDK.APISecret,
			cfg.EsimSDK.BaseURL,
			cfg.EsimSDK.TimezoneOffset,
		)
		appLogger.Info("eSIM service initialized")
	} else {
		appLogger.Warn("eSIM service not configured, orders will be created without provider integration")
	}

	productService := services.NewProductService(db.GetProductRepository())
	orderService := services.NewOrderService(
		db.GetOrderRepository(),
		db.GetOrderDetailRepository(),
		db.GetProductRepository(),
		walletService,
		esimService,
	)

	// 初始化订单同步服务
	var orderSyncService services.OrderSyncService
	if cfg.EsimSDK.APIKey != "" && cfg.EsimSDK.APIKey != "${ESIM_API_KEY}" {
		// 创建 eSIM Client 用于订单同步
		esimClient := esim.NewClient(esim.Config{
			APIKey:         cfg.EsimSDK.APIKey,
			APISecret:      cfg.EsimSDK.APISecret,
			BaseURL:        cfg.EsimSDK.BaseURL,
			TimezoneOffset: cfg.EsimSDK.TimezoneOffset,
		})

		// 创建订单同步服务
		orderSyncService = services.NewOrderSyncService(
			db.GetOrderRepository(),
			orderService,
			esimClient,
		)
		appLogger.Info("OrderSyncService initialized successfully")
	} else {
		appLogger.Warn("eSIM SDK not configured, OrderSyncService will not be initialized")
	}

	// 初始化 Telegram Bot
	telegramBot, err := bot.NewBot(&cfg.Telegram, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize bot: %v", err)
		log.Fatalf("Failed to initialize bot: %v", err)
	}
	notificationService := services.NewNotificationService(telegramBot.GetAPI(), appLogger)

	// 创建充值服务
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

	// 创建 HTTP 服务器
	httpServer := server.NewMiniAppHTTPServer(cfg, productService, walletService, orderService, walletHistoryService, rechargeService)

	// 启动区块链监控定时任务
	go func() {
		log.Println("Starting blockchain monitoring task...")
		startBlockchainMonitoring(rechargeService)
	}()

	// 启动充值订单过期定时任务
	go func() {
		log.Println("Starting recharge order expiration monitoring task...")
		startExpireOldOrdersMonitoring(rechargeService)
	}()

	// 启动订单同步定时任务
	if orderSyncService != nil {
		go func() {
			log.Println("Starting order sync task...")
			startOrderSyncTask(orderSyncService, appLogger)
		}()
	}

	// 启动服务器
	go func() {
		log.Printf("Starting Mini App HTTP server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// startBlockchainMonitoring 启动区块链监控定时任务
func startBlockchainMonitoring(rechargeService services.RechargeService) {
	// 每30秒执行一次监控任务
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Blockchain monitoring started, checking every 30 seconds")

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

			// 处理待支付的充值订单
			if err := rechargeService.ProcessPendingRecharges(ctx); err != nil {
				log.Printf("Error processing pending recharges: %v", err)
			}

			cancel()
		}
	}
}

// startExpireOldOrdersMonitoring 启动充值订单过期监控定时任务
func startExpireOldOrdersMonitoring(rechargeService services.RechargeService) {
	// 每30秒执行一次监控任务
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Recharge order expiration monitoring started, checking every 30 seconds")

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			// 清理过期订单
			if err := rechargeService.ExpireOldOrders(ctx); err != nil {
				log.Printf("Error expiring old orders: %v", err)
			}
			cancel()
		}
	}
}

// startOrderSyncTask 启动订单同步定时任务
func startOrderSyncTask(orderSyncService services.OrderSyncService, appLogger *logger.Logger) {
	// 每10秒执行一次同步任务
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Order sync task started, checking every 10 seconds")

	for {
		select {
		case <-ticker.C:
			// 创建带超时的上下文
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

			// 处理待同步订单
			if err := orderSyncService.ProcessPendingOrders(ctx); err != nil {
				appLogger.Error("Error processing pending orders: %v", err)
			}

			cancel()
		}
	}
}
