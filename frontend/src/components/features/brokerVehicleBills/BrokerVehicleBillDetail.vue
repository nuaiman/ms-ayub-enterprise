<!-- src/components/features/brokerVehicleBills/BrokerVehicleBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ bill.vehicle_number }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ bill.broker_name }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">Transport #{{ bill.transport_id }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-blue)">{{ formatCurrency(bill.bill_amount)
                    }}</span>
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
            <span class="text-xs text-(--color-text-secondary)">Vehicle ID: {{ bill.vehicle_id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(bill.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(bill.updated_at) }}</span>
        </div>

        <!-- Bill Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
                <!-- Vehicle -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicle</p>
                    <p class="text-sm text-(--color-text-primary)">{{ bill.vehicle_number }}</p>
                </div>

                <!-- Broker -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Broker</p>
                    <p class="text-sm text-(--color-text-primary)">{{ bill.broker_name }}</p>
                </div>

                <!-- Transport -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Transport</p>
                    <p class="text-sm text-(--color-text-primary)">#{{ bill.transport_id }}</p>
                </div>
            </div>

            <div class="space-y-4">
                <!-- Joma Cost -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Joma Cost</p>
                    <p class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(bill.joma_cost) }}
                    </p>
                </div>

                <!-- Vehicle Cost -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicle Cost
                    </p>
                    <p class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(bill.vehicle_cost) }}
                    </p>
                </div>

                <!-- Bill Amount -->
                <div class="border-t border-(--color-border) pt-3">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Bill
                        Amount</p>
                    <p class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(bill.bill_amount) }}</p>
                    <p class="text-xs text-(--color-text-secondary">{{ formatCurrency(bill.joma_cost) }} + {{
                        formatCurrency(bill.vehicle_cost) }}</p>
                </div>

                <!-- Paid Amount -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Paid Amount
                    </p>
                    <p class="text-lg font-bold"
                        :class="bill.status === 'paid' ? 'text-(--color-green)' : 'text-(--color-red)'">
                        {{ formatCurrency(bill.paid_amount) }}
                    </p>
                </div>
            </div>
        </div>

        <!-- Notes -->
        <div v-if="bill.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ bill.notes }}</p>
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
import type { BrokerVehicleBill } from '@/types/brokerVehicleBill'
import { useBrokerVehicleBillsStore } from '@/stores/brokerVehicleBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: BrokerVehicleBill | null
}>()

const emit = defineEmits<{
    'close': []
    'pay': [bill: BrokerVehicleBill]
    'cancel': [bill: BrokerVehicleBill]
}>()

const brokerVehicleBillsStore = useBrokerVehicleBillsStore()

const getStatusBadgeClass = (status: string): string => {
    return brokerVehicleBillsStore.getStatusBadgeClass(status)
}

const getStatusDotClass = (status: string): string => {
    return brokerVehicleBillsStore.getStatusDotClass(status)
}

const getStatusLabel = (status: string): string => {
    return brokerVehicleBillsStore.getStatusLabel(status)
}

const formatDate = (dateStr: string | null): string => {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}
</script>