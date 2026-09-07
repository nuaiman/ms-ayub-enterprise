<!-- src/components/features/godownStoreBills/GodownStoreBillDetail.vue -->
<template>
    <div v-if="store" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ store.lot_name }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ store.godown_name }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ store.customer_name }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="getStatusBadgeClass(store)">
                        <span class="w-1.5 h-1.5 rounded-full"
                            :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                        {{ getStatusLabel(store) }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">Store ID: {{ store.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Billing Start: {{ formatDate(store.billing_start)
            }}</span>
            <span v-if="store.billing_end" class="text-xs text-(--color-text-secondary)">
                Billing End: {{ formatDate(store.billing_end) }}
            </span>
        </div>

        <!-- Billing Details -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Store Details -->
            <div class="md:col-span-3">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">Store Details
                </p>
                <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Bill Type</p>
                        <p class="text-sm font-medium text-(--color-text-primary) capitalize">{{ store.store_bill_type
                        }}</p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Rate</p>
                        <p class="text-sm font-medium text-(--color-text-primary)">{{ formatCurrency(store.godown_cut)
                        }}</p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Quantity</p>
                        <p class="text-sm font-medium text-(--color-text-primary)">{{ store.quantity }} {{
                            store.quantity_unit }}</p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Weight</p>
                        <p class="text-sm font-medium text-(--color-text-primary)">{{ store.weight }} {{
                            store.weight_unit }}</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Financial Summary -->
        <div class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">Financial Summary
            </p>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Monthly Bill</p>
                    <p class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(store.monthly_bill) }}</p>
                </div>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Total Paid</p>
                    <p class="text-lg font-bold text-(--color-green)">{{ formatCurrency(store.total_billed) }}</p>
                </div>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                    <p class="text-lg font-bold text-(--color-red)">{{ formatCurrency(store.outstanding) }}</p>
                </div>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Last Paid Through</p>
                    <p class="text-sm font-medium text-(--color-text-primary)">{{ store.last_paid_through ?
                        formatDate(store.last_paid_through) : '—' }}</p>
                </div>
            </div>
        </div>

        <!-- Notes -->
        <div v-if="store.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ store.notes }}</p>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', store)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Update Payment
            </button>
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { GodownStoreBillStore } from '@/types/godownStoreBill'
import { useGodownStoreBillsStore } from '@/stores/godownStoreBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    store: GodownStoreBillStore | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [store: GodownStoreBillStore]
}>()

const godownStoreBillsStore = useGodownStoreBillsStore()

const getStatusBadgeClass = (store: GodownStoreBillStore): string => {
    return godownStoreBillsStore.getStatusBadgeClass(store)
}

const getStatusLabel = (store: GodownStoreBillStore): string => {
    return godownStoreBillsStore.getStatusLabel(store)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    })
}
</script>