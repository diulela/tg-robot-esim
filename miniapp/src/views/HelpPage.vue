<template>
  <div class="help-page">
    <h1 class="page-title">帮助中心</h1>

    <!-- 搜索框 -->
    <div class="search-section">
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索问题..."
          class="search-input"
          @input="handleSearch"
        />
        <button class="search-btn">🔍</button>
      </div>
    </div>

    <!-- 快捷入口 -->
    <div class="quick-actions">
      <div class="action-card" @click="contactSupport">
        <div class="action-icon">💬</div>
        <div class="action-title">联系客服</div>
        <div class="action-desc">在线客服为您解答</div>
      </div>
      
      <div class="action-card" @click="reportProblem">
        <div class="action-icon">🐛</div>
        <div class="action-title">问题反馈</div>
        <div class="action-desc">报告使用问题</div>
      </div>
    </div>

    <!-- 常见问题 -->
    <div class="faq-section">
      <h3 class="section-title">常见问题</h3>
      
      <div v-if="filteredFaqs.length > 0" class="faq-list">
        <div 
          v-for="faq in filteredFaqs" 
          :key="faq.id"
          class="faq-item"
          @click="toggleFaq(faq.id)"
        >
          <div class="faq-question">
            <span class="question-text">{{ faq.question }}</span>
            <span class="toggle-icon" :class="{ expanded: expandedFaqs.includes(faq.id) }">
              ▼
            </span>
          </div>
          
          <div v-if="expandedFaqs.includes(faq.id)" class="faq-answer">
            <div v-html="faq.answer"></div>
          </div>
        </div>
      </div>
      
      <div v-else class="no-results">
        <div class="no-results-icon">🔍</div>
        <p>没有找到相关问题</p>
        <button @click="clearSearch" class="clear-search-btn">清除搜索</button>
      </div>
    </div>

    <!-- 使用指南 -->
    <div class="guide-section">
      <h3 class="section-title">使用指南</h3>
      <div class="guide-list">
        <div class="guide-item" @click="openGuide('purchase')">
          <div class="guide-icon">📱</div>
          <div class="guide-content">
            <div class="guide-title">如何购买 eSIM</div>
            <div class="guide-desc">详细的购买流程说明</div>
          </div>
          <div class="guide-arrow">›</div>
        </div>
        
        <div class="guide-item" @click="openGuide('activation')">
          <div class="guide-icon">⚡</div>
          <div class="guide-content">
            <div class="guide-title">eSIM 激活教程</div>
            <div class="guide-desc">如何激活和使用 eSIM</div>
          </div>
          <div class="guide-arrow">›</div>
        </div>
        
        <div class="guide-item" @click="openGuide('troubleshooting')">
          <div class="guide-icon">🔧</div>
          <div class="guide-content">
            <div class="guide-title">故障排除</div>
            <div class="guide-desc">常见问题解决方案</div>
          </div>
          <div class="guide-arrow">›</div>
        </div>
        
        <div class="guide-item" @click="openGuide('payment')">
          <div class="guide-icon">💳</div>
          <div class="guide-content">
            <div class="guide-title">支付帮助</div>
            <div class="guide-desc">支付方式和问题解决</div>
          </div>
          <div class="guide-arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 联系方式 -->
    <div class="contact-section">
      <h3 class="section-title">联系我们</h3>
      <div class="contact-methods">
        <div class="contact-item">
          <div class="contact-icon">📧</div>
          <div class="contact-info">
            <div class="contact-title">邮箱支持</div>
            <div class="contact-value">support@esim-store.com</div>
          </div>
        </div>
        
        <div class="contact-item">
          <div class="contact-icon">⏰</div>
          <div class="contact-info">
            <div class="contact-title">服务时间</div>
            <div class="contact-value">周一至周日 9:00-21:00</div>
          </div>
        </div>
        
        <div class="contact-item">
          <div class="contact-icon">🌐</div>
          <div class="contact-info">
            <div class="contact-title">官方网站</div>
            <div class="contact-value">www.esim-store.com</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/app'

export default {
  name: 'HelpPage',
  setup() {
    const appStore = useAppStore()
    
    const searchQuery = ref('')
    const expandedFaqs = ref([])
    
    // FAQ 数据
    const faqs = ref([
      {
        id: 1,
        question: '什么是 eSIM？',
        answer: 'eSIM（嵌入式 SIM 卡）是一种数字化的 SIM 卡，无需物理插卡即可使用移动网络服务。它直接嵌入在设备中，可以通过软件进行激活和管理。',
        category: 'basic'
      },
      {
        id: 2,
        question: '如何激活 eSIM？',
        answer: '购买成功后，您将收到包含二维码的激活信息。在设备的"设置"中找到"蜂窝网络"或"移动数据"选项，选择"添加蜂窝套餐"，然后扫描二维码即可激活。',
        category: 'activation'
      },
      {
        id: 3,
        question: '支持哪些设备？',
        answer: '支持大部分支持 eSIM 的设备，包括：<br/>• iPhone XS 及以上型号<br/>• Samsung Galaxy S20 及以上<br/>• Google Pixel 3 及以上<br/>• iPad Pro、iPad Air 等<br/>具体兼容性请查看设备说明。',
        category: 'device'
      },
      {
        id: 4,
        question: '可以使用多长时间？',
        answer: '套餐有效期根据您购买的具体套餐而定，通常为 7-30 天。有效期从激活时开始计算，过期后需要重新购买。',
        category: 'usage'
      },
      {
        id: 5,
        question: '支付失败怎么办？',
        answer: '如果支付失败，请检查：<br/>• 网络连接是否正常<br/>• 支付账户余额是否充足<br/>• 支付方式是否有效<br/>如问题持续，请联系客服协助处理。',
        category: 'payment'
      },
      {
        id: 6,
        question: '可以退款吗？',
        answer: '未激活的 eSIM 套餐可在购买后 24 小时内申请退款。已激活的套餐由于技术特性无法退款，请在购买前仔细确认套餐信息。',
        category: 'refund'
      },
      {
        id: 7,
        question: '网络速度如何？',
        answer: '网络速度取决于当地运营商网络质量和信号强度。我们的 eSIM 支持 4G/5G 网络，在信号良好的区域可以获得高速上网体验。',
        category: 'network'
      },
      {
        id: 8,
        question: '可以分享热点吗？',
        answer: '大部分套餐支持热点分享功能，但可能会消耗更多流量。具体是否支持请查看套餐详情说明。',
        category: 'hotspot'
      }
    ])
    
    // 计算属性
    const filteredFaqs = computed(() => {
      if (!searchQuery.value.trim()) {
        return faqs.value
      }
      
      const query = searchQuery.value.toLowerCase()
      return faqs.value.filter(faq => 
        faq.question.toLowerCase().includes(query) ||
        faq.answer.toLowerCase().includes(query)
      )
    })
    
    // 方法
    const handleSearch = () => {
      // 搜索时清除展开状态
      expandedFaqs.value = []
    }
    
    const clearSearch = () => {
      searchQuery.value = ''
    }
    
    const toggleFaq = (faqId) => {
      const index = expandedFaqs.value.indexOf(faqId)
      if (index > -1) {
        expandedFaqs.value.splice(index, 1)
      } else {
        expandedFaqs.value.push(faqId)
      }
    }
    
    const contactSupport = () => {
      appStore.showInfo('正在连接客服，请稍候...')
    }
    
    const reportProblem = () => {
      appStore.showInfo('问题反馈功能开发中')
    }
    
    const openGuide = (type) => {
      const guides = {
        purchase: '购买指南功能开发中',
        activation: 'eSIM 激活教程开发中',
        troubleshooting: '故障排除指南开发中',
        payment: '支付帮助功能开发中'
      }
      
      appStore.showInfo(guides[type] || '指南功能开发中')
    }
    
    // 生命周期
    onMounted(() => {
      // 可以在这里加载更多帮助内容
    })
    
    return {
      searchQuery,
      expandedFaqs,
      filteredFaqs,
      handleSearch,
      clearSearch,
      toggleFaq,
      contactSupport,
      reportProblem,
      openGuide
    }
  }
}
</script>

<style scoped>
.help-page {
  padding: 16px;
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 24px 0;
  text-align: center;
}

.search-section {
  margin-bottom: 24px;
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  padding: 4px;
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: none;
  background: transparent;
  font-size: 16px;
  color: var(--tg-theme-text-color, #000000);
}

.search-input:focus {
  outline: none;
}

.search-btn {
  padding: 8px 12px;
  background: var(--tg-theme-button-color, #0088cc);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.action-card {
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.action-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.action-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.action-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.faq-section,
.guide-section,
.contact-section {
  margin-bottom: 24px;
}

.faq-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.faq-item {
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s ease;
}

.faq-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.faq-question {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}

.question-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
}

.toggle-icon {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  transition: transform 0.2s ease;
}

.toggle-icon.expanded {
  transform: rotate(180deg);
}

.faq-answer {
  padding: 0 16px 16px 16px;
  border-top: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
}

.faq-answer div {
  font-size: 14px;
  color: var(--tg-theme-text-color, #000000);
  line-height: 1.5;
  padding-top: 12px;
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  text-align: center;
}

.no-results-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-results p {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 16px 0;
}

.clear-search-btn {
  padding: 8px 16px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
}

.guide-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.guide-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.guide-item:hover {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
}

.guide-icon {
  font-size: 24px;
  margin-right: 16px;
}

.guide-content {
  flex: 1;
}

.guide-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.guide-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.guide-arrow {
  font-size: 18px;
  color: var(--tg-theme-hint-color, #666666);
}

.contact-methods {
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
  font-size: 24px;
  margin-right: 16px;
}

.contact-info {
  flex: 1;
}

.contact-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.contact-value {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

@media (max-width: 480px) {
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .action-card {
    padding: 12px;
  }
  
  .action-icon {
    font-size: 28px;
  }
}
</style>