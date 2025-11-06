# eSIM 订单 API 测试指南

## API 端点

### 1. 创建 eSIM 订单

```http
POST /api/miniapp/esim/orders
Content-Type: application/json
X-Telegram-Init-Data: user_id=123456

{
    "product_id": 1,
    "quantity": 1,
    "total_amount": "27.7200",
    "remark": "测试订单"
}
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "order_id": 1,
        "order_no": "ORD1730800000001",
        "status": "processing",
        "total_amount": "27.7200",
        "created_at": "2024-11-05T10:00:00Z"
    }
}
```

### 2. 获取用户 eSIM 订单列表

```http
GET /api/miniapp/esim/orders?limit=20&offset=0&status=processing
X-Telegram-Init-Data: user_id=123456
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "orders": [
            {
                "order_id": 1,
                "order_no": "ORD1730800000001",
                "product_id": 1,
                "product_name": "中国 无限流量/15天",
                "quantity": 1,
                "unit_price": "27.7200",
                "total_amount": "27.7200",
                "status": "processing",
                "provider_order_id": "12345",
                "created_at": "2024-11-05T10:00:00Z",
                "updated_at": "2024-11-05T10:00:00Z",
                "completed_at": null
            }
        ],
        "total": 1,
        "limit": 20,
        "offset": 0
    }
}
```

### 3. 获取 eSIM 订单详情

```http
GET /api/miniapp/esim/orders/1
X-Telegram-Init-Data: user_id=123456
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "order_id": 1,
        "order_no": "ORD1730800000001",
        "user_id": 123456,
        "product_id": 1,
        "product_name": "中国 无限流量/15天",
        "quantity": 1,
        "unit_price": "27.7200",
        "total_amount": "27.7200",
        "status": "completed",
        "provider_order_id": "12345",
        "order_items": [
            {
                "id": 1,
                "product_id": 1,
                "product_name": "中国 无限流量/15天",
                "quantity": 1,
                "unit_price": 27.72,
                "subtotal": 27.72,
                "data_size": 0,
                "valid_days": 15
            }
        ],
        "esims": [
            {
                "id": 1,
                "iccid": "8901234567890123456",
                "status": "active",
                "has_activation_code": true,
                "has_qr_code": true
            }
        ],
        "created_at": "2024-11-05T10:00:00Z",
        "updated_at": "2024-11-05T10:01:00Z",
        "completed_at": "2024-11-05T10:01:00Z"
    }
}
```

## 错误响应

### 余额不足
```json
{
    "code": 40001,
    "message": "余额不足，请先充值",
    "details": ""
}
```

### 订单不存在
```json
{
    "code": 40003,
    "message": "订单不存在",
    "details": ""
}
```

### 未授权访问
```json
{
    "code": 401,
    "message": "Unauthorized",
    "details": "Invalid user ID"
}
```

## 测试流程

1. **准备测试数据**
   - 确保数据库中有产品数据
   - 确保用户有足够的余额

2. **创建订单**
   - 调用创建订单 API
   - 验证订单状态为 "processing"
   - 验证余额被冻结

3. **查询订单**
   - 调用订单列表 API
   - 调用订单详情 API
   - 验证返回的订单信息

4. **模拟订单完成**
   - 手动调用订单同步服务
   - 验证订单状态变为 "completed"
   - 验证余额被正确扣除
   - 验证订单详情包含 eSIM 信息

## 注意事项

1. **用户认证**: 在生产环境中需要验证 Telegram Web App 数据
2. **余额检查**: 确保用户有足够的可用余额
3. **产品验证**: 确保产品存在且可用
4. **并发安全**: 系统使用数据库事务确保操作原子性
5. **错误处理**: 所有错误都有相应的错误码和中文提示