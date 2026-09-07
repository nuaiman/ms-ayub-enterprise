<!-- src/components/features/godownStoreBills/GodownStoreBillRow.vue -->
<template>
    <div
        class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30">
        <!-- Lot - 3 columns -->
        <div class="col-span-3 min-w-0 cursor-pointer" @click="handleView">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ store.lot_name }}
            </div>
            <div class="text-xs text-(--color-text-secondary) truncate">
                {{ store.item_name }}
            </div>
        </div>

        <!-- Godown - 2 columns -->
        <div class="col-span-2 min-w-0 cursor-pointer" @click="handleView">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ store.godown_name }}
            </span>
        </div>

        <!-- Customer - 2 columns -->
        <div class="col-span-2 min-w-0 cursor-pointer" @click="handleView">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ store.customer_name }}
            </span>
        </div>

        <!-- Monthly Bill - 1 column -->
        <div class="col-span-1 cursor-pointer" @click="handleView">
            <span class="text-sm font-semibold text-(--color-text-primary)">
                {{ formatCurrency(store.monthly_bill) }}
            </span>
        </div>

        <!-- Outstanding - 1 column -->
        <div class="col-span-1 cursor-pointer" @click="handleView">
            <span class="text-sm font-semibold"
                :class="(store.outstanding || 0) > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                {{ formatCurrency(store.outstanding) }}
            </span>
        </div>

        <!-- Status - 1 column -->
        <div class="col-span-1 cursor-pointer" @click="handleView">
            <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="getStatusBadgeClass(store)">
                <span class="w-1.5 h-1.5 rounded-full"
                    :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                {{ getStatusLabel(store) }}
            </span>
        </div>

        <!-- Actions - 2 columns -->
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
                    <button v-if="store.is_active && (store.outstanding || 0) > 0" @click="handlePay"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                        Record Payment
                    </button>

                    <!-- End Billing -->
                    <button v-if="store.is_active && !store.billing_end" @click="handleEndBilling"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-yellow) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                        End Billing
                    </button>

                    <!-- Reactivate Billing -->
                    <button v-if="!store.is_active || store.billing_end" @click="handleReactivate"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                        Reactivate Billing
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
import type { GodownStoreBillStore } from '@/types/godownStoreBill'
import { useGodownStoreBillsStore } from '@/stores/godownStoreBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    store: GodownStoreBillStore
}>()

const emit = defineEmits<{
    'view': [store: GodownStoreBillStore]
    'pay': [store: GodownStoreBillStore]
    'end-billing': [store: GodownStoreBillStore]
    'reactivate': [store: GodownStoreBillStore]
}>()

const godownStoreBillsStore = useGodownStoreBillsStore()
const isOpen = ref(false)

const getStatusBadgeClass = (store: GodownStoreBillStore): string => {
    return godownStoreBillsStore.getStatusBadgeClass(store)
}

const getStatusLabel = (store: GodownStoreBillStore): string => {
    return godownStoreBillsStore.getStatusLabel(store)
}

const toggleMenu = () => {
    isOpen.value = !isOpen.value
}

const closeMenu = () => {
    isOpen.value = false
}

const handleView = () => {
    closeMenu()
    emit('view', props.store)
}

const handlePay = () => {
    closeMenu()
    emit('pay', props.store)
}

const handleEndBilling = () => {
    closeMenu()
    emit('end-billing', props.store)
}

const handleReactivate = () => {
    closeMenu()
    emit('reactivate', props.store)
}
</script>