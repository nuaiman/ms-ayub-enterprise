<!-- src/components/features/lots/LotFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Customer & Lot Number -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Customer -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Customer <span class="text-(--color-red)">*</span>
                </label>
                <select :value="customerId ?? ''"
                    @change="$emit('update:customerId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled || !canEditCustomer" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a customer</option>
                    <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                        {{ getCustomerDisplayName(customer) }}
                    </option>
                </select>
                <p v-if="!canEditCustomer" class="text-xs text-(--color-text-secondary) mt-1">
                    Customer cannot be changed
                </p>
            </div>

            <!-- Lot Number -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Lot Number
                </label>
                <input :value="lotNumber || ''" type="number" placeholder="Auto-generated" disabled
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed" />
                <p class="text-xs text-(--color-text-secondary) mt-1">
                    Auto-incremented per customer
                </p>
            </div>
        </div>

        <!-- Product Name & Category -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Product Name
                </label>
                <input :value="productName"
                    @input="$emit('update:productName', ($event.target as HTMLInputElement).value)" type="text"
                    placeholder="Enter product name" :disabled="disabled || !canEditProduct"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                <p v-if="!canEditProduct" class="text-xs text-(--color-text-secondary) mt-1">
                    Product name cannot be changed
                </p>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Category
                </label>
                <input :value="category" @input="$emit('update:category', ($event.target as HTMLInputElement).value)"
                    type="text" placeholder="Enter category" :disabled="disabled || !canEditProduct"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Charge Settings -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Customer Charge Type <span class="text-(--color-red)">*</span>
                </label>
                <select :value="customerChargeType"
                    @change="$emit('update:customerChargeType', ($event.target as HTMLSelectElement).value as CustomerChargeType)"
                    :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="quantity">Quantity</option>
                    <option value="weight">Weight</option>
                </select>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Majhi Bill Type <span class="text-(--color-red)">*</span>
                </label>
                <select :value="majhiBillType"
                    @change="$emit('update:majhiBillType', ($event.target as HTMLSelectElement).value as MajhiBillType)"
                    :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="quantity">Quantity</option>
                    <option value="weight">Weight</option>
                    <option value="job">Job</option>
                </select>
            </div>
        </div>

        <!-- Rates & Cuts -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Storage Rate
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="customerStorageRate"
                        @input="$emit('update:customerStorageRate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Unload Rate
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="unloadRate"
                        @input="$emit('update:unloadRate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
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

        <!-- Assigned Majhi -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Majhi
            </label>
            <select :value="majhiId ?? ''"
                @change="$emit('update:majhiId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option value="">Select a majhi</option>
                <option v-for="majhi in majhiOptions" :key="majhi.id" :value="majhi.id">
                    {{ majhi.name }}
                </option>
            </select>
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
                rows="3" placeholder="Enter any notes about this lot" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CustomerChargeType, MajhiBillType } from '@/types/lot'
import type { Customer } from '@/types/customer'
import type { Majhi } from '@/types/majhi'

defineProps<{
    customerId: number | null
    productName: string
    category: string
    lotNumber: number
    customerChargeType: CustomerChargeType
    majhiBillType: MajhiBillType
    customerStorageRate: number
    unloadRate: number
    majhiId: number | null
    majhiCut: number
    isActive: boolean
    notes: string
    customerOptions?: Customer[]
    majhiOptions?: Majhi[]
    disabled?: boolean
    canEditCustomer?: boolean
    canEditProduct?: boolean
    showActive?: boolean
}>()

defineEmits<{
    (e: 'update:customerId', value: number | null): void
    (e: 'update:productName', value: string): void
    (e: 'update:category', value: string): void
    (e: 'update:lotNumber', value: number): void
    (e: 'update:customerChargeType', value: CustomerChargeType): void
    (e: 'update:majhiBillType', value: MajhiBillType): void
    (e: 'update:customerStorageRate', value: number): void
    (e: 'update:unloadRate', value: number): void
    (e: 'update:majhiId', value: number | null): void
    (e: 'update:majhiCut', value: number): void
    (e: 'update:isActive', value: boolean): void
    (e: 'update:notes', value: string): void
}>()

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}
</script>