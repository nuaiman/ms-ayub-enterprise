<!-- src/components/features/customerDeliveryBills/CustomerDeliveryBillForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Delivery Bill Information
            </h3>

            <div v-if="props.bill" class="space-y-4 mb-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                        <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedCustomerLabel }}</p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Delivery Item</p>
                        <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedItemLabel }}</p>
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Weight at Billing
                            <span v-if="unitLabel('weight')" class="text-xs font-normal text-(--color-text-secondary)">
                                ({{ unitLabel('weight') }})
                            </span>
                        </label>
                        <input v-model.number="form.weight_at_billing" type="number" step="0.01" min="0"
                            placeholder="0.00" :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Quantity at Billing
                            <span v-if="unitLabel('quantity')"
                                class="text-xs font-normal text-(--color-text-secondary)">
                                ({{ unitLabel('quantity') }})
                            </span>
                        </label>
                        <input v-model.number="form.quantity_at_billing" type="number" step="0.01" min="0"
                            placeholder="0.00" :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>

            <CustomerDeliveryBillFields v-model:bill-type="form.bill_type" v-model:rate="form.rate"
                :disabled="submitting" />
        </div>

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
                    Saving...
                </span>
                <span v-else>Save Changes</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { CustomerDeliveryBill, CustomerDeliveryBillType } from '@/types/customerDeliveryBill'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useCustomersStore } from '@/stores/customers'
import { useDeliveriesStore } from '@/stores/deliveries'
import { push } from 'notivue'
import CustomerDeliveryBillFields from './CustomerDeliveryBillFields.vue'

const props = defineProps<{
    bill: CustomerDeliveryBill | null
}>()

const emit = defineEmits<{
    'bill-updated': []
    'cancel': []
}>()

const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
const customersStore = useCustomersStore()
const deliveriesStore = useDeliveriesStore()

const submitting = ref(false)

const form = ref({
    bill_type: 'quantity' as CustomerDeliveryBillType,
    rate: 0,
    weight_at_billing: 0,
    quantity_at_billing: 0,
})

const lockedCustomerLabel = computed(() => {
    if (!props.bill) return ''
    return customersStore.getCustomerName(props.bill.customer_id)
})

const lockedItemLabel = computed(() => {
    if (!props.bill) return ''
    const item = findItemById(props.bill.delivery_item_id)
    if (!item) return `Item #${props.bill.delivery_item_id}`
    return `Item #${item.id} · Delivery #${item.delivery_id}`
})

const findItemById = (itemId: number) => {
    for (const list of Object.values(deliveriesStore.itemsByDelivery)) {
        const found = list.find(i => i.id === itemId)
        if (found) return found
    }
    return null
}

const unitLabel = (kind: 'weight' | 'quantity'): string => {
    if (!props.bill) return ''
    if (kind === 'weight') return props.bill.weight_unit_at_billing || ''
    return props.bill.quantity_unit_at_billing || ''
}

const initialize = () => {
    if (props.bill) {
        form.value = {
            bill_type: props.bill.bill_type,
            rate: props.bill.rate,
            weight_at_billing: props.bill.weight_at_billing,
            quantity_at_billing: props.bill.quantity_at_billing,
        }
    }
}

watch(() => props.bill, initialize, { immediate: true })

const submit = async () => {
    if (!props.bill) return
    if (form.value.rate < 0) {
        push.error('Rate cannot be negative')
        return
    }

    submitting.value = true
    try {
        const success = await customerDeliveryBillsStore.updateCustomerDeliveryBill(props.bill.id, {
            bill_type: form.value.bill_type,
            rate: form.value.rate,
            weight_at_billing: form.value.weight_at_billing,
            quantity_at_billing: form.value.quantity_at_billing,
        })
        if (success) {
            emit('bill-updated')
        }
    } finally {
        submitting.value = false
    }
}
</script>