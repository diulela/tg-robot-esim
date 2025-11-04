<template>
  <div class="product-list-secondary-page">
    <!-- 页面头部 -->
    <!-- <PageHeader
      :title="pageTitle"
      :show-back="true"
      @back="handleBack"
    /> -->

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
const countryCode = ref<string>('')
const countryName = ref<string>('')
const category = ref<string>('')
const isLoadingMore = ref(false)

// 计算属性
const pageTitle = computed(() => countryName.value || '商品列表')

const products = computed(() => {
  return productsStore.getHotItemProducts(countryCode.value)
})

const isLoading = computed(() => productsStore.isLoading)

const error = computed(() => productsStore.error)

const canLoadMore = computed(() => {
  const pagination = productsStore.getHotItemPagination(countryCode.value)
  return pagination.hasNext && !isLoading.value
})

// 方法
const fetchProducts = async () => {
  if (!countryCode.value) return

  try {
    await productsStore.fetchProductsByHotItem(countryCode.value)
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
  if (!countryCode.value || isLoadingMore.value) return

  isLoadingMore.value = true
  try {
    await productsStore.loadMoreHotItemProducts(countryCode.value)
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
  countryCode.value = route.params['countryCode'] as string
  countryName.value = (route.query['name'] as string) || countryCode.value
  category.value = (route.query['category'] as string) || 'hot'

  // 设置页面状态
  appStore.setCurrentPage('product-list-secondary', countryName.value)
  appStore.setBackButton(true)
  appStore.setBottomNav(false)
  appStore.recordPageLoadTime()

  // 设置当前热门项
  productsStore.setCurrentHotItem({
    code: countryCode.value,
    name: countryName.value
  })

  // 加载商品数据
  await fetchProducts()
})

onUnmounted(() => {
  // 恢复底部导航
  appStore.setBackButton(false)
  appStore.setBottomNav(true)
})
</script>

<style scoped>
.product-list-secondary-page {
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

.product-list-secondary-page {
  animation: fadeIn 0.3s ease-out;
}
</style>
