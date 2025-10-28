# eSIM SDK for Go

eSIM 代理 API 的 Go 语言 SDK，提供完整的 API 调用封装。

## 安装

```bash
go get tg-robot-sim/pkg/esim
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    "tg-robot-sim/pkg/esim"
)

func main() {
    // 初始化客户端
    client := esim.NewClient(esim.Config{
        APIKey:    "your-api-key",
        APISecret: "your-api-secret",
        BaseURL:   "https://api.your-domain.com",
    })

    // 获取产品列表
    products, err := client.GetProducts(&esim.ProductParams{
        Country: "US",
        Limit:   10,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(products)
}
```

## 功能模块

### 产品管理
- `GetProducts()` - 获取产品列表
- `GetProduct()` - 获取产品详情
- `GetCountries()` - 获取支持的国家列表

### 订单管理
- `CreateOrder()` - 创建订单
- `GetOrders()` - 获取订单列表
- `GetOrder()` - 获取订单详情

### eSIM 管理
- `GetEsims()` - 获取 eSIM 列表
- `GetEsim()` - 获取 eSIM 详情
- `TopupEsim()` - eSIM 充值
- `GetEsimUsage()` - 获取 eSIM 使用统计

### 账户管理
- `GetAccount()` - 获取账户信息
- `GetBalance()` - 获取账户余额
- `GetFinanceRecords()` - 获取财务记录

## API 文档

详细 API 文档请访问: https://your-domain.com/api-docs
