<!-- src/views/LotsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Lots</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ lotsStore.lots.length }}</p>
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
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Majhi</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withMajhiCount }}</p>
                </div>
            </div>

            <!-- Lot List -->
            <div class="flex-1 min-h-0 mt-6">
                <LotList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import AppLayout from '@/components/layouts/AppLayout.vue'
import LotList from '@/components/features/lots/LotList.vue'

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()

const activeCount = computed(() => lotsStore.lots.filter(l => l.is_active).length)
const inactiveCount = computed(() => lotsStore.lots.filter(l => !l.is_active).length)
const withMajhiCount = computed(() => lotsStore.lots.filter(l => l.majhi_id).length)

onMounted(async () => {
    await Promise.all([
        lotsStore.fetchLots(),
        customersStore.fetchCustomers(),
        majhisStore.fetchMajhis(),
    ])
})
</script>