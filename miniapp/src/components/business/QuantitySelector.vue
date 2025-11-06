<template>
    <div class="quantity-selector">
        <label class="quantity-label">购买数量：</label>

        <div class="quantity-controls">
            <v-btn icon variant="outlined" size="small" color="primary" :disabled="disabled || modelValue <= min"
                @click="decrease" class="quantity-btn decrease-btn">
                <v-icon>mdi-minus</v-icon>
            </v-btn>

            <div class="quantity-display">
                <span class="quantity-value">{{ modelValue }}</span>
            </div>

            <v-btn icon variant="outlined" size="small" color="primary"
                :disabled="disabled || (max !== undefined && modelValue >= max)" @click="increase"
                class="quantity-btn increase-btn">
                <v-icon>mdi-plus</v-icon>
            </v-btn>
        </div>
    </div>
</template>

<script setup lang="ts">
import { telegramService } from '@/services/telegram'

// Props 定义
interface Props {
    modelValue: number
    min?: number
    max?: number
    disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: 1,
    min: 1,
    disabled: false
})

// Emits 定义
const emit = defineEmits<{
    'update:modelValue': [value: number]
}>()

// 方法
const decrease = () => {
    if (props.disabled || props.modelValue <= props.min) return

    const newValue = props.modelValue - 1
    emit('update:modelValue', newValue)

    // 触觉反馈
    telegramService.impactFeedback('light')
}

const increase = () => {
    if (props.disabled || (props.max !== undefined && props.modelValue >= props.max)) return

    const newValue = props.modelValue + 1
    emit('update:modelValue', newValue)

    // 触觉反馈
    telegramService.impactFeedback('light')
}
</script>

<style scoped lang="scss">
.quantity-selector {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .quantity-label {
        font-size: 0.875rem;
        font-weight: 500;
        color: rgba(var(--v-theme-on-surface), 0.8);
    }

    .quantity-controls {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 16px;

        .quantity-btn {
            width: 40px;
            height: 40px;
            border-radius: 8px;
            transition: all 0.2s ease;

            &:hover:not(:disabled) {
                transform: scale(1.05);
            }

            &:active:not(:disabled) {
                transform: scale(0.95);
            }

            &:disabled {
                opacity: 0.4;
            }

            .v-icon {
                font-size: 18px;
            }
        }

        .quantity-display {
            display: flex;
            align-items: center;
            justify-content: center;
            min-width: 60px;
            height: 40px;
            background: rgba(var(--v-theme-primary), 0.1);
            border: 1px solid rgba(var(--v-theme-primary), 0.3);
            border-radius: 8px;
            transition: all 0.2s ease;

            .quantity-value {
                font-size: 1.125rem;
                font-weight: 600;
                color: rgb(var(--v-theme-primary));
                user-select: none;
            }
        }
    }
}

// 响应式适配
@media (max-width: 360px) {
    .quantity-selector {
        .quantity-controls {
            gap: 12px;

            .quantity-btn {
                width: 36px;
                height: 36px;

                .v-icon {
                    font-size: 16px;
                }
            }

            .quantity-display {
                min-width: 50px;
                height: 36px;

                .quantity-value {
                    font-size: 1rem;
                }
            }
        }
    }
}
</style>