<template>
  <div class="settings-page">
    <h1 class="page-title">设置</h1>

    <!-- 通知设置 -->
    <div class="settings-section">
      <h3 class="section-title">通知设置</h3>
      <div class="settings-items">
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">订单通知</div>
            <div class="setting-desc">订单状态变更时推送通知</div>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input v-model="settings.orderNotification" type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">促销通知</div>
            <div class="setting-desc">优惠活动和促销信息推送</div>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input v-model="settings.promotionNotification" type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">系统通知</div>
            <div class="setting-desc">系统维护和重要公告</div>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input v-model="settings.systemNotification" type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- 显示设置 -->
    <div class="settings-section">
      <h3 class="section-title">显示设置</h3>
      <div class="settings-items">
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">主题模式</div>
            <div class="setting-desc">跟随系统或手动选择</div>
          </div>
          <div class="setting-control">
            <select v-model="settings.theme" class="theme-select">
              <option value="auto">跟随系统</option>
              <option value="light">浅色模式</option>
              <option value="dark">深色模式</option>
            </select>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">语言</div>
            <div class="setting-desc">选择应用显示语言</div>
          </div>
          <div class="setting-control">
            <select v-model="settings.language" class="language-select">
              <option value="zh-CN">简体中文</option>
              <option value="zh-TW">繁體中文</option>
              <option value="en">English</option>
            </select>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">货币单位</div>
            <div class="setting-desc">价格显示的货币单位</div>
          </div>
          <div class="setting-control">
            <select v-model="settings.currency" class="currency-select">
              <option value="CNY">人民币 (¥)</option>
              <option value="USD">美元 ($)</option>
              <option value="EUR">欧元 (€)</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <!-- 隐私设置 -->
    <div class="settings-section">
      <h3 class="section-title">隐私设置</h3>
      <div class="settings-items">
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">数据分析</div>
            <div class="setting-desc">帮助改进应用体验</div>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input v-model="settings.analytics" type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">个性化推荐</div>
            <div class="setting-desc">基于使用习惯推荐商品</div>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input v-model="settings.personalization" type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- 缓存设置 -->
    <div class="settings-section">
      <h3 class="section-title">存储设置</h3>
      <div class="settings-items">
        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">缓存大小</div>
            <div class="setting-desc">当前缓存: {{ cacheSize }}</div>
          </div>
          <div class="setting-control">
            <button @click="clearCache" class="clear-cache-btn">清理缓存</button>
          </div>
        </div>

        <div class="setting-item">
          <div class="setting-info">
            <div class="setting-title">自动清理</div>
            <div class="setting-desc">定期清理过期缓存</div>
          </div>
          <div class="setting-control">
            <label class="switch">
              <input v-model="settings.autoClearCache" type="checkbox" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- 关于信息 -->
    <div class="settings-section">
      <h3 class="section-title">关于</h3>
      <div class="settings-items">
        <div class="setting-item" @click="checkUpdate">
          <div class="setting-info">
            <div class="setting-title">检查更新</div>
            <div class="setting-desc">当前版本: {{ appVersion }}</div>
          </div>
          <div class="setting-control">
            <span class="arrow">›</span>
          </div>
        </div>

        <div class="setting-item" @click="goToHelp">
          <div class="setting-info">
            <div class="setting-title">帮助中心</div>
            <div class="setting-desc">常见问题与使用指南</div>
          </div>
          <div class="setting-control">
            <span class="arrow">›</span>
          </div>
        </div>

        <div class="setting-item" @click="goToAbout">
          <div class="setting-info">
            <div class="setting-title">关于我们</div>
            <div class="setting-desc">了解更多信息</div>
          </div>
          <div class="setting-control">
            <span class="arrow">›</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 保存按钮 -->
    <div class="save-section">
      <button @click="saveSettings" :disabled="!hasChanges" class="save-btn">
        保存设置
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'

export default {
  name: 'SettingsPage',
  setup() {
    const router = useRouter()
    const appStore = useAppStore()
    
    const settings = ref({
      orderNotification: true,
      promotionNotification: false,
      systemNotification: true,
      theme: 'auto',
      language: 'zh-CN',
      currency: 'CNY',
      analytics: true,
      personalization: true,
      autoClearCache: false
    })
    
    const originalSettings = ref({})
    const cacheSize = ref('12.5 MB')
    const appVersion = ref('1.0.0')
    
    // 计算属性
    const hasChanges = computed(() => {
      return JSON.stringify(settings.value) !== JSON.stringify(originalSettings.value)
    })
    
    // 方法
    const loadSettings = () => {
      try {
        const savedSettings = localStorage.getItem('app-settings')
        if (savedSettings) {
          const parsed = JSON.parse(savedSettings)
          settings.value = { ...settings.value, ...parsed }
        }
        originalSettings.value = { ...settings.value }
      } catch (error) {
        console.error('加载设置失败:', error)
      }
    }
    
    const saveSettings = async () => {
      try {
        localStorage.setItem('app-settings', JSON.stringify(settings.value))
        originalSettings.value = { ...settings.value }
        
        // 应用主题设置
        applyThemeSettings()
        
        appStore.showSuccess('设置已保存')
      } catch (error) {
        console.error('保存设置失败:', error)
        appStore.showError('保存设置失败，请稍后重试')
      }
    }
    
    const applyThemeSettings = () => {
      const { theme } = settings.value
      
      if (theme === 'dark') {
        document.documentElement.classList.add('dark-theme')
      } else if (theme === 'light') {
        document.documentElement.classList.remove('dark-theme')
      } else {
        // 跟随系统
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        if (prefersDark) {
          document.documentElement.classList.add('dark-theme')
        } else {
          document.documentElement.classList.remove('dark-theme')
        }
      }
    }
    
    const clearCache = async () => {
      if (confirm('确定要清理缓存吗？这将删除所有本地缓存数据。')) {
        try {
          // 清理各种缓存
          if ('caches' in window) {
            const cacheNames = await caches.keys()
            await Promise.all(
              cacheNames.map(cacheName => caches.delete(cacheName))
            )
          }
          
          // 清理 localStorage 中的缓存数据（保留设置）
          const keysToKeep = ['app-settings', 'user-data']
          const allKeys = Object.keys(localStorage)
          allKeys.forEach(key => {
            if (!keysToKeep.includes(key)) {
              localStorage.removeItem(key)
            }
          })
          
          cacheSize.value = '0 MB'
          appStore.showSuccess('缓存已清理')
        } catch (error) {
          console.error('清理缓存失败:', error)
          appStore.showError('清理缓存失败，请稍后重试')
        }
      }
    }
    
    const checkUpdate = () => {
      appStore.showInfo('当前已是最新版本')
    }
    
    const goToHelp = () => {
      router.push({ name: 'Help' })
    }
    
    const goToAbout = () => {
      router.push({ name: 'About' })
    }
    
    // 监听设置变化，自动应用某些设置
    watch(() => settings.value.theme, () => {
      applyThemeSettings()
    })
    
    // 生命周期
    onMounted(() => {
      loadSettings()
      applyThemeSettings()
    })
    
    return {
      settings,
      cacheSize,
      appVersion,
      hasChanges,
      saveSettings,
      clearCache,
      checkUpdate,
      goToHelp,
      goToAbout
    }
  }
}
</script>

<style scoped>
.settings-page {
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

.settings-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 12px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.settings-items {
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  overflow: hidden;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  cursor: pointer;
  transition: all 0.2s ease;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item:hover {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
}

.setting-info {
  flex: 1;
}

.setting-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.setting-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.setting-control {
  display: flex;
  align-items: center;
}

/* 开关样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 28px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--tg-theme-hint-color, #ccc);
  transition: 0.3s;
  border-radius: 28px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: var(--tg-theme-button-color, #0088cc);
}

input:checked + .slider:before {
  transform: translateX(20px);
}

/* 选择框样式 */
.theme-select,
.language-select,
.currency-select {
  padding: 8px 12px;
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 6px;
  background: var(--tg-theme-bg-color, #ffffff);
  color: var(--tg-theme-text-color, #000000);
  font-size: 14px;
  min-width: 120px;
}

.clear-cache-btn {
  padding: 8px 16px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.clear-cache-btn:hover {
  opacity: 0.9;
}

.arrow {
  font-size: 18px;
  color: var(--tg-theme-hint-color, #666666);
}

.save-section {
  margin-top: 32px;
  padding-bottom: 32px;
}

.save-btn {
  width: 100%;
  padding: 16px;
  background: var(--tg-theme-button-color, #0088cc);
  color: var(--tg-theme-button-text-color, #ffffff);
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.save-btn:not(:disabled):hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

/* 深色主题支持 */
:global(.dark-theme) .settings-page {
  background: #1a1a1a;
  color: #ffffff;
}

:global(.dark-theme) .settings-items {
  background: #2a2a2a;
  border-color: #404040;
}

:global(.dark-theme) .setting-item {
  border-color: #404040;
}

:global(.dark-theme) .setting-item:hover {
  background: #333333;
}

:global(.dark-theme) .theme-select,
:global(.dark-theme) .language-select,
:global(.dark-theme) .currency-select {
  background: #2a2a2a;
  border-color: #404040;
  color: #ffffff;
}

@media (max-width: 480px) {
  .setting-item {
    padding: 12px 16px;
  }
  
  .theme-select,
  .language-select,
  .currency-select {
    min-width: 100px;
    font-size: 12px;
  }
}
</style>