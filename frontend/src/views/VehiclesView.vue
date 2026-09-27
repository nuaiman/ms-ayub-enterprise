<!-- src/views/VehiclesView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Vehicles</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ vehiclesStore.totalVehicles }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Joma Cost</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(vehiclesStore.totalJomaCost) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Vehicle Cost</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{
                        formatCurrency(vehiclesStore.totalVehicleCost) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Broker Paid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{
                        formatCurrency(vehiclesStore.totalBrokerPaid) }}</p>
                </div>
            </div>

            <!-- Vehicle List -->
            <div class="flex-1 min-h-0 mt-6">
                <VehicleList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useVehiclesStore } from '@/stores/vehicles'
import { useTransportsStore } from '@/stores/transports'
import { useBrokersStore } from '@/stores/brokers'
import AppLayout from '@/components/layouts/AppLayout.vue'
import VehicleList from '@/components/features/vehicles/VehicleList.vue'
import { formatCurrency } from '@/utils/currency'

const vehiclesStore = useVehiclesStore()
const transportsStore = useTransportsStore()
const brokersStore = useBrokersStore()

onMounted(async () => {
    await Promise.all([
        vehiclesStore.fetchVehicles(),
        transportsStore.fetchTransports(),
        brokersStore.fetchBrokers(),
    ])
})
</script>