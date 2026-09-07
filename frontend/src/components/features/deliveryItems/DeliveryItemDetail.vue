<!-- src/components/features/deliveryItems/DeliveryItemDetail.vue -->
<template>
    <div v-if="deliveryItem" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Delivery Item #{{ deliveryItem.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">Delivery #{{ deliveryItem.delivery_id }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ getItemName(deliveryItem.item_id) }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ deliveryItem.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDateTime(deliveryItem.created_at)
                }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDateTime(deliveryItem.updated_at)
                }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery</p>
                <p class="text-sm text-(--color-text-primary)">#{{ deliveryItem.delivery_id }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                <p class="text-sm text-(--color-text-primary)">{{ getStoreDisplayName(store) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Item</p>
                <p class="text-sm text-(--color-text-primary)">{{ getItemName(deliveryItem.item_id) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
                <p class="text-sm text-(--color-text-primary)">{{ getLotName(deliveryItem.lot_id) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.quantity }} {{ deliveryItem.quantity_unit
                    }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.weight }} {{ deliveryItem.weight_unit }}
                </p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                <p class="text-sm text-(--color-text-primary)">{{ getMajhiName(deliveryItem.majhi_id) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Loading Rate</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(deliveryItem.loading_rate) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi Cut</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(deliveryItem.majhi_cut) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicle Number</p>
                <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.vehicle_number || '—' }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Driver Number</p>
                <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.driver_number || '—' }}</p>
            </div>

            <!-- Billing Info (read-only) -->
            <div class="md:col-span-2 border-t border-(--color-border) pt-4 mt-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">Billing Info
                    (Snapshot at Delivery)</p>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Customer Charge Type</p>
                        <p class="text-sm text-(--color-text-primary)">{{
                            getChargeTypeLabel(deliveryItem.customer_charge_type) }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Majhi Bill Type</p>
                        <p class="text-sm text-(--color-text-primary)">{{
                            getMajhiBillTypeLabel(deliveryItem.majhi_bill_type) }}</p>
                    </div>
                </div>
            </div>

            <div v-if="deliveryItem.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ deliveryItem.notes }}</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', deliveryItem)"
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
import type { DeliveryItem } from '@/types/deliveryItem'
import type { Store } from '@/types/store'
import { useStoresStore } from '@/stores/stores'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    deliveryItem: DeliveryItem | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [deliveryItem: DeliveryItem]
    'updated': []
}>()

const storesStore = useStoresStore()
const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const deliveryItemsStore = useDeliveryItemsStore()

const store = computed(() => {
    if (!props.deliveryItem) return undefined
    return storesStore.getStoreById(props.deliveryItem.store_id)
})

const getStoreDisplayName = (store: Store | undefined): string => {
    if (!store) return `Store #${props.deliveryItem?.store_id || 'Unknown'}`
    return storesStore.getStoreDisplayName(store)
}

const getItemName = (id: number): string => itemsStore.getItemName(id)
const getLotName = (id: number): string => lotsStore.getLotName(id)

const getMajhiName = (id: number | null): string => {
    if (!id) return '—'
    return majhisStore.getMajhiName(id)
}

const getChargeTypeLabel = (chargeType: 'weight' | 'quantity'): string => {
    return deliveryItemsStore.getChargeTypeLabel(chargeType)
}

const getMajhiBillTypeLabel = (billType: 'weight' | 'quantity' | 'job'): string => {
    return deliveryItemsStore.getMajhiBillTypeLabel(billType)
}

const formatDateTime = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}
</script>