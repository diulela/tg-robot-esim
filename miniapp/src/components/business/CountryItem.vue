<template>
  <v-list-item
    class="country-item"
    :class="itemClasses"
    @click="handleClick"
  >
    <!-- 国家代码 -->
    <template #prepend>
      <div class="country-code">
        <span class="code-text">{{ country.code }}</span>
      </div>
    </template>

    <!-- 国家信息 -->
    <v-list-item-title class="country-name">
      {{ country.name }}
    </v-list-item-title>

    <v-list-item-subtitle v-if="showFlag" class="country-subtitle">
      <span v-if="country.flag" class="country-flag">{{ country.flag }}</span>
      <span v-if="productCount > 0" class="product-count">
        {{ productCount }} 个套餐
      </span>
      <span v-else class="no-products">暂无套餐</span>
    </v-list-item-subtitle>

    <!-- 箭头图标 -->
    <template #append>
      <div class="item-actions">
        <!-- 热门标签 -->
        <v-chip
          v-if="country.isPopular"
          color="error"
          size="x-small"
          variant="flat"
          class="popular-chip"
        >
          热门
        </v-chip>
        
        <!-- 箭头 -->
        <v-icon
          v-if="showArrow"
          :color="arrowColor"
          size="20"
        >
          mdi-chevron-right
        </v-icon>
      </div>
    </template>

    <!-- 加载状态覆盖层 -->
    <div v-if="loading" class="loading-overlay">
      <v-progress-circular
        indeterminate
        size="20"
        width="2"
        color="primary"
      />
    </div>
  </v-list-item>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { telegramService } from '@/services/telegram'
import type { Country } from '@/types'

// Props
interface Props {
  country: Country
  showFlag?: boolean
  showArrow?: boolean
  loading?: boolean
  disabled?: boolean
  variant?: 'default' | 'compact' | 'detailed'
}

const props = withDefaults(defineProps<Props>(), {
  showFlag: true,
  showArrow: true,
  loading: false,
  disabled: false,
  variant: 'default'
})

// Emits
const emit = defineEmits<{
  click: [country: Country]
}>()

// 计算属性
const productCount = computed(() => {
  return props.country.products?.length || 0
})

const arrowColor = computed(() => {
  return props.disabled ? 'disabled' : 'on-surface-variant'
})

const itemClasses = computed(() => {
  const classes = []
  
  if (!props.disabled && !props.loading) {
    classes.push('clickable')
  }
  
  if (props.loading) {
    classes.push('loading')
  }
  
  if (props.disabled) {
    classes.push('disabled')
  }
  
  if (props.country.isPopular) {
    classes.push('popular')
  }
  
  classes.push(`variant-${props.variant}`)
  
  return classes
})

// 方法
const handleClick = () => {
  if (!props.disabled && !props.loading) {
    telegramService.selectionFeedback()
    emit('click', props.country)
  }
}
</script>

<style scoped lang="scss">
.country-item {
  position: relative;
  transition: all 0.2s ease;
  border-radius: 8px;
  margin-bottom: 1px;
  
  &.clickable {
    cursor: pointer;
    
    &:hover {
      background: rgba(var(--v-theme-on-surface), 0.04);
    }
    
    &:active {
      background: rgba(var(--v-theme-on-surface), 0.08);
    }
  }
  
  &.loading {
    pointer-events: none;
  }
  
  &.disabled {
    opacity: 0.6;
    pointer-events: none;
  }
  
  &.popular {
    background: rgba(var(--v-theme-primary), 0.02);
    border-left: 3px solid rgb(var(--v-theme-primary));
  }
  
  .country-code {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    background: rgba(var(--v-theme-primary), 0.1);
    border-radius: 8px;
    margin-right: 12px;
    
    .code-text {
      font-size: 0.75rem;
      font-weight: 700;
      color: rgb(var(--v-theme-primary));
      text-transform: uppercase;
    }
  }
  
  .country-name {
    font-size: 0.875rem;
    font-weight: 500;
    color: rgb(var(--v-theme-on-surface));
    line-height: 1.3;
  }
  
  .country-subtitle {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 2px;
    
    .country-flag {
      font-size: 1rem;
    }
    
    .product-count {
      font-size: 0.75rem;
      color: rgba(var(--v-theme-on-surface), 0.6);
    }
    
    .no-products {
      font-size: 0.75rem;
      color: rgba(var(--v-theme-on-surface), 0.4);
      font-style: italic;
    }
  }
  
  .item-actions {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .popular-chip {
      font-size: 0.6rem;
      height: 16px;
      min-width: 28px;
    }
  }
  
  .loading-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(var(--v-theme-surface), 0.8);
    border-radius: inherit;
  }
  
  // 变体样式
  &.variant-compact {
    min-height: 48px;
    
    .country-code {
      width: 32px;
      height: 32px;
      margin-right: 8px;
      
      .code-text {
        font-size: 0.7rem;
      }
    }
    
    .country-name {
      font-size: 0.8125rem;
    }
    
    .country-subtitle {
      .country-flag {
        font-size: 0.875rem;
      }
      
      .product-count,
      .no-products {
        font-size: 0.7rem;
      }
    }
  }
  
  &.variant-detailed {
    min-height: 72px;
    
    .country-code {
      width: 48px;
      height: 48px;
      margin-right: 16px;
      
      .code-text {
        font-size: 0.8125rem;
      }
    }
    
    .country-name {
      font-size: 0.9375rem;
    }
    
    .country-subtitle {
      margin-top: 4px;
      
      .country-flag {
        font-size: 1.125rem;
      }
      
      .product-count,
      .no-products {
        font-size: 0.8125rem;
      }
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .country-item {
    min-height: 56px;
    
    .country-code {
      width: 36px;
      height: 36px;
      margin-right: 10px;
      
      .code-text {
        font-size: 0.7rem;
      }
    }
    
    .country-name {
      font-size: 0.8125rem;
    }
    
    .country-subtitle {
      .country-flag {
        font-size: 0.875rem;
      }
      
      .product-count,
      .no-products {
        font-size: 0.7rem;
      }
    }
    
    .item-actions {
      gap: 6px;
      
      .popular-chip {
        font-size: 0.55rem;
        height: 14px;
        min-width: 24px;
      }
    }
  }
}

@media (min-width: 481px) {
  .country-item {
    min-height: 68px;
    
    .country-code {
      width: 44px;
      height: 44px;
      margin-right: 14px;
      
      .code-text {
        font-size: 0.8125rem;
      }
    }
    
    .country-name {
      font-size: 0.9375rem;
    }
    
    .country-subtitle {
      .country-flag {
        font-size: 1.125rem;
      }
      
      .product-count,
      .no-products {
        font-size: 0.8125rem;
      }
    }
    
    .item-actions {
      gap: 10px;
      
      .popular-chip {
        font-size: 0.65rem;
        height: 18px;
        min-width: 32px;
      }
    }
  }
}

// 主题适配
.v-theme--dark {
  .country-item {
    &.popular {
      background: rgba(var(--v-theme-primary), 0.05);
    }
    
    .country-code {
      background: rgba(var(--v-theme-primary), 0.15);
    }
    
    .loading-overlay {
      background: rgba(var(--v-theme-surface), 0.9);
    }
  }
}

// 动画效果
.country-item {
  animation: slideInLeft 0.3s ease;
}

@keyframes slideInLeft {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

// 热门国家特殊效果
.country-item.popular {
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(90deg, rgba(var(--v-theme-primary), 0.03) 0%, transparent 100%);
    border-radius: inherit;
    pointer-events: none;
  }
}

// 列表项间距优化
.country-item + .country-item {
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}

// 无产品状态样式
.country-item:has(.no-products) {
  opacity: 0.7;
  
  .country-code {
    background: rgba(var(--v-theme-on-surface), 0.1);
    
    .code-text {
      color: rgba(var(--v-theme-on-surface), 0.6);
    }
  }
}
</style>