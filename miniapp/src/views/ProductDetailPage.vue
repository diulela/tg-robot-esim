<template>
  <div class="product-detail-page">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>正在加载商品详情...</p>
    </div>

    <!-- 商品详情 -->
    <div v-else-if="product" class="product-detail">
      <!-- 商品图片和基本信息 -->
      <div class="product-hero">
        <div class="product-image">
          <img :src="product.image || '/images/default-product.png'" :alt="product.name" />
          <div v-if="product.isPopular" class="popular-badge">热门</div>
        </div>
        
        <div class="product-basic-info">
          <h1 class="product-name">{{ product.name }}</h1>
          <p class="product-description">{{ product.description }}</p>
          
          <div class="price-section">
            <span v-if="product.originalPrice && product.originalPrice > product.price" class="original-price">
              原价 ¥{{ product.originalPrice }}
            </span>
            <span class="current-price">¥{{ product.price }}</span>
            <span v-if="product.discount" class="discount-badge">
              {{ product.discount }}折
            </span>
          </div>
        </div>
      </div>

      <!-- 商品规格 -->
      <div class="product-specs-section">
        <h3 class="section-title">套餐详情</h3>
        <div class="specs-grid">
          <div v-if="product.dataAmount" class="spec-item">
            <div class="spec-icon">📶</div>
            <div class="spec-content">
              <div class="spec-label">流量</div>
              <div class="spec-value">{{ formatDataAmount(product.dataAmount) }}</div>
            </div>
          </div>
          
          <div v-if="product.validityDays" class="spec-item">
            <div class="spec-icon">📅</div>
            <div class="spec-content">
              <div class="spec-label">有效期</div>
              <div class="spec-value">{{ product.validityDays }}天</div>
            </div>
          </div>
          
          <div v-if="product.coverage" class="spec-item">
            <div class="spec-icon">🌍</div>
            <div class="spec-content">
              <div class="spec-label">覆盖范围</div>
              <div class="spec-value">{{ product.coverage }}</div>
            </div>
          </div>
          
          <div v-if="product.speed" class="spec-item">
            <div class="spec-icon">⚡</div>
            <div class="spec-content">
              <div class="spec-label">网络速度</div>
              <div class="spec-value">{{ product.speed }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 数量选择 -->
      <div class="quantity-section">
        <h3 class="section-title">购买数量</h3>
        <div class="quantity-selector">
          <button @click="decreaseQuantity" :disabled="quantity <= 1" class="quantity-btn">-</button>
          <span class="quantity-display">{{ quantity }}</span>
          <button @click="increaseQuantity" :disabled="quantity >= 10" class="quantity-btn">+</button>
        </div>
        <p class="quantity-note">单次最多购买 10 个</p>
      </div>

      <!-- 总价显示 -->
      <div class="total-price-section">
        <div class="total-price-content">
          <span class="total-label">总计</span>
          <span class="total-price">¥{{ totalPrice }}</span>
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else class="error-state">
      <div class="error-icon">❌</div>
      <h3>商品不存在</h3>
      <p>抱歉，您查看的商品不存在或已下架</p>
      <button @click="goBack" class="back-btn">返回商品列表</button>
    </div>

    <!-- 底部操作栏 -->
    <div v-if="product" class="bottom-actions">
      <button @click="addToCart" class="add-to-cart-btn">
        加入购物车
      </button>
      <button @click="buyNow" class="buy-now-btn">
        立即购买
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductsStore } from '@/stores/products'
import { useAppStore } from '@/stores/app'

export default {
  name: 'ProductDetailPage',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const productsStore = useProductsStore()
    const appStore = useAppStore()
    
    const loading = ref(false)
    const product = ref(null)
    const quantity = ref(1)
    
    // 计算属性
    const totalPrice = computed(() => {
      return product.value ? (product.value.price * quantity.value).toFixed(2) : '0.00'
    })
    
    // 方法
    const loadProduct = async () => {
      loading.value = true
      try {
        const productId = route.params.id
        product.value = await productsStore.getProductById(productId)
        
        if (!product.value) {
          console.error('商品不存在:', productId)
        }
      } catch (error) {
        console.error('加载商品详情失败:', error)
        appStore.showError('加载商品详情失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
    
    const formatDataAmount = (amount) => {
      if (amount >= 1024) {
        return `${(amount / 1024).toFixed(1)}GB`
      }
      return `${amount}MB`
    }
    
    const increaseQuantity = () => {
      if (quantity.value < 10) {
        quantity.value++
      }
    }
    
    const decreaseQuantity = () => {
      if (quantity.value > 1) {
        quantity.value--
      }
    }
    
    const addToCart = () => {
      if (!product.value) return
      
      for (let i = 0; i < quantity.value; i++) {
        productsStore.addToCart(product.value)
      }
      
      appStore.showSuccess(`已添加 ${quantity.value} 个商品到购物车`)
    }
    
    const buyNow = () => {
      if (!product.value) return
      
      // 添加到购物车
      for (let i = 0; i < quantity.value; i++) {
        productsStore.addToCart(product.value)
      }
      
      // 跳转到订单页面
      router.push({ name: 'Orders' })
    }
    
    const goBack = () => {
      router.back()
    }
    
    // 生命周期
    onMounted(() => {
      loadProduct()
    })
    
    return {
      loading,
      product,
      quantity,
      totalPrice,
      formatDataAmount,
      increaseQuantity,
      decreaseQuantity,
      addToCart,
      buyNow,
      goBack
    }
  }
}
</script>

<style scoped>
.product-detail-page {
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
  padding-bottom: 80px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--tg-theme-hint-color, #e0e0e0);
  border-top: 3px solid var(--tg-theme-button-color, #0088cc);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.product-hero {
  padding: 20px;
  border-bottom: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.product-image {
  position: relative;
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.product-image img {
  width: 120px;
  height: 120px;
  object-fit: contain;
}

.popular-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  background: #ff4757;
  color: white;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 12px;
  font-weight: bold;
}

.product-name {
  font-size: 24px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 12px 0;
  line-height: 1.3;
}

.product-description {
  font-size: 16px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 20px 0;
  line-height: 1.5;
}

.price-section {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.original-price {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #999999);
  text-decoration: line-through;
}

.current-price {
  font-size: 28px;
  font-weight: bold;
  color: #ff4757;
}

.discount-badge {
  background: #ff4757;
  color: white;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: bold;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.product-specs-section,
.quantity-section {
  padding: 20px;
  border-bottom: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.specs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
}

.spec-icon {
  font-size: 24px;
}

.spec-content {
  flex: 1;
}

.spec-label {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 4px;
}

.spec-value {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
}

.quantity-selector {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 8px;
}

.quantity-btn {
  width: 40px;
  height: 40px;
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  background: var(--tg-theme-bg-color, #ffffff);
  color: var(--tg-theme-text-color, #000000);
  border-radius: 8px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quantity-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity-display {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  min-width: 40px;
  text-align: center;
}

.quantity-note {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0;
}

.total-price-section {
  padding: 20px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
}

.total-price-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total-label {
  font-size: 16px;
  color: var(--tg-theme-text-color, #000000);
}

.total-price {
  font-size: 24px;
  font-weight: bold;
  color: #ff4757;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-state h3 {
  font-size: 18px;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 8px 0;
}

.error-state p {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 20px 0;
}

.back-btn {
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--tg-theme-bg-color, #ffffff);
  border-top: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  padding: 16px;
  display: flex;
  gap: 12px;
  z-index: 1000;
}

.add-to-cart-btn,
.buy-now-btn {
  flex: 1;
  padding: 16px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-to-cart-btn {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  color: var(--tg-theme-text-color, #000000);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.buy-now-btn {
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
}

.add-to-cart-btn:hover,
.buy-now-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

@media (max-width: 480px) {
  .specs-grid {
    grid-template-columns: 1fr;
  }
  
  .product-hero,
  .product-specs-section,
  .quantity-section,
  .total-price-section {
    padding: 16px;
  }
}
</style>