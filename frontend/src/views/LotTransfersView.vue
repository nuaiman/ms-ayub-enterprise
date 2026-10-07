<!-- src/views/LotTransfersView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Transfers</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ lotTransfersStore.totalTransfers }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ thisMonthCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Lots Moved</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ distinctLotsCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customers Involved</p>
                <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ distinctCustomersCount }}</p>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 min-h-0 mt-6">
            <LotTransferList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useLotTransfersStore } from '@/stores/lotTransfers'
import LotTransferList from '@/components/features/lotTransfers/LotTransferList.vue'

const lotTransfersStore = useLotTransfersStore()

const thisMonthCount = computed(() => {
    const now = new Date()
    const y = now.getFullYear()
    const m = now.getMonth()
    return lotTransfersStore.transfers.filter(t => {
        const d = new Date(t.transferred_at)
        return d.getFullYear() === y && d.getMonth() === m
    }).length
})

const distinctLotsCount = computed(() => {
    const set = new Set<number>()
    lotTransfersStore.transfers.forEach(t => set.add(t.lot_id))
    return set.size
})

const distinctCustomersCount = computed(() => {
    const set = new Set<number>()
    lotTransfersStore.transfers.forEach(t => {
        set.add(t.from_customer_id)
        set.add(t.to_customer_id)
    })
    return set.size
})

onMounted(() => {
    lotTransfersStore.fetchAllTransfers()
})
</script>