<!-- src/views/CustomersView.vue -->
<template>
    <AppLayout>
        <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
            <!-- Stats -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Customers</p>
                    <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ customersStore.customers.length }}
                    </p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Company</p>
                    <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ withCompanyCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Contact Person</p>
                    <p class="text-2xl font-bold text-(--color-green) mt-1">{{ withContactPersonCount }}</p>
                </div>
                <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">With Email</p>
                    <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ withEmailCount }}</p>
                </div>
            </div>

            <!-- Customer List -->
            <div class="flex-1 min-h-0 mt-6">
                <CustomerList />
            </div>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useCustomersStore } from '@/stores/customers'
import AppLayout from '@/components/layouts/AppLayout.vue'
import CustomerList from '@/components/features/customers/CustomerList.vue'

const customersStore = useCustomersStore()

const withCompanyCount = computed(() => {
    return customersStore.customers.filter(c => c.company_name).length
})

const withContactPersonCount = computed(() => {
    return customersStore.customers.filter(c => c.contact_person).length
})

const withEmailCount = computed(() => {
    return customersStore.customers.filter(c => c.email).length
})

onMounted(() => {
    customersStore.fetchCustomers()
})
</script>