<!-- src/components/features/salaries/SalaryList.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-200px)]">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
            <div class="flex items-center gap-3">
                <h2 class="text-lg font-semibold text-(--color-text-primary)">Salaries</h2>
                <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
                    {{ filteredSalaries.length }}
                </span>
            </div>

            <div class="flex items-center gap-2 flex-wrap">
                <!-- Search -->
                <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
                    <input :value="searchQuery" @input="handleSearch" type="text" placeholder="Search salaries..."
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

                <!-- Status Filter -->
                <select v-model="statusFilter"
                    class="px-3 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                    <option value="">All Status</option>
                    <option value="draft">Draft</option>
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

                <!-- Create -->
                <button v-if="canManageSalaries" @click="createDialogOpen = true"
                    class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    <span class="hidden sm:inline">Create</span>
                </button>
            </div>
        </div>

        <!-- Rest of the template remains the same -->
        <div class="flex-1 min-h-0 overflow-auto">
            <div class="min-w-225">
                <!-- Header Row -->
                <div
                    class="grid grid-cols-12 items-center py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
                    <div class="col-span-3 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('employee_id')">
                        <span class="flex items-center gap-1">
                            Employee
                            <svg v-if="sortField === 'employee_id'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-2 cursor-pointer hover:text-(--color-text-primary) transition-colors"
                        @click="toggleSort('month_year')">
                        <span class="flex items-center gap-1">
                            Month
                            <svg v-if="sortField === 'month_year'" class="w-3 h-3"
                                :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor"
                                viewBox="0 0 24 24">
                                <path d="M7 10l5 5 5-5z" />
                            </svg>
                        </span>
                    </div>
                    <div class="col-span-1 text-right">Basic</div>
                    <div class="col-span-1 text-right">Bonus</div>
                    <div class="col-span-1 text-right">Ded.</div>
                    <div class="col-span-1 text-right">Total</div>
                    <div class="col-span-1">Status</div>
                    <div class="col-span-2 text-right">Actions</div>
                </div>

                <!-- Loading -->
                <div v-if="loading" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-4">
                        <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                        </svg>
                        <p class="text-sm text-(--color-text-secondary)">Loading salaries...</p>
                    </div>
                </div>

                <!-- Empty -->
                <div v-else-if="filteredSalaries.length === 0" class="flex items-center justify-center py-12">
                    <div class="text-center space-y-3">
                        <div
                            class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                            <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                                viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                    d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v1m0 4v1m0-1v1m0-1h.01M12 15v1" />
                            </svg>
                        </div>
                        <p class="text-sm font-medium text-(--color-text-primary)">No salaries found</p>
                        <p class="text-xs text-(--color-text-secondary)">{{ searchQuery ? 'Try adjusting your search' :
                            'Create a new salary to get started' }}</p>
                    </div>
                </div>

                <!-- Rows -->
                <div v-else>
                    <SalaryRow v-for="salary in filteredSalaries" :key="salary.id" :salary="salary"
                        @view="openDetailDialog" @edit="handleEditSalary" @delete="handleDeleteSalary"
                        @pay="handleMarkAsPaid" @cancel="handleCancelSalary" @updated="fetchSalaries" />
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div v-if="!loading && filteredSalaries.length > 0"
            class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0 mt-auto">
            <p class="text-xs text-(--color-text-secondary)">Showing {{ filteredSalaries.length }} of {{
                salariesStore.totalSalaries }} salaries</p>
        </div>

        <!-- Dialogs -->
        <!-- Create Dialog -->
        <BaseDialog v-model="createDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Create New Salary</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Add a salary record for an employee</p>
            </div>
            <SalaryForm mode="create" @salary-created="handleSalaryCreated" @cancel="createDialogOpen = false" />
        </BaseDialog>

        <!-- Detail Dialog -->
        <BaseDialog v-model="detailDialogOpen" max-width="3xl">
            <SalaryDetail v-if="selectedSalary" :salary="selectedSalary" @close="closeDetailDialog"
                @edit="handleEditSalaryFromDetail" @updated="fetchSalaries" />
        </BaseDialog>

        <!-- Edit Dialog -->
        <BaseDialog v-model="editDialogOpen" max-width="3xl">
            <div class="mb-6">
                <h2 class="text-xl font-bold text-(--color-text-primary)">Edit Salary</h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">Update salary details</p>
            </div>
            <SalaryForm v-if="selectedSalary" mode="edit" :salary="selectedSalary" @salary-updated="handleSalaryUpdated"
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
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Delete Salary</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to delete the salary record for <span
                    class="font-medium text-(--color-text-primary)">{{ getEmployeeName(selectedSalary?.employee_id)
                    }}</span>?
            </p>
            <template #actions>
                <button @click="deleteDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmDelete"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-colors">Delete</button>
            </template>
        </BaseDialog>

        <!-- Mark as Paid Dialog -->
        <BaseDialog v-model="payDialogOpen" max-width="sm">
            <div class="flex items-center gap-3">
                <div
                    class="w-10 h-10 rounded-full bg-(--color-green)/10 text-(--color-green) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Mark as Paid</h2>
                    <p class="text-xs text-(--color-text-secondary)">Record payment for this salary</p>
                </div>
            </div>

            <div class="space-y-4 mt-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Payment Method</label>
                    <select v-model="paymentForm.method"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                        <option value="cash">Cash</option>
                        <option value="bank_transfer">Bank Transfer</option>
                        <option value="check">Check</option>
                        <option value="mobile_banking">Mobile Banking</option>
                    </select>
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Reference Number</label>
                    <input v-model="paymentForm.reference" type="text" placeholder="Enter reference number"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Payment Date</label>
                    <input v-model="paymentForm.date" type="date"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
            </div>

            <template #actions>
                <button @click="payDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmPay"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-green) text-white rounded-lg hover:opacity-90 transition-colors">Confirm
                    Payment</button>
            </template>
        </BaseDialog>

        <!-- Cancel Confirmation Dialog -->
        <BaseDialog v-model="cancelDialogOpen" max-width="sm">
            <div class="flex items-center gap-3">
                <div
                    class="w-10 h-10 rounded-full bg-(--color-yellow)/10 text-(--color-yellow) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-lg font-bold text-(--color-text-primary)">Cancel Salary</h2>
                    <p class="text-xs text-(--color-text-secondary)">This action cannot be undone</p>
                </div>
            </div>
            <p class="text-sm text-(--color-text-secondary) mt-4">
                Are you sure you want to cancel the salary record for <span
                    class="font-medium text-(--color-text-primary)">{{
                        getEmployeeName(selectedSalary?.employee_id) }}</span>?
            </p>
            <template #actions>
                <button @click="cancelDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
                <button @click="confirmCancel"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-yellow) text-white rounded-lg hover:opacity-90 transition-colors">Confirm
                    Cancel</button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useSalariesStore } from '@/stores/salaries'
import { useUsersStore } from '@/stores/users'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { Salary, PaymentMethod } from '@/types/salary'
import SalaryRow from './SalaryRow.vue'
import SalaryForm from './SalaryForm.vue'
import SalaryDetail from './SalaryDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import MonthFilter from '@/components/ui/MonthFilter.vue'
import { push } from 'notivue'

const salariesStore = useSalariesStore()
const usersStore = useUsersStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const searchQuery = ref('')
const monthFilter = ref('')
const statusFilter = ref('')
const sortField = ref<'employee_id' | 'month_year' | 'status'>('month_year')
const sortDirection = ref<'asc' | 'desc'>('desc')

// Dialogs
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const payDialogOpen = ref(false)
const cancelDialogOpen = ref(false)

const selectedSalary = ref<Salary | null>(null)

// Payment form
const paymentForm = ref({
    method: 'cash' as PaymentMethod,
    reference: '',
    date: new Date().toISOString().slice(0, 10)
})

const canManageSalaries = computed(() => {
    const role = auth.user?.role
    return role === 'admin' || role === 'manager'
})

// Get unique months from salaries
const availableMonths = computed(() => {
    const months = new Set<string>()
    salariesStore.salaries.forEach(s => {
        months.add(s.month_year)
    })
    return Array.from(months).sort((a, b) => b.localeCompare(a))
})

const filteredSalaries = computed(() => {
    let result = [...salariesStore.salaries]

    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(s => {
            const employeeName = usersStore.getUserName(s.employee_id).toLowerCase()
            return employeeName.includes(query) ||
                s.month_year.includes(query) ||
                s.status.includes(query)
        })
    }

    if (monthFilter.value) {
        result = result.filter(s => s.month_year === monthFilter.value)
    }

    if (statusFilter.value) {
        result = result.filter(s => s.status === statusFilter.value)
    }

    result.sort((a, b) => {
        let comparison = 0
        switch (sortField.value) {
            case 'employee_id':
                comparison = a.employee_id - b.employee_id
                break
            case 'month_year':
                comparison = a.month_year.localeCompare(b.month_year)
                break
            case 'status':
                comparison = a.status.localeCompare(b.status)
                break
            default:
                comparison = 0
        }
        return sortDirection.value === 'desc' ? -comparison : comparison
    })

    return result
})

const getEmployeeName = (id?: number): string => {
    if (!id) return 'Unknown'
    return usersStore.getUserName(id)
}

const fetchSalaries = async () => {
    loading.value = true
    try {
        await Promise.all([
            salariesStore.fetchSalaries(),
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

const toggleSort = (field: 'employee_id' | 'month_year' | 'status') => {
    if (sortField.value === field) {
        sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
    } else {
        sortField.value = field
        sortDirection.value = 'desc'
    }
}

const handleCopyToClipboard = async () => {
    const headers = 'Employee\tMonth\tBasic\tBonus\tDeductions\tTotal\tStatus'
    const rows = filteredSalaries.value.map(s => {
        const basic = salariesStore.getEmployeeBasicSalary(s.employee_id)
        const total = basic + s.bonus - s.deductions
        return `${usersStore.getUserName(s.employee_id)}\t${s.month_year}\t${basic.toFixed(2)}\t${s.bonus.toFixed(2)}\t${s.deductions.toFixed(2)}\t${total.toFixed(2)}\t${s.status}`
    })
    await clipboardStore.copyToClipboard(headers + '\n' + rows.join('\n'))
}

// Dialog handlers
const openDetailDialog = (salary: Salary) => {
    selectedSalary.value = salary
    detailDialogOpen.value = true
}

const closeDetailDialog = () => {
    detailDialogOpen.value = false
    setTimeout(() => { selectedSalary.value = null }, 300)
}

const handleEditSalary = (salary: Salary) => {
    selectedSalary.value = salary
    editDialogOpen.value = true
}

const handleEditSalaryFromDetail = (salary: Salary) => {
    detailDialogOpen.value = false
    setTimeout(() => {
        selectedSalary.value = salary
        editDialogOpen.value = true
    }, 300)
}

const handleDeleteSalary = (salary: Salary) => {
    selectedSalary.value = salary
    deleteDialogOpen.value = true
}

const handleMarkAsPaid = (salary: Salary) => {
    selectedSalary.value = salary
    paymentForm.value = {
        method: 'cash',
        reference: '',
        date: new Date().toISOString().slice(0, 10)
    }
    payDialogOpen.value = true
}

const handleCancelSalary = (salary: Salary) => {
    selectedSalary.value = salary
    cancelDialogOpen.value = true
}

const confirmDelete = async () => {
    if (!selectedSalary.value) return
    const success = await salariesStore.deleteSalary(selectedSalary.value.id)
    if (success) {
        deleteDialogOpen.value = false
        await fetchSalaries()
    }
}

const confirmPay = async () => {
    if (!selectedSalary.value) return

    let paymentDateISO = undefined
    if (paymentForm.value.date) {
        const dateObj = new Date(paymentForm.value.date + 'T00:00:00')
        paymentDateISO = dateObj.toISOString()
    }

    const result = await salariesStore.markSalaryAsPaid(selectedSalary.value.id, {
        payment_method: paymentForm.value.method,
        reference_number: paymentForm.value.reference || null,
        payment_date: paymentDateISO,
    })

    if (result) {
        payDialogOpen.value = false
        await fetchSalaries()
    }
}

const confirmCancel = async () => {
    if (!selectedSalary.value) return
    const result = await salariesStore.cancelSalary(selectedSalary.value.id)
    if (result) {
        cancelDialogOpen.value = false
        await fetchSalaries()
    }
}

const handleSalaryCreated = async () => {
    createDialogOpen.value = false
    await fetchSalaries()
}

const handleSalaryUpdated = async () => {
    editDialogOpen.value = false
    detailDialogOpen.value = false
    await fetchSalaries()
}

// Auto-select current month if no filter
watch(() => salariesStore.salaries, (newSalaries) => {
    if (newSalaries.length > 0 && !monthFilter.value) {
        const currentMonth = new Date().toISOString().slice(0, 7)
        const hasCurrentMonth = newSalaries.some(s => s.month_year === currentMonth)
        if (hasCurrentMonth) {
            monthFilter.value = currentMonth
        }
    }
}, { immediate: true })

onMounted(() => {
    fetchSalaries()
})
</script>