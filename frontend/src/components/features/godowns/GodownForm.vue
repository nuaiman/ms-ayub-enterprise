<!-- src/components/features/godowns/GodownForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Godown Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Godown Information
            </h3>
            <div class="space-y-4">
                <!-- Name -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Name <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.name" type="text" placeholder="Enter godown name" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Phone -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Phone
                    </label>
                    <input v-model="form.phone" type="tel" placeholder="Enter phone number"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Monthly Rent -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Monthly Rent
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.monthly_rent" type="number" step="0.01" min="0" placeholder="0.00"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <!-- Is Active -->
                <div class="flex items-center gap-3">
                    <label class="text-sm font-medium text-(--color-text-primary)">Active</label>
                    <div @click="form.is_active = !form.is_active"
                        class="w-11 h-6 rounded-full cursor-pointer transition-colors duration-200 flex items-center px-0.5"
                        :class="form.is_active ? 'bg-(--color-green)' : 'bg-(--color-muted-bg)'">
                        <div class="w-5 h-5 rounded-full bg-white shadow-sm transition-transform duration-200"
                            :class="form.is_active ? 'translate-x-5' : 'translate-x-0'"></div>
                    </div>
                </div>

                <!-- Notes -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this godown"
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Godown' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Godown } from '@/types/godown'
import { useGodownsStore } from '@/stores/godowns'
import { push } from 'notivue'

const props = defineProps<{
    godown?: Godown | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'godown-created': []
    'godown-updated': []
    'cancel': []
}>()

const godownsStore = useGodownsStore()
const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit' || !!props.godown)

const form = ref({
    name: '',
    phone: '',
    monthly_rent: 0,
    is_active: true,
    notes: '',
})

const initializeForm = () => {
    if (props.godown) {
        form.value = {
            name: props.godown.name || '',
            phone: props.godown.phone || '',
            monthly_rent: props.godown.monthly_rent || 0,
            is_active: props.godown.is_active !== undefined ? props.godown.is_active : true,
            notes: props.godown.notes || '',
        }
    } else {
        form.value = {
            name: '',
            phone: '',
            monthly_rent: 0,
            is_active: true,
            notes: '',
        }
    }
}

watch(() => props.godown, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.godown) {
        initializeForm()
    } else {
        form.value = {
            name: '',
            phone: '',
            monthly_rent: 0,
            is_active: true,
            notes: '',
        }
    }
}

const submit = async () => {
    if (!form.value.name.trim()) {
        push.error('Name is required')
        return
    }

    if (form.value.monthly_rent < 0) {
        push.error('Monthly rent cannot be negative')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.godown) {
            const success = await godownsStore.updateGodown(props.godown.id, {
                name: form.value.name.trim(),
                phone: form.value.phone.trim() || null,
                monthly_rent: form.value.monthly_rent,
                is_active: form.value.is_active,
                notes: form.value.notes.trim() || null,
            })

            if (success) {
                push.success('Godown updated successfully!')
                emit('godown-updated')
            }
        } else {
            const newGodown = await godownsStore.createGodown({
                name: form.value.name.trim(),
                phone: form.value.phone.trim() || null,
                monthly_rent: form.value.monthly_rent,
                is_active: form.value.is_active,
                notes: form.value.notes.trim() || null,
            })

            if (newGodown) {
                push.success('Godown created successfully!')
                resetForm()
                emit('godown-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update godown' : 'Failed to create godown')
    } finally {
        submitting.value = false
    }
}
</script>