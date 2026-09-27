<!-- src/components/features/vehicles/VehicleFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Transport (only when standalone=true) -->
        <div v-if="standalone">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Transport <span v-if="required" class="text-(--color-red)">*</span>
            </label>
            <select :value="transportId"
                @change="$emit('update:transportId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled || !canEditTransport"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a transport</option>
                <option v-for="transport in (transportOptions || [])" :key="transport.id" :value="transport.id">
                    {{ getTransportDisplayName(transport) }}
                </option>
            </select>
            <p v-if="!canEditTransport" class="text-xs text-(--color-text-secondary) mt-1">
                Transport cannot be changed
            </p>
            <p v-if="transportOptions && transportOptions.length === 0" class="text-xs text-(--color-yellow) mt-1">
                No transports available. Please create a transport first.
            </p>
        </div>

        <!-- Vehicle Number -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Vehicle Number <span v-if="required" class="text-(--color-red)">*</span>
            </label>
            <input :value="vehicleNumber"
                @input="$emit('update:vehicleNumber', ($event.target as HTMLInputElement).value)" type="text"
                placeholder="Enter vehicle number" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <!-- Broker -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Broker <span v-if="required" class="text-(--color-red)">*</span>
            </label>
            <select :value="brokerId"
                @change="$emit('update:brokerId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a broker</option>
                <option v-for="broker in (brokerOptions || [])" :key="broker.id" :value="broker.id">
                    {{ broker.name }}
                </option>
            </select>
        </div>

        <!-- Costs - Row 1 -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Joma Cost
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="jomaCost"
                        @input="$emit('update:jomaCost', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Vehicle Cost
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="vehicleCost"
                        @input="$emit('update:vehicleCost', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>

        <!-- Expense Costs -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Other Cost <span class="text-xs font-normal text-(--color-text-secondary)">(expense)</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="otherCost"
                        @input="$emit('update:otherCost', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <p class="text-[10px] text-(--color-text-secondary) mt-0.5">Creates an expense record</p>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Labour Cost <span class="text-xs font-normal text-(--color-text-secondary)">(expense)</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="labourCost"
                        @input="$emit('update:labourCost', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <p class="text-[10px] text-(--color-text-secondary) mt-0.5">Creates an expense record</p>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Demarage Cost <span class="text-xs font-normal text-(--color-text-secondary)">(expense)</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="demarageCost"
                        @input="$emit('update:demarageCost', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
                <p class="text-[10px] text-(--color-text-secondary) mt-0.5">Creates an expense record</p>
            </div>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea :value="notes" @input="$emit('update:notes', ($event.target as HTMLTextAreaElement).value)"
                rows="2" placeholder="Enter any notes about this vehicle" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Transport } from '@/types/transport'
import type { Broker } from '@/types/broker'

defineProps<{
    transportId?: number | null
    vehicleNumber: string
    brokerId?: number | null
    jomaCost: number
    vehicleCost: number
    otherCost: number
    labourCost: number
    demarageCost: number
    notes?: string
    transportOptions?: Transport[]
    brokerOptions?: Broker[]
    disabled?: boolean
    required?: boolean
    standalone?: boolean
    canEditTransport?: boolean
}>()

defineEmits<{
    (e: 'update:transportId', value: number | null): void
    (e: 'update:vehicleNumber', value: string): void
    (e: 'update:brokerId', value: number | null): void
    (e: 'update:jomaCost', value: number): void
    (e: 'update:vehicleCost', value: number): void
    (e: 'update:otherCost', value: number): void
    (e: 'update:labourCost', value: number): void
    (e: 'update:demarageCost', value: number): void
    (e: 'update:notes', value: string): void
}>()

const getTransportDisplayName = (transport: Transport): string => {
    const fromTo = `${transport.from_location}${transport.to_location ? ` ৳ ${transport.to_location}` : ''}`
    return `${fromTo} (${transport.vehicle_quantity} vehicles)`
}
</script>