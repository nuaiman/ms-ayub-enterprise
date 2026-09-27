<!-- src/components/features/lots/LotForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Lot Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Lot Information
            </h3>

            <LotFields v-model:customer-id="form.customer_id" v-model:product-name="form.product_name"
                v-model:category="form.category" v-model:lot-number="form.lot_number"
                v-model:customer-charge-type="form.customer_charge_type" v-model:majhi-bill-type="form.majhi_bill_type"
                v-model:customer-storage-rate="form.customer_storage_rate" v-model:unload-rate="form.unload_rate"
                v-model:majhi-id="form.majhi_id" v-model:majhi-cut="form.majhi_cut" v-model:is-active="form.is_active"
                v-model:notes="form.notes" :customer-options="customerOptions" :majhi-options="majhiOptions"
                :disabled="submitting" :can-edit-customer="!isEditMode && !isReaddMode"
                :can-edit-product="!isEditMode && !isReaddMode" :show-active="true" />
        </div>

        <!-- Stores Section (hidden when retrying image on created lot) -->
        <div v-if="!justCreatedId" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Stores
                </h3>
                <span class="text-xs text-(--color-text-secondary)">{{ stores.length }} store(s)</span>
            </div>

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

            <div v-else
                class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No stores added yet.
            </div>

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

        <!-- Image (all modes) -->
        <div class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Lot Image <span class="normal-case text-xs font-normal">(Optional)</span>
            </h3>
            <div class="flex flex-col sm:flex-row items-start gap-4">
                <div class="shrink-0">
                    <div v-if="imagePreview"
                        class="relative w-24 h-24 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                        <button type="button" @click="removeImage"
                            class="absolute top-1 right-1 w-5 h-5 bg-(--color-red) text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">৳—</button>
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
                        Lot was created but the image upload failed. Retry below or remove the image.
                    </p>
                    <p v-else-if="isEditMode && props.lot?.image_url" class="text-xs text-(--color-text-secondary)">
                        Current image will be replaced
                    </p>
                </div>
            </div>
        </div>

        <!-- Summary -->
        <div v-if="stores.length > 0 && !justCreatedId" class="border-t border-(--color-border) pt-4">
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-sm font-medium text-(--color-text-primary)">Summary</p>
                <div class="mt-2 space-y-1 text-sm text-(--color-text-secondary)">
                    <p>৳œ“ {{ stores.length }} store(s) will be created</p>
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
            <button type="submit" :disabled="submitting || uploading"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting || uploading" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ justCreatedId ? 'Retrying image...' : (isEditMode ? 'Saving...' : 'Creating...') }}
                </span>
                <span v-else>{{ justCreatedId ? 'Retry Image Upload' : (isEditMode ? 'Save Changes' : 'Create Lot')
                    }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Lot, CustomerChargeType, MajhiBillType } from '@/types/lot'
import type { StoreBillType } from '@/types/store'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useStoresStore } from '@/stores/stores'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
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
    prefill?: Lot | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'lot-created': []
    'lot-updated': []
    'cancel': []
}>()

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const storesStore = useStoresStore()

const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit')
const isReaddMode = computed(() => !isEditMode.value && !!props.prefill)

const customerOptions = computed(() => customersStore.customers)
const majhiOptions = computed(() => majhisStore.majhis)
const godownOptions = computed(() => godownsStore.godowns.filter(g => g.is_active))

const stores = ref<StoreForm[]>([])

// Tracks the ID of a freshly created lot when image upload failed.
// When set, the form skips lot/store creation and only retries the image.
const justCreatedId = ref<number | null>(null)

// ============================================================
// IMAGE STATE & HELPERS
// ============================================================
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

// ============================================================
// FORM STATE
// ============================================================
const getNextLotNumber = (customerId: number): number => {
    const existingLots = lotsStore.lots.filter(l => l.customer_id === customerId)
    if (existingLots.length === 0) return 1
    const maxLotNumber = Math.max(...existingLots.map(l => l.lot_number))
    return maxLotNumber + 1
}

const form = ref({
    customer_id: null as number | null,
    product_name: '',
    category: '',
    lot_number: 0,
    customer_charge_type: 'quantity' as CustomerChargeType,
    majhi_bill_type: 'quantity' as MajhiBillType,
    customer_storage_rate: 0,
    unload_rate: 0,
    majhi_id: null as number | null,
    majhi_cut: 0,
    is_active: true,
    notes: '',
    customer_last_paid_through: null as string | null,
    customer_last_paid_amount: 0,
    customer_paid_unload_amount: 0,
    majhi_total_paid: 0,
})

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

const resetToBlank = () => {
    form.value = {
        customer_id: null,
        product_name: '',
        category: '',
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
    stores.value = []
    justCreatedId.value = null
    clearImageState()
}

const seedFromLot = (source: Lot, zeroStock: boolean) => {
    form.value = {
        customer_id: source.customer_id,
        product_name: source.product_name || '',
        category: source.category || '',
        lot_number: zeroStock ? (source.customer_id ? getNextLotNumber(source.customer_id) : 0) : source.lot_number,
        customer_charge_type: source.customer_charge_type,
        majhi_bill_type: source.majhi_bill_type,
        customer_storage_rate: source.customer_storage_rate,
        unload_rate: source.unload_rate,
        majhi_id: source.majhi_id ?? null,
        majhi_cut: source.majhi_cut,
        is_active: source.is_active !== undefined ? source.is_active : true,
        notes: source.notes || '',
        customer_last_paid_through: source.customer_last_paid_through ?? null,
        customer_last_paid_amount: source.customer_last_paid_amount ?? 0,
        customer_paid_unload_amount: source.customer_paid_unload_amount ?? 0,
        majhi_total_paid: source.majhi_total_paid ?? 0,
    }

    const lotStores = storesStore.getStoresByLotId(source.id)
    stores.value = lotStores.map(s => {
        const billingStart = s.billing_start
            ? new Date(s.billing_start).toISOString().slice(0, 10)
            : new Date().toISOString().slice(0, 10)
        const billingEnd = s.billing_end
            ? new Date(s.billing_end).toISOString().slice(0, 10)
            : null

        return {
            lot_id: null,
            godown_id: s.godown_id,
            store_bill_type: s.store_bill_type,
            godown_cut: s.godown_cut,
            quantity: zeroStock ? 0 : s.quantity,
            quantity_unit: s.quantity_unit || 'units',
            weight: zeroStock ? 0 : s.weight,
            weight_unit: s.weight_unit || 'kg',
            is_active: s.is_active,
            billing_start: billingStart,
            billing_end: billingEnd,
            notes: s.notes || '',
        }
    })

    justCreatedId.value = null
    clearImageState()
    if (isEditMode.value && source.image_url) {
        imagePreview.value = getImageUrl(source.image_url) || null
    }
}

const initialize = () => {
    if (isEditMode.value && props.lot) {
        seedFromLot(props.lot, false)
    } else if (props.prefill) {
        seedFromLot(props.prefill, true)
    } else {
        resetToBlank()
    }
}

const ensureDependenciesLoaded = async () => {
    const tasks: Promise<unknown>[] = []

    if (customersStore.customers.length === 0) tasks.push(customersStore.fetchCustomers())
    if (majhisStore.majhis.length === 0) tasks.push(majhisStore.fetchMajhis())
    if (godownsStore.godowns.length === 0) tasks.push(godownsStore.fetchGodowns())
    if (lotsStore.lots.length === 0) tasks.push(lotsStore.fetchLots())
    if (storesStore.stores.length === 0) tasks.push(storesStore.fetchStores())

    if (tasks.length > 0) {
        await Promise.all(tasks)
    }
}

watch(() => form.value.customer_id, (newCustomerId, oldCustomerId) => {
    if (isEditMode.value) return
    if (justCreatedId.value) return
    if (newCustomerId === oldCustomerId) return
    if (!newCustomerId) {
        form.value.lot_number = 0
        return
    }
    form.value.lot_number = getNextLotNumber(newCustomerId)
})

watch([() => props.lot, () => props.prefill, () => props.mode], initialize, { immediate: true })

const resetForm = () => {
    initialize()
}

const submit = async () => {
    // ============================================================
    // RETRY IMAGE PATH — lot already created, only upload remains
    // ============================================================
    if (justCreatedId.value !== null) {
        if (!imageFile.value) {
            // Nothing to retry; just close via emit
            resetForm()
            emit('lot-created')
            return
        }
        uploading.value = true
        uploadError.value = null
        uploadProgress.value = 0
        try {
            const url = await uploadImage('lots', justCreatedId.value, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                return
            }
            push.success('Lot and image created successfully!')
            resetForm()
            emit('lot-created')
        } finally {
            uploading.value = false
        }
        return
    }

    // ============================================================
    // NORMAL CREATE / EDIT PATH
    // ============================================================
    if (!form.value.customer_id) {
        push.error('Please select a customer')
        return
    }

    if (!form.value.lot_number || form.value.lot_number <= 0) {
        push.error('Lot number is required')
        return
    }

    for (let i = 0; i < stores.value.length; i++) {
        const s = stores.value[i]
        if (!s) {
            push.error(`Store #${i + 1}: invalid entry`)
            return
        }
        if (!s.godown_id) {
            push.error(`Store #${i + 1}: godown is required`)
            return
        }
        if (!s.quantity || s.quantity <= 0) {
            push.error(`Store #${i + 1}: quantity must be greater than 0`)
            return
        }
        if (!s.weight || s.weight <= 0) {
            push.error(`Store #${i + 1}: weight must be greater than 0`)
            return
        }
    }

    submitting.value = true

    try {
        let lotId: number

        if (isEditMode.value && props.lot) {
            const success = await lotsStore.updateLot(props.lot.id, {
                customer_id: form.value.customer_id,
                product_name: form.value.product_name.trim() || null,
                category: form.value.category.trim() || null,
                customer_charge_type: form.value.customer_charge_type,
                majhi_bill_type: form.value.majhi_bill_type,
                customer_storage_rate: form.value.customer_storage_rate,
                unload_rate: form.value.unload_rate,
                majhi_id: form.value.majhi_id,
                majhi_cut: form.value.majhi_cut,
                notes: form.value.notes.trim() || null,
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
        } else {
            const newLot = await lotsStore.createLot({
                customer_id: form.value.customer_id,
                product_name: form.value.product_name.trim() || null,
                category: form.value.category.trim() || null,
                lot_number: form.value.lot_number,
                customer_charge_type: form.value.customer_charge_type,
                majhi_bill_type: form.value.majhi_bill_type,
                customer_storage_rate: form.value.customer_storage_rate,
                unload_rate: form.value.unload_rate,
                majhi_id: form.value.majhi_id,
                majhi_cut: form.value.majhi_cut,
                notes: form.value.notes.trim() || null,
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
        }

        // Create stores (only on fresh create, not on edit — edit keeps existing stores untouched here)
        if (!isEditMode.value && stores.value.length > 0) {
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

        // ============================================================
        // IMAGE HANDLING
        // ============================================================
        if (imageFile.value) {
            uploading.value = true
            uploadProgress.value = 0
            uploadError.value = null
            const url = await uploadImage('lots', lotId, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            uploading.value = false
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                // Keep dialog open, remember the created ID for retry
                if (!isEditMode.value) {
                    justCreatedId.value = lotId
                }
                push.warning('Lot saved but image upload failed. Retry below.')
                await lotsStore.fetchLots()
                return
            }
        } else if (isEditMode.value && props.lot && imageToDelete.value && props.lot.image_url) {
            await deleteImage('lots', props.lot.id, props.lot.image_url)
        }

        // Success — close dialog
        if (isEditMode.value) {
            emit('lot-updated')
        } else {
            emit('lot-created')
        }

        resetForm()
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update lot' : 'Failed to create lot')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await ensureDependenciesLoaded()
    initialize()
})
</script>