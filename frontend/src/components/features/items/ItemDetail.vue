<!-- src/components/features/items/ItemDetail.vue -->
<template>
    <div v-if="item" class="space-y-6">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="shrink-0">
                <div
                    class="w-20 h-20 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-10 h-10 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
                <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ item.product_name || 'Unnamed Item' }}
                </h2>
                <div class="flex items-center gap-2 flex-wrap mt-1">
                    <span v-if="item.category" class="text-sm text-(--color-text-secondary)">{{ item.category }}</span>
                    <span v-if="item.category" class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="text-sm text-(--color-text-secondary)">{{ getCustomerName(item.customer_id) }}</span>
                    <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
                    <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                        :class="item.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                        <span class="w-1.5 h-1.5 rounded-full"
                            :class="item.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                        {{ item.is_active ? 'Active' : 'Inactive' }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Meta -->
        <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
            <span class="text-xs text-(--color-text-secondary)">ID: {{ item.id }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Created: {{ formatDate(item.created_at) }}</span>
            <span class="w-px h-4 bg-(--color-border)"></span>
            <span class="text-xs text-(--color-text-secondary)">Updated: {{ formatDate(item.updated_at) }}</span>
        </div>

        <!-- Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Product Name -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Product Name</p>
                <p class="text-sm text-(--color-text-primary)">{{ item.product_name || '—' }}</p>
            </div>

            <!-- Category -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Category</p>
                <p class="text-sm text-(--color-text-primary)">{{ item.category || '—' }}</p>
            </div>

            <!-- Customer -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Customer</p>
                <p class="text-sm text-(--color-text-primary)">{{ getCustomerName(item.customer_id) }}</p>
            </div>

            <!-- Status -->
            <div>
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Status</p>
                <span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium border"
                    :class="item.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-red) text-(--color-red)'">
                    <span class="w-1.5 h-1.5 rounded-full"
                        :class="item.is_active ? 'bg-(--color-green)' : 'bg-(--color-red)'"></span>
                    {{ item.is_active ? 'Active' : 'Inactive' }}
                </span>
            </div>

            <!-- Notes -->
            <div v-if="item.notes" class="md:col-span-2">
                <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Notes</p>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
                    <p class="text-sm text-(--color-text-secondary) whitespace-pre-wrap">{{ item.notes }}</p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('edit', item)"
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
import type { Item } from '@/types/item'
import { useCustomersStore } from '@/stores/customers'

const props = defineProps<{
    item: Item | null
}>()

const emit = defineEmits<{
    'close': []
    'edit': [item: Item]
    'updated': []
}>()

const customersStore = useCustomersStore()

const getCustomerName = (id: number | null): string => {
    if (!id) return '—'
    return customersStore.getCustomerName(id)
}

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}
</script>