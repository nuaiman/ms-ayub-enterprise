<!-- src/components/features/godownStoreBills/GodownStoreBillList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Godown Store Bills</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ godownStoreBillsStore.filteredStoreBills.length }}
                </span>
                <span v-if="godownStoreBillsStore.totalOutstanding > 0"
                    class="text-sm text-(--color-red) bg-(--color-red)/10 px-2 py-0.5 rounded-md">
                    Outstanding: {{ formatCurrency(godownStoreBillsStore.totalOutstanding) }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
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

                <!-- Month Filter -->
                <MonthFilter v-model="monthFilter" :items="availableMonths" />

                <!-- Refresh -->
                <button @click="refreshData"
                    class="h-9 w-9 flex items-center justify-center border border-(--color-border) rounded-lg text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors relative shrink-0"
                    title="Refresh data">
                    <svg class="w-4 h-4" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor"
                        viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
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
                        @click="toggleSort('lot_name')">
                        <span class="flex items-center gap-1">
                            Lot / Item
                            <svg v-if="sortField === 'lot_name'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('godown_name')">
                        <span class="flex items-center gap-1">
                            Godown
                            <svg v-if="sortField === 'godown_name'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2">Customer</div>
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('monthly_bill')">
                        <span class="flex items-center gap-1">
                            Monthly
                            <svg v-if="sortField === 'monthly_bill'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('outstanding')">
                        <span class="flex items-center gap-1">
                            Outstanding
                            <svg v-if="sortField === 'outstanding'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1">Status</div>
                    <div class="col-span-2 flex items-center justify-end">Actions</div>
                </div>

                <!-- Loading -->
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

                <!-- Empty -->
                <div v-else-if="filteredStoreBills.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No active stores with billing</p>
                        <p class="text-xs text-(--color-text-secondary)">Stores with godown_cut > 0 and inventory will
                            appear here</p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <GodownStoreBillRow v-for="store in filteredStoreBills" :key="store.id" :store="store"
                        @view="openDetailDialog" @pay="openPayDialog" @end-billing="openEndBillingDialog"
                        @reactivate="openReactivateDialog" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredStoreBills.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredStoreBills.length }} stores</p>
            <div class="flex items-center gap-4 text-xs text-(--color-text-secondary)">
                <span>Total Monthly: {{ formatCurrency(godownStoreBillsStore.totalMonthlyBill) }}</span>
                <span>Total Paid: {{ formatCurrency(godownStoreBillsStore.totalPaid) }}</span>
                <span class="font-semibold text-(--color-red)">Outstanding: {{
                    formatCurrency(godownStoreBillsStore.totalOutstanding) }}</span>
            </div>
        </div>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <GodownStoreBillDetail v-if="selectedStore" :store="selectedStore" @close="closeDetailDialog"
                @updated="fetchData" />
        </BaseDialog>

        <!-- Pay Dialog -->
        <BaseDialog v-model="payDialogOpen" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-green)/10 text-(--color-green) flex items-center justify-center shrink-0">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Record Payment</h2>
                        <p class="text-xs text-(--color-text-secondary)">Update payment for {{ selectedStore?.lot_name
                        }}</p>
                    </div>
                </div>

                <div class="space-y-3">
                    <div class="grid grid-cols-2 gap-3 text-sm bg-(--color-muted-bg)/30 p-3 rounded-lg">
                        <div>
                            <p class="text-xs text-(--color-text-secondary)">Monthly Bill</p>
                            <p class="font-medium text-(--color-blue)">{{ formatCurrency(selectedStore?.monthly_bill) }}
                            </p>
                        </div>
                        <div>
                            <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                            <p class="font-medium text-(--color-red)">{{ formatCurrency(selectedStore?.outstanding) }}
                            </p>
                        </div>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Payment Amount <span class="text-(--color-red)">*</span>
                        </label>
                        <div class="relative">
                            <span
                                class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                            <input v-model.number="payForm.amount" type="number" step="0.01" min="0"
                                :max="selectedStore?.outstanding" placeholder="0.00" required
                                class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        </div>
                        <p class="text-xs text-(--color-text-secondary) mt-1">Max: {{
                            formatCurrency(selectedStore?.outstanding) }}</p>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Paid Through Date
                        </label>
                        <input v-model="payForm.paid_through" type="date"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        <p class="text-xs text-(--color-text-secondary) mt-1">Date up to which payment is made</p>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Payment Method
                        </label>
                        <select v-model="payForm.method"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                            <option value="cash">Cash</option>
                            <option value="bank_transfer">Bank Transfer</option>
                            <option value="check">Check</option>
                            <option value="mobile_banking">Mobile Banking</option>
                        </select>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Reference Number
                        </label>
                        <input v-model="payForm.reference" type="text" placeholder="Enter reference number"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Notes
                        </label>
                        <textarea v-model="payForm.notes" rows="2" placeholder="Enter payment notes"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                    </div>
                </div>
            </div>

            <template #actions>
                <button @click="payDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmPay" :disabled="!payForm.amount || payForm.amount <= 0"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-green) text-white rounded-lg hover:opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    Confirm Payment
                </button>
            </template>
        </BaseDialog>

        <!-- End Billing Dialog -->
        <BaseDialog v-model="endBillingDialogOpen" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-yellow)/10 text-(--color-yellow) flex items-center justify-center shrink-0">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">End Billing</h2>
                        <p class="text-xs text-(--color-text-secondary)">Stop billing for this store</p>
                    </div>
                </div>

                <div class="bg-(--color-yellow)/5 border border-(--color-yellow)/20 rounded-lg p-4">
                    <p class="text-sm text-(--color-text-secondary)">
                        Are you sure you want to end billing for <span
                            class="font-medium text-(--color-text-primary)">{{
                                selectedStore?.lot_name }}</span>?
                    </p>
                    <p class="text-xs text-(--color-text-secondary) mt-2">
                        This will stop generating bills for this store. You can reactivate it later.
                    </p>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Billing End Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="endBillingForm.date" type="date" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="endBillingForm.notes" rows="2" placeholder="Reason for ending billing"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>
            </div>

            <template #actions>
                <button @click="endBillingDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmEndBilling"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-yellow) text-white rounded-lg hover:opacity-90 transition-colors">
                    End Billing
                </button>
            </template>
        </BaseDialog>

        <!-- Reactivate Dialog -->
        <BaseDialog v-model="reactivateDialogOpen" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center shrink-0">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Reactivate Billing</h2>
                        <p class="text-xs text-(--color-text-secondary)">Resume billing for this store</p>
                    </div>
                </div>

                <div class="bg-(--color-blue)/5 border border-(--color-blue)/20 rounded-lg p-4">
                    <p class="text-sm text-(--color-text-secondary)">
                        Are you sure you want to reactivate billing for <span
                            class="font-medium text-(--color-text-primary)">{{
                                selectedStore?.lot_name }}</span>?
                    </p>
                    <p class="text-xs text-(--color-text-secondary) mt-2">
                        This will resume generating bills for this store from today.
                    </p>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        New Billing Start Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="reactivateForm.date" type="date" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="reactivateForm.notes" rows="2" placeholder="Reason for reactivation"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>
            </div>

            <template #actions>
                <button @click="reactivateDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmReactivate"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-colors">
                    Reactivate
                </button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useGodownStoreBillsStore } from '@/stores/godownStoreBills'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useGodownsStore } from '@/stores/godowns'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import type { GodownStoreBillStore, GodownStoreBillSortField } from '@/types/godownStoreBill'
import GodownStoreBillRow from './GodownStoreBillRow.vue'
import GodownStoreBillDetail from './GodownStoreBillDetail.vue'
import MonthFilter from '@/components/ui/MonthFilter.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'

const godownStoreBillsStore = useGodownStoreBillsStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const godownsStore = useGodownsStore()
const itemsStore = useItemsStore()
const customersStore = useCustomersStore()

const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const monthFilter = ref('')
const sortField = ref<GodownStoreBillSortField>('lot_name')
const sortDirection = ref<'asc' | 'desc'>('asc')

// Dialogs
const detailDialogOpen = ref(false)
const payDialogOpen = ref(false)
const endBillingDialogOpen = ref(false)
const reactivateDialogOpen = ref(false)
const selectedStore = ref<GodownStoreBillStore | null>(null)

// Form states
const payForm = ref({
    amount: 0,
    paid_through: '',
    method: 'cash',
    reference: '',
    notes: '',
})

const endBillingForm = ref({
    date: '',
    notes: '',
})

const reactivateForm = ref({
    date: '',
    notes: '',
})

const availableMonths = computed(() => godownStoreBillsStore.availableMonths)
const filteredStoreBills = computed(() => godownStoreBillsStore.filteredStoreBills)

const fetchData = async () => {
    loading.value = true
    try {
        await Promise.all([
            storesStore.fetchStores(),
            lotsStore.fetchLots(),
            godownsStore.fetchGodowns(),
            itemsStore.fetchItems(),
            customersStore.fetchCustomers()
        ])
    } finally {
        loading.value = false
    }
}

const refreshData = async () => {
    refreshing.value = true
    try {
        await fetchData()
    } finally {
        refreshing.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    searchQuery.value = target.value
    godownStoreBillsStore.setSearchQuery(searchQuery.value)
}

const clearSearch = () => {
    searchQuery.value = ''
    godownStoreBillsStore.clearSearch()
}

const toggleSort = (field: GodownStoreBillSortField) => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
    godownStoreBillsStore.setSort(field)
}

// Dialog handlers
const openDetailDialog = (store: GodownStoreBillStore) => {
    selectedStore.value = store
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedStore.value = null }, 300)
}

const openPayDialog = (store: GodownStoreBillStore) => {
    selectedStore.value = store
    const today = new Date().toISOString().slice(0, 10)
    payForm.value = {
        amount: store.outstanding || 0,
        paid_through: today,
        method: 'cash',
        reference: '',
        notes: '',
    }
    payDialogOpen.value = true
}

const confirmPay = async () => {
    if (!selectedStore.value) return
    if (payForm.value.amount <= 0) {
        push.error('Payment amount must be greater than 0')
        return
    }
    if (payForm.value.amount > (selectedStore.value.outstanding || 0)) {
        push.error('Payment amount exceeds outstanding balance')
        return
    }

    try {
        const newTotalPaid = (selectedStore.value.last_paid_amount || 0) + payForm.value.amount

        // Format date for backend
        const paidThrough = payForm.value.paid_through
            ? formatDateForBackend(payForm.value.paid_through)
            : null

        await storesStore.updateStore(selectedStore.value.id, {
            last_paid_amount: newTotalPaid,
            last_paid_through: paidThrough,
            notes: payForm.value.notes || null,
        })

        push.success('Payment recorded successfully!')
        payDialogOpen.value = false
        await fetchData()
    } catch (error) {
        console.error('Error recording payment:', error)
        push.error('Failed to record payment')
    }
}

const openEndBillingDialog = (store: GodownStoreBillStore) => {
    selectedStore.value = store
    const today = new Date().toISOString().slice(0, 10)
    endBillingForm.value = {
        date: today,
        notes: '',
    }
    endBillingDialogOpen.value = true
}

const confirmEndBilling = async () => {
    if (!selectedStore.value) return
    if (!endBillingForm.value.date) {
        push.error('Please select an end date')
        return
    }

    try {
        // Format date for backend
        const endDate = formatDateForBackend(endBillingForm.value.date)

        await storesStore.updateStore(selectedStore.value.id, {
            billing_end: endDate,
            notes: endBillingForm.value.notes || null,
        })

        push.success('Billing ended successfully!')
        endBillingDialogOpen.value = false
        await fetchData()
    } catch (error) {
        console.error('Error ending billing:', error)
        push.error('Failed to end billing')
    }
}

const openReactivateDialog = (store: GodownStoreBillStore) => {
    selectedStore.value = store
    const today = new Date().toISOString().slice(0, 10)
    reactivateForm.value = {
        date: today,
        notes: '',
    }
    reactivateDialogOpen.value = true
}

const confirmReactivate = async () => {
    if (!selectedStore.value) return
    if (!reactivateForm.value.date) {
        push.error('Please select a start date')
        return
    }

    try {
        // Format date for backend
        const startDate = formatDateForBackend(reactivateForm.value.date)

        await storesStore.updateStore(selectedStore.value.id, {
            billing_start: startDate,
            billing_end: null,
            notes: reactivateForm.value.notes || null,
        })

        push.success('Billing reactivated successfully!')
        reactivateDialogOpen.value = false
        await fetchData()
    } catch (error) {
        console.error('Error reactivating billing:', error)
        push.error('Failed to reactivate billing')
    }
}

onMounted(() => {
    fetchData()
})
</script>