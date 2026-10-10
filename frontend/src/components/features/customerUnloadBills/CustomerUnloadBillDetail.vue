<!-- src/components/features/customerUnloadBills/CustomerUnloadBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Customer Unload Bill #{{ bill.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ customerLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ storeLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary) capitalize">{{ bill.bill_type }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(bill.rate)
                    }}</span>
                </div>
            </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ bill.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(bill.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(bill.updated_at) }}</span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1 capitalize">{{ bill.bill_type }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Rate</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ formatCurrency(bill.rate) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                <p class="text-lg font-bold text-(--color-green) mt-1">{{ formatCurrency(bill.total_paid) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                <p class="text-lg font-bold mt-1"
                    :class="bill.remaining > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(bill.remaining) }}
                </p>
            </div>
        </div>

        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Bill Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerLabel }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                    <p class="text-sm text-(--color-text-primary)">{{ storeLabel }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ bill.bill_type }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Rate</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.rate) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight at
                        Billing</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.weight_at_billing }} {{ bill.weight_unit_at_billing || '' }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity at
                        Billing</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.quantity_at_billing }} {{ bill.quantity_unit_at_billing || '' }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Amount
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.total_amount) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.total_paid) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid
                        Through</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.total_paid_through ? formatDateShort(bill.total_paid_through) : '-' }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(bill.created_at) }}</p>
                </div>
            </div>
        </section>

        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('payment', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-green) text-white hover:opacity-90 transition-all duration-200">
                Add Payment
            </button>
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { CustomerUnloadBill } from '@/types/customerUnloadBill'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerUnloadBill | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [bill: CustomerUnloadBill]
    'payment': [bill: CustomerUnloadBill]
    'updated': []
}>()

const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()

const customerLabel = computed(() => {
    if (!props.bill) return '-'
    return customersStore.getCustomerName(props.bill.customer_id)
})

const storeLabel = computed(() => {
    if (!props.bill) return '-'
    const store = storesStore.getStoreById(props.bill.store_id)
    if (store) {
        const lot = lotsStore.getLotById(store.lot_id)
        if (lot) return `Store #${store.id} - ${lot.product_name} (Lot ${lot.lot_number})`
    }
    return `Store #${props.bill.store_id}`
})

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
}

const formatDateShort = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })
}
</script>