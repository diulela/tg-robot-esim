<template>
  <div class="order-note-input">
    <v-textarea
      v-model="internalValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :maxlength="maxLength"
      variant="outlined"
      rows="3"
      auto-grow
      no-resize
      hide-details="auto"
      class="note-textarea"
      @input="handleInput"
      @focus="handleFocus"
      @blur="handleBlur"
    />
    
    <div class="character-count">
      <span class="count-text" :class="{ 'count-warning': isNearLimit }">
        {{ characterCount }}/{{ maxLength }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { telegramService } from '@/services/telegram'

// Props 定义
interface Props {
  modelValue: string
  placeholder?: string
  maxLength?: number
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入订单备注（选填）',
  maxLength: 200,
  disabled: false
})

// Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: string]
  'focus': [event: FocusEvent]
  'blur': [event: FocusEvent]
}>()

// 响应式状态
const isFocused = ref(false)

// 计算属性
const internalValue = computed({
  get: () => props.modelValue,
  set: (value: string) => {
    // 限制字符长度
    const trimmedValue = value.slice(0, props.maxLength)
    emit('update:modelValue', trimmedValue)
  }
})

const characterCount = computed(() => props.modelValue.length)

const isNearLimit = computed(() => {
  const ratio = characterCount.value / props.maxLength
  return ratio >= 0.8 // 当字符数达到80%时显示警告样式
})

// 方法
const handleInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement
  const value = target.value
  
  // 如果超出限制，截断并提供触觉反馈
  if (value.length > props.maxLength) {
    telegramService.impactFeedback('medium')
  }
}

const handleFocus = (event: FocusEvent) => {
  isFocused.value = true
  emit('focus', event)
  telegramService.impactFeedback('light')
}

const handleBlur = (event: FocusEvent) => {
  isFocused.value = false
  emit('blur', event)
}

// 监听字符数变化，提供反馈
watch(characterCount, (newCount, oldCount) => {
  // 当接近限制时提供触觉反馈
  if (newCount >= props.maxLength && oldCount < props.maxLength) {
    telegramService.impactFeedback('medium')
  }
})
</script>

<style scoped lang="scss">
.order-note-input {
  position: relative;

  .note-textarea {
    :deep(.v-field) {
      border-radius: 12px;
      transition: all 0.2s ease;

      &:hover {
        box-shadow: 0 2px 8px rgba(var(--v-theme-primary), 0.1);
      }

      &.v-field--focused {
        box-shadow: 0 4px 12px rgba(var(--v-theme-primary), 0.2);
      }
    }

    :deep(.v-field__input) {
      font-size: 0.875rem;
      line-height: 1.4;
      padding: 12px 16px;
    }

    :deep(.v-field__outline) {
      --v-field-border-opacity: 0.3;
    }

    :deep(.v-field--focused .v-field__outline) {
      --v-field-border-width: 2px;
    }
  }

  .character-count {
    display: flex;
    justify-content: flex-end;
    margin-top: 8px;
    padding-right: 4px;

    .count-text {
      font-size: 0.75rem;
      color: rgba(var(--v-theme-on-surface), 0.6);
      transition: color 0.2s ease;

      &.count-warning {
        color: rgb(var(--v-theme-warning));
        font-weight: 500;
      }
    }
  }
}

// 禁用状态样式
.order-note-input:has(.note-textarea:disabled) {
  .note-textarea {
    :deep(.v-field) {
      opacity: 0.6;
    }
  }

  .character-count {
    opacity: 0.6;
  }
}

// 响应式适配
@media (max-width: 360px) {
  .order-note-input {
    .note-textarea {
      :deep(.v-field__input) {
        font-size: 0.8rem;
        padding: 10px 14px;
      }
    }

    .character-count {
      margin-top: 6px;

      .count-text {
        font-size: 0.7rem;
      }
    }
  }
}
</style>