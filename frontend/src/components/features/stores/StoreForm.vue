<!-- src/components/features/stores/StoreForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Store Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Store Information
            </h3>

            <StoreFields v-model:lot-id="form.lot_id" v-model:godown-id="form.godown_id"
                v-model:store-bill-type="form.store_bill_type" v-model:godown-cut="form.godown_cut"
                v-model:quantity="form.quantity" v-model:quantity-unit="form.quantity_unit" v-model:weight="form.weight"
                v-model:weight-unit="form.weight_unit" v-model:is-active="form.is_active"
                v-model:billing-start="form.billing_start" v-model:billing-end="form.billing_end"
                v-model:notes="form.notes" :lot-options="lotOptions" :godown-options="godownOptions"
                :disabled="submitting" :required="true" :standalone="false" :can-edit-lot="!isEditMode && !isReaddMode"
                :can-edit-godown="!isEditMode && !isReaddMode" :show-active="true" />
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Store' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Store, StoreBillType } from '@/types/store'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useGodownsStore } from '@/stores/godowns'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'
import StoreFields from './StoreFields.vue'

const props = defineProps<{
    store?: Store | null
    prefill?: Store | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'store-created': []
    'store-updated': []
    'cancel': []
}>()

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()

const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit')
const isReaddMode = computed(() => !isEditMode.value && !!props.prefill)

const lotOptions = computed(() => {
    return lotsStore.lots.filter(l => l.is_active)
})

const godownOptions = computed(() => {
    return godownsStore.godowns.filter(g => g.is_active)
})

const form = ref({
    lot_id: null as number | null,
    godown_id: null as number | null,
    store_bill_type: 'quantity' as StoreBillType,
    godown_cut: 0,
    quantity: 0,
    quantity_unit: 'units',
    weight: 0,
    weight_unit: 'kg',
    is_active: true,
    billing_start: '',
    billing_end: null as string | null,
    notes: '',
})

const resetToBlank = () => {
    const today = new Date().toISOString().slice(0, 10)
    form.value = {
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

const seedFromStore = (source: Store, zeroStock: boolean) => {
    const billingStart = source.billing_start
        ? new Date(source.billing_start).toISOString().slice(0, 10)
        : new Date().toISOString().slice(0, 10)
    const billingEnd = source.billing_end
        ? new Date(source.billing_end).toISOString().slice(0, 10)
        : null

    form.value = {
        lot_id: source.lot_id,
        godown_id: source.godown_id,
        store_bill_type: source.store_bill_type,
        godown_cut: source.godown_cut,
        quantity: zeroStock ? 0 : source.quantity,
        quantity_unit: source.quantity_unit || 'units',
        weight: zeroStock ? 0 : source.weight,
        weight_unit: source.weight_unit || 'kg',
        is_active: source.is_active !== undefined ? source.is_active : true,
        billing_start: billingStart,
        billing_end: billingEnd,
        notes: source.notes || '',
    }
}

const initialize = () => {
    if (isEditMode.value && props.store) {
        seedFromStore(props.store, false)
    } else if (props.prefill) {
        seedFromStore(props.prefill, true)
    } else {
        resetToBlank()
    }
}

const ensureDependenciesLoaded = async () => {
    const tasks: Promise<unknown>[] = []

    if (lotsStore.lots.length === 0) tasks.push(lotsStore.fetchLots())
    if (godownsStore.godowns.length === 0) tasks.push(godownsStore.fetchGodowns())

    if (tasks.length > 0) {
        await Promise.all(tasks)
    }
}

watch([() => props.store, () => props.prefill, () => props.mode], initialize, { immediate: true })

const resetForm = () => {
    initialize()
}

const submit = async () => {
    if (!form.value.lot_id) {
        push.error('Please select a lot')
        return
    }

    if (!form.value.godown_id) {
        push.error('Please select a godown')
        return
    }

    if (!form.value.billing_start) {
        push.error('Billing start date is required')
        return
    }

    if (form.value.billing_end && form.value.billing_start > form.value.billing_end) {
        push.error('Billing start date must be before billing end date')
        return
    }

    submitting.value = true

    const billingStart = formatDateForBackend(form.value.billing_start)
    const billingEnd = formatDateForBackend(form.value.billing_end)

    try {
        if (isEditMode.value && props.store) {
            const success = await storesStore.updateStore(props.store.id, {
                store_bill_type: form.value.store_bill_type,
                godown_cut: form.value.godown_cut,
                quantity: form.value.quantity,
                quantity_unit: form.value.quantity_unit,
                weight: form.value.weight,
                weight_unit: form.value.weight_unit,
                billing_start: billingStart,
                billing_end: billingEnd,
                notes: form.value.notes.trim() || null,
            })

            if (success) {
                push.success('Store updated successfully!')
                emit('store-updated')
            }
        } else {
            const newStore = await storesStore.createStore({
                lot_id: form.value.lot_id,
                godown_id: form.value.godown_id,
                store_bill_type: form.value.store_bill_type,
                godown_cut: form.value.godown_cut,
                quantity: form.value.quantity,
                quantity_unit: form.value.quantity_unit,
                weight: form.value.weight,
                weight_unit: form.value.weight_unit,
                billing_start: billingStart,
                billing_end: billingEnd,
                notes: form.value.notes.trim() || null,
            })

            if (newStore) {
                push.success('Store created successfully!')
                resetForm()
                emit('store-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update store' : 'Failed to create store')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await ensureDependenciesLoaded()
    initialize()
})
</script>