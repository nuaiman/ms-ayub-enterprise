<!-- src/components/features/deliveryItems/DeliveryItemFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Store / Lot Hierarchy -->
        <div>
            <h4 class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">
                Select Store → Lot
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Store -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Store <span v-if="required" class="text-(--color-red)">*</span>
                    </label>
                    <select :value="storeId" @change="onStoreChange" :disabled="disabled || !canEditStore"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">
                            {{ !customerId ? 'Select customer first' : 'Select a store' }}
                        </option>
                        <option v-for="store in availableStores" :key="store.id" :value="store.id">
                            {{ getStoreOptionLabel(store) }}
                        </option>
                    </select>
                    <p v-if="!customerId" class="text-xs text-(--color-yellow) mt-1">
                        Please select a customer first
                    </p>
                    <p v-if="availableStores.length === 0 && customerId" class="text-xs text-(--color-yellow) mt-1">
                        No stores found for this customer
                    </p>
                    <p v-if="!canEditStore" class="text-xs text-(--color-text-secondary) mt-1">
                        Store cannot be changed
                    </p>
                </div>

                <!-- Lot (auto-derived from store) -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Lot
                    </label>
                    <select :value="lotId" @change="onLotChange" :disabled="disabled || !storeId || !canEditLot"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">
                            {{ !storeId ? 'Select store first' : 'Lot for this store' }}
                        </option>
                        <option v-for="lot in filteredLots" :key="lot.id" :value="lot.id">
                            {{ getLotDisplayName(lot) }}
                        </option>
                    </select>
                    <p v-if="!canEditLot" class="text-xs text-(--color-text-secondary) mt-1">
                        Lot cannot be changed
                    </p>
                </div>
            </div>
        </div>

        <!-- Store Info -->
        <div v-if="selectedStore" class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Quantity</p>
                    <p class="font-medium text-(--color-text-primary)">{{ selectedStore.quantity }} {{
                        selectedStore.quantity_unit }}</p>
                </div>
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Weight</p>
                    <p class="font-medium text-(--color-text-primary)">{{ selectedStore.weight }} {{
                        selectedStore.weight_unit }}</p>
                </div>
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Bill Type</p>
                    <p class="font-medium capitalize text-(--color-text-primary)">{{ selectedStore.store_bill_type }}
                    </p>
                </div>
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Godown Cut</p>
                    <p class="font-medium text-(--color-text-primary)">{{ formatCurrency(selectedStore.godown_cut) }}
                    </p>
                </div>
            </div>
        </div>

        <!-- Majhi -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Majhi
            </label>
            <select :value="majhiId"
                @change="$emit('update:majhiId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a majhi</option>
                <option v-for="majhi in majhiOptions" :key="majhi.id" :value="majhi.id">
                    {{ majhi.name }}
                </option>
            </select>
        </div>

        <!-- Quantity & Weight -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Quantity <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <div class="flex gap-2">
                    <input :value="quantity"
                        @input="$emit('update:quantity', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" :max="selectedStore?.quantity || 0" placeholder="0"
                        :disabled="disabled || !selectedStore"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    <input :value="quantityUnit"
                        @input="$emit('update:quantityUnit', ($event.target as HTMLInputElement).value)" type="text"
                        placeholder="Unit" :disabled="disabled"
                        class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <p v-if="selectedStore" class="text-xs text-(--color-text-secondary) mt-1">
                    Available: {{ selectedStore.quantity }} {{ selectedStore.quantity_unit }}
                </p>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Weight <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <div class="flex gap-2">
                    <input :value="weight"
                        @input="$emit('update:weight', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" :max="selectedStore?.weight || 0" placeholder="0"
                        :disabled="disabled || !selectedStore"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    <input :value="weightUnit"
                        @input="$emit('update:weightUnit', ($event.target as HTMLInputElement).value)" type="text"
                        placeholder="Unit" :disabled="disabled"
                        class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <p v-if="selectedStore" class="text-xs text-(--color-text-secondary) mt-1">
                    Available: {{ selectedStore.weight }} {{ selectedStore.weight_unit }}
                </p>
            </div>
        </div>

        <!-- Rates -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Loading Rate
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="loadingRate"
                        @input="$emit('update:loadingRate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Majhi Cut
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="majhiCut"
                        @input="$emit('update:majhiCut', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Vehicle & Driver -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Vehicle Number
                </label>
                <input :value="vehicleNumber"
                    @input="$emit('update:vehicleNumber', ($event.target as HTMLInputElement).value)" type="text"
                    placeholder="Enter vehicle number" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Driver Number
                </label>
                <input :value="driverNumber"
                    @input="$emit('update:driverNumber', ($event.target as HTMLInputElement).value)" type="tel"
                    placeholder="Enter driver phone" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea :value="notes" @input="$emit('update:notes', ($event.target as HTMLTextAreaElement).value)"
                rows="2" placeholder="Enter any notes about this delivery item" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Lot } from '@/types/lot'
import type { Store } from '@/types/store'
import type { Majhi } from '@/types/majhi'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    customerId: number | null
    lotId: number | null
    storeId: number | null
    majhiId: number | null
    quantity: number
    quantityUnit: string
    weight: number
    weightUnit: string
    loadingRate: number
    majhiCut: number
    vehicleNumber: string
    driverNumber: string
    notes: string
    lotOptions?: Lot[]
    storeOptions?: Store[]
    majhiOptions?: Majhi[]
    disabled?: boolean
    required?: boolean
    standalone?: boolean
    canEditLot?: boolean
    canEditStore?: boolean
}>()

const emit = defineEmits<{
    (e: 'update:lotId', value: number | null): void
    (e: 'update:storeId', value: number | null): void
    (e: 'update:majhiId', value: number | null): void
    (e: 'update:quantity', value: number): void
    (e: 'update:quantityUnit', value: string): void
    (e: 'update:weight', value: number): void
    (e: 'update:weightUnit', value: string): void
    (e: 'update:loadingRate', value: number): void
    (e: 'update:majhiCut', value: number): void
    (e: 'update:vehicleNumber', value: string): void
    (e: 'update:driverNumber', value: string): void
    (e: 'update:notes', value: string): void
}>()

const lotsStore = useLotsStore()
const storesStore = useStoresStore()

const availableStores = computed(() => {
    let stores = props.storeOptions && props.storeOptions.length > 0
        ? props.storeOptions
        : storesStore.stores

    stores = stores.filter(s => s.is_active && (s.quantity > 0 || s.weight > 0))

    if (props.customerId) {
        stores = stores.filter(s => {
            const lot = lotsStore.getLotById(s.lot_id)
            if (!lot) return false
            return lot.customer_id === props.customerId
        })
    }

    return stores
})

const filteredLots = computed(() => {
    if (!props.storeId) return []
    const store = storesStore.getStoreById(props.storeId)
    if (!store) return []
    const lot = lotsStore.getLotById(store.lot_id)
    if (!lot) return []
    return [lot]
})

const selectedStore = computed(() => {
    if (!props.storeId) return null
    const stores = props.storeOptions && props.storeOptions.length > 0
        ? props.storeOptions
        : storesStore.stores
    return stores.find(s => s.id === props.storeId) || null
})

const getLotDisplayName = (lot: Lot): string => {
    return `${lotsStore.getLotDisplayName(lot)} — Lot #${lot.lot_number}`
}

const getStoreDisplayName = (store: Store): string => {
    return storesStore.getStoreDisplayName(store)
}

const getStoreOptionLabel = (store: Store): string => {
    const base = getStoreDisplayName(store)
    const qty = `${store.quantity} ${store.quantity_unit}`
    const wt = `${store.weight} ${store.weight_unit}`
    return `${base} (Qty: ${qty}, Wt: ${wt})`
}

const onStoreChange = (event: Event) => {
    const target = event.target as HTMLSelectElement
    const value = target.value ? parseInt(target.value) : null
    emit('update:storeId', value)

    if (value) {
        const store = storesStore.getStoreById(value)
        if (store) {
            emit('update:lotId', store.lot_id)
        } else {
            emit('update:lotId', null)
        }
    } else {
        emit('update:lotId', null)
    }

    emit('update:quantity', 0)
    emit('update:weight', 0)
}

const onLotChange = (event: Event) => {
    const target = event.target as HTMLSelectElement
    const value = target.value ? parseInt(target.value) : null
    emit('update:lotId', value)

    emit('update:quantity', 0)
    emit('update:weight', 0)
}
</script>