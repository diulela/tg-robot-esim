// 应用全局状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { Notification } from '@/types'

export const useAppStore = defineStore('app', () => {
  // 状态
  const isLoading = ref(false)
  const notifications = ref<Notification[]>([])
  const theme = ref<'light' | 'dark' | 'auto'>('auto')
  const language = ref('zh-CN')
  const isOnline = ref(navigator.onLine)
  const error = ref<string | null>(null)

  // 计算属性
  const hasNotifications = computed(() => notifications.value.length > 0)
  
  const activeNotifications = computed(() => 
    notifications.value.filter(n => !n.persistent || Date.now() - n.createdAt.getTime() < (n.duration || 5000))
  )

  const currentTheme = computed(() => {
    if (theme.value === 'auto') {
      // 检查 Telegram 主题或系统主题
      if (window.Telegram?.WebApp?.colorScheme) {
        return window.Telegram.WebApp.colorScheme
      }
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }
    return theme.value
  })

  // 操作方法
  const setLoading = (loading: boolean) => {
    isLoading.value = loading
  }

  const showNotification = (notification: Omit<Notification, 'id' | 'createdAt'>) => {
    const newNotification: Notification = {
      id: `notification-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      createdAt: new Date(),
      duration: 5000,
      ...notification
    }

    notifications.value.push(newNotification)

    // 自动移除非持久化通知
    if (!newNotification.persistent && newNotification.duration) {
      setTimeout(() => {
        removeNotification(newNotification.id)
      }, newNotification.duration)
    }

    return newNotification.id
  }

  const removeNotification = (id: string) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  const clearNotifications = () => {
    notifications.value = []
  }

  const showSuccess = (message: string, duration = 3000) => {
    return showNotification({
      type: 'success',
      message,
      duration
    })
  }

  const showError = (message: string, duration = 5000) => {
    return showNotification({
      type: 'error',
      message,
      duration
    })
  }

  const showWarning = (message: string, duration = 4000) => {
    return showNotification({
      type: 'warning',
      message,
      duration
    })
  }

  const showInfo = (message: string, duration = 3000) => {
    return showNotification({
      type: 'info',
      message,
      duration
    })
  }

  const setTheme = (newTheme: 'light' | 'dark' | 'auto') => {
    theme.value = newTheme
    localStorage.setItem('app-theme', newTheme)
  }

  const setLanguage = (newLanguage: string) => {
    language.value = newLanguage
    localStorage.setItem('app-language', newLanguage)
  }

  const setOnlineStatus = (online: boolean) => {
    isOnline.value = online
  }

  const setError = (errorMessage: string | null) => {
    error.value = errorMessage
  }

  const clearError = () => {
    error.value = null
  }

  // 初始化方法
  const initialize = () => {
    // 从本地存储恢复设置
    const savedTheme = localStorage.getItem('app-theme') as 'light' | 'dark' | 'auto'
    if (savedTheme) {
      theme.value = savedTheme
    }

    const savedLanguage = localStorage.getItem('app-language')
    if (savedLanguage) {
      language.value = savedLanguage
    }

    // 监听网络状态变化
    window.addEventListener('online', () => setOnlineStatus(true))
    window.addEventListener('offline', () => setOnlineStatus(false))

    // 监听系统主题变化
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', () => {
      // 如果是自动主题，触发重新计算
      if (theme.value === 'auto') {
        // 触发响应式更新
        theme.value = 'auto'
      }
    })

    console.log('[App] 应用状态初始化完成')
  }

  // 工具方法
  const formatError = (error: any): string => {
    if (typeof error === 'string') {
      return error
    }
    
    if (error instanceof Error) {
      return error.message
    }
    
    if (error?.message) {
      return error.message
    }
    
    return '未知错误'
  }

  const handleGlobalError = (error: any) => {
    const errorMessage = formatError(error)
    console.error('[App] 全局错误:', error)
    
    setError(errorMessage)
    showError(errorMessage)
  }

  const retry = async (fn: () => Promise<any>, maxRetries = 3, delay = 1000) => {
    let lastError: any
    
    for (let i = 0; i < maxRetries; i++) {
      try {
        return await fn()
      } catch (error) {
        lastError = error
        
        if (i < maxRetries - 1) {
          await new Promise(resolve => setTimeout(resolve, delay * (i + 1)))
        }
      }
    }
    
    throw lastError
  }

  // 返回状态和方法
  return {
    // 只读状态
    isLoading: readonly(isLoading),
    notifications: readonly(notifications),
    theme: readonly(theme),
    language: readonly(language),
    isOnline: readonly(isOnline),
    error: readonly(error),

    // 计算属性
    hasNotifications,
    activeNotifications,
    currentTheme,

    // 操作方法
    setLoading,
    showNotification,
    removeNotification,
    clearNotifications,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    setTheme,
    setLanguage,
    setOnlineStatus,
    setError,
    clearError,
    initialize,
    handleGlobalError,
    retry,
    formatError
  }
})