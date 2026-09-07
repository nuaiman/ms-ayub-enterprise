<!-- src/components/features/items/ItemForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Item Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Item Information
            </h3>

            <ItemFields v-model:product-name="form.product_name" v-model:category="form.category"
                v-model:customer-id="form.customer_id" v-model:is-active="form.is_active" v-model:notes="form.notes"
                :customer-options="customerOptions" :disabled="submitting" />
        </div>

        <!-- Lot Information (optional) -->
        <div v-if="showLotFields" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Lot Information
                </h3>
                <button type="button" @click="showLotFields = false"
                    class="text-xs text-(--color-red) hover:opacity-75 transition-colors">
                    Remove Lot
                </button>
            </div>

            <LotFields v-model:item-id="lotForm.item_id" v-model:lot-number="lotForm.lot_number"
                v-model:customer-charge-type="lotForm.customer_charge_type"
                v-model:majhi-bill-type="lotForm.majhi_bill_type"
                v-model:customer-storage-rate="lotForm.customer_storage_rate" v-model:unload-rate="lotForm.unload_rate"
                v-model:majhi-id="lotForm.majhi_id" v-model:majhi-cut="lotForm.majhi_cut"
                v-model:is-active="lotForm.is_active" v-model:notes="lotForm.notes" :item-options="[]"
                :majhi-options="majhiOptions" :disabled="submitting" :required="false" :standalone="true"
                :can-edit-item="false" :can-edit-lot-number="true" :show-active="true" />
        </div>

        <!-- Add Lot Button -->
        <div v-if="!showLotFields" class="border-t border-(--color-border) pt-4">
            <button type="button" @click="showLotFields = true"
                class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                </svg>
                Add Lot
            </button>
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

            <!-- Add Store Button - ALWAYS at the bottom -->
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
        <div v-if="showLotFields || stores.length > 0" class="border-t border-(--color-border) pt-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-sm font-medium text-(--color-text-primary)">Summary</p>
                <div class="mt-2 space-y-1 text-sm text-(--color-text-secondary)">
                    <p v-if="showLotFields">✓ Lot will be created</p>
                    <p v-if="stores.length > 0">✓ {{ stores.length }} store(s) will be created</p>
                    <p class="text-xs text-(--color-text-secondary) mt-2">Item will be created first, then lot, then
                        stores</p>
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Item' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Item } from '@/types/item'
import type { Customer } from '@/types/customer'
import type { StoreBillType } from '@/types/store'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { push } from 'notivue'
import ItemFields from './ItemFields.vue'
import LotFields from '@/components/features/lots/LotFields.vue'
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
    item?: Item | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'item-created': []
    'item-updated': []
    'cancel': []
}>()

const itemsStore = useItemsStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.item)

// Show/hide optional sections
const showLotFields = ref(false)
const stores = ref<StoreForm[]>([])

const customerOptions = computed(() => {
    return customersStore.customers
})

const majhiOptions = computed(() => {
    return majhisStore.majhis
})

const godownOptions = computed(() => {
    return godownsStore.godowns.filter(g => g.is_active)
})

// Item form
const form = ref({
    product_name: '',
    category: '',
    customer_id: null as number | null,
    is_active: true,
    notes: '',
})

// Lot form (when adding lot)
const lotForm = ref({
    item_id: null as number | null,
    lot_number: 0,
    customer_charge_type: 'quantity' as any,
    majhi_bill_type: 'quantity' as any,
    customer_storage_rate: 0,
    unload_rate: 0,
    majhi_id: null as number | null,
    majhi_cut: 0,
    is_active: true,
    notes: '',
})

// Calculate next lot number for a given item
const getNextLotNumber = (itemId: number): number => {
    const existingLots = lotsStore.lots.filter(l => l.item_id === itemId)
    if (existingLots.length === 0) return 1
    const maxLotNumber = Math.max(...existingLots.map(l => l.lot_number))
    return maxLotNumber + 1
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
    const today = new Date().toISOString().slice(0, 10)

    if (props.item) {
        form.value = {
            product_name: props.item.product_name || '',
            category: props.item.category || '',
            customer_id: props.item.customer_id || null,
            is_active: props.item.is_active !== undefined ? props.item.is_active : true,
            notes: props.item.notes || '',
        }
    } else {
        form.value = {
            product_name: '',
            category: '',
            customer_id: null,
            is_active: true,
            notes: '',
        }
    }

    // Reset optional sections
    showLotFields.value = false
    stores.value = []

    // Reset lot form
    lotForm.value = {
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
    }
}

watch(() => props.item, initializeForm, { immediate: true })

const resetForm = () => {
    initializeForm()
}

const submit = async () => {
    // Validate: at least product_name or category
    if (!form.value.product_name.trim() && !form.value.category.trim()) {
        push.error('Either product name or category is required')
        return
    }

    submitting.value = true

    try {
        // 1. Create the item first
        let itemId: number

        if (isEditMode.value && props.item) {
            // Update existing item
            const updated = await itemsStore.updateItem(props.item.id, {
                product_name: form.value.product_name.trim() || null,
                category: form.value.category.trim() || null,
                customer_id: form.value.customer_id,
                notes: form.value.notes.trim() || null,
            })

            if (!updated) {
                submitting.value = false
                return
            }

            itemId = props.item.id
            push.success('Item updated successfully!')
            emit('item-updated')
        } else {
            // Create new item
            const newItem = await itemsStore.createItem({
                product_name: form.value.product_name.trim() || null,
                category: form.value.category.trim() || null,
                customer_id: form.value.customer_id,
                notes: form.value.notes.trim() || null,
            })

            if (!newItem) {
                submitting.value = false
                return
            }

            itemId = newItem.id
            push.success('Item created successfully!')
        }

        // 2. If lot fields are shown, create the lot
        let lotId: number | null = null

        if (showLotFields.value) {
            // Set the item_id from the created/updated item
            lotForm.value.item_id = itemId

            // Auto-generate lot number if not set
            if (!lotForm.value.lot_number || lotForm.value.lot_number <= 0) {
                lotForm.value.lot_number = getNextLotNumber(itemId)
            }

            const newLot = await lotsStore.createLot({
                item_id: lotForm.value.item_id,
                lot_number: lotForm.value.lot_number,
                customer_charge_type: lotForm.value.customer_charge_type,
                majhi_bill_type: lotForm.value.majhi_bill_type,
                customer_storage_rate: lotForm.value.customer_storage_rate,
                unload_rate: lotForm.value.unload_rate,
                majhi_id: lotForm.value.majhi_id,
                majhi_cut: lotForm.value.majhi_cut,
                notes: lotForm.value.notes.trim() || null,
            })

            if (newLot) {
                lotId = newLot.id
                push.success('Lot created successfully!')
            } else {
                push.warning('Item created, but Lot creation failed')
                // Continue without lot
            }
        }

        // 3. Create all stores
        if (stores.value.length > 0 && lotId) {
            let createdCount = 0
            let failedCount = 0

            for (const storeData of stores.value) {
                // Skip stores without godown
                if (!storeData.godown_id) {
                    failedCount++
                    continue
                }

                // Format dates for backend
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
        } else if (stores.value.length > 0 && !lotId) {
            push.warning('Stores require a lot. Please add a lot first.')
        }

        resetForm()
        emit('item-created')

    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update item' : 'Failed to create item')
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    if (customersStore.customers.length === 0) {
        customersStore.fetchCustomers()
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