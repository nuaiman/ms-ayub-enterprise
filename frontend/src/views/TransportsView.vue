<!-- src/views/TransportsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Transports</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ transportsStore.transports.length
                        }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Customer</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withCustomerCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ thisMonthCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Commission</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{
                        formatCurrency(transportsStore.totalCommission) }}</p>
                </div>
            </div>

            <!-- Transport List -->
            <div class="flex-1 min-h-0 mt-6">
                <TransportList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useTransportsStore } from '@/stores/transports'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import AppLayout from '@/components/layouts/AppLayout.vue'
import TransportList from '@/components/features/transports/TransportList.vue'
import { formatCurrency } from '@/utils/currency'

const transportsStore = useTransportsStore()
const customersStore = useCustomersStore()
const usersStore = useUsersStore()

const withCustomerCount = computed(() => {
    return transportsStore.transports.filter(t => t.customer_id).length
})

const thisMonthCount = computed(() => {
    const now = new Date()
    const currentMonth = now.getMonth()
    const currentYear = now.getFullYear()
    return transportsStore.transports.filter(t => {
        const date = new Date(t.transport_date)
        return date.getMonth() === currentMonth && date.getFullYear() === currentYear
    }).length
})

onMounted(async () => {
    await Promise.all([
        transportsStore.fetchTransports(),
        customersStore.fetchCustomers(),
        usersStore.fetchUsers()
    ])
})
</script>