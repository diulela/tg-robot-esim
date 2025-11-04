<template>
  <div class="product-list-page">
    <!-- 国家选择视图 -->
    <div v-if="!selectedCountry" class="country-selection-view">
      <!-- 页面头部 -->
      <div class="page-header">
        <h1 class="page-title">{{ pageTitle }}</h1>
      </div>

      <!-- 栏目导航 -->
      <div class="category-tabs">
        <div v-for="tab in categoryTabs" :key="tab.key" @click="switchCategory(tab.key)"
          :class="['tab-item', { active: activeCategory === tab.key }]">
          {{ tab.label }}
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>正在加载...</p>
      </div>

      <!-- 国家/地区列表 -->
      <div v-else class="countries-grid">
        <div v-for="country in currentCountries" :key="country.code" class="country-card" @click="selectCountry(country)">
          <div class="country-flag">{{ country.code }}</div>
          <div class="country-name">{{ country.name }}</div>
        </div>
      </div>
    </div>

    <!-- 商品列表视图 -->
    <div v-else class="product-list-view">
      <!-- 商品列表头部 -->
      <div class="product-header">
        <button @click="goBack" class="back-btn">
          <span class="back-icon">‹</span>
        </button>
        <h1 class="country-title">{{ selectedCountry.name }}</h1>
        <div class="header-spacer"></div>
      </div>

      <!-- 栏目导航 -->
      <div class="category-tabs">
        <div v-for="tab in categoryTabs" :key="tab.key" @click="switchCategory(tab.key)"
          :class="['tab-item', { active: activeCategory === tab.key }]">
          {{ tab.label }}
        </div>
      </div>

      <!-- 商品加载状态 -->
      <div v-if="loadingProducts" class="loading-container">
        <div class="loading-spinner"></div>
        <p>正在加载商品...</p>
      </div>

      <!-- 商品列表 -->
      <div v-else class="products-container">
        <div v-for="product in products" :key="product.id" class="product-card" @click="goToProductDetail(product.id)">
          <!-- 商品图标 -->
          <div class="product-icon">
            <img :src="product.icon" :alt="product.name" />
          </div>

          <!-- 商品信息 -->
          <div class="product-info">
            <h3 class="product-name">{{ product.name }}</h3>
            <p class="product-description">{{ product.description }}</p>
            
            <!-- 商品规格 -->
            <div class="product-specs">
              <span class="spec-item">
                <span class="spec-icon">📶</span>
                {{ product.data }}
              </span>
              <span class="spec-item">
                <span class="spec-icon">⏰</span>
                {{ product.validity }}
              </span>
              <span class="spec-item">
                <span class="spec-icon">🌍</span>
                {{ product.coverage }}
              </span>
            </div>

            <!-- 商品特性 -->
            <div class="product-features">
              <span v-for="feature in product.features" :key="feature" class="feature-tag">
                <span class="feature-icon">✓</span>
                {{ feature }}
              </span>
            </div>

            <!-- 价格和购买 -->
            <div class="product-footer">
              <div class="price-section">
                <div class="current-price">{{ product.currentPrice }}</div>
                <div v-if="product.originalPrice" class="original-price">{{ product.originalPrice }}</div>
              </div>
              <button @click.stop="buyProduct(product)" class="buy-btn">
                <span class="cart-icon">🛒</span>
                立即购买
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

export default {
  name: 'ProductListPage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()

    const loading = ref(false)
    const activeCategory = ref('hot')

    // 栏目配置
    const categoryTabs = [
      { key: 'hot', label: '热门' },
      { key: 'local', label: '本地' },
      { key: 'region', label: '区域' },
      { key: 'global', label: '全球' }
    ]

    // 国家/地区数据
    const countriesData = {
      hot: [
        { code: 'CN', name: '中国' },
        { code: 'HK', name: '香港' },
        { code: 'TW', name: '台湾' },
        { code: 'JP', name: '日本' },
        { code: 'VN', name: '越南' },
        { code: 'US', name: '美国' },
        { code: 'MO', name: '澳门' },
        { code: 'TH', name: '泰国' },
        { code: 'KR', name: '韩国' },
        { code: 'SG', name: '新加坡' },
        { code: 'MY', name: '马来西亚' },
        { code: 'AU', name: '澳大利亚' },
        { code: 'GB', name: '英国' }
      ],
      local: [
        { code: 'CN', name: '中国' },
        { code: 'HK', name: '香港' },
        { code: 'MO', name: '澳门' },
        { code: 'TW', name: '台湾' }
      ],
      region: [
        { code: 'JP', name: '日本' },
        { code: 'KR', name: '韩国' },
        { code: 'TH', name: '泰国' },
        { code: 'VN', name: '越南' },
        { code: 'SG', name: '新加坡' },
        { code: 'MY', name: '马来西亚' },
        { code: 'PH', name: '菲律宾' },
        { code: 'ID', name: '印度尼西亚' }
      ],
      global: [
        { code: 'US', name: '美国' },
        { code: 'CA', name: '加拿大' },
        { code: 'GB', name: '英国' },
        { code: 'FR', name: '法国' },
        { code: 'DE', name: '德国' },
        { code: 'IT', name: '意大利' },
        { code: 'ES', name: '西班牙' },
        { code: 'AU', name: '澳大利亚' },
        { code: 'NZ', name: '新西兰' }
      ]
    }

    // 计算属性
    const currentCountries = computed(() => {
      return countriesData[activeCategory.value] || []
    })

    const pageTitle = computed(() => {
      const titles = {
        hot: '热门国家',
        local: '本地区域',
        region: '亚洲地区',
        global: '全球覆盖'
      }
      return titles[activeCategory.value] || '热门国家'
    })

    // 方法
    const switchCategory = (category) => {
      activeCategory.value = category
    }

    const selectCountry = (country) => {
      // 跳转到该国家的商品列表
      router.push({
        name: 'Countries',
        params: { region: country.code.toLowerCase() },
        query: { country: country.name }
      })
    }

    const loadData = async () => {
      loading.value = true
      try {
        // 模拟加载数据
        await new Promise(resolve => setTimeout(resolve, 500))
      } catch (error) {
        console.error('加载数据失败:', error)
        appStore.showError('加载数据失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }

    // 生命周期
    onMounted(() => {
      loadData()
    })

    return {
      loading,
      activeCategory,
      categoryTabs,
      currentCountries,
      pageTitle,
      switchCategory,
      selectCountry
    }
  }
}
</script>

<style scoped>
.product-list-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.page-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 16px;
  text-align: center;
}

.page-title {
  font-size: 20px;
  font-weight: bold;
  color: white;
  margin: 0;
}

.category-tabs {
  background: white;
  display: flex;
  padding: 0;
  margin: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: white;
  margin: 0 16px;
  border-radius: 12px;
  margin-top: 16px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e0e0e0;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.countries-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  background: #e0e0e0;
  margin: 0;
  padding: 0;
}

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
}

.country-card:hover {
  background: #f5f5f5;
  transform: scale(1.02);
}

.country-card:active {
  background: #e8e8e8;
  transform: scale(0.98);
}

.country-flag {
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

/* 特殊处理最后一行不满3个的情况 */
.country-card:nth-last-child(1):nth-child(3n-1) {
  grid-column: span 2;
}

.country-card:nth-last-child(1):nth-child(3n-2) {
  grid-column: span 3;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .page-header {
    padding: 16px;
  }

  .page-title {
    font-size: 18px;
  }

  .tab-item {
    padding: 14px 8px;
    font-size: 14px;
  }

  .countries-grid {
    grid-template-columns: repeat(3, 1fr);
  }

  .country-card {
    padding: 24px 12px;
    min-height: 100px;
  }

  .country-flag {
    font-size: 28px;
  }

  .country-name {
    font-size: 12px;
  }
}

@media (max-width: 360px) {
  .country-card {
    padding: 20px 8px;
    min-height: 90px;
  }

  .country-flag {
    font-size: 24px;
  }

  .country-name {
    font-size: 11px;
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

.countries-grid {
  animation: fadeInUp 0.5s ease-out;
}

.country-card {
  animation: fadeInUp 0.5s ease-out;
}

.country-card:nth-child(1) {
  animation-delay: 0.1s;
}

.country-card:nth-child(2) {
  animation-delay: 0.15s;
}

.country-card:nth-child(3) {
  animation-delay: 0.2s;
}

.country-card:nth-child(4) {
  animation-delay: 0.25s;
}

.country-card:nth-child(5) {
  animation-delay: 0.3s;
}

.country-card:nth-child(6) {
  animation-delay: 0.35s;
}

.country-card:nth-child(7) {
  animation-delay: 0.4s;
}

.country-card:nth-child(8) {
  animation-delay: 0.45s;
}

.country-card:nth-child(9) {
  animation-delay: 0.5s;
}
</style>