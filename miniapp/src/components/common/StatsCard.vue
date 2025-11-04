<template>
  <v-card
    :color="cardColor"
    :variant="variant"
    class="stats-card"
    :class="cardClasses"
    @click="handleClick"
  >
    <v-card-text class="card-content">
      <!-- 图标 -->
      <div class="icon-container">
        <v-icon
          :icon="icon"
          :color="iconColor"
          :size="iconSize"
        />
      </div>
      
      <!-- 数值和标题 -->
      <div class="stats-content">
        <div class="stats-value">
          {{ formattedValue }}
        </div>
        
        <div class="stats-title">
          {{ title }}
        </div>
        
        <!-- 变化趋势 -->
        <div v-if="change !== undefined" class="stats-change">
          <v-icon
            :icon="changeIcon"
            :color="changeColor"
            size="14"
          />
          <span :class="`text-${changeColor}`">
            {{ Math.abs(change) }}%
          </span>
        </div>
      </div>
      
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-overlay">
        <v-progress-circular
          indeterminate
          size="24"
          width="2"
          color="primary"
        />
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { telegramService } from '@/services/telegram'

// Props
interface Props {
  title: string
  value: string | number
  change?: number
  changeType?: 'increase' | 'decrease' | 'neutral'
  icon?: string
  color?: string
  variant?: 'elevated' | 'flat' | 'tonal' | 'outlined' | 'text' | 'plain'
  loading?: boolean
  clickable?: boolean
  size?: 'small' | 'default' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'mdi-chart-line',
  color: 'primary',
  variant: 'elevated',
  loading: false,
  clickable: true,
  size: 'default'
})

// Emits
const emit = defineEmits<{
  click: []
}>()

// 计算属性
const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    // 格式化数字显示
    if (props.value >= 1000000) {
      return `${(props.value / 1000000).toFixed(1)}M`
    } else if (props.value >= 1000) {
      return `${(props.value / 1000).toFixed(1)}K`
    } else {
      return props.value.toString()
    }
  }
  return props.value
})

const cardColor = computed(() => {
  if (props.variant === 'elevated' || props.variant === 'flat') {
    return 'surface'
  }
  return props.color
})

const iconColor = computed(() => {
  if (props.variant === 'tonal') {
    return props.color
  }
  return props.color
})

const iconSize = computed(() => {
  const sizeMap = {
    small: 20,
    default: 24,
    large: 28
  }
  return sizeMap[props.size]
})

const changeIcon = computed(() => {
  if (props.changeType === 'increase') {
    return 'mdi-trending-up'
  } else if (props.changeType === 'decrease') {
    return 'mdi-trending-down'
  } else {
    return 'mdi-trending-neutral'
  }
})

const changeColor = computed(() => {
  if (props.changeType === 'increase') {
    return 'success'
  } else if (props.changeType === 'decrease') {
    return 'error'
  } else {
    return 'warning'
  }
})

const cardClasses = computed(() => {
  const classes = []
  
  if (props.clickable) {
    classes.push('clickable')
  }
  
  if (props.loading) {
    classes.push('loading')
  }
  
  classes.push(`size-${props.size}`)
  
  return classes
})

// 方法
const handleClick = () => {
  if (props.clickable && !props.loading) {
    telegramService.selectionFeedback()
    emit('click')
  }
}
</script>

<style scoped lang="scss">
.stats-card {
  position: relative;
  transition: all 0.2s ease;
  
  &.clickable {
    cursor: pointer;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
  
  &.loading {
    pointer-events: none;
  }
  
  .card-content {
    padding: 12px !important;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    position: relative;
    min-height: 80px;
    
    .icon-container {
      margin-bottom: 8px;
      opacity: 0.8;
    }
    
    .stats-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      
      .stats-value {
        font-size: 1.25rem;
        font-weight: 700;
        color: rgb(var(--v-theme-on-surface));
        line-height: 1.2;
        margin-bottom: 4px;
      }
      
      .stats-title {
        font-size: 0.75rem;
        font-weight: 500;
        color: rgba(var(--v-theme-on-surface), 0.7);
        line-height: 1.2;
        margin-bottom: 6px;
      }
      
      .stats-change {
        display: flex;
        align-items: center;
        gap: 2px;
        font-size: 0.7rem;
        font-weight: 600;
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
  }
  
  // 尺寸变体
  &.size-small {
    .card-content {
      padding: 8px !important;
      min-height: 60px;
      
      .stats-content {
        .stats-value {
          font-size: 1rem;
        }
        
        .stats-title {
          font-size: 0.7rem;
        }
        
        .stats-change {
          font-size: 0.65rem;
        }
      }
    }
  }
  
  &.size-large {
    .card-content {
      padding: 16px !important;
      min-height: 100px;
      
      .stats-content {
        .stats-value {
          font-size: 1.5rem;
        }
        
        .stats-title {
          font-size: 0.875rem;
        }
        
        .stats-change {
          font-size: 0.75rem;
        }
      }
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .stats-card {
    .card-content {
      padding: 8px !important;
      min-height: 70px;
      
      .stats-content {
        .stats-value {
          font-size: 1.125rem;
        }
        
        .stats-title {
          font-size: 0.7rem;
        }
        
        .stats-change {
          font-size: 0.65rem;
        }
      }
    }
  }
}

@media (min-width: 481px) {
  .stats-card {
    .card-content {
      padding: 16px !important;
      min-height: 90px;
      
      .stats-content {
        .stats-value {
          font-size: 1.375rem;
        }
        
        .stats-title {
          font-size: 0.8125rem;
        }
        
        .stats-change {
          font-size: 0.75rem;
        }
      }
    }
  }
}

// 主题适配
.v-theme--dark {
  .stats-card {
    .loading-overlay {
      background: rgba(var(--v-theme-surface), 0.9);
    }
  }
}
</style>