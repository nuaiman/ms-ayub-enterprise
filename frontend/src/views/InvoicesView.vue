<!-- src/views/InvoicesView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Invoices</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ invoicesStore.totalInvoices }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Subtotal</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ formatCurrency(invoicesStore.totalSubtotal) }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Discounts</p>
                <p class="text-2xl font-bold text-(--color-red) mt-1">
                    − {{ formatCurrency(invoicesStore.totalDiscount) }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Grand Total</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">
                    {{ formatCurrency(invoicesStore.totalAmount) }}
                </p>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 min-h-0 mt-6">
            <InvoiceList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useInvoicesStore } from '@/stores/invoices'
import InvoiceList from '@/components/features/invoices/InvoiceList.vue'
import { formatCurrency } from '@/utils/currency'

const invoicesStore = useInvoicesStore()

onMounted(() => {
    invoicesStore.fetchInvoices()
})
</script>