<template>
  <div class="profile-page">
    <!-- 用户信息卡片 -->
    <div class="user-info-card">
      <div class="user-avatar">
        <img v-if="user.avatar" :src="user.avatar" :alt="user.name" />
        <div v-else class="avatar-placeholder">{{ getAvatarText(user.name) }}</div>
      </div>
      
      <div class="user-details">
        <h2 class="user-name">{{ user.name || '未设置昵称' }}</h2>
        <p class="user-id">ID: {{ user.id }}</p>
        <div class="user-stats">
          <div class="stat-item">
            <span class="stat-value">{{ user.orderCount || 0 }}</span>
            <span class="stat-label">订单数</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">¥{{ formatAmount(user.totalSpent || 0) }}</span>
            <span class="stat-label">总消费</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 功能菜单 -->
    <div class="menu-sections">
      <!-- 账户管理 -->
      <div class="menu-section">
        <h3 class="section-title">账户管理</h3>
        <div class="menu-items">
          <div class="menu-item" @click="goToWallet">
            <div class="menu-icon">💰</div>
            <div class="menu-content">
              <div class="menu-title">我的钱包</div>
              <div class="menu-desc">余额：¥{{ formatAmount(user.balance || 0) }}</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
          
          <div class="menu-item" @click="goToOrders">
            <div class="menu-icon">📦</div>
            <div class="menu-content">
              <div class="menu-title">我的订单</div>
              <div class="menu-desc">查看订单历史</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
          
          <div class="menu-item" @click="editProfile">
            <div class="menu-icon">👤</div>
            <div class="menu-content">
              <div class="menu-title">个人资料</div>
              <div class="menu-desc">编辑个人信息</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
        </div>
      </div>

      <!-- 应用设置 -->
      <div class="menu-section">
        <h3 class="section-title">应用设置</h3>
        <div class="menu-items">
          <div class="menu-item" @click="goToSettings">
            <div class="menu-icon">⚙️</div>
            <div class="menu-content">
              <div class="menu-title">设置</div>
              <div class="menu-desc">通知、语言等设置</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
          
          <div class="menu-item" @click="goToHelp">
            <div class="menu-icon">❓</div>
            <div class="menu-content">
              <div class="menu-title">帮助中心</div>
              <div class="menu-desc">常见问题与客服</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
          
          <div class="menu-item" @click="goToAbout">
            <div class="menu-icon">ℹ️</div>
            <div class="menu-content">
              <div class="menu-title">关于我们</div>
              <div class="menu-desc">版本信息与条款</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
        </div>
      </div>

      <!-- 其他功能 -->
      <div class="menu-section">
        <h3 class="section-title">其他</h3>
        <div class="menu-items">
          <div class="menu-item" @click="shareApp">
            <div class="menu-icon">📤</div>
            <div class="menu-content">
              <div class="menu-title">分享应用</div>
              <div class="menu-desc">推荐给朋友</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
          
          <div class="menu-item" @click="feedback">
            <div class="menu-icon">💬</div>
            <div class="menu-content">
              <div class="menu-title">意见反馈</div>
              <div class="menu-desc">帮助我们改进</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
          
          <div class="menu-item danger" @click="logout">
            <div class="menu-icon">🚪</div>
            <div class="menu-content">
              <div class="menu-title">退出登录</div>
              <div class="menu-desc">清除本地数据</div>
            </div>
            <div class="menu-arrow">›</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 版本信息 -->
    <div class="version-info">
      <p>版本 {{ appVersion }}</p>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

export default {
  name: 'ProfilePage',
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    const appStore = useAppStore()
    
    const user = ref({
      id: '123456789',
      name: '测试用户',
      avatar: '',
      orderCount: 5,
      totalSpent: 299.50,
      balance: 128.50
    })
    
    const appVersion = ref('1.0.0')
    
    // 方法
    const formatAmount = (amount) => {
      return amount.toFixed(2)
    }
    
    const getAvatarText = (name) => {
      if (!name) return '?'
      return name.charAt(0).toUpperCase()
    }
    
    const goToWallet = () => {
      router.push({ name: 'Wallet' })
    }
    
    const goToOrders = () => {
      router.push({ name: 'Orders' })
    }
    
    const editProfile = () => {
      appStore.showInfo('个人资料编辑功能开发中')
    }
    
    const goToSettings = () => {
      router.push({ name: 'Settings' })
    }
    
    const goToHelp = () => {
      router.push({ name: 'Help' })
    }
    
    const goToAbout = () => {
      router.push({ name: 'About' })
    }
    
    const shareApp = () => {
      if (navigator.share) {
        navigator.share({
          title: 'eSIM 商城',
          text: '便捷的 eSIM 购买平台',
          url: window.location.origin
        })
      } else {
        appStore.showInfo('分享功能开发中')
      }
    }
    
    const feedback = () => {
      appStore.showInfo('意见反馈功能开发中')
    }
    
    const logout = () => {
      if (confirm('确定要退出登录吗？')) {
        userStore.logout()
        appStore.showSuccess('已退出登录')
        router.push({ name: 'Home' })
      }
    }
    
    const loadUserData = async () => {
      try {
        // 模拟加载用户数据
        const userData = await userStore.getCurrentUser()
        if (userData) {
          user.value = { ...user.value, ...userData }
        }
      } catch (error) {
        console.error('加载用户数据失败:', error)
      }
    }
    
    // 生命周期
    onMounted(() => {
      loadUserData()
    })
    
    return {
      user,
      appVersion,
      formatAmount,
      getAvatarText,
      goToWallet,
      goToOrders,
      editProfile,
      goToSettings,
      goToHelp,
      goToAbout,
      shareApp,
      feedback,
      logout
    }
  }
}
</script>

<style scoped>
.profile-page {
  padding: 16px;
  min-height: 100vh;
  background: var(--tg-theme-bg-color, #ffffff);
}

.user-info-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 24px;
  color: white;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 20px;
  font-weight: bold;
  margin: 0 0 4px 0;
}

.user-id {
  font-size: 12px;
  opacity: 0.8;
  margin: 0 0 16px 0;
}

.user-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  opacity: 0.8;
}

.menu-sections {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-title {
  font-size: 14px;
  font-weight: bold;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0 0 12px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.menu-items {
  background: var(--tg-theme-bg-color, #ffffff);
  border: 1px solid var(--tg-theme-hint-color, #e0e0e0);
  border-radius: 12px;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-bottom: 1px solid var(--tg-theme-hint-color, #e0e0e0);
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:hover {
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
}

.menu-item.danger {
  color: #f44336;
}

.menu-item.danger .menu-title {
  color: #f44336;
}

.menu-icon {
  font-size: 20px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--tg-theme-secondary-bg-color, #f5f5f5);
  border-radius: 50%;
  margin-right: 16px;
  flex-shrink: 0;
}

.menu-content {
  flex: 1;
}

.menu-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--tg-theme-text-color, #000000);
  margin-bottom: 4px;
}

.menu-desc {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
}

.menu-arrow {
  font-size: 18px;
  color: var(--tg-theme-hint-color, #666666);
  margin-left: 8px;
}

.version-info {
  text-align: center;
  margin-top: 32px;
  padding: 16px;
}

.version-info p {
  font-size: 12px;
  color: var(--tg-theme-hint-color, #666666);
  margin: 0;
}

@media (max-width: 480px) {
  .user-info-card {
    flex-direction: column;
    text-align: center;
  }
  
  .user-stats {
    justify-content: center;
  }
  
  .menu-item {
    padding: 12px 16px;
  }
  
  .menu-icon {
    width: 36px;
    height: 36px;
    margin-right: 12px;
  }
}
</style>