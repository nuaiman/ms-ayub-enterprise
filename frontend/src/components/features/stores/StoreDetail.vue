<!-- src/components/features/stores/StoreDetail.vue -->
<template>
    <div v-if="store" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ godownName }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ lotDisplayName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">Lot #{{ lotNumber }}</span>
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
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Stock on Hand</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1 leading-tight">
                    {{ formatNumber(store.quantity) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ store.quantity_unit }}</span>
                </p>
                <p class="text-lg font-bold text-(--color-text-primary) leading-tight">
                    {{ formatNumber(store.weight) }}
                    <span class="text-xs font-medium text-(--color-text-secondary)">{{ store.weight_unit }}</span>
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Monthly Godown Bill</p>
                <p class="text-lg font-bold text-(--color-blue) mt-1">{{ formatCurrency(monthlyBill) }}</p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5 capitalize">{{ store.store_bill_type }} ৳— {{
                    formatCurrency(store.godown_cut) }}</p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Last Paid</p>
                <p class="text-lg font-bold text-(--color-green) mt-1">{{ formatCurrency(store.last_paid_amount) }}</p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5">
                    {{ store.last_paid_through ? `Through ${formatDateShort(store.last_paid_through)}` : 'Never paid' }}
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                <p class="text-lg font-bold mt-1"
                    :class="outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(outstanding) }}
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
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
                    <p class="text-sm text-(--color-text-primary)">{{ godownName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lotDisplayName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot Number
                    </p>
                    <p class="text-sm text-(--color-text-primary)">#{{ lotNumber }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Bill Type</p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ store.store_bill_type }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown Cut
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(store.godown_cut) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Billing Start
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDateShort(store.billing_start) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Billing End
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ store.billing_end ?
                        formatDateShort(store.billing_end) : '৳' }}</p>
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
                <div v-if="store.notes" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) mt-1">
                        <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ store.notes }}</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Godown store bill -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Godown Store
                    Bill</h3>
            </div>
            <div class="p-4 grid grid-cols-3 gap-3">
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Billed</p>
                    <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ formatCurrency(monthlyBill) }}</p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Paid</p>
                    <p class="text-lg font-bold text-(--color-green) mt-1">{{ formatCurrency(store.last_paid_amount)
                        }}</p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Due</p>
                    <p class="text-lg font-bold mt-1"
                        :class="outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                        {{ formatCurrency(outstanding) }}
                    </p>
                </div>
            </div>
        </section>

        <!-- Deliveries -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Deliveries</h3>
                <span class="text-xs text-(--color-text-secondary)">{{ deliveries.length }}</span>
            </div>

            <div v-if="deliveries.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No deliveries from this store yet.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="di in deliveries" :key="di.id" class="p-4">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div>
                            <p class="text-sm font-semibold text-(--color-text-primary)">
                                Delivery ৳ {{ formatDate(di.created_at) }}
                            </p>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-1 text-xs text-(--color-text-secondary)">
                                <span>Quantity: <span class="text-(--color-text-primary)">{{ di.quantity }} {{
                                    di.quantity_unit }}</span></span>
                                <span>Weight: <span class="text-(--color-text-primary)">{{ di.weight }} {{
                                    di.weight_unit }}</span></span>
                                <span v-if="di.vehicle_number">Vehicle: <span class="text-(--color-text-primary)">{{
                                    di.vehicle_number }}</span></span>
                                <span v-if="di.driver_number">Driver: <span class="text-(--color-text-primary)">{{
                                    di.driver_number }}</span></span>
                                <span v-if="di.majhi_id">Majhi: <span class="text-(--color-text-primary)">{{
                                    getMajhiName(di.majhi_id) }}</span></span>
                            </div>
                            <div v-if="di.notes" class="text-xs text-(--color-text-secondary) mt-2 whitespace-pre-wrap">
                                {{ di.notes }}
                            </div>
                        </div>
                        <div class="text-right text-xs shrink-0">
                            <p class="text-(--color-text-secondary)">Loading Bill</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(deliveryBillAmount(di)) }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Damages -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Damages</h3>
                <span class="text-xs text-(--color-text-secondary)">{{ damages.length }}</span>
            </div>

            <div v-if="damages.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No damages recorded for this store.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="d in damages" :key="d.id" class="p-4">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div>
                            <p class="text-sm font-semibold text-(--color-text-primary)">{{ d.reason }}</p>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-1 text-xs text-(--color-text-secondary)">
                                <span v-if="d.quantity > 0">Quantity: <span class="text-(--color-text-primary)">{{
                                    d.quantity }} {{ d.quantity_unit }}</span></span>
                                <span v-if="d.weight > 0">Weight: <span class="text-(--color-text-primary)">{{
                                    d.weight }} {{ d.weight_unit }}</span></span>
                                <span>Date: <span class="text-(--color-text-primary)">{{
                                    formatDateShort(d.damage_date) }}</span></span>
                            </div>
                            <div v-if="d.notes" class="text-xs text-(--color-text-secondary) mt-2 whitespace-pre-wrap">
                                {{ d.notes }}
                            </div>
                        </div>
                        <div class="text-right text-xs shrink-0">
                            <p class="text-(--color-text-secondary)">Amount</p>
                            <p class="font-semibold text-(--color-red)">{{ formatCurrency(d.amount) }}</p>
                        </div>
                    </div>
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
import { useMajhisStore } from '@/stores/majhis'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDamagesStore } from '@/stores/damages'
import { formatCurrency } from '@/utils/currency'

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
const majhisStore = useMajhisStore()
const deliveryItemsStore = useDeliveryItemsStore()
const damagesStore = useDamagesStore()

const lot = computed(() => {
    if (!props.store) return null
    return lotsStore.getLotById(props.store.lot_id) || null
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
    if (!props.store) return '৳'
    return godownsStore.getGodownName(props.store.godown_id)
})

const deliveries = computed(() => {
    if (!props.store) return []
    return deliveryItemsStore.deliveryItems.filter(di => di.store_id === props.store!.id)
})

const damages = computed(() => {
    if (!props.store) return []
    return damagesStore.damages.filter(d => d.store_id === props.store!.id)
})

const monthlyBill = computed(() => {
    if (!props.store) return 0
    if (props.store.store_bill_type === 'weight') {
        return (props.store.weight || 0) * (props.store.godown_cut || 0)
    }
    return (props.store.quantity || 0) * (props.store.godown_cut || 0)
})

const outstanding = computed(() => {
    if (!props.store) return 0
    return Math.max(0, monthlyBill.value - (props.store.last_paid_amount || 0))
})

const getMajhiName = (id: number | null): string => {
    if (!id) return '৳'
    return majhisStore.getMajhiName(id)
}

const deliveryBillAmount = (di: any): number => {
    const rate = di.loading_rate || 0
    if (di.customer_charge_type === 'quantity') return rate * (di.quantity || 0)
    return rate * (di.weight || 0)
}

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