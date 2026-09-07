<!-- src/components/features/customerStorageBills/CustomerStorageBillDetail.vue -->
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
                    <span class="text-sm text-(--color-text-secondary)">Lot #{{ bill.lot_id }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold"
                        :class="bill.outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                        {{ formatCurrency(bill.outstanding) }} outstanding
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">Lot ID: {{ bill.lot_id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Billing: {{ formatDate(bill.billing_start) }} {{
                bill.billing_end ? '→ ' + formatDate(bill.billing_end) : '→ Active' }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">{{ bill.months_billed }} month(s)</span>
        </div>

        <!-- Financial Summary -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Monthly Bill</p>
                <p class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(bill.monthly_bill) }}</p>
                <p class="text-xs text-(--color-text-secondary)">{{ bill.customer_charge_type }} × {{
                    bill.customer_storage_rate }}</p>
            </div>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Total Billed</p>
                <p class="text-lg font-bold text-(--color-text-primary)">{{ formatCurrency(bill.total_billed) }}</p>
                <p class="text-xs text-(--color-text-secondary)">{{ bill.months_billed }} months</p>
            </div>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Total Paid</p>
                <p class="text-lg font-bold text-(--color-green)">{{ formatCurrency(bill.total_paid) }}</p>
                <p class="text-xs text-(--color-text-secondary)">Last: {{ bill.last_paid_through ?
                    formatDate(bill.last_paid_through) : 'Never' }}</p>
            </div>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                <p class="text-lg font-bold"
                    :class="bill.outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(bill.outstanding) }}
                </p>
            </div>
        </div>

        <!-- Inventory Details -->
        <div class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">Inventory</p>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Quantity</p>
                    <p class="text-sm font-medium text-(--color-text-primary)">{{ bill.quantity }} {{ bill.quantity_unit
                        }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Weight</p>
                    <p class="text-sm font-medium text-(--color-text-primary)">{{ bill.weight }} {{ bill.weight_unit }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Charge Type</p>
                    <p class="text-sm font-medium capitalize text-(--color-text-primary)">{{ bill.customer_charge_type
                        }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary)">Storage Rate</p>
                    <p class="text-sm font-medium text-(--color-text-primary)">{{
                        formatCurrency(bill.customer_storage_rate) }}</p>
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
            <button @click="emit('pay', bill)"
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
import type { CustomerStorageBill } from '@/types/customerStorageBill'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerStorageBill | null
}>()

const emit = defineEmits<{
    'close': []
    'pay': [bill: CustomerStorageBill]
}>()

const formatDate = (dateStr: string | null): string => {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    })
}
</script>