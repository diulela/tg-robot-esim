<template>
  <div class="wallet-page">
    <!-- 钱包余额卡片 -->
    <div class="wallet-balance-card">
      <div class="balance-header">
        <h2 class="balance-title">钱包余额</h2>
        <button @click="refreshBalance" class="refresh-btn" :disabled="loading">
          <span class="refresh-icon">🔄</span>
        </button>
      </div>
      
      <div class="balance-amount">
        <span class="currency">¥</span>
        <span class="amount">{{ formatAmount(walletBalance) }}</span>
      </div>
      
      <div class="balance-actions">
        <button @click="goToRecharge" class="recharge-btn">
          充值
        </button>
        <button @click="showWithdrawDialog" class="withdraw-btn">
          提现
        </button>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="quick-actions">
      <div class="action-item" @click="goToRecharge">
        <div class="action-icon">💳</div>
        <span class="action-text">充值</span>
      </div>
      
      <div class="action-item" @click="showTransactionHistory">
        <div class="action-icon">📊</div>
        <span class="action-text">交易记录</span>
      </div>
      
      <div class="action-item" @click="showWithdrawDialog">
        <div class="action-icon">💰</div>
        <span class="action-text">提现</span>
      </div>
      
      <div class="action-item" @click="showSettings">
        <div class="action-icon">⚙️</div>
        <span class="action-text">设置</span>
      </div>
    </div>

    <!-- 最近交易 -->
    <div class="recent-transactions">
      <div class="section-header">
        <h3 class="section-title">最近交易</h3>
        <button @click="showAllTransactions" class="view-all-btn">查看全部</button>
      </div>
      
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>正在加载交易记录...</p>
      </div>
      
      <div v-else-if="recentTransactions.length > 0" class="transactions-list">
        <div 
          v-for="transaction in recentTransactions" 
          :key="transaction.id"
          class="transaction-item"
          @click="showTransactionDetail(transaction)"
        >
          <div class="transaction-icon">
            <span v-if="transaction.type === 'recharge'">💳</span>
            <span v-else-if="transaction.type === 'purchase'">🛒</span>
            <span v-else-if="transaction.type === 'withdraw'">💰</span>
            <span v-else>💸</span>
          </div>
          
          <div class="transaction-info">
            <div class="transaction-title">{{ getTransactionTitle(transaction) }}</div>
            <div class="transaction-time">{{ formatTime(transaction.createdAt) }}</div>
          </div>
          
          <div class="transaction-amount" :class="getAmountClass(transaction)">
            {{ getAmountText(transaction) }}
          </div>
        </div>
      </div>
      
      <div v-else class="empty-transactions">
        <div class="empty-icon">📝</div>
        <p>暂无交易记录</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

export default {
  name: 'WalletPage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()
    
    const loading = ref(false)
    const walletBalance = ref(0)
    const recentTransactions = ref([])
    
    // 方法
    const loadWalletData = async () => {
      loading.value = true
      try {
        // 模拟加载钱包数据
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        walletBalance.value = 128.50
        recentTransactions.value = [
          {
            id: '1',
            type: 'purchase',
            title: '购买 eSIM 套餐',
            amount: -29.90,
            status: 'completed',
            createdAt: new Date(Date.now() - 2 * 60 * 60 * 1000)
          },
          {
            id: '2',
            type: 'recharge',
            title: '钱包充值',
            amount: 100.00,
            status: 'completed',
            createdAt: new Date(Date.now() - 24 * 60 * 60 * 1000)
          }
        ]
      } catch (error) {
        console.error('加载钱包数据失败:', error)
        appStore.showError('加载钱包数据失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
    
    const refreshBalance = async () => {
      await loadWalletData()
      appStore.showSuccess('余额已刷新')
    }
    
    const formatAmount = (amount) => {
      return amount.toFixed(2)
    }
    
    const formatTime = (date) => {
      const now = new Date()
      const diff = now - date
      const hours = Math.floor(diff / (1000 * 60 * 60))
      const days = Math.floor(hours / 24)
      
      if (days > 0) {
        return `${days}天前`
      } else if (hours > 0) {
        return `${hours}小时前`
      } else {
        return '刚刚'
      }
    }
    
    const getTransactionTitle = (transaction) => {
      return transaction.title || '未知交易'
    }
    
    const getAmountClass = (transaction) => {
      return transaction.amount > 0 ? 'amount-positive' : 'amount-negative'
    }
    
    const getAmountText = (transaction) => {
      const prefix = transaction.amount > 0 ? '+' : ''
      return `${prefix}¥${Math.abs(transaction.amount).toFixed(2)}`
    }
    
    const goToRecharge = () => {
      router.push({ name: 'WalletRecharge' })
    }
    
    const showWithdrawDialog = () => {
      appStore.showInfo('提现功能开发中')
    }
    
    const showTransactionHistory = () => {
      appStore.showInfo('交易历史功能开发中')
    }
    
    const showAllTransactions = () => {
      showTransactionHistory()
    }
    
    const showTransactionDetail = (transaction) => {
      appStore.showInfo(`交易详情: ${transaction.title}`)
    }
    
    const showSettings = () => {
      router.push({ name: 'Settings' })
    }
    
    // 生命周期
    onMounted(() => {
      loadWalletData()
    })
    
    return {
      loading,
      walletBalance,
      recentTransactions,
      refreshBalance,
      formatAmount,
      formatTime,
      getTransactionTitle,
      getAmountClass,
      getAmountText,
      goToRecharge,
      showWithdrawDialog,
      showTransactionHistory,
      showAllTransactions,
      showTransactionDetail,
      showSettings
    }
  }
}
</script>

<style scoped>
.wallet-page {
  padding: 16px;
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
}

.wallet-balance-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 24px;
  color: white;
  margin-bottom: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.balance-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.balance-title {
  font-size: 16px;
  font-weight: 500;
  margin: 0;
  opacity: 0.9;
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 8px;
  padding: 8px;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
}

.refresh-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-icon {
  font-size: 16px;
}

.balance-amount {
  display: flex;
  align-items: baseline;
  margin-bottom: 24px;
}

.currency {
  font-size: 24px;
  font-weight: 500;
  margin-right: 4px;
}

.amount {
  font-size: 36px;
  font-weight: bold;
}

.balance-actions {
  display: flex;
  gap: 12px;
}

.recharge-btn,
.withdraw-btn {
  flex: 1;
  padding: 12px 16px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recharge-btn {
  background: white;
  color: #667eea;
}

.withdraw-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.recharge-btn:hover,
.withdraw-btn:hover {
  transform: translateY(-1px);
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 8px;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.action-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.action-text {
  font-size: 12px;
  color: var(--tg-theme-text-color, #000000);
  text-align: center;
}

.recent-transactions {
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  padding: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0;
}

.view-all-btn {
  background: none;
  border: none;
  color: var(--tg-theme-button-color, #0088cc);
  font-size: 14px;
  cursor: pointer;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  text-align: center;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 2px solid var(--tg-theme-hint-color, #e0e0e0);
  border-top: 2px solid var(--tg-theme-button-color, #0088cc);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.transactions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.transaction-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.transaction-item:hover {
  background: var(--tg-theme-hint-color, #e0e0e0);
}

.transaction-icon {
  font-size: 20px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--tg-theme-bg-color, #ffffff);
  border-radius: 50%;
}

.transaction-info {
  flex: 1;
}

.transaction-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.transaction-time {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.transaction-amount {
  font-size: 14px;
  font-weight: bold;
}

.amount-positive {
  color: #4caf50;
}

.amount-negative {
  color: #f44336;
}

.empty-transactions {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.empty-transactions p {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0;
}

@media (max-width: 480px) {
  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>