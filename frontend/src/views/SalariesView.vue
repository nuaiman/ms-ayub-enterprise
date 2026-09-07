<!-- src/views/SalariesView.vue -->
<template>
    <AppLayout>
        <div class="space-y-6">


            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Salaries</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ salariesStore.totalSalaries }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Draft</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ salariesStore.draftCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Paid</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ salariesStore.statusCounts.paid || 0 }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Cancelled</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{ salariesStore.statusCounts.cancelled || 0
                    }}</p>
                </div>
            </div>

            <!-- Salary List -->
            <div class="flex-1 min-h-0 mt-6">
                <SalaryList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSalariesStore } from '@/stores/salaries'
import AppLayout from '@/components/layouts/AppLayout.vue'
import SalaryList from '@/components/features/salaries/SalaryList.vue'
import SalaryForm from '@/components/features/salaries/SalaryForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'

const auth = useAuthStore()
const salariesStore = useSalariesStore()

const showCreateDialog = ref(false)

const canManageSalaries = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const handleSalaryCreated = async () => {
    showCreateDialog.value = false
    await salariesStore.fetchSalaries()
}

onMounted(async () => {
    await salariesStore.fetchSalaries()
})
</script>