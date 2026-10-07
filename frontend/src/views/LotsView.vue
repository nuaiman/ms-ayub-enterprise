<!-- src/views/LotsView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Lots</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ lotsStore.lots.length }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Stores</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ totalStores }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Active Stores</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ activeStores }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Customers</p>
                <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ distinctCustomers }}</p>
            </div>
        </div>

        <div class="flex-1 min-h-0 mt-6">
            <LotList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import LotList from '@/components/features/lots/LotList.vue'

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()

const totalStores = computed(() => storesStore.stores.length)
const activeStores = computed(() => storesStore.stores.filter(s => s.is_active).length)

const distinctCustomers = computed(() => {
    const set = new Set<number>()
    lotsStore.lots.forEach(l => set.add(l.customer_id))
    return set.size
})

onMounted(async () => {
    await Promise.all([
        lotsStore.fetchLots(),
        customersStore.fetchCustomers(),
        storesStore.fetchStores(),
    ])
})
</script>