<!-- src/views/ItemsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Items</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ itemsStore.items.length }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Active</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ activeCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Inactive</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{ inactiveCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Customer</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withCustomerCount }}</p>
                </div>
            </div>

            <!-- Item List -->
            <div class="flex-1 min-h-0 mt-6">
                <ItemList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import AppLayout from '@/components/layouts/AppLayout.vue'
import ItemList from '@/components/features/items/ItemList.vue'

const itemsStore = useItemsStore()
const customersStore = useCustomersStore()

const activeCount = computed(() => {
    return itemsStore.items.filter(i => i.is_active).length
})

const inactiveCount = computed(() => {
    return itemsStore.items.filter(i => !i.is_active).length
})

const withCustomerCount = computed(() => {
    return itemsStore.items.filter(i => i.customer_id).length
})

onMounted(async () => {
    await Promise.all([
        itemsStore.fetchItems(),
        customersStore.fetchCustomers()
    ])
})
</script>