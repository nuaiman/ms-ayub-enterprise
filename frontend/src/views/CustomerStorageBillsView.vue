<!-- src/views/CustomerStorageBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Active Lots</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{
                        customerStorageBillsStore.customerBillData.length }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Monthly Bill</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(customerStorageBillsStore.totalMonthlyBill) }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Billed</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{
                        formatCurrency(customerStorageBillsStore.totalBilled) }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(customerStorageBillsStore.totalOutstanding) }}
                    </p>
                </div>
            </div>

            <!-- Customer Storage Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <CustomerStorageBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCustomerStorageBillsStore } from '@/stores/customerStorageBills'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import AppLayout from '@/components/layouts/AppLayout.vue'
import CustomerStorageBillList from '@/components/features/customerStorageBills/CustomerStorageBillList.vue'
import { formatCurrency } from '@/utils/currency'

const customerStorageBillsStore = useCustomerStorageBillsStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const itemsStore = useItemsStore()
const customersStore = useCustomersStore()

onMounted(async () => {
    await Promise.all([
        lotsStore.fetchLots(),
        storesStore.fetchStores(),
        itemsStore.fetchItems(),
        customersStore.fetchCustomers(),
    ])
})
</script>