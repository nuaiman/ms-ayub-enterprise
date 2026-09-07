<!-- src/components/features/godowns/GodownList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Godowns</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredGodowns.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search godowns..."
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
                <button v-if="canManageGodowns" @click="createDialogOpen = true"
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
                    <!-- Name - 3 -->
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('name')">
                        <span class="flex items-center gap-1">
                            Name
                            <svg v-if="sortField === 'name'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Phone - 2 -->
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('phone')">
                        <span class="flex items-center gap-1">
                            Phone
                            <svg v-if="sortField === 'phone'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Notes - 3 -->
                    <div class="col-span-3">Notes</div>

                    <!-- Status - 2 -->
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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

                    <!-- Monthly Rent - 1 -->
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('monthly_rent')">
                        <span class="flex items-center gap-1">
                            Rent
                            <svg v-if="sortField === 'monthly_rent'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Actions - 1 -->
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
                        <p class="text-sm text-(--color-text-secondary)">Loading godowns...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredGodowns.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No godowns found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new godown to get started' }}</p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <GodownRow v-for="godown in filteredGodowns" :key="godown.id" :godown="godown"
                        @view="openDetailDialog" @edit="handleEditGodown" @delete="handleDeleteGodown"
                        @toggle-active="handleToggleActive" @updated="fetchGodowns" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredGodowns.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredGodowns.length }} of {{
                godownsStore.godowns.length }} godowns</p>
        </div>

        <!-- Dialogs -->
        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Godown</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new godown to the system</p>
            </div>
            <GodownForm mode="create" @godown-created="handleGodownCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <GodownDetail v-if="selectedGodown" :godown="selectedGodown" @close="closeDetailDialog"
                @edit="handleEditGodownFromDetail" @updated="fetchGodowns" />
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Godown</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update godown information</p>
            </div>
            <GodownForm v-if="selectedGodown" mode="edit" :godown="selectedGodown" @godown-updated="handleGodownUpdated"
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Godown</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete the godown "<span class="font-medium text-(--color-text-primary)">{{
                    selectedGodown?.name }}</span>"?
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
import { useGodownsStore } from '@/stores/godowns'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Godown, GodownSortField, SortDirection } from '@/types/godown'
import GodownRow from './GodownRow.vue'
import GodownForm from './GodownForm.vue'
import GodownDetail from './GodownDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'
import { formatCurrency } from '@/utils/currency'

const godownsStore = useGodownsStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const showActiveOnly = ref(false)
const sortField = ref<GodownSortField>('name')
const sortDirection = ref<SortDirection>('asc')

// Dialogs
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedGodown = ref<Godown | null>(null)

const canManageGodowns = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const filteredGodowns = computed(() => {
    let result = [...godownsStore.godowns]

    // Search
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(g =>
            g.name.toLowerCase().includes(query) ||
            (g.phone && g.phone.toLowerCase().includes(query)) ||
            (g.notes && g.notes.toLowerCase().includes(query))
        )
    }

    // Filter: Active only
    if (showActiveOnly.value) {
        result = result.filter(g => g.is_active)
    }

    // Sort
    result.sort((a, b) => {
        let comparison = 0
        switch (sortField.value) {
            case 'name':
                comparison = a.name.localeCompare(b.name)
                break
            case 'phone':
                comparison = (a.phone || '').localeCompare(b.phone || '')
                break
            case 'is_active':
                comparison = (a.is_active === b.is_active) ? 0 : a.is_active ? -1 : 1
                break
            case 'monthly_rent':
                comparison = a.monthly_rent - b.monthly_rent
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

const fetchGodowns = async () => {
    loading.value = true
    try {
        await godownsStore.fetchGodowns()
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

const toggleSort = (field: GodownSortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Name\tPhone\tNotes\tStatus\tMonthly Rent'
    const rows = filteredGodowns.value.map(g => {
        return `${g.name}\t${g.phone || ''}\t${g.notes || ''}\t${g.is_active ? 'Active' : 'Inactive'}\t${g.monthly_rent.toFixed(2)}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const handleToggleActive = async (godown: Godown) => {
    const success = await godownsStore.toggleGodownActive(godown.id)
    if (success) {
        await fetchGodowns()
    }
}

// Dialog handlers
const openDetailDialog = (godown: Godown) => {
    selectedGodown.value = godown
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedGodown.value = null }, 300)
}

const handleEditGodown = (godown: Godown) => {
    selectedGodown.value = godown
    editDialogOpen.value = true
}

const handleEditGodownFromDetail = (godown: Godown) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedGodown.value = godown
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteGodown = (godown: Godown) => {
    selectedGodown.value = godown
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedGodown.value) return
    const success = await godownsStore.deleteGodown(selectedGodown.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchGodowns()
    }
}

const handleGodownCreated = async () => {
    createDialogOpen.value = false
    await fetchGodowns()
}

const handleGodownUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchGodowns()
}

onMounted(() => {
    fetchGodowns()
})
</script>