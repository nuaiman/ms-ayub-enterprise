<!-- src/components/features/deliveryItems/DeliveryItemList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">
                    Delivery Items
                </h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredDeliveryItems.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search items..."
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

                <!-- Delivery Filter -->
                <select v-model="deliveryFilter"
                    class="px-3 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                    <option value="">All Deliveries</option>
                    <option v-for="delivery in deliveryOptions" :key="delivery.id" :value="delivery.id">
                        #{{ delivery.id }} - {{ getDeliveryDisplayName(delivery) }}
                    </option>
                </select>

                <!-- Copy -->
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
                <button v-if="canManageDeliveryItems" @click="createDialogOpen = true"
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
                    <div class="col-span-2 min-w-0 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('delivery_id')">
                        <span class="flex items-center gap-1">
                            Delivery
                            <svg v-if="sortField === 'delivery_id'" class="w-3 h-3 shrink-0"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-3 min-w-0 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('item_id')">
                        <span class="flex items-center gap-1">
                            Item
                            <svg v-if="sortField === 'item_id'" class="w-3 h-3 shrink-0"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 min-w-0 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('lot_id')">
                        <span class="flex items-center gap-1">
                            Lot
                            <svg v-if="sortField === 'lot_id'" class="w-3 h-3 shrink-0"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 min-w-0 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('quantity')">
                        <span class="flex items-center gap-1">
                            Qty
                            <svg v-if="sortField === 'quantity'" class="w-3 h-3 shrink-0"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 min-w-0 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('weight')">
                        <span class="flex items-center gap-1">
                            Wgt
                            <svg v-if="sortField === 'weight'" class="w-3 h-3 shrink-0"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 min-w-0">Majhi</div>
                    <div class="col-span-1 flex items-center justify-end">Actions</div>
                </div>

                <!-- Loading -->
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading delivery items...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredDeliveryItems.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No delivery items found</p>
                        <p class="text-xs text-(--color-text-secondary)">
                            {{ searchQuery ? 'Try adjusting your search' : 'Create a new delivery item to get started'
                            }}
                        </p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <DeliveryItemRow v-for="item in filteredDeliveryItems" :key="item.id" :deliveryItem="item"
                        @view="openDetailDialog" @edit="handleEditDeliveryItem" @delete="handleDeleteDeliveryItem"
                        @updated="fetchDeliveryItems" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredDeliveryItems.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">
                Showing {{ filteredDeliveryItems.length }} of {{ deliveryItemsStore.deliveryItems.length }} items
            </p>
        </div>

        <!-- Dialogs -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create Delivery Item</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add an item to a delivery</p>
            </div>
            <DeliveryItemForm mode="create" @delivery-item-created="handleDeliveryItemCreated"
                @cancel="createDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <DeliveryItemDetail v-if="selectedDeliveryItem" :deliveryItem="selectedDeliveryItem"
                @close="closeDetailDialog" @edit="handleEditDeliveryItemFromDetail" @updated="fetchDeliveryItems" />
        </BaseDialog>

        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Delivery Item</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update delivery item information</p>
            </div>
            <DeliveryItemForm v-if="selectedDeliveryItem" mode="edit" :deliveryItem="selectedDeliveryItem"
                @delivery-item-updated="handleDeliveryItemUpdated" @cancel="editDialogOpen = false" />
        </BaseDialog>

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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Delivery Item</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">Are you sure you want to delete this delivery item?
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
import { ref, computed, onMounted } from 'vue'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useItemsStore } from '@/stores/items'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { DeliveryItem, DeliveryItemSortField, SortDirection } from '@/types/deliveryItem'
import type { Delivery } from '@/types/delivery'
import DeliveryItemRow from './DeliveryItemRow.vue'
import DeliveryItemForm from './DeliveryItemForm.vue'
import DeliveryItemDetail from './DeliveryItemDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'

const deliveryItemsStore = useDeliveryItemsStore()
const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const itemsStore = useItemsStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const deliveryFilter = ref('')
const sortField = ref<DeliveryItemSortField>('created_at')
const sortDirection = ref<SortDirection>('desc')

const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const selectedDeliveryItem = ref<DeliveryItem | null>(null)

const canManageDeliveryItems = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const deliveryOptions = computed(() => deliveriesStore.deliveries)

const getDeliveryDisplayName = (delivery: Delivery): string => {
    const customerName = customersStore.getCustomerName(delivery.customer_id!)
    const date = new Date(delivery.delivery_date).toLocaleDateString()
    return `${customerName} - ${date}`
}

const filteredDeliveryItems = computed(() => {
    let result = [...deliveryItemsStore.deliveryItems]

    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(item =>
            itemsStore.getItemName(item.item_id).toLowerCase().includes(query) ||
            lotsStore.getLotName(item.lot_id).toLowerCase().includes(query) ||
            (item.majhi_id && majhisStore.getMajhiName(item.majhi_id).toLowerCase().includes(query)) ||
            String(item.quantity).includes(query) ||
            String(item.weight).includes(query) ||
            (item.vehicle_number && item.vehicle_number.toLowerCase().includes(query)) ||
            (item.driver_number && item.driver_number.toLowerCase().includes(query)) ||
            (item.notes && item.notes.toLowerCase().includes(query))
        )
    }

    if (deliveryFilter.value) {
        const filterId = parseInt(deliveryFilter.value)
        result = result.filter(item => item.delivery_id === filterId)
    }

    result.sort((a, b) => {
        let comparison = 0
        switch (sortField.value) {
            case 'delivery_id': comparison = a.delivery_id - b.delivery_id; break
            case 'store_id': comparison = a.store_id - b.store_id; break
            case 'item_id': comparison = a.item_id - b.item_id; break
            case 'lot_id': comparison = a.lot_id - b.lot_id; break
            case 'quantity': comparison = a.quantity - b.quantity; break
            case 'weight': comparison = a.weight - b.weight; break
            case 'created_at': comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime(); break
            default: comparison = 0
        }
        return sortDirection.value === 'desc' ? -comparison : comparison
    })

    return result
})

const fetchDeliveryItems = async () => {
    loading.value = true
    try {
        await Promise.all([
            deliveryItemsStore.fetchDeliveryItems(),
            deliveriesStore.fetchDeliveries(),
            customersStore.fetchCustomers(),
            itemsStore.fetchItems(),
            lotsStore.fetchLots(),
            majhisStore.fetchMajhis()
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    searchQuery.value = target.value
}

const clearSearch = () => { searchQuery.value = '' }

const toggleSort = (field: DeliveryItemSortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Delivery\tItem\tLot\tQuantity\tWeight\tMajhi'
    const rows = filteredDeliveryItems.value.map(item => {
        return `${item.delivery_id}\t${itemsStore.getItemName(item.item_id)}\t${lotsStore.getLotName(item.lot_id)}\t${item.quantity} ${item.quantity_unit}\t${item.weight} ${item.weight_unit}\t${item.majhi_id ? majhisStore.getMajhiName(item.majhi_id) : ''}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const openDetailDialog = (item: DeliveryItem) => {
    selectedDeliveryItem.value = item
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedDeliveryItem.value = null }, 300)
}

const handleEditDeliveryItem = (item: DeliveryItem) => {
    selectedDeliveryItem.value = item
    editDialogOpen.value = true
}

const handleEditDeliveryItemFromDetail = (item: DeliveryItem) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedDeliveryItem.value = item
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteDeliveryItem = (item: DeliveryItem) => {
    selectedDeliveryItem.value = item
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedDeliveryItem.value) return
    const success = await deliveryItemsStore.deleteDeliveryItem(selectedDeliveryItem.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchDeliveryItems()
    }
}

const handleDeliveryItemCreated = async () => {
    createDialogOpen.value = false
    await fetchDeliveryItems()
}

const handleDeliveryItemUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchDeliveryItems()
}

onMounted(() => {
    fetchDeliveryItems()
})
</script>