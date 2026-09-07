<!-- src/components/features/customerStorageBills/CustomerStorageBillRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Customer - 3 columns -->
        <div class="col-span-3 min-w-0">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ bill.customer_name }}
            </div>
        </div>

        <!-- Item / Lot - 3 columns -->
        <div class="col-span-3 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ bill.item_name }}
            </span>
            <span class="text-xs text-(--color-text-secondary) truncate block">
                Lot #{{ bill.lot_id }}
            </span>
        </div>

        <!-- Monthly Bill - 2 columns -->
        <div class="col-span-2">
            <span class="text-sm font-semibold text-(--color-text-primary)">
                {{ formatCurrency(bill.monthly_bill) }}
            </span>
            <span class="text-xs text-(--color-text-secondary) block">
                {{ bill.customer_charge_type }}
            </span>
        </div>

        <!-- Outstanding - 2 columns -->
        <div class="col-span-2">
            <span class="text-sm font-semibold"
                :class="bill.outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                {{ formatCurrency(bill.outstanding) }}
            </span>
            <span class="text-xs text-(--color-text-secondary) block">
                {{ bill.months_billed }} month(s)
            </span>
        </div>

        <!-- Actions - 2 columns, right aligned -->
        <div class="col-span-2 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <!-- Dropdown -->
            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
                    <!-- View Details -->
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

                    <!-- Record Payment -->
                    <button @click="handlePay"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                        Record Payment
                    </button>
                </div>
            </Transition>

            <!-- Backdrop -->
            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { CustomerStorageBill } from '@/types/customerStorageBill'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerStorageBill
}>()

const emit = defineEmits<{
    'view': [bill: CustomerStorageBill]
    'pay': [bill: CustomerStorageBill]
}>()

const isOpen = ref(false)

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
</script>