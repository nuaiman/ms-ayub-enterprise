<!-- src/views/MajhiLoadingBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ majhiLoadingBillsStore.totalBills
                        }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(majhiLoadingBillsStore.totalAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ majhiLoadingBillsStore.totalUnpaid }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(majhiLoadingBillsStore.totalUnpaidAmount) }}</p>
                </div>
            </div>

            <!-- Majhi Loading Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <MajhiLoadingBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMajhiLoadingBillsStore } from '@/stores/majhiLoadingBills'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useItemsStore } from '@/stores/items'
import { useMajhisStore } from '@/stores/majhis'
import { useLotsStore } from '@/stores/lots'
import AppLayout from '@/components/layouts/AppLayout.vue'
import MajhiLoadingBillList from '@/components/features/majhiLoadingBills/MajhiLoadingBillList.vue'
import { formatCurrency } from '@/utils/currency'

const majhiLoadingBillsStore = useMajhiLoadingBillsStore()
const deliveryItemsStore = useDeliveryItemsStore()
const deliveriesStore = useDeliveriesStore()
const itemsStore = useItemsStore()
const majhisStore = useMajhisStore()
const lotsStore = useLotsStore()

onMounted(async () => {
    await Promise.all([
        deliveryItemsStore.fetchDeliveryItems(),
        deliveriesStore.fetchDeliveries(),
        itemsStore.fetchItems(),
        majhisStore.fetchMajhis(),
        lotsStore.fetchLots(),
    ])
})
</script>