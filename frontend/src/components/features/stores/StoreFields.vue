<!-- src/components/features/stores/StoreFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Lot (hidden when standalone - used in LotForm where lot is already created) -->
        <div v-if="!standalone" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Lot -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Lot <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <select :value="lotId"
                    @change="$emit('update:lotId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled || !canEditLot" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a lot</option>
                    <option v-for="lot in lotOptions" :key="lot.id" :value="lot.id">
                        {{ getLotDisplayName(lot) }}
                    </option>
                </select>
                <p v-if="!canEditLot" class="text-xs text-(--color-text-secondary) mt-1">
                    Lot cannot be changed
                </p>
            </div>

            <!-- Godown -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Godown <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <select :value="godownId"
                    @change="$emit('update:godownId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled || !canEditGodown" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a godown</option>
                    <option v-for="godown in godownOptions" :key="godown.id" :value="godown.id">
                        {{ godown.name }}
                    </option>
                </select>
                <p v-if="!canEditGodown" class="text-xs text-(--color-text-secondary) mt-1">
                    Godown cannot be changed
                </p>
            </div>
        </div>

        <!-- Godown ONLY (when standalone is true - used in LotForm) -->
        <div v-if="standalone" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Godown <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <select :value="godownId"
                    @change="$emit('update:godownId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled || !canEditGodown" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a godown</option>
                    <option v-for="godown in godownOptions" :key="godown.id" :value="godown.id">
                        {{ godown.name }}
                    </option>
                </select>
                <p v-if="!canEditGodown" class="text-xs text-(--color-text-secondary) mt-1">
                    Godown cannot be changed
                </p>
            </div>
        </div>

        <!-- Store Settings -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Bill Type
                </label>
                <select :value="storeBillType"
                    @change="$emit('update:storeBillType', ($event.target as HTMLSelectElement).value as StoreBillType)"
                    :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="quantity">Quantity</option>
                    <option value="weight">Weight</option>
                </select>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Godown Cut
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="godownCut"
                        @input="$emit('update:godownCut', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Inventory -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Quantity
                </label>
                <div class="flex gap-2">
                    <input :value="quantity"
                        @input="$emit('update:quantity', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0" :disabled="disabled"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    <input :value="quantityUnit"
                        @input="$emit('update:quantityUnit', ($event.target as HTMLInputElement).value)" type="text"
                        placeholder="Unit" :disabled="disabled"
                        class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Weight
                </label>
                <div class="flex gap-2">
                    <input :value="weight"
                        @input="$emit('update:weight', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0" :disabled="disabled"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    <input :value="weightUnit"
                        @input="$emit('update:weightUnit', ($event.target as HTMLInputElement).value)" type="text"
                        placeholder="Unit" :disabled="disabled"
                        class="w-24 px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Billing Period -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Billing Start <span v-if="required" class="text-(--color-red)">*</span>
                </label>
                <input :value="billingStart"
                    @input="$emit('update:billingStart', ($event.target as HTMLInputElement).value)" type="date"
                    :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Billing End
                </label>
                <input :value="billingEnd"
                    @input="$emit('update:billingEnd', ($event.target as HTMLInputElement).value || null)" type="date"
                    :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Is Active -->
        <div v-if="showActive" class="flex items-center gap-3">
            <label class="text-sm font-medium text-(--color-text-primary)">Active</label>
            <div @click="$emit('update:isActive', !isActive)"
                class="w-11 h-6 rounded-full cursor-pointer transition-colors duration-200 flex items-center px-0.5"
                :class="isActive ? 'bg-(--color-green)' : 'bg-(--color-muted-bg)'">
                <div class="w-5 h-5 rounded-full bg-white shadow-sm transition-transform duration-200"
                    :class="isActive ? 'translate-x-5' : 'translate-x-0'"></div>
            </div>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea :value="notes" @input="$emit('update:notes', ($event.target as HTMLTextAreaElement).value)"
                rows="2" placeholder="Enter any notes about this store" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { StoreBillType } from '@/types/store'
import type { Lot } from '@/types/lot'
import type { Godown } from '@/types/godown'
import { useLotsStore } from '@/stores/lots'

defineProps<{
    lotId: number | null
    godownId: number | null
    storeBillType: StoreBillType
    godownCut: number
    quantity: number
    quantityUnit: string
    weight: number
    weightUnit: string
    isActive: boolean
    billingStart: string
    billingEnd: string | null
    notes: string
    lotOptions?: Lot[]
    godownOptions?: Godown[]
    disabled?: boolean
    required?: boolean
    standalone?: boolean
    canEditLot?: boolean
    canEditGodown?: boolean
    showActive?: boolean
}>()

defineEmits<{
    (e: 'update:lotId', value: number | null): void
    (e: 'update:godownId', value: number | null): void
    (e: 'update:storeBillType', value: StoreBillType): void
    (e: 'update:godownCut', value: number): void
    (e: 'update:quantity', value: number): void
    (e: 'update:quantityUnit', value: string): void
    (e: 'update:weight', value: number): void
    (e: 'update:weightUnit', value: string): void
    (e: 'update:isActive', value: boolean): void
    (e: 'update:billingStart', value: string): void
    (e: 'update:billingEnd', value: string | null): void
    (e: 'update:notes', value: string): void
}>()

const lotsStore = useLotsStore()

const getLotDisplayName = (lot: Lot): string => {
    return `${lotsStore.getLotDisplayName(lot)} — Lot #${lot.lot_number}`
}
</script>