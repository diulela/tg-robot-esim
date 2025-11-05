package server

import (
	"log"
	"net/http"
	"tg-robot-sim/api"
	"tg-robot-sim/config"
	"tg-robot-sim/pkg/logger"
	"tg-robot-sim/server/middleware"
	"tg-robot-sim/services"
)

func NewMiniAppHTTPServer(
	cfg *config.Config,
	productService services.ProductService,
	walletService services.WalletService,
	orderService services.OrderService,
	transactionService services.TransactionService,
	rechargeService services.RechargeService,
) *http.Server {
	mux := http.NewServeMux()

	// 创建 Mini App 处理器
	miniAppHandler := api.NewMiniAppApiService(
		productService,
		walletService,
		orderService,
		transactionService,
		rechargeService,
	)

	// 注册路由
	miniAppHandler.RegisterRoutes(mux)

	// 静态文件服务（Mini App 前端）
	fs := http.FileServer(http.Dir("miniapp/dist"))
	mux.Handle("/miniapp/", http.StripPrefix("/miniapp/", fs))

	// 初始化日志系统
	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	// 应用中间件
	var handler http.Handler = mux
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.TelegramWebAppMiddleware(cfg.Telegram.BotToken)(handler)

	return &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}
