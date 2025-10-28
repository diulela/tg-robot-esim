---
inclusion: fileMatch
fileMatchPattern: 'backend/**'
---

# Go 后端开发规范 (Telegram Bot + eSIM 电商)

## 架构模式

### 严格分层架构
- **handlers/**: 仅处理请求/响应，调用服务层，无业务逻辑
- **services/**: 业务逻辑核心，必须先定义接口再实现
- **storage/repository/**: 数据访问层，抽象数据库操作
- **storage/models/**: GORM 模型定义，包含关联关系

### 依赖注入要求
- 所有服务接口定义在 `services/interfaces.go`
- 构造函数模式：`NewXxxService(deps...) XxxService`
- 在 `main.go` 统一组装依赖，避免循环依赖

## 关键开发约定

### 新功能开发流程
1. 在 `services/interfaces.go` 定义服务接口
2. 在 `storage/models/` 定义数据模型 (GORM)
3. 在 `storage/repository/` 实现数据访问层
4. 在 `services/` 实现业务逻辑
5. 在 `handlers/` 创建请求处理器
6. 在相应的 `main.go` 中注册路由和依赖

### 必须遵循的规范
- 所有函数返回 `error` 作为最后一个参数
- 使用 `context.Context` 作为第一个参数传递上下文
- 错误信息使用中文，便于用户理解
- 数据库操作必须包含事务处理
- Telegram Bot 交互使用 InlineKeyboard

## 代码风格要求

### 函数签名模式
```go
// 标准服务方法签名
func (s *serviceImpl) MethodName(ctx context.Context, req *RequestType) (*ResponseType, error)

// 仓储方法签名
func (r *repoImpl) MethodName(ctx context.Context, params ...interface{}) (*Model, error)

// 处理器方法签名
func (h *handler) HandleMethod(w http.ResponseWriter, r *http.Request)
```

### 错误处理模式
```go
// 包装错误，提供中文消息
if err := someOperation(); err != nil {
    return fmt.Errorf("操作失败: %w", err)
}

// 验证错误，直接返回中文消息
if req.TelegramID == 0 {
    return errors.New("Telegram ID 不能为空")
}

// 数据库错误处理
if err := db.Create(&model).Error; err != nil {
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return errors.New("记录已存在")
    }
    return fmt.Errorf("数据库操作失败: %w", err)
}
```

### 注释规范
- 所有导出函数必须有中文注释
- 注释格式：`// FunctionName 功能描述`
- 复杂业务逻辑添加行内中文注释

## 服务层实现模式

### 接口定义要求 (services/interfaces.go)
```go
// 所有服务接口必须在此文件定义
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error)
    GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
}

type OrderService interface {
    CreateOrder(ctx context.Context, req *CreateOrderRequest) (*models.Order, error)
    ProcessPayment(ctx context.Context, orderID int64, txHash string) error
}
```

### 服务实现模式
```go
// 私有结构体实现接口
type userService struct {
    userRepo repository.UserRepository
    logger   Logger
}

// 构造函数返回接口类型
func NewUserService(userRepo repository.UserRepository, logger Logger) UserService {
    return &userService{userRepo: userRepo, logger: logger}
}

// 方法实现：验证 -> 业务逻辑 -> 数据操作 -> 日志记录
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
    // 1. 输入验证
    if req.TelegramID == 0 {
        return nil, errors.New("Telegram ID 不能为空")
    }
    
    // 2. 业务逻辑检查
    if exists, _ := s.userRepo.ExistsByTelegramID(ctx, req.TelegramID); exists {
        return nil, errors.New("用户已存在")
    }
    
    // 3. 数据操作
    user := &models.User{TelegramID: req.TelegramID, Username: req.Username}
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }
    
    // 4. 日志记录
    s.logger.Info("用户创建成功", "user_id", user.ID)
    return user, nil
}
```

## 数据层实现规范

### GORM 模型定义
```go
// 标准模型结构
type User struct {
    ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    TelegramID int64          `gorm:"uniqueIndex;not null" json:"telegram_id"`
    Username   string         `gorm:"size:255" json:"username"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 关联关系
    Orders []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// 必须实现 TableName 方法
func (User) TableName() string { return "users" }
```

### 仓储层模式
```go
// 仓储接口 (storage/repository/interfaces.go)
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
    ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error)
}

// 仓储实现
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// 标准 CRUD 操作模式
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

## Telegram Bot 开发模式

### Bot 处理器结构
```go
// 标准 Bot 处理器
type BotHandler struct {
    bot         *tgbotapi.BotAPI
    userService services.UserService
    orderService services.OrderService
}

// 更新消息分发
func (h *BotHandler) HandleUpdate(update tgbotapi.Update) error {
    switch {
    case update.Message != nil:
        return h.handleMessage(update.Message)
    case update.CallbackQuery != nil:
        return h.handleCallbackQuery(update.CallbackQuery)
    }
    return nil
}
```

### 交互设计原则
- 所有用户交互使用 InlineKeyboard，避免 ReplyKeyboard
- 按钮回调数据格式：`action:param1:param2`
- 消息文本使用中文，包含 emoji 提升用户体验
- 错误消息友好提示，不暴露技术细节

### 命令处理模式
```go
// 命令分发
func (h *BotHandler) handleMessage(message *tgbotapi.Message) error {
    // 1. 确保用户存在
    user, err := h.ensureUserExists(message.From)
    if err != nil {
        return h.sendErrorMessage(message.Chat.ID, "系统错误，请稍后重试")
    }
    
    // 2. 命令分发
    switch message.Command() {
    case "start":
        return h.handleStart(message, user)
    case "products":
        return h.handleProducts(message, user)
    default:
        return h.handleUnknown(message, user)
    }
}

// 回调查询处理
func (h *BotHandler) handleCallbackQuery(callback *tgbotapi.CallbackQuery) error {
    parts := strings.Split(callback.Data, ":")
    action := parts[0]
    
    switch action {
    case "buy_product":
        return h.handleBuyProduct(callback, parts[1])
    case "confirm_order":
        return h.handleConfirmOrder(callback, parts[1])
    }
    return nil
}
```

## HTTP API 开发规范

### 路由注册模式
```go
// 在 main.go 中注册路由
func setupRoutes(mux *http.ServeMux, handlers *HTTPHandler) {
    // Mini App API
    mux.HandleFunc("GET /api/products", handlers.GetProducts)
    mux.HandleFunc("POST /api/orders", handlers.CreateOrder)
    mux.HandleFunc("GET /api/wallet/{telegramId}", handlers.GetWallet)
    
    // 健康检查
    mux.HandleFunc("GET /health", handlers.Health)
}
```

### 统一响应格式
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// 标准响应方法
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(APIResponse{
        Success: statusCode < 400,
        Data:    data,
    })
}

func SendError(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(APIResponse{
        Success: false,
        Error:   message,
    })
}
```

### 请求处理模式
```go
func (h *HTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // 1. 解析请求
    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        SendError(w, 400, "请求格式错误")
        return
    }
    
    // 2. 验证 Telegram Web App 数据
    if !h.validateTelegramWebAppData(r.Header.Get("Authorization")) {
        SendError(w, 401, "未授权访问")
        return
    }
    
    // 3. 调用服务层
    order, err := h.orderService.CreateOrder(r.Context(), &req)
    if err != nil {
        SendError(w, 500, "创建订单失败")
        return
    }
    
    // 4. 返回结果
    SendJSON(w, 201, order)
}
```

## 关键业务逻辑

### 区块链支付集成
- 使用 TRON 网络 USDT-TRC20
- 异步监控交易确认状态
- 支付确认前不发放 eSIM 产品
- 实现交易失败的退款机制

### 用户状态管理
- 通过 Telegram ID 唯一标识用户
- 维护用户会话状态支持多步骤对话
- 实现用户钱包余额管理

### 订单处理流程
1. 创建订单 -> 等待支付状态
2. 检测到支付 -> 处理中状态  
3. 发放产品 -> 完成状态
4. 支付失败 -> 取消状态

### 错误处理策略
- 所有用户可见错误使用中文提示
- 系统错误记录详细日志但不暴露给用户
- 网络错误实现重试机制
- 数据库操作失败进行事务回滚

## 配置和环境管理

### 配置加载模式
```go
// 配置结构
type Config struct {
    Server     ServerConfig     `json:"server"`
    Database   DatabaseConfig   `json:"database"`
    Telegram   TelegramConfig   `json:"telegram"`
    Blockchain BlockchainConfig `json:"blockchain"`
}

// 环境变量优先级高于配置文件
func LoadConfig() (*Config, error) {
    config := &Config{}
    
    // 从 JSON 文件加载基础配置
    if data, err := os.ReadFile("config/config.json"); err == nil {
        json.Unmarshal(data, config)
    }
    
    // 环境变量覆盖
    if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
        config.Telegram.BotToken = token
    }
    
    return config, nil
}
```

### 必需环境变量
- `TELEGRAM_BOT_TOKEN`: Telegram Bot 令牌
- `DATABASE_URL`: 数据库连接字符串  
- `TRON_API_KEY`: TRON 网络 API 密钥
- `WEBHOOK_URL`: Telegram Webhook URL (生产环境)

### 日志记录要求
- 使用结构化日志 (slog)
- 敏感信息不记录到日志
- 错误日志包含足够上下文信息
- 生产环境日志级别设为 INFO