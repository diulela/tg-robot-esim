<template>
  <v-card
    :color="cardColor"
    :variant="variant"
    class="action-card"
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
      
      <!-- 内容 -->
      <div class="action-content">
        <h4 class="action-title">{{ title }}</h4>
        <p v-if="description" class="action-description">{{ description }}</p>
      </div>
      
      <!-- 箭头图标 -->
      <div v-if="showArrow" class="arrow-container">
        <v-icon
          icon="mdi-chevron-right"
          :color="arrowColor"
          size="20"
        />
      </div>
      
      <!-- 徽章 -->
      <v-badge
        v-if="badge"
        :content="badge"
        :color="badgeColor"
        class="action-badge"
        floating
        offset-x="8"
        offset-y="8"
      />
      
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
  description?: string
  icon: string
  color?: string
  variant?: 'elevated' | 'flat' | 'tonal' | 'outlined' | 'text' | 'plain'
  loading?: boolean
  disabled?: boolean
  showArrow?: boolean
  badge?: string | number
  badgeColor?: string
  size?: 'small' | 'default' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  variant: 'elevated',
  loading: false,
  disabled: false,
  showArrow: true,
  badgeColor: 'error',
  size: 'default'
})

// Emits
const emit = defineEmits<{
  click: []
}>()

// 计算属性
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

const arrowColor = computed(() => {
  return 'on-surface-variant'
})

const iconSize = computed(() => {
  const sizeMap = {
    small: 20,
    default: 24,
    large: 28
  }
  return sizeMap[props.size]
})

const cardClasses = computed(() => {
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
  
  classes.push(`size-${props.size}`)
  
  return classes
})

// 方法
const handleClick = () => {
  if (!props.disabled && !props.loading) {
    telegramService.selectionFeedback()
    emit('click')
  }
}
</script>

<style scoped lang="scss">
.action-card {
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
  
  &.disabled {
    opacity: 0.6;
    pointer-events: none;
  }
  
  .card-content {
    padding: 16px !important;
    display: flex;
    align-items: center;
    gap: 12px;
    position: relative;
    min-height: 72px;
    
    .icon-container {
      flex-shrink: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      border-radius: 12px;
      background: rgba(var(--v-theme-primary), 0.1);
    }
    
    .action-content {
      flex: 1;
      min-width: 0; // 防止文字溢出
      
      .action-title {
        font-size: 0.875rem;
        font-weight: 600;
        color: rgb(var(--v-theme-on-surface));
        margin: 0 0 4px 0;
        line-height: 1.2;
      }
      
      .action-description {
        font-size: 0.75rem;
        color: rgba(var(--v-theme-on-surface), 0.7);
        margin: 0;
        line-height: 1.3;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
    
    .arrow-container {
      flex-shrink: 0;
      opacity: 0.6;
    }
    
    .action-badge {
      position: absolute;
      top: 8px;
      right: 8px;
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
      padding: 12px !important;
      min-height: 56px;
      gap: 8px;
      
      .icon-container {
        width: 32px;
        height: 32px;
        border-radius: 8px;
      }
      
      .action-content {
        .action-title {
          font-size: 0.8125rem;
        }
        
        .action-description {
          font-size: 0.7rem;
        }
      }
    }
  }
  
  &.size-large {
    .card-content {
      padding: 20px !important;
      min-height: 88px;
      gap: 16px;
      
      .icon-container {
        width: 48px;
        height: 48px;
        border-radius: 16px;
      }
      
      .action-content {
        .action-title {
          font-size: 1rem;
        }
        
        .action-description {
          font-size: 0.8125rem;
        }
      }
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .action-card {
    .card-content {
      padding: 12px !important;
      min-height: 64px;
      gap: 10px;
      
      .icon-container {
        width: 36px;
        height: 36px;
        border-radius: 10px;
      }
      
      .action-content {
        .action-title {
          font-size: 0.8125rem;
        }
        
        .action-description {
          font-size: 0.7rem;
        }
      }
    }
  }
}

@media (min-width: 481px) {
  .action-card {
    .card-content {
      padding: 18px !important;
      min-height: 80px;
      gap: 14px;
      
      .icon-container {
        width: 44px;
        height: 44px;
        border-radius: 14px;
      }
      
      .action-content {
        .action-title {
          font-size: 0.9375rem;
        }
        
        .action-description {
          font-size: 0.8125rem;
        }
      }
    }
  }
}

// 主题适配
.v-theme--dark {
  .action-card {
    .card-content {
      .icon-container {
        background: rgba(var(--v-theme-primary), 0.2);
      }
    }
    
    .loading-overlay {
      background: rgba(var(--v-theme-surface), 0.9);
    }
  }
}
</style>