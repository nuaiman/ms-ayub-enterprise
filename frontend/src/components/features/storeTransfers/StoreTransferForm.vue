<!-- src/components/features/storeTransfers/StoreTransferForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Store picker (create only, not locked) -->
        <div v-if="!lockedStore">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Store <span class="text-(--color-red)">*</span>
            </label>
            <select v-model="form.store_id" required :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a store</option>
                <option v-for="store in activeStores" :key="store.id" :value="store.id">
                    {{ storeLabel(store) }}
                </option>
            </select>
        </div>

        <!-- Locked store display -->
        <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Store</p>
            <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedStoreLabel }}</p>
        </div>

        <!-- Current state snapshot -->
        <div v-if="currentStore" class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Current Godown</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                    {{ currentGodownName }}
                </p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Current Weight</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                    {{ formatNumber(currentStore.weight) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ weightUnit }}</span>
                </p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Current Quantity</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                    {{ formatNumber(currentStore.quantity) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ quantityUnit }}</span>
                </p>
            </div>
        </div>

        <!-- Target godown -->
        <div v-if="currentStore">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Transfer To <span class="text-(--color-red)">*</span>
            </label>
            <select v-model="form.to_godown_id" required :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select target godown</option>
                <option v-for="g in availableGodowns" :key="g.id" :value="g.id">
                    {{ g.name }}
                </option>
            </select>
            <p v-if="availableGodowns.length === 0" class="text-xs text-(--color-yellow) mt-1">
                No other godowns available.
            </p>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea v-model="form.notes" rows="3" placeholder="Optional notes about this transfer"
                :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || !form.store_id || !form.to_godown_id"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    Transferring...
                </span>
                <span v-else>Transfer Store</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import type { Store } from '@/types/store'
import { useStoreTransfersStore } from '@/stores/storeTransfers'
import { useStoresStore } from '@/stores/stores'
import { useGodownsStore } from '@/stores/godowns'
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
    'transfer-created': []
    'cancel': []
}>()

const storeTransfersStore = useStoreTransfersStore()
const storesStore = useStoresStore()
const godownsStore = useGodownsStore()
const lotsStore = useLotsStore()

const submitting = ref(false)

const form = ref({
    store_id: props.presetStoreId as number | null,
    to_godown_id: null as number | null,
    notes: '',
})

const activeStores = computed(() =>
    storesStore.stores.filter(s => s.is_active)
)

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

const currentGodownName = computed(() => {
    if (!currentStore.value) return '—'
    return godownsStore.getGodownName(currentStore.value.godown_id)
})

const availableGodowns = computed(() => {
    if (!currentStore.value) return godownsStore.godowns
    return godownsStore.godowns.filter(g => g.id !== currentStore.value!.godown_id)
})

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

// If the store changes and the previously selected target godown is now the
// current godown, reset the target.
watch(currentStore, (store) => {
    if (store && form.value.to_godown_id === store.godown_id) {
        form.value.to_godown_id = null
    }
})

const submit = async () => {
    if (!form.value.store_id) {
        push.error('Please select a store')
        return
    }
    if (!form.value.to_godown_id) {
        push.error('Please select a target godown')
        return
    }

    submitting.value = true

    try {
        const result = await storeTransfersStore.createTransfer(
            form.value.store_id,
            {
                to_godown_id: form.value.to_godown_id,
                notes: form.value.notes.trim() || null,
            }
        )

        if (result) {
            await storesStore.fetchStores()
            resetForm()
            emit('transfer-created')
        }
    } finally {
        submitting.value = false
    }
}

const resetForm = () => {
    form.value = {
        store_id: props.presetStoreId ?? null,
        to_godown_id: null,
        notes: '',
    }
}

onMounted(async () => {
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (godownsStore.godowns.length === 0) await godownsStore.fetchGodowns()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
})
</script>