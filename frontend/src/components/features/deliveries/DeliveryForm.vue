<!-- src/components/features/deliveries/DeliveryForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Customer Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Information
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Customer <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.customer_id" required :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option :value="null">Select a customer</option>
                        <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                            {{ getCustomerDisplayName(customer) }}
                        </option>
                    </select>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Delivery Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.delivery_date" type="date" required :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
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
                        :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Receiver Phone
                    </label>
                    <input v-model="form.receiver_phone" type="tel" placeholder="Enter receiver phone"
                        :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
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
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">From</label>
                    <input v-model="form.from_location" type="text" placeholder="Enter origin location"
                        :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">To</label>
                    <input v-model="form.to_location" type="text" placeholder="Enter destination location"
                        :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Delivery Items Section (hidden when retrying image) -->
        <div v-if="!justCreatedId" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Delivery Items
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ deliveryItems.length }} item(s)</span>
            </div>

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

                    <DeliveryItemFields :customer-id="form.customer_id" v-model:lot-id="item.lot_id"
                        v-model:store-id="item.store_id" v-model:majhi-id="item.majhi_id"
                        v-model:quantity="item.quantity" v-model:quantity-unit="item.quantity_unit"
                        v-model:weight="item.weight" v-model:weight-unit="item.weight_unit"
                        v-model:loading-rate="item.loading_rate" v-model:majhi-cut="item.majhi_cut"
                        v-model:vehicle-number="item.vehicle_number" v-model:driver-number="item.driver_number"
                        v-model:notes="item.notes" :lot-options="lotOptions" :store-options="storeOptions"
                        :majhi-options="majhiOptions" :disabled="submitting" :required="false" :standalone="true"
                        :can-edit-lot="true" :can-edit-store="true" />
                </div>
            </div>

            <div v-else
                class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No delivery items added yet.
            </div>

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
                :disabled="!!justCreatedId"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>

        <!-- Image -->
        <div class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Delivery Image <span class="normal-case text-xs font-normal">(Optional)</span>
            </h3>
            <div class="flex flex-col sm:flex-row items-start gap-4">
                <div class="shrink-0">
                    <div v-if="imagePreview"
                        class="relative w-24 h-24 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                        <button type="button" @click="removeImage"
                            class="absolute top-1 right-1 w-5 h-5 bg-(--color-red) text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">Ã—</button>
                    </div>
                    <div v-else
                        class="w-24 h-24 rounded-lg border-2 border-dashed border-(--color-border) flex items-center justify-center bg-(--color-muted-bg)">
                        <svg class="w-10 h-10 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                    </div>
                </div>

                <div class="flex-1 space-y-3">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Upload
                            Image</label>
                        <div class="flex flex-wrap gap-2">
                            <label
                                class="px-4 py-2 text-sm font-medium rounded-lg cursor-pointer bg-(--color-muted-bg) border border-(--color-border) hover:bg-(--color-muted-bg)/70 transition-all duration-200 inline-flex items-center gap-2"
                                :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                                </svg>
                                {{ uploading ? 'Uploading...' : 'Choose File' }}
                                <input type="file" accept="image/*" class="hidden" @change="handleFileSelect"
                                    :disabled="uploading" />
                            </label>
                            <button v-if="imagePreview" type="button" @click="removeImage"
                                class="px-4 py-2 text-sm font-medium rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-all duration-200">
                                Remove
                            </button>
                        </div>
                    </div>

                    <div v-if="uploading" class="w-full bg-(--color-muted-bg) rounded-full h-1.5 overflow-hidden">
                        <div class="bg-(--color-blue) h-full rounded-full transition-all duration-300"
                            :style="{ width: uploadProgress + '%' }"></div>
                    </div>

                    <p v-if="uploadError" class="text-xs text-(--color-red)">{{ uploadError }}</p>
                    <p v-else-if="justCreatedId" class="text-xs text-(--color-yellow)">
                        Delivery was created but the image upload failed. Retry below or remove the image.
                    </p>
                    <p v-else-if="isEditMode && props.delivery?.image_url"
                        class="text-xs text-(--color-text-secondary)">
                        Current image will be replaced
                    </p>
                </div>
            </div>
        </div>

        <!-- Summary -->
        <div v-if="deliveryItems.length > 0 && !justCreatedId" class="border-t border-(--color-border) pt-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-sm font-medium text-(--color-text-primary)">Summary</p>
                <div class="mt-2 space-y-1 text-sm text-(--color-text-secondary)">
                    <p>âœ“ {{ deliveryItems.length }} item(s) will be added</p>
                    <p class="text-xs text-(--color-text-secondary) mt-2">Delivery will be created first, then items</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting || uploading"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || uploading || (!justCreatedId && !canSubmit)"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting || uploading" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ justCreatedId ? 'Retrying image...' : (isEditMode ? 'Saving...' : 'Creating...') }}
                </span>
                <span v-else>{{ justCreatedId ? 'Retry Image Upload' : (isEditMode ? 'Save Changes' : 'Create Delivery')
                    }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Delivery } from '@/types/delivery'
import type { Customer } from '@/types/customer'
import type { Lot } from '@/types/lot'
import type { Store } from '@/types/store'
import type { Majhi } from '@/types/majhi'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'
import DeliveryItemFields from '@/components/features/deliveryItems/DeliveryItemFields.vue'

interface DeliveryItemForm {
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
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const deliveryItemsStore = useDeliveryItemsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.delivery)

const customerOptions = computed(() => customersStore.customers)

const lotOptions = computed(() => lotsStore.lots.filter(l => l.is_active))
const storeOptions = computed(() => storesStore.stores.filter(s => s.is_active && (s.quantity > 0 || s.weight > 0)))
const majhiOptions = computed(() => majhisStore.majhis)

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

const justCreatedId = ref<number | null>(null)

// Image state + helpers
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)
const imageToDelete = ref(false)

const clearImageState = () => {
    imageFile.value = null
    imagePreview.value = null
    uploadError.value = null
    uploadProgress.value = 0
    imageToDelete.value = false
    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null
    if (fileInput) fileInput.value = ''
}

const handleFileSelect = (event: Event) => {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (!file) return

    if (!file.type.startsWith('image/')) {
        uploadError.value = 'Please select an image file'
        return
    }

    if (file.size > 5 * 1024 * 1024) {
        uploadError.value = 'Image size should be less than 5MB'
        return
    }

    uploadError.value = null
    imageFile.value = file
    imageToDelete.value = false

    const reader = new FileReader()
    reader.onload = (e) => {
        imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
}

const removeImage = () => {
    imageFile.value = null
    imagePreview.value = null
    uploadError.value = null
    uploadProgress.value = 0
    imageToDelete.value = true
    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null
    if (fileInput) fileInput.value = ''
}

const createEmptyDeliveryItem = (): DeliveryItemForm => ({
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

        justCreatedId.value = null
        clearImageState()
        if (isEditMode.value && props.delivery.image_url) {
            imagePreview.value = getImageUrl(props.delivery.image_url) || null
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
        justCreatedId.value = null
        clearImageState()
    }
}

watch(() => props.delivery, initializeForm, { immediate: true })

const resetForm = () => {
    initializeForm()
}

const submit = async () => {
    // RETRY IMAGE PATH
    if (justCreatedId.value !== null) {
        if (!imageFile.value) {
            resetForm()
            emit('delivery-created')
            return
        }
        uploading.value = true
        uploadError.value = null
        uploadProgress.value = 0
        try {
            const url = await uploadImage('deliveries', justCreatedId.value, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                return
            }
            push.success('Delivery and image created successfully!')
            resetForm()
            emit('delivery-created')
        } finally {
            uploading.value = false
        }
        return
    }

    // NORMAL VALIDATION
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
        }

        // Delete existing delivery items if editing
        if (isEditMode.value) {
            const existingItems = deliveryItemsStore.getDeliveryItemsByDeliveryId(deliveryId)
            for (const item of existingItems) {
                await deliveryItemsStore.deleteDeliveryItem(item.id)
            }
        }

        // Create delivery items
        let createdCount = 0
        let failedCount = 0

        for (const item of deliveryItems.value) {
            if (!item.store_id) {
                failedCount++
                continue
            }

            const store = storesStore.getStoreById(item.store_id)
            if (!store) {
                failedCount++
                continue
            }

            const newItem = await deliveryItemsStore.createDeliveryItem({
                delivery_id: deliveryId,
                store_id: item.store_id,
                majhi_id: item.majhi_id,
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
        }

        // IMAGE HANDLING
        if (imageFile.value) {
            uploading.value = true
            uploadProgress.value = 0
            uploadError.value = null
            const url = await uploadImage('deliveries', deliveryId, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            uploading.value = false
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                if (!isEditMode.value) {
                    justCreatedId.value = deliveryId
                }
                push.warning('Delivery saved but image upload failed. Retry below.')
                await deliveriesStore.fetchDeliveries()
                return
            }
        } else if (isEditMode.value && props.delivery && imageToDelete.value && props.delivery.image_url) {
            await deleteImage('deliveries', props.delivery.id, props.delivery.image_url)
        }

        if (isEditMode.value) {
            emit('delivery-updated')
        } else {
            emit('delivery-created')
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
    if (customersStore.customers.length === 0) customersStore.fetchCustomers()
    if (storesStore.stores.length === 0) storesStore.fetchStores()
    if (lotsStore.lots.length === 0) lotsStore.fetchLots()
    if (majhisStore.majhis.length === 0) majhisStore.fetchMajhis()
    if (deliveryItemsStore.deliveryItems.length === 0) deliveryItemsStore.fetchDeliveryItems()
})
</script>