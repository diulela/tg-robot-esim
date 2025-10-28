# Telegram eSIM 电商系统 - 后端服务

基于 Go 语言开发的 Telegram Bot 和 HTTP API 服务，为 eSIM 电商平台提供完整的后端支持。

## 🚀 功能特性

- 🤖 **Telegram Bot 服务**: 智能对话处理、交互式菜单、命令路由
- 🌐 **HTTP API 服务**: RESTful API，支持 Mini App 前端
- 💰 **区块链支付**: TRON 网络 USDT-TRC20 支付集成
- 📱 **eSIM 管理**: 产品管理、订单处理、激活码生成
- 🔐 **用户系统**: 用户注册、认证、钱包管理
- 📊 **数据持久化**: SQLite/MySQL 数据库支持
- 🔔 **实时通知**: 交易状态、订单更新通知

## 🛠 技术栈

- **语言**: Go 1.24.2+
- **Web框架**: 标准库 `net/http`
- **ORM**: GORM v1.30.0
- **Telegram**: telegram-bot-api v5.5.1
- **数据库**: SQLite (开发) / MySQL (生产)
- **区块链**: TRON 网络 API

## 📁 项目结构

```
backend/
├── cmd/                    # 应用程序入口点
│   ├── bot/               # Telegram Bot 主程序
│   ├── miniapp/           # Mini App HTTP 服务器
│   └── gm/                # 管理工具
├── config/                # 配置管理
│   ├── config.go         # 配置结构定义
│   ├── config.json       # 配置文件
│   └── config.example.json # 配置模板
├── handlers/              # 请求处理器
│   ├── bot/              # Bot 消息处理器
│   ├── miniapp/          # HTTP API 处理器
│   ├── middleware/       # 中间件
│   ├── interfaces.go     # 处理器接口定义
│   └── registry.go       # 处理器注册
├── services/              # 业务逻辑层
│   ├── interfaces.go     # 服务接口定义
│   ├── dialog_service.go # 对话服务
│   ├── menu_service.go   # 菜单服务
│   ├── product_service.go # 产品服务
│   ├── order_service.go  # 订单服务
│   ├── wallet_service.go # 钱包服务
│   └── blockchain_service.go # 区块链服务
├── storage/               # 数据存储层
│   ├── data/             # 数据库连接和迁移
│   ├── models/           # 数据模型
│   └── repository/       # 数据访问层
├── pkg/                   # 可复用包和工具
│   ├── bot/              # Bot 相关工具
│   ├── tron/             # TRON 区块链集成
│   ├── logger/           # 日志工具
│   ├── retry/            # 重试机制
│   └── sdk/              # 第三方 SDK
├── utils/                 # 通用工具函数
└── scripts/               # 构建和部署脚本
```## ⚡ 快
速开始

### 环境要求

- **Go**: 1.24.2 或更高版本
- **数据库**: SQLite 3.x 或 MySQL 8.0+
- **Telegram Bot Token**: 从 @BotFather 获取
- **TRON API Key**: 用于区块链交易监控

### 安装和配置

1. **克隆项目并安装依赖**
```bash
cd backend
go mod download
```

2. **配置环境**
```bash
# 复制配置模板
cp config/config.example.json config/config.json

# 设置环境变量
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export TRON_API_KEY="your_tron_api_key_here"
export DATABASE_URL="sqlite://./data/bot.db"
```

3. **初始化数据库**
```bash
# SQLite (开发环境)
mkdir -p data
touch data/bot.db

# MySQL (生产环境)
mysql -u root -p -e "CREATE DATABASE telegram_bot CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 运行服务

#### 开发环境
```bash
# 启动 Telegram Bot 服务
go run cmd/bot/main.go

# 启动 Mini App HTTP 服务
go run cmd/miniapp/main.go

# 启动管理工具
go run cmd/gm/main.go
```

#### 生产环境
```bash
# 构建所有服务
./scripts/build.sh

# 运行 Bot 服务
./bin/bot

# 运行 Mini App 服务
./bin/miniapp
```

#### 使用 Docker
```bash
# 构建并启动所有服务
cd ../docker
docker-compose up -d

# 查看服务状态
docker-compose ps
```#
# ⚙️ 配置说明

### 环境变量

| 变量名 | 描述 | 必需 | 默认值 |
|--------|------|------|--------|
| `TELEGRAM_BOT_TOKEN` | Telegram Bot 令牌 | ✅ | - |
| `DATABASE_URL` | 数据库连接字符串 | ❌ | `sqlite://./data/bot.db` |
| `TRON_API_KEY` | TRON 网络 API 密钥 | ✅ | - |
| `ESIM_API_KEY` | eSIM 服务 API 密钥 | ❌ | - |
| `ESIM_API_SECRET` | eSIM 服务 API 密钥 | ❌ | - |
| `LOG_LEVEL` | 日志级别 | ❌ | `info` |
| `DEBUG` | 调试模式 | ❌ | `false` |

### 配置文件结构

```json
{
  "telegram": {
    "bot_token": "${TELEGRAM_BOT_TOKEN}",
    "webhook_url": "",
    "timeout": "60s",
    "debug": false
  },
  "database": {
    "type": "sqlite",
    "dsn": "bot.db",
    "max_connections": 10,
    "max_idle_conns": 5,
    "conn_max_life": "1h"
  },
  "blockchain": {
    "tron_api_key": "${TRON_API_KEY}",
    "tron_endpoint": "https://api.trongrid.io",
    "monitor_interval": "30s",
    "required_confirmations": 12,
    "wallet_address": ""
  },
  "server": {
    "port": 8080,
    "read_timeout": "30s",
    "write_timeout": "30s",
    "idle_timeout": "120s"
  },
  "esim_sdk": {
    "api_key": "${ESIM_API_KEY}",
    "api_secret": "${ESIM_API_SECRET}",
    "base_url": "https://api.your-domain.com",
    "timezone_offset": 0
  },
  "logging": {
    "level": "info",
    "file": "bot.log",
    "max_size": 100,
    "max_age": 30,
    "compress": true
  }
}
```

### 数据库配置

#### SQLite (开发环境)
```json
{
  "database": {
    "type": "sqlite",
    "dsn": "./data/bot.db"
  }
}
```

#### MySQL (生产环境)
```json
{
  "database": {
    "type": "mysql",
    "dsn": "user:password@tcp(localhost:3306)/telegram_bot?charset=utf8mb4&parseTime=True&loc=Local",
    "max_connections": 20,
    "max_idle_conns": 10,
    "conn_max_life": "1h"
  }
}
```##
 🏗 架构设计

### 分层架构

```
┌─────────────────────────────────────────┐
│              表示层 (Handlers)            │
│  ┌─────────────┐  ┌─────────────────────┐ │
│  │ Bot Handler │  │ HTTP API Handler    │ │
│  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│              业务层 (Services)            │
│  ┌─────────┐ ┌─────────┐ ┌─────────────┐ │
│  │ Dialog  │ │ Product │ │ Blockchain  │ │
│  │ Service │ │ Service │ │   Service   │ │
│  └─────────┘ └─────────┘ └─────────────┘ │
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│            数据访问层 (Repository)         │
│  ┌─────────┐ ┌─────────┐ ┌─────────────┐ │
│  │  User   │ │ Product │ │ Transaction │ │
│  │  Repo   │ │  Repo   │ │    Repo     │ │
│  └─────────┘ └─────────┘ └─────────────┘ │
└─────────────────────────────────────────┘
                    │
┌─────────────────────────────────────────┐
│              存储层 (Storage)             │
│        ┌─────────┐  ┌─────────────┐      │
│        │ SQLite  │  │   MySQL     │      │
│        └─────────┘  └─────────────┘      │
└─────────────────────────────────────────┘
```

### 核心组件

#### 1. 处理器层 (Handlers)
- **MessageHandler**: 处理 Telegram 文本消息
- **CallbackHandler**: 处理内联键盘回调
- **CommandHandler**: 处理 Bot 命令
- **HTTPHandler**: 处理 HTTP API 请求

#### 2. 服务层 (Services)
- **DialogService**: 对话管理和会话状态
- **MenuService**: 交互式菜单系统
- **ProductService**: eSIM 产品管理
- **OrderService**: 订单处理和管理
- **WalletService**: 用户钱包和余额
- **BlockchainService**: 区块链交易监控
- **NotificationService**: 消息通知服务

#### 3. 数据层 (Repository)
- **UserRepository**: 用户数据访问
- **ProductRepository**: 产品数据访问
- **OrderRepository**: 订单数据访问
- **TransactionRepository**: 交易数据访问

### 依赖注入

所有服务通过接口定义，支持依赖注入模式：

```go
// 服务接口定义
type ProductService interface {
    GetProducts(ctx context.Context) ([]*models.Product, error)
    GetProductByID(ctx context.Context, id int64) (*models.Product, error)
}

// 服务实现
type productService struct {
    productRepo repository.ProductRepository
    logger      Logger
}

// 依赖注入
func NewProductService(repo repository.ProductRepository, logger Logger) ProductService {
    return &productService{
        productRepo: repo,
        logger:      logger,
    }
}
```## 🔧 开发指
南

### 添加新功能

#### 1. 定义服务接口
在 `services/interfaces.go` 中定义新的服务接口：

```go
// NewFeatureService 新功能服务接口
type NewFeatureService interface {
    ProcessFeature(ctx context.Context, data *FeatureData) (*FeatureResult, error)
    ValidateFeature(data *FeatureData) error
}
```

#### 2. 实现服务
创建 `services/new_feature_service.go`：

```go
package services

import (
    "context"
    "tg-robot-sim/storage/repository"
)

type newFeatureService struct {
    repo   repository.FeatureRepository
    logger Logger
}

func NewNewFeatureService(repo repository.FeatureRepository, logger Logger) NewFeatureService {
    return &newFeatureService{
        repo:   repo,
        logger: logger,
    }
}

func (s *newFeatureService) ProcessFeature(ctx context.Context, data *FeatureData) (*FeatureResult, error) {
    // 实现业务逻辑
    s.logger.Info("处理新功能", "data", data)
    
    // 数据验证
    if err := s.ValidateFeature(data); err != nil {
        return nil, fmt.Errorf("数据验证失败: %w", err)
    }
    
    // 处理逻辑
    result, err := s.repo.ProcessData(ctx, data)
    if err != nil {
        return nil, fmt.Errorf("处理失败: %w", err)
    }
    
    return result, nil
}
```

#### 3. 创建处理器
创建 `handlers/bot/new_feature_handler.go`：

```go
package bot

import (
    "context"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "tg-robot-sim/services"
)

type NewFeatureHandler struct {
    bot            *tgbotapi.BotAPI
    featureService services.NewFeatureService
}

func NewNewFeatureHandler(bot *tgbotapi.BotAPI, service services.NewFeatureService) *NewFeatureHandler {
    return &NewFeatureHandler{
        bot:            bot,
        featureService: service,
    }
}

func (h *NewFeatureHandler) HandleCommand(ctx context.Context, message *tgbotapi.Message) error {
    // 处理命令逻辑
    result, err := h.featureService.ProcessFeature(ctx, &services.FeatureData{
        UserID: message.From.ID,
        Data:   message.Text,
    })
    
    if err != nil {
        return fmt.Errorf("处理功能失败: %w", err)
    }
    
    // 发送响应
    msg := tgbotapi.NewMessage(message.Chat.ID, result.Message)
    _, err = h.bot.Send(msg)
    return err
}

func (h *NewFeatureHandler) GetCommand() string {
    return "newfeature"
}

func (h *NewFeatureHandler) GetDescription() string {
    return "新功能命令"
}
```

#### 4. 注册处理器
在 `handlers/registry.go` 中注册：

```go
func (r *handlerRegistry) RegisterHandlers() error {
    // 注册新功能处理器
    newFeatureHandler := bot.NewNewFeatureHandler(r.bot, r.services.NewFeature)
    if err := r.RegisterCommandHandler(newFeatureHandler); err != nil {
        return fmt.Errorf("注册新功能处理器失败: %w", err)
    }
    
    return nil
}
```

### 数据模型定义

在 `storage/models/` 中定义数据模型：

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Feature 功能模型
type Feature struct {
    ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID    int64          `gorm:"not null;index" json:"user_id"`
    Name      string         `gorm:"size:255;not null" json:"name"`
    Data      string         `gorm:"type:text" json:"data"`
    Status    string         `gorm:"size:50;default:'active'" json:"status"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Feature) TableName() string {
    return "features"
}
```

### 数据库迁移

在 `storage/data/migrations.go` 中添加迁移：

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Product{},
        &models.Order{},
        &models.Feature{}, // 新增模型
    )
}
```## 🧪 测试


### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./services/...
go test ./handlers/...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 测试示例

创建 `services/product_service_test.go`：

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

// MockProductRepository 模拟产品仓储
type MockProductRepository struct {
    mock.Mock
}

func (m *MockProductRepository) GetAll(ctx context.Context) ([]*models.Product, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*models.Product), args.Error(1)
}

// TestProductService_GetProducts 测试获取产品列表
func TestProductService_GetProducts(t *testing.T) {
    // 准备测试数据
    mockRepo := new(MockProductRepository)
    mockLogger := new(MockLogger)
    productService := services.NewProductService(mockRepo, mockLogger)
    
    expectedProducts := []*models.Product{
        {ID: 1, Name: "测试产品1", Price: "10.00"},
        {ID: 2, Name: "测试产品2", Price: "20.00"},
    }
    
    // 设置模拟期望
    mockRepo.On("GetAll", mock.Anything).Return(expectedProducts, nil)
    
    // 执行测试
    products, err := productService.GetProducts(context.Background())
    
    // 验证结果
    assert.NoError(t, err)
    assert.Len(t, products, 2)
    assert.Equal(t, "测试产品1", products[0].Name)
    
    // 验证模拟调用
    mockRepo.AssertExpectations(t)
}
```

### 集成测试

创建 `tests/integration_test.go`：

```go
package tests

import (
    "context"
    "testing"
    "tg-robot-sim/config"
    "tg-robot-sim/storage/data"
)

func TestDatabaseConnection(t *testing.T) {
    // 加载测试配置
    cfg, err := config.LoadConfig("../config/config.test.json")
    assert.NoError(t, err)
    
    // 连接数据库
    db, err := data.NewDatabase(cfg.Database)
    assert.NoError(t, err)
    defer db.Close()
    
    // 测试数据库操作
    var count int64
    err = db.Model(&models.User{}).Count(&count).Error
    assert.NoError(t, err)
}
```

## 📊 监控和日志

### 结构化日志

使用结构化日志记录关键信息：

```go
import "log/slog"

// 创建日志实例
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// 记录日志
logger.Info("用户创建成功",
    "user_id", user.ID,
    "telegram_id", user.TelegramID,
    "username", user.Username,
)

logger.Error("数据库连接失败",
    "error", err,
    "database_type", cfg.Database.Type,
    "retry_count", retryCount,
)
```

### 性能监控

```go
// 监控函数执行时间
func (s *productService) GetProducts(ctx context.Context) ([]*models.Product, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        s.logger.Info("获取产品列表完成",
            "duration", duration,
            "duration_ms", duration.Milliseconds(),
        )
    }()
    
    // 业务逻辑
    return s.productRepo.GetAll(ctx)
}
```

### 健康检查

实现健康检查端点：

```go
// handlers/health.go
func (h *HTTPHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Unix(),
        "version":   "1.0.0",
        "services": map[string]string{
            "database":   "healthy",
            "blockchain": "healthy",
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
```## 🚀 部署


### 本地部署

```bash
# 构建所有服务
./scripts/build.sh

# 启动服务
./bin/bot &
./bin/miniapp &

# 检查服务状态
ps aux | grep -E "(bot|miniapp)"
```

### Docker 部署

```bash
# 构建镜像
cd ../docker
docker build -f Dockerfile -t tg-robot-sim:latest ..

# 启动服务栈
docker-compose up -d

# 查看日志
docker-compose logs -f telegram-bot
```

### 生产环境部署

1. **环境准备**
```bash
# 创建部署目录
mkdir -p /opt/tg-robot-sim/{bin,config,data,logs}

# 设置权限
useradd -r -s /bin/false tg-robot
chown -R tg-robot:tg-robot /opt/tg-robot-sim
```

2. **配置 systemd 服务**

创建 `/etc/systemd/system/tg-robot-bot.service`：
```ini
[Unit]
Description=Telegram Robot Bot Service
After=network.target

[Service]
Type=simple
User=tg-robot
Group=tg-robot
WorkingDirectory=/opt/tg-robot-sim
ExecStart=/opt/tg-robot-sim/bin/bot
Restart=always
RestartSec=5
Environment=TELEGRAM_BOT_TOKEN=your_token_here
Environment=DATABASE_URL=mysql://user:pass@localhost/db

[Install]
WantedBy=multi-user.target
```

3. **启动服务**
```bash
# 重载 systemd 配置
systemctl daemon-reload

# 启动并启用服务
systemctl enable tg-robot-bot
systemctl start tg-robot-bot

# 检查状态
systemctl status tg-robot-bot
```

### 备份和恢复

#### 数据库备份
```bash
# MySQL 备份
mysqldump -u root -p telegram_bot > backup_$(date +%Y%m%d_%H%M%S).sql

# SQLite 备份
cp data/bot.db backup/bot_$(date +%Y%m%d_%H%M%S).db
```

#### 配置备份
```bash
# 备份配置文件
tar -czf config_backup_$(date +%Y%m%d_%H%M%S).tar.gz config/
```

## 🔒 安全

### 环境变量管理

使用 `.env` 文件管理敏感信息：

```bash
# .env
TELEGRAM_BOT_TOKEN=1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ
TRON_API_KEY=your-tron-api-key
DATABASE_URL=mysql://user:password@localhost:3306/telegram_bot
ESIM_API_KEY=your-esim-api-key
ESIM_API_SECRET=your-esim-api-secret
```

### API 安全

1. **Telegram Web App 数据验证**
```go
func ValidateTelegramWebAppData(initData string, botToken string) bool {
    // 验证 Telegram Web App 初始化数据
    // 实现 HMAC-SHA256 验证逻辑
    return true
}
```

2. **请求限制**
```go
// 实现请求频率限制
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 实现限流逻辑
        next.ServeHTTP(w, r)
    })
}
```

3. **输入验证**
```go
func ValidateUserInput(input string) error {
    if len(input) > 1000 {
        return errors.New("输入内容过长")
    }
    
    // 检查恶意内容
    if containsMaliciousContent(input) {
        return errors.New("输入内容包含非法字符")
    }
    
    return nil
}
```

## 📚 API 文档

### Telegram Bot 命令

| 命令 | 描述 | 示例 |
|------|------|------|
| `/start` | 启动机器人，显示主菜单 | `/start` |
| `/products` | 查看产品列表 | `/products` |
| `/wallet` | 查看钱包余额 | `/wallet` |
| `/orders` | 查看订单历史 | `/orders` |
| `/help` | 显示帮助信息 | `/help` |

### HTTP API 端点

#### 产品相关
```
GET    /api/products          # 获取产品列表
GET    /api/products/{id}     # 获取产品详情
POST   /api/products          # 创建产品 (管理员)
PUT    /api/products/{id}     # 更新产品 (管理员)
DELETE /api/products/{id}     # 删除产品 (管理员)
```

#### 订单相关
```
GET    /api/orders            # 获取用户订单列表
POST   /api/orders            # 创建订单
GET    /api/orders/{id}       # 获取订单详情
PUT    /api/orders/{id}/pay   # 支付订单
```

#### 钱包相关
```
GET    /api/wallet            # 获取钱包信息
POST   /api/wallet/recharge   # 钱包充值
GET    /api/wallet/transactions # 获取交易记录
```

## 🤝 贡献指南

### 代码规范

1. **提交信息格式**
```
功能: 添加产品管理功能
修复: 修复订单状态更新问题
优化: 优化数据库查询性能
文档: 更新 API 文档
```

2. **代码审查清单**
- [ ] 代码遵循 Go 官方规范
- [ ] 包含完整的中文注释
- [ ] 实现了错误处理
- [ ] 添加了单元测试
- [ ] 更新了相关文档

### 开发流程

1. Fork 项目并创建功能分支
2. 实现功能并添加测试
3. 确保所有测试通过
4. 提交 Pull Request
5. 代码审查和合并

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE) 文件

## 📞 支持

如有问题或建议，请通过以下方式联系：

- 创建 GitHub Issue
- 发送邮件至项目维护者
- 加入项目讨论群组

---

**注意**: 本项目仅供学习和研究使用，请遵守相关法律法规和平台服务条款。