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
	"tg-robot-sim/server"
	"tg-robot-sim/services"
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

	// 创建服务
	// 注意：这里需要实现 BlockchainService，暂时传 nil
	walletService := services.NewWalletService(db.GetWalletRepository(), db.GetRechargeOrderRepository(), nil)
	productService := services.NewProductService(db.GetProductRepository())
	orderService := services.NewOrderService(db.GetOrderRepository(), db.GetProductRepository(), walletService)
	transactionService := services.NewTransactionService(db.GetTransactionRepository())

	// 创建 HTTP 服务器
	httpServer := server.NewMiniAppHTTPServer(cfg, productService, walletService, orderService, transactionService)

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
