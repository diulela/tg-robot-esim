# eSIM 产品功能集成示例

## 完整集成示例

以下是在 `cmd/bot/main.go` 中集成 eSIM 产品功能的完整示例：

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-robot-sim/config"
	"tg-robot-sim/handlers"
	"tg-robot-sim/pkg/logger"
	"tg-robot-sim/services"
	"tg-robot-sim/storage"
	"tg-robot-sim/storage/repository"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// 2. 初始化日志
	appLogger := logger.NewLogger(cfg.Logging.Level, cfg.Logging.File)
	appLogger.Info("Starting Telegram Bot...")

	// 3. 初始化数据库
	db, err := storage.InitDatabase(cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 4. 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// 5. 初始化服务层
	sessionService := services.NewSessionService(sessionRepo, appLogger)
	dialogService := services.NewDialogService(sessionService, appLogger)
	menuService := services.NewMenuService(sessionService, appLogger)
	notificationService := services.NewNotificationService(appLogger)
	
	// 初始化 eSIM 服务
	esimService := services.NewEsimService(
		cfg.EsimSDK.APIKey,
		cfg.EsimSDK.APISecret,
		cfg.EsimSDK.BaseURL,
	)

	// 6. 初始化 Telegram Bot
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		appLogger.Fatal("Failed to create bot: %v", err)
	}

	bot.Debug = cfg.Telegram.Debug
	appLogger.Info("Authorized on account %s", bot.Self.UserName)

	// 7. 创建处理器注册表
	registry := handlers.NewRegistry(bot, appLogger)

	// 8. 注册命令处理器
	startHandler := handlers.NewStartHandler(bot, userRepo, dialogService)
	registry.RegisterCommandHandler(startHandler)

	helpHandler := handlers.NewHelpHandler(bot)
	registry.RegisterCommandHandler(helpHandler)

	// 注册产品处理器（命令和回调）
	productsHandler := handlers.NewProductsHandler(bot, esimService, appLogger)
	registry.RegisterCommandHandler(productsHandler)
	registry.RegisterCallbackHandler(productsHandler)

	// 9. 注册回调处理器
	callbackHandler := handlers.NewCallbackQueryHandler(bot, menuService, appLogger)
	registry.RegisterCallbackHandler(callbackHandler)

	// 10. 启动消息处理循环
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置更新配置
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// 11. 处理系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	appLogger.Info("Bot is running. Press Ctrl+C to stop.")

	// 12. 主循环
	for {
		select {
		case <-sigChan:
			appLogger.Info("Shutting down gracefully...")
			cancel()
			return

		case update := <-updates:
			go handleUpdate(ctx, registry, update, appLogger)
		}
	}
}

// handleUpdate 处理更新
func handleUpdate(ctx context.Context, registry handlers.HandlerRegistry, update tgbotapi.Update, logger services.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic in update handler: %v", r)
		}
	}()

	// 处理消息
	if update.Message != nil {
		if err := registry.RouteMessage(ctx, update.Message); err != nil {
			logger.Error("Failed to route message: %v", err)
		}
		return
	}

	// 处理回调查询
	if update.CallbackQuery != nil {
		if err := registry.RouteCallback(ctx, update.CallbackQuery); err != nil {
			logger.Error("Failed to route callback: %v", err)
		}
		return
	}
}
```

## 关键集成点

### 1. 配置加载

```go
cfg, err := config.LoadConfig("config/config.json")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

配置文件应包含 eSIM SDK 配置：

```json
{
  "esim_sdk": {
    "api_key": "your-api-key",
    "api_secret": "your-api-secret",
    "base_url": "https://api.your-domain.com"
  }
}
```

### 2. 服务初始化

```go
// 初始化 eSIM 服务
esimService := services.NewEsimService(
    cfg.EsimSDK.APIKey,
    cfg.EsimSDK.APISecret,
    cfg.EsimSDK.BaseURL,
)
```

### 3. 处理器注册

```go
// 创建产品处理器
productsHandler := handlers.NewProductsHandler(bot, esimService, appLogger)

// 注册为命令处理器（处理 /products 命令）
registry.RegisterCommandHandler(productsHandler)

// 注册为回调处理器（处理按钮点击）
registry.RegisterCallbackHandler(productsHandler)
```

### 4. 更新主菜单

主菜单已自动更新，包含 "🛍️ 浏览产品" 按钮。

## 环境变量配置

创建 `.env` 文件：

```bash
# Telegram Bot
TELEGRAM_BOT_TOKEN=your-telegram-bot-token

# Database
DATABASE_URL=bot.db

# eSIM API
ESIM_API_KEY=your-esim-api-key
ESIM_API_SECRET=your-esim-api-secret
ESIM_BASE_URL=https://api.your-domain.com

# Logging
LOG_LEVEL=info
DEBUG=false
```

## 测试步骤

### 1. 配置测试

```bash
# 设置环境变量
export TELEGRAM_BOT_TOKEN="your-token"
export ESIM_API_KEY="your-api-key"
export ESIM_API_SECRET="your-api-secret"

# 运行程序
go run cmd/bot/main.go
```

### 2. 功能测试

在 Telegram 中测试以下功能：

1. **启动机器人**
   ```
   /start
   ```
   应该看到主菜单，包含 "🛍️ 浏览产品" 按钮

2. **浏览产品**
   - 点击 "🛍️ 浏览产品"
   - 选择产品类型（本地/区域/全球）
   - 查看产品列表
   - 点击查看产品详情

3. **搜索产品**
   ```
   /products CN
   ```
   应该显示中国相关的产品

4. **查看产品详情**
   - 在产品列表中点击 "查看详情"
   - 应该显示完整的产品信息

## 错误处理

### API 调用失败

```go
products, err := esimService.GetProducts(ctx, params)
if err != nil {
    logger.Error("Failed to get products: %v", err)
    // 向用户发送友好的错误消息
    return sendError(chatID, "获取产品列表失败，请稍后重试")
}
```

### 配置验证

```go
if cfg.EsimSDK.APIKey == "" || cfg.EsimSDK.APIKey == "${ESIM_API_KEY}" {
    return fmt.Errorf("eSIM API key is required")
}
```

## 性能优化

### 1. 使用上下文超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

products, err := esimService.GetProducts(ctx, params)
```

### 2. 缓存产品列表

```go
// 可以添加缓存层
type CachedEsimService struct {
    service EsimService
    cache   *cache.Cache
}

func (s *CachedEsimService) GetProducts(ctx context.Context, params *esim.ProductParams) (*esim.ProductListResponse, error) {
    // 检查缓存
    cacheKey := fmt.Sprintf("products:%s:%d", params.Type, params.Page)
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*esim.ProductListResponse), nil
    }

    // 调用实际服务
    products, err := s.service.GetProducts(ctx, params)
    if err != nil {
        return nil, err
    }

    // 缓存结果（5分钟）
    s.cache.Set(cacheKey, products, 5*time.Minute)
    return products, nil
}
```

### 3. 并发处理

```go
// 在 handleUpdate 中使用 goroutine
go handleUpdate(ctx, registry, update, appLogger)
```

## 监控和日志

### 添加关键日志

```go
logger.Info("User %d browsing products, type: %s, page: %d", 
    userID, productType, page)

logger.Info("User %d viewing product detail, productID: %d", 
    userID, productID)

logger.Info("User %d starting purchase, productID: %d", 
    userID, productID)
```

### 监控指标

可以添加以下监控指标：
- 产品浏览次数
- 产品详情查看次数
- 购买转化率
- API 调用成功率
- API 响应时间

## 下一步

1. 实现订单管理功能
2. 添加支付集成
3. 实现 eSIM 使用情况查询
4. 添加充值功能
5. 实现多语言支持
