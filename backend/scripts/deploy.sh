#!/bin/bash

# Mini App 部署脚本

set -e

echo "=== Mini App Deployment Script ==="

# 检查环境变量
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "Error: TELEGRAM_BOT_TOKEN is not set"
    exit 1
fi

# 构建前端
echo "Building frontend..."
cd miniapp
npm install
npm run build
cd ..

# 构建后端
echo "Building backend..."
go build -o miniapp-server ./cmd/miniapp

# 运行数据库迁移
echo "Running database migrations..."
./miniapp-server migrate

# 重启服务
echo "Restarting service..."
if command -v systemctl &> /dev/null; then
    sudo systemctl restart miniapp
else
    echo "Please restart the service manually"
fi

echo "=== Deployment completed ==="
