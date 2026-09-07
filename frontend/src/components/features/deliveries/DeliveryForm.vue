<!-- src/components/features/deliveries/DeliveryForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Customer Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Information
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Customer -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Customer <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.customer_id" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                        <option :value="null">Select a customer</option>
                        <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                            {{ getCustomerDisplayName(customer) }}
                        </option>
                    </select>
                    <p class="text-xs text-(--color-text-secondary) mt-1">Items will be filtered by this customer</p>
                </div>

                <!-- Delivery Date -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Delivery Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.delivery_date" type="date" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
            </div>
        </div>

        <!-- Receiver Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Receiver Information
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Receiver Name
                    </label>
                    <input v-model="form.receiver_name" type="text" placeholder="Enter receiver name"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Receiver Phone
                    </label>
                    <input v-model="form.receiver_phone" type="tel" placeholder="Enter receiver phone"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
            </div>
        </div>

        <!-- Location -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Location
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        From
                    </label>
                    <input v-model="form.from_location" type="text" placeholder="Enter origin location"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        To
                    </label>
                    <input v-model="form.to_location" type="text" placeholder="Enter destination location"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
            </div>
        </div>

        <!-- Delivery Items Section -->
        <div class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Delivery Items
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ deliveryItems.length }} item(s)</span>
            </div>

            <!-- Existing Delivery Items -->
            <div v-if="deliveryItems.length > 0" class="space-y-4">
                <div v-for="(item, index) in deliveryItems" :key="index"
                    class="relative p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10">
                    <button type="button" @click="removeDeliveryItem(index)"
                        class="absolute top-2 right-2 p-1 rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-colors">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">Item #{{ index + 1 }}</h4>

                    <DeliveryItemFields :customer-id="form.customer_id" v-model:item-id="item.item_id"
                        v-model:lot-id="item.lot_id" v-model:store-id="item.store_id" v-model:majhi-id="item.majhi_id"
                        v-model:quantity="item.quantity" v-model:quantity-unit="item.quantity_unit"
                        v-model:weight="item.weight" v-model:weight-unit="item.weight_unit"
                        v-model:loading-rate="item.loading_rate" v-model:majhi-cut="item.majhi_cut"
                        v-model:vehicle-number="item.vehicle_number" v-model:driver-number="item.driver_number"
                        v-model:notes="item.notes" :item-options="itemOptions" :lot-options="lotOptions"
                        :store-options="storeOptions" :majhi-options="majhiOptions" :disabled="submitting"
                        :required="false" :standalone="true" :can-edit-item="true" :can-edit-lot="true"
                        :can-edit-store="true" />
                </div>
            </div>

            <!-- Empty State -->
            <div v-else
                class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No delivery items added yet.
            </div>

            <!-- Add Delivery Item Button -->
            <div class="mt-4">
                <button type="button" @click="addDeliveryItem"
                    class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    Add Delivery Item
                </button>
            </div>
        </div>

        <!-- Delivery Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Delivery Notes
            </label>
            <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this delivery"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
        </div>

        <!-- Summary -->
        <div v-if="deliveryItems.length > 0" class="border-t border-(--color-border) pt-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-sm font-medium text-(--color-text-primary)">Summary</p>
                <div class="mt-2 space-y-1 text-sm text-(--color-text-secondary)">
                    <p>✓ {{ deliveryItems.length }} item(s) will be added</p>
                    <p class="text-xs text-(--color-text-secondary) mt-2">Delivery will be created first, then items</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || !canSubmit"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ isEditMode ? 'Saving...' : 'Creating...' }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Delivery' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Delivery } from '@/types/delivery'
import type { Customer } from '@/types/customer'
import type { Item } from '@/types/item'
import type { Lot } from '@/types/lot'
import type { Store } from '@/types/store'
import type { Majhi } from '@/types/majhi'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'
import DeliveryItemFields from '@/components/features/deliveryItems/DeliveryItemFields.vue'

interface DeliveryItemForm {
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
    delivery?: Delivery | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'delivery-created': []
    'delivery-updated': []
    'cancel': []
}>()

const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const deliveryItemsStore = useDeliveryItemsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.delivery)

const customerOptions = computed(() => {
    return customersStore.customers
})

const itemOptions = computed(() => {
    return itemsStore.items.filter(i => i.is_active)
})

const lotOptions = computed(() => {
    return lotsStore.lots.filter(l => l.is_active)
})

const storeOptions = computed(() => {
    return storesStore.stores.filter(s => s.is_active && (s.quantity > 0 || s.weight > 0))
})

const majhiOptions = computed(() => {
    return majhisStore.majhis
})

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}

const canSubmit = computed(() => {
    if (!form.value.customer_id) return false
    if (deliveryItems.value.length === 0) return false
    const hasValidItem = deliveryItems.value.some(item => item.store_id !== null)
    return hasValidItem
})

const today = new Date().toISOString().slice(0, 10)

const form = ref({
    customer_id: null as number | null,
    delivery_date: today,
    receiver_name: '',
    receiver_phone: '',
    from_location: '',
    to_location: '',
    notes: '',
})

const deliveryItems = ref<DeliveryItemForm[]>([])

const createEmptyDeliveryItem = (): DeliveryItemForm => ({
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

const addDeliveryItem = () => {
    deliveryItems.value.push(createEmptyDeliveryItem())
}

const removeDeliveryItem = (index: number) => {
    deliveryItems.value.splice(index, 1)
}

const initializeForm = () => {
    if (props.delivery) {
        const date = props.delivery.delivery_date ? new Date(props.delivery.delivery_date).toISOString().slice(0, 10) : today

        form.value = {
            customer_id: props.delivery.customer_id || null,
            delivery_date: date,
            receiver_name: props.delivery.receiver_name || '',
            receiver_phone: props.delivery.receiver_phone || '',
            from_location: props.delivery.from_location || '',
            to_location: props.delivery.to_location || '',
            notes: props.delivery.notes || '',
        }

        if (isEditMode.value) {
            const existingItems = deliveryItemsStore.getDeliveryItemsByDeliveryId(props.delivery.id)
            if (existingItems.length > 0) {
                deliveryItems.value = existingItems.map(item => ({
                    item_id: item.item_id,
                    lot_id: item.lot_id,
                    store_id: item.store_id,
                    majhi_id: item.majhi_id || null,
                    quantity: item.quantity || 0,
                    quantity_unit: item.quantity_unit || 'units',
                    weight: item.weight || 0,
                    weight_unit: item.weight_unit || 'kg',
                    loading_rate: item.loading_rate || 0,
                    majhi_cut: item.majhi_cut || 0,
                    vehicle_number: item.vehicle_number || '',
                    driver_number: item.driver_number || '',
                    notes: item.notes || '',
                }))
            } else {
                deliveryItems.value = [createEmptyDeliveryItem()]
            }
        }
    } else {
        form.value = {
            customer_id: null,
            delivery_date: today,
            receiver_name: '',
            receiver_phone: '',
            from_location: '',
            to_location: '',
            notes: '',
        }
        deliveryItems.value = [createEmptyDeliveryItem()]
    }
}

watch(() => props.delivery, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.delivery) {
        initializeForm()
    } else {
        form.value = {
            customer_id: null,
            delivery_date: today,
            receiver_name: '',
            receiver_phone: '',
            from_location: '',
            to_location: '',
            notes: '',
        }
        deliveryItems.value = [createEmptyDeliveryItem()]
    }
}

const submit = async () => {
    if (!form.value.customer_id) {
        push.error('Please select a customer')
        return
    }

    if (!form.value.delivery_date) {
        push.error('Delivery date is required')
        return
    }

    if (deliveryItems.value.length === 0) {
        push.error('At least one delivery item is required')
        return
    }

    const invalidItems = deliveryItems.value.filter(item => item.store_id === null)
    if (invalidItems.length > 0) {
        push.error('Please select a store for all delivery items')
        return
    }

    submitting.value = true

    const deliveryDate = formatDateForBackend(form.value.delivery_date)

    try {
        let deliveryId: number

        if (isEditMode.value && props.delivery) {
            const success = await deliveriesStore.updateDelivery(props.delivery.id, {
                customer_id: form.value.customer_id,
                delivery_date: deliveryDate,
                receiver_name: form.value.receiver_name.trim() || null,
                receiver_phone: form.value.receiver_phone.trim() || null,
                from_location: form.value.from_location.trim() || null,
                to_location: form.value.to_location.trim() || null,
                notes: form.value.notes.trim() || null,
            })

            if (!success) {
                submitting.value = false
                return
            }

            deliveryId = props.delivery.id
            push.success('Delivery updated successfully!')
            emit('delivery-updated')
        } else {
            const newDelivery = await deliveriesStore.createDelivery({
                customer_id: form.value.customer_id,
                delivery_date: deliveryDate,
                receiver_name: form.value.receiver_name.trim() || null,
                receiver_phone: form.value.receiver_phone.trim() || null,
                from_location: form.value.from_location.trim() || null,
                to_location: form.value.to_location.trim() || null,
                notes: form.value.notes.trim() || null,
            })

            if (!newDelivery) {
                submitting.value = false
                return
            }

            deliveryId = newDelivery.id
            push.success('Delivery created successfully!')
            emit('delivery-created')
        }

        // Delete existing delivery items if editing
        if (isEditMode.value) {
            const existingItems = deliveryItemsStore.getDeliveryItemsByDeliveryId(deliveryId)
            for (const item of existingItems) {
                await deliveryItemsStore.deleteDeliveryItem(item.id)
            }
        }

        // Create new delivery items
        let createdCount = 0
        let failedCount = 0

        for (const item of deliveryItems.value) {
            if (!item.store_id) {
                failedCount++
                continue
            }

            // Get lot to get customer_charge_type, customer_paid_unload_amount,
            // majhi_bill_type, and majhi_total_paid (these will be auto-populated by backend)
            const store = storesStore.getStoreById(item.store_id)
            if (!store) {
                failedCount++
                continue
            }

            const newItem = await deliveryItemsStore.createDeliveryItem({
                delivery_id: deliveryId,
                store_id: item.store_id,
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
                // The backend will populate these from the lot:
                // customer_charge_type, customer_paid_unload_amount,
                // majhi_bill_type, majhi_total_paid
            })

            if (newItem) {
                createdCount++
            } else {
                failedCount++
            }
        }

        if (createdCount > 0 && failedCount === 0) {
            push.success(`${createdCount} delivery item(s) added successfully!`)
        } else if (createdCount > 0 && failedCount > 0) {
            push.warning(`${createdCount} item(s) added, ${failedCount} failed`)
        } else {
            push.warning('No delivery items were added')
        }

        resetForm()

    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update delivery' : 'Failed to create delivery')
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    if (customersStore.customers.length === 0) {
        customersStore.fetchCustomers()
    }
    if (storesStore.stores.length === 0) {
        storesStore.fetchStores()
    }
    if (itemsStore.items.length === 0) {
        itemsStore.fetchItems()
    }
    if (lotsStore.lots.length === 0) {
        lotsStore.fetchLots()
    }
    if (majhisStore.majhis.length === 0) {
        majhisStore.fetchMajhis()
    }
    if (deliveryItemsStore.deliveryItems.length === 0) {
        deliveryItemsStore.fetchDeliveryItems()
    }
})
</script>