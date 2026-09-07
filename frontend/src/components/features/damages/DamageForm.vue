<!-- src/components/features/damages/DamageForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Damage Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Damage Information
            </h3>
            <div class="space-y-4">
                <!-- Store -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Store <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.store_id" required :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a store</option>
                        <option v-for="store in storeOptions" :key="store.id" :value="store.id">
                            {{ getStoreDisplayName(store) }}
                        </option>
                    </select>
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Store cannot be changed</p>
                </div>

                <!-- Reason -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Reason <span class="text-(--color-red)">*</span>
                    </label>
                    <textarea v-model="form.reason" rows="2" placeholder="Enter the reason for damage" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>

                <!-- Damage Date -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Damage Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.damage_date" type="date" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
            </div>
        </div>

        <!-- Inventory & Amount -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Inventory & Amount
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <!-- Quantity -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Quantity
                    </label>
                    <div class="flex gap-2">
                        <input v-model.number="form.quantity" type="number" step="0.01" min="0" placeholder="0"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        <input v-model="form.quantity_unit" type="text" placeholder="Unit"
                            class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <!-- Weight -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Weight
                    </label>
                    <div class="flex gap-2">
                        <input v-model.number="form.weight" type="number" step="0.01" min="0" placeholder="0"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        <input v-model="form.weight_unit" type="text" placeholder="Unit"
                            class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
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
            </div>
            <p class="text-xs text-(--color-text-secondary) mt-2">Either quantity or weight must be greater than 0</p>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea v-model="form.notes" rows="2" placeholder="Enter any additional notes"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Damage Record' }}</span>
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

const storeOptions = computed(() => {
    return storesStore.stores
})

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
    }
}

watch(() => props.damage, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.damage) {
        initializeForm()
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
    }
}

const submit = async () => {
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

    // Format date for backend
    const damageDate = formatDateForBackend(form.value.damage_date)

    try {
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

            if (success) {
                push.success('Damage record updated successfully!')
                emit('damage-updated')
            }
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

            if (newDamage) {
                push.success('Damage record created successfully!')
                resetForm()
                emit('damage-created')
            }
        }
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