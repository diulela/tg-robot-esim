<template>
  <v-card
    class="order-card"
    :class="cardClasses"
    :variant="variant"
    @click="handleClick"
  >
    <v-card-text class="card-content">
      <!-- 订单头部 -->
      <div class="order-header">
        <div class="order-info">
          <h4 class="order-title">{{ order.productName }}</h4>
          <p class="order-number">{{ order.orderNumber }}</p>
        </div>
        
        <v-chip
          :color="statusColor"
          :variant="chipVariant"
          size="small"
          class="status-chip"
        >
          {{ statusText }}
        </v-chip>
      </div>
      
      <!-- 订单详情 -->
      <div v-if="!compact" class="order-details">
        <div class="detail-row">
          <span class="detail-label">下单时间</span>
          <span class="detail-value">{{ formatDate(order.createdAt) }}</span>
        </div>
        
        <div v-if="order.paidAt" class="detail-row">
          <span class="detail-label">支付时间</span>
          <span class="detail-value">{{ formatDate(order.paidAt) }}</span>
        </div>
        
        <div class="detail-row">
          <span class="detail-label">订单金额</span>
          <span class="detail-value price">{{ formatAmount(order.amount, order.currency) }}</span>
        </div>
      </div>
      
      <!-- 紧凑模式的简化信息 -->
      <div v-else class="order-summary">
        <span class="summary-date">{{ formatDate(order.createdAt) }}</span>
        <span class="summary-amount">{{ formatAmount(order.amount, order.currency) }}</span>
      </div>
      
      <!-- 操作按钮 -->
      <div v-if="showActions" class="order-actions">
        <v-btn
          v-if="canViewDetails"
          variant="outlined"
          size="small"
          color="primary"
          @click.stop="viewDetails"
        >
          查看详情
        </v-btn>
        
        <v-btn
          v-if="canCopyOrderNumber"
          variant="text"
          size="small"
          color="primary"
          @click.stop="copyOrderNumber"
        >
          <v-icon start size="16">mdi-content-copy</v-icon>
          复制订单号
        </v-btn>
        
        <v-btn
          v-if="canCancel"
          variant="text"
          size="small"
          color="error"
          @click.stop="cancelOrder"
        >
          取消订单
        </v-btn>
        
        <v-btn
          v-if="canRetryPayment"
          variant="elevated"
          size="small"
          color="primary"
          @click.stop="retryPayment"
        >
          重新支付
        </v-btn>
      </div>
      
      <!-- eSIM 信息预览 -->
      <div v-if="order.esimInfo && !compact" class="esim-preview">
        <v-divider class="my-3" />
        <div class="esim-info">
          <v-icon color="success" size="16">mdi-sim</v-icon>
          <span class="esim-text">eSIM 已激活</span>
          <v-btn
            variant="text"
            size="x-small"
            color="primary"
            @click.stop="viewESIMDetails"
          >
            查看详情
          </v-btn>
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useOrdersStore } from '@/stores/orders'
import { telegramService } from '@/services/telegram'
import type { Order, OrderStatus } from '@/types'

// Props
interface Props {
  order: Order
  showActions?: boolean
  compact?: boolean
  variant?: 'elevated' | 'flat' | 'tonal' | 'outlined' | 'text' | 'plain'
  clickable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showActions: false,
  compact: false,
  variant: 'elevated',
  clickable: true
})

// Emits
const emit = defineEmits<{
  click: [order: Order]
  viewDetails: [orderId: string]
  copyOrderNumber: [orderNumber: string]
}>()

// 组合式 API
const router = useRouter()
const appStore = useAppStore()
const ordersStore = useOrdersStore()

// 计算属性
const statusText = computed(() => {
  return ordersStore.getOrderStatusText(props.order.status)
})

const statusColor = computed(() => {
  return ordersStore.getOrderStatusColor(props.order.status)
})

const chipVariant = computed(() => {
  return props.order.status === 'completed' ? 'flat' : 'tonal'
})

const cardClasses = computed(() => {
  const classes = []
  
  if (props.clickable) {
    classes.push('clickable')
  }
  
  if (props.compact) {
    classes.push('compact')
  }
  
  return classes
})

const canViewDetails = computed(() => {
  return true // 所有订单都可以查看详情
})

const canCopyOrderNumber = computed(() => {
  return true // 所有订单都可以复制订单号
})

const canCancel = computed(() => {
  return props.order.status === 'pending'
})

const canRetryPayment = computed(() => {
  return props.order.status === 'pending' || props.order.status === 'cancelled'
})

// 方法
const formatDate = (dateString: string): string => {
  return ordersStore.formatOrderDate(dateString)
}

const formatAmount = (amount: number, currency: string): string => {
  return `${currency === 'USD' ? '$' : '¥'}${amount.toFixed(2)}`
}

const handleClick = () => {
  if (props.clickable) {
    telegramService.selectionFeedback()
    emit('click', props.order)
  }
}

const viewDetails = () => {
  emit('viewDetails', props.order.id)
  router.push({ name: 'OrderDetail', params: { id: props.order.id } })
}

const copyOrderNumber = async () => {
  try {
    await navigator.clipboard.writeText(props.order.orderNumber)
    
    appStore.showNotification({
      type: 'success',
      message: '订单号已复制到剪贴板',
      duration: 2000
    })
    
    telegramService.impactFeedback('light')
    emit('copyOrderNumber', props.order.orderNumber)
  } catch (error) {
    console.error('复制订单号失败:', error)
    
    appStore.showNotification({
      type: 'error',
      message: '复制失败，请手动复制',
      duration: 3000
    })
  }
}

const cancelOrder = async () => {
  try {
    const confirmed = await telegramService.showConfirm('确定要取消这个订单吗？')
    
    if (confirmed) {
      await ordersStore.cancelOrder(props.order.id, '用户主动取消')
      
      appStore.showNotification({
        type: 'success',
        message: '订单已取消',
        duration: 3000
      })
      
      telegramService.notificationFeedback('success')
    }
  } catch (error) {
    console.error('取消订单失败:', error)
    
    appStore.showNotification({
      type: 'error',
      message: '取消订单失败，请重试',
      duration: 3000
    })
    
    telegramService.notificationFeedback('error')
  }
}

const retryPayment = async () => {
  try {
    const paymentUrl = await ordersStore.retryPayment(props.order.id)
    
    // 打开支付链接
    telegramService.openLink(paymentUrl)
    
    appStore.showNotification({
      type: 'info',
      message: '正在跳转到支付页面...',
      duration: 3000
    })
  } catch (error) {
    console.error('重新支付失败:', error)
    
    appStore.showNotification({
      type: 'error',
      message: '获取支付链接失败，请重试',
      duration: 3000
    })
    
    telegramService.notificationFeedback('error')
  }
}

const viewESIMDetails = () => {
  router.push({ 
    name: 'OrderDetail', 
    params: { id: props.order.id },
    hash: '#esim-info'
  })
}
</script>

<style scoped lang="scss">
.order-card {
  transition: all 0.2s ease;
  
  &.clickable {
    cursor: pointer;
    
    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
  
  &.compact {
    .card-content {
      padding: 12px !important;
    }
  }
  
  .card-content {
    padding: 16px !important;
    
    .order-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 12px;
      
      .order-info {
        flex: 1;
        min-width: 0;
        
        .order-title {
          font-size: 0.875rem;
          font-weight: 600;
          color: rgb(var(--v-theme-on-surface));
          margin: 0 0 4px 0;
          line-height: 1.3;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        
        .order-number {
          font-size: 0.75rem;
          color: rgba(var(--v-theme-on-surface), 0.6);
          margin: 0;
          font-family: 'Roboto Mono', monospace;
        }
      }
      
      .status-chip {
        flex-shrink: 0;
      }
    }
    
    .order-details {
      .detail-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        .detail-label {
          font-size: 0.75rem;
          color: rgba(var(--v-theme-on-surface), 0.6);
        }
        
        .detail-value {
          font-size: 0.75rem;
          color: rgb(var(--v-theme-on-surface));
          font-weight: 500;
          
          &.price {
            color: rgb(var(--v-theme-primary));
            font-weight: 600;
          }
        }
      }
    }
    
    .order-summary {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .summary-date {
        font-size: 0.75rem;
        color: rgba(var(--v-theme-on-surface), 0.6);
      }
      
      .summary-amount {
        font-size: 0.875rem;
        color: rgb(var(--v-theme-primary));
        font-weight: 600;
      }
    }
    
    .order-actions {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: 12px;
      
      .v-btn {
        flex: 1;
        min-width: 0;
        
        &:only-child {
          flex: none;
          min-width: 120px;
        }
      }
    }
    
    .esim-preview {
      .esim-info {
        display: flex;
        align-items: center;
        gap: 8px;
        
        .esim-text {
          flex: 1;
          font-size: 0.75rem;
          color: rgb(var(--v-theme-success));
          font-weight: 500;
        }
      }
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .order-card {
    .card-content {
      padding: 12px !important;
      
      .order-header {
        gap: 8px;
        margin-bottom: 10px;
        
        .order-info {
          .order-title {
            font-size: 0.8125rem;
          }
          
          .order-number {
            font-size: 0.7rem;
          }
        }
      }
      
      .order-actions {
        gap: 6px;
        margin-top: 10px;
        
        .v-btn {
          font-size: 0.75rem;
        }
      }
    }
  }
}

@media (min-width: 481px) {
  .order-card {
    .card-content {
      padding: 18px !important;
      
      .order-header {
        gap: 16px;
        margin-bottom: 16px;
        
        .order-info {
          .order-title {
            font-size: 0.9375rem;
          }
          
          .order-number {
            font-size: 0.8125rem;
          }
        }
      }
      
      .order-details {
        .detail-row {
          margin-bottom: 10px;
          
          .detail-label,
          .detail-value {
            font-size: 0.8125rem;
          }
        }
      }
      
      .order-summary {
        .summary-date {
          font-size: 0.8125rem;
        }
        
        .summary-amount {
          font-size: 0.9375rem;
        }
      }
      
      .order-actions {
        gap: 10px;
        margin-top: 16px;
      }
    }
  }
}
</style>