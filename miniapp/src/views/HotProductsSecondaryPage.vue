<template>
  <div class="hot-products-secondary-page">
    <!-- 页面头部 -->
    <PageHeader
      :title="pageTitle"
      :show-back="true"
      @back="handleBack"
    />

    <!-- 商品列表 -->
    <ProductList
      :products="products"
      :loading="isLoading"
      :error="error"
      @buy-product="handleBuyProduct"
      @product-click="handleProductClick"
    />

    <!-- 加载更多 -->
    <LoadMoreButton
      v-if="canLoadMore"
      :loading="isLoadingMore"
      @load-more="handleLoadMore"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductsStore } from '@/stores/products'
import { useAppStore } from '@/stores/app'
import PageHeader from '@/components/PageHeader.vue'
import ProductList from '@/components/ProductList.vue'
import LoadMoreButton from '@/components/LoadMoreButton.vue'
import type { Product } from '@/types'

// 路由和Store
const route = useRoute()
const router = useRouter()
const productsStore = useProductsStore()
const appStore = useAppStore()

// 状态
const hotItemCode = ref<string>('')
const hotItemName = ref<string>('')
const isLoadingMore = ref(false)

// 计算属性
const pageTitle = computed(() => hotItemName.value || '热门商品')

const products = computed(() => {
  return productsStore.getHotItemProducts(hotItemCode.value)
})

const isLoading = computed(() => productsStore.isLoading)

const error = computed(() => productsStore.error)

const canLoadMore = computed(() => {
  const pagination = productsStore.getHotItemPagination(hotItemCode.value)
  return pagination.hasNext && !isLoading.value
})

// 方法
const fetchProducts = async () => {
  if (!hotItemCode.value) return

  try {
    await productsStore.fetchProductsByHotItem(hotItemCode.value)
  } catch (err) {
    console.error('获取商品失败:', err)
    appStore.showNotification({
      type: 'error',
      message: '加载商品失败，请稍后重试'
    })
  }
}

const handleBack = () => {
  router.back()
}

const handleBuyProduct = async (product: Product) => {
  try {
    // 设置当前产品
    productsStore.setCurrentProduct(product)
    
    // 跳转到购买页面
    await router.push({
      name: 'ProductDetail',
      params: { id: product.id }
    })
  } catch (err) {
    console.error('购买失败:', err)
    appStore.showNotification({
      type: 'error',
      message: '操作失败，请稍后重试'
    })
  }
}

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

const handleLoadMore = async () => {
  if (!hotItemCode.value || isLoadingMore.value) return

  isLoadingMore.value = true
  try {
    await productsStore.loadMoreHotItemProducts(hotItemCode.value)
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
  // 获取路由参数
  hotItemCode.value = route.params['hotItemCode'] as string
  hotItemName.value = (route.query['name'] as string) || hotItemCode.value

  // 页面初始化完成

  // 设置当前热门项
  productsStore.setCurrentHotItem({
    code: hotItemCode.value,
    name: hotItemName.value
  })

  // 加载商品数据
  await fetchProducts()
})

onUnmounted(() => {
  // 页面卸载清理
})
</script>

<style scoped>
.hot-products-secondary-page {
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

.hot-products-secondary-page {
  animation: fadeIn 0.3s ease-out;
}
</style>
