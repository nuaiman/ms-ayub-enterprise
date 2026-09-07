<!-- src/views/GodownStoreBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Active Stores</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{
                        godownStoreBillsStore.storeBillData.length }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Monthly Bill</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(godownStoreBillsStore.totalMonthlyBill) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{
                        formatCurrency(godownStoreBillsStore.totalPaid)
                        }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(godownStoreBillsStore.totalOutstanding) }}</p>
                </div>
            </div>

            <!-- Godown Store Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <GodownStoreBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { useGodownStoreBillsStore } from '@/stores/godownStoreBills'
import AppLayout from '@/components/layouts/AppLayout.vue'
import GodownStoreBillList from '@/components/features/godownStoreBills/GodownStoreBillList.vue'
import { formatCurrency } from '@/utils/currency'

const godownStoreBillsStore = useGodownStoreBillsStore()
</script>