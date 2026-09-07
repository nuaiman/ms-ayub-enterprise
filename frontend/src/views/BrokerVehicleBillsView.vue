<!-- src/views/BrokerVehicleBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ brokerVehicleBillsStore.totalBills
                        }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(brokerVehicleBillsStore.totalAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ brokerVehicleBillsStore.totalUnpaid }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(brokerVehicleBillsStore.totalUnpaidAmount) }}</p>
                </div>
            </div>

            <!-- Broker Vehicle Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <BrokerVehicleBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useBrokerVehicleBillsStore } from '@/stores/brokerVehicleBills'
import { useVehiclesStore } from '@/stores/vehicles'
import { useBrokersStore } from '@/stores/brokers'
import { useTransportsStore } from '@/stores/transports'
import AppLayout from '@/components/layouts/AppLayout.vue'
import BrokerVehicleBillList from '@/components/features/brokerVehicleBills/BrokerVehicleBillList.vue'
import { formatCurrency } from '@/utils/currency'

const brokerVehicleBillsStore = useBrokerVehicleBillsStore()
const vehiclesStore = useVehiclesStore()
const brokersStore = useBrokersStore()
const transportsStore = useTransportsStore()

onMounted(async () => {
    await Promise.all([
        vehiclesStore.fetchVehicles(),
        brokersStore.fetchBrokers(),
        transportsStore.fetchTransports(),
    ])
})
</script>