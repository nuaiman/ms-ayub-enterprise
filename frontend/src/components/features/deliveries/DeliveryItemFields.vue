<!-- src/components/features/deliveries/DeliveryItemFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Store -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Store <span class="text-(--color-red)">*</span>
            </label>
            <select :value="storeId ?? ''"
                @change="$emit('update:storeId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled" required
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option value="">Select a store</option>
                <option v-for="store in storeOptions" :key="store.id" :value="store.id">
                    {{ storeLabel(store) }}
                </option>
            </select>
        </div>

        <!-- Auto-derived lot display -->
        <div v-if="derivedLot" class="p-3 rounded-lg bg-(--color-blue)/5 border border-(--color-blue)/20">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Lot (derived from store)</p>
            <p class="text-sm font-semibold text-(--color-text-primary) mt-0.5">
                {{ derivedLot.product_name }} — Lot {{ derivedLot.lot_number }}
            </p>
        </div>

        <!-- Vehicle + Driver -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Vehicle Number
                </label>
                <input :value="vehicleNumber"
                    @input="$emit('update:vehicleNumber', ($event.target as HTMLInputElement).value)" type="text"
                    placeholder="e.g. DHAKA-METRO-1234" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Driver Number
                </label>
                <input :value="driverNumber"
                    @input="$emit('update:driverNumber', ($event.target as HTMLInputElement).value)" type="text"
                    placeholder="e.g. 01700000000" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Quantity + Weight -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Quantity
                    <span v-if="quantityUnit" class="text-xs font-normal text-(--color-text-secondary)">
                        ({{ quantityUnit }})
                    </span>
                </label>
                <input :value="quantity"
                    @input="$emit('update:quantity', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Weight
                    <span v-if="weightUnit" class="text-xs font-normal text-(--color-text-secondary)">
                        ({{ weightUnit }})
                    </span>
                </label>
                <input :value="weight"
                    @input="$emit('update:weight', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Customer Delivery Bill -->
        <div class="border-t border-(--color-border)/60 pt-4">
            <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">
                Customer Delivery Bill
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Bill Type <span class="text-(--color-red)">*</span>
                    </label>
                    <select :value="cdBillType"
                        @change="$emit('update:cdBillType', ($event.target as HTMLSelectElement).value as CustomerDeliveryBillType)"
                        :disabled="disabled" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="quantity">Quantity</option>
                        <option value="weight">Weight</option>
                    </select>
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Rate <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input :value="cdBillRate"
                            @input="$emit('update:cdBillRate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                            type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>
        </div>

        <!-- Majhi Bill -->
        <div class="border-t border-(--color-border)/60 pt-4">
            <h4 class="text-sm font-medium text-(--color-text-primary) mb-3">
                Majhi Bill
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Majhi <span class="text-(--color-red)">*</span>
                    </label>
                    <select :value="majhiId ?? ''"
                        @change="$emit('update:majhiId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                        :disabled="disabled" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select a majhi</option>
                        <option v-for="m in majhiOptions" :key="m.id" :value="m.id">
                            {{ m.name }}
                        </option>
                    </select>
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Bill Type <span class="text-(--color-red)">*</span>
                    </label>
                    <select :value="majhiBillType"
                        @change="$emit('update:majhiBillType', ($event.target as HTMLSelectElement).value as MajhiBillType)"
                        :disabled="disabled" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="quantity">Quantity</option>
                        <option value="weight">Weight</option>
                        <option value="job">Job (Fixed)</option>
                    </select>
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Rate <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input :value="majhiBillRate"
                            @input="$emit('update:majhiBillRate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                            type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Store } from '@/types/store'
import type { Majhi } from '@/types/majhi'
import type { CustomerDeliveryBillType } from '@/types/customerDeliveryBill'
import type { MajhiBillType } from '@/types/majhiBill'
import { useLotsStore } from '@/stores/lots'

const props = withDefaults(defineProps<{
    storeId?: number | null
    majhiId?: number | null
    vehicleNumber?: string
    driverNumber?: string
    quantity: number
    weight: number

    cdBillType: CustomerDeliveryBillType
    cdBillRate: number

    majhiBillType: MajhiBillType
    majhiBillRate: number

    storeOptions?: Store[]
    majhiOptions?: Majhi[]
    disabled?: boolean
}>(), {
    storeId: null,
    majhiId: null,
    vehicleNumber: '',
    driverNumber: '',
    storeOptions: () => [],
    majhiOptions: () => [],
    disabled: false,
})

defineEmits<{
    (e: 'update:storeId', value: number | null): void
    (e: 'update:majhiId', value: number | null): void
    (e: 'update:vehicleNumber', value: string): void
    (e: 'update:driverNumber', value: string): void
    (e: 'update:quantity', value: number): void
    (e: 'update:weight', value: number): void
    (e: 'update:cdBillType', value: CustomerDeliveryBillType): void
    (e: 'update:cdBillRate', value: number): void
    (e: 'update:majhiBillType', value: MajhiBillType): void
    (e: 'update:majhiBillRate', value: number): void
}>()

const lotsStore = useLotsStore()

const selectedStore = computed<Store | null>(() => {
    if (!props.storeId) return null
    return props.storeOptions.find(s => s.id === props.storeId) ?? null
})

const derivedLot = computed(() => {
    if (!selectedStore.value) return null
    return lotsStore.getLotById(selectedStore.value.lot_id) ?? null
})

const weightUnit = computed(() => derivedLot.value?.weight_unit ?? '')
const quantityUnit = computed(() => derivedLot.value?.quantity_unit ?? '')

const storeLabel = (store: Store): string => {
    const lot = lotsStore.getLotById(store.lot_id)
    const lotName = lot ? lot.product_name : '—'
    const lotNum = lot ? lot.lot_number : '—'
    return `Store #${store.id} — ${lotName} (Lot ${lotNum})`
}
</script>