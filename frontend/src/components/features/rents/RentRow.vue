<!-- src/components/features/rents/RentRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Godown - 3 columns -->
        <div class="col-span-3 min-w-0">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ getGodownName(rent.godown_id) }}
            </div>
        </div>

        <!-- Month - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ formatMonthYear(rent.month_year) }}
            </span>
        </div>

        <!-- Amount - 2 columns -->
        <div class="col-span-2">
            <span class="text-sm font-semibold text-(--color-text-primary)">
                {{ formatCurrency(rent.amount) }}
            </span>
        </div>

        <!-- Status - 2 columns -->
        <div class="col-span-2">
            <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="getStatusBadgeClass(rent.status)">
                <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(rent.status)"></span>
                {{ rent.status.charAt(0).toUpperCase() + rent.status.slice(1) }}
            </span>
        </div>

        <!-- Notes - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ rent.notes || '—' }}
            </span>
        </div>

        <!-- Actions - 1 column -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <!-- Dropdown -->
            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
                    <!-- View -->
                    <button @click="handleView"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                        </svg>
                        View Details
                    </button>

                    <!-- Edit - only draft -->
                    <button v-if="rent.status === 'draft'" @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>

                    <!-- Mark Paid - only draft -->
                    <button v-if="rent.status === 'draft'" @click="handlePay"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                        Mark as Paid
                    </button>

                    <!-- Delete - only draft -->
                    <button v-if="rent.status === 'draft'" @click="handleDelete"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-red) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        Delete
                    </button>

                    <!-- Cancel - draft or paid -->
                    <button v-if="rent.status === 'draft' || rent.status === 'paid'" @click="handleCancel"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-yellow) hover:bg-(--color-muted-bg) transition-colors border-t border-(--color-border) mt-1 pt-1">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                        Cancel
                    </button>
                </div>
            </Transition>

            <!-- Backdrop -->
            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Rent, RentStatus } from '@/types/rent'
import { useGodownsStore } from '@/stores/godowns'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    rent: Rent
}>()

const emit = defineEmits<{
    'view': [rent: Rent]
    'edit': [rent: Rent]
    'delete': [rent: Rent]
    'pay': [rent: Rent]
    'cancel': [rent: Rent]
}>()

const godownsStore = useGodownsStore()
const isOpen = ref(false)

const getGodownName = (id: number): string => {
    return godownsStore.getGodownName(id)
}

const formatMonthYear = (monthYear: string): string => {
    const [year, month] = monthYear.split('-')
    if (!year || !month) return monthYear
    const date = new Date(parseInt(year), parseInt(month) - 1)
    return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
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

const toggleMenu = () => {
    isOpen.value = !isOpen.value
}

const closeMenu = () => {
    isOpen.value = false
}

const handleView = () => {
    closeMenu()
    emit('view', props.rent)
}

const handleEdit = () => {
    closeMenu()
    emit('edit', props.rent)
}

const handleDelete = () => {
    closeMenu()
    emit('delete', props.rent)
}

const handlePay = () => {
    closeMenu()
    emit('pay', props.rent)
}

const handleCancel = () => {
    closeMenu()
    emit('cancel', props.rent)
}
</script>