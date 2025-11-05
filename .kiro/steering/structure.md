---
inclusion: always
---

# 项目架构和开发规范

## 项目概述

这是一个基于 Telegram 的 eSIM 电商系统，包含：
- **后端**: Go 语言开发的 Telegram Bot 和 HTTP API 服务
- **前端**: Vue 3.0 + TypeScript + Vuetify 的 Telegram Mini App
- **数据库**: SQLite/MySQL with GORM
- **区块链**: TRON 网络 USDT-TRC20 支付集成

## 项目结构规范

### 后端目录结构 (backend/)
```
backend/
├── cmd/                    # 应用程序入口点
│   ├── bot/               # Telegram Bot 主程序
│   ├── gm/                # 管理工具程序
│   └── miniapp/           # Mini App HTTP 服务器
├── api/                   # API 路由定义
├── config/                # 配置管理
├── handlers/              # HTTP 和 Bot 请求处理器
│   └── bot/              # Bot 专用处理器
├── services/              # 业务逻辑层 (接口定义和实现)
├── storage/               # 数据存储层
│   ├── data/             # 数据库连接和迁移
│   ├── models/           # GORM 数据模型
│   └── repository/       # 数据访问层
├── server/                # HTTP 服务器配置
│   └── middleware/       # 中间件
├── pkg/                   # 可复用包和工具
│   ├── bot/              # Bot 相关工具
│   ├── logger/           # 日志工具
│   ├── telegram/         # Telegram SDK 封装
│   └── tron/             # TRON 区块链工具
├── scripts/               # 构建和部署脚本
├── docs/                  # 项目文档
└── utils/                 # 通用工具函数
```

### 前端目录结构 (miniapp/)
```
miniapp/
├── src/
│   ├── components/        # 可复用 UI 组件
│   │   └── layout/       # 布局组件
│   ├── views/            # 页面视图组件
│   ├── router/           # Vue Router 配置
│   │   └── guards.ts     # 路由守卫
│   ├── stores/           # Pinia 状态管理
│   ├── services/         # API 服务层
│   │   └── api/          # API 接口定义
│   ├── plugins/          # Vue 插件配置
│   ├── types/            # TypeScript 类型定义
│   ├── utils/            # 前端工具函数
│   └── styles/           # 样式文件
├── docs/                 # 前端文档
├── styles/               # 全局样式文件
├── index.html           # 入口 HTML
├── vite.config.ts       # Vite 构建配置
├── tsconfig.json        # TypeScript 配置
└── package.json         # 依赖管理
```

## 架构设计原则

### 分层架构
1. **表示层** (Handlers): 处理 HTTP 请求和 Telegram 消息
2. **业务层** (Services): 核心业务逻辑，定义接口
3. **数据层** (Repository): 数据访问抽象
4. **存储层** (Storage): 数据库操作实现

### 依赖注入
- 所有服务通过接口定义，支持依赖注入
- 在 `main.go` 中组装所有依赖关系
- 避免硬编码依赖，提高可测试性

### 接口优先设计
- 所有服务必须先定义接口 (在 `services/interfaces.go`)
- 实现与接口分离，支持多种实现方式
- 便于单元测试和模拟

## 开发工作流程

### 新功能开发步骤
1. **需求分析**: 明确功能需求和业务逻辑
2. **接口设计**: 在 `services/interfaces.go` 定义服务接口
3. **数据模型**: 在 `storage/models/` 定义 GORM 数据模型
4. **仓储层**: 在 `storage/repository/` 实现数据访问层
5. **服务实现**: 在 `services/` 实现业务逻辑服务
6. **处理器实现**: 在 `handlers/` 创建 HTTP/Bot 请求处理器
7. **路由注册**: 在 `api/` 或相应的注册文件中添加路由
8. **前端类型**: 在 `miniapp/src/types/` 定义 TypeScript 类型
9. **前端服务**: 在 `miniapp/src/services/` 实现 API 调用
10. **前端组件**: 在 `miniapp/src/components/` 或 `views/` 实现 UI
11. **状态管理**: 在 `miniapp/src/stores/` 添加 Pinia 状态
12. **路由配置**: 在 `miniapp/src/router/` 配置页面路由
13. **测试验证**: 进行功能测试和集成测试

### 代码提交规范
- 每个功能模块独立提交
- 提交信息使用中文，格式：`功能: 简短描述`
- 确保代码通过 `go fmt` 和基本测试

## 技术栈约定

### 后端技术栈
- **语言**: Go 1.24.2+
- **Web框架**: 标准库 `net/http`
- **ORM**: GORM v1.30.0
- **Telegram**: telegram-bot-api v5.5.1
- **数据库**: SQLite (开发) / MySQL (生产)
- **区块链**: TRON 网络集成

### 前端技术栈
- **框架**: Vue 3.0 + Composition API
- **语言**: TypeScript 5.3+
- **UI库**: Vuetify 3.4+ (Material Design)
- **状态管理**: Pinia 2.1+
- **路由**: Vue Router 4.2+
- **构建工具**: Vite 5.0+
- **HTTP客户端**: Axios 1.6+
- **工具库**: Day.js, QRCode.js
- **开发工具**: ESLint, Prettier, VConsole

### 部署和运维
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **环境管理**: 环境变量 + JSON 配置文件
- **构建优化**: 代码分割、资源压缩、缓存策略

## 安全和性能规范

### 安全要求
- 所有敏感信息通过环境变量管理
- API 接口必须验证 Telegram Web App 数据
- 数据库操作必须防止 SQL 注入
- 用户输入必须进行验证和清理

### 性能要求
- 数据库查询必须使用索引优化
- API 响应时间控制在 200ms 内
- 实现适当的缓存策略
- 区块链交易采用异步处理

### 错误处理
- 所有函数必须返回错误信息
- 使用结构化日志记录
- 向用户提供友好的中文错误提示
- 实现优雅的服务降级

## 测试策略

### 单元测试
- 所有服务层必须有单元测试
- 使用接口模拟进行隔离测试
- 测试覆盖率目标 80%+

### 集成测试
- API 端点集成测试
- 数据库操作测试
- Telegram Bot 交互测试

### 手动测试
- Telegram 环境功能测试
- 支付流程端到端测试
- 用户体验测试
## 核
心服务接口规范

### 已实现的核心服务接口
- **DialogService**: 对话服务，处理用户消息和会话管理
- **MenuService**: 菜单服务，管理交互式菜单系统
- **BlockchainService**: 区块链服务，监控和验证 TRON 交易
- **NotificationService**: 通知服务，发送各种类型的通知
- **RechargeService**: 充值服务，处理用户充值业务逻辑
- **WalletHistoryService**: 钱包历史服务，管理交易历史记录

### 服务接口设计原则
- 所有服务接口必须在 `services/interfaces.go` 中定义
- 使用 `context.Context` 作为第一个参数
- 返回值最后一个参数必须是 `error`
- 结构体字段使用 JSON 标签便于序列化
- 枚举类型使用自定义类型和常量定义

## 前端架构特点

### Vue 3.0 + Composition API
- 使用 `<script setup>` 语法糖
- 组合式 API 替代选项式 API
- TypeScript 类型安全
- 响应式状态管理

### Vuetify 3.4 Material Design
- 现代化 UI 组件库
- 响应式设计支持
- 主题定制能力
- 移动端优化

### Pinia 状态管理
- 替代 Vuex 的现代状态管理
- TypeScript 原生支持
- 模块化状态设计
- 开发工具集成

### 关键页面组件
- **HomePage**: 首页，产品展示和导航
- **ProductPage/ProductDetailPage**: 产品列表和详情
- **OrderPage/OrderDetailPage**: 订单管理
- **WalletPage**: 钱包管理和充值
- **ProfilePage**: 用户个人资料
- **各类产品页面**: 热门、本地、全球、区域产品

## Docker 容器化部署

### 容器配置
- **Dockerfile**: 后端 Go 应用容器
- **Dockerfile.miniapp**: 前端 Vue 应用容器
- **docker-compose.yml**: 完整服务编排
- **docker-compose.miniapp.yml**: 仅前端服务
- **nginx.conf**: 反向代理配置

### 部署环境
- 开发环境：本地 Docker Compose
- 生产环境：容器化部署
- 静态资源：Nginx 代理
- API 服务：Go HTTP 服务器

## 开发环境配置

### 后端开发
```bash
cd backend
go mod tidy
go run cmd/bot/main.go      # 启动 Bot 服务
go run cmd/miniapp/main.go  # 启动 HTTP API 服务
```

### 前端开发
```bash
cd miniapp
npm install
npm run dev     # 启动开发服务器 (localhost:8082)
npm run build   # 构建生产版本
```

### 环境变量配置
- 后端：`backend/config/config.json` + 环境变量
- 前端：`.env.development` / `.env.example`
- Docker：`docker-compose.yml` 环境变量配置

## 代码质量保证

### 后端代码规范
- Go 标准格式化：`go fmt`
- 错误处理：所有函数返回 error
- 中文注释：导出函数必须有中文注释
- 接口优先：先定义接口再实现

### 前端代码规范
- TypeScript 类型检查：`vue-tsc --noEmit`
- ESLint 代码检查：`eslint . --fix`
- Prettier 代码格式化：`prettier --write src/`
- Vue 3 组合式 API 规范

### 测试策略
- 后端：单元测试 + 集成测试
- 前端：组件测试 + E2E 测试
- API：接口测试和文档
- 用户体验：Telegram 环境测试