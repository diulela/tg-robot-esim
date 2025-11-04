<template>
  <PageWrapper
    :loading="isLoading"
    loading-text="正在加载区域信息..."
    :error="error"
    @retry="handleRetry"
    class="region-page"
  >
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">选择区域</h2>
      <p class="page-subtitle">选择您需要的 eSIM 服务区域</p>
    </div>

    <!-- 标签页导航 -->
    <div class="tabs-container">
      <v-tabs
        v-model="activeTab"
        color="primary"
        align-tabs="center"
        class="region-tabs"
      >
        <v-tab value="hot" class="tab-item">
          <v-icon start size="18">mdi-fire</v-icon>
          热门
        </v-tab>
        <v-tab value="local" class="tab-item">
          <v-icon start size="18">mdi-map-marker</v-icon>
          本地
        </v-tab>
        <v-tab value="region" class="tab-item">
          <v-icon start size="18">mdi-earth</v-icon>
          区域
        </v-tab>
        <v-tab value="global" class="tab-item">
          <v-icon start size="18">mdi-web</v-icon>
          全球
        </v-tab>
      </v-tabs>
    </div>

    <!-- 标签页内容 -->
    <div class="tabs-content">
      <v-tabs-window v-model="activeTab" class="tabs-window">
        <!-- 热门区域 -->
        <v-tabs-window-item value="hot" class="tab-content">
          <RegionGrid
            :regions="popularRegions"
            :loading="regionsLoading"
            @select-region="handleRegionSelect"
          />
        </v-tabs-window-item>

        <!-- 本地区域 -->
        <v-tabs-window-item value="local" class="tab-content">
          <RegionGrid
            :regions="localRegions"
            :loading="regionsLoading"
            @select-region="handleRegionSelect"
          />
        </v-tabs-window-item>

        <!-- 所有区域 -->
        <v-tabs-window-item value="region" class="tab-content">
          <RegionGrid
            :regions="allRegions"
            :loading="regionsLoading"
            @select-region="handleRegionSelect"
          />
        </v-tabs-window-item>

        <!-- 全球套餐 -->
        <v-tabs-window-item value="global" class="tab-content">
          <RegionGrid
            :regions="globalRegions"
            :loading="regionsLoading"
            @select-region="handleRegionSelect"
          />
        </v-tabs-window-item>
      </v-tabs-window>
    </div>

    <!-- 搜索建议 -->
    <div v-if="showSearchSuggestion" class="search-suggestion">
      <v-card variant="tonal" color="info" class="suggestion-card">
        <v-card-text class="suggestion-content">
          <v-icon color="info" size="20">mdi-lightbulb-outline</v-icon>
          <span class="suggestion-text">找不到合适的区域？试试搜索具体国家</span>
          <v-btn
            variant="text"
            color="info"
            size="small"
            @click="navigateToCountries"
          >
            搜索国家
          </v-btn>
        </v-card-text>
      </v-card>
    </div>
  </PageWrapper>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useProductsStore } from '@/stores/products'
import { telegramService } from '@/services/telegram'
import type { Region } from '@/types'

import PageWrapper from '@/components/layout/PageWrapper.vue'
import RegionGrid from '@/components/business/RegionGrid.vue'

// 组合式 API
const router = useRouter()
const appStore = useAppStore()
const productsStore = useProductsStore()

// 响应式状态
const isLoading = ref(false)
const error = ref<string | null>(null)
const regionsLoading = ref(false)
const activeTab = ref('hot')
const showSearchSuggestion = ref(false)

// 计算属性
const popularRegions = computed(() => {
  return productsStore.popularRegions
})

const localRegions = computed(() => {
  // 根据用户位置或 Telegram 语言设置推荐本地区域
  const userLanguage = appStore.currentLanguage
  
  if (userLanguage.startsWith('zh')) {
    // 中文用户推荐亚洲区域
    return productsStore.regions.filter(region => 
      ['asia', 'china', 'hongkong', 'taiwan'].includes(region.code.toLowerCase())
    )
  } else {
    // 其他用户推荐欧美区域
    return productsStore.regions.filter(region => 
      ['europe', 'north-america', 'usa'].includes(region.code.toLowerCase())
    )
  }
})

const allRegions = computed(() => {
  return productsStore.regions.filter(region => 
    !region.code.toLowerCase().includes('global')
  )
})

const globalRegions = computed(() => {
  return productsStore.regions.filter(region => 
    region.code.toLowerCase().includes('global') || 
    region.name.toLowerCase().includes('全球')
  )
})

// 方法
const loadRegions = async () => {
  isLoading.value = true
  error.value = null
  regionsLoading.value = true

  try {
    await productsStore.fetchRegions()
    console.log('[RegionPage] 区域数据加载成功')
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : '加载区域数据失败'
    error.value = errorMessage
    console.error('[RegionPage] 加载区域数据失败:', err)
  } finally {
    isLoading.value = false
    regionsLoading.value = false
  }
}

const handleRetry = () => {
  loadRegions()
}

const handleRegionSelect = async (region: Region) => {
  try {
    // 设置选中的区域
    await productsStore.setSelectedRegion(region)
    
    // 触觉反馈
    telegramService.selectionFeedback()
    
    // 导航到国家列表页面
    router.push({ 
      name: 'Countries', 
      params: { region: region.code } 
    })
    
    console.log('[RegionPage] 选择区域:', region.name)
  } catch (err) {
    console.error('[RegionPage] 选择区域失败:', err)
    
    appStore.showNotification({
      type: 'error',
      message: '选择区域失败，请重试',
      duration: 3000
    })
  }
}

const navigateToCountries = () => {
  router.push({ name: 'Countries' })
  telegramService.selectionFeedback()
}

// 监听标签页变化
watch(activeTab, (newTab) => {
  console.log('[RegionPage] 切换到标签页:', newTab)
  
  // 根据标签页内容决定是否显示搜索建议
  showSearchSuggestion.value = newTab === 'region' && allRegions.value.length > 8
  
  // 触觉反馈
  telegramService.selectionFeedback()
})

// 生命周期
onMounted(async () => {
  console.log('[RegionPage] 组件挂载')
  
  // 设置页面标题
  appStore.setCurrentPage('Regions', '选择区域')
  
  // 加载区域数据
  await loadRegions()
  
  // 检查是否需要显示搜索建议
  showSearchSuggestion.value = activeTab.value === 'region' && allRegions.value.length > 8
})
</script>

<style scoped lang="scss">
.region-page {
  .page-header {
    text-align: center;
    margin-bottom: 24px;
    
    .page-title {
      font-size: 1.5rem;
      font-weight: 600;
      color: rgb(var(--v-theme-on-surface));
      margin: 0 0 8px 0;
    }
    
    .page-subtitle {
      font-size: 0.875rem;
      color: rgba(var(--v-theme-on-surface), 0.7);
      margin: 0;
      line-height: 1.4;
    }
  }
  
  .tabs-container {
    margin-bottom: 24px;
    
    .region-tabs {
      .tab-item {
        font-size: 0.875rem;
        font-weight: 500;
        min-width: 0;
        flex: 1;
        
        .v-icon {
          margin-right: 4px;
        }
      }
    }
  }
  
  .tabs-content {
    .tabs-window {
      .tab-content {
        padding-top: 8px;
      }
    }
  }
  
  .search-suggestion {
    margin-top: 32px;
    
    .suggestion-card {
      .suggestion-content {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 16px !important;
        
        .suggestion-text {
          flex: 1;
          font-size: 0.875rem;
          line-height: 1.4;
        }
      }
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .region-page {
    .page-header {
      margin-bottom: 20px;
      
      .page-title {
        font-size: 1.25rem;
      }
      
      .page-subtitle {
        font-size: 0.8125rem;
      }
    }
    
    .tabs-container {
      margin-bottom: 20px;
      
      .region-tabs {
        .tab-item {
          font-size: 0.8125rem;
          
          .v-icon {
            margin-right: 2px;
          }
        }
      }
    }
    
    .search-suggestion {
      margin-top: 24px;
      
      .suggestion-card {
        .suggestion-content {
          gap: 8px;
          padding: 12px !important;
          
          .suggestion-text {
            font-size: 0.8125rem;
          }
        }
      }
    }
  }
}

@media (min-width: 481px) {
  .region-page {
    .page-header {
      margin-bottom: 32px;
      
      .page-title {
        font-size: 1.75rem;
      }
      
      .page-subtitle {
        font-size: 0.9375rem;
      }
    }
    
    .tabs-container {
      margin-bottom: 32px;
      
      .region-tabs {
        .tab-item {
          font-size: 0.9375rem;
          
          .v-icon {
            margin-right: 6px;
          }
        }
      }
    }
    
    .search-suggestion {
      margin-top: 40px;
      
      .suggestion-card {
        .suggestion-content {
          gap: 16px;
          padding: 20px !important;
          
          .suggestion-text {
            font-size: 0.9375rem;
          }
        }
      }
    }
  }
}

// 标签页过渡动画
.v-tabs-window {
  .v-tabs-window-item {
    transition: all 0.3s ease;
  }
}

// 深色主题适配
.v-theme--dark {
  .region-page {
    .page-header {
      .page-title {
        color: rgb(var(--v-theme-on-surface));
      }
      
      .page-subtitle {
        color: rgba(var(--v-theme-on-surface), 0.7);
      }
    }
  }
}
</style>