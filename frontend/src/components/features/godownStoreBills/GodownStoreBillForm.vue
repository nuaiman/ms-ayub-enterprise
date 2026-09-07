<!-- src/components/features/godownStoreBills/GodownStoreBillForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Update Payment Information
            </h3>

            <!-- Store Info (read-only) -->
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border) mb-4">
                <div class="grid grid-cols-2 gap-4 text-sm">
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Store</p>
                        <p class="font-medium text-(--color-text-primary)">{{ store?.lot_name }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Monthly Bill</p>
                        <p class="font-medium text-(--color-blue)">{{ formatCurrency(store?.monthly_bill) }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Current Paid</p>
                        <p class="font-medium text-(--color-green)">{{ formatCurrency(store?.total_billed) }}</p>
                    </div>
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                        <p class="font-medium text-(--color-red)">{{ formatCurrency(store?.outstanding) }}</p>
                    </div>
                </div>
            </div>

            <!-- Last Paid Amount -->
            <div class="mb-4">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Last Paid Amount <span class="text-(--color-red)">*</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input v-model.number="form.last_paid_amount" type="number" step="0.01" min="0" placeholder="0.00"
                        required
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
                <p class="text-xs text-(--color-text-secondary) mt-1">Total amount paid so far for this store</p>
            </div>

            <!-- Last Paid Through -->
            <div class="mb-4">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Last Paid Through
                </label>
                <input v-model="form.last_paid_through" type="date"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                <p class="text-xs text-(--color-text-secondary) mt-1">Date up to which payment has been made</p>
            </div>

            <!-- Notes -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Notes
                </label>
                <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this payment"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
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
                    Saving...
                </span>
                <span v-else>Update Payment</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { GodownStoreBillStore } from '@/types/godownStoreBill'
import { useStoresStore } from '@/stores/stores'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'

const props = defineProps<{
    store: GodownStoreBillStore | null
}>()

const emit = defineEmits<{
    'updated': []
    'cancel': []
}>()

const storesStore = useStoresStore()
const submitting = ref(false)

const form = ref({
    last_paid_amount: 0,
    last_paid_through: '',
    notes: '',
})

// Initialize form with store data
const initializeForm = () => {
    if (props.store) {
        form.value = {
            last_paid_amount: props.store.last_paid_amount || 0,
            last_paid_through: props.store.last_paid_through || '',
            notes: props.store.notes || '',
        }
    }
}

watch(() => props.store, initializeForm, { immediate: true })

const submit = async () => {
    if (!props.store) return

    if (form.value.last_paid_amount < 0) {
        push.error('Last paid amount cannot be negative')
        return
    }

    submitting.value = true

    try {
        // Update the store with new payment information
        const success = await storesStore.updateStore(props.store.id, {
            last_paid_amount: form.value.last_paid_amount,
            last_paid_through: form.value.last_paid_through || null,
            notes: form.value.notes || null,
        })

        if (success) {
            push.success('Payment information updated successfully!')
            emit('updated')
        }
    } catch (error) {
        console.error('Error updating payment:', error)
        push.error('Failed to update payment information')
    } finally {
        submitting.value = false
    }
}
</script>