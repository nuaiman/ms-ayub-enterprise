<!-- src/components/features/storeAdjustments/StoreAdjustmentForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Store picker (create only) -->
        <div v-if="!lockedStore">
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

        <!-- Locked store display -->
        <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Store</p>
            <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedStoreLabel }}</p>
        </div>

        <!-- Current snapshot -->
        <div v-if="currentStore" class="grid grid-cols-2 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Current Weight</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">
                    {{ formatNumber(currentStore.weight) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ weightUnit }}</span>
                </p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Current Quantity</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">
                    {{ formatNumber(currentStore.quantity) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ quantityUnit }}</span>
                </p>
            </div>
        </div>

        <!-- Adjustment type toggle -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Adjustment Type <span class="text-(--color-red)">*</span>
            </label>
            <div class="grid grid-cols-2 gap-2 p-1 rounded-lg bg-(--color-muted-bg) border border-(--color-border)">
                <button type="button" :disabled="submitting" @click="form.adjustment_type = 'delta'"
                    class="px-3 py-2 rounded-md text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="form.adjustment_type === 'delta'
                        ? 'bg-(--color-surface) text-(--color-text-primary) shadow-sm'
                        : 'text-(--color-text-secondary) hover:text-(--color-text-primary)'">
                    Change By (Delta)
                </button>
                <button type="button" :disabled="submitting" @click="form.adjustment_type = 'absolute'"
                    class="px-3 py-2 rounded-md text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="form.adjustment_type === 'absolute'
                        ? 'bg-(--color-surface) text-(--color-text-primary) shadow-sm'
                        : 'text-(--color-text-secondary) hover:text-(--color-text-primary)'">
                    Set To (Absolute)
                </button>
            </div>
            <p class="text-xs text-(--color-text-secondary) mt-1.5">
                {{ form.adjustment_type === 'delta'
                    ? 'Enter the change to apply. Use negative values to decrease.'
                    : 'Enter the new target value. Zero means "leave unchanged".' }}
            </p>
        </div>

        <!-- Weight + Quantity -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Weight
                    <span v-if="weightUnit" class="text-xs font-normal text-(--color-text-secondary)">
                        ({{ weightUnit }})
                    </span>
                </label>
                <input v-model.number="form.input_weight" type="number" step="0.01"
                    :min="form.adjustment_type === 'absolute' ? 0 : undefined" placeholder="0.00" :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Quantity
                    <span v-if="quantityUnit" class="text-xs font-normal text-(--color-text-secondary)">
                        ({{ quantityUnit }})
                    </span>
                </label>
                <input v-model.number="form.input_quantity" type="number" step="0.01"
                    :min="form.adjustment_type === 'absolute' ? 0 : undefined" placeholder="0.00" :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Reason -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Reason
            </label>
            <input v-model="form.reason" type="text" placeholder="e.g. Damage, spillage, recount" :disabled="submitting"
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
                    Recording...
                </span>
                <span v-else>Record Adjustment</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import type { Store } from '@/types/store'
import type { StoreAdjustmentType } from '@/types/storeAdjustment'
import { useStoreAdjustmentsStore } from '@/stores/storeAdjustments'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { push } from 'notivue'

const props = withDefaults(defineProps<{
    presetStoreId?: number | null
    lockedStore?: boolean
}>(), {
    presetStoreId: null,
    lockedStore: false,
})

const emit = defineEmits<{
    'adjustment-created': []
    'cancel': []
}>()

const storeAdjustmentsStore = useStoreAdjustmentsStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()

const submitting = ref(false)

const form = ref({
    store_id: props.presetStoreId as number | null,
    adjustment_type: 'delta' as StoreAdjustmentType,
    input_weight: 0,
    input_quantity: 0,
    reason: '',
    notes: '',
})

const storeOptions = computed(() => storesStore.stores)

const currentStore = computed<Store | null>(() => {
    if (!form.value.store_id) return null
    return storesStore.getStoreById(form.value.store_id) ?? null
})

const currentLot = computed(() => {
    if (!currentStore.value) return null
    return lotsStore.getLotById(currentStore.value.lot_id) ?? null
})

const weightUnit = computed(() => currentLot.value?.weight_unit ?? '')
const quantityUnit = computed(() => currentLot.value?.quantity_unit ?? '')

const lockedStoreLabel = computed(() => {
    if (!currentStore.value) return '—'
    const lot = currentLot.value
    if (!lot) return `Store #${currentStore.value.id}`
    return `Store #${currentStore.value.id} — ${lot.product_name} (Lot ${lot.lot_number})`
})

const storeLabel = (store: Store): string => {
    const lot = lotsStore.getLotById(store.lot_id)
    if (!lot) return `Store #${store.id}`
    return `Store #${store.id} — ${lot.product_name} (Lot ${lot.lot_number})`
}

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)

watch(() => props.presetStoreId, (val) => {
    if (val) form.value.store_id = val
}, { immediate: true })

const submit = async () => {
    if (!form.value.store_id) {
        push.error('Please select a store')
        return
    }

    if (form.value.input_weight === 0 && form.value.input_quantity === 0) {
        push.error('At least one of weight or quantity must be non-zero')
        return
    }

    if (form.value.adjustment_type === 'absolute') {
        if (form.value.input_weight < 0 || form.value.input_quantity < 0) {
            push.error('Absolute values cannot be negative')
            return
        }
    }

    submitting.value = true

    try {
        const result = await storeAdjustmentsStore.createAdjustment(
            form.value.store_id,
            {
                adjustment_type: form.value.adjustment_type,
                input_weight: form.value.input_weight,
                input_quantity: form.value.input_quantity,
                reason: form.value.reason.trim() || null,
                notes: form.value.notes.trim() || null,
            }
        )

        if (result) {
            // Refresh stores so the snapshot reflects the new values.
            await storesStore.fetchStores()
            resetForm()
            emit('adjustment-created')
        }
    } finally {
        submitting.value = false
    }
}

const resetForm = () => {
    form.value = {
        store_id: props.presetStoreId ?? null,
        adjustment_type: 'delta',
        input_weight: 0,
        input_quantity: 0,
        reason: '',
        notes: '',
    }
}

onMounted(async () => {
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
})
</script>