<!-- src/components/features/stores/StoreList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Stores</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredStores.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search stores..."
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

                <button @click="showActiveOnly = !showActiveOnly"
                    class="h-9 px-3 flex items-center gap-1.5 rounded-lg text-sm border border-(--color-border) transition-colors"
                    :class="showActiveOnly ? 'bg-(--color-blue)/10 border-(--color-blue) text-(--color-blue)' : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)'">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Active Only
                </button>

                <button @click="showInventoryOnly = !showInventoryOnly"
                    class="h-9 px-3 flex items-center gap-1.5 rounded-lg text-sm border border-(--color-border) transition-colors"
                    :class="showInventoryOnly ? 'bg-(--color-green)/10 border-(--color-green) text-(--color-green)' : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)'">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    Has Inventory
                </button>

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

                <button v-if="canManageStores" @click="createDialogOpen = true"
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
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('lot_id')">
                        <span class="flex items-center gap-1">
                            Lot / Product
                            <svg v-if="sortField === 'lot_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-3">Customer</div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('godown_id')">
                        <span class="flex items-center gap-1">
                            Godown
                            <svg v-if="sortField === 'godown_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2">Stock</div>
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('is_active')">
                        <span class="flex items-center gap-1">
                            Status
                            <svg v-if="sortField === 'is_active'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 flex items-center justify-end">Actions</div>
                </div>

                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading stores...</p>
                    </div>
                </div>

                <div v-else-if="filteredStores.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No stores found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new store to get started' }}</p>
                    </div>
                </div>

                <div v-else>
                    <StoreRow v-for="store in filteredStores" :key="store.id" :store="store" @view="openDetailDialog"
                        @edit="handleEditStore" @readd="handleReaddStore" @delete="handleDeleteStore"
                        @toggle-active="handleToggleActive" @updated="fetchStores" />
                </div>
            </div>
        </div>

        <div v-if="!loading && filteredStores.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredStores.length }} of {{
                storesStore.stores.length }} stores</p>
        </div>

        <!-- Dialogs -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new store to the system</p>
            </div>
            <StoreForm mode="create" @store-created="handleStoreCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="readdDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Re-add Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Create a new store based on this one</p>
            </div>
            <StoreForm v-if="readdStore" mode="create" :prefill="readdStore" @store-created="handleStoreCreated"
                @cancel="readdDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <StoreDetail v-if="selectedStore" :store="selectedStore" @close="closeDetailDialog"
                @edit="handleEditStoreFromDetail" @updated="fetchStores" />
        </BaseDialog>

        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update store information</p>
            </div>
            <StoreForm v-if="selectedStore" mode="edit" :store="selectedStore" @store-updated="handleStoreUpdated"
                @cancel="editDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Store</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete this store for "<span class="font-medium text-(--color-text-primary)">{{
                    getLotDisplayName(selectedStore?.lot_id) }}</span>"?
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
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useGodownsStore } from '@/stores/godowns'
import { useCustomersStore } from '@/stores/customers'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Store, StoreSortField, SortDirection } from '@/types/store'
import StoreRow from './StoreRow.vue'
import StoreForm from './StoreForm.vue'
import StoreDetail from './StoreDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()
const customersStore = useCustomersStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const showActiveOnly = ref(false)
const showInventoryOnly = ref(false)
const sortField = ref<StoreSortField>('lot_id')
const sortDirection = ref<SortDirection>('asc')

const createDialogOpen = ref(false)
const readdDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedStore = ref<Store | null>(null)
const readdStore = ref<Store | null>(null)

const canManageStores = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const getLotDisplayName = (lotId?: number): string => {
    if (!lotId) return 'Unknown'
    const lot = lotsStore.getLotById(lotId)
    if (!lot) return `Lot #${lotId}`
    return lotsStore.getLotDisplayName(lot)
}

const filteredStores = computed(() => {
    let result = [...storesStore.stores]

    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(s => {
            const lot = lotsStore.getLotById(s.lot_id)
            const customerName = lot?.customer_id
                ? customersStore.getCustomerName(lot.customer_id).toLowerCase()
                : ''
            const lotName = lot ? lotsStore.getLotDisplayName(lot).toLowerCase() : ''
            return (
                lotName.includes(query) ||
                godownsStore.getGodownName(s.godown_id).toLowerCase().includes(query) ||
                customerName.includes(query) ||
                String(s.quantity).includes(query) ||
                String(s.weight).includes(query) ||
                (s.notes && s.notes.toLowerCase().includes(query))
            )
        })
    }

    if (showActiveOnly.value) {
        result = result.filter(s => s.is_active)
    }

    if (showInventoryOnly.value) {
        result = result.filter(s => s.quantity > 0 || s.weight > 0)
    }

    result.sort((a, b) => {
        let comparison = 0
        switch (sortField.value) {
            case 'lot_id':
                comparison = a.lot_id - b.lot_id
                break
            case 'godown_id':
                comparison = a.godown_id - b.godown_id
                break
            case 'quantity':
                comparison = a.quantity - b.quantity
                break
            case 'weight':
                comparison = a.weight - b.weight
                break
            case 'is_active':
                comparison = (a.is_active === b.is_active) ? 0 : a.is_active ? -1 : 1
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

const fetchStores = async () => {
    loading.value = true
    try {
        await Promise.all([
            storesStore.fetchStores(),
            lotsStore.fetchLots(),
            godownsStore.fetchGodowns(),
            customersStore.fetchCustomers(),
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

const toggleSort = (field: StoreSortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Lot\tGodown\tQuantity\tWeight\tStatus'
    const rows = filteredStores.value.map(s => {
        const lotName = getLotDisplayName(s.lot_id)
        const godownName = godownsStore.getGodownName(s.godown_id)
        return `${lotName}\t${godownName}\t${s.quantity} ${s.quantity_unit}\t${s.weight} ${s.weight_unit}\t${s.is_active ? 'Active' : 'Inactive'}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const handleToggleActive = async (store: Store) => {
    const success = await storesStore.toggleStoreActive(store.id)
    if (success) {
        await fetchStores()
    }
}

const openDetailDialog = (store: Store) => {
    selectedStore.value = store
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedStore.value = null }, 300)
}

const handleEditStore = (store: Store) => {
    selectedStore.value = store
    editDialogOpen.value = true
}

const handleEditStoreFromDetail = (store: Store) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedStore.value = store
        editDialogOpen.value = true
    }, 300)
}

const handleReaddStore = (store: Store) => {
    readdStore.value = store
    readdDialogOpen.value = true
}

const handleDeleteStore = (store: Store) => {
    selectedStore.value = store
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedStore.value) return
    const success = await storesStore.deleteStore(selectedStore.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchStores()
    }
}

const handleStoreCreated = async () => {
    createDialogOpen.value = false
    readdDialogOpen.value = false
    readdStore.value = null
    await fetchStores()
}

const handleStoreUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchStores()
}

onMounted(() => {
    fetchStores()
})
</script>