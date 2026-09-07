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
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Broker</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withBrokerCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Driver</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ withDriverCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Joma Cost</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{
                        formatCurrency(vehiclesStore.totalJomaCost) }}</p>
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
import { computed, onMounted } from 'vue'
import { useVehiclesStore } from '@/stores/vehicles'
import { useBrokersStore } from '@/stores/brokers'
import { useTransportsStore } from '@/stores/transports'
import { useUsersStore } from '@/stores/users'
import AppLayout from '@/components/layouts/AppLayout.vue'
import VehicleList from '@/components/features/vehicles/VehicleList.vue'
import { formatCurrency } from '@/utils/currency'

const vehiclesStore = useVehiclesStore()
const brokersStore = useBrokersStore()
const transportsStore = useTransportsStore()
const usersStore = useUsersStore()

const withBrokerCount = computed(() => {
    return vehiclesStore.vehicles.filter(v => v.broker_id).length
})

const withDriverCount = computed(() => {
    return vehiclesStore.vehicles.filter(v => v.driver_name).length
})

onMounted(async () => {
    await Promise.all([
        vehiclesStore.fetchVehicles(),
        brokersStore.fetchBrokers(),
        transportsStore.fetchTransports(),
        usersStore.fetchUsers()
    ])
})
</script>