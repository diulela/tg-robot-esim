// 全局类型声明和导出

// 导出所有类型
export * from './api'
export * from './telegram'
export * from './components'
export * from './esim-order'

// 通用工具类型
export type Nullable<T> = T | null
export type Optional<T> = T | undefined
export type Maybe<T> = T | null | undefined

// 深度可选类型
export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P]
}

// 深度只读类型
export type DeepReadonly<T> = {
  readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P]
}

// 选择性必需类型
export type RequiredKeys<T, K extends keyof T> = T & Required<Pick<T, K>>

// 选择性可选类型
export type OptionalKeys<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>

// 值类型提取
export type ValueOf<T> = T[keyof T]

// 数组元素类型提取
export type ArrayElement<T> = T extends (infer U)[] ? U : never

// Promise 结果类型提取
export type PromiseResult<T> = T extends Promise<infer U> ? U : T

// 函数参数类型提取
export type FunctionArgs<T> = T extends (...args: infer U) => any ? U : never

// 函数返回类型提取
export type FunctionReturn<T> = T extends (...args: any[]) => infer U ? U : never

// 键值对类型
export interface KeyValuePair<K = string, V = any> {
  key: K
  value: V
}

// 选项类型
export interface SelectOption<T = any> {
  label: string
  value: T
  disabled?: boolean
  icon?: string
  description?: string
}

// 菜单项类型
export interface MenuItem {
  id: string
  label: string
  icon?: string
  route?: string
  action?: () => void
  children?: MenuItem[]
  disabled?: boolean
  divider?: boolean
}

// 面包屑项类型
export interface BreadcrumbItem {
  title: string
  to?: string
  disabled?: boolean
  exact?: boolean
}

// 分页信息类型
export interface PaginationInfo {
  page: number
  pageSize: number
  total: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}

// 排序信息类型
export interface SortInfo {
  key: string
  order: 'asc' | 'desc'
}

// 筛选信息类型
export interface FilterInfo {
  key: string
  value: any
  operator?: 'eq' | 'ne' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'nin' | 'like'
}

// 搜索参数类型
export interface SearchParams {
  query?: string
  filters?: FilterInfo[]
  sort?: SortInfo
  pagination?: Pick<PaginationInfo, 'page' | 'pageSize'>
}

// 表格列定义类型
export interface TableColumn<T = any> {
  key: keyof T | string
  title: string
  sortable?: boolean
  filterable?: boolean
  width?: string | number
  minWidth?: string | number
  maxWidth?: string | number
  align?: 'left' | 'center' | 'right'
  fixed?: 'left' | 'right'
  render?: (value: any, record: T, index: number) => any
  formatter?: (value: any) => string
}

// 表单字段类型
export interface FormField {
  name: string
  label: string
  type: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search' | 'textarea' | 'select' | 'checkbox' | 'radio' | 'date' | 'time' | 'datetime'
  placeholder?: string
  required?: boolean
  disabled?: boolean
  readonly?: boolean
  options?: SelectOption[]
  validation?: {
    required?: boolean
    min?: number
    max?: number
    minLength?: number
    maxLength?: number
    pattern?: RegExp
    custom?: (value: any) => boolean | string
  }
  defaultValue?: any
  description?: string
}

// 表单验证结果类型
export interface ValidationResult {
  valid: boolean
  errors: Record<string, string[]>
}

// 通知类型
export interface Notification {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title?: string
  message: string
  duration?: number
  persistent?: boolean
  actions?: Array<{
    label: string
    action: () => void
    color?: string
  }>
  createdAt: Date
}

// 应用配置类型
export interface AppConfig {
  name: string
  version: string
  apiBaseUrl: string
  telegramBotToken?: string
  telegramWebAppUrl?: string
  defaultLanguage: string
  supportedLanguages: string[]
  theme: {
    defaultMode: 'light' | 'dark' | 'auto'
    primaryColor: string
    secondaryColor: string
  }
  features: {
    enableNotifications: boolean
    enableAnalytics: boolean
    enableErrorReporting: boolean
    enableOfflineMode: boolean
  }
}

// 环境变量类型
export interface EnvironmentVariables {
  NODE_ENV: 'development' | 'production' | 'test'
  VITE_API_BASE_URL: string
  VITE_TELEGRAM_BOT_TOKEN?: string
  VITE_TELEGRAM_WEB_APP_URL?: string
  VITE_ENABLE_MOCK?: string
  VITE_ENABLE_DEBUG?: string
}

// 错误边界类型
export interface ErrorBoundaryState {
  hasError: boolean
  error?: Error
  errorInfo?: {
    componentStack: string
  }
}

// 加载状态类型
export interface LoadingState {
  loading: boolean
  error?: string | Error
  data?: any
}

// 异步操作状态类型
export type AsyncState<T = any> = {
  loading: boolean
  data?: T
  error?: string | Error
  lastUpdated?: Date
}

// 缓存项类型
export interface CacheItem<T = any> {
  key: string
  value: T
  expiresAt?: Date
  createdAt: Date
  accessCount: number
  lastAccessed: Date
}

// 存储适配器接口
export interface StorageAdapter {
  getItem(key: string): string | null
  setItem(key: string, value: string): void
  removeItem(key: string): void
  clear(): void
  key(index: number): string | null
  readonly length: number
}

// 事件总线类型
export interface EventBus {
  on<T = any>(event: string, handler: (data: T) => void): void
  off<T = any>(event: string, handler: (data: T) => void): void
  emit<T = any>(event: string, data?: T): void
  once<T = any>(event: string, handler: (data: T) => void): void
  clear(): void
}

// 插件接口
export interface Plugin {
  name: string
  version: string
  install: (app: any, options?: any) => void
  uninstall?: (app: any) => void
}

// 中间件类型
export type Middleware<T = any> = (context: T, next: () => void) => void | Promise<void>

// 守卫类型
export type Guard<T = any> = (context: T) => boolean | Promise<boolean>

// 钩子函数类型
export type Hook<T = any> = (context: T) => void | Promise<void>

// 工厂函数类型
export type Factory<T = any> = (...args: any[]) => T

// 单例类型
export type Singleton<T = any> = {
  getInstance: () => T
}

// 观察者模式类型
export interface Observer<T = any> {
  update(data: T): void
}

export interface Observable<T = any> {
  subscribe(observer: Observer<T>): void
  unsubscribe(observer: Observer<T>): void
  notify(data: T): void
}

// 策略模式类型
export interface Strategy<T = any, R = any> {
  execute(context: T): R
}

// 命令模式类型
export interface Command {
  execute(): void | Promise<void>
  undo?(): void | Promise<void>
  canExecute?(): boolean
}

// 状态机类型
export interface StateMachine<S = string, E = string> {
  currentState: S
  transition(event: E): void
  canTransition(event: E): boolean
  getAvailableTransitions(): E[]
}