<!-- src/views/DeliveriesView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Deliveries</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ deliveriesStore.totalDeliveries }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ thisMonthCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Items</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ deliveriesStore.totalItemCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customers Served</p>
                <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ distinctCustomers }}</p>
            </div>
        </div>

        <div class="flex-1 min-h-0 mt-6">
            <DeliveryList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDeliveriesStore } from '@/stores/deliveries'
import DeliveryList from '@/components/features/deliveries/DeliveryList.vue'

const deliveriesStore = useDeliveriesStore()

const thisMonthCount = computed(() => {
    const now = new Date()
    const y = now.getFullYear()
    const m = now.getMonth()
    return deliveriesStore.deliveries.filter(d => {
        const date = new Date(d.delivery_date)
        return date.getFullYear() === y && date.getMonth() === m
    }).length
})

const distinctCustomers = computed(() => {
    const set = new Set<number>()
    deliveriesStore.deliveries.forEach(d => {
        if (d.customer_id) set.add(d.customer_id)
    })
    return set.size
})

onMounted(() => {
    deliveriesStore.fetchDeliveries()
})
</script>