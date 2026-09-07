<!-- src/views/RentsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Rents</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ rentsStore.totalRents }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Draft</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ rentsStore.draftCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Paid</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ rentsStore.statusCounts.paid || 0 }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{ formatCurrency(rentsStore.totalOutstanding)
                        }}</p>
                </div>
            </div>

            <!-- Rent List -->
            <div class="flex-1 min-h-0 mt-6">
                <RentList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRentsStore } from '@/stores/rents'
import { useGodownsStore } from '@/stores/godowns'
import { useUsersStore } from '@/stores/users'
import AppLayout from '@/components/layouts/AppLayout.vue'
import RentList from '@/components/features/rents/RentList.vue'
import { formatCurrency } from '@/utils/currency'

const rentsStore = useRentsStore()
const godownsStore = useGodownsStore()
const usersStore = useUsersStore()

onMounted(async () => {
    await Promise.all([
        rentsStore.fetchCurrentMonthRents(),
        godownsStore.fetchGodowns(),
        usersStore.fetchUsers()
    ])
})
</script>