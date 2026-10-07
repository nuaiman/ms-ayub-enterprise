<!-- src/components/features/stores/StoreDetail.vue -->
<template>
    <div v-if="store" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div v-if="store.image_url"
                    class="w-20 h-20 rounded-full overflow-hidden border-2 border-(--color-border)">
                    <img :src="getImageUrl(store.image_url)" :alt="lotDisplayName" class="w-full h-full object-cover" />
                </div>
                <div v-else
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Store #{{ store.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ lotDisplayName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">Lot {{ lotNumber }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="store.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                        <span class="w-1.5 h-1.5 rounded-full"
                            :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                        {{ store.is_active ? 'Active' : 'Inactive' }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Summary strip -->
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1 leading-tight">
                    {{ formatNumber(store.quantity) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ quantityUnit }}</span>
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1 leading-tight">
                    {{ formatNumber(store.weight) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ weightUnit }}</span>
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Start Date</p>
                <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                    {{ formatDateShort(store.start_date) }}
                </p>
            </div>
        </div>

        <!-- Store info -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Store
                    Information</h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lotDisplayName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot Number
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ lotNumber }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
                    <p class="text-sm text-(--color-text-primary)">{{ godownName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatNumber(store.quantity) }} {{ quantityUnit }}
                    </p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatNumber(store.weight) }} {{ weightUnit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Start Date
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDateShort(store.start_date) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Status</p>
                    <p class="text-sm text-(--color-text-primary)">{{ store.is_active ? 'Active' : 'Inactive' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(store.created_at) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Updated</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(store.updated_at) }}</p>
                </div>
            </div>
        </section>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', store)"
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
import type { Store } from '@/types/store'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useGodownsStore } from '@/stores/godowns'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    store: Store | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [store: Store]
    'updated': []
}>()

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const godownsStore = useGodownsStore()

const lot = computed(() => {
    if (!props.store) return null
    return lotsStore.getLotById(props.store.lot_id) || null
})

const lotNumber = computed(() => lot.value?.lot_number ?? '—')
const lotDisplayName = computed(() => lot.value?.product_name ?? '—')
const quantityUnit = computed(() => lot.value?.quantity_unit ?? 'units')
const weightUnit = computed(() => lot.value?.weight_unit ?? 'kg')

const customerName = computed(() => {
    if (!lot.value) return '—'
    return customersStore.getCustomerName(lot.value.customer_id)
})

const godownName = computed(() => {
    if (!props.store) return '—'
    return godownsStore.getGodownName(props.store.godown_id)
})

const formatNumber = (n: number): string => {
    return new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
}

const formatDateShort = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })
}
</script>