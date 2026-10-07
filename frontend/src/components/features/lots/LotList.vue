<!-- src/components/features/lots/LotList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Lots</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredLots.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search lots or stores..."
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

                <!-- Sort -->
                <select v-model="sortKey"
                    class="px-3 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                    <option value="lot_number_asc">Lot # (A→Z)</option>
                    <option value="lot_number_desc">Lot # (Z→A)</option>
                    <option value="product_asc">Product (A→Z)</option>
                    <option value="product_desc">Product (Z→A)</option>
                    <option value="customer_asc">Customer (A→Z)</option>
                    <option value="customer_desc">Customer (Z→A)</option>
                    <option value="created_desc">Newest first</option>
                    <option value="created_asc">Oldest first</option>
                </select>

                <!-- Copy -->
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

                <!-- Create -->
                <button v-if="canManageLots" @click="createDialogOpen = true"
                    class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    <span class="hidden sm:inline">Create</span>
                </button>
            </div>
        </div>

        <!-- Tree -->
        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading lots...</p>
                    </div>
                </div>

                <div v-else-if="filteredLots.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No lots found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new lot to get started' }}</p>
                    </div>
                </div>

                <div v-else class="space-y-4">
                    <!-- Column header -->
                    <div
                        class="grid grid-cols-12 items-center w-full py-2.5 px-3 rounded-lg bg-(--color-muted-bg)/30 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                        <div class="col-span-4 min-w-0 pr-3">Customer</div>
                        <div class="col-span-4 min-w-0 pr-3">Lot / Product</div>
                        <div class="col-span-2 min-w-0 pr-3">Units</div>
                        <div class="col-span-1 min-w-0 pr-3">Stores</div>
                        <div class="col-span-1 text-right">Actions</div>
                    </div>

                    <div v-for="lot in filteredLots" :key="lot.id"
                        class="rounded-xl border border-(--color-border) bg-(--color-surface)">
                        <!-- Lot header row -->
                        <LotRow :lot="lot" :stores-count="storesByLot(lot.id).length" @view="openDetailDialog"
                            @edit="handleEditLot" @delete="handleDeleteLot" @add-store="handleAddStore"
                            @transfer="handleTransferLot" @updated="fetchLots" />

                        <!-- Stores section -->
                        <div class="border-t border-(--color-border)">
                            <!-- Empty state -->
                            <div v-if="visibleStores(lot).length === 0"
                                class="px-3 py-3 text-xs text-(--color-text-secondary) italic">
                                No stores for this lot
                            </div>

                            <!-- Store children -->
                            <div v-else>
                                <StoreRow v-for="store in visibleStores(lot)" :key="store.id" :store="store" hide-lot
                                    @view="openStoreDetailDialog" @edit="handleEditStore" @delete="handleDeleteStore"
                                    @adjust="handleAdjustStore" @transfer="handleTransferStore" @updated="fetchLots" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredLots.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredLots.length }} of {{
                lotsStore.lots.length }} lots · {{ totalVisibleStores }} stores</p>
        </div>

        <!-- ===================================================== -->
        <!-- LOT DIALOGS                                            -->
        <!-- ===================================================== -->

        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Lot</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new lot to the system</p>
            </div>
            <LotForm mode="create" @lot-created="handleLotCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <LotDetail v-if="selectedLot" :lot="selectedLot" @close="closeDetailDialog" @edit="handleEditLotFromDetail"
                @updated="fetchLots" />
        </BaseDialog>

        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Lot</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update lot information</p>
            </div>
            <LotForm v-if="selectedLot" mode="edit" :lot="selectedLot" @lot-updated="handleLotUpdated"
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Lot</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete the lot "<span class="font-medium text-(--color-text-primary)">{{
                    selectedLot?.product_name }}</span>" (Lot {{ selectedLot?.lot_number }})?
            </p>
            <template #actions>
                <button @click="deleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>

        <!-- Lot Transfer Dialog (Account Transfer) -->
        <BaseDialog v-model="lotTransferDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Transfer Lot to Another Customer</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Reassign this lot and optionally set up customer store billing
                </p>
            </div>
            <LotTransferForm v-if="selectedLot" :preset-lot-id="selectedLot.id" :locked-lot="true"
                @transfer-created="handleLotTransferCreated" @cancel="lotTransferDialogOpen = false" />
        </BaseDialog>

        <!-- ===================================================== -->
        <!-- STORE DIALOGS                                          -->
        <!-- ===================================================== -->

        <!-- Add Store dialog (from LotRow action) -->
        <BaseDialog v-model="storeCreateDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Add Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Add a new store for <span class="font-medium text-(--color-text-primary)">{{
                        selectedLot?.product_name }}</span> (Lot {{ selectedLot?.lot_number }})
                </p>
            </div>
            <StoreForm v-if="selectedLot" mode="create" :preset-lot-id="selectedLot.id"
                @store-created="handleStoreCreated" @cancel="storeCreateDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="storeDetailDialogOpen" max-width="3xl">
            <StoreDetail v-if="selectedStore" :store="selectedStore" @close="closeStoreDetailDialog"
                @edit="handleEditStoreFromDetail" @updated="fetchLots" />
        </BaseDialog>

        <BaseDialog v-model="storeEditDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update store information</p>
            </div>
            <StoreForm v-if="selectedStore" mode="edit" :store="selectedStore" @store-updated="handleStoreUpdated"
                @cancel="storeEditDialogOpen = false" />
        </BaseDialog>

        <!-- Adjust Stock Dialog -->
        <BaseDialog v-model="adjustDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Adjust Store Stock</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Record a weight or quantity adjustment for this store
                </p>
            </div>
            <StoreAdjustmentForm v-if="selectedStore" :preset-store-id="selectedStore.id" :locked-store="true"
                @adjustment-created="handleAdjustmentCreated" @cancel="adjustDialogOpen = false" />
        </BaseDialog>

        <!-- Transfer Store Dialog -->
        <BaseDialog v-model="transferDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Transfer Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Move this store to a different godown
                </p>
            </div>
            <StoreTransferForm v-if="selectedStore" :preset-store-id="selectedStore.id" :locked-store="true"
                @transfer-created="handleTransferCreated" @cancel="transferDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="storeDeleteDialogOpen" max-width="sm">
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
                Are you sure you want to delete store #<span class="font-medium text-(--color-text-primary)">{{
                    selectedStore?.id }}</span>?
            </p>
            <template #actions>
                <button @click="storeDeleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmStoreDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Lot } from '@/types/lot'
import type { Store } from '@/types/store'
import LotRow from './LotRow.vue'
import LotForm from './LotForm.vue'
import LotDetail from './LotDetail.vue'
import StoreRow from '@/components/features/stores/StoreRow.vue'
import StoreForm from '@/components/features/stores/StoreForm.vue'
import StoreDetail from '@/components/features/stores/StoreDetail.vue'
import StoreAdjustmentForm from '@/components/features/storeAdjustments/StoreAdjustmentForm.vue'
import StoreTransferForm from '@/components/features/storeTransfers/StoreTransferForm.vue'
import LotTransferForm from '@/components/features/lotTransfers/LotTransferForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'

type SortKey =
    | 'lot_number_asc' | 'lot_number_desc'
    | 'product_asc' | 'product_desc'
    | 'customer_asc' | 'customer_desc'
    | 'created_desc' | 'created_asc'

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const sortKey = ref<SortKey>('lot_number_asc')

// Lot dialogs
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const lotTransferDialogOpen = ref(false)
const selectedLot = ref<Lot | null>(null)

// Store dialogs
const storeCreateDialogOpen = ref(false)
const storeDetailDialogOpen = ref(false)
const storeEditDialogOpen = ref(false)
const storeDeleteDialogOpen = ref(false)
const adjustDialogOpen = ref(false)
const transferDialogOpen = ref(false)
const selectedStore = ref<Store | null>(null)

const canManageLots = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const storesByLot = (lotID: number): Store[] => storesStore.getStoresByLotId(lotID)

const storeMatches = (store: Store, q: string): boolean => {
    const lot = lotsStore.getLotById(store.lot_id)
    const lotName = lot ? lot.product_name.toLowerCase() : ''
    const lotNum = lot ? lot.lot_number.toLowerCase() : ''
    const customerName = lot ? customersStore.getCustomerName(lot.customer_id).toLowerCase() : ''
    return (
        lotName.includes(q) ||
        lotNum.includes(q) ||
        customerName.includes(q) ||
        String(store.quantity).includes(q) ||
        String(store.weight).includes(q)
    )
}

const visibleStores = (lot: Lot): Store[] => {
    const stores = storesByLot(lot.id)
    if (!searchQuery.value) return stores

    const q = searchQuery.value.toLowerCase()

    const lotMatches =
        lot.product_name.toLowerCase().includes(q) ||
        lot.lot_number.toLowerCase().includes(q) ||
        customersStore.getCustomerName(lot.customer_id).toLowerCase().includes(q)

    if (lotMatches) return stores

    return stores.filter(s => storeMatches(s, q))
}

const filteredLots = computed(() => {
    let result = [...lotsStore.lots]

    if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase()
        result = result.filter(lot => {
            const lotMatches =
                lot.product_name.toLowerCase().includes(q) ||
                lot.lot_number.toLowerCase().includes(q) ||
                customersStore.getCustomerName(lot.customer_id).toLowerCase().includes(q)

            if (lotMatches) return true

            return storesByLot(lot.id).some(s => storeMatches(s, q))
        })
    }

    result.sort((a, b) => {
        switch (sortKey.value) {
            case 'lot_number_asc':
                return a.lot_number.localeCompare(b.lot_number)
            case 'lot_number_desc':
                return b.lot_number.localeCompare(a.lot_number)
            case 'product_asc':
                return a.product_name.localeCompare(b.product_name)
            case 'product_desc':
                return b.product_name.localeCompare(a.product_name)
            case 'customer_asc':
                return customersStore.getCustomerName(a.customer_id)
                    .localeCompare(customersStore.getCustomerName(b.customer_id))
            case 'customer_desc':
                return customersStore.getCustomerName(b.customer_id)
                    .localeCompare(customersStore.getCustomerName(a.customer_id))
            case 'created_asc':
                return new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
            case 'created_desc':
                return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        }
    })

    return result
})

const totalVisibleStores = computed(() =>
    filteredLots.value.reduce((sum, lot) => sum + visibleStores(lot).length, 0)
)

const fetchLots = async () => {
    loading.value = true
    try {
        await Promise.all([
            lotsStore.fetchLots(),
            customersStore.fetchCustomers(),
            storesStore.fetchStores(),
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

const handleCopyToClipboard = async () => {
    const lines: string[] = []
    lines.push('Product\tLot #\tCustomer\tUnits\tStores')

    for (const lot of filteredLots.value) {
        const stores = visibleStores(lot)

        lines.push(
            `${lot.product_name}\t${lot.lot_number}\t${customersStore.getCustomerName(lot.customer_id)}\t` +
            `Wt:${lot.weight_unit} Qty:${lot.quantity_unit}\t${stores.length}`
        )

        for (const s of stores) {
            lines.push(
                `\t\tStore #${s.id}\t${s.quantity} ${lot.quantity_unit} / ${s.weight} ${lot.weight_unit}\t` +
                `${s.is_active ? 'Active' : 'Inactive'}`
            )
        }
    }

    await clipboardStore.copyToClipboard(lines.join('\n'))
}

// ---------- LOT handlers ----------
const openDetailDialog = (lot: Lot) => {
    selectedLot.value = lot
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedLot.value = null }, 300)
}

const handleEditLot = (lot: Lot) => {
    selectedLot.value = lot
    editDialogOpen.value = true
}

const handleEditLotFromDetail = (lot: Lot) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedLot.value = lot
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteLot = (lot: Lot) => {
    selectedLot.value = lot
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedLot.value) return
    const success = await lotsStore.deleteLot(selectedLot.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchLots()
    }
}

const handleLotCreated = async () => {
    createDialogOpen.value = false
    await fetchLots()
}

const handleLotUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchLots()
}

const handleTransferLot = (lot: Lot) => {
    selectedLot.value = lot
    lotTransferDialogOpen.value = true
}

const handleLotTransferCreated = async () => {
    lotTransferDialogOpen.value = false
    await fetchLots()
}

// ---------- STORE handlers ----------
const handleAddStore = (lot: Lot) => {
    selectedLot.value = lot
    storeCreateDialogOpen.value = true
}

const handleStoreCreated = async () => {
    storeCreateDialogOpen.value = false
    await fetchLots()
}

const openStoreDetailDialog = (store: Store) => {
    selectedStore.value = store
    storeDetailDialogOpen.value = true
}

const closeStoreDetailDialog = () => {
    storeDetailDialogOpen.value = false
    setTimeout(() => { selectedStore.value = null }, 300)
}

const handleEditStore = (store: Store) => {
    selectedStore.value = store
    storeEditDialogOpen.value = true
}

const handleEditStoreFromDetail = (store: Store) => {
    storeDetailDialogOpen.value = false
    setTimeout(() => {
        selectedStore.value = store
        storeEditDialogOpen.value = true
    }, 300)
}

const handleDeleteStore = (store: Store) => {
    selectedStore.value = store
    storeDeleteDialogOpen.value = true
}

const handleAdjustStore = (store: Store) => {
    selectedStore.value = store
    adjustDialogOpen.value = true
}

const handleAdjustmentCreated = async () => {
    adjustDialogOpen.value = false
    await fetchLots()
}

const handleTransferStore = (store: Store) => {
    selectedStore.value = store
    transferDialogOpen.value = true
}

const handleTransferCreated = async () => {
    transferDialogOpen.value = false
    await fetchLots()
}

const confirmStoreDelete = async () => {
    if (!selectedStore.value) return
    const success = await storesStore.deleteStore(selectedStore.value.id)
    if (success) {
        storeDeleteDialogOpen.value = false
        await fetchLots()
    }
}

const handleStoreUpdated = async () => {
    storeEditDialogOpen.value = false
    storeDetailDialogOpen.value = false
    await fetchLots()
}

onMounted(() => { fetchLots() })
</script>