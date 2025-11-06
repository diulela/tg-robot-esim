<template>
  <div class="price-calculator">
    <div class="price-container" :class="{ 'price-updating': isUpdating }">
      <div class="price-content">
        <span class="price-label">总计：</span>
        <span class="price-value">{{ formattedPrice }}</span>
      </div>
      
      <!-- 价格变化动画指示器 -->
      <div v-if="isUpdating" class="update-indicator">
        <v-progress-circular
          indeterminate
          size="16"
          width="2"
          color="primary"
        />
      </div>
    </div>
    
    <!-- 价格明细（可选显示） -->
    <div v-if="showDetails" class="price-details">
      <div class="detail-item">
        <span class="detail-label">单价：</span>
        <span class="detail-value">{{ formatCurrency(unitPrice, currency) }}</span>
      </div>
      <div class="detail-item">
        <span class="detail-label">数量：</span>
        <span class="detail-value">{{ quantity }}</span>
      </div>
      <div class="detail-divider"></div>
      <div class="detail-item total-item">
        <span class="detail-label">小计：</span>
        <span class="detail-value">{{ formattedPrice }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

// Props 定义
interface Props {
  unitPrice: number
  quantity: number
  currency?: string
  showDetails?: boolean
  animationDuration?: number
}

const props = withDefaults(defineProps<Props>(), {
  unitPrice: 0,
  quantity: 1,
  currency: 'USD',
  showDetails: false,
  animationDuration: 300
})

// 响应式状态
const isUpdating = ref(false)
const displayPrice = ref(0)

// 计算属性
const totalPrice = computed(() => props.unitPrice * props.quantity)

const formattedPrice = computed(() => formatCurrency(displayPrice.value, props.currency))

// 方法
const formatCurrency = (amount: number, currency: string): string => {
  const currencySymbols: Record<string, string> = {
    'USD': '$',
    'CNY': '¥',
    'EUR': '€',
    'GBP': '£'
  }
  
  const symbol = currencySymbols[currency] || '$'
  return `${symbol}${amount.toFixed(2)}`
}

const animatePrice = async (from: number, to: number) => {
  if (from === to) return
  
  isUpdating.value = true
  
  const duration = props.animationDuration
  const steps = 20
  const stepDuration = duration / steps
  const stepValue = (to - from) / steps
  
  for (let i = 0; i <= steps; i++) {
    displayPrice.value = from + (stepValue * i)
    
    if (i < steps) {
      await new Promise(resolve => setTimeout(resolve, stepDuration))
    }
  }
  
  // 确保最终值准确
  displayPrice.value = to
  
  await nextTick()
  isUpdating.value = false
}

// 监听价格变化
watch(totalPrice, (newPrice, oldPrice) => {
  if (oldPrice !== undefined) {
    animatePrice(oldPrice, newPrice)
  } else {
    displayPrice.value = newPrice
  }
}, { immediate: true })

// 初始化显示价格
displayPrice.value = totalPrice.value
</script>

<style scoped lang="scss">
.price-calculator {
  .price-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    background: linear-gradient(135deg, rgba(var(--v-theme-primary), 0.1) 0%, rgba(var(--v-theme-secondary), 0.1) 100%);
    border: 2px solid rgb(var(--v-theme-primary));
    border-radius: 12px;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: -100%;
      width: 100%;
      height: 100%;
      background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
      transition: left 0.5s ease;
    }

    &.price-updating {
      &::before {
        left: 100%;
      }
    }

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(var(--v-theme-primary), 0.2);
    }

    .price-content {
      display: flex;
      align-items: center;
      gap: 8px;
      flex: 1;

      .price-label {
        font-size: 1rem;
        font-weight: 500;
        color: rgb(var(--v-theme-primary));
      }

      .price-value {
        font-size: 1.375rem;
        font-weight: 700;
        color: rgb(var(--v-theme-primary));
        letter-spacing: -0.02em;
        transition: all 0.2s ease;
      }
    }

    .update-indicator {
      margin-left: 8px;
    }
  }

  .price-details {
    margin-top: 12px;
    padding: 12px 16px;
    background: rgba(var(--v-theme-surface), 0.5);
    border-radius: 8px;
    border: 1px solid rgba(var(--v-theme-outline), 0.2);

    .detail-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;

      &:last-child {
        margin-bottom: 0;
      }

      &.total-item {
        font-weight: 600;
        color: rgb(var(--v-theme-primary));
      }

      .detail-label {
        font-size: 0.875rem;
        color: rgba(var(--v-theme-on-surface), 0.7);
      }

      .detail-value {
        font-size: 0.875rem;
        font-weight: 500;
        color: rgb(var(--v-theme-on-surface));
      }
    }

    .detail-divider {
      height: 1px;
      background: rgba(var(--v-theme-outline), 0.2);
      margin: 8px 0;
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .price-calculator {
    .price-container {
      padding: 14px 16px;

      .price-content {
        .price-label {
          font-size: 0.9rem;
        }

        .price-value {
          font-size: 1.25rem;
        }
      }
    }

    .price-details {
      padding: 10px 14px;

      .detail-item {
        .detail-label,
        .detail-value {
          font-size: 0.8rem;
        }
      }
    }
  }
}

// 暗色主题适配
@media (prefers-color-scheme: dark) {
  .price-calculator {
    .price-container {
      background: linear-gradient(135deg, rgba(var(--v-theme-primary), 0.15) 0%, rgba(var(--v-theme-secondary), 0.15) 100%);
      border-color: rgba(var(--v-theme-primary), 0.8);
    }

    .price-details {
      background: rgba(var(--v-theme-surface), 0.3);
      border-color: rgba(var(--v-theme-outline), 0.3);
    }
  }
}
</style>