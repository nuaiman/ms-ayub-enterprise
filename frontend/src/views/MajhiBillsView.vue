<!-- src/views/MajhiBillsView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ majhiBillsStore.totalBills }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Billed</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">
                    {{ formatCurrency(majhiBillsStore.totalBilled) }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">
                    {{ formatCurrency(majhiBillsStore.totalPaid) }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                <p class="text-2xl font-bold mt-1"
                    :class="outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(outstanding) }}
                </p>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 min-h-0 mt-6">
            <MajhiBillList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import MajhiBillList from '@/components/features/majhiBills/MajhiBillList.vue'
import { formatCurrency } from '@/utils/currency'

const majhiBillsStore = useMajhiBillsStore()

const outstanding = computed(() =>
    Math.max(0, majhiBillsStore.totalBilled - majhiBillsStore.totalPaid)
)

onMounted(() => {
    majhiBillsStore.fetchMajhiBills()
})
</script>