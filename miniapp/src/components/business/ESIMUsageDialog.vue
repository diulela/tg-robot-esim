<template>
  <v-dialog v-model="dialog" max-width="400" persistent>
    <v-card>
      <v-card-title class="dialog-title">
        <v-icon start>mdi-chart-donut</v-icon>
        eSIM 使用情况
        <v-spacer />
        <v-btn icon variant="text" @click="closeDialog">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-card-text v-if="!esimStore.isLoadingUsage && esimStore.usageInfo">
        <!-- eSIM 信息 -->
        <div class="esim-info-section">
          <h4>eSIM信息</h4>
          <div class="info-item">
            <span class="label">ICCID</span>
            <div class="value-container">
              <span class="value">{{ esimStore.usageInfo.iccid }}</span>
              <v-btn icon variant="text" size="small" @click="copyIccid">
                <v-icon size="16">mdi-content-copy</v-icon>
              </v-btn>
            </div>
          </div>
        </div>

        <!-- 流量使用 -->
        <div class="usage-section">
          <h4>流量使用</h4>
          <div class="usage-stats">
            <div class="stat-item">
              <span class="stat-label">总流量</span>
              <span class="stat-value">{{ esimStore.formatDataSize(esimStore.usageInfo.dataTotal) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">已使用</span>
              <span class="stat-value">{{ esimStore.formatDataSize(esimStore.usageInfo.dataUsed) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">剩余流量</span>
              <span class="stat-value">{{ esimStore.formatDataSize(esimStore.usageInfo.dataRemaining) }}</span>
            </div>
          </div>

          <!-- 使用进度 -->
          <div class="usage-progress">
            <v-progress-linear
              :model-value="parseFloat(esimStore.usageInfo.usagePercentage)"
              height="8"
              color="primary"
              bg-color="grey-lighten-3"
              rounded
            />
            <div class="progress-text">
              <span>使用进度</span>
              <span>{{ esimStore.usageInfo.usagePercentage }}%</span>
            </div>
          </div>

          <div class="usage-tip">
            <v-icon start size="16" :color="getTipColor()">{{ getTipIcon() }}</v-icon>
            <span>{{ getTipText() }}</span>
          </div>
        </div>

        <!-- 有效期 -->
        <div class="validity-section">
          <div class="info-item">
            <span class="label">有效期</span>
            <span class="value">{{ esimStore.formatDateTime(esimStore.usageInfo.expireTime) }}</span>
          </div>
        </div>
      </v-card-text>

      <v-card-text v-else-if="esimStore.isLoadingUsage">
        <div class="loading-container">
          <v-progress-circular indeterminate color="primary" />
          <p>正在获取使用情况...</p>
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

      <v-card-actions v-if="!esimStore.isLoadingUsage && esimStore.usageInfo">
        <v-spacer />
        <v-btn color="primary" variant="elevated" @click="openTopupDialog">
          流量充值
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useESIMStore } from '@/stores/esim'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'

// Props 定义
interface Props {
  modelValue: boolean
  orderId?: number
}

const props = defineProps<Props>()

// Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'open-topup': []
}>()

// 组合式 API
const esimStore = useESIMStore()
const appStore = useAppStore()

// 响应式状态
const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 监听弹窗打开，自动获取数据
watch(dialog, async (newValue) => {
  if (newValue && props.orderId) {
    try {
      await esimStore.fetchUsage(props.orderId)
    } catch (error) {
      console.error('获取使用情况失败:', error)
    }
  }
})

// 方法
const closeDialog = () => {
  dialog.value = false
  esimStore.clearError()
}

const copyIccid = async () => {
  if (!esimStore.usageInfo) return

  try {
    await navigator.clipboard.writeText(esimStore.usageInfo.iccid)
    
    appStore.showNotification({
      type: 'success',
      message: 'ICCID 已复制到剪贴板',
      duration: 2000
    })
    
    telegramService.impactFeedback('light')
  } catch (err) {
    appStore.showNotification({
      type: 'error',
      message: '复制失败',
      duration: 2000
    })
  }
}

const retryFetch = async () => {
  if (!props.orderId) return
  
  try {
    esimStore.clearError()
    await esimStore.fetchUsage(props.orderId, true)
  } catch (error) {
    console.error('重试获取使用情况失败:', error)
  }
}

const openTopupDialog = () => {
  emit('open-topup')
  closeDialog()
}

// 计算属性
const getTipColor = () => {
  if (!esimStore.usageInfo) return 'info'
  
  const percentage = parseFloat(esimStore.usageInfo.usagePercentage)
  if (percentage >= 90) return 'error'
  if (percentage >= 70) return 'warning'
  return 'success'
}

const getTipIcon = () => {
  if (!esimStore.usageInfo) return 'mdi-information'
  
  const percentage = parseFloat(esimStore.usageInfo.usagePercentage)
  if (percentage >= 90) return 'mdi-alert'
  if (percentage >= 70) return 'mdi-alert-outline'
  return 'mdi-check-circle'
}

const getTipText = () => {
  if (!esimStore.usageInfo) return '流量充足，请放心使用'
  
  const percentage = parseFloat(esimStore.usageInfo.usagePercentage)
  if (percentage >= 90) return '流量即将用完，建议及时充值'
  if (percentage >= 70) return '流量使用较多，请注意剩余量'
  return '流量充足，请放心使用'
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

.esim-info-section,
.usage-section,
.validity-section {
  margin-bottom: 24px;

  h4 {
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 12px;
    color: rgb(var(--v-theme-on-surface));
  }
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;

  .label {
    font-size: 0.875rem;
    font-weight: 500;
    color: rgba(var(--v-theme-on-surface), 0.7);
  }

  .value {
    font-size: 0.875rem;
    font-weight: 500;
    color: rgb(var(--v-theme-on-surface));
    font-family: 'Courier New', monospace;
  }

  .value-container {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.usage-stats {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;

  .stat-item {
    text-align: center;

    .stat-label {
      display: block;
      font-size: 0.75rem;
      color: rgba(var(--v-theme-on-surface), 0.6);
      margin-bottom: 4px;
    }

    .stat-value {
      display: block;
      font-size: 0.875rem;
      font-weight: 600;
      color: rgb(var(--v-theme-primary));
    }
  }
}

.usage-progress {
  margin-bottom: 16px;

  .progress-text {
    display: flex;
    justify-content: space-between;
    margin-top: 8px;
    font-size: 0.75rem;
    color: rgba(var(--v-theme-on-surface), 0.7);
  }
}

.usage-tip {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: rgba(var(--v-theme-surface-variant), 0.5);
  border-radius: 8px;
  font-size: 0.75rem;

  .v-icon {
    margin-right: 8px;
  }
}

.loading-container,
.error-container {
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
  .usage-stats {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .stat-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    text-align: left;

    .stat-label,
    .stat-value {
      display: inline;
    }
  }
}
</style>