# 技术栈和构建系统

## 技术栈

### 核心技术
- **语言**: Go 1.24.2+
- **框架**: Telegram Bot API v5 (github.com/go-telegram-bot-api/telegram-bot-api/v5)
- **数据库**: SQLite/MySQL with GORM ORM
- **区块链**: TRON 网络 (TRC20) API

### 主要依赖
- `gorm.io/gorm` - ORM 框架
- `gorm.io/driver/mysql` - MySQL 驱动
- `gorm.io/driver/sqlite` - SQLite 驱动
- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API

## 构建系统

### 开发环境
```bash
# 安装依赖
go mod download

# 开发运行
go run cmd/bot/main.go

# 运行测试
go test ./...

# 测试覆盖率
go test -cover ./...
```

### 构建命令
```bash
# 本地构建
go build -o bot cmd/bot/main.go

# Linux 交叉编译
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o main-linux-amd64 ./cmd/bot/main.go

# 使用构建脚本
chmod +x scripts/build.sh
./scripts/build.sh
```

### Docker 构建
```bash
# 构建镜像
docker build -f docker/Dockerfile -t tg-robot-sim .

# Docker Compose 部署
docker-compose -f docker/docker-compose.yml up -d
```

## 环境配置

### 必需环境变量
- `TELEGRAM_BOT_TOKEN`: Telegram Bot 令牌 (必需)
- `TRON_API_KEY`: TRON 网络 API 密钥 (可选)

### 可选环境变量
- `DATABASE_URL`: 数据库连接字符串
- `DEBUG`: 调试模式开关
- `LOG_LEVEL`: 日志级别 (debug/info/warn/error)

### 配置文件
- 主配置: `config/config.json`
- 配置模板: `config/config.example.json`
- 配置结构: `config/config.go`

## 代码规范

### Go 代码风格
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 接口定义在 `interfaces.go` 文件中
- 错误处理遵循 Go 惯例

### 项目约定
- 所有服务接口定义在各自的 `interfaces.go` 文件中
- 使用依赖注入模式
- 配置通过环境变量和配置文件管理
- 日志使用结构化日志格式