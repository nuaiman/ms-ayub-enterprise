<!-- src/components/features/majhiBills/MajhiBillForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Majhi Bill Information
            </h3>

            <!-- Frozen snapshot overrides (edit mode only) -->
            <div v-if="isEditMode" class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Weight at Billing
                        <span v-if="unitLabel('weight')" class="text-xs font-normal text-(--color-text-secondary)">
                            ({{ unitLabel('weight') }})
                        </span>
                    </label>
                    <input v-model.number="form.weight_at_billing" type="number" step="0.01" min="0" placeholder="0.00"
                        :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Quantity at Billing
                        <span v-if="unitLabel('quantity')" class="text-xs font-normal text-(--color-text-secondary)">
                            ({{ unitLabel('quantity') }})
                        </span>
                    </label>
                    <input v-model.number="form.quantity_at_billing" type="number" step="0.01" min="0"
                        placeholder="0.00" :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <MajhiBillFields v-model:majhi-id="form.majhi_id" v-model:store-id="form.store_id"
                v-model:bill-type="form.bill_type" v-model:rate="form.rate" :majhi-options="majhiOptions"
                :store-options="storeOptions" :disabled="submitting" :locked="isEditMode"
                :locked-majhi-label="lockedMajhiLabel" :locked-store-label="lockedStoreLabel" />
        </div>

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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Bill' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { MajhiBill, MajhiBillType } from '@/types/majhiBill'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import { useMajhisStore } from '@/stores/majhis'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { push } from 'notivue'
import MajhiBillFields from './MajhiBillFields.vue'

const props = defineProps<{
    bill?: MajhiBill | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'bill-created': []
    'bill-updated': []
    'cancel': []
}>()

const majhiBillsStore = useMajhiBillsStore()
const majhisStore = useMajhisStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.bill)

const majhiOptions = computed(() => majhisStore.majhis)
const storeOptions = computed(() => storesStore.stores)

const form = ref({
    majhi_id: null as number | null,
    store_id: null as number | null,
    bill_type: 'quantity' as MajhiBillType,
    rate: 0,
    weight_at_billing: 0,
    quantity_at_billing: 0,
})

const lockedMajhiLabel = computed(() => {
    if (!props.bill) return ''
    return majhisStore.getMajhiName(props.bill.majhi_id)
})

const lockedStoreLabel = computed(() => {
    if (!props.bill) return ''

    const storeId = props.bill.store_id
    if (storeId === null || storeId === undefined) {
        return props.bill.delivery_item_id ? `Delivery Item #${props.bill.delivery_item_id}` : 'â€”'
    }

    const store = storesStore.getStoreById(storeId)
    if (!store) return `Store #${storeId}`
    const lot = lotsStore.getLotById(store.lot_id)
    if (lot) return `Store #${store.id} â€” ${lot.product_name} (Lot ${lot.lot_number})`
    return `Store #${store.id}`
})

const unitLabel = (kind: 'weight' | 'quantity'): string => {
    if (!props.bill) return ''
    if (kind === 'weight') return props.bill.weight_unit_at_billing || ''
    return props.bill.quantity_unit_at_billing || ''
}

const initialize = () => {
    if (props.bill) {
        form.value = {
            majhi_id: props.bill.majhi_id,
            store_id: props.bill.store_id,
            bill_type: props.bill.bill_type,
            rate: props.bill.rate,
            weight_at_billing: props.bill.weight_at_billing,
            quantity_at_billing: props.bill.quantity_at_billing,
        }
    } else {
        form.value = {
            majhi_id: null,
            store_id: null,
            bill_type: 'quantity',
            rate: 0,
            weight_at_billing: 0,
            quantity_at_billing: 0,
        }
    }
}

watch(() => props.bill, initialize, { immediate: true })

const resetForm = () => { initialize() }

const submit = async () => {
    if (!isEditMode.value) {
        if (!form.value.majhi_id) {
            push.error('Please select a majhi')
            return
        }
        if (!form.value.store_id) {
            push.error('Please select a store')
            return
        }
    }
    if (form.value.rate < 0) {
        push.error('Rate cannot be negative')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.bill) {
            const success = await majhiBillsStore.updateMajhiBill(props.bill.id, {
                bill_type: form.value.bill_type,
                rate: form.value.rate,
                weight_at_billing: form.value.weight_at_billing,
                quantity_at_billing: form.value.quantity_at_billing,
            })
            if (success) {
                push.success('Majhi bill updated successfully!')
                emit('bill-updated')
            }
        } else {
            const newBill = await majhiBillsStore.createMajhiBill({
                majhi_id: form.value.majhi_id!,
                store_id: form.value.store_id!,
                bill_type: form.value.bill_type,
                rate: form.value.rate,
            })
            if (newBill) {
                push.success('Majhi bill created successfully!')
                resetForm()
                emit('bill-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update majhi bill' : 'Failed to create majhi bill')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
    initialize()
})
</script>