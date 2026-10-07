<!-- src/components/features/storeAdjustments/StoreAdjustmentRow.vue -->
<template>
    <div
        class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30">
        <!-- Store - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ storeLabel }}
            </div>
            <div v-if="lotLabel" class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                {{ lotLabel }}
            </div>
        </div>

        <!-- Type - 1 column -->
        <div class="col-span-1 min-w-0 pr-3">
            <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize" :class="adjustment.adjustment_type === 'delta'
                ? 'bg-(--color-blue)/10 text-(--color-blue)'
                : 'bg-(--color-yellow)/10 text-(--color-yellow)'">
                {{ adjustment.adjustment_type }}
            </span>
        </div>

        <!-- Delta (weight / quantity) - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <div class="text-sm font-semibold flex items-center gap-2">
                <span v-if="adjustment.weight_delta !== 0" :class="deltaClass(adjustment.weight_delta)">
                    {{ formatDelta(adjustment.weight_delta) }} {{ weightUnit }}
                </span>
                <span v-if="adjustment.quantity_delta !== 0" :class="deltaClass(adjustment.quantity_delta)">
                    {{ formatDelta(adjustment.quantity_delta) }} {{ quantityUnit }}
                </span>
            </div>
        </div>

        <!-- Reason - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ adjustment.reason || '—' }}
            </span>
        </div>

        <!-- By / When - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <div class="text-xs text-(--color-text-secondary) leading-tight truncate">
                {{ userName }}
            </div>
            <div class="text-xs text-(--color-text-secondary)/70 leading-tight">
                {{ formatDate(adjustment.adjusted_at) }}
            </div>
        </div>

        <!-- Actions - 1 column -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="handleDelete" :disabled="deleting"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md text-(--color-red) hover:bg-(--color-red)/10 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Delete adjustment">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { StoreAdjustment } from '@/types/storeAdjustment'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useUsersStore } from '@/stores/users'

const props = defineProps<{
    adjustment: StoreAdjustment
    deleting?: boolean
}>()

const emit = defineEmits<{
    'delete': [adjustment: StoreAdjustment]
}>()

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const usersStore = useUsersStore()

const store = computed(() => storesStore.getStoreById(props.adjustment.store_id))
const lot = computed(() => {
    if (!store.value) return null
    return lotsStore.getLotById(store.value.lot_id) ?? null
})

const storeLabel = computed(() =>
    store.value ? `Store #${store.value.id}` : `Store #${props.adjustment.store_id}`
)

const lotLabel = computed(() => {
    if (!lot.value) return ''
    return `${lot.value.product_name} · Lot ${lot.value.lot_number}`
})

const weightUnit = computed(() => lot.value?.weight_unit ?? '')
const quantityUnit = computed(() => lot.value?.quantity_unit ?? '')

const userName = computed(() => usersStore.getUserName(props.adjustment.user_id))

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(Math.abs(n))

const formatDelta = (n: number): string => `${n > 0 ? '+' : '−'}${formatNumber(n)}`

const deltaClass = (n: number): string =>
    n > 0 ? 'text-(--color-green)' : 'text-(--color-red)'

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit',
    })

const handleDelete = () => {
    emit('delete', props.adjustment)
}
</script>