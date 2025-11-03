<template>
  <v-bottom-navigation
    v-model="currentValue"
    :color="color"
    :grow="grow"
    :height="height"
    class="bottom-navigation"
    :class="navigationClasses"
  >
    <v-btn
      v-for="item in navigationItems"
      :key="item.value"
      :value="item.value"
      :disabled="item.disabled"
      class="nav-item"
      :class="{ 'nav-item--active': currentValue === item.value }"
      @click="handleItemClick(item)"
    >
      <!-- 图标 -->
      <v-icon :size="iconSize">{{ item.icon }}</v-icon>
      
      <!-- 标签 -->
      <span class="nav-label">{{ item.title }}</span>
      
      <!-- 徽章 -->
      <v-badge
        v-if="item.badge"
        :content="item.badge"
        :color="badgeColor"
        class="nav-badge"
        floating
        offset-x="8"
        offset-y="8"
      />
    </v-btn>
  </v-bottom-navigation>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { telegramService } from '@/services/telegram'

// 导航项接口
interface NavigationItem {
  value: string
  title: string
  icon: string
  route?: string
  badge?: string | number
  disabled?: boolean
  action?: () => void
}

// Props
interface Props {
  modelValue?: string | number
  items?: NavigationItem[]
  color?: string
  grow?: boolean
  height?: string | number
  iconSize?: string | number
  badgeColor?: string
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 'home',
  items: () => [
    {
      value: 'home',
      title: '首页',
      icon: 'mdi-home',
      route: 'Home'
    },
    {
      value: 'products',
      title: '商品',
      icon: 'mdi-shopping',
      route: 'Products'
    },
    {
      value: 'orders',
      title: '订单',
      icon: 'mdi-receipt',
      route: 'Orders'
    },
    {
      value: 'profile',
      title: '我的',
      icon: 'mdi-account',
      route: 'Profile'
    }
  ],
  color: 'primary',
  grow: true,
  height: 64,
  iconSize: 24,
  badgeColor: 'error'
})

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  'change': [value: string | number, item: NavigationItem]
  'click': [item: NavigationItem]
}>()

// 组合式 API
const router = useRouter()
const route = useRoute()

// 计算属性
const currentValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const navigationItems = computed(() => props.items)

const navigationClasses = computed(() => {
  const classes = []
  
  if (props.class) {
    classes.push(props.class)
  }
  
  return classes
})

// 方法
const handleItemClick = async (item: NavigationItem) => {
  // 触觉反馈
  telegramService.selectionFeedback()
  
  // 更新当前值
  currentValue.value = item.value
  
  // 发射事件
  emit('change', item.value, item)
  emit('click', item)
  
  // 执行自定义动作或路由导航
  if (item.action) {
    item.action()
  } else if (item.route && !item.disabled) {
    try {
      await router.push({ name: item.route })
    } catch (error) {
      console.error('导航失败:', error)
    }
  }
}

// 根据当前路由自动更新选中项
const updateActiveItem = () => {
  const routeItemMap: Record<string, string> = {
    'Home': 'home',
    'Products': 'products',
    'ProductDetail': 'products',
    'Orders': 'orders',
    'OrderDetail': 'orders',
    'Profile': 'profile',
    'Wallet': 'profile',
    'Settings': 'profile'
  }
  
  const activeValue = routeItemMap[route.name as string]
  if (activeValue && activeValue !== currentValue.value) {
    currentValue.value = activeValue
  }
}

// 监听路由变化
watch(
  () => route.name,
  () => updateActiveItem(),
  { immediate: true }
)
</script>

<style scoped lang="scss">
.bottom-navigation {
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  
  .nav-item {
    position: relative;
    min-width: 64px;
    transition: all 0.2s ease;
    
    .v-icon {
      margin-bottom: 4px;
      transition: all 0.2s ease;
    }
    
    .nav-label {
      font-size: 0.75rem;
      line-height: 1;
      font-weight: 500;
      transition: all 0.2s ease;
    }
    
    .nav-badge {
      position: absolute;
      top: 8px;
      right: 12px;
    }
    
    // 激活状态
    &--active {
      .v-icon {
        transform: scale(1.1);
      }
      
      .nav-label {
        font-weight: 600;
      }
    }
    
    // 禁用状态
    &:disabled {
      opacity: 0.5;
      
      .v-icon,
      .nav-label {
        color: rgba(var(--v-theme-on-surface), 0.38) !important;
      }
    }
    
    // 悬停效果 (桌面端)
    @media (hover: hover) {
      &:hover:not(:disabled) {
        .v-icon {
          transform: scale(1.05);
        }
      }
    }
    
    // 点击效果
    &:active:not(:disabled) {
      transform: scale(0.98);
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .bottom-navigation {
    .nav-item {
      min-width: 56px;
      
      .nav-label {
        font-size: 0.7rem;
      }
      
      .v-icon {
        font-size: 20px !important;
      }
    }
  }
}

@media (min-width: 481px) {
  .bottom-navigation {
    .nav-item {
      min-width: 72px;
      
      .nav-label {
        font-size: 0.8rem;
      }
      
      .v-icon {
        font-size: 26px !important;
      }
    }
  }
}

// 深色主题适配
.v-theme--dark {
  .bottom-navigation {
    border-top-color: rgba(var(--v-theme-on-surface), 0.2);
  }
}

// 安全区域适配
@supports (padding-bottom: env(safe-area-inset-bottom)) {
  .bottom-navigation {
    padding-bottom: env(safe-area-inset-bottom);
  }
}
</style>