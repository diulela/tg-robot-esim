# 产品交互流程 - Inline Mode

## 功能概述

现在支持两种产品浏览方式：
1. **传统方式**：使用 `/products` 命令查看产品列表
2. **Inline Mode**：在任何聊天中输入 `@botname 产品` 来快速浏览和分享产品

## Inline Mode 使用方法

### 1. 基本使用
在任何聊天中输入：
- `@botname` - 显示产品列表
- `@botname 产品` - 显示产品列表  
- `@botname product` - 显示产品列表
- `@botname 详情 1` - 显示产品1的详情
- `@botname 亚洲` - 搜索包含"亚洲"的产品

### 2. Inline 结果展示
用户会看到一个结果列表：
- **产品列表总览** - 点击可发送完整产品列表
- **单个产品卡片** - 每个产品一个独立的结果项

### 3. 分享功能
用户可以将产品信息直接分享到：
- 私聊对话
- 群组聊天  
- 频道
- 保存的消息

## 传统交互流程

### 1. 显示产品列表
```
🌏 亚洲区域产品

📄 第 1/3 页 (共 15 个产品)
━━━━━━━━━━━━━━━━━━

1. 亚洲多国通用 eSIM
   📊 5.0GB  ⏰ 30天  💰 15.99 USDT

2. 亚洲商务套餐 eSIM  
   📊 10.0GB  ⏰ 30天  💰 25.99 USDT

3. 亚洲旅游专享 eSIM
   📊 3.0GB  ⏰ 15天  💰 12.99 USDT

💡 点击下方「选择产品」按钮，然后回复产品编号查看详情

[🔍 选择产品]
[⬅️ 上一页] [📄 1/3] [下一页 ➡️]
[🔙 返回主菜单]
```

### 2. 用户点击"选择产品"
```
🔍 选择产品

请回复您想查看的产品编号
例如：回复 1 查看产品1的详情

💡 提示：直接输入数字即可

[🔙 返回产品列表]
```

### 3. 用户输入数字（如 "2"）
系统自动显示对应产品的详情页面。

## 技术实现

### 主要组件

1. **ProductsHandler** - 主要的产品处理器
   - `buildSimpleProductList()` - 构建简洁的产品列表
   - `buildSimpleProductListKeyboard()` - 构建包含"选择产品"按钮的键盘
   - `promptProductSelection()` - 显示产品选择提示
   - `handleProductSelection()` - 处理产品选择逻辑

2. **ProductNumberHandler** - 专门处理数字输入的消息处理器
   - `HandleMessage()` - 处理用户输入的数字
   - `CanHandle()` - 判断是否为1-5的数字消息

3. **InlineHandler** - Inline 查询处理器
   - `HandleInlineQuery()` - 处理 Inline 查询
   - `buildProductListResults()` - 构建产品列表结果
   - `buildProductDetailResults()` - 构建产品详情结果
   - `searchProducts()` - 搜索产品功能

### 回调处理

- `product_select` - 用户点击"选择产品"按钮
- `products_page:N` - 翻页到第 N 页  
- `product_detail:ID` - 查看产品详情
- `product_buy:ID` - 购买产品
- `products_back` - 返回产品列表
- `main_menu` - 返回主菜单

## 优势

### Inline Mode 优势
1. **跨聊天使用** - 可在任何聊天中快速查看产品
2. **即时分享** - 直接将产品信息分享给他人
3. **搜索功能** - 支持产品名称搜索
4. **无需切换** - 不需要离开当前聊天

### 传统模式优势  
1. **交互简洁** - 只需点击一个按钮，然后输入数字
2. **视觉清晰** - 产品列表简洁明了，信息一目了然
3. **操作直观** - 数字编号对应清晰，不会混淆
4. **扩展性好** - 可以轻松支持更多产品和页面

## 使用场景

- **个人查看**：使用 `/products` 命令
- **分享推荐**：使用 Inline Mode `@botname 产品`
- **群组讨论**：在群组中使用 Inline Mode 分享特定产品
- **客服支持**：快速向客户展示产品选项