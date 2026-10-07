<!-- src/views/IncomesView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Incomes</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ incomesStore.totalIncomes }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ formatCurrency(incomesStore.totalAmount)
                    }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ thisMonthCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month Amount</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ formatCurrency(thisMonthAmount) }}</p>
            </div>
        </div>

        <!-- Income List -->
        <div class="flex-1 min-h-0 mt-6">
            <IncomeList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useIncomesStore } from '@/stores/incomes'
import IncomeList from '@/components/features/incomes/IncomeList.vue'
import { formatCurrency } from '@/utils/currency'

const incomesStore = useIncomesStore()

const thisMonth = computed(() => {
    const now = new Date()
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
})

const thisMonthCount = computed(() => {
    return incomesStore.incomes.filter(e => {
        const date = new Date(e.income_date)
        const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
        return month === thisMonth.value
    }).length
})

const thisMonthAmount = computed(() => {
    return incomesStore.incomes
        .filter(e => {
            const date = new Date(e.income_date)
            const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
            return month === thisMonth.value
        })
        .reduce((sum, e) => sum + e.amount, 0)
})

onMounted(() => {
    incomesStore.fetchIncomes()
})
</script>