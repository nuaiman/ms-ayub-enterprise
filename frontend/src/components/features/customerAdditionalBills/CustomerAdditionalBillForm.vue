<!-- src/components/features/customerAdditionalBills/CustomerAdditionalBillForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Additional Bill Information
            </h3>

            <!-- Create mode: pick customer -->
            <div v-if="!isEditMode" class="mb-4">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Customer <span class="text-(--color-red)">*</span>
                </label>
                <select v-model="form.customer_id" required :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option :value="null">Select a customer</option>
                    <option v-for="c in customerOptions" :key="c.id" :value="c.id">
                        {{ customerLabel(c) }}
                    </option>
                </select>
            </div>

            <!-- Edit mode: locked customer -->
            <div v-else class="mb-4 p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedCustomerLabel }}</p>
            </div>

            <!-- Amount -->
            <div class="mb-4">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Amount <span class="text-(--color-red)">*</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input v-model.number="form.amount" type="number" step="0.01" min="0.01" placeholder="0.00" required
                        :disabled="submitting"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <!-- Description -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Description <span class="text-(--color-red)">*</span>
                </label>
                <input v-model="form.description" type="text" placeholder="e.g. Late delivery penalty, extra handling"
                    required :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ isEditMode ? 'Saving...' : 'Creating...' }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Bill' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { CustomerAdditionalBill } from '@/types/customerAdditionalBill'
import type { Customer } from '@/types/customer'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { useCustomersStore } from '@/stores/customers'
import { push } from 'notivue'

const props = defineProps<{
    bill?: CustomerAdditionalBill | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'bill-created': []
    'bill-updated': []
    'cancel': []
}>()

const customerAdditionalBillsStore = useCustomerAdditionalBillsStore()
const customersStore = useCustomersStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.bill)

const customerOptions = computed(() => customersStore.customers)

const form = ref({
    customer_id: null as number | null,
    amount: null as number | null,
    description: '',
})

const lockedCustomerLabel = computed(() => {
    if (!props.bill) return ''
    return customersStore.getCustomerName(props.bill.customer_id)
})

const customerLabel = (c: Customer): string =>
    c.company_name || c.contact_person || `Customer #${c.id}`

const initialize = () => {
    if (props.bill) {
        form.value = {
            customer_id: props.bill.customer_id,
            amount: props.bill.amount,
            description: props.bill.description,
        }
    } else {
        form.value = {
            customer_id: null,
            amount: null,
            description: '',
        }
    }
}

watch(() => props.bill, initialize, { immediate: true })

const resetForm = () => { initialize() }

const submit = async () => {
    if (!isEditMode.value && !form.value.customer_id) {
        push.error('Please select a customer')
        return
    }
    if (!form.value.amount || form.value.amount <= 0) {
        push.error('Amount must be greater than 0')
        return
    }
    if (!form.value.description.trim()) {
        push.error('Description is required')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.bill) {
            const success = await customerAdditionalBillsStore.updateCustomerAdditionalBill(props.bill.id, {
                amount: form.value.amount,
                description: form.value.description.trim(),
            })
            if (success) {
                push.success('Customer additional bill updated successfully!')
                emit('bill-updated')
            }
        } else {
            const newBill = await customerAdditionalBillsStore.createCustomerAdditionalBill({
                customer_id: form.value.customer_id!,
                amount: form.value.amount,
                description: form.value.description.trim(),
            })
            if (newBill) {
                push.success('Customer additional bill created successfully!')
                resetForm()
                emit('bill-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update customer additional bill' : 'Failed to create customer additional bill')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
})
</script>