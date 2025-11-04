<template>
  <div class="product-page">
    <!-- 栏目导航 -->
    <div class="category-tabs">
      <div 
        v-for="tab in categoryTabs" 
        :key="tab.key" 
        @click="switchCategory(tab.key)"
        :class="['tab-item', { active: activeCategory === tab.key }]"
      >
        {{ tab.label }}
      </div>
    </div>

    <!-- 动态组件 -->
    <keep-alive>
      <component 
        :is="currentComponent" 
        class="category-content"
      />
    </keep-alive>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProductsStore } from '@/stores/products'
import HotProductsPage from './HotProductsPage.vue'
import LocalProductsPage from './LocalProductsPage.vue'
import RegionProductsPage from './RegionProductsPage.vue'
import GlobalProductsPage from './GlobalProductsPage.vue'

// Store
const productsStore = useProductsStore()

// 状态
const activeCategory = ref<'hot' | 'local' | 'region' | 'global'>('hot')

// 栏目配置
const categoryTabs: Array<{ key: 'hot' | 'local' | 'region' | 'global'; label: string }> = [
  { key: 'hot', label: '热门' },
  { key: 'local', label: '本地' },
  { key: 'region', label: '区域' },
  { key: 'global', label: '全球' }
]

// 组件映射
const componentMap = {
  hot: HotProductsPage,
  local: LocalProductsPage,
  region: RegionProductsPage,
  global: GlobalProductsPage
}

// 计算属性 - 当前组件
const currentComponent = computed(() => {
  return componentMap[activeCategory.value]
})

// 方法 - 切换栏目
const switchCategory = (category: 'hot' | 'local' | 'region' | 'global') => {
  activeCategory.value = category
  productsStore.setCurrentCategory(category)
}
</script>

<style scoped>
.product-page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #f5f5f5;
}

.category-tabs {
  background: white;
  display: flex;
  padding: 0;
  margin: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.tab-item {
  flex: 1;
  padding: 16px 12px;
  text-align: center;
  font-size: 16px;
  color: #666666;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
  border-bottom: 3px solid transparent;
}

.tab-item.active {
  color: #667eea;
  border-bottom-color: #667eea;
  font-weight: 600;
}

.tab-item:hover {
  background: #f5f5f5;
}

.category-content {
  flex: 1;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .tab-item {
    padding: 14px 8px;
    font-size: 14px;
  }
}

@media (max-width: 360px) {
  .tab-item {
    padding: 12px 6px;
    font-size: 13px;
  }
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>