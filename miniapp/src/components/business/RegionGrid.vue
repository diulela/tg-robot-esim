<template>
  <div class="region-grid" :class="gridClasses">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-grid">
      <div
        v-for="i in skeletonCount"
        :key="`skeleton-${i}`"
        class="region-skeleton"
      >
        <v-skeleton-loader
          type="avatar, text"
          class="skeleton-item"
        />
      </div>
    </div>

    <!-- 区域卡片 -->
    <div v-else-if="regions.length > 0" class="regions-container">
      <RegionCard
        v-for="region in regions"
        :key="region.id"
        :region="region"
        @click="handleRegionClick"
      />
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <v-icon size="48" color="grey-lighten-1">mdi-earth-off</v-icon>
      <h4 class="empty-title">暂无区域数据</h4>
      <p class="empty-subtitle">请稍后重试或联系客服</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Region } from '@/types'
import RegionCard from './RegionCard.vue'

// Props
interface Props {
  regions: Region[]
  loading?: boolean
  columns?: number
  gap?: string | number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  columns: 2,
  gap: 16
})

// Emits
const emit = defineEmits<{
  selectRegion: [region: Region]
}>()

// 计算属性
const gridClasses = computed(() => {
  const classes = []
  
  if (props.loading) {
    classes.push('loading')
  }
  
  classes.push(`columns-${props.columns}`)
  
  return classes
})

const skeletonCount = computed(() => {
  // 根据列数显示骨架屏数量
  return props.columns * 3 // 显示3行的骨架屏
})

const gridStyle = computed(() => {
  return {
    '--grid-columns': props.columns,
    '--grid-gap': typeof props.gap === 'number' ? `${props.gap}px` : props.gap
  }
})

// 方法
const handleRegionClick = (region: Region) => {
  emit('selectRegion', region)
}
</script>

<style scoped lang="scss">
.region-grid {
  --grid-columns: 2;
  --grid-gap: 16px;
  
  .loading-grid {
    display: grid;
    grid-template-columns: repeat(var(--grid-columns), 1fr);
    gap: var(--grid-gap);
    
    .region-skeleton {
      .skeleton-item {
        height: 120px;
        border-radius: 16px;
      }
    }
  }
  
  .regions-container {
    display: grid;
    grid-template-columns: repeat(var(--grid-columns), 1fr);
    gap: var(--grid-gap);
  }
  
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 48px 24px;
    text-align: center;
    
    .empty-title {
      margin: 16px 0 8px;
      color: rgba(var(--v-theme-on-surface), 0.8);
      font-size: 1.125rem;
      font-weight: 600;
    }
    
    .empty-subtitle {
      margin: 0;
      color: rgba(var(--v-theme-on-surface), 0.6);
      font-size: 0.875rem;
      line-height: 1.5;
    }
  }
  
  // 响应式列数调整
  &.columns-1 {
    --grid-columns: 1;
  }
  
  &.columns-2 {
    --grid-columns: 2;
  }
  
  &.columns-3 {
    --grid-columns: 3;
  }
  
  &.columns-4 {
    --grid-columns: 4;
  }
}

// 响应式适配
@media (max-width: 360px) {
  .region-grid {
    --grid-gap: 12px;
    
    .empty-state {
      padding: 32px 16px;
      
      .empty-title {
        font-size: 1rem;
      }
      
      .empty-subtitle {
        font-size: 0.8125rem;
      }
    }
  }
}

@media (min-width: 481px) {
  .region-grid {
    --grid-gap: 20px;
    
    .empty-state {
      padding: 64px 32px;
      
      .empty-title {
        font-size: 1.25rem;
      }
      
      .empty-subtitle {
        font-size: 0.9375rem;
      }
    }
  }
}

// 动画效果
.regions-container {
  .region-card {
    animation: fadeInUp 0.3s ease;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 为不同列数优化间距
@media (max-width: 480px) {
  .region-grid {
    &.columns-3,
    &.columns-4 {
      --grid-columns: 2; // 小屏幕强制使用2列
    }
  }
}

@media (min-width: 768px) {
  .region-grid {
    &.columns-2 {
      --grid-columns: 3; // 大屏幕可以显示更多列
    }
    
    &.columns-3 {
      --grid-columns: 4;
    }
  }
}
</style>