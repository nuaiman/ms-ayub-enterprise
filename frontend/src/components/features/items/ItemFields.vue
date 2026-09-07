<!-- src/components/features/items/ItemFields.vue -->
<template>
    <div class="space-y-4">
        <!-- Product Name -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Product Name
            </label>
            <input v-model="localProductName" type="text" placeholder="Enter product name" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed"
                @input="$emit('update:productName', localProductName)" />
        </div>

        <!-- Category -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Category
            </label>
            <input v-model="localCategory" type="text" placeholder="Enter category" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed"
                @input="$emit('update:category', localCategory)" />
        </div>

        <!-- Customer -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Customer
            </label>
            <select :value="customerId"
                @change="$emit('update:customerId', parseInt(($event.target as HTMLSelectElement).value) || null)"
                :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option :value="null">Select a customer</option>
                <option v-for="customer in customerOptions" :key="customer.id" :value="customer.id">
                    {{ getCustomerDisplayName(customer) }}
                </option>
            </select>
        </div>

        <!-- Is Active -->
        <div class="flex items-center gap-3">
            <label class="text-sm font-medium text-(--color-text-primary)">Active</label>
            <div @click="$emit('update:isActive', !isActive)"
                class="w-11 h-6 rounded-full cursor-pointer transition-colors duration-200 flex items-center px-0.5"
                :class="isActive ? 'bg-(--color-green)' : 'bg-(--color-muted-bg)'">
                <div class="w-5 h-5 rounded-full bg-white shadow-sm transition-transform duration-200"
                    :class="isActive ? 'translate-x-5' : 'translate-x-0'"></div>
            </div>
        </div>

        <!-- Notes -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea v-model="localNotes" rows="3" placeholder="Enter any notes about this item" :disabled="disabled"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"
                @input="$emit('update:notes', localNotes)"></textarea>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Customer } from '@/types/customer'

const props = defineProps<{
    productName: string
    category: string
    customerId: number | null
    isActive: boolean
    notes: string
    customerOptions?: Customer[]
    disabled?: boolean
}>()

const emit = defineEmits<{
    (e: 'update:productName', value: string): void
    (e: 'update:category', value: string): void
    (e: 'update:customerId', value: number | null): void
    (e: 'update:isActive', value: boolean): void
    (e: 'update:notes', value: string): void
}>()

const localProductName = ref(props.productName)
const localCategory = ref(props.category)
const localNotes = ref(props.notes)

watch(() => props.productName, (val) => { localProductName.value = val })
watch(() => props.category, (val) => { localCategory.value = val })
watch(() => props.notes, (val) => { localNotes.value = val })

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}
</script>