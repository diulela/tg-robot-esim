<template>
  <div class="order-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">我的订单</h1>
    </div>

    <!-- 订单列表 -->
    <div class="orders-container">
      <!-- 加载状态 -->
      <div v-if="isLoading && orders.length === 0" class="loading-state">
        <v-progress-circular indeterminate color="primary" size="48" />
        <p class="loading-text">正在加载订单...</p>
      </div>

      <!-- 订单卡片列表 -->
      <div v-else-if="orders.length > 0" class="orders-list">
        <div
          v-for="order in orders"
          :key="order.id"
          class="order-card"
          @click="navigateToDetail(order.id)"
        >
          <!-- 产品信息区 -->
          <div class="product-section">
            <div class="product-info">
              <h3 class="product-name">{{ order.productName }}</h3>
              <div class="product-meta">
                <span class="quantity">× {{ order.quantity }}</span>
                <span class="data-size">
                  <v-icon size="14">mdi-database-outline</v-icon>
                  {{ extractDataSize(order.productName) }}
                </span>
                <span class="valid-days">
                  <v-icon size="14">mdi-clock-outline</v-icon>
                  {{ extractValidDays(order.productName) }}
                </span>
              </div>
            </div>
            <div class="product-price">${{ order.unitPrice }}</div>
          </div>

          <!-- 订单详情区 -->
          <div class="order-details">
            <div class="detail-row">
              <span class="detail-label">订单编号</span>
              <span class="detail-value">{{ order.orderNo }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">下单时间</span>
              <span class="detail-value">{{ formatDateTime(order.createdAt) }}</span>
              <span class="actual-payment">实付：<span class="amount">${{ order.totalAmount }}</span></span>
            </div>
            <div v-if="order.providerOrderId" class="detail-row">
              <span class="detail-label">ICCID</span>
              <span class="detail-value iccid">{{ formatICCID(order.providerOrderId) }}</span>
            </div>
          </div>

          <!-- 状态和操作区 -->
          <div class="order-footer">
            <div class="status-badge" :class="`status-${order.status}`">
              {{ getStatusText(order.status) }}
            </div>
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              @click.stop="navigateToDetail(order.id)"
            >
              查看详情
            </v-btn>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!isLoading" class="empty-state">
        <v-icon size="80" color="grey-lighten-2">mdi-receipt-text-outline</v-icon>
        <h3 class="empty-title">暂无订单</h3>
        <p class="empty-subtitle">您还没有任何订单记录</p>
        <v-btn
          color="primary"
          variant="elevated"
          size="large"
          class="mt-6"
          @click="navigateToProducts"
        >
          <v-icon start>mdi-shopping</v-icon>
          开始购买
        </v-btn>
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore && orders.length > 0" class="load-more">
        <v-btn
          :loading="isLoadingMore"
          variant="outlined"
          color="primary"
          block
          @click="loadMore"
        >
          加载更多
        </v-btn>
      </div>

      <!-- 错误提示 -->
      <v-snackbar
        v-model="showError"
        :timeout="3000"
        color="error"
        location="top"
      >
        {{ errorMessage }}
        <template #actions>
          <v-btn variant="text" @click="showError = false">关闭</v-btn>
        </template>
      </v-snackbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { orderApi } from '@/services/api/order'
import { telegramService } from '@/services/telegram'
import type { EsimOrder } from '@/types/esim-order'

// 路由
const router = useRouter()

// 响应式状态
const orders = ref<EsimOrder[]>([])
const isLoading = ref(false)
const isLoadingMore = ref(false)
const showError = ref(false)
const errorMessage = ref('')
const currentOffset = ref(0)
const pageSize = 20
const totalOrders = ref(0)

// 计算属性
const hasMore = computed(() => {
  return orders.value.length < totalOrders.value
})

// 方法
const loadOrders = async (append = false) => {
  if (append) {
    isLoadingMore.value = true
  } else {
    isLoading.value = true
    currentOffset.value = 0
  }

  try {
    const response = await orderApi.getEsimOrders({
      limit: pageSize,
      offset: currentOffset.value
    })

    if (append) {
      orders.value = [...orders.value, ...response.items]
    } else {
      orders.value = response.items
    }

    totalOrders.value = response.total
    currentOffset.value += response.items.length

    console.log('[OrderPage] 订单加载成功:', response.items.length, '个订单')
  } catch (error) {
    console.error('[OrderPage] 加载订单失败:', error)
    errorMessage.value = '加载订单失败，请重试'
    showError.value = true
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

const loadMore = async () => {
  if (isLoadingMore.value || !hasMore.value) return
  await loadOrders(true)
}

const navigateToDetail = (orderId: string) => {
  router.push({ name: 'OrderDetail', params: { id: orderId } })
  telegramService.impactFeedback('light')
}

const navigateToProducts = () => {
  router.push({ name: 'Products' })
  telegramService.impactFeedback('medium')
}

// 格式化日期时间
const formatDateTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 格式化 ICCID（显示部分）
const formatICCID = (iccid: string): string => {
  if (!iccid) return '未分配'
  // 如果太长，只显示前面部分
  return iccid.length > 20 ? iccid.substring(0, 20) : iccid
}

// 从产品名称提取流量大小
const extractDataSize = (productName: string): string => {
  const match = productName.match(/(\d+)(MB|GB)/i)
  return match ? match[0] : '200MB'
}

// 从产品名称提取有效天数
const extractValidDays = (productName: string): string => {
  const match = productName.match(/(\d+)天/)
  return match ? match[0] : '30天'
}

// 获取状态文本
const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    pending: '待支付',
    processing: '处理中',
    paid: '已付款',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

// 生命周期
onMounted(async () => {
  console.log('[OrderPage] 页面挂载')
  await loadOrders()
})
</script>

<style scoped lang="scss">
.order-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding-bottom: 80px;

  .page-header {
    padding: 20px 16px;
    text-align: center;

    .page-title {
      font-size: 1.5rem;
      font-weight: 600;
      color: white;
      margin: 0;
    }
  }

  .orders-container {
    padding: 0 16px 16px;

    .loading-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 80px 20px;
      gap: 16px;

      .loading-text {
        color: white;
        font-size: 0.875rem;
        margin: 0;
      }
    }

    .orders-list {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .order-card {
        background: white;
        border-radius: 12px;
        padding: 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        cursor: pointer;
        transition: all 0.2s ease;

        &:hover {
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
          transform: translateY(-2px);
        }

        &:active {
          transform: translateY(0);
        }

        .product-section {
          display: flex;
          justify-content: space-between;
          align-items: flex-start;
          margin-bottom: 16px;
          padding-bottom: 12px;
          border-bottom: 1px solid #f0f0f0;

          .product-info {
            flex: 1;

            .product-name {
              font-size: 0.9375rem;
              font-weight: 500;
              color: #333;
              margin: 0 0 8px 0;
              line-height: 1.4;
            }

            .product-meta {
              display: flex;
              align-items: center;
              gap: 12px;
              font-size: 0.8125rem;
              color: #666;

              .quantity {
                color: #999;
              }

              .data-size,
              .valid-days {
                display: flex;
                align-items: center;
                gap: 4px;

                .v-icon {
                  color: #999;
                }
              }
            }
          }

          .product-price {
            font-size: 1.125rem;
            font-weight: 600;
            color: #667eea;
            white-space: nowrap;
            margin-left: 12px;
          }
        }

        .order-details {
          margin-bottom: 16px;

          .detail-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 6px 0;
            font-size: 0.8125rem;

            .detail-label {
              color: #999;
              min-width: 70px;
            }

            .detail-value {
              flex: 1;
              color: #666;
              text-align: left;
              margin-left: 12px;

              &.iccid {
                font-family: monospace;
                font-size: 0.75rem;
              }
            }

            .actual-payment {
              color: #999;
              margin-left: auto;
              white-space: nowrap;

              .amount {
                color: #f56c6c;
                font-weight: 600;
                font-size: 0.9375rem;
              }
            }
          }
        }

        .order-footer {
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 12px;

          .status-badge {
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.75rem;
            font-weight: 500;
            white-space: nowrap;

            &.status-pending {
              background: #fff3e0;
              color: #f57c00;
            }

            &.status-processing {
              background: #e3f2fd;
              color: #1976d2;
            }

            &.status-paid {
              background: #f3e5f5;
              color: #7b1fa2;
            }

            &.status-completed {
              background: #e8f5e9;
              color: #388e3c;
            }

            &.status-failed {
              background: #ffebee;
              color: #d32f2f;
            }

            &.status-cancelled {
              background: #f5f5f5;
              color: #757575;
            }
          }

          .v-btn {
            flex-shrink: 0;
          }
        }
      }
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 80px 20px;
      text-align: center;

      .empty-title {
        font-size: 1.125rem;
        font-weight: 600;
        color: white;
        margin: 16px 0 8px;
      }

      .empty-subtitle {
        font-size: 0.875rem;
        color: rgba(255, 255, 255, 0.8);
        margin: 0;
      }
    }

    .load-more {
      margin-top: 16px;
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .order-page {
    .page-header {
      padding: 16px 12px;

      .page-title {
        font-size: 1.25rem;
      }
    }

    .orders-container {
      padding: 0 12px 12px;

      .orders-list {
        gap: 10px;

        .order-card {
          padding: 12px;

          .product-section {
            .product-info {
              .product-name {
                font-size: 0.875rem;
              }

              .product-meta {
                gap: 8px;
                font-size: 0.75rem;
              }
            }

            .product-price {
              font-size: 1rem;
            }
          }

          .order-details {
            .detail-row {
              font-size: 0.75rem;

              .detail-value.iccid {
                font-size: 0.6875rem;
              }

              .actual-payment .amount {
                font-size: 0.875rem;
              }
            }
          }
        }
      }
    }
  }
}

@media (min-width: 481px) {
  .order-page {
    .orders-container {
      max-width: 600px;
      margin: 0 auto;
    }
  }
}

// 深色主题适配
.v-theme--dark {
  .order-page {
    background: linear-gradient(135deg, #1a237e 0%, #4a148c 100%);
  }
}
</style>
