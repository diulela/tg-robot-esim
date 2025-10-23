# 项目结构和组织

## 目录结构

```
tg-robot-sim/
├── cmd/                    # 应用程序入口点
│   └── bot/               # 机器人主程序
│       └── main.go        # 主入口文件
├── config/                # 配置管理
│   ├── config.go          # 配置结构定义
│   ├── config.json        # 实际配置文件
│   └── config.example.json # 配置模板
├── docker/                # Docker 相关文件
├── handlers/              # HTTP/Telegram 请求处理器
│   ├── interfaces.go      # 处理器接口定义
│   ├── registry.go        # 处理器注册
│   ├── middleware.go      # 中间件
│   └── *_handler.go       # 具体处理器实现
├── pkg/                   # 可复用的公共包
│   ├── bot/              # Bot 相关工具
│   ├── logger/           # 日志工具
│   ├── retry/            # 重试机制
│   └── tron/             # TRON 区块链工具
├── services/              # 业务逻辑层
│   ├── interfaces.go      # 服务接口定义
│   └── *_service.go       # 具体服务实现
├── storage/               # 数据存储层
├── utils/                 # 工具函数
└── scripts/               # 构建和部署脚本
```

## 架构模式

### 分层架构
1. **入口层** (`cmd/`): 应用程序启动和初始化
2. **处理器层** (`handlers/`): 请求处理和路由
3. **服务层** (`services/`): 业务逻辑实现
4. **存储层** (`storage/`): 数据持久化
5. **工具层** (`pkg/`, `utils/`): 可复用组件

### 设计原则
- **接口分离**: 每个模块都有对应的 `interfaces.go` 文件
- **依赖注入**: 通过接口进行依赖管理
- **单一职责**: 每个服务专注于特定业务领域
- **模块化**: 功能按模块组织，便于维护和扩展

## 核心模块

### 处理器模块 (`handlers/`)
- `start_handler.go`: 启动命令处理
- `help_handler.go`: 帮助命令处理
- `message_handler.go`: 消息处理
- `callback_handler.go`: 回调查询处理
- `middleware.go`: 中间件逻辑
- `registry.go`: 处理器注册管理

### 服务模块 (`services/`)
- `dialog_service.go`: 对话管理服务
- `menu_service.go`: 菜单系统服务
- `session_service.go`: 会话管理服务
- `blockchain_service.go`: 区块链交易服务
- `notification_service.go`: 通知服务

### 公共包 (`pkg/`)
- `bot/`: Telegram Bot 相关工具
- `logger/`: 日志记录工具
- `retry/`: 重试机制实现
- `tron/`: TRON 区块链集成

## 文件命名约定

### Go 文件
- 服务文件: `*_service.go`
- 处理器文件: `*_handler.go`
- 接口定义: `interfaces.go`
- 测试文件: `*_test.go`

### 配置文件
- 主配置: `config.json`
- 配置模板: `config.example.json`
- 环境配置: `.env`

## 扩展指南

### 添加新处理器
1. 在 `handlers/` 目录创建 `new_handler.go`
2. 实现 `CommandHandler` 接口
3. 在 `registry.go` 中注册处理器

### 添加新服务
1. 在 `services/interfaces.go` 中定义接口
2. 在 `services/` 目录创建 `new_service.go`
3. 实现服务接口
4. 在需要的处理器中注入服务

### 添加新的可复用包
1. 在 `pkg/` 目录创建新包目录
2. 定义包的公共接口
3. 实现具体功能
4. 编写相应的测试