<!-- src/components/features/stores/StoreRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Lot / Product - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ lotDisplayName }}
            </div>
            <div class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                Lot #{{ lotNumber }}
            </div>
        </div>

        <!-- Customer - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ customerName }}
            </span>
            <span v-if="customerPhone" class="text-xs text-(--color-text-secondary)/70 truncate block mt-0.5">
                {{ customerPhone }}
            </span>
        </div>

        <!-- Godown - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ godownName }}
            </span>
            <span v-if="godownPhone" class="text-xs text-(--color-text-secondary)/70 truncate block mt-0.5">
                {{ godownPhone }}
            </span>
        </div>

        <!-- Stock - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-primary) block truncate">
                {{ formatNumber(store.quantity) }} {{ store.quantity_unit }}
            </span>
            <span class="text-xs text-(--color-text-secondary)/70 block truncate mt-0.5">
                {{ formatNumber(store.weight) }} {{ store.weight_unit }}
            </span>
        </div>

        <!-- Status - 1 column -->
        <div class="col-span-1 min-w-0 pr-3">
            <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="store.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                <span class="w-1.5 h-1.5 rounded-full"
                    :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                {{ store.is_active ? 'Active' : 'Inactive' }}
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

                    <button @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>

                    <button @click="handleReadd"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                        </svg>
                        Re-add Store
                    </button>

                    <button @click="handleToggleActive"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-(--color-muted-bg) transition-colors"
                        :class="store.is_active ? 'text-(--color-yellow)' : 'text-(--color-green)'">
                        <svg v-if="store.is_active" class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                        </svg>
                        <svg v-else class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        {{ store.is_active ? 'Deactivate' : 'Activate' }}
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
import type { Store } from '@/types/store'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useGodownsStore } from '@/stores/godowns'

const props = defineProps<{
    store: Store
}>()

const emit = defineEmits<{
    'view': [store: Store]
    'edit': [store: Store]
    'readd': [store: Store]
    'delete': [store: Store]
    'toggle-active': [store: Store]
    'updated': []
}>()

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const godownsStore = useGodownsStore()
const isOpen = ref(false)

const lot = computed(() => lotsStore.getLotById(props.store.lot_id))
const lotNumber = computed(() => lot.value?.lot_number ?? '৳')

const lotDisplayName = computed(() => {
    if (!lot.value) return '৳'
    return lotsStore.getLotDisplayName(lot.value)
})

const customerName = computed(() => {
    if (!lot.value?.customer_id) return '৳'
    return customersStore.getCustomerName(lot.value.customer_id)
})

const customerPhone = computed(() => {
    if (!lot.value?.customer_id) return ''
    const c = customersStore.getCustomerById(lot.value.customer_id)
    return c?.phone || ''
})

const godownName = computed(() => godownsStore.getGodownName(props.store.godown_id))

const godownPhone = computed(() => {
    const g = godownsStore.getGodownById(props.store.godown_id)
    return g?.phone || ''
})

const formatNumber = (n: number): string => {
    return new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)
}

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.store) }
const handleEdit = () => { closeMenu(); emit('edit', props.store) }
const handleReadd = () => { closeMenu(); emit('readd', props.store) }
const handleToggleActive = () => { closeMenu(); emit('toggle-active', props.store) }
const handleDelete = () => { closeMenu(); emit('delete', props.store) }
</script>