<!-- src/components/features/godownBills/GodownBillFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Godown picker (hidden when locked or derived) -->
        <div v-if="!locked && !derived">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Godown <span class="text-(--color-red)">*</span>
            </label>
            <select :value="godownId ?? ''"
                @change="$emit('update:godownId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled" required
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option value="">Select a godown</option>
                <option v-for="godown in godownOptions" :key="godown.id" :value="godown.id">
                    {{ godown.name }}
                </option>
            </select>
        </div>

        <!-- Locked display (edit mode) -->
        <div v-else-if="locked && !derived"
            class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
            <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedGodownLabel }}</p>
        </div>

        <!-- Bill Type + Rate (two-column) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Bill Type -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Bill Type <span class="text-(--color-red)">*</span>
                </label>
                <select :value="billType"
                    @change="$emit('update:billType', ($event.target as HTMLSelectElement).value as GodownBillType)"
                    :disabled="disabled" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="fixed">Fixed</option>
                    <option value="quantity">Quantity</option>
                    <option value="weight">Weight</option>
                </select>
            </div>

            <!-- Rate -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Rate <span class="text-(--color-red)">*</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="rate"
                        @input="$emit('update:rate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Godown } from '@/types/godown'
import type { GodownBillType } from '@/types/godownBill'

withDefaults(defineProps<{
    godownId?: number | null
    billType: GodownBillType
    rate: number
    godownOptions?: Godown[]
    disabled?: boolean
    locked?: boolean
    lockedGodownLabel?: string
    derived?: boolean
}>(), {
    godownId: null,
    godownOptions: () => [],
    disabled: false,
    locked: false,
    lockedGodownLabel: '',
    derived: false,
})

defineEmits<{
    (e: 'update:godownId', value: number | null): void
    (e: 'update:billType', value: GodownBillType): void
    (e: 'update:rate', value: number): void
}>()
</script>