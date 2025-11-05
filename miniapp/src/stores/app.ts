// 应用全局状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly, watch } from 'vue'
import type { Notification } from '@/types'
import { telegramService } from '@/services/telegram'

// 应用配置接口
interface AppConfig {
  name: string
  version: string
  apiBaseUrl: string
  enableMock: boolean
  enableDebug: boolean
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

// 默认配置
const defaultConfig: AppConfig = {
  name: 'eSIM Mini App',
  version: '1.0.0',
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  enableMock: import.meta.env.VITE_ENABLE_MOCK === 'true',
  enableDebug: import.meta.env.VITE_ENABLE_DEBUG === 'true',
  defaultLanguage: 'zh-cn',
  supportedLanguages: ['zh-cn', 'en', 'zh-tw'],
  theme: {
    defaultMode: 'auto',
    primaryColor: '#6366F1',
    secondaryColor: '#EC4899'
  },
  features: {
    enableNotifications: true,
    enableAnalytics: false,
    enableErrorReporting: true,
    enableOfflineMode: false
  }
}

export const useAppStore = defineStore('app', () => {
  // 状态
  const config = ref<AppConfig>(defaultConfig)
  const isInitialized = ref(false)
  const isLoading = ref(false)
  const currentTheme = ref<'light' | 'dark'>('light')
  const currentLanguage = ref('zh-cn')
  const isOnline = ref(navigator.onLine)
  const notifications = ref<Notification[]>([])
  const globalError = ref<string | null>(null)
  const debugMode = ref(false)
  const performanceMetrics = ref({
    appStartTime: Date.now(),
    lastPageLoadTime: 0,
    apiResponseTimes: [] as number[]
  })

  // 页面状态
  const currentPage = ref('')
  const pageTitle = ref('')
  const showBackButton = ref(false)
  const showBottomNav = ref(true)
  const pageLoading = ref(false)

  // 计算属性
  const isDarkMode = computed(() => currentTheme.value === 'dark')

  const isLightMode = computed(() => currentTheme.value === 'light')

  const supportedLanguageOptions = computed(() => 
    config.value.supportedLanguages.map(lang => ({
      value: lang,
      label: getLanguageLabel(lang)
    }))
  )

  const hasNotifications = computed(() => notifications.value.length > 0)

  const unreadNotifications = computed(() => 
    notifications.value.filter(n => !n.persistent)
  )

  const persistentNotifications = computed(() => 
    notifications.value.filter(n => n.persistent)
  )

  const isOffline = computed(() => !isOnline.value)

  const canUseFeature = computed(() => (feature: keyof AppConfig['features']) => 
    config.value.features[feature]
  )

  const averageApiResponseTime = computed(() => {
    const times = performanceMetrics.value.apiResponseTimes
    if (times.length === 0) return 0
    return times.reduce((sum, time) => sum + time, 0) / times.length
  })

  const appUptime = computed(() => 
    Date.now() - performanceMetrics.value.appStartTime
  )

  // 操作方法
  const initialize = async (): Promise<void> => {
    if (isInitialized.value) return

    isLoading.value = true

    try {
      // 初始化主题
      await initializeTheme()

      // 初始化语言
      await initializeLanguage()

      // 设置网络状态监听
      setupNetworkListeners()

      // 设置 Telegram 主题监听
      setupTelegramThemeListener()

      // 加载持久化配置
      await loadPersistedConfig()

      // 设置调试模式
      debugMode.value = config.value.enableDebug

      isInitialized.value = true
      console.log('[App] 应用初始化成功')
    } catch (error) {
      console.error('[App] 应用初始化失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const initializeTheme = async (): Promise<void> => {
    try {
      // 从 Telegram 获取主题
      const telegramTheme = telegramService.getColorScheme()
      
      // 从本地存储获取用户偏好
      const savedTheme = localStorage.getItem('app_theme') as 'light' | 'dark' | null
      
      if (savedTheme) {
        currentTheme.value = savedTheme
      } else if (config.value.theme.defaultMode === 'auto') {
        currentTheme.value = telegramTheme
      } else {
        currentTheme.value = config.value.theme.defaultMode as 'light' | 'dark'
      }

      // 应用主题到 DOM
      applyThemeToDOM()
    } catch (error) {
      console.warn('[App] 主题初始化失败:', error)
      currentTheme.value = 'light'
    }
  }

  const initializeLanguage = async (): Promise<void> => {
    try {
      // 从 Telegram 获取语言
      const telegramUser = telegramService.getUser()
      const telegramLang = telegramUser?.language_code

      // 从本地存储获取用户偏好
      const savedLang = localStorage.getItem('app_language')

      if (savedLang && config.value.supportedLanguages.includes(savedLang)) {
        currentLanguage.value = savedLang
      } else if (telegramLang && config.value.supportedLanguages.includes(telegramLang)) {
        currentLanguage.value = telegramLang
      } else {
        currentLanguage.value = config.value.defaultLanguage
      }
    } catch (error) {
      console.warn('[App] 语言初始化失败:', error)
      currentLanguage.value = config.value.defaultLanguage
    }
  }

  const setupNetworkListeners = (): void => {
    const updateOnlineStatus = () => {
      isOnline.value = navigator.onLine
      
      if (isOnline.value) {
        showNotification({
          type: 'success',
          message: '网络连接已恢复',
          duration: 3000
        })
      } else {
        showNotification({
          type: 'warning',
          message: '网络连接已断开',
          persistent: true
        })
      }
    }

    window.addEventListener('online', updateOnlineStatus)
    window.addEventListener('offline', updateOnlineStatus)
  }

  const setupTelegramThemeListener = (): void => {
    telegramService.on('themeChanged', (themeParams) => {
      if (config.value.theme.defaultMode === 'auto') {
        const newTheme = telegramService.getColorScheme()
        if (newTheme !== currentTheme.value) {
          setTheme(newTheme)
        }
      }
    })
  }

  const loadPersistedConfig = async (): Promise<void> => {
    try {
      const savedConfig = localStorage.getItem('app_config')
      if (savedConfig) {
        const parsedConfig = JSON.parse(savedConfig)
        config.value = { ...config.value, ...parsedConfig }
      }
    } catch (error) {
      console.warn('[App] 加载持久化配置失败:', error)
    }
  }

  const saveConfig = async (): Promise<void> => {
    try {
      localStorage.setItem('app_config', JSON.stringify(config.value))
    } catch (error) {
      console.warn('[App] 保存配置失败:', error)
    }
  }

  const setTheme = async (theme: 'light' | 'dark'): Promise<void> => {
    currentTheme.value = theme
    applyThemeToDOM()
    
    try {
      localStorage.setItem('app_theme', theme)
    } catch (error) {
      console.warn('[App] 保存主题设置失败:', error)
    }
  }

  const toggleTheme = async (): Promise<void> => {
    const newTheme = currentTheme.value === 'light' ? 'dark' : 'light'
    await setTheme(newTheme)
  }

  const setLanguage = async (language: string): Promise<void> => {
    if (!config.value.supportedLanguages.includes(language)) {
      throw new Error(`不支持的语言: ${language}`)
    }

    currentLanguage.value = language
    
    try {
      localStorage.setItem('app_language', language)
    } catch (error) {
      console.warn('[App] 保存语言设置失败:', error)
    }
  }

  const applyThemeToDOM = (): void => {
    document.documentElement.setAttribute('data-theme', currentTheme.value)
    
    // 更新 meta 标签
    const themeColorMeta = document.querySelector('meta[name="theme-color"]')
    if (themeColorMeta) {
      const color = currentTheme.value === 'dark' ? '#0F172A' : '#FFFFFF'
      themeColorMeta.setAttribute('content', color)
    }
  }

  // 通知管理
  const showNotification = (notification: Omit<Notification, 'id' | 'createdAt'>): string => {
    const id = `notification_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    const newNotification: Notification = {
      id,
      createdAt: new Date(),
      ...notification
    }

    notifications.value.push(newNotification)

    // 自动移除非持久化通知
    if (!notification.persistent && notification.duration !== 0) {
      const duration = notification.duration || 5000
      setTimeout(() => {
        removeNotification(id)
      }, duration)
    }

    return id
  }

  const removeNotification = (id: string): void => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  const clearNotifications = (): void => {
    notifications.value = []
  }

  const clearPersistentNotifications = (): void => {
    notifications.value = notifications.value.filter(n => !n.persistent)
  }

  // 页面状态管理
  const setCurrentPage = (page: string, title?: string): void => {
    currentPage.value = page
    if (title) {
      pageTitle.value = title
      document.title = `${title} - ${config.value.name}`
    }
  }

  const setPageLoading = (loading: boolean): void => {
    pageLoading.value = loading
  }

  const setBackButton = (show: boolean): void => {
    showBackButton.value = show
    
    if (show) {
      telegramService.showBackButton()
    } else {
      telegramService.hideBackButton()
    }
  }

  const setBottomNav = (show: boolean): void => {
    showBottomNav.value = show
  }

  // 错误处理
  const setGlobalError = (error: string | null): void => {
    globalError.value = error
    
    if (error) {
      showNotification({
        type: 'error',
        message: error,
        persistent: true
      })
    }
  }

  const clearGlobalError = (): void => {
    globalError.value = null
  }

  // 性能监控
  const recordPageLoadTime = (): void => {
    performanceMetrics.value.lastPageLoadTime = Date.now()
  }

  const recordApiResponseTime = (responseTime: number): void => {
    performanceMetrics.value.apiResponseTimes.push(responseTime)
    
    // 只保留最近 100 次记录
    if (performanceMetrics.value.apiResponseTimes.length > 100) {
      performanceMetrics.value.apiResponseTimes.shift()
    }
  }

  // 便捷通知方法
  const showSuccess = (message: string, duration = 3000): string => {
    return showNotification({
      type: 'success',
      message,
      duration
    })
  }

  const showError = (message: string, duration = 5000): string => {
    return showNotification({
      type: 'error',
      message,
      duration
    })
  }

  const showWarning = (message: string, duration = 4000): string => {
    return showNotification({
      type: 'warning',
      message,
      duration
    })
  }

  const showInfo = (message: string, duration = 3000): string => {
    return showNotification({
      type: 'info',
      message,
      duration
    })
  }

  // 工具方法
  const getLanguageLabel = (langCode: string): string => {
    const labels: Record<string, string> = {
      'zh-cn': '简体中文',
      'zh-tw': '繁體中文',
      'en': 'English'
    }
    return labels[langCode] || langCode
  }

  const formatUptime = (): string => {
    const uptime = appUptime.value
    const hours = Math.floor(uptime / (1000 * 60 * 60))
    const minutes = Math.floor((uptime % (1000 * 60 * 60)) / (1000 * 60))
    const seconds = Math.floor((uptime % (1000 * 60)) / 1000)
    
    if (hours > 0) {
      return `${hours}小时${minutes}分钟`
    } else if (minutes > 0) {
      return `${minutes}分钟${seconds}秒`
    } else {
      return `${seconds}秒`
    }
  }

  const getAppInfo = () => {
    return {
      name: config.value.name,
      version: config.value.version,
      theme: currentTheme.value,
      language: currentLanguage.value,
      isOnline: isOnline.value,
      uptime: formatUptime(),
      averageApiResponseTime: Math.round(averageApiResponseTime.value),
      debugMode: debugMode.value
    }
  }

  // 监听主题变化
  watch(currentTheme, (newTheme) => {
    console.log('[App] 主题已切换:', newTheme)
  })

  // 监听语言变化
  watch(currentLanguage, (newLang) => {
    console.log('[App] 语言已切换:', newLang)
  })

  // 返回状态和方法
  return {
    // 只读状态
    config: readonly(config),
    isInitialized: readonly(isInitialized),
    isLoading: readonly(isLoading),
    currentTheme: readonly(currentTheme),
    currentLanguage: readonly(currentLanguage),
    isOnline: readonly(isOnline),
    notifications: readonly(notifications),
    globalError: readonly(globalError),
    debugMode: readonly(debugMode),
    performanceMetrics: readonly(performanceMetrics),
    currentPage: readonly(currentPage),
    pageTitle: readonly(pageTitle),
    showBackButton: readonly(showBackButton),
    showBottomNav: readonly(showBottomNav),
    pageLoading: readonly(pageLoading),

    // 计算属性
    isDarkMode,
    isLightMode,
    supportedLanguageOptions,
    hasNotifications,
    unreadNotifications,
    persistentNotifications,
    isOffline,
    canUseFeature,
    averageApiResponseTime,
    appUptime,

    // 操作方法
    initialize,
    setTheme,
    toggleTheme,
    setLanguage,
    showNotification,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    removeNotification,
    clearNotifications,
    clearPersistentNotifications,
    setCurrentPage,
    setPageLoading,
    setBackButton,
    setBottomNav,
    setGlobalError,
    clearGlobalError,
    recordPageLoadTime,
    recordApiResponseTime,
    saveConfig,

    // 工具方法
    getLanguageLabel,
    formatUptime,
    getAppInfo
  }
})