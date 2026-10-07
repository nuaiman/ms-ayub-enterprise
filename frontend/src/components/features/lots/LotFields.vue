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
                    @change="$emit('update:customerId', parseInt(($event.target as HTMLSelectElement).value) || 0)"
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
                    Lot Number <span class="text-(--color-red)">*</span>
                </label>
                <input :value="lotNumber" @input="$emit('update:lotNumber', ($event.target as HTMLInputElement).value)"
                    type="text" placeholder="e.g. A-001" :disabled="disabled" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                <p class="text-xs text-(--color-text-secondary) mt-1">
                    Must be unique per customer
                </p>
            </div>
        </div>

        <!-- Product Name -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Product Name <span class="text-(--color-red)">*</span>
            </label>
            <input :value="productName" @input="$emit('update:productName', ($event.target as HTMLInputElement).value)"
                type="text" placeholder="Enter product name" :disabled="disabled" required
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <!-- Units -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Weight Unit
                </label>
                <input :value="weightUnit"
                    @input="$emit('update:weightUnit', ($event.target as HTMLInputElement).value)" type="text"
                    placeholder="kg" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Quantity Unit
                </label>
                <input :value="quantityUnit"
                    @input="$emit('update:quantityUnit', ($event.target as HTMLInputElement).value)" type="text"
                    placeholder="units" :disabled="disabled"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Customer } from '@/types/customer'

defineProps<{
    customerId: number | null
    lotNumber: string
    productName: string
    weightUnit: string
    quantityUnit: string
    customerOptions?: Customer[]
    disabled?: boolean
    canEditCustomer?: boolean
}>()

defineEmits<{
    (e: 'update:customerId', value: number | null): void
    (e: 'update:lotNumber', value: string): void
    (e: 'update:productName', value: string): void
    (e: 'update:weightUnit', value: string): void
    (e: 'update:quantityUnit', value: string): void
}>()

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}
</script>