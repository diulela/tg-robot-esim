<template>
  <PageWrapper
    :loading="isLoading && !hasOrders"
    loading-text="正在加载订单..."
    :error="error"
    @retry="handleRetry"
    class="order-page"
  >
    <!-- 页面标题和统计 -->
    <div class="page-header">
      <h2 class="page-title">我的订单</h2>
      <div v-if="hasOrders" class="order-stats">
        <span class="stats-text">
          共 {{ totalOrders }} 个订单，总消费 {{ formatCurrency(totalSpent) }}
        </span>
      </div>
    </div>

    <!-- 订单筛选 -->
    <div class="filter-section">
      <v-chip-group
        v-model="selectedStatus"
        color="primary"
        variant="tonal"
        class="status-chips"
        @update:model-value="handleStatusFilter"
      >
        <v-chip value="" size="small">全部</v-chip>
        <v-chip value="pending" size="small">
          <v-icon start size="14">mdi-clock-outline</v-icon>
          待支付
        </v-chip>
        <v-chip value="paid" size="small">
          <v-icon start size="14">mdi-credit-card-check</v-icon>
          已支付
        </v-chip>
        <v-chip value="completed" size="small">
          <v-icon start size="14">mdi-check-circle</v-icon>
          已完成
        </v-chip>
        <v-chip value="cancelled" size="small">
          <v-icon start size="14">mdi-close-circle</v-icon>
          已取消
        </v-chip>
      </v-chip-group>
    </div>

    <!-- 订单列表 -->
    <div class="orders-content">
      <!-- 下拉刷新提示 -->
      <div v-if="refreshing" class="refresh-indicator">
        <v-progress-circular
          indeterminate
          size="24"
          width="2"
          color="primary"
        />
        <span class="refresh-text">正在刷新...</span>
      </div>

      <!-- 订单卡片列表 -->
      <div v-if="displayedOrders.length > 0" class="orders-list">
        <OrderCard
          v-for="order in displayedOrders"
          :key="order.id"
          :order="order"
          :show-actions="true"
          @click="navigateToOrderDetail"
          @view-details="navigateToOrderDetail"
          @copy-order-number="handleCopyOrderNumber"
        />
      </div>

      <!-- 加载更多 -->
      <div v-if="canLoadMore" class="load-more-section">
        <v-btn
          :loading="loadingMore"
          variant="outlined"
          color="primary"
          block
          @click="loadMore"
        >
          <v-icon start>mdi-refresh</v-icon>
          加载更多
        </v-btn>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!hasOrders && !isLoading" class="empty-orders">
        <v-icon size="64" color="grey-lighten-1">mdi-receipt-outline</v-icon>
        <h4 class="empty-title">暂无订单</h4>
        <p class="empty-subtitle">您还没有任何订单记录</p>
        <v-btn
          color="primary"
          variant="elevated"
          @click="navigateToProducts"
          class="mt-4"
        >
          <v-icon start>mdi-shopping</v-icon>
          开始购买
        </v-btn>
      </div>

      <!-- 筛选无结果 -->
      <div v-else-if="selectedStatus && displayedOrders.length === 0 && !isLoading" class="no-filtered-results">
        <v-icon size="48" color="grey-lighten-1">mdi-filter-off</v-icon>
        <h4 class="no-results-title">没有找到相关订单</h4>
        <p class="no-results-subtitle">
          当前筛选条件下没有订单，试试其他状态
        </p>
        <v-btn
          variant="outlined"
          color="primary"
          @click="clearFilter"
          class="mt-4"
        >
          查看全部订单
        </v-btn>
      </div>
    </div>

    <!-- 浮动操作按钮 -->
    <v-fab
      icon="mdi-refresh"
      color="primary"
      size="small"
      location="bottom end"
      class="refresh-fab"
      @click="refreshOrders"
      :loading="refreshing"
    />
  </PageWrapper>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useOrdersStore } from '@/stores/orders'
import { telegramService } from '@/services/telegram'
import type { OrderStatus } from '@/types'

import PageWrapper from '@/components/layout/PageWrapper.vue'
import OrderCard from '@/components/business/OrderCard.vue'

// 组合式 API
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const ordersStore = useOrdersStore()

// 响应式状态
const isLoading = ref(false)
const error = ref<string | null>(null)
const refreshing = ref(false)
const loadingMore = ref(false)
const selectedStatus = ref<string>('')

// 计算属性
const hasOrders = computed(() => ordersStore.hasOrders)
const totalOrders = computed(() => ordersStore.totalOrders)
const totalSpent = computed(() => ordersStore.totalSpent)
const canLoadMore = computed(() => ordersStore.canLoadMore)

const displayedOrders = computed(() => {
  if (!selectedStatus.value) {
    return ordersStore.orders
  }
  
  return ordersStore.ordersByStatus[selectedStatus.value as OrderStatus] || []
})

// 方法
const loadOrders = async (append = false) => {
  if (!append) {
    isLoading.value = true
  } else {
    loadingMore.value = true
  }
  
  error.value = null

  try {
    const filters = selectedStatus.value 
      ? { status: selectedStatus.value as OrderStatus }
      : {}
    
    await ordersStore.fetchOrders(filters, append)
    console.log('[OrderPage] 订单数据加载成功')
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : '加载订单失败'
    error.value = errorMessage
    console.error('[OrderPage] 加载订单失败:', err)
  } finally {
    isLoading.value = false
    loadingMore.value = false
  }
}

const refreshOrders = async () => {
  if (refreshing.value) return

  refreshing.value = true
  
  try {
    await ordersStore.refresh()
    
    // 触觉反馈
    telegramService.impactFeedback('light')
    
    appStore.showNotification({
      type: 'success',
      message: '订单列表已刷新',
      duration: 2000
    })
  } catch (err) {
    console.error('[OrderPage] 刷新订单失败:', err)
    
    appStore.showNotification({
      type: 'error',
      message: '刷新失败，请重试',
      duration: 3000
    })
  } finally {
    refreshing.value = false
  }
}

const loadMore = async () => {
  if (loadingMore.value || !canLoadMore.value) return
  
  await loadOrders(true)
}

const handleRetry = () => {
  loadOrders()
}

const handleStatusFilter = (status: string | undefined) => {
  selectedStatus.value = status || ''
  
  // 触觉反馈
  telegramService.selectionFeedback()
  
  // 重新加载订单
  loadOrders()
}

const clearFilter = () => {
  selectedStatus.value = ''
  loadOrders()
}

const navigateToOrderDetail = (orderId: string) => {
  router.push({ name: 'OrderDetail', params: { id: orderId } })
  telegramService.selectionFeedback()
}

const navigateToProducts = () => {
  router.push({ name: 'Products' })
  telegramService.selectionFeedback()
}

const handleCopyOrderNumber = (orderNumber: string) => {
  console.log('[OrderPage] 复制订单号:', orderNumber)
}

const formatCurrency = (amount: number): string => {
  return `$${amount.toFixed(2)}`
}

// 监听路由查询参数
watch(
  () => route.query.status,
  (newStatus) => {
    if (newStatus && newStatus !== selectedStatus.value) {
      selectedStatus.value = newStatus as string
      loadOrders()
    }
  },
  { immediate: true }
)

// 生命周期
onMounted(async () => {
  console.log('[OrderPage] 组件挂载')
  
  // 页面初始化完成
  
  // 从路由查询参数获取初始状态
  const initialStatus = route.query.status as string
  if (initialStatus) {
    selectedStatus.value = initialStatus
  }
  
  // 加载订单数据
  await loadOrders()
})
</script>

<style scoped lang="scss">
.order-page {
  .page-header {
    text-align: center;
    margin-bottom: 24px;
    
    .page-title {
      font-size: 1.5rem;
      font-weight: 600;
      color: rgb(var(--v-theme-on-surface));
      margin: 0 0 8px 0;
    }
    
    .order-stats {
      .stats-text {
        font-size: 0.875rem;
        color: rgba(var(--v-theme-on-surface), 0.7);
      }
    }
  }
  
  .filter-section {
    margin-bottom: 24px;
    
    .status-chips {
      justify-content: flex-start;
      
      .v-chip {
        margin-right: 8px;
        margin-bottom: 8px;
      }
    }
  }
  
  .orders-content {
    .refresh-indicator {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 12px;
      padding: 16px;
      margin-bottom: 16px;
      
      .refresh-text {
        font-size: 0.875rem;
        color: rgba(var(--v-theme-on-surface), 0.7);
      }
    }
    
    .orders-list {
      display: flex;
      flex-direction: column;
      gap: 12px;
    }
    
    .load-more-section {
      margin-top: 24px;
      padding: 0 16px;
    }
    
    .empty-orders,
    .no-filtered-results {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 48px 24px;
      text-align: center;
      
      .empty-title,
      .no-results-title {
        margin: 16px 0 8px;
        color: rgba(var(--v-theme-on-surface), 0.8);
        font-size: 1.125rem;
        font-weight: 600;
      }
      
      .empty-subtitle,
      .no-results-subtitle {
        margin: 0;
        color: rgba(var(--v-theme-on-surface), 0.6);
        font-size: 0.875rem;
        line-height: 1.5;
      }
    }
  }
  
  .refresh-fab {
    bottom: 80px !important;
    right: 16px !important;
  }
}

// 响应式适配
@media (max-width: 360px) {
  .order-page {
    .page-header {
      margin-bottom: 20px;
      
      .page-title {
        font-size: 1.25rem;
      }
      
      .order-stats {
        .stats-text {
          font-size: 0.8125rem;
        }
      }
    }
    
    .filter-section {
      margin-bottom: 20px;
      
      .status-chips {
        .v-chip {
          margin-right: 6px;
          margin-bottom: 6px;
        }
      }
    }
    
    .orders-content {
      .orders-list {
        gap: 10px;
      }
      
      .load-more-section {
        margin-top: 20px;
        padding: 0 12px;
      }
      
      .empty-orders,
      .no-filtered-results {
        padding: 32px 16px;
        
        .empty-title,
        .no-results-title {
          font-size: 1rem;
        }
        
        .empty-subtitle,
        .no-results-subtitle {
          font-size: 0.8125rem;
        }
      }
    }
    
    .refresh-fab {
      bottom: 76px !important;
      right: 12px !important;
    }
  }
}

@media (min-width: 481px) {
  .order-page {
    .page-header {
      margin-bottom: 32px;
      
      .page-title {
        font-size: 1.75rem;
      }
      
      .order-stats {
        .stats-text {
          font-size: 0.9375rem;
        }
      }
    }
    
    .filter-section {
      margin-bottom: 32px;
      
      .status-chips {
        .v-chip {
          margin-right: 10px;
          margin-bottom: 10px;
        }
      }
    }
    
    .orders-content {
      .orders-list {
        gap: 16px;
      }
      
      .load-more-section {
        margin-top: 32px;
        padding: 0 20px;
      }
      
      .empty-orders,
      .no-filtered-results {
        padding: 64px 32px;
        
        .empty-title,
        .no-results-title {
          font-size: 1.25rem;
        }
        
        .empty-subtitle,
        .no-results-subtitle {
          font-size: 0.9375rem;
        }
      }
    }
  }
}

// 下拉刷新动画
.refresh-indicator {
  animation: fadeInDown 0.3s ease;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 订单列表动画
.orders-list {
  .order-card {
    animation: slideInUp 0.3s ease;
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 筛选芯片样式优化
.status-chips {
  .v-chip {
    transition: all 0.2s ease;
    
    &:hover {
      transform: translateY(-1px);
    }
    
    &.v-chip--selected {
      box-shadow: 0 2px 4px rgba(var(--v-theme-primary), 0.3);
    }
  }
}

// 深色主题适配
.v-theme--dark {
  .order-page {
    .page-header {
      .page-title {
        color: rgb(var(--v-theme-on-surface));
      }
      
      .order-stats {
        .stats-text {
          color: rgba(var(--v-theme-on-surface), 0.7);
        }
      }
    }
  }
}
</style>