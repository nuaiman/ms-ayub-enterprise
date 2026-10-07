<!-- src/views/MajhisView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Majhis</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ majhisStore.majhis.length }}</p>
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

        <!-- Majhi List -->
        <div class="flex-1 min-h-0 mt-6">
            <MajhiList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useMajhisStore } from '@/stores/majhis'
import MajhiList from '@/components/features/majhis/MajhiList.vue'

const majhisStore = useMajhisStore()

const withPhoneCount = computed(() => {
    return majhisStore.majhis.filter(m => m.phone).length
})

const withNotesCount = computed(() => {
    return majhisStore.majhis.filter(m => m.notes).length
})

const createdThisMonthCount = computed(() => {
    const now = new Date()
    const currentMonth = now.getMonth()
    const currentYear = now.getFullYear()
    return majhisStore.majhis.filter(m => {
        const date = new Date(m.created_at)
        return date.getMonth() === currentMonth && date.getFullYear() === currentYear
    }).length
})

onMounted(() => {
    majhisStore.fetchMajhis()
})
</script>