<!-- src/components/features/invoices/InvoiceList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Invoices</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredInvoices.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="invoicesStore.searchQuery" @input="handleSearch" type="text"
                        placeholder="Search invoices..."
                        class="w-full sm:w-56 pl-9 pr-8 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--color-text-secondary)"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <button v-if="invoicesStore.searchQuery" @click="clearSearch" type="button"
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
                    <div class="col-span-4">Entity</div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('subtotal')">
                        <span class="flex items-center gap-1">
                            Subtotal
                            <svg v-if="invoicesStore.sortField === 'subtotal'" class="w-3 h-3"
                                :class="{ 'rotate-180': invoicesStore.sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('discount_amount')">
                        <span class="flex items-center gap-1">
                            Discount
                            <svg v-if="invoicesStore.sortField === 'discount_amount'" class="w-3 h-3"
                                :class="{ 'rotate-180': invoicesStore.sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('total')">
                        <span class="flex items-center gap-1">
                            Total
                            <svg v-if="invoicesStore.sortField === 'total'" class="w-3 h-3"
                                :class="{ 'rotate-180': invoicesStore.sortDirection === 'desc' }" fill="currentColor"
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
                        <p class="text-sm text-(--color-text-secondary)">Loading invoices...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredInvoices.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No invoices found</p>
                        <p class="text-xs text-(--color-text-secondary)">
                            {{
                                invoicesStore.searchQuery
                                    ? "Try adjusting your search"
                                    : "Create a new invoice to get started"
                            }}
                        </p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <InvoiceRow v-for="inv in filteredInvoices" :key="inv.id" :invoice="inv" @view="openDetailDialog"
                        @print="handlePrintRow" @edit="handleEditInvoice" @delete="handleDeleteInvoice"
                        @updated="fetchInvoices" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredInvoices.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">
                Showing {{ filteredInvoices.length }} of {{ invoicesStore.totalInvoices }} invoices
            </p>
            <p class="text-xs text-(--color-text-secondary)">
                Total: {{ formatCurrency(invoicesStore.totalAmount) }}
            </p>
        </div>

        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Invoice</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Select an entity, pick bills, and apply discounts
                </p>
            </div>
            <InvoiceForm mode="create" @invoice-created="handleInvoiceCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <InvoiceDetail v-if="selectedInvoice" :invoice="selectedInvoice" @close="closeDetailDialog"
                @edit="handleEditInvoiceFromDetail" @print="handlePrintInvoice" @updated="fetchInvoices" />
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Invoice</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update discounts or notes</p>
            </div>
            <InvoiceForm v-if="selectedInvoice" mode="edit" :invoice="selectedInvoice"
                @invoice-updated="handleInvoiceUpdated" @cancel="editDialogOpen = false" />
        </BaseDialog>

        <!-- Preview / Print Dialog -->
        <BaseDialog v-model="previewDialogOpen" max-width="full">
            <InvoicePreview v-if="previewInvoice" :invoice="previewInvoice" @back="previewDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Invoice</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete invoice #<span class="font-medium text-(--color-text-primary)">{{
                    selectedInvoice?.id }}</span>?
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
import { useInvoicesStore } from '@/stores/invoices'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useBrokersStore } from '@/stores/brokers'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Invoice, InvoiceSortField } from '@/types/invoice'
import type { InvoiceDetail as InvoiceDetailType } from '@/types/invoice'
import InvoiceRow from './InvoiceRow.vue'
import InvoiceForm from './InvoiceForm.vue'
import InvoiceDetail from './InvoiceDetail.vue'
import InvoicePreview from './InvoicePreview.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'

const invoicesStore = useInvoicesStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const brokersStore = useBrokersStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const previewDialogOpen = ref(false)
const deleteDialogOpen = ref(false)

const selectedInvoice = ref<InvoiceDetailType | null>(null)
const previewInvoice = ref<(Invoice & Partial<InvoiceDetailType>) | null>(null)

const canManage = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const filteredInvoices = computed(() => invoicesStore.filteredInvoices)

const fetchInvoices = async () => {
    loading.value = true
    try {
        await Promise.all([
            invoicesStore.fetchInvoices(),
            customersStore.fetchCustomers(),
            majhisStore.fetchMajhis(),
            godownsStore.fetchGodowns(),
            brokersStore.fetchBrokers(),
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    invoicesStore.setSearchQuery(target.value)
}

const clearSearch = () => {
    invoicesStore.clearSearch()
}

const toggleSort = (field: InvoiceSortField) => {
    invoicesStore.setSort(field)
}

const handleCopyToClipboard = async () => {
    const headers = 'ID\tEntity Type\tEntity\tSubtotal\tDiscount\tTotal'
    const rows = filteredInvoices.value.map((i) => {
        let entity = `#${i.entity_id}`
        if (i.entity_type === 'customer') {
            const c = customersStore.getCustomerById(i.entity_id)
            entity = c ? (c.company_name || c.contact_person || `Customer #${c.id}`) : `Customer #${i.entity_id}`
        } else if (i.entity_type === 'majhi') {
            entity = majhisStore.getMajhiName(i.entity_id)
        } else if (i.entity_type === 'godown') {
            entity = godownsStore.getGodownName(i.entity_id)
        } else if (i.entity_type === 'broker') {
            entity = brokersStore.getBrokerName(i.entity_id)
        }
        return `${i.id}\t${i.entity_type}\t${entity}\t${i.subtotal.toFixed(2)}\t${i.discount_amount.toFixed(2)}\t${i.total.toFixed(2)}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

// =========================================================================
// DETAIL DIALOG
// =========================================================================

const openDetailDialog = async (invoice: Invoice) => {
    const detail = await invoicesStore.fetchInvoiceDetail(invoice.id)
    if (!detail) return
    selectedInvoice.value = detail
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedInvoice.value = null }, 300)
}

// =========================================================================
// EDIT DIALOG
// =========================================================================

const handleEditInvoice = async (invoice: Invoice) => {
    const detail = await invoicesStore.fetchInvoiceDetail(invoice.id)
    if (!detail) return
    selectedInvoice.value = detail
    editDialogOpen.value = true
}

const handleEditInvoiceFromDetail = (invoice: InvoiceDetailType) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedInvoice.value = invoice
        editDialogOpen.value = true
    }, 300)
}

// =========================================================================
// PRINT / PREVIEW
// =========================================================================

const handlePrintInvoice = (invoice: InvoiceDetailType) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        previewInvoice.value = invoice
        previewDialogOpen.value = true
    }, 300)
}

const handlePrintRow = async (invoice: Invoice) => {
    const detail = await invoicesStore.fetchInvoiceDetail(invoice.id)
    if (!detail) return
    previewInvoice.value = detail
    previewDialogOpen.value = true
}

// =========================================================================
// DELETE DIALOG
// =========================================================================

const handleDeleteInvoice = (invoice: Invoice) => {
    selectedInvoice.value = {
        ...invoice,
        items: [],
        discounts: [],
    }
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedInvoice.value) return
    const ok = await invoicesStore.deleteInvoice(selectedInvoice.value.id)
    if (ok) {
        deleteDialogOpen.value = false
        await fetchInvoices()
    }
}

// =========================================================================
// CREATE / UPDATE CALLBACKS
// =========================================================================

const handleInvoiceCreated = async (created: InvoiceDetailType) => {
    createDialogOpen.value = false
    await fetchInvoices()
    previewInvoice.value = created
    previewDialogOpen.value = true
}

const handleInvoiceUpdated = async (updated: InvoiceDetailType) => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchInvoices()
    previewInvoice.value = updated
    previewDialogOpen.value = true
}

onMounted(() => {
    fetchInvoices()
})
</script>