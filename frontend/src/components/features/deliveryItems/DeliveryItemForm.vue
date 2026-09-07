<!-- src/components/features/deliveryItems/DeliveryItemForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Delivery Selection -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Delivery Information
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Delivery -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Delivery <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.delivery_id" @change="onDeliveryChange" required
                        :disabled="isEditMode || isSubmitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a delivery</option>
                        <option v-for="delivery in deliveryOptions" :key="delivery.id" :value="delivery.id">
                            #{{ delivery.id }} - {{ getDeliveryDisplayName(delivery) }}
                        </option>
                    </select>
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Delivery cannot be changed
                    </p>
                </div>

                <!-- Delivery Info -->
                <div v-if="selectedDelivery">
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Delivery Info
                    </label>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border) text-sm">
                        <p class="text-(--color-text-secondary)">Customer: <span class="text-(--color-text-primary)">{{
                            getCustomerName(selectedDelivery.customer_id) }}</span></p>
                        <p class="text-(--color-text-secondary)">Date: <span class="text-(--color-text-primary)">{{
                            formatDeliveryDate(selectedDelivery.delivery_date) }}</span></p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Delivery Items Section -->
        <div v-if="!isEditMode && form.delivery_id" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Delivery Items
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ form.items.length }} item(s)</span>
            </div>

            <!-- Item Rows -->
            <div class="space-y-4 max-h-125 overflow-y-auto pr-1">
                <div v-for="(item, index) in form.items" :key="index"
                    class="p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10 hover:border-(--color-blue)/30 transition-all duration-200">

                    <div class="flex items-center justify-between mb-3">
                        <h4 class="text-sm font-medium text-(--color-text-primary) flex items-center gap-2">
                            <span
                                class="w-6 h-6 rounded-full bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center text-xs font-bold">
                                {{ index + 1 }}
                            </span>
                            Item #{{ index + 1 }}
                        </h4>
                        <span v-if="item.store_id"
                            class="text-xs text-(--color-green) bg-(--color-green)/10 px-2 py-0.5 rounded-full">
                            ✓ Store Selected
                        </span>
                    </div>

                    <DeliveryItemFields :customer-id="selectedDelivery?.customer_id || null"
                        v-model:item-id="item.item_id" v-model:lot-id="item.lot_id" v-model:store-id="item.store_id"
                        v-model:majhi-id="item.majhi_id" v-model:quantity="item.quantity"
                        v-model:quantity-unit="item.quantity_unit" v-model:weight="item.weight"
                        v-model:weight-unit="item.weight_unit" v-model:loading-rate="item.loading_rate"
                        v-model:majhi-cut="item.majhi_cut" v-model:vehicle-number="item.vehicle_number"
                        v-model:driver-number="item.driver_number" v-model:notes="item.notes"
                        :item-options="itemOptions" :lot-options="lotOptions" :store-options="storeOptions"
                        :majhi-options="majhiOptions" :disabled="isSubmitting" :required="false" :standalone="true"
                        :can-edit-item="true" :can-edit-lot="true" :can-edit-store="true" />

                    <div class="mt-3 flex items-center justify-end">
                        <button type="button" @click="removeItemRow(index)" :disabled="form.items.length <= 1"
                            class="text-xs text-(--color-red) hover:bg-(--color-red)/10 px-3 py-1 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed">
                            Remove Item
                        </button>
                    </div>
                </div>
            </div>

            <!-- Add Item Button -->
            <div class="mt-4">
                <button type="button" @click="addItemRow"
                    class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    Add Item
                </button>
                <p class="text-xs text-(--color-text-secondary) mt-2">
                    {{ form.items.length }} item(s) to add
                </p>
            </div>
        </div>

        <!-- Edit Mode - Single Item -->
        <div v-if="isEditMode && props.deliveryItem" class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Edit Delivery Item
            </h3>

            <DeliveryItemFields :customer-id="selectedDelivery?.customer_id || null" v-model:item-id="editForm.item_id"
                v-model:lot-id="editForm.lot_id" v-model:store-id="editForm.store_id"
                v-model:majhi-id="editForm.majhi_id" v-model:quantity="editForm.quantity"
                v-model:quantity-unit="editForm.quantity_unit" v-model:weight="editForm.weight"
                v-model:weight-unit="editForm.weight_unit" v-model:loading-rate="editForm.loading_rate"
                v-model:majhi-cut="editForm.majhi_cut" v-model:vehicle-number="editForm.vehicle_number"
                v-model:driver-number="editForm.driver_number" v-model:notes="editForm.notes"
                :item-options="itemOptions" :lot-options="lotOptions" :store-options="storeOptions"
                :majhi-options="majhiOptions" :disabled="isSubmitting" :required="true" :standalone="true"
                :can-edit-item="false" :can-edit-lot="false" :can-edit-store="false" />
        </div>

        <!-- Summary -->
        <div v-if="!isEditMode && form.items.length > 0 && form.delivery_id"
            class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
            <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-(--color-text-primary)">Total Items</span>
                <span class="text-xl font-bold text-(--color-blue)">{{ form.items.length }} item(s)</span>
            </div>
            <p class="text-xs text-(--color-text-secondary) mt-1">{{form.items.filter(i => i.store_id).length}}
                store(s) selected</p>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="isSubmitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="isSubmitting || !canSubmit"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="isSubmitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ isEditMode ? 'Saving...' : `Adding ${form.items.length} item(s)...` }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : `Add ${form.items.length} Item(s)` }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { push } from 'notivue'
import type { DeliveryItem } from '@/types/deliveryItem'
import type { Delivery } from '@/types/delivery'
import type { Item } from '@/types/item'
import type { Lot } from '@/types/lot'
import type { Store } from '@/types/store'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useStoresStore } from '@/stores/stores'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomersStore } from '@/stores/customers'
import DeliveryItemFields from './DeliveryItemFields.vue'

interface DeliveryItemRow {
    item_id: number | null
    lot_id: number | null
    store_id: number | null
    majhi_id: number | null
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    loading_rate: number
    majhi_cut: number
    vehicle_number: string
    driver_number: string
    notes: string
}

const props = defineProps<{
    mode?: 'create' | 'edit'
    deliveryItem?: DeliveryItem | null
}>()

const emit = defineEmits<{
    'delivery-item-created': []
    'delivery-item-updated': []
    'cancel': []
}>()

const deliveryItemsStore = useDeliveryItemsStore()
const deliveriesStore = useDeliveriesStore()
const storesStore = useStoresStore()
const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const customersStore = useCustomersStore()

const isSubmitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.deliveryItem)

// ============= OPTIONS =============
const deliveryOptions = computed(() => deliveriesStore.deliveries)

const itemOptions = computed(() => {
    return itemsStore.items.filter(i => i.is_active)
})

const lotOptions = computed(() => {
    return lotsStore.lots.filter(l => l.is_active)
})

const storeOptions = computed(() => {
    return storesStore.stores.filter(s => s.is_active && (s.quantity > 0 || s.weight > 0))
})

const majhiOptions = computed(() => majhisStore.majhis)

// ============= SELECTED ENTITIES =============
const selectedDelivery = computed(() => {
    if (!form.value.delivery_id) return null
    return deliveriesStore.getDeliveryById(form.value.delivery_id)
})

// ============= HELPERS =============
const getDeliveryDisplayName = (delivery: Delivery): string => {
    const customerName = delivery.customer_id ? customersStore.getCustomerName(delivery.customer_id) : 'No customer'
    const date = new Date(delivery.delivery_date).toLocaleDateString()
    return `${customerName} - ${date}`
}

const getCustomerName = (id: number | null): string => {
    if (!id) return 'No customer'
    return customersStore.getCustomerName(id)
}

const formatDeliveryDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    })
}

// ============= CHANGE HANDLERS =============
const onDeliveryChange = () => {
    form.value.items = [createEmptyItem()]
}

// ============= FORM STATE =============
const form = ref({
    delivery_id: null as number | null,
    items: [] as DeliveryItemRow[],
})

// Edit mode form
const editForm = ref({
    item_id: null as number | null,
    lot_id: null as number | null,
    store_id: null as number | null,
    majhi_id: null as number | null,
    quantity: 0,
    quantity_unit: 'units',
    weight: 0,
    weight_unit: 'kg',
    loading_rate: 0,
    majhi_cut: 0,
    vehicle_number: '',
    driver_number: '',
    notes: '',
})

const canSubmit = computed(() => {
    if (isEditMode.value) {
        return !!editForm.value.store_id
    }

    if (!form.value.delivery_id) return false

    const hasValidItem = form.value.items.some(item => item.store_id !== null)
    if (!hasValidItem) return false

    return true
})

const createEmptyItem = (): DeliveryItemRow => ({
    item_id: null,
    lot_id: null,
    store_id: null,
    majhi_id: null,
    quantity: 0,
    quantity_unit: 'units',
    weight: 0,
    weight_unit: 'kg',
    loading_rate: 0,
    majhi_cut: 0,
    vehicle_number: '',
    driver_number: '',
    notes: '',
})

const addItemRow = () => {
    form.value.items.push(createEmptyItem())
}

const removeItemRow = (index: number) => {
    if (form.value.items.length > 1) {
        form.value.items.splice(index, 1)
    }
}

const resetForm = () => {
    if (isEditMode.value && props.deliveryItem) {
        const store = storesStore.getStoreById(props.deliveryItem.store_id)
        const lot = store ? lotsStore.getLotById(store.lot_id) : null
        const item = lot ? itemsStore.getItemById(lot.item_id) : null

        editForm.value = {
            item_id: item?.id || null,
            lot_id: lot?.id || null,
            store_id: props.deliveryItem.store_id,
            majhi_id: props.deliveryItem.majhi_id || null,
            quantity: props.deliveryItem.quantity || 0,
            quantity_unit: props.deliveryItem.quantity_unit || 'units',
            weight: props.deliveryItem.weight || 0,
            weight_unit: props.deliveryItem.weight_unit || 'kg',
            loading_rate: props.deliveryItem.loading_rate || 0,
            majhi_cut: props.deliveryItem.majhi_cut || 0,
            vehicle_number: props.deliveryItem.vehicle_number || '',
            driver_number: props.deliveryItem.driver_number || '',
            notes: props.deliveryItem.notes || '',
        }
        form.value.delivery_id = props.deliveryItem.delivery_id
    } else {
        form.value = {
            delivery_id: null,
            items: [],
        }
        form.value.items.push(createEmptyItem())
    }
}

// Initialize on mount
watch(() => props.deliveryItem, resetForm, { immediate: true })

// ============= SUBMIT =============
const submit = async () => {
    if (isEditMode.value && props.deliveryItem) {
        if (!editForm.value.store_id) {
            push.error('Store is required')
            return
        }

        isSubmitting.value = true

        try {
            const success = await deliveryItemsStore.updateDeliveryItem(props.deliveryItem.id, {
                store_id: editForm.value.store_id,
                majhi_id: editForm.value.majhi_id,
                vehicle_number: editForm.value.vehicle_number.trim() || null,
                driver_number: editForm.value.driver_number.trim() || null,
                quantity: editForm.value.quantity,
                quantity_unit: editForm.value.quantity_unit,
                weight: editForm.value.weight,
                weight_unit: editForm.value.weight_unit,
                loading_rate: editForm.value.loading_rate,
                majhi_cut: editForm.value.majhi_cut,
                notes: editForm.value.notes.trim() || null,
            })

            if (success) {
                push.success('Delivery item updated successfully!')
                emit('delivery-item-updated')
                emit('cancel')
            }
        } catch (error) {
            console.error('Error:', error)
            push.error('Failed to update delivery item')
        } finally {
            isSubmitting.value = false
        }
        return
    }

    if (!form.value.delivery_id) {
        push.error('Delivery is required')
        return
    }

    const validItems = form.value.items.filter(item => item.store_id !== null)
    if (validItems.length === 0) {
        push.error('At least one item with a store is required')
        return
    }

    isSubmitting.value = true

    try {
        let successCount = 0
        let failCount = 0

        for (const item of validItems) {
            const result = await deliveryItemsStore.createDeliveryItem({
                delivery_id: form.value.delivery_id,
                store_id: item.store_id!,
                majhi_id: item.majhi_id,
                item_id: item.item_id!,
                lot_id: item.lot_id!,
                vehicle_number: item.vehicle_number.trim() || null,
                driver_number: item.driver_number.trim() || null,
                quantity: item.quantity || 0,
                quantity_unit: item.quantity_unit || 'units',
                weight: item.weight || 0,
                weight_unit: item.weight_unit || 'kg',
                loading_rate: item.loading_rate || 0,
                majhi_cut: item.majhi_cut || 0,
                notes: item.notes.trim() || null,
            })

            if (result) {
                successCount++
            } else {
                failCount++
            }
        }

        if (successCount > 0 && failCount === 0) {
            push.success(`${successCount} delivery item(s) added successfully!`)
            resetForm()
            emit('delivery-item-created')
        } else if (successCount > 0 && failCount > 0) {
            push.warning(`${successCount} item(s) added, ${failCount} failed`)
            emit('delivery-item-created')
        } else {
            push.error('Failed to add delivery items')
        }
    } catch (error) {
        console.error('Error:', error)
        push.error('Failed to add delivery items')
    } finally {
        isSubmitting.value = false
    }
}

// ============= MOUNT =============
onMounted(async () => {
    try {
        await Promise.all([
            deliveriesStore.fetchDeliveries(),
            storesStore.fetchStores(),
            itemsStore.fetchItems(),
            lotsStore.fetchLots(),
            majhisStore.fetchMajhis(),
            customersStore.fetchCustomers(),
        ])
    } catch (error) {
        console.error('Failed to load form data:', error)
        push.error('Failed to load form data')
    }
})
</script>