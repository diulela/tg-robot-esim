<template>
  <PageWrapper
    :loading="isLoading"
    loading-text="正在加载数据..."
    :error="error"
    @retry="handleRetry"
    class="home-page"
  >
    <!-- 用户问候区域 -->
    <div class="greeting-section">
      <div class="greeting-background">
        <div class="greeting-content">
          <!-- 用户头像和问候 -->
          <div class="user-greeting">
            <v-avatar size="48" class="user-avatar">
              <v-img
                v-if="userStore.avatarUrl"
                :src="userStore.avatarUrl"
                :alt="userStore.displayName"
              />
              <v-icon v-else size="32" color="white">mdi-account-circle</v-icon>
            </v-avatar>
            
            <div class="greeting-text">
              <h2 class="greeting-title">{{ greetingMessage }}</h2>
              <p class="greeting-subtitle">{{ currentTimeText }}</p>
            </div>
            
            <!-- 用户按钮 -->
            <v-btn
              icon
              variant="text"
              color="white"
              @click="navigateToProfile"
              class="profile-btn"
            >
              <v-icon>mdi-account</v-icon>
            </v-btn>
          </div>
        </div>
      </div>
    </div>

    <!-- 今日概览 -->
    <div class="overview-section">
      <div class="section-header">
        <h3 class="section-title">今日概览</h3>
        <v-btn
          variant="text"
          color="primary"
          size="small"
          @click="refreshStats"
          :loading="statsLoading"
        >
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </div>

      <div class="stats-grid">
        <StatsCard
          title="待订单"
          :value="dashboardStats.pendingOrders"
          :change="12"
          change-type="increase"
          icon="mdi-clock-outline"
          color="warning"
          @click="navigateToOrders('pending')"
        />
        
        <StatsCard
          title="待付款"
          :value="dashboardStats.totalOrders - dashboardStats.completedOrders"
          :change="-5"
          change-type="decrease"
          icon="mdi-credit-card-outline"
          color="error"
          @click="navigateToOrders('paid')"
        />
        
        <StatsCard
          title="已完成"
          :value="dashboardStats.completedOrders"
          :change="8"
          change-type="increase"
          icon="mdi-check-circle-outline"
          color="success"
          @click="navigateToOrders('completed')"
        />
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="actions-section">
      <div class="section-header">
        <h3 class="section-title">快捷操作</h3>
        <v-btn
          variant="text"
          color="primary"
          size="small"
          @click="showMoreActions"
        >
          查看更多
        </v-btn>
      </div>

      <div class="action-cards">
        <ActionCard
          title="浏览商品"
          description="查看最新产品"
          icon="mdi-shopping"
          color="primary"
          @click="navigateToProducts"
        />
        
        <ActionCard
          title="我的订单"
          description="查看订单状态"
          icon="mdi-receipt"
          color="secondary"
          @click="navigateToOrders()"
        />
      </div>
    </div>

    <!-- 最近订单 -->
    <div v-if="recentOrders.length > 0" class="recent-section">
      <div class="section-header">
        <h3 class="section-title">最近订单</h3>
        <v-btn
          variant="text"
          color="primary"
          size="small"
          @click="navigateToOrders()"
        >
          查看全部
        </v-btn>
      </div>

      <div class="recent-orders">
        <OrderCard
          v-for="order in recentOrders.slice(0, 3)"
          :key="order.id"
          :order="order"
          compact
          @click="navigateToOrderDetail(order.id)"
        />
      </div>
    </div>

    <!-- 空状态提示 -->
    <div v-else class="empty-orders">
      <v-icon size="64" color="grey-lighten-1">mdi-inbox-outline</v-icon>
      <h4 class="empty-title">暂无订单记录</h4>
      <p class="empty-subtitle">开始浏览商品，创建您的第一个订单</p>
      <v-btn
        color="primary"
        variant="elevated"
        @click="navigateToProducts"
        class="mt-4"
      >
        <v-icon start>mdi-shopping</v-icon>
        浏览商品
      </v-btn>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { useOrdersStore } from '@/stores/orders'
import { userApi } from '@/services/api'
import { telegramService } from '@/services/telegram'
import type { DashboardStats, OrderStatus } from '@/types'

import PageWrapper from '@/components/layout/PageWrapper.vue'
import StatsCard from '@/components/common/StatsCard.vue'
import ActionCard from '@/components/common/ActionCard.vue'
import OrderCard from '@/components/business/OrderCard.vue'

// 组合式 API
const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()
const ordersStore = useOrdersStore()

// 响应式状态
const isLoading = ref(false)
const error = ref<string | null>(null)
const statsLoading = ref(false)
const currentTime = ref(new Date())
const dashboardStats = ref<DashboardStats>({
  totalOrders: 0,
  pendingOrders: 0,
  completedOrders: 0,
  totalSpent: 0,
  currency: 'USD',
  orderTrends: [],
  recentOrders: []
})

// 定时器
let timeInterval: NodeJS.Timeout | null = null

// 计算属性
const greetingMessage = computed(() => {
  if (userStore.isAuthenticated) {
    return userStore.getGreeting()
  }
  return '欢迎使用 eSIM 商城'
})

const currentTimeText = computed(() => {
  const now = currentTime.value
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  }
  return now.toLocaleDateString('zh-CN', options)
})

const recentOrders = computed(() => {
  return ordersStore.recentOrders
})

// 方法
const loadDashboardData = async () => {
  if (!userStore.isAuthenticated) return

  isLoading.value = true
  error.value = null

  try {
    // 并行加载数据
    const [stats, orders] = await Promise.all([
      userApi.getStats(),
      ordersStore.fetchOrders({ pageSize: 5 })
    ])

    dashboardStats.value = stats
    console.log('[HomePage] 仪表板数据加载成功')
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : '加载数据失败'
    error.value = errorMessage
    console.error('[HomePage] 加载仪表板数据失败:', err)
  } finally {
    isLoading.value = false
  }
}

const refreshStats = async () => {
  if (statsLoading.value) return

  statsLoading.value = true
  
  try {
    const stats = await userApi.getStats()
    dashboardStats.value = stats
    
    // 触觉反馈
    telegramService.impactFeedback('light')
    
    appStore.showNotification({
      type: 'success',
      message: '数据已刷新',
      duration: 2000
    })
  } catch (err) {
    console.error('[HomePage] 刷新统计数据失败:', err)
    appStore.showNotification({
      type: 'error',
      message: '刷新失败，请重试',
      duration: 3000
    })
  } finally {
    statsLoading.value = false
  }
}

const handleRetry = () => {
  loadDashboardData()
}

const updateCurrentTime = () => {
  currentTime.value = new Date()
}

// 导航方法
const navigateToProfile = () => {
  router.push({ name: 'Profile' })
  telegramService.selectionFeedback()
}

const navigateToProducts = () => {
  router.push({ name: 'Products' })
  telegramService.selectionFeedback()
}

const navigateToOrders = (status?: OrderStatus | string) => {
  const query = status ? { status } : {}
  router.push({ name: 'Orders', query })
  telegramService.selectionFeedback()
}

const navigateToOrderDetail = (orderId: string) => {
  router.push({ name: 'OrderDetail', params: { id: orderId } })
  telegramService.selectionFeedback()
}

const showMoreActions = () => {
  // 显示更多操作的底部抽屉或菜单
  appStore.showNotification({
    type: 'info',
    message: '更多功能即将上线',
    duration: 2000
  })
}

// 生命周期
onMounted(async () => {
  console.log('[HomePage] 组件挂载')
  
  // 设置页面标题
  appStore.setCurrentPage('Home', '首页')
  
  // 启动时间更新定时器
  timeInterval = setInterval(updateCurrentTime, 60000) // 每分钟更新一次
  
  // 加载数据
  await loadDashboardData()
})

onUnmounted(() => {
  console.log('[HomePage] 组件卸载')
  
  // 清理定时器
  if (timeInterval) {
    clearInterval(timeInterval)
    timeInterval = null
  }
})
</script>

<style scoped lang="scss">
.home-page {
  .greeting-section {
    margin: -16px -16px 24px -16px;
    
    .greeting-background {
      background: linear-gradient(135deg, #6366F1 0%, #8B5CF6 50%, #A855F7 100%);
      padding: 24px 16px;
      border-radius: 0 0 24px 24px;
      
      .greeting-content {
        .user-greeting {
          display: flex;
          align-items: center;
          gap: 16px;
          
          .user-avatar {
            border: 2px solid rgba(255, 255, 255, 0.3);
          }
          
          .greeting-text {
            flex: 1;
            
            .greeting-title {
              color: white;
              font-size: 1.25rem;
              font-weight: 600;
              margin: 0 0 4px 0;
            }
            
            .greeting-subtitle {
              color: rgba(255, 255, 255, 0.8);
              font-size: 0.875rem;
              margin: 0;
            }
          }
          
          .profile-btn {
            opacity: 0.8;
            
            &:hover {
              opacity: 1;
            }
          }
        }
      }
    }
  }
  
  .overview-section,
  .actions-section,
  .recent-section {
    margin-bottom: 32px;
    
    .section-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 16px;
      
      .section-title {
        font-size: 1.125rem;
        font-weight: 600;
        color: rgb(var(--v-theme-on-surface));
        margin: 0;
      }
    }
  }
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    
    @media (max-width: 360px) {
      gap: 8px;
    }
  }
  
  .action-cards {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    
    @media (max-width: 360px) {
      gap: 12px;
    }
  }
  
  .recent-orders {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .empty-orders {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    text-align: center;
    
    .empty-title {
      margin: 16px 0 8px;
      color: rgba(var(--v-theme-on-surface), 0.8);
      font-size: 1.125rem;
      font-weight: 600;
    }
    
    .empty-subtitle {
      margin-bottom: 0;
      color: rgba(var(--v-theme-on-surface), 0.6);
      font-size: 0.875rem;
      line-height: 1.5;
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .home-page {
    .greeting-section {
      margin: -12px -12px 20px -12px;
      
      .greeting-background {
        padding: 20px 12px;
        border-radius: 0 0 20px 20px;
        
        .user-greeting {
          gap: 12px;
          
          .greeting-text {
            .greeting-title {
              font-size: 1.125rem;
            }
            
            .greeting-subtitle {
              font-size: 0.8125rem;
            }
          }
        }
      }
    }
    
    .overview-section,
    .actions-section,
    .recent-section {
      margin-bottom: 24px;
    }
    
    .empty-orders {
      padding: 32px 16px;
    }
  }
}

@media (min-width: 481px) {
  .home-page {
    .greeting-section {
      margin: -20px -20px 32px -20px;
      
      .greeting-background {
        padding: 32px 20px;
        
        .user-greeting {
          gap: 20px;
          
          .greeting-text {
            .greeting-title {
              font-size: 1.375rem;
            }
          }
        }
      }
    }
    
    .stats-grid {
      gap: 16px;
    }
    
    .action-cards {
      gap: 20px;
    }
  }
}
</style>