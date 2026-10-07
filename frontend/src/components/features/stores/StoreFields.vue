<!-- src/components/features/stores/StoreFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Weight & Quantity -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Weight <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <input :value="weight"
                    @input="$emit('update:weight', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Quantity <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <input :value="quantity"
                    @input="$emit('update:quantity', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Start Date -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Start Date
            </label>
            <input :value="startDate" @input="$emit('update:startDate', ($event.target as HTMLInputElement).value)"
                type="date" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    weight: number
    quantity: number
    startDate: string
    disabled?: boolean
    required?: boolean
}>()

defineEmits<{
    (e: 'update:weight', value: number): void
    (e: 'update:quantity', value: number): void
    (e: 'update:startDate', value: string): void
}>()
</script>