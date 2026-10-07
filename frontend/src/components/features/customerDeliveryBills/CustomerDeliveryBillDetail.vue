<!-- src/components/features/customerDeliveryBills/CustomerDeliveryBillDetail.vue -->
<template>
    <div v-if="bill" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Customer Delivery Bill #{{ bill.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ customerLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">Item #{{ bill.delivery_item_id }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary) capitalize">{{ bill.bill_type }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(bill.rate)
                        }}</span>
                </div>
            </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ bill.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(bill.created_at) }}</span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1 capitalize">{{ bill.bill_type }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Rate</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ formatCurrency(bill.rate) }}</p>
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
        </div>

        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Bill Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerLabel }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery Item
                    </p>
                    <p class="text-sm text-(--color-text-primary)">Item #{{ bill.delivery_item_id }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ bill.bill_type }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Rate</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.rate) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight at
                        Billing</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.weight_at_billing }} {{ bill.weight_unit_at_billing || '' }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity at
                        Billing</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.quantity_at_billing }} {{ bill.quantity_unit_at_billing || '' }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Amount
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.total_amount) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(bill.total_paid) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid
                        Through</p>
                    <p class="text-sm text-(--color-text-primary)">
                        {{ bill.total_paid_through ? formatDateShort(bill.total_paid_through) : '—' }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(bill.created_at) }}</p>
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
import type { CustomerDeliveryBill } from '@/types/customerDeliveryBill'
import { useCustomersStore } from '@/stores/customers'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    bill: CustomerDeliveryBill | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [bill: CustomerDeliveryBill]
    'payment': [bill: CustomerDeliveryBill]
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