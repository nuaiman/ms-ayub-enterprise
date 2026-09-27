<!-- src/components/features/customerTransportBills/CustomerTransportBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">
                    Transport #{{ bill.id }}
                </h2>

                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">
                        {{ bill.customer_name }}
                    </span>

                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>

                    <span class="text-sm text-(--color-text-secondary)">
                        {{ bill.from_location }}
                        <span v-if="bill.to_location">
                            ৳ {{ bill.to_location }}
                        </span>
                    </span>

                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>

                    <span class="text-sm font-semibold text-(--color-blue)">
                        {{ formatCurrency(bill.bill_amount) }}
                    </span>

                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>

                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="getStatusBadgeClass(bill.status)">
                        <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(bill.status)"></span>
                        {{ getStatusLabel(bill.status) }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">
                Transport ID: {{ bill.id }}
            </span>

            <span class="w-px h-4 bg-(--color-border)"></span>

            <span class="text-xs text-(--color-text-secondary)">
                Vehicles: {{ bill.total_vehicles }}
            </span>

            <span class="w-px h-4 bg-(--color-border)"></span>

            <span class="text-xs text-(--color-text-secondary)">
                Created: {{ formatDate(bill.created_at) }}
            </span>

            <span class="w-px h-4 bg-(--color-border)"></span>

            <span class="text-xs text-(--color-text-secondary)">
                Updated: {{ formatDate(bill.updated_at) }}
            </span>
        </div>

        <!-- Financial Summary -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Charge Unit</p>
                <p class="text-lg font-bold text-(--color-text-primary) capitalize">
                    {{ bill.customer_charge_unit }}
                </p>
                <p class="text-xs text-(--color-text-secondary) mt-1">
                    {{ formatCurrency(bill.customer_charge_per_unit) }} per unit ৳‚৳· {{ bill.customer_total_unit }}
                    total
                </p>
            </div>

            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Total Bill</p>
                <p class="text-lg font-bold text-(--color-blue)">
                    {{ formatCurrency(bill.bill_amount) }}
                </p>
                <p class="text-xs text-(--color-text-secondary) mt-1">
                    {{ bill.vehicle_quantity }} vehicle(s)
                </p>
            </div>

            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Total Paid</p>
                <p class="text-lg font-bold"
                    :class="bill.status === 'paid' ? 'text-(--color-green)' : 'text-(--color-red)'">
                    {{ formatCurrency(bill.paid_amount) }}
                </p>
                <p class="text-xs text-(--color-text-secondary) mt-1">
                    {{ bill.payment_date ? `Paid through ${formatDate(bill.payment_date)}` : 'Not paid yet' }}
                </p>
            </div>
        </div>

        <!-- Vehicle Breakdown -->
        <div class="border-t border-(--color-border) pt-4">
            <div class="flex items-center justify-between mb-3">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">
                    Vehicle Breakdown
                </p>

                <span class="text-xs text-(--color-text-secondary)">
                    {{ bill.vehicles.length }} vehicles
                </span>
            </div>

            <div class="space-y-2">
                <!-- Header -->
                <div
                    class="grid grid-cols-12 gap-2 py-2 px-3 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">
                    <div class="col-span-4">Vehicle</div>
                    <div class="col-span-4">Broker</div>
                    <div class="col-span-2 text-right">Joma</div>
                    <div class="col-span-2 text-right">Vehicle Cost</div>
                </div>

                <!-- Rows -->
                <div v-for="vehicle in bill.vehicles" :key="vehicle.vehicle_id"
                    class="grid grid-cols-12 gap-2 py-2 px-3 rounded-lg bg-(--color-muted-bg)/20 border border-(--color-border) text-sm hover:bg-(--color-muted-bg)/40 transition-colors">
                    <div class="col-span-4 font-medium text-(--color-text-primary) truncate">
                        {{ vehicle.vehicle_number }}
                    </div>

                    <div class="col-span-4 text-(--color-text-secondary) truncate">
                        {{ vehicle.broker_name || '৳' }}
                    </div>

                    <div class="col-span-2 text-right text-(--color-text-secondary)">
                        {{ formatCurrency(vehicle.joma_cost) }}
                    </div>

                    <div class="col-span-2 text-right text-(--color-text-secondary)">
                        {{ formatCurrency(vehicle.vehicle_cost) }}
                    </div>
                </div>

                <!-- Total Row -->
                <div
                    class="grid grid-cols-12 gap-2 py-2 px-3 rounded-lg bg-(--color-blue)/5 border border-(--color-blue)/20 text-sm font-semibold">
                    <div class="col-span-9 text-right text-(--color-text-secondary)">
                        Total Bill
                    </div>

                    <div class="col-span-3 text-right text-(--color-blue)">
                        {{ formatCurrency(bill.bill_amount) }}
                    </div>
                </div>
            </div>
        </div>

        <!-- Notes -->
        <div v-if="bill.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">
                Notes
            </p>

            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">
                    {{ bill.notes }}
                </p>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button v-if="bill.status === 'unpaid'" @click="emit('pay', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-green) text-white hover:opacity-90 transition-all duration-200">
                Record Payment
            </button>

            <button v-if="bill.status === 'unpaid'" @click="emit('cancel', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-yellow) text-white hover:opacity-90 transition-all duration-200">
                Cancel Bill
            </button>

            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CustomerTransportBill } from '@/types/customerTransportBill'
import { useCustomerTransportBillsStore } from '@/stores/customerTransportBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerTransportBill | null
}>()

const emit = defineEmits<{
    'close': []
    'pay': [bill: CustomerTransportBill]
    'cancel': [bill: CustomerTransportBill]
}>()

const store = useCustomerTransportBillsStore()

const getStatusBadgeClass = (status: string): string => {
    return store.getStatusBadgeClass(status)
}

const getStatusDotClass = (status: string): string => {
    return store.getStatusDotClass(status)
}

const getStatusLabel = (status: string): string => {
    return store.getStatusLabel(status)
}

const formatDate = (dateStr: string | null): string => {
    if (!dateStr) return '৳'
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}
</script>