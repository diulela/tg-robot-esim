# 购买流程添加邮箱验证 - 修改说明

## 修改概述

为了满足第三方 eSIM API 的要求，在用户购买 eSIM 产品时添加了邮箱地址的必填验证。邮箱将用于接收 eSIM 激活信息。

## 修改文件清单

### 后端修改 (3个文件)

#### 1. `backend/services/interfaces.go`
- ✅ 在 `CreateEsimOrderRequest` 结构体中添加 `CustomerEmail` 字段
- ✅ 添加 `validate:"required,email"` 标签

```go
type CreateEsimOrderRequest struct {
    UserID        int64  `json:"user_id" validate:"required"`
    ProductID     int    `json:"product_id" validate:"required"`
    Quantity      int    `json:"quantity" validate:"required,min=1"`
    TotalAmount   string `json:"total_amount" validate:"required"`
    CustomerEmail string `json:"customer_email" validate:"required,email"` // 新增
    Remark        string `json:"remark,omitempty"`
}
```

#### 2. `backend/services/order_service.go`
- ✅ 添加 `isValidEmail()` 邮箱格式验证函数
- ✅ 在 `CreateEsimOrder()` 中添加邮箱验证逻辑
- ✅ 修改 `createProviderOrder()` 方法，使用用户提供的邮箱

**关键修改点：**
```go
// 验证邮箱
if req.CustomerEmail == "" {
    return nil, errors.New("客户邮箱不能为空")
}
if !isValidEmail(req.CustomerEmail) {
    return nil, errors.New("邮箱格式不正确")
}

// 调用第三方 API 时使用真实邮箱
providerOrderID, err := s.createProviderOrder(ctx, order, product, req.CustomerEmail)
```

#### 3. `backend/api/esim_orders.go`
- ✅ 在 `handleCreateEsimOrder()` 中添加邮箱必填验证
- ✅ 添加邮箱相关的错误处理

```go
if req.CustomerEmail == "" {
    h.sendError(w, http.StatusBadRequest, "Customer email is required", "")
    return
}
```

### 前端修改 (4个文件)

#### 4. `miniapp/src/types/esim-order.ts`
- ✅ 在 `CreateEsimOrderRequest` 接口中添加 `customerEmail` 字段

```typescript
export interface CreateEsimOrderRequest {
  productId: number
  quantity: number
  totalAmount: string
  customerEmail: string  // 新增
  remark?: string
}
```

#### 5. `miniapp/src/components/business/EmailInput.vue` (新建)
- ✅ 创建可复用的邮箱输入组件
- ✅ 实时邮箱格式验证
- ✅ 支持禁用状态和错误提示
- ✅ Material Design 风格，集成 Vuetify

**功能特性：**
- 实时格式验证（正则表达式）
- 必填验证
- 成功状态图标显示
- 友好的提示信息
- 响应式设计

#### 6. `miniapp/src/components/business/PurchaseModal.vue`
- ✅ 导入 `EmailInput` 组件
- ✅ 添加邮箱输入框（位于数量选择器和备注之间）
- ✅ 添加邮箱验证逻辑
- ✅ 在提交订单时包含邮箱参数
- ✅ 添加邮箱相关的错误提示

**关键修改点：**
```vue
<!-- 邮箱输入 -->
<div class="email-section">
  <EmailInput 
    v-model="customerEmail" 
    :disabled="isLoading" 
    :required="true"
    hint="eSIM 激活信息将发送到此邮箱"
    @validation="handleEmailValidation"
  />
</div>
```

```typescript
// 验证邮箱
if (!customerEmail.value || !isEmailValid.value) {
  error.value = '请输入有效的邮箱地址'
  return
}

// 创建订单请求
const orderRequest = {
  productId: Number(props.product.id),
  quantity: quantity.value,
  totalAmount: (props.product.price * quantity.value).toFixed(4),
  customerEmail: customerEmail.value  // 新增
}
```

## 验证规则

### 后端验证
1. **必填验证**：邮箱不能为空
2. **格式验证**：使用正则表达式 `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
3. **错误消息**：返回中文错误提示

### 前端验证
1. **实时验证**：用户输入时实时检查格式
2. **提交验证**：点击购买按钮前再次验证
3. **视觉反馈**：有效邮箱显示绿色勾选图标
4. **错误提示**：显示友好的中文错误消息

## 用户体验优化

1. **清晰的提示**：输入框下方显示"eSIM 激活信息将发送到此邮箱"
2. **即时反馈**：输入有效邮箱后显示绿色勾选图标
3. **错误处理**：邮箱格式错误时显示红色错误提示
4. **触觉反馈**：集成 Telegram 触觉反馈
5. **响应式设计**：适配小屏幕设备（320px+）

## 测试建议

### 后端测试
```bash
# 测试邮箱必填
curl -X POST http://localhost:8080/api/miniapp/esim/orders \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1, "total_amount": "10.00"}'
# 预期：返回 "Customer email is required"

# 测试邮箱格式
curl -X POST http://localhost:8080/api/miniapp/esim/orders \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1, "total_amount": "10.00", "customer_email": "invalid"}'
# 预期：返回 "邮箱格式不正确"

# 测试正常流程
curl -X POST http://localhost:8080/api/miniapp/esim/orders \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1, "total_amount": "10.00", "customer_email": "user@example.com"}'
# 预期：成功创建订单
```

### 前端测试
1. 打开购买弹窗
2. 不填写邮箱，点击"确认支付" → 应显示"请输入邮箱地址"
3. 输入无效邮箱（如 "test"），点击"确认支付" → 应显示"邮箱格式不正确"
4. 输入有效邮箱（如 "test@example.com"） → 应显示绿色勾选图标
5. 完成购买流程 → 订单应成功创建

## 兼容性说明

- ✅ 向后兼容：旧版本前端将无法创建订单（需要更新）
- ✅ 数据库无需修改：邮箱仅用于第三方 API 调用
- ✅ TypeScript 类型安全：所有类型定义已更新
- ✅ 响应式设计：支持所有移动设备尺寸

## 部署注意事项

1. **前后端同步部署**：必须同时部署前端和后端更新
2. **API 版本**：如需保持向后兼容，考虑创建新的 API 版本
3. **用户通知**：建议在更新后通知用户新的邮箱要求
4. **测试环境**：先在测试环境验证完整流程

## 完成状态

- ✅ 后端类型定义更新
- ✅ 后端验证逻辑实现
- ✅ 后端 API 接口更新
- ✅ 前端类型定义更新
- ✅ 前端邮箱输入组件创建
- ✅ 前端购买流程更新
- ✅ 错误处理和用户提示
- ✅ 代码诊断检查通过

所有修改已完成，代码可以直接使用！
