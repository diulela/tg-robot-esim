---
inclusion: fileMatch
fileMatchPattern: 'backend/**'
---

# 后端开发规范 (Go + Telegram Bot)

## Go 开发环境

### 版本要求
- **Go版本**: 1.24.2+
- **模块管理**: Go Modules (`go.mod`)
- **代码格式化**: `go fmt` (必须)
- **依赖管理**: `go mod tidy` (定期执行)

### 开发工具链
```bash
# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test ./...

# 构建应用
go build -o bin/bot ./cmd/bot/main.go
go build -o bin/miniapp ./cmd/miniapp/main.go

# 交叉编译 (Linux)
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o main-linux-amd64 ./cmd/bot/main.go
```

## 项目架构规范

### 目录结构约定
```
backend/
├── cmd/                    # 应用程序入口点
│   ├── bot/               # Telegram Bot 服务
│   └── miniapp/           # HTTP API 服务
├── config/                # 配置管理
│   ├── config.go         # 配置结构定义
│   └── config.json       # 配置文件
├── handlers/              # 请求处理器
│   ├── bot/              # Bot 消息处理器
│   ├── http/             # HTTP 请求处理器
│   └── middleware.go     # 中间件
├── services/              # 业务逻辑层
│   ├── interfaces.go     # 服务接口定义
│   ├── product.go        # 产品服务实现
│   ├── wallet.go         # 钱包服务实现
│   └── order.go          # 订单服务实现
├── storage/               # 数据存储层
│   ├── data/             # 数据库连接和迁移
│   ├── repository/       # 数据访问层
│   └── models/           # 数据模型
├── pkg/                   # 可复用包
│   ├── blockchain/       # 区块链相关
│   ├── telegram/         # Telegram 工具
│   └── utils/            # 通用工具
└── utils/                 # 项目特定工具
```

### 包命名规范
- 包名使用小写字母，避免下划线
- 包名应该简洁且具有描述性
- 避免使用 `common`、`util`、`base` 等通用名称
- 每个包应该有明确的职责

## 代码风格规范

### 命名约定
```go
// 常量：大写字母 + 下划线
const (
    DEFAULT_TIMEOUT = 30 * time.Second
    MAX_RETRY_COUNT = 3
)

// 变量和函数：驼峰命名
var userService UserService
func getUserByID(id int64) (*models.User, error) { }

// 结构体和接口：帕斯卡命名
type UserService interface { }
type ProductRepository struct { }

// 私有成员：小写开头
type userImpl struct {
    db *gorm.DB
    logger *log.Logger
}
```

### 函数和方法规范
```go
// 函数注释：使用中文，说明功能、参数和返回值
// CreateUser 创建新用户
// 参数：userData 用户数据
// 返回：创建的用户信息和可能的错误
func (s *userService) CreateUser(userData *CreateUserRequest) (*models.User, error) {
    // 参数验证
    if userData == nil {
        return nil, errors.New("用户数据不能为空")
    }
    
    // 业务逻辑
    user := &models.User{
        TelegramID: userData.TelegramID,
        Username:   userData.Username,
        CreatedAt:  time.Now(),
    }
    
    // 数据库操作
    if err := s.db.Create(user).Error; err != nil {
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }
    
    return user, nil
}
```

### 错误处理规范
```go
// 使用 errors.New 创建简单错误
func validateInput(input string) error {
    if input == "" {
        return errors.New("输入不能为空")
    }
    return nil
}

// 使用 fmt.Errorf 包装错误
func processData(data []byte) error {
    if err := validateData(data); err != nil {
        return fmt.Errorf("数据验证失败: %w", err)
    }
    return nil
}

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("字段 %s 验证失败: %s", e.Field, e.Message)
}
```

## 服务层设计规范

### 接口定义 (services/interfaces.go)
```go
package services

import (
    "context"
    "tg-robot-sim/storage/models"
)

// UserService 用户服务接口
type UserService interface {
    // CreateUser 创建用户
    CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error)
    
    // GetUserByTelegramID 根据 Telegram ID 获取用户
    GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
    
    // UpdateUser 更新用户信息
    UpdateUser(ctx context.Context, userID int64, req *UpdateUserRequest) error
}

// ProductService 产品服务接口
type ProductService interface {
    // GetProducts 获取产品列表
    GetProducts(ctx context.Context, filter *ProductFilter) ([]*models.Product, error)
    
    // GetProductByID 根据ID获取产品
    GetProductByID(ctx context.Context, id int64) (*models.Product, error)
}
```

### 服务实现规范
```go
package services

import (
    "context"
    "fmt"
    "tg-robot-sim/storage/repository"
    "tg-robot-sim/storage/models"
)

// userService 用户服务实现
type userService struct {
    userRepo repository.UserRepository
    logger   Logger
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, logger Logger) UserService {
    return &userService{
        userRepo: userRepo,
        logger:   logger,
    }
}

// CreateUser 实现用户创建逻辑
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    // 输入验证
    if err := s.validateCreateUserRequest(req); err != nil {
        return nil, fmt.Errorf("请求验证失败: %w", err)
    }
    
    // 检查用户是否已存在
    existingUser, err := s.userRepo.GetByTelegramID(ctx, req.TelegramID)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("检查用户存在性失败: %w", err)
    }
    
    if existingUser != nil {
        return nil, errors.New("用户已存在")
    }
    
    // 创建用户
    user := &models.User{
        TelegramID: req.TelegramID,
        Username:   req.Username,
        FirstName:  req.FirstName,
        LastName:   req.LastName,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        s.logger.Error("创建用户失败", "error", err, "telegram_id", req.TelegramID)
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }
    
    s.logger.Info("用户创建成功", "user_id", user.ID, "telegram_id", user.TelegramID)
    return user, nil
}
```

## 数据层设计规范

### 模型定义 (storage/models/)
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// User 用户模型
type User struct {
    ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    TelegramID int64          `gorm:"uniqueIndex;not null" json:"telegram_id"`
    Username   string         `gorm:"size:255" json:"username"`
    FirstName  string         `gorm:"size:255" json:"first_name"`
    LastName   string         `gorm:"size:255" json:"last_name"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 关联关系
    Wallet *Wallet `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
    Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

### 仓储接口定义 (storage/repository/interfaces.go)
```go
package repository

import (
    "context"
    "tg-robot-sim/storage/models"
)

// UserRepository 用户仓储接口
type UserRepository interface {
    // Create 创建用户
    Create(ctx context.Context, user *models.User) error
    
    // GetByID 根据ID获取用户
    GetByID(ctx context.Context, id int64) (*models.User, error)
    
    // GetByTelegramID 根据Telegram ID获取用户
    GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
    
    // Update 更新用户
    Update(ctx context.Context, user *models.User) error
    
    // Delete 删除用户
    Delete(ctx context.Context, id int64) error
}
```

### 仓储实现规范
```go
package repository

import (
    "context"
    "fmt"
    "gorm.io/gorm"
    "tg-robot-sim/storage/models"
)

// userRepository 用户仓储实现
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
        return fmt.Errorf("创建用户失败: %w", err)
    }
    return nil
}

// GetByTelegramID 根据Telegram ID获取用户
func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).
        Where("telegram_id = ?", telegramID).
        First(&user).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("用户不存在: telegram_id=%d", telegramID)
        }
        return nil, fmt.Errorf("查询用户失败: %w", err)
    }
    
    return &user, nil
}
```

## Telegram Bot 开发规范

### Bot 处理器结构
```go
package handlers

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "tg-robot-sim/services"
)

// BotHandler Bot消息处理器
type BotHandler struct {
    bot            *tgbotapi.BotAPI
    userService    services.UserService
    productService services.ProductService
    orderService   services.OrderService
}

// NewBotHandler 创建Bot处理器
func NewBotHandler(
    bot *tgbotapi.BotAPI,
    userService services.UserService,
    productService services.ProductService,
    orderService services.OrderService,
) *BotHandler {
    return &BotHandler{
        bot:            bot,
        userService:    userService,
        productService: productService,
        orderService:   orderService,
    }
}

// HandleUpdate 处理更新消息
func (h *BotHandler) HandleUpdate(update tgbotapi.Update) error {
    switch {
    case update.Message != nil:
        return h.handleMessage(update.Message)
    case update.CallbackQuery != nil:
        return h.handleCallbackQuery(update.CallbackQuery)
    default:
        return nil
    }
}
```

### 命令处理规范
```go
// handleMessage 处理文本消息
func (h *BotHandler) handleMessage(message *tgbotapi.Message) error {
    // 确保用户存在
    user, err := h.ensureUser(message.From)
    if err != nil {
        return fmt.Errorf("确保用户存在失败: %w", err)
    }
    
    // 根据命令分发处理
    switch message.Command() {
    case "start":
        return h.handleStartCommand(message, user)
    case "products":
        return h.handleProductsCommand(message, user)
    case "wallet":
        return h.handleWalletCommand(message, user)
    default:
        return h.handleUnknownCommand(message, user)
    }
}

// handleStartCommand 处理开始命令
func (h *BotHandler) handleStartCommand(message *tgbotapi.Message, user *models.User) error {
    welcomeText := fmt.Sprintf("欢迎使用 eSIM 商城，%s！\n\n请选择您需要的服务：", user.FirstName)
    
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("📱 浏览产品", "products"),
            tgbotapi.NewInlineKeyboardButtonData("💰 我的钱包", "wallet"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonWebApp("🛒 打开商城", tgbotapi.WebApp{
                URL: "https://your-domain.com/miniapp",
            }),
        ),
    )
    
    msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
    msg.ReplyMarkup = keyboard
    
    _, err := h.bot.Send(msg)
    return err
}
```

## HTTP API 开发规范

### 路由注册
```go
package handlers

import (
    "net/http"
    "tg-robot-sim/services"
)

// HTTPHandler HTTP请求处理器
type HTTPHandler struct {
    productService services.ProductService
    walletService  services.WalletService
    orderService   services.OrderService
}

// RegisterRoutes 注册路由
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
    // API 路由
    mux.HandleFunc("GET /api/products", h.handleGetProducts)
    mux.HandleFunc("GET /api/products/{id}", h.handleGetProduct)
    mux.HandleFunc("POST /api/orders", h.handleCreateOrder)
    mux.HandleFunc("GET /api/wallet", h.handleGetWallet)
    mux.HandleFunc("POST /api/wallet/recharge", h.handleRecharge)
    
    // 健康检查
    mux.HandleFunc("GET /health", h.handleHealth)
}
```

### API 响应格式
```go
// APIResponse 统一API响应格式
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Code    int         `json:"code"`
}

// SendJSONResponse 发送JSON响应
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := APIResponse{
        Success: statusCode < 400,
        Data:    data,
        Code:    statusCode,
    }
    
    json.NewEncoder(w).Encode(response)
}

// SendErrorResponse 发送错误响应
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := APIResponse{
        Success: false,
        Error:   message,
        Code:    statusCode,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

## 配置管理规范

### 配置结构定义
```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

// Config 应用配置
type Config struct {
    Server    ServerConfig    `json:"server"`
    Database  DatabaseConfig  `json:"database"`
    Telegram  TelegramConfig  `json:"telegram"`
    Blockchain BlockchainConfig `json:"blockchain"`
    Logging   LoggingConfig   `json:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Port         int           `json:"port"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
    IdleTimeout  time.Duration `json:"idle_timeout"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
    // 读取配置文件
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析配置文件失败: %w", err)
    }
    
    // 从环境变量覆盖配置
    config.overrideFromEnv()
    
    return &config, nil
}

// overrideFromEnv 从环境变量覆盖配置
func (c *Config) overrideFromEnv() {
    if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
        c.Telegram.BotToken = token
    }
    
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        c.Database.URL = dbURL
    }
}
```

## 日志记录规范

### 结构化日志
```go
package utils

import (
    "log/slog"
    "os"
)

// Logger 日志接口
type Logger interface {
    Info(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
}

// NewLogger 创建日志实例
func NewLogger(level string) Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    opts := &slog.HandlerOptions{
        Level: logLevel,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// 使用示例
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    s.logger.Info("开始创建用户", 
        "telegram_id", req.TelegramID,
        "username", req.Username,
    )
    
    // ... 业务逻辑
    
    if err != nil {
        s.logger.Error("创建用户失败",
            "error", err,
            "telegram_id", req.TelegramID,
        )
        return nil, err
    }
    
    s.logger.Info("用户创建成功",
        "user_id", user.ID,
        "telegram_id", user.TelegramID,
    )
    
    return user, nil
}
```

## 测试规范

### 单元测试
```go
package services_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "tg-robot-sim/services"
    "tg-robot-sim/storage/models"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

// TestUserService_CreateUser 测试用户创建
func TestUserService_CreateUser(t *testing.T) {
    // 准备测试数据
    mockRepo := new(MockUserRepository)
    mockLogger := new(MockLogger)
    userService := services.NewUserService(mockRepo, mockLogger)
    
    req := &services.CreateUserRequest{
        TelegramID: 123456789,
        Username:   "testuser",
        FirstName:  "Test",
    }
    
    // 设置模拟期望
    mockRepo.On("GetByTelegramID", mock.Anything, req.TelegramID).
        Return(nil, gorm.ErrRecordNotFound)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
        Return(nil)
    
    // 执行测试
    user, err := userService.CreateUser(context.Background(), req)
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.TelegramID, user.TelegramID)
    assert.Equal(t, req.Username, user.Username)
    
    // 验证模拟调用
    mockRepo.AssertExpectations(t)
}
```

## 部署和运维

### Docker 配置
```dockerfile
# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/bot/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

CMD ["./main"]
```

### 环境变量管理
```bash
# .env 文件示例
TELEGRAM_BOT_TOKEN=your_bot_token_here
DATABASE_URL=sqlite://./data.db
TRON_API_KEY=your_tron_api_key
LOG_LEVEL=info
DEBUG=false
```