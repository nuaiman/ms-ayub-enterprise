<!-- src/views/StoresView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Stores</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ storesStore.stores.length }}</p>
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
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Inventory</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withInventoryCount }}</p>
                </div>
            </div>

            <!-- Store List -->
            <div class="flex-1 min-h-0 mt-6">
                <StoreList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useGodownsStore } from '@/stores/godowns'
import AppLayout from '@/components/layouts/AppLayout.vue'
import StoreList from '@/components/features/stores/StoreList.vue'

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()

const activeCount = computed(() => {
    return storesStore.stores.filter(s => s.is_active).length
})

const inactiveCount = computed(() => {
    return storesStore.stores.filter(s => !s.is_active).length
})

const withInventoryCount = computed(() => {
    return storesStore.stores.filter(s => s.quantity > 0 || s.weight > 0).length
})

onMounted(async () => {
    await Promise.all([
        storesStore.fetchStores(),
        lotsStore.fetchLots(),
        godownsStore.fetchGodowns()
    ])
})
</script>