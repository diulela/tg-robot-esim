---
inclusion: fileMatch
fileMatchPattern: 'miniapp/**'
---

# 前端开发规范 (Telegram Mini App)

## 开发环境和工具

### 技术栈
- **构建工具**: Vite 5.0+ (快速开发和热重载)
- **样式**: 原生 CSS，使用 CSS 变量和现代布局
- **JavaScript**: ES6+ 模块化，无框架依赖
- **Telegram**: Telegram Web App SDK

### 开发服务器
```bash
cd miniapp
npm install
npm run dev  # 启动开发服务器 (localhost:3000)
npm run build  # 构建生产版本
```

## 代码组织规范

### 文件命名约定
- 页面文件: `kebab-case.js` (如 `product-list.js`)
- 组件文件: `PascalCase.js` (如 `ProductCard.js`)
- 工具文件: `camelCase.js` (如 `apiClient.js`)
- 样式文件: `kebab-case.css` (如 `product-list.css`)

### 模块导入导出
```javascript
// 使用命名导出
export const ProductCard = {
  render: (product) => { /* ... */ }
};

// 使用默认导出用于主要功能
export default class ProductService {
  // ...
}

// 导入时保持一致性
import ProductService from './services/ProductService.js';
import { ProductCard, PriceFormatter } from './components/index.js';
```

## UI/UX 设计规范

### Telegram 主题适配
- 使用 Telegram Web App 提供的 CSS 变量
- 支持明暗主题自动切换
- 遵循 Telegram 的设计语言

### 响应式设计
- 移动优先设计 (Mobile First)
- 支持不同屏幕尺寸 (320px - 480px)
- 使用 Flexbox 和 CSS Grid 布局
- 触摸友好的交互元素 (最小 44px 点击区域)

### CSS 变量系统
```css
:root {
  /* Telegram 主题变量 */
  --tg-theme-bg-color: var(--tg-theme-bg-color, #ffffff);
  --tg-theme-text-color: var(--tg-theme-text-color, #000000);
  
  /* 自定义设计变量 */
  --primary-color: #0088cc;
  --success-color: #4caf50;
  --error-color: #f44336;
  --border-radius: 8px;
  --spacing-unit: 8px;
}
```

## JavaScript 开发规范

### 代码风格
- 使用 ES6+ 语法 (const/let, 箭头函数, 模板字符串)
- 函数和变量使用 camelCase
- 常量使用 UPPER_SNAKE_CASE
- 类名使用 PascalCase

### 错误处理
```javascript
// API 调用错误处理
async function fetchProducts() {
  try {
    const response = await apiClient.get('/api/products');
    return response.data;
  } catch (error) {
    console.error('获取产品列表失败:', error);
    showErrorMessage('加载产品失败，请稍后重试');
    return [];
  }
}

// 用户友好的错误提示
function showErrorMessage(message) {
  // 显示中文错误信息
  const toast = document.createElement('div');
  toast.className = 'error-toast';
  toast.textContent = message;
  document.body.appendChild(toast);
}
```

### 状态管理
- 使用 localStorage 持久化用户数据
- 实现简单的状态管理器 (state.js)
- 避免全局变量，使用模块化状态

```javascript
// 状态管理示例
const AppState = {
  user: null,
  cart: [],
  
  setUser(userData) {
    this.user = userData;
    localStorage.setItem('user', JSON.stringify(userData));
  },
  
  addToCart(product) {
    this.cart.push(product);
    this.saveCart();
  }
};
```

## Telegram Web App 集成

### 初始化和配置
```javascript
// telegram.js - Telegram Web App 集成
export const TelegramWebApp = {
  init() {
    if (window.Telegram?.WebApp) {
      const tg = window.Telegram.WebApp;
      tg.ready();
      tg.expand();
      return tg;
    }
    // 开发环境模拟
    return this.createMockTelegram();
  },
  
  // 开发环境模拟对象
  createMockTelegram() {
    return {
      initDataUnsafe: { user: { id: 123, first_name: '测试用户' } },
      MainButton: { show: () => {}, hide: () => {} },
      close: () => console.log('关闭 Mini App')
    };
  }
};
```

### 用户数据处理
- 验证 Telegram 用户数据的完整性
- 处理用户授权和身份验证
- 安全地传输用户信息到后端

### 主按钮 (MainButton) 使用
```javascript
// 动态控制 Telegram 主按钮
function updateMainButton(text, callback) {
  const tg = window.Telegram.WebApp;
  tg.MainButton.text = text;
  tg.MainButton.show();
  tg.MainButton.onClick(callback);
}

// 示例：购买按钮
updateMainButton('立即购买 ¥99', () => {
  // 处理购买逻辑
  processPurchase();
});
```

## API 集成规范

### HTTP 客户端
```javascript
// apiClient.js - 统一的 API 客户端
class ApiClient {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }
  
  async request(method, endpoint, data = null) {
    const url = `${this.baseURL}${endpoint}`;
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json',
        'X-Telegram-Init-Data': this.getTelegramInitData()
      }
    };
    
    if (data) {
      options.body = JSON.stringify(data);
    }
    
    const response = await fetch(url, options);
    
    if (!response.ok) {
      throw new Error(`API 请求失败: ${response.status}`);
    }
    
    return response.json();
  }
}
```

### 数据验证
- 客户端输入验证 (validator.js)
- 服务端数据验证
- 用户友好的验证错误提示

## 性能优化

### 资源加载
- 图片懒加载
- 代码分割和按需加载
- 静态资源缓存策略

### 用户体验
- 加载状态指示器
- 骨架屏 (Skeleton Loading)
- 平滑的页面转场动画
- 离线状态处理

### 构建优化
```javascript
// vite.config.js - 构建优化配置
export default {
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['qrcode'],
          utils: ['src/utils/formatter.js', 'src/utils/validator.js']
        }
      }
    },
    minify: 'terser',
    sourcemap: false
  }
};
```

## 测试和调试

### 开发调试
- 使用浏览器开发者工具
- Telegram Web App 调试模式
- 移动设备调试 (Chrome DevTools)

### 测试策略
- 功能测试：核心用户流程
- 兼容性测试：不同 Telegram 客户端
- 性能测试：加载速度和响应时间

### 错误监控
```javascript
// 全局错误处理
window.addEventListener('error', (event) => {
  console.error('JavaScript 错误:', event.error);
  // 发送错误报告到后端
  reportError(event.error);
});

window.addEventListener('unhandledrejection', (event) => {
  console.error('未处理的 Promise 拒绝:', event.reason);
  reportError(event.reason);
});
```