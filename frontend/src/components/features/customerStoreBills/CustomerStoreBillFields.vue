<!-- src/components/features/customerStoreBills/CustomerStoreBillFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Customer + Store pickers (only when not locked AND not derived) -->
        <div v-if="!locked && !derived" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Customer <span class="text-(--color-red)">*</span>
                </label>
                <select :value="customerId ?? ''"
                    @change="$emit('update:customerId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a customer</option>
                    <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                        {{ getCustomerLabel(customer) }}
                    </option>
                </select>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Store <span class="text-(--color-red)">*</span>
                </label>
                <select :value="storeId ?? ''"
                    @change="$emit('update:storeId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                    :disabled="disabled" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="">Select a store</option>
                    <option v-for="store in storeOptions" :key="store.id" :value="store.id">
                        {{ getStoreLabel(store) }}
                    </option>
                </select>
            </div>
        </div>

        <!-- Locked display (edit mode) -->
        <div v-else-if="locked" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedCustomerLabel }}</p>
            </div>
            <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Store</p>
                <p class="text-sm text-(--color-text-primary) mt-0.5">{{ lockedStoreLabel }}</p>
            </div>
        </div>

        <!-- Bill Type + Rate (two-column) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Bill Type -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Bill Type <span class="text-(--color-red)">*</span>
                </label>
                <select :value="billType"
                    @change="$emit('update:billType', ($event.target as HTMLSelectElement).value as CustomerStoreBillType)"
                    :disabled="disabled" required
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option value="quantity">Quantity</option>
                    <option value="weight">Weight</option>
                </select>
            </div>

            <!-- Rate -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Rate <span class="text-(--color-red)">*</span>
                </label>
                <div class="relative">
                    <span
                        class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                    <input :value="rate"
                        @input="$emit('update:rate', parseFloat(($event.target as HTMLInputElement).value) || 0)"
                        type="number" step="0.01" min="0" placeholder="0.00" :disabled="disabled"
                        class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Customer } from '@/types/customer'
import type { Store } from '@/types/store'
import type { CustomerStoreBillType } from '@/types/customerStoreBill'
import { useLotsStore } from '@/stores/lots'

withDefaults(defineProps<{
    customerId?: number | null
    storeId?: number | null
    billType: CustomerStoreBillType
    rate: number
    customerOptions?: Customer[]
    storeOptions?: Store[]
    disabled?: boolean
    locked?: boolean
    lockedCustomerLabel?: string
    lockedStoreLabel?: string
    derived?: boolean
}>(), {
    customerId: null,
    storeId: null,
    customerOptions: () => [],
    storeOptions: () => [],
    disabled: false,
    locked: false,
    lockedCustomerLabel: '',
    lockedStoreLabel: '',
    derived: false,
})

defineEmits<{
    (e: 'update:customerId', value: number | null): void
    (e: 'update:storeId', value: number | null): void
    (e: 'update:billType', value: CustomerStoreBillType): void
    (e: 'update:rate', value: number): void
}>()

const lotsStore = useLotsStore()

const getCustomerLabel = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}

const getStoreLabel = (store: Store): string => {
    const lot = lotsStore.getLotById(store.lot_id)
    const lotName = lot ? lot.product_name : '—'
    const lotNum = lot ? lot.lot_number : '—'
    return `Store #${store.id} — ${lotName} (Lot ${lotNum})`
}
</script>