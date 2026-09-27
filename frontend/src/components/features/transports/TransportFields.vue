<!-- src/components/features/transports/TransportFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Customer -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Customer <span v-if="required" class="text-(--color-red)">*</span>
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

        <!-- Transport Type -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Transport Type
            </label>
            <select :value="transportType" @change="onTransportTypeChange" :disabled="disabled"
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
import type { TransportType } from '@/types/transport'

defineProps<{
    customerId: number | null
    fromLocation: string
    toLocation: string
    vehicleQuantity: number
    transportType: TransportType | null
    transportDate: string
    officeCommissionAmount: number
    notes: string
    customerOptions?: Customer[]
    disabled?: boolean
    required?: boolean
}>()

const emit = defineEmits<{
    (e: 'update:customerId', value: number | null): void
    (e: 'update:fromLocation', value: string): void
    (e: 'update:toLocation', value: string): void
    (e: 'update:vehicleQuantity', value: number): void
    (e: 'update:transportType', value: TransportType | null): void
    (e: 'update:transportDate', value: string): void
    (e: 'update:officeCommissionAmount', value: number): void
    (e: 'update:notes', value: string): void
}>()

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}

const onTransportTypeChange = (event: Event) => {
    const value = (event.target as HTMLSelectElement).value
    if (value === 'local' || value === 'district') {
        emit('update:transportType', value)
    } else {
        emit('update:transportType', null)
    }
}
</script>