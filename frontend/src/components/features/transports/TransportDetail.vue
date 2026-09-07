<!-- src/components/features/transports/TransportDetail.vue -->
<template>
    <div v-if="transport" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Transport #{{ transport.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ getCustomerName(transport.customer_id)
                    }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ transport.from_location }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDate(transport.transport_date)
                    }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ transport.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDateTime(transport.created_at)
            }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDateTime(transport.updated_at)
            }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Customer -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                <p class="text-sm text-(--color-text-primary)">{{ getCustomerName(transport.customer_id) }}</p>
            </div>

            <!-- From -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">From</p>
                <p class="text-sm text-(--color-text-primary)">{{ transport.from_location }}</p>
            </div>

            <!-- To -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">To</p>
                <p class="text-sm text-(--color-text-primary)">{{ transport.to_location || '—' }}</p>
            </div>

            <!-- Delivery Type -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery Type</p>
                <p class="text-sm text-(--color-text-primary) capitalize">{{ transport.delivery_type || '—' }}</p>
            </div>

            <!-- Vehicle Quantity -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicles</p>
                <p class="text-sm text-(--color-text-primary)">{{ transport.vehicle_quantity }}</p>
            </div>

            <!-- Office Commission -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Commission</p>
                <p class="text-lg font-semibold text-(--color-text-primary)">{{
                    formatCurrency(transport.office_commission_amount) }}</p>
            </div>

            <!-- Customer Total Paid (NEW) -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer Total
                    Paid</p>
                <p class="text-lg font-semibold text-(--color-blue)">{{ formatCurrency(transport.customer_total_paid ||
                    0) }}</p>
            </div>

            <!-- Transport Date -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Transport Date</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatDateTime(transport.transport_date) }}</p>
            </div>

            <!-- Created By -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created By</p>
                <p class="text-sm text-(--color-text-primary)">{{ getUserName(transport.user_id) }}</p>
            </div>

            <!-- Notes -->
            <div v-if="transport.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ transport.notes }}</p>
                </div>
            </div>

            <!-- Image -->
            <div v-if="transport.image_url" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Attachment</p>
                <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md">
                    <img :src="getImageUrl(transport.image_url)" alt="Transport attachment"
                        class="w-full object-cover max-h-64" />
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', transport)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('manage-vehicles', transport)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-green) text-white hover:opacity-90 transition-all duration-200">
                Manage Vehicles
            </button>
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Transport } from '@/types/transport'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    transport: Transport | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [transport: Transport]
    'manage-vehicles': [transport: Transport]
    'updated': []
}>()

const customersStore = useCustomersStore()
const usersStore = useUsersStore()

const getCustomerName = (id: number | null): string => {
    if (!id) return '—'
    return customersStore.getCustomerName(id)
}

const getUserName = (id: number): string => {
    return usersStore.getUserName(id)
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