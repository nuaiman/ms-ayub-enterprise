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
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search lots..."
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

                <!-- Filter: Active Only -->
                <button @click="showActiveOnly = !showActiveOnly"
                    class="h-9 px-3 flex items-center gap-1.5 rounded-lg text-sm border border-(--color-border) transition-colors"
                    :class="showActiveOnly ? 'bg-(--color-blue)/10 border-(--color-blue) text-(--color-blue)' : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)'">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Active Only
                </button>

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
                <button v-if="canManageLots" @click="createDialogOpen = true"
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
                            Product / Lot
                            <svg v-if="sortField === 'customer_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-3">Customer</div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('customer_charge_type')">
                        <span class="flex items-center gap-1">
                            Charge
                            <svg v-if="sortField === 'customer_charge_type'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2">Majhi</div>
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

                <!-- Loading -->
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

                <!-- Empty -->
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

                <!-- Rows -->
                <div v-else>
                    <LotRow v-for="lot in filteredLots" :key="lot.id" :lot="lot" @view="openDetailDialog"
                        @edit="handleEditLot" @readd="handleReaddLot" @delete="handleDeleteLot"
                        @toggle-active="handleToggleActive" @updated="fetchLots" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredLots.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredLots.length }} of {{
                lotsStore.lots.length }} lots</p>
        </div>

        <!-- Dialogs -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Lot</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new lot to the system</p>
            </div>
            <LotForm mode="create" @lot-created="handleLotCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="readdDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Re-add Lot</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Create a new lot based on this one</p>
            </div>
            <LotForm v-if="readdLot" mode="create" :prefill="readdLot" @lot-created="handleLotCreated"
                @cancel="readdDialogOpen = false" />
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
                Are you sure you want to delete the lot for "<span class="font-medium text-(--color-text-primary)">{{
                    getLotDisplayName(selectedLot) }}</span>" (Lot #{{ selectedLot?.lot_number }})?
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
import { useLotsStore } from '@/stores/lots'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useStoresStore } from '@/stores/stores'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Lot, LotSortField, SortDirection } from '@/types/lot'
import LotRow from './LotRow.vue'
import LotForm from './LotForm.vue'
import LotDetail from './LotDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const lotsStore = useLotsStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const storesStore = useStoresStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const showActiveOnly = ref(false)
const sortField = ref<LotSortField>('lot_number')
const sortDirection = ref<SortDirection>('asc')

const createDialogOpen = ref(false)
const readdDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedLot = ref<Lot | null>(null)
const readdLot = ref<Lot | null>(null)

const canManageLots = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const getLotDisplayName = (lot?: Lot | null): string => {
    if (!lot) return 'Unknown'
    return lotsStore.getLotDisplayName(lot)
}

const filteredLots = computed(() => {
    let result = [...lotsStore.lots]

    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(l =>
            lotsStore.getLotDisplayName(l).toLowerCase().includes(query) ||
            (l.product_name && l.product_name.toLowerCase().includes(query)) ||
            (l.category && l.category.toLowerCase().includes(query)) ||
            String(l.lot_number).includes(query) ||
            l.customer_charge_type.toLowerCase().includes(query) ||
            (l.customer_id && customersStore.getCustomerName(l.customer_id).toLowerCase().includes(query)) ||
            (l.majhi_id && majhisStore.getMajhiName(l.majhi_id).toLowerCase().includes(query)) ||
            (l.notes && l.notes.toLowerCase().includes(query))
        )
    }

    if (showActiveOnly.value) {
        result = result.filter(l => l.is_active)
    }

    result.sort((a, b) => {
        let comparison = 0
        switch (sortField.value) {
            case 'customer_id':
                comparison = (a.customer_id || 0) - (b.customer_id || 0)
                break
            case 'lot_number':
                comparison = a.lot_number - b.lot_number
                break
            case 'customer_charge_type':
                comparison = a.customer_charge_type.localeCompare(b.customer_charge_type)
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

const fetchLots = async () => {
    loading.value = true
    try {
        await Promise.all([
            lotsStore.fetchLots(),
            customersStore.fetchCustomers(),
            majhisStore.fetchMajhis(),
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

const clearSearch = () => {
    searchQuery.value = ''
}

const toggleSort = (field: LotSortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Product\tLot #\tCustomer\tCharge Type\tMajhi\tStatus'
    const rows = filteredLots.value.map(l => {
        const productName = lotsStore.getLotDisplayName(l)
        const customerName = l.customer_id ? customersStore.getCustomerName(l.customer_id) : ''
        const majhiName = l.majhi_id ? majhisStore.getMajhiName(l.majhi_id) : ''
        return `${productName}\t${l.lot_number}\t${customerName}\t${l.customer_charge_type}\t${majhiName}\t${l.is_active ? 'Active' : 'Inactive'}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const handleToggleActive = async (lot: Lot) => {
    const success = await lotsStore.toggleLotActive(lot.id)
    if (success) {
        await fetchLots()
    }
}

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

const handleReaddLot = (lot: Lot) => {
    readdLot.value = lot
    readdDialogOpen.value = true
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
    readdDialogOpen.value = false
    readdLot.value = null
    await fetchLots()
}

const handleLotUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchLots()
}

onMounted(() => {
    fetchLots()
})
</script>