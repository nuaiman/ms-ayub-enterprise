<!-- src/components/features/lots/LotForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Lot Information
            </h3>

            <LotFields v-model:customer-id="form.customer_id" v-model:lot-number="form.lot_number"
                v-model:product-name="form.product_name" v-model:weight-unit="form.weight_unit"
                v-model:quantity-unit="form.quantity_unit" :customer-options="customerOptions" :disabled="submitting"
                :can-edit-customer="!isEditMode" />
        </div>

        <!-- Stores Section (create only) -->
        <div v-if="!isEditMode" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Stores
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ stores.length }} store(s)</span>
            </div>

            <div v-if="stores.length > 0" class="space-y-6">
                <div v-for="(store, index) in stores" :key="store.key"
                    class="relative p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10 space-y-6">
                    <button type="button" @click="removeStore(index)" :disabled="submitting"
                        class="absolute top-2 right-2 p-1 rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <div>
                        <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">Store #{{ index + 1 }}</h4>

                        <!-- Godown -->
                        <div class="mb-4">
                            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                                Godown <span class="text-(--color-red)">*</span>
                            </label>
                            <select v-model="store.godown_id" required :disabled="submitting"
                                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                                <option :value="null">Select a godown</option>
                                <option v-for="godown in godownOptions" :key="godown.id" :value="godown.id">
                                    {{ godown.name }}
                                </option>
                            </select>
                        </div>

                        <StoreFields v-model:weight="store.weight" v-model:quantity="store.quantity"
                            v-model:start-date="store.start_date" :disabled="submitting" :required="false" />
                    </div>

                    <!-- Customer Store Bill (billing month auto-derived from store.start_date) -->
                    <div class="border-t border-(--color-border)/60 pt-4">
                        <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">
                            Customer Store Bill
                        </h4>
                        <CustomerStoreBillFields v-model:bill-type="store.bill_type" v-model:rate="store.bill_rate"
                            :disabled="submitting" :derived="true" />
                    </div>

                    <!-- Majhi Bill -->
                    <div class="border-t border-(--color-border)/60 pt-4">
                        <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">
                            Majhi Bill
                        </h4>
                        <MajhiBillFields v-model:majhi-id="store.majhi_id" v-model:bill-type="store.majhi_bill_type"
                            v-model:rate="store.majhi_bill_rate" :majhi-options="majhiOptions" :disabled="submitting"
                            :derived="true" />
                    </div>

                    <!-- Godown Bill (billing month auto-derived from store.start_date) -->
                    <div class="border-t border-(--color-border)/60 pt-4">
                        <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">
                            Godown Bill
                        </h4>
                        <GodownBillFields v-model:bill-type="store.godown_bill_type"
                            v-model:rate="store.godown_bill_rate" :disabled="submitting" :derived="true" />
                    </div>
                </div>
            </div>

            <div v-else
                class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No stores added yet.
            </div>

            <div class="mt-4">
                <button type="button" @click="addStore" :disabled="submitting"
                    class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    Add Store
                </button>
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Lot' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Lot } from '@/types/lot'
import type { CustomerStoreBillType } from '@/types/customerStoreBill'
import type { MajhiBillType } from '@/types/majhiBill'
import type { GodownBillType } from '@/types/godownBill'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useGodownsStore } from '@/stores/godowns'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomerStoreBillsStore } from '@/stores/customerStoreBills'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import { useGodownBillsStore } from '@/stores/godownBills'
import { push } from 'notivue'
import LotFields from './LotFields.vue'
import StoreFields from '@/components/features/stores/StoreFields.vue'
import CustomerStoreBillFields from '@/components/features/customerStoreBills/CustomerStoreBillFields.vue'
import MajhiBillFields from '@/components/features/majhiBills/MajhiBillFields.vue'
import GodownBillFields from '@/components/features/godownBills/GodownBillFields.vue'

interface StoreForm {
    key: number
    storeId: number | null

    godown_id: number | null
    weight: number
    quantity: number
    start_date: string

    // Customer store bill (month_year auto-derived by backend)
    bill_type: CustomerStoreBillType
    bill_rate: number

    // Majhi bill
    majhi_id: number | null
    majhi_bill_type: MajhiBillType
    majhi_bill_rate: number

    // Godown bill (month_year auto-derived by backend)
    godown_bill_type: GodownBillType
    godown_bill_rate: number
}

const props = defineProps<{
    lot?: Lot | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'lot-created': []
    'lot-updated': []
    'cancel': []
}>()

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const godownsStore = useGodownsStore()
const majhisStore = useMajhisStore()
const customerStoreBillsStore = useCustomerStoreBillsStore()
const majhiBillsStore = useMajhiBillsStore()
const godownBillsStore = useGodownBillsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.lot)

const customerOptions = computed(() => customersStore.customers)
const godownOptions = computed(() => godownsStore.godowns)
const majhiOptions = computed(() => majhisStore.majhis)
const stores = ref<StoreForm[]>([])

let nextKey = 1

const today = new Date().toISOString().slice(0, 10)

const form = ref({
    customer_id: null as number | null,
    lot_number: '',
    product_name: '',
    weight_unit: 'kg',
    quantity_unit: 'units',
})

const createEmptyStore = (): StoreForm => ({
    key: nextKey++,
    storeId: null,
    godown_id: null,
    weight: 0,
    quantity: 0,
    start_date: today,
    bill_type: 'quantity',
    bill_rate: 0,
    majhi_id: null,
    majhi_bill_type: 'quantity',
    majhi_bill_rate: 0,
    godown_bill_type: 'fixed',
    godown_bill_rate: 0,
})

const addStore = () => { stores.value.push(createEmptyStore()) }
const removeStore = (index: number) => { stores.value.splice(index, 1) }

const initialize = () => {
    if (isEditMode.value && props.lot) {
        form.value = {
            customer_id: props.lot.customer_id,
            lot_number: props.lot.lot_number,
            product_name: props.lot.product_name,
            weight_unit: props.lot.weight_unit,
            quantity_unit: props.lot.quantity_unit,
        }
        stores.value = []
    } else {
        form.value = {
            customer_id: null,
            lot_number: '',
            product_name: '',
            weight_unit: 'kg',
            quantity_unit: 'units',
        }
        stores.value = []
    }
}

watch([() => props.lot, () => props.mode], () => {
    initialize()
}, { immediate: true })

const resetForm = () => { initialize() }

const validateStores = (): boolean => {
    for (let i = 0; i < stores.value.length; i++) {
        const s = stores.value[i]
        if (!s) continue

        if (!s.godown_id) {
            push.error(`Store #${i + 1}: godown is required`)
            return false
        }
        if (s.weight < 0 || s.quantity < 0) {
            push.error(`Store #${i + 1}: weight/quantity cannot be negative`)
            return false
        }
        if (s.weight === 0 && s.quantity === 0) {
            push.error(`Store #${i + 1}: weight or quantity must be greater than 0`)
            return false
        }

        // Customer store bill
        if (!s.bill_rate || s.bill_rate <= 0) {
            push.error(`Store #${i + 1}: customer store bill rate must be greater than 0`)
            return false
        }

        // Majhi bill
        if (!s.majhi_id) {
            push.error(`Store #${i + 1}: majhi is required`)
            return false
        }
        if (!s.majhi_bill_rate || s.majhi_bill_rate <= 0) {
            push.error(`Store #${i + 1}: majhi bill rate must be greater than 0`)
            return false
        }

        // Godown bill
        if (!s.godown_bill_rate || s.godown_bill_rate <= 0) {
            push.error(`Store #${i + 1}: godown bill rate must be greater than 0`)
            return false
        }
    }
    return true
}

const submit = async () => {
    if (!form.value.customer_id) {
        push.error('Please select a customer')
        return
    }
    if (!form.value.lot_number.trim()) {
        push.error('Lot number is required')
        return
    }
    if (!form.value.product_name.trim()) {
        push.error('Product name is required')
        return
    }

    if (!isEditMode.value && stores.value.length > 0) {
        if (!validateStores()) return
    }

    submitting.value = true

    try {
        // ============================================================
        // LOT
        // ============================================================
        let lotId: number

        if (isEditMode.value && props.lot) {
            const success = await lotsStore.updateLot(props.lot.id, {
                customer_id: form.value.customer_id,
                lot_number: form.value.lot_number.trim(),
                product_name: form.value.product_name.trim(),
                weight_unit: form.value.weight_unit,
                quantity_unit: form.value.quantity_unit,
            })
            if (!success) return
            lotId = props.lot.id
        } else {
            const newLot = await lotsStore.createLot({
                customer_id: form.value.customer_id,
                lot_number: form.value.lot_number.trim(),
                product_name: form.value.product_name.trim(),
                weight_unit: form.value.weight_unit,
                quantity_unit: form.value.quantity_unit,
            })
            if (!newLot) return
            lotId = newLot.id
        }

        // ============================================================
        // STORES + BILLS (create only)
        // Month_year is intentionally omitted — backend derives it
        // from each newly-created store's start_date.
        // ============================================================
        if (!isEditMode.value) {
            for (const store of stores.value) {
                const newStore = await storesStore.createStore({
                    lot_id: lotId,
                    godown_id: store.godown_id!,
                    weight: store.weight,
                    quantity: store.quantity,
                    start_date: store.start_date ? `${store.start_date}T00:00:00Z` : undefined,
                })
                if (!newStore) continue

                await customerStoreBillsStore.createCustomerStoreBill({
                    customer_id: form.value.customer_id,
                    store_id: newStore.id,
                    bill_type: store.bill_type,
                    rate: store.bill_rate,
                })

                await majhiBillsStore.createMajhiBill({
                    majhi_id: store.majhi_id!,
                    store_id: newStore.id,
                    bill_type: store.majhi_bill_type,
                    rate: store.majhi_bill_rate,
                })

                await godownBillsStore.createGodownBill({
                    godown_id: newStore.godown_id,
                    store_id: newStore.id,
                    bill_type: store.godown_bill_type,
                    rate: store.godown_bill_rate,
                })
            }
        }

        if (isEditMode.value) {
            push.success('Lot updated successfully!')
            emit('lot-updated')
        } else {
            push.success('Lot, stores and bills created successfully!')
            resetForm()
            emit('lot-created')
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update lot' : 'Failed to create lot')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
    if (godownsStore.godowns.length === 0) await godownsStore.fetchGodowns()
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
})
</script>