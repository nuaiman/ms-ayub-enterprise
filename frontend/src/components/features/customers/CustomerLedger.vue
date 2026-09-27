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
                    {{ customerName }} ৳ Ledger
                </h2>
                <p class="text-sm text-(--color-text-secondary) mt-1">
                    {{ ledger.summary.eventCount }} event(s)
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

                <button @click="handleExportPDF" :disabled="exporting"
                    class="px-3 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200 inline-flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed">
                    <svg v-if="!exporting" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                    </svg>
                    <span class="hidden sm:inline">{{ exporting ? 'Generating...' : 'Export PDF' }}</span>
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

                <button @click="ledger.resetFilters"
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

        <!-- Timeline (on-screen) -->
        <div v-if="ledger.filteredEvents.length === 0" class="text-center py-12 text-sm text-(--color-text-secondary)">
            No events match the current filters.
        </div>

        <div v-else class="relative">
            <div class="absolute left-5 top-2 bottom-2 w-px bg-(--color-border)" aria-hidden="true"></div>

            <div class="space-y-3">
                <div v-for="event in ledger.filteredEvents" :key="event.id" class="relative flex gap-4 pl-0">
                    <div class="shrink-0 z-10">
                        <div
                            class="w-10 h-10 rounded-full border-2 border-(--color-border) bg-(--color-surface) flex items-center justify-center text-base">
                            {{ event.icon }}
                        </div>
                    </div>

                    <div
                        class="flex-1 min-w-0 rounded-lg border border-(--color-border) bg-(--color-surface) p-3 hover:border-(--color-blue)/30 transition-colors">
                        <div class="flex items-start justify-between gap-3 flex-wrap">
                            <div class="min-w-0 flex-1">
                                <div class="flex items-center gap-2 flex-wrap">
                                    <span class="text-xs font-semibold uppercase tracking-wider px-2 py-0.5 rounded-md"
                                        :class="typeChipClass(event.type)">
                                        {{ ledger.typeLabel(event.type) }}
                                    </span>
                                    <span v-if="event.reference"
                                        class="text-xs text-(--color-text-secondary) font-mono">
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
        </div>

        <!-- Footer actions -->
        <div class="flex flex-wrap items-center justify-end gap-3 pt-2 border-t border-(--color-border)">
            <button @click="emit('close')"
                class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Close
            </button>
        </div>

        <!-- Hidden print view host -->
        <div class="pdf-render-host" aria-hidden="true">
            <CustomerLedgerPrintView ref="printViewRef" :customer-name="customerName" :customer-phone="customer?.phone"
                :customer-address="customer?.address" :events="ledger.filteredEvents" :summary="ledger.summary"
                :date-from="ledger.dateFrom" :date-to="ledger.dateTo" :customer-view="ledger.customerView" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useCustomerLedgerStore, type LedgerEventType } from '@/stores/customerLedger'
import { useCustomersStore } from '@/stores/customers'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'
import jsPDF from 'jspdf'
import html2canvas from 'html2canvas'
import CustomerLedgerPrintView from './CustomerLedgerPrintView.vue'

const props = defineProps<{
    customerId: number
}>()

const emit = defineEmits<{
    'close': []
}>()

const ledger = useCustomerLedgerStore()
const customersStore = useCustomersStore()

const printViewRef = ref<InstanceType<typeof CustomerLedgerPrintView> | null>(null)
const exporting = ref(false)

const customerName = computed(() => customersStore.getCustomerName(props.customerId))
const customer = computed(() => customersStore.getCustomerById(props.customerId))

const allTypes: LedgerEventType[] = [
    'customer_created',
    'lot_created',
    'store_created',
    'delivery_created',
    'delivery_item_added',
    'damage_recorded',
    'transport_created',
    'vehicle_added',
    'storage_payment',
    'unload_payment',
    'delivery_payment',
    'transport_payment',
    'additional_charge_created',
    'additional_charge_payment',
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

const typeChipClass = (type: LedgerEventType): string => {
    const map: Record<LedgerEventType, string> = {
        customer_created: 'bg-blue-500/10 text-blue-600',
        lot_created: 'bg-indigo-500/10 text-indigo-600',
        store_created: 'bg-teal-500/10 text-teal-600',
        delivery_created: 'bg-green-500/10 text-green-600',
        delivery_item_added: 'bg-lime-500/10 text-lime-600',
        damage_recorded: 'bg-red-500/10 text-red-600',
        transport_created: 'bg-orange-500/10 text-orange-600',
        vehicle_added: 'bg-amber-500/10 text-amber-600',
        storage_payment: 'bg-emerald-500/10 text-emerald-600',
        unload_payment: 'bg-emerald-500/10 text-emerald-600',
        delivery_payment: 'bg-emerald-500/10 text-emerald-600',
        transport_payment: 'bg-emerald-500/10 text-emerald-600',
        additional_charge_created: 'bg-purple-500/10 text-purple-600',
        additional_charge_payment: 'bg-emerald-500/10 text-emerald-600',
    }
    return map[type] || 'bg-(--color-muted-bg) text-(--color-text-secondary)'
}

const amountClass = (kind?: 'debit' | 'credit' | 'neutral'): string => {
    if (kind === 'debit') return 'text-(--color-blue)'
    if (kind === 'credit') return 'text-(--color-green)'
    return 'text-(--color-text-secondary)'
}

const humanKey = (key: string): string => {
    return key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

const csvEscape = (value: unknown): string => {
    if (value === null || value === undefined) return ''
    const s = String(value).replace(/"/g, '""')
    return `"${s}"`
}

const handleExportCSV = () => {
    const headers = ['Date', 'Type', 'Reference', 'Title', 'Description', 'Amount', 'Amount Label', 'Meta']
    const rows = ledger.filteredEvents.map(e => {
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

    const lines = [
        headers.map(csvEscape).join(','),
        ...rows.map(r => r.map(csvEscape).join(',')),
    ]
    const csv = lines.join('\n')

    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `ledger_${slugify(customerName.value)}_${timestamp()}.csv`
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

const handleExportPDF = async () => {
    if (exporting.value) return
    exporting.value = true

    try {
        await nextTick()
        await new Promise(resolve => setTimeout(resolve, 200))

        const pageEls = printViewRef.value?.pageRefs?.filter(Boolean) as HTMLElement[] | undefined
        if (!pageEls || pageEls.length === 0) {
            push.error('Failed to prepare PDF')
            exporting.value = false
            return
        }

        const pdf = new jsPDF({
            orientation: 'portrait',
            unit: 'px',
            format: [794, 1123],
        })

        for (let i = 0; i < pageEls.length; i++) {
            const el = pageEls[i]
            if (!el) continue

            const canvas = await html2canvas(el, {
                scale: 2,
                useCORS: true,
                logging: false,
                backgroundColor: '#ffffff',
                width: 794,
                height: 1123,
                windowWidth: 794,
                windowHeight: 1123,
            })

            const imgData = canvas.toDataURL('image/png')
            if (i > 0) pdf.addPage()
            pdf.addImage(imgData, 'PNG', 0, 0, 794, 1123)
        }

        pdf.save(`ledger_${slugify(customerName.value)}_${timestamp()}.pdf`)
        push.success('PDF exported')
    } catch (err) {
        console.error('[LEDGER] PDF export error:', err)
        push.error('Failed to generate PDF')
    } finally {
        exporting.value = false
    }
}

watch(() => props.customerId, (id) => {
    ledger.setCustomer(id)
}, { immediate: true })

onMounted(() => {
    ledger.setCustomer(props.customerId)
})

onUnmounted(() => {
    ledger.setCustomer(null)
    ledger.resetFilters()
    ledger.setCustomerView(false)
})
</script>

<style scoped>
.pdf-render-host {
    position: fixed;
    left: -10000px;
    top: 0;
    width: 210mm;
    pointer-events: none;
    opacity: 0;
}
</style>