<!-- src/views/DamagesView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Damages</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ damagesStore.damages.length }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(damagesStore.totalDamageAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ thisMonthCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month Amount</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ formatCurrency(thisMonthAmount) }}</p>
                </div>
            </div>

            <!-- Damage List -->
            <div class="flex-1 min-h-0 mt-6">
                <DamageList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDamagesStore } from '@/stores/damages'
import { useStoresStore } from '@/stores/stores'
import { useUsersStore } from '@/stores/users'
import AppLayout from '@/components/layouts/AppLayout.vue'
import DamageList from '@/components/features/damages/DamageList.vue'
import { formatCurrency } from '@/utils/currency'

const damagesStore = useDamagesStore()
const storesStore = useStoresStore()
const usersStore = useUsersStore()

const thisMonth = computed(() => {
    const now = new Date()
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
})

const thisMonthCount = computed(() => {
    return damagesStore.damages.filter(d => {
        const date = new Date(d.damage_date)
        const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
        return month === thisMonth.value
    }).length
})

const thisMonthAmount = computed(() => {
    return damagesStore.damages
        .filter(d => {
            const date = new Date(d.damage_date)
            const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
            return month === thisMonth.value
        })
        .reduce((sum, d) => sum + d.amount, 0)
})

onMounted(async () => {
    await Promise.all([
        damagesStore.fetchDamages(),
        storesStore.fetchStores(),
        usersStore.fetchUsers()
    ])
})
</script>