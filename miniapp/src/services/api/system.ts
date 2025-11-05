// 系统相关 API
import { apiClient } from './client'

// 健康检查响应接口
export interface HealthCheckResponse {
  status: string
  timestamp: string
}

// 系统配置响应接口
export interface SystemConfigResponse {
  version: string
  features: Record<string, boolean>
}

export class SystemApi {
  // 健康检查
  async healthCheck(): Promise<HealthCheckResponse> {
    return apiClient.get('/health')
  }

  // 获取系统配置
  async getConfig(): Promise<SystemConfigResponse> {
    return apiClient.get('/config')
  }
}

// 创建系统 API 实例
export const systemApi = new SystemApi()