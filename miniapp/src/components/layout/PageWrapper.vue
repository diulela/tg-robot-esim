<template>
  <div class="page-wrapper" :class="wrapperClasses">
    <!-- 页面头部 -->
    <div v-if="$slots.header" class="page-header">
      <slot name="header" />
    </div>
    
    <!-- 页面内容 -->
    <div class="page-content" :class="contentClasses">
      <MobileContainer
        :fluid="fluid"
        :max-width="maxWidth"
        :padding="padding"
        :safe-area="safeArea"
        :full-height="fullHeight"
        :scrollable="scrollable"
      >
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <v-progress-circular
            indeterminate
            color="primary"
            size="48"
          />
          <p class="loading-text">{{ loadingText }}</p>
        </div>
        
        <!-- 错误状态 -->
        <div v-else-if="error" class="error-container">
          <v-icon size="64" color="error">mdi-alert-circle</v-icon>
          <h3 class="error-title">出错了</h3>
          <p class="error-message">{{ error }}</p>
          <v-btn
            v-if="showRetry"
            color="primary"
            @click="$emit('retry')"
            class="retry-btn"
          >
            重试
          </v-btn>
        </div>
        
        <!-- 空状态 -->
        <div v-else-if="empty" class="empty-container">
          <v-icon size="64" color="grey">{{ emptyIcon }}</v-icon>
          <h3 class="empty-title">{{ emptyTitle }}</h3>
          <p class="empty-message">{{ emptyMessage }}</p>
          <v-btn
            v-if="emptyAction"
            color="primary"
            @click="$emit('emptyAction')"
            class="empty-action-btn"
          >
            {{ emptyActionText }}
          </v-btn>
        </div>
        
        <!-- 正常内容 -->
        <div v-else class="content-wrapper">
          <slot />
        </div>
      </MobileContainer>
    </div>
    
    <!-- 页面底部 -->
    <div v-if="$slots.footer" class="page-footer">
      <slot name="footer" />
    </div>
    
    <!-- 浮动操作按钮 -->
    <v-fab
      v-if="showFab && fabIcon"
      :icon="fabIcon"
      :color="fabColor"
      :size="fabSize"
      location="bottom end"
      class="page-fab"
      @click="$emit('fabClick')"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import MobileContainer from './MobileContainer.vue'

// Props
interface Props {
  // 容器属性
  fluid?: boolean
  maxWidth?: string | number
  padding?: string | number
  safeArea?: boolean
  fullHeight?: boolean
  scrollable?: boolean
  
  // 页面状态
  loading?: boolean
  loadingText?: string
  error?: string | null
  showRetry?: boolean
  empty?: boolean
  emptyIcon?: string
  emptyTitle?: string
  emptyMessage?: string
  emptyAction?: boolean
  emptyActionText?: string
  
  // 浮动按钮
  showFab?: boolean
  fabIcon?: string
  fabColor?: string
  fabSize?: string | number
  
  // 样式
  class?: string
  contentClass?: string
  background?: string
}

const props = withDefaults(defineProps<Props>(), {
  fluid: true,
  maxWidth: '480px',
  padding: '16px',
  safeArea: true,
  fullHeight: false,
  scrollable: false,
  
  loading: false,
  loadingText: '加载中...',
  error: null,
  showRetry: true,
  empty: false,
  emptyIcon: 'mdi-inbox',
  emptyTitle: '暂无数据',
  emptyMessage: '这里还没有任何内容',
  emptyAction: false,
  emptyActionText: '刷新',
  
  showFab: false,
  fabColor: 'primary',
  fabSize: 'default'
})

// Emits
defineEmits<{
  retry: []
  emptyAction: []
  fabClick: []
}>()

// 计算属性
const wrapperClasses = computed(() => {
  const classes = []
  
  if (props.class) {
    classes.push(props.class)
  }
  
  if (props.background) {
    classes.push(`bg-${props.background}`)
  }
  
  return classes
})

const contentClasses = computed(() => {
  const classes = []
  
  if (props.contentClass) {
    classes.push(props.contentClass)
  }
  
  if (props.loading) {
    classes.push('loading-state')
  }
  
  if (props.error) {
    classes.push('error-state')
  }
  
  if (props.empty) {
    classes.push('empty-state')
  }
  
  return classes
})
</script>

<style scoped lang="scss">
.page-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-height: 100dvh;
  
  .page-header {
    flex-shrink: 0;
    z-index: 10;
  }
  
  .page-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    
    &.loading-state,
    &.error-state,
    &.empty-state {
      justify-content: center;
      align-items: center;
    }
  }
  
  .page-footer {
    flex-shrink: 0;
    z-index: 10;
  }
}

// 加载状态
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  
  .loading-text {
    margin-top: 16px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    font-size: 0.875rem;
  }
}

// 错误状态
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  
  .error-title {
    margin: 16px 0 8px;
    color: rgb(var(--v-theme-error));
    font-size: 1.25rem;
    font-weight: 600;
  }
  
  .error-message {
    margin-bottom: 24px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    font-size: 0.875rem;
    line-height: 1.5;
  }
  
  .retry-btn {
    min-width: 120px;
  }
}

// 空状态
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  
  .empty-title {
    margin: 16px 0 8px;
    color: rgba(var(--v-theme-on-surface), 0.8);
    font-size: 1.25rem;
    font-weight: 600;
  }
  
  .empty-message {
    margin-bottom: 24px;
    color: rgba(var(--v-theme-on-surface), 0.6);
    font-size: 0.875rem;
    line-height: 1.5;
  }
  
  .empty-action-btn {
    min-width: 120px;
  }
}

// 内容包装器
.content-wrapper {
  flex: 1;
  width: 100%;
}

// 浮动按钮
.page-fab {
  position: fixed !important;
  bottom: 80px; // 避免与底部导航重叠
  right: 16px;
  z-index: 1000;
  
  @media (max-width: 360px) {
    bottom: 76px;
    right: 12px;
  }
}

// 背景样式
.bg-gradient {
  background: linear-gradient(135deg, rgb(var(--v-theme-primary)) 0%, rgb(var(--v-theme-secondary)) 100%);
}

.bg-surface {
  background: rgb(var(--v-theme-surface));
}

.bg-background {
  background: rgb(var(--v-theme-background));
}

// 响应式适配
@media (max-width: 360px) {
  .loading-container,
  .error-container,
  .empty-container {
    padding: 32px 16px;
  }
}

@media (min-width: 481px) {
  .loading-container,
  .error-container,
  .empty-container {
    padding: 48px 24px;
  }
}
</style>