<!-- src/components/features/damages/DamageDetail.vue -->
<template>
    <div v-if="damage" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div v-if="damage.image_url"
                    class="w-20 h-20 rounded-full overflow-hidden border-2 border-(--color-border)">
                    <img :src="getImageUrl(damage.image_url)" :alt="damage.reason" class="w-full h-full object-cover" />
                </div>
                <div v-else
                    class="w-20 h-20 rounded-full bg-(--color-red)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-red)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Damage #{{ damage.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ storeLabel }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ damage.reason }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-red)">{{ formatCurrency(damage.amount) }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ damage.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Date: {{ formatDateShort(damage.damage_date) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">By: {{ userName }}</span>
        </div>

        <!-- Summary cards -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">
                    {{ formatNumber(damage.quantity) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ damage.quantity_unit }}</span>
                </p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">
                    {{ formatNumber(damage.weight) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ damage.weight_unit }}</span>
                </p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                <p class="text-lg font-bold text-(--color-red) mt-1">{{ formatCurrency(damage.amount) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Source</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1">Store</p>
            </div>
        </div>

        <!-- Details -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">
                    Damage Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                    <p class="text-sm text-(--color-text-primary)">{{ storeLabel }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Reason</p>
                    <p class="text-sm text-(--color-text-primary)">{{ damage.reason }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Damage Date
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDateShort(damage.damage_date) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Recorded By
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ userName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDateTime(damage.created_at) }}</p>
                </div>
            </div>
        </section>

        <!-- Notes -->
        <div v-if="damage.notes" class="border-t border-(--color-border) pt-4">
            <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">Notes</p>
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ damage.notes }}</p>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', damage)"
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
import { computed } from 'vue'
import type { Damage } from '@/types/damage'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    damage: Damage | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [damage: Damage]
    'updated': []
}>()

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const usersStore = useUsersStore()

const storeLabel = computed(() => {
    if (!props.damage) return '—'
    const store = storesStore.getStoreById(props.damage.store_id)
    if (!store) return `Store #${props.damage.store_id}`
    const lot = lotsStore.getLotById(store.lot_id)
    if (!lot) return `Store #${store.id}`
    return `Store #${store.id} — ${lot.product_name} (Lot ${lot.lot_number})`
})

const userName = computed(() => {
    if (!props.damage) return '—'
    return usersStore.getUserName(props.damage.user_id)
})

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)

const formatDateShort = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })

const formatDateTime = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
</script>