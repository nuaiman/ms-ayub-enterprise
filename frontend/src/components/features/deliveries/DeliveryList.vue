<!-- src/components/features/deliveries/DeliveryList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Deliveries</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredDeliveries.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search deliveries..."
                        class="w-full sm:w-56 pl-9 pr-8 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--color-text-secondary)"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <button v-if="searchQuery" @click="clearSearch" type="button"
                        class="absolute right-2.5 top-1/2 -translate-y-1/2 text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Month Filter -->
                <MonthFilter v-model="monthFilter" :items="availableMonths" />

                <!-- Copy Button -->
                <button @click="handleCopyToClipboard"
                    class="h-9 w-9 flex items-center justify-center border border-(--color-border) rounded-lg text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors relative shrink-0"
                    title="Copy table to clipboard">
                    <svg v-if="!clipboardStore.copied" class="w-4 h-4" fill="none" stroke="currentColor"
                        viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                    <svg v-else class="w-4 h-4 text-(--color-green)" fill="none" stroke="currentColor"
                        viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                </button>

                <!-- Create -->
                <button v-if="canManageDeliveries" @click="createDialogOpen = true"
                    class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    <span class="hidden sm:inline">Create</span>
                </button>
            </div>
        </div>

        <!-- Table -->
        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <!-- Header Row -->
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('customer_id')">
                        <span class="flex items-center gap-1">
                            Customer
                            <svg v-if="sortField === 'customer_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2">Receiver</div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('from_location')">
                        <span class="flex items-center gap-1">
                            From
                            <svg v-if="sortField === 'from_location'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('to_location')">
                        <span class="flex items-center gap-1">
                            To
                            <svg v-if="sortField === 'to_location'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('delivery_date')">
                        <span class="flex items-center gap-1">
                            Date
                            <svg v-if="sortField === 'delivery_date'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 flex items-center justify-end">Actions</div>
                </div>

                <!-- Loading -->
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading deliveries...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredDeliveries.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No deliveries found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new delivery to get started' }}</p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <DeliveryRow v-for="delivery in filteredDeliveries" :key="delivery.id" :delivery="delivery"
                        @view="openDetailDialog" @edit="handleEditDelivery" @delete="handleDeleteDelivery"
                        @manage-items="handleManageItems" @updated="fetchDeliveries" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredDeliveries.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredDeliveries.length }} of {{
                deliveriesStore.deliveries.length }} deliveries</p>
        </div>

        <!-- Dialogs -->
        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Delivery</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new delivery to the system</p>
            </div>
            <DeliveryForm mode="create" @delivery-created="handleDeliveryCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <DeliveryDetail v-if="selectedDelivery" :delivery="selectedDelivery" @close="closeDetailDialog"
                @edit="handleEditDeliveryFromDetail" @manage-items="handleManageItemsFromDetail"
                @updated="fetchDeliveries" />
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Delivery</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update delivery information</p>
            </div>
            <DeliveryForm v-if="selectedDelivery" mode="edit" :delivery="selectedDelivery"
                @delivery-updated="handleDeliveryUpdated" @cancel="editDialogOpen = false" />
        </BaseDialog>

        <!-- Delete Confirmation Dialog -->
        <BaseDialog v-model="deleteDialogOpen" max-width="sm">
            <div class="flex items-center gap-3">
                <div
                    class="w-10 h-10 rounded-full bg-(--color-red)/10 text-(--color-red) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Delivery</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete delivery #<span class="font-medium text-(--color-text-primary)">{{
                    selectedDelivery?.id }}</span>?
            </p>
            <template #actions>
                <button @click="deleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Delivery, DeliverySortField, SortDirection } from '@/types/delivery'
import DeliveryRow from './DeliveryRow.vue'
import DeliveryForm from './DeliveryForm.vue'
import DeliveryDetail from './DeliveryDetail.vue'
import MonthFilter from '@/components/ui/MonthFilter.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const usersStore = useUsersStore()
const deliveryItemsStore = useDeliveryItemsStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const monthFilter = ref('')
const sortField = ref<DeliverySortField>('delivery_date')
const sortDirection = ref<SortDirection>('desc')

// Dialogs
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedDelivery = ref<Delivery | null>(null)

const canManageDeliveries = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

// Get unique months from deliveries
const availableMonths = computed(() => {
    const months = new Set<string>()
    deliveriesStore.deliveries.forEach(d => {
        const date = new Date(d.delivery_date)
        const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
        months.add(month)
    })
    return Array.from(months).sort((a, b) => b.localeCompare(a))
})

const getCustomerName = (id: number | null): string => {
    if (!id) return '—'
    return customersStore.getCustomerName(id)
}

const filteredDeliveries = computed(() => {
    let result = [...deliveriesStore.deliveries]

    // Search
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(d =>
            (d.receiver_name && d.receiver_name.toLowerCase().includes(query)) ||
            (d.receiver_phone && d.receiver_phone.toLowerCase().includes(query)) ||
            (d.from_location && d.from_location.toLowerCase().includes(query)) ||
            (d.to_location && d.to_location.toLowerCase().includes(query)) ||
            (d.notes && d.notes.toLowerCase().includes(query)) ||
            getCustomerName(d.customer_id).toLowerCase().includes(query) ||
            usersStore.getUserName(d.user_id).toLowerCase().includes(query)
        )
    }

    // Filter by month
    if (monthFilter.value) {
        result = result.filter(d => {
            const date = new Date(d.delivery_date)
            const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
            return month === monthFilter.value
        })
    }

    // Sort
    result.sort((a, b) => {
        let comparison = 0
        switch (sortField.value) {
            case 'customer_id':
                comparison = (a.customer_id || 0) - (b.customer_id || 0)
                break
            case 'delivery_date':
                comparison = new Date(a.delivery_date).getTime() - new Date(b.delivery_date).getTime()
                break
            case 'from_location':
                comparison = (a.from_location || '').localeCompare(b.from_location || '')
                break
            case 'to_location':
                comparison = (a.to_location || '').localeCompare(b.to_location || '')
                break
            case 'created_at':
                comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
                break
            default:
                comparison = 0
        }
        return sortDirection.value === 'desc' ? -comparison : comparison
    })

    return result
})

const fetchDeliveries = async () => {
    loading.value = true
    try {
        await Promise.all([
            deliveriesStore.fetchDeliveries(),
            customersStore.fetchCustomers(),
            usersStore.fetchUsers(),
            deliveryItemsStore.fetchDeliveryItems()
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    searchQuery.value = target.value
}

const clearSearch = () => {
    searchQuery.value = ''
}

const toggleSort = (field: DeliverySortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Customer\tReceiver\tFrom\tTo\tDate'
    const rows = filteredDeliveries.value.map(d => {
        return `${getCustomerName(d.customer_id)}\t${d.receiver_name || ''}\t${d.from_location || ''}\t${d.to_location || ''}\t${new Date(d.delivery_date).toLocaleDateString()}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

// Dialog handlers
const openDetailDialog = (delivery: Delivery) => {
    selectedDelivery.value = delivery
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedDelivery.value = null }, 300)
}

const handleEditDelivery = (delivery: Delivery) => {
    selectedDelivery.value = delivery
    editDialogOpen.value = true
}

const handleEditDeliveryFromDetail = (delivery: Delivery) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedDelivery.value = delivery
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteDelivery = (delivery: Delivery) => {
    selectedDelivery.value = delivery
    deleteDialogOpen.value = true
}

const handleManageItems = (delivery: Delivery) => {
    selectedDelivery.value = delivery
    // Close detail dialog if open
    detailDialogOpen.value = false
    // Emit event to parent to open delivery items management
    // This will be handled in the parent view
    push.info(`Manage items for delivery #${delivery.id}`)
}

const handleManageItemsFromDetail = (delivery: Delivery) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        handleManageItems(delivery)
    }, 300)
}

const confirmDelete = async () => {
    if (!selectedDelivery.value) return
    const success = await deliveriesStore.deleteDelivery(selectedDelivery.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchDeliveries()
    }
}

const handleDeliveryCreated = async () => {
    createDialogOpen.value = false
    await fetchDeliveries()
}

const handleDeliveryUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchDeliveries()
}

// Auto-select current month
watch(() => deliveriesStore.deliveries, (newDeliveries) => {
    if (newDeliveries.length > 0 && !monthFilter.value) {
        const currentMonth = new Date().toISOString().slice(0, 7)
        const hasCurrentMonth = newDeliveries.some(d => {
            const date = new Date(d.delivery_date)
            const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
            return month === currentMonth
        })
        if (hasCurrentMonth) {
            monthFilter.value = currentMonth
        }
    }
}, { immediate: true })

onMounted(() => {
    fetchDeliveries()
})
</script>