<template>
  <div class="load-more-container">
    <v-btn
      :loading="loading"
      :disabled="disabled || loading"
      color="primary"
      variant="outlined"
      size="large"
      block
      @click="handleLoadMore"
    >
      <v-icon start>mdi-refresh</v-icon>
      {{ loading ? '加载中...' : '加载更多' }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
import type { LoadMoreButtonProps, LoadMoreButtonEmits } from '@/types'

// Props
withDefaults(defineProps<LoadMoreButtonProps>(), {
  loading: false,
  disabled: false
})

// Emits
const emit = defineEmits<LoadMoreButtonEmits>()

// 方法
const handleLoadMore = () => {
  emit('loadMore')
}
</script>

<style scoped>
.load-more-container {
  padding: 16px;
  display: flex;
  justify-content: center;
}

.v-btn {
  transition: all 0.3s ease;
}

.v-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.v-btn:active:not(:disabled) {
  transform: translateY(0);
}
</style>
