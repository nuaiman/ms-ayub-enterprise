<!-- src/components/features/damages/DamageDetail.vue -->
<template>
    <div v-if="damage" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-red)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-red)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ damage.reason }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ godownName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ lotDisplayName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">Lot #{{ lotNumber }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm font-semibold text-(--color-red)">{{ formatCurrency(damage.amount) }}</span>
                </div>
            </div>
        </div>

        <!-- Damage Information -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Damage
                    Information</h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Reason</p>
                    <p class="text-sm text-(--color-text-primary)">{{ damage.reason }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Amount</p>
                    <p class="text-sm font-semibold text-(--color-red)">{{ formatCurrency(damage.amount) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                    <p class="text-sm text-(--color-text-primary)">{{ damage.quantity }} {{ damage.quantity_unit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                    <p class="text-sm text-(--color-text-primary)">{{ damage.weight }} {{ damage.weight_unit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Damage Date
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(damage.damage_date) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Recorded By
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ recordedBy }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(damage.created_at) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Updated</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(damage.updated_at) }}</p>
                </div>
                <div v-if="damage.notes" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) mt-1">
                        <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ damage.notes }}</p>
                    </div>
                </div>
                <div v-if="damage.image_url" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Attachment
                    </p>
                    <div class="rounded-lg overflow-hidden border border-(--color-border) max-w-md mt-1">
                        <img :src="getImageUrl(damage.image_url)" alt="Damage attachment"
                            class="w-full object-cover max-h-64" />
                    </div>
                </div>
            </div>
        </section>

        <!-- Related context -->
        <section v-if="store" class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Related
                    Context</h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Product</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lotDisplayName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot Number
                    </p>
                    <p class="text-sm text-(--color-text-primary)">#{{ lotNumber }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
                    <p class="text-sm text-(--color-text-primary)">{{ godownName }}</p>
                </div>
                <div v-if="lot">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Charge Type
                    </p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ lot.customer_charge_type }}</p>
                </div>
                <div v-if="lot">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Storage Rate
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.customer_storage_rate) }}</p>
                </div>
            </div>
        </section>

        <!-- Store snapshot -->
        <section v-if="store" class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Store Snapshot
                </h3>
            </div>
            <div class="p-4 grid grid-cols-2 md:grid-cols-4 gap-3">
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                    <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ store.quantity }} {{
                        store.quantity_unit }}</p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                    <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ store.weight }} {{
                        store.weight_unit }}</p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                    <p class="text-lg font-bold text-(--color-text-primary) mt-1 capitalize">{{ store.store_bill_type
                    }}</p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Status</p>
                    <p class="text-lg font-bold mt-1"
                        :class="store.is_active ? 'text-(--color-green)' : 'text-(--color-red)'">
                        {{ store.is_active ? 'Active' : 'Inactive' }}
                    </p>
                </div>
            </div>
        </section>

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
import { useCustomersStore } from '@/stores/customers'
import { useGodownsStore } from '@/stores/godowns'
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
const customersStore = useCustomersStore()
const godownsStore = useGodownsStore()
const usersStore = useUsersStore()

const store = computed(() => {
    if (!props.damage) return null
    return storesStore.getStoreById(props.damage.store_id) || null
})

const lot = computed(() => {
    if (!store.value) return null
    return lotsStore.getLotById(store.value.lot_id) || null
})

const lotNumber = computed(() => lot.value?.lot_number ?? '৳')

const lotDisplayName = computed(() => {
    if (!lot.value) return '৳'
    return lotsStore.getLotDisplayName(lot.value)
})

const customerName = computed(() => {
    if (!lot.value?.customer_id) return '৳'
    return customersStore.getCustomerName(lot.value.customer_id)
})

const godownName = computed(() => {
    if (!store.value) return '৳'
    return godownsStore.getGodownName(store.value.godown_id)
})

const recordedBy = computed(() => {
    if (!props.damage) return '৳'
    return usersStore.getUserName(props.damage.user_id)
})

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    })
}
</script>