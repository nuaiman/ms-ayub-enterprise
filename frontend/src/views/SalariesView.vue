<!-- src/views/SalariesView.vue -->
<template>
    <div class="space-y-6">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Salaries</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ salariesStore.totalSalaries }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Draft</p>
                <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ salariesStore.draftCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Paid</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ salariesStore.statusCounts.paid || 0 }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Cancelled</p>
                <p class="text-2xl font-bold text-(--color-red) mt-1">{{ salariesStore.statusCounts.cancelled || 0
                }}</p>
            </div>
        </div>

        <!-- Salary List -->
        <div class="flex-1 min-h-0 mt-6">
            <SalaryList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSalariesStore } from '@/stores/salaries'
import SalaryList from '@/components/features/salaries/SalaryList.vue'

const salariesStore = useSalariesStore()

onMounted(async () => {
    await salariesStore.fetchSalaries()
})
</script>