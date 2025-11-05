<template>
  <div class="usdt-recharge-page">
    <h1 class="page-title">USDT 充值</h1>
    <!-- 快捷金额选择 -->
    <div class="quick-amounts">
      <h4 class="quick-title">快捷选择</h4>
      <div class="quick-options">
        <button v-for="quickAmount in quickAmounts" :key="quickAmount" @click="selectQuickAmount(quickAmount)"
          :class="['quick-option', { active: amount == quickAmount }]">
          {{ quickAmount }} USDT
        </button>
      </div>
    </div>

    <!-- 充值金额输入 -->
    <div class="amount-section">
      <h3 class="section-title">充值金额</h3>
      <div class="amount-input-container">
        <input v-model="amount" type="number" placeholder="请输入充值金额" class="amount-input" @input="validateAmount"
          min="10" max="10000" step="0.01" />
        <span class="currency-label">USDT</span>
      </div>
      <div class="amount-tips">
        <p class="tip-text">• 最小充值金额：10 USDT</p>
        <p class="tip-text">• 最大充值金额：10,000 USDT</p>
        <p class="tip-text">• 仅支持 USDT-TRC20 网络</p>
      </div>
      <div v-if="amountError" class="error-message">
        {{ amountError }}
      </div>
    </div>

    <!-- 充值按钮 -->
    <div class="recharge-actions">
      <button @click="createRechargeOrder" :disabled="!canRecharge || loading" class="recharge-btn">
        <span v-if="loading">创建订单中...</span>
        <span v-else>创建充值订单</span>
      </button>
    </div>

    <!-- 充值说明 -->
    <div class="recharge-info">
      <h4 class="info-title">充值说明</h4>
      <div class="info-content">
        <div class="info-item">
          <span class="info-icon">🔒</span>
          <div class="info-text">
            <div class="info-label">安全保障</div>
            <div class="info-desc">采用区块链技术，资金安全有保障</div>
          </div>
        </div>
        <div class="info-item">
          <span class="info-icon">⚡</span>
          <div class="info-text">
            <div class="info-label">快速到账</div>
            <div class="info-desc">19个区块确认后自动到账，约10-30分钟</div>
          </div>
        </div>
        <div class="info-item">
          <span class="info-icon">💰</span>
          <div class="info-text">
            <div class="info-label">低手续费</div>
            <div class="info-desc">TRON网络手续费低，节省成本</div>
          </div>
        </div>
      </div>
    </div>



    <!-- 充值注意事项 -->
    <div class="recharge-notes">
      <h4>重要提醒</h4>
      <ul>
        <li>请务必使用 TRON (TRC20) 网络转账</li>
        <li>转账金额必须与订单金额完全一致</li>
        <li>充值订单有效期为30分钟</li>
        <li>请勿向充值地址转入其他代币</li>
        <li>如有问题请及时联系客服</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import api from '@/services/api'

export default {
  name: 'USDTRechargePage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()

    const amount = ref('')
    const amountError = ref('')
    const loading = ref(false)

    // 快捷金额选项
    const quickAmounts = [50, 100, 200, 500]

    // 计算属性
    const canRecharge = computed(() => {
      const amountNum = parseFloat(amount.value)
      return amountNum >= 10 && amountNum <= 10000 && !amountError.value
    })

    // 方法
    const validateAmount = () => {
      const amountNum = parseFloat(amount.value)
      amountError.value = ''

      if (!amount.value) {
        return
      }

      if (isNaN(amountNum) || amountNum <= 0) {
        amountError.value = '请输入有效的金额'
        return
      }

      if (amountNum < 10) {
        amountError.value = '充值金额不能低于 10 USDT'
        return
      }

      if (amountNum > 10000) {
        amountError.value = '充值金额不能超过 10,000 USDT'
        return
      }

      // 检查小数位数
      const decimalPlaces = (amount.value.split('.')[1] || '').length
      if (decimalPlaces > 2) {
        amountError.value = '金额最多支持2位小数'
        return
      }
    }

    const selectQuickAmount = (quickAmount) => {
      amount.value = quickAmount.toString()
      validateAmount()
    }
    
    const testRouteJump = async () => {
      console.log('测试路由跳转...')
      const testOrderNo = 'RCH17623472201700'
      
      try {
        await router.push({
          name: 'USDTRechargeDetail',
          params: { orderNo: testOrderNo }
        })
        console.log('测试跳转成功')
      } catch (error) {
        console.error('测试跳转失败:', error)
        appStore.showError('测试跳转失败: ' + error.message)
      }
    }

    const createRechargeOrder = async () => {
      if (!canRecharge.value || loading.value) return

      loading.value = true

      try {
        // 调用API创建充值订单
        const response = await api.wallet.createRechargeOrder({
          amount: amount.value
        })

        console.log('充值订单创建响应:', response)
        console.log('响应类型:', typeof response)
        console.log('order_no 字段:', response.order_no)
        console.log('order_no 类型:', typeof response.order_no)

        // 检查响应数据
        if (!response) {
          console.error('响应为空:', response)
          appStore.showError('创建充值订单失败：服务器无响应')
          return
        }

        if (!response.order_no) {
          console.error('响应数据格式错误，缺少 order_no 字段:', response)
          console.error('响应对象的所有键:', Object.keys(response))
          appStore.showError('创建充值订单失败：响应数据格式错误')
          return
        }

        // 显示成功提示
        appStore.showSuccess('充值订单创建成功')

        // 跳转到充值详情页面
        console.log('准备跳转到订单详情页面，订单号:', response.order_no)
        
        // 构建跳转路径
        const targetPath = `/wallet/recharge/detail/${response.order_no}`
        console.log('目标路径:', targetPath)
        
        try {
          // 尝试使用路径跳转
          await router.push(targetPath)
          console.log('路由跳转成功')
        } catch (routerError) {
          console.error('路由跳转失败:', routerError)
          
          // 尝试使用路由名称跳转
          try {
            console.log('尝试使用路由名称跳转...')
            await router.push({
              name: 'USDTRechargeDetail',
              params: { orderNo: response.order_no }
            })
            console.log('使用路由名称跳转成功')
          } catch (nameRouterError) {
            console.error('使用路由名称跳转也失败:', nameRouterError)
            appStore.showError('页面跳转失败，请手动查看充值订单')
          }
        }

      } catch (error) {
        console.error('创建充值订单失败:', error)

        // 根据错误类型显示不同的提示
        if (error.code === '40001') {
          amountError.value = error.message || '充值金额低于最小限额'
        } else if (error.code === '40002') {
          amountError.value = error.message || '充值金额格式错误'
        } else if (error.code === 'NETWORK_ERROR') {
          appStore.showError('网络连接失败，请检查网络设置')
        } else {
          appStore.showError(error.message || '创建充值订单失败，请稍后重试')
        }
      } finally {
        loading.value = false
      }
    }

    return {
      amount,
      amountError,
      loading,
      quickAmounts,
      canRecharge,
      validateAmount,
      selectQuickAmount,
      createRechargeOrder,
      testRouteJump
    }
  }
}
</script>

<style scoped>
.usdt-recharge-page {
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

.amount-section {
  margin-bottom: 24px;
}

.amount-input-container {
  position: relative;
  margin-bottom: 12px;
}

.amount-input {
  width: 100%;
  padding: 16px 60px 16px 16px;
  border: 2px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  font-size: 18px;
  font-weight: bold;
  background: var(--tg-theme-bg-color, #ffffff);
  color: var(--tg-theme-text-color, #000000);
  text-align: center;
  box-sizing: border-box;
}

.amount-input:focus {
  outline: none;
  border-color: var(--tg-theme-button-color, #0088cc);
}

.currency-label {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-hint-color, #666666);
}

.amount-tips {
  margin-bottom: 8px;
}

.tip-text {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 4px 0;
  line-height: 1.4;
}

.error-message {
  font-size: 12px;
  color: #ff4757;
  margin-top: 8px;
  padding: 8px 12px;
  background: rgba(255, 71, 87, 0.1);
  border-radius: 8px;
}

.quick-amounts {
  margin-bottom: 24px;
}

.quick-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 12px 0;
}

.quick-options {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.quick-option {
  padding: 12px 8px;
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 8px;
  background: var(--tg-theme-bg-color, #ffffff);
  color: var(--tg-theme-text-color, #000000);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-option.active {
  border-color: var(--tg-theme-button-color, #0088cc);
  background: rgba(0, 136, 204, 0.1);
  color: var(--tg-theme-button-color, #0088cc);
}

.quick-option:hover {
  border-color: var(--tg-theme-button-color, #0088cc);
}

.recharge-info {
  margin-bottom: 24px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.info-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.info-icon {
  font-size: 20px;
  margin-top: 2px;
}

.info-text {
  flex: 1;
}

.info-label {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.info-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  line-height: 1.4;
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
  margin-bottom: 6px;
  line-height: 1.4;
}

@media (max-width: 480px) {
  .quick-options {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>