<!-- src/components/features/damages/DamageList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Damages</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredDamages.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search damages..."
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
                <button v-if="canManage" @click="createDialogOpen = true"
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
                    <div class="col-span-4 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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
                    <div class="col-span-3">Reason</div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('quantity')">
                        <span class="flex items-center gap-1">
                            Quantity / Weight
                            <svg v-if="sortField === 'quantity'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('amount')">
                        <span class="flex items-center gap-1">
                            Amount
                            <svg v-if="sortField === 'amount'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('damage_date')">
                        <span class="flex items-center gap-1">
                            Date
                            <svg v-if="sortField === 'damage_date'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 text-right">Actions</div>
                </div>

                <!-- Loading -->
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading damages...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredDamages.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No damages found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Record a damage to get started' }}</p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <DamageRow v-for="damage in filteredDamages" :key="damage.id" :damage="damage"
                        @view="openDetailDialog" @edit="handleEditDamage" @delete="handleDeleteDamage"
                        @updated="fetchDamages" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredDamages.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredDamages.length }} of {{ totalDamages }}
                damages</p>
            <p class="text-xs text-(--color-text-secondary)">Total: {{ formatCurrency(damagesStore.totalDamageAmount) }}
            </p>
        </div>

        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Record Damage</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Record a new damage against a store</p>
            </div>
            <DamageForm mode="create" @damage-created="handleDamageCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <DamageDetail v-if="selectedDamage" :damage="selectedDamage" @close="closeDetailDialog"
                @edit="handleEditDamageFromDetail" @updated="fetchDamages" />
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Damage</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update damage information</p>
            </div>
            <DamageForm v-if="selectedDamage" mode="edit" :damage="selectedDamage" @damage-updated="handleDamageUpdated"
                @cancel="editDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Damage</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete damage #<span class="font-medium text-(--color-text-primary)">{{
                    selectedDamage?.id }}</span>?
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
import { useDamagesStore } from '@/stores/damages'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Damage } from '@/types/damage'
import DamageRow from './DamageRow.vue'
import DamageForm from './DamageForm.vue'
import DamageDetail from './DamageDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'

const damagesStore = useDamagesStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedDamage = ref<Damage | null>(null)

const searchQuery = computed(() => damagesStore.searchQuery)
const sortField = computed(() => damagesStore.sortField)
const sortDirection = computed(() => damagesStore.sortDirection)
const filteredDamages = computed(() => damagesStore.filteredDamages)
const totalDamages = computed(() => damagesStore.totalDamages)

const canManage = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const fetchDamages = async () => {
    loading.value = true
    try {
        await Promise.all([
            damagesStore.fetchDamages(),
            storesStore.fetchStores(),
            lotsStore.fetchLots(),
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    damagesStore.setSearchQuery(target.value)
}

const clearSearch = () => {
    damagesStore.clearSearch()
}

const toggleSort = (field: Parameters<typeof damagesStore.setSort>[0]) => {
    damagesStore.setSort(field)
}

const handleCopyToClipboard = async () => {
    const headers = 'Store ID\tReason\tQuantity\tWeight\tAmount\tDate'
    const rows = filteredDamages.value.map(d =>
        `${d.store_id}\t${d.reason}\t${d.quantity.toFixed(2)} ${d.quantity_unit}\t${d.weight.toFixed(2)} ${d.weight_unit}\t${d.amount.toFixed(2)}\t${new Date(d.damage_date).toLocaleDateString()}`
    )
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const openDetailDialog = (damage: Damage) => {
    selectedDamage.value = damage
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedDamage.value = null }, 300)
}

const handleEditDamage = (damage: Damage) => {
    selectedDamage.value = damage
    editDialogOpen.value = true
}

const handleEditDamageFromDetail = (damage: Damage) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedDamage.value = damage
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteDamage = (damage: Damage) => {
    selectedDamage.value = damage
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedDamage.value) return
    const success = await damagesStore.deleteDamage(selectedDamage.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchDamages()
    }
}

const handleDamageCreated = async () => {
    createDialogOpen.value = false
    await fetchDamages()
}

const handleDamageUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchDamages()
}

onMounted(() => { fetchDamages() })
</script>