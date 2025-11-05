// API 服务层实现
import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import type {
  ApiResponse,
  PaginatedResponse,
  Product,
  Order,
  Region,
  Country,
  Wallet,
  WalletTransaction,
  DashboardStats,
  ProductQueryParams,
  OrderQueryParams,
  CreateOrderRequest,
  WalletRechargeRequest,
  ApiError,
  // HttpStatusCode
} from '@/types'
import { OrderStatus } from '@/types'

// API 客户端配置
interface ApiClientConfig {
  baseURL: string
  timeout: number
  retryAttempts: number
  retryDelay: number
}

// 默认配置
const defaultConfig: ApiClientConfig = {
  baseURL: import.meta.env["VITE_API_BASE_URL"] || 'http://localhost:8080/api',
  timeout: 10000,
  retryAttempts: 3,
  retryDelay: 1000
}

class ApiClient {
  private client: AxiosInstance
  private config: ApiClientConfig

  constructor(config: Partial<ApiClientConfig> = {}) {
    this.config = { ...defaultConfig, ...config }
    this.client = this.createAxiosInstance()
    this.setupInterceptors()
  }

  private createAxiosInstance(): AxiosInstance {
    return axios.create({
      baseURL: this.config.baseURL,
      timeout: this.config.timeout,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    })
  }

  private setupInterceptors(): void {
    // 请求拦截器
    this.client.interceptors.request.use(
      (config) => {
        // 添加 Telegram 初始化数据
        const telegramInitData = this.getTelegramInitData()
        if (telegramInitData) {
          config.headers['X-Telegram-Init-Data'] = telegramInitData
        }

        // 添加时间戳防止缓存
        if (config.method === 'get') {
          config.params = {
            ...config.params,
            _t: Date.now()
          }
        }

        console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`, config.params || config.data)
        return config
      },
      (error) => {
        console.error('[API] Request error:', error)
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.client.interceptors.response.use(
      (response) => {
        console.log(`[API] Response:`, response.data)
        return response
      },
      async (error) => {
        const originalRequest = error.config

        // 重试逻辑
        if (this.shouldRetry(error) && !originalRequest._retry) {
          originalRequest._retry = true
          originalRequest._retryCount = (originalRequest._retryCount || 0) + 1

          if (originalRequest._retryCount <= this.config.retryAttempts) {
            console.log(`[API] Retrying request (${originalRequest._retryCount}/${this.config.retryAttempts})`)
            await this.delay(this.config.retryDelay * originalRequest._retryCount)
            return this.client(originalRequest)
          }
        }

        // 错误处理
        const apiError = this.handleError(error)
        console.error('[API] Response error:', apiError)
        return Promise.reject(apiError)
      }
    )
  }

  private getTelegramInitData(): string | null {
    if (typeof window !== 'undefined' && window.Telegram?.WebApp) {
      return window.Telegram.WebApp.initData
    }
    return null
  }

  private shouldRetry(error: any): boolean {
    // 网络错误或 5xx 服务器错误时重试
    return !error.response || (error.response.status >= 500 && error.response.status < 600)
  }

  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  private handleError(error: any): ApiError {
    const timestamp = new Date().toISOString()

    if (error.response) {
      // 服务器响应错误
      const { status, data } = error.response
      return {
        code: `HTTP_${status}`,
        message: data?.message || data?.error || `HTTP Error ${status}`,
        details: data,
        timestamp
      }
    } else if (error.request) {
      // 网络错误
      return {
        code: 'NETWORK_ERROR',
        message: '网络连接失败，请检查网络设置',
        timestamp
      }
    } else {
      // 其他错误
      return {
        code: 'UNKNOWN_ERROR',
        message: error.message || '未知错误',
        timestamp
      }
    }
  }

  // 通用请求方法
  private async request<T>(config: AxiosRequestConfig): Promise<T> {
    try {
      const response: AxiosResponse<any> = await this.client(config)

      // 适配后端响应格式: { code: 0, message: "success", data: {...} }
      if (response.data.code === 0) {
        return response.data.data as T
      } else {
        throw new Error(response.data.message || '请求失败')
      }
    } catch (error) {
      throw error
    }
  }

  // GET 请求
  async get<T>(url: string, params?: any): Promise<T> {
    return this.request<T>({ method: 'GET', url, params })
  }

  // POST 请求
  async post<T>(url: string, data?: any): Promise<T> {
    return this.request<T>({ method: 'POST', url, data })
  }

  // PUT 请求
  async put<T>(url: string, data?: any): Promise<T> {
    return this.request<T>({ method: 'PUT', url, data })
  }

  // DELETE 请求
  async delete<T>(url: string): Promise<T> {
    return this.request<T>({ method: 'DELETE', url })
  }

  // PATCH 请求
  async patch<T>(url: string, data?: any): Promise<T> {
    return this.request<T>({ method: 'PATCH', url, data })
  }
}

// 创建 API 客户端实例
const apiClient = new ApiClient()

// 产品相关 API
export const productApi = {
  // 转换后端产品数据为前端格式
  transformProduct(backendProduct: any): Product {
    // 解析 countries JSON 字符串
    let coverageAreas: string[] = []
    let countryCode = ''
    let country = ''

    try {
      if (backendProduct.countries) {
        const countriesData = JSON.parse(backendProduct.countries)
        if (Array.isArray(countriesData) && countriesData.length > 0) {
          coverageAreas = countriesData.map((c: any) => c.cn || c.name || '')
          countryCode = countriesData[0].code || ''
          country = countriesData[0].cn || countriesData[0].name || ''
        }
      }
    } catch (e) {
      console.warn('解析国家数据失败:', e)
    }

    // 解析 features JSON 字符串
    let features: string[] = []
    try {
      if (backendProduct.features) {
        features = JSON.parse(backendProduct.features)
      }
    } catch (e) {
      console.warn('解析特性数据失败:', e)
    }

    // 格式化数据量
    const dataAmount = backendProduct.data_size >= 1024
      ? `${(backendProduct.data_size / 1024).toFixed(1)}GB`
      : `${backendProduct.data_size}MB`

    return {
      id: String(backendProduct.id),
      name: backendProduct.name || '',
      description: backendProduct.description || '',
      price: backendProduct.price || 0,
      originalPrice: backendProduct.retail_price || undefined,
      currency: 'USD',
      region: backendProduct.type || 'global',
      country: country,
      countryCode: countryCode,
      dataAmount: dataAmount,
      validDays: backendProduct.valid_days || 0,
      coverage: backendProduct.name || '',
      coverageAreas: coverageAreas,
      features: features,
      icon: backendProduct.image || undefined,
      isActive: backendProduct.status === 'active',
      isPopular: backendProduct.is_hot || false,
      createdAt: backendProduct.created_at || new Date().toISOString(),
      updatedAt: backendProduct.updated_at || new Date().toISOString()
    }
  },

  // 转换后端响应为前端格式
  transformProductResponse(backendData: any): PaginatedResponse<Product> {
    const { products = [], total = 0, limit = 20, offset = 0 } = backendData
    const page = Math.floor(offset / limit) + 1
    const totalPages = Math.ceil(total / limit)

    // 转换每个产品数据
    const transformedProducts = products.map((p: any) => this.transformProduct(p))

    return {
      items: transformedProducts,
      total,
      page,
      pageSize: limit,
      totalPages,
      hasNext: offset + limit < total,
      hasPrev: offset > 0
    }
  },

  // 获取产品列表
  async getProducts(params?: ProductQueryParams): Promise<PaginatedResponse<Product>> {
    // 转换前端参数为后端格式
    const backendParams: any = {}
    if (params?.type) backendParams.type = params.type
    if (params?.country) backendParams.country = params.country
    if (params?.search) backendParams.search = params.search
    if (params?.page) {
      const pageSize = params.pageSize || 20
      backendParams.limit = pageSize
      backendParams.offset = (params.page - 1) * pageSize
    } else {
      backendParams.limit = params?.pageSize || 20
      backendParams.offset = params?.offset || 0
    }

    const data = await apiClient.get('/miniapp/products', backendParams)
    return this.transformProductResponse(data)
  },

  // 获取产品详情
  async getProduct(id: string): Promise<Product> {
    const data = await apiClient.get(`/miniapp/products/${id}`)
    return this.transformProduct(data)
  },

  // 根据类型获取产品
  async getProductsByType(type: string, params?: Omit<ProductQueryParams, 'type'>): Promise<PaginatedResponse<Product>> {
    return this.getProducts({ ...params, type })
  },

  // 根据国家获取产品
  async getProductsByCountry(country: string, params?: Omit<ProductQueryParams, 'country'>): Promise<PaginatedResponse<Product>> {
    return this.getProducts({ ...params, country })
  },

  // 搜索产品
  async searchProducts(query: string, params?: ProductQueryParams): Promise<PaginatedResponse<Product>> {
    return this.getProducts({ ...params, search: query })
  }
}

// 订单相关 API
export const orderApi = {
  // 转换后端订单数据为前端格式
  transformOrder(backendOrder: any): Order {
    const order: Order = {
      id: String(backendOrder.id),
      orderNumber: backendOrder.order_number || '',
      productId: String(backendOrder.product_id),
      productName: backendOrder.product_name || '',
      amount: backendOrder.amount || 0,
      currency: backendOrder.currency || 'USD',
      status: backendOrder.status || OrderStatus.PENDING,
      createdAt: backendOrder.created_at || new Date().toISOString()
    }

    // 可选字段
    if (backendOrder.payment_method) order.paymentMethod = backendOrder.payment_method
    if (backendOrder.paid_at) order.paidAt = backendOrder.paid_at
    if (backendOrder.completed_at) order.completedAt = backendOrder.completed_at
    if (backendOrder.transaction_hash) order.transactionHash = backendOrder.transaction_hash
    if (backendOrder.refund_reason) order.refundReason = backendOrder.refund_reason

    // eSIM 信息
    if (backendOrder.esim_info) {
      order.esimInfo = {
        iccid: backendOrder.esim_info.iccid || '',
        activationCode: backendOrder.esim_info.activation_code || '',
        qrCode: backendOrder.esim_info.qr_code || '',
        apnType: backendOrder.esim_info.apn_type || 'automatic',
        isRoaming: backendOrder.esim_info.is_roaming || false,
        activatedAt: backendOrder.esim_info.activated_at,
        expiresAt: backendOrder.esim_info.expires_at,
        dataUsed: backendOrder.esim_info.data_used,
        dataRemaining: backendOrder.esim_info.data_remaining
      }
    }

    return order
  },

  // 转换后端订单响应为前端格式
  transformOrderResponse(backendData: any): PaginatedResponse<Order> {
    const { orders = [], stats, limit = 20, offset = 0 } = backendData
    const total = stats?.total_orders || 0
    const page = Math.floor(offset / limit) + 1
    const totalPages = Math.ceil(total / limit)

    // 转换每个订单数据
    const transformedOrders = orders.map((o: any) => this.transformOrder(o))

    return {
      items: transformedOrders,
      total,
      page,
      pageSize: limit,
      totalPages,
      hasNext: offset + limit < total,
      hasPrev: offset > 0
    }
  },

  // 获取订单列表
  async getOrders(params?: OrderQueryParams): Promise<PaginatedResponse<Order>> {
    // 转换前端参数为后端格式
    const backendParams: any = {}
    if (params?.status) backendParams.status = params.status
    if (params?.startDate) backendParams.start_date = params.startDate
    if (params?.endDate) backendParams.end_date = params.endDate
    if (params?.page) {
      const pageSize = params.pageSize || 20
      backendParams.limit = pageSize
      backendParams.offset = (params.page - 1) * pageSize
    } else {
      backendParams.limit = params?.pageSize || 20
      backendParams.offset = 0
    }

    const data = await apiClient.get('/miniapp/orders', backendParams)
    return this.transformOrderResponse(data)
  },

  // 获取订单详情
  async getOrder(id: string): Promise<Order> {
    const data = await apiClient.get(`/miniapp/orders/${id}`)
    return this.transformOrder(data)
  },

  // 创建订单 (购买产品)
  async createOrder(data: CreateOrderRequest): Promise<Order> {
    const result = await apiClient.post('/miniapp/purchase', data)
    return this.transformOrder(result)
  },

  // 取消订单
  async cancelOrder(id: string, reason?: string): Promise<Order> {
    const result = await apiClient.patch(`/miniapp/orders/${id}/cancel`, { reason })
    return this.transformOrder(result)
  },

  // 获取订单状态
  async getOrderStatus(id: string): Promise<{ status: string; updatedAt: string }> {
    return apiClient.get(`/miniapp/orders/${id}/status`)
  },

  // 重新支付订单
  async retryPayment(id: string): Promise<{ paymentUrl: string }> {
    return apiClient.post(`/miniapp/orders/${id}/retry-payment`)
  }
}

// 区域和国家相关 API
export const regionApi = {
  // 获取所有区域
  async getRegions(): Promise<Region[]> {
    return apiClient.get('/regions')
  },

  // 获取热门区域
  async getPopularRegions(): Promise<Region[]> {
    return apiClient.get('/regions/popular')
  },

  // 获取区域详情
  async getRegion(id: string): Promise<Region> {
    return apiClient.get(`/regions/${id}`)
  },

  // 获取所有国家
  async getCountries(): Promise<Country[]> {
    return apiClient.get('/countries')
  },

  // 根据区域获取国家
  async getCountriesByRegion(region: string): Promise<Country[]> {
    return apiClient.get(`/countries/region/${region}`)
  },

  // 获取热门国家
  async getPopularCountries(): Promise<Country[]> {
    return apiClient.get('/countries/popular')
  },

  // 搜索国家
  async searchCountries(query: string): Promise<Country[]> {
    return apiClient.get('/countries/search', { q: query })
  }
}

// 钱包相关 API
export const walletApi = {
  // 获取钱包余额
  async getWallet(): Promise<Wallet> {
    return apiClient.get('/miniapp/wallet/balance')
  },

  // 获取交易记录
  async getTransactions(params?: { limit?: number; offset?: number; type?: string }): Promise<PaginatedResponse<WalletTransaction>> {
    return apiClient.get('/miniapp/transactions', params)
  },

  // 钱包充值（旧版本，保持兼容性）
  async recharge(data: WalletRechargeRequest): Promise<{ paymentUrl: string; transactionId: string }> {
    return apiClient.post('/miniapp/wallet/recharge', data)
  },

  // 获取充值状态（旧版本）
  async getRechargeStatus(transactionId: string): Promise<{ status: string; amount?: number }> {
    return apiClient.get(`/miniapp/wallet/recharge/${transactionId}/status`)
  },

  // 创建 USDT 充值订单
  async createRechargeOrder(data: { amount: string }): Promise<{
    order_no: string;
    amount: string;
    exact_amount: string;
    wallet_address: string;
    status: string;
    expires_at: string;
    created_at: string;
  }> {
    return apiClient.post('/miniapp/wallet/recharge', data)
  },

  // 获取充值订单详情
  async getRechargeOrder(orderNo: string): Promise<{
    order_no: string;
    amount: string;
    exact_amount: string;
    wallet_address: string;
    status: string;
    tx_hash: string;
    confirmations: number;
    expires_at: string;
    confirmed_at?: string;
    created_at: string;
  }> {
    return apiClient.get(`/miniapp/wallet/recharge/${orderNo}`)
  },

  // 手动检查充值状态
  async checkRechargeStatus(orderNo: string): Promise<{
    order_no: string;
    status: string;
    tx_hash: string;
    confirmations: number;
    confirmed_at?: string;
  }> {
    return apiClient.post(`/miniapp/wallet/recharge/${orderNo}/check`)
  },

  // 获取充值历史
  async getRechargeHistory(params?: { limit?: number; offset?: number }): Promise<{
    orders: Array<{
      order_no: string;
      amount: string;
      status: string;
      tx_hash: string;
      created_at: string;
      confirmed_at?: string;
    }>;
    total: number;
    limit: number;
    offset: number;
  }> {
    return apiClient.get('/miniapp/wallet/recharge/history', params)
  }
}

// 用户相关 API
export const userApi = {
  // 获取用户信息
  async getProfile(): Promise<{ id: string; telegramId: number; firstName: string; lastName?: string; username?: string }> {
    return apiClient.get('/user/profile')
  },

  // 更新用户信息
  async updateProfile(data: { firstName?: string; lastName?: string }): Promise<void> {
    return apiClient.put('/user/profile', data)
  },

  // 获取用户统计
  async getStats(): Promise<DashboardStats> {
    return apiClient.get('/user/stats')
  }
}

// 系统相关 API
export const systemApi = {
  // 健康检查
  async healthCheck(): Promise<{ status: string; timestamp: string }> {
    return apiClient.get('/health')
  },

  // 获取系统配置
  async getConfig(): Promise<{ version: string; features: Record<string, boolean> }> {
    return apiClient.get('/config')
  }
}

// 导出 API 客户端
export { apiClient }
export default {
  product: productApi,
  order: orderApi,
  region: regionApi,
  wallet: walletApi,
  user: userApi,
  system: systemApi
}