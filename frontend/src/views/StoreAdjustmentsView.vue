<!-- src/views/StoreAdjustmentsView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Adjustments</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ storeAdjustmentsStore.totalAdjustments }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ thisMonthCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Increases</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ increasesCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Decreases</p>
                <p class="text-2xl font-bold text-(--color-red) mt-1">{{ decreasesCount }}</p>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 min-h-0 mt-6">
            <StoreAdjustmentList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useStoreAdjustmentsStore } from '@/stores/storeAdjustments'
import StoreAdjustmentList from '@/components/features/storeAdjustments/StoreAdjustmentList.vue'

const storeAdjustmentsStore = useStoreAdjustmentsStore()

const thisMonthCount = computed(() => {
    const now = new Date()
    const y = now.getFullYear()
    const m = now.getMonth()
    return storeAdjustmentsStore.adjustments.filter(a => {
        const d = new Date(a.adjusted_at)
        return d.getFullYear() === y && d.getMonth() === m
    }).length
})

const increasesCount = computed(() =>
    storeAdjustmentsStore.adjustments.filter(
        a => a.weight_delta > 0 || a.quantity_delta > 0
    ).length
)

const decreasesCount = computed(() =>
    storeAdjustmentsStore.adjustments.filter(
        a => a.weight_delta < 0 || a.quantity_delta < 0
    ).length
)

onMounted(() => {
    storeAdjustmentsStore.fetchAllAdjustments()
})
</script>