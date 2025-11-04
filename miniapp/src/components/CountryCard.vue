<template>
  <div class="country-card" @click="handleClick">
    <div class="country-code">{{ country.code }}</div>
    <div class="country-name">{{ country.name }}</div>
  </div>
</template>

<script setup lang="ts">
// Props
interface Country {
  code: string
  name: string
  en?: string
}

interface Props {
  country: Country
}

defineProps<Props>()

// Emits
const emit = defineEmits<{
  'click': []
}>()

// 处理点击
const handleClick = () => {
  emit('click')
}
</script>

<style scoped>
.country-card {
  background: white;
  padding: 32px 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  min-height: 120px;
  border: 1px solid transparent;
}

.country-card:hover {
  background: #f5f5f5;
  transform: scale(1.02);
}

.country-card:active {
  background: #e8e8e8;
  transform: scale(0.98);
}

.country-code {
  font-size: 32px;
  font-weight: bold;
  color: #333333;
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.country-name {
  font-size: 14px;
  color: #666666;
  text-align: center;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 360px) {
  .country-card {
    padding: 20px 8px;
    min-height: 90px;
  }

  .country-code {
    font-size: 24px;
  }

  .country-name {
    font-size: 11px;
  }
}

@media (min-width: 360px) and (max-width: 480px) {
  .country-card {
    padding: 24px 12px;
    min-height: 100px;
  }

  .country-code {
    font-size: 28px;
  }

  .country-name {
    font-size: 12px;
  }
}

/* 页面进入动画 */
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

.country-card {
  animation: fadeInUp 0.5s ease-out;
}
</style>
