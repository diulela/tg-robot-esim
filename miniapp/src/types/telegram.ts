// Telegram Web App 集成类型定义

// Telegram 用户类型
export interface TelegramUser {
  id: number
  first_name: string
  last_name?: string
  username?: string
  language_code?: string
  is_premium?: boolean
  photo_url?: string
  allows_write_to_pm?: boolean
}

// Telegram 聊天类型
export interface TelegramChat {
  id: number
  type: 'private' | 'group' | 'supergroup' | 'channel'
  title?: string
  username?: string
  photo_url?: string
}

// Telegram 初始化数据
export interface TelegramInitData {
  user?: TelegramUser
  receiver?: TelegramUser
  chat?: TelegramChat
  start_param?: string
  can_send_after?: number
  auth_date: number
  hash: string
}

// Telegram 主题参数
export interface TelegramThemeParams {
  bg_color?: string
  text_color?: string
  hint_color?: string
  link_color?: string
  button_color?: string
  button_text_color?: string
  secondary_bg_color?: string
  header_bg_color?: string
  accent_text_color?: string
  section_bg_color?: string
  section_header_text_color?: string
  subtitle_text_color?: string
  destructive_text_color?: string
}

// Telegram 后退按钮
export interface TelegramBackButton {
  isVisible: boolean
  onClick: (callback: () => void) => void
  offClick: (callback: () => void) => void
  show: () => void
  hide: () => void
}

// Telegram 主按钮参数
export interface TelegramMainButtonParams {
  text?: string
  color?: string
  text_color?: string
  is_active?: boolean
  is_visible?: boolean
}

// Telegram 主按钮
export interface TelegramMainButton {
  text: string
  color: string
  textColor: string
  isVisible: boolean
  isActive: boolean
  isProgressVisible: boolean
  setText: (text: string) => void
  onClick: (callback: () => void) => void
  offClick: (callback: () => void) => void
  show: () => void
  hide: () => void
  enable: () => void
  disable: () => void
  showProgress: (leaveActive?: boolean) => void
  hideProgress: () => void
  setParams: (params: TelegramMainButtonParams) => void
}

// Telegram 触觉反馈
export interface TelegramHapticFeedback {
  impactOccurred: (style: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft') => void
  notificationOccurred: (type: 'error' | 'success' | 'warning') => void
  selectionChanged: () => void
}

// Telegram 弹窗按钮
export interface TelegramPopupButton {
  id?: string
  type?: 'default' | 'ok' | 'close' | 'cancel' | 'destructive'
  text?: string
}

// Telegram 弹窗参数
export interface TelegramPopupParams {
  title?: string
  message: string
  buttons?: TelegramPopupButton[]
}

// Telegram 扫码结果
export interface TelegramScanQrResult {
  data: string
}

// Telegram Web App 主接口
export interface TelegramWebApp {
  // 基础属性
  initData: string
  initDataUnsafe: TelegramInitData
  version: string
  platform: string
  colorScheme: 'light' | 'dark'
  themeParams: TelegramThemeParams
  isExpanded: boolean
  viewportHeight: number
  viewportStableHeight: number
  headerColor: string
  backgroundColor: string
  
  // 组件
  BackButton: TelegramBackButton
  MainButton: TelegramMainButton
  HapticFeedback: TelegramHapticFeedback
  
  // 方法
  ready: () => void
  expand: () => void
  close: () => void
  sendData: (data: string) => void
  switchInlineQuery: (query: string, choose_chat_types?: string[]) => void
  openLink: (url: string, options?: { try_instant_view?: boolean }) => void
  openTelegramLink: (url: string) => void
  openInvoice: (url: string, callback?: (status: string) => void) => void
  showPopup: (params: TelegramPopupParams, callback?: (button_id: string) => void) => void
  showAlert: (message: string, callback?: () => void) => void
  showConfirm: (message: string, callback?: (confirmed: boolean) => void) => void
  showScanQrPopup: (params: { text?: string }, callback?: (result: TelegramScanQrResult | null) => void) => void
  closeScanQrPopup: () => void
  readTextFromClipboard: (callback?: (text: string | null) => void) => void
  requestWriteAccess: (callback?: (granted: boolean) => void) => void
  requestContact: (callback?: (granted: boolean) => void) => void
  
  // 事件监听
  onEvent: (eventType: string, eventHandler: (...args: any[]) => void) => void
  offEvent: (eventType: string, eventHandler: (...args: any[]) => void) => void
}

// Telegram Web App 事件类型
export type TelegramEventType = 
  | 'themeChanged'
  | 'viewportChanged'
  | 'mainButtonClicked'
  | 'backButtonClicked'
  | 'settingsButtonClicked'
  | 'invoiceClosed'
  | 'popupClosed'
  | 'qrTextReceived'
  | 'clipboardTextReceived'
  | 'writeAccessRequested'
  | 'contactRequested'

// Telegram Web App 事件处理器
export type TelegramEventHandler<T = any> = (data: T) => void

// Telegram 验证数据
export interface TelegramAuthData {
  initData: string
  hash: string
  isValid: boolean
  user?: TelegramUser
  timestamp: number
}

// Telegram Bot 命令
export interface TelegramBotCommand {
  command: string
  description: string
}

// Telegram 内联键盘按钮
export interface TelegramInlineKeyboardButton {
  text: string
  url?: string
  callback_data?: string
  web_app?: {
    url: string
  }
  switch_inline_query?: string
  switch_inline_query_current_chat?: string
}

// Telegram 内联键盘
export interface TelegramInlineKeyboard {
  inline_keyboard: TelegramInlineKeyboardButton[][]
}

// Telegram 消息类型
export interface TelegramMessage {
  message_id: number
  from?: TelegramUser
  chat: TelegramChat
  date: number
  text?: string
  photo?: any[]
  document?: any
  reply_markup?: TelegramInlineKeyboard
}

// Telegram Bot API 响应
export interface TelegramBotResponse<T = any> {
  ok: boolean
  result?: T
  error_code?: number
  description?: string
}

// 全局 Window 接口扩展
declare global {
  interface Window {
    Telegram?: {
      WebApp: TelegramWebApp
    }
  }
}