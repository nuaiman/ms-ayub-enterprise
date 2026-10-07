<!-- src/components/features/invoices/InvoiceDetail.vue -->
<template>
    <div v-if="invoice" class="space-y-6">
        <!-- Header -->
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Invoice #{{ invoice.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span
                        class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize bg-(--color-blue)/10 text-(--color-blue)">
                        {{ invoice.entity_type }}
                    </span>
                    <span class="text-sm text-(--color-text-secondary)">{{ entityLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(invoice.total)
                    }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ invoice.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(invoice.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(invoice.updated_at) }}</span>
        </div>

        <!-- Summary cards -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Subtotal</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ formatCurrency(invoice.subtotal) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Discount</p>
                <p class="text-lg font-bold text-(--color-red) mt-1">− {{ formatCurrency(invoice.discount_amount) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total</p>
                <p class="text-lg font-bold text-(--color-blue) mt-1">{{ formatCurrency(invoice.total) }}</p>
            </div>
        </div>

        <!-- Bills Covered -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Bills Covered
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ invoice.items.length }}</span>
            </div>

            <div v-if="invoice.items.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No bills recorded on this invoice.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="it in invoice.items" :key="it.id" class="p-4 flex items-start justify-between gap-3">
                    <div class="min-w-0">
                        <div class="text-sm font-semibold text-(--color-text-primary) capitalize">
                            {{ billTypeLabel(it.bill_type) }}
                        </div>
                        <div class="text-xs text-(--color-text-secondary) mt-0.5">
                            {{ billReference(it) }}
                        </div>
                    </div>
                    <div class="text-sm font-semibold text-(--color-text-primary) shrink-0">
                        {{ formatCurrency(it.amount) }}
                    </div>
                </div>
            </div>
        </section>

        <!-- Discounts -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Discounts</h3>
                <span class="text-xs text-(--color-text-secondary)">{{ invoice.discounts.length }}</span>
            </div>

            <div v-if="invoice.discounts.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No discounts applied to this invoice.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="d in invoice.discounts" :key="d.id" class="p-4 flex items-start justify-between gap-3">
                    <div class="min-w-0">
                        <div class="text-sm font-semibold text-(--color-text-primary) capitalize">
                            {{ d.type }} discount
                            <span class="text-(--color-text-secondary) font-normal">
                                ({{ d.type === 'percent' ? d.value + '%' : formatCurrency(d.value) }})
                            </span>
                        </div>
                        <div v-if="d.reason" class="text-xs text-(--color-text-secondary) mt-0.5">
                            {{ d.reason }}
                        </div>
                    </div>
                    <div class="text-sm font-semibold text-(--color-red) shrink-0">
                        − {{ formatCurrency(d.computed_amount) }}
                    </div>
                </div>
            </div>
        </section>

        <!-- Notes -->
        <div v-if="invoice.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ invoice.notes }}</p>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('print', invoice)"
                class="px-4 py-2 text-sm font-medium rounded-lg border border-(--color-border) hover:bg-(--color-muted-bg) transition-all duration-200 flex items-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                </svg>
                Print
            </button>
            <button @click="emit('edit', invoice)"
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
import type { InvoiceDetail as InvoiceDetailType, InvoiceItem, UnbilledBillType } from '@/types/invoice'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useBrokersStore } from '@/stores/brokers'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    invoice: InvoiceDetailType | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [invoice: InvoiceDetailType]
    'print': [invoice: InvoiceDetailType]
    'updated': []
}>()

const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const brokersStore = useBrokersStore()

const entityLabel = computed(() => {
    if (!props.invoice) return ''
    const inv = props.invoice
    switch (inv.entity_type) {
        case 'customer': {
            const c = customersStore.getCustomerById(inv.entity_id)
            return c ? (c.company_name || c.contact_person || `Customer #${c.id}`) : `Customer #${inv.entity_id}`
        }
        case 'majhi':
            return majhisStore.getMajhiName(inv.entity_id)
        case 'godown':
            return godownsStore.getGodownName(inv.entity_id)
        case 'broker':
            return brokersStore.getBrokerName(inv.entity_id)
        default:
            return `#${inv.entity_id}`
    }
})

const billTypeLabel = (t: UnbilledBillType): string => {
    const map: Record<UnbilledBillType, string> = {
        customer_store: 'Customer Store Bill',
        customer_delivery: 'Customer Delivery Bill',
        customer_additional: 'Customer Additional Bill',
        majhi: 'Majhi Bill',
        godown: 'Godown Bill',
    }
    return map[t] || t
}

const billReference = (it: InvoiceItem): string => {
    return `#${it.bill_id}`
}

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
</script>