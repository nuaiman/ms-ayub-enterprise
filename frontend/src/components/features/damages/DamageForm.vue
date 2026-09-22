<!-- src/components/features/damages/DamageForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Damage Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Damage Information
            </h3>
            <div class="space-y-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Store <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.store_id" required :disabled="isEditMode || !!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a store</option>
                        <option v-for="store in storeOptions" :key="store.id" :value="store.id">
                            {{ getStoreDisplayName(store) }}
                        </option>
                    </select>
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Store cannot be changed</p>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Reason <span class="text-(--color-red)">*</span>
                    </label>
                    <textarea v-model="form.reason" rows="2" placeholder="Enter the reason for damage" required
                        :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Damage Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.damage_date" type="date" required :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Inventory & Amount -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Inventory & Amount
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Quantity
                    </label>
                    <div class="flex gap-2">
                        <input v-model.number="form.quantity" type="number" step="0.01" min="0" placeholder="0"
                            :disabled="!!justCreatedId"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                        <input v-model="form.quantity_unit" type="text" placeholder="Unit" :disabled="!!justCreatedId"
                            class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Weight
                    </label>
                    <div class="flex gap-2">
                        <input v-model.number="form.weight" type="number" step="0.01" min="0" placeholder="0"
                            :disabled="!!justCreatedId"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                        <input v-model="form.weight_unit" type="text" placeholder="Unit" :disabled="!!justCreatedId"
                            class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Amount <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.amount" type="number" step="0.01" min="0" placeholder="0.00"
                            required :disabled="!!justCreatedId"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>
            <p class="text-xs text-(--color-text-secondary) mt-2">Either quantity or weight must be greater than 0</p>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea v-model="form.notes" rows="2" placeholder="Enter any additional notes" :disabled="!!justCreatedId"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
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
                        Damage record was created but the image upload failed. Retry below or remove the image.
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
                <span
                    class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">&#2547;</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Damage } from '@/types/damage'
import type { Store } from '@/types/store'
import { useDamagesStore } from '@/stores/damages'
import { useStoresStore } from '@/stores/stores'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'

const props = defineProps<{
    damage?: Damage | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'damage-created': []
    'damage-updated': []
    'cancel': []
}>()

const damagesStore = useDamagesStore()
const storesStore = useStoresStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.damage)

const storeOptions = computed(() => storesStore.stores)

const getStoreDisplayName = (store: Store): string => {
    return storesStore.getStoreDisplayName(store)
}

const today = new Date().toISOString().slice(0, 10)

const form = ref({
    store_id: null as number | null,
    quantity: 0,
    quantity_unit: 'units',
    weight: 0,
    weight_unit: 'kg',
    damage_date: today,
    reason: '',
    amount: 0,
    notes: '',
})

// Tracks the ID of a freshly created damage when image upload failed.
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
// FORM INIT
// ============================================================
const initializeForm = () => {
    if (props.damage) {
        const damageDate = props.damage.damage_date ? new Date(props.damage.damage_date).toISOString().slice(0, 10) : today

        form.value = {
            store_id: props.damage.store_id,
            quantity: props.damage.quantity || 0,
            quantity_unit: props.damage.quantity_unit || 'units',
            weight: props.damage.weight || 0,
            weight_unit: props.damage.weight_unit || 'kg',
            damage_date: damageDate,
            reason: props.damage.reason || '',
            amount: props.damage.amount || 0,
            notes: props.damage.notes || '',
        }

        justCreatedId.value = null
        clearImageState()
        if (isEditMode.value && props.damage.image_url) {
            imagePreview.value = getImageUrl(props.damage.image_url) || null
        }
    } else {
        form.value = {
            store_id: null,
            quantity: 0,
            quantity_unit: 'units',
            weight: 0,
            weight_unit: 'kg',
            damage_date: today,
            reason: '',
            amount: 0,
            notes: '',
        }
        justCreatedId.value = null
        clearImageState()
    }
}

watch(() => props.damage, initializeForm, { immediate: true })

const resetForm = () => {
    initializeForm()
}

const submit = async () => {
    // ============================================================
    // RETRY IMAGE PATH
    // ============================================================
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
            push.success('Damage record and image created successfully!')
            resetForm()
            emit('damage-created')
        } finally {
            uploading.value = false
        }
        return
    }

    // ============================================================
    // NORMAL VALIDATION
    // ============================================================
    if (!form.value.store_id) {
        push.error('Please select a store')
        return
    }

    if (!form.value.reason.trim()) {
        push.error('Please enter a reason')
        return
    }

    if (form.value.quantity < 0 || form.value.weight < 0) {
        push.error('Quantity and weight cannot be negative')
        return
    }

    if (form.value.quantity === 0 && form.value.weight === 0) {
        push.error('Either quantity or weight must be greater than 0')
        return
    }

    if (form.value.amount < 0) {
        push.error('Amount cannot be negative')
        return
    }

    if (!form.value.damage_date) {
        push.error('Damage date is required')
        return
    }

    submitting.value = true

    const damageDate = formatDateForBackend(form.value.damage_date)

    try {
        let damageId: number

        if (isEditMode.value && props.damage) {
            const success = await damagesStore.updateDamage(props.damage.id, {
                quantity: form.value.quantity,
                quantity_unit: form.value.quantity_unit,
                weight: form.value.weight,
                weight_unit: form.value.weight_unit,
                damage_date: damageDate,
                reason: form.value.reason.trim(),
                amount: form.value.amount,
                notes: form.value.notes.trim() || null,
            })

            if (!success) {
                submitting.value = false
                return
            }

            damageId = props.damage.id
            push.success('Damage record updated successfully!')
        } else {
            const newDamage = await damagesStore.createDamage({
                store_id: form.value.store_id,
                quantity: form.value.quantity,
                quantity_unit: form.value.quantity_unit,
                weight: form.value.weight,
                weight_unit: form.value.weight_unit,
                damage_date: damageDate,
                reason: form.value.reason.trim(),
                amount: form.value.amount,
                notes: form.value.notes.trim() || null,
            })

            if (!newDamage) {
                submitting.value = false
                return
            }

            damageId = newDamage.id
            push.success('Damage record created successfully!')
        }

        // ============================================================
        // IMAGE HANDLING
        // ============================================================
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
            emit('damage-created')
        }

        resetForm()
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update damage record' : 'Failed to create damage record')
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    if (storesStore.stores.length === 0) {
        storesStore.fetchStores()
    }
})
</script>