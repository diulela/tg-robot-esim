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
  HttpStatusCode
} from '@/types'

// API 客户端配置
interface ApiClientConfig {
  baseURL: string
  timeout: number
  retryAttempts: number
  retryDelay: number
}

// 默认配置
const defaultConfig: ApiClientConfig = {
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
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
      const response: AxiosResponse<ApiResponse<T>> = await this.client(config)
      
      if (response.data.success) {
        return response.data.data as T
      } else {
        throw new Error(response.data.message || response.data.error || '请求失败')
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
  // 获取产品列表
  async getProducts(params?: ProductQueryParams): Promise<PaginatedResponse<Product>> {
    return apiClient.get('/products', params)
  },

  // 获取产品详情
  async getProduct(id: string): Promise<Product> {
    return apiClient.get(`/products/${id}`)
  },

  // 根据区域获取产品
  async getProductsByRegion(region: string, params?: Omit<ProductQueryParams, 'region'>): Promise<PaginatedResponse<Product>> {
    return apiClient.get(`/products/region/${region}`, params)
  },

  // 根据国家获取产品
  async getProductsByCountry(country: string, params?: Omit<ProductQueryParams, 'country'>): Promise<PaginatedResponse<Product>> {
    return apiClient.get(`/products/country/${country}`, params)
  },

  // 搜索产品
  async searchProducts(query: string, params?: ProductQueryParams): Promise<PaginatedResponse<Product>> {
    return apiClient.get('/products/search', { ...params, q: query })
  }
}

// 订单相关 API
export const orderApi = {
  // 获取订单列表
  async getOrders(params?: OrderQueryParams): Promise<PaginatedResponse<Order>> {
    return apiClient.get('/orders', params)
  },

  // 获取订单详情
  async getOrder(id: string): Promise<Order> {
    return apiClient.get(`/orders/${id}`)
  },

  // 创建订单
  async createOrder(data: CreateOrderRequest): Promise<Order> {
    return apiClient.post('/orders', data)
  },

  // 取消订单
  async cancelOrder(id: string, reason?: string): Promise<Order> {
    return apiClient.patch(`/orders/${id}/cancel`, { reason })
  },

  // 获取订单状态
  async getOrderStatus(id: string): Promise<{ status: string; updatedAt: string }> {
    return apiClient.get(`/orders/${id}/status`)
  },

  // 重新支付订单
  async retryPayment(id: string): Promise<{ paymentUrl: string }> {
    return apiClient.post(`/orders/${id}/retry-payment`)
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
  // 获取钱包信息
  async getWallet(): Promise<Wallet> {
    return apiClient.get('/wallet')
  },

  // 获取钱包交易记录
  async getTransactions(params?: { page?: number; pageSize?: number; type?: string }): Promise<PaginatedResponse<WalletTransaction>> {
    return apiClient.get('/wallet/transactions', params)
  },

  // 钱包充值
  async recharge(data: WalletRechargeRequest): Promise<{ paymentUrl: string; transactionId: string }> {
    return apiClient.post('/wallet/recharge', data)
  },

  // 获取充值状态
  async getRechargeStatus(transactionId: string): Promise<{ status: string; amount?: number }> {
    return apiClient.get(`/wallet/recharge/${transactionId}/status`)
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