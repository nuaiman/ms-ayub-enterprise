<!-- src/components/ui/BaseDialog.vue -->
<template>
    <Teleport to="body">
        <div v-if="modelValue"
            class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm transition-opacity p-4"
            @click.self="handleOutsideClick">
            <div class="bg-(--color-surface) rounded-2xl shadow-2xl border border-(--color-border) max-h-[90vh] overflow-y-auto w-full transition-all duration-200 relative hide-scrollbar"
                :class="maxWidthClass" @click.stop>
                <!-- Floating Close Button - Centered with content -->
                <button @click="handleClose"
                    class="absolute top-4 right-4 z-20 p-2 rounded-full hover:bg-(--color-muted-bg) transition-colors"
                    aria-label="Close dialog">
                    <svg class="w-5 h-5 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                        viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>

                <!-- Content -->
                <div class="px-4 sm:px-6 py-5 sm:py-6">
                    <slot>
                        <p v-if="message" class="text-sm text-(--color-text-secondary)">{{ message }}</p>
                    </slot>
                </div>

                <!-- Footer Actions -->
                <div v-if="$slots.actions"
                    class="sticky bottom-0 bg-(--color-surface) rounded-b-2xl px-4 sm:px-6 py-3 sm:py-4 border-t border-(--color-border)/50 flex flex-col sm:flex-row items-center gap-2 sm:gap-3 justify-end">
                    <slot name="actions" />
                </div>
            </div>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type DialogVariant = 'danger' | 'warning' | 'info' | 'success'

const props = withDefaults(defineProps<{
    modelValue: boolean
    title?: string
    message?: string
    icon?: unknown
    variant?: DialogVariant
    maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full' | '2xl' | '3xl'
    closeOnOutsideClick?: boolean
}>(), {
    variant: 'info',
    maxWidth: '3xl',
    closeOnOutsideClick: false,
})

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'close'): void
}>()

const maxWidthClass = computed(() => {
    const map = {
        sm: 'max-w-sm',
        md: 'max-w-md',
        lg: 'max-w-lg',
        xl: 'max-w-xl',
        full: 'max-w-4xl',
        '2xl': 'max-w-5xl',
        '3xl': 'max-w-6xl',
    }
    return map[props.maxWidth] || 'max-w-3xl'
})

const handleOutsideClick = () => {
    if (props.closeOnOutsideClick) {
        handleClose()
    }
}

const handleClose = () => {
    emit('update:modelValue', false)
    emit('close')
}
</script>