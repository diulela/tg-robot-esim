# Telegram Mini App - eSIM 商城

这是一个基于 Telegram Mini App 的 eSIM 电商前端应用。

## 项目结构

```
miniapp/
├── src/
│   ├── components/      # 可复用组件（待实现）
│   ├── pages/          # 页面组件
│   ├── services/       # API 服务（待实现）
│   └── utils/          # 工具函数
│       ├── telegram.js # Telegram Web App 集成
│       ├── router.js   # 路由系统
│       ├── state.js    # 状态管理
│       ├── formatter.js # 数据格式化
│       └── validator.js # 输入验证
├── styles/
│   ├── main.css        # 主样式
│   └── components.css  # 组件样式
├── index.html          # 入口 HTML
├── package.json        # 项目配置
└── vite.config.js      # Vite 构建配置
```

## 开发环境设置

### 安装依赖

```bash
cd miniapp
npm install
```

### 启动开发服务器

```bash
npm run dev
```

开发服务器将在 http://localhost:3000 启动。

### 构建生产版本

```bash
npm run build
```

构建输出将在 `dist/` 目录中。

## 功能特性

- ✅ Telegram Web App SDK 集成
- ✅ 基于 Hash 的路由系统
- ✅ 状态管理（localStorage）
- ✅ 响应式设计
- ✅ 主题适配（Telegram 主题）
- ⏳ 产品列表和详情（待实现）
- ⏳ 购买流程（待实现）
- ⏳ 钱包管理（待实现）
- ⏳ 充值功能（待实现）

## 技术栈

- **构建工具**: Vite
- **样式**: 原生 CSS（CSS 变量 + Flexbox/Grid）
- **JavaScript**: ES6+ 模块
- **Telegram**: Telegram Web App SDK

## 开发说明

### 在 Telegram 中测试

1. 创建一个 Telegram Bot
2. 设置 Mini App URL 为你的开发服务器地址
3. 在 Telegram 中打开 Mini App 进行测试

### 本地开发（无 Telegram）

应用包含模拟的 Telegram Web App 环境，可以在浏览器中直接开发和测试。

## 下一步

根据实施计划，接下来将实现：

1. 后端数据模型和 API
2. 前端 API 服务层
3. UI 组件
4. 页面功能
5. 系统集成

详见 `.kiro/specs/tg-mini-app-ecommerce/tasks.md`
