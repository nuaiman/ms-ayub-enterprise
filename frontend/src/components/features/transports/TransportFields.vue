<!-- src/components/features/transports/TransportFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Customer -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Customer
            </label>
            <select :value="customerId"
                @change="$emit('update:customerId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a customer</option>
                <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                    {{ getCustomerDisplayName(customer) }}
                </option>
            </select>
        </div>

        <!-- From Location -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                From Location <span v-if="required" class="text-(--color-red)">*</span>
            </label>
            <input :value="fromLocation"
                @input="$emit('update:fromLocation', ($event.target as HTMLInputElement).value)" type="text"
                placeholder="Enter origin location" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <!-- To Location -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                To Location
            </label>
            <input :value="toLocation" @input="$emit('update:toLocation', ($event.target as HTMLInputElement).value)"
                type="text" placeholder="Enter destination location" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <!-- Vehicle Quantity -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Vehicle Quantity <span v-if="required" class="text-(--color-red)">*</span>
            </label>
            <input :value="vehicleQuantity"
                @input="$emit('update:vehicleQuantity', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                type="number" step="1" min="0" placeholder="0" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <!-- Delivery Type -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Delivery Type
            </label>
            <select :value="deliveryType"
                @change="$emit('update:deliveryType', ($event.target as HTMLSelectElement).value || null)"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select type</option>
                <option value="local">Local</option>
                <option value="district">District</option>
            </select>
        </div>

        <!-- Transport Date -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Transport Date <span v-if="required" class="text-(--color-red)">*</span>
            </label>
            <input :value="transportDate"
                @input="$emit('update:transportDate', ($event.target as HTMLInputElement).value)" type="date"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <!-- Office Commission Amount -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Office Commission Amount
            </label>
            <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                <input :value="officeCommissionAmount"
                    @input="$emit('update:officeCommissionAmount', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                    class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
        </div>

        <!-- Customer Total Paid - ONLY shown in standalone mode (for future billing) -->
        <div v-if="standalone" class="border-t border-(--color-border) pt-4">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Customer Total Paid
            </label>
            <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                <input :value="customerTotalPaid"
                    @input="$emit('update:customerTotalPaid', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                    type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                    class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
            <p class="text-[10px] text-(--color-text-secondary) mt-0.5">Total paid by customer for this transport</p>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea :value="notes" @input="$emit('update:notes', ($event.target as HTMLTextAreaElement).value)"
                rows="3" placeholder="Enter any notes about this transport" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Customer } from '@/types/customer'

// Props
defineProps<{
    customerId: number | null
    fromLocation: string
    toLocation: string
    vehicleQuantity: number
    deliveryType: string | null
    transportDate: string
    officeCommissionAmount: number
    customerTotalPaid?: number  // Only shown in standalone mode
    notes: string
    customerOptions?: Customer[]
    disabled?: boolean
    required?: boolean
    standalone?: boolean  // true = shows customer_total_paid field
}>()

// Emits
defineEmits<{
    (e: 'update:customerId', value: number | null): void
    (e: 'update:fromLocation', value: string): void
    (e: 'update:toLocation', value: string): void
    (e: 'update:vehicleQuantity', value: number): void
    (e: 'update:deliveryType', value: string | null): void
    (e: 'update:transportDate', value: string): void
    (e: 'update:officeCommissionAmount', value: number): void
    (e: 'update:customerTotalPaid', value: number): void
    (e: 'update:notes', value: string): void
}>()

// Helper
const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}
</script>