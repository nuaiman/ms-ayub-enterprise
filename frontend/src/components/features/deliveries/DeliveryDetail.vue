<!-- src/components/features/deliveries/DeliveryDetail.vue -->
<template>
    <div v-if="delivery" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Delivery #{{ delivery.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ getCustomerName(delivery.customer_id)
                        }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDate(delivery.delivery_date) }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ delivery.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDateTime(delivery.created_at)
                }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDateTime(delivery.updated_at)
                }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Customer -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                <p class="text-sm text-(--color-text-primary)">{{ getCustomerName(delivery.customer_id) }}</p>
            </div>

            <!-- Delivery Date -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery Date</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatDateTime(delivery.delivery_date) }}</p>
            </div>

            <!-- Receiver Name -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Receiver</p>
                <p class="text-sm text-(--color-text-primary)">{{ delivery.receiver_name || '—' }}</p>
            </div>

            <!-- Receiver Phone -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Receiver Phone</p>
                <p class="text-sm text-(--color-text-primary)">{{ delivery.receiver_phone || '—' }}</p>
            </div>

            <!-- From -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">From</p>
                <p class="text-sm text-(--color-text-primary)">{{ delivery.from_location || '—' }}</p>
            </div>

            <!-- To -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">To</p>
                <p class="text-sm text-(--color-text-primary)">{{ delivery.to_location || '—' }}</p>
            </div>

            <!-- Notes -->
            <div v-if="delivery.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ delivery.notes }}</p>
                </div>
            </div>

            <!-- Image -->
            <div v-if="delivery.image_url" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Attachment</p>
                <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md">
                    <img :src="getImageUrl(delivery.image_url)" alt="Delivery attachment"
                        class="w-full object-cover max-h-64" />
                </div>
            </div>
        </div>

        <!-- Delivery Items -->
        <div class="border-t border-(--color-border) pt-4">
            <div class="flex items-center justify-between mb-3">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery Items</p>
                <span class="text-xs text-(--color-text-secondary)">{{ deliveryItems.length }} items</span>
            </div>

            <div v-if="deliveryItems.length === 0" class="text-center py-4 text-sm text-(--color-text-secondary)">
                No items in this delivery
            </div>

            <div v-else class="space-y-2">
                <div v-for="item in deliveryItems" :key="item.id"
                    class="grid grid-cols-12 gap-2 py-2 px-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border) text-sm">
                    <div class="col-span-3 font-medium text-(--color-text-primary) truncate">
                        {{ getItemName(item.item_id) }}
                    </div>
                    <div class="col-span-3 text-(--color-text-secondary) truncate">
                        Lot: {{ getLotName(item.lot_id) }}
                    </div>
                    <div class="col-span-2 text-(--color-text-secondary)">
                        {{ item.quantity }} {{ item.quantity_unit }}
                    </div>
                    <div class="col-span-2 text-(--color-text-secondary)">
                        {{ item.weight }} {{ item.weight_unit }}
                    </div>
                    <div class="col-span-2 text-(--color-text-secondary) truncate">
                        {{ getMajhiName(item.majhi_id) }}
                    </div>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', delivery)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('manage-items', delivery)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-green) text-white hover:opacity-90 transition-all duration-200">
                Manage Items
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
import type { Delivery } from '@/types/delivery'
import type { DeliveryItem } from '@/types/deliveryItem'
import { useCustomersStore } from '@/stores/customers'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    delivery: Delivery | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [delivery: Delivery]
    'manage-items': [delivery: Delivery]
    'updated': []
}>()

const customersStore = useCustomersStore()
const deliveryItemsStore = useDeliveryItemsStore()
const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()

const deliveryItems = computed(() => {
    if (!props.delivery) return []
    return deliveryItemsStore.getDeliveryItemsByDeliveryId(props.delivery.id)
})

const getCustomerName = (id: number | null): string => {
    if (!id) return '—'
    return customersStore.getCustomerName(id)
}

const getItemName = (id: number): string => {
    return itemsStore.getItemName(id)
}

const getLotName = (id: number): string => {
    return lotsStore.getLotName(id)
}

const getMajhiName = (id: number | null): string => {
    if (!id) return '—'
    return majhisStore.getMajhiName(id)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'long',
        day: 'numeric',
        year: 'numeric'
    })
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