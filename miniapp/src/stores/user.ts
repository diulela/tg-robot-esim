// 用户状态管理
import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { TelegramUser, TelegramAuthData } from '@/types'
import { telegramService } from '@/services/telegram'
import { userApi } from '@/services/api'

// 用户信息接口
interface UserProfile {
  id: string
  telegramId: number
  firstName: string
  lastName?: string
  username?: string
  languageCode?: string
  isPremium?: boolean
  photoUrl?: string
  createdAt: string
  updatedAt: string
}

// 用户偏好设置
interface UserPreferences {
  language: string
  currency: string
  theme: 'light' | 'dark' | 'auto'
  notifications: {
    orderUpdates: boolean
    promotions: boolean
    systemMessages: boolean
  }
  privacy: {
    shareProfile: boolean
    allowAnalytics: boolean
  }
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const telegramUser = ref<TelegramUser | null>(null)
  const userProfile = ref<UserProfile | null>(null)
  const preferences = ref<UserPreferences>({
    language: 'zh-cn',
    currency: 'USD',
    theme: 'auto',
    notifications: {
      orderUpdates: true,
      promotions: true,
      systemMessages: true
    },
    privacy: {
      shareProfile: false,
      allowAnalytics: true
    }
  })
  const authData = ref<TelegramAuthData | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => {
    return !!(telegramUser.value && authData.value?.isValid)
  })

  const displayName = computed(() => {
    if (!telegramUser.value) return '未知用户'
    
    const { first_name, last_name, username } = telegramUser.value
    if (first_name && last_name) {
      return `${first_name} ${last_name}`
    }
    return first_name || username || '未知用户'
  })

  const avatarUrl = computed(() => {
    return telegramUser.value?.photo_url || null
  })

  const isPremium = computed(() => {
    return telegramUser.value?.is_premium || false
  })

  const currentLanguage = computed(() => {
    return preferences.value.language || telegramUser.value?.language_code || 'zh-cn'
  })

  const currentTheme = computed(() => {
    if (preferences.value.theme === 'auto') {
      return telegramService.getColorScheme()
    }
    return preferences.value.theme
  })

  // 操作方法
  const initializeUser = async (): Promise<void> => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      // 获取 Telegram 用户信息
      const tgUser = telegramService.getUser()
      if (!tgUser) {
        throw new Error('无法获取 Telegram 用户信息')
      }

      // 验证 Telegram 数据
      const authResult = telegramService.validateAuthData()
      if (!authResult.isValid) {
        throw new Error('Telegram 数据验证失败')
      }

      // 设置用户信息
      telegramUser.value = tgUser
      authData.value = authResult

      // 从本地存储加载偏好设置
      await loadPreferences()

      // 获取用户资料
      try {
        userProfile.value = await userApi.getProfile()
      } catch (profileError) {
        console.warn('获取用户资料失败:', profileError)
        // 用户资料获取失败不影响登录流程
      }

      console.log('[User] 用户初始化成功:', {
        telegramId: tgUser.id,
        name: displayName.value,
        isAuthenticated: isAuthenticated.value
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '用户初始化失败'
      error.value = errorMessage
      console.error('[User] 用户初始化失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const updateProfile = async (updates: Partial<Pick<UserProfile, 'firstName' | 'lastName'>>): Promise<void> => {
    if (!isAuthenticated.value) {
      throw new Error('用户未认证')
    }

    isLoading.value = true
    error.value = null

    try {
      await userApi.updateProfile(updates)
      
      // 更新本地状态
      if (userProfile.value) {
        userProfile.value = {
          ...userProfile.value,
          ...updates,
          updatedAt: new Date().toISOString()
        }
      }

      console.log('[User] 用户资料更新成功')
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '更新用户资料失败'
      error.value = errorMessage
      console.error('[User] 更新用户资料失败:', err)
      throw new Error(errorMessage)
    } finally {
      isLoading.value = false
    }
  }

  const updatePreferences = async (updates: Partial<UserPreferences>): Promise<void> => {
    try {
      // 更新偏好设置
      preferences.value = {
        ...preferences.value,
        ...updates
      }

      // 保存到本地存储
      await savePreferences()

      console.log('[User] 用户偏好设置更新成功')
    } catch (err) {
      console.error('[User] 更新用户偏好设置失败:', err)
      throw err
    }
  }

  const setLanguage = async (language: string): Promise<void> => {
    await updatePreferences({ language })
  }

  const setCurrency = async (currency: string): Promise<void> => {
    await updatePreferences({ currency })
  }

  const setTheme = async (theme: 'light' | 'dark' | 'auto'): Promise<void> => {
    await updatePreferences({ theme })
  }

  const updateNotificationSettings = async (notifications: Partial<UserPreferences['notifications']>): Promise<void> => {
    await updatePreferences({
      notifications: {
        ...preferences.value.notifications,
        ...notifications
      }
    })
  }

  const updatePrivacySettings = async (privacy: Partial<UserPreferences['privacy']>): Promise<void> => {
    await updatePreferences({
      privacy: {
        ...preferences.value.privacy,
        ...privacy
      }
    })
  }

  const logout = async (): Promise<void> => {
    try {
      // 清除状态
      telegramUser.value = null
      userProfile.value = null
      authData.value = null
      error.value = null

      // 清除本地存储
      await clearStoredData()

      console.log('[User] 用户登出成功')
    } catch (err) {
      console.error('[User] 登出失败:', err)
    }
  }

  const refreshAuthData = (): void => {
    if (telegramUser.value) {
      authData.value = telegramService.validateAuthData()
    }
  }

  // 本地存储相关方法
  const loadPreferences = async (): Promise<void> => {
    try {
      const stored = localStorage.getItem('user_preferences')
      if (stored) {
        const parsedPreferences = JSON.parse(stored)
        preferences.value = {
          ...preferences.value,
          ...parsedPreferences
        }
      }
    } catch (err) {
      console.warn('[User] 加载用户偏好设置失败:', err)
    }
  }

  const savePreferences = async (): Promise<void> => {
    try {
      localStorage.setItem('user_preferences', JSON.stringify(preferences.value))
    } catch (err) {
      console.warn('[User] 保存用户偏好设置失败:', err)
    }
  }

  const clearStoredData = async (): Promise<void> => {
    try {
      localStorage.removeItem('user_preferences')
      localStorage.removeItem('user_profile')
    } catch (err) {
      console.warn('[User] 清除本地数据失败:', err)
    }
  }

  // 工具方法
  const hasPermission = (permission: string): boolean => {
    // 这里可以根据用户角色和权限进行判断
    // 目前简化处理，所有认证用户都有基本权限
    return isAuthenticated.value
  }

  const getGreeting = (): string => {
    const hour = new Date().getHours()
    const name = displayName.value

    if (hour < 6) {
      return `夜深了，${name}`
    } else if (hour < 12) {
      return `早上好，${name}`
    } else if (hour < 18) {
      return `下午好，${name}`
    } else {
      return `晚上好，${name}`
    }
  }

  const formatUserInfo = () => {
    return {
      id: userProfile.value?.id,
      telegramId: telegramUser.value?.id,
      name: displayName.value,
      username: telegramUser.value?.username,
      language: currentLanguage.value,
      theme: currentTheme.value,
      isPremium: isPremium.value,
      isAuthenticated: isAuthenticated.value
    }
  }

  // 返回状态和方法
  return {
    // 只读状态
    telegramUser: readonly(telegramUser),
    userProfile: readonly(userProfile),
    preferences: readonly(preferences),
    authData: readonly(authData),
    isLoading: readonly(isLoading),
    error: readonly(error),

    // 计算属性
    isAuthenticated,
    displayName,
    avatarUrl,
    isPremium,
    currentLanguage,
    currentTheme,

    // 操作方法
    initializeUser,
    updateProfile,
    updatePreferences,
    setLanguage,
    setCurrency,
    setTheme,
    updateNotificationSettings,
    updatePrivacySettings,
    logout,
    refreshAuthData,

    // 工具方法
    hasPermission,
    getGreeting,
    formatUserInfo
  }
})