// eSIM 相关 API
import { apiClient } from './client'
import type { 
  ESIMUsageResponse, 
  PackageHistoryItem, 
  TopupPackagesResponse, 
  TopupRequest, 
  TopupResponse 
} from '@/types'

export class ESIMApi {
  // 获取 eSIM 使用情况
  async getUsage(orderId: number): Promise<ESIMUsageResponse> {
    try {
      const response = await apiClient.get<ESIMUsageResponse>(`/miniapp/esim/usage/${orderId}`)
      return response
    } catch (error: any) {
      console.error('[ESIM] 获取使用情况失败:', error)
      
      if (error.code === 'HTTP_404') {
        throw new Error('eSIM 使用信息暂未生成，请稍后再试')
      } else if (error.code === 'NETWORK_ERROR') {
        throw new Error('网络连接失败，请检查网络后重试')
      } else {
        throw new Error(error.message || '获取 eSIM 使用情况失败')
      }
    }
  }

  // 获取套餐历史
  async getHistory(orderId: number): Promise<PackageHistoryItem[]> {
    try {
      const response = await apiClient.get<{ packages: PackageHistoryItem[] }>(`/miniapp/esim/history/${orderId}`)
      return response.packages || []
    } catch (error: any) {
      console.error('[ESIM] 获取套餐历史失败:', error)
      
      if (error.code === 'HTTP_404') {
        throw new Error('暂无套餐历史记录')
      } else if (error.code === 'NETWORK_ERROR') {
        throw new Error('网络连接失败，请检查网络后重试')
      } else {
        throw new Error(error.message || '获取套餐历史失败')
      }
    }
  }

  // 获取充值套餐
  async getTopupPackages(orderId: number): Promise<TopupPackagesResponse> {
    try {
      const response = await apiClient.get<TopupPackagesResponse>(`/miniapp/esim/topup-packages/${orderId}`)
      return response
    } catch (error: any) {
      console.error('[ESIM] 获取充值套餐失败:', error)
      
      if (error.code === 'HTTP_404') {
        throw new Error('暂无可用充值套餐')
      } else if (error.code === 'NETWORK_ERROR') {
        throw new Error('网络连接失败，请检查网络后重试')
      } else {
        throw new Error(error.message || '获取充值套餐失败')
      }
    }
  }

  // 充值流量
  async topupEsim(orderId: number, request: TopupRequest): Promise<TopupResponse> {
    try {
      const response = await apiClient.post<TopupResponse>(`/miniapp/esim/topup/${orderId}`, request)
      return response
    } catch (error: any) {
      console.error('[ESIM] 充值失败:', error)
      
      if (error.code === 'HTTP_400') {
        throw new Error('充值参数错误，请重新选择套餐')
      } else if (error.code === 'HTTP_402') {
        throw new Error('余额不足，请先充值')
      } else if (error.code === 'HTTP_409') {
        throw new Error('充值冲突，请稍后重试')
      } else if (error.code === 'NETWORK_ERROR') {
        throw new Error('网络连接失败，请检查网络后重试')
      } else {
        throw new Error(error.message || 'eSIM 充值失败')
      }
    }
  }

  // 导出 PDF
  async exportPDF(orderId: number): Promise<Blob> {
    try {
      // 直接使用环境变量中的 API 基础 URL
      const baseURL = import.meta.env['VITE_API_BASE_URL'] || 'http://localhost:8080/api'
      const response = await fetch(`${baseURL}/miniapp/esim/pdf/${orderId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token') || ''}`,
          'X-Telegram-Init-Data': localStorage.getItem('telegram_init_data') || ''
        }
      })
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }
      
      return await response.blob()
    } catch (error) {
      console.error('[ESIM] PDF 导出失败:', error)
      throw new Error('PDF 导出失败')
    }
  }

  // 批量导出所有 eSIM 的 PDF
  async exportAllPDF(orderId: number): Promise<Blob> {
    try {
      const baseURL = import.meta.env['VITE_API_BASE_URL'] || 'http://localhost:8080/api'
      const response = await fetch(`${baseURL}/miniapp/esim/pdf-all/${orderId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('auth_token') || ''}`,
          'X-Telegram-Init-Data': localStorage.getItem('telegram_init_data') || ''
        }
      })
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`)
      }
      
      return await response.blob()
    } catch (error) {
      console.error('[ESIM] 批量 PDF 导出失败:', error)
      throw new Error('批量 PDF 导出失败')
    }
  }

  // 获取安装说明
  async getInstallGuide(orderId: number): Promise<{ guide: string; qrCode: string }> {
    try {
      const response = await apiClient.get<{ guide: string; qrCode: string }>(`/miniapp/esim/install-guide/${orderId}`)
      return response
    } catch (error) {
      console.error('[ESIM] 获取安装说明失败:', error)
      throw new Error('获取安装说明失败')
    }
  }
}

// 创建 eSIM API 实例
export const esimApi = new ESIMApi()