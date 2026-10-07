<!-- src/components/features/lotTransfers/LotTransferRow.vue -->
<template>
    <div
        class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30">
        <!-- Lot - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ lotLabel }}
            </div>
            <div v-if="lotNumberLabel" class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                {{ lotNumberLabel }}
            </div>
        </div>

        <!-- From → To - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="flex items-center gap-2 text-sm">
                <span class="text-(--color-text-secondary) truncate">
                    {{ fromCustomerName }}
                </span>
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary) shrink-0" fill="none" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
                <span class="font-semibold text-(--color-text-primary) truncate">
                    {{ toCustomerName }}
                </span>
            </div>
        </div>

        <!-- By / When - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
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
import type { LotTransfer } from '@/types/lotTransfer'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'

const props = defineProps<{
    transfer: LotTransfer
    deleting?: boolean
    isLatest?: boolean
}>()

const emit = defineEmits<{
    'delete': [transfer: LotTransfer]
}>()

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const usersStore = useUsersStore()

const lot = computed(() => lotsStore.getLotById(props.transfer.lot_id))

const lotLabel = computed(() =>
    lot.value ? lot.value.product_name : `Lot #${props.transfer.lot_id}`
)

const lotNumberLabel = computed(() =>
    lot.value ? `Lot ${lot.value.lot_number}` : ''
)

const fromCustomerName = computed(() =>
    customersStore.getCustomerName(props.transfer.from_customer_id)
)

const toCustomerName = computed(() =>
    customersStore.getCustomerName(props.transfer.to_customer_id)
)

const userName = computed(() => usersStore.getUserName(props.transfer.user_id))

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit',
    })

const handleDelete = () => {
    emit('delete', props.transfer)
}
</script>