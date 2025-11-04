<template>
  <div class="page-header">
    <v-btn
      v-if="showBack"
      icon
      variant="text"
      class="back-button"
      @click="handleBack"
    >
      <v-icon>mdi-chevron-left</v-icon>
    </v-btn>
    <h1 class="page-title">{{ title }}</h1>
    <div class="header-spacer"></div>
  </div>
</template>

<script setup lang="ts">
import type { PageHeaderProps, PageHeaderEmits } from '@/types'

// Props
const props = withDefaults(defineProps<PageHeaderProps>(), {
  showBack: true
})

// Emits
const emit = defineEmits<PageHeaderEmits>()

// 方法
const handleBack = () => {
  emit('back')
}
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  min-height: 56px;
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.back-button {
  color: white !important;
  flex-shrink: 0;
}

.page-title {
  flex: 1;
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  padding: 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-spacer {
  width: 48px;
  flex-shrink: 0;
}

/* 响应式设计 */
@media (min-width: 960px) {
  .page-header {
    min-height: 64px;
  }

  .page-title {
    font-size: 20px;
  }
}

/* 动画效果 */
.back-button {
  transition: transform 0.2s ease;
}

.back-button:hover {
  transform: scale(1.1);
}

.back-button:active {
  transform: scale(0.95);
}
</style>
