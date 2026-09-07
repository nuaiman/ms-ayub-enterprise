<!-- src/views/LogsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Logs</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ logsStore.totalLogs }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Today</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ todayCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Week</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ thisWeekCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">This Month</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ thisMonthCount }}</p>
                </div>
            </div>

            <!-- Log List -->
            <div class="flex-1 min-h-0 mt-6">
                <LogList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useLogsStore } from '@/stores/logs'
import AppLayout from '@/components/layouts/AppLayout.vue'
import LogList from '@/components/features/logs/LogList.vue'

const logsStore = useLogsStore()

const todayCount = computed(() => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    return logsStore.logs.filter(log => {
        const date = new Date(log.created_at)
        return date >= today
    }).length
})

const thisWeekCount = computed(() => {
    const now = new Date()
    const startOfWeek = new Date(now)
    startOfWeek.setDate(now.getDate() - now.getDay())
    startOfWeek.setHours(0, 0, 0, 0)
    return logsStore.logs.filter(log => {
        const date = new Date(log.created_at)
        return date >= startOfWeek
    }).length
})

const thisMonthCount = computed(() => {
    const now = new Date()
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    return logsStore.logs.filter(log => {
        const date = new Date(log.created_at)
        return date >= startOfMonth
    }).length
})

onMounted(() => {
    logsStore.fetchLogs()
})
</script>