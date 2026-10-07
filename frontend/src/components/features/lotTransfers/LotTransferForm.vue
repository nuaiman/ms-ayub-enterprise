<!-- src/components/features/lotTransfers/LotTransferForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Lot picker (create only, not locked) -->
        <div v-if="!lockedLot">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Lot <span class="text-(--color-red)">*</span>
            </label>
            <select v-model="form.lot_id" required :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a lot</option>
                <option v-for="lot in lotOptions" :key="lot.id" :value="lot.id">
                    {{ lotLabel(lot) }}
                </option>
            </select>
        </div>

        <!-- Locked lot display -->
        <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
            <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedLotLabel }}</p>
        </div>

        <!-- Current customer snapshot -->
        <div v-if="currentLot" class="p-3 rounded-lg bg-(--color-surface) border border-(--color-border)">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Current Customer</p>
            <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                {{ currentCustomerName }}
            </p>
        </div>

        <!-- Target customer -->
        <div v-if="currentLot">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Transfer To <span class="text-(--color-red)">*</span>
            </label>
            <select v-model="form.to_customer_id" required :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select target customer</option>
                <option v-for="c in availableCustomers" :key="c.id" :value="c.id">
                    {{ customerLabel(c) }}
                </option>
            </select>
            <p v-if="availableCustomers.length === 0" class="text-xs text-(--color-yellow) mt-1">
                No other customers available.
            </p>
        </div>

        <!-- Billing block -->
        <div v-if="currentLot" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <div>
                    <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                        Customer Store Billing
                    </h3>
                    <p class="text-xs text-(--color-text-secondary) mt-1">
                        Optional. If a rate is provided, a customer store bill is created
                        for the current month for each active store under this lot, under
                        the new customer.
                    </p>
                </div>
            </div>

            <CustomerStoreBillFields v-model:bill-type="form.bill_type" v-model:rate="form.rate" :disabled="submitting"
                :derived="true" />
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea v-model="form.notes" rows="3" placeholder="Optional notes about this transfer"
                :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || !form.lot_id || !form.to_customer_id"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    Transferring...
                </span>
                <span v-else>Transfer Lot</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import type { Lot } from '@/types/lot'
import type { Customer } from '@/types/customer'
import type { LotBillType } from '@/types/lotTransfer'
import { useLotTransfersStore } from '@/stores/lotTransfers'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { push } from 'notivue'
import CustomerStoreBillFields from '@/components/features/customerStoreBills/CustomerStoreBillFields.vue'

const props = withDefaults(defineProps<{
    presetLotId?: number | null
    lockedLot?: boolean
}>(), {
    presetLotId: null,
    lockedLot: false,
})

const emit = defineEmits<{
    'transfer-created': []
    'cancel': []
}>()

const lotTransfersStore = useLotTransfersStore()
const lotsStore = useLotsStore()
const customersStore = useCustomersStore()

const submitting = ref(false)

const form = ref({
    lot_id: props.presetLotId as number | null,
    to_customer_id: null as number | null,
    bill_type: 'quantity' as LotBillType,
    rate: 0,
    notes: '',
})

const lotOptions = computed(() => lotsStore.lots)

const currentLot = computed<Lot | null>(() => {
    if (!form.value.lot_id) return null
    return lotsStore.getLotById(form.value.lot_id) ?? null
})

const currentCustomerName = computed(() => {
    if (!currentLot.value) return '—'
    return customersStore.getCustomerName(currentLot.value.customer_id)
})

const availableCustomers = computed(() => {
    if (!currentLot.value) return customersStore.customers
    return customersStore.customers.filter(c => c.id !== currentLot.value!.customer_id)
})

const lockedLotLabel = computed(() => {
    if (!currentLot.value) return '—'
    return `${currentLot.value.product_name} — Lot ${currentLot.value.lot_number}`
})

const customerLabel = (c: Customer): string =>
    c.company_name || c.contact_person || `Customer #${c.id}`

const lotLabel = (lot: Lot): string =>
    `${lot.product_name} — Lot ${lot.lot_number} (${customersStore.getCustomerName(lot.customer_id)})`

watch(() => props.presetLotId, (val) => {
    if (val) form.value.lot_id = val
}, { immediate: true })

// If the lot changes and the previously selected target customer is now the
// current customer, reset the target.
watch(currentLot, (lot) => {
    if (lot && form.value.to_customer_id === lot.customer_id) {
        form.value.to_customer_id = null
    }
})

const submit = async () => {
    if (!form.value.lot_id) {
        push.error('Please select a lot')
        return
    }
    if (!form.value.to_customer_id) {
        push.error('Please select a target customer')
        return
    }
    if (form.value.rate < 0) {
        push.error('Rate cannot be negative')
        return
    }

    submitting.value = true

    try {
        const payload = {
            to_customer_id: form.value.to_customer_id,
            notes: form.value.notes.trim() || null,
        } as {
            to_customer_id: number
            notes: string | null
            bill_type?: LotBillType
            rate?: number
        }

        // Only include billing fields if a positive rate is set.
        if (form.value.rate > 0) {
            payload.bill_type = form.value.bill_type
            payload.rate = form.value.rate
        }

        const result = await lotTransfersStore.createTransfer(form.value.lot_id, payload)

        if (result) {
            await Promise.all([
                lotsStore.fetchLots(),
                customersStore.fetchCustomers(),
            ])
            resetForm()
            emit('transfer-created')
        }
    } finally {
        submitting.value = false
    }
}

const resetForm = () => {
    form.value = {
        lot_id: props.presetLotId ?? null,
        to_customer_id: null,
        bill_type: 'quantity',
        rate: 0,
        notes: '',
    }
}

onMounted(async () => {
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
})
</script>