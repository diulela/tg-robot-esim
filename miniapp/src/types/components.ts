// Vue 组件 Props 和 Emit 类型定义

import type { Order, Product, Region, Country, ESIMInfo, OrderStatus } from './api'

// 基础组件 Props
export interface BaseComponentProps {
  loading?: boolean
  disabled?: boolean
  variant?: 'default' | 'outlined' | 'text' | 'elevated' | 'flat' | 'tonal'
  size?: 'x-small' | 'small' | 'default' | 'large' | 'x-large'
  color?: string
  rounded?: boolean | string | number
}

// 布局组件 Props
export interface AppLayoutProps {
  title?: string
  showBackButton?: boolean
  showUserButton?: boolean
  showBottomNav?: boolean
  currentTab?: string
}

export interface MobileContainerProps {
  fluid?: boolean
  maxWidth?: string | number
  padding?: string | number
  class?: string
}

// 业务组件 Props
export interface OrderCardProps {
  order: Order
  showActions?: boolean
  compact?: boolean
}

export interface OrderCardEmits {
  click: [order: Order]
  viewDetails: [orderId: string]
  copyOrderNumber: [orderNumber: string]
}

export interface ProductCardProps {
  product: Product
  showPrice?: boolean
  showDescription?: boolean
  compact?: boolean
}

export interface ProductCardEmits {
  click: [product: Product]
  buy: [product: Product]
  addToCart: [productId: string]
  viewDetails: [productId: string]
}

// 热门商品页面组件 Props
export interface PageHeaderProps {
  title: string
  showBack: boolean
}

export interface PageHeaderEmits {
  back: []
}

export interface ProductListProps {
  products: Product[]
  loading: boolean
  error: string | null
}

export interface ProductListEmits {
  buyProduct: [product: Product]
  productClick: [product: Product]
}

export interface LoadMoreButtonProps {
  loading: boolean
  disabled?: boolean
}

export interface LoadMoreButtonEmits {
  loadMore: []
}

export interface RegionGridProps {
  regions: Region[]
  columns?: number
  loading?: boolean
}

export interface RegionGridEmits {
  selectRegion: [region: Region]
}

export interface CountryListProps {
  countries: Country[]
  searchable?: boolean
  showIndex?: boolean
  loading?: boolean
}

export interface CountryListEmits {
  selectCountry: [country: Country]
  search: [query: string]
}

export interface CountryItemProps {
  country: Country
  showFlag?: boolean
  showArrow?: boolean
}

export interface CountryItemEmits {
  click: [country: Country]
}

// eSIM 相关组件 Props
export interface ESIMInfoProps {
  esimInfo: ESIMInfo
  showQRCode?: boolean
  allowCopy?: boolean
}

export interface ESIMInfoEmits {
  copyIccid: [iccid: string]
  copyActivationCode: [code: string]
  downloadQR: [qrCode: string]
  shareQR: [qrCode: string]
}

export interface QRCodeDisplayProps {
  value: string
  size?: number
  level?: 'L' | 'M' | 'Q' | 'H'
  includeMargin?: boolean
  backgroundColor?: string
  foregroundColor?: string
}

export interface QRCodeDisplayEmits {
  generated: [dataUrl: string]
  error: [error: Error]
}

// 状态组件 Props
export interface StatusChipProps {
  status: OrderStatus | string
  variant?: 'default' | 'outlined' | 'text'
  size?: 'small' | 'default' | 'large'
}

export interface LoadingSpinnerProps {
  size?: string | number
  color?: string
  indeterminate?: boolean
  width?: string | number
}

export interface EmptyStateProps {
  icon?: string
  title?: string
  description?: string
  actionText?: string
  showAction?: boolean
}

export interface EmptyStateEmits {
  action: []
}

// 表单组件 Props
export interface SearchBarProps {
  modelValue: string
  placeholder?: string
  clearable?: boolean
  loading?: boolean
  debounce?: number
}

export interface SearchBarEmits {
  'update:modelValue': [value: string]
  search: [query: string]
  clear: []
}

export interface FilterChipsProps {
  filters: Array<{
    key: string
    label: string
    value: any
    active?: boolean
  }>
  multiple?: boolean
}

export interface FilterChipsEmits {
  change: [filters: Array<{ key: string; value: any }>]
}

// 导航组件 Props
export interface TabsProps {
  modelValue: string | number
  items: Array<{
    value: string | number
    title: string
    icon?: string
    disabled?: boolean
    badge?: string | number
  }>
  variant?: 'default' | 'pills' | 'underline'
  color?: string
  centered?: boolean
}

export interface TabsEmits {
  'update:modelValue': [value: string | number]
  change: [value: string | number]
}

export interface BottomNavigationProps {
  modelValue: string | number
  items: Array<{
    value: string | number
    title: string
    icon: string
    badge?: string | number
    disabled?: boolean
  }>
  color?: string
  grow?: boolean
}

export interface BottomNavigationEmits {
  'update:modelValue': [value: string | number]
  change: [value: string | number]
}

// 数据展示组件 Props
export interface StatsCardProps {
  title: string
  value: string | number
  change?: number
  changeType?: 'increase' | 'decrease' | 'neutral'
  icon?: string
  color?: string
  loading?: boolean
}

export interface DataTableProps<T = any> {
  items: T[]
  headers: Array<{
    key: string
    title: string
    sortable?: boolean
    width?: string | number
    align?: 'start' | 'center' | 'end'
  }>
  loading?: boolean
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  itemsPerPage?: number
  page?: number
  showSelect?: boolean
  showExpand?: boolean
}

export interface DataTableEmits<T = any> {
  'update:sortBy': [key: string]
  'update:sortOrder': [order: 'asc' | 'desc']
  'update:page': [page: number]
  'update:itemsPerPage': [count: number]
  'update:selected': [items: T[]]
  'update:expanded': [items: T[]]
  'click:row': [item: T, index: number]
}

// 反馈组件 Props
export interface SnackbarProps {
  modelValue: boolean
  text: string
  color?: string
  timeout?: number
  location?: string
  variant?: 'default' | 'outlined' | 'text' | 'elevated' | 'flat' | 'tonal'
  closable?: boolean
  multiLine?: boolean
}

export interface SnackbarEmits {
  'update:modelValue': [value: boolean]
  close: []
}

export interface DialogProps {
  modelValue: boolean
  title?: string
  text?: string
  persistent?: boolean
  maxWidth?: string | number
  fullscreen?: boolean
  scrollable?: boolean
}

export interface DialogEmits {
  'update:modelValue': [value: boolean]
  confirm: []
  cancel: []
}

// 通用事件类型
export interface ClickEvent {
  originalEvent: MouseEvent
  stop: () => void
  prevent: () => void
}

export interface InputEvent {
  target: HTMLInputElement
  value: string
}

export interface SelectEvent<T = any> {
  item: T
  value: any
  index: number
}

// 路由相关类型
export interface RouteMetaData {
  title?: string
  showBackButton?: boolean
  showBottomNav?: boolean
  requiresAuth?: boolean
  keepAlive?: boolean
  transition?: string
}

// 主题相关类型
export interface ThemeConfig {
  primary: string
  secondary: string
  accent: string
  error: string
  warning: string
  info: string
  success: string
  surface: string
  background: string
}

export interface ThemeVariables {
  '--v-theme-primary': string
  '--v-theme-secondary': string
  '--v-theme-accent': string
  '--v-theme-error': string
  '--v-theme-warning': string
  '--v-theme-info': string
  '--v-theme-success': string
  '--v-theme-surface': string
  '--v-theme-background': string
}