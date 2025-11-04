<template>
  <v-card class="esim-info-card" variant="elevated">
    <v-card-title class="card-title">
      <v-icon start>mdi-sim</v-icon>
      eSIM 信息
      <v-spacer />
      <v-chip 
        :color="props.esimInfo.isRoaming ? 'warning' : 'success'" 
        size="small" 
        variant="flat"
      >
        {{ props.esimInfo.isRoaming ? '漫游' : '本地' }}
      </v-chip>
    </v-card-title>

    <v-card-text class="card-content">
      <!-- eSIM 标签 -->
      <div class="esim-tags">
        <v-chip 
          color="primary" 
          variant="elevated" 
          size="small" 
          class="esim-tag"
        >
          eSIM 1
        </v-chip>

        <v-chip 
          v-if="props.esimInfo.activatedAt" 
          color="success" 
          variant="flat" 
          size="small" 
          class="esim-tag"
        >
          已激活
        </v-chip>

        <v-chip 
          v-else 
          color="warning" 
          variant="flat" 
          size="small" 
          class="esim-tag"
        >
          未激活
        </v-chip>
      </div>

      <!-- ICCID -->
      <div class="info-section">
        <div class="section-header">
          <h4 class="section-title">ICCID</h4>
          <v-btn 
            v-if="props.allowCopy" 
            icon 
            variant="text" 
            size="small" 
            @click="copyICCID"
          >
            <v-icon size="16">mdi-content-copy</v-icon>
          </v-btn>
        </div>
        <p class="section-value">{{ props.esimInfo.iccid }}</p>
      </div>

      <!-- 激活码 -->
      <div class="info-section">
        <div class="section-header">
          <h4 class="section-title">激活码</h4>
          <v-btn 
            v-if="props.allowCopy" 
            icon 
            variant="text" 
            size="small" 
            @click="copyActivationCode"
          >
            <v-icon size="16">mdi-content-copy</v-icon>
          </v-btn>
        </div>
        <p class="section-value activation-code">{{ props.esimInfo.activationCode }}</p>
      </div>

      <!-- APN 类型 -->
      <div class="info-section">
        <div class="section-header">
          <h4 class="section-title">APN 类型</h4>
        </div>
        <p class="section-value">{{ apnTypeText }}</p>
      </div>

      <!-- 二维码 -->
      <div v-if="props.showQRCode" class="qr-section">
        <div class="section-header">
          <h4 class="section-title">激活二维码</h4>
          <div class="qr-actions">
            <v-btn 
              icon 
              variant="text" 
              size="small" 
              @click="downloadQR"
            >
              <v-icon size="16">mdi-download</v-icon>
            </v-btn>
            <v-btn 
              icon 
              variant="text" 
              size="small" 
              @click="shareQR"
            >
              <v-icon size="16">mdi-share</v-icon>
            </v-btn>
          </div>
        </div>

        <div class="qr-container">
          <QRCodeDisplay 
            :value="props.esimInfo.activationCode" 
            :size="200" 
            background-color="#ffffff"
            foreground-color="#000000" 
            @generated="handleQRGenerated" 
            @error="handleQRError" 
          />
        </div>

        <p class="qr-hint">
          使用手机相机扫描此二维码即可快速添加 eSIM
        </p>
      </div>
    </v-card-text>
  </v-card>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'
import QRCodeDisplay from '@/components/common/QRCodeDisplay.vue'
import type { ESIMInfo } from '@/types'

// Props
interface Props {
    esimInfo: ESIMInfo
    showQRCode?: boolean
    allowCopy?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    showQRCode: true,
    allowCopy: true
})

// Emits
const emit = defineEmits<{
    copyIccid: [iccid: string]
    copyActivationCode: [code: string]
    downloadQR: [qrCode: string]
    shareQR: [qrCode: string]
}>()

// 组合式 API
const appStore = useAppStore()

// 响应式状态
const qrDataUrl = ref<string>('')

// 计算属性
const apnTypeText = computed(() => {
    return props.esimInfo.apnType === 'manual' ? '手动配置' : '自动配置'
})

// 方法
const copyICCID = async () => {
    try {
        await navigator.clipboard.writeText(props.esimInfo.iccid)
        appStore.showNotification({
            type: 'success',
            message: 'ICCID 已复制到剪贴板',
            duration: 2000
        })
        telegramService.impactFeedback('light')
        emit('copyIccid', props.esimInfo.iccid)
    } catch (err) {
        appStore.showNotification({
            type: 'error',
            message: '复制失败',
            duration: 2000
        })
    }
}

const copyActivationCode = async () => {
    try {
        await navigator.clipboard.writeText(props.esimInfo.activationCode)
        appStore.showNotification({
            type: 'success',
            message: '激活码已复制到剪贴板',
            duration: 2000
        })
        telegramService.impactFeedback('light')
        emit('copyActivationCode', props.esimInfo.activationCode)
    } catch (err) {
        appStore.showNotification({
            type: 'error',
            message: '复制失败',
            duration: 2000
        })
    }
}

const downloadQR = () => {
    if (!qrDataUrl.value) return

    try {
        const link = document.createElement('a')
        link.download = `esim-qr-${props.esimInfo.iccid}.png`
        link.href = qrDataUrl.value
        link.click()

        appStore.showNotification({
            type: 'success',
            message: '二维码已下载',
            duration: 2000
        })
        telegramService.impactFeedback('medium')
        emit('downloadQR', qrDataUrl.value)
    } catch (err) {
        appStore.showNotification({
            type: 'error',
            message: '下载失败',
            duration: 2000
        })
    }
}

const shareQR = async () => {
    if (!qrDataUrl.value) return

    try {
        if (navigator.share) {
            // 使用 Web Share API
            const response = await fetch(qrDataUrl.value)
            const blob = await response.blob()
            const file = new File([blob], `esim-qr-${props.esimInfo.iccid}.png`, { type: 'image/png' })

            await navigator.share({
                title: 'eSIM 激活二维码',
                text: '扫描此二维码激活 eSIM',
                files: [file]
            })
        } else {
            // 降级到复制链接
            await navigator.clipboard.writeText(props.esimInfo.activationCode)
            appStore.showNotification({
                type: 'info',
                message: '激活码已复制，您可以手动分享',
                duration: 3000
            })
        }

        telegramService.impactFeedback('medium')
        emit('shareQR', qrDataUrl.value)
    } catch (err) {
        appStore.showNotification({
            type: 'error',
            message: '分享失败',
            duration: 2000
        })
    }
}

const handleQRGenerated = (dataUrl: string) => {
    qrDataUrl.value = dataUrl
}

const handleQRError = (error: Error) => {
    console.error('二维码生成失败:', error)
    appStore.showNotification({
        type: 'error',
        message: '二维码生成失败',
        duration: 2000
    })
}
</script><style scoped lang="scss">
.esim-info-card {
  .card-title {
    font-size: 1.1rem;
    font-weight: 600;
    padding: 16px 16px 8px 16px;
    
    .v-icon {
      margin-right: 8px;
    }
  }
  
  .card-content {
    padding: 8px 16px 16px 16px;
  }
}

.esim-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  
  .esim-tag {
    font-weight: 500;
  }
}

.info-section {
  margin-bottom: 20px;
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .section-title {
      font-size: 0.875rem;
      font-weight: 600;
      color: rgba(var(--v-theme-on-surface), 0.8);
      margin: 0;
    }
  }
  
  .section-value {
    font-size: 0.875rem;
    color: rgb(var(--v-theme-on-surface));
    margin: 0;
    word-break: break-all;
    line-height: 1.4;
    
    &.activation-code {
      font-family: 'Courier New', monospace;
      background: rgba(var(--v-theme-primary), 0.1);
      padding: 8px;
      border-radius: 4px;
      border: 1px solid rgba(var(--v-theme-primary), 0.2);
    }
  }
}

.qr-section {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  
  .section-header {
    margin-bottom: 16px;
    
    .qr-actions {
      display: flex;
      gap: 4px;
    }
  }
  
  .qr-container {
    display: flex;
    justify-content: center;
    margin-bottom: 12px;
    padding: 16px;
    background: #ffffff;
    border-radius: 8px;
    border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  }
  
  .qr-hint {
    font-size: 0.75rem;
    color: rgba(var(--v-theme-on-surface), 0.6);
    text-align: center;
    margin: 0;
    line-height: 1.4;
  }
}

// 响应式适配
@media (max-width: 360px) {
  .esim-tags {
    flex-wrap: wrap;
  }
  
  .info-section {
    .section-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 4px;
    }
  }
  
  .qr-container {
    padding: 12px;
  }
}
</style>