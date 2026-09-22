<!-- src/views/CustomerLotBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ customerLotBillsStore.totalBills
                    }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(customerLotBillsStore.totalAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ customerLotBillsStore.totalUnpaid }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(customerLotBillsStore.totalUnpaidAmount) }}</p>
                </div>
            </div>

            <!-- Customer Lot Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <CustomerLotBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCustomerLotBillsStore } from '@/stores/customerLotBills'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import AppLayout from '@/components/layouts/AppLayout.vue'
import CustomerLotBillList from '@/components/features/customerLotBills/CustomerLotBillList.vue'
import { formatCurrency } from '@/utils/currency'

const customerLotBillsStore = useCustomerLotBillsStore()
const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()

onMounted(async () => {
    await Promise.all([
        lotsStore.fetchLots(),
        customersStore.fetchCustomers(),
        storesStore.fetchStores(),
    ])
})
</script>