<template>
  <v-dialog v-model="dialog" max-width="500" persistent>
    <v-card>
      <v-card-title class="dialog-title">
        <v-icon start>mdi-history</v-icon>
        eSIM 套餐历史
        <v-spacer />
        <v-btn icon variant="text" @click="closeDialog">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-card-text v-if="!esimStore.isLoadingHistory && esimStore.hasPackageHistory">
        <div class="history-list">
          <div
            v-for="(item, index) in esimStore.packageHistory"
            :key="item.id"
            class="history-item"
          >
            <!-- 套餐标签 -->
            <div class="package-header">
              <v-chip
                :color="esimStore.getStatusColor(item.status)"
                variant="elevated"
                size="small"
                class="package-chip"
              >
                <v-icon start size="14">mdi-package-variant</v-icon>
                套餐 {{ index + 1 }}
              </v-chip>
              <span class="package-date">{{ formatDate(item.createdAt) }}</span>
            </div>

            <!-- 套餐详情 -->
            <div class="package-details">
              <div class="detail-row">
                <span class="detail-label">套餐名称</span>
                <span class="detail-value">{{ item.packageName }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">流量大小</span>
                <span class="detail-value">{{ item.dataSize }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">有效期</span>
                <span class="detail-value">{{ item.validDays }}天</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">价格</span>
                <span class="detail-value">${{ item.price }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">状态</span>
                <span class="detail-value">{{ esimStore.getStatusText(item.status) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">剩余流量</span>
                <span class="detail-value">{{ item.remainingData }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">激活时间</span>
                <span class="detail-value">{{ esimStore.formatDateTime(item.activationTime) }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">到期时间</span>
                <span class="detail-value">{{ esimStore.formatDateTime(item.expireTime) }}</span>
              </div>
            </div>
          </div>
        </div>
      </v-card-text>

      <v-card-text v-else-if="esimStore.isLoadingHistory">
        <div class="loading-container">
          <v-progress-circular indeterminate color="primary" />
          <p>正在获取套餐历史...</p>
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
          <p>暂无套餐历史记录</p>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="outlined" @click="closeDialog">
          关闭
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useESIMStore } from '@/stores/esim'

// Props 定义
interface Props {
  modelValue: boolean
  orderId?: number
}

const props = defineProps<Props>()

// Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

// 组合式 API
const esimStore = useESIMStore()

// 响应式状态
const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 监听弹窗打开，自动获取数据
watch(dialog, async (newValue) => {
  if (newValue && props.orderId) {
    try {
      await esimStore.fetchHistory(props.orderId)
    } catch (error) {
      console.error('获取套餐历史失败:', error)
    }
  }
})

// 方法
const closeDialog = () => {
  dialog.value = false
  esimStore.clearError()
}

const retryFetch = async () => {
  if (!props.orderId) return
  
  try {
    esimStore.clearError()
    await esimStore.fetchHistory(props.orderId, true)
  } catch (error) {
    console.error('重试获取套餐历史失败:', error)
  }
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
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

.history-list {
  max-height: 400px;
  overflow-y: auto;
}

.history-item {
  margin-bottom: 20px;
  padding: 16px;
  border: 1px solid rgba(var(--v-theme-outline), 0.2);
  border-radius: 8px;
  background: rgba(var(--v-theme-surface-variant), 0.3);

  &:last-child {
    margin-bottom: 0;
  }
}

.package-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;

  .package-chip {
    font-size: 0.75rem;
    height: 28px;
  }

  .package-date {
    font-size: 0.75rem;
    color: rgba(var(--v-theme-on-surface), 0.6);
  }
}

.package-details {
  .detail-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 0;
    border-bottom: 1px solid rgba(var(--v-theme-outline), 0.1);

    &:last-child {
      border-bottom: none;
    }

    .detail-label {
      font-size: 0.875rem;
      font-weight: 500;
      color: rgba(var(--v-theme-on-surface), 0.7);
      min-width: 80px;
    }

    .detail-value {
      font-size: 0.875rem;
      font-weight: 500;
      color: rgb(var(--v-theme-on-surface));
      text-align: right;
      flex: 1;
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

// 响应式适配
@media (max-width: 360px) {
  .history-item {
    padding: 12px;
  }

  .package-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .detail-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;

    .detail-value {
      text-align: left;
    }
  }
}
</style>