<template>
  <v-card class="esim-info-card" variant="elevated" id="esim-info">
    <v-card-title class="card-title">
      <v-icon start>mdi-sim</v-icon>
      eSIM信息 (共1张卡)
    </v-card-title>

    <v-card-text class="card-content">
      <!-- eSIM 状态标签 -->
      <div class="esim-status">
        <v-chip
          color="primary"
          variant="elevated"
          size="small"
          class="status-chip"
        >
          <v-icon start size="14">mdi-sim</v-icon>
          eSIM 1
        </v-chip>
        
        <v-chip
          color="info"
          variant="tonal"
          size="small"
          class="status-chip"
          @click="openUsageDialog"
        >
          <v-icon start size="14">mdi-chart-donut</v-icon>
          查看用量
        </v-chip>
        
        <v-chip
          color="success"
          variant="tonal"
          size="small"
          class="status-chip"
          @click="openHistoryDialog"
        >
          <v-icon start size="14">mdi-history</v-icon>
          套餐历史
        </v-chip>
        
        <v-chip
          v-if="esimInfo.status && esimInfo.usagePercentage && parseFloat(esimInfo.usagePercentage) > 80"
          color="warning"
          variant="tonal"
          size="small"
          class="status-chip"
        >
          <v-icon start size="14">mdi-alert</v-icon>
          流量不足
        </v-chip>
        
        <v-chip
          color="orange"
          variant="elevated"
          size="small"
          class="status-chip"
          @click="openTopupDialog"
        >
          <v-icon start size="14">mdi-plus-circle</v-icon>
          充值
        </v-chip>
      </div>

      <!-- eSIM 详细信息 -->
      <div class="esim-details">
        <!-- ICCID -->
        <div class="detail-item">
          <span class="detail-label">ICCID</span>
          <div class="detail-value-container">
            <span class="detail-value">{{ esimInfo.iccid }}</span>
            <v-btn
              icon
              variant="text"
              size="small"
              @click="copyIccid"
              class="copy-btn"
            >
              <v-icon size="16">mdi-content-copy</v-icon>
            </v-btn>
          </div>
        </div>

        <!-- 激活码 -->
        <div class="detail-item">
          <span class="detail-label">激活码</span>
          <div class="detail-value-container">
            <span class="detail-value">{{ esimInfo.activationCode }}</span>
            <v-btn
              icon
              variant="text"
              size="small"
              @click="copyActivationCode"
              class="copy-btn"
            >
              <v-icon size="16">mdi-content-copy</v-icon>
            </v-btn>
          </div>
        </div>

        <!-- LPA地址 -->
        <div v-if="esimInfo.lpaAddress" class="detail-item">
          <span class="detail-label">LPA地址</span>
          <div class="detail-value-container">
            <span class="detail-value">{{ esimInfo.lpaAddress }}</span>
            <v-btn
              icon
              variant="text"
              size="small"
              @click="copyLpaAddress"
              class="copy-btn"
            >
              <v-icon size="16">mdi-content-copy</v-icon>
            </v-btn>
          </div>
        </div>

        <!-- APN类型 -->
        <div class="detail-item">
          <span class="detail-label">APN类型</span>
          <span class="detail-value">{{ esimInfo.apnType === 'manual' ? '手动' : '自动' }}</span>
        </div>

        <!-- APN值 -->
        <div v-if="esimInfo.apnValue" class="detail-item">
          <span class="detail-label">APN值</span>
          <span class="detail-value">{{ esimInfo.apnValue }}</span>
        </div>

        <!-- 是否漫游 -->
        <div class="detail-item">
          <span class="detail-label">是否漫游</span>
          <span class="detail-value">{{ esimInfo.isRoaming ? '是' : '否' }}</span>
        </div>

        <!-- eSIM状态 -->
        <div v-if="esimInfo.status" class="detail-item">
          <span class="detail-label">状态</span>
          <v-chip 
            :color="getStatusColor(esimInfo.status)" 
            size="small" 
            variant="tonal"
          >
            {{ getStatusText(esimInfo.status) }}
          </v-chip>
        </div>

        <!-- 使用进度 -->
        <div v-if="esimInfo.usagePercentage" class="detail-item usage-item">
          <span class="detail-label">使用进度</span>
          <div class="usage-progress">
            <v-progress-linear
              :model-value="parseFloat(esimInfo.usagePercentage)"
              height="6"
              color="primary"
              bg-color="grey-lighten-3"
              rounded
            />
            <span class="usage-text">{{ esimInfo.usagePercentage }}%</span>
          </div>
        </div>
      </div>

      <!-- 二维码区域 -->
      <div v-if="showQRCode && esimInfo.qrCode" class="qr-code-section">
        <div class="qr-code-container">
          <img
            :src="esimInfo.qrCode"
            alt="eSIM QR Code"
            class="qr-code-image"
          />
        </div>
        
        <div class="qr-actions">
          <v-btn
            variant="outlined"
            color="primary"
            size="small"
            @click="downloadQR"
            class="qr-btn"
          >
            <v-icon start size="16">mdi-download</v-icon>
            下载二维码
          </v-btn>
          
          <v-btn
            variant="outlined"
            color="secondary"
            size="small"
            @click="shareQR"
            class="qr-btn"
          >
            <v-icon start size="16">mdi-share</v-icon>
            分享二维码
          </v-btn>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="action-section">
        <div class="action-buttons">
          <v-btn
            color="primary"
            variant="elevated"
            @click="exportPDF"
            :loading="isExporting"
            class="action-btn"
          >
            <v-icon start>mdi-file-pdf-box</v-icon>
            导出PDF
          </v-btn>
          
          <v-btn
            color="secondary"
            variant="outlined"
            @click="openInstallGuide"
            class="action-btn"
          >
            <v-icon start>mdi-cellphone-cog</v-icon>
            安装说明
          </v-btn>
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAppStore } from '@/stores/app'
import { useESIMStore } from '@/stores/esim'
import { telegramService } from '@/services/telegram'
import type { ESIMInfo, ESIMStatus } from '@/types'

// Props 定义
interface Props {
  esimInfo: ESIMInfo
  orderId?: number
  showQRCode?: boolean
  allowCopy?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showQRCode: true,
  allowCopy: true
})

// Emits 定义
const emit = defineEmits<{
  'copy-iccid': [iccid: string]
  'copy-activation-code': [code: string]
  'download-qr': [qrCode: string]
  'share-qr': [qrCode: string]
  'open-usage-dialog': []
  'open-history-dialog': []
  'open-topup-dialog': []
  'open-install-guide': []
}>()

// 组合式 API
const appStore = useAppStore()
const esimStore = useESIMStore()

// 响应式状态
const isExporting = ref(false)

// 方法
const copyIccid = async () => {
  if (!props.allowCopy) return

  try {
    await navigator.clipboard.writeText(props.esimInfo.iccid)
    
    appStore.showNotification({
      type: 'success',
      message: 'ICCID 已复制到剪贴板',
      duration: 2000
    })
    
    telegramService.impactFeedback('light')
    emit('copy-iccid', props.esimInfo.iccid)
  } catch (err) {
    appStore.showNotification({
      type: 'error',
      message: '复制失败',
      duration: 2000
    })
  }
}

const copyActivationCode = async () => {
  if (!props.allowCopy) return

  try {
    await navigator.clipboard.writeText(props.esimInfo.activationCode)
    
    appStore.showNotification({
      type: 'success',
      message: '激活码已复制到剪贴板',
      duration: 2000
    })
    
    telegramService.impactFeedback('light')
    emit('copy-activation-code', props.esimInfo.activationCode)
  } catch (err) {
    appStore.showNotification({
      type: 'error',
      message: '复制失败',
      duration: 2000
    })
  }
}

const downloadQR = () => {
  telegramService.impactFeedback('medium')
  emit('download-qr', props.esimInfo.qrCode)
  
  // 创建下载链接
  const link = document.createElement('a')
  link.href = props.esimInfo.qrCode
  link.download = `esim-qr-${props.esimInfo.iccid}.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  appStore.showNotification({
    type: 'success',
    message: '二维码下载成功',
    duration: 2000
  })
}

const shareQR = () => {
  telegramService.impactFeedback('medium')
  emit('share-qr', props.esimInfo.qrCode)
  
  if (navigator.share) {
    navigator.share({
      title: 'eSIM 二维码',
      text: `ICCID: ${props.esimInfo.iccid}`,
      url: props.esimInfo.qrCode
    }).catch(err => {
      console.error('分享失败:', err)
    })
  } else {
    // 降级处理：复制到剪贴板
    copyActivationCode()
  }
}

const copyLpaAddress = async () => {
  if (!props.allowCopy || !props.esimInfo.lpaAddress) return

  try {
    await navigator.clipboard.writeText(props.esimInfo.lpaAddress)
    
    appStore.showNotification({
      type: 'success',
      message: 'LPA地址已复制到剪贴板',
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

const openUsageDialog = () => {
  emit('open-usage-dialog')
  telegramService.impactFeedback('light')
}

const openHistoryDialog = () => {
  emit('open-history-dialog')
  telegramService.impactFeedback('light')
}

const openTopupDialog = () => {
  emit('open-topup-dialog')
  telegramService.impactFeedback('light')
}

const openInstallGuide = () => {
  emit('open-install-guide')
  telegramService.impactFeedback('light')
}

const exportPDF = async () => {
  if (!props.orderId) {
    appStore.showNotification({
      type: 'error',
      message: '缺少订单信息',
      duration: 2000
    })
    return
  }

  isExporting.value = true
  
  try {
    await esimStore.exportPDF(props.orderId)
    
    appStore.showNotification({
      type: 'success',
      message: 'PDF 导出成功',
      duration: 2000
    })
    
    telegramService.impactFeedback('medium')
  } catch (error) {
    appStore.showNotification({
      type: 'error',
      message: error instanceof Error ? error.message : 'PDF 导出失败',
      duration: 3000
    })
  } finally {
    isExporting.value = false
  }
}

const getStatusColor = (status: ESIMStatus): string => {
  const colorMap: Record<ESIMStatus, string> = {
    pending: 'warning',
    active: 'success',
    expired: 'error',
    suspended: 'warning',
    terminated: 'error'
  }
  return colorMap[status] || 'default'
}

const getStatusText = (status: ESIMStatus): string => {
  const textMap: Record<ESIMStatus, string> = {
    pending: '待激活',
    active: '已激活',
    expired: '已过期',
    suspended: '已暂停',
    terminated: '已终止'
  }
  return textMap[status] || status
}
</script>

<style scoped lang="scss">
.esim-info-card {
  .card-title {
    font-size: 1.1rem;
    font-weight: 600;
    padding: 16px 16px 8px 16px;
    color: rgb(var(--v-theme-primary));

    .v-icon {
      margin-right: 8px;
    }
  }

  .card-content {
    padding: 8px 16px 16px 16px;
  }
}

.esim-status {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 20px;

  .status-chip {
    font-size: 0.75rem;
    height: 28px;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    }
  }
}

.esim-details {
  margin-bottom: 24px;

  .detail-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid rgba(var(--v-theme-outline), 0.12);

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
      font-family: 'Courier New', monospace;
      word-break: break-all;
    }

    .detail-value-container {
      display: flex;
      align-items: center;
      gap: 8px;
      flex: 1;
      justify-content: flex-end;

      .copy-btn {
        opacity: 0.7;
        transition: opacity 0.2s ease;

        &:hover {
          opacity: 1;
        }
      }
    }

    &.usage-item {
      .usage-progress {
        display: flex;
        align-items: center;
        gap: 8px;
        flex: 1;
        justify-content: flex-end;

        .v-progress-linear {
          width: 80px;
        }

        .usage-text {
          font-size: 0.75rem;
          font-weight: 600;
          color: rgb(var(--v-theme-primary));
          min-width: 35px;
        }
      }
    }
  }
}

.qr-code-section {
  text-align: center;
  margin-bottom: 24px;

  .qr-code-container {
    display: flex;
    justify-content: center;
    margin-bottom: 16px;

    .qr-code-image {
      width: 200px;
      height: 200px;
      border: 1px solid rgba(var(--v-theme-outline), 0.2);
      border-radius: 8px;
      background: white;
    }
  }

  .qr-actions {
    display: flex;
    gap: 12px;
    justify-content: center;

    .qr-btn {
      min-width: 120px;
    }
  }
}

.action-section {
  .action-buttons {
    display: flex;
    gap: 12px;

    .action-btn {
      flex: 1;
      height: 44px;
      font-weight: 600;
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .esim-info-card {
    .card-content {
      padding: 6px 12px 12px 12px;
    }
  }

  .esim-status {
    gap: 6px;

    .status-chip {
      font-size: 0.7rem;
      height: 26px;
    }
  }

  .esim-details {
    .detail-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;

      .detail-value-container {
        width: 100%;
        justify-content: space-between;
      }
    }
  }

  .qr-code-section {
    .qr-code-container {
      .qr-code-image {
        width: 160px;
        height: 160px;
      }
    }

    .qr-actions {
      flex-direction: column;

      .qr-btn {
        width: 100%;
      }
    }
  }

  .action-section {
    .action-buttons {
      flex-direction: column;

      .action-btn {
        width: 100%;
      }
    }
  }
}
</style>