<template>
  <div class="qr-code-display">
    <canvas
      ref="canvasRef"
      :width="size"
      :height="size"
      class="qr-canvas"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import QRCode from 'qrcode'

// Props
interface Props {
  value: string
  size?: number
  level?: 'L' | 'M' | 'Q' | 'H'
  includeMargin?: boolean
  backgroundColor?: string
  foregroundColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 200,
  level: 'M',
  includeMargin: true,
  backgroundColor: '#ffffff',
  foregroundColor: '#000000'
})

// Emits
const emit = defineEmits<{
  generated: [dataUrl: string]
  error: [error: Error]
}>()

// 响应式状态
const canvasRef = ref<HTMLCanvasElement | null>(null)

// 方法
const generateQRCode = async () => {
  if (!canvasRef.value || !props.value) return

  try {
    await QRCode.toCanvas(canvasRef.value, props.value, {
      width: props.size,
      margin: props.includeMargin ? 2 : 0,
      color: {
        dark: props.foregroundColor,
        light: props.backgroundColor
      },
      errorCorrectionLevel: props.level
    })

    // 获取 data URL 并发射事件
    const dataUrl = canvasRef.value.toDataURL('image/png')
    emit('generated', dataUrl)
  } catch (error) {
    console.error('二维码生成失败:', error)
    emit('error', error as Error)
  }
}

// 监听属性变化
watch(
  () => [props.value, props.size, props.level, props.backgroundColor, props.foregroundColor],
  () => generateQRCode(),
  { deep: true }
)

// 生命周期
onMounted(() => {
  generateQRCode()
})
</script>

<style scoped lang="scss">
.qr-code-display {
  display: inline-block;
  
  .qr-canvas {
    display: block;
    border-radius: 4px;
  }
}
</style>