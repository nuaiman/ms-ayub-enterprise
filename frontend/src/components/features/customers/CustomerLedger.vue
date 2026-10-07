<!-- src/components/features/customers/CustomerLedger.vue -->
<template>
    <div class="space-y-5">
        <!-- Header -->
        <div class="flex items-start gap-4">
            <div class="shrink-0">
                <div
                    class="w-16 h-16 rounded-full bg-(--color-blue)/10 border-2 border-(--color-border) flex items-center justify-center">
                    <svg class="w-8 h-8 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                </div>
            </div>

            <div class="flex-1 min-w-0">
                <h2 class="text-xl font-bold text-(--color-text-primary)">
                    {{ customerName }} — Ledger
                </h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    {{ ledger.summary.eventCount }} event(s)
                    <template v-if="ledger.hasSelection">
                        · <span class="text-(--color-blue) font-medium">{{ ledger.selectedCount }} selected</span>
                    </template>
                </p>
            </div>

            <!-- Export actions -->
            <div class="flex items-center gap-2 shrink-0">
                <button @click="handleExportCSV"
                    class="px-3 py-2 text-sm font-medium rounded-lg border border-(--color-border) hover:bg-(--color-muted-bg) transition-all duration-200 inline-flex items-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span class="hidden sm:inline">Export CSV</span>
                </button>

                <button @click="handlePrint"
                    class="px-3 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200 inline-flex items-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                    </svg>
                    <span class="hidden sm:inline">Print</span>
                </button>
            </div>
        </div>

        <!-- Summary -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Billed</p>
                <p class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(ledger.summary.totalBilled) }}</p>
            </div>
            <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                <p class="text-lg font-bold text-(--color-green)">{{ formatCurrency(ledger.summary.totalPaid) }}</p>
            </div>
            <div class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                <p class="text-lg font-bold"
                    :class="ledger.summary.outstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                    {{ formatCurrency(ledger.summary.outstanding) }}
                </p>
            </div>
        </div>

        <!-- Filters -->
        <div class="space-y-3 p-3 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10">
            <div class="flex flex-col sm:flex-row gap-3">
                <div class="relative flex-1">
                    <input :value="ledger.searchQuery" @input="handleSearch" type="text" placeholder="Search events..."
                        class="w-full pl-9 pr-8 py-2 rounded-lg text-sm bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--color-text-secondary)"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                </div>

                <input :value="ledger.dateFrom" @input="handleDateFrom" type="date"
                    class="px-3 py-2 rounded-lg text-sm bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />

                <input :value="ledger.dateTo" @input="handleDateTo" type="date"
                    class="px-3 py-2 rounded-lg text-sm bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />

                <button @click="ledger.resetFilters()"
                    class="px-3 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors whitespace-nowrap">
                    Reset
                </button>
            </div>

            <div class="flex flex-wrap items-center gap-1.5">
                <button v-for="t in allTypes" :key="t" @click="ledger.toggleType(t)" :class="[
                    'px-2.5 py-1 text-xs rounded-full border transition-colors',
                    ledger.typeFilter.includes(t)
                        ? 'bg-(--color-blue)/10 border-(--color-blue) text-(--color-blue)'
                        : 'border-(--color-border) text-(--color-text-secondary) hover:bg-(--color-muted-bg)'
                ]">
                    {{ ledger.typeLabel(t) }}
                </button>

                <label
                    class="ml-auto inline-flex items-center gap-2 px-3 py-1 rounded-full border border-(--color-border) cursor-pointer hover:bg-(--color-muted-bg) transition-colors">
                    <input type="checkbox" :checked="ledger.customerView" @change="handleCustomerViewToggle"
                        class="w-3.5 h-3.5 rounded border-(--color-border) text-(--color-green) focus:ring-2 focus:ring-(--color-green)/20 focus:ring-offset-0" />
                    <span class="text-xs font-medium text-(--color-text-secondary)">Customer view</span>
                </label>
            </div>
        </div>

        <!-- Selection toolbar -->
        <div v-if="!ledger.isLoading && ledger.filteredEvents.length > 0"
            class="flex flex-wrap items-center justify-between gap-3 p-3 rounded-lg border border-(--color-border) bg-(--color-surface)">
            <div class="flex items-center gap-3">
                <label class="inline-flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" :checked="ledger.allVisibleSelected"
                        :indeterminate.prop="ledger.someVisibleSelected" @change="handleSelectAllToggle"
                        class="w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 cursor-pointer" />
                    <span class="text-xs font-medium text-(--color-text-secondary) select-none">
                        <template v-if="ledger.hasSelection">
                            {{ ledger.selectedCount }} of {{ ledger.filteredEvents.length }} selected
                        </template>
                        <template v-else>
                            Select events to include (none = all)
                        </template>
                    </span>
                </label>
            </div>

            <div class="flex items-center gap-2">
                <button v-if="!ledger.allVisibleSelected" @click="ledger.selectAllVisible()"
                    class="text-xs font-medium text-(--color-blue) hover:opacity-80 transition-opacity">
                    Select all visible
                </button>
                <button v-if="ledger.hasSelection" @click="ledger.clearSelection()"
                    class="text-xs font-medium text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
                    Clear selection
                </button>
            </div>
        </div>

        <!-- Timeline -->
        <div v-if="ledger.isLoading" class="text-center py-12 text-sm text-(--color-text-secondary)">
            Loading ledger…
        </div>

        <div v-else-if="ledger.filteredEvents.length === 0"
            class="text-center py-12 text-sm text-(--color-text-secondary)">
            No events match the current filters.
        </div>

        <!-- Events list (no icon column, no timeline spine) -->
        <div v-else class="space-y-2">
            <div v-for="event in ledger.filteredEvents" :key="event.id" class="flex items-stretch gap-3">
                <!-- Checkbox -->
                <div class="shrink-0 flex items-start pt-3.5">
                    <input type="checkbox" :checked="ledger.isEventSelected(event.id)"
                        @change="ledger.toggleEvent(event.id)"
                        class="w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 cursor-pointer" />
                </div>

                <!-- Event card -->
                <div class="flex-1 min-w-0 rounded-lg border border-(--color-border) bg-(--color-surface) p-3 hover:border-(--color-blue)/30 transition-colors"
                    :class="{ 'ring-1 ring-(--color-blue)/40': ledger.isEventSelected(event.id) }">
                    <div class="flex items-start justify-between gap-3 flex-wrap">
                        <div class="min-w-0 flex-1">
                            <div class="flex items-center gap-2 flex-wrap">
                                <span class="text-xs font-semibold uppercase tracking-wider px-2 py-0.5 rounded-md"
                                    :class="typeChipClass(event.type)">
                                    {{ ledger.typeLabel(event.type) }}
                                </span>
                                <span v-if="event.reference" class="text-xs text-(--color-text-secondary) font-mono">
                                    {{ event.reference }}
                                </span>
                            </div>
                            <p class="text-sm font-medium text-(--color-text-primary) mt-1">
                                {{ event.title }}
                            </p>
                            <p class="text-xs text-(--color-text-secondary) mt-0.5">
                                {{ event.description }}
                            </p>
                        </div>

                        <div class="text-right shrink-0">
                            <p v-if="typeof event.amount === 'number' && event.amount !== 0"
                                class="text-sm font-semibold" :class="amountClass(event.amountKind)">
                                {{ formatCurrency(event.amount) }}
                            </p>
                            <p v-if="event.amountLabel" class="text-[10px] uppercase tracking-wider"
                                :class="amountClass(event.amountKind)">
                                {{ event.amountLabel }}
                            </p>
                            <p class="text-[10px] text-(--color-text-secondary) mt-1">
                                {{ ledger.formatEventDate(event.date) }}
                            </p>
                        </div>
                    </div>

                    <div v-if="event.meta && Object.keys(event.meta).length > 0"
                        class="mt-2 pt-2 border-t border-(--color-border)/60 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-x-4 gap-y-1">
                        <div v-for="(value, key) in event.meta" :key="key" class="text-xs min-w-0">
                            <span class="text-(--color-text-secondary)">{{ humanKey(String(key)) }}:</span>
                            <span class="text-(--color-text-primary) ml-1">{{ value ?? '—' }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Footer actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-2 border-t border-(--color-border)">
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>

        <!-- Hidden print host -->
        <Teleport to="body">
            <div v-if="printMode" class="ledger-print-root">
                <CustomerLedgerView :customer-name="customerName" :customer-phone="customer?.phone"
                    :customer-address="customer?.address" :events="ledger.activeEvents" :summary="ledger.summary"
                    :date-from="ledger.dateFrom" :date-to="ledger.dateTo" />
            </div>
        </Teleport>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useCustomerLedgerStore, type LedgerEventType } from '@/stores/customerLedger'
import { useCustomersStore } from '@/stores/customers'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'
import CustomerLedgerView from './CustomerLedgerView.vue'

const props = defineProps<{
    customerId: number
}>()

const emit = defineEmits<{
    'close': []
}>()

const ledger = useCustomerLedgerStore()
const customersStore = useCustomersStore()

const printMode = ref(false)

const customerName = computed(() => customersStore.getCustomerName(props.customerId))
const customer = computed(() => customersStore.getCustomerById(props.customerId))

const allTypes: LedgerEventType[] = [
    'customer_created',
    'lot_created',
    'store_created',
    'delivery_created',
    'delivery_item_added',
    'damage_recorded',
    'lot_transferred_out',
    'lot_transferred_in',
    'customer_store_bill',
    'customer_delivery_bill',
    'customer_additional_bill',
    'invoice_created',
    'payment_received',
]

const handleSearch = (e: Event) => {
    ledger.setSearchQuery((e.target as HTMLInputElement).value)
}

const handleDateFrom = (e: Event) => {
    ledger.setDateRange((e.target as HTMLInputElement).value, ledger.dateTo)
}

const handleDateTo = (e: Event) => {
    ledger.setDateRange(ledger.dateFrom, (e.target as HTMLInputElement).value)
}

const handleCustomerViewToggle = (e: Event) => {
    ledger.setCustomerView((e.target as HTMLInputElement).checked)
}

const handleSelectAllToggle = () => {
    if (ledger.allVisibleSelected) {
        ledger.clearSelection()
    } else {
        ledger.selectAllVisible()
    }
}

const typeChipClass = (type: LedgerEventType): string => {
    const map: Record<LedgerEventType, string> = {
        customer_created: 'bg-blue-500/10 text-blue-600',
        lot_created: 'bg-indigo-500/10 text-indigo-600',
        store_created: 'bg-teal-500/10 text-teal-600',
        delivery_created: 'bg-green-500/10 text-green-600',
        delivery_item_added: 'bg-lime-500/10 text-lime-600',
        damage_recorded: 'bg-red-500/10 text-red-600',
        lot_transferred_out: 'bg-orange-500/10 text-orange-600',
        lot_transferred_in: 'bg-amber-500/10 text-amber-600',
        customer_store_bill: 'bg-purple-500/10 text-purple-600',
        customer_delivery_bill: 'bg-purple-500/10 text-purple-600',
        customer_additional_bill: 'bg-purple-500/10 text-purple-600',
        invoice_created: 'bg-fuchsia-500/10 text-fuchsia-600',
        payment_received: 'bg-emerald-500/10 text-emerald-600',
    }
    return map[type] || 'bg-(--color-muted-bg) text-(--color-text-secondary)'
}

const amountClass = (kind?: 'debit' | 'credit' | 'neutral'): string => {
    if (kind === 'debit') return 'text-(--color-blue)'
    if (kind === 'credit') return 'text-(--color-green)'
    return 'text-(--color-text-secondary)'
}

const humanKey = (key: string): string => {
    return key.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
}

const csvEscape = (value: unknown): string => {
    if (value === null || value === undefined) return ''
    const s = String(value).replace(/"/g, '""')
    return `"${s}"`
}

const handleExportCSV = () => {
    const events = ledger.activeEvents

    const headers = ['Date', 'Type', 'Reference', 'Title', 'Description', 'Amount', 'Amount Label', 'Meta']
    const rows = events.map((e) => {
        const metaPairs: string[] = []
        if (e.meta) {
            for (const [k, v] of Object.entries(e.meta)) {
                metaPairs.push(`${k}=${v ?? ''}`)
            }
        }
        return [
            new Date(e.date).toISOString(),
            e.type,
            e.reference || '',
            e.title,
            e.description,
            typeof e.amount === 'number' ? e.amount.toFixed(2) : '',
            e.amountLabel || '',
            metaPairs.join(' | '),
        ]
    })

    const lines = [headers.map(csvEscape).join(','), ...rows.map((r) => r.map(csvEscape).join(','))]
    const csv = lines.join('\n')

    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const suffix = ledger.hasSelection ? '_selected' : ''
    a.download = `ledger_${slugify(customerName.value)}${suffix}_${timestamp()}.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)

    push.success('CSV exported')
}

const slugify = (s: string): string =>
    s.toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_+|_+$/g, '').slice(0, 40)

const timestamp = (): string => {
    const d = new Date()
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}_${String(d.getHours()).padStart(2, '0')}${String(d.getMinutes()).padStart(2, '0')}`
}

const handlePrint = async () => {
    printMode.value = true
    await nextTick()
    await new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(() => r(null))))

    const cleanup = () => {
        printMode.value = false
        window.removeEventListener('afterprint', cleanup)
    }
    window.addEventListener('afterprint', cleanup)
    window.print()
    setTimeout(cleanup, 5000)
}

watch(
    () => props.customerId,
    async (id) => {
        await ledger.setCustomer(id)
    },
    { immediate: true }
)

onMounted(async () => {
    await ledger.setCustomer(props.customerId)
})

onUnmounted(() => {
    ledger.setCustomer(null)
    ledger.resetFilters()
    ledger.setCustomerView(false)
    ledger.clearSelection()
})
</script>