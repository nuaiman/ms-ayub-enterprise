<!-- src/components/features/stores/StoreForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Store Information
            </h3>

            <!-- Lot (create only) -->
            <div v-if="!isEditMode" class="mb-4">
                <!-- Locked display when lot is preset (e.g. Add Store from LotRow) -->
                <div v-if="presetLotId" class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
                    <p class="text-sm text-(--color-text-primary) mt-0.5">
                        {{ presetLotLabel }}
                    </p>
                </div>

                <!-- Editable dropdown otherwise -->
                <template v-else>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Lot <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.lot_id" required :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option :value="null">Select a lot</option>
                        <option v-for="lot in lotOptions" :key="lot.id" :value="lot.id">
                            {{ lot.product_name }} — Lot {{ lot.lot_number }}
                        </option>
                    </select>
                </template>
            </div>

            <!-- Edit mode: lot name -->
            <div v-else-if="props.store"
                class="mb-4 p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Lot</p>
                <p class="text-sm text-(--color-text-primary)">{{ currentLot?.product_name ?? '—' }} — Lot {{
                    currentLot?.lot_number ?? '—' }}</p>
            </div>

            <!-- Godown (create only) -->
            <div v-if="!isEditMode" class="mb-4">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Godown <span class="text-(--color-red)">*</span>
                </label>
                <select v-model="form.godown_id" required :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option :value="null">Select a godown</option>
                    <option v-for="godown in godownOptions" :key="godown.id" :value="godown.id">
                        {{ godown.name }}
                    </option>
                </select>
            </div>

            <!-- Edit mode: godown name -->
            <div v-else-if="props.store"
                class="mb-4 p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary)">Godown</p>
                <p class="text-sm text-(--color-text-primary)">{{ currentGodown?.name ?? '—' }}</p>
            </div>

            <StoreFields v-model:weight="form.weight" v-model:quantity="form.quantity"
                v-model:start-date="form.start_date" :disabled="submitting" :required="true" />

            <!-- Refresh current-month bills (edit mode only) -->
            <div v-if="isEditMode" class="mt-4 flex items-center gap-3">
                <button type="button" @click="handleRefreshBills" :disabled="submitting || refreshingBills"
                    class="px-3 py-1.5 text-xs font-medium rounded-lg border border-(--color-border) text-(--color-text-secondary) hover:bg-(--color-muted-bg) hover:text-(--color-text-primary) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                    {{ refreshingBills ? 'Refreshing…' : 'Refresh current-month bills' }}
                </button>
                <span class="text-xs text-(--color-text-secondary)">
                    Updates this month's bills with the current weight / quantity above.
                </span>
            </div>
        </div>

        <!-- Bills (create only) -->
        <template v-if="!isEditMode">
            <!-- Customer Store Bill -->
            <div class="border-t border-(--color-border) pt-6">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                    Customer Store Bill
                </h3>

                <CustomerStoreBillFields v-model:bill-type="billForm.bill_type" v-model:rate="billForm.rate"
                    :disabled="submitting" :derived="true" />
            </div>

            <!-- Majhi Bill -->
            <div class="border-t border-(--color-border) pt-6">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                    Majhi Bill
                </h3>

                <MajhiBillFields v-model:majhi-id="majhiBillForm.majhi_id" v-model:bill-type="majhiBillForm.bill_type"
                    v-model:rate="majhiBillForm.rate" :majhi-options="majhiOptions" :disabled="submitting"
                    :derived="true" />
            </div>

            <!-- Godown Bill -->
            <div class="border-t border-(--color-border) pt-6">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                    Godown Bill </h3>

                <GodownBillFields v-model:bill-type="godownBillForm.bill_type" v-model:rate="godownBillForm.rate"
                    :disabled="submitting" :derived="true" />
            </div>
        </template>

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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Store' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Store } from '@/types/store'
import type { CustomerStoreBillType } from '@/types/customerStoreBill'
import type { MajhiBillType } from '@/types/majhiBill'
import type { GodownBillType } from '@/types/godownBill'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useGodownsStore } from '@/stores/godowns'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomerStoreBillsStore } from '@/stores/customerStoreBills'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import { useGodownBillsStore } from '@/stores/godownBills'
import { push } from 'notivue'
import StoreFields from './StoreFields.vue'
import CustomerStoreBillFields from '@/components/features/customerStoreBills/CustomerStoreBillFields.vue'
import MajhiBillFields from '@/components/features/majhiBills/MajhiBillFields.vue'
import GodownBillFields from '@/components/features/godownBills/GodownBillFields.vue'

const props = withDefaults(defineProps<{
    store?: Store | null
    mode?: 'create' | 'edit'
    presetLotId?: number | null
}>(), {
    presetLotId: null,
})

const emit = defineEmits<{
    'store-created': []
    'store-updated': []
    'cancel': []
}>()

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()
const majhisStore = useMajhisStore()
const customerStoreBillsStore = useCustomerStoreBillsStore()
const majhiBillsStore = useMajhiBillsStore()
const godownBillsStore = useGodownBillsStore()

const submitting = ref(false)
const refreshingBills = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.store)

const lotOptions = computed(() => lotsStore.lots)
const godownOptions = computed(() => godownsStore.godowns)
const majhiOptions = computed(() => majhisStore.majhis)

const presetLotLabel = computed(() => {
    if (!props.presetLotId) return ''
    const lot = lotsStore.getLotById(props.presetLotId)
    if (!lot) return `Lot #${props.presetLotId}`
    return `${lot.product_name} — Lot ${lot.lot_number}`
})

const currentLot = computed(() => {
    if (!props.store) return null
    return lotsStore.getLotById(props.store.lot_id) || null
})

const currentGodown = computed(() => {
    if (!props.store) return null
    return godownsStore.getGodownById(props.store.godown_id) || null
})

const today = new Date().toISOString().slice(0, 10)

const form = ref({
    lot_id: null as number | null,
    godown_id: null as number | null,
    weight: 0,
    quantity: 0,
    start_date: today,
})

const billForm = ref({
    bill_type: 'quantity' as CustomerStoreBillType,
    rate: 0,
})

const majhiBillForm = ref({
    majhi_id: null as number | null,
    bill_type: 'quantity' as MajhiBillType,
    rate: 0,
})

const godownBillForm = ref({
    bill_type: 'fixed' as GodownBillType,
    rate: 0,
})

const initialize = () => {
    if (isEditMode.value && props.store) {
        const startDate = props.store.start_date
            ? new Date(props.store.start_date).toISOString().slice(0, 10)
            : today

        form.value = {
            lot_id: props.store.lot_id,
            godown_id: props.store.godown_id,
            weight: props.store.weight,
            quantity: props.store.quantity,
            start_date: startDate,
        }
    } else {
        form.value = {
            lot_id: props.presetLotId ?? null,
            godown_id: null,
            weight: 0,
            quantity: 0,
            start_date: today,
        }
    }

    billForm.value = { bill_type: 'quantity', rate: 0 }
    majhiBillForm.value = { majhi_id: null, bill_type: 'quantity', rate: 0 }
    godownBillForm.value = { bill_type: 'fixed', rate: 0 }
}

watch([() => props.store, () => props.mode, () => props.presetLotId], () => { initialize() }, { immediate: true })

const resetForm = () => { initialize() }

const submit = async () => {
    if (!isEditMode.value && !form.value.lot_id) {
        push.error('Please select a lot')
        return
    }
    if (!isEditMode.value && !form.value.godown_id) {
        push.error('Please select a godown')
        return
    }

    if (form.value.weight < 0) {
        push.error('Weight cannot be negative')
        return
    }
    if (form.value.quantity < 0) {
        push.error('Quantity cannot be negative')
        return
    }
    if (!form.value.start_date) {
        push.error('Start date is required')
        return
    }

    submitting.value = true

    const startDateISO = `${form.value.start_date}T00:00:00Z`

    try {
        // ============================================================
        // EDIT MODE
        // ============================================================
        if (isEditMode.value && props.store) {
            const storeSuccess = await storesStore.updateStore(props.store.id, {
                weight: form.value.weight,
                quantity: form.value.quantity,
                start_date: startDateISO,
            })
            if (!storeSuccess) return

            push.success('Store updated successfully!')
            emit('store-updated')
            return
        }

        // ============================================================
        // CREATE MODE
        // ============================================================
        if (!validateBills()) return

        const newStore = await storesStore.createStore({
            lot_id: form.value.lot_id!,
            godown_id: form.value.godown_id!,
            weight: form.value.weight,
            quantity: form.value.quantity,
            start_date: startDateISO,
        })
        if (!newStore) return

        const lot = lotsStore.getLotById(form.value.lot_id!)
        if (!lot) {
            push.error('Lot not found')
            return
        }

        await customerStoreBillsStore.createCustomerStoreBill({
            customer_id: lot.customer_id,
            store_id: newStore.id,
            bill_type: billForm.value.bill_type,
            rate: billForm.value.rate,
        })

        await majhiBillsStore.createMajhiBill({
            majhi_id: majhiBillForm.value.majhi_id!,
            store_id: newStore.id,
            bill_type: majhiBillForm.value.bill_type,
            rate: majhiBillForm.value.rate,
        })

        await godownBillsStore.createGodownBill({
            godown_id: newStore.godown_id,
            store_id: newStore.id,
            bill_type: godownBillForm.value.bill_type,
            rate: godownBillForm.value.rate,
        })

        push.success('Store and bills created successfully!')
        resetForm()
        emit('store-created')
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update store' : 'Failed to create store')
    } finally {
        submitting.value = false
    }
}

const validateBills = (): boolean => {
    if (!billForm.value.rate || billForm.value.rate <= 0) {
        push.error('Customer store bill: rate must be greater than 0')
        return false
    }
    if (!majhiBillForm.value.majhi_id) {
        push.error('Majhi bill: majhi is required')
        return false
    }
    if (!majhiBillForm.value.rate || majhiBillForm.value.rate <= 0) {
        push.error('Majhi bill: rate must be greater than 0')
        return false
    }
    if (!godownBillForm.value.rate || godownBillForm.value.rate <= 0) {
        push.error('Godown bill: rate must be greater than 0')
        return false
    }
    return true
}

const handleRefreshBills = async () => {
    if (!props.store) return
    refreshingBills.value = true
    try {
        const ok = await storesStore.refreshStoreBills(props.store.id)
        if (ok) {
            await Promise.all([
                godownBillsStore.fetchGodownBills(),
                majhiBillsStore.fetchMajhiBills(),
                customerStoreBillsStore.fetchCustomerStoreBills(),
            ])
        }
    } finally {
        refreshingBills.value = false
    }
}

onMounted(async () => {
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
    if (godownsStore.godowns.length === 0) await godownsStore.fetchGodowns()
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
})
</script>