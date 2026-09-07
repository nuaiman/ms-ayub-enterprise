<!-- src/views/GodownsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Godowns</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ godownsStore.godowns.length }}</p>
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
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Monthly Rent</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(godownsStore.getTotalMonthlyRent()) }}</p>
                </div>
            </div>

            <!-- Godown List -->
            <div class="flex-1 min-h-0 mt-6">
                <GodownList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useGodownsStore } from '@/stores/godowns'
import AppLayout from '@/components/layouts/AppLayout.vue'
import GodownList from '@/components/features/godowns/GodownList.vue'
import { formatCurrency } from '@/utils/currency'

const godownsStore = useGodownsStore()

const activeCount = computed(() => {
    return godownsStore.godowns.filter(g => g.is_active).length
})

const inactiveCount = computed(() => {
    return godownsStore.godowns.filter(g => !g.is_active).length
})

onMounted(() => {
    godownsStore.fetchGodowns()
})
</script>