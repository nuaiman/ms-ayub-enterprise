<!-- src/components/features/deliveryItems/DeliveryItemDetail.vue -->
<template>
    <div v-if="deliveryItem" class="space-y-6">
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
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ lotDisplayName }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">Lot #{{ lotNumber }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDateShort(deliveryItem.created_at)
                        }}</span>
                    <span v-if="majhiName !== '৳'" class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span v-if="majhiName !== '৳'" class="text-sm text-(--color-text-secondary)">{{ majhiName
                        }}</span>
                </div>
            </div>
        </div>

        <!-- Summary strip -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ deliveryItem.quantity }}</p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5">{{ deliveryItem.quantity_unit }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">{{ deliveryItem.weight }}</p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5">{{ deliveryItem.weight_unit }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Loading Bill</p>
                <p class="text-lg font-bold text-(--color-blue) mt-1">{{ formatCurrency(loadingBillAmount) }}</p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5 capitalize">{{ deliveryItem.customer_charge_type
                    }} ৳— {{ formatCurrency(deliveryItem.loading_rate) }}</p>
            </div>
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Paid</p>
                <p class="text-lg font-bold text-(--color-green) mt-1">{{
                    formatCurrency(deliveryItem.customer_paid_unload_amount) }}</p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5">toward loading bill</p>
            </div>
        </div>

        <!-- Delivery item information -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Delivery Item
                    Information</h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Product</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lotDisplayName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot</p>
                    <p class="text-sm text-(--color-text-primary)">#{{ lotNumber }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Godown</p>
                    <p class="text-sm text-(--color-text-primary)">{{ godownName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                    <p class="text-sm text-(--color-text-primary)">{{ majhiName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Vehicle
                        Number</p>
                    <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.vehicle_number || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Driver
                        Number</p>
                    <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.driver_number || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Quantity</p>
                    <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.quantity }} {{
                        deliveryItem.quantity_unit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Weight</p>
                    <p class="text-sm text-(--color-text-primary)">{{ deliveryItem.weight }} {{
                        deliveryItem.weight_unit }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(deliveryItem.created_at) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Updated</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(deliveryItem.updated_at) }}</p>
                </div>
                <div v-if="deliveryItem.notes" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) mt-1">
                        <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ deliveryItem.notes }}
                        </p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Billing snapshot -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Billing Snapshot
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-3">
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Customer
                        Loading</p>
                    <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Charge Type</p>
                            <p class="font-semibold text-(--color-text-primary) capitalize">{{
                                deliveryItem.customer_charge_type }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Loading Rate</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(deliveryItem.loading_rate) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Bill Amount</p>
                            <p class="font-semibold text-(--color-blue)">{{ formatCurrency(loadingBillAmount) }}</p>
                        </div>
                    </div>
                    <div class="mt-3 pt-3 border-t border-(--color-border) grid grid-cols-2 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Paid</p>
                            <p class="font-semibold text-(--color-green)">{{
                                formatCurrency(deliveryItem.customer_paid_unload_amount) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Due</p>
                            <p class="font-semibold"
                                :class="customerDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                                {{ formatCurrency(customerDue) }}</p>
                        </div>
                    </div>
                </div>

                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Majhi Cut
                    </p>
                    <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Bill Type</p>
                            <p class="font-semibold text-(--color-text-primary) capitalize">{{
                                deliveryItem.majhi_bill_type }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Majhi Cut</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(deliveryItem.majhi_cut) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Bill Amount</p>
                            <p class="font-semibold text-(--color-blue)">{{ formatCurrency(majhiBillAmount) }}</p>
                        </div>
                    </div>
                    <div class="mt-3 pt-3 border-t border-(--color-border) grid grid-cols-2 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Paid</p>
                            <p class="font-semibold text-(--color-green)">{{
                                formatCurrency(deliveryItem.majhi_total_paid) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Due</p>
                            <p class="font-semibold"
                                :class="majhiDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                                {{ formatCurrency(majhiDue) }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', deliveryItem)"
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
import type { DeliveryItem } from '@/types/deliveryItem'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useGodownsStore } from '@/stores/godowns'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    deliveryItem: DeliveryItem | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [deliveryItem: DeliveryItem]
    'updated': []
}>()

const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const godownsStore = useGodownsStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()

const lot = computed(() => {
    if (!props.deliveryItem) return null
    return lotsStore.getLotById(props.deliveryItem.lot_id) || null
})

const lotNumber = computed(() => lot.value?.lot_number ?? '৳')

const lotDisplayName = computed(() => {
    if (!lot.value) return '৳'
    return lotsStore.getLotDisplayName(lot.value)
})

const store = computed(() => {
    if (!props.deliveryItem) return null
    return storesStore.getStoreById(props.deliveryItem.store_id) || null
})

const godownName = computed(() => {
    if (!store.value) return '৳'
    return godownsStore.getGodownName(store.value.godown_id)
})

const customerName = computed(() => {
    if (!lot.value?.customer_id) return '৳'
    return customersStore.getCustomerName(lot.value.customer_id)
})

const majhiName = computed(() => {
    if (!props.deliveryItem?.majhi_id) return '৳'
    return majhisStore.getMajhiName(props.deliveryItem.majhi_id)
})

const loadingBillAmount = computed(() => {
    if (!props.deliveryItem) return 0
    const rate = props.deliveryItem.loading_rate || 0
    if (props.deliveryItem.customer_charge_type === 'quantity') {
        return rate * (props.deliveryItem.quantity || 0)
    }
    return rate * (props.deliveryItem.weight || 0)
})

const majhiBillAmount = computed(() => {
    if (!props.deliveryItem) return 0
    const cut = props.deliveryItem.majhi_cut || 0
    switch (props.deliveryItem.majhi_bill_type) {
        case 'quantity':
            return cut * (props.deliveryItem.quantity || 0)
        case 'weight':
            return cut * (props.deliveryItem.weight || 0)
        case 'job':
            return cut
        default:
            return 0
    }
})

const customerDue = computed(() => {
    if (!props.deliveryItem) return 0
    return Math.max(0, loadingBillAmount.value - (props.deliveryItem.customer_paid_unload_amount || 0))
})

const majhiDue = computed(() => {
    if (!props.deliveryItem) return 0
    return Math.max(0, majhiBillAmount.value - (props.deliveryItem.majhi_total_paid || 0))
})

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