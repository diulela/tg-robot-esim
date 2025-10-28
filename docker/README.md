# Docker 部署指南

## 文件说明

### 核心文件
- `Dockerfile` - Telegram Bot 服务镜像构建文件
- `Dockerfile.miniapp` - Mini App 服务镜像构建文件（包含前后端）
- `docker-compose.yml` - 完整服务栈（Bot + MySQL + Redis）
- `docker-compose.miniapp.yml` - Mini App 服务栈（HTTP API + Nginx）
- `nginx.conf` - Nginx 反向代理配置

## 部署方式

### 1. 开发环境 - 仅 Bot 服务
```bash
# 构建并启动 Bot 服务
docker-compose up -d

# 查看日志
docker-compose logs -f telegram-bot

# 停止服务
docker-compose down
```

### 2. 生产环境 - Mini App 服务
```bash
# 构建前端资源
cd ../miniapp
npm run build

# 启动 Mini App 服务
cd ../docker
docker-compose -f docker-compose.miniapp.yml up -d

# 查看服务状态
docker-compose -f docker-compose.miniapp.yml ps
```

### 3. 完整部署 - 所有服务
```bash
# 构建前端
cd ../miniapp && npm run build && cd ../docker

# 启动所有服务
docker-compose up -d
docker-compose -f docker-compose.miniapp.yml up -d
```

## 环境变量配置

创建 `.env` 文件：
```bash
# Telegram 配置
TELEGRAM_BOT_TOKEN=your_bot_token_here

# 数据库配置
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=telegram_bot
MYSQL_USER=botuser
MYSQL_PASSWORD=your_bot_password

# 区块链配置
TRON_API_KEY=your_tron_api_key

# 应用配置
LOG_LEVEL=info
DATABASE_URL=mysql://botuser:your_bot_password@mysql:3306/telegram_bot
```

## 服务端口

### 开发环境
- MySQL: `3306`
- Redis: `6379`
- Bot 健康检查: `8080`

### 生产环境
- HTTP API: `8080`
- Nginx: `80` (HTTP), `443` (HTTPS)

## 数据持久化

### 数据卷
- `bot_data` - Bot 应用数据
- `bot_logs` - 应用日志
- `mysql_data` - MySQL 数据库
- `redis_data` - Redis 数据

### 备份数据
```bash
# 备份 MySQL
docker exec tg-robot-mysql mysqldump -u root -p telegram_bot > backup.sql

# 备份应用数据
docker cp tg-robot-sim:/app/data ./backup/data
```

## 健康检查

### Bot 服务
```bash
# 检查 Bot 进程
docker exec tg-robot-sim pgrep main

# 查看健康状态
docker inspect tg-robot-sim --format='{{.State.Health.Status}}'
```

### Mini App 服务
```bash
# 检查 API 健康
curl http://localhost:8080/health

# 检查 Nginx 状态
curl http://localhost/health
```

## 故障排除

### 常见问题

1. **Bot 无法启动**
   - 检查 `TELEGRAM_BOT_TOKEN` 是否正确
   - 确认网络连接正常
   - 查看容器日志：`docker logs tg-robot-sim`

2. **数据库连接失败**
   - 检查 MySQL 容器是否正常运行
   - 验证数据库连接字符串
   - 确认用户权限设置

3. **前端资源加载失败**
   - 确认前端已正确构建：`npm run build`
   - 检查 Nginx 配置和静态文件路径
   - 验证文件权限设置

### 日志查看
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs telegram-bot
docker-compose logs mysql

# 实时跟踪日志
docker-compose logs -f telegram-bot
```

## 性能优化

### 资源限制
在 `docker-compose.yml` 中添加资源限制：
```yaml
services:
  telegram-bot:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
```

### 缓存优化
- 使用多阶段构建减少镜像大小
- 合理配置 `.dockerignore` 文件
- 启用 Docker BuildKit 加速构建

## 安全配置

### 生产环境建议
1. 使用非 root 用户运行容器
2. 配置 HTTPS 证书
3. 设置防火墙规则
4. 定期更新基础镜像
5. 使用 Docker secrets 管理敏感信息

### SSL 证书配置
将证书文件放置在 `./ssl/` 目录：
- `cert.pem` - SSL 证书
- `key.pem` - 私钥文件

然后取消注释 `nginx.conf` 中的 HTTPS 配置。