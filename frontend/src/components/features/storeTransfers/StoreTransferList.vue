<!-- src/components/features/storeTransfers/StoreTransferList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Store Transfers</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredTransfers.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search transfers..."
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
                <button v-if="canManage" @click="openCreateDialog"
                    class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    <span class="hidden sm:inline">Transfer Store</span>
                </button>
            </div>
        </div>

        <!-- Table -->
        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('store_id')">
                        <span class="flex items-center gap-1">
                            Store
                            <svg v-if="sortField === 'store_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-4">From → To</div>
                    <div class="col-span-2">Snapshot</div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('transferred_at')">
                        <span class="flex items-center gap-1">
                            By / When
                            <svg v-if="sortField === 'transferred_at'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 text-right">Actions</div>
                </div>

                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading transfers...</p>
                    </div>
                </div>

                <div v-else-if="filteredTransfers.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M14 5l7 7m0 0l-7 7m7-7H3" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No transfers found</p>
                        <p class="text-xs text-(--color-text-secondary)">
                            {{ searchQuery ? 'Try adjusting your search' : 'Transfer a store to get started' }}
                        </p>
                    </div>
                </div>

                <div v-else>
                    <StoreTransferRow v-for="t in filteredTransfers" :key="t.id" :transfer="t"
                        :deleting="deletingId === t.id" :is-latest="latestTransferIdByStore[t.store_id] === t.id"
                        @delete="handleDelete" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredTransfers.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">
                Showing {{ filteredTransfers.length }} of {{ storeTransfersStore.totalTransfers }} transfers
            </p>
        </div>

        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Transfer Store</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Move a store from its current godown to another
                </p>
            </div>
            <StoreTransferForm @transfer-created="handleTransferCreated" @cancel="createDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Transfer</h2>
                    <p class="text-xs text-(--color-text-secondary)">
                        This will move the store back to its previous godown
                    </p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete this transfer?
                <span v-if="selectedTransfer" class="block mt-2 text-(--color-text-primary) font-medium">
                    {{ godownsStore.getGodownName(selectedTransfer.from_godown_id) }}
                    →
                    {{ godownsStore.getGodownName(selectedTransfer.to_godown_id) }}
                </span>
            </p>
            <template #actions>
                <button @click="deleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmDelete" :disabled="!!deletingId"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    {{ deletingId ? 'Deleting...' : 'Delete' }}
                </button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useStoreTransfersStore } from '@/stores/storeTransfers'
import { useStoresStore } from '@/stores/stores'
import { useGodownsStore } from '@/stores/godowns'
import { useLotsStore } from '@/stores/lots'
import { useUsersStore } from '@/stores/users'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { StoreTransfer } from '@/types/storeTransfer'
import StoreTransferRow from './StoreTransferRow.vue'
import StoreTransferForm from './StoreTransferForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'

const storeTransfersStore = useStoreTransfersStore()
const storesStore = useStoresStore()
const godownsStore = useGodownsStore()
const lotsStore = useLotsStore()
const usersStore = useUsersStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const createDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const selectedTransfer = ref<StoreTransfer | null>(null)
const deletingId = ref<number | null>(null)

const searchQuery = computed(() => storeTransfersStore.searchQuery)
const sortField = computed(() => storeTransfersStore.sortField)
const sortDirection = computed(() => storeTransfersStore.sortDirection)
const filteredTransfers = computed(() => storeTransfersStore.filteredTransfers)

const canManage = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

// Map of store_id -> latest transfer id, for enabling delete only on the latest.
const latestTransferIdByStore = computed(() => {
    const map: Record<number, number> = {}
    for (const t of storeTransfersStore.transfers) {
        const existing = map[t.store_id]
        if (existing === undefined) {
            map[t.store_id] = t.id
            continue
        }
        // We want the highest id (most recently created) per store.
        if (t.id > existing) map[t.store_id] = t.id
    }
    return map
})

const fetchTransfers = async () => {
    loading.value = true
    try {
        await Promise.all([
            storeTransfersStore.fetchAllTransfers(),
            storesStore.fetchStores(),
            godownsStore.fetchGodowns(),
            lotsStore.fetchLots(),
            usersStore.fetchUsers(),
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    storeTransfersStore.setSearchQuery(target.value)
}

const clearSearch = () => {
    storeTransfersStore.clearSearch()
}

const toggleSort = (field: Parameters<typeof storeTransfersStore.setSort>[0]) => {
    storeTransfersStore.setSort(field)
}

const openCreateDialog = () => {
    createDialogOpen.value = true
}

const handleCopyToClipboard = async () => {
    const headers = 'Store ID\tFrom Godown\tTo Godown\tWeight\tQuantity\tUser\tTransferred At'
    const rows = filteredTransfers.value.map(t => {
        return `${t.store_id}\t${godownsStore.getGodownName(t.from_godown_id)}\t${godownsStore.getGodownName(t.to_godown_id)}\t${t.weight_at_transfer.toFixed(2)}\t${t.quantity_at_transfer.toFixed(2)}\t${usersStore.getUserName(t.user_id)}\t${new Date(t.transferred_at).toLocaleString()}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const handleTransferCreated = async () => {
    createDialogOpen.value = false
    await fetchTransfers()
}

const handleDelete = (transfer: StoreTransfer) => {
    selectedTransfer.value = transfer
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedTransfer.value) return
    const id = selectedTransfer.value.id
    deletingId.value = id
    try {
        const success = await storeTransfersStore.deleteTransfer(id)
        if (success) {
            deleteDialogOpen.value = false
            selectedTransfer.value = null
            await storesStore.fetchStores()
        }
    } finally {
        deletingId.value = null
    }
}

onMounted(() => {
    fetchTransfers()
})
</script>