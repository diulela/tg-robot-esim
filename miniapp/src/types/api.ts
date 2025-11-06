// API 和业务数据类型定义

// 产品相关类型
export interface Product {
  id: string
  name: string
  description: string
  price: number
  originalPrice?: number // 原价（用于显示折扣）
  currency: string
  region: string
  country: string
  countryCode: string
  dataAmount: string
  validDays: number
  coverage: string // 覆盖地区
  coverageAreas: string[]
  features: string[] // 特性标签
  icon?: string // 商品图标
  isActive: boolean
  isPopular: boolean // 是否热门
  createdAt: string
  updatedAt: string
}

// 热门项类型（前端硬编码数据）
export interface HotItem {
  code: string // 国家/地区代码
  name: string // 显示名称
}

// 订单状态枚举
export enum OrderStatus {
  PENDING = 'pending',
  PAID = 'paid',
  PROCESSING = 'processing',
  COMPLETED = 'completed',
  CANCELLED = 'cancelled',
  REFUNDED = 'refunded'
}

// 支付方式枚举
export enum PaymentMethod {
  USDT_TRC20 = 'usdt_trc20',
  TELEGRAM_STARS = 'telegram_stars',
  CREDIT_CARD = 'credit_card'
}

// 订单类型
export interface Order {
  id: string
  orderNumber: string
  productId: string
  productName: string
  amount: number
  currency: string
  status: OrderStatus
  paymentMethod?: PaymentMethod
  createdAt: string
  paidAt?: string
  completedAt?: string
  esimInfo?: ESIMInfo
  transactionHash?: string
  refundReason?: string
  
  // 新增 eSIM 相关字段
  esimUsage?: ESIMUsageInfo
  packageHistory?: PackageHistoryItem[]
  canTopup?: boolean
  canExportPDF?: boolean
}

// eSIM 状态枚举
export enum ESIMStatus {
  PENDING = 'pending',
  ACTIVE = 'active',
  EXPIRED = 'expired',
  SUSPENDED = 'suspended',
  TERMINATED = 'terminated'
}

// eSIM 信息类型
export interface ESIMInfo {
  iccid: string
  activationCode: string
  qrCode: string
  apnType: 'manual' | 'automatic'
  apnValue?: string
  isRoaming: boolean
  activatedAt?: string
  expiresAt?: string
  dataUsed?: number
  dataRemaining?: number
  dataTotal?: number
  usagePercentage?: string
  status?: ESIMStatus
  lpaAddress?: string
  confirmationCode?: string
}

// eSIM 使用情况信息
export interface ESIMUsageInfo {
  iccid: string
  status: ESIMStatus
  activationTime: string
  expireTime: string
  dataUsed: number // MB
  dataTotal: number // MB
  dataRemaining: number // MB
  usagePercentage: string
}

// eSIM 使用情况响应
export interface ESIMUsageResponse {
  success: boolean
  message: string
  data: {
    orderId: number
    esim: ESIMUsageInfo
  }
}

// 套餐历史项
export interface PackageHistoryItem {
  id: string
  packageName: string
  dataSize: string
  validDays: number
  price: number
  status: 'FINISHED' | 'ACTIVE' | 'EXPIRED'
  remainingData: string
  activationTime: string
  expireTime: string
  createdAt: string
}

// 充值套餐
export interface TopupPackage {
  id: string
  title: string
  data: string
  price: number
  validity: number
  description: string
}

// 充值套餐响应
export interface TopupPackagesResponse {
  success: boolean
  message: string
  data: {
    orderId: number
    packages: TopupPackage[]
  }
}

// 充值请求
export interface TopupRequest {
  packageId: string
  description?: string
}

// 充值响应
export interface TopupResponse {
  success: boolean
  message: string
  data: {
    topupOrderId: number
    orderId: number
    packageId: string
    amount: number
    status: string
  }
}

// 区域类型
export interface Region {
  id: string
  name: string
  code: string
  icon: string
  countries: Country[]
  isPopular: boolean
  displayOrder: number
}

// 国家类型
export interface Country {
  id: string
  name: string
  code: string
  flag: string
  region: string
  products: Product[]
  isPopular: boolean
}

// 用户钱包类型
export interface Wallet {
  id: string
  userId: string
  balance: number
  currency: string
  frozenAmount: number
  totalRecharge: number
  totalSpent: number
  createdAt: string
  updatedAt: string
}

// 钱包交易类型
export interface WalletTransaction {
  id: number
  user_id: number
  type: 'recharge' | 'payment' | 'refund' | 'freeze' | 'unfreeze'
  amount: string
  balance_before: string
  balance_after: string
  status: 'pending' | 'completed' | 'failed' | 'cancelled'
  description: string
  related_type: string
  related_id: string
  tx_hash: string
  metadata?: string
  created_at: string
  updated_at: string
}

// 钱包历史统计类型
export interface WalletHistoryStats {
  total_records: number
  total_income: string
  total_expense: string
  pending_amount: string
  completed_amount: string
}

// API 响应基础类型
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
  code?: number
}

// 分页响应类型
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
  hasNext: boolean
  hasPrev: boolean
}

// API 请求参数类型
export interface ApiRequestParams {
  [key: string]: string | number | boolean | undefined
}

// 产品查询参数
export interface ProductQueryParams extends ApiRequestParams {
  type?: string // local, regional, global
  country?: string // 国家代码筛选
  search?: string // 搜索关键词
  minPrice?: number
  maxPrice?: number
  limit?: number // 每页数量
  offset?: number // 偏移量
  page?: number
  pageSize?: number
  sortBy?: 'price' | 'validDays' | 'dataAmount' | 'createdAt'
  sortOrder?: 'asc' | 'desc'
}

// 订单查询参数
export interface OrderQueryParams extends ApiRequestParams {
  status?: OrderStatus
  startDate?: string
  endDate?: string
  page?: number
  pageSize?: number
  sortBy?: 'createdAt' | 'amount' | 'status'
  sortOrder?: 'asc' | 'desc'
}

// 创建订单请求
export interface CreateOrderRequest {
  product_id: number // 后端期望的是 number 类型的 product_id
  paymentMethod?: PaymentMethod
  quantity?: number
  order_note?: string // 订单备注
}

// 购买请求数据结构
export interface PurchaseRequest {
  product_id: number
  quantity: number
  order_note?: string
  payment_method?: PaymentMethod
}

// 钱包充值请求
export interface WalletRechargeRequest {
  amount: number
  currency: string
  paymentMethod: PaymentMethod
}

// 统计数据类型
export interface DashboardStats {
  totalOrders: number
  pendingOrders: number
  completedOrders: number
  totalSpent: number
  currency: string
  orderTrends: {
    period: string
    count: number
    change: number
  }[]
  recentOrders: Order[]
}

// 错误类型
export interface ApiError {
  code: string
  message: string
  details?: Record<string, any>
  timestamp: string
}

// HTTP 错误状态码
export enum HttpStatusCode {
  OK = 200,
  CREATED = 201,
  NO_CONTENT = 204,
  BAD_REQUEST = 400,
  UNAUTHORIZED = 401,
  FORBIDDEN = 403,
  NOT_FOUND = 404,
  CONFLICT = 409,
  UNPROCESSABLE_ENTITY = 422,
  INTERNAL_SERVER_ERROR = 500,
  BAD_GATEWAY = 502,
  SERVICE_UNAVAILABLE = 503
}