<!-- src/components/features/deliveries/DeliveryRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Customer - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <div class="flex items-center gap-3">
                <div class="shrink-0">
                    <div v-if="delivery.image_url"
                        class="w-9 h-9 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="getImageUrl(delivery.image_url)" :alt="customerName"
                            class="w-full h-full object-cover" />
                    </div>
                    <div v-else
                        class="w-9 h-9 rounded-lg bg-(--color-blue)/10 border border-(--color-border) flex items-center justify-center">
                        <svg class="w-4 h-4 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                        </svg>
                    </div>
                </div>
                <div class="min-w-0">
                    <div class="font-medium text-(--color-text-primary) truncate text-sm">
                        {{ customerName }}
                    </div>
                    <div v-if="customerPhone" class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                        {{ customerPhone }}
                    </div>
                </div>
            </div>
        </div>

        <!-- Receiver - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ delivery.receiver_name || 'â€”' }}
            </span>
            <span v-if="delivery.receiver_phone" class="text-xs text-(--color-text-secondary)/70 truncate block mt-0.5">
                {{ delivery.receiver_phone }}
            </span>
        </div>

        <!-- From - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ delivery.from_location || 'â€”' }}
            </span>
        </div>

        <!-- To - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ delivery.to_location || 'â€”' }}
            </span>
        </div>

        <!-- Items / Date - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-primary) block truncate">
                {{ itemsCount }} item(s)
            </span>
            <span class="text-xs text-(--color-text-secondary)/70 block truncate mt-0.5">
                {{ formatDateShort(delivery.delivery_date) }}
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

                    <button @click="handleManageItems"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                        </svg>
                        Manage Items
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
import type { Delivery } from '@/types/delivery'
import { useCustomersStore } from '@/stores/customers'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    delivery: Delivery
}>()

const emit = defineEmits<{
    'view': [delivery: Delivery]
    'edit': [delivery: Delivery]
    'delete': [delivery: Delivery]
    'manage-items': [delivery: Delivery]
    'updated': []
}>()

const customersStore = useCustomersStore()
const deliveryItemsStore = useDeliveryItemsStore()
const isOpen = ref(false)

const customerName = computed(() => {
    if (!props.delivery.customer_id) return 'â€”'
    return customersStore.getCustomerName(props.delivery.customer_id)
})

const customerPhone = computed(() => {
    if (!props.delivery.customer_id) return ''
    const c = customersStore.getCustomerById(props.delivery.customer_id)
    return c?.phone || ''
})

const itemsCount = computed(() => {
    return deliveryItemsStore.getDeliveryItemsByDeliveryId(props.delivery.id).length
})

const formatDateShort = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })
}

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.delivery) }
const handleEdit = () => { closeMenu(); emit('edit', props.delivery) }
const handleManageItems = () => { closeMenu(); emit('manage-items', props.delivery) }
const handleDelete = () => { closeMenu(); emit('delete', props.delivery) }
</script>