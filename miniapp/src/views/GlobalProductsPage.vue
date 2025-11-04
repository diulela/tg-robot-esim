<template>
  <div class="global-products-page">
    <!-- 商品列表 -->
    <ProductList
      :products="products"
      :loading="isLoading"
      :error="error"
      @buy-product="handleBuyProduct"
      @product-click="handleProductClick"
    />

    <!-- 加载更多按钮 -->
    <LoadMoreButton
      v-if="canLoadMore"
      :loading="isLoadingMore"
      @load-more="handleLoadMore"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProductsStore } from '@/stores/products'
import { useAppStore } from '@/stores/app'
import ProductList from '@/components/ProductList.vue'
import LoadMoreButton from '@/components/LoadMoreButton.vue'
import type { Product } from '@/types'

// 路由和 Store
const router = useRouter()
const productsStore = useProductsStore()
const appStore = useAppStore()

// 状态
const isLoadingMore = ref(false)

// 计算属性
const products = computed(() => productsStore.getGlobalProducts())

const isLoading = computed(() => productsStore.isLoading)

const error = computed(() => productsStore.error)

const canLoadMore = computed(() => {
  const pagination = productsStore.getGlobalProductsPagination()
  return pagination.hasNext && !isLoading.value
})

// 方法 - 获取全球商品
const fetchProducts = async () => {
  try {
    await productsStore.fetchGlobalProducts()
  } catch (err) {
    console.error('获取全球商品失败:', err)
    appStore.showNotification({
      type: 'error',
      message: '加载商品失败，请稍后重试'
    })
  }
}

// 方法 - 处理购买按钮点击
const handleBuyProduct = async (product: Product) => {
  try {
    // 设置当前产品
    productsStore.setCurrentProduct(product)
    
    // 跳转到商品详情页面
    await router.push({
      name: 'ProductDetail',
      params: { id: product.id }
    })
  } catch (err) {
    console.error('跳转失败:', err)
    appStore.showNotification({
      type: 'error',
      message: '操作失败，请稍后重试'
    })
  }
}

// 方法 - 处理商品点击
const handleProductClick = async (product: Product) => {
  try {
    // 设置当前产品
    productsStore.setCurrentProduct(product)
    
    // 跳转到商品详情页面
    await router.push({
      name: 'ProductDetail',
      params: { id: product.id }
    })
  } catch (err) {
    console.error('跳转失败:', err)
  }
}

// 方法 - 加载更多商品
const handleLoadMore = async () => {
  if (isLoadingMore.value) return

  isLoadingMore.value = true
  try {
    await productsStore.loadMoreGlobalProducts()
  } catch (err) {
    console.error('加载更多失败:', err)
    appStore.showNotification({
      type: 'error',
      message: '加载更多失败，请稍后重试'
    })
  } finally {
    isLoadingMore.value = false
  }
}

// 生命周期
onMounted(async () => {
  // 设置当前栏目
  productsStore.setCurrentCategory('global')
  await fetchProducts()
})
</script>

<style scoped>
.global-products-page {
  min-height: 100vh;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
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

.global-products-page {
  animation: fadeIn 0.3s ease-out;
}
</style>
