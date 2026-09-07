<!-- src/views/ExpensesView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Expenses</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ expensesStore.totalExpenses }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{ formatCurrency(expensesStore.totalAmount)
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

            <!-- Expense List -->
            <div class="flex-1 min-h-0 mt-6">
                <ExpenseList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useExpensesStore } from '@/stores/expenses'
import AppLayout from '@/components/layouts/AppLayout.vue'
import ExpenseList from '@/components/features/expenses/ExpenseList.vue'
import { formatCurrency } from '@/utils/currency'

const expensesStore = useExpensesStore()

const thisMonth = computed(() => {
    const now = new Date()
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
})

const thisMonthCount = computed(() => {
    return expensesStore.expenses.filter(e => {
        const date = new Date(e.expense_date)
        const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
        return month === thisMonth.value
    }).length
})

const thisMonthAmount = computed(() => {
    return expensesStore.expenses
        .filter(e => {
            const date = new Date(e.expense_date)
            const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
            return month === thisMonth.value
        })
        .reduce((sum, e) => sum + e.amount, 0)
})

onMounted(() => {
    expensesStore.fetchExpenses()
})
</script>