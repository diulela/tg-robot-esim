# eSIM 管理工具 (GM Tool)

管理员运维工具，用于管理 eSIM 产品数据和系统维护。

## 功能

- **产品同步**: 从 eSIM API 同步产品数据到本地数据库
- **产品查询**: 查看本地数据库中的产品信息

## 编译

```bash
# Windows
go build -o gm.exe cmd/gm/main.go

# Linux/Mac
go build -o gm cmd/gm/main.go
```

## 使用方法

### 1. 同步产品数据

从 API 同步所有产品到本地数据库：

```bash
gm -cmd sync-products
```

只同步特定类型的产品：

```bash
# 同步本地产品
gm -cmd sync-products -type local

# 同步区域产品
gm -cmd sync-products -type regional

# 同步全球产品
gm -cmd sync-products -type global
```

限制同步数量（用于测试）：

```bash
# 只同步前 10 个产品
gm -cmd sync-products -limit 10
```

使用自定义配置文件：

```bash
gm -cmd sync-products -config /path/to/config.json
```

### 2. 查看本地产品

列出所有本地产品：

```bash
gm -cmd list-products
```

列出特定类型的产品：

```bash
gm -cmd list-products -type local
```

### 3. 帮助信息

```bash
gm -cmd help
```

## 命令参数

### 全局参数

- `-cmd <command>`: 要执行的命令
- `-config <path>`: 配置文件路径（默认: `config/config.json`）

### sync-products 参数

- `-type <type>`: 产品类型（`local`, `regional`, `global`）
- `-limit <n>`: 限制同步数量（0 表示全部）

### list-products 参数

- `-type <type>`: 产品类型过滤

## 数据库表结构

产品数据存储在 `products` 表中，包含以下字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 主键 |
| third_party_id | string | 第三方产品ID（唯一） |
| name | string | 产品中文名称 |
| name_en | string | 产品英文名称 |
| description | text | 产品描述 |
| type | string | 产品类型 |
| countries | text | 支持的国家（JSON） |
| data_size | int | 流量大小（MB） |
| valid_days | int | 有效天数 |
| features | text | 产品特性（JSON） |
| image | string | 产品图片URL |
| price | decimal | 价格 |
| cost_price | decimal | 成本价 |
| retail_price | decimal | 零售价 |
| agent_price | decimal | 代理价 |
| platform_profit | decimal | 平台利润 |
| is_hot | boolean | 是否热门 |
| is_recommend | boolean | 是否推荐 |
| sort_order | int | 排序 |
| status | string | 状态 |
| synced_at | timestamp | 同步时间 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

## 使用场景

### 场景 1: 首次部署

首次部署时，需要同步所有产品数据：

```bash
# 1. 确保配置文件正确
cat config/config.json

# 2. 同步所有产品
gm -cmd sync-products

# 3. 验证同步结果
gm -cmd list-products
```

### 场景 2: 定期更新

建议每天定期同步产品数据，保持数据最新：

```bash
# 可以添加到 cron 或 Windows 任务计划
0 2 * * * /path/to/gm -cmd sync-products
```

### 场景 3: 测试新产品

测试时只同步少量产品：

```bash
# 同步 5 个产品进行测试
gm -cmd sync-products -limit 5 -type local
```

### 场景 4: 数据验证

同步后验证数据：

```bash
# 查看所有产品
gm -cmd list-products

# 查看特定类型
gm -cmd list-products -type global
```

## 注意事项

1. **API 限流**: 同步时会自动添加延迟（500ms），避免触发 API 限流
2. **数据更新**: 使用 `Upsert` 操作，已存在的产品会被更新
3. **唯一标识**: 使用 `third_party_id` 作为唯一标识，避免重复
4. **事务处理**: 批量操作使用事务，确保数据一致性
5. **错误处理**: 单个产品失败不会影响其他产品的同步

## 故障排查

### 问题 1: 连接数据库失败

```
错误: 初始化数据库失败: failed to connect to database
```

解决方案：
- 检查 `config/config.json` 中的数据库配置
- 确保数据库服务正在运行
- 检查数据库连接权限

### 问题 2: API 认证失败

```
错误: API error 401: Unauthorized
```

解决方案：
- 检查 `config/config.json` 中的 API 密钥配置
- 确保 API Key 和 API Secret 正确
- 检查时区偏移设置

### 问题 3: 同步失败

```
错误: 获取产品列表失败: timeout
```

解决方案：
- 检查网络连接
- 增加超时时间
- 检查 API 服务状态

## 扩展功能

未来可以添加的功能：

- [ ] 产品删除命令
- [ ] 产品更新命令
- [ ] 数据导出功能
- [ ] 数据统计报告
- [ ] 价格批量调整
- [ ] 产品状态管理
- [ ] 同步日志记录
- [ ] 增量同步优化

## 技术支持

如有问题，请联系技术支持或查看项目文档。
