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
        >
          <v-icon start size="14">mdi-check-circle</v-icon>
          查看中继
        </v-chip>
        
        <v-chip
          color="success"
          variant="tonal"
          size="small"
          class="status-chip"
        >
          <v-icon start size="14">mdi-file-pdf-box</v-icon>
          查看PDF
        </v-chip>
        
        <v-chip
          color="warning"
          variant="tonal"
          size="small"
          class="status-chip"
        >
          <v-icon start size="14">mdi-shield-check</v-icon>
          套餐用完
        </v-chip>
        
        <v-chip
          color="orange"
          variant="elevated"
          size="small"
          class="status-chip"
        >
          <v-icon start size="14">mdi-history</v-icon>
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

        <!-- APN类型 -->
        <div class="detail-item">
          <span class="detail-label">APN类型</span>
          <span class="detail-value">{{ esimInfo.apnType }}</span>
        </div>

        <!-- 是否漫游 -->
        <div class="detail-item">
          <span class="detail-label">是否漫游</span>
          <span class="detail-value">{{ esimInfo.isRoaming ? '是' : '否' }}</span>
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

      <!-- 使用说明按钮 -->
      <div class="action-section">
        <v-btn
          color="primary"
          variant="elevated"
          block
          size="large"
          @click="downloadManual"
          class="manual-btn"
        >
          <v-icon start>mdi-file-pdf-box</v-icon>
          导出全部eSIM卡PDF
        </v-btn>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'
import type { ESIMInfo } from '@/types'

// Props 定义
interface Props {
  esimInfo: ESIMInfo
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
}>()

// 组合式 API
const appStore = useAppStore()

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

const downloadManual = () => {
  telegramService.impactFeedback('medium')
  
  appStore.showNotification({
    type: 'info',
    message: '正在生成 PDF 文档...',
    duration: 2000
  })
  
  // 这里可以调用 API 生成 PDF
  console.log('下载 eSIM 使用手册 PDF')
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
  .manual-btn {
    height: 48px;
    font-weight: 600;
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
}
</style>