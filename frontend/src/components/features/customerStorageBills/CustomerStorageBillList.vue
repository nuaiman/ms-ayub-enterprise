<!-- src/components/features/customerStorageBills/CustomerStorageBillList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Customer Storage Bills</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredCustomerBills.length }}
                </span>
                <span v-if="customerStorageBillsStore.totalOutstanding > 0"
                    class="text-sm text-(--color-red) bg-(--color-red)/10 px-2 py-0.5 rounded-md">
                    Outstanding: {{ formatCurrency(customerStorageBillsStore.totalOutstanding) }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="customerStorageBillsStore.searchQuery" @input="handleSearch" type="text"
                        placeholder="Search customers..."
                        class="w-full sm:w-56 pl-9 pr-8 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--color-text-secondary)"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <button v-if="customerStorageBillsStore.searchQuery"
                        @click="customerStorageBillsStore.clearSearch()" type="button"
                        class="absolute right-2.5 top-1/2 -translate-y-1/2 text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

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
            </div>
        </div>

        <!-- Table -->
        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <!-- Header Row -->
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('item_name')">
                        <span class="flex items-center gap-1">
                            Item / Lot
                            <svg v-if="sortField === 'item_name'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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
                        <p class="text-sm text-(--color-text-secondary)">Loading customer storage bills...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredCustomerBills.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No customer storage bills found</p>
                        <p class="text-xs text-(--color-text-secondary)">
                            {{ customerStorageBillsStore.searchQuery
                                ? 'Try adjusting your search'
                                : 'Active lots with customer billing will appear here'
                            }}
                        </p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <CustomerStorageBillRow v-for="bill in filteredCustomerBills" :key="bill.id" :bill="bill"
                        @view="openDetailDialog" @pay="openPaymentDialog" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredCustomerBills.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredCustomerBills.length }} customer storage
                bills
            </p>
            <div class="flex items-center gap-4 text-xs text-(--color-text-secondary)">
                <span>Total Monthly: {{ formatCurrency(customerStorageBillsStore.totalMonthlyBill) }}</span>
                <span>Total Billed: {{ formatCurrency(customerStorageBillsStore.totalBilled) }}</span>
                <span>Total Paid: {{ formatCurrency(customerStorageBillsStore.totalPaid) }}</span>
                <span class="font-semibold text-(--color-red)">Outstanding: {{
                    formatCurrency(customerStorageBillsStore.totalOutstanding)
                }}</span>
            </div>
        </div>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <CustomerStorageBillDetail v-if="selectedBill" :bill="selectedBill" @close="detailDialogOpen = false"
                @pay="handlePayFromDetail" />
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
                            <p class="text-xs text-(--color-text-secondary)">Total Billed</p>
                            <p class="font-medium text-(--color-text-primary)">{{
                                formatCurrency(selectedBill?.total_billed) }}</p>
                        </div>
                        <div>
                            <p class="text-xs text-(--color-text-secondary)">Outstanding</p>
                            <p class="font-medium text-(--color-red)">{{ formatCurrency(selectedBill?.outstanding) }}
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
                            <input v-model.number="paymentForm.amount" type="number" step="0.01" min="0"
                                :max="selectedBill?.outstanding" placeholder="0.00" required
                                class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        </div>
                        <p class="text-xs text-(--color-text-secondary) mt-1">Max: {{
                            formatCurrency(selectedBill?.outstanding) }}</p>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Paid Through Date
                        </label>
                        <input v-model="paymentForm.paid_through" type="date"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        <p class="text-xs text-(--color-text-secondary) mt-1">Date up to which payment is made</p>
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
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useCustomerStorageBillsStore } from '@/stores/customerStorageBills'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import { useClipboardStore } from '@/stores/clipboard'
import type { CustomerStorageBill, CustomerStorageBillSortField } from '@/types/customerStorageBill'
import CustomerStorageBillRow from './CustomerStorageBillRow.vue'
import CustomerStorageBillDetail from './CustomerStorageBillDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'

const customerStorageBillsStore = useCustomerStorageBillsStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const itemsStore = useItemsStore()
const customersStore = useCustomersStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const detailDialogOpen = ref(false)
const paymentDialogOpen = ref(false)
const selectedBill = ref<CustomerStorageBill | null>(null)

const paymentForm = ref({
    amount: 0,
    paid_through: '',
})

const filteredCustomerBills = computed(() => customerStorageBillsStore.filteredCustomerBills)

const sortField = computed(() => customerStorageBillsStore.sortField)
const sortDirection = computed(() => customerStorageBillsStore.sortDirection)

const fetchData = async () => {
    loading.value = true
    try {
        await Promise.all([
            lotsStore.fetchLots(),
            storesStore.fetchStores(),
            itemsStore.fetchItems(),
            customersStore.fetchCustomers(),
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    customerStorageBillsStore.setSearchQuery(target.value)
}

const toggleSort = (field: CustomerStorageBillSortField) => {
    customerStorageBillsStore.setSort(field)
}

const handleCopyToClipboard = async () => {
    const headers = 'Customer\tItem/Lot\tMonthly Bill\tTotal Billed\tTotal Paid\tOutstanding'
    const rows = filteredCustomerBills.value.map((b: CustomerStorageBill) => {
        return `${b.customer_name}\t${b.item_name} (Lot #${b.lot_id})\t${b.monthly_bill.toFixed(2)}\t${b.total_billed.toFixed(2)}\t${b.total_paid.toFixed(2)}\t${b.outstanding.toFixed(2)}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

const openDetailDialog = (bill: CustomerStorageBill) => {
    selectedBill.value = bill
    detailDialogOpen.value = true
}

const openPaymentDialog = (bill: CustomerStorageBill) => {
    selectedBill.value = bill
    paymentForm.value = {
        amount: bill.outstanding > 0 ? bill.outstanding : 0,
        paid_through: new Date().toISOString().slice(0, 10),
    }
    paymentDialogOpen.value = true
}

const handlePayFromDetail = (bill: CustomerStorageBill) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        openPaymentDialog(bill)
    }, 300)
}

const confirmPayment = async () => {
    if (!selectedBill.value) return
    if (paymentForm.value.amount <= 0) {
        push.error('Payment amount must be greater than 0')
        return
    }
    if (paymentForm.value.amount > selectedBill.value.outstanding) {
        push.error('Payment amount exceeds outstanding balance')
        return
    }

    // Format date for backend (converts "2026-09-06" to "2026-09-06T00:00:00Z")
    const paidThrough = formatDateForBackend(paymentForm.value.paid_through)

    const success = await customerStorageBillsStore.recordCustomerPayment(
        selectedBill.value.lot_id,
        paymentForm.value.amount,
        paidThrough || null
    )

    if (success) {
        push.success('Payment recorded successfully!')
        paymentDialogOpen.value = false
        await fetchData()
    }
}

// Auto-select current month
watch(() => customerStorageBillsStore.customerBillData, (data) => {
    if (data.length > 0 && !customerStorageBillsStore.monthFilter) {
        const currentMonth = new Date().toISOString().slice(0, 7)
        const hasCurrentMonth = data.some((b: CustomerStorageBill) => {
            if (!b.billing_start) return false
            const start = new Date(b.billing_start)
            const month = `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}`
            return month === currentMonth
        })
        if (hasCurrentMonth) {
            customerStorageBillsStore.monthFilter = currentMonth
        }
    }
}, { immediate: true })

onMounted(() => {
    fetchData()
})
</script>