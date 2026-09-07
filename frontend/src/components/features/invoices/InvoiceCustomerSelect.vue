<!-- src/components/features/invoices/InvoiceCustomerSelect.vue -->
<template>
    <div class="space-y-6">
        <div class="p-6 rounded-xl bg-(--color-muted-bg)/30 border border-(--color-border)">
            <div class="flex items-start gap-3 mb-4">
                <div
                    class="w-10 h-10 rounded-xl bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-base font-semibold text-(--color-text-primary)">Select Customer & Invoice Type</h2>
                    <p class="text-sm text-(--color-text-secondary)">Choose the customer and type of invoice</p>
                </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Customer <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="localCustomerId"
                        class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200 appearance-none">
                        <option :value="null">Select a customer...</option>
                        <option v-for="customer in customersStore.customers" :key="customer.id" :value="customer.id">
                            {{ getCustomerDisplayName(customer) }}
                        </option>
                    </select>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Invoice Type <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="localInvoiceType"
                        class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200 appearance-none">
                        <option value="godown">Godown Invoice</option>
                        <option value="transport">Transport Invoice</option>
                    </select>
                </div>
            </div>

            <div v-if="localCustomerId" class="mt-4 pt-4 border-t border-(--color-border)">
                <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm">
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Customer</p>
                        <p class="font-medium text-(--color-text-primary)">{{ getCustomerDisplayName(selectedCustomer)
                            }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Phone</p>
                        <p class="font-medium text-(--color-text-primary)">{{ selectedCustomer?.phone || 'N/A' }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Email</p>
                        <p class="font-medium text-(--color-text-primary)">{{ selectedCustomer?.email || 'N/A' }}</p>
                    </div>
                </div>
                <p class="text-xs text-(--color-text-secondary) mt-3">
                    All unpaid and partially paid bills will be available for selection
                </p>
            </div>
        </div>

        <div class="flex justify-end">
            <button @click="handleNext" :disabled="!localCustomerId"
                class="px-6 py-2.5 bg-(--color-blue) text-white rounded-xl text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100 flex items-center gap-2">
                Next Step
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useCustomersStore } from '@/stores/customers'
import type { Customer } from '@/types/customer'
import type { InvoiceType } from '@/types/invoice'

const props = defineProps<{
    customerId: number | null
    invoiceType: InvoiceType
}>()

const emit = defineEmits<{
    'update:customerId': [value: number | null]
    'update:invoiceType': [value: InvoiceType]
    'next': []
}>()

const customersStore = useCustomersStore()

const localCustomerId = ref(props.customerId)
const localInvoiceType = ref(props.invoiceType)

const selectedCustomer = computed<Customer | null>(() => {
    if (!localCustomerId.value) return null
    const customer = customersStore.getCustomerById(localCustomerId.value)
    return customer || null
})

watch(() => localCustomerId.value, (val) => {
    emit('update:customerId', val)
})

watch(() => localInvoiceType.value, (val) => {
    emit('update:invoiceType', val)
})

const getCustomerDisplayName = (customer: Customer | null): string => {
    if (!customer) return 'Unknown'
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}

const handleNext = () => {
    if (localCustomerId.value) {
        emit('next')
    }
}
</script>