# eSIM 产品功能使用指南

## 功能概述

Telegram 机器人现已集成 eSIM 产品浏览和购买功能，用户可以通过以下方式浏览和购买 eSIM 产品：

### 主要功能

1. **产品分类浏览**
   - 本地产品（单个国家）
   - 区域产品（多个国家）
   - 全球产品（全球通用）

2. **产品搜索**
   - 通过国家代码搜索产品
   - 查看产品详细信息

3. **产品购买**
   - 创建订单
   - 查看订单状态
   - 获取 eSIM 信息

## 使用方式

### 1. 通过主菜单浏览

用户发送 `/start` 命令后，点击 "🛍️ 浏览产品" 按钮，可以看到三种产品类型：

- 🏠 本地产品
- 🌏 区域产品
- 🌍 全球产品

点击任意类型即可查看该类型的产品列表。

### 2. 通过命令搜索

用户可以使用 `/products` 命令搜索特定国家的产品：

```
/products          # 显示产品类型选择菜单
/products CN       # 搜索中国的产品
/products US       # 搜索美国的产品
/products JP       # 搜索日本的产品
```

### 3. 产品列表功能

产品列表页面显示：
- 产品名称
- 支持的国家
- 流量大小和有效期
- 零售价和代理价
- 分页导航（每页5个产品）

用户可以点击 "查看详情" 按钮查看产品完整信息。

### 4. 产品详情

产品详情页面显示：
- 产品完整名称（中英文）
- 产品类型
- 支持的国家列表
- 流量大小
- 有效期
- 价格信息
- 产品特性

用户可以点击 "🛒 立即购买" 开始购买流程。

## 配置说明

### 环境变量

在 `.env` 文件或系统环境变量中配置：

```bash
# eSIM API 配置
ESIM_API_KEY=your-api-key
ESIM_API_SECRET=your-api-secret
ESIM_BASE_URL=https://api.your-domain.com
```

### 配置文件

在 `config/config.json` 中配置：

```json
{
  "esim_sdk": {
    "api_key": "your-api-key",
    "api_secret": "your-api-secret",
    "base_url": "https://api.your-domain.com"
  }
}
```

## 代码集成

### 初始化 eSIM 服务

```go
import (
    "tg-robot-sim/services"
    "tg-robot-sim/config"
)

// 创建 eSIM 服务
esimService := services.NewEsimService(
    config.EsimSDK.APIKey,
    config.EsimSDK.APISecret,
    config.EsimSDK.BaseURL,
)
```

### 注册处理器

```go
import (
    "tg-robot-sim/handlers"
)

// 创建产品处理器
productsHandler := handlers.NewProductsHandler(bot, esimService, logger)

// 注册命令处理器
registry.RegisterCommandHandler(productsHandler)

// 注册回调处理器
registry.RegisterCallbackHandler(productsHandler)
```

## API 接口

### 获取产品列表

```go
products, err := esimService.GetProducts(ctx, &esim.ProductParams{
    Type:  esim.ProductTypeLocal,
    Page:  1,
    Limit: 20,
})
```

### 搜索产品

```go
products, err := esimService.GetProducts(ctx, &esim.ProductParams{
    Country: "CN",
    Limit:   10,
})
```

### 获取产品详情

```go
product, err := esimService.GetProduct(ctx, productID)
```

### 创建订单

```go
order, err := esimService.CreateOrder(ctx, esim.CreateOrderRequest{
    ProductID:     1,
    CustomerEmail: "customer@example.com",
    CustomerPhone: "+86 138 0000 0000",
    Quantity:      1,
})
```

## 用户交互流程

### 浏览产品流程

1. 用户发送 `/start` 或 `/products`
2. 选择产品类型（本地/区域/全球）
3. 浏览产品列表
4. 点击查看产品详情
5. 点击立即购买

### 搜索产品流程

1. 用户发送 `/products CN`
2. 显示中国相关的产品列表
3. 点击查看产品详情
4. 点击立即购买

### 购买流程

1. 点击 "🛒 立即购买"
2. 输入客户信息（邮箱、手机号、数量）
3. 确认订单信息
4. 完成支付
5. 获取 eSIM 信息

## 注意事项

1. **API 密钥安全**：请妥善保管 API 密钥，不要提交到代码仓库
2. **错误处理**：所有 API 调用都应该有适当的错误处理
3. **用户体验**：产品列表使用分页，避免一次加载过多数据
4. **价格显示**：同时显示零售价和代理价，方便代理商了解利润空间

## 扩展功能

未来可以添加的功能：

- [ ] 订单管理（查看历史订单）
- [ ] eSIM 使用情况查询
- [ ] eSIM 充值功能
- [ ] 收藏夹功能
- [ ] 产品推荐
- [ ] 多语言支持
- [ ] 支付集成

## 故障排查

### 产品列表为空

检查：
1. API 密钥是否正确配置
2. API 服务是否正常运行
3. 网络连接是否正常
4. 查看日志获取详细错误信息

### 无法创建订单

检查：
1. 产品 ID 是否有效
2. 客户信息格式是否正确
3. API 账户余额是否充足
4. 查看 API 返回的错误信息

## 技术支持

如有问题，请联系技术支持或查看 API 文档：
- API 文档：https://your-domain.com/api-docs
- 技术支持：support@your-domain.com
