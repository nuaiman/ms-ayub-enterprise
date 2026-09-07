<!-- src/components/features/rents/RentDetail.vue -->
<template>
    <div v-if="rent" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ getGodownName(rent.godown_id) }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ formatMonthYear(rent.month_year) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="getStatusBadgeClass(rent.status)">
                        <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(rent.status)"></span>
                        {{ rent.status.charAt(0).toUpperCase() + rent.status.slice(1) }}
                    </span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(rent.amount)
                        }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ rent.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(rent.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(rent.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
                <!-- Godown -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
                    <p class="text-sm text-(--color-text-primary)">{{ getGodownName(rent.godown_id) }}</p>
                </div>

                <!-- Month -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Month</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatMonthYear(rent.month_year) }}</p>
                </div>

                <!-- Amount -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                    <p class="text-lg font-semibold text-(--color-text-primary)">{{ formatCurrency(rent.amount) }}</p>
                </div>
            </div>

            <div class="space-y-4">
                <!-- Status -->
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Status</p>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="getStatusBadgeClass(rent.status)">
                        <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(rent.status)"></span>
                        {{ rent.status.charAt(0).toUpperCase() + rent.status.slice(1) }}
                    </span>
                </div>

                <!-- Payment Details (if paid) -->
                <div v-if="rent.status === 'paid'" class="pt-2 border-t border-(--color-border)">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Payment
                        Details</p>
                    <div class="space-y-1 mt-1">
                        <p class="text-sm text-(--color-text-primary)">Method: {{ rent.payment_method ?
                            formatPaymentMethod(rent.payment_method) : '—' }}</p>
                        <p class="text-sm text-(--color-text-primary)">Date: {{ rent.payment_date ?
                            formatDate(rent.payment_date) : '—' }}</p>
                        <p v-if="rent.reference_number" class="text-sm text-(--color-text-primary)">Reference: {{
                            rent.reference_number }}</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Notes -->
        <div v-if="rent.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ rent.notes }}</p>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button v-if="rent.status === 'draft'" @click="emit('edit', rent)"
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
import type { Rent, RentStatus, PaymentMethod } from '@/types/rent'
import { useGodownsStore } from '@/stores/godowns'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    rent: Rent | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [rent: Rent]
    'updated': []
}>()

const godownsStore = useGodownsStore()

const getGodownName = (id: number): string => {
    return godownsStore.getGodownName(id)
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
        minute: '2-digit'
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

const getStatusBadgeClass = (status: RentStatus): string => {
    switch (status) {
        case 'draft': return 'border-(--color-yellow) text-(--color-yellow)'
        case 'paid': return 'border-(--color-green) text-(--color-green)'
        case 'cancelled': return 'border-(--color-red) text-(--color-red)'
        default: return 'border-(--color-border) text-(--color-text-secondary)'
    }
}

const getStatusDotClass = (status: RentStatus): string => {
    switch (status) {
        case 'draft': return 'bg-(--color-yellow)'
        case 'paid': return 'bg-(--color-green)'
        case 'cancelled': return 'bg-(--color-red)'
        default: return 'bg-(--color-text-secondary)'
    }
}
</script>