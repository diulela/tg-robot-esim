<template>
  <div class="error-page">
    <div class="error-container">
      <!-- 错误图标 -->
      <div class="error-icon">
        <span v-if="errorType === 'network'">📡</span>
        <span v-else-if="errorType === 'server'">🔧</span>
        <span v-else-if="errorType === 'permission'">🔒</span>
        <span v-else>❌</span>
      </div>
      
      <!-- 错误标题 -->
      <h1 class="error-title">{{ errorTitle }}</h1>
      
      <!-- 错误描述 -->
      <p class="error-description">{{ errorDescription }}</p>
      
      <!-- 错误代码 -->
      <div v-if="errorCode" class="error-code">
        错误代码: {{ errorCode }}
      </div>
      
      <!-- 操作按钮 -->
      <div class="error-actions">
        <button @click="retry" class="retry-btn">
          重试
        </button>
        
        <button @click="goHome" class="home-btn">
          返回首页
        </button>
        
        <button @click="contactSupport" class="support-btn">
          联系客服
        </button>
      </div>
      
      <!-- 错误详情（开发模式） -->
      <div v-if="isDevelopment && errorDetails" class="error-details">
        <details>
          <summary>错误详情</summary>
          <pre>{{ errorDetails }}</pre>
        </details>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

export default {
  name: 'ErrorPage',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const appStore = useAppStore()
    
    const errorType = ref('unknown')
    const errorCode = ref('')
    const errorDetails = ref('')
    
    // 计算属性
    const isDevelopment = computed(() => {
      return import.meta.env.DEV
    })
    
    const errorTitle = computed(() => {
      const titles = {
        network: '网络连接失败',
        server: '服务器错误',
        permission: '权限不足',
        notfound: '页面不存在',
        timeout: '请求超时',
        unknown: '出现了一些问题'
      }
      
      return titles[errorType.value] || titles.unknown
    })
    
    const errorDescription = computed(() => {
      const descriptions = {
        network: '请检查您的网络连接，然后重试。',
        server: '服务器暂时无法响应，请稍后重试。',
        permission: '您没有访问此内容的权限。',
        notfound: '您访问的页面不存在或已被删除。',
        timeout: '请求处理时间过长，请稍后重试。',
        unknown: '系统遇到了未知错误，我们正在努力修复。'
      }
      
      return descriptions[errorType.value] || descriptions.unknown
    })
    
    // 方法
    const parseError = () => {
      const query = route.query
      
      // 从查询参数获取错误信息
      if (query.type) {
        errorType.value = query.type
      }
      
      if (query.code) {
        errorCode.value = query.code
      }
      
      if (query.message) {
        errorDetails.value = query.message
      }
      
      // 从路由状态获取错误信息
      if (route.params.error) {
        try {
          const errorInfo = JSON.parse(route.params.error)
          errorType.value = errorInfo.type || errorType.value
          errorCode.value = errorInfo.code || errorCode.value
          errorDetails.value = errorInfo.details || errorDetails.value
        } catch (e) {
          console.error('解析错误信息失败:', e)
        }
      }
    }
    
    const retry = () => {
      // 尝试重新加载上一个页面
      if (window.history.length > 1) {
        router.back()
      } else {
        router.push({ name: 'Home' })
      }
    }
    
    const goHome = () => {
      router.push({ name: 'Home' })
    }
    
    const contactSupport = () => {
      // 构建错误报告
      const errorReport = {
        type: errorType.value,
        code: errorCode.value,
        details: errorDetails.value,
        url: window.location.href,
        userAgent: navigator.userAgent,
        timestamp: new Date().toISOString()
      }
      
      // 这里可以发送错误报告到客服系统
      console.log('错误报告:', errorReport)
      
      appStore.showInfo('错误报告已发送，客服将尽快联系您')
    }
    
    // 生命周期
    onMounted(() => {
      parseError()
      
      // 记录错误到分析系统
      if (errorType.value !== 'unknown') {
        console.error('页面错误:', {
          type: errorType.value,
          code: errorCode.value,
          details: errorDetails.value
        })
      }
    })
    
    return {
      errorType,
      errorCode,
      errorDetails,
      isDevelopment,
      errorTitle,
      errorDescription,
      retry,
      goHome,
      contactSupport
    }
  }
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.error-container {
  text-align: center;
  max-width: 400px;
  width: 100%;
}

.error-icon {
  font-size: 80px;
  margin-bottom: 24px;
  opacity: 0.8;
}

.error-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--tg-theme-text-color, #000000);
  margin: 0 0 16px 0;
}

.error-description {
  font-size: 16px;
  color: var(--tg-theme-hint-color, #666666);
  line-height: 1.5;
  margin: 0 0 20px 0;
}

.error-code {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #999999);
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  padding: 8px 12px;
  border-radius: 6px;
  margin-bottom: 32px;
  font-family: monospace;
}

.error-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.retry-btn,
.home-btn,
.support-btn {
  padding: 14px 20px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.retry-btn {
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
}

.home-btn {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  color: var(--tg-theme-text-color, #000000);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.support-btn {
  background: transparent;
  color: var(--tg-theme-button-color, #0088cc);
  border: 1px solid var(--tg-theme-button-color, #0088cc);
}

.retry-btn:hover,
.home-btn:hover,
.support-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.error-details {
  margin-top: 24px;
  text-align: left;
}

.error-details summary {
  font-size: 14px;
  color: var(--tg-theme-hint-color, #666666);
  cursor: pointer;
  padding: 8px 0;
}

.error-details pre {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 6px;
  padding: 12px;
  font-size: 12px;
  color: var(--tg-theme-text-color, #000000);
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-word;
}

@media (max-width: 480px) {
  .error-page {
    padding: 16px;
  }
  
  .error-icon {
    font-size: 64px;
  }
  
  .error-title {
    font-size: 20px;
  }
  
  .error-description {
    font-size: 14px;
  }
  
  .error-actions {
    gap: 10px;
  }
  
  .retry-btn,
  .home-btn,
  .support-btn {
    padding: 12px 16px;
    font-size: 14px;
  }
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.error-container {
  animation: fadeIn 0.5s ease-out;
}

/* 错误类型特定样式 */
.error-page[data-error-type="network"] .error-icon {
  color: #ff9800;
}

.error-page[data-error-type="server"] .error-icon {
  color: #f44336;
}

.error-page[data-error-type="permission"] .error-icon {
  color: #9c27b0;
}

.error-page[data-error-type="notfound"] .error-icon {
  color: #607d8b;
}
</style>