<template>
  <div class="country-list">
    <!-- 分组显示 -->
    <template v-if="showGroupHeaders && grouped">
      <div 
        v-for="(items, letter) in grouped" 
        :key="letter" 
        class="country-group"
      >
        <div 
          class="group-header" 
          :id="`letter-${letter}`"
        >
          {{ letter }}
        </div>
        <div class="countries-grid">
          <CountryCard
            v-for="country in items"
            :key="country.code"
            :country="country"
            @click="handleCountryClick(country)"
          />
        </div>
      </div>
    </template>
    
    <!-- 普通显示 -->
    <template v-else>
      <div class="countries-grid">
        <CountryCard
          v-for="country in countries"
          :key="country.code"
          :country="country"
          @click="handleCountryClick(country)"
        />
      </div>
    </template>

    <!-- 空状态 -->
    <div v-if="countries.length === 0" class="empty-state">
      <span class="empty-icon">🔍</span>
      <p class="empty-text">未找到匹配的国家或地区</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import CountryCard from './CountryCard.vue'

// 类型定义
interface Country {
  code: string
  name: string
  en?: string
}

// Props
interface Props {
  countries: Country[]
  grouped?: Record<string, Country[]>
  showGroupHeaders?: boolean
  columns?: number
}

withDefaults(defineProps<Props>(), {
  showGroupHeaders: false,
  columns: 3
})

// Emits
const emit = defineEmits<{
  'country-click': [country: Country]
}>()

// 处理国家点击
const handleCountryClick = (country: Country) => {
  emit('country-click', country)
}
</script>

<style scoped>
.country-list {
  flex: 1;
  overflow-y: auto;
  background: #f5f5f5;
}

/* 分组标题 */
.country-group {
  margin-bottom: 0;
}

.group-header {
  position: sticky;
  top: 0;
  background: #f5f5f5;
  padding: 12px 16px;
  font-size: 16px;
  font-weight: 600;
  color: #667eea;
  z-index: 10;
  border-bottom: 1px solid #e0e0e0;
}

/* 国家网格 */
.countries-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  background: #e0e0e0;
  margin: 0;
  padding: 0;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-text {
  font-size: 16px;
  color: #999999;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 360px) {
  .group-header {
    padding: 10px 12px;
    font-size: 14px;
  }

  .empty-icon {
    font-size: 40px;
  }

  .empty-text {
    font-size: 14px;
  }
}

@media (min-width: 480px) {
  .group-header {
    padding: 14px 20px;
    font-size: 18px;
  }
}

/* 滚动优化 */
.country-list {
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
}

/* 页面进入动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.country-list {
  animation: fadeIn 0.3s ease-out;
}
</style>
