<template>
  <div class="wallet-recharge-page">
    <h1 class="page-title">钱包充值</h1>
    
    <!-- 充值金额选择 -->
    <div class="amount-section">
      <h3 class="section-title">选择充值金额</h3>
      <div class="amount-options">
        <button 
          v-for="option in amountOptions" 
          :key="option.value"
          @click="selectAmount(option.value)"
          :class="['amount-option', { active: selectedAmount === option.value }]"
        >
          <span class="amount-value">¥{{ option.value }}</span>
          <span v-if="option.bonus" class="amount-bonus">送¥{{ option.bonus }}</span>
        </button>
      </div>
      
      <div class="custom-amount">
        <input 
          v-model="customAmount" 
          type="number" 
          placeholder="自定义金额"
          class="custom-input"
          @input="selectCustomAmount"
          min="1"
          max="10000"
        />
      </div>
    </div>

    <!-- 支付方式选择 -->
    <div class="payment-section">
      <h3 class="section-title">选择支付方式</h3>
      <div class="payment-methods">
        <div 
          v-for="method in paymentMethods" 
          :key="method.id"
          @click="selectPaymentMethod(method.id)"
          :class="['payment-method', { active: selectedPaymentMethod === method.id }]"
        >
          <div class="method-icon">{{ method.icon }}</div>
          <div class="method-info">
            <div class="method-name">{{ method.name }}</div>
            <div class="method-desc">{{ method.description }}</div>
          </div>
          <div class="method-check">
            <span v-if="selectedPaymentMethod === method.id">✓</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 充值信息确认 -->
    <div class="confirm-section">
      <div class="confirm-item">
        <span class="confirm-label">充值金额</span>
        <span class="confirm-value">¥{{ finalAmount }}</span>
      </div>
      <div v-if="bonusAmount > 0" class="confirm-item">
        <span class="confirm-label">赠送金额</span>
        <span class="confirm-value bonus">+¥{{ bonusAmount }}</span>
      </div>
      <div class="confirm-item total">
        <span class="confirm-label">到账金额</span>
        <span class="confirm-value">¥{{ totalAmount }}</span>
      </div>
    </div>

    <!-- 充值按钮 -->
    <div class="recharge-actions">
      <button 
        @click="processRecharge" 
        :disabled="!canRecharge || processing"
        class="recharge-btn"
      >
        <span v-if="processing">处理中...</span>
        <span v-else>立即充值 ¥{{ finalAmount }}</span>
      </button>
    </div>

    <!-- 充值说明 -->
    <div class="recharge-notes">
      <h4>充值说明</h4>
      <ul>
        <li>充值金额将实时到账</li>
        <li>单次充值金额范围：¥1 - ¥10,000</li>
        <li>充值过程中请勿关闭页面</li>
        <li>如遇问题请联系客服</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

export default {
  name: 'WalletRechargePage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()
    
    const selectedAmount = ref(0)
    const customAmount = ref('')
    const selectedPaymentMethod = ref('alipay')
    const processing = ref(false)
    
    // 充值金额选项
    const amountOptions = [
      { value: 50, bonus: 0 },
      { value: 100, bonus: 5 },
      { value: 200, bonus: 15 },
      { value: 500, bonus: 50 },
      { value: 1000, bonus: 120 },
      { value: 2000, bonus: 300 }
    ]
    
    // 支付方式
    const paymentMethods = [
      {
        id: 'alipay',
        name: '支付宝',
        description: '推荐使用，到账快速',
        icon: '💙'
      },
      {
        id: 'wechat',
        name: '微信支付',
        description: '微信扫码支付',
        icon: '💚'
      },
      {
        id: 'usdt',
        name: 'USDT-TRC20',
        description: '区块链支付，手续费低',
        icon: '₿'
      }
    ]
    
    // 计算属性
    const finalAmount = computed(() => {
      if (customAmount.value) {
        return parseFloat(customAmount.value) || 0
      }
      return selectedAmount.value
    })
    
    const bonusAmount = computed(() => {
      if (customAmount.value) {
        return 0 // 自定义金额没有赠送
      }
      const option = amountOptions.find(opt => opt.value === selectedAmount.value)
      return option ? option.bonus : 0
    })
    
    const totalAmount = computed(() => {
      return finalAmount.value + bonusAmount.value
    })
    
    const canRecharge = computed(() => {
      return finalAmount.value > 0 && selectedPaymentMethod.value
    })
    
    // 方法
    const selectAmount = (amount) => {
      selectedAmount.value = amount
      customAmount.value = ''
    }
    
    const selectCustomAmount = () => {
      selectedAmount.value = 0
    }
    
    const selectPaymentMethod = (methodId) => {
      selectedPaymentMethod.value = methodId
    }
    
    const processRecharge = async () => {
      if (!canRecharge.value || processing.value) return
      
      processing.value = true
      
      try {
        // 模拟支付处理
        await new Promise(resolve => setTimeout(resolve, 2000))
        
        // 根据支付方式处理
        if (selectedPaymentMethod.value === 'usdt') {
          // USDT 支付跳转到区块链支付页面
          appStore.showInfo('正在跳转到 USDT 支付页面...')
        } else {
          // 其他支付方式
          appStore.showSuccess(`充值成功！到账金额：¥${totalAmount.value}`)
          router.push({ name: 'Wallet' })
        }
      } catch (error) {
        console.error('充值失败:', error)
        appStore.showError('充值失败，请稍后重试')
      } finally {
        processing.value = false
      }
    }
    
    return {
      selectedAmount,
      customAmount,
      selectedPaymentMethod,
      processing,
      amountOptions,
      paymentMethods,
      finalAmount,
      bonusAmount,
      totalAmount,
      canRecharge,
      selectAmount,
      selectCustomAmount,
      selectPaymentMethod,
      processRecharge
    }
  }
}
</script>

<style scoped>
.wallet-recharge-page {
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

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.amount-section,
.payment-section,
.confirm-section {
  margin-bottom: 24px;
}

.amount-options {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.amount-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  border: 2px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  background: var(--tg-theme-bg-color, #ffffff);
  cursor: pointer;
  transition: all 0.2s ease;
}

.amount-option.active {
  border-color: var(--tg-theme-button-color, #0088cc);
  background: rgba(0, 136, 204, 0.1);
}

.amount-value {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.amount-bonus {
  font-size: 12px;
  color: #ff4757;
  font-weight: 500;
}

.custom-amount {
  margin-top: 12px;
}

.custom-input {
  width: 100%;
  padding: 16px;
  border: 2px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  font-size: 16px;
  background: var(--tg-theme-bg-color, #ffffff);
  color: var(--tg-theme-text-color, #000000);
  text-align: center;
}

.custom-input:focus {
  outline: none;
  border-color: var(--tg-theme-button-color, #0088cc);
}

.payment-methods {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.payment-method {
  display: flex;
  align-items: center;
  padding: 16px;
  border: 2px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  background: var(--tg-theme-bg-color, #ffffff);
  cursor: pointer;
  transition: all 0.2s ease;
}

.payment-method.active {
  border-color: var(--tg-theme-button-color, #0088cc);
  background: rgba(0, 136, 204, 0.1);
}

.method-icon {
  font-size: 24px;
  margin-right: 16px;
}

.method-info {
  flex: 1;
}

.method-name {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.method-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.method-check {
  font-size: 18px;
  color: var(--tg-theme-button-color, #0088cc);
  font-weight: bold;
}

.confirm-section {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.confirm-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.confirm-item.total {
  border-top: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  margin-top: 8px;
  padding-top: 16px;
  font-weight: bold;
}

.confirm-label {
  font-size: 14px;
  color: var(--tg-theme-text-color, #000000);
}

.confirm-value {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
}

.confirm-value.bonus {
  color: #ff4757;
}

.recharge-actions {
  margin-bottom: 24px;
}

.recharge-btn {
  width: 100%;
  padding: 16px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
}

.recharge-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.recharge-btn:not(:disabled):hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.recharge-notes {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.recharge-notes h4 {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 12px 0;
}

.recharge-notes ul {
  margin: 0;
  padding-left: 16px;
}

.recharge-notes li {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 4px;
  line-height: 1.4;
}

@media (max-width: 480px) {
  .amount-options {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .amount-option {
    padding: 12px 8px;
  }
  
  .amount-value {
    font-size: 16px;
  }
}
</style>