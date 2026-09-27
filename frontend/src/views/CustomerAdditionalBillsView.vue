<!-- src/views/CustomerAdditionalBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ store.totalBills }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ formatCurrency(store.totalAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ store.totalUnpaid }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{ formatCurrency(store.totalUnpaidAmount) }}
                    </p>
                </div>
            </div>

            <!-- List -->
            <div class="flex-1 min-h-0 mt-6">
                <CustomerAdditionalBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { useCustomersStore } from '@/stores/customers'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useTransportsStore } from '@/stores/transports'
import { useDamagesStore } from '@/stores/damages'
import { useGodownsStore } from '@/stores/godowns'
import AppLayout from '@/components/layouts/AppLayout.vue'
import CustomerAdditionalBillList from '@/components/features/customerAdditionalBills/CustomerAdditionalBillList.vue'
import { formatCurrency } from '@/utils/currency'

const store = useCustomerAdditionalBillsStore()
const customersStore = useCustomersStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const deliveriesStore = useDeliveriesStore()
const transportsStore = useTransportsStore()
const damagesStore = useDamagesStore()
const godownsStore = useGodownsStore()

onMounted(async () => {
    await Promise.all([
        customersStore.fetchCustomers(),
        lotsStore.fetchLots(),
        storesStore.fetchStores(),
        deliveriesStore.fetchDeliveries(),
        transportsStore.fetchTransports(),
        damagesStore.fetchDamages(),
        godownsStore.fetchGodowns(),
        store.fetchCharges(),
    ])
})
</script>