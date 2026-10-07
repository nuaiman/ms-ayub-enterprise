<!-- src/views/StoreTransfersView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Transfers</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ storeTransfersStore.totalTransfers }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ thisMonthCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Stores Moved</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ distinctStoresCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Godowns Involved</p>
                <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ distinctGodownsCount }}</p>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 min-h-0 mt-6">
            <StoreTransferList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useStoreTransfersStore } from '@/stores/storeTransfers'
import StoreTransferList from '@/components/features/storeTransfers/StoreTransferList.vue'

const storeTransfersStore = useStoreTransfersStore()

const thisMonthCount = computed(() => {
    const now = new Date()
    const y = now.getFullYear()
    const m = now.getMonth()
    return storeTransfersStore.transfers.filter(t => {
        const d = new Date(t.transferred_at)
        return d.getFullYear() === y && d.getMonth() === m
    }).length
})

const distinctStoresCount = computed(() => {
    const set = new Set<number>()
    storeTransfersStore.transfers.forEach(t => set.add(t.store_id))
    return set.size
})

const distinctGodownsCount = computed(() => {
    const set = new Set<number>()
    storeTransfersStore.transfers.forEach(t => {
        set.add(t.from_godown_id)
        set.add(t.to_godown_id)
    })
    return set.size
})

onMounted(() => {
    storeTransfersStore.fetchAllTransfers()
})
</script>