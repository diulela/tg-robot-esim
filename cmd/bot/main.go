package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"tg-robot-sim/config"
	"tg-robot-sim/handlers"
	"tg-robot-sim/pkg/bot"
	"tg-robot-sim/pkg/logger"
	"tg-robot-sim/pkg/tron"
	"tg-robot-sim/services"
	"tg-robot-sim/storage/data"
)

const (
	defaultConfigPath = "config/config.json"
)

func main() {
	// 添加全局 panic 处理
	defer func() {
		if r := recover(); r != nil {
			// 获取详细的堆栈跟踪信息
			stack := make([]byte, 4096)
			length := runtime.Stack(stack, false)
			log.Printf("FATAL PANIC: %v\nStack trace:\n%s", r, string(stack[:length]))
			os.Exit(1)
		}
	}()

	// 加载配置
	cfg, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 初始化日志系统
	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	appLogger.Info("Bot starting with config:")
	appLogger.Info("- Database: %s", cfg.Database.Type)
	appLogger.Info("- Log Level: %s", cfg.Logging.Level)
	appLogger.Info("- Monitor Interval: %v", cfg.Blockchain.MonitorInterval)

	// 初始化数据库
	db, err := data.NewDatabase(&cfg.Database)
	if err != nil {
		appLogger.Error("Failed to initialize database: %v", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 自动迁移数据库表
	if err := db.AutoMigrate(); err != nil {
		appLogger.Error("Failed to migrate database: %v", err)
		log.Fatalf("Failed to migrate database: %v", err)
	}
	appLogger.Info("Database initialized and migrated successfully")

	// 初始化会话服务
	sessionTimeout := 30 * time.Minute
	sessionService := services.NewSessionService(db.GetSessionRepository(), sessionTimeout)
	appLogger.Info("Session service initialized")

	// 初始化 TRON 客户端
	tronClient := tron.NewClient(cfg.Blockchain.TronEndpoint, cfg.Blockchain.TronAPIKey, appLogger)
	appLogger.Info("TRON client initialized")

	// 初始化区块链服务
	blockchainService := services.NewBlockchainService(tronClient, db.GetTransactionRepository(), &cfg.Blockchain, appLogger)
	appLogger.Info("Blockchain service initialized")

	// 初始化菜单服务
	menuService := services.NewMenuService(sessionService, appLogger)
	appLogger.Info("Menu service initialized")

	// 初始化对话服务
	dialogService := services.NewDialogService(sessionService, db.GetUserRepository(), menuService, appLogger)
	appLogger.Info("Dialog service initialized")

	// 初始化 eSIM 服务
	var esimService services.EsimService
	if cfg.EsimSDK.APIKey != "" && cfg.EsimSDK.APIKey != "${ESIM_API_KEY}" {
		esimService = services.NewEsimService(
			cfg.EsimSDK.APIKey,
			cfg.EsimSDK.APISecret,
			cfg.EsimSDK.BaseURL,
			cfg.EsimSDK.TimezoneOffset,
		)
		appLogger.Info("eSIM service initialized with timezone offset: %d hours", cfg.EsimSDK.TimezoneOffset)
	} else {
		appLogger.Warn("eSIM service not configured, product features will be disabled")
	}

	// 初始化 Telegram Bot
	telegramBot, err := bot.NewBot(&cfg.Telegram, appLogger)
	if err != nil {
		appLogger.Error("Failed to initialize bot: %v", err)
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// 初始化通知服务
	// notificationService := services.NewNotificationService(telegramBot.GetAPI(), appLogger)
	// appLogger.Info("Notification service initialized")

	// 注册中间件
	registry := telegramBot.GetRegistry()

	// 注册日志中间件
	loggingMiddleware := handlers.NewLoggingMiddleware(appLogger)
	if err := registry.RegisterMiddleware(loggingMiddleware); err != nil {
		appLogger.Error("Failed to register logging middleware: %v", err)
		log.Fatalf("Failed to register logging middleware: %v", err)
	}

	// 注册限流中间件
	rateLimitMiddleware := handlers.NewRateLimitMiddleware(1 * time.Second)
	if err := registry.RegisterMiddleware(rateLimitMiddleware); err != nil {
		appLogger.Error("Failed to register rate limit middleware: %v", err)
		log.Fatalf("Failed to register rate limit middleware: %v", err)
	}

	// 注册命令处理器
	startHandler := handlers.NewStartHandler(telegramBot.GetAPI(), db.GetUserRepository(), dialogService)
	if err := registry.RegisterCommandHandler(startHandler); err != nil {
		appLogger.Error("Failed to register start handler: %v", err)
		log.Fatalf("Failed to register start handler: %v", err)
	}

	helpHandler := handlers.NewHelpHandler(telegramBot.GetAPI(), dialogService)
	if err := registry.RegisterCommandHandler(helpHandler); err != nil {
		appLogger.Error("Failed to register help handler: %v", err)
		log.Fatalf("Failed to register help handler: %v", err)
	}

	menuHandler := handlers.NewMenuHandler(telegramBot.GetAPI(), menuService)
	if err := registry.RegisterCommandHandler(menuHandler); err != nil {
		appLogger.Error("Failed to register menu handler: %v", err)
		log.Fatalf("Failed to register menu handler: %v", err)
	}

	// 注册产品处理器（如果 eSIM 服务已配置）
	if esimService != nil {
		productsHandler := handlers.NewProductsHandler(
			telegramBot.GetAPI(),
			esimService,
			db.GetProductRepository(),
			db.GetProductDetailRepository(),
			appLogger,
		)
		if err := registry.RegisterCommandHandler(productsHandler); err != nil {
			appLogger.Error("Failed to register products command handler: %v", err)
			log.Fatalf("Failed to register products command handler: %v", err)
		}
		if err := registry.RegisterCallbackHandler(productsHandler); err != nil {
			appLogger.Error("Failed to register products callback handler: %v", err)
			log.Fatalf("Failed to register products callback handler: %v", err)
		}

		// 注册 Inline 查询处理器
		inlineHandler := handlers.NewInlineHandler(
			telegramBot.GetAPI(),
			db.GetProductRepository(),
			appLogger,
		)
		if err := registry.RegisterInlineHandler(inlineHandler); err != nil {
			appLogger.Error("Failed to register inline handler: %v", err)
			log.Fatalf("Failed to register inline handler: %v", err)
		}

		appLogger.Info("Products and inline handlers registered successfully")
	}

	// 注册消息处理器
	messageHandler := handlers.NewGeneralMessageHandler(telegramBot.GetAPI(), dialogService)
	if err := registry.RegisterMessageHandler(messageHandler); err != nil {
		appLogger.Error("Failed to register message handler: %v", err)
		log.Fatalf("Failed to register message handler: %v", err)
	}

	// 注册回调处理器
	callbackHandler := handlers.NewCallbackQueryHandler(telegramBot.GetAPI(), menuService, appLogger)
	if err := registry.RegisterCallbackHandler(callbackHandler); err != nil {
		appLogger.Error("Failed to register callback handler: %v", err)
		log.Fatalf("Failed to register callback handler: %v", err)
	}

	appLogger.Info("All handlers and middleware registered successfully")

	// 启动区块链监控（如果配置了钱包地址）
	if cfg.Blockchain.WalletAddress != "" {
		if err := blockchainService.MonitorAddress(cfg.Blockchain.WalletAddress); err != nil {
			appLogger.Error("Failed to add wallet address to monitoring: %v", err)
		} else {
			appLogger.Info("Added wallet address to monitoring: %s", cfg.Blockchain.WalletAddress)
		}

		if err := blockchainService.StartMonitoring(ctx); err != nil {
			appLogger.Error("Failed to start blockchain monitoring: %v", err)
		} else {
			appLogger.Info("Blockchain monitoring started")
		}
	}

	// 等待信号
	go func() {
		<-sigChan
		appLogger.Info("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// 启动机器人
	if err := telegramBot.Start(ctx); err != nil {
		appLogger.Error("Failed to start bot: %v", err)
		log.Fatalf("Failed to start bot: %v", err)
	}

	// 等待关闭信号
	<-ctx.Done()
	appLogger.Info("Context cancelled, shutting down...")

	// 优雅关闭
	if cfg.Blockchain.WalletAddress != "" {
		if err := blockchainService.StopMonitoring(); err != nil {
			appLogger.Error("Failed to stop blockchain monitoring: %v", err)
		} else {
			appLogger.Info("Blockchain monitoring stopped")
		}
	}

	telegramBot.Stop()
	appLogger.Info("Bot shutdown complete")
}
