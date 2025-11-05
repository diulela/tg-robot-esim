# API 服务模块

这个目录包含了按功能模块拆分的 API 服务。

## 文件结构

```
api/
├── client.ts      # API 客户端基础类，包含请求拦截器、错误处理等
├── product.ts     # 产品相关 API
├── order.ts       # 订单相关 API  
├── wallet.ts      # 钱包相关 API
├── region.ts      # 区域和国家相关 API
├── user.ts        # 用户相关 API
├── system.ts      # 系统相关 API
├── index.ts       # 统一导出文件
└── README.md      # 说明文档
```

## 使用方式

### 导入整个 API 对象（推荐）
```typescript
import api from '@/services/api'

// 使用钱包 API
const wallet = await api.wallet.getWallet()

// 使用产品 API
const products = await api.product.getProducts()
```

### 导入特定的 API 模块
```typescript
import { walletApi, productApi } from '@/services/api'

// 直接使用模块
const wallet = await walletApi.getWallet()
const products = await productApi.getProducts()
```

### 导入 API 客户端
```typescript
import { apiClient } from '@/services/api'

// 直接使用客户端进行请求
const data = await apiClient.get('/custom/endpoint')
```

## 模块说明

### ApiClient (client.ts)
- 基于 Axios 的 HTTP 客户端
- 自动添加 Telegram 初始化数据
- 请求重试机制
- 统一错误处理
- 请求/响应拦截器

### ProductApi (product.ts)
- 产品列表查询
- 产品详情获取
- 产品搜索
- 数据格式转换

### OrderApi (order.ts)
- 订单创建和管理
- 订单状态查询
- 订单历史记录
- 支付重试

### WalletApi (wallet.ts)
- 钱包余额查询
- 交易记录获取
- USDT 充值订单管理
- 充值历史查询

### RegionApi (region.ts)
- 区域列表获取
- 国家信息查询
- 地理位置相关功能

### UserApi (user.ts)
- 用户信息管理
- 用户统计数据
- 个人资料更新

### SystemApi (system.ts)
- 系统健康检查
- 系统配置获取
- 版本信息查询

## 向后兼容

原有的 `@/services/api` 导入方式仍然有效，所有现有代码无需修改。

## 扩展新模块

1. 在 `api/` 目录下创建新的模块文件
2. 实现相应的 API 类
3. 在 `index.ts` 中导出新模块
4. 在统一的 API 对象中添加新模块