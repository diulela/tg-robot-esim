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
    
    // 记录页面加载时间
    const appStore = useAppStore()
    appStore.recordPageLoadTime()
    
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

  // 3. 数据预加载守卫
  if (!(await preloadGuard(to, from, next))) return

  // 4. 网络状态守卫
  if (!(await networkGuard(to, from, next))) return

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

// 数据预加载守卫
async function preloadGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
): Promise<boolean> {
  const appStore = useAppStore()
  
  try {
    // 根据路由预加载必要数据
    switch (to.name) {
      case 'Orders':
      case 'OrderDetail':
        // 预加载订单数据
        await preloadOrderData(to)
        break
        
      case 'Products':
      case 'ProductDetail':
        // 预加载产品数据
        await preloadProductData(to)
        break
        
      case 'Regions':
      case 'Countries':
        // 预加载区域数据
        await preloadRegionData(to)
        break
        
      case 'Wallet':
      case 'WalletRecharge':
        // 预加载钱包数据
        await preloadWalletData(to)
        break
    }
  } catch (error) {
    console.warn('[Router] 数据预加载失败:', error)
    
    // 数据预加载失败不阻止导航，但显示警告
    appStore.showNotification({
      type: 'warning',
      message: '数据加载失败，页面可能显示不完整',
      duration: 3000
    })
  }

  return true
}

// 网络状态守卫
async function networkGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
): Promise<boolean> {
  const appStore = useAppStore()
  
  // 检查网络状态
  if (appStore.isOffline) {
    // 检查目标页面是否支持离线访问
    const offlinePages = ['Home', 'Help', 'About', 'Settings']
    
    if (!offlinePages.includes(to.name as string)) {
      console.warn('[Router] 离线状态，无法访问需要网络的页面')
      
      appStore.showNotification({
        type: 'warning',
        message: '网络连接不可用，无法访问此页面',
        duration: 3000
      })
      
      // 如果当前页面支持离线，则停留在当前页面
      if (offlinePages.includes(from.name as string)) {
        next(false)
      } else {
        // 否则重定向到首页
        next({ name: 'Home' })
      }
      return false
    }
  }

  return true
}

// 更新页面状态
function updatePageState(to: RouteLocationNormalized) {
  const appStore = useAppStore()
  
  // 更新当前页面信息
  appStore.setCurrentPage(to.name as string, to.meta.title)
  
  // 更新页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - eSIM Mini App`
  }
  
  // 更新导航状态
  appStore.setBackButton(to.meta.showBackButton || false)
  appStore.setBottomNav(to.meta.showBottomNav || false)
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

// 数据预加载函数
async function preloadOrderData(to: RouteLocationNormalized) {
  const { useOrdersStore } = await import('@/stores/orders')
  const ordersStore = useOrdersStore()
  
  if (to.name === 'OrderDetail' && to.params.id) {
    // 预加载订单详情
    await ordersStore.fetchOrderById(to.params.id as string)
  } else if (to.name === 'Orders' && !ordersStore.hasOrders) {
    // 预加载订单列表
    await ordersStore.fetchOrders()
  }
}

async function preloadProductData(to: RouteLocationNormalized) {
  const { useProductsStore } = await import('@/stores/products')
  const productsStore = useProductsStore()
  
  if (to.name === 'ProductDetail' && to.params.id) {
    // 预加载产品详情
    await productsStore.fetchProductById(to.params.id as string)
  }
  // 移除 Products 页面的预加载，因为它现在是一个容器组件
  // 各个栏目会在需要时自行加载数据
}

async function preloadRegionData(to: RouteLocationNormalized) {
  const { useProductsStore } = await import('@/stores/products')
  const productsStore = useProductsStore()
  
  // 预加载区域数据
  await productsStore.fetchRegions()
  
  if (to.name === 'Countries' && to.params.region) {
    // 预加载指定区域的国家数据
    await productsStore.fetchCountries(to.params.region as string)
  } else if (to.name === 'Countries') {
    // 预加载所有国家数据
    await productsStore.fetchCountries()
  }
}

async function preloadWalletData(to: RouteLocationNormalized) {
  // 钱包数据预加载逻辑
  // 这里可以添加钱包相关的数据预加载
  console.log('[Router] 预加载钱包数据')
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