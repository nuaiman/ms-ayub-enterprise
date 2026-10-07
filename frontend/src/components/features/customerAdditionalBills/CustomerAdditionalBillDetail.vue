<!-- src/components/features/customerAdditionalBills/CustomerAdditionalBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Additional Bill #{{ bill.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ customerLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ bill.description }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(bill.amount)
                        }}</span>
                </div>
            </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ bill.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(bill.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(bill.updated_at) }}</span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ formatCurrency(bill.amount) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                <p class="text-lg font-bold text-(--color-green) mt-1">{{ formatCurrency(bill.total_paid) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                <p class="text-lg font-bold mt-1"
                    :class="bill.remaining > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(bill.remaining) }}
                </p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Paid Through</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                    {{ bill.total_paid_through ? formatDateShort(bill.total_paid_through) : '—' }}
                </p>
            </div>
        </div>

        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">
                    Bill Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerLabel }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Description
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ bill.description }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.amount) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.total_paid) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Remaining</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.remaining) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid
                        Through</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.total_paid_through ? formatDateShort(bill.total_paid_through) : '—' }}
                    </p>
                </div>
            </div>
        </section>

        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('payment', bill)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-green) text-white hover:opacity-90 transition-all duration-200">
                Add Payment
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
import type { CustomerAdditionalBill } from '@/types/customerAdditionalBill'
import { useCustomersStore } from '@/stores/customers'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerAdditionalBill | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [bill: CustomerAdditionalBill]
    'payment': [bill: CustomerAdditionalBill]
    'updated': []
}>()

const customersStore = useCustomersStore()

const customerLabel = computed(() => {
    if (!props.bill) return '—'
    return customersStore.getCustomerName(props.bill.customer_id)
})

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })

const formatDateShort = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })
</script>