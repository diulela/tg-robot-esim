// Pinia 状态管理配置和导出
import { createPinia } from 'pinia'
import type { App } from 'vue'

// 创建 Pinia 实例
export const pinia = createPinia()

// 安装插件函数
export function setupStore(app: App) {
  app.use(pinia)
}

// 导出所有 store
export { useAppStore } from './app'
export { useUserStore } from './user'
export { useOrdersStore } from './orders'
export { useProductsStore } from './products'

// 导出类型
export type * from './app'
export type * from './user'
export type * from './orders'
export type * from './products'

// 持久化插件 (简单实现)
export function createPersistedState() {
  return ({ store }: { store: any }) => {
    // 从 localStorage 恢复状态
    const storageKey = `pinia_${store.$id}`
    const savedState = localStorage.getItem(storageKey)
    
    if (savedState) {
      try {
        const parsedState = JSON.parse(savedState)
        store.$patch(parsedState)
      } catch (error) {
        console.warn(`Failed to restore state for store ${store.$id}:`, error)
      }
    }

    // 监听状态变化并保存到 localStorage
    store.$subscribe((mutation: any, state: any) => {
      try {
        // 只保存需要持久化的状态
        const persistedState = getPersistableState(store.$id, state)
        if (persistedState) {
          localStorage.setItem(storageKey, JSON.stringify(persistedState))
        }
      } catch (error) {
        console.warn(`Failed to persist state for store ${store.$id}:`, error)
      }
    })
  }
}

// 获取需要持久化的状态
function getPersistableState(storeId: string, state: any) {
  switch (storeId) {
    case 'app':
      return {
        currentTheme: state.currentTheme,
        currentLanguage: state.currentLanguage,
        config: state.config
      }
    
    case 'user':
      return {
        preferences: state.preferences
      }
    
    case 'products':
      return {
        selectedRegion: state.selectedRegion,
        selectedCountry: state.selectedCountry,
        filters: state.filters
      }
    
    case 'orders':
      return {
        filters: state.filters
      }
    
    default:
      return null
  }
}

// Store 工具函数
export function resetAllStores() {
  const stores = [
    'app',
    'user', 
    'orders',
    'products'
  ]
  
  stores.forEach(storeId => {
    try {
      localStorage.removeItem(`pinia_${storeId}`)
    } catch (error) {
      console.warn(`Failed to clear persisted state for store ${storeId}:`, error)
    }
  })
}

export function getStoreState(storeId: string) {
  try {
    const savedState = localStorage.getItem(`pinia_${storeId}`)
    return savedState ? JSON.parse(savedState) : null
  } catch (error) {
    console.warn(`Failed to get state for store ${storeId}:`, error)
    return null
  }
}

export function setStoreState(storeId: string, state: any) {
  try {
    localStorage.setItem(`pinia_${storeId}`, JSON.stringify(state))
  } catch (error) {
    console.warn(`Failed to set state for store ${storeId}:`, error)
  }
}

// 开发工具
export function logAllStores() {
  if (import.meta.env.NODE_ENV === 'development') {
    const stores = [
      'app',
      'user',
      'orders', 
      'products'
    ]
    
    console.group('🏪 Pinia Stores State')
    stores.forEach(storeId => {
      const state = getStoreState(storeId)
      if (state) {
        console.log(`${storeId}:`, state)
      }
    })
    console.groupEnd()
  }
}

// 默认导出
export default pinia