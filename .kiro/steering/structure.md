---
inclusion: always
---

# 项目架构和开发规范

## 项目概述

这是一个基于 Telegram 的 eSIM 电商系统，包含：
- **后端**: Go 语言开发的 Telegram Bot 和 HTTP API 服务
- **前端**: Telegram Mini App (JavaScript/HTML/CSS)
- **数据库**: SQLite/MySQL with GORM
- **区块链**: TRON 网络 USDT-TRC20 支付集成

## 项目结构规范

### 后端目录结构 (backend/)
```
backend/
├── cmd/                    # 应用程序入口点
│   ├── bot/               # Telegram Bot 主程序
│   └── miniapp/           # Mini App HTTP 服务器
├── config/                # 配置管理
├── handlers/              # HTTP 和 Bot 请求处理器
├── services/              # 业务逻辑层 (接口定义)
├── storage/               # 数据存储层
│   ├── data/             # 数据库连接和迁移
│   └── repository/       # 数据访问层
├── pkg/                   # 可复用包和工具
└── utils/                 # 通用工具函数
```

### 前端目录结构 (miniapp/)
```
miniapp/
├── src/
│   ├── components/        # 可复用 UI 组件
│   ├── pages/            # 页面组件
│   ├── services/         # API 服务层
│   └── utils/            # 前端工具函数
├── styles/               # CSS 样式文件
├── index.html           # 入口 HTML
└── vite.config.js       # 构建配置
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
3. **数据模型**: 在相应目录定义数据结构和数据库模型
4. **服务实现**: 实现业务逻辑服务
5. **处理器实现**: 创建 HTTP/Bot 请求处理器
6. **路由注册**: 在相应的注册文件中添加路由
7. **前端集成**: 实现前端页面和 API 调用
8. **测试验证**: 进行功能测试和集成测试

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

### 前端技术栈
- **构建工具**: Vite 5.0+
- **样式**: 原生 CSS (CSS变量 + Flexbox/Grid)
- **JavaScript**: ES6+ 模块化
- **Telegram SDK**: Telegram Web App API

### 部署和运维
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **环境管理**: 环境变量 + JSON 配置文件

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