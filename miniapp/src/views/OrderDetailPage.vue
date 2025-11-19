<template>
  <PageWrapper :loading="isLoading" loading-text="正在加载订单详情..." :error="error" @retry="handleRetry"
    class="order-detail-page">
    <div v-if="order" class="order-content">
      <!-- 订单状态头部 -->
      <div class="status-header">
        <div class="status-icon-container">
          <v-icon :icon="statusIcon" :color="statusColor" size="48" class="status-icon" />
        </div>

        <div class="status-info">
          <h2 class="status-title">{{ statusTitle }}</h2>
          <p class="status-subtitle">{{ statusSubtitle }}</p>
        </div>

        <div class="status-actions">
          <v-btn
            icon
            variant="text"
            size="small"
            @click="refreshData"
            :loading="isLoading"
            class="refresh-btn"
          >
            <v-icon>mdi-refresh</v-icon>
          </v-btn>
        </div>
      </div>

      <!-- 订单信息卡片 -->
      <v-card class="info-card" variant="elevated">
        <v-card-title class="card-title">
          <v-icon start>mdi-receipt</v-icon>
          订单信息
        </v-card-title>

        <v-card-text class="card-content">
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">订单编号</span>
              <div class="info-value-container">
                <span class="info-value">{{ order.orderNo }}</span>
                <v-btn icon variant="text" size="small" @click="copyOrderNumber">
                  <v-icon size="16">mdi-content-copy</v-icon>
                </v-btn>
              </div>
            </div>

            <div class="info-item">
              <span class="info-label">下单时间</span>
              <span class="info-value">{{ formatDateTime(order.createdAt) }}</span>
            </div>

            <div v-if="order.completedAt" class="info-item">
              <span class="info-label">完成时间</span>
              <span class="info-value">{{ formatDateTime(order.completedAt) }}</span>
            </div>

            <div class="info-item">
              <span class="info-label">支付方式</span>
              <span class="info-value">{{ paymentMethodText }}</span>
            </div>

            <div class="info-item">
              <span class="info-label">订单金额</span>
              <span class="info-value price">${{ order.totalAmount }}</span>
            </div>
          </div>
        </v-card-text>
      </v-card>

      <!-- 商品信息卡片 -->
      <v-card class="info-card" variant="elevated">
        <v-card-title class="card-title">
          <v-icon start>mdi-package-variant</v-icon>
          商品信息
        </v-card-title>

        <v-card-text class="card-content">
          <div class="product-info">
            <h4 class="product-name">{{ order.productName }}</h4>
            <div class="product-details">
              <div class="detail-chip">
                <v-icon start size="14">mdi-database</v-icon>
                1GB
              </div>
              <div class="detail-chip">
                <v-icon start size="14">mdi-calendar</v-icon>
                7天
              </div>
              <div class="detail-chip">
                <v-icon start size="14">mdi-map-marker</v-icon>
                亚洲
              </div>
            </div>
          </div>
        </v-card-text>
      </v-card>

      <!-- eSIM 信息卡片 -->
      <div v-if="order.esims && order.esims.length > 0" class="esim-section">
        <v-card class="info-card" variant="elevated">
          <v-card-title class="card-title">
            <v-icon start>mdi-sim</v-icon>
            eSIM 信息
          </v-card-title>
          <v-card-text class="card-content">
            <div v-for="(esim, index) in order.esims" :key="esim.iccid" class="esim-item">
              <div class="esim-header">
                <h4 class="esim-title">eSIM #{{ index + 1 }}</h4>
                <div class="status-badge" :class="`status-${esim.status}`">
                  {{ esim.status }}
                </div>
              </div>
              <div class="info-grid">
                <div class="info-item">
                  <span class="info-label">ICCID</span>
                  <span class="info-value">{{ esim.iccid }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">激活码</span>
                  <span class="info-value">{{ esim.activationCode }}</span>
                </div>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <v-btn v-if="canRetryPayment" color="primary" variant="elevated" block size="large" @click="retryPayment"
          :loading="retryingPayment">
          重新支付
        </v-btn>

        <v-btn v-if="order.status === EsimOrderStatus.COMPLETED && order.esims && order.esims.length > 0" 
          color="success" variant="outlined" block size="large"
          @click="scrollToESIM">
          查看 eSIM 信息
        </v-btn>

        <v-btn color="grey" variant="outlined" block size="large" @click="goBack">
          返回订单列表
        </v-btn>
      </div>
    </div>

    <!-- eSIM 相关弹窗 - 暂时注释，等待组件适配 -->
    <!--
    <ESIMUsageDialog 
      v-if="order"
      v-model="showUsageDialog" 
      :order-id="parseInt(order.id)"
      @open-topup="showTopupDialog = true"
    />
    
    <ESIMHistoryDialog 
      v-if="order"
      v-model="showHistoryDialog" 
      :order-id="parseInt(order.id)"
    />
    
    <ESIMTopupDialog 
      v-if="order"
      v-model="showTopupDialog" 
      :order-id="parseInt(order.id)"
      @topup-success="handleTopupSuccess"
    />
    
    <ESIMInstallGuide 
      v-if="order && order.esims && order.esims.length > 0"
      v-model="showInstallGuide" 
      :order-id="parseInt(order.id)"
      :esim-info="order.esims[0]"
    />
    -->
  </PageWrapper>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import { useESIMStore } from '@/stores/esim'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'
import PageWrapper from '@/components/layout/PageWrapper.vue'
import { EsimOrderStatus } from '@/types/esim-order'

// Props
interface Props {
  id?: string
}

const props = defineProps<Props>()

// 组合式 API
const route = useRoute()
const router = useRouter()
const ordersStore = useOrdersStore()
const esimStore = useESIMStore()
const appStore = useAppStore()

// 响应式状态
const isLoading = ref(false)
const error = ref<string | null>(null)
const retryingPayment = ref(false)

// 计算属性
const orderId = computed(() => props.id || route.params['id'] as string)

const order = computed(() => {
  if (orderId.value) {
    return ordersStore.findOrderById(orderId.value) || ordersStore.currentOrder
  }
  return null
})

const statusIcon = computed(() => {
  if (!order.value) return 'mdi-help-circle'

  const iconMap: Record<EsimOrderStatus, string> = {
    [EsimOrderStatus.PENDING]: 'mdi-clock-outline',
    [EsimOrderStatus.PAID]: 'mdi-check-circle',
    [EsimOrderStatus.PROCESSING]: 'mdi-cog',
    [EsimOrderStatus.COMPLETED]: 'mdi-check-circle',
    [EsimOrderStatus.CANCELLED]: 'mdi-close-circle',
    [EsimOrderStatus.FAILED]: 'mdi-alert-circle'
  }

  return iconMap[order.value.status] || 'mdi-help-circle'
})

const statusColor = computed(() => {
  if (!order.value) return 'grey'
  return ordersStore.getOrderStatusColor(order.value.status)
})

const statusTitle = computed(() => {
  if (!order.value) return '未知状态'
  return ordersStore.getOrderStatusText(order.value.status)
})

const statusSubtitle = computed(() => {
  if (!order.value) return ''

  switch (order.value.status) {
    case EsimOrderStatus.PENDING:
      return '请完成支付以激活您的 eSIM'
    case EsimOrderStatus.PAID:
      return '正在处理您的订单'
    case EsimOrderStatus.PROCESSING:
      return '正在为您准备 eSIM'
    case EsimOrderStatus.COMPLETED:
      return 'eSIM 已准备就绪'
    case EsimOrderStatus.CANCELLED:
      return '订单已取消'
    case EsimOrderStatus.FAILED:
      return '订单处理失败'
    default:
      return ''
  }
})

const canRetryPayment = computed(() => {
  return order.value?.status === EsimOrderStatus.PENDING
})

const paymentMethodText = computed(() => {
  // EsimOrder 类型暂时没有 paymentMethod 字段，默认显示钱包余额
  return '钱包余额'
})

// 方法
const fetchOrderDetail = async () => {
  if (!orderId.value) {
    error.value = '订单ID不存在'
    return
  }

  isLoading.value = true
  error.value = null

  try {
    await ordersStore.fetchOrderById(orderId.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : '获取订单详情失败'
  } finally {
    isLoading.value = false
  }
}

const handleRetry = async () => {
  await fetchOrderDetail()
}

const copyOrderNumber = async () => {
  if (!order.value) return

  try {
    await navigator.clipboard.writeText(order.value.orderNo)
    appStore.showNotification({
      type: 'success',
      message: '订单号已复制到剪贴板',
      duration: 2000
    })
    telegramService.impactFeedback('light')
  } catch (err) {
    appStore.showNotification({
      type: 'error',
      message: '复制失败',
      duration: 2000
    })
  }
}

const retryPayment = async () => {
  if (!order.value || retryingPayment.value) return

  retryingPayment.value = true

  try {
    // TODO: 实现重新支付逻辑
    appStore.showNotification({
      type: 'info',
      message: '重新支付功能开发中...',
      duration: 3000
    })
  } catch (err) {
    appStore.showNotification({
      type: 'error',
      message: err instanceof Error ? err.message : '重新支付失败',
      duration: 3000
    })
  } finally {
    retryingPayment.value = false
  }
}

const formatDateTime = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 移除未使用的 formatAmount 函数

const scrollToESIM = () => {
  const esimElement = document.querySelector('.esim-section')
  if (esimElement) {
    esimElement.scrollIntoView({ behavior: 'smooth' })
  }
}

const goBack = () => {
  router.push({ name: 'Orders' })
}

// 刷新页面数据
const refreshData = async () => {
  if (!order.value) return
  
  try {
    // 清除 eSIM 相关缓存
    esimStore.clearCache(parseInt(order.value.id))
    
    // 重新获取订单详情
    await fetchOrderDetail()
    
    appStore.showNotification({
      type: 'success',
      message: '数据已刷新',
      duration: 2000
    })
  } catch (error) {
    appStore.showNotification({
      type: 'error',
      message: '刷新失败，请重试',
      duration: 2000
    })
  }
}

// 生命周期
onMounted(async () => {
  await fetchOrderDetail()
})
</script>
<style scoped lang="scss">
  .order-detail-page {
  .order-content {
  padding: 0;
  }
  }

  .status-header {
  display: flex;
  align-items: center;
  padding: 24px 16px;
  background: linear-gradient(135deg, rgb(var(--v-theme-primary)) 0%, rgb(var(--v-theme-secondary)) 100%);
  color: white;
  margin: -16px -16px 24px -16px;

  .status-icon-container {
  margin-right: 16px;

  .status-icon {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  padding: 8px;
  }
  }

  .status-info {
  flex: 1;

  .status-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 4px 0;
  }

  .status-subtitle {
  font-size: 0.875rem;
  opacity: 0.9;
  margin: 0;
  }
  }

  .status-actions {
  .refresh-btn {
  color: white;
  background: rgba(255, 255, 255, 0.2);
  
  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
  }
  }
  }

  .info-card {
  margin-bottom: 16px;

  .card-title {
  font-size: 1.1rem;
  font-weight: 600;
  padding: 16px 16px 8px 16px;

  .v-icon {
  margin-right: 8px;
  }
  }

  .card-content {
  padding: 8px 16px 16px 16px;
  }
  }

  .info-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
  }

  .info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .info-label {
  font-size: 0.875rem;
  color: rgba(var(--v-theme-on-surface), 0.6);
  font-weight: 500;
  }

  .info-value {
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(var(--v-theme-on-surface));

  &.price {
  color: rgb(var(--v-theme-primary));
  font-weight: 600;
  font-size: 1rem;
  }
  }

  .info-value-container {
  display: flex;
  align-items: center;
  gap: 8px;
  }
  }

  .product-info {
  .product-name {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0 0 12px 0;
  color: rgb(var(--v-theme-on-surface));
  }

  .product-details {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;

  .detail-chip {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  background: rgba(var(--v-theme-primary), 0.1);
  color: rgb(var(--v-theme-primary));
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;

  .v-icon {
  margin-right: 4px;
  }
  }
  }
  }

  .esim-section {
  margin-bottom: 16px;

  .esim-item {
    padding: 12px 0;
    border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.08);

    &:last-child {
      border-bottom: none;
    }

    .esim-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .esim-title {
        font-size: 0.9375rem;
        font-weight: 600;
        margin: 0;
        color: rgb(var(--v-theme-on-surface));
      }

      .status-badge {
        padding: 4px 12px;
        border-radius: 12px;
        font-size: 0.75rem;
        font-weight: 500;
        text-transform: capitalize;

        &.status-active {
          background: #e8f5e9;
          color: #388e3c;
        }

        &.status-inactive {
          background: #f5f5f5;
          color: #757575;
        }

        &.status-pending {
          background: #fff3e0;
          color: #f57c00;
        }
      }
    }
  }
  }

  .action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  }

  // 响应式适配
  @media (max-width: 360px) {
  .status-header {
  padding: 20px 12px;
  margin: -12px -12px 20px -12px;

  .status-info {
  .status-title {
  font-size: 1.25rem;
  }

  .status-subtitle {
  font-size: 0.8rem;
  }
  }
  }

  .info-item {
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;

  .info-value-container {
  align-self: flex-end;
  }
  }
  }
  </style>