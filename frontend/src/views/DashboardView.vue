<!-- src/views/DashboardView.vue -->
<template>
    <AppLayout>
        <!-- Loading State -->
        <DashboardLoading v-if="dashboardStore.isLoading" />

        <!-- Dashboard Content -->
        <div v-else-if="dashboardStore.isLoaded" class="flex flex-col h-full">
            <!-- Stats -->
            <DashboardStats :stats="dashboardStore.metrics" />

            <!-- Outstanding Snapshot -->
            <div class="mt-4">
                <OutstandingSnapshot :total-due="dashboardStore.metrics.totalDue"
                    :outstanding-by-customer="dashboardStore.metrics.outstandingByCustomer"
                    :outstanding-by-type="dashboardStore.metrics.outstandingByType"
                    :majhi-unpaid-amount="dashboardStore.metrics.majhiUnpaidAmount"
                    :broker-unpaid-amount="dashboardStore.metrics.brokerUnpaidAmount"
                    :rent-unpaid-amount="dashboardStore.metrics.rentUnpaidAmount" />
            </div>

            <!-- Customer Revenue Table -->
            <div class="mt-4 flex-1 min-h-0">
                <CustomerRevenueTable v-model="dashboardStore.selectedMonth" :customers="dashboardStore.customerRevenue"
                    :months="dashboardStore.availableMonths" />
            </div>
        </div>

        <!-- Error / Empty State -->
        <div v-else class="flex flex-col items-center justify-center h-full">
            <p class="text-(--color-text-secondary)">No data available</p>
            <button @click="loadData"
                class="mt-4 px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-opacity">
                Retry
            </button>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import AppLayout from '@/components/layouts/AppLayout.vue'
import DashboardLoading from '@/components/features/dashboard/DashboardLoading.vue'
import DashboardStats from '@/components/features/dashboard/DashboardStats.vue'
import OutstandingSnapshot from '@/components/features/dashboard/OutstandingSnapshot.vue'
import CustomerRevenueTable from '@/components/features/dashboard/CustomerRevenueTable.vue'
import { useDashboardStore } from '@/stores/dashboards'

const dashboardStore = useDashboardStore()

const loadData = async () => {
    await dashboardStore.loadDashboardData()
}

onMounted(() => {
    loadData()
})
</script>