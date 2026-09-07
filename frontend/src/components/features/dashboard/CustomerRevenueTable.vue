<!-- src/components/features/dashboard/CustomerRevenueTable.vue -->
<template>
    <div
        class="rounded-2xl bg-(--color-surface) border border-(--color-border) overflow-hidden flex-1 flex flex-col min-h-0">
        <div class="p-4 pb-0 flex flex-col sm:flex-row sm:items-center justify-between gap-2 shrink-0">
            <div>
                <h3 class="text-sm font-semibold text-(--color-text-primary)">👥 Customer Revenue Details</h3>
                <p class="text-xs text-(--color-text-secondary)">Complete breakdown by customer</p>
            </div>
            <select :value="modelValue" @change="handleMonthChange"
                class="px-3 py-1.5 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue)">
                <option v-for="month in months" :key="month.value" :value="month.value">
                    {{ month.label }}
                </option>
            </select>
        </div>

        <div class="overflow-y-auto flex-1 min-h-0 mt-2">
            <table class="w-full text-sm">
                <thead class="sticky top-0 z-10">
                    <tr class="border-y border-(--color-border) bg-(--color-muted-bg)">
                        <th
                            class="text-left py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Customer</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Storage</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Unload</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Delivery</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Transport</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Total</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Paid</th>
                        <th
                            class="text-right py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Due</th>
                        <th
                            class="text-center py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            %</th>
                        <th
                            class="text-center py-2.5 px-4 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                            Status</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(customer, index) in customers" :key="index"
                        class="border-b border-(--color-border) hover:bg-(--color-muted-bg)/20 transition-colors">
                        <td class="py-2.5 px-4">
                            <div class="flex items-center gap-2.5">
                                <div
                                    class="w-7 h-7 rounded-full bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center text-[10px] font-bold shrink-0">
                                    {{ customer.initials }}
                                </div>
                                <div>
                                    <p class="font-medium text-(--color-text-primary) text-sm">{{ customer.name }}</p>
                                    <p class="text-[10px] text-(--color-text-secondary)">{{ customer.phone }}</p>
                                </div>
                            </div>
                        </td>
                        <td class="py-2.5 px-4 text-right text-(--color-text-secondary) text-sm">{{
                            formatCurrency(customer.storage) }}</td>
                        <td class="py-2.5 px-4 text-right text-(--color-text-secondary) text-sm">{{
                            formatCurrency(customer.unload) }}</td>
                        <td class="py-2.5 px-4 text-right text-(--color-text-secondary) text-sm">{{
                            formatCurrency(customer.delivery) }}</td>
                        <td class="py-2.5 px-4 text-right text-(--color-text-secondary) text-sm">{{
                            formatCurrency(customer.transport) }}</td>
                        <td class="py-2.5 px-4 text-right font-medium text-(--color-text-primary) text-sm">{{
                            formatCurrency(customer.total) }}</td>
                        <td class="py-2.5 px-4 text-right font-medium text-(--color-green) text-sm">{{
                            formatCurrency(customer.paid) }}</td>
                        <td class="py-2.5 px-4 text-right font-medium text-sm"
                            :class="customer.due > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                            {{ formatCurrency(customer.due) }}
                        </td>
                        <td class="py-2.5 px-4">
                            <div class="flex items-center justify-center gap-1.5">
                                <div class="w-12 h-1 bg-(--color-muted-bg) rounded-full overflow-hidden">
                                    <div class="h-full rounded-full"
                                        :class="customer.paymentPercentage >= 80 ? 'bg-(--color-green)' : customer.paymentPercentage >= 50 ? 'bg-(--color-yellow)' : 'bg-(--color-red)'"
                                        :style="{ width: customer.paymentPercentage + '%' }"></div>
                                </div>
                                <span class="text-[10px] text-(--color-text-secondary) w-7 text-right">{{
                                    customer.paymentPercentage }}%</span>
                            </div>
                        </td>
                        <td class="py-2.5 px-4 text-center">
                            <span
                                class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium border"
                                :class="customer.due > 0 ? 'border-(--color-yellow) text-(--color-yellow)' : 'border-(--color-green) text-(--color-green)'">
                                <span class="w-1 h-1 rounded-full"
                                    :class="customer.due > 0 ? 'bg-(--color-yellow)' : 'bg-(--color-green)'"></span>
                                {{ customer.due > 0 ? 'Partial' : 'Paid' }}
                            </span>
                        </td>
                    </tr>
                    <tr v-if="customers.length === 0">
                        <td colspan="10" class="text-center py-8 text-(--color-text-secondary)">
                            No customer data available
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { CustomerRevenue, MonthOption } from '@/types/dashboard'
import { formatCurrency } from '@/utils/currency'

defineProps<{
    customers: CustomerRevenue[]
    months: MonthOption[]
    modelValue: string
}>()

const emit = defineEmits<{
    'update:modelValue': [value: string]
}>()

const handleMonthChange = (event: Event) => {
    const target = event.target as HTMLSelectElement
    emit('update:modelValue', target.value)
}
</script>