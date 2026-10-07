<!-- src/components/features/damages/DamageRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Store - 4 columns -->
        <div class="col-span-4 min-w-0 pr-3">
            <div class="flex items-center gap-3">
                <div class="shrink-0">
                    <div v-if="damage.image_url"
                        class="w-9 h-9 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="getImageUrl(damage.image_url)" :alt="damage.reason"
                            class="w-full h-full object-cover" />
                    </div>
                    <div v-else
                        class="w-9 h-9 rounded-lg bg-(--color-red)/10 border border-(--color-border) flex items-center justify-center">
                        <svg class="w-4 h-4 text-(--color-red)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                    </div>
                </div>
                <div class="min-w-0">
                    <div class="font-medium text-(--color-text-primary) truncate text-sm">
                        {{ storeLabel }}
                    </div>
                    <div v-if="lotLabel" class="text-xs text-(--color-text-secondary) truncate mt-0.5">
                        {{ lotLabel }}
                    </div>
                </div>
            </div>
        </div>

        <!-- Reason - 3 columns -->
        <div class="col-span-3 min-w-0 pr-3">
            <div class="text-sm text-(--color-text-primary) truncate">{{ damage.reason }}</div>
        </div>

        <!-- Quantity / Weight - 2 columns -->
        <div class="col-span-2 min-w-0 pr-3">
            <span class="text-xs text-(--color-text-secondary) truncate block">
                Qty: {{ formatNumber(damage.quantity) }} {{ damage.quantity_unit }}
            </span>
            <span class="text-xs text-(--color-text-secondary) truncate block mt-0.5">
                Wt: {{ formatNumber(damage.weight) }} {{ damage.weight_unit }}
            </span>
        </div>

        <!-- Amount - 1 column -->
        <div class="col-span-1 min-w-0 pr-3">
            <span class="text-sm font-semibold text-(--color-red) truncate block">
                {{ formatCurrency(damage.amount) }}
            </span>
        </div>

        <!-- Date - 1 column -->
        <div class="col-span-1 min-w-0 pr-3">
            <span class="text-xs text-(--color-text-secondary) truncate block">
                {{ formatDateShort(damage.damage_date) }}
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
import type { Damage } from '@/types/damage'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    damage: Damage
}>()

const emit = defineEmits<{
    'view': [damage: Damage]
    'edit': [damage: Damage]
    'delete': [damage: Damage]
    'updated': []
}>()

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const isOpen = ref(false)

const store = computed(() => storesStore.getStoreById(props.damage.store_id))

const storeLabel = computed(() =>
    store.value ? `Store #${store.value.id}` : `Store #${props.damage.store_id}`
)

const lot = computed(() => {
    if (!store.value) return null
    return lotsStore.getLotById(store.value.lot_id) ?? null
})

const lotLabel = computed(() => {
    if (!lot.value) return ''
    return `${lot.value.product_name} · Lot ${lot.value.lot_number}`
})

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)

const formatDateShort = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.damage) }
const handleEdit = () => { closeMenu(); emit('edit', props.damage) }
const handleDelete = () => { closeMenu(); emit('delete', props.damage) }
</script>