<!-- src/components/features/stores/StoreRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Lot / Product - 4 columns (image only when nested) -->
        <div class="col-span-4 min-w-0 pr-3">
            <!-- Nested: show image + placeholder -->
            <div v-if="hideLot" class="flex items-center gap-3">
                <div class="shrink-0">
                    <div v-if="store.image_url"
                        class="w-9 h-9 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="getImageUrl(store.image_url)" :alt="lotDisplayName"
                            class="w-full h-full object-cover" />
                    </div>
                    <div v-else
                        class="w-9 h-9 rounded-lg bg-(--color-muted-bg) border border-(--color-border) flex items-center justify-center">
                        <svg class="w-4 h-4 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                        </svg>
                    </div>
                </div>
            </div>

            <!-- Standalone: lot name + number -->
            <div v-else class="min-w-0">
                <div class="font-medium text-(--color-text-primary) truncate text-sm">
                    {{ lotDisplayName }}
                </div>
                <div class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                    Lot {{ lotNumber }}
                </div>
            </div>
        </div>

        <!-- Godown - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ godownName }}
            </span>
        </div>

        <!-- Stock - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-xs text-(--color-text-secondary) truncate block">
                Wt: {{ formatNumber(store.weight) }} {{ weightUnit }}
            </span>
            <span class="text-xs text-(--color-text-secondary) truncate block mt-0.5">
                Qty: {{ formatNumber(store.quantity) }} {{ quantityUnit }}
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

                    <!-- Adjust Stock -->
                    <button @click="handleAdjust"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-blue) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        Adjust Stock
                    </button>

                    <!-- Transfer Store -->
                    <button @click="handleTransfer" :disabled="!store.is_active"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-green) hover:bg-(--color-muted-bg) transition-colors disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-transparent">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M14 5l7 7m0 0l-7 7m7-7H3" />
                        </svg>
                        Transfer Store
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
import { useGodownsStore } from '@/stores/godowns'
import { getImageUrl } from '@/utils/image'

const props = withDefaults(defineProps<{
    store: Store
    hideLot?: boolean
}>(), {
    hideLot: false,
})

const emit = defineEmits<{
    'view': [store: Store]
    'edit': [store: Store]
    'delete': [store: Store]
    'adjust': [store: Store]
    'transfer': [store: Store]
    'updated': []
}>()

const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()
const isOpen = ref(false)

const lot = computed(() => lotsStore.getLotById(props.store.lot_id))
const lotNumber = computed(() => lot.value?.lot_number ?? '—')
const lotDisplayName = computed(() => lot.value ? lot.value.product_name : '—')
const quantityUnit = computed(() => lot.value?.quantity_unit ?? 'units')
const weightUnit = computed(() => lot.value?.weight_unit ?? 'kg')

const godownName = computed(() => godownsStore.getGodownName(props.store.godown_id))

const formatNumber = (n: number): string => {
    return new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)
}

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.store) }
const handleEdit = () => { closeMenu(); emit('edit', props.store) }
const handleAdjust = () => { closeMenu(); emit('adjust', props.store) }
const handleTransfer = () => { closeMenu(); emit('transfer', props.store) }
const handleDelete = () => { closeMenu(); emit('delete', props.store) }
</script>