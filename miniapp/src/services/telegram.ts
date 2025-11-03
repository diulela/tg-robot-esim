// Telegram Web App 服务实现
import type {
  TelegramWebApp,
  TelegramUser,
  TelegramThemeParams,
  TelegramAuthData,
  TelegramEventType,
  TelegramEventHandler
} from '@/types'

// Telegram 服务配置
interface TelegramServiceConfig {
  enableMock: boolean
  mockUser?: TelegramUser
  debug: boolean
}

// 默认配置
const defaultConfig: TelegramServiceConfig = {
  enableMock: import.meta.env.VITE_ENABLE_MOCK === 'true' || import.meta.env.NODE_ENV === 'development',
  debug: import.meta.env.VITE_ENABLE_DEBUG === 'true' || import.meta.env.NODE_ENV === 'development',
  mockUser: {
    id: 123456789,
    first_name: '测试用户',
    last_name: 'Test',
    username: 'testuser',
    language_code: 'zh-cn',
    is_premium: false
  }
}

class TelegramService {
  private webApp: TelegramWebApp | null = null
  private config: TelegramServiceConfig
  private eventHandlers: Map<string, Set<TelegramEventHandler>> = new Map()
  private isInitialized = false

  constructor(config: Partial<TelegramServiceConfig> = {}) {
    this.config = { ...defaultConfig, ...config }
    this.init()
  }

  private init(): void {
    if (this.isInitialized) return

    try {
      // 检查是否在 Telegram 环境中
      if (typeof window !== 'undefined' && window.Telegram?.WebApp) {
        this.webApp = window.Telegram.WebApp
        this.setupTelegramWebApp()
        this.log('Telegram Web App initialized successfully')
      } else if (this.config.enableMock) {
        this.webApp = this.createMockTelegramWebApp()
        this.log('Mock Telegram Web App initialized')
      } else {
        throw new Error('Telegram Web App not available')
      }

      this.isInitialized = true
    } catch (error) {
      console.error('[Telegram] Initialization failed:', error)
      throw error
    }
  }

  private setupTelegramWebApp(): void {
    if (!this.webApp) return

    // 准备 Web App
    this.webApp.ready()
    
    // 展开 Web App
    this.webApp.expand()

    // 设置事件监听
    this.setupEventListeners()

    this.log('Telegram Web App setup completed')
  }

  private setupEventListeners(): void {
    if (!this.webApp) return

    // 主题变化事件
    this.webApp.onEvent('themeChanged', () => {
      this.emit('themeChanged', this.getThemeParams())
    })

    // 视窗变化事件
    this.webApp.onEvent('viewportChanged', () => {
      this.emit('viewportChanged', {
        height: this.webApp!.viewportHeight,
        stableHeight: this.webApp!.viewportStableHeight,
        isExpanded: this.webApp!.isExpanded
      })
    })

    // 后退按钮点击事件
    this.webApp.BackButton.onClick(() => {
      this.emit('backButtonClicked', {})
    })

    // 主按钮点击事件
    this.webApp.MainButton.onClick(() => {
      this.emit('mainButtonClicked', {})
    })
  }

  private createMockTelegramWebApp(): TelegramWebApp {
    const mockThemeParams: TelegramThemeParams = {
      bg_color: '#ffffff',
      text_color: '#000000',
      hint_color: '#999999',
      link_color: '#6366f1',
      button_color: '#6366f1',
      button_text_color: '#ffffff',
      secondary_bg_color: '#f8fafc'
    }

    return {
      initData: 'mock_init_data',
      initDataUnsafe: {
        user: this.config.mockUser,
        auth_date: Math.floor(Date.now() / 1000),
        hash: 'mock_hash'
      },
      version: '6.0',
      platform: 'web',
      colorScheme: 'light',
      themeParams: mockThemeParams,
      isExpanded: true,
      viewportHeight: window.innerHeight,
      viewportStableHeight: window.innerHeight,
      headerColor: mockThemeParams.bg_color || '#ffffff',
      backgroundColor: mockThemeParams.bg_color || '#ffffff',
      BackButton: {
        isVisible: false,
        onClick: (callback: () => void) => {
          // Mock implementation
        },
        offClick: (callback: () => void) => {
          // Mock implementation
        },
        show: () => {
          this.log('BackButton.show()')
        },
        hide: () => {
          this.log('BackButton.hide()')
        }
      },
      MainButton: {
        text: '',
        color: mockThemeParams.button_color || '#6366f1',
        textColor: mockThemeParams.button_text_color || '#ffffff',
        isVisible: false,
        isActive: true,
        isProgressVisible: false,
        setText: (text: string) => {
          this.log(`MainButton.setText("${text}")`)
        },
        onClick: (callback: () => void) => {
          // Mock implementation
        },
        offClick: (callback: () => void) => {
          // Mock implementation
        },
        show: () => {
          this.log('MainButton.show()')
        },
        hide: () => {
          this.log('MainButton.hide()')
        },
        enable: () => {
          this.log('MainButton.enable()')
        },
        disable: () => {
          this.log('MainButton.disable()')
        },
        showProgress: (leaveActive?: boolean) => {
          this.log(`MainButton.showProgress(${leaveActive})`)
        },
        hideProgress: () => {
          this.log('MainButton.hideProgress()')
        },
        setParams: (params: any) => {
          this.log('MainButton.setParams()', params)
        }
      },
      HapticFeedback: {
        impactOccurred: (style: any) => {
          this.log(`HapticFeedback.impactOccurred("${style}")`)
        },
        notificationOccurred: (type: any) => {
          this.log(`HapticFeedback.notificationOccurred("${type}")`)
        },
        selectionChanged: () => {
          this.log('HapticFeedback.selectionChanged()')
        }
      },
      ready: () => {
        this.log('WebApp.ready()')
      },
      expand: () => {
        this.log('WebApp.expand()')
      },
      close: () => {
        this.log('WebApp.close()')
      },
      sendData: (data: string) => {
        this.log(`WebApp.sendData("${data}")`)
      },
      switchInlineQuery: (query: string, choose_chat_types?: string[]) => {
        this.log(`WebApp.switchInlineQuery("${query}", ${JSON.stringify(choose_chat_types)})`)
      },
      openLink: (url: string, options?: any) => {
        this.log(`WebApp.openLink("${url}", ${JSON.stringify(options)})`)
        window.open(url, '_blank')
      },
      openTelegramLink: (url: string) => {
        this.log(`WebApp.openTelegramLink("${url}")`)
      },
      openInvoice: (url: string, callback?: any) => {
        this.log(`WebApp.openInvoice("${url}")`)
      },
      showPopup: (params: any, callback?: any) => {
        this.log('WebApp.showPopup()', params)
        if (callback) callback('ok')
      },
      showAlert: (message: string, callback?: any) => {
        this.log(`WebApp.showAlert("${message}")`)
        alert(message)
        if (callback) callback()
      },
      showConfirm: (message: string, callback?: any) => {
        this.log(`WebApp.showConfirm("${message}")`)
        const result = confirm(message)
        if (callback) callback(result)
      },
      showScanQrPopup: (params: any, callback?: any) => {
        this.log('WebApp.showScanQrPopup()', params)
      },
      closeScanQrPopup: () => {
        this.log('WebApp.closeScanQrPopup()')
      },
      readTextFromClipboard: (callback?: any) => {
        this.log('WebApp.readTextFromClipboard()')
        if (callback) callback(null)
      },
      requestWriteAccess: (callback?: any) => {
        this.log('WebApp.requestWriteAccess()')
        if (callback) callback(true)
      },
      requestContact: (callback?: any) => {
        this.log('WebApp.requestContact()')
        if (callback) callback(true)
      },
      onEvent: (eventType: string, eventHandler: any) => {
        this.log(`WebApp.onEvent("${eventType}")`)
      },
      offEvent: (eventType: string, eventHandler: any) => {
        this.log(`WebApp.offEvent("${eventType}")`)
      }
    }
  }

  private log(message: string, ...args: any[]): void {
    if (this.config.debug) {
      console.log(`[Telegram] ${message}`, ...args)
    }
  }

  // 事件系统
  on<T = any>(event: TelegramEventType, handler: TelegramEventHandler<T>): void {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, new Set())
    }
    this.eventHandlers.get(event)!.add(handler)
  }

  off<T = any>(event: TelegramEventType, handler: TelegramEventHandler<T>): void {
    const handlers = this.eventHandlers.get(event)
    if (handlers) {
      handlers.delete(handler)
    }
  }

  emit<T = any>(event: TelegramEventType, data: T): void {
    const handlers = this.eventHandlers.get(event)
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(data)
        } catch (error) {
          console.error(`[Telegram] Event handler error for "${event}":`, error)
        }
      })
    }
  }

  // 公共方法
  isAvailable(): boolean {
    return this.webApp !== null
  }

  getWebApp(): TelegramWebApp | null {
    return this.webApp
  }

  getUser(): TelegramUser | null {
    return this.webApp?.initDataUnsafe.user || null
  }

  getInitData(): string {
    return this.webApp?.initData || ''
  }

  getThemeParams(): TelegramThemeParams {
    return this.webApp?.themeParams || {}
  }

  getColorScheme(): 'light' | 'dark' {
    return this.webApp?.colorScheme || 'light'
  }

  getViewportHeight(): number {
    return this.webApp?.viewportHeight || window.innerHeight
  }

  isExpanded(): boolean {
    return this.webApp?.isExpanded || false
  }

  // 主题相关方法
  isDarkMode(): boolean {
    return this.getColorScheme() === 'dark'
  }

  getPrimaryColor(): string {
    const themeParams = this.getThemeParams()
    return themeParams.button_color || '#6366f1'
  }

  getBackgroundColor(): string {
    const themeParams = this.getThemeParams()
    return themeParams.bg_color || '#ffffff'
  }

  getTextColor(): string {
    const themeParams = this.getThemeParams()
    return themeParams.text_color || '#000000'
  }

  // 按钮控制方法
  showBackButton(): void {
    this.webApp?.BackButton.show()
  }

  hideBackButton(): void {
    this.webApp?.BackButton.hide()
  }

  setMainButton(text: string, callback?: () => void): void {
    if (!this.webApp) return

    this.webApp.MainButton.setText(text)
    this.webApp.MainButton.show()
    
    if (callback) {
      this.webApp.MainButton.onClick(callback)
    }
  }

  hideMainButton(): void {
    this.webApp?.MainButton.hide()
  }

  showMainButtonProgress(): void {
    this.webApp?.MainButton.showProgress()
  }

  hideMainButtonProgress(): void {
    this.webApp?.MainButton.hideProgress()
  }

  // 触觉反馈方法
  impactFeedback(style: 'light' | 'medium' | 'heavy' = 'medium'): void {
    this.webApp?.HapticFeedback.impactOccurred(style)
  }

  notificationFeedback(type: 'error' | 'success' | 'warning'): void {
    this.webApp?.HapticFeedback.notificationOccurred(type)
  }

  selectionFeedback(): void {
    this.webApp?.HapticFeedback.selectionChanged()
  }

  // 弹窗方法
  showAlert(message: string): Promise<void> {
    return new Promise((resolve) => {
      this.webApp?.showAlert(message, () => resolve())
    })
  }

  showConfirm(message: string): Promise<boolean> {
    return new Promise((resolve) => {
      this.webApp?.showConfirm(message, (confirmed) => resolve(confirmed))
    })
  }

  // 链接和导航方法
  openLink(url: string, tryInstantView = false): void {
    this.webApp?.openLink(url, { try_instant_view: tryInstantView })
  }

  openTelegramLink(url: string): void {
    this.webApp?.openTelegramLink(url)
  }

  close(): void {
    this.webApp?.close()
  }

  // 数据验证方法
  validateAuthData(): TelegramAuthData {
    const initData = this.getInitData()
    const user = this.getUser()
    
    // 在实际应用中，这里应该验证 hash 的有效性
    // 目前简化处理
    return {
      initData,
      hash: this.webApp?.initDataUnsafe.hash || '',
      isValid: !!initData && !!user,
      user: user || undefined,
      timestamp: this.webApp?.initDataUnsafe.auth_date || Math.floor(Date.now() / 1000)
    }
  }

  // 工具方法
  ready(): void {
    this.webApp?.ready()
  }

  expand(): void {
    this.webApp?.expand()
  }

  sendData(data: any): void {
    const jsonData = typeof data === 'string' ? data : JSON.stringify(data)
    this.webApp?.sendData(jsonData)
  }
}

// 创建 Telegram 服务实例
export const telegramService = new TelegramService()

// 导出服务类
export { TelegramService }
export default telegramService