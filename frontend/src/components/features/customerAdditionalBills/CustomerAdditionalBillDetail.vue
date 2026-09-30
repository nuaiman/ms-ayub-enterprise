<!-- src/components/features/customerAdditionalBills/CustomerAdditionalBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ bill.customer_name }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm font-semibold text-(--color-blue)">{{ formatCurrency(bill.amount) }}</span>
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
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(bill.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(bill.updated_at) }}</span>
        </div>

        <!-- Bill Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ bill.customer_name }}</p>
                </div>

                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Description
                    </p>
                    <p class="text-sm text-(--color-text-primary) whitespace-pre-wrap">{{ bill.description }}</p>
                </div>
            </div>

            <div class="space-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                    <p class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(bill.amount) }}</p>
                </div>

                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Paid</p>
                    <p class="text-lg font-bold text-(--color-green)">{{ formatCurrency(bill.paid_amount) }}</p>
                </div>

                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Outstanding
                    </p>
                    <p class="text-lg font-bold"
                        :class="bill.outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                        {{ formatCurrency(bill.outstanding) }}
                    </p>
                </div>

                <div v-if="bill.payment_date">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Paid Through
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(bill.payment_date) }}</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button v-if="bill.status === 'unpaid'" @click="emit('pay', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-green) text-white hover:opacity-90 transition-all duration-200">
                Record Payment
            </button>
            <button @click="emit('edit', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('delete', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-all duration-200">
                Delete
            </button>
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CustomerAdditionalBill } from '@/types/customerAdditionalBill'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { formatCurrency } from '@/utils/currency'

defineProps<{
    bill: CustomerAdditionalBill | null
}>()

const emit = defineEmits<{
    'close': []
    'pay': [bill: CustomerAdditionalBill]
    'edit': [bill: CustomerAdditionalBill]
    'delete': [bill: CustomerAdditionalBill]
}>()

const store = useCustomerAdditionalBillsStore()

const getStatusBadgeClass = (status: string): string => store.getStatusBadgeClass(status)
const getStatusDotClass = (status: string): string => store.getStatusDotClass(status)
const getStatusLabel = (status: string): string => store.getStatusLabel(status)

const formatDate = (dateStr: string | null): string => {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}
</script>