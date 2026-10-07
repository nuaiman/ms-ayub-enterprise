<!-- src/components/features/incomes/IncomeDetail.vue -->
<template>
    <div v-if="income" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-green)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-green)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ income.title }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDate(income.income_date) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-green)">{{ formatCurrency(income.amount) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">By {{ getUserName(income.user_id) }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ income.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDateTime(income.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDateTime(income.updated_at) }}</span>
        </div>

        <!-- Notes -->
        <div v-if="income.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ income.notes }}</p>
            </div>
        </div>

        <!-- Image -->
        <div v-if="income.image_url" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Attachment</p>
            <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md">
                <img :src="getImageUrl(income.image_url)" alt="Income attachment"
                    class="w-full object-cover max-h-64" />
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', income)"
                class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
                Edit
            </button>
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Income } from '@/types/income'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    income: Income | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [income: Income]
    'updated': []
}>()

const usersStore = useUsersStore()

const getUserName = (userId: number): string => {
    return usersStore.getUserName(userId)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'long',
        day: 'numeric',
        year: 'numeric'
    })
}

const formatDateTime = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}
</script>