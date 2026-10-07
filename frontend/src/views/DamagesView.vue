<!-- src/views/DamagesView.vue -->
<template>
    <div class="space-y-4">
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Damages</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ damagesStore.totalDamages }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                <p class="text-2xl font-bold text-(--color-red) mt-1">{{ formatCurrency(damagesStore.totalDamageAmount)
                    }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ thisMonthCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Image</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ withImageCount }}</p>
            </div>
        </div>

        <!-- List -->
        <DamageList />
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useDamagesStore } from '@/stores/damages'
import DamageList from '@/components/features/damages/DamageList.vue'
import { formatCurrency } from '@/utils/currency'

const damagesStore = useDamagesStore()

const thisMonthCount = computed(() => {
    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth()
    return damagesStore.damages.filter(d => {
        const date = new Date(d.damage_date)
        return date.getFullYear() === year && date.getMonth() === month
    }).length
})

const withImageCount = computed(() =>
    damagesStore.damages.filter(d => d.image_url !== null && d.image_url !== '').length
)
</script>