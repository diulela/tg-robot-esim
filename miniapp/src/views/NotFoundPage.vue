<template>
  <div class="not-found-page">
    <div class="not-found-container">
      <!-- 404 图标 -->
      <div class="not-found-icon">
        <div class="icon-404">404</div>
      </div>
      
      <!-- 标题 -->
      <h1 class="not-found-title">页面不存在</h1>
      
      <!-- 描述 -->
      <p class="not-found-description">
        抱歉，您访问的页面不存在或已被删除。
      </p>
      
      <!-- 可能的原因 -->
      <div class="possible-reasons">
        <h3>可能的原因：</h3>
        <ul>
          <li>页面地址输入错误</li>
          <li>页面已被移动或删除</li>
          <li>链接已过期</li>
          <li>您没有访问权限</li>
        </ul>
      </div>
      
      <!-- 建议操作 -->
      <div class="suggestions">
        <h3>您可以尝试：</h3>
        <div class="suggestion-actions">
          <button @click="goHome" class="suggestion-btn primary">
            <span class="btn-icon">🏠</span>
            返回首页
          </button>
          
          <button @click="goBack" class="suggestion-btn secondary">
            <span class="btn-icon">⬅️</span>
            返回上页
          </button>
          
          <button @click="searchProducts" class="suggestion-btn secondary">
            <span class="btn-icon">🔍</span>
            搜索商品
          </button>
          
          <button @click="contactSupport" class="suggestion-btn secondary">
            <span class="btn-icon">💬</span>
            联系客服
          </button>
        </div>
      </div>
      
      <!-- 热门页面推荐 -->
      <div class="popular-pages">
        <h3>热门页面：</h3>
        <div class="page-links">
          <div class="page-link" @click="goToProducts">
            <div class="link-icon">📱</div>
            <div class="link-info">
              <div class="link-title">商品列表</div>
              <div class="link-desc">浏览所有 eSIM 套餐</div>
            </div>
          </div>
          
          <div class="page-link" @click="goToRegions">
            <div class="link-icon">🌍</div>
            <div class="link-info">
              <div class="link-title">选择地区</div>
              <div class="link-desc">按地区查找套餐</div>
            </div>
          </div>
          
          <div class="page-link" @click="goToHelp">
            <div class="link-icon">❓</div>
            <div class="link-info">
              <div class="link-title">帮助中心</div>
              <div class="link-desc">常见问题解答</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 搜索框 -->
      <div class="search-section">
        <div class="search-box">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索您需要的内容..."
            class="search-input"
            @keyup.enter="performSearch"
          />
          <button @click="performSearch" class="search-btn">
            搜索
          </button>
        </div>
      </div>
    </div>
    
    <!-- 装饰性元素 -->
    <div class="decoration-elements">
      <div class="floating-element element-1">📱</div>
      <div class="floating-element element-2">🌐</div>
      <div class="floating-element element-3">⚡</div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

export default {
  name: 'NotFoundPage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()
    
    const searchQuery = ref('')
    
    // 方法
    const goHome = () => {
      router.push({ name: 'Home' })
    }
    
    const goBack = () => {
      if (window.history.length > 1) {
        router.back()
      } else {
        goHome()
      }
    }
    
    const searchProducts = () => {
      router.push({ name: 'Products' })
    }
    
    const contactSupport = () => {
      router.push({ name: 'Help' })
    }
    
    const goToProducts = () => {
      router.push({ name: 'Products' })
    }
    
    const goToRegions = () => {
      router.push({ name: 'Regions' })
    }
    
    const goToHelp = () => {
      router.push({ name: 'Help' })
    }
    
    const performSearch = () => {
      if (searchQuery.value.trim()) {
        // 这里可以实现搜索功能
        appStore.showInfo(`搜索功能开发中，搜索词：${searchQuery.value}`)
        // 或者跳转到商品列表页面
        router.push({ name: 'Products' })
      }
    }
    
    return {
      searchQuery,
      goHome,
      goBack,
      searchProducts,
      contactSupport,
      goToProducts,
      goToRegions,
      goToHelp,
      performSearch
    }
  }
}
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.not-found-container {
  background: var(--tg-theme-bg-color, #ffffff);
  border-radius: 20px;
  padding: 40px 32px;
  max-width: 500px;
  width: 100%;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
}

.not-found-icon {
  margin-bottom: 24px;
}

.icon-404 {
  font-size: 72px;
  font-weight: bold;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.not-found-title {
  font-size: 28px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.not-found-description {
  font-size: 16px;
  color: var(--tg-theme-hint-color, #666666);
  line-height: 1.5;
  margin: 0 0 32px 0;
}

.possible-reasons,
.suggestions,
.popular-pages {
  text-align: left;
  margin-bottom: 24px;
}

.possible-reasons h3,
.suggestions h3,
.popular-pages h3 {
  font-size: 16px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 12px 0;
}

.possible-reasons ul {
  margin: 0;
  padding-left: 20px;
}

.possible-reasons li {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 4px;
  line-height: 1.4;
}

.suggestion-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.suggestion-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.suggestion-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  grid-column: 1 / -1;
}

.suggestion-btn.secondary {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  color: var(--tg-theme-text-color, #000000);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.suggestion-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.btn-icon {
  font-size: 16px;
}

.page-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.page-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-link:hover {
  background: var(--tg-theme-hint-color, #e0e0e0);
  transform: translateX(4px);
}

.link-icon {
  font-size: 20px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.link-info {
  flex: 1;
}

.link-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 2px;
}

.link-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.search-section {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.search-box {
  display: flex;
  gap: 8px;
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  font-size: 14px;
  background: var(--tg-theme-bg-color, #ffffff);
  color: var(--tg-theme-text-color, #000000);
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.search-btn {
  padding: 12px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.search-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.decoration-elements {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.floating-element {
  position: absolute;
  font-size: 24px;
  opacity: 0.1;
  animation: float 6s ease-in-out infinite;
}

.element-1 {
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.element-2 {
  top: 60%;
  right: 15%;
  animation-delay: 2s;
}

.element-3 {
  bottom: 30%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

/* 页面进入动画 */
@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.not-found-container {
  animation: slideUp 0.6s ease-out;
}

@media (max-width: 480px) {
  .not-found-page {
    padding: 16px;
  }
  
  .not-found-container {
    padding: 32px 24px;
  }
  
  .icon-404 {
    font-size: 56px;
  }
  
  .not-found-title {
    font-size: 24px;
  }
  
  .not-found-description {
    font-size: 14px;
  }
  
  .suggestion-actions {
    grid-template-columns: 1fr;
  }
  
  .suggestion-btn.primary {
    grid-column: 1;
  }
  
  .search-box {
    flex-direction: column;
  }
  
  .search-input,
  .search-btn {
    width: 100%;
  }
}
</style>