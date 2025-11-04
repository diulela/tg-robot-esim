<template>
  <v-card
    class="region-card"
    :class="cardClasses"
    :variant="variant"
    @click="handleClick"
  >
    <v-card-text class="card-content">
      <!-- 区域图标 -->
      <div class="region-icon">
        <v-avatar
          :size="iconSize"
          :color="iconColor"
          class="icon-avatar"
        >
          <v-img
            v-if="region.icon && region.icon.startsWith('http')"
            :src="region.icon"
            :alt="region.name"
          />
          <v-icon
            v-else
            :icon="regionIcon"
            :size="iconSize * 0.6"
            color="white"
          />
        </v-avatar>
        
        <!-- 热门标签 -->
        <v-chip
          v-if="region.isPopular"
          color="error"
          size="x-small"
          variant="flat"
          class="popular-badge"
        >
          热门
        </v-chip>
      </div>
      
      <!-- 区域信息 -->
      <div class="region-info">
        <h4 class="region-name">{{ region.name }}</h4>
        <p v-if="showCountryCount" class="country-count">
          {{ countryCountText }}
        </p>
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
import type { Region } from '@/types'

// Props
interface Props {
  region: Region
  variant?: 'elevated' | 'flat' | 'tonal' | 'outlined' | 'text' | 'plain'
  size?: 'small' | 'default' | 'large'
  loading?: boolean
  disabled?: boolean
  showCountryCount?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'elevated',
  size: 'default',
  loading: false,
  disabled: false,
  showCountryCount: true
})

// Emits
const emit = defineEmits<{
  click: [region: Region]
}>()

// 计算属性
const regionIcon = computed(() => {
  // 根据区域代码返回对应图标
  const iconMap: Record<string, string> = {
    'asia': 'mdi-earth',
    'europe': 'mdi-earth',
    'north-america': 'mdi-earth',
    'south-america': 'mdi-earth',
    'africa': 'mdi-earth',
    'oceania': 'mdi-earth',
    'middle-east': 'mdi-mosque',
    'global': 'mdi-web',
    'china': 'mdi-flag',
    'usa': 'mdi-flag',
    'eu': 'mdi-flag'
  }
  
  const code = props.region.code.toLowerCase()
  return iconMap[code] || 'mdi-earth'
})

const iconColor = computed(() => {
  // 根据区域类型返回不同颜色
  if (props.region.isPopular) {
    return 'primary'
  }
  
  const colorMap: Record<string, string> = {
    'asia': 'orange',
    'europe': 'blue',
    'north-america': 'green',
    'south-america': 'teal',
    'africa': 'brown',
    'oceania': 'cyan',
    'middle-east': 'purple',
    'global': 'indigo'
  }
  
  const code = props.region.code.toLowerCase()
  return colorMap[code] || 'primary'
})

const iconSize = computed(() => {
  const sizeMap = {
    small: 40,
    default: 48,
    large: 56
  }
  return sizeMap[props.size]
})

const countryCountText = computed(() => {
  const count = props.region.countries?.length || 0
  return count > 0 ? `${count} 个国家` : '敬请期待'
})

const cardClasses = computed(() => {
  const classes = ['region-card']
  
  if (!props.disabled && !props.loading) {
    classes.push('clickable')
  }
  
  if (props.loading) {
    classes.push('loading')
  }
  
  if (props.disabled) {
    classes.push('disabled')
  }
  
  if (props.region.isPopular) {
    classes.push('popular')
  }
  
  classes.push(`size-${props.size}`)
  
  return classes
})

// 方法
const handleClick = () => {
  if (!props.disabled && !props.loading) {
    telegramService.selectionFeedback()
    emit('click', props.region)
  }
}
</script>

<style scoped lang="scss">
.region-card {
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
  
  &.popular {
    border: 2px solid rgb(var(--v-theme-primary));
  }
  
  .card-content {
    padding: 16px !important;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    position: relative;
    min-height: 120px;
    
    .region-icon {
      position: relative;
      margin-bottom: 12px;
      
      .icon-avatar {
        border: 2px solid rgba(var(--v-theme-on-surface), 0.1);
      }
      
      .popular-badge {
        position: absolute;
        top: -4px;
        right: -8px;
        font-size: 0.6rem;
        height: 16px;
        min-width: 28px;
      }
    }
    
    .region-info {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: center;
      
      .region-name {
        font-size: 0.875rem;
        font-weight: 600;
        color: rgb(var(--v-theme-on-surface));
        margin: 0 0 4px 0;
        line-height: 1.3;
      }
      
      .country-count {
        font-size: 0.75rem;
        color: rgba(var(--v-theme-on-surface), 0.6);
        margin: 0;
        line-height: 1.2;
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
      padding: 12px !important;
      min-height: 100px;
      
      .region-icon {
        margin-bottom: 8px;
      }
      
      .region-info {
        .region-name {
          font-size: 0.8125rem;
        }
        
        .country-count {
          font-size: 0.7rem;
        }
      }
    }
  }
  
  &.size-large {
    .card-content {
      padding: 20px !important;
      min-height: 140px;
      
      .region-icon {
        margin-bottom: 16px;
      }
      
      .region-info {
        .region-name {
          font-size: 1rem;
        }
        
        .country-count {
          font-size: 0.8125rem;
        }
      }
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .region-card {
    .card-content {
      padding: 12px !important;
      min-height: 100px;
      
      .region-icon {
        margin-bottom: 8px;
        
        .popular-badge {
          font-size: 0.55rem;
          height: 14px;
          min-width: 24px;
        }
      }
      
      .region-info {
        .region-name {
          font-size: 0.8125rem;
        }
        
        .country-count {
          font-size: 0.7rem;
        }
      }
    }
  }
}

@media (min-width: 481px) {
  .region-card {
    .card-content {
      padding: 18px !important;
      min-height: 130px;
      
      .region-icon {
        margin-bottom: 14px;
        
        .popular-badge {
          font-size: 0.65rem;
          height: 18px;
          min-width: 32px;
        }
      }
      
      .region-info {
        .region-name {
          font-size: 0.9375rem;
        }
        
        .country-count {
          font-size: 0.8125rem;
        }
      }
    }
  }
}

// 主题适配
.v-theme--dark {
  .region-card {
    .card-content {
      .region-icon {
        .icon-avatar {
          border-color: rgba(var(--v-theme-on-surface), 0.2);
        }
      }
    }
    
    .loading-overlay {
      background: rgba(var(--v-theme-surface), 0.9);
    }
  }
}

// 动画效果
.region-card {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

// 热门区域特殊效果
.region-card.popular {
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(var(--v-theme-primary), 0.05) 0%, rgba(var(--v-theme-secondary), 0.05) 100%);
    border-radius: inherit;
    pointer-events: none;
  }
}
</style>