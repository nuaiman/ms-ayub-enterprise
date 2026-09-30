<!-- src/components/features/customerAdditionalBills/CustomerAdditionalBillList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Customer Additional Bills</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredBills.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search bills..."
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

                <!-- Status Filter -->
                <select v-model="statusFilter" @change="handleStatusFilterChange"
                    class="px-3 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                    <option value="">All Status</option>
                    <option value="unpaid">Unpaid</option>
                    <option value="paid">Paid</option>
                </select>

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
                <button @click="openCustomerPicker"
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
                        @click="toggleSort('customer_name')">
                        <span class="flex items-center gap-1">
                            Customer
                            <svg v-if="sortField === 'customer_name'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-4 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('description')">
                        <span class="flex items-center gap-1">
                            Description
                            <svg v-if="sortField === 'description'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('status')">
                        <span class="flex items-center gap-1">
                            Status
                            <svg v-if="sortField === 'status'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                </div>

                <!-- Loading -->
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading bills...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredBills.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No additional bills found</p>
                        <p class="text-xs text-(--color-text-secondary)">
                            {{ searchQuery
                                ? 'Try adjusting your search'
                                : 'Create a new additional charge to get started'
                            }}
                        </p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <CustomerAdditionalBillRow v-for="bill in filteredBills" :key="bill.id" :bill="bill"
                        @view="openDetailDialog" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredBills.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredBills.length }} bills</p>
            <div class="flex items-center gap-4 text-xs text-(--color-text-secondary)">
                <span>Total Amount: {{ formatCurrency(store.totalAmount) }}</span>
                <span class="text-(--color-yellow)">Unpaid: {{ formatCurrency(store.totalUnpaidAmount) }}</span>
            </div>
        </div>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <CustomerAdditionalBillDetail v-if="selectedBill" :bill="selectedBill" @close="detailDialogOpen = false"
                @pay="openPaymentDialog" @edit="openEditDialog" @delete="openDeleteDialog" />
        </BaseDialog>

        <!-- Payment Dialog -->
        <BaseDialog v-model="paymentDialogOpen" max-width="sm">
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
                        <p class="text-xs text-(--color-text-secondary)">For {{ selectedBill?.customer_name }}</p>
                    </div>
                </div>

                <div class="space-y-3">
                    <div class="grid grid-cols-2 gap-3 text-sm bg-(--color-muted-bg)/30 p-3 rounded-lg">
                        <div>
                            <p class="text-xs text-(--color-text-secondary)">Amount</p>
                            <p class="font-medium text-(--color-blue)">{{ formatCurrency(selectedBill?.amount) }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-(--color-text-secondary)">Already Paid</p>
                            <p class="font-medium text-(--color-green)">{{ formatCurrency(selectedBill?.paid_amount) }}
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
                            <input v-model.number="paymentForm.amount" type="number" step="0.01" min="0.01"
                                :max="maxPaymentAmount" placeholder="0.00" required
                                class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        </div>
                        <p class="text-xs text-(--color-text-secondary) mt-1">
                            Remaining: {{ formatCurrency(maxPaymentAmount) }}
                        </p>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Payment Date <span class="text-(--color-red)">*</span>
                        </label>
                        <input v-model="paymentForm.payment_date" type="date" required
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>
            </div>

            <template #actions>
                <button @click="paymentDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>
                <button @click="confirmPayment" :disabled="!paymentForm.amount || paymentForm.amount <= 0"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-green) text-white rounded-lg hover:opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    Confirm Payment
                </button>
            </template>
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Additional Charge</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update amount and description</p>
            </div>
            <AdditionalChargeForm v-if="selectedBill" mode="edit" :bill="selectedBill"
                @charge-updated="handleChargeUpdated" @cancel="editDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Additional Charge</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete this charge for "<span
                    class="font-medium text-(--color-text-primary)">{{
                        selectedBill?.customer_name }}</span>"?
            </p>
            <template #actions>
                <button @click="deleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>

        <!-- Customer Picker Dialog -->
        <BaseDialog v-model="customerPickerDialogOpen" max-width="sm">
            <div class="flex items-center gap-3 mb-4">
                <div
                    class="w-10 h-10 rounded-full bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Select Customer</h2>
                    <p class="text-xs text-(--color-text-secondary)">Choose who this charge is for</p>
                </div>
            </div>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Customer <span class="text-(--color-red)">*</span>
                </label>
                <select v-model="pickedCustomerId"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                    <option :value="null">Select a customer...</option>
                    <option v-for="customer in customersStore.customers" :key="customer.id" :value="customer.id">
                        {{ getCustomerDisplayName(customer) }}
                    </option>
                </select>
            </div>

            <template #actions>
                <button @click="customerPickerDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>
                <button @click="confirmCustomerPick" :disabled="!pickedCustomerId"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    Continue
                </button>
            </template>
        </BaseDialog>

        <!-- Create Charge Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Add Additional Charge</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    Create a new charge for {{ pickedCustomerName }}
                </p>
            </div>
            <AdditionalChargeForm v-if="createDialogOpen && pickedCustomerId" :locked-customer-id="pickedCustomerId"
                @charge-created="handleChargeCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { useCustomersStore } from '@/stores/customers'
import { useClipboardStore } from '@/stores/clipboard'
import type { CustomerAdditionalBill, CustomerAdditionalBillSortField } from '@/types/customerAdditionalBill'
import type { Customer } from '@/types/customer'
import CustomerAdditionalBillRow from './CustomerAdditionalBillRow.vue'
import CustomerAdditionalBillDetail from './CustomerAdditionalBillDetail.vue'
import AdditionalChargeForm from './AdditionalChargeForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'

const store = useCustomerAdditionalBillsStore()
const customersStore = useCustomersStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const detailDialogOpen = ref(false)
const paymentDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const customerPickerDialogOpen = ref(false)
const createDialogOpen = ref(false)
const selectedBill = ref<CustomerAdditionalBill | null>(null)
const pickedCustomerId = ref<number | null>(null)

const paymentForm = ref({
    amount: 0,
    payment_date: new Date().toISOString().slice(0, 10),
})

const filteredBills = computed(() => store.filteredBills)
const searchQuery = computed(() => store.searchQuery)
const statusFilter = computed(() => store.statusFilter)
const sortField = computed(() => store.sortField)
const sortDirection = computed(() => store.sortDirection)

const maxPaymentAmount = computed(() => {
    if (!selectedBill.value) return 0
    return selectedBill.value.amount - selectedBill.value.paid_amount
})

const pickedCustomerName = computed(() => {
    if (!pickedCustomerId.value) return ''
    const c = customersStore.getCustomerById(pickedCustomerId.value)
    if (!c) return ''
    return c.company_name || c.contact_person || `Customer #${c.id}`
})

const getCustomerDisplayName = (customer: Customer): string => {
    return customer.company_name || customer.contact_person || `Customer #${customer.id}`
}

const fetchData = async () => {
    loading.value = true
    try {
        await Promise.all([
            customersStore.fetchCustomers(),
            store.fetchCharges(),
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    store.setSearchQuery((e.target as HTMLInputElement).value)
}

const clearSearch = () => {
    store.clearSearch()
}

const handleStatusFilterChange = () => {
    store.setStatusFilter(statusFilter.value as 'unpaid' | 'paid' | '')
}

const toggleSort = (field: CustomerAdditionalBillSortField) => {
    store.setSort(field)
}

const handleCopyToClipboard = async () => {
    const headers = 'Customer\tDescription\tAmount\tPaid\tStatus'
    const rows = filteredBills.value.map(b => {
        return `${b.customer_name}\t${b.description}\t${b.amount.toFixed(2)}\t${b.paid_amount.toFixed(2)}\t${b.status}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const openDetailDialog = (bill: CustomerAdditionalBill) => {
    selectedBill.value = bill
    detailDialogOpen.value = true
}

const openPaymentDialog = (bill: CustomerAdditionalBill) => {
    selectedBill.value = bill
    const remaining = bill.amount - bill.paid_amount
    paymentForm.value = {
        amount: remaining > 0 ? remaining : 0,
        payment_date: new Date().toISOString().slice(0, 10),
    }
    detailDialogOpen.value = false
    setTimeout(() => {
        paymentDialogOpen.value = true
    }, 300)
}

const openEditDialog = (bill: CustomerAdditionalBill) => {
    selectedBill.value = bill
    detailDialogOpen.value = false
    setTimeout(() => {
        editDialogOpen.value = true
    }, 300)
}

const openDeleteDialog = (bill: CustomerAdditionalBill) => {
    selectedBill.value = bill
    detailDialogOpen.value = false
    setTimeout(() => {
        deleteDialogOpen.value = true
    }, 300)
}

const openCustomerPicker = () => {
    pickedCustomerId.value = null
    customerPickerDialogOpen.value = true
}

const confirmCustomerPick = () => {
    if (!pickedCustomerId.value) return
    customerPickerDialogOpen.value = false
    setTimeout(() => {
        createDialogOpen.value = true
    }, 300)
}

const confirmPayment = async () => {
    if (!selectedBill.value) return

    if (paymentForm.value.amount <= 0) {
        push.error('Payment amount must be greater than 0')
        return
    }

    const paymentDate = formatDateForBackend(paymentForm.value.payment_date)

    const result = await store.markBillAsPaid(selectedBill.value.id, {
        amount: paymentForm.value.amount,
        payment_date: paymentDate,
    })

    if (result) {
        paymentDialogOpen.value = false
        await fetchData()
    }
}

const confirmDelete = async () => {
    if (!selectedBill.value) return
    const success = await store.deleteCharge(selectedBill.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchData()
    }
}

const handleChargeUpdated = async () => {
    editDialogOpen.value = false
    await fetchData()
}

const handleChargeCreated = async () => {
    createDialogOpen.value = false
    pickedCustomerId.value = null
    await fetchData()
}

onMounted(() => {
    fetchData()
})
</script>