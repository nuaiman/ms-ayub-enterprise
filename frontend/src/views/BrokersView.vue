<!-- src/views/BrokersView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Brokers</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ brokersStore.brokers.length }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Phone</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withPhoneCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Notes</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ withNotesCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Created This Month</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ createdThisMonthCount }}</p>
                </div>
            </div>

            <!-- Broker List -->
            <div class="flex-1 min-h-0 mt-6">
                <BrokerList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useBrokersStore } from '@/stores/brokers'
import AppLayout from '@/components/layouts/AppLayout.vue'
import BrokerList from '@/components/features/brokers/BrokerList.vue'

const brokersStore = useBrokersStore()

const withPhoneCount = computed(() => {
    return brokersStore.brokers.filter(b => b.phone).length
})

const withNotesCount = computed(() => {
    return brokersStore.brokers.filter(b => b.notes).length
})

const createdThisMonthCount = computed(() => {
    const now = new Date()
    const currentMonth = now.getMonth()
    const currentYear = now.getFullYear()
    return brokersStore.brokers.filter(b => {
        const date = new Date(b.created_at)
        return date.getMonth() === currentMonth && date.getFullYear() === currentYear
    }).length
})

onMounted(() => {
    brokersStore.fetchBrokers()
})
</script>