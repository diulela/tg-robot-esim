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
              <span class="info-label">订单金额</span>
              <span class="info-value price">{{ formatAmount(order.totalAmount) }}</span>
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
                <v-icon start size="14">mdi-package-variant</v-icon>
                数量: {{ order.quantity }}
              </div>
              <div class="detail-chip">
                <v-icon start size="14">mdi-currency-usd</v-icon>
                单价: {{ formatAmount(order.unitPrice) }}
              </div>
              <div class="detail-chip">
                <v-icon start size="14">mdi-calculator</v-icon>
                小计: {{ formatAmount(order.totalAmount) }}
              </div>
            </div>
          </div>
        </v-card-text>
      </v-card>

      <!-- eSIM 信息卡片列表 -->
      <div v-if="esimList.length > 0" class="esim-section">
        <v-card class="info-card" variant="elevated">
          <v-card-title class="card-title">
            <v-icon start>mdi-sim</v-icon>
            eSIM 信息 (共{{ esimList.length }}张卡)
          </v-card-title>
          <v-card-text class="card-content">
            <div v-for="(esim, index) in esimList" :key="esim.iccid" class="esim-item">
              <div class="esim-header">
                <v-chip color="primary" size="small">
                  <v-icon start size="14">mdi-sim</v-icon>
                  eSIM {{ index + 1 }}
                </v-chip>
              </div>
              <div class="esim-details">
                <div class="detail-row">
                  <span class="detail-label">ICCID</span>
                  <div class="detail-value-container">
                    <span class="detail-value">{{ esim.iccid }}</span>
                    <v-btn icon variant="text" size="small" @click="copyText(esim.iccid, 'ICCID')">
                      <v-icon size="16">mdi-content-copy</v-icon>
                    </v-btn>
                  </div>
                </div>
                <div class="detail-row">
                  <span class="detail-label">激活码</span>
                  <div class="detail-value-container">
                    <span class="detail-value">{{ esim.activationCode }}</span>
                    <v-btn icon variant="text" size="small" @click="copyText(esim.activationCode, '激活码')">
                      <v-icon size="16">mdi-content-copy</v-icon>
                    </v-btn>
                  </div>
                </div>
                <div v-if="esim.apnType" class="detail-row">
                  <span class="detail-label">APN类型</span>
                  <span class="detail-value">{{ esim.apnType === 'manual' ? '手动' : '自动' }}</span>
                </div>
                <div v-if="esim.isRoaming !== undefined" class="detail-row">
                  <span class="detail-label">是否漫游</span>
                  <span class="detail-value">{{ esim.isRoaming ? '是' : '否' }}</span>
                </div>
              </div>
              <!-- 二维码 -->
              <div v-if="esim.qrCode" class="qr-code-section">
                <img :src="esim.qrCode" alt="eSIM QR Code" class="qr-code-image" />
              </div>
            </div>
          </v-card-text>
        </v-card>
      </div>

      <!-- 费用信息卡片 -->
      <v-card class="info-card" variant="elevated">
        <v-card-title class="card-title">
          <v-icon start>mdi-cash</v-icon>
          费用信息
        </v-card-title>

        <v-card-text class="card-content">
          <div class="cost-info">
            <div class="cost-item">
              <span class="cost-label">商品总额</span>
              <span class="cost-value">{{ formatAmount(order.totalAmount) }}</span>
            </div>
            <v-divider class="my-2" />
            <div class="cost-item total">
              <span class="cost-label">实付金额</span>
              <span class="cost-value price">{{ formatAmount(order.totalAmount) }}</span>
            </div>
          </div>
        </v-card-text>
      </v-card>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <v-btn 
          v-if="esimList.length > 1" 
          color="primary" 
          variant="elevated" 
          block 
          size="large" 
          @click="exportAllPDF"
          :loading="isExportingAll"
        >
          <v-icon start>mdi-file-pdf-box</v-icon>
          导出全部 eSIM 卡为 PDF
        </v-btn>

        <v-btn 
          v-if="canRefund" 
          color="warning" 
          variant="outlined" 
          block 
          size="large" 
          @click="showRefundDialog = true"
        >
          <v-icon start>mdi-undo</v-icon>
          申请退款
        </v-btn>

        <v-btn v-if="canRetryPayment" color="primary" variant="elevated" block size="large" @click="retryPayment"
          :loading="retryingPayment">
          重新支付
        </v-btn>

        <v-btn v-if="order.status === 'completed'" color="success" variant="outlined" block size="large"
          @click="scrollToESIM">
          查看 eSIM 信息
        </v-btn>

        <v-btn color="grey" variant="outlined" block size="large" @click="goBack">
          返回订单列表
        </v-btn>
      </div>
    </div>

    <!-- 退款确认对话框 -->
    <v-dialog v-model="showRefundDialog" max-width="400" persistent>
      <v-card>
        <v-card-title class="dialog-title">
          <v-icon start color="warning">mdi-alert</v-icon>
          确认退款
        </v-card-title>
        
        <v-card-text>
          <p>您确定要申请退款吗？退款后订单将被取消。</p>
          <v-textarea
            v-model="refundReason"
            label="退款原因（可选）"
            placeholder="请输入退款原因..."
            rows="3"
            variant="outlined"
            class="mt-4"
          />
        </v-card-text>
        
        <v-card-actions>
          <v-spacer />
          <v-btn variant="outlined" @click="showRefundDialog = false" :disabled="isRefunding">
            取消
          </v-btn>
          <v-btn color="warning" variant="elevated" @click="confirmRefund" :loading="isRefunding">
            确认退款
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </PageWrapper>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import { useESIMStore } from '@/stores/esim'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'
import { orderApi } from '@/services/api'
import PageWrapper from '@/components/layout/PageWrapper.vue'

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
const isExportingAll = ref(false)
const showRefundDialog = ref(false)
const refundReason = ref('')
const isRefunding = ref(false)

// 计算属性
const orderId = computed(() => props.id || route.params['id'] as string)

const order = computed(() => {
  if (orderId.value) {
    return ordersStore.currentOrder
  }
  return null
})

const esimList = computed(() => order.value?.esims || [])

const statusIcon = computed(() => {
  if (!order.value) return 'mdi-help-circle'

  const iconMap: Record<string, string> = {
    pending: 'mdi-clock-outline',
    paid: 'mdi-check-circle',
    processing: 'mdi-cog',
    completed: 'mdi-check-circle',
    cancelled: 'mdi-close-circle',
    failed: 'mdi-alert-circle'
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
    case 'pending':
      return '请完成支付以激活您的 eSIM'
    case 'paid':
      return '正在处理您的订单'
    case 'processing':
      return '正在为您准备 eSIM'
    case 'completed':
      return 'eSIM 已准备就绪'
    case 'cancelled':
      return '订单已取消'
    case 'failed':
      return '订单处理失败'
    default:
      return ''
  }
})

const canRetryPayment = computed(() => {
  return order.value?.status === 'pending'
})

const canRefund = computed(() => {
  return order.value?.status === 'paid' || order.value?.status === 'completed'
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

const copyText = async (text: string, label: string) => {
  try {
    await navigator.clipboard.writeText(text)
    appStore.showNotification({
      type: 'success',
      message: `${label}已复制到剪贴板`,
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

const formatAmount = (amount: string): string => {
  return `$${parseFloat(amount).toFixed(2)}`
}

const scrollToESIM = () => {
  const esimElement = document.getElementById('esim-info')
  if (esimElement) {
    esimElement.scrollIntoView({ behavior: 'smooth' })
  }
}

const goBack = () => {
  router.push({ name: 'Orders' })
}

// 批量导出 PDF
const exportAllPDF = async () => {
  if (!order.value || isExportingAll.value) return

  isExportingAll.value = true
  
  try {
    await esimStore.exportAllPDF(parseInt(order.value.id))
    
    appStore.showNotification({
      type: 'success',
      message: '批量 PDF 导出成功',
      duration: 2000
    })
    
    telegramService.impactFeedback('medium')
  } catch (error) {
    appStore.showNotification({
      type: 'error',
      message: error instanceof Error ? error.message : '批量 PDF 导出失败',
      duration: 3000
    })
  } finally {
    isExportingAll.value = false
  }
}

// 确认退款
const confirmRefund = async () => {
  if (!order.value || isRefunding.value) return

  isRefunding.value = true
  
  try {
    const result = await orderApi.refundOrder(order.value.id, refundReason.value || '用户申请退款')
    
    if (result.success) {
      appStore.showNotification({
        type: 'success',
        message: '退款申请已提交',
        duration: 3000
      })
      
      telegramService.impactFeedback('heavy')
      
      // 关闭对话框并刷新订单详情
      showRefundDialog.value = false
      refundReason.value = ''
      await fetchOrderDetail()
    } else {
      throw new Error(result.message || '退款申请失败')
    }
  } catch (error) {
    appStore.showNotification({
      type: 'error',
      message: error instanceof Error ? error.message : '退款申请失败',
      duration: 3000
    })
  } finally {
    isRefunding.value = false
  }
}

// 刷新页面数据
const refreshData = async () => {
  if (!order.value) return
  
  try {
    // 清除 eSIM 相关缓存
    esimStore.clearCache(parseInt(order.value['id'] as string))
    
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
    padding: 16px 0;
    border-bottom: 1px solid rgba(var(--v-theme-outline), 0.12);

    &:last-child {
      border-bottom: none;
    }

    .esim-header {
      margin-bottom: 12px;
    }

    .esim-details {
      .detail-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 8px 0;

        .detail-label {
          font-size: 0.875rem;
          color: rgba(var(--v-theme-on-surface), 0.6);
          font-weight: 500;
        }

        .detail-value {
          font-size: 0.875rem;
          font-weight: 500;
          color: rgb(var(--v-theme-on-surface));
          font-family: 'Courier New', monospace;
          word-break: break-all;
        }

        .detail-value-container {
          display: flex;
          align-items: center;
          gap: 8px;
          flex: 1;
          justify-content: flex-end;
        }
      }
    }

    .qr-code-section {
      text-align: center;
      margin-top: 16px;

      .qr-code-image {
        width: 200px;
        height: 200px;
        border: 1px solid rgba(var(--v-theme-outline), 0.2);
        border-radius: 8px;
        background: white;
      }
    }
  }
  }

  .cost-info {
    .cost-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;

      .cost-label {
        font-size: 0.875rem;
        color: rgba(var(--v-theme-on-surface), 0.6);
        font-weight: 500;
      }

      .cost-value {
        font-size: 0.875rem;
        font-weight: 500;
        color: rgb(var(--v-theme-on-surface));

        &.price {
          color: rgb(var(--v-theme-primary));
          font-weight: 600;
          font-size: 1.1rem;
        }
      }

      &.total {
        padding-top: 12px;

        .cost-label {
          font-size: 1rem;
          font-weight: 600;
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