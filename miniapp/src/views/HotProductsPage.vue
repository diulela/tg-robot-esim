<template>
  <div class="hot-products-page">
    <!-- 国家/地区列表 -->
    <CountryListComponent
      :countries="countries"
      :show-group-headers="false"
      @country-click="handleCountryClick"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
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
const countries = ref<Country[]>([])

// 初始化国家数据
onMounted(() => {
  countries.value = HOT_ITEMS_DATA.hot as Country[]
})

// 方法 - 处理国家点击
const handleCountryClick = (country: Country) => {
  router.push({
    name: 'ProductListSecondary',
    params: { countryCode: country.code },
    query: { 
      name: country.name,
      category: 'hot'
    }
  })
}
</script>

<style scoped>
.hot-products-page {
  flex: 1;
  overflow-y: auto;
  background: #f5f5f5;
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

.hot-products-page {
  animation: fadeIn 0.3s ease-out;
}
</style>
