<!-- src/views/DeliveryItemsView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Items</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{
                        deliveryItemsStore.deliveryItems.length }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Quantity</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ deliveryItemsStore.totalQuantity }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Weight</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ deliveryItemsStore.totalWeight }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Majhi</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ withMajhiCount }}</p>
                </div>
            </div>

            <!-- Delivery Item List -->
            <div class="flex-1 min-h-0 mt-6">
                <DeliveryItemList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import AppLayout from '@/components/layouts/AppLayout.vue'
import DeliveryItemList from '@/components/features/deliveryItems/DeliveryItemList.vue'

const deliveryItemsStore = useDeliveryItemsStore()
const deliveriesStore = useDeliveriesStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()

const withMajhiCount = computed(() => {
    return deliveryItemsStore.deliveryItems.filter(d => d.majhi_id).length
})

onMounted(async () => {
    await Promise.all([
        deliveryItemsStore.fetchDeliveryItems(),
        deliveriesStore.fetchDeliveries(),
        storesStore.fetchStores(),
        lotsStore.fetchLots(),
        majhisStore.fetchMajhis()
    ])
})
</script>