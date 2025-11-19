// eSIM 状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { 
  ESIMUsageInfo, 
  PackageHistoryItem, 
  TopupPackage,
  ESIMUsageResponse,
  TopupPackagesResponse,
  TopupResponse
} from '@/types'
import { esimApi } from '@/services/api'

// 错误处理工具类
export class ESIMErrorHandler {
  static handleApiError(error: any): string {
    if (error.code === 'NETWORK_TIMEOUT') {
      return '网络超时，请检查网络连接后重试'
    }
    
    if (error.code === 'ESIM_NOT_FOUND') {
      return 'eSIM 信息暂未生成，请稍后再试'
    }
    
    if (error.code === 'TOPUP_PACKAGES_UNAVAILABLE') {
      return '暂无可用充值套餐，请联系客服'
    }
    
    return error.message || '操作失败，请重试'
  }

  static showErrorNotification(message: string, appStore: any) {
    appStore.showNotification({
      type: 'error',
      message,
      duration: 4000
    })
  }
}

export const useESIMStore = defineStore('esim', () => {
  // 状态
  const usageInfo = ref<ESIMUsageInfo | null>(null)
  const packageHistory = ref<PackageHistoryItem[]>([])
  const topupPackages = ref<TopupPackage[]>([])
  const isLoadingUsage = ref(false)
  const isLoadingHistory = ref(false)
  const isLoadingTopup = ref(false)
  const isTopupping = ref(false)
  const error = ref<string | null>(null)

  // 缓存管理
  const usageCache = ref<{ [orderId: number]: { data: ESIMUsageInfo; timestamp: number } }>({})
  const historyCache = ref<{ [orderId: number]: { data: PackageHistoryItem[]; timestamp: number } }>({})
  const packagesCache = ref<{ [orderId: number]: { data: TopupPackage[]; timestamp: number } }>({})

  // 缓存过期时间（毫秒）
  const CACHE_TIMEOUT = {
    usage: 5 * 60 * 1000, // 5分钟
    history: 30 * 60 * 1000, // 30分钟
    packages: 10 * 60 * 1000 // 10分钟
  }

  // 计算属性
  const hasUsageInfo = computed(() => !!usageInfo.value)
  const hasPackageHistory = computed(() => packageHistory.value.length > 0)
  const hasTopupPackages = computed(() => topupPackages.value.length > 0)

  // 工具方法
  const isCacheValid = (timestamp: number, timeout: number): boolean => {
    return Date.now() - timestamp < timeout
  }

  // 操作方法
  const fetchUsage = async (orderId: number, forceRefresh = false): Promise<void> => {
    // 检查缓存
    const cached = usageCache.value[orderId]
    if (!forceRefresh && cached && isCacheValid(cached.timestamp, CACHE_TIMEOUT.usage)) {
      usageInfo.value = cached.data
      return
    }

    isLoadingUsage.value = true
    error.value = null

    try {
      const response: ESIMUsageResponse = await esimApi.getUsage(orderId)
      if (response.success) {
        usageInfo.value = response.data.esim
        
        // 更新缓存
        usageCache.value[orderId] = {
          data: response.data.esim,
          timestamp: Date.now()
        }
      } else {
        throw new Error(response.message || '获取使用情况失败')
      }
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] 获取使用情况失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoadingUsage.value = false
    }
  }

  const fetchHistory = async (orderId: number, forceRefresh = false): Promise<void> => {
    // 检查缓存
    const cached = historyCache.value[orderId]
    if (!forceRefresh && cached && isCacheValid(cached.timestamp, CACHE_TIMEOUT.history)) {
      packageHistory.value = cached.data
      return
    }

    isLoadingHistory.value = true
    error.value = null

    try {
      const history = await esimApi.getHistory(orderId)
      packageHistory.value = history
      
      // 更新缓存
      historyCache.value[orderId] = {
        data: history,
        timestamp: Date.now()
      }
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] 获取套餐历史失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoadingHistory.value = false
    }
  }

  const fetchTopupPackages = async (orderId: number, forceRefresh = false): Promise<void> => {
    // 检查缓存
    const cached = packagesCache.value[orderId]
    if (!forceRefresh && cached && isCacheValid(cached.timestamp, CACHE_TIMEOUT.packages)) {
      topupPackages.value = cached.data
      return
    }

    isLoadingTopup.value = true
    error.value = null

    try {
      const response: TopupPackagesResponse = await esimApi.getTopupPackages(orderId)
      if (response.success) {
        topupPackages.value = response.data.packages
        
        // 更新缓存
        packagesCache.value[orderId] = {
          data: response.data.packages,
          timestamp: Date.now()
        }
      } else {
        throw new Error(response.message || '获取充值套餐失败')
      }
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] 获取充值套餐失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoadingTopup.value = false
    }
  }

  const topupEsim = async (orderId: number, packageId: string): Promise<void> => {
    isTopupping.value = true
    error.value = null

    try {
      const response: TopupResponse = await esimApi.topupEsim(orderId, { packageId })
      if (!response.success) {
        throw new Error(response.message || '充值失败')
      }
      
      // 充值成功后刷新使用情况和清除相关缓存
      delete usageCache.value[orderId]
      delete historyCache.value[orderId]
      await fetchUsage(orderId, true)
      
      console.log('[ESIM Store] 充值成功:', response.data)
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] 充值失败:', err)
      throw new Error(errorMessage)
    } finally {
      isTopupping.value = false
    }
  }

  const exportPDF = async (orderId: number): Promise<void> => {
    try {
      const blob = await esimApi.exportPDF(orderId)
      
      // 创建下载链接
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `esim-order-${orderId}.pdf`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      console.log('[ESIM Store] PDF 导出成功')
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] PDF 导出失败:', err)
      throw new Error(errorMessage)
    }
  }

  const exportAllPDF = async (orderId: number): Promise<void> => {
    try {
      const blob = await esimApi.exportAllPDF(orderId)
      
      // 创建下载链接
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `esim-order-all-${orderId}.pdf`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      console.log('[ESIM Store] 批量 PDF 导出成功')
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] 批量 PDF 导出失败:', err)
      throw new Error(errorMessage)
    }
  }

  const getInstallGuide = async (orderId: number): Promise<{ guide: string; qrCode: string }> => {
    try {
      const guide = await esimApi.getInstallGuide(orderId)
      console.log('[ESIM Store] 获取安装说明成功')
      return guide
    } catch (err) {
      const errorMessage = ESIMErrorHandler.handleApiError(err)
      error.value = errorMessage
      console.error('[ESIM Store] 获取安装说明失败:', err)
      throw new Error(errorMessage)
    }
  }

  // 工具方法
  const clearError = (): void => {
    error.value = null
  }

  const reset = (): void => {
    usageInfo.value = null
    packageHistory.value = []
    topupPackages.value = []
    error.value = null
  }

  const clearCache = (orderId?: number): void => {
    if (orderId) {
      delete usageCache.value[orderId]
      delete historyCache.value[orderId]
      delete packagesCache.value[orderId]
    } else {
      usageCache.value = {}
      historyCache.value = {}
      packagesCache.value = {}
    }
  }

  // 格式化工具方法
  const formatDataSize = (sizeInMB: number): string => {
    if (sizeInMB >= 1024) {
      return `${(sizeInMB / 1024).toFixed(2)} GB`
    }
    return `${sizeInMB.toFixed(0)} MB`
  }

  const formatDateTime = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getStatusColor = (status: string): string => {
    const colorMap: Record<string, string> = {
      'ACTIVE': 'success',
      'FINISHED': 'info',
      'EXPIRED': 'warning'
    }
    return colorMap[status] || 'default'
  }

  const getStatusText = (status: string): string => {
    const textMap: Record<string, string> = {
      'ACTIVE': '使用中',
      'FINISHED': '已用完',
      'EXPIRED': '已过期'
    }
    return textMap[status] || status
  }

  return {
    // 只读状态
    usageInfo: readonly(usageInfo),
    packageHistory: readonly(packageHistory),
    topupPackages: readonly(topupPackages),
    isLoadingUsage: readonly(isLoadingUsage),
    isLoadingHistory: readonly(isLoadingHistory),
    isLoadingTopup: readonly(isLoadingTopup),
    isTopupping: readonly(isTopupping),
    error: readonly(error),

    // 计算属性
    hasUsageInfo,
    hasPackageHistory,
    hasTopupPackages,

    // 操作方法
    fetchUsage,
    fetchHistory,
    fetchTopupPackages,
    topupEsim,
    exportPDF,
    exportAllPDF,
    getInstallGuide,
    clearError,
    reset,
    clearCache,

    // 工具方法
    formatDataSize,
    formatDateTime,
    getStatusColor,
    getStatusText
  }
})