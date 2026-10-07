<!-- src/components/features/customerStoreBills/CustomerStoreBillForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Store Bill Information
            </h3>

            <!-- Create mode -->
            <div v-if="!isEditMode" class="space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
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

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Store <span class="text-(--color-red)">*</span>
                        </label>
                        <select v-model="form.store_id" required :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                            <option :value="null">Select a store</option>
                            <option v-for="s in storeOptions" :key="s.id" :value="s.id">
                                {{ storeLabel(s) }}
                            </option>
                        </select>
                    </div>
                </div>

                <div v-if="selectedStore" class="p-3 rounded-lg bg-(--color-blue)/5 border border-(--color-blue)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Billing Month</p>
                    <p class="text-sm font-semibold text-(--color-text-primary) mt-0.5">
                        {{ billingMonthLabel }}
                    </p>
                    <p class="text-xs text-(--color-text-secondary) mt-1">
                        Auto-derived from this store's start date ({{ storeStartLabel }}).
                    </p>
                </div>
            </div>

            <!-- Edit mode: locked + frozen snapshot overrides -->
            <div v-else class="space-y-4 mb-4">
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                        <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedCustomerLabel }}</p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                        <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedStoreLabel }}</p>
                    </div>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                        <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Month</p>
                        <p class="text-sm text-(--color-text-primary) mt-0.5">{{ props.bill?.month_year || '—' }}</p>
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

            <CustomerStoreBillFields v-model:bill-type="form.bill_type" v-model:rate="form.rate" :disabled="submitting"
                :derived="true" />
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
                    {{ isEditMode ? 'Saving...' : 'Creating...' }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Bill' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { CustomerStoreBill, CustomerStoreBillType } from '@/types/customerStoreBill'
import type { Customer } from '@/types/customer'
import type { Store } from '@/types/store'
import { useCustomerStoreBillsStore } from '@/stores/customerStoreBills'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { push } from 'notivue'
import CustomerStoreBillFields from './CustomerStoreBillFields.vue'

const props = defineProps<{
    bill?: CustomerStoreBill | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'bill-created': []
    'bill-updated': []
    'cancel': []
}>()

const customerStoreBillsStore = useCustomerStoreBillsStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.bill)

const customerOptions = computed(() => customersStore.customers)
const storeOptions = computed(() => storesStore.stores)

const form = ref({
    customer_id: null as number | null,
    store_id: null as number | null,
    bill_type: 'quantity' as CustomerStoreBillType,
    rate: 0,
    weight_at_billing: 0,
    quantity_at_billing: 0,
})

const selectedStore = computed(() => {
    if (!form.value.store_id) return null
    return storesStore.getStoreById(form.value.store_id) ?? null
})

const storeStartLabel = computed(() => {
    if (!selectedStore.value) return ''
    const d = new Date(selectedStore.value.start_date)
    if (isNaN(d.getTime())) return ''
    return d.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
})

const billingMonthLabel = computed(() => {
    if (!selectedStore.value) return ''
    const start = new Date(selectedStore.value.start_date)
    if (isNaN(start.getTime())) return ''

    const startYear = start.getFullYear()
    const startMon = start.getMonth() + 1

    const now = new Date()
    const nowYear = now.getFullYear()
    const nowMon = now.getMonth() + 1

    const startMonth = startYear * 12 + startMon
    const currentMonth = nowYear * 12 + nowMon

    const target = currentMonth >= startMonth
        ? { year: nowYear, month: nowMon }
        : { year: startYear, month: startMon }

    return new Date(target.year, target.month - 1).toLocaleDateString('en-US', {
        month: 'long',
        year: 'numeric',
    })
})

const lockedCustomerLabel = computed(() => {
    if (!props.bill) return ''
    return customersStore.getCustomerName(props.bill.customer_id)
})

const lockedStoreLabel = computed(() => {
    if (!props.bill) return ''
    const store = storesStore.getStoreById(props.bill.store_id)
    if (!store) return `Store #${props.bill.store_id}`
    const lot = lotsStore.getLotById(store.lot_id)
    if (lot) return `Store #${store.id} — ${lot.product_name} (Lot ${lot.lot_number})`
    return `Store #${store.id}`
})

const customerLabel = (c: Customer): string =>
    c.company_name || c.contact_person || `Customer #${c.id}`

const storeLabel = (s: Store): string => {
    const lot = lotsStore.getLotById(s.lot_id)
    const lotName = lot ? lot.product_name : '—'
    const lotNum = lot ? lot.lot_number : '—'
    return `Store #${s.id} — ${lotName} (Lot ${lotNum})`
}

const unitLabel = (kind: 'weight' | 'quantity'): string => {
    if (!props.bill) return ''
    if (kind === 'weight') return props.bill.weight_unit_at_billing || ''
    return props.bill.quantity_unit_at_billing || ''
}

const initialize = () => {
    if (props.bill) {
        form.value = {
            customer_id: props.bill.customer_id,
            store_id: props.bill.store_id,
            bill_type: props.bill.bill_type,
            rate: props.bill.rate,
            weight_at_billing: props.bill.weight_at_billing,
            quantity_at_billing: props.bill.quantity_at_billing,
        }
    } else {
        form.value = {
            customer_id: null,
            store_id: null,
            bill_type: 'quantity',
            rate: 0,
            weight_at_billing: 0,
            quantity_at_billing: 0,
        }
    }
}

watch(() => props.bill, initialize, { immediate: true })

const resetForm = () => { initialize() }

const submit = async () => {
    if (!isEditMode.value) {
        if (!form.value.customer_id) {
            push.error('Please select a customer')
            return
        }
        if (!form.value.store_id) {
            push.error('Please select a store')
            return
        }
    }
    if (form.value.rate < 0) {
        push.error('Rate cannot be negative')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.bill) {
            const success = await customerStoreBillsStore.updateCustomerStoreBill(props.bill.id, {
                bill_type: form.value.bill_type,
                rate: form.value.rate,
                weight_at_billing: form.value.weight_at_billing,
                quantity_at_billing: form.value.quantity_at_billing,
            })
            if (success) {
                push.success('Customer store bill updated successfully!')
                emit('bill-updated')
            }
        } else {
            const newBill = await customerStoreBillsStore.createCustomerStoreBill({
                customer_id: form.value.customer_id!,
                store_id: form.value.store_id!,
                bill_type: form.value.bill_type,
                rate: form.value.rate,
            })
            if (newBill) {
                push.success('Customer store bill created successfully!')
                resetForm()
                emit('bill-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update customer store bill' : 'Failed to create customer store bill')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
    initialize()
})
</script>