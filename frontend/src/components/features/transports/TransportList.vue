<!-- src/components/features/transports/TransportList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Transports</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredTransports.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search transports..."
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
                <button v-if="canManageTransports" @click="createDialogOpen = true"
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
                <!-- Header Row - 12 columns -->
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <!-- Customer - 2 columns -->
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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

                    <!-- From - 2 columns -->
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

                    <!-- To - 2 columns -->
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

                    <!-- Vehicles - 1 column -->
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('vehicle_quantity')">
                        <span class="flex items-center gap-1">
                            Veh.
                            <svg v-if="sortField === 'vehicle_quantity'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Date - 2 columns -->
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('transport_date')">
                        <span class="flex items-center gap-1">
                            Date
                            <svg v-if="sortField === 'transport_date'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Commission - 1 column -->
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('office_commission_amount')">
                        <span class="flex items-center gap-1">
                            Comm.
                            <svg v-if="sortField === 'office_commission_amount'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Actions - 1 column, right aligned -->
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
                        <p class="text-sm text-(--color-text-secondary)">Loading transports...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredTransports.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No transports found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new transport to get started' }}</p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <TransportRow v-for="transport in filteredTransports" :key="transport.id" :transport="transport"
                        @view="openDetailDialog" @edit="handleEditTransport" @delete="handleDeleteTransport"
                        @manage-vehicles="handleManageVehicles" @updated="fetchTransports" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredTransports.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredTransports.length }} of {{
                transportsStore.transports.length }} transports</p>
            <div class="flex items-center gap-4 text-xs text-(--color-text-secondary)">
                <span>Total Commission: {{ formatCurrency(transportsStore.totalCommission) }}</span>
                <span class="text-(--color-blue)">Total Customer Paid: {{
                    formatCurrency(transportsStore.totalCustomerPaid || 0) }}</span>
            </div>
        </div>

        <!-- Dialogs -->
        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Transport</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new transport to the system</p>
            </div>
            <TransportForm mode="create" @transport-created="handleTransportCreated"
                @cancel="createDialogOpen = false" />
        </BaseDialog>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <TransportDetail v-if="selectedTransport" :transport="selectedTransport" @close="closeDetailDialog"
                @edit="handleEditTransportFromDetail" @manage-vehicles="handleManageVehiclesFromDetail"
                @updated="fetchTransports" />
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Transport</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update transport information</p>
            </div>
            <TransportForm v-if="selectedTransport" mode="edit" :transport="selectedTransport"
                @transport-updated="handleTransportUpdated" @cancel="editDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Transport</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete transport #<span class="font-medium text-(--color-text-primary)">{{
                    selectedTransport?.id }}</span>?
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
import { useTransportsStore } from '@/stores/transports'
import { useCustomersStore } from '@/stores/customers'
import { useUsersStore } from '@/stores/users'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Transport, TransportSortField, SortDirection } from '@/types/transport'
import TransportRow from './TransportRow.vue'
import TransportForm from './TransportForm.vue'
import TransportDetail from './TransportDetail.vue'
import MonthFilter from '@/components/ui/MonthFilter.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'
import { formatCurrency } from '@/utils/currency'

const transportsStore = useTransportsStore()
const customersStore = useCustomersStore()
const usersStore = useUsersStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const monthFilter = ref('')
const sortField = ref<TransportSortField>('transport_date')
const sortDirection = ref<SortDirection>('desc')

// Dialogs
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedTransport = ref<Transport | null>(null)

const canManageTransports = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

// Get unique months from transports
const availableMonths = computed(() => {
    const months = new Set<string>()
    transportsStore.transports.forEach(t => {
        const date = new Date(t.transport_date)
        const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
        months.add(month)
    })
    return Array.from(months).sort((a, b) => b.localeCompare(a))
})

const getCustomerName = (id: number | null): string => {
    if (!id) return '—'
    return customersStore.getCustomerName(id)
}

const filteredTransports = computed(() => {
    let result = [...transportsStore.transports]

    // Search
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(t =>
            t.from_location.toLowerCase().includes(query) ||
            (t.to_location && t.to_location.toLowerCase().includes(query)) ||
            (t.notes && t.notes.toLowerCase().includes(query)) ||
            String(t.vehicle_quantity).includes(query) ||
            (t.delivery_type && t.delivery_type.toLowerCase().includes(query)) ||
            String(t.office_commission_amount).includes(query) ||
            String(t.customer_total_paid || 0).includes(query) ||
            getCustomerName(t.customer_id).toLowerCase().includes(query) ||
            usersStore.getUserName(t.user_id).toLowerCase().includes(query)
        )
    }

    // Filter by month
    if (monthFilter.value) {
        result = result.filter(t => {
            const date = new Date(t.transport_date)
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
            case 'from_location':
                comparison = a.from_location.localeCompare(b.from_location)
                break
            case 'to_location':
                comparison = (a.to_location || '').localeCompare(b.to_location || '')
                break
            case 'vehicle_quantity':
                comparison = a.vehicle_quantity - b.vehicle_quantity
                break
            case 'delivery_type':
                comparison = (a.delivery_type || '').localeCompare(b.delivery_type || '')
                break
            case 'transport_date':
                comparison = new Date(a.transport_date).getTime() - new Date(b.transport_date).getTime()
                break
            case 'office_commission_amount':
                comparison = a.office_commission_amount - b.office_commission_amount
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

const fetchTransports = async () => {
    loading.value = true
    try {
        await Promise.all([
            transportsStore.fetchTransports(),
            customersStore.fetchCustomers(),
            usersStore.fetchUsers()
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

const toggleSort = (field: TransportSortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Customer\tFrom\tTo\tVehicles\tDate\tCommission'
    const rows = filteredTransports.value.map(t => {
        return `${getCustomerName(t.customer_id)}\t${t.from_location}\t${t.to_location || ''}\t${t.vehicle_quantity}\t${new Date(t.transport_date).toLocaleDateString()}\t${t.office_commission_amount.toFixed(2)}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

// Dialog handlers
const openDetailDialog = (transport: Transport) => {
    selectedTransport.value = transport
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedTransport.value = null }, 300)
}

const handleEditTransport = (transport: Transport) => {
    selectedTransport.value = transport
    editDialogOpen.value = true
}

const handleEditTransportFromDetail = (transport: Transport) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedTransport.value = transport
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteTransport = (transport: Transport) => {
    selectedTransport.value = transport
    deleteDialogOpen.value = true
}

const handleManageVehicles = (transport: Transport) => {
    selectedTransport.value = transport
    // Close detail dialog if open
    detailDialogOpen.value = false
    // This will be handled by the parent view
    push.info(`Manage vehicles for transport #${transport.id}`)
}

const handleManageVehiclesFromDetail = (transport: Transport) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        handleManageVehicles(transport)
    }, 300)
}

const confirmDelete = async () => {
    if (!selectedTransport.value) return
    const success = await transportsStore.deleteTransport(selectedTransport.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchTransports()
    }
}

const handleTransportCreated = async () => {
    createDialogOpen.value = false
    await fetchTransports()
}

const handleTransportUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchTransports()
}

// Auto-select current month
watch(() => transportsStore.transports, (newTransports) => {
    if (newTransports.length > 0 && !monthFilter.value) {
        const currentMonth = new Date().toISOString().slice(0, 7)
        const hasCurrentMonth = newTransports.some(t => {
            const date = new Date(t.transport_date)
            const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
            return month === currentMonth
        })
        if (hasCurrentMonth) {
            monthFilter.value = currentMonth
        }
    }
}, { immediate: true })

onMounted(() => {
    fetchTransports()
})
</script>