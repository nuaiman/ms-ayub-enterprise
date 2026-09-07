<!-- src/components/features/stores/StoreDetail.vue -->
<template>
    <div v-if="store" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ getLotName(store.lot_id) }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">@ {{ getGodownName(store.godown_id) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="store.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                        <span class="w-1.5 h-1.5 rounded-full"
                            :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                        {{ store.is_active ? 'Active' : 'Inactive' }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ store.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(store.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(store.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Lot -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
                <p class="text-sm text-(--color-text-primary)">{{ getLotName(store.lot_id) }}</p>
            </div>

            <!-- Godown -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
                <p class="text-sm text-(--color-text-primary)">{{ getGodownName(store.godown_id) }}</p>
            </div>

            <!-- Store Bill Type -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                <p class="text-sm text-(--color-text-primary) capitalize">{{ store.store_bill_type }}</p>
            </div>

            <!-- Godown Cut -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown Cut</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(store.godown_cut) }}</p>
            </div>

            <!-- Quantity -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                <p class="text-sm text-(--color-text-primary)">{{ store.quantity }} {{ store.quantity_unit }}</p>
            </div>

            <!-- Weight -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                <p class="text-sm text-(--color-text-primary)">{{ store.weight }} {{ store.weight_unit }}</p>
            </div>

            <!-- Billing Start -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Billing Start</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatDateOnly(store.billing_start) }}</p>
            </div>

            <!-- Billing End -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Billing End</p>
                <p class="text-sm text-(--color-text-primary)">{{ store.billing_end ? formatDateOnly(store.billing_end)
                    : '—' }}</p>
            </div>

            <!-- Last Paid Through -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Last Paid Through
                </p>
                <p class="text-sm text-(--color-text-primary)">{{ store.last_paid_through ?
                    formatDateOnly(store.last_paid_through) : '—' }}</p>
            </div>

            <!-- Last Paid Amount -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Last Paid Amount
                </p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(store.last_paid_amount) }}</p>
            </div>

            <!-- Notes -->
            <div v-if="store.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ store.notes }}</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', store)"
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
import type { Store } from '@/types/store'
import { useLotsStore } from '@/stores/lots'
import { useGodownsStore } from '@/stores/godowns'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    store: Store | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [store: Store]
    'updated': []
}>()

const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()

const getLotName = (id: number): string => {
    return lotsStore.getLotName(id)
}

const getGodownName = (id: number): string => {
    return godownsStore.getGodownName(id)
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
</script>