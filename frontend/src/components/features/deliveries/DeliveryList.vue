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
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search deliveries..."
                        class="w-full sm:w-64 pl-9 pr-8 py-2 rounded-lg text-sm
                               bg-(--color-muted-bg) border border-(--color-border)
                               text-(--color-text-primary) placeholder:text-(--color-text-secondary)
                               focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
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

                <button @click="handleCopyToClipboard"
                    class="h-9 w-9 flex items-center justify-center border border-(--color-border) rounded-lg text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors relative shrink-0"
                    title="Copy to clipboard">
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

                <button v-if="canManage" @click="createDialogOpen = true"
                    class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    <span class="hidden sm:inline">Create</span>
                </button>
            </div>
        </div>

        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
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

                <div v-else-if="filteredDeliveries.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M9 17a2 2 0 11-4 0 2 2 0 014 0zM20 17a2 2 0 11-4 0 2 2 0 014 0zM13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a1 1 0 001 1h1" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No deliveries found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new delivery to get started' }}</p>
                    </div>
                </div>

                <div v-else class="space-y-4">
                    <div
                        class="grid grid-cols-12 items-center w-full py-2.5 px-3 rounded-lg bg-(--color-muted-bg)/30 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                        <div class="col-span-4 min-w-0 pr-3">Customer</div>
                        <div class="col-span-3 min-w-0 pr-3">Date / Receiver</div>
                        <div class="col-span-3 min-w-0 pr-3">Route</div>
                        <div class="col-span-2 text-right">Items</div>
                    </div>

                    <div v-for="delivery in filteredDeliveries" :key="delivery.id"
                        class="rounded-xl border border-(--color-border) bg-(--color-surface)">
                        <DeliveryRow :delivery="delivery" :items-count="getItemsFor(delivery.id).length"
                            @view="openDetailDialog" @edit="handleEditDelivery" @delete="handleDeleteDelivery"
                            @add-item="handleAddItem" />

                        <div class="border-t border-(--color-border)">
                            <div v-if="itemsLoading[delivery.id]"
                                class="px-3 py-3 text-xs text-(--color-text-secondary) italic">
                                Loading itemsâ€¦
                            </div>
                            <div v-else-if="getItemsFor(delivery.id).length === 0"
                                class="px-3 py-3 text-xs text-(--color-text-secondary) italic">
                                No items for this delivery
                            </div>
                            <div v-else>
                                <DeliveryItemRow v-for="it in getItemsFor(delivery.id)" :key="it.id" :item="it"
                                    @edit="handleEditItem" @delete="handleDeleteItem" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="!loading && filteredDeliveries.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredDeliveries.length }} of {{
                deliveriesStore.deliveries.length }} deliveries Â· {{ deliveriesStore.totalItemCount }} items</p>
        </div>

        <!-- ===================================================== -->
        <!-- DIALOGS                                                -->
        <!-- ===================================================== -->

        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Delivery</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Record a new delivery and its items</p>
            </div>
            <DeliveryForm mode="create" @delivery-created="handleDeliveryCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <DeliveryDetail v-if="selectedDelivery" :delivery="selectedDelivery" @close="closeDetailDialog"
                @edit="handleEditDeliveryFromDetail" @updated="fetchDeliveries" />
        </BaseDialog>

        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Delivery</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update delivery information</p>
            </div>
            <DeliveryForm v-if="selectedDelivery" mode="edit" :delivery="selectedDelivery"
                @delivery-updated="handleDeliveryUpdated" @cancel="editDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="addItemDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Add Delivery Item</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    For delivery #{{ selectedDelivery?.id }}
                </p>
            </div>
            <DeliveryItemAddForm v-if="selectedDelivery" :delivery-id="selectedDelivery.id"
                :customer-id="selectedDelivery.customer_id" @item-created="handleItemCreated"
                @cancel="addItemDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="editItemDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Delivery Item</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update item details</p>
            </div>
            <DeliveryItemEditForm v-if="selectedItem" :item="selectedItem" @item-updated="handleItemUpdated"
                @cancel="editItemDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Delivery</h2>
                    <p class="text-xs text-(--color-text-secondary)">This will also delete all its items</p>
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

        <BaseDialog v-model="deleteItemDialogOpen" max-width="sm">
            <div class="flex items-center gap-3">
                <div
                    class="w-10 h-10 rounded-full bg-(--color-red)/10 text-(--color-red) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Item</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Delete this item from delivery #<span class="font-medium text-(--color-text-primary)">{{
                    selectedItem?.delivery_id }}</span>?
            </p>
            <template #actions>
                <button @click="deleteItemDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmItemDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useMajhisStore } from '@/stores/majhis'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Delivery, DeliveryItem } from '@/types/delivery'
import DeliveryRow from './DeliveryRow.vue'
import DeliveryItemRow from './DeliveryItemRow.vue'
import DeliveryForm from './DeliveryForm.vue'
import DeliveryDetail from './DeliveryDetail.vue'
import DeliveryItemEditForm from './DeliveryItemEditForm.vue'
import DeliveryItemAddForm from './DeliveryItemAddForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'

const deliveriesStore = useDeliveriesStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const majhisStore = useMajhisStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = computed(() => deliveriesStore.searchQuery)
const itemsLoading = computed(() => deliveriesStore.itemsLoading)

const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const addItemDialogOpen = ref(false)
const editItemDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const deleteItemDialogOpen = ref(false)

const selectedDelivery = ref<Delivery | null>(null)
const selectedItem = ref<DeliveryItem | null>(null)

const canManage = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const filteredDeliveries = computed(() => deliveriesStore.filteredDeliveries)

const getItemsFor = (id: number): DeliveryItem[] => deliveriesStore.getItemsFor(id)

// =============================================================================
// Fetch
// =============================================================================

const fetchDeliveries = async () => {
    loading.value = true
    try {
        await Promise.all([
            deliveriesStore.fetchDeliveries(),
            customersStore.fetchCustomers(),
            storesStore.fetchStores(),
            lotsStore.fetchLots(),
            majhisStore.fetchMajhis(),
        ])
        // Load items for each delivery (uses store cache).
        await Promise.all(deliveriesStore.deliveries.map(d => deliveriesStore.loadItemsFor(d.id)))
    } finally {
        loading.value = false
    }
}

// =============================================================================
// Search / Copy
// =============================================================================

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    deliveriesStore.setSearchQuery(target.value)
}

const clearSearch = () => { deliveriesStore.clearSearch() }

const handleCopyToClipboard = async () => {
    const lines: string[] = []
    lines.push('ID\tCustomer\tDate\tReceiver\tRoute\tItems')

    for (const d of filteredDeliveries.value) {
        const items = getItemsFor(d.id)
        const customer = d.customer_id ? customersStore.getCustomerName(d.customer_id) : ''
        const receiver = [d.receiver_name, d.receiver_phone].filter(Boolean).join(' / ')
        const route = `${d.from_location || ''} â†’ ${d.to_location || ''}`
        lines.push(
            `${d.id}\t${customer}\t${new Date(d.delivery_date).toLocaleDateString()}\t${receiver}\t${route}\t${items.length}`
        )
        for (const it of items) {
            lines.push(
                `\t\t\t\tStore #${it.store_id}\tQty:${it.quantity} Wt:${it.weight}`
            )
        }
    }

    await clipboardStore.copyToClipboard(lines.join('\n'))
}

// =============================================================================
// Dialog handlers
// =============================================================================

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

const handleAddItem = (delivery: Delivery) => {
    selectedDelivery.value = delivery
    addItemDialogOpen.value = true
}

const handleEditItem = (item: DeliveryItem) => {
    selectedItem.value = item
    editItemDialogOpen.value = true
}

const handleDeleteItem = (item: DeliveryItem) => {
    selectedItem.value = item
    deleteItemDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedDelivery.value) return
    const ok = await deliveriesStore.deleteDelivery(selectedDelivery.value.id)
    if (ok) {
        deleteDialogOpen.value = false
        await fetchDeliveries()
    }
}

const confirmItemDelete = async () => {
    if (!selectedItem.value) return
    const ok = await deliveriesStore.deleteDeliveryItem(selectedItem.value.id)
    if (ok) {
        deleteItemDialogOpen.value = false
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

const handleItemCreated = async () => {
    addItemDialogOpen.value = false
    if (selectedDelivery.value) {
        await deliveriesStore.loadItemsFor(selectedDelivery.value.id, true)
    }
}

const handleItemUpdated = async () => {
    editItemDialogOpen.value = false
    if (selectedItem.value) {
        await deliveriesStore.loadItemsFor(selectedItem.value.delivery_id, true)
    }
}

onMounted(() => { fetchDeliveries() })
</script>