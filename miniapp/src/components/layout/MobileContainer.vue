<template>
  <v-container
    :fluid="fluid"
    :class="containerClasses"
    class="mobile-container"
  >
    <slot />
  </v-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// Props
interface Props {
  fluid?: boolean
  maxWidth?: string | number
  padding?: string | number
  class?: string
  safeArea?: boolean
  fullHeight?: boolean
  scrollable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  fluid: true,
  maxWidth: '480px',
  padding: '16px',
  safeArea: true,
  fullHeight: false,
  scrollable: false
})

// 计算属性
const containerClasses = computed(() => {
  const classes = []
  
  if (props.class) {
    classes.push(props.class)
  }
  
  if (props.safeArea) {
    classes.push('safe-area')
  }
  
  if (props.fullHeight) {
    classes.push('full-height')
  }
  
  if (props.scrollable) {
    classes.push('scrollable')
  }
  
  return classes
})

const containerStyle = computed(() => {
  const style: Record<string, string> = {}
  
  if (props.maxWidth) {
    style.maxWidth = typeof props.maxWidth === 'number' 
      ? `${props.maxWidth}px` 
      : props.maxWidth
  }
  
  if (props.padding) {
    style.padding = typeof props.padding === 'number' 
      ? `${props.padding}px` 
      : props.padding
  }
  
  return style
})
</script>

<style scoped lang="scss">
.mobile-container {
  max-width: v-bind('containerStyle.maxWidth');
  margin: 0 auto;
  padding: v-bind('containerStyle.padding');
  
  // 安全区域适配
  &.safe-area {
    padding-top: max(env(safe-area-inset-top), v-bind('containerStyle.padding'));
    padding-bottom: max(env(safe-area-inset-bottom), v-bind('containerStyle.padding'));
    padding-left: max(env(safe-area-inset-left), v-bind('containerStyle.padding'));
    padding-right: max(env(safe-area-inset-right), v-bind('containerStyle.padding'));
  }
  
  // 全高度
  &.full-height {
    min-height: 100vh;
    min-height: 100dvh; // 动态视口高度
  }
  
  // 可滚动
  &.scrollable {
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    
    // 自定义滚动条
    &::-webkit-scrollbar {
      width: 4px;
    }
    
    &::-webkit-scrollbar-track {
      background: rgba(0, 0, 0, 0.1);
      border-radius: 2px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: rgba(0, 0, 0, 0.3);
      border-radius: 2px;
      
      &:hover {
        background: rgba(0, 0, 0, 0.5);
      }
    }
  }
  
  // 响应式适配
  @media (max-width: 360px) {
    padding: 12px;
    
    &.safe-area {
      padding-top: max(env(safe-area-inset-top), 12px);
      padding-bottom: max(env(safe-area-inset-bottom), 12px);
      padding-left: max(env(safe-area-inset-left), 12px);
      padding-right: max(env(safe-area-inset-right), 12px);
    }
  }
  
  @media (min-width: 481px) {
    padding: 20px;
    
    &.safe-area {
      padding-top: max(env(safe-area-inset-top), 20px);
      padding-bottom: max(env(safe-area-inset-bottom), 20px);
      padding-left: max(env(safe-area-inset-left), 20px);
      padding-right: max(env(safe-area-inset-right), 20px);
    }
  }
}
</style>