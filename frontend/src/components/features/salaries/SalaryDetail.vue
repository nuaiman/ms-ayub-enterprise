<!-- src/components/features/salaries/SalaryDetail.vue -->
<template>
    <div v-if="salary" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Avatar -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ getEmployeeName(salary.employee_id) }}
                </h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ formatMonthYear(salary.month_year) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="getStatusBadgeClass(salary.status)">
                        <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(salary.status)"></span>
                        {{ salary.status.charAt(0).toUpperCase() + salary.status.slice(1) }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ salary.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(salary.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(salary.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
                <!-- Basic Salary -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Basic Salary
                    </p>
                    <p class="text-sm font-medium text-(--color-text-primary)">{{ formatCurrency(salary.basic_salary ||
                        0) }}</p>
                    <p class="text-xs text-(--color-text-secondary)">From employee profile</p>
                </div>

                <!-- Bonus -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Bonus</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(salary.bonus) }}</p>
                </div>

                <!-- Deductions -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Deductions</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(salary.deductions) }}</p>
                </div>
            </div>

            <div class="space-y-4">
                <!-- Total -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Total Salary
                    </p>
                    <p class="text-2xl font-bold text-(--color-blue)">{{ formatCurrency(salary.total_salary || 0) }}</p>
                    <p class="text-xs text-(--color-text-secondary)">
                        {{ formatCurrency(salary.basic_salary || 0) }} + {{ formatCurrency(salary.bonus) }} - {{
                            formatCurrency(salary.deductions) }}
                    </p>
                </div>

                <!-- Payment Details -->
                <div v-if="salary.status === 'paid'">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Payment
                        Details</p>
                    <div class="space-y-1">
                        <p class="text-sm text-(--color-text-primary)">
                            Method: {{ salary.payment_method ? formatPaymentMethod(salary.payment_method) : '—' }}
                        </p>
                        <p class="text-sm text-(--color-text-primary)">
                            Date: {{ salary.payment_date ? formatDate(salary.payment_date) : '—' }}
                        </p>
                        <p v-if="salary.reference_number" class="text-sm text-(--color-text-primary)">
                            Reference: {{ salary.reference_number }}
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Notes -->
        <div v-if="salary.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ salary.notes }}</p>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button v-if="salary.status === 'draft'" @click="emit('edit', salary)"
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
import { computed } from 'vue'
import type { Salary, SalaryStatus, PaymentMethod } from '@/types/salary'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    salary: Salary | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [salary: Salary]
    'updated': []
}>()

const usersStore = useUsersStore()

const getEmployeeName = (id?: number): string => {
    if (!id) return 'Unknown'
    return usersStore.getUserName(id)
}

const formatMonthYear = (monthYear: string): string => {
    const [year, month] = monthYear.split('-')
    if (!year || !month) return monthYear
    const date = new Date(parseInt(year), parseInt(month) - 1)
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    })
}

const formatPaymentMethod = (method: PaymentMethod): string => {
    const labels: Record<PaymentMethod, string> = {
        cash: 'Cash',
        bank_transfer: 'Bank Transfer',
        check: 'Check',
        mobile_banking: 'Mobile Banking'
    }
    return labels[method] || method
}

const getStatusBadgeClass = (status: SalaryStatus): string => {
    switch (status) {
        case 'draft':
            return 'border-(--color-yellow) text-(--color-yellow)'
        case 'paid':
            return 'border-(--color-green) text-(--color-green)'
        case 'cancelled':
            return 'border-(--color-red) text-(--color-red)'
        default:
            return 'border-(--color-border) text-(--color-text-secondary)'
    }
}

const getStatusDotClass = (status: SalaryStatus): string => {
    switch (status) {
        case 'draft':
            return 'bg-(--color-yellow)'
        case 'paid':
            return 'bg-(--color-green)'
        case 'cancelled':
            return 'bg-(--color-red)'
        default:
            return 'bg-(--color-text-secondary)'
    }
}
</script>