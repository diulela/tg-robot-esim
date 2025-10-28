# Telegram Robot Simulator

基于 Go 语言开发的 Telegram 机器人核心系统，提供对话处理、交互式菜单和区块链交易监控功能。

## 功能特性

- 🤖 **智能对话处理**: 支持命令处理、自定义键盘、内联按钮和回调查询
- 📱 **交互式菜单系统**: 提供用户友好的导航界面
- 💰 **区块链交易监控**: USDT-TRC20 充值系统，支持区块链交易验证
- 🔔 **实时通知服务**: 多种类型的消息通知支持
- 📊 **会话管理**: 用户会话状态管理和持久化

## 技术栈

- **语言**: Go 1.23.5+
- **框架**: Telegram Bot API v5
- **数据库**: SQLite/MySQL with GORM
- **区块链**: TRON 网络 (TRC20) API

## 项目结构

```
tg-robot-sim/
├── cmd/                    # 应用程序入口点
│   └── bot/               # 机器人主程序
├── config/                # 配置文件
├── handlers/              # 请求处理器
├── models/                # 数据模型
├── services/              # 业务逻辑层
├── storage/               # 数据存储层
├── pkg/                   # 可复用包
└── utils/                 # 工具函数
```

## 快速开始

### 1. 环境要求

- Go 1.23.5 或更高版本
- SQLite 或 MySQL 数据库
- Telegram Bot Token

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置

复制配置文件模板：

```bash
cp config/config.example.json config/config.json
```

设置环境变量：

```bash
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export TRON_API_KEY="your_tron_api_key_here"
```

### 4. 运行

```bash
# 开发环境
go run cmd/bot/main.go

# 编译运行
go build -o bot cmd/bot/main.go
./bot
```

### 5. 交叉编译 (Linux)

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o main-linux-amd64 ./cmd/bot/main.go
```

## 配置说明

### 环境变量

- `TELEGRAM_BOT_TOKEN`: Telegram Bot 令牌 (必需)
- `DATABASE_URL`: 数据库连接字符串
- `TRON_API_KEY`: TRON 网络 API 密钥
- `DEBUG`: 调试模式开关
- `LOG_LEVEL`: 日志级别

### 配置文件

详细配置请参考 `config/config.example.json`

## 开发指南

### 添加新的命令处理器

1. 在 `handlers/` 目录创建新的处理器文件
2. 实现 `CommandHandler` 接口
3. 在 `main.go` 中注册处理器

### 添加新的业务服务

1. 在 `services/` 目录创建服务文件
2. 定义服务接口和实现
3. 在需要的处理器中注入服务

### 添加新的数据模型

1. 在 `models/` 目录定义结构体
2. 添加 GORM 标签进行数据库映射
3. 在数据库迁移中添加表结构

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./services/...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 许可证

MIT License