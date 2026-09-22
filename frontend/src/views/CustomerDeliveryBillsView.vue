<!-- src/views/CustomerDeliveryBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{
                        customerDeliveryBillsStore.totalBills }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(customerDeliveryBillsStore.totalAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ customerDeliveryBillsStore.totalUnpaid
                    }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(customerDeliveryBillsStore.totalUnpaidAmount) }}</p>
                </div>
            </div>

            <!-- Customer Delivery Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <CustomerDeliveryBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useLotsStore } from '@/stores/lots'
import AppLayout from '@/components/layouts/AppLayout.vue'
import CustomerDeliveryBillList from '@/components/features/customerDeliveryBills/CustomerDeliveryBillList.vue'
import { formatCurrency } from '@/utils/currency'

const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
const deliveryItemsStore = useDeliveryItemsStore()
const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const lotsStore = useLotsStore()

onMounted(async () => {
    await Promise.all([
        deliveryItemsStore.fetchDeliveryItems(),
        deliveriesStore.fetchDeliveries(),
        customersStore.fetchCustomers(),
        lotsStore.fetchLots(),
    ])
})
</script>