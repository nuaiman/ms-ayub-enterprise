<!-- src/components/features/lots/LotRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 bg-(--color-muted-bg)/30 cursor-pointer transition-colors duration-150"
        @click="handleView">
        <!-- Customer - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ customerName }}
            </span>
            <span v-if="customerPhone" class="text-xs text-(--color-text-secondary)/70 truncate block mt-0.5">
                {{ customerPhone }}
            </span>
        </div>

        <!-- Lot / Product - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="min-w-0">
                <div class="font-semibold text-(--color-text-primary) truncate text-sm">
                    {{ lot.product_name }}
                </div>
                <div class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                    Lot {{ lot.lot_number }}
                </div>
            </div>
        </div>

        <!-- Units - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-xs text-(--color-text-secondary) truncate block">
                Wt: {{ lot.weight_unit }}
            </span>
            <span class="text-xs text-(--color-text-secondary) truncate block mt-0.5">
                Qty: {{ lot.quantity_unit }}
            </span>
        </div>

        <!-- Stores count - 1 column -->
        <div class="col-span-1 min-w-0 pr-3">
            <span
                class="inline-flex items-center justify-center min-w-6 h-6 px-1.5 rounded-md bg-(--color-blue)/10 text-(--color-blue) text-xs font-semibold">
                {{ storesCount }}
            </span>
        </div>

        <!-- Actions - 1 column -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md bg-(--color-surface) hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
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

                    <button @click="handleAddStore"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 5v14M5 12h14" />
                        </svg>
                        Add Store
                    </button>

                    <!-- Transfer Account -->
                    <button @click="handleTransfer"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M14 5l7 7m0 0l-7 7m7-7H3" />
                        </svg>
                        Transfer Account
                    </button>

                    <button @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>

                    <button @click="handleDelete"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-red) hover:bg-(--color-muted-bg) transition-colors border-t border-(--color-border) mt-1 pt-1">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        Delete
                    </button>
                </div>
            </Transition>

            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Lot } from '@/types/lot'
import { useCustomersStore } from '@/stores/customers'

const props = defineProps<{
    lot: Lot
    storesCount: number
}>()

const emit = defineEmits<{
    'view': [lot: Lot]
    'edit': [lot: Lot]
    'delete': [lot: Lot]
    'add-store': [lot: Lot]
    'transfer': [lot: Lot]
    'updated': []
}>()

const customersStore = useCustomersStore()
const isOpen = ref(false)

const customerName = computed(() =>
    customersStore.getCustomerName(props.lot.customer_id)
)

const customerPhone = computed(() => {
    const c = customersStore.getCustomerById(props.lot.customer_id)
    return c?.phone || ''
})

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.lot) }
const handleEdit = () => { closeMenu(); emit('edit', props.lot) }
const handleDelete = () => { closeMenu(); emit('delete', props.lot) }
const handleAddStore = () => { closeMenu(); emit('add-store', props.lot) }
const handleTransfer = () => { closeMenu(); emit('transfer', props.lot) }
</script>