<!-- src/components/features/vehicles/VehicleDetail.vue -->
<template>
    <div v-if="vehicle" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ vehicle.vehicle_number }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">Transport #{{ vehicle.transport_id }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ brokerName }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ vehicle.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(vehicle.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(vehicle.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Transport</p>
                <p class="text-sm text-(--color-text-primary)">#{{ vehicle.transport_id }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicle Number</p>
                <p class="text-sm text-(--color-text-primary)">{{ vehicle.vehicle_number }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Broker</p>
                <p class="text-sm text-(--color-text-primary)">{{ brokerName }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created By</p>
                <p class="text-sm text-(--color-text-primary)">{{ getUserName(vehicle.user_id) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Joma Cost</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(vehicle.joma_cost) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicle Cost</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(vehicle.vehicle_cost) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Paid to
                    Broker</p>
                <p class="text-sm font-semibold text-(--color-green)">{{ formatCurrency(vehicle.total_paid_to_broker) }}
                </p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Other Cost</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(vehicle.other_cost) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Labour Cost</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(vehicle.labour_cost) }}</p>
            </div>

            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Demarage Cost</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(vehicle.demarage_cost) }}</p>
            </div>

            <div v-if="vehicle.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) mt-1">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ vehicle.notes }}</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', vehicle)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Vehicle } from '@/types/vehicle'
import { useBrokersStore } from '@/stores/brokers'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    vehicle: Vehicle | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [vehicle: Vehicle]
    'updated': []
}>()

const brokersStore = useBrokersStore()
const usersStore = useUsersStore()

const brokerName = computed(() => {
    if (!props.vehicle?.broker_id) return '৳'
    return brokersStore.getBrokerName(props.vehicle.broker_id)
})

const getUserName = (id: number): string => usersStore.getUserName(id)

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}
</script>