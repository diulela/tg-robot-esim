#!/bin/bash

# 构建脚本
set -e

# 项目信息
PROJECT_NAME="tg-robot-sim"
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=${GIT_COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}

# 构建目录
BUILD_DIR="./bin"
mkdir -p $BUILD_DIR

echo "Building $PROJECT_NAME v$VERSION..."
echo "Build time: $BUILD_TIME"
echo "Git commit: $GIT_COMMIT"

# 设置构建标志
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

# 本地构建
echo "Building for local platform..."
go build -ldflags "$LDFLAGS" -o $BUILD_DIR/$PROJECT_NAME ./cmd/bot/main.go

# Linux AMD64 交叉编译
echo "Cross-compiling for Linux AMD64..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -a -o $BUILD_DIR/$PROJECT_NAME-linux-amd64 ./cmd/bot/main.go

# Windows AMD64 交叉编译
echo "Cross-compiling for Windows AMD64..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -a -o $BUILD_DIR/$PROJECT_NAME-windows-amd64.exe ./cmd/bot/main.go

echo "Build completed successfully!"
echo "Binaries are available in $BUILD_DIR/"
ls -la $BUILD_DIR/