<!-- src/components/features/dashboard/OutstandingSnapshot.vue -->
<template>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <!-- Top Outstanding Customers -->
        <div class="rounded-2xl p-4 bg-(--color-surface) border border-(--color-border)">
            <div class="flex items-center justify-between mb-2">
                <p class="text-sm font-medium text-(--color-text-primary)">📋 Top Due Customers</p>
                <span class="text-xs text-(--color-red)">{{ formatCurrency(totalDue) }}</span>
            </div>
            <div class="space-y-1.5">
                <div v-for="item in outstandingByCustomer" :key="item.name"
                    class="flex items-center justify-between text-sm">
                    <span class="text-(--color-text-secondary) truncate max-w-28">{{ item.name }}</span>
                    <span class="font-medium text-(--color-red)">{{ formatCurrency(item.amount) }}</span>
                </div>
                <div v-if="outstandingByCustomer.length === 0"
                    class="text-sm text-(--color-text-secondary) text-center py-2">
                    No outstanding customers
                </div>
            </div>
        </div>

        <!-- Outstanding by Type -->
        <div class="rounded-2xl p-4 bg-(--color-surface) border border-(--color-border)">
            <div class="flex items-center justify-between mb-2">
                <p class="text-sm font-medium text-(--color-text-primary)">By Bill Type</p>
                <span class="text-xs text-(--color-text-secondary)">{{ outstandingByType.length }} types</span>
            </div>
            <div class="space-y-2">
                <div v-for="item in outstandingByType" :key="item.name" class="space-y-0.5">
                    <div class="flex items-center justify-between text-xs">
                        <span class="text-(--color-text-secondary)">{{ item.name }}</span>
                        <span class="font-medium text-(--color-red)">{{ formatCurrency(item.amount) }}</span>
                    </div>
                    <div class="w-full h-1 bg-(--color-muted-bg) rounded-full overflow-hidden">
                        <div class="h-full rounded-full"
                            :style="{ width: item.percentage + '%', backgroundColor: item.color }"></div>
                    </div>
                </div>
                <div v-if="outstandingByType.length === 0"
                    class="text-sm text-(--color-text-secondary) text-center py-2">
                    No outstanding bills
                </div>
            </div>
        </div>

        <!-- Payment Obligations -->
        <div class="rounded-2xl p-4 bg-(--color-surface) border border-(--color-border)">
            <p class="text-sm font-medium text-(--color-text-primary) mb-2">👷 Due Payments</p>
            <div class="space-y-2">
                <div class="flex items-center justify-between text-sm">
                    <span class="text-(--color-text-secondary)">Majhi</span>
                    <span class="font-medium text-(--color-red)">{{ formatCurrency(majhiUnpaidAmount) }}</span>
                </div>
                <div class="flex items-center justify-between text-sm">
                    <span class="text-(--color-text-secondary)">Brokers</span>
                    <span class="font-medium text-(--color-red)">{{ formatCurrency(brokerUnpaidAmount) }}</span>
                </div>
                <div class="flex items-center justify-between text-sm">
                    <span class="text-(--color-text-secondary)">Rent</span>
                    <span class="font-medium text-(--color-red)">{{ formatCurrency(rentUnpaidAmount) }}</span>
                </div>
                <div class="flex items-center justify-between text-sm pt-1 border-t border-(--color-border)">
                    <span class="text-(--color-text-secondary) font-medium">Total Due</span>
                    <span class="font-bold text-(--color-red)">{{ formatCurrency(totalDuePayments) }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { OutstandingCustomer, OutstandingByType } from '@/types/dashboard'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    totalDue: number
    outstandingByCustomer: OutstandingCustomer[]
    outstandingByType: OutstandingByType[]
    majhiUnpaidAmount: number
    brokerUnpaidAmount: number
    rentUnpaidAmount: number
}>()

// Fix: Use computed instead of a function
const totalDuePayments = computed(() => {
    return props.majhiUnpaidAmount + props.brokerUnpaidAmount + props.rentUnpaidAmount
})
</script>