// API 客户端基础类
import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import type { ApiError } from '@/types'

// API 客户端配置
export interface ApiClientConfig {
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

export class ApiClient {
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

// 创建默认的 API 客户端实例
export const apiClient = new ApiClient()