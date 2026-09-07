<!-- src/components/features/customerDeliveryBills/CustomerDeliveryBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ bill.customer_name }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ bill.item_name }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">Delivery #{{ bill.delivery_id }}</span>
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
            <span class="text-xs text-(--color-text-secondary)">Bill ID: {{ bill.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Delivery Item ID: {{ bill.delivery_item_id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Delivery Date: {{ formatDate(bill.delivery_date)
                }}</span>
        </div>

        <!-- Bill Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
                <!-- Customer -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ bill.customer_name }}</p>
                </div>

                <!-- Item -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Item</p>
                    <p class="text-sm text-(--color-text-primary)">{{ bill.item_name }}</p>
                </div>

                <!-- Delivery -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery</p>
                    <p class="text-sm text-(--color-text-primary)">#{{ bill.delivery_id }}</p>
                    <p class="text-xs text-(--color-text-secondary)">{{ formatDate(bill.delivery_date) }}</p>
                </div>

                <!-- Charge Type -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Charge Type
                    </p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ bill.customer_charge_type }}</p>
                </div>
            </div>

            <div class="space-y-4">
                <!-- Loading Rate -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Loading Rate
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.loading_rate) }}</p>
                </div>

                <!-- Quantity / Weight -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">
                        {{ bill.customer_charge_type === 'quantity' ? 'Quantity' : 'Weight' }}
                    </p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.customer_charge_type === 'quantity' ? bill.quantity : bill.weight }}
                        {{ bill.customer_charge_type === 'quantity' ? bill.quantity_unit : bill.weight_unit }}
                    </p>
                </div>

                <!-- Bill Amount -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Bill Amount
                    </p>
                    <p class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(bill.bill_amount) }}</p>
                    <p class="text-xs text-(--color-text-secondary)">{{ bill.loading_rate }} × {{
                        bill.customer_charge_type === 'quantity' ? bill.quantity : bill.weight }}</p>
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
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CustomerDeliveryBill } from '@/types/customerDeliveryBill'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerDeliveryBill | null
}>()

const emit = defineEmits<{
    'close': []
    'pay': [bill: CustomerDeliveryBill]
}>()

const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()

const getStatusBadgeClass = (status: string): string => {
    return customerDeliveryBillsStore.getStatusBadgeClass(status)
}

const getStatusDotClass = (status: string): string => {
    return customerDeliveryBillsStore.getStatusDotClass(status)
}

const getStatusLabel = (status: string): string => {
    return customerDeliveryBillsStore.getStatusLabel(status)
}

const formatDate = (dateStr: string | null): string => {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    })
}
</script>