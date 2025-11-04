<template>
  <div class="product-list">
    <!-- 加载状态 -->
    <div v-if="loading && products.length === 0" class="loading-container">
      <v-progress-circular
        indeterminate
        color="primary"
        size="64"
      />
      <p class="loading-text">正在加载商品...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-container">
      <v-icon size="64" color="error">mdi-alert-circle-outline</v-icon>
      <p class="error-text">{{ error }}</p>
      <v-btn
        color="primary"
        variant="elevated"
        @click="handleRetry"
      >
        重试
      </v-btn>
    </div>

    <!-- 空状态 -->
    <div v-else-if="products.length === 0" class="empty-container">
      <v-icon size="64" color="grey">mdi-package-variant</v-icon>
      <p class="empty-text">暂无商品</p>
      <p class="empty-hint">请稍后再试或选择其他地区</p>
    </div>

    <!-- 商品网格 -->
    <div v-else class="products-grid">
      <ProductCard
        v-for="product in products"
        :key="product.id"
        :product="product"
        @click="handleProductClick"
        @buy="handleBuyProduct"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue'
import type { ProductListProps, ProductListEmits } from '@/types'

// 懒加载 ProductCard 组件
const ProductCard = defineAsyncComponent(() => import('./ProductCard.vue'))

// Props
defineProps<ProductListProps>()

// Emits
const emit = defineEmits<ProductListEmits>()

// 方法
const handleProductClick = (product: any) => {
  emit('productClick', product)
}

const handleBuyProduct = (product: any) => {
  emit('buyProduct', product)
}

const handleRetry = () => {
  // 触发父组件重新加载
  window.location.reload()
}
</script>

<style scoped>
.product-list {
  padding: 16px;
  min-height: 400px;
}

.loading-container,
.error-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  gap: 16px;
}

.loading-text,
.error-text,
.empty-text {
  font-size: 16px;
  font-weight: 500;
  color: #666;
  margin: 0;
}

.empty-hint {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.products-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: 1fr;
}

/* 响应式设计 */
@media (min-width: 600px) {
  .products-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 960px) {
  .products-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* 动画效果 */
.products-grid > * {
  animation: fadeInUp 0.5s ease-out;
}

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

/* 为每个卡片添加延迟动画 */
.products-grid > *:nth-child(1) { animation-delay: 0.05s; }
.products-grid > *:nth-child(2) { animation-delay: 0.1s; }
.products-grid > *:nth-child(3) { animation-delay: 0.15s; }
.products-grid > *:nth-child(4) { animation-delay: 0.2s; }
.products-grid > *:nth-child(5) { animation-delay: 0.25s; }
.products-grid > *:nth-child(6) { animation-delay: 0.3s; }
</style>
