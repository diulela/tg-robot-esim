<template>
  <div class="about-page">
    <!-- 应用信息 -->
    <div class="app-info-section">
      <div class="app-logo">
        <div class="logo-icon">📱</div>
      </div>
      
      <h1 class="app-name">eSIM 商城</h1>
      <p class="app-tagline">便捷的全球 eSIM 服务平台</p>
      
      <div class="version-info">
        <span class="version-label">版本</span>
        <span class="version-number">{{ appVersion }}</span>
      </div>
    </div>

    <!-- 公司介绍 -->
    <div class="company-section">
      <h3 class="section-title">关于我们</h3>
      <div class="company-content">
        <p>
          我们致力于为全球用户提供便捷、可靠的 eSIM 服务。通过先进的技术和优质的服务，
          让您在世界各地都能享受到稳定的移动网络连接。
        </p>
        
        <div class="features-grid">
          <div class="feature-item">
            <div class="feature-icon">🌍</div>
            <div class="feature-title">全球覆盖</div>
            <div class="feature-desc">覆盖 200+ 国家和地区</div>
          </div>
          
          <div class="feature-item">
            <div class="feature-icon">⚡</div>
            <div class="feature-title">即时激活</div>
            <div class="feature-desc">购买后立即可用</div>
          </div>
          
          <div class="feature-item">
            <div class="feature-icon">🔒</div>
            <div class="feature-title">安全可靠</div>
            <div class="feature-desc">企业级安全保障</div>
          </div>
          
          <div class="feature-item">
            <div class="feature-icon">💰</div>
            <div class="feature-title">价格透明</div>
            <div class="feature-desc">无隐藏费用</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 服务条款 -->
    <div class="terms-section">
      <h3 class="section-title">服务与条款</h3>
      <div class="terms-list">
        <div class="term-item" @click="openTerms('privacy')">
          <div class="term-icon">🔐</div>
          <div class="term-content">
            <div class="term-title">隐私政策</div>
            <div class="term-desc">了解我们如何保护您的隐私</div>
          </div>
          <div class="term-arrow">›</div>
        </div>
        
        <div class="term-item" @click="openTerms('service')">
          <div class="term-icon">📋</div>
          <div class="term-content">
            <div class="term-title">服务条款</div>
            <div class="term-desc">使用服务的相关条款</div>
          </div>
          <div class="term-arrow">›</div>
        </div>
        
        <div class="term-item" @click="openTerms('refund')">
          <div class="term-icon">💳</div>
          <div class="term-content">
            <div class="term-title">退款政策</div>
            <div class="term-desc">退款相关规定说明</div>
          </div>
          <div class="term-arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 联系信息 -->
    <div class="contact-section">
      <h3 class="section-title">联系我们</h3>
      <div class="contact-info">
        <div class="contact-item">
          <div class="contact-icon">🌐</div>
          <div class="contact-details">
            <div class="contact-label">官方网站</div>
            <div class="contact-value">www.esim-store.com</div>
          </div>
        </div>
        
        <div class="contact-item">
          <div class="contact-icon">📧</div>
          <div class="contact-details">
            <div class="contact-label">客服邮箱</div>
            <div class="contact-value">support@esim-store.com</div>
          </div>
        </div>
        
        <div class="contact-item">
          <div class="contact-icon">💬</div>
          <div class="contact-details">
            <div class="contact-label">在线客服</div>
            <div class="contact-value">24/7 在线支持</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 技术信息 -->
    <div class="tech-section">
      <h3 class="section-title">技术信息</h3>
      <div class="tech-info">
        <div class="tech-item">
          <span class="tech-label">构建版本</span>
          <span class="tech-value">{{ buildNumber }}</span>
        </div>
        
        <div class="tech-item">
          <span class="tech-label">构建时间</span>
          <span class="tech-value">{{ buildTime }}</span>
        </div>
        
        <div class="tech-item">
          <span class="tech-label">运行环境</span>
          <span class="tech-value">{{ environment }}</span>
        </div>
        
        <div class="tech-item">
          <span class="tech-label">API 版本</span>
          <span class="tech-value">{{ apiVersion }}</span>
        </div>
      </div>
    </div>

    <!-- 开源许可 -->
    <div class="license-section">
      <h3 class="section-title">开源许可</h3>
      <div class="license-content">
        <p>本应用使用了以下开源项目：</p>
        <div class="license-list">
          <div class="license-item">
            <div class="license-name">Vue.js</div>
            <div class="license-type">MIT License</div>
          </div>
          
          <div class="license-item">
            <div class="license-name">Vue Router</div>
            <div class="license-type">MIT License</div>
          </div>
          
          <div class="license-item">
            <div class="license-name">Pinia</div>
            <div class="license-type">MIT License</div>
          </div>
          
          <div class="license-item">
            <div class="license-name">Vite</div>
            <div class="license-type">MIT License</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 版权信息 -->
    <div class="copyright-section">
      <p class="copyright-text">
        © {{ currentYear }} eSIM Store. All rights reserved.
      </p>
      <p class="company-info">
        深圳市某某科技有限公司
      </p>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useAppStore } from '@/stores/app'

export default {
  name: 'AboutPage',
  setup() {
    const appStore = useAppStore()
    
    const appVersion = ref('1.0.0')
    const buildNumber = ref('20241104001')
    const buildTime = ref('2024-11-04 15:30:00')
    const environment = ref('Production')
    const apiVersion = ref('v1.0')
    
    // 计算属性
    const currentYear = computed(() => {
      return new Date().getFullYear()
    })
    
    // 方法
    const openTerms = (type) => {
      const terms = {
        privacy: '隐私政策页面开发中',
        service: '服务条款页面开发中',
        refund: '退款政策页面开发中'
      }
      
      appStore.showInfo(terms[type] || '页面开发中')
    }
    
    return {
      appVersion,
      buildNumber,
      buildTime,
      environment,
      apiVersion,
      currentYear,
      openTerms
    }
  }
}
</script>

<style scoped>
.about-page {
  padding: 16px;
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
}

.app-info-section {
  text-align: center;
  padding: 32px 16px;
  margin-bottom: 24px;
}

.app-logo {
  margin-bottom: 16px;
}

.logo-icon {
  font-size: 64px;
  display: inline-block;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.app-name {
  font-size: 28px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 8px 0;
}

.app-tagline {
  font-size: 16px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 20px 0;
}

.version-info {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 20px;
}

.version-label {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.version-number {
  font-size: 12px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.company-section,
.terms-section,
.contact-section,
.tech-section,
.license-section {
  margin-bottom: 24px;
}

.company-content p {
  font-size: 14px;
  color: var(--tg-theme-text-color, #000000);
  line-height: 1.6;
  margin: 0 0 20px 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.feature-item {
  text-align: center;
  padding: 16px;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
}

.feature-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.feature-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.feature-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.terms-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.term-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.term-item:hover {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
}

.term-icon {
  font-size: 20px;
  margin-right: 16px;
}

.term-content {
  flex: 1;
}

.term-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.term-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.term-arrow {
  font-size: 18px;
  color: var(--tg-theme-hint-color, #666666);
}

.contact-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
}

.contact-icon {
  font-size: 20px;
  margin-right: 16px;
}

.contact-details {
  flex: 1;
}

.contact-label {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin-bottom: 4px;
}

.contact-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
}

.tech-info {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 12px;
  padding: 16px;
}

.tech-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.tech-item:last-child {
  border-bottom: none;
}

.tech-label {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.tech-value {
  font-size: 12px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
}

.license-content p {
  font-size: 14px;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.license-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.license-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 8px;
}

.license-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
}

.license-type {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.copyright-section {
  text-align: center;
  padding: 24px 16px;
  border-top: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  margin-top: 24px;
}

.copyright-text {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 8px 0;
}

.company-info {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0;
}

@media (max-width: 480px) {
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .app-info-section {
    padding: 24px 16px;
  }
  
  .logo-icon {
    font-size: 48px;
    padding: 12px;
  }
  
  .app-name {
    font-size: 24px;
  }
}
</style>