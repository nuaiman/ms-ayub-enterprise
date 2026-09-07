<!-- src/components/features/customerTransportBills/CustomerTransportBillRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Transport - 1 column -->
        <div class="col-span-1">
            <span class="text-sm font-medium text-(--color-text-primary)">
                #{{ bill.transport_id }}
            </span>
        </div>

        <!-- Customer - 2 columns -->
        <div class="col-span-2 min-w-0">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ bill.customer_name }}
            </div>
        </div>

        <!-- From - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ bill.from_location }}
            </span>
        </div>

        <!-- To - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ bill.to_location || '—' }}
            </span>
        </div>

        <!-- Vehicles - 1 column -->
        <div class="col-span-1">
            <span class="text-sm text-(--color-text-secondary)">
                {{ bill.total_vehicles }}
            </span>
        </div>

        <!-- Amount - 2 columns -->
        <div class="col-span-2">
            <span class="text-sm font-semibold text-(--color-text-primary)">
                {{ formatCurrency(bill.bill_amount) }}
            </span>
            <span v-if="outstanding > 0" class="text-xs text-(--color-red) block">
                {{ formatCurrency(outstanding) }} outstanding
            </span>
        </div>

        <!-- Status - 1 column -->
        <div class="col-span-1">
            <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="getStatusBadgeClass(bill.status)">
                <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(bill.status)"></span>
                {{ getStatusLabel(bill.status) }}
            </span>
        </div>

        <!-- Actions - 1 column, right aligned -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
                    <button @click="handleView"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                        </svg>
                        View Details
                    </button>

                    <button v-if="outstanding > 0" @click="handlePay"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                        Record Payment
                    </button>

                    <button v-if="outstanding > 0" @click="handleCancel"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-yellow) hover:bg-(--color-muted-bg) transition-colors border-t border-(--color-border) mt-1 pt-1">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                        Cancel Bill
                    </button>
                </div>
            </Transition>

            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { CustomerTransportBill } from '@/types/customerTransportBill'
import { useCustomerTransportBillsStore } from '@/stores/customerTransportBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerTransportBill
}>()

const emit = defineEmits<{
    'view': [bill: CustomerTransportBill]
    'pay': [bill: CustomerTransportBill]
    'cancel': [bill: CustomerTransportBill]
}>()

const customerTransportBillsStore = useCustomerTransportBillsStore()
const isOpen = ref(false)

const outstanding = computed(() => {
    return props.bill.bill_amount - props.bill.paid_amount
})

const getStatusBadgeClass = (status: string): string => {
    return customerTransportBillsStore.getStatusBadgeClass(status)
}

const getStatusDotClass = (status: string): string => {
    return customerTransportBillsStore.getStatusDotClass(status)
}

const getStatusLabel = (status: string): string => {
    return customerTransportBillsStore.getStatusLabel(status)
}

const toggleMenu = () => {
    isOpen.value = !isOpen.value
}

const closeMenu = () => {
    isOpen.value = false
}

const handleView = () => {
    closeMenu()
    emit('view', props.bill)
}

const handlePay = () => {
    closeMenu()
    emit('pay', props.bill)
}

const handleCancel = () => {
    closeMenu()
    emit('cancel', props.bill)
}
</script>