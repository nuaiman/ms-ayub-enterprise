<!-- src/components/features/incomes/IncomeForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Income Details -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Income Information
            </h3>
            <div class="space-y-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Title <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative" ref="dropdownRef">
                        <input ref="inputRef" v-model="searchTerm" type="text" placeholder="Enter income title" required
                            :disabled="!!justCreatedId" @focus="openDropdown" @input="onInput"
                            @keydown.down="selectNext" @keydown.up="selectPrevious" @keydown.enter="handleEnter"
                            @keydown.escape="closeDropdown"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />

                        <button type="button" @click="toggleDropdown" :disabled="!!justCreatedId"
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 9l-7 7-7-7" />
                            </svg>
                        </button>

                        <Transition enter-active-class="transition ease-out duration-200"
                            enter-from-class="opacity-0 -translate-y-1" enter-to-class="opacity-100 translate-y-0"
                            leave-active-class="transition ease-in duration-150"
                            leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 -translate-y-1">
                            <div v-if="isDropdownOpen && filteredSuggestions.length > 0 && !justCreatedId"
                                class="absolute left-0 right-0 top-full mt-1 max-h-48 overflow-y-auto bg-(--color-surface) border border-(--color-border) rounded-lg shadow-lg z-50 py-1">
                                <button v-for="(suggestion, index) in filteredSuggestions" :key="suggestion"
                                    type="button" @click="selectSuggestion(suggestion)"
                                    @mouseenter="highlightedIndex = index"
                                    class="w-full px-3 py-2 text-sm text-left transition-colors" :class="[
                                        highlightedIndex === index
                                            ? 'bg-(--color-blue)/10 text-(--color-blue)'
                                            : 'text-(--color-text-primary) hover:bg-(--color-muted-bg)'
                                    ]">
                                    <span class="flex items-center gap-2">
                                        <span>{{ suggestion }}</span>
                                        <span class="text-xs text-(--color-text-secondary) ml-auto">
                                            {{ getUsageCount(suggestion) }}×
                                        </span>
                                    </span>
                                </button>
                            </div>
                        </Transition>
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

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Income Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.income_date" type="date" required :disabled="!!justCreatedId"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Additional Information
            </h3>
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Notes</label>
                <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this income"
                    :disabled="!!justCreatedId"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
            </div>
        </div>

        <!-- Image -->
        <div class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Income Image <span class="normal-case text-xs font-normal">(Optional)</span>
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
                    <p v-else-if="justCreatedId" class="text-xs text-(--color-yellow)">
                        Income was created but the image upload failed. Retry below or remove the image.
                    </p>
                    <p v-else-if="isEditMode && props.income?.image_url" class="text-xs text-(--color-text-secondary)">
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
                <span v-else>{{ justCreatedId ? 'Retry Image Upload' : (isEditMode ? 'Save Changes' : 'Create Income')
                }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick, onUnmounted } from 'vue'
import type { Income } from '@/types/income'
import { useIncomesStore } from '@/stores/incomes'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const props = defineProps<{
    income?: Income | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'income-created': []
    'income-updated': []
    'cancel': []
}>()

const incomesStore = useIncomesStore()
const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit' || !!props.income)

// Autocomplete state
const dropdownRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const searchTerm = ref('')
const isDropdownOpen = ref(false)
const highlightedIndex = ref(-1)

const uniqueTitles = computed(() => {
    return incomesStore.incomes
        .map(e => e.title)
        .filter(t => t && t.trim())
        .filter((t, i, self) => self.indexOf(t) === i)
        .sort((a, b) => a.localeCompare(b))
})

const filteredSuggestions = computed(() => {
    if (!searchTerm.value.trim()) return uniqueTitles.value.slice(0, 10)
    const term = searchTerm.value.toLowerCase()
    return uniqueTitles.value.filter(t => t.toLowerCase().includes(term)).slice(0, 10)
})

const getUsageCount = (title: string): number => {
    if (!title) return 0
    return incomesStore.incomes.filter(e => e.title === title).length
}

const openDropdown = () => {
    if (!isEditMode.value && !justCreatedId.value && filteredSuggestions.value.length > 0) {
        isDropdownOpen.value = true
        highlightedIndex.value = -1
    }
}

const closeDropdown = () => {
    isDropdownOpen.value = false
    highlightedIndex.value = -1
}

const toggleDropdown = () => {
    if (isDropdownOpen.value) closeDropdown()
    else openDropdown()
}

const selectSuggestion = (suggestion: string) => {
    if (!suggestion) return
    searchTerm.value = suggestion
    form.value.title = suggestion
    closeDropdown()
    nextTick(() => {
        const amountInput = document.querySelector('input[type="number"]') as HTMLInputElement
        if (amountInput) amountInput.focus()
    })
}

const onInput = () => {
    form.value.title = searchTerm.value
    if (searchTerm.value.trim()) {
        isDropdownOpen.value = true
        highlightedIndex.value = -1
    } else {
        closeDropdown()
    }
}

const selectNext = () => {
    if (highlightedIndex.value < filteredSuggestions.value.length - 1) highlightedIndex.value++
}
const selectPrevious = () => {
    if (highlightedIndex.value > 0) highlightedIndex.value--
}
const handleEnter = () => {
    if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredSuggestions.value.length) {
        const s = filteredSuggestions.value[highlightedIndex.value]
        if (s) selectSuggestion(s)
    } else {
        closeDropdown()
    }
}

// FORM STATE
const form = ref({
    title: '',
    amount: null as number | null,
    income_date: '',
    notes: '',
})

const justCreatedId = ref<number | null>(null)

// IMAGE STATE + HELPERS
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

const initializeForm = () => {
    if (props.income) {
        const date = new Date(props.income.income_date)
        form.value = {
            title: props.income.title || '',
            amount: props.income.amount || null,
            income_date: date.toISOString().slice(0, 10),
            notes: props.income.notes || '',
        }
        searchTerm.value = props.income.title || ''

        justCreatedId.value = null
        clearImageState()
        if (isEditMode.value && props.income.image_url) {
            imagePreview.value = getImageUrl(props.income.image_url) || null
        }
    } else {
        const now = new Date()
        form.value = {
            title: '',
            amount: null,
            income_date: now.toISOString().slice(0, 10),
            notes: '',
        }
        searchTerm.value = ''
        justCreatedId.value = null
        clearImageState()
    }
}

watch(() => props.income, initializeForm, { immediate: true })

const resetForm = () => {
    initializeForm()
}

const submit = async () => {
    // RETRY IMAGE PATH
    if (justCreatedId.value !== null) {
        if (!imageFile.value) {
            resetForm()
            emit('income-created')
            return
        }
        uploading.value = true
        uploadError.value = null
        uploadProgress.value = 0
        try {
            const url = await uploadImage('incomes', justCreatedId.value, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                return
            }
            push.success('Income and image created successfully!')
            resetForm()
            emit('income-created')
        } finally {
            uploading.value = false
        }
        return
    }

    // NORMAL VALIDATION
    const title = form.value.title.trim()
    if (!title) {
        push.error('Title is required')
        return
    }
    if (!form.value.amount || form.value.amount <= 0) {
        push.error('Amount must be greater than 0')
        return
    }
    if (!form.value.income_date) {
        push.error('Income date is required')
        return
    }

    submitting.value = true

    try {
        let incomeId: number

        if (isEditMode.value && props.income) {
            const success = await incomesStore.updateIncome(props.income.id, {
                title,
                amount: form.value.amount,
                income_date: form.value.income_date,
                notes: form.value.notes.trim() || null,
            })

            if (!success) {
                submitting.value = false
                return
            }

            incomeId = props.income.id
            push.success('Income updated successfully!')
        } else {
            const newIncome = await incomesStore.createIncome({
                title,
                amount: form.value.amount,
                income_date: form.value.income_date,
                notes: form.value.notes.trim() || null,
            })

            if (!newIncome) {
                submitting.value = false
                return
            }

            incomeId = newIncome.id
            push.success('Income created successfully!')
        }

        // IMAGE HANDLING
        if (imageFile.value) {
            uploading.value = true
            uploadProgress.value = 0
            uploadError.value = null
            const url = await uploadImage('incomes', incomeId, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            uploading.value = false
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                if (!isEditMode.value) {
                    justCreatedId.value = incomeId
                }
                push.warning('Income saved but image upload failed. Retry below.')
                await incomesStore.fetchIncomes()
                return
            }
        } else if (isEditMode.value && props.income && imageToDelete.value && props.income.image_url) {
            await deleteImage('incomes', props.income.id, props.income.image_url)
        }

        if (isEditMode.value) {
            emit('income-updated')
        } else {
            emit('income-created')
        }

        resetForm()
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update income' : 'Failed to create income')
    } finally {
        submitting.value = false
    }
}

const handleClickOutside = (event: MouseEvent) => {
    if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
        closeDropdown()
    }
}

onMounted(() => {
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>