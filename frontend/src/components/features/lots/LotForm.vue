<!-- src/components/features/lots/LotForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Lot Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Lot Information
            </h3>

            <LotFields v-model:item-id="form.item_id" v-model:lot-number="form.lot_number"
                v-model:customer-charge-type="form.customer_charge_type" v-model:majhi-bill-type="form.majhi_bill_type"
                v-model:customer-storage-rate="form.customer_storage_rate" v-model:unload-rate="form.unload_rate"
                v-model:majhi-id="form.majhi_id" v-model:majhi-cut="form.majhi_cut" v-model:is-active="form.is_active"
                v-model:notes="form.notes" :item-options="itemOptions" :majhi-options="majhiOptions"
                :disabled="submitting" :required="true" :can-edit-item="!isEditMode" :can-edit-lot-number="false"
                :show-active="true" />
        </div>

        <!-- Stores Section -->
        <div class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Stores
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ stores.length }} store(s)</span>
            </div>

            <!-- Existing Stores -->
            <div v-if="stores.length > 0" class="space-y-4">
                <div v-for="(store, index) in stores" :key="index"
                    class="relative p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10">
                    <button type="button" @click="removeStore(index)"
                        class="absolute top-2 right-2 p-1 rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-colors">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">Store #{{ index + 1 }}</h4>

                    <StoreFields v-model:lot-id="store.lot_id" v-model:godown-id="store.godown_id"
                        v-model:store-bill-type="store.store_bill_type" v-model:godown-cut="store.godown_cut"
                        v-model:quantity="store.quantity" v-model:quantity-unit="store.quantity_unit"
                        v-model:weight="store.weight" v-model:weight-unit="store.weight_unit"
                        v-model:is-active="store.is_active" v-model:billing-start="store.billing_start"
                        v-model:billing-end="store.billing_end" v-model:notes="store.notes" :lot-options="[]"
                        :godown-options="godownOptions" :disabled="submitting" :required="false" :standalone="true"
                        :can-edit-lot="false" :can-edit-godown="true" :show-active="true" />
                </div>
            </div>

            <!-- Empty State -->
            <div v-else
                class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No stores added yet.
            </div>

            <!-- Add Store Button -->
            <div class="mt-4">
                <button type="button" @click="addStore"
                    class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    Add Store
                </button>
            </div>
        </div>

        <!-- Summary -->
        <div v-if="stores.length > 0" class="border-t border-(--color-border) pt-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-sm font-medium text-(--color-text-primary)">Summary</p>
                <div class="mt-2 space-y-1 text-sm text-(--color-text-secondary)">
                    <p>✓ {{ stores.length }} store(s) will be created</p>
                    <p class="text-xs text-(--color-text-secondary) mt-2">Lot will be created first, then stores</p>
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Lot' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Lot, CustomerChargeType, MajhiBillType } from '@/types/lot'
import type { Item } from '@/types/item'
import type { StoreBillType } from '@/types/store'
import { useLotsStore } from '@/stores/lots'
import { useItemsStore } from '@/stores/items'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useStoresStore } from '@/stores/stores'
import { push } from 'notivue'
import LotFields from './LotFields.vue'
import StoreFields from '@/components/features/stores/StoreFields.vue'

interface StoreForm {
    lot_id: number | null
    godown_id: number | null
    store_bill_type: StoreBillType
    godown_cut: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    is_active: boolean
    billing_start: string
    billing_end: string | null
    notes: string
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
const itemsStore = useItemsStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const storesStore = useStoresStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.lot)

const itemOptions = computed(() => {
    return itemsStore.items.filter(i => i.is_active)
})

const majhiOptions = computed(() => {
    return majhisStore.majhis
})

const godownOptions = computed(() => {
    return godownsStore.godowns.filter(g => g.is_active)
})

// Store list
const stores = ref<StoreForm[]>([])

// Calculate next lot number for a given item
const getNextLotNumber = (itemId: number): number => {
    const existingLots = lotsStore.lots.filter(l => l.item_id === itemId)
    if (existingLots.length === 0) return 1
    const maxLotNumber = Math.max(...existingLots.map(l => l.lot_number))
    return maxLotNumber + 1
}

// Form includes new fields with defaults
const form = ref({
    item_id: null as number | null,
    lot_number: 0,
    customer_charge_type: 'quantity' as CustomerChargeType,
    majhi_bill_type: 'quantity' as MajhiBillType,
    customer_storage_rate: 0,
    unload_rate: 0,
    majhi_id: null as number | null,
    majhi_cut: 0,
    is_active: true,
    notes: '',
    // New fields - not displayed in UI
    customer_last_paid_through: null as string | null,
    customer_last_paid_amount: 0,
    customer_paid_unload_amount: 0,
    majhi_total_paid: 0,
})

// When item changes, auto-generate lot number
const onItemChange = () => {
    if (form.value.item_id && !isEditMode.value) {
        form.value.lot_number = getNextLotNumber(form.value.item_id)
    }
}

const createEmptyStore = (): StoreForm => {
    const today = new Date().toISOString().slice(0, 10)
    return {
        lot_id: null,
        godown_id: null,
        store_bill_type: 'quantity',
        godown_cut: 0,
        quantity: 0,
        quantity_unit: 'units',
        weight: 0,
        weight_unit: 'kg',
        is_active: true,
        billing_start: today,
        billing_end: null,
        notes: '',
    }
}

const addStore = () => {
    stores.value.push(createEmptyStore())
}

const removeStore = (index: number) => {
    stores.value.splice(index, 1)
}

const initializeForm = () => {
    if (props.lot) {
        form.value = {
            item_id: props.lot.item_id,
            lot_number: props.lot.lot_number,
            customer_charge_type: props.lot.customer_charge_type,
            majhi_bill_type: props.lot.majhi_bill_type,
            customer_storage_rate: props.lot.customer_storage_rate,
            unload_rate: props.lot.unload_rate,
            majhi_id: props.lot.majhi_id || null,
            majhi_cut: props.lot.majhi_cut,
            is_active: props.lot.is_active !== undefined ? props.lot.is_active : true,
            notes: props.lot.notes || '',
            customer_last_paid_through: props.lot.customer_last_paid_through || null,
            customer_last_paid_amount: props.lot.customer_last_paid_amount || 0,
            customer_paid_unload_amount: props.lot.customer_paid_unload_amount || 0,
            majhi_total_paid: props.lot.majhi_total_paid || 0,
        }
    } else {
        form.value = {
            item_id: null,
            lot_number: 0,
            customer_charge_type: 'quantity',
            majhi_bill_type: 'quantity',
            customer_storage_rate: 0,
            unload_rate: 0,
            majhi_id: null,
            majhi_cut: 0,
            is_active: true,
            notes: '',
            customer_last_paid_through: null,
            customer_last_paid_amount: 0,
            customer_paid_unload_amount: 0,
            majhi_total_paid: 0,
        }
    }

    // Reset stores
    stores.value = []

    // If editing, load existing stores for this lot
    if (props.lot) {
        const existingStores = storesStore.getStoresByLotId(props.lot.id)
        if (existingStores.length > 0) {
            existingStores.forEach(store => {
                const billingStart = store.billing_start ? new Date(store.billing_start).toISOString().slice(0, 10) : new Date().toISOString().slice(0, 10)
                const billingEnd = store.billing_end ? new Date(store.billing_end).toISOString().slice(0, 10) : null

                stores.value.push({
                    lot_id: store.lot_id,
                    godown_id: store.godown_id,
                    store_bill_type: store.store_bill_type,
                    godown_cut: store.godown_cut,
                    quantity: store.quantity,
                    quantity_unit: store.quantity_unit || 'units',
                    weight: store.weight,
                    weight_unit: store.weight_unit || 'kg',
                    is_active: store.is_active,
                    billing_start: billingStart,
                    billing_end: billingEnd,
                    notes: store.notes || '',
                })
            })
        }
    }
}

watch(() => props.lot, initializeForm, { immediate: true })
watch(() => form.value.item_id, onItemChange)

const resetForm = () => {
    initializeForm()
}

const submit = async () => {
    if (!form.value.item_id) {
        push.error('Please select an item')
        return
    }

    if (!form.value.lot_number || form.value.lot_number <= 0) {
        push.error('Lot number is required')
        return
    }

    submitting.value = true

    try {
        let lotId: number

        if (isEditMode.value && props.lot) {
            // Update existing lot
            const success = await lotsStore.updateLot(props.lot.id, {
                customer_charge_type: form.value.customer_charge_type,
                majhi_bill_type: form.value.majhi_bill_type,
                customer_storage_rate: form.value.customer_storage_rate,
                unload_rate: form.value.unload_rate,
                majhi_id: form.value.majhi_id,
                majhi_cut: form.value.majhi_cut,
                notes: form.value.notes.trim() || null,
                // New fields - not displayed but can be updated
                customer_last_paid_through: form.value.customer_last_paid_through,
                customer_last_paid_amount: form.value.customer_last_paid_amount,
                customer_paid_unload_amount: form.value.customer_paid_unload_amount,
                majhi_total_paid: form.value.majhi_total_paid,
            })

            if (!success) {
                submitting.value = false
                return
            }

            lotId = props.lot.id
            push.success('Lot updated successfully!')
            emit('lot-updated')
        } else {
            // Create new lot
            const newLot = await lotsStore.createLot({
                item_id: form.value.item_id,
                lot_number: form.value.lot_number,
                customer_charge_type: form.value.customer_charge_type,
                majhi_bill_type: form.value.majhi_bill_type,
                customer_storage_rate: form.value.customer_storage_rate,
                unload_rate: form.value.unload_rate,
                majhi_id: form.value.majhi_id,
                majhi_cut: form.value.majhi_cut,
                notes: form.value.notes.trim() || null,
                // New fields - defaults
                customer_last_paid_through: null,
                customer_last_paid_amount: 0,
                customer_paid_unload_amount: 0,
                majhi_total_paid: 0,
            })

            if (!newLot) {
                submitting.value = false
                return
            }

            lotId = newLot.id
            push.success('Lot created successfully!')
            emit('lot-created')
        }

        // Create stores
        if (stores.value.length > 0) {
            let createdCount = 0
            let failedCount = 0

            for (const storeData of stores.value) {
                if (!storeData.godown_id) {
                    failedCount++
                    continue
                }

                const billingStart = storeData.billing_start ? `${storeData.billing_start}T00:00:00Z` : undefined
                const billingEnd = storeData.billing_end ? `${storeData.billing_end}T00:00:00Z` : null

                const newStore = await storesStore.createStore({
                    lot_id: lotId,
                    godown_id: storeData.godown_id,
                    store_bill_type: storeData.store_bill_type,
                    godown_cut: storeData.godown_cut,
                    quantity: storeData.quantity,
                    quantity_unit: storeData.quantity_unit,
                    weight: storeData.weight,
                    weight_unit: storeData.weight_unit,
                    billing_start: billingStart,
                    billing_end: billingEnd,
                    notes: storeData.notes.trim() || null,
                })

                if (newStore) {
                    createdCount++
                } else {
                    failedCount++
                }
            }

            if (createdCount > 0 && failedCount === 0) {
                push.success(`${createdCount} store(s) created successfully!`)
            } else if (createdCount > 0 && failedCount > 0) {
                push.warning(`${createdCount} store(s) created, ${failedCount} failed`)
            } else {
                push.warning('No stores were created')
            }
        }

        resetForm()

    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update lot' : 'Failed to create lot')
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    if (itemsStore.items.length === 0) {
        itemsStore.fetchItems()
    }
    if (majhisStore.majhis.length === 0) {
        majhisStore.fetchMajhis()
    }
    if (godownsStore.godowns.length === 0) {
        godownsStore.fetchGodowns()
    }
    if (lotsStore.lots.length === 0) {
        lotsStore.fetchLots()
    }
    if (storesStore.stores.length === 0) {
        storesStore.fetchStores()
    }
})
</script>