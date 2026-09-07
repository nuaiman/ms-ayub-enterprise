<!-- src/views/DeliveriesView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Deliveries</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ deliveriesStore.deliveries.length
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
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Items</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ withItemsCount }}</p>
                </div>
            </div>

            <!-- Delivery List -->
            <div class="flex-1 min-h-0 mt-6">
                <DeliveryList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import AppLayout from '@/components/layouts/AppLayout.vue'
import DeliveryList from '@/components/features/deliveries/DeliveryList.vue'

const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const usersStore = useUsersStore()
const deliveryItemsStore = useDeliveryItemsStore()

const withCustomerCount = computed(() => {
    return deliveriesStore.deliveries.filter(d => d.customer_id).length
})

const thisMonthCount = computed(() => {
    const now = new Date()
    const currentMonth = now.getMonth()
    const currentYear = now.getFullYear()
    return deliveriesStore.deliveries.filter(d => {
        const date = new Date(d.delivery_date)
        return date.getMonth() === currentMonth && date.getFullYear() === currentYear
    }).length
})

const withItemsCount = computed(() => {
    return deliveriesStore.deliveries.filter(d => deliveryItemsStore.getDeliveryItemsByDeliveryId(d.id).length > 0).length
})

onMounted(async () => {
    await Promise.all([
        deliveriesStore.fetchDeliveries(),
        customersStore.fetchCustomers(),
        usersStore.fetchUsers(),
        deliveryItemsStore.fetchDeliveryItems()
    ])
})
</script>