<template>
  <div class="alphabet-index">
    <div
      v-for="letter in alphabet"
      :key="letter"
      :class="[
        'index-letter',
        { 
          active: letter === currentLetter,
          disabled: !availableLetters.includes(letter)
        }
      ]"
      @click="handleLetterClick(letter)"
      @touchstart.prevent="handleTouchStart(letter)"
      @touchmove.prevent="handleTouchMove"
      @touchend.prevent="handleTouchEnd"
    >
      {{ letter }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// Props
interface Props {
  availableLetters: string[]
  currentLetter?: string
}

const props = withDefaults(defineProps<Props>(), {
  currentLetter: ''
})

// Emits
const emit = defineEmits<{
  'letter-click': [letter: string]
}>()

// 字母表
const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('')

// 触摸状态
const isTouching = ref(false)

// 处理字母点击
const handleLetterClick = (letter: string) => {
  if (props.availableLetters.includes(letter)) {
    emit('letter-click', letter)
  }
}

// 处理触摸开始
const handleTouchStart = (letter: string) => {
  isTouching.value = true
  if (props.availableLetters.includes(letter)) {
    emit('letter-click', letter)
  }
}

// 处理触摸移动
const handleTouchMove = (event: TouchEvent) => {
  if (!isTouching.value) return

  const touch = event.touches[0]
  const element = document.elementFromPoint(touch.clientX, touch.clientY)

  if (element && element.classList.contains('index-letter')) {
    const letter = element.textContent?.trim()
    if (letter && props.availableLetters.includes(letter)) {
      emit('letter-click', letter)
    }
  }
}

// 处理触摸结束
const handleTouchEnd = () => {
  isTouching.value = false
}
</script>

<style scoped>
.alphabet-index {
  position: fixed;
  right: 4px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  z-index: 100;
  padding: 8px 4px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.index-letter {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  font-size: 11px;
  font-weight: 600;
  color: #667eea;
  cursor: pointer;
  user-select: none;
  transition: all 0.2s ease;
  border-radius: 50%;
}

.index-letter:hover:not(.disabled) {
  background: rgba(102, 126, 234, 0.1);
  transform: scale(1.2);
}

.index-letter.active {
  background: #667eea;
  color: #ffffff;
  transform: scale(1.3);
}

.index-letter.disabled {
  color: #cccccc;
  cursor: not-allowed;
  opacity: 0.5;
}

/* 触摸反馈 */
.index-letter:active:not(.disabled) {
  background: rgba(102, 126, 234, 0.2);
}

/* 响应式设计 */
@media (max-width: 360px) {
  .alphabet-index {
    right: 2px;
    padding: 6px 2px;
    gap: 1px;
  }

  .index-letter {
    width: 18px;
    height: 18px;
    font-size: 10px;
  }
}

@media (min-width: 480px) {
  .alphabet-index {
    right: 8px;
    padding: 10px 6px;
    gap: 3px;
  }

  .index-letter {
    width: 24px;
    height: 24px;
    font-size: 12px;
  }
}
</style>
