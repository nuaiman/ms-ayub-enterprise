<!-- src/components/features/storeTransfers/StoreTransferRow.vue -->
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

        <!-- From → To - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="flex items-center gap-2 text-sm">
                <span class="text-(--color-text-secondary) truncate">
                    {{ fromGodownName }}
                </span>
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary) shrink-0" fill="none" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
                <span class="font-semibold text-(--color-text-primary) truncate">
                    {{ toGodownName }}
                </span>
            </div>
        </div>

        <!-- Snapshot - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <div class="text-xs text-(--color-text-secondary) truncate">
                Wt: {{ formatNumber(transfer.weight_at_transfer) }} {{ weightUnit }}
            </div>
            <div class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                Qty: {{ formatNumber(transfer.quantity_at_transfer) }} {{ quantityUnit }}
            </div>
        </div>

        <!-- By / When - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <div class="text-xs text-(--color-text-secondary) leading-tight truncate">
                {{ userName }}
            </div>
            <div class="text-xs text-(--color-text-secondary)/70 leading-tight">
                {{ formatDate(transfer.transferred_at) }}
            </div>
        </div>

        <!-- Actions - 1 column -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="handleDelete" :disabled="deleting || !isLatest"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md text-(--color-red) hover:bg-(--color-red)/10 transition-all duration-200 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:bg-transparent"
                :title="isLatest ? 'Delete transfer (reverses it)' : 'Cannot delete: a later transfer exists'">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { StoreTransfer } from '@/types/storeTransfer'
import { useStoresStore } from '@/stores/stores'
import { useGodownsStore } from '@/stores/godowns'
import { useLotsStore } from '@/stores/lots'
import { useUsersStore } from '@/stores/users'

const props = defineProps<{
    transfer: StoreTransfer
    deleting?: boolean
    isLatest?: boolean
}>()

const emit = defineEmits<{
    'delete': [transfer: StoreTransfer]
}>()

const storesStore = useStoresStore()
const godownsStore = useGodownsStore()
const lotsStore = useLotsStore()
const usersStore = useUsersStore()

const store = computed(() => storesStore.getStoreById(props.transfer.store_id))
const lot = computed(() => {
    if (!store.value) return null
    return lotsStore.getLotById(store.value.lot_id) ?? null
})

const storeLabel = computed(() =>
    store.value ? `Store #${store.value.id}` : `Store #${props.transfer.store_id}`
)

const lotLabel = computed(() => {
    if (!lot.value) return ''
    return `${lot.value.product_name} · Lot ${lot.value.lot_number}`
})

const weightUnit = computed(() => lot.value?.weight_unit ?? '')
const quantityUnit = computed(() => lot.value?.quantity_unit ?? '')

const fromGodownName = computed(() =>
    godownsStore.getGodownName(props.transfer.from_godown_id)
)

const toGodownName = computed(() =>
    godownsStore.getGodownName(props.transfer.to_godown_id)
)

const userName = computed(() => usersStore.getUserName(props.transfer.user_id))

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit',
    })

const handleDelete = () => {
    emit('delete', props.transfer)
}
</script>