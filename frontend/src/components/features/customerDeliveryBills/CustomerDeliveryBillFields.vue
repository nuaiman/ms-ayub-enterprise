<!-- src/components/features/customerDeliveryBills/CustomerDeliveryBillFields.vue -->
<template>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Bill Type <span class="text-(--color-red)">*</span>
            </label>
            <select :value="billType"
                @change="$emit('update:billType', ($event.target as HTMLSelectElement).value as CustomerDeliveryBillType)"
                :disabled="disabled" required
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option value="quantity">Quantity</option>
                <option value="weight">Weight</option>
            </select>
        </div>
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Rate <span class="text-(--color-red)">*</span>
            </label>
            <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                <input :value="rate"
                    @input="$emit('update:rate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                    class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CustomerDeliveryBillType } from '@/types/customerDeliveryBill'

withDefaults(defineProps<{
    billType: CustomerDeliveryBillType
    rate: number
    disabled?: boolean
}>(), {
    disabled: false,
})

defineEmits<{
    (e: 'update:billType', value: CustomerDeliveryBillType): void
    (e: 'update:rate', value: number): void
}>()
</script>