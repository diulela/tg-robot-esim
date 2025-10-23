@echo off
REM Windows 构建脚本

setlocal enabledelayedexpansion

REM 项目信息
set PROJECT_NAME=tg-robot-sim
if "%VERSION%"=="" set VERSION=1.0.0
for /f "tokens=*" %%i in ('powershell -command "Get-Date -UFormat '%%Y-%%m-%%d_%%H:%%M:%%S'"') do set BUILD_TIME=%%i
set GIT_COMMIT=unknown

REM 构建目录
set BUILD_DIR=.\bin
if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo Building %PROJECT_NAME% v%VERSION%...
echo Build time: %BUILD_TIME%
echo Git commit: %GIT_COMMIT%

REM 设置构建标志
set LDFLAGS=-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%

REM 本地构建
echo Building for local platform...
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%PROJECT_NAME%.exe .\cmd\bot\main.go

REM Linux AMD64 交叉编译
echo Cross-compiling for Linux AMD64...
set CGO_ENABLED=1
set GOOS=linux
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -a -o %BUILD_DIR%\%PROJECT_NAME%-linux-amd64 .\cmd\bot\main.go

echo Build completed successfully!
echo Binaries are available in %BUILD_DIR%\
dir %BUILD_DIR%

endlocal