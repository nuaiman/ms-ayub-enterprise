<!-- src/components/features/dashboard/DashboardStats.vue -->
<template>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <!-- Total Received from Customers -->
        <div class="rounded-2xl p-5 bg-(--color-surface) border border-(--color-border)">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Received from
                        Customers</p>
                    <p class="text-2xl font-bold text-(--color-green)">{{ formatCurrency(stats.totalReceived) }}</p>
                </div>
                <div
                    class="w-12 h-12 rounded-xl bg-(--color-green)/10 text-(--color-green) flex items-center justify-center">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                </div>
            </div>
            <div class="flex items-center gap-4 mt-2 text-sm">
                <div>
                    <p class="text-xs text-(--color-text-secondary)">This Month</p>
                    <p class="font-semibold text-(--color-green)">{{ formatCurrency(stats.monthlyReceived) }}</p>
                </div>
                <div class="w-px h-8 bg-(--color-border)"></div>
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Collection Rate</p>
                    <p class="font-semibold"
                        :class="stats.collectionRate >= 70 ? 'text-(--color-green)' : 'text-(--color-yellow)'">
                        {{ stats.collectionRate }}%
                    </p>
                </div>
            </div>
        </div>

        <!-- Total Due from Customers -->
        <div class="rounded-2xl p-5 bg-(--color-surface) border border-(--color-border)">
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Due from Customers
                    </p>
                    <p class="text-2xl font-bold text-(--color-red)">{{ formatCurrency(stats.totalDue) }}</p>
                </div>
                <div
                    class="w-12 h-12 rounded-xl bg-(--color-red)/10 text-(--color-red) flex items-center justify-center">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                </div>
            </div>
            <div class="flex items-center gap-4 mt-2 text-sm">
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Overdue</p>
                    <p class="font-semibold text-(--color-red)">{{ formatCurrency(stats.overdue) }}</p>
                </div>
                <div class="w-px h-8 bg-(--color-border)"></div>
                <div>
                    <p class="text-xs text-(--color-text-secondary)">Customers</p>
                    <p class="font-semibold text-(--color-red)">{{ stats.outstandingCustomers }} customers</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { DashboardMetrics } from '@/types/dashboard'
import { formatCurrency } from '@/utils/currency'

defineProps<{
    stats: DashboardMetrics
}>()
</script>