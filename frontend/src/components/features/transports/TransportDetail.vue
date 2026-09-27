<!-- src/components/features/transports/TransportDetail.vue -->
<template>
    <div v-if="transport" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Transport #{{ transport.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
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

        <!-- Transport Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Transport Type
                </p>
                <p class="text-sm text-(--color-text-primary) capitalize">{{ transport.transport_type || '৳' }}
                </p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">From</p>
                <p class="text-sm text-(--color-text-primary)">{{ transport.from_location }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">To</p>
                <p class="text-sm text-(--color-text-primary)">{{ transport.to_location || '৳' }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicles</p>
                <p class="text-sm text-(--color-text-primary)">{{ transport.vehicle_quantity }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Office
                    Commission</p>
                <p class="text-lg font-semibold text-(--color-text-primary)">{{
                    formatCurrency(transport.office_commission_amount) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Transport Date</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatDateTime(transport.transport_date) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created By</p>
                <p class="text-sm text-(--color-text-primary)">{{ getUserName(transport.user_id) }}</p>
            </div>
        </div>

        <!-- Customer Billing -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Customer
                    Billing</h3>
            </div>
            <div class="p-4 grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Charge Unit</p>
                    <p class="text-sm font-medium text-(--color-text-primary) capitalize">
                        {{ transport.customer_charge_unit }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Total Unit</p>
                    <p class="text-sm font-medium text-(--color-text-primary)">
                        {{ transport.customer_total_unit }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Per Unit</p>
                    <p class="text-sm font-medium text-(--color-text-primary)">
                        {{ formatCurrency(transport.customer_charge_per_unit) }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Total Charge</p>
                    <p class="text-sm font-bold text-(--color-blue)">
                        {{ formatCurrency(transport.customer_total_charge) }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Total Paid</p>
                    <p class="text-sm font-bold text-(--color-green)">
                        {{ formatCurrency(transport.customer_total_paid) }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Paid Through</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ transport.customer_total_paid_through ? formatDate(transport.customer_total_paid_through) :
                            '৳' }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                    <p class="text-sm font-bold"
                        :class="customerOutstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                        {{ formatCurrency(customerOutstanding) }}
                    </p>
                </div>
            </div>
        </section>

        <!-- Broker Payment -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Broker Payment
                </h3>
            </div>
            <div class="p-4 grid grid-cols-2 md:grid-cols-3 gap-4">
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Bill Amount</p>
                    <p class="text-sm font-bold text-(--color-blue)">{{ formatCurrency(brokerBillAmount) }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Total Paid to Broker</p>
                    <p class="text-sm font-bold text-(--color-green)">
                        {{ formatCurrency(totalPaidToBroker) }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                    <p class="text-sm font-bold"
                        :class="brokerOutstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                        {{ formatCurrency(brokerOutstanding) }}
                    </p>
                </div>
            </div>
        </section>

        <!-- Notes -->
        <div v-if="transport.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ transport.notes }}</p>
            </div>
        </div>

        <!-- Image -->
        <div v-if="transport.image_url" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Attachment</p>
            <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md">
                <img :src="getImageUrl(transport.image_url)" alt="Transport attachment"
                    class="w-full object-cover max-h-64" />
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', transport)"
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
import type { Transport } from '@/types/transport'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import { useVehiclesStore } from '@/stores/vehicles'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    transport: Transport | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [transport: Transport]
    'updated': []
}>()

const customersStore = useCustomersStore()
const usersStore = useUsersStore()
const vehiclesStore = useVehiclesStore()

const customerName = computed(() => {
    if (!props.transport?.customer_id) return '৳'
    return customersStore.getCustomerName(props.transport.customer_id)
})

const getUserName = (id: number): string => usersStore.getUserName(id)

const brokerBillAmount = computed(() => {
    if (!props.transport) return 0
    const vehicles = vehiclesStore.getVehiclesByTransportId(props.transport.id)
    return vehicles.reduce((sum, v) => sum + v.joma_cost + v.vehicle_cost, 0)
})

const totalPaidToBroker = computed(() => {
    if (!props.transport) return 0
    const vehicles = vehiclesStore.getVehiclesByTransportId(props.transport.id)
    return vehicles.reduce((sum, v) => sum + (v.total_paid_to_broker || 0), 0)
})

const brokerOutstanding = computed(() => {
    return Math.max(0, brokerBillAmount.value - totalPaidToBroker.value)
})

const customerOutstanding = computed(() => {
    if (!props.transport) return 0
    return Math.max(0, (props.transport.customer_total_charge || 0) - (props.transport.customer_total_paid || 0))
})

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'long',
        day: 'numeric',
        year: 'numeric',
    })
}

const formatDateTime = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}
</script>