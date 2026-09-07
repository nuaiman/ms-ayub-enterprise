<!-- src/components/features/deliveryItems/DeliveryItemRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Delivery - 2 -->
        <div class="col-span-2 min-w-0 pr-3">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                #{{ deliveryItem.delivery_id }}
            </div>
            <span class="text-xs text-(--color-text-secondary) truncate block">
                {{ getDeliveryDate(deliveryItem.delivery_id) }}
            </span>
        </div>

        <!-- Item - 3 -->
        <div class="col-span-3 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ getItemName(deliveryItem.item_id) }}
            </span>
        </div>

        <!-- Lot - 2 -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ getLotName(deliveryItem.lot_id) }}
            </span>
        </div>

        <!-- Quantity - 1 -->
        <div class="col-span-1 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ deliveryItem.quantity }} {{ deliveryItem.quantity_unit }}
            </span>
        </div>

        <!-- Weight - 1 -->
        <div class="col-span-1 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ deliveryItem.weight }} {{ deliveryItem.weight_unit }}
            </span>
        </div>

        <!-- Majhi - 2 -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ getMajhiName(deliveryItem.majhi_id) }}
            </span>
        </div>

        <!-- Actions - 1 -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200 shrink-0">
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
import { ref } from 'vue'
import type { DeliveryItem } from '@/types/deliveryItem'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'

const props = defineProps<{
    deliveryItem: DeliveryItem
}>()

const emit = defineEmits<{
    'view': [deliveryItem: DeliveryItem]
    'edit': [deliveryItem: DeliveryItem]
    'delete': [deliveryItem: DeliveryItem]
    'updated': []
}>()

const deliveriesStore = useDeliveriesStore()
const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const isOpen = ref(false)

const getDeliveryDate = (id: number): string => {
    const delivery = deliveriesStore.getDeliveryById(id)
    if (!delivery) return ''
    return new Date(delivery.delivery_date).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric'
    })
}

const getItemName = (id: number): string => itemsStore.getItemName(id)
const getLotName = (id: number): string => lotsStore.getLotName(id)

const getMajhiName = (id: number | null): string => {
    if (!id) return '—'
    return majhisStore.getMajhiName(id)
}

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.deliveryItem) }
const handleEdit = () => { closeMenu(); emit('edit', props.deliveryItem) }
const handleDelete = () => { closeMenu(); emit('delete', props.deliveryItem) }
</script>