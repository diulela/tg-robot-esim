<template>
  <v-dialog v-model="dialog" max-width="500" persistent scrollable>
    <v-card>
      <v-card-title class="dialog-title">
        <v-icon start>mdi-cellphone-cog</v-icon>
        eSIM 安装说明
        <v-spacer />
        <v-btn icon variant="text" @click="closeDialog">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-card-text class="guide-content">
        <!-- 设备选择标签 -->
        <v-tabs v-model="activeTab" color="primary" class="device-tabs">
          <v-tab value="iphone">
            <v-icon start>mdi-apple</v-icon>
            iPhone
          </v-tab>
          <v-tab value="android">
            <v-icon start>mdi-android</v-icon>
            Android
          </v-tab>
        </v-tabs>

        <v-tabs-window v-model="activeTab" class="tabs-content">
          <!-- iPhone 安装说明 -->
          <v-tabs-window-item value="iphone">
            <div class="install-steps">
              <h4 class="steps-title">iPhone eSIM 安装步骤</h4>
              
              <div class="step-item" v-for="(step, index) in iphoneSteps" :key="index">
                <div class="step-number">{{ index + 1 }}</div>
                <div class="step-content">
                  <h5 class="step-title">{{ step.title }}</h5>
                  <p class="step-description">{{ step.description }}</p>
                  <div v-if="step.note" class="step-note">
                    <v-icon start size="16" color="info">mdi-information</v-icon>
                    {{ step.note }}
                  </div>
                </div>
              </div>

              <!-- 系统要求 -->
              <div class="requirements">
                <h5>系统要求</h5>
                <ul>
                  <li>iOS 12.1 或更高版本</li>
                  <li>支持 eSIM 的 iPhone 机型</li>
                  <li>稳定的网络连接</li>
                </ul>
              </div>
            </div>
          </v-tabs-window-item>

          <!-- Android 安装说明 -->
          <v-tabs-window-item value="android">
            <div class="install-steps">
              <h4 class="steps-title">Android eSIM 安装步骤</h4>
              
              <div class="step-item" v-for="(step, index) in androidSteps" :key="index">
                <div class="step-number">{{ index + 1 }}</div>
                <div class="step-content">
                  <h5 class="step-title">{{ step.title }}</h5>
                  <p class="step-description">{{ step.description }}</p>
                  <div v-if="step.note" class="step-note">
                    <v-icon start size="16" color="info">mdi-information</v-icon>
                    {{ step.note }}
                  </div>
                </div>
              </div>

              <!-- 系统要求 -->
              <div class="requirements">
                <h5>系统要求</h5>
                <ul>
                  <li>Android 9.0 或更高版本</li>
                  <li>支持 eSIM 的 Android 设备</li>
                  <li>稳定的网络连接</li>
                </ul>
              </div>
            </div>
          </v-tabs-window-item>
        </v-tabs-window>

        <!-- 常见问题 -->
        <div class="faq-section">
          <h4>常见问题</h4>
          <v-expansion-panels variant="accordion">
            <v-expansion-panel
              v-for="(faq, index) in faqs"
              :key="index"
              :title="faq.question"
            >
              <v-expansion-panel-text>
                {{ faq.answer }}
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </div>

        <!-- 注意事项 -->
        <div class="warnings">
          <h5>重要提示</h5>
          <div class="warning-item" v-for="warning in warnings" :key="warning">
            <v-icon start size="16" color="warning">mdi-alert</v-icon>
            {{ warning }}
          </div>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn 
          v-if="isIOS" 
          color="success" 
          variant="elevated" 
          @click="directInstall"
          :loading="isInstalling"
        >
          <v-icon start>mdi-apple</v-icon>
          在iPhone上直接安装
        </v-btn>
        <v-btn variant="outlined" @click="closeDialog">
          关闭
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/app'
import { telegramService } from '@/services/telegram'

// Props 定义
interface Props {
  modelValue: boolean
  orderId?: number
  esimInfo?: {
    iccid: string
    activationCode: string
    qrCode: string
  }
}

const props = defineProps<Props>()

// Emits 定义
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

// 组合式 API
const appStore = useAppStore()

// 响应式状态
const activeTab = ref('iphone')
const isInstalling = ref(false)
const userAgent = ref('')

const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 检测是否为 iOS 设备
const isIOS = computed(() => {
  return /iPad|iPhone|iPod/.test(userAgent.value) || 
         (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
})

// iPhone 安装步骤
const iphoneSteps = [
  {
    title: '打开设置',
    description: '在 iPhone 上打开"设置"应用',
    note: '确保设备已连接到稳定的 Wi-Fi 网络'
  },
  {
    title: '选择蜂窝网络',
    description: '点击"蜂窝网络"选项'
  },
  {
    title: '添加蜂窝套餐',
    description: '点击"添加蜂窝套餐"或"添加 eSIM"'
  },
  {
    title: '扫描二维码',
    description: '使用相机扫描提供的 eSIM 二维码，或手动输入激活码'
  },
  {
    title: '确认安装',
    description: '按照屏幕提示完成 eSIM 安装和激活'
  },
  {
    title: '设置标签',
    description: '为新的 eSIM 设置一个便于识别的标签名称'
  }
]

// Android 安装步骤
const androidSteps = [
  {
    title: '打开设置',
    description: '在 Android 设备上打开"设置"应用',
    note: '不同品牌的 Android 设备界面可能略有差异'
  },
  {
    title: '网络和互联网',
    description: '点击"网络和互联网"或"连接"选项'
  },
  {
    title: '移动网络',
    description: '选择"移动网络"或"SIM 卡管理"'
  },
  {
    title: '添加运营商',
    description: '点击"添加运营商"或"下载 SIM 卡"'
  },
  {
    title: '扫描二维码',
    description: '选择"扫描二维码"并扫描提供的 eSIM 二维码'
  },
  {
    title: '完成激活',
    description: '按照提示完成 eSIM 的下载和激活过程'
  }
]

// 常见问题
const faqs = [
  {
    question: '我的设备支持 eSIM 吗？',
    answer: '大部分 2018 年后发布的 iPhone 和支持 eSIM 的 Android 设备都可以使用。具体请查看设备规格或联系设备制造商。'
  },
  {
    question: 'eSIM 安装失败怎么办？',
    answer: '请确保设备连接稳定的网络，检查二维码是否清晰完整，或尝试手动输入激活码。如仍有问题请联系客服。'
  },
  {
    question: '可以同时使用多个 eSIM 吗？',
    answer: '大部分设备支持存储多个 eSIM，但通常只能同时激活 1-2 个。具体数量取决于设备型号。'
  },
  {
    question: 'eSIM 可以转移到其他设备吗？',
    answer: 'eSIM 通常与特定设备绑定，无法直接转移。如需更换设备，请联系客服获取新的 eSIM。'
  }
]

// 注意事项
const warnings = [
  '安装前请确保设备连接稳定的网络',
  '每个 eSIM 只能安装一次，请谨慎操作',
  '安装过程中请勿关闭应用或断开网络',
  '如遇问题请及时联系客服获取帮助'
]

// 方法
const closeDialog = () => {
  dialog.value = false
}

const directInstall = async () => {
  if (!props.esimInfo?.activationCode) {
    appStore.showNotification({
      type: 'error',
      message: '缺少激活码信息',
      duration: 2000
    })
    return
  }

  isInstalling.value = true
  
  try {
    // 构建 eSIM 安装 URL
    const installUrl = `https://esimsetup.apple.com/esim_qrcode_provisioning?carddata=${encodeURIComponent(props.esimInfo.activationCode)}`
    
    // 在 iOS 设备上打开安装链接
    if (isIOS.value) {
      window.location.href = installUrl
    } else {
      // 非 iOS 设备提示用户
      appStore.showNotification({
        type: 'info',
        message: '此功能仅支持 iOS 设备',
        duration: 3000
      })
    }
    
    telegramService.impactFeedback('medium')
    
  } catch (error) {
    appStore.showNotification({
      type: 'error',
      message: '启动安装失败，请手动安装',
      duration: 3000
    })
  } finally {
    isInstalling.value = false
  }
}

// 生命周期
onMounted(() => {
  userAgent.value = navigator.userAgent
  
  // 根据设备类型设置默认标签
  if (isIOS.value) {
    activeTab.value = 'iphone'
  } else {
    activeTab.value = 'android'
  }
})
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

.guide-content {
  padding: 0 16px 16px 16px;
}

.device-tabs {
  margin-bottom: 20px;
}

.tabs-content {
  min-height: 300px;
}

.install-steps {
  .steps-title {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 16px;
    color: rgb(var(--v-theme-on-surface));
  }
}

.step-item {
  display: flex;
  margin-bottom: 20px;
  align-items: flex-start;

  .step-number {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgb(var(--v-theme-primary));
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 0.875rem;
    margin-right: 16px;
    flex-shrink: 0;
  }

  .step-content {
    flex: 1;

    .step-title {
      font-size: 1rem;
      font-weight: 600;
      margin: 0 0 4px 0;
      color: rgb(var(--v-theme-on-surface));
    }

    .step-description {
      font-size: 0.875rem;
      color: rgba(var(--v-theme-on-surface), 0.7);
      margin: 0 0 8px 0;
      line-height: 1.4;
    }

    .step-note {
      display: flex;
      align-items: center;
      font-size: 0.75rem;
      color: rgb(var(--v-theme-info));
      background: rgba(var(--v-theme-info), 0.1);
      padding: 6px 8px;
      border-radius: 4px;

      .v-icon {
        margin-right: 6px;
      }
    }
  }
}

.requirements {
  margin-top: 24px;
  padding: 16px;
  background: rgba(var(--v-theme-surface-variant), 0.5);
  border-radius: 8px;

  h5 {
    font-size: 1rem;
    font-weight: 600;
    margin: 0 0 8px 0;
    color: rgb(var(--v-theme-on-surface));
  }

  ul {
    margin: 0;
    padding-left: 20px;

    li {
      font-size: 0.875rem;
      color: rgba(var(--v-theme-on-surface), 0.7);
      margin-bottom: 4px;
    }
  }
}

.faq-section {
  margin-top: 24px;

  h4 {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 12px;
    color: rgb(var(--v-theme-on-surface));
  }
}

.warnings {
  margin-top: 24px;
  padding: 16px;
  background: rgba(var(--v-theme-warning), 0.1);
  border-radius: 8px;
  border-left: 4px solid rgb(var(--v-theme-warning));

  h5 {
    font-size: 1rem;
    font-weight: 600;
    margin: 0 0 12px 0;
    color: rgb(var(--v-theme-warning));
  }

  .warning-item {
    display: flex;
    align-items: center;
    font-size: 0.875rem;
    color: rgba(var(--v-theme-on-surface), 0.8);
    margin-bottom: 8px;

    &:last-child {
      margin-bottom: 0;
    }

    .v-icon {
      margin-right: 8px;
    }
  }
}

// 响应式适配
@media (max-width: 360px) {
  .step-item {
    .step-number {
      width: 28px;
      height: 28px;
      font-size: 0.75rem;
      margin-right: 12px;
    }

    .step-content {
      .step-title {
        font-size: 0.9rem;
      }

      .step-description {
        font-size: 0.8rem;
      }
    }
  }

  .requirements,
  .warnings {
    padding: 12px;
  }
}
</style>