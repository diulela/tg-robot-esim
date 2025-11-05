---
inclusion: fileMatch
fileMatchPattern: 'miniapp/**'
---

# 前端开发规范 (Vue 3.0 + TypeScript Telegram Mini App)

## 技术栈
- **框架**: Vue 3.0 + Composition API + `<script setup>`
- **语言**: TypeScript 5.3+ (类型安全)
- **UI库**: Vuetify 3.4+ (Material Design)
- **状态管理**: Pinia 2.1+ (现代状态管理)
- **路由**: Vue Router 4.2+
- **构建工具**: Vite 5.0+
- **HTTP客户端**: Axios 1.6+
- **开发工具**: ESLint, Prettier, VConsole

### 开发命令
```bash
npm run dev        # 开发服务器 (localhost:8082)
npm run build      # 构建生产版本
npm run type-check # TypeScript 检查
npm run lint       # 代码检查
```

## 代码组织规范

### 文件命名
- **Vue 组件**: `PascalCase.vue` (如 `ProductCard.vue`)
- **页面组件**: `PascalCasePage.vue` (如 `ProductDetailPage.vue`)
- **TypeScript 文件**: `camelCase.ts` (如 `apiClient.ts`)
- **样式文件**: `kebab-case.scss` (如 `global.scss`)

### Vue 3 组件结构
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// Props 定义
interface Props {
  title: string
  count?: number
}

const props = withDefaults(defineProps<Props>(), {
  count: 0
})

// Emits 定义
const emit = defineEmits<{
  update: [value: string]
  click: [event: MouseEvent]
}>()

// 响应式状态
const isLoading = ref(false)
const displayTitle = computed(() => `${props.title} (${props.count})`)

onMounted(() => {
  // 初始化逻辑
})
</script>

<style scoped lang="scss">
// 组件样式
</style>
```

## UI/UX 设计规范

### Vuetify Material Design
- 使用 Vuetify 3.4+ 组件库，遵循 Material Design 规范
- 支持明暗主题自动切换，集成 Telegram 主题
- 移动优先设计 (320px - 480px)，最小 44px 点击区域
- 虚拟滚动、懒加载、代码分割优化

### 主题配置
```typescript
// plugins/vuetify.ts
const vuetify = createVuetify({
  theme: {
    defaultTheme: telegramService.getColorScheme(),
    themes: {
      light: {
        colors: {
          primary: '#6366F1',
          secondary: '#EC4899'
        }
      },
      dark: {
        colors: {
          primary: '#818CF8',
          secondary: '#F472B6'
        }
      }
    }
  }
})
```

## TypeScript 开发规范

### 代码风格
- **严格模式**: 启用 TypeScript 严格模式
- **命名规范**: 变量/函数 camelCase，类型/接口 PascalCase
- **类型导入**: 使用 `import type` 导入类型定义

### 错误处理
```typescript
async function fetchProducts(): Promise<Product[]> {
  try {
    const response = await productApi.getProducts()
    return response.data
  } catch (error) {
    console.error('获取产品列表失败:', error)
    const appStore = useAppStore()
    appStore.showError('加载产品失败，请稍后重试')
    return []
  }
}
```

### Pinia 状态管理
```typescript
// stores/user.ts
export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  
  const isAuthenticated = computed(() => !!user.value)
  const displayName = computed(() => user.value?.firstName || '未知用户')
  
  const setUser = (userData: User) => {
    user.value = userData
    localStorage.setItem('user', JSON.stringify(userData))
  }
  
  const fetchUserProfile = async (): Promise<void> => {
    isLoading.value = true
    try {
      const profile = await userApi.getProfile()
      setUser(profile)
    } finally {
      isLoading.value = false
    }
  }
  
  return {
    user: readonly(user),
    isLoading: readonly(isLoading),
    isAuthenticated,
    displayName,
    setUser,
    fetchUserProfile
  }
})
```

## Telegram Web App 集成

### 服务层封装
```typescript
// services/telegram.ts
class TelegramService {
  private tg: any
  
  constructor() {
    this.tg = this.initializeTelegram()
  }
  
  private initializeTelegram() {
    if (window.Telegram?.WebApp) {
      const tg = window.Telegram.WebApp
      tg.ready()
      tg.expand()
      return tg
    }
    return this.createMockTelegram() // 开发环境模拟
  }
  
  getUser(): TelegramUser | null {
    return this.tg.initDataUnsafe?.user || null
  }
  
  getColorScheme(): 'light' | 'dark' {
    return this.tg.colorScheme || 'light'
  }
  
  impactFeedback(style: 'light' | 'medium' | 'heavy' = 'medium') {
    this.tg.HapticFeedback?.impactOccurred(style)
  }
  
  showBackButton() {
    this.tg.BackButton?.show()
  }
  
  hideBackButton() {
    this.tg.BackButton?.hide()
  }
}

export const telegramService = new TelegramService()
```

### 主按钮控制
```typescript
// 组件中使用
<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

const handlePurchase = () => {
  telegramService.impactFeedback('medium')
}

onMounted(() => {
  const tg = window.Telegram?.WebApp
  if (tg?.MainButton) {
    tg.MainButton.text = '立即购买 ¥99'
    tg.MainButton.show()
    tg.MainButton.onClick(handlePurchase)
  }
})

onUnmounted(() => {
  const tg = window.Telegram?.WebApp
  tg?.MainButton?.hide()
})
</script>
```

## API 集成规范

### API 客户端
```typescript
// services/api/client.ts
export class ApiClient {
  private instance: AxiosInstance
  
  constructor(config: ApiClientConfig) {
    this.instance = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout,
      headers: { 'Content-Type': 'application/json' }
    })
    
    this.setupInterceptors(config.enableAuth)
  }
  
  private setupInterceptors(enableAuth: boolean) {
    this.instance.interceptors.request.use((config) => {
      if (enableAuth) {
        const authHeader = TelegramAuthService.generateAuthHeader()
        if (authHeader) {
          config.headers['X-Telegram-Init-Data'] = authHeader
        }
      }
      return config
    })
    
    this.instance.interceptors.response.use(
      (response) => response,
      (error) => Promise.reject(this.handleApiError(error))
    )
  }
  
  async get<T>(url: string): Promise<T> {
    const response = await this.instance.get(url)
    return response.data
  }
  
  async post<T>(url: string, data?: any): Promise<T> {
    const response = await this.instance.post(url, data)
    return response.data
  }
}

export const apiClient = new ApiClient({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000,
  enableAuth: true
})
```

### 专用 API 服务
```typescript
// services/api/product.ts
export class ProductApi {
  async getProducts(filters?: ProductFilters): Promise<ProductListResponse> {
    const params = new URLSearchParams()
    if (filters?.countryCode) params.append('country', filters.countryCode)
    return apiClient.get(`/products?${params.toString()}`)
  }
  
  async getProduct(id: string): Promise<Product> {
    return apiClient.get(`/products/${id}`)
  }
}

export const productApi = new ProductApi()
```

## 性能优化

### Vue 3 优化
```typescript
// 组件懒加载
const ProductDetailPage = defineAsyncComponent(() => 
  import('@/views/ProductDetailPage.vue')
)

// v-memo 优化列表渲染
<template>
  <div v-for="product in products" :key="product.id" v-memo="[product.id, product.price]">
    <ProductCard :product="product" />
  </div>
</template>
```

### 用户体验优化
```vue
<template>
  <!-- 骨架屏 -->
  <v-skeleton-loader v-if="isLoading" type="card" />
  
  <!-- 无限滚动 -->
  <v-infinite-scroll @load="loadMore" :loading="isLoadingMore">
    <ProductCard v-for="product in products" :key="product.id" :product="product" />
  </v-infinite-scroll>
</template>
```

### Vite 构建优化
```typescript
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          vuetify: ['vuetify'],
          utils: ['axios', 'dayjs', 'qrcode']
        }
      }
    }
  }
})
```

## 测试和调试

### 开发调试
```typescript
// VConsole 移动端调试 (plugins/vconsole.ts)
export function initVConsole() {
  const isDev = import.meta.env.DEV
  const hasDebugParam = new URLSearchParams(window.location.search).get('debug') === 'true'
  
  if (isDev || hasDebugParam) {
    new VConsole({ theme: 'dark' })
  }
}
```

### 错误监控
```typescript
// 全局错误处理
window.addEventListener('error', (event) => {
  console.error('[Global] JavaScript 错误:', event.error)
  reportError({
    type: 'javascript_error',
    message: event.message,
    stack: event.error?.stack
  })
})

window.addEventListener('unhandledrejection', (event) => {
  console.error('[Global] Promise 拒绝:', event.reason)
  reportError({
    type: 'promise_rejection',
    reason: event.reason
  })
})

async function reportError(errorData: any) {
  try {
    await fetch('/api/errors', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(errorData)
    })
  } catch (error) {
    console.error('发送错误报告失败:', error)
  }
}
```

### 组件测试
```typescript
// 组件单元测试
import { mount } from '@vue/test-utils'
import ProductCard from '@/components/ProductCard.vue'

describe('ProductCard', () => {
  it('正确渲染产品信息', () => {
    const wrapper = mount(ProductCard, {
      props: { product: { id: '1', name: '测试产品', price: 99.99 } }
    })
    
    expect(wrapper.text()).toContain('测试产品')
    expect(wrapper.text()).toContain('¥99.99')
  })
})
```