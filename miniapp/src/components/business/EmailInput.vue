<template>
  <div class="email-input-wrapper">
    <v-text-field
      v-model="internalValue"
      label="邮箱地址"
      placeholder="请输入接收 eSIM 的邮箱"
      type="email"
      variant="outlined"
      density="comfortable"
      :disabled="disabled"
      :error-messages="errorMessage"
      :rules="rules"
      prepend-inner-icon="mdi-email-outline"
      clearable
      @update:model-value="handleInput"
      @blur="handleBlur"
      class="email-input"
    >
      <template #append-inner>
        <v-icon v-if="isValid && internalValue" color="success" size="small">
          mdi-check-circle
        </v-icon>
      </template>
    </v-text-field>
    <p v-if="hint" class="email-hint">
      {{ hint }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

// Props 定义
interface Props {
  modelValue: string
  disabled?: boolean
  required?: boolean
  hint?: string
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: true,
  hint: 'eSIM 激活信息将发送到此邮箱'
})

// Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: string]
  'validation': [isValid: boolean]
}>()

// 响应式状态
const internalValue = ref(props.modelValue)
const errorMessage = ref<string>('')
const touched = ref(false)

// 邮箱验证正则表达式
const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/

// 验证规则
const rules = computed(() => {
  const ruleList: Array<(v: string) => boolean | string> = []
  
  if (props.required) {
    ruleList.push((v: string) => !!v || '邮箱地址不能为空')
  }
  
  ruleList.push((v: string) => {
    if (!v && !props.required) return true
    return emailRegex.test(v) || '邮箱格式不正确'
  })
  
  return ruleList
})

// 计算属性
const isValid = computed(() => {
  if (!internalValue.value && !props.required) return true
  if (!internalValue.value && props.required) return false
  return emailRegex.test(internalValue.value)
})

// 方法
const handleInput = (value: string) => {
  internalValue.value = value
  emit('update:modelValue', value)
  
  // 实时验证
  if (touched.value) {
    validateEmail()
  }
  
  // 发送验证状态
  emit('validation', isValid.value)
}

const handleBlur = () => {
  touched.value = true
  validateEmail()
}

const validateEmail = () => {
  if (!internalValue.value && props.required) {
    errorMessage.value = '邮箱地址不能为空'
    return false
  }
  
  if (internalValue.value && !emailRegex.test(internalValue.value)) {
    errorMessage.value = '邮箱格式不正确'
    return false
  }
  
  errorMessage.value = ''
  return true
}

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  internalValue.value = newValue
})

// 暴露验证方法
defineExpose({
  validate: validateEmail,
  isValid
})
</script>

<style scoped lang="scss">
.email-input-wrapper {
  width: 100%;
  
  .email-input {
    :deep(.v-field) {
      border-radius: 12px;
    }
    
    :deep(.v-field__input) {
      font-size: 0.9375rem;
    }
  }
  
  .email-hint {
    margin-top: -12px;
    margin-bottom: 8px;
    padding-left: 16px;
    font-size: 0.75rem;
    color: rgba(var(--v-theme-on-surface), 0.6);
    line-height: 1.4;
  }
}

// 响应式适配
@media (max-width: 360px) {
  .email-input-wrapper {
    .email-input {
      :deep(.v-field__input) {
        font-size: 0.875rem;
      }
    }
    
    .email-hint {
      font-size: 0.6875rem;
    }
  }
}
</style>
