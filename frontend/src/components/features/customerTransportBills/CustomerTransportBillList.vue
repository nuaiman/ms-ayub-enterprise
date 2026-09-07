<!-- src/components/features/customerTransportBills/CustomerTransportBillList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">
                    Customer Transport Bills
                </h2>

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
                    <option value="cancelled">Cancelled</option>
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
            </div>
        </div>

        <!-- Table -->
        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <!-- Header Row -->
                <div
                    class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <!-- Transport - 1 columns -->
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('transport_id')">
                        <span class="flex items-center gap-1">
                            Trp.
                            <svg v-if="sortField === 'transport_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Customer - 2 columns -->
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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

                    <!-- Vehicles - 1 columns -->
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

                    <!-- Amount - 1 columns -->
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('bill_amount')">
                        <span class="flex items-center gap-1">
                            Amount
                            <svg v-if="sortField === 'bill_amount'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>

                    <!-- Status - 1 columns -->
                    <div class="col-span-1 cursor-pointer hover:text-(--color-text-primary) transition-colors"
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

                    <!-- Actions - 2 columns -->
                    <div class="col-span-2 flex items-center justify-end">
                        <span>Actions</span>
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

                        <p class="text-sm text-(--color-text-secondary)">
                            Loading bills...
                        </p>
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
                                    d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                            </svg>
                        </div>

                        <p class="text-sm font-medium text-(--color-text-primary)">
                            No bills found
                        </p>

                        <p class="text-xs text-(--color-text-secondary)">
                            {{ searchQuery
                                ? 'Try adjusting your search'
                                : 'Transports with vehicles having customer_charge > 0 will appear here'
                            }}
                        </p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <CustomerTransportBillRow v-for="bill in filteredBills" :key="bill.id" :bill="bill"
                        @view="openDetailDialog" @pay="openPaymentDialog" @cancel="handleCancelBill" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredBills.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">
                Showing {{ filteredBills.length }} bills
            </p>

            <div class="flex items-center gap-4 text-xs text-(--color-text-secondary)">
                <span>
                    Total: {{ formatCurrency(store.totalAmount) }}
                </span>

                <span class="text-(--color-yellow)">
                    Unpaid: {{ formatCurrency(store.totalUnpaidAmount) }}
                </span>
            </div>
        </div>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <CustomerTransportBillDetail v-if="selectedBill" :bill="selectedBill" @close="detailDialogOpen = false"
                @pay="handlePayFromDetail" @cancel="handleCancelFromDetail" />
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
                        <h2 class="text-lg font-bold text-(--color-text-primary)">
                            Record Payment
                        </h2>

                        <p class="text-xs text-(--color-text-secondary)">
                            For Transport #{{ selectedBill?.transport_id }}
                        </p>
                    </div>
                </div>

                <div class="space-y-3">
                    <div class="grid grid-cols-2 gap-3 text-sm bg-(--color-muted-bg)/30 p-3 rounded-lg">
                        <div>
                            <p class="text-xs text-(--color-text-secondary)">
                                Bill Amount
                            </p>

                            <p class="font-medium text-(--color-blue)">
                                {{ formatCurrency(selectedBill?.bill_amount) }}
                            </p>
                        </div>

                        <div>
                            <p class="text-xs text-(--color-text-secondary)">
                                Already Paid
                            </p>

                            <p class="font-medium text-(--color-green)">
                                {{ formatCurrency(selectedBill?.paid_amount) }}
                            </p>
                        </div>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Payment Amount
                            <span class="text-(--color-red)">*</span>
                        </label>

                        <div class="relative">
                            <span
                                class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">
                                ৳
                            </span>

                            <input v-model.number="paymentForm.amount" type="number" step="0.01" min="0.01"
                                :max="selectedBill?.bill_amount" placeholder="0.00" required
                                class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                        </div>

                        <p class="text-xs text-(--color-text-secondary) mt-1">
                            Max: {{ formatCurrency(selectedBill?.bill_amount) }}
                        </p>
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Payment Date
                            <span class="text-(--color-red)">*</span>
                        </label>

                        <input v-model="paymentForm.payment_date" type="date" required
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Notes
                        </label>

                        <textarea v-model="paymentForm.notes" rows="2" placeholder="Optional notes"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
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

        <!-- Cancel Confirmation Dialog -->
        <BaseDialog v-model="cancelDialogOpen" max-width="sm">
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
                        <h2 class="text-lg font-bold text-(--color-text-primary)">
                            Cancel Bill
                        </h2>

                        <p class="text-xs text-(--color-text-secondary)">
                            This will set customer_charge to 0 on all vehicles
                        </p>
                    </div>
                </div>

                <p class="text-sm text-(--color-text-secondary) mt-4">
                    Are you sure you want to cancel this bill for
                    "<span class="font-medium text-(--color-text-primary)">
                        Transport #{{ selectedBill?.transport_id }}
                    </span>"?
                </p>

                <p class="text-xs text-(--color-text-secondary) mt-2">
                    This will set customer_charge to 0 on all {{ selectedBill?.total_vehicles }} vehicles.
                </p>
            </div>

            <template #actions>
                <button @click="cancelDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>

                <button @click="confirmCancel"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-yellow) text-white rounded-lg hover:opacity-90 transition-colors">
                    Confirm Cancel
                </button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCustomerTransportBillsStore } from '@/stores/customerTransportBills'
import { useTransportsStore } from '@/stores/transports'
import { useVehiclesStore } from '@/stores/vehicles'
import { useCustomersStore } from '@/stores/customers'
import { useBrokersStore } from '@/stores/brokers'
import { useClipboardStore } from '@/stores/clipboard'
import type {
    CustomerTransportBill,
    CustomerTransportBillSortField,
} from '@/types/customerTransportBill'
import CustomerTransportBillRow from './CustomerTransportBillRow.vue'
import CustomerTransportBillDetail from './CustomerTransportBillDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { formatCurrency } from '@/utils/currency'
import { formatDateForBackend } from '@/utils/date'
import { push } from 'notivue'

const store = useCustomerTransportBillsStore()
const transportsStore = useTransportsStore()
const vehiclesStore = useVehiclesStore()
const customersStore = useCustomersStore()
const brokersStore = useBrokersStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const statusFilter = ref('')
const detailDialogOpen = ref(false)
const paymentDialogOpen = ref(false)
const cancelDialogOpen = ref(false)
const selectedBill = ref<CustomerTransportBill | null>(null)

const paymentForm = ref({
    amount: 0,
    payment_date: new Date().toISOString().slice(0, 10),
    notes: '',
})

const filteredBills = computed(() => store.filteredBills)
const sortField = computed(() => store.sortField)
const sortDirection = computed(() => store.sortDirection)

const fetchData = async () => {
    loading.value = true

    try {
        await Promise.all([
            transportsStore.fetchTransports(),
            vehiclesStore.fetchVehicles(),
            customersStore.fetchCustomers(),
            brokersStore.fetchBrokers(),
        ])
    } finally {
        loading.value = false
    }
}

const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement
    searchQuery.value = target.value
    store.setSearchQuery(searchQuery.value)
}

const clearSearch = () => {
    searchQuery.value = ''
    store.clearSearch()
}

const handleStatusFilterChange = () => {
    store.setStatusFilter(
        statusFilter.value as 'unpaid' | 'paid' | 'cancelled' | ''
    )
}

const toggleSort = (field: CustomerTransportBillSortField) => {
    store.setSort(field)
}

const handleCopyToClipboard = async () => {
    const headers = 'Transport\tCustomer\tFrom\tTo\tVehicles\tAmount\tStatus'

    const rows = filteredBills.value.map((b: CustomerTransportBill) => {
        return `#${b.transport_id}\t${b.customer_name}\t${b.from_location}\t${b.to_location || ''}\t${b.total_vehicles}\t${b.bill_amount.toFixed(2)}\t${b.status}`
    })

    await clipboardStore.copyToClipboard(
        headers + '\n' + rows.join('\n')
    )
}

const openDetailDialog = (bill: CustomerTransportBill) => {
    selectedBill.value = bill
    detailDialogOpen.value = true
}

const openPaymentDialog = (bill: CustomerTransportBill) => {
    selectedBill.value = bill

    paymentForm.value = {
        amount: bill.bill_amount - bill.paid_amount,
        payment_date: new Date().toISOString().slice(0, 10),
        notes: '',
    }

    paymentDialogOpen.value = true
}

const handlePayFromDetail = (bill: CustomerTransportBill) => {
    detailDialogOpen.value = false

    setTimeout(() => {
        openPaymentDialog(bill)
    }, 300)
}

const handleCancelBill = (bill: CustomerTransportBill) => {
    selectedBill.value = bill
    cancelDialogOpen.value = true
}

const handleCancelFromDetail = (bill: CustomerTransportBill) => {
    detailDialogOpen.value = false

    setTimeout(() => {
        selectedBill.value = bill
        cancelDialogOpen.value = true
    }, 300)
}

const confirmPayment = async () => {
    if (!selectedBill.value) return

    if (paymentForm.value.amount <= 0) {
        push.error('Payment amount must be greater than 0')
        return
    }

    if (
        paymentForm.value.amount >
        selectedBill.value.bill_amount - selectedBill.value.paid_amount
    ) {
        push.error('Payment amount exceeds remaining balance')
        return
    }

    const paymentDate = formatDateForBackend(
        paymentForm.value.payment_date
    )

    const result = await store.markBillAsPaid(
        selectedBill.value.transport_id,
        {
            payment_date: paymentDate,
            notes: paymentForm.value.notes?.trim() || null,
        }
    )

    if (result) {
        push.success('Payment recorded successfully!')
        paymentDialogOpen.value = false
        await fetchData()
    }
}

const confirmCancel = async () => {
    if (!selectedBill.value) return

    const result = await store.cancelBill(
        selectedBill.value.transport_id
    )

    if (result) {
        push.success('Bill cancelled successfully')
        cancelDialogOpen.value = false
        await fetchData()
    }
}

onMounted(() => {
    fetchData()
})
</script>