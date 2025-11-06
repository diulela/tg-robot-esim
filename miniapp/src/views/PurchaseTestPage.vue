<template>
  <div class="purchase-test-page">
    <v-container>
      <v-row>
        <v-col cols="12">
          <h1>购买弹窗测试页面</h1>
          <p>这是一个测试页面，用于验证购买弹窗功能。</p>
        </v-col>
      </v-row>
      
      <v-row>
        <v-col cols="12" md="6">
          <ProductCard
            :product="testProduct"
            @click="handleProductClick"
            @buy="handleProductBuy"
          />
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ProductCard from '@/components/ProductCard.vue'
import type { Product } from '@/types'

// 测试商品数据
const testProduct = ref<Product>({
  id: '1',
  name: '中国 无限流量/15天',
  description: '中国大陆地区15天无限流量eSIM卡，支持4G/5G网络',
  price: 27.72,
  originalPrice: 35.00,
  currency: 'USD',
  region: 'local',
  country: '中国',
  countryCode: 'CN',
  dataAmount: '无限',
  validDays: 15,
  coverage: '中国大陆',
  coverageAreas: ['中国大陆'],
  features: ['4G/5G网络', '无限流量', '即时激活'],
  isActive: true,
  isPopular: true,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString()
})

// 事件处理
const handleProductClick = (product: Product) => {
  console.log('产品点击:', product)
}

const handleProductBuy = (product: Product) => {
  console.log('产品购买:', product)
}
</script>

<style scoped lang="scss">
.purchase-test-page {
  padding: 20px;
  
  h1 {
    color: rgb(var(--v-theme-primary));
    margin-bottom: 16px;
  }
  
  p {
    color: rgba(var(--v-theme-on-surface), 0.7);
    margin-bottom: 24px;
  }
}
</style>