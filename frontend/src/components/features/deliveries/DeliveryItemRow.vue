<!-- src/components/features/deliveries/DeliveryItemRow.vue -->
<template>
    <div
        class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30">
        <!-- Lot / Store - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="min-w-0">
                <div class="font-medium text-(--color-text-primary) truncate text-sm">
                    Store #{{ item.store_id }}
                </div>
                <div v-if="lotLabel" class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                    {{ lotLabel }}
                </div>
            </div>
        </div>

        <!-- Majhi / Vehicle - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <div class="text-xs text-(--color-text-secondary) truncate">
                {{ majhiName || '—' }}
            </div>
            <div v-if="item.vehicle_number" class="text-xs text-(--color-text-secondary)/70 truncate mt-0.5">
                {{ item.vehicle_number }}
            </div>
        </div>

        <!-- Quantity + Weight - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <span class="text-xs text-(--color-text-secondary) truncate block">
                Qty: {{ formatNumber(item.quantity) }} {{ quantityUnit }}
            </span>
            <span class="text-xs text-(--color-text-secondary) truncate block mt-0.5">
                Wt: {{ formatNumber(item.weight) }} {{ weightUnit }}
            </span>
        </div>

        <!-- Actions - 2 columns -->
        <div class="col-span-2 flex items-center justify-end relative" @click.stop>
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
                    class="absolute right-0 top-9 w-44 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
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
import type { DeliveryItem } from '@/types/delivery'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'

const props = defineProps<{
    item: DeliveryItem
}>()

const emit = defineEmits<{
    'edit': [item: DeliveryItem]
    'delete': [item: DeliveryItem]
}>()

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const isOpen = ref(false)

const store = computed(() => storesStore.getStoreById(props.item.store_id))
const lot = computed(() => {
    if (!store.value) return null
    return lotsStore.getLotById(store.value.lot_id) ?? null
})

const lotLabel = computed(() => {
    if (!lot.value) return ''
    return `${lot.value.product_name} · Lot ${lot.value.lot_number}`
})

const weightUnit = computed(() => lot.value?.weight_unit ?? '')
const quantityUnit = computed(() => lot.value?.quantity_unit ?? '')

const majhiName = computed(() => {
    if (!props.item.majhi_id) return ''
    return majhisStore.getMajhiName(props.item.majhi_id)
})

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleEdit = () => { closeMenu(); emit('edit', props.item) }
const handleDelete = () => { closeMenu(); emit('delete', props.item) }
</script>