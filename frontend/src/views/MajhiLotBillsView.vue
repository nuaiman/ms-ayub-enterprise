<!-- src/views/MajhiLotBillsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ majhiLotBillsStore.totalBills }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Amount</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{
                        formatCurrency(majhiLotBillsStore.totalAmount) }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ majhiLotBillsStore.totalUnpaid }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Unpaid Amount</p>
                    <p class="text-2xl font-bold text-(--color-red) mt-1">{{
                        formatCurrency(majhiLotBillsStore.totalUnpaidAmount) }}</p>
                </div>
            </div>

            <!-- Majhi Lot Bill List -->
            <div class="flex-1 min-h-0 mt-6">
                <MajhiLotBillList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useMajhiLotBillsStore } from '@/stores/majhiLotBills'
import { useLotsStore } from '@/stores/lots'
import { useItemsStore } from '@/stores/items'
import { useMajhisStore } from '@/stores/majhis'
import { useStoresStore } from '@/stores/stores'
import AppLayout from '@/components/layouts/AppLayout.vue'
import MajhiLotBillList from '@/components/features/majhiLotBills/MajhiLotBillList.vue'
import { formatCurrency } from '@/utils/currency'

const majhiLotBillsStore = useMajhiLotBillsStore()
const lotsStore = useLotsStore()
const itemsStore = useItemsStore()
const majhisStore = useMajhisStore()
const storesStore = useStoresStore()

onMounted(async () => {
    await Promise.all([
        lotsStore.fetchLots(),
        itemsStore.fetchItems(),
        majhisStore.fetchMajhis(),
        storesStore.fetchStores(),
    ])
})
</script>