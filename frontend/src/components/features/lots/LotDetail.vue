<!-- src/components/features/lots/LotDetail.vue -->
<template>
    <div v-if="lot" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ getItemName(lot.item_id) }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">Lot #{{ lot.lot_number }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatChargeType(lot.customer_charge_type)
                        }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="lot.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                        <span class="w-1.5 h-1.5 rounded-full"
                            :class="lot.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                        {{ lot.is_active ? 'Active' : 'Inactive' }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ lot.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(lot.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(lot.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Item -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Item</p>
                <p class="text-sm text-(--color-text-primary)">{{ getItemName(lot.item_id) }}</p>
            </div>

            <!-- Lot Number -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot Number</p>
                <p class="text-sm text-(--color-text-primary)">#{{ lot.lot_number }}</p>
            </div>

            <!-- Customer Charge Type -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer Charge
                    Type</p>
                <p class="text-sm text-(--color-text-primary) capitalize">{{ lot.customer_charge_type }}</p>
            </div>

            <!-- Majhi Bill Type -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi Bill Type
                </p>
                <p class="text-sm text-(--color-text-primary) capitalize">{{ lot.majhi_bill_type }}</p>
            </div>

            <!-- Customer Storage Rate -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Storage Rate</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.customer_storage_rate) }}</p>
            </div>

            <!-- Unload Rate -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Unload Rate</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.unload_rate) }}</p>
            </div>

            <!-- Majhi -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                <p class="text-sm text-(--color-text-primary)">{{ getMajhiName(lot.majhi_id) }}</p>
            </div>

            <!-- Majhi Cut -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi Cut</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.majhi_cut) }}</p>
            </div>

            <!-- Status -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Status</p>
                <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                    :class="lot.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                    <span class="w-1.5 h-1.5 rounded-full"
                        :class="lot.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                    {{ lot.is_active ? 'Active' : 'Inactive' }}
                </span>
            </div>

            <!-- Payment Summary Section -->
            <div class="md:col-span-2 border-t border-(--color-border) pt-4 mt-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">Payment
                    Summary</p>
                <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <!-- Customer Storage Payment -->
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Storage Bill - Paid Through</p>
                        <p class="text-sm font-medium text-(--color-text-primary)">
                            {{ lot.customer_last_paid_through
                                ? formatDateOnly(lot.customer_last_paid_through)
                                : 'Not paid yet'
                            }}
                        </p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Storage Bill - Total Paid</p>
                        <p class="text-sm font-semibold text-(--color-green)">
                            {{ formatCurrency(lot.customer_last_paid_amount) }}
                        </p>
                    </div>
                    <!-- Customer Unload Payment (NEW) -->
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Unload Bill - Paid Amount</p>
                        <p class="text-sm font-semibold" :class="getUnloadStatusColor(lot)">
                            {{ formatCurrency(lot.customer_paid_unload_amount) }}
                        </p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary)">Unload Bill - Status</p>
                        <span
                            class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                            :class="getUnloadStatusBadge(lot)">
                            <span class="w-1.5 h-1.5 rounded-full" :class="getUnloadStatusDot(lot)"></span>
                            {{ getUnloadStatusLabel(lot) }}
                        </span>
                    </div>
                    <!-- Majhi Payment -->
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border) md:col-span-4">
                        <p class="text-xs text-(--color-text-secondary)">Majhi Total Paid</p>
                        <p class="text-sm font-semibold text-(--color-blue)">
                            {{ formatCurrency(lot.majhi_total_paid) }}
                        </p>
                    </div>
                </div>
            </div>

            <!-- Notes -->
            <div v-if="lot.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ lot.notes }}</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', lot)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
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
import type { Lot, CustomerChargeType } from '@/types/lot'
import { useItemsStore } from '@/stores/items'
import { useMajhisStore } from '@/stores/majhis'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    lot: Lot | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [lot: Lot]
    'updated': []
}>()

const itemsStore = useItemsStore()
const majhisStore = useMajhisStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()

const getItemName = (id: number): string => {
    return itemsStore.getItemName(id)
}

const getMajhiName = (id: number | null): string => {
    if (!id) return '—'
    return majhisStore.getMajhiName(id)
}

const formatChargeType = (type: CustomerChargeType): string => {
    return type.charAt(0).toUpperCase() + type.slice(1)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

const formatDateOnly = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    })
}

// Unload bill helpers
const getUnloadBillAmount = (lot: Lot): number => {
    const stores = storesStore.getStoresByLotId(lot.id)
    let totalQuantity = 0
    let totalWeight = 0
    for (const store of stores) {
        totalQuantity += store.quantity
        totalWeight += store.weight
    }
    if (lot.customer_charge_type === 'quantity') {
        return lot.unload_rate * totalQuantity
    } else {
        return lot.unload_rate * totalWeight
    }
}

const isUnloadPaid = (lot: Lot): boolean => {
    const billAmount = getUnloadBillAmount(lot)
    if (billAmount === 0) return true
    return (lot.customer_paid_unload_amount || 0) >= billAmount
}

const getUnloadStatusLabel = (lot: Lot): string => {
    if (lot.unload_rate === 0) return 'No Bill'
    return isUnloadPaid(lot) ? 'Paid' : 'Unpaid'
}

const getUnloadStatusColor = (lot: Lot): string => {
    if (lot.unload_rate === 0) return 'text-(--color-text-secondary)'
    return isUnloadPaid(lot) ? 'text-(--color-green)' : 'text-(--color-red)'
}

const getUnloadStatusBadge = (lot: Lot): string => {
    if (lot.unload_rate === 0) return 'border-(--color-text-secondary) text-(--color-text-secondary)'
    return isUnloadPaid(lot)
        ? 'border-(--color-green) text-(--color-green)'
        : 'border-(--color-red) text-(--color-red)'
}

const getUnloadStatusDot = (lot: Lot): string => {
    if (lot.unload_rate === 0) return 'bg-(--color-text-secondary)'
    return isUnloadPaid(lot) ? 'bg-(--color-green)' : 'bg-(--color-red)'
}
</script>