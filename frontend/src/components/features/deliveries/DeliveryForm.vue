<!-- src/components/features/deliveries/DeliveryForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Delivery Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Delivery Information
            </h3>

            <div class="space-y-4">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Customer
                        </label>
                        <select v-model="form.customer_id" :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                            <option :value="null">— None —</option>
                            <option v-for="c in customerOptions" :key="c.id" :value="c.id">
                                {{ customerLabel(c) }}
                            </option>
                        </select>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Delivery Date <span class="text-(--color-red)">*</span>
                        </label>
                        <input v-model="form.delivery_date" type="date" required :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Receiver Name
                        </label>
                        <input v-model="form.receiver_name" type="text" placeholder="Enter receiver name"
                            :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Receiver Phone
                        </label>
                        <input v-model="form.receiver_phone" type="tel" placeholder="Enter receiver phone"
                            :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            From Location
                        </label>
                        <input v-model="form.from_location" type="text" placeholder="Enter from location"
                            :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            To Location
                        </label>
                        <input v-model="form.to_location" type="text" placeholder="Enter to location"
                            :disabled="submitting"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="form.notes" rows="2" placeholder="Additional details" :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
                </div>
            </div>
        </div>

        <!-- Items (create mode only) -->
        <div v-if="!isEditMode" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Items
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ items.length }} item(s)</span>
            </div>

            <div v-if="items.length > 0" class="space-y-4">
                <div v-for="(item, index) in items" :key="item.key"
                    class="relative p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10">
                    <button type="button" @click="removeItem(index)" :disabled="submitting"
                        class="absolute top-2 right-2 p-1 rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">
                        Item #{{ index + 1 }}
                    </h4>

                    <DeliveryItemFields v-model:store-id="item.store_id" v-model:majhi-id="item.majhi_id"
                        v-model:vehicle-number="item.vehicle_number" v-model:driver-number="item.driver_number"
                        v-model:quantity="item.quantity" v-model:weight="item.weight"
                        v-model:cd-bill-type="item.cd_bill_type" v-model:cd-bill-rate="item.cd_bill_rate"
                        v-model:majhi-bill-type="item.majhi_bill_type" v-model:majhi-bill-rate="item.majhi_bill_rate"
                        :store-options="storeOptions" :majhi-options="majhiOptions" :disabled="submitting" />
                </div>
            </div>

            <div v-else
                class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No items added yet.
            </div>

            <div class="mt-4">
                <button type="button" @click="addItem" :disabled="submitting"
                    class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    Add Item
                </button>
            </div>
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
                            class="absolute top-1 right-1 w-5 h-5 bg-(--color-red) text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">×</button>
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
                    <p v-else-if="isEditMode && props.delivery?.image_url"
                        class="text-xs text-(--color-text-secondary)">
                        Current image will be replaced
                    </p>
                </div>
            </div>
        </div>

        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting || uploading"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || uploading"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting || uploading" class="inline-flex items-center justify-center gap-2">
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
import type { CustomerDeliveryBillType } from '@/types/customerDeliveryBill'
import type { MajhiBillType } from '@/types/majhiBill'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { push } from 'notivue'
import DeliveryItemFields from './DeliveryItemFields.vue'

interface ItemForm {
    key: number
    store_id: number | null
    majhi_id: number | null
    vehicle_number: string
    driver_number: string
    quantity: number
    weight: number

    cd_bill_type: CustomerDeliveryBillType
    cd_bill_rate: number

    majhi_bill_type: MajhiBillType
    majhi_bill_rate: number
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
const majhisStore = useMajhisStore()
const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
const majhiBillsStore = useMajhiBillsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.delivery)

const customerOptions = computed(() => customersStore.customers)
const storeOptions = computed(() => storesStore.stores)
const majhiOptions = computed(() => majhisStore.majhis)

const items = ref<ItemForm[]>([])
let nextKey = 1

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

const createEmptyItem = (): ItemForm => ({
    key: nextKey++,
    store_id: null,
    majhi_id: null,
    vehicle_number: '',
    driver_number: '',
    quantity: 0,
    weight: 0,
    cd_bill_type: 'quantity',
    cd_bill_rate: 0,
    majhi_bill_type: 'quantity',
    majhi_bill_rate: 0,
})

const addItem = () => { items.value.push(createEmptyItem()) }
const removeItem = (index: number) => { items.value.splice(index, 1) }

// IMAGE STATE
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

const initialize = () => {
    if (props.delivery) {
        const date = new Date(props.delivery.delivery_date)
        form.value = {
            customer_id: props.delivery.customer_id,
            delivery_date: date.toISOString().slice(0, 10),
            receiver_name: props.delivery.receiver_name || '',
            receiver_phone: props.delivery.receiver_phone || '',
            from_location: props.delivery.from_location || '',
            to_location: props.delivery.to_location || '',
            notes: props.delivery.notes || '',
        }
        items.value = []
        clearImageState()
        if (props.delivery.image_url) {
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
        items.value = []
        clearImageState()
    }
}

watch(() => props.delivery, initialize, { immediate: true })

const resetForm = () => { initialize() }

const customerLabel = (c: Customer): string =>
    c.company_name || c.contact_person || `Customer #${c.id}`

/**
 * Validates every item row. Returns the index of the first invalid item
 * (via `problem.index`) and a human-readable reason (via `problem.reason`),
 * or `null` if everything checks out.
 */
const findFirstItemProblem = (): { index: number; reason: string } | null => {
    for (let i = 0; i < items.value.length; i++) {
        const it = items.value[i]
        if (!it) continue

        if (!it.store_id) return { index: i, reason: 'store is required' }
        if (!it.majhi_id) return { index: i, reason: 'majhi is required' }
        if (it.quantity < 0 || it.weight < 0) {
            return { index: i, reason: 'quantity/weight cannot be negative' }
        }
        if (it.quantity === 0 && it.weight === 0) {
            return { index: i, reason: 'quantity or weight must be greater than 0' }
        }
        if (!it.cd_bill_rate || it.cd_bill_rate <= 0) {
            return { index: i, reason: 'customer delivery bill rate must be greater than 0' }
        }
        if (!it.majhi_bill_rate || it.majhi_bill_rate <= 0) {
            return { index: i, reason: 'majhi bill rate must be greater than 0' }
        }
    }
    return null
}

const submit = async () => {
    if (!form.value.delivery_date) {
        push.error('Delivery date is required')
        return
    }

    // ---- Create mode: require at least 1 item, and every item valid ----
    if (!isEditMode.value) {
        if (items.value.length === 0) {
            push.error('Please add at least one item')
            return
        }
        const problem = findFirstItemProblem()
        if (problem) {
            push.error(`Item #${problem.index + 1}: ${problem.reason}`)
            return
        }
    }

    submitting.value = true

    try {
        const deliveryDateISO = `${form.value.delivery_date}T00:00:00Z`

        // ============================================================
        // EDIT MODE — master only.
        // ============================================================
        if (isEditMode.value && props.delivery) {
            const success = await deliveriesStore.updateDelivery(props.delivery.id, {
                customer_id: form.value.customer_id,
                delivery_date: deliveryDateISO,
                receiver_name: form.value.receiver_name.trim() || null,
                receiver_phone: form.value.receiver_phone.trim() || null,
                from_location: form.value.from_location.trim() || null,
                to_location: form.value.to_location.trim() || null,
                notes: form.value.notes.trim() || null,
            })
            if (!success) return

            if (imageFile.value) {
                uploading.value = true
                uploadProgress.value = 0
                uploadError.value = null
                const url = await uploadImage('deliveries', props.delivery.id, imageFile.value, (p) => {
                    uploadProgress.value = p
                })
                uploading.value = false
                if (!url) {
                    uploadError.value = 'Failed to upload image. Please try again.'
                    push.warning('Delivery saved but image upload failed.')
                    return
                }
            } else if (imageToDelete.value && props.delivery.image_url) {
                await deleteImage('deliveries', props.delivery.id, props.delivery.image_url)
            }

            await deliveriesStore.refreshDelivery(props.delivery.id)

            push.success('Delivery updated successfully!')
            emit('delivery-updated')
            return
        }

        // ============================================================
        // CREATE MODE.
        // ============================================================

        // 1. Create the master record.
        const newDelivery = await deliveriesStore.createDelivery({
            customer_id: form.value.customer_id,
            delivery_date: deliveryDateISO,
            receiver_name: form.value.receiver_name.trim() || null,
            receiver_phone: form.value.receiver_phone.trim() || null,
            from_location: form.value.from_location.trim() || null,
            to_location: form.value.to_location.trim() || null,
            notes: form.value.notes.trim() || null,
        })
        if (!newDelivery) return

        // 2. Upload image (if any).
        if (imageFile.value) {
            uploading.value = true
            uploadProgress.value = 0
            uploadError.value = null
            const url = await uploadImage('deliveries', newDelivery.id, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            uploading.value = false
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                push.warning('Delivery saved but image upload failed.')
            }
        }

        // 3. Create every item, then its bills.
        //    IMPORTANT: if any item fails, we surface the actual reason
        //    instead of silently skipping it.
        for (let i = 0; i < items.value.length; i++) {
            const it = items.value[i]
            if (!it) continue

            const createdItem = await deliveriesStore.createDeliveryItem(newDelivery.id, {
                store_id: it.store_id!,
                majhi_id: it.majhi_id!,
                vehicle_number: it.vehicle_number.trim() || null,
                driver_number: it.driver_number.trim() || null,
                quantity: it.quantity,
                weight: it.weight,
            })

            if (!createdItem) {
                // The item itself failed to save. The store already pushed
                // the exact backend reason to the toast stack. Add one
                // more contextual toast so the user knows which item and
                // what to do next.
                push.error(
                    `Item #${i + 1} failed to save. The delivery was created — ` +
                    `open it and add the missing item manually.`
                )
                await deliveriesStore.refreshDelivery(newDelivery.id)
                emit('delivery-created')
                return
            }

            // Customer delivery bill (only when the delivery has a customer).
            if (newDelivery.customer_id) {
                const cdBill = await customerDeliveryBillsStore.createCustomerDeliveryBill({
                    customer_id: newDelivery.customer_id,
                    delivery_item_id: createdItem.id,
                    bill_type: it.cd_bill_type,
                    rate: it.cd_bill_rate,
                })
                if (!cdBill) {
                    push.warning(`Item #${i + 1}: customer delivery bill could not be created.`)
                }
            }

            // Majhi bill.
            const mb = await majhiBillsStore.createMajhiBill({
                majhi_id: it.majhi_id!,
                delivery_item_id: createdItem.id,
                bill_type: it.majhi_bill_type,
                rate: it.majhi_bill_rate,
            })
            if (!mb) {
                push.warning(`Item #${i + 1}: majhi bill could not be created.`)
            }
        }

        // 4. Pull the fresh master (with its item count populated) into the store.
        await deliveriesStore.refreshDelivery(newDelivery.id)
        await deliveriesStore.loadItemsFor(newDelivery.id, true)

        push.success('Delivery created successfully!')
        resetForm()
        emit('delivery-created')
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update delivery' : 'Failed to create delivery')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
})
</script>