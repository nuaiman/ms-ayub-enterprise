<!-- src/components/features/customers/CustomerRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Company / Contact - 3 columns -->
        <div class="col-span-3">
            <div class="min-w-0">
                <div class="font-medium text-(--color-text-primary) truncate text-sm">
                    {{ customer.company_name || customer.contact_person || 'Unnamed' }}
                </div>
                <div v-if="customer.contact_person && customer.company_name"
                    class="text-xs text-(--color-text-secondary) truncate">
                    {{ customer.contact_person }}
                </div>
            </div>
        </div>

        <!-- Phone - 2 columns -->
        <div class="col-span-2">
            <span class="text-sm text-(--color-text-secondary)">{{ customer.phone }}</span>
        </div>

        <!-- Email - 3 columns -->
        <div class="col-span-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">{{ customer.email || '—' }}</span>
        </div>

        <!-- Address - 3 columns -->
        <div class="col-span-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">{{ customer.address || '—' }}</span>
        </div>

        <!-- Actions - 1 column, right aligned -->
        <div class="col-span-1 flex justify-end relative" @click.stop>
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
                    class="absolute right-0 top-9 w-52 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
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

                    <button @click="handleViewLedger"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        View Ledger
                    </button>

                    <button @click="handleAddCharge"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 5v14M5 12h14" />
                        </svg>
                        Add Charge
                    </button>

                    <button @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
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
import type { Customer } from '@/types/customer'

const props = defineProps<{
    customer: Customer
}>()

const emit = defineEmits<{
    'view': [customer: Customer]
    'edit': [customer: Customer]
    'view-ledger': [customer: Customer]
    'add-charge': [customer: Customer]
    'updated': []
}>()

const isOpen = ref(false)

const toggleMenu = () => {
    isOpen.value = !isOpen.value
}

const closeMenu = () => {
    isOpen.value = false
}

const handleView = () => {
    closeMenu()
    emit('view', props.customer)
}

const handleViewLedger = () => {
    closeMenu()
    emit('view-ledger', props.customer)
}

const handleAddCharge = () => {
    closeMenu()
    emit('add-charge', props.customer)
}

const handleEdit = () => {
    closeMenu()
    emit('edit', props.customer)
}
</script>