<!-- src/components/features/invoices/InvoiceRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Entity - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="flex items-center gap-2">
                <span
                    class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize bg-(--color-blue)/10 text-(--color-blue)">
                    {{ invoice.entity_type }}
                </span>
                <span class="font-medium text-(--color-text-primary) truncate text-sm">
                    {{ entityLabel }}
                </span>
            </div>
            <div v-if="invoice.notes" class="text-xs text-(--color-text-secondary) truncate mt-1">
                {{ invoice.notes }}
            </div>
        </div>

        <!-- Subtotal - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) block">
                {{ formatCurrency(invoice.subtotal) }}
            </span>
        </div>

        <!-- Discount - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm font-semibold text-(--color-red) block"
                :class="{ 'text-(--color-text-secondary)': invoice.discount_amount === 0 }">
                {{ invoice.discount_amount > 0 ? '− ' + formatCurrency(invoice.discount_amount) : '—' }}
            </span>
        </div>

        <!-- Total - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <span class="text-sm font-bold text-(--color-text-primary) block">
                {{ formatCurrency(invoice.total) }}
            </span>
            <span class="text-xs text-(--color-text-secondary)/70 block mt-0.5">
                {{ formatDateShort(invoice.created_at) }}
            </span>
        </div>

        <!-- Actions - 1 column -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
                    <button @click="handleView"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                        </svg>
                        View Details
                    </button>

                    <button @click="handlePrint"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                        </svg>
                        Print
                    </button>

                    <button @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>

                    <button @click="handleDelete"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-red) hover:bg-(--color-muted-bg) transition-colors border-t border-(--color-border) mt-1 pt-1">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        Delete
                    </button>
                </div>
            </Transition>

            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Invoice } from '@/types/invoice'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useBrokersStore } from '@/stores/brokers'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    invoice: Invoice
}>()

const emit = defineEmits<{
    'view': [invoice: Invoice]
    'print': [invoice: Invoice]
    'edit': [invoice: Invoice]
    'delete': [invoice: Invoice]
    'updated': []
}>()

const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const brokersStore = useBrokersStore()
const isOpen = ref(false)

const entityLabel = computed(() => {
    switch (props.invoice.entity_type) {
        case 'customer': {
            const c = customersStore.getCustomerById(props.invoice.entity_id)
            return c ? (c.company_name || c.contact_person || `Customer #${c.id}`) : `Customer #${props.invoice.entity_id}`
        }
        case 'majhi':
            return majhisStore.getMajhiName(props.invoice.entity_id)
        case 'godown':
            return godownsStore.getGodownName(props.invoice.entity_id)
        case 'broker':
            return brokersStore.getBrokerName(props.invoice.entity_id)
        default:
            return `#${props.invoice.entity_id}`
    }
})

const formatDateShort = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.invoice) }
const handlePrint = () => { closeMenu(); emit('print', props.invoice) }
const handleEdit = () => { closeMenu(); emit('edit', props.invoice) }
const handleDelete = () => { closeMenu(); emit('delete', props.invoice) }
</script>