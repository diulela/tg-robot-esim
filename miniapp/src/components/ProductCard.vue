<template>
  <v-card
    class="product-card"
    :class="{ compact }"
    elevation="2"
    @click="handleClick"
  >
    <div class="product-content">
      <!-- 商品图标 -->
      <div class="product-icon">
        <v-img
          v-if="product.icon"
          :src="product.icon"
          :alt="product.name"
          width="80"
          height="80"
          cover
        />
        <v-icon v-else size="80" color="primary">mdi-sim</v-icon>
      </div>

      <!-- 商品信息 -->
      <div class="product-info">
        <h3 class="product-name">{{ product.name }}</h3>
        <p class="product-description">{{ product.description }}</p>

        <!-- 商品规格 -->
        <div class="product-specs">
          <span class="spec-item">
            <v-icon size="16" class="spec-icon">mdi-database</v-icon>
            {{ product.dataAmount }}
          </span>
          <span class="spec-item">
            <v-icon size="16" class="spec-icon">mdi-clock-outline</v-icon>
            {{ formatValidDays(product.validDays) }}
          </span>
          <span class="spec-item">
            <v-icon size="16" class="spec-icon">mdi-earth</v-icon>
            {{ product.coverage || product.country }}
          </span>
        </div>

        <!-- 商品特性 -->
        <div v-if="product.features && product.features.length > 0" class="product-features">
          <v-chip
            v-for="feature in product.features.slice(0, 3)"
            :key="feature"
            size="small"
            variant="tonal"
            color="success"
            class="feature-chip"
          >
            <v-icon start size="14">mdi-check</v-icon>
            {{ feature }}
          </v-chip>
        </div>

        <!-- 价格和购买 -->
        <div class="product-footer">
          <div class="price-section">
            <span class="current-price">${{ product.price.toFixed(2) }}</span>
            <span v-if="product.originalPrice && product.originalPrice > product.price" class="original-price">
              ${{ product.originalPrice.toFixed(2) }}
            </span>
          </div>
          <v-btn
            color="primary"
            variant="elevated"
            class="buy-button"
            @click.stop="handleBuy"
          >
            <v-icon start>mdi-cart</v-icon>
            立即购买
          </v-btn>
        </div>
      </div>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import type { ProductCardProps, ProductCardEmits } from '@/types'

// Props
const props = withDefaults(defineProps<ProductCardProps>(), {
  compact: false
})

// Emits
const emit = defineEmits<ProductCardEmits>()

// 方法
const formatValidDays = (days: number): string => {
  if (days === 1) return '1天'
  if (days === 7) return '1周'
  if (days === 30) return '1个月'
  if (days === 365) return '1年'
  return `${days}天`
}

const handleClick = () => {
  emit('click', props.product)
}

const handleBuy = () => {
  emit('buy', props.product)
}
</script>

<style scoped>
.product-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
}

.product-card:active {
  transform: translateY(-2px);
}

.product-content {
  display: flex;
  gap: 16px;
  padding: 16px;
}

.product-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

.product-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-description {
  font-size: 13px;
  color: #666;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

.product-specs {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 4px;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #666;
}

.spec-icon {
  color: #667eea;
}

.product-features {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.feature-chip {
  font-size: 11px;
  height: 24px;
}

.product-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 8px;
}

.price-section {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.current-price {
  font-size: 20px;
  font-weight: 700;
  color: #667eea;
}

.original-price {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
}

.buy-button {
  flex-shrink: 0;
  transition: transform 0.2s ease;
}

.buy-button:hover {
  transform: scale(1.05);
}

.buy-button:active {
  transform: scale(0.98);
}

/* Compact 模式 */
.product-card.compact .product-content {
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.product-card.compact .product-icon {
  width: 60px;
  height: 60px;
}

.product-card.compact .product-name {
  font-size: 14px;
}

.product-card.compact .product-description {
  font-size: 12px;
}

/* 响应式设计 */
@media (max-width: 600px) {
  .product-content {
    flex-direction: column;
  }

  .product-icon {
    width: 100%;
    height: 120px;
  }

  .product-footer {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .buy-button {
    width: 100%;
  }
}
</style>
