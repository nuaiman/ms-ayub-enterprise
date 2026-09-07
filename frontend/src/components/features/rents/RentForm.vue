<!-- src/components/features/rents/RentForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Rent Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Rent Information
            </h3>
            <div class="space-y-4">
                <!-- Godown -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Godown <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.godown_id" required :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a godown</option>
                        <option v-for="godown in godownOptions" :key="godown.id" :value="godown.id">
                            {{ godown.name }}
                        </option>
                    </select>
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Godown cannot be changed</p>
                    <p v-if="selectedGodown && !isEditMode" class="text-xs text-(--color-text-secondary) mt-1">
                        Monthly rent: {{ formatCurrency(selectedGodown.monthly_rent) }}
                    </p>
                </div>

                <!-- Month -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Month <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.month_year" type="month" required :max="currentMonth" :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Month cannot be changed</p>
                </div>

                <!-- Amount -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Amount <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.amount" type="number" step="0.01" min="0" placeholder="0.00"
                            required
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <!-- Notes -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this rent"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Rent' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Rent } from '@/types/rent'
import { useRentsStore } from '@/stores/rents'
import { useGodownsStore } from '@/stores/godowns'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'

const props = defineProps<{
    rent?: Rent | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'rent-created': []
    'rent-updated': []
    'cancel': []
}>()

const rentsStore = useRentsStore()
const godownsStore = useGodownsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.rent)
const currentMonth = new Date().toISOString().slice(0, 7)

const godownOptions = computed(() => {
    return godownsStore.activeGodowns
})

const selectedGodown = computed(() => {
    if (!form.value.godown_id) return null
    return godownsStore.getGodownById(form.value.godown_id)
})

const form = ref({
    godown_id: null as number | null,
    month_year: '',
    amount: 0,
    notes: '',
})

const initializeForm = () => {
    if (props.rent) {
        form.value = {
            godown_id: props.rent.godown_id,
            month_year: props.rent.month_year,
            amount: props.rent.amount,
            notes: props.rent.notes || '',
        }
    } else {
        form.value = {
            godown_id: null,
            month_year: currentMonth,
            amount: 0,
            notes: '',
        }
    }
}

watch(() => props.rent, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.rent) {
        initializeForm()
    } else {
        form.value = {
            godown_id: null,
            month_year: currentMonth,
            amount: 0,
            notes: '',
        }
    }
}

const submit = async () => {
    if (!form.value.godown_id) {
        push.error('Please select a godown')
        return
    }

    if (!form.value.month_year) {
        push.error('Please select a month')
        return
    }

    if (form.value.amount < 0) {
        push.error('Amount cannot be negative')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.rent) {
            const success = await rentsStore.updateRent(props.rent.id, {
                amount: form.value.amount,
                notes: form.value.notes.trim() || null,
            })

            if (success) {
                push.success('Rent updated successfully!')
                emit('rent-updated')
            }
        } else {
            const newRent = await rentsStore.createRent({
                godown_id: form.value.godown_id,
                month_year: form.value.month_year,
                amount: form.value.amount,
                notes: form.value.notes.trim() || null,
            })

            if (newRent) {
                push.success('Rent created successfully!')
                resetForm()
                emit('rent-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update rent' : 'Failed to create rent')
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    if (godownsStore.godowns.length === 0) {
        godownsStore.fetchGodowns()
    }
})
</script>