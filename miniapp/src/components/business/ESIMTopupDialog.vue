<template>
  <v-dialog v-model="dialog" max-width="450" persistent>
    <v-card>
      <v-card-title class="dialog-title">
        <v-icon start>mdi-plus-circle</v-icon>
        流量充值
        <v-spacer />
        <v-btn icon variant="text" @click="closeDialog">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-card-text v-if="!esimStore.isLoadingTopup && esimStore.hasTopupPackages">
        <div class="topup-packages">
          <div
            v-for="pkg in esimStore.topupPackages"
            :key="pkg.id"
            class="package-option"
            :class="{ selected: selectedPackage?.id === pkg.id }"
            @click="selectPackage(pkg)"
          >
            <div class="package-info">
              <h4 class="package-title">{{ pkg.title }}</h4>
              <p class="package-description">{{ pkg.description }}</p>
              <div class="package-specs">
                <span class="spec-item">
                  <v-icon start size="14">mdi-database</v-icon>
                  {{ pkg.data }}
                </span>
                <span class="spec-item">
                  <v-icon start size="14">mdi-calendar</v-icon>
                  {{ pkg.validity }}天
                </span>
              </div>
            </div>
            <div class="package-price">
              <span class="price-amount">${{ pkg.price }}</span>
            </div>
            <div class="package-selection">
              <v-icon 
                :color="selectedPackage?.id === pkg.id ? 'primary' : 'grey'"
                size="20"
              >
                {{ selectedPackage?.id === pkg.id ? 'mdi-check-circle' : 'mdi-circle-outline' }}
              </v-icon>
            </div>
          </div>
        </div>

        <!-- 选中套餐信息 -->
        <div v-if="selectedPackage" class="selected-info">
          <v-divider class="my-4" />
          <div class="selected-summary">
            <h4>充值详情</h4>
            <div class="summary-item">
              <span>套餐名称：</span>
              <span>{{ selectedPackage.title }}</span>
            </div>
            <div class="summary-item">
              <span>流量大小：</span>
              <span>{{ selectedPackage.data }}</span>
            </div>
            <div class="summary-item">
              <span>有效期：</span>
              <span>{{ selectedPackage.validity }}天</span>
            </div>
            <div class="summary-item total">
              <span>充值金额：</span>
              <span class="price">${{ selectedPackage.price }}</span>
            </div>
          </div>
        </div>
      </v-card-text>

      <v-card-text v-else-if="esimStore.isLoadingTopup">
        <div class="loading-container">
          <v-progress-circular indeterminate color="primary" />
          <p>正在获取充值套餐...</p>
        </div>
      </v-card-text>

      <v-card-text v-else-if="esimStore.error">
        <div class="error-container">
          <v-icon size="48" color="error">mdi-alert-circle</v-icon>
          <p class="error-message">{{ esimStore.error }}</p>
          <v-btn color="primary" variant="outlined" @click="retryFetch">
            重试
          </v-btn>
        </div>
      </v-card-text>

      <v-card-text v-else>
        <div class="empty-state">
          <v-icon size="48" color="grey">mdi-package-variant-closed</v-icon>
          <p>暂无可用充值套餐</p>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="outlined" @click="closeDialog" :disabled="esimStore.isTopupping">
          取消
        </v-btn>
        <v-btn
          color="primary"
          variant="elevated"
          :disabled="!selectedPackage"
          :loading="esimStore.isTopupping"
          @click="confirmTopup"
        >
          确认充值
        </v-btn>
      </v-card-actions>
    </v-card>

    <!-- 充值确认对话框 -->
    <v-dialog v-model="showConfirmDialog" max-width="350" persistent>
      <v-card>
        <v-card-title class="confirm-title">
          <v-icon start color="warning">mdi-alert</v-icon>
          确认充值
        </v-card-title>
        
        <v-card-text>
          <p>您确定要为此 eSIM 充值以下套餐吗？</p>
          <div class="confirm-details" v-if="selectedPackage">
            <div class="confirm-item">
              <strong>{{ selectedPackage.title }}</strong>
            </div>
            <div class="confirm-item">
              流量：{{ selectedPackage.data }}
            </div>
            <div class="confirm-item">
              有效期：{{ selectedPackage.validity }}天
            </div>
            <div class="confirm-item price">
              金额：${{ selectedPackage.price }}
            </div>
          </div>
        </v-card-text>
        
        <v-card-actions>
          <v-spacer />
          <v-btn variant="outlined" @click="showConfirmDialog = false">
            取消
          </v-btn>
          <v-btn color="primary" variant="elevated" @click="executeTopup">
            确认充值
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useESIMStore } from '@/stores/esim'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'
import type { TopupPackage } from '@/types'

// Props 定义
interface Props {
  modelValue: boolean
  orderId?: number
}

const props = defineProps<Props>()

// Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'topup-success': []
}>()

// 组合式 API
const esimStore = useESIMStore()
const appStore = useAppStore()

// 响应式状态
const selectedPackage = ref<TopupPackage | null>(null)
const showConfirmDialog = ref(false)

const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 监听弹窗打开，自动获取数据
watch(dialog, async (newValue) => {
  if (newValue && props.orderId) {
    try {
      await esimStore.fetchTopupPackages(props.orderId)
    } catch (error) {
      console.error('获取充值套餐失败:', error)
    }
  } else if (!newValue) {
    // 弹窗关闭时重置状态
    selectedPackage.value = null
    showConfirmDialog.value = false
  }
})

// 方法
const closeDialog = () => {
  dialog.value = false
  esimStore.clearError()
}

const selectPackage = (pkg: TopupPackage) => {
  selectedPackage.value = pkg
  telegramService.impactFeedback('light')
}

const retryFetch = async () => {
  if (!props.orderId) return
  
  try {
    esimStore.clearError()
    await esimStore.fetchTopupPackages(props.orderId, true)
  } catch (error) {
    console.error('重试获取充值套餐失败:', error)
  }
}

const confirmTopup = () => {
  if (!selectedPackage.value) {
    appStore.showNotification({
      type: 'warning',
      message: '请选择充值套餐',
      duration: 2000
    })
    return
  }
  
  showConfirmDialog.value = true
  telegramService.impactFeedback('medium')
}

const executeTopup = async () => {
  if (!selectedPackage.value || !props.orderId) return

  try {
    showConfirmDialog.value = false
    
    await esimStore.topupEsim(props.orderId, selectedPackage.value.id)
    
    appStore.showNotification({
      type: 'success',
      message: '充值成功！流量已添加到您的 eSIM',
      duration: 3000
    })
    
    telegramService.impactFeedback('heavy')
    emit('topup-success')
    closeDialog()
    
  } catch (error) {
    appStore.showNotification({
      type: 'error',
      message: error instanceof Error ? error.message : '充值失败，请重试',
      duration: 3000
    })
  }
}
</script>

<style scoped lang="scss">
.dialog-title {
  font-size: 1.1rem;
  font-weight: 600;
  padding: 16px 16px 8px 16px;
  color: rgb(var(--v-theme-primary));

  .v-icon {
    margin-right: 8px;
  }
}

.topup-packages {
  max-height: 350px;
  overflow-y: auto;
}

.package-option {
  display: flex;
  align-items: center;
  padding: 16px;
  margin-bottom: 12px;
  border: 2px solid rgba(var(--v-theme-outline), 0.2);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: rgba(var(--v-theme-primary), 0.5);
    background: rgba(var(--v-theme-primary), 0.05);
  }

  &.selected {
    border-color: rgb(var(--v-theme-primary));
    background: rgba(var(--v-theme-primary), 0.1);
  }

  &:last-child {
    margin-bottom: 0;
  }
}

.package-info {
  flex: 1;

  .package-title {
    font-size: 1rem;
    font-weight: 600;
    margin: 0 0 4px 0;
    color: rgb(var(--v-theme-on-surface));
  }

  .package-description {
    font-size: 0.75rem;
    color: rgba(var(--v-theme-on-surface), 0.6);
    margin: 0 0 8px 0;
  }

  .package-specs {
    display: flex;
    gap: 12px;

    .spec-item {
      display: flex;
      align-items: center;
      font-size: 0.75rem;
      color: rgba(var(--v-theme-on-surface), 0.7);

      .v-icon {
        margin-right: 4px;
      }
    }
  }
}

.package-price {
  margin: 0 16px;

  .price-amount {
    font-size: 1.25rem;
    font-weight: 700;
    color: rgb(var(--v-theme-primary));
  }
}

.package-selection {
  display: flex;
  align-items: center;
}

.selected-info {
  .selected-summary {
    h4 {
      font-size: 1rem;
      font-weight: 600;
      margin-bottom: 12px;
      color: rgb(var(--v-theme-on-surface));
    }

    .summary-item {
      display: flex;
      justify-content: space-between;
      padding: 6px 0;
      font-size: 0.875rem;

      &.total {
        font-weight: 600;
        border-top: 1px solid rgba(var(--v-theme-outline), 0.2);
        margin-top: 8px;
        padding-top: 12px;

        .price {
          color: rgb(var(--v-theme-primary));
          font-size: 1rem;
        }
      }
    }
  }
}

.loading-container,
.error-container,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 16px;
  text-align: center;

  p {
    margin: 16px 0;
    color: rgba(var(--v-theme-on-surface), 0.7);
  }

  .error-message {
    color: rgb(var(--v-theme-error));
  }
}

.confirm-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: rgb(var(--v-theme-warning));
}

.confirm-details {
  margin-top: 16px;
  padding: 12px;
  background: rgba(var(--v-theme-surface-variant), 0.5);
  border-radius: 8px;

  .confirm-item {
    padding: 4px 0;
    font-size: 0.875rem;

    &.price {
      font-weight: 600;
      color: rgb(var(--v-theme-primary));
      font-size: 1rem;
      margin-top: 8px;
      padding-top: 8px;
      border-top: 1px solid rgba(var(--v-theme-outline), 0.2);
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .package-option {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;

    .package-price {
      margin: 0;
      align-self: flex-end;
    }

    .package-selection {
      align-self: center;
    }
  }

  .package-specs {
    flex-direction: column;
    gap: 4px;
  }
}
</style>