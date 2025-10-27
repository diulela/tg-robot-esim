@echo off
echo Building GM Tool...
echo.

REM 设置输出目录
set OUTPUT_DIR=bin
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

REM 构建 Windows 版本
echo Building for Windows...
go build -o %OUTPUT_DIR%\gm.exe cmd\gm\main.go

if %errorlevel% equ 0 (
    echo.
    echo ✓ Build successful!
    echo.
    echo Output: %OUTPUT_DIR%\gm.exe
    echo.
    echo Usage:
    echo   %OUTPUT_DIR%\gm.exe -cmd sync-products
    echo   %OUTPUT_DIR%\gm.exe -cmd list-products
    echo   %OUTPUT_DIR%\gm.exe -cmd help
    echo.
) else (
    echo.
    echo ✗ Build failed!
    echo.
)

pause
