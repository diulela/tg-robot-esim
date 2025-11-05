// Vue Router 配置
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { RouteMetaData } from '@/types'

// 扩展路由元信息类型
declare module 'vue-router' {
  interface RouteMeta extends RouteMetaData {}
}

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/ProductPage.vue'),
    meta: {
      title: '商品列表',
      showBottomNav: true,
      keepAlive: true
    }
  },
  {
    path: '/regions',
    name: 'Regions',
    component: () => import('@/views/RegionPage.vue'),
    meta: {
      title: '选择区域',
      showBackButton: true,
      keepAlive: true
    }
  },
  {
    path: '/countries/:region?',
    name: 'Countries',
    component: () => import('@/views/CountryPage.vue'),
    meta: {
      title: '选择国家',
      showBackButton: true,
      keepAlive: true
    },
    props: true
  },
  {
    path: '/products',
    name: 'Products',
    component: () => import('@/views/ProductPage.vue'),
    meta: {
      title: '商品列表',
      showBackButton: true,
      showBottomNav: true,
      keepAlive: true
    }
  },
  {
    path: '/products/list/:countryCode',
    name: 'ProductListSecondary',
    component: () => import('@/views/ProductListSecondaryPage.vue'),
    meta: {
      title: '商品列表',
      showBackButton: true,
      showBottomNav: false,
      keepAlive: false
    },
    props: true
  },
  // 保留旧路由以兼容
  {
    path: '/hot-products/:hotItemCode',
    redirect: to => {
      return {
        name: 'ProductListSecondary',
        params: { countryCode: to.params.hotItemCode },
        query: to.query
      }
    }
  },
  {
    path: '/products/:id',
    name: 'ProductDetail',
    component: () => import('@/views/ProductDetailPage.vue'),
    meta: {
      title: '商品详情',
      showBackButton: true
    },
    props: true
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/OrderPage.vue'),
    meta: {
      title: '我的订单',
      showBottomNav: true,
      requiresAuth: true,
      keepAlive: true
    }
  },
  {
    path: '/orders/:id',
    name: 'OrderDetail',
    component: () => import('@/views/OrderDetailPage.vue'),
    meta: {
      title: '订单详情',
      showBackButton: true,
      requiresAuth: true
    },
    props: true
  },
  {
    path: '/wallet',
    name: 'Wallet',
    component: () => import('@/views/WalletPage.vue'),
    meta: {
      title: '我的钱包',
      showBackButton: true,
      showBottomNav: true,
      requiresAuth: true,
      keepAlive: true
    }
  },
  {
    path: '/wallet/recharge',
    name: 'WalletRecharge',
    component: () => import('@/views/WalletRechargePage.vue'),
    meta: {
      title: '钱包充值',
      showBackButton: true,
      requiresAuth: true
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfilePage.vue'),
    meta: {
      title: '个人中心',
      showBackButton: true,
      showBottomNav: true,
      requiresAuth: true,
      keepAlive: true
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsPage.vue'),
    meta: {
      title: '设置',
      showBackButton: true,
      requiresAuth: true
    }
  },
  {
    path: '/help',
    name: 'Help',
    component: () => import('@/views/HelpPage.vue'),
    meta: {
      title: '帮助中心',
      showBackButton: true
    }
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/AboutPage.vue'),
    meta: {
      title: '关于我们',
      showBackButton: true
    }
  },
  // 错误页面
  {
    path: '/error',
    name: 'Error',
    component: () => import('@/views/ErrorPage.vue'),
    meta: {
      title: '出错了',
      showBackButton: true
    }
  },
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundPage.vue'),
    meta: {
      title: '页面不存在',
      showBackButton: true
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // 如果有保存的滚动位置，恢复它
    if (savedPosition) {
      return savedPosition
    }
    
    // 如果是锚点链接，滚动到对应元素
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    }
    
    // 默认滚动到顶部
    return { top: 0 }
  }
})

// 导出路由实例
export default router

// 导出路由配置用于其他地方使用
export { routes }

// 路由工具函数
export function getRouteByName(name: string) {
  return routes.find(route => route.name === name)
}

export function getRouteTitle(name: string): string {
  const route = getRouteByName(name)
  return route?.meta?.title || '未知页面'
}

export function isRouteRequiresAuth(name: string): boolean {
  const route = getRouteByName(name)
  return route?.meta?.requiresAuth || false
}

export function shouldShowBackButton(name: string): boolean {
  const route = getRouteByName(name)
  return route?.meta?.showBackButton || false
}

export function shouldShowBottomNav(name: string): boolean {
  const route = getRouteByName(name)
  return route?.meta?.showBottomNav || false
}

export function shouldKeepAlive(name: string): boolean {
  const route = getRouteByName(name)
  return route?.meta?.keepAlive || false
}

// 路由导航方法
export function navigateToHome() {
  return router.push({ name: 'Home' })
}

export function navigateToRegions() {
  return router.push({ name: 'Regions' })
}

export function navigateToCountries(region?: string) {
  return router.push({ 
    name: 'Countries', 
    params: region ? { region } : {} 
  })
}

export function navigateToProducts() {
  return router.push({ name: 'Products' })
}

export function navigateToHotProducts(hotItemCode: string, hotItemName?: string) {
  return router.push({ 
    name: 'HotProductsSecondary', 
    params: { hotItemCode },
    query: hotItemName ? { name: hotItemName } : {}
  })
}

export function navigateToProductDetail(id: string) {
  return router.push({ 
    name: 'ProductDetail', 
    params: { id } 
  })
}

export function navigateToOrders() {
  return router.push({ name: 'Orders' })
}

export function navigateToOrderDetail(id: string) {
  return router.push({ 
    name: 'OrderDetail', 
    params: { id } 
  })
}

export function navigateToWallet() {
  return router.push({ name: 'Wallet' })
}

export function navigateToWalletRecharge() {
  return router.push({ name: 'WalletRecharge' })
}

export function navigateToProfile() {
  return router.push({ name: 'Profile' })
}

export function navigateToSettings() {
  return router.push({ name: 'Settings' })
}

export function navigateToHelp() {
  return router.push({ name: 'Help' })
}

export function navigateToAbout() {
  return router.push({ name: 'About' })
}

export function navigateBack() {
  return router.back()
}

export function navigateToError(error?: string) {
  return router.push({ 
    name: 'Error', 
    query: error ? { error } : {} 
  })
}