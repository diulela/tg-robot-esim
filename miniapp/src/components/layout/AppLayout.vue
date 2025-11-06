<template>
  <v-app :theme="currentTheme">
    <!-- 顶部应用栏 -->
    <v-app-bar
      v-if="showAppBar"
      :color="appBarColor"
      density="compact"
      flat
      class="app-bar"
    >
      <!-- 后退按钮 -->
      <v-app-bar-nav-icon
        v-if="showBackButton"
        @click="handleBack"
        class="touch-target"
      >
        <v-icon>mdi-arrow-left</v-icon>
      </v-app-bar-nav-icon>

      <!-- 页面标题 -->
      <v-app-bar-title class="app-bar-title">
        {{ pageTitle }}
      </v-app-bar-title>

      <v-spacer />

      <!-- 用户按钮 -->
      <v-btn
        v-if="showUserButton && user"
        icon
        @click="showUserMenu"
        class="touch-target"
      >
        <v-avatar size="32">
          <v-img
            v-if="user.avatarUrl"
            :src="user.avatarUrl"
            :alt="user.displayName"
          />
          <v-icon v-else>mdi-account-circle</v-icon>
        </v-avatar>
      </v-btn>

      <!-- 主题切换按钮 (仅开发模式) -->
      <v-btn
        v-if="debugMode"
        icon
        @click="toggleTheme"
        class="touch-target"
      >
        <v-icon>
          {{ isDarkMode ? 'mdi-weather-sunny' : 'mdi-weather-night' }}
        </v-icon>
      </v-btn>
    </v-app-bar>

    <!-- 主内容区域 -->
    <v-main class="main-content">
      <div class="page-container">
        <!-- 页面加载指示器 -->
        <v-progress-linear
          v-if="pageLoading"
          indeterminate
          color="primary"
          class="page-loading"
        />

        <!-- 离线提示 -->
        <v-banner
          v-if="isOffline"
          color="warning"
          icon="mdi-wifi-off"
          class="offline-banner"
        >
          <template #text>
            网络连接不可用，部分功能可能受限
          </template>
        </v-banner>

        <!-- 路由视图 -->
        <router-view v-slot="{ Component, route }">
          <transition
            :name="getTransitionName(route)"
            mode="out-in"
            @enter="onPageEnter"
            @leave="onPageLeave"
          >
            <keep-alive :include="keepAlivePages">
              <component :is="Component" :key="route.fullPath" />
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </v-main>

    <!-- 底部导航栏 -->
    <v-bottom-navigation
      v-if="showBottomNav"
      v-model="currentTab"
      color="primary"
      grow
      class="bottom-nav"
    >
      <!-- <v-btn
        value="home"
        @click="navigateTo('Home')"
        class="nav-btn"
      >
        <v-icon>mdi-home</v-icon>
        <span>首页</span>
      </v-btn> -->

      <v-btn
        value="products"
        @click="navigateTo('Products')"
        class="nav-btn"
      >
        <v-icon>mdi-shopping</v-icon>
        <span>商品</span>
      </v-btn>

      <!-- <v-btn
        value="orders"
        @click="navigateTo('Orders')"
        class="nav-btn"
      >
        <v-icon>mdi-receipt</v-icon>
        <span>订单</span>
      </v-btn> -->

      <v-btn
        value="profile"
        @click="navigateTo('Profile')"
        class="nav-btn"
      >
        <v-icon>mdi-account</v-icon>
        <span>我的</span>
      </v-btn>
    </v-bottom-navigation>

    <!-- 用户菜单 -->
    <v-menu
      v-model="userMenuOpen"
      :activator="userMenuActivator"
      location="bottom end"
      offset="8"
    >
      <v-list>
        <v-list-item
          v-if="user"
          :title="user.displayName"
          :subtitle="user.isPremium ? 'Telegram Premium' : 'Telegram'"
        >
          <template #prepend>
            <v-avatar size="40">
              <v-img
                v-if="user.avatarUrl"
                :src="user.avatarUrl"
                :alt="user.displayName"
              />
              <v-icon v-else>mdi-account-circle</v-icon>
            </v-avatar>
          </template>
        </v-list-item>

        <v-divider />

        <v-list-item
          title="个人中心"
          prepend-icon="mdi-account"
          @click="navigateTo('Profile')"
        />

        <v-list-item
          title="我的钱包"
          prepend-icon="mdi-wallet"
          @click="navigateTo('Wallet')"
        />

        <v-list-item
          title="设置"
          prepend-icon="mdi-cog"
          @click="navigateTo('Settings')"
        />

        <v-divider />

        <v-list-item
          title="帮助中心"
          prepend-icon="mdi-help-circle"
          @click="navigateTo('Help')"
        />
      </v-list>
    </v-menu>

    <!-- 全局通知 -->
    <div class="notifications-container">
      <v-snackbar
        v-for="notification in notifications"
        :key="notification.id"
        v-model="notification.show"
        :color="notification.type"
        :timeout="notification.duration || 5000"
        :persistent="notification.persistent"
        location="top"
        class="notification-snackbar"
      >
        {{ notification.message }}
        
        <template #actions>
          <v-btn
            v-if="notification.persistent"
            icon="mdi-close"
            @click="removeNotification(notification.id)"
          />
        </template>
      </v-snackbar>
    </div>
  </v-app>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { telegramService } from '@/services/telegram'
import type { Notification } from '@/types'

// Props
interface Props {
  showAppBar?: boolean
  showUserButton?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showAppBar: true,
  showUserButton: true
})

// 组合式 API
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

// 响应式状态
const userMenuOpen = ref(false)
const userMenuActivator = ref<HTMLElement | null>(null)
const currentTab = ref('home')
const notifications = ref<Array<Notification & { show: boolean }>>([])

// 计算属性
const currentTheme = computed(() => appStore.currentTheme)
const isDarkMode = computed(() => appStore.isDarkMode)
const pageTitle = computed(() => appStore.pageTitle || route.meta.title || '首页')
const showBackButton = computed(() => appStore.showBackButton)
const showBottomNav = computed(() => appStore.showBottomNav)
const pageLoading = computed(() => appStore.pageLoading)
const isOffline = computed(() => appStore.isOffline)
const debugMode = computed(() => appStore.debugMode)
const user = computed(() => userStore.isAuthenticated ? {
  displayName: userStore.displayName,
  avatarUrl: userStore.avatarUrl,
  isPremium: userStore.isPremium
} : null)

const appBarColor = computed(() => {
  if (isDarkMode.value) {
    return 'surface'
  }
  return 'primary'
})

const keepAlivePages = computed(() => {
  // 需要缓存的页面组件名称
  return ['HomePage', 'ProductListPage', 'OrderPage', 'ProfilePage']
})

// 方法
const handleBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push({ name: 'Home' })
  }
  
  // 触觉反馈
  telegramService.impactFeedback('light')
}

const showUserMenu = (event: Event) => {
  userMenuActivator.value = event.currentTarget as HTMLElement
  userMenuOpen.value = true
}

const toggleTheme = async () => {
  await appStore.toggleTheme()
  telegramService.impactFeedback('light')
}

const navigateTo = async (routeName: string) => {
  try {
    await router.push({ name: routeName })
    telegramService.selectionFeedback()
  } catch (error) {
    console.error('导航失败:', error)
  }
}

const getTransitionName = (route: any) => {
  // 根据路由层级决定过渡动画
  const routeDepth = route.path.split('/').length
  const currentDepth = router.currentRoute.value.path.split('/').length
  
  if (routeDepth > currentDepth) {
    return 'slide-left'
  } else if (routeDepth < currentDepth) {
    return 'slide-right'
  } else {
    return 'fade'
  }
}

const onPageEnter = () => {
  appStore.setLoading(false)
}

const onPageLeave = () => {
  appStore.setLoading(true)
}

const removeNotification = (id: string) => {
  appStore.removeNotification(id)
}

// 监听通知变化
watch(
  () => appStore.notifications,
  (newNotifications) => {
    notifications.value = newNotifications.map(n => ({
      ...n,
      show: true
    }))
  },
  { deep: true, immediate: true }
)

// 监听路由变化更新底部导航
watch(
  () => route.name,
  (newRouteName) => {
    const routeTabMap: Record<string, string> = {
      'Home': 'home',
      'Products': 'products',
      'Orders': 'orders',
      'Profile': 'profile',
      'Wallet': 'profile'
    }
    
    currentTab.value = routeTabMap[newRouteName as string] || 'home'
  },
  { immediate: true }
)

// 监听 Telegram 后退按钮
const handleTelegramBack = () => {
  handleBack()
}

// 生命周期
onMounted(() => {
  // 监听 Telegram 后退按钮
  telegramService.on('backButtonClicked', handleTelegramBack)
  
  // 初始化当前标签页
  const routeTabMap: Record<string, string> = {
    'Home': 'home',
    'Products': 'products', 
    'Orders': 'orders',
    'Profile': 'profile'
  }
  currentTab.value = routeTabMap[route.name as string] || 'home'
})

onUnmounted(() => {
  // 清理事件监听
  telegramService.off('backButtonClicked', handleTelegramBack)
})
</script>

<style scoped lang="scss">
.app-bar {
  .app-bar-title {
    font-weight: 600;
    font-size: 1.1rem;
  }
}

.main-content {
  .page-container {
    min-height: 100vh;
    position: relative;
  }
  
  .page-loading {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    z-index: 1000;
  }
  
  .offline-banner {
    position: sticky;
    top: 0;
    z-index: 999;
  }
}

.bottom-nav {
  .nav-btn {
    min-width: 64px;
    
    .v-icon {
      margin-bottom: 4px;
    }
    
    span {
      font-size: 0.75rem;
      line-height: 1;
    }
  }
}

.notifications-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 2000;
  pointer-events: none;
  
  .notification-snackbar {
    pointer-events: auto;
  }
}

// 页面过渡动画
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.3s ease;
}

.slide-left-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.slide-left-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-right-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-right-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 响应式适配
@media (max-width: 360px) {
  .app-bar {
    .app-bar-title {
      font-size: 1rem;
    }
  }
  
  .bottom-nav {
    .nav-btn span {
      font-size: 0.7rem;
    }
  }
}
</style>