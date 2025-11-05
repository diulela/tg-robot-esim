<template>
  <div class="usdt-recharge-history-page">
    
    <!-- 加载状态 -->
    <div v-if="loading && orders.length === 0" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载充值历史...</p>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!loading && orders.length === 0" class="empty-container">
      <div class="empty-icon">📋</div>
      <div class="empty-text">暂无充值记录</div>
      <button @click="goToRecharge" class="recharge-btn">
        立即充值
      </button>
    </div>

    <!-- 充值记录列表 -->
    <div v-else class="orders-container">

      <!-- 订单列表 -->
      <div class="orders-list">
        <div 
          v-for="order in orders" 
          :key="order.order_no"
          @click="viewOrderDetail(order)"
          class="order-item"
        >
          <div class="order-header">
            <div class="order-status">
              <span :class="['status-badge', getStatusClass(order.status)]">
                {{ getStatusText(order.status) }}
              </span>
            </div>
            <div class="order-time">
              {{ formatTime(order.created_at) }}
            </div>
          </div>
          
          <div class="order-content">
            <div class="order-info">
              <div class="order-amount">
                <span class="amount-label">充值金额</span>
                <span class="amount-value">{{ order.amount }} USDT</span>
              </div>
              <div class="order-no">
                <span class="no-label">订单号</span>
                <span class="no-value">{{ order.order_no }}</span>
              </div>
            </div>
            
            <div class="order-actions">
              <button 
                v-if="order.status === 'confirmed' && order.tx_hash"
                @click.stop="viewTransaction(order.tx_hash)"
                class="action-btn view-tx-btn"
              >
                查看交易
              </button>
            </div>
          </div>
          
          <div v-if="order.confirmed_at" class="order-footer">
            <span class="confirmed-time">
              到账时间：{{ formatTime(order.confirmed_at) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore" class="load-more-container">
        <button 
          @click="loadMore" 
          :disabled="loadingMore"
          class="load-more-btn"
        >
          <span v-if="loadingMore">加载中...</span>
          <span v-else>加载更多</span>
        </button>
      </div>
      
      <div v-else-if="orders.length > 0" class="no-more-text">
        已显示全部记录
      </div>
    </div>

    <!-- 刷新按钮 -->
    <div class="refresh-container">
      <button @click="refreshList" :disabled="loading" class="refresh-btn">
        <span class="refresh-icon">🔄</span>
        刷新
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import api from '@/services/api'

export default {
  name: 'USDTRechargeHistoryPage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()
    
    const orders = ref([])
    const loading = ref(false)
    const loadingMore = ref(false)
    const total = ref(0)
    const currentPage = ref(0)
    const pageSize = ref(20)
    
    // 计算属性
    const hasMore = computed(() => {
      return orders.value.length < total.value
    })
    
    const successCount = computed(() => {
      return orders.value.filter(order => order.status === 'confirmed').length
    })
    
    // 方法
    const loadOrders = async (isLoadMore = false) => {
      try {
        if (isLoadMore) {
          loadingMore.value = true
        } else {
          loading.value = true
          currentPage.value = 0
          orders.value = []
        }
        
        const offset = currentPage.value * pageSize.value
        const response = await api.wallet.getRechargeHistory({
          limit: pageSize.value,
          offset: offset
        })
        
        if (isLoadMore) {
          orders.value.push(...response.orders)
        } else {
          orders.value = response.orders
        }
        
        total.value = response.total
        currentPage.value++
        
      } catch (error) {
        console.error('加载充值历史失败:', error)
        appStore.showError(error.message || '加载充值历史失败')
      } finally {
        loading.value = false
        loadingMore.value = false
      }
    }
    
    const loadMore = () => {
      if (!hasMore.value || loadingMore.value) return
      loadOrders(true)
    }
    
    const refreshList = () => {
      loadOrders(false)
    }
    
    const viewOrderDetail = (order) => {
      router.push({
        name: 'USDTRechargeDetail',
        params: { orderNo: order.order_no }
      })
    }
    
    const checkOrderStatus = async (order) => {
      try {
        const response = await api.wallet.checkRechargeStatus(order.order_no)
        
        // 更新本地订单状态
        const index = orders.value.findIndex(o => o.order_no === order.order_no)
        if (index !== -1) {
          orders.value[index].status = response.status
          orders.value[index].tx_hash = response.tx_hash
          orders.value[index].confirmed_at = response.confirmed_at
        }
        
        if (response.status === 'confirmed') {
          appStore.showSuccess('充值成功！')
        } else if (response.status === 'expired') {
          appStore.showWarning('订单已过期')
        } else {
          appStore.showInfo('暂未检测到转账')
        }
        
      } catch (error) {
        console.error('检查订单状态失败:', error)
        appStore.showError(error.message || '检查状态失败')
      }
    }
    
    const viewTransaction = (txHash) => {
      if (txHash) {
        const url = `https://tronscan.org/#/transaction/${txHash}`
        window.open(url, '_blank')
      }
    }
    
    const retryRecharge = (order) => {
      router.push({ name: 'USDTRecharge' })
    }
    
    const goToRecharge = () => {
      router.push({ name: 'USDTRecharge' })
    }
    
    const getStatusClass = (status) => {
      switch (status) {
        case 'pending': return 'status-pending'
        case 'confirmed': return 'status-confirmed'
        case 'expired': return 'status-expired'
        default: return 'status-pending'
      }
    }
    
    const getStatusText = (status) => {
      switch (status) {
        case 'pending': return '等待转账'
        case 'confirmed': return '充值成功'
        case 'expired': return '已过期'
        default: return '未知状态'
      }
    }
    
    const formatTime = (timeStr) => {
      const date = new Date(timeStr)
      const now = new Date()
      const diff = now.getTime() - date.getTime()
      
      // 如果是今天
      if (diff < 24 * 60 * 60 * 1000 && date.getDate() === now.getDate()) {
        return date.toLocaleTimeString('zh-CN', {
          hour: '2-digit',
          minute: '2-digit'
        })
      }
      
      // 如果是今年
      if (date.getFullYear() === now.getFullYear()) {
        return date.toLocaleDateString('zh-CN', {
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        })
      }
      
      // 其他情况显示完整日期
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      })
    }
    
    // 生命周期
    onMounted(() => {
      loadOrders()
    })
    
    return {
      orders,
      loading,
      loadingMore,
      total,
      hasMore,
      successCount,
      loadOrders,
      loadMore,
      refreshList,
      viewOrderDetail,
      checkOrderStatus,
      viewTransaction,
      retryRecharge,
      goToRecharge,
      getStatusClass,
      getStatusText,
      formatTime
    }
  }
}
</script>

<style scoped>
.usdt-recharge-history-page {
  padding: 16px;
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 24px 0;
  text-align: center;
}

.loading-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  text-align: center;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--tg-theme-hint-color, #e0e0e0);
  border-top: 4px solid var(--tg-theme-button-color, #0088cc);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-text {
  font-size: 16px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 24px;
}

.recharge-btn {
  padding: 12px 24px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stats-item {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
}

.stats-label {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 8px;
}

.stats-value {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
}

.orders-list {
  margin-bottom: 24px;
}

.order-item {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.order-item:hover {
  border-color: var(--tg-theme-button-color, #0088cc);
  transform: translateY(-1px);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
}

.status-pending {
  background: rgba(255, 193, 7, 0.2);
  color: #ffc107;
}

.status-confirmed {
  background: rgba(40, 167, 69, 0.2);
  color: #28a745;
}

.status-expired {
  background: rgba(220, 53, 69, 0.2);
  color: #dc3545;
}

.order-time {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.order-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.order-info {
  flex: 1;
}

.order-amount,
.order-no {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.amount-label,
.no-label {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
}

.amount-value {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
}

.no-value {
  font-size: 12px;
  color: var(--tg-theme-text-color, #000000);
  font-family: 'Courier New', monospace;
}

.order-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.view-tx-btn {
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
}

.check-btn {
  background: #ffc107;
  color: #000000;
}

.retry-btn {
  background: #dc3545;
  color: #ffffff;
}

.action-btn:hover {
  opacity: 0.8;
  transform: translateY(-1px);
}

.order-footer {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.confirmed-time {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.load-more-container {
  text-align: center;
  margin-bottom: 16px;
}

.load-more-btn {
  padding: 12px 24px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.load-more-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.no-more-text {
  text-align: center;
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 16px;
}

.refresh-container {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 100;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 24px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.2s ease;
}

.refresh-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-icon {
  font-size: 16px;
}

@media (max-width: 480px) {
  .order-content {
    flex-direction: column;
    gap: 12px;
  }
  
  .order-actions {
    flex-direction: row;
    justify-content: flex-end;
  }
  
  .stats-section {
    grid-template-columns: 1fr;
  }
}
</style>