<!-- src/views/GodownsView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Godowns</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ godownsStore.godowns.length }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Phone</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withPhoneCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Notes</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ withNotesCount }}</p>
            </div>
        </div>

        <div class="flex-1 min-h-0 mt-6">
            <GodownList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useGodownsStore } from '@/stores/godowns'
import GodownList from '@/components/features/godowns/GodownList.vue'

const godownsStore = useGodownsStore()

const withPhoneCount = computed(() => godownsStore.godowns.filter(g => g.phone).length)
const withNotesCount = computed(() => godownsStore.godowns.filter(g => g.notes).length)

onMounted(() => {
    godownsStore.fetchGodowns()
})
</script>