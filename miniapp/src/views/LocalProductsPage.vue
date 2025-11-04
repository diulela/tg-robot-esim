<template>
  <div class="local-products-page">
    <!-- 搜索框 -->
    <SearchBar
      v-model="searchQuery"
      placeholder="搜索国家或地区"
      @input="handleSearch"
    />

    <!-- 国家/地区列表 -->
    <CountryListComponent
      :countries="filteredCountries"
      :grouped="groupedCountries"
      :show-group-headers="true"
      @country-click="handleCountryClick"
    />

    <!-- 字母索引 -->
    <AlphabetIndex
      :available-letters="availableLetters"
      :current-letter="currentLetter"
      @letter-click="scrollToLetter"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import SearchBar from '@/components/SearchBar.vue'
import AlphabetIndex from '@/components/AlphabetIndex.vue'
import CountryListComponent from '@/components/CountryListComponent.vue'
import { HOT_ITEMS_DATA } from '@/stores/products'

// 类型定义
interface Country {
  code: string
  name: string
  en?: string
}

// 路由
const router = useRouter()

// 状态
const searchQuery = ref('')
const currentLetter = ref('')
const countries = ref<Country[]>([])

// 初始化国家数据
onMounted(() => {
  countries.value = HOT_ITEMS_DATA.local as Country[]
})

// 计算属性 - 过滤后的国家列表
const filteredCountries = computed(() => {
  if (!searchQuery.value.trim()) {
    return countries.value
  }

  const query = searchQuery.value.toLowerCase()
  return countries.value.filter(country => 
    country.name.toLowerCase().includes(query) ||
    (country.en && country.en.toLowerCase().includes(query))
  )
})

// 计算属性 - 按字母分组的国家列表
const groupedCountries = computed(() => {
  const grouped: Record<string, Country[]> = {}

  filteredCountries.value.forEach(country => {
    // 使用英文名的首字母进行分组
    const firstLetter = (country.en || country.code).charAt(0).toUpperCase()
    
    if (!grouped[firstLetter]) {
      grouped[firstLetter] = []
    }
    
    grouped[firstLetter].push(country)
  })

  // 对每个分组内的国家按英文名排序
  Object.keys(grouped).forEach(letter => {
    const countries = grouped[letter]
    if (countries) {
      countries.sort((a, b) => {
        const nameA = a.en || a.name
        const nameB = b.en || b.name
        return nameA.localeCompare(nameB)
      })
    }
  })

  return grouped
})

// 计算属性 - 可用的字母列表
const availableLetters = computed(() => {
  return Object.keys(groupedCountries.value).sort()
})

// 方法 - 处理搜索
const handleSearch = (query: string) => {
  searchQuery.value = query
  // 搜索时重置当前字母
  currentLetter.value = ''
}

// 方法 - 处理国家点击
const handleCountryClick = (country: Country) => {
  router.push({
    name: 'ProductListSecondary',
    params: { countryCode: country.code },
    query: { 
      name: country.name,
      category: 'local'
    }
  })
}

// 方法 - 滚动到指定字母
const scrollToLetter = (letter: string) => {
  currentLetter.value = letter
  
  const element = document.getElementById(`letter-${letter}`)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}
</script>

<style scoped>
.local-products-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f5f5f5;
  position: relative;
}

/* 页面进入动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.local-products-page {
  animation: fadeIn 0.3s ease-out;
}
</style>
