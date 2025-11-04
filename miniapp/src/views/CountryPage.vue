<template>
  <PageWrapper
    :loading="isLoading"
    loading-text="正在加载国家列表..."
    :error="error"
    @retry="handleRetry"
    class="country-page"
  >
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        {{ selectedRegion ? `${selectedRegion.name}国家` : '选择国家' }}
      </h2>
      <p class="page-subtitle">选择您需要 eSIM 服务的国家</p>
    </div>

    <!-- 搜索栏 -->
    <div class="search-section">
      <v-text-field
        v-model="searchQuery"
        placeholder="搜索国家或地区"
        prepend-inner-icon="mdi-magnify"
        variant="outlined"
        density="comfortable"
        clearable
        hide-details
        class="search-field"
        @input="handleSearch"
        @clear="handleSearchClear"
      />
    </div>

    <!-- 字母索引 -->
    <div class="content-container">
      <div class="countries-list">
        <!-- 搜索结果提示 -->
        <div v-if="searchQuery && filteredCountries.length > 0" class="search-results">
          <p class="results-text">
            找到 {{ filteredCountries.length }} 个匹配的国家
          </p>
        </div>

        <!-- 国家列表 -->
        <div v-if="groupedCountries.length > 0" class="countries-content">
          <div
            v-for="group in groupedCountries"
            :key="group.letter"
            class="country-group"
          >
            <!-- 字母分组标题 -->
            <div class="group-header">
              <h3 class="group-letter">{{ group.letter }}</h3>
            </div>

            <!-- 国家项目列表 -->
            <div class="group-countries">
              <CountryItem
                v-for="country in group.countries"
                :key="country.id"
                :country="country"
                @click="handleCountrySelect"
              />
            </div>
          </div>
        </div>

        <!-- 搜索无结果 -->
        <div v-else-if="searchQuery" class="no-results">
          <v-icon size="48" color="grey-lighten-1">mdi-magnify</v-icon>
          <h4 class="no-results-title">未找到匹配的国家</h4>
          <p class="no-results-subtitle">
            请尝试其他关键词或浏览所有国家
          </p>
          <v-btn
            variant="outlined"
            color="primary"
            @click="clearSearch"
            class="mt-4"
          >
            浏览所有国家
          </v-btn>
        </div>

        <!-- 空状态 -->
        <div v-else-if="!isLoading" class="empty-state">
          <v-icon size="48" color="grey-lighten-1">mdi-earth-off</v-icon>
          <h4 class="empty-title">暂无国家数据</h4>
          <p class="empty-subtitle">请稍后重试或选择其他区域</p>
        </div>
      </div>

      <!-- 字母索引导航 -->
      <div class="alphabet-index">
        <div class="index-container">
          <button
            v-for="letter in alphabetIndex"
            :key="letter"
            :class="['index-letter', { active: activeLetter === letter }]"
            @click="scrollToLetter(letter)"
          >
            {{ letter }}
          </button>
        </div>
      </div>
    </div>

    <!-- 浮动操作按钮 -->
    <v-fab
      icon="mdi-arrow-up"
      color="primary"
      size="small"
      location="bottom start"
      class="scroll-top-fab"
      @click="scrollToTop"
      v-show="showScrollTop"
    />
  </PageWrapper>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useProductsStore } from '@/stores/products'
import { telegramService } from '@/services/telegram'
import type { Country } from '@/types'

import PageWrapper from '@/components/layout/PageWrapper.vue'
import CountryItem from '@/components/business/CountryItem.vue'

// Props
interface Props {
  region?: string
}

const props = defineProps<Props>()

// 组合式 API
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const productsStore = useProductsStore()

// 响应式状态
const isLoading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const activeLetter = ref('A')
const showScrollTop = ref(false)

// 计算属性
const selectedRegion = computed(() => {
  return productsStore.selectedRegion
})

const allCountries = computed(() => {
  if (props.region || route.params.region) {
    // 如果指定了区域，只显示该区域的国家
    const regionCode = props.region || route.params.region as string
    return productsStore.countries.filter(country => 
      country.region.toLowerCase() === regionCode.toLowerCase()
    )
  }
  return productsStore.countries
})

const filteredCountries = computed(() => {
  if (!searchQuery.value) {
    return allCountries.value
  }

  const query = searchQuery.value.toLowerCase()
  return allCountries.value.filter(country =>
    country.name.toLowerCase().includes(query) ||
    country.code.toLowerCase().includes(query)
  )
})

const groupedCountries = computed(() => {
  const groups: Array<{ letter: string; countries: Country[] }> = []
  const countryMap = new Map<string, Country[]>()

  // 按首字母分组
  filteredCountries.value.forEach(country => {
    const firstLetter = country.name.charAt(0).toUpperCase()
    if (!countryMap.has(firstLetter)) {
      countryMap.set(firstLetter, [])
    }
    countryMap.get(firstLetter)!.push(country)
  })

  // 转换为数组并排序
  Array.from(countryMap.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .forEach(([letter, countries]) => {
      groups.push({
        letter,
        countries: countries.sort((a, b) => a.name.localeCompare(b.name))
      })
    })

  return groups
})

const alphabetIndex = computed(() => {
  // 生成字母索引，只显示有国家的字母
  const letters = new Set<string>()
  filteredCountries.value.forEach(country => {
    letters.add(country.name.charAt(0).toUpperCase())
  })
  return Array.from(letters).sort()
})

// 方法
const loadCountries = async () => {
  isLoading.value = true
  error.value = null

  try {
    const regionCode = props.region || route.params.region as string
    
    if (regionCode) {
      // 加载指定区域的国家
      await productsStore.fetchCountries(regionCode)
    } else {
      // 加载所有国家
      await productsStore.fetchCountries()
    }
    
    console.log('[CountryPage] 国家数据加载成功')
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : '加载国家数据失败'
    error.value = errorMessage
    console.error('[CountryPage] 加载国家数据失败:', err)
  } finally {
    isLoading.value = false
  }
}

const handleRetry = () => {
  loadCountries()
}

const handleSearch = async (value: string) => {
  searchQuery.value = value
  
  // 如果有搜索内容，进行远程搜索
  if (value && value.length >= 2) {
    try {
      const results = await productsStore.searchCountries(value)
      console.log('[CountryPage] 搜索结果:', results.length)
    } catch (err) {
      console.warn('[CountryPage] 搜索失败:', err)
    }
  }
}

const handleSearchClear = () => {
  searchQuery.value = ''
}

const clearSearch = () => {
  searchQuery.value = ''
}

const handleCountrySelect = async (country: Country) => {
  try {
    // 设置选中的国家
    await productsStore.setSelectedCountry(country)
    
    // 触觉反馈
    telegramService.selectionFeedback()
    
    // 导航到产品列表页面
    router.push({ 
      name: 'Products',
      query: { country: country.code }
    })
    
    console.log('[CountryPage] 选择国家:', country.name)
  } catch (err) {
    console.error('[CountryPage] 选择国家失败:', err)
    
    appStore.showNotification({
      type: 'error',
      message: '选择国家失败，请重试',
      duration: 3000
    })
  }
}

const scrollToLetter = (letter: string) => {
  const element = document.querySelector(`[data-letter="${letter}"]`)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeLetter.value = letter
    telegramService.selectionFeedback()
  }
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
  telegramService.impactFeedback('light')
}

// 监听滚动事件
const handleScroll = () => {
  showScrollTop.value = window.scrollY > 300
  
  // 更新当前激活的字母
  const groups = document.querySelectorAll('.country-group')
  let currentLetter = 'A'
  
  groups.forEach(group => {
    const rect = group.getBoundingClientRect()
    if (rect.top <= 100) {
      const letter = group.querySelector('.group-letter')?.textContent
      if (letter) {
        currentLetter = letter
      }
    }
  })
  
  activeLetter.value = currentLetter
}

// 监听路由参数变化
watch(
  () => route.params.region,
  (newRegion) => {
    if (newRegion !== props.region) {
      loadCountries()
    }
  }
)

// 生命周期
onMounted(async () => {
  console.log('[CountryPage] 组件挂载')
  
  // 设置页面标题
  const title = selectedRegion.value 
    ? `${selectedRegion.value.name}国家` 
    : '选择国家'
  appStore.setCurrentPage('Countries', title)
  
  // 加载国家数据
  await loadCountries()
  
  // 添加滚动监听
  window.addEventListener('scroll', handleScroll)
})

// 清理
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="scss">
.country-page {
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
  
  .search-section {
    margin-bottom: 24px;
    
    .search-field {
      .v-field__input {
        font-size: 0.875rem;
      }
    }
  }
  
  .content-container {
    display: flex;
    gap: 16px;
    position: relative;
    
    .countries-list {
      flex: 1;
      min-width: 0;
      
      .search-results {
        margin-bottom: 16px;
        
        .results-text {
          font-size: 0.875rem;
          color: rgba(var(--v-theme-on-surface), 0.7);
          margin: 0;
        }
      }
      
      .countries-content {
        .country-group {
          margin-bottom: 24px;
          
          .group-header {
            position: sticky;
            top: 0;
            background: rgb(var(--v-theme-background));
            z-index: 10;
            padding: 8px 0;
            margin-bottom: 8px;
            
            .group-letter {
              font-size: 1.125rem;
              font-weight: 600;
              color: rgb(var(--v-theme-primary));
              margin: 0;
            }
          }
          
          .group-countries {
            display: flex;
            flex-direction: column;
            gap: 1px;
          }
        }
      }
      
      .no-results,
      .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 48px 24px;
        text-align: center;
        
        .no-results-title,
        .empty-title {
          margin: 16px 0 8px;
          color: rgba(var(--v-theme-on-surface), 0.8);
          font-size: 1.125rem;
          font-weight: 600;
        }
        
        .no-results-subtitle,
        .empty-subtitle {
          margin: 0;
          color: rgba(var(--v-theme-on-surface), 0.6);
          font-size: 0.875rem;
          line-height: 1.5;
        }
      }
    }
    
    .alphabet-index {
      position: sticky;
      top: 100px;
      height: fit-content;
      
      .index-container {
        display: flex;
        flex-direction: column;
        gap: 2px;
        padding: 8px 4px;
        background: rgba(var(--v-theme-surface), 0.8);
        border-radius: 12px;
        backdrop-filter: blur(8px);
        
        .index-letter {
          width: 24px;
          height: 24px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 0.75rem;
          font-weight: 600;
          color: rgba(var(--v-theme-on-surface), 0.6);
          background: none;
          border: none;
          border-radius: 4px;
          cursor: pointer;
          transition: all 0.2s ease;
          
          &:hover {
            background: rgba(var(--v-theme-primary), 0.1);
            color: rgb(var(--v-theme-primary));
          }
          
          &.active {
            background: rgb(var(--v-theme-primary));
            color: white;
          }
        }
      }
    }
  }
  
  .scroll-top-fab {
    bottom: 80px !important;
    left: 16px !important;
  }
}

// 响应式适配
@media (max-width: 360px) {
  .country-page {
    .page-header {
      margin-bottom: 20px;
      
      .page-title {
        font-size: 1.25rem;
      }
      
      .page-subtitle {
        font-size: 0.8125rem;
      }
    }
    
    .search-section {
      margin-bottom: 20px;
    }
    
    .content-container {
      gap: 12px;
      
      .countries-list {
        .countries-content {
          .country-group {
            margin-bottom: 20px;
            
            .group-header {
              .group-letter {
                font-size: 1rem;
              }
            }
          }
        }
        
        .no-results,
        .empty-state {
          padding: 32px 16px;
        }
      }
      
      .alphabet-index {
        .index-container {
          padding: 6px 3px;
          
          .index-letter {
            width: 20px;
            height: 20px;
            font-size: 0.7rem;
          }
        }
      }
    }
    
    .scroll-top-fab {
      bottom: 76px !important;
      left: 12px !important;
    }
  }
}

@media (min-width: 481px) {
  .country-page {
    .page-header {
      margin-bottom: 32px;
      
      .page-title {
        font-size: 1.75rem;
      }
      
      .page-subtitle {
        font-size: 0.9375rem;
      }
    }
    
    .search-section {
      margin-bottom: 32px;
    }
    
    .content-container {
      gap: 20px;
      
      .countries-list {
        .countries-content {
          .country-group {
            margin-bottom: 32px;
            
            .group-header {
              .group-letter {
                font-size: 1.25rem;
              }
            }
          }
        }
        
        .no-results,
        .empty-state {
          padding: 64px 32px;
        }
      }
      
      .alphabet-index {
        .index-container {
          padding: 10px 6px;
          
          .index-letter {
            width: 28px;
            height: 28px;
            font-size: 0.8125rem;
          }
        }
      }
    }
  }
}

// 深色主题适配
.v-theme--dark {
  .country-page {
    .content-container {
      .countries-list {
        .countries-content {
          .country-group {
            .group-header {
              background: rgb(var(--v-theme-background));
            }
          }
        }
      }
      
      .alphabet-index {
        .index-container {
          background: rgba(var(--v-theme-surface), 0.9);
        }
      }
    }
  }
}
</style>