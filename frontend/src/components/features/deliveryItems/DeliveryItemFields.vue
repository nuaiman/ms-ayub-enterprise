<!-- src/components/features/deliveryItems/DeliveryItemFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Item, Lot, Store Hierarchy -->
        <div>
            <h4 class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-3">
                Select Item → Lot → Store
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <!-- Item -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Item <span v-if="required" class="text-(--color-red)">*</span>
                    </label>
                    <select :value="itemId" @change="onItemChange" :disabled="disabled || !canEditItem"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">
                            {{ !customerId ? 'Select customer first' : 'Select an item' }}
                        </option>
                        <option v-for="item in availableItems" :key="item.id" :value="item.id">
                            {{ getItemDisplayName(item) }}
                        </option>
                    </select>
                    <p v-if="!customerId" class="text-xs text-(--color-yellow) mt-1">
                        Please select a customer first
                    </p>
                    <p v-if="availableItems.length === 0 && customerId" class="text-xs text-(--color-yellow) mt-1">
                        No items found for this customer
                    </p>
                    <p v-if="!canEditItem" class="text-xs text-(--color-text-secondary) mt-1">
                        Item cannot be changed
                    </p>
                </div>

                <!-- Lot -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Lot <span v-if="required" class="text-(--color-red)">*</span>
                    </label>
                    <select :value="lotId" @change="onLotChange" :disabled="disabled || !itemId || !canEditLot"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">
                            {{ !itemId ? 'Select item first' : 'Select a lot' }}
                        </option>
                        <option v-for="lot in filteredLots" :key="lot.id" :value="lot.id">
                            {{ getLotDisplayName(lot) }}
                        </option>
                    </select>
                    <p v-if="!canEditLot" class="text-xs text-(--color-text-secondary) mt-1">
                        Lot cannot be changed
                    </p>
                    <p v-if="filteredLots.length === 0 && itemId" class="text-xs text-(--color-yellow) mt-1">
                        No active lots for this item
                    </p>
                </div>

                <!-- Store -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Store <span v-if="required" class="text-(--color-red)">*</span>
                    </label>
                    <select :value="storeId" @change="onStoreChange" :disabled="disabled || !lotId || !canEditStore"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">
                            {{ !lotId ? 'Select lot first' : 'Select a store' }}
                        </option>
                        <option v-for="store in filteredStores" :key="store.id" :value="store.id">
                            {{ getStoreDisplayName(store) }}
                        </option>
                    </select>
                    <p v-if="!canEditStore" class="text-xs text-(--color-text-secondary) mt-1">
                        Store cannot be changed
                    </p>
                    <p v-if="filteredStores.length === 0 && lotId" class="text-xs text-(--color-yellow) mt-1">
                        No active stores with inventory for this lot
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
            <!-- Quantity -->
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

            <!-- Weight -->
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
            <!-- Loading Rate -->
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

            <!-- Majhi Cut -->
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

        <!-- HIDDEN FIELDS: customer_charge_type, customer_paid_unload_amount, majhi_bill_type, majhi_total_paid -->
        <!-- These are automatically populated from the lot and NOT displayed in the form -->
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Item } from '@/types/item'
import type { Lot } from '@/types/lot'
import type { Store } from '@/types/store'
import type { Majhi } from '@/types/majhi'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { formatCurrency } from '@/utils/currency'

// Props
const props = defineProps<{
    customerId: number | null
    itemId: number | null
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
    itemOptions?: Item[]
    lotOptions?: Lot[]
    storeOptions?: Store[]
    majhiOptions?: Majhi[]
    disabled?: boolean
    required?: boolean
    standalone?: boolean
    canEditItem?: boolean
    canEditLot?: boolean
    canEditStore?: boolean
}>()

// Emits
const emit = defineEmits<{
    (e: 'update:itemId', value: number | null): void
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

const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()

// Available items (filtered by customer)
const availableItems = computed(() => {
    let items = props.itemOptions && props.itemOptions.length > 0
        ? props.itemOptions
        : itemsStore.items

    // Filter by customer if customerId is provided
    if (props.customerId) {
        items = items.filter(i => i.customer_id === props.customerId && i.is_active)
    } else {
        items = items.filter(i => i.is_active)
    }

    return items
})

// Filtered lots based on selected item
const filteredLots = computed(() => {
    if (!props.itemId) return []
    const lots = props.lotOptions && props.lotOptions.length > 0
        ? props.lotOptions
        : lotsStore.lots
    return lots.filter(l => l.item_id === props.itemId && l.is_active)
})

// Filtered stores based on selected lot
const filteredStores = computed(() => {
    if (!props.lotId) return []
    const stores = props.storeOptions && props.storeOptions.length > 0
        ? props.storeOptions
        : storesStore.stores
    return stores.filter(s => s.lot_id === props.lotId && s.is_active && (s.quantity > 0 || s.weight > 0))
})

// Selected store details
const selectedStore = computed(() => {
    if (!props.storeId) return null
    const stores = props.storeOptions && props.storeOptions.length > 0
        ? props.storeOptions
        : storesStore.stores
    return stores.find(s => s.id === props.storeId) || null
})

// Get display names
const getItemDisplayName = (item: Item): string => {
    return itemsStore.getItemDisplayName(item)
}

const getLotDisplayName = (lot: Lot): string => {
    return lotsStore.getLotDisplayName(lot)
}

const getStoreDisplayName = (store: Store): string => {
    return storesStore.getStoreDisplayName(store)
}

// Change handlers
const onItemChange = (event: Event) => {
    const target = event.target as HTMLSelectElement
    const value = target.value ? parseInt(target.value) : null
    emit('update:itemId', value)
    emit('update:lotId', null)
    emit('update:storeId', null)
    emit('update:quantity', 0)
    emit('update:weight', 0)
}

const onLotChange = (event: Event) => {
    const target = event.target as HTMLSelectElement
    const value = target.value ? parseInt(target.value) : null
    emit('update:lotId', value)
    emit('update:storeId', null)
    emit('update:quantity', 0)
    emit('update:weight', 0)
}

const onStoreChange = (event: Event) => {
    const target = event.target as HTMLSelectElement
    const value = target.value ? parseInt(target.value) : null
    emit('update:storeId', value)
}
</script>