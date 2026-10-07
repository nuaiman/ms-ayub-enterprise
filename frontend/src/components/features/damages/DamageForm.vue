<!-- src/components/features/damages/DamageForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Damage Information
            </h3>

            <!-- Create mode: Store picker -->
            <template v-if="!isEditMode">
                <div class="mb-4">
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Store <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.store_id" required :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option :value="null">Select a store</option>
                        <option v-for="store in storeOptions" :key="store.id" :value="store.id">
                            {{ storeLabel(store) }}
                        </option>
                    </select>
                </div>
            </template>

            <!-- Edit mode: locked store -->
            <div v-else class="mb-4 p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedStoreLabel }}</p>
            </div>

            <!-- Quantity + Weight (at least one > 0) -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Quantity
                        <span v-if="unitLabel('quantity')" class="text-xs font-normal text-(--color-text-secondary)">
                            ({{ unitLabel('quantity') }})
                        </span>
                    </label>
                    <input v-model.number="form.quantity" type="number" step="0.01" min="0" placeholder="0"
                        :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Weight
                        <span v-if="unitLabel('weight')" class="text-xs font-normal text-(--color-text-secondary)">
                            ({{ unitLabel('weight') }})
                        </span>
                    </label>
                    <input v-model.number="form.weight" type="number" step="0.01" min="0" placeholder="0"
                        :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <!-- Damage date + Amount -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Damage Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.damage_date" type="date" required :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Amount
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.amount" type="number" step="0.01" min="0" placeholder="0.00"
                            :disabled="submitting"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>

            <!-- Reason -->
            <div class="mb-4">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Reason <span class="text-(--color-red)">*</span>
                </label>
                <input v-model="form.reason" type="text" placeholder="e.g. Broken bags, spillage" required
                    :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>

            <!-- Notes -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Notes
                </label>
                <textarea v-model="form.notes" rows="3" placeholder="Additional details" :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
            </div>
        </div>

        <!-- Image -->
        <div class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Damage Image <span class="normal-case text-xs font-normal">(Optional)</span>
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
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Upload Image</label>
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
                        Damage was created but the image upload failed. Retry below or remove the image.
                    </p>
                    <p v-else-if="isEditMode && props.damage?.image_url" class="text-xs text-(--color-text-secondary)">
                        Current image will be replaced
                    </p>
                </div>
            </div>
        </div>

        <!-- Actions -->
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
                    {{ justCreatedId ? 'Retrying image...' : (isEditMode ? 'Saving...' : 'Creating...') }}
                </span>
                <span v-else>{{ justCreatedId ? 'Retry Image Upload' : (isEditMode ? 'Save Changes' : 'Create Damage')
                }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Damage } from '@/types/damage'
import { useDamagesStore } from '@/stores/damages'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { push } from 'notivue'
import type { Store } from '@/types/store'

const props = defineProps<{
    damage?: Damage | null
    mode?: 'create' | 'edit'
    presetStoreId?: number | null
    lockedSource?: boolean
}>()

const emit = defineEmits<{
    'damage-created': []
    'damage-updated': []
    'cancel': []
}>()

const damagesStore = useDamagesStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.damage)
const lockedSource = computed(() => props.lockedSource === true || isEditMode.value)

const storeOptions = computed(() => storesStore.stores)

const today = new Date().toISOString().slice(0, 10)

const form = ref({
    store_id: (props.presetStoreId ?? null) as number | null,
    quantity: 0,
    weight: 0,
    damage_date: today,
    reason: '',
    amount: 0,
    notes: '',
})

// Units derived from the selected store's parent lot.
const selectedStore = computed<Store | null>(() => {
    if (!form.value.store_id) return null
    return storesStore.getStoreById(form.value.store_id) ?? null
})

const selectedLot = computed(() => {
    if (!selectedStore.value) return null
    return lotsStore.getLotById(selectedStore.value.lot_id) ?? null
})

const unitLabel = (kind: 'quantity' | 'weight'): string => {
    if (kind === 'quantity') return selectedLot.value?.quantity_unit ?? ''
    return selectedLot.value?.weight_unit ?? ''
}

const lockedStoreLabel = computed(() => {
    if (!selectedStore.value) return '—'
    const lot = selectedLot.value
    if (!lot) return `Store #${selectedStore.value.id}`
    return `Store #${selectedStore.value.id} — ${lot.product_name} (Lot ${lot.lot_number})`
})

const storeLabel = (store: Store): string => {
    const lot = lotsStore.getLotById(store.lot_id)
    if (!lot) return `Store #${store.id}`
    return `Store #${store.id} — ${lot.product_name} (Lot ${lot.lot_number})`
}

// Image state
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

const justCreatedId = ref<number | null>(null)

const initialize = () => {
    if (props.damage) {
        const date = new Date(props.damage.damage_date)
        form.value = {
            store_id: props.damage.store_id,
            quantity: props.damage.quantity,
            weight: props.damage.weight,
            damage_date: date.toISOString().slice(0, 10),
            reason: props.damage.reason,
            amount: props.damage.amount,
            notes: props.damage.notes || '',
        }
        justCreatedId.value = null
        clearImageState()
        if (props.damage.image_url) {
            imagePreview.value = getImageUrl(props.damage.image_url) || null
        }
    } else {
        form.value = {
            store_id: props.presetStoreId ?? null,
            quantity: 0,
            weight: 0,
            damage_date: today,
            reason: '',
            amount: 0,
            notes: '',
        }
        justCreatedId.value = null
        clearImageState()
    }
}

watch(() => props.damage, initialize, { immediate: true })

const resetForm = () => { initialize() }

const submit = async () => {
    // Retry image path
    if (justCreatedId.value !== null) {
        if (!imageFile.value) {
            resetForm()
            emit('damage-created')
            return
        }
        uploading.value = true
        uploadError.value = null
        uploadProgress.value = 0
        try {
            const url = await uploadImage('damages', justCreatedId.value, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                return
            }
            push.success('Damage and image created successfully!')
            resetForm()
            emit('damage-created')
        } finally {
            uploading.value = false
        }
        return
    }

    if (!form.value.store_id) {
        push.error('Please select a store')
        return
    }
    if (!form.value.reason.trim()) {
        push.error('Reason is required')
        return
    }
    if ((form.value.quantity || 0) <= 0 && (form.value.weight || 0) <= 0) {
        push.error('At least one of quantity or weight must be greater than 0')
        return
    }
    if ((form.value.amount || 0) < 0) {
        push.error('Amount cannot be negative')
        return
    }
    if (!form.value.damage_date) {
        push.error('Damage date is required')
        return
    }

    submitting.value = true

    try {
        let damageId: number

        if (isEditMode.value && props.damage) {
            const updated = await damagesStore.updateDamage(props.damage.id, {
                quantity: form.value.quantity,
                weight: form.value.weight,
                damage_date: `${form.value.damage_date}T00:00:00Z`,
                reason: form.value.reason.trim(),
                amount: form.value.amount,
                notes: form.value.notes.trim() || null,
            })
            if (!updated) { submitting.value = false; return }
            damageId = props.damage.id
        } else {
            const created = await damagesStore.createDamage({
                store_id: form.value.store_id,
                quantity: form.value.quantity,
                weight: form.value.weight,
                damage_date: `${form.value.damage_date}T00:00:00Z`,
                reason: form.value.reason.trim(),
                amount: form.value.amount,
                notes: form.value.notes.trim() || null,
            })
            if (!created) { submitting.value = false; return }
            damageId = created.id
        }

        // Image handling
        if (imageFile.value) {
            uploading.value = true
            uploadProgress.value = 0
            uploadError.value = null
            const url = await uploadImage('damages', damageId, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            uploading.value = false
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                if (!isEditMode.value) {
                    justCreatedId.value = damageId
                }
                push.warning('Damage saved but image upload failed. Retry below.')
                await damagesStore.fetchDamages()
                return
            }
        } else if (isEditMode.value && props.damage && imageToDelete.value && props.damage.image_url) {
            await deleteImage('damages', props.damage.id, props.damage.image_url)
        }

        if (isEditMode.value) {
            emit('damage-updated')
        } else {
            resetForm()
            emit('damage-created')
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update damage' : 'Failed to create damage')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
})
</script>