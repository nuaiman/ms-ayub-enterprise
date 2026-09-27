<!-- src/components/features/deliveries/DeliveryDetail.vue -->
<template>
    <div v-if="delivery" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Delivery</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDate(delivery.delivery_date)
                        }}</span>
                    <span v-if="delivery.from_location" class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span v-if="delivery.from_location" class="text-sm text-(--color-text-secondary)">{{
                        delivery.from_location }}{{ delivery.to_location ? `৳’ ${delivery.to_location}` : '' }}</span>
                </div>
            </div>
        </div>

        <!-- Summary strip -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Items</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ deliveryItems.length }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Quantity</p>
                <div class="mt-1">
                    <p v-if="totalQuantityByUnit.length === 0" class="text-lg font-bold text-(--color-text-primary)">৳
                    </p>
                    <p v-else v-for="(line, idx) in totalQuantityByUnit" :key="idx"
                        class="text-lg font-bold text-(--color-text-primary) leading-tight">
                        {{ line }}
                    </p>
                </div>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Weight</p>
                <div class="mt-1">
                    <p v-if="totalWeightByUnit.length === 0" class="text-lg font-bold text-(--color-text-primary)">৳
                    </p>
                    <p v-else v-for="(line, idx) in totalWeightByUnit" :key="idx"
                        class="text-lg font-bold text-(--color-text-primary) leading-tight">
                        {{ line }}
                    </p>
                </div>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Receiver</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1 truncate">
                    {{ delivery.receiver_name || '৳' }}
                </p>
                <p v-if="delivery.receiver_phone" class="text-xs text-(--color-text-secondary) mt-0.5 truncate">
                    {{ delivery.receiver_phone }}
                </p>
            </div>
        </div>

        <!-- Delivery information -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Delivery
                    Information</h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery Date
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(delivery.delivery_date) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">From</p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.from_location || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">To</p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.to_location || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Receiver</p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.receiver_name || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Receiver
                        Phone</p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.receiver_phone || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Recorded By
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ recordedBy }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(delivery.created_at) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Updated</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(delivery.updated_at) }}</p>
                </div>
                <div v-if="delivery.notes" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) mt-1">
                        <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ delivery.notes }}</p>
                    </div>
                </div>
                <div v-if="delivery.image_url" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Attachment
                    </p>
                    <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md mt-1">
                        <img :src="getImageUrl(delivery.image_url)" alt="Delivery attachment"
                            class="w-full object-cover max-h-64" />
                    </div>
                </div>
            </div>
        </section>

        <!-- Delivery items -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Delivery Items
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ deliveryItems.length }}</span>
            </div>

            <div v-if="deliveryItems.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No items in this delivery.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="di in deliveryItems" :key="di.id" class="p-4">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div class="min-w-0">
                            <div class="flex items-center gap-2 flex-wrap">
                                <span class="text-sm font-semibold text-(--color-text-primary)">{{
                                    getLotName(di.lot_id) }}</span>
                                <span class="text-xs text-(--color-text-secondary)">Lot #{{
                                    getLotNumber(di.lot_id) }}</span>
                            </div>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-xs text-(--color-text-secondary)">
                                <span>Quantity: <span class="text-(--color-text-primary)">{{ di.quantity }} {{
                                    di.quantity_unit }}</span></span>
                                <span>Weight: <span class="text-(--color-text-primary)">{{ di.weight }} {{
                                    di.weight_unit }}</span></span>
                                <span v-if="di.vehicle_number">Vehicle: <span class="text-(--color-text-primary)">{{
                                    di.vehicle_number }}</span></span>
                                <span v-if="di.driver_number">Driver: <span class="text-(--color-text-primary)">{{
                                    di.driver_number }}</span></span>
                                <span v-if="di.majhi_id">Majhi: <span class="text-(--color-text-primary)">{{
                                    getMajhiName(di.majhi_id) }}</span></span>
                            </div>
                            <div v-if="di.notes" class="text-xs text-(--color-text-secondary) mt-2 whitespace-pre-wrap">
                                {{ di.notes }}
                            </div>
                        </div>
                        <div class="text-right text-xs shrink-0">
                            <p class="text-(--color-text-secondary)">Loading Bill</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(deliveryBillAmount(di)) }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

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
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { formatCurrency } from '@/utils/currency'
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
const usersStore = useUsersStore()
const deliveryItemsStore = useDeliveryItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()

const customerName = computed(() => {
    if (!props.delivery?.customer_id) return '৳'
    return customersStore.getCustomerName(props.delivery.customer_id)
})

const recordedBy = computed(() => {
    if (!props.delivery) return '৳'
    return usersStore.getUserName(props.delivery.user_id)
})

const deliveryItems = computed(() => {
    if (!props.delivery) return []
    return deliveryItemsStore.getDeliveryItemsByDeliveryId(props.delivery.id)
})

const totalQuantityByUnit = computed<string[]>(() => {
    const byUnit: Record<string, number> = {}
    for (const di of deliveryItems.value) {
        if (di.quantity > 0) {
            const u = di.quantity_unit || 'units'
            byUnit[u] = (byUnit[u] || 0) + di.quantity
        }
    }
    return Object.entries(byUnit).map(([u, v]) => {
        const n = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(v)
        return `${n} ${u}`
    })
})

const totalWeightByUnit = computed<string[]>(() => {
    const byUnit: Record<string, number> = {}
    for (const di of deliveryItems.value) {
        if (di.weight > 0) {
            const u = di.weight_unit || 'kg'
            byUnit[u] = (byUnit[u] || 0) + di.weight
        }
    }
    return Object.entries(byUnit).map(([u, v]) => {
        const n = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(v)
        return `${n} ${u}`
    })
})

const getLotName = (id: number): string => lotsStore.getLotName(id)
const getLotNumber = (lotId: number): number | string => {
    const lot = lotsStore.getLotById(lotId)
    return lot ? lot.lot_number : '৳'
}
const getMajhiName = (id: number | null): string => {
    if (!id) return '৳'
    return majhisStore.getMajhiName(id)
}

const deliveryBillAmount = (di: any): number => {
    const rate = di.loading_rate || 0
    if (di.customer_charge_type === 'quantity') return rate * (di.quantity || 0)
    return rate * (di.weight || 0)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
}
</script>