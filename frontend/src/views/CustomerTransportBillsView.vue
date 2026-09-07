<!-- src/views/CustomerTransportBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">
                        Total Bills
                    </p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                        {{ store.totalBills }}
                    </p>
                </div>

                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">
                        Total Amount
                    </p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">
                        {{ formatCurrency(store.totalAmount) }}
                    </p>
                </div>

                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">
                        Unpaid
                    </p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">
                        {{ store.totalUnpaid }}
                    </p>
                </div>

                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">
                        Unpaid Amount
                    </p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">
                        {{ formatCurrency(store.totalUnpaidAmount) }}
                    </p>
                </div>
            </div>

            <!-- Customer Transport Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <CustomerTransportBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCustomerTransportBillsStore } from '@/stores/customerTransportBills'
import { useTransportsStore } from '@/stores/transports'
import { useVehiclesStore } from '@/stores/vehicles'
import { useCustomersStore } from '@/stores/customers'
import { useBrokersStore } from '@/stores/brokers'
import AppLayout from '@/components/layouts/AppLayout.vue'
import CustomerTransportBillList from '@/components/features/customerTransportBills/CustomerTransportBillList.vue'
import { formatCurrency } from '@/utils/currency'

const store = useCustomerTransportBillsStore()
const transportsStore = useTransportsStore()
const vehiclesStore = useVehiclesStore()
const customersStore = useCustomersStore()
const brokersStore = useBrokersStore()

onMounted(async () => {
    await Promise.all([
        transportsStore.fetchTransports(),
        vehiclesStore.fetchVehicles(),
        customersStore.fetchCustomers(),
        brokersStore.fetchBrokers(),
    ])
})
</script>