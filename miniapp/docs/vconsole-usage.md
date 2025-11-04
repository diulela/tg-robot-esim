# vConsole 调试控制台使用指南

## 简介

vConsole 是一个轻量级、可扩展的移动端网页调试工具，专为移动端 Web 开发设计。在 Telegram Mini App 中，它可以帮助我们查看控制台日志、网络请求、元素信息等。

## 启用方式

### 1. 开发环境自动启用

在开发环境（`npm run dev`）中，vConsole 会自动启用。

### 2. 生产环境通过 URL 参数启用

在生产环境中，可以通过在 URL 后添加 `?debug=true` 参数来启用：

```
https://your-app.com/?debug=true
```

### 3. 通过环境变量启用

在 `.env.development` 或 `.env.production` 文件中设置：

```env
VITE_ENABLE_VCONSOLE=true
```

## 功能特性

### 1. Console 面板
- 查看 `console.log`、`console.warn`、`console.error` 等日志
- 支持日志过滤和搜索
- 显示日志时间戳和来源

### 2. Network 面板
- 查看所有网络请求（XHR、Fetch）
- 查看请求和响应的详细信息
- 查看请求耗时和状态码

### 3. Element 面板
- 查看 DOM 树结构
- 查看元素的样式和属性
- 实时修改元素（仅用于调试）

### 4. Storage 面板
- 查看 localStorage
- 查看 sessionStorage
- 查看 Cookies

### 5. System 面板
- 查看设备信息
- 查看浏览器信息
- 查看屏幕信息

## 使用技巧

### 1. 快速打开/关闭

点击页面右下角的绿色按钮即可打开/关闭 vConsole 面板。

### 2. 清除日志

在 Console 面板中点击"Clear"按钮可以清除所有日志。

### 3. 过滤日志

在 Console 面板顶部的搜索框中输入关键词可以过滤日志。

### 4. 查看网络请求详情

在 Network 面板中点击任意请求可以查看详细的请求和响应信息。

### 5. 复制日志

长按日志可以复制日志内容。

## 编程接口

### 初始化

```typescript
import { initVConsole } from '@/plugins/vconsole'

// 初始化 vConsole
initVConsole()
```

### 显示/隐藏

```typescript
import { showVConsole, hideVConsole } from '@/plugins/vconsole'

// 显示 vConsole
showVConsole()

// 隐藏 vConsole
hideVConsole()
```

### 销毁

```typescript
import { destroyVConsole } from '@/plugins/vconsole'

// 销毁 vConsole 实例
destroyVConsole()
```

### 获取实例

```typescript
import { getVConsole } from '@/plugins/vconsole'

// 获取 vConsole 实例
const vConsole = getVConsole()
```

## 调试技巧

### 1. 调试 API 请求

在 Network 面板中可以查看所有 API 请求：
- 请求 URL
- 请求方法（GET、POST 等）
- 请求头
- 请求体
- 响应状态码
- 响应数据
- 请求耗时

### 2. 调试状态管理

使用 `console.log` 输出 Pinia store 的状态：

```typescript
import { useProductsStore } from '@/stores/products'

const productsStore = useProductsStore()
console.log('Products:', productsStore.products)
console.log('Loading:', productsStore.isLoading)
```

### 3. 调试路由

在路由守卫中添加日志：

```typescript
router.beforeEach((to, from, next) => {
  console.log('导航到:', to.name, to.path)
  console.log('来自:', from.name, from.path)
  next()
})
```

### 4. 调试 Telegram WebApp

查看 Telegram WebApp 的信息：

```typescript
console.log('Telegram WebApp:', window.Telegram?.WebApp)
console.log('User:', window.Telegram?.WebApp?.initDataUnsafe?.user)
console.log('Theme:', window.Telegram?.WebApp?.themeParams)
```

## 性能影响

vConsole 对性能的影响很小：
- 打包后大小约 100KB（gzip 后约 30KB）
- 只在需要时加载
- 不影响生产环境性能（默认不启用）

## 注意事项

1. **生产环境**: 默认情况下，vConsole 在生产环境中不会启用，除非通过 URL 参数 `?debug=true` 手动启用
2. **敏感信息**: 注意不要在日志中输出敏感信息（如密码、token 等）
3. **性能**: 大量日志可能会影响性能，建议在生产环境中谨慎使用
4. **移除**: 如果不需要 vConsole，可以通过设置环境变量 `VITE_ENABLE_VCONSOLE=false` 来禁用

## 替代方案

如果 vConsole 不满足需求，可以考虑以下替代方案：

1. **Eruda**: 功能更强大的移动端调试工具
2. **Chrome DevTools**: 通过 USB 连接手机使用 Chrome 远程调试
3. **Safari Web Inspector**: iOS 设备可以使用 Safari 的 Web Inspector

## 相关链接

- [vConsole GitHub](https://github.com/Tencent/vConsole)
- [vConsole 文档](https://github.com/Tencent/vConsole/blob/dev/README_CN.md)
