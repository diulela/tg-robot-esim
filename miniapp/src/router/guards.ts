// 路由守卫配置
import type { Router, NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { telegramService } from '@/services/telegram'

// 设置路由守卫
export function setupRouterGuards(router: Router) {
  // 全局前置守卫
  router.beforeEach(async (to, from, next) => {
    console.log(`[Router] 导航到: ${to.name} (${to.path})`)
    
    try {
      // 执行各种守卫检查
      await executeGuards(to, from, next)
    } catch (error) {
      console.error('[Router] 路由守卫执行失败:', error)
      next({ name: 'Error', query: { error: '页面加载失败' } })
    }
  })

  // 全局后置钩子
  router.afterEach((to, from) => {
    console.log(`[Router] 导航完成: ${from.name} -> ${to.name}`)
    
    // 更新页面状态
    updatePageState(to)
    
    // 更新 Telegram 按钮状态
    updateTelegramButtons(to)
  })

  // 路由错误处理
  router.onError((error) => {
    console.error('[Router] 路由错误:', error)
    router.push({ name: 'Error', query: { error: '路由加载失败' } })
  })
}

// 执行守卫检查
async function executeGuards(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  // 1. 认证守卫
  if (!(await authGuard(to, from, next))) return

  // 2. 权限守卫
  if (!(await permissionGuard(to, from, next))) return

  // 所有守卫通过，继续导航
  next()
}

// 认证守卫
async function authGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
): Promise<boolean> {
  // 检查路由是否需要认证
  if (!to.meta.requiresAuth) {
    return true
  }

  const userStore = useUserStore()
  
  // 如果用户未初始化，尝试初始化
  if (!userStore.isAuthenticated) {
    try {
      await userStore.initializeUser()
    } catch (error) {
      console.warn('[Router] 用户认证失败:', error)
      
      // 显示认证失败提示
      const appStore = useAppStore()
      appStore.showNotification({
        type: 'error',
        message: '用户认证失败，请重新打开应用',
        persistent: true
      })
      
      // 重定向到首页
      next({ name: 'Home' })
      return false
    }
  }

  // 检查认证状态
  if (!userStore.isAuthenticated) {
    console.warn('[Router] 用户未认证，重定向到首页')
    
    const appStore = useAppStore()
    appStore.showNotification({
      type: 'warning',
      message: '请先完成用户认证',
      duration: 3000
    })
    
    next({ name: 'Home' })
    return false
  }

  return true
}

// 权限守卫
async function permissionGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
): Promise<boolean> {
  const userStore = useUserStore()
  
  // 检查用户是否有访问权限
  // 这里可以根据用户角色、权限等进行更复杂的权限检查
  // 目前简化处理，认证用户都有基本权限
  
  if (to.meta.requiresAuth && !userStore.hasPermission('basic')) {
    console.warn('[Router] 用户权限不足')
    
    const appStore = useAppStore()
    appStore.showNotification({
      type: 'error',
      message: '您没有访问此页面的权限',
      duration: 3000
    })
    
    next({ name: 'Home' })
    return false
  }

  return true
}



// 更新页面状态
function updatePageState(to: RouteLocationNormalized) {
  

  // 更新页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - eSIM Mini App`
  }

}

// 更新 Telegram 按钮状态
function updateTelegramButtons(to: RouteLocationNormalized) {
  try {
    // 更新后退按钮
    if (to.meta.showBackButton) {
      telegramService.showBackButton()
    } else {
      telegramService.hideBackButton()
    }
    
    // 隐藏主按钮（页面组件会根据需要显示）
    telegramService.hideMainButton()
    
    // 触觉反馈
    telegramService.selectionFeedback()
  } catch (error) {
    console.warn('[Router] 更新 Telegram 按钮状态失败:', error)
  }
}


// 导航辅助函数
export function createNavigationHelper(router: Router) {
  return {
    // 安全导航 - 带错误处理
    async safePush(to: any) {
      try {
        await router.push(to)
      } catch (error) {
        console.error('[Router] 导航失败:', error)
        router.push({ name: 'Error', query: { error: '页面跳转失败' } })
      }
    },
    
    // 安全替换 - 带错误处理
    async safeReplace(to: any) {
      try {
        await router.replace(to)
      } catch (error) {
        console.error('[Router] 页面替换失败:', error)
        router.push({ name: 'Error', query: { error: '页面跳转失败' } })
      }
    },
    
    // 安全后退
    safeBack() {
      try {
        if (window.history.length > 1) {
          router.back()
        } else {
          router.push({ name: 'Home' })
        }
      } catch (error) {
        console.error('[Router] 后退失败:', error)
        router.push({ name: 'Home' })
      }
    },
    
    // 检查是否可以后退
    canGoBack(): boolean {
      return window.history.length > 1
    },
    
    // 获取当前路由信息
    getCurrentRoute() {
      return router.currentRoute.value
    },
    
    // 检查是否在指定路由
    isCurrentRoute(name: string): boolean {
      return router.currentRoute.value.name === name
    }
  }
}