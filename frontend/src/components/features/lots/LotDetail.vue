<!-- src/components/features/lots/LotDetail.vue -->
<template>
    <div v-if="lot" class="space-y-6">
        <!-- Header -->
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ lot.product_name }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">Lot {{ lot.lot_number }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
                </div>
            </div>
        </div>

        <!-- Summary strip -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Stock on Hand</p>
                <div class="mt-1">
                    <p v-if="stockTotals.length === 0" class="text-lg font-bold text-(--color-text-primary)">—</p>
                    <p v-else v-for="(line, idx) in stockTotals" :key="idx"
                        class="text-lg font-bold text-(--color-text-primary) leading-tight">
                        {{ line }}
                    </p>
                </div>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Stores</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">
                    {{ activeStoresCount }} <span class="text-xs font-medium text-(--color-text-secondary)">/
                        {{ stores.length }} active</span>
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Weight Unit</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ lot.weight_unit }}</p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Quantity Unit</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ lot.quantity_unit }}</p>
            </div>
        </div>

        <!-- Lot info -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Lot Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot Number
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.lot_number }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Product Name
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.product_name }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight Unit
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.weight_unit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity
                        Unit</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.quantity_unit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(lot.created_at) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Updated</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(lot.updated_at) }}</p>
                </div>
            </div>
        </section>

        <!-- Stores -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Stores</h3>
                <span class="text-xs text-(--color-text-secondary)">{{ stores.length }}</span>
            </div>

            <div v-if="stores.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No stores for this lot.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="store in stores" :key="store.id" class="p-4">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div class="min-w-0">
                            <div class="flex items-center gap-2 flex-wrap">
                                <span class="text-sm font-semibold text-(--color-text-primary)">
                                    Store #{{ store.id }}
                                </span>
                                <span
                                    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                                    :class="store.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                                    <span class="w-1.5 h-1.5 rounded-full"
                                        :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                                    {{ store.is_active ? 'Active' : 'Inactive' }}
                                </span>
                            </div>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-xs text-(--color-text-secondary)">
                                <span>Quantity: <span class="text-(--color-text-primary)">{{ store.quantity }} {{
                                    lot.quantity_unit }}</span></span>
                                <span>Weight: <span class="text-(--color-text-primary)">{{ store.weight }} {{
                                    lot.weight_unit }}</span></span>
                                <span>Start: <span class="text-(--color-text-primary)">{{
                                    formatDateShort(store.start_date) }}</span></span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

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
import type { Lot } from '@/types/lot'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'

const props = defineProps<{
    lot: Lot | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [lot: Lot]
    'updated': []
}>()

const customersStore = useCustomersStore()
const storesStore = useStoresStore()

const customerName = computed(() => {
    if (!props.lot) return '—'
    return customersStore.getCustomerName(props.lot.customer_id)
})

const stores = computed(() => {
    if (!props.lot) return []
    return storesStore.getStoresByLotId(props.lot.id)
})

const activeStoresCount = computed(() => stores.value.filter(s => s.is_active).length)

const stockTotals = computed<string[]>(() => {
    const byUnit: Record<string, number> = {}
    for (const s of stores.value) {
        if (!s.is_active) continue
        if (s.quantity > 0) {
            const u = props.lot?.quantity_unit || 'units'
            byUnit[u] = (byUnit[u] || 0) + s.quantity
        }
        if (s.weight > 0) {
            const u = props.lot?.weight_unit || 'kg'
            byUnit[u] = (byUnit[u] || 0) + s.weight
        }
    }
    const lines: string[] = []
    for (const [unit, total] of Object.entries(byUnit)) {
        const n = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(total)
        lines.push(`${n} ${unit}`)
    }
    return lines
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