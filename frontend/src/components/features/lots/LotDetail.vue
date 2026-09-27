<!-- src/components/features/lots/LotDetail.vue -->
<template>
    <div v-if="lot" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center overflow-hidden">
                    <img v-if="lot.image_url" :src="getImageUrl(lot.image_url)" alt="lot"
                        class="w-full h-full object-cover" />
                    <svg v-else class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor"
                        viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ lotDisplayName }}</h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span class="text-sm text-(--color-text-secondary)">Lot #{{ lot.lot_number }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary) capitalize">{{ lot.customer_charge_type
                        }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="lot.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                        <span class="w-1.5 h-1.5 rounded-full"
                            :class="lot.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                        {{ lot.is_active ? 'Active' : 'Inactive' }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Summary strip -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Stock on Hand</p>
                <div class="mt-1">
                    <p v-if="stockTotals.length === 0" class="text-lg font-bold text-(--color-text-primary)">৳</p>
                    <p v-else v-for="(line, idx) in stockTotals" :key="idx"
                        class="text-lg font-bold text-(--color-text-primary) leading-tight">
                        {{ line }}
                    </p>
                </div>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Stores</p>
                <p class="text-lg font-bold text-(--color-text-primary) mt-1">
                    {{ activeStoresCount }} <span class="text-xs font-medium text-(--color-text-secondary)">/
                        {{ stores.length }} active</span>
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customer Due</p>
                <p class="text-lg font-bold mt-1"
                    :class="billingSummary.customerDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(billingSummary.customerDue) }}
                </p>
            </div>

            <div class="rounded-xl p-3 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Majhi Due</p>
                <p class="text-lg font-bold mt-1"
                    :class="billingSummary.majhiDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(billingSummary.majhiDue) }}
                </p>
            </div>
        </div>

        <!-- Lot info -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Lot Information
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                    <p class="text-sm text-(--color-text-primary)">{{ customerName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Lot Number
                    </p>
                    <p class="text-sm text-(--color-text-primary)">#{{ lot.lot_number }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Product Name
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.product_name || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Category</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.category || '৳' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer
                        Charge Type</p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ lot.customer_charge_type }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi Bill
                        Type</p>
                    <p class="text-sm text-(--color-text-primary) capitalize">{{ lot.majhi_bill_type }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Storage Rate
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.customer_storage_rate) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Unload Rate
                    </p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.unload_rate) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi</p>
                    <p class="text-sm text-(--color-text-primary)">{{ majhiName }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Majhi Cut</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatCurrency(lot.majhi_cut) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Status</p>
                    <p class="text-sm text-(--color-text-primary)">{{ lot.is_active ? 'Active' : 'Inactive' }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Created</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(lot.created_at) }}</p>
                </div>
                <div>
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Updated</p>
                    <p class="text-sm text-(--color-text-primary)">{{ formatDate(lot.updated_at) }}</p>
                </div>
                <div v-if="lot.notes" class="md:col-span-2">
                    <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                    <div class="p-3 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) mt-1">
                        <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ lot.notes }}</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Payment tracking -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Payment Tracking
                </h3>
            </div>
            <div class="p-4 grid grid-cols-1 md:grid-cols-3 gap-3">
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Storage Paid
                    </p>
                    <p class="text-lg font-bold text-(--color-green) mt-1">{{
                        formatCurrency(lot.customer_last_paid_amount) }}</p>
                    <p class="text-xs text-(--color-text-secondary) mt-1">
                        {{ lot.customer_last_paid_through ? `Through ${formatDateShort(lot.customer_last_paid_through)}`
                            : 'Not paid yet' }}
                    </p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Unload Paid
                    </p>
                    <p class="text-lg font-bold text-(--color-green) mt-1">{{
                        formatCurrency(lot.customer_paid_unload_amount) }}</p>
                    <p class="text-xs text-(--color-text-secondary) mt-1">Total paid toward unload bill</p>
                </div>
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Majhi Paid
                    </p>
                    <p class="text-lg font-bold text-(--color-blue) mt-1">{{ formatCurrency(lot.majhi_total_paid) }}</p>
                    <p class="text-xs text-(--color-text-secondary) mt-1">Total paid to majhi</p>
                </div>
            </div>
        </section>

        <!-- Billing summary -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border)">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Billing Summary
                </h3>
            </div>

            <div class="p-4 grid grid-cols-1 md:grid-cols-2 gap-3">
                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Customer
                        Storage</p>
                    <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Billed</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(billingSummary.storageBilled) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Paid</p>
                            <p class="font-semibold text-(--color-green)">{{
                                formatCurrency(billingSummary.storagePaid) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Due</p>
                            <p class="font-semibold"
                                :class="billingSummary.storageDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                                {{ formatCurrency(billingSummary.storageDue) }}</p>
                        </div>
                    </div>
                </div>

                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Customer
                        Unload</p>
                    <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Billed</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(billingSummary.unloadBilled) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Paid</p>
                            <p class="font-semibold text-(--color-green)">{{
                                formatCurrency(billingSummary.unloadPaid) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Due</p>
                            <p class="font-semibold"
                                :class="billingSummary.unloadDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                                {{ formatCurrency(billingSummary.unloadDue) }}</p>
                        </div>
                    </div>
                </div>

                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Customer
                        Delivery</p>
                    <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Billed</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(billingSummary.deliveryBilled) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Paid</p>
                            <p class="font-semibold text-(--color-green)">{{
                                formatCurrency(billingSummary.deliveryPaid) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Due</p>
                            <p class="font-semibold"
                                :class="billingSummary.deliveryDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                                {{ formatCurrency(billingSummary.deliveryDue) }}</p>
                        </div>
                    </div>
                </div>

                <div class="p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/20">
                    <p class="text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">Majhi Cut
                    </p>
                    <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
                        <div>
                            <p class="text-(--color-text-secondary)">Billed</p>
                            <p class="font-semibold text-(--color-text-primary)">{{
                                formatCurrency(billingSummary.majhiBilled) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Paid</p>
                            <p class="font-semibold text-(--color-green)">{{
                                formatCurrency(billingSummary.majhiPaid) }}</p>
                        </div>
                        <div>
                            <p class="text-(--color-text-secondary)">Due</p>
                            <p class="font-semibold"
                                :class="billingSummary.majhiDue > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                                {{ formatCurrency(billingSummary.majhiDue) }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Stores -->
        <section class="rounded-xl border border-(--color-border) bg-(--color-surface)">
            <div class="px-4 py-3 border-b border-(--color-border) flex items-center justify-between">
                <h3 class="text-sm font-semibold text-(--color-text-primary) uppercase tracking-wider">Stores</h3>
                <span class="text-xs text-(--color-text-secondary)">{{ stores.length }}</span>
            </div>

            <div v-if="stores.length === 0" class="p-4 text-sm text-(--color-text-secondary)">
                No stores for this lot.
            </div>

            <div v-else class="divide-y divide-(--color-border)">
                <div v-for="store in stores" :key="store.id" class="p-4">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div class="min-w-0">
                            <div class="flex items-center gap-2 flex-wrap">
                                <span class="text-sm font-semibold text-(--color-text-primary)">{{
                                    getGodownName(store.godown_id) }}</span>
                                <span
                                    class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                                    :class="store.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                                    <span class="w-1.5 h-1.5 rounded-full"
                                        :class="store.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                                    {{ store.is_active ? 'Active' : 'Inactive' }}
                                </span>
                            </div>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-xs text-(--color-text-secondary)">
                                <span>Quantity: <span class="text-(--color-text-primary)">{{ store.quantity }} {{
                                    store.quantity_unit }}</span></span>
                                <span>Weight: <span class="text-(--color-text-primary)">{{ store.weight }} {{
                                    store.weight_unit }}</span></span>
                                <span>Bill Type: <span class="text-(--color-text-primary) capitalize">{{
                                    store.store_bill_type }}</span></span>
                                <span>Godown Cut: <span class="text-(--color-text-primary)">{{
                                    formatCurrency(store.godown_cut) }}</span></span>
                            </div>
                            <div class="flex flex-wrap gap-x-4 gap-y-1 mt-1 text-xs text-(--color-text-secondary)">
                                <span v-if="store.billing_start">Billing: <span class="text-(--color-text-primary)">{{
                                    formatDateShort(store.billing_start) }}</span></span>
                                <span v-if="store.billing_end">৳†’ <span class="text-(--color-text-primary)">{{
                                    formatDateShort(store.billing_end) }}</span></span>
                            </div>
                        </div>

                        <div class="text-right text-xs space-y-1 shrink-0">
                            <div v-if="(store.last_paid_amount || 0) > 0">
                                <p class="text-(--color-text-secondary)">Last Paid</p>
                                <p class="font-semibold text-(--color-green)">{{ formatCurrency(store.last_paid_amount)
                                    }}</p>
                                <p v-if="store.last_paid_through" class="text-(--color-text-secondary)">through {{
                                    formatDateShort(store.last_paid_through) }}</p>
                            </div>
                        </div>
                    </div>
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
                No deliveries for this lot.
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

                        <div class="text-right text-xs space-y-1 shrink-0">
                            <div>
                                <p class="text-(--color-text-secondary)">Loading Bill</p>
                                <p class="font-semibold text-(--color-text-primary)">{{
                                    formatCurrency(deliveryBillAmount(di)) }}</p>
                            </div>
                            <div v-if="(di.customer_paid_unload_amount || 0) > 0">
                                <p class="text-(--color-text-secondary)">Paid</p>
                                <p class="font-semibold text-(--color-green)">{{
                                    formatCurrency(di.customer_paid_unload_amount) }}</p>
                            </div>
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
                No damages recorded for this lot.
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
            <button @click="emit('edit', lot)"
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
import type { Lot } from '@/types/lot'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDamagesStore } from '@/stores/damages'
import { useCustomerStorageBillsStore } from '@/stores/customerStorageBills'
import { useCustomerLotBillsStore } from '@/stores/customerLotBills'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useMajhiLotBillsStore } from '@/stores/majhiLotBills'
import { useMajhiLoadingBillsStore } from '@/stores/majhiLoadingBills'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
    lot: Lot | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [lot: Lot]
    'updated': []
}>()

const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const deliveryItemsStore = useDeliveryItemsStore()
const damagesStore = useDamagesStore()
const customerStorageBillsStore = useCustomerStorageBillsStore()
const customerLotBillsStore = useCustomerLotBillsStore()
const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
const majhiLotBillsStore = useMajhiLotBillsStore()
const majhiLoadingBillsStore = useMajhiLoadingBillsStore()

const lotDisplayName = computed(() => {
    if (!props.lot) return '৳'
    return lotsStore.getLotDisplayName(props.lot)
})

const customerName = computed(() => {
    if (!props.lot?.customer_id) return '৳'
    return customersStore.getCustomerName(props.lot.customer_id)
})

const majhiName = computed(() => {
    if (!props.lot?.majhi_id) return '৳'
    return majhisStore.getMajhiName(props.lot.majhi_id)
})

const stores = computed(() => {
    if (!props.lot) return []
    return storesStore.getStoresByLotId(props.lot.id)
})

const storeIds = computed(() => new Set(stores.value.map(s => s.id)))

const deliveries = computed(() => {
    if (!props.lot) return []
    return deliveryItemsStore.deliveryItems.filter(di => di.lot_id === props.lot!.id)
})

const damages = computed(() => {
    if (!props.lot) return []
    return damagesStore.damages.filter(d => storeIds.value.has(d.store_id))
})

const activeStoresCount = computed(() => stores.value.filter(s => s.is_active).length)

const stockTotals = computed<string[]>(() => {
    const byUnit: Record<string, number> = {}
    for (const s of stores.value) {
        if (!s.is_active) continue
        if (s.quantity > 0) {
            const u = s.quantity_unit || 'units'
            byUnit[u] = (byUnit[u] || 0) + s.quantity
        }
        if (s.weight > 0) {
            const u = s.weight_unit || 'kg'
            byUnit[u] = (byUnit[u] || 0) + s.weight
        }
    }

    const lines: string[] = []
    for (const [unit, total] of Object.entries(byUnit)) {
        const n = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(total)
        lines.push(`${n} ${unit}`)
    }
    return lines
})

const billingSummary = computed(() => {
    const lotId = props.lot?.id

    let storageBilled = 0
    let storagePaid = 0
    const storageBills = customerStorageBillsStore.customerBillData.filter(b => b.lot_id === lotId)
    for (const b of storageBills) {
        storageBilled += b.total_billed || 0
        storagePaid += b.total_paid || 0
    }
    const storageDue = Math.max(0, storageBilled - storagePaid)

    let unloadBilled = 0
    let unloadPaid = 0
    const unloadBills = customerLotBillsStore.bills.filter(b => b.lot_id === lotId)
    for (const b of unloadBills) {
        unloadBilled += b.bill_amount || 0
        unloadPaid += b.paid_amount || 0
    }
    const unloadDue = Math.max(0, unloadBilled - unloadPaid)

    const deliveryItemIds = new Set(deliveries.value.map(d => d.id))
    let deliveryBilled = 0
    let deliveryPaid = 0
    const deliveryBills = customerDeliveryBillsStore.bills.filter(b => deliveryItemIds.has(b.delivery_item_id))
    for (const b of deliveryBills) {
        deliveryBilled += b.bill_amount || 0
        deliveryPaid += b.paid_amount || 0
    }
    const deliveryDue = Math.max(0, deliveryBilled - deliveryPaid)

    let majhiBilled = 0
    let majhiPaid = 0
    const majhiLotBills = majhiLotBillsStore.bills.filter(b => b.lot_id === lotId)
    for (const b of majhiLotBills) {
        majhiBilled += b.bill_amount || 0
        majhiPaid += b.paid_amount || 0
    }
    const majhiLoadingBills = majhiLoadingBillsStore.bills.filter(b => deliveryItemIds.has(b.delivery_item_id))
    for (const b of majhiLoadingBills) {
        majhiBilled += b.bill_amount || 0
        majhiPaid += b.paid_amount || 0
    }
    const majhiDue = Math.max(0, majhiBilled - majhiPaid)

    const customerDue = storageDue + unloadDue + deliveryDue

    return {
        storageBilled, storagePaid, storageDue,
        unloadBilled, unloadPaid, unloadDue,
        deliveryBilled, deliveryPaid, deliveryDue,
        majhiBilled, majhiPaid, majhiDue,
        customerDue,
    }
})

const getGodownName = (id: number): string => godownsStore.getGodownName(id)

const getMajhiName = (id: number | null): string => {
    if (!id) return '৳'
    return majhisStore.getMajhiName(id)
}

const deliveryBillAmount = (di: any): number => {
    const rate = di.loading_rate || 0
    if (di.customer_charge_type === 'quantity') return rate * (di.quantity || 0)
    return rate * (di.weight || 0)
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