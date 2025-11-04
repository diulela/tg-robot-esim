<template>
  <div class="search-bar">
    <span class="search-icon">🔍</span>
    <input
      ref="inputRef"
      type="text"
      class="search-input"
      :value="modelValue"
      :placeholder="placeholder"
      @input="handleInput"
      @focus="handleFocus"
      @blur="handleBlur"
    />
    <button
      v-if="modelValue"
      class="clear-btn"
      @click="handleClear"
      aria-label="清除搜索"
    >
      ✕
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// Props
interface Props {
  modelValue: string
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '搜索'
})

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: string]
  'input': [value: string]
  'clear': []
  'focus': []
  'blur': []
}>()

// Refs
const inputRef = ref<HTMLInputElement>()

// 防抖定时器
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// 处理输入
const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value

  // 立即更新 v-model
  emit('update:modelValue', value)

  // 防抖触发 input 事件（300ms）
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    emit('input', value)
  }, 300)
}

// 处理清除
const handleClear = () => {
  emit('update:modelValue', '')
  emit('input', '')
  emit('clear')

  // 清除防抖定时器
  if (debounceTimer) {
    clearTimeout(debounceTimer)
    debounceTimer = null
  }

  // 聚焦输入框
  inputRef.value?.focus()
}

// 处理聚焦
const handleFocus = () => {
  emit('focus')
}

// 处理失焦
const handleBlur = () => {
  emit('blur')
}
</script>

<style scoped>
.search-bar {
  position: relative;
  display: flex;
  align-items: center;
  background: #ffffff;
  border-radius: 8px;
  padding: 12px 16px;
  margin: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.2s ease;
}

.search-bar:focus-within {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.search-icon {
  font-size: 18px;
  margin-right: 8px;
  color: #999999;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 16px;
  color: #333333;
  background: transparent;
  min-width: 0;
}

.search-input::placeholder {
  color: #999999;
}

.clear-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: none;
  background: #e0e0e0;
  border-radius: 50%;
  color: #666666;
  font-size: 14px;
  cursor: pointer;
  flex-shrink: 0;
  margin-left: 8px;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  background: #d0d0d0;
}

.clear-btn:active {
  transform: scale(0.9);
}

/* 响应式设计 */
@media (max-width: 360px) {
  .search-bar {
    padding: 10px 12px;
    margin: 12px;
  }

  .search-input {
    font-size: 14px;
  }
}
</style>
