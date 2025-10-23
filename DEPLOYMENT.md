# 部署指南

本文档描述了如何部署 Telegram 机器人系统。

## 环境要求

### 系统要求
- Linux/Windows/macOS
- Go 1.23.5+ (开发环境)
- SQLite 或 MySQL 数据库
- 网络连接 (访问 Telegram API 和 TRON 网络)

### 必需的配置
- Telegram Bot Token (通过 @BotFather 获取)
- TRON API Key (可选，用于区块链功能)

## 部署方式

### 1. 二进制部署

#### 构建
```bash
# Linux/macOS
chmod +x scripts/build.sh
./scripts/build.sh

# Windows
scripts\build.bat
```

#### 配置
```bash
# 复制配置文件
cp config/config.example.json config/config.json

# 编辑配置文件
vim config/config.json
```

#### 运行
```bash
# 设置环境变量
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
export TRON_API_KEY="your_tron_api_key_here"

# 运行
./bin/tg-robot-sim
```

### 2. Docker 部署

#### 使用 Docker Compose (推荐)
```bash
# 创建环境变量文件
cat > .env << EOF
TELEGRAM_BOT_TOKEN=your_bot_token_here
TRON_API_KEY=your_tron_api_key_here
MYSQL_ROOT_PASSWORD=secure_root_password
MYSQL_DATABASE=telegram_bot
MYSQL_USER=botuser
MYSQL_PASSWORD=secure_bot_password
EOF

# 启动服务
docker-compose -f docker/docker-compose.yml up -d

# 查看日志
docker-compose -f docker/docker-compose.yml logs -f telegram-bot
```

#### 单独使用 Docker
```bash
# 构建镜像
docker build -f docker/Dockerfile -t tg-robot-sim .

# 运行容器
docker run -d \
  --name tg-robot-sim \
  -e TELEGRAM_BOT_TOKEN="your_bot_token_here" \
  -e TRON_API_KEY="your_tron_api_key_here" \
  -v $(pwd)/data:/app/data \
  tg-robot-sim
```

### 3. 系统服务部署

#### Linux Systemd
创建服务文件 `/etc/systemd/system/tg-robot-sim.service`:

```ini
[Unit]
Description=Telegram Robot Simulator
After=network.target

[Service]
Type=simple
User=telegram-bot
Group=telegram-bot
WorkingDirectory=/opt/tg-robot-sim
ExecStart=/opt/tg-robot-sim/bin/tg-robot-sim
Restart=always
RestartSec=10
Environment=TELEGRAM_BOT_TOKEN=your_bot_token_here
Environment=TRON_API_KEY=your_tron_api_key_here

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable tg-robot-sim
sudo systemctl start tg-robot-sim
sudo systemctl status tg-robot-sim
```

## 配置说明

### 环境变量
- `TELEGRAM_BOT_TOKEN`: Telegram Bot 令牌 (必需)
- `TRON_API_KEY`: TRON 网络 API 密钥 (可选)
- `DATABASE_URL`: 数据库连接字符串
- `DEBUG`: 调试模式 (true/false)
- `LOG_LEVEL`: 日志级别 (debug/info/warn/error)

### 配置文件
主要配置在 `config/config.json` 中，包括:
- Telegram 配置
- 数据库配置
- 区块链配置
- 日志配置

## 监控和维护

### 健康检查
系统提供内置的健康检查功能，可以通过以下方式监控:

```bash
# 检查进程状态
ps aux | grep tg-robot-sim

# 查看日志
tail -f bot.log

# Docker 环境
docker-compose logs -f telegram-bot
```

### 日志管理
- 日志文件位置: `bot.log`
- 日志级别可通过配置调整
- 支持日志轮转 (通过配置)

### 数据备份
```bash
# SQLite 备份
cp bot.db bot.db.backup.$(date +%Y%m%d_%H%M%S)

# MySQL 备份
mysqldump -u botuser -p telegram_bot > backup_$(date +%Y%m%d_%H%M%S).sql
```

## 故障排除

### 常见问题

1. **Bot Token 无效**
   - 检查环境变量设置
   - 确认 Token 格式正确
   - 通过 @BotFather 验证 Token

2. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接字符串
   - 确认数据库权限

3. **区块链 API 错误**
   - 检查网络连接
   - 验证 API Key
   - 查看 API 限制

### 日志分析
```bash
# 查看错误日志
grep "ERROR" bot.log

# 查看最近的日志
tail -n 100 bot.log

# 实时监控日志
tail -f bot.log | grep -E "(ERROR|WARN)"
```

## 性能优化

### 系统资源
- 推荐内存: 512MB+
- 推荐 CPU: 1 核心+
- 磁盘空间: 1GB+ (用于日志和数据库)

### 配置优化
- 调整数据库连接池大小
- 设置合适的日志级别
- 配置会话超时时间

## 安全建议

1. **环境变量安全**
   - 不要在代码中硬编码敏感信息
   - 使用环境变量或安全的配置管理

2. **网络安全**
   - 使用 HTTPS/TLS 连接
   - 限制网络访问权限
   - 定期更新依赖

3. **数据安全**
   - 定期备份数据
   - 加密敏感数据
   - 限制数据库访问权限