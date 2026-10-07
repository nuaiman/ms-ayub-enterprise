<!-- src/components/features/customerStoreBills/CustomerStoreBillList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Customer Store Bills</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredBills.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
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

                <button v-if="canManage" @click="createDialogOpen = true"
                    class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    <span class="hidden sm:inline">Create</span>
                </button>
            </div>
        </div>

        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('bill_type')">
                        <span class="flex items-center gap-1">
                            Type
                            <svg v-if="sortField === 'bill_type'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('rate')">
                        <span class="flex items-center gap-1">
                            Rate
                            <svg v-if="sortField === 'rate'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('total_paid')">
                        <span class="flex items-center gap-1">
                            Total Paid
                            <svg v-if="sortField === 'total_paid'" class="w-3 h-3"
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
                        <p class="text-sm text-(--color-text-secondary)">Loading customer store bills...</p>
                    </div>
                </div>

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
                        <p class="text-sm font-medium text-(--color-text-primary)">No customer store bills found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new customer store bill to get started' }}</p>
                    </div>
                </div>

                <div v-else>
                    <CustomerStoreBillRow v-for="bill in filteredBills" :key="bill.id" :bill="bill"
                        @view="openDetailDialog" @edit="handleEditBill" @payment="handlePayment"
                        @delete="handleDeleteBill" @updated="fetchBills" />
                </div>
            </div>
        </div>

        <div v-if="!loading && filteredBills.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredBills.length }} of {{ totalBills }}
                bills</p>
            <p class="text-xs text-(--color-text-secondary)">Total Paid: {{ formatCurrency(totalPaid) }}</p>
        </div>

        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Customer Store Bill</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a new customer store bill</p>
            </div>
            <CustomerStoreBillForm mode="create" @bill-created="handleBillCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <CustomerStoreBillDetail v-if="selectedBill" :bill="selectedBill" @close="closeDetailDialog"
                @edit="handleEditBillFromDetail" @payment="handlePaymentFromDetail" @updated="fetchBills" />
        </BaseDialog>

        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Customer Store Bill</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update customer store bill information</p>
            </div>
            <CustomerStoreBillForm v-if="selectedBill" mode="edit" :bill="selectedBill"
                @bill-updated="handleBillUpdated" @cancel="editDialogOpen = false" />
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Customer Store Bill</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete customer store bill #<span
                    class="font-medium text-(--color-text-primary)">{{ selectedBill?.id }}</span>?
            </p>
            <template #actions>
                <button @click="deleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>

        <!-- Add Payment Dialog -->
        <BaseDialog v-model="paymentDialogOpen" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-green)/10 text-(--color-green) flex items-center justify-center shrink-0">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                        </svg>
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Add Payment</h2>
                        <p class="text-xs text-(--color-text-secondary)">Bill #{{ selectedBill?.id }}</p>
                    </div>
                </div>

                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <div class="flex items-center justify-between text-sm">
                        <span class="text-(--color-text-secondary)">Rate</span>
                        <span class="font-semibold text-(--color-text-primary)">{{ formatCurrency(selectedBill?.rate ||
                            0)
                            }}</span>
                    </div>
                    <div class="flex items-center justify-between text-sm mt-1">
                        <span class="text-(--color-text-secondary)">Paid So Far</span>
                        <span class="font-semibold text-(--color-green)">{{ formatCurrency(selectedBill?.total_paid ||
                            0)
                            }}</span>
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Amount <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="paymentForm.amount" type="number" step="0.01" min="0.01"
                            placeholder="0.00"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Payment Method</label>
                    <select v-model="paymentForm.payment_method"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                        <option value="cash">Cash</option>
                        <option value="bank_transfer">Bank Transfer</option>
                        <option value="check">Check</option>
                        <option value="mobile_banking">Mobile Banking</option>
                    </select>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Reference Number</label>
                    <input v-model="paymentForm.reference_number" type="text" placeholder="Optional"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Payment Date</label>
                    <input v-model="paymentForm.payment_date" type="date"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Notes</label>
                    <textarea v-model="paymentForm.notes" rows="2" placeholder="Optional"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>
            </div>
            <template #actions>
                <button @click="paymentDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmPayment" :disabled="paymentSubmitting"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-green) text-white rounded-lg hover:opacity-90 transition-colors disabled:opacity-50">
                    {{ paymentSubmitting ? 'Saving...' : 'Record Payment' }}
                </button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCustomerStoreBillsStore } from '@/stores/customerStoreBills'
import { useCustomersStore } from '@/stores/customers'
import { useStoresStore } from '@/stores/stores'
import { useLotsStore } from '@/stores/lots'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { CustomerStoreBill, CustomerStoreBillSortField } from '@/types/customerStoreBill'
import CustomerStoreBillRow from './CustomerStoreBillRow.vue'
import CustomerStoreBillForm from './CustomerStoreBillForm.vue'
import CustomerStoreBillDetail from './CustomerStoreBillDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'

const customerStoreBillsStore = useCustomerStoreBillsStore()
const customersStore = useCustomersStore()
const storesStore = useStoresStore()
const lotsStore = useLotsStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')

const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const paymentDialogOpen = ref(false)

const selectedBill = ref<CustomerStoreBill | null>(null)

const paymentSubmitting = ref(false)
const paymentForm = ref({
    amount: 0,
    payment_method: 'cash' as 'cash' | 'bank_transfer' | 'check' | 'mobile_banking',
    reference_number: '',
    payment_date: new Date().toISOString().slice(0, 10),
    notes: '',
})

const canManage = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

const filteredBills = computed(() => {
    let result = [...customerStoreBillsStore.bills]
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(b =>
            b.bill_type.toLowerCase().includes(query) ||
            String(b.rate).includes(query) ||
            String(b.total_paid).includes(query) ||
            customersStore.getCustomerName(b.customer_id).toLowerCase().includes(query)
        )
    }
    return result
})

const totalBills = computed(() => customerStoreBillsStore.bills.length)
const totalPaid = computed(() => customerStoreBillsStore.bills.reduce((s, b) => s + b.total_paid, 0))

const sortField = computed(() => customerStoreBillsStore.sortField)
const sortDirection = computed(() => customerStoreBillsStore.sortDirection)

const fetchBills = async () => {
    loading.value = true
    try {
        await Promise.all([
            customerStoreBillsStore.fetchCustomerStoreBills(),
            customersStore.fetchCustomers(),
            storesStore.fetchStores(),
            lotsStore.fetchLots(),
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

const toggleSort = (field: CustomerStoreBillSortField) => {
    customerStoreBillsStore.setSort(field)
}

const handleCopyToClipboard = async () => {
    const headers = 'Customer\tStore ID\tBill Type\tRate\tTotal Paid\tPaid Through'
    const rows = filteredBills.value.map(b =>
        `${customersStore.getCustomerName(b.customer_id)}\t${b.store_id}\t${b.bill_type}\t${b.rate.toFixed(2)}\t${b.total_paid.toFixed(2)}\t${b.total_paid_through || ''}`
    )
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const openDetailDialog = (bill: CustomerStoreBill) => {
    selectedBill.value = bill
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedBill.value = null }, 300)
}

const handleEditBill = (bill: CustomerStoreBill) => {
    selectedBill.value = bill
    editDialogOpen.value = true
}

const handleEditBillFromDetail = (bill: CustomerStoreBill) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedBill.value = bill
        editDialogOpen.value = true
    }, 300)
}

const handlePayment = (bill: CustomerStoreBill) => {
    selectedBill.value = bill
    paymentForm.value = {
        amount: 0,
        payment_method: 'cash',
        reference_number: '',
        payment_date: new Date().toISOString().slice(0, 10),
        notes: '',
    }
    paymentDialogOpen.value = true
}

const handlePaymentFromDetail = (bill: CustomerStoreBill) => {
    detailDialogOpen.value = false
    setTimeout(() => { handlePayment(bill) }, 300)
}

const handleDeleteBill = (bill: CustomerStoreBill) => {
    selectedBill.value = bill
    deleteDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedBill.value) return
    const success = await customerStoreBillsStore.deleteCustomerStoreBill(selectedBill.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchBills()
    }
}

const confirmPayment = async () => {
    if (!selectedBill.value) return
    if (!paymentForm.value.amount || paymentForm.value.amount <= 0) {
        push.error('Amount must be greater than 0')
        return
    }

    paymentSubmitting.value = true
    try {
        const paidThroughISO = paymentForm.value.payment_date
            ? `${paymentForm.value.payment_date}T00:00:00Z`
            : undefined

        const ok = await customerStoreBillsStore.createCustomerStoreBillPayment(selectedBill.value.id, {
            amount: paymentForm.value.amount,
            payment_date: paidThroughISO,
            payment_method: paymentForm.value.payment_method,
            reference_number: paymentForm.value.reference_number.trim() || null,
            notes: paymentForm.value.notes.trim() || null,
        })

        if (ok) {
            paymentDialogOpen.value = false
            await fetchBills()
        }
    } finally {
        paymentSubmitting.value = false
    }
}

const handleBillCreated = async () => {
    createDialogOpen.value = false
    await fetchBills()
}

const handleBillUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchBills()
}

onMounted(() => { fetchBills() })
</script>