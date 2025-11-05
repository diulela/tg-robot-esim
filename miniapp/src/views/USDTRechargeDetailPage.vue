<template>
  <div class="usdt-recharge-detail-page">
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载订单信息...</p>
    </div>

    <div v-else-if="order" class="order-content">
      <!-- 订单状态 -->
      <div class="status-section">
        <div :class="['status-badge', statusClass]">
          <span class="status-icon">{{ statusIcon }}</span>
          <span class="status-text">{{ statusText }}</span>
        </div>
        <div v-if="order.status === 'pending'" class="countdown-container">
          <div class="countdown-label">订单有效期</div>
          <div class="countdown-time">{{ countdownText }}</div>
        </div>
      </div>

      <!-- 充值信息 -->
      <div class="recharge-info">
        <h3 class="section-title">充值信息</h3>

        <!-- 精确金额 -->
        <div class="info-item">
          <div class="info-label">转账金额</div>
          <div class="info-value amount-value">
            <span class="amount-text">{{ order.exact_amount }} USDT</span>
            <button @click="copyAmount" class="copy-btn">
              <span class="copy-icon">📋</span>
              复制
            </button>
          </div>
          <div class="info-note">请务必转账此精确金额，多转或少转都无法到账</div>
        </div>

        <!-- 收款地址 -->
        <div class="info-item">
          <div class="info-label">收款地址</div>
          <div class="info-value address-value">
            <span class="address-text">{{ order.wallet_address }}</span>
            <button @click="copyAddress" class="copy-btn">
              <span class="copy-icon">📋</span>
              复制
            </button>
          </div>
          <div class="info-note">TRON (TRC20) 网络地址</div>
        </div>
        <!-- 订单信息 -->
        <div class="order-info">
          <div class="order-item">
            <span class="order-label">订单号</span>
            <span class="order-value">{{ order.order_no }}</span>
          </div>
          <div class="order-item">
            <span class="order-label">创建时间</span>
            <span class="order-value">{{ formatTime(order.created_at) }}</span>
          </div>
          <div v-if="order.confirmed_at" class="order-item">
            <span class="order-label">到账时间</span>
            <span class="order-value">{{ formatTime(order.confirmed_at) }}</span>
          </div>
          <div v-if="order.tx_hash" class="order-item">
            <span class="order-label">交易哈希</span>
            <span class="order-value hash-value">
              <span class="hash-text">{{ shortHash(order.tx_hash) }}</span>
              <button @click="viewTransaction" class="view-btn">查看</button>
            </span>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="actions-section">
        <button v-if="order.status === 'pending'" @click="checkStatus" :disabled="checking" class="check-btn">
          <span v-if="checking">检查中...</span>
          <span v-else>我已转账，检查状态</span>
        </button>

        <button v-if="order.status === 'confirmed'" @click="goToWallet" class="wallet-btn">
          查看钱包
        </button>

      </div>

      <!-- 充值说明 -->
      <div class="instructions">
        <h4>转账步骤</h4>
        <ol>
          <li>复制上方的精确金额和收款地址</li>
          <li>打开您的 USDT 钱包应用</li>
          <li>选择 TRON (TRC20) 网络</li>
          <li>输入收款地址和精确金额</li>
          <li>确认并发送转账</li>
          <li>转账完成后点击"我已转账"按钮</li>
        </ol>

        <div class="warning-box">
          <h5>⚠️ 重要提醒</h5>
          <ul>
            <li>必须使用 TRON (TRC20) 网络，其他网络无法到账</li>
            <li>转账金额必须与显示的精确金额完全一致</li>
            <li>订单有效期为30分钟，过期后需重新创建</li>
            <li>到账需要19个区块确认，约10-30分钟</li>
          </ul>
        </div>
      </div>
    </div>

    <div v-else class="error-container">
      <div class="error-icon">❌</div>
      <div class="error-message">{{ errorMessage }}</div>
      <button @click="goBack" class="back-btn">返回</button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import api from '@/services/api'
import QRCode from 'qrcode'

export default {
  name: 'USDTRechargeDetailPage',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const appStore = useAppStore()

    const order = ref(null)
    const loading = ref(true)
    const checking = ref(false)
    const errorMessage = ref('')
    const qrCanvas = ref(null)
    const countdownTimer = ref(null)
    const remainingTime = ref(0)

    // 计算属性
    const statusClass = computed(() => {
      if (!order.value) return ''
      switch (order.value.status) {
        case 'pending': return 'status-pending'
        case 'confirmed': return 'status-confirmed'
        case 'expired': return 'status-expired'
        default: return 'status-pending'
      }
    })

    const statusIcon = computed(() => {
      if (!order.value) return ''
      switch (order.value.status) {
        case 'pending': return '⏳'
        case 'confirmed': return '✅'
        case 'expired': return '⏰'
        default: return '⏳'
      }
    })

    const statusText = computed(() => {
      if (!order.value) return ''
      switch (order.value.status) {
        case 'pending': return '等待转账'
        case 'confirmed': return '充值成功'
        case 'expired': return '订单已过期'
        default: return '等待转账'
      }
    })

    const countdownText = computed(() => {
      if (remainingTime.value <= 0) return '已过期'

      const minutes = Math.floor(remainingTime.value / 60)
      const seconds = remainingTime.value % 60
      return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
    })

    // 方法
    const loadOrderDetail = async () => {
      try {
        loading.value = true
        const orderNo = route.params.orderNo

        if (!orderNo) {
          throw new Error('订单号不能为空')
        }

        const response = await api.wallet.getRechargeOrder(orderNo)
        order.value = response

        // 计算剩余时间
        if (response.status === 'pending') {
          const expiresAt = new Date(response.expires_at).getTime()
          const now = Date.now()
          remainingTime.value = Math.max(0, Math.floor((expiresAt - now) / 1000))

          // 启动倒计时
          startCountdown()
        }

      } catch (error) {
        console.error('加载订单详情失败:', error)
        errorMessage.value = error.message || '加载订单详情失败'
      } finally {
        loading.value = false
      }
    }

    const startCountdown = () => {
      if (countdownTimer.value) {
        clearInterval(countdownTimer.value)
      }

      countdownTimer.value = setInterval(() => {
        if (remainingTime.value > 0) {
          remainingTime.value--
        } else {
          clearInterval(countdownTimer.value)
          // 订单过期，刷新状态
          loadOrderDetail()
        }
      }, 1000)
    }


    const copyAmount = async () => {
      try {
        await navigator.clipboard.writeText(order.value.exact_amount)
        appStore.showSuccess('金额已复制到剪贴板')
      } catch (error) {
        console.error('复制失败:', error)
        appStore.showError('复制失败，请手动复制')
      }
    }

    const copyAddress = async () => {
      try {
        await navigator.clipboard.writeText(order.value.wallet_address)
        appStore.showSuccess('地址已复制到剪贴板')
      } catch (error) {
        console.error('复制失败:', error)
        appStore.showError('复制失败，请手动复制')
      }
    }

    const checkStatus = async () => {
      if (checking.value) return

      try {
        checking.value = true
        const response = await api.wallet.checkRechargeStatus(order.value.order_no)

        // 更新订单状态
        order.value.status = response.status
        order.value.tx_hash = response.tx_hash
        order.value.confirmations = response.confirmations
        order.value.confirmed_at = response.confirmed_at

        if (response.status === 'confirmed') {
          appStore.showSuccess('充值成功！余额已更新')
          clearInterval(countdownTimer.value)
        } else if (response.status === 'expired') {
          appStore.showWarning('订单已过期，请重新创建充值订单')
          clearInterval(countdownTimer.value)
        } else {
          appStore.showInfo('暂未检测到转账，请稍后再试')
        }

      } catch (error) {
        console.error('检查状态失败:', error)
        appStore.showError(error.message || '检查状态失败，请稍后重试')
      } finally {
        checking.value = false
      }
    }

    const viewTransaction = () => {
      if (order.value.tx_hash) {
        // 打开 TRON 区块链浏览器
        const url = `https://tronscan.org/#/transaction/${order.value.tx_hash}`
        window.open(url, '_blank')
      }
    }

    const goToWallet = () => {
      router.push({ name: 'Wallet' })
    }

    const createNewOrder = () => {
      router.push({ name: 'USDTRecharge' })
    }

    const goBack = () => {
      router.back()
    }

    const formatTime = (timeStr) => {
      const date = new Date(timeStr)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }

    const shortHash = (hash) => {
      if (!hash) return ''
      return `${hash.slice(0, 8)}...${hash.slice(-8)}`
    }

    // 生命周期
    onMounted(() => {
      loadOrderDetail()
    })

    onUnmounted(() => {
      if (countdownTimer.value) {
        clearInterval(countdownTimer.value)
      }
    })

    return {
      order,
      loading,
      checking,
      errorMessage,
      qrCanvas,
      remainingTime,
      statusClass,
      statusIcon,
      statusText,
      countdownText,
      loadOrderDetail,
      copyAmount,
      copyAddress,
      checkStatus,
      viewTransaction,
      goToWallet,
      createNewOrder,
      goBack,
      formatTime,
      shortHash
    }
  }
}
</script>

<style scoped>
.usdt-recharge-detail-page {
  padding: 16px;
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
}

.loading-container,
.error-container {
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
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-message {
  font-size: 16px;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 24px;
}

.back-btn {
  padding: 12px 24px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.status-section {
  text-align: center;
  margin-bottom: 24px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 20px;
  font-weight: bold;
  margin-bottom: 16px;
}

.status-pending {
  background: rgba(255, 193, 7, 0.1);
  color: #ffc107;
  border: 2px solid #ffc107;
}

.status-confirmed {
  background: rgba(40, 167, 69, 0.1);
  color: #28a745;
  border: 2px solid #28a745;
}

.status-expired {
  background: rgba(220, 53, 69, 0.1);
  color: #dc3545;
  border: 2px solid #dc3545;
}

.countdown-container {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.countdown-label {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 8px;
}

.countdown-time {
  font-size: 24px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  font-family: 'Courier New', monospace;
}

.recharge-info {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.info-item {
  margin-bottom: 24px;
  padding: 16px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
}

.info-label {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 8px;
}

.info-value {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.amount-text,
.address-text {
  flex: 1;
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  word-break: break-all;
  font-family: 'Courier New', monospace;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  white-space: nowrap;
}

.copy-icon {
  font-size: 14px;
}

.info-note {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  line-height: 1.4;
}

.qr-section {
  text-align: center;
  margin-bottom: 24px;
  padding: 20px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
}

.qr-container {
  display: inline-block;
  padding: 16px;
  background: white;
  border-radius: 12px;
  margin-bottom: 12px;
}

.qr-code {
  display: block;
}

.qr-label {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
}

.order-info {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.order-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.order-item:last-child {
  border-bottom: none;
}

.order-label {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
}

.order-value {
  font-size: 14px;
  color: var(--tg-theme-text-color, #000000);
  text-align: right;
}

.hash-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hash-text {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.view-btn {
  padding: 4px 8px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 4px;
  font-size: 10px;
  cursor: pointer;
}

.actions-section {
  margin-bottom: 24px;
}

.check-btn,
.wallet-btn,
.retry-btn {
  width: 100%;
  padding: 16px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
}

.check-btn {
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
}

.wallet-btn {
  background: #28a745;
  color: white;
}

.retry-btn {
  background: #ffc107;
  color: #000000;
}

.check-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.instructions {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.instructions h4 {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 12px 0;
}

.instructions ol {
  margin: 0 0 16px 0;
  padding-left: 20px;
}

.instructions li {
  font-size: 14px;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 8px;
  line-height: 1.4;
}

.warning-box {
  background: rgba(255, 193, 7, 0.1);
  border: 1px solid #ffc107;
  border-radius: 8px;
  padding: 12px;
}

.warning-box h5 {
  font-size: 14px;
  font-weight: bold;
  color: #ffc107;
  margin: 0 0 8px 0;
}

.warning-box ul {
  margin: 0;
  padding-left: 16px;
}

.warning-box li {
  font-size: 12px;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
  line-height: 1.4;
}

@media (max-width: 480px) {
  .info-value {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .copy-btn {
    align-self: flex-end;
  }
}
</style>