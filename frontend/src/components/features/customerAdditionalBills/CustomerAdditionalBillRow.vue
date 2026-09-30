<!-- src/components/features/customerAdditionalBills/CustomerAdditionalBillRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Customer - 4 columns -->
        <div class="col-span-4 min-w-0">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ bill.customer_name }}
            </div>
        </div>

        <!-- Description - 4 columns -->
        <div class="col-span-4 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ bill.description }}
            </span>
        </div>

        <!-- Amount - 2 columns -->
        <div class="col-span-2">
            <span class="text-sm font-semibold text-(--color-text-primary)">
                {{ formatCurrency(bill.amount) }}
            </span>
            <span v-if="outstanding > 0" class="text-xs text-(--color-red) block">
                {{ formatCurrency(outstanding) }} outstanding
            </span>
        </div>

        <!-- Status - 2 columns -->
        <div class="col-span-2">
            <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="getStatusBadgeClass(bill.status)">
                <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(bill.status)"></span>
                {{ getStatusLabel(bill.status) }}
            </span>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { CustomerAdditionalBill } from '@/types/customerAdditionalBill'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerAdditionalBill
}>()

const emit = defineEmits<{
    'view': [bill: CustomerAdditionalBill]
}>()

const store = useCustomerAdditionalBillsStore()

const outstanding = computed(() => props.bill.outstanding)

const getStatusBadgeClass = (status: string): string => store.getStatusBadgeClass(status)
const getStatusDotClass = (status: string): string => store.getStatusDotClass(status)
const getStatusLabel = (status: string): string => store.getStatusLabel(status)

const handleView = () => {
    emit('view', props.bill)
}
</script>