<!-- src/components/features/customerAdditionalBills/AdditionalChargeForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Charge Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Charge Information
            </h3>
            <div class="space-y-4">
                <!-- Customer -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Customer <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.customer_id" required :disabled="isEditMode || lockedCustomerId !== null"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option :value="null">Select a customer</option>
                        <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                            {{ getCustomerDisplayName(customer) }}
                        </option>
                    </select>
                    <p v-if="isEditMode || lockedCustomerId !== null"
                        class="text-xs text-(--color-text-secondary) mt-1">
                        Customer cannot be changed
                    </p>
                </div>

                <!-- Entity Type -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Entity Type <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.entity_type" required :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="lot">Lot</option>
                        <option value="store">Store</option>
                        <option value="delivery">Delivery</option>
                        <option value="transport">Transport</option>
                        <option value="damage">Damage</option>
                        <option value="godown">Godown</option>
                    </select>
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">
                        Entity cannot be changed
                    </p>
                </div>

                <!-- Entity ID -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Entity <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.entity_id" required :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option :value="null">Select an entity</option>
                        <option v-for="option in entityOptions" :key="option.id" :value="option.id">
                            {{ option.label }}
                        </option>
                    </select>
                    <p v-if="entityOptions.length === 0" class="text-xs text-(--color-yellow) mt-1">
                        No {{ form.entity_type }} records found
                    </p>
                </div>

                <!-- Amount -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Amount <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.amount" type="number" step="0.01" min="0" placeholder="0.00"
                            required
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <!-- Description -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Description <span class="text-(--color-red)">*</span>
                    </label>
                    <textarea v-model="form.description" rows="3" placeholder="Enter a description for this charge"
                        required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Charge' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Customer } from '@/types/customer'
import type { AdditionalChargeEntityType } from '@/types/customerAdditionalBill'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { useCustomersStore } from '@/stores/customers'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useTransportsStore } from '@/stores/transports'
import { useDamagesStore } from '@/stores/damages'
import { useGodownsStore } from '@/stores/godowns'
import { push } from 'notivue'

interface ChargeBillShape {
    id: number
    customer_id: number
    entity_type: AdditionalChargeEntityType
    entity_id: number
    amount: number
    description: string
}

const props = defineProps<{
    mode?: 'create' | 'edit'
    bill?: ChargeBillShape | null
    lockedCustomerId?: number | null
}>()

const emit = defineEmits<{
    'charge-created': []
    'charge-updated': []
    'cancel': []
}>()

const store = useCustomerAdditionalBillsStore()
const customersStore = useCustomersStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const deliveriesStore = useDeliveriesStore()
const transportsStore = useTransportsStore()
const damagesStore = useDamagesStore()
const godownsStore = useGodownsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.bill)

const customerOptions = computed(() => customersStore.customers)

const entityOptions = computed<{ id: number; label: string }[]>(() => {
    const cid = form.value.customer_id
    switch (form.value.entity_type) {
        case 'lot':
            return lotsStore.lots
                .filter(l => !cid || l.customer_id === cid)
                .map(l => ({ id: l.id, label: `${lotsStore.getLotDisplayName(l)} — Lot #${l.lot_number}` }))
        case 'store':
            return storesStore.stores.map(s => ({
                id: s.id,
                label: storesStore.getStoreDisplayName(s),
            }))
        case 'delivery':
            return deliveriesStore.deliveries
                .filter(d => !cid || d.customer_id === cid)
                .map(d => ({
                    id: d.id,
                    label: `Delivery #${d.id}${d.from_location ? ` — ${d.from_location}` : ''}`,
                }))
        case 'transport':
            return transportsStore.transports
                .filter(t => !cid || t.customer_id === cid)
                .map(t => ({
                    id: t.id,
                    label: `Transport #${t.id}${t.from_location ? ` — ${t.from_location}` : ''}`,
                }))
        case 'damage':
            return damagesStore.damages.map(d => ({
                id: d.id,
                label: `Damage #${d.id} — ${d.reason}`,
            }))
        case 'godown':
            return godownsStore.godowns.map(g => ({ id: g.id, label: g.name }))
        default:
            return []
    }
})

const form = ref<{
    customer_id: number | null
    entity_type: AdditionalChargeEntityType
    entity_id: number | null
    amount: number
    description: string
}>({
    customer_id: null,
    entity_type: 'lot',
    entity_id: null,
    amount: 0,
    description: '',
})

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}

const initializeForm = () => {
    if (props.bill) {
        form.value = {
            customer_id: props.bill.customer_id,
            entity_type: props.bill.entity_type,
            entity_id: props.bill.entity_id,
            amount: props.bill.amount,
            description: props.bill.description,
        }
    } else {
        form.value = {
            customer_id: props.lockedCustomerId ?? null,
            entity_type: 'lot',
            entity_id: null,
            amount: 0,
            description: '',
        }
    }
}

watch(() => props.bill, initializeForm, { immediate: true })

watch(() => props.lockedCustomerId, (val) => {
    if (!props.bill && val != null) {
        form.value.customer_id = val
    }
})

// Reset entity selection if entity type changes
watch(() => form.value.entity_type, () => {
    form.value.entity_id = null
})

// Reset entity selection if customer changes (since entity options filter by customer)
watch(() => form.value.customer_id, () => {
    form.value.entity_id = null
})

const resetForm = () => {
    if (isEditMode.value && props.bill) {
        initializeForm()
    } else {
        form.value = {
            customer_id: props.lockedCustomerId ?? null,
            entity_type: 'lot',
            entity_id: null,
            amount: 0,
            description: '',
        }
    }
}

const submit = async () => {
    if (!form.value.customer_id) {
        push.error('Customer is required')
        return
    }
    if (!form.value.entity_id) {
        push.error('Entity is required')
        return
    }
    if (form.value.amount < 0) {
        push.error('Amount cannot be negative')
        return
    }
    if (!form.value.description.trim()) {
        push.error('Description is required')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.bill) {
            const result = await store.updateCharge(props.bill.id, {
                amount: form.value.amount,
                description: form.value.description.trim(),
            })
            if (result) {
                emit('charge-updated')
            }
        } else {
            const result = await store.createCharge({
                customer_id: form.value.customer_id,
                entity_type: form.value.entity_type,
                entity_id: form.value.entity_id,
                amount: form.value.amount,
                description: form.value.description.trim(),
            })
            if (result) {
                resetForm()
                emit('charge-created')
            }
        }
    } finally {
        submitting.value = false
    }
}
</script>