<!-- src/components/features/majhiBills/MajhiBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Majhi Bill #{{ bill.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ majhiLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ storeLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary) capitalize">
                        {{ bill.bill_type === 'job' ? 'Job (Fixed)' : bill.bill_type }}
                    </span>
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
                <p class="text-lg font-bold text-(--color-text-primary) mt-1 capitalize">
                    {{ bill.bill_type === 'job' ? 'Job' : bill.bill_type }}
                </p>
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
                    :class="outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(outstanding) }}
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
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                    <p class="text-sm text-(--color-text-primary)">{{ majhiLabel }}</p>
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
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.total_paid) }}</p>
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
import type { MajhiBill } from '@/types/majhiBill'
import { useMajhisStore } from '@/stores/majhis'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: MajhiBill | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [bill: MajhiBill]
    'payment': [bill: MajhiBill]
    'updated': []
}>()

const majhisStore = useMajhisStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()

const majhiLabel = computed(() => {
    if (!props.bill) return 'â€”'
    return majhisStore.getMajhiName(props.bill.majhi_id)
})

const storeLabel = computed(() => {
    if (!props.bill) return 'â€”'

    const storeId = props.bill.store_id
    if (storeId === null || storeId === undefined) {
        return props.bill.delivery_item_id ? `Delivery Item #${props.bill.delivery_item_id}` : 'â€”'
    }

    const store = storesStore.getStoreById(storeId)
    if (!store) return `Store #${storeId}`
    const lot = lotsStore.getLotById(store.lot_id)
    if (lot) return `Store #${store.id} â€” ${lot.product_name} (Lot ${lot.lot_number})`
    return `Store #${store.id}`
})

const outstanding = computed(() => {
    if (!props.bill) return 0
    return Math.max(0, (props.bill.rate || 0) - (props.bill.total_paid || 0))
})

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
}
</script>