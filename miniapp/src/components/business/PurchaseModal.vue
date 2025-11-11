<template>
  <v-dialog v-model="internalVisible" max-width="400" persistent class="purchase-modal">
    <v-card class="purchase-card">
      <!-- 弹窗标题 -->
      <v-card-title class="modal-title">
        <span>选择购买数量</span>
        <v-btn icon variant="text" size="small" @click="handleClose" class="close-btn">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-card-text class="modal-content">
        <!-- 商品信息 -->
        <div v-if="product" class="product-info">
          <h3 class="product-name">{{ product.name }}</h3>
          <p class="product-price">单价：${{ product.price.toFixed(2) }}</p>
        </div>

        <!-- 数量选择器 -->
        <div class="quantity-section">
          <QuantitySelector v-model="quantity" :disabled="isLoading" @update:model-value="handleQuantityChange" />
        </div>

        <!-- 邮箱输入 -->
        <div class="email-section">
          <EmailInput 
            v-model="customerEmail" 
            :disabled="isLoading" 
            :required="true"
            hint="eSIM 激活信息将发送到此邮箱"
            @validation="handleEmailValidation"
          />
        </div>

        <!-- 订单备注 -->
        <div class="note-section">
          <OrderNoteInput v-model="orderNote" :disabled="isLoading" placeholder="请输入订单备注（选填）" :max-length="200" />
        </div>

        <!-- 价格计算器 -->
        <div class="price-section">
          <PriceCalculator :unit-price="product?.price || 0" :quantity="quantity" currency="USD" />
        </div>

        <!-- 错误提示 -->
        <v-alert v-if="error" type="error" variant="tonal" class="error-alert" closable @click:close="clearError">
          {{ error }}
        </v-alert>
      </v-card-text>

      <!-- 操作按钮 -->
      <v-card-actions class="modal-actions">
        <v-btn variant="outlined" color="grey" @click="handleClose" :disabled="isLoading" class="cancel-btn">
          取消
        </v-btn>

        <v-btn variant="elevated" color="primary" @click="handlePurchase" :loading="isLoading"
          :disabled="!product || isLoading" class="purchase-btn">
          确认支付
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useOrdersStore } from '@/stores/orders'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'
import type { Product } from '@/types'
import type { EsimOrder } from '@/types/esim-order'
import QuantitySelector from './QuantitySelector.vue'
import EmailInput from './EmailInput.vue'
import OrderNoteInput from './OrderNoteInput.vue'
import PriceCalculator from './PriceCalculator.vue'

// Props 定义
interface Props {
  visible: boolean
  product: Product | null
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  product: null
})

// Emits 定义
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'purchase-success': [order: EsimOrder]
  'purchase-error': [error: string]
}>()

// 组合式 API
const router = useRouter()
const ordersStore = useOrdersStore()
const appStore = useAppStore()

// 响应式状态
const quantity = ref(1)
const customerEmail = ref('')
const isEmailValid = ref(false)
const orderNote = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)

// 计算属性
const internalVisible = computed({
  get: () => props.visible,
  set: (value: boolean) => emit('update:visible', value)
})

// 移除未使用的 totalPrice 计算属性，价格计算由 PriceCalculator 组件处理

// 方法
const handleClose = () => {
  if (isLoading.value) return

  // 重置状态
  resetModalState()

  // 关闭弹窗
  emit('update:visible', false)

  // 触觉反馈
  telegramService.impactFeedback('light')
}

const handleQuantityChange = (newQuantity: number) => {
  quantity.value = newQuantity
  clearError()
}

const handleEmailValidation = (isValid: boolean) => {
  isEmailValid.value = isValid
  clearError()
}

const handlePurchase = async () => {
  if (!props.product || isLoading.value) return

  // 验证邮箱
  if (!customerEmail.value) {
    error.value = '请输入邮箱地址'
    appStore.showNotification({
      type: 'error',
      message: '请输入邮箱地址',
      duration: 2000
    })
    return
  }

  if (!isEmailValid.value) {
    error.value = '邮箱格式不正确'
    appStore.showNotification({
      type: 'error',
      message: '邮箱格式不正确',
      duration: 2000
    })
    return
  }

  isLoading.value = true
  error.value = null

  try {
    // 触觉反馈
    telegramService.impactFeedback('medium')

    // 创建订单请求
    const orderRequest: any = {
      productId: Number(props.product.id),
      quantity: quantity.value,
      totalAmount: (props.product.price * quantity.value).toFixed(4),
      customerEmail: customerEmail.value
    }

    // 只在有备注时添加 remark 字段
    if (orderNote.value) {
      orderRequest.remark = orderNote.value
    }

    // 调用订单创建服务
    const newOrder = await ordersStore.createOrder(orderRequest)

    // 发送成功事件
    emit('purchase-success', newOrder)

    // 显示成功提示
    appStore.showNotification({
      type: 'success',
      message: '订单创建成功！',
      duration: 2000
    })

    // 关闭弹窗
    emit('update:visible', false)

    // 跳转到订单详情页面
    router.push({
      name: 'OrderDetail',
      params: { id: newOrder.id }
    })

    // 重置状态
    resetModalState()

  } catch (err) {
    const errorMessage = handlePurchaseError(err)
    error.value = errorMessage
    emit('purchase-error', errorMessage)

    // 显示错误提示
    appStore.showNotification({
      type: 'error',
      message: errorMessage,
      duration: 3000
    })
  } finally {
    isLoading.value = false
  }
}

const handlePurchaseError = (error: any): string => {
  let userMessage = '购买失败，请重试'

  if (error instanceof Error) {
    const message = error.message.toLowerCase()

    if (message.includes('余额不足') || message.includes('balance')) {
      userMessage = '余额不足，请先充值'
    } else if (message.includes('库存') || message.includes('stock')) {
      userMessage = '商品库存不足'
    } else if (message.includes('网络') || message.includes('network')) {
      userMessage = '网络连接失败，请重试'
    } else if (message.includes('登录') || message.includes('auth')) {
      userMessage = '请先登录'
    } else if (message.includes('邮箱') || message.includes('email')) {
      userMessage = '邮箱格式不正确，请检查'
    }
  }

  return userMessage
}

const clearError = () => {
  error.value = null
}

const resetModalState = () => {
  quantity.value = 1
  customerEmail.value = ''
  isEmailValid.value = false
  orderNote.value = ''
  error.value = null
  isLoading.value = false
}

// 监听弹窗显示状态
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    // 弹窗打开时重置状态
    resetModalState()
  }
})

// 监听产品变化
watch(() => props.product, (newProduct) => {
  if (newProduct) {
    // 产品变化时重置状态
    resetModalState()
  }
})
</script>

<style scoped lang="scss">
.purchase-modal {
  .purchase-card {
    border-radius: 16px;
    overflow: hidden;
  }

  .modal-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px 16px 24px;
    font-size: 1.25rem;
    font-weight: 600;
    background: linear-gradient(135deg, rgb(var(--v-theme-primary)) 0%, rgb(var(--v-theme-secondary)) 100%);
    color: white;

    .close-btn {
      color: white;
      opacity: 0.8;

      &:hover {
        opacity: 1;
      }
    }
  }

  .modal-content {
    padding: 24px;
  }

  .product-info {
    text-align: center;
    margin-bottom: 24px;

    .product-name {
      font-size: 1.125rem;
      font-weight: 600;
      margin: 0 0 8px 0;
      color: rgb(var(--v-theme-on-surface));
    }

    .product-price {
      font-size: 1rem;
      color: rgba(var(--v-theme-on-surface), 0.7);
      margin: 0;
    }
  }

  .quantity-section {
    margin-bottom: 24px;
  }

  .email-section {
    margin-bottom: 24px;
  }

  .note-section {
    margin-bottom: 24px;
  }

  .price-section {
    margin-bottom: 16px;
  }

  .error-alert {
    margin-bottom: 16px;
  }

  .modal-actions {
    padding: 16px 24px 24px 24px;
    gap: 12px;

    .cancel-btn {
      flex: 1;
      height: 48px;
    }

    .purchase-btn {
      flex: 2;
      height: 48px;
      font-weight: 600;
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .purchase-modal {
    .modal-title {
      padding: 16px 20px 12px 20px;
      font-size: 1.125rem;
    }

    .modal-content {
      padding: 20px;
    }

    .modal-actions {
      padding: 12px 20px 20px 20px;
      flex-direction: column;

      .cancel-btn,
      .purchase-btn {
        width: 100%;
        flex: none;
      }
    }
  }
}
</style>