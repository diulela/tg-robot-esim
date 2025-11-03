<template>
  <AppLayout>
    <!-- 应用内容通过路由视图渲染 -->
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { telegramService } from '@/services/telegram'

// 组合式 API
const appStore = useAppStore()
const userStore = useUserStore()

// 应用初始化
const initializeApp = async () => {
  try {
    console.log('[App] 开始初始化应用...')
    
    // 1. 初始化应用状态
    await appStore.initialize()
    
    // 2. 初始化用户状态
    await userStore.initializeUser()
    
    console.log('[App] 应用初始化完成')
  } catch (error) {
    console.error('[App] 应用初始化失败:', error)
    
    // 显示错误通知
    appStore.showNotification({
      type: 'error',
      message: '应用初始化失败，请刷新重试',
      persistent: true
    })
  }
}

// 处理应用错误
const handleError = (error: Error) => {
  console.error('[App] 全局错误:', error)
  
  appStore.showNotification({
    type: 'error',
    message: '应用出现错误，请刷新重试',
    persistent: true
  })
}

// 处理未捕获的 Promise 拒绝
const handleUnhandledRejection = (event: PromiseRejectionEvent) => {
  console.error('[App] 未处理的 Promise 拒绝:', event.reason)
  
  appStore.showNotification({
    type: 'error',
    message: '操作失败，请重试',
    duration: 3000
  })
  
  // 阻止默认的错误处理
  event.preventDefault()
}

// 生命周期
onMounted(async () => {
  // 设置全局错误处理
  window.addEventListener('error', (event) => handleError(event.error))
  window.addEventListener('unhandledrejection', handleUnhandledRejection)
  
  // 初始化应用
  await initializeApp()
})

onUnmounted(() => {
  // 清理事件监听
  window.removeEventListener('error', (event) => handleError(event.error))
  window.removeEventListener('unhandledrejection', handleUnhandledRejection)
})
</script>

<style lang="scss">
// 全局样式导入
@import '@/styles/global.scss';
@import '@/styles/mobile.scss';

// 应用根样式
#app {
  font-family: 'Roboto', 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: rgb(var(--v-theme-on-background));
  background: rgb(var(--v-theme-background));
}

// 确保应用占满全屏
html, body {
  height: 100%;
  margin: 0;
  padding: 0;
  overflow-x: hidden;
}

// Telegram Web App 特定样式
body.tg-web-app {
  // 隐藏默认的滚动条
  scrollbar-width: none;
  -ms-overflow-style: none;
  
  &::-webkit-scrollbar {
    display: none;
  }
}

// 禁用文本选择 (移动端优化)
* {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

// 允许输入框文本选择
input, textarea, [contenteditable] {
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
  user-select: text;
}

// 禁用长按菜单 (移动端)
* {
  -webkit-touch-callout: none;
  -webkit-tap-highlight-color: transparent;
}

// 优化触摸滚动
* {
  -webkit-overflow-scrolling: touch;
}
</style>