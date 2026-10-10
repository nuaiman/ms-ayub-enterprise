<!-- src/views/CustomerUnloadBillsView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Bills</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">
                    {{ customerUnloadBillsStore.totalBills }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Billed</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">
                    {{ formatCurrency(customerUnloadBillsStore.totalBilled) }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">
                    {{ formatCurrency(customerUnloadBillsStore.totalPaid) }}
                </p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                <p class="text-2xl font-bold mt-1"
                    :class="outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(outstanding) }}
                </p>
            </div>
        </div>

        <!-- List -->
        <div class="flex-1 min-h-0 mt-6">
            <CustomerUnloadBillList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useCustomerUnloadBillsStore } from '@/stores/customerUnloadBills'
import CustomerUnloadBillList from '@/components/features/customerUnloadBills/CustomerUnloadBillList.vue'
import { formatCurrency } from '@/utils/currency'

const customerUnloadBillsStore = useCustomerUnloadBillsStore()

const outstanding = computed(() =>
    Math.max(0, customerUnloadBillsStore.totalBilled - customerUnloadBillsStore.totalPaid)
)

onMounted(() => {
    customerUnloadBillsStore.fetchCustomerUnloadBills()
})
</script>