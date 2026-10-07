<!-- src/components/features/deliveries/DeliveryDetail.vue -->
<template>
    <div v-if="delivery" class="space-y-6">
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 17a2 2 0 11-4 0 2 2 0 014 0zM20 17a2 2 0 11-4 0 2 2 0 014 0zM13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">Delivery #{{ delivery.id }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">{{ customerName }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ formatDate(delivery.delivery_date) }}</span>
                </div>
            </div>
        </div>

        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ delivery.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(delivery.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(delivery.updated_at) }}</span>
        </div>

        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">
                    Delivery Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Delivery Date
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(delivery.delivery_date) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Receiver Name
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.receiver_name || '—' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Receiver Phone
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.receiver_phone || '—' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">From</p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.from_location || '—' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">To</p>
                    <p class="text-sm text-(--color-text-primary)">{{ delivery.to_location || '—' }}</p>
                </div>
            </div>
        </section>

        <section v-if="delivery.notes" class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Notes</h3>
            </div>
            <div class="p-4">
                <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ delivery.notes }}</p>
            </div>
        </section>

        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Items</h3>
                <span class="text-xs text-(--color-text-secondary)">{{ items.length }}</span>
            </div>

            <div v-if="loading" class="p-4 text-sm text-(--color-text-secondary)">Loading items…</div>
            <div v-else-if="items.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No items for this delivery.
            </div>
            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="it in items" :key="it.id" class="p-4">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div class="min-w-0 flex-1">
                            <div class="text-sm font-semibold text-(--color-text-primary)">
                                Store #{{ it.store_id }}
                            </div>
                            <div v-if="lotLabelFor(it)" class="text-xs text-(--color-text-secondary) mt-0.5">
                                {{ lotLabelFor(it) }}
                            </div>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-xs text-(--color-text-secondary)">
                                <span>Qty: <span class="text-(--color-text-primary)">{{ formatNumber(it.quantity)
                                        }}</span></span>
                                <span>Wt: <span class="text-(--color-text-primary)">{{ formatNumber(it.weight)
                                        }}</span></span>
                                <span v-if="it.vehicle_number">Vehicle: <span class="text-(--color-text-primary)">{{
                                        it.vehicle_number }}</span></span>
                                <span v-if="it.driver_number">Driver: <span class="text-(--color-text-primary)">{{
                                        it.driver_number }}</span></span>
                                <span>Majhi: <span class="text-(--color-text-primary)">{{
                                        majhisStore.getMajhiName(it.majhi_id) }}</span></span>
                            </div>

                            <div class="mt-2 flex flex-wrap gap-2">
                                <span v-if="customerBillFor(it.id)"
                                    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md text-xs font-medium bg-(--color-blue)/10 text-(--color-blue)">
                                    Customer Bill: {{ formatCurrency(customerBillFor(it.id)!.total_amount) }}
                                </span>
                                <span v-if="majhiBillFor(it.id)"
                                    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md text-xs font-medium bg-(--color-yellow)/10 text-(--color-yellow)">
                                    Majhi Bill: {{ formatCurrency(majhiBillFor(it.id)!.total_amount) }}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', delivery)"
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
import { ref, computed, onMounted, watch } from 'vue'
import type { Delivery, DeliveryItem } from '@/types/delivery'
import type { CustomerDeliveryBill } from '@/types/customerDeliveryBill'
import type { MajhiBill } from '@/types/majhiBill'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useMajhiBillsStore } from '@/stores/majhiBills'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    delivery: Delivery | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [delivery: Delivery]
    'updated': []
}>()

const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
const majhiBillsStore = useMajhiBillsStore()

const items = ref<DeliveryItem[]>([])
const loading = ref(false)
const customerBillsByItem = ref<Record<number, CustomerDeliveryBill>>({})
const majhiBillsByItem = ref<Record<number, MajhiBill>>({})

const customerName = computed(() => {
    if (!props.delivery?.customer_id) return '—'
    return customersStore.getCustomerName(props.delivery.customer_id)
})

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric', month: 'short', day: 'numeric',
    })

const formatNumber = (n: number): string =>
    new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(n)

const lotLabelFor = (it: DeliveryItem): string => {
    const store = storesStore.getStoreById(it.store_id)
    if (!store) return ''
    const lot = lotsStore.getLotById(store.lot_id)
    if (!lot) return ''
    return `${lot.product_name} · Lot ${lot.lot_number}`
}

const customerBillFor = (itemId: number) => customerBillsByItem.value[itemId] ?? null
const majhiBillFor = (itemId: number) => majhiBillsByItem.value[itemId] ?? null

const loadItems = async () => {
    if (!props.delivery) {
        items.value = []
        return
    }
    loading.value = true
    try {
        items.value = await deliveriesStore.fetchDeliveryItems(props.delivery.id)

        const cdMap: Record<number, CustomerDeliveryBill> = {}
        const mbMap: Record<number, MajhiBill> = {}

        await Promise.all(items.value.map(async it => {
            const cd = await customerDeliveryBillsStore.fetchCustomerDeliveryBillByItemId(it.id)
            if (cd) cdMap[it.id] = cd
            const mb = await majhiBillsStore.fetchMajhiBillsByDeliveryItemId(it.id)
            if (mb[0]) mbMap[it.id] = mb[0]
        }))

        customerBillsByItem.value = cdMap
        majhiBillsByItem.value = mbMap
    } finally {
        loading.value = false
    }
}

watch(() => props.delivery?.id, loadItems, { immediate: true })

onMounted(async () => {
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
    if (storesStore.stores.length === 0) await storesStore.fetchStores()
    if (lotsStore.lots.length === 0) await lotsStore.fetchLots()
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
})
</script>