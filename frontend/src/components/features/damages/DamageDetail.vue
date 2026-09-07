<!-- src/components/features/damages/DamageDetail.vue -->
<template>
    <div v-if="damage" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-red)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-red)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ damage.reason }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ getStoreDisplayName(store) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-red)">{{ formatCurrency(damage.amount) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDate(damage.damage_date) }}</span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ damage.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDateTime(damage.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDateTime(damage.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Store -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                <p class="text-sm text-(--color-text-primary)">{{ getStoreDisplayName(store) }}</p>
            </div>

            <!-- Reason -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Reason</p>
                <p class="text-sm text-(--color-text-primary)">{{ damage.reason }}</p>
            </div>

            <!-- Quantity -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                <p class="text-sm text-(--color-text-primary)">{{ damage.quantity }} {{ damage.quantity_unit }}</p>
            </div>

            <!-- Weight -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                <p class="text-sm text-(--color-text-primary)">{{ damage.weight }} {{ damage.weight_unit }}</p>
            </div>

            <!-- Amount -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                <p class="text-lg font-semibold text-(--color-red)">{{ formatCurrency(damage.amount) }}</p>
            </div>

            <!-- Damage Date -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Damage Date</p>
                <p class="text-sm text-(--color-text-primary)">{{ formatDate(damage.damage_date) }}</p>
            </div>

            <!-- User -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Recorded By</p>
                <p class="text-sm text-(--color-text-primary)">{{ getUserName(damage.user_id) }}</p>
            </div>

            <!-- Notes -->
            <div v-if="damage.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ damage.notes }}</p>
                </div>
            </div>

            <!-- Image -->
            <div v-if="damage.image_url" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Attachment</p>
                <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md">
                    <img :src="getImageUrl(damage.image_url)" alt="Damage attachment"
                        class="w-full object-cover max-h-64" />
                </div>
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
import type { Store } from '@/types/store'
import { useStoresStore } from '@/stores/stores'
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
const usersStore = useUsersStore()

const store = computed(() => {
    if (!props.damage) return undefined
    return storesStore.getStoreById(props.damage.store_id)
})

const getStoreDisplayName = (store: Store | undefined): string => {
    if (!store) return `Store #${props.damage?.store_id || 'Unknown'}`
    return storesStore.getStoreDisplayName(store)
}

const getUserName = (userId: number): string => {
    return usersStore.getUserName(userId)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'long',
        day: 'numeric',
        year: 'numeric'
    })
}

const formatDateTime = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}
</script>