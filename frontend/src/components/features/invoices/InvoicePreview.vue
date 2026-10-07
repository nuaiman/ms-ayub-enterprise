<!-- src/components/features/invoices/InvoicePreview.vue -->
<template>
    <div class="space-y-6">
        <!-- ============================================================ -->
        <!-- SUMMARY / DETAILS PANEL                                        -->
        <!-- ============================================================ -->
        <div class="p-6 rounded-xl bg-(--color-muted-bg)/20 border border-(--color-border)">
            <div class="flex items-start gap-3 mb-4">
                <div
                    class="w-10 h-10 rounded-xl bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-base font-semibold text-(--color-text-primary)">Invoice Details</h2>
                    <p class="text-sm text-(--color-text-secondary)">Review and print the invoice</p>
                </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="p-3 rounded-lg bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Invoice Number</p>
                    <p class="text-sm font-semibold text-(--color-text-primary) mt-1">INV-{{ invoice?.id ?? '—' }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-surface) border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Invoice Date</p>
                    <p class="text-sm font-semibold text-(--color-text-primary) mt-1">
                        {{ formatDate(invoice?.created_at) }}
                    </p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-surface) border border-(--color-border) sm:col-span-2">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Bill To</p>
                    <p class="text-sm font-semibold text-(--color-text-primary) mt-1">{{ party.name }}</p>
                    <p v-if="party.phone" class="text-xs text-(--color-text-secondary) mt-0.5">{{ party.phone }}</p>
                    <p v-if="party.address" class="text-xs text-(--color-text-secondary) mt-0.5">{{ party.address }}</p>
                </div>
            </div>
        </div>

        <!-- ============================================================ -->
        <!-- PREVIEW PANEL                                                  -->
        <!-- ============================================================ -->
        <div class="rounded-xl border border-(--color-border) bg-(--color-surface) overflow-hidden">
            <div
                class="flex items-center justify-between p-4 border-b border-(--color-border) bg-(--color-muted-bg)/10">
                <div>
                    <h2 class="text-base font-semibold text-(--color-text-primary)">Invoice Preview</h2>
                    <p class="text-sm text-(--color-text-secondary)">A4 · {{ totalPages }} page{{ totalPages === 1 ? ''
                        :
                        's' }}</p>
                </div>
                <span
                    class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium border border-(--color-green) text-(--color-green)">
                    <span class="w-1.5 h-1.5 rounded-full bg-(--color-green)"></span>
                    Ready
                </span>
            </div>

            <div class="p-4 bg-(--color-muted-bg)/10 overflow-auto">
                <div class="invoice-print-host">
                    <div v-for="(page, pageIdx) in pages" :key="pageIdx" class="invoice-page">
                        <!-- ============= HEADER (every page) ============= -->
                        <header class="invoice-header">
                            <div class="invoice-header-left">
                                <div class="invoice-logo">
                                    <img src="@/assets/logo.png" alt="Logo" />
                                </div>
                                <div class="invoice-header-text">
                                    <h1 class="invoice-company-name">{{ company.name }}</h1>
                                    <p class="invoice-company-line">{{ company.address }}</p>
                                    <p class="invoice-company-line">{{ company.contact }}</p>
                                    <p class="invoice-company-line">{{ company.email }}</p>
                                </div>
                            </div>
                            <div class="invoice-header-right">
                                <div class="invoice-title">INVOICE</div>
                                <div class="invoice-meta-line">
                                    <span class="invoice-meta-key">Invoice #</span>
                                    <span class="invoice-meta-val">INV-{{ invoice?.id ?? '—' }}</span>
                                </div>
                                <div class="invoice-meta-line">
                                    <span class="invoice-meta-key">Date</span>
                                    <span class="invoice-meta-val">{{ formatDate(invoice?.created_at) }}</span>
                                </div>
                                <div class="invoice-meta-line">
                                    <span class="invoice-meta-key">Page</span>
                                    <span class="invoice-meta-val">{{ pageIdx + 1 }} of {{ totalPages }}</span>
                                </div>
                            </div>
                        </header>

                        <!-- ============= BILL TO (page 1, full width) ============= -->
                        <section v-if="pageIdx === 0" class="invoice-parties">
                            <div class="invoice-party-block">
                                <div class="invoice-party-label">Bill To</div>
                                <div class="invoice-party-name">{{ party.name }}</div>
                                <div class="invoice-party-details">
                                    <div v-if="party.phone" class="invoice-party-line">
                                        <span class="invoice-party-line-key">Phone</span>
                                        <span>{{ party.phone }}</span>
                                    </div>
                                    <div v-if="party.email" class="invoice-party-line">
                                        <span class="invoice-party-line-key">Email</span>
                                        <span>{{ party.email }}</span>
                                    </div>
                                    <div v-if="party.address" class="invoice-party-line invoice-party-line-block">
                                        <span class="invoice-party-line-key">Address</span>
                                        <span>{{ party.address }}</span>
                                    </div>
                                </div>
                            </div>
                        </section>

                        <!-- Continuation banner (page 2+) -->
                        <div v-else class="invoice-continuation">
                            <div class="invoice-continuation-left">
                                <div class="invoice-continuation-title">
                                    Invoice — {{ party.name }}
                                </div>
                                <div class="invoice-continuation-sub">
                                    <span v-if="invoice?.id">Ref: INV-{{ invoice.id }}</span>
                                </div>
                            </div>
                            <div class="invoice-continuation-right">
                                Page {{ pageIdx + 1 }} of {{ totalPages }}
                            </div>
                        </div>

                        <!-- ============= ITEMS TABLE (stretches) ============= -->
                        <section class="invoice-table-wrap">
                            <table class="invoice-table">
                                <thead>
                                    <tr>
                                        <th class="col-num">#</th>
                                        <th class="col-type">Type</th>
                                        <th class="col-ref">Reference</th>
                                        <th class="col-amt">Amount</th>
                                    </tr>
                                </thead>
                            </table>
                            <div class="invoice-table-body">
                                <table class="invoice-table">
                                    <tbody>
                                        <tr v-if="page.items.length === 0">
                                            <td colspan="4" class="invoice-table-empty">
                                                No bills recorded on this invoice.
                                            </td>
                                        </tr>
                                        <tr v-for="row in page.items" :key="row.index">
                                            <td class="col-num">{{ row.index + 1 }}</td>
                                            <td class="col-type capitalize">{{ billTypeLabel(row.item.bill_type) }}
                                            </td>
                                            <td class="col-ref">{{ billReference(row.item) }}</td>
                                            <td class="col-amt">{{ formatCurrency(row.item.amount) }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </section>

                        <!-- ============= DISCOUNTS (last page) ============= -->
                        <section v-if="page.isLast && discounts.length > 0" class="invoice-discounts">
                            <h2 class="invoice-section-title">Discounts</h2>
                            <table class="invoice-table invoice-table-compact">
                                <thead>
                                    <tr>
                                        <th class="col-num">#</th>
                                        <th>Description</th>
                                        <th class="col-disc-type">Type</th>
                                        <th class="col-amt">Amount</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="(d, idx) in discounts" :key="d.id">
                                        <td class="col-num">{{ idx + 1 }}</td>
                                        <td class="desc-cell">{{ d.reason || 'Discount' }}</td>
                                        <td class="capitalize">
                                            {{ d.type === 'percent' ? `${d.value}%` : 'Flat' }}
                                        </td>
                                        <td class="col-amt">− {{ formatCurrency(d.computed_amount) }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </section>

                        <!-- ============= BOTTOM ROW (last page) ============= -->
                        <section v-if="page.isLast" class="invoice-bottom-row">
                            <div class="invoice-bottom-left">
                                <div class="invoice-signature-slot">
                                    <div class="invoice-signature-line"></div>
                                    <div class="invoice-signature-label">Authorized Signature</div>
                                    <div v-if="preparedByName" class="invoice-signature-sub">
                                        Prepared by {{ preparedByName }}
                                    </div>
                                </div>
                            </div>

                            <div class="invoice-bottom-right">
                                <div class="invoice-totals-box">
                                    <div class="invoice-totals-row">
                                        <span class="invoice-totals-label">Subtotal</span>
                                        <span class="invoice-totals-value">
                                            {{ formatCurrency(invoice?.subtotal || 0) }}
                                        </span>
                                    </div>
                                    <div class="invoice-totals-row">
                                        <span class="invoice-totals-label">Discount</span>
                                        <span class="invoice-totals-value invoice-totals-value-red">
                                            − {{ formatCurrency(invoice?.discount_amount || 0) }}
                                        </span>
                                    </div>
                                    <div class="invoice-totals-row invoice-totals-row-grand">
                                        <span class="invoice-totals-label-lg">Grand Total</span>
                                        <span class="invoice-totals-value-lg">
                                            {{ formatCurrency(invoice?.total || 0) }}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </section>

                        <!-- ============= NOTES (last page) ============= -->
                        <section v-if="page.isLast && invoice?.notes" class="invoice-notes-section">
                            <h2 class="invoice-section-title">Notes</h2>
                            <p class="invoice-notes">{{ invoice.notes }}</p>
                        </section>

                        <!-- ============= FOOTER (every page) ============= -->
                        <footer class="invoice-footer">
                            <p class="invoice-footer-main">{{ company.name }} · {{ company.subtitle }}</p>
                            <p class="invoice-footer-sub">Generated on {{ formatDate(new Date().toISOString()) }}</p>
                        </footer>
                    </div>
                </div>
            </div>
        </div>

        <!-- ============================================================ -->
        <!-- ACTIONS                                                        -->
        <!-- ============================================================ -->
        <div class="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('back')"
                class="w-full sm:w-auto px-5 py-2.5 text-sm font-medium rounded-xl border border-(--color-border) hover:bg-(--color-muted-bg) transition-all duration-200 flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                Back
            </button>

            <button @click="handlePrint"
                class="w-full sm:w-auto px-6 py-2.5 bg-(--color-blue) text-white rounded-xl text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                </svg>
                Print Invoice
            </button>
        </div>
    </div>

    <!-- ============================================================================ -->
    <!-- PRINT-ONLY PORTAL                                                            -->
    <!-- ============================================================================ -->
    <Teleport to="body">
        <div v-if="printMode" class="invoice-print-root">
            <div v-for="(page, pageIdx) in pages" :key="'p-' + pageIdx" class="invoice-page">
                <header class="invoice-header">
                    <div class="invoice-header-left">
                        <div class="invoice-logo">
                            <img src="@/assets/logo.png" alt="Logo" />
                        </div>
                        <div class="invoice-header-text">
                            <h1 class="invoice-company-name">{{ company.name }}</h1>
                            <p class="invoice-company-line">{{ company.address }}</p>
                            <p class="invoice-company-line">{{ company.contact }}</p>
                            <p class="invoice-company-line">{{ company.email }}</p>
                        </div>
                    </div>
                    <div class="invoice-header-right">
                        <div class="invoice-title">INVOICE</div>
                        <div class="invoice-meta-line">
                            <span class="invoice-meta-key">Invoice #</span>
                            <span class="invoice-meta-val">INV-{{ invoice?.id ?? '—' }}</span>
                        </div>
                        <div class="invoice-meta-line">
                            <span class="invoice-meta-key">Date</span>
                            <span class="invoice-meta-val">{{ formatDate(invoice?.created_at) }}</span>
                        </div>
                        <div class="invoice-meta-line">
                            <span class="invoice-meta-key">Page</span>
                            <span class="invoice-meta-val">{{ pageIdx + 1 }} of {{ totalPages }}</span>
                        </div>
                    </div>
                </header>

                <!-- ============= BILL TO (page 1, full width) ============= -->
                <section v-if="pageIdx === 0" class="invoice-parties">
                    <div class="invoice-party-block">
                        <div class="invoice-party-label">Bill To</div>
                        <div class="invoice-party-name">{{ party.name }}</div>
                        <div class="invoice-party-details">
                            <div v-if="party.phone" class="invoice-party-line">
                                <span class="invoice-party-line-key">Phone</span>
                                <span>{{ party.phone }}</span>
                            </div>
                            <div v-if="party.email" class="invoice-party-line">
                                <span class="invoice-party-line-key">Email</span>
                                <span>{{ party.email }}</span>
                            </div>
                            <div v-if="party.address" class="invoice-party-line invoice-party-line-block">
                                <span class="invoice-party-line-key">Address</span>
                                <span>{{ party.address }}</span>
                            </div>
                        </div>
                    </div>
                </section>

                <div v-else class="invoice-continuation">
                    <div class="invoice-continuation-left">
                        <div class="invoice-continuation-title">
                            Invoice — {{ party.name }}
                        </div>
                        <div class="invoice-continuation-sub">
                            <span v-if="invoice?.id">Ref: INV-{{ invoice.id }}</span>
                        </div>
                    </div>
                    <div class="invoice-continuation-right">
                        Page {{ pageIdx + 1 }} of {{ totalPages }}
                    </div>
                </div>

                <section class="invoice-table-wrap">
                    <table class="invoice-table">
                        <thead>
                            <tr>
                                <th class="col-num">#</th>
                                <th class="col-type">Type</th>
                                <th class="col-ref">Reference</th>
                                <th class="col-amt">Amount</th>
                            </tr>
                        </thead>
                    </table>
                    <div class="invoice-table-body">
                        <table class="invoice-table">
                            <tbody>
                                <tr v-if="page.items.length === 0">
                                    <td colspan="4" class="invoice-table-empty">
                                        No bills recorded on this invoice.
                                    </td>
                                </tr>
                                <tr v-for="row in page.items" :key="'pr-' + row.index">
                                    <td class="col-num">{{ row.index + 1 }}</td>
                                    <td class="col-type capitalize">{{ billTypeLabel(row.item.bill_type) }}</td>
                                    <td class="col-ref">{{ billReference(row.item) }}</td>
                                    <td class="col-amt">{{ formatCurrency(row.item.amount) }}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </section>

                <section v-if="page.isLast && discounts.length > 0" class="invoice-discounts">
                    <h2 class="invoice-section-title">Discounts</h2>
                    <table class="invoice-table invoice-table-compact">
                        <thead>
                            <tr>
                                <th class="col-num">#</th>
                                <th>Description</th>
                                <th class="col-disc-type">Type</th>
                                <th class="col-amt">Amount</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(d, idx) in discounts" :key="'pd-' + d.id">
                                <td class="col-num">{{ idx + 1 }}</td>
                                <td class="desc-cell">{{ d.reason || 'Discount' }}</td>
                                <td class="capitalize">
                                    {{ d.type === 'percent' ? `${d.value}%` : 'Flat' }}
                                </td>
                                <td class="col-amt">− {{ formatCurrency(d.computed_amount) }}</td>
                            </tr>
                        </tbody>
                    </table>
                </section>

                <section v-if="page.isLast" class="invoice-bottom-row">
                    <div class="invoice-bottom-left">
                        <div class="invoice-signature-slot">
                            <div class="invoice-signature-line"></div>
                            <div class="invoice-signature-label">Authorized Signature</div>
                            <div v-if="preparedByName" class="invoice-signature-sub">
                                Prepared by {{ preparedByName }}
                            </div>
                        </div>
                    </div>
                    <div class="invoice-bottom-right">
                        <div class="invoice-totals-box">
                            <div class="invoice-totals-row">
                                <span class="invoice-totals-label">Subtotal</span>
                                <span class="invoice-totals-value">
                                    {{ formatCurrency(invoice?.subtotal || 0) }}
                                </span>
                            </div>
                            <div class="invoice-totals-row">
                                <span class="invoice-totals-label">Discount</span>
                                <span class="invoice-totals-value invoice-totals-value-red">
                                    − {{ formatCurrency(invoice?.discount_amount || 0) }}
                                </span>
                            </div>
                            <div class="invoice-totals-row invoice-totals-row-grand">
                                <span class="invoice-totals-label-lg">Grand Total</span>
                                <span class="invoice-totals-value-lg">
                                    {{ formatCurrency(invoice?.total || 0) }}
                                </span>
                            </div>
                        </div>
                    </div>
                </section>

                <section v-if="page.isLast && invoice?.notes" class="invoice-notes-section">
                    <h2 class="invoice-section-title">Notes</h2>
                    <p class="invoice-notes">{{ invoice.notes }}</p>
                </section>

                <footer class="invoice-footer">
                    <p class="invoice-footer-main">{{ company.name }} · {{ company.subtitle }}</p>
                    <p class="invoice-footer-sub">Generated on {{ formatDate(new Date().toISOString()) }}</p>
                </footer>
            </div>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Invoice, InvoiceDetail, InvoiceItem, UnbilledBillType } from '@/types/invoice'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useBrokersStore } from '@/stores/brokers'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
    invoice: (Invoice & Partial<InvoiceDetail>) | null
}>()

const emit = defineEmits<{
    'back': []
}>()

const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const brokersStore = useBrokersStore()
const usersStore = useUsersStore()

const printMode = ref(false)

const company = {
    name: 'M/S. AYUB ENTERPRISE',
    subtitle: 'Commission Agent — Godown & Logistics',
    address: 'Ali Hossain Chairman Building (1st Floor), 958/27, Strand Road, Mazirghat, Chattogram.',
    contact: 'Cell : 01813-397288, 01705-727492, 01793-287709, 01864-106406',
    email: 'E-mail : mdayubenterprise@gmail.com',
}

const items = computed<InvoiceItem[]>(() => props.invoice?.items ?? [])
const discounts = computed(() => props.invoice?.discounts ?? [])

const billTypeLabel = (t: UnbilledBillType): string => {
    const map: Record<UnbilledBillType, string> = {
        customer_store: 'Store',
        customer_delivery: 'Delivery',
        customer_additional: 'Additional',
        majhi: 'Majhi',
        godown: 'Godown',
    }
    return map[t] || t
}

const billReference = (it: InvoiceItem): string => {
    switch (it.bill_type) {
        case 'customer_store':
            return `Customer Store Bill #${it.bill_id}`
        case 'customer_delivery':
            return `Customer Delivery Bill #${it.bill_id}`
        case 'customer_additional':
            return `Customer Additional Bill #${it.bill_id}`
        case 'majhi':
            return `Majhi Bill #${it.bill_id}`
        case 'godown':
            return `Godown Bill #${it.bill_id}`
        default:
            return `#${it.bill_id}`
    }
}

const party = computed(() => {
    const inv = props.invoice
    if (!inv) return { name: '—', phone: '', address: '', email: '' }

    switch (inv.entity_type) {
        case 'customer': {
            const c = customersStore.getCustomerById(inv.entity_id)
            if (!c) return { name: `Customer #${inv.entity_id}`, phone: '', address: '', email: '' }
            return {
                name: c.company_name || c.contact_person || `Customer #${c.id}`,
                phone: c.phone || '',
                address: c.address || '',
                email: c.email || '',
            }
        }
        case 'majhi': {
            const m = majhisStore.majhis.find((x) => x.id === inv.entity_id)
            return {
                name: m?.name || `Majhi #${inv.entity_id}`,
                phone: m?.phone || '',
                address: '',
                email: '',
            }
        }
        case 'godown': {
            const g = godownsStore.godowns.find((x) => x.id === inv.entity_id)
            return {
                name: g?.name || `Godown #${inv.entity_id}`,
                phone: g?.phone || '',
                address: '',
                email: '',
            }
        }
        case 'broker': {
            const b = brokersStore.brokers.find((x) => x.id === inv.entity_id)
            return {
                name: b?.name || `Broker #${inv.entity_id}`,
                phone: b?.phone || '',
                address: '',
                email: '',
            }
        }
        default:
            return { name: `#${inv.entity_id}`, phone: '', address: '', email: '' }
    }
})

const preparedByName = computed(() => {
    const uid = props.invoice?.user_id
    if (!uid) return ''
    const name = usersStore.getUserName(uid)
    if (/^User #\d+$/.test(name)) return ''
    return name
})

const formatDate = (dateStr?: string | null): string => {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    })
}

// ============================================================================
// PAGINATION
// ============================================================================

const ROWS_FIRST_PAGE = 13
const ROWS_CONTINUATION = 20
const ROWS_LAST_PAGE_TAIL = 11

interface ItemRow { index: number; item: InvoiceItem }
interface Page { items: ItemRow[]; isLast: boolean }

const pages = computed<Page[]>(() => {
    const allItems = items.value
    const rows: ItemRow[] = allItems.map((it, i) => ({ index: i, item: it }))

    if (rows.length === 0) {
        return [{ items: [], isLast: true }]
    }

    const result: Page[] = []

    if (rows.length <= ROWS_LAST_PAGE_TAIL) {
        return [{ items: rows, isLast: true }]
    }

    let cursor = 0
    const first = rows.slice(cursor, cursor + ROWS_FIRST_PAGE)
    cursor += first.length
    result.push({ items: first, isLast: false })

    while (cursor < rows.length) {
        const remaining = rows.length - cursor

        if (remaining <= ROWS_LAST_PAGE_TAIL) {
            result.push({ items: rows.slice(cursor, cursor + remaining), isLast: true })
            break
        }

        const slice = rows.slice(cursor, cursor + ROWS_CONTINUATION)
        cursor += slice.length
        const isLast = cursor >= rows.length
        result.push({ items: slice, isLast })
    }

    return result
})

const totalPages = computed(() => pages.value.length)

// ============================================================================
// PRINT
// ============================================================================
const handlePrint = async () => {
    printMode.value = true
    await new Promise((r) => requestAnimationFrame(() => requestAnimationFrame(() => r(null))))
    window.print()
    printMode.value = false
}
</script>

<style scoped>
.invoice-print-host {
    width: 210mm;
    margin: 0 auto;
}
</style>

<style>
/* ============================================================================
   Global styles — print portal lives outside the component subtree.
   ============================================================================ */

.invoice-page {
    width: 210mm;
    height: 297mm;
    padding: 11mm 14mm 0 14mm;
    margin: 0 auto 8mm;
    background: #ffffff;
    color: #111111;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    font-size: 11.5px;
    line-height: 1.45;
    box-sizing: border-box;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
    display: flex;
    flex-direction: column;
    page-break-after: always;
    break-after: page;
    overflow: hidden;
}

.invoice-page:last-child {
    page-break-after: auto;
    break-after: auto;
    margin-bottom: 0;
}

/* ---------- Header ---------- */
.invoice-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 20px;
    padding-bottom: 10px;
    border-bottom: 2px solid #111111;
    margin-bottom: 12px;
    flex-shrink: 0;
}

.invoice-header-left {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    min-width: 0;
}

.invoice-logo {
    width: 46px;
    height: 46px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
}

.invoice-logo img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    display: block;
}

.invoice-header-text {
    min-width: 0;
}

.invoice-company-name {
    font-size: 15px;
    font-weight: 800;
    letter-spacing: 0.04em;
    margin: 0;
    color: #111111;
    text-transform: uppercase;
}

.invoice-company-line {
    font-size: 10px;
    color: #555555;
    margin: 1px 0 0;
    line-height: 1.35;
}

.invoice-header-right {
    text-align: right;
    flex-shrink: 0;
    min-width: 160px;
}

.invoice-title {
    font-size: 24px;
    font-weight: 900;
    letter-spacing: 0.14em;
    color: #111111;
    line-height: 1;
    margin-bottom: 6px;
}

.invoice-meta-line {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    font-size: 10px;
    line-height: 1.5;
}

.invoice-meta-key {
    color: #777777;
}

.invoice-meta-val {
    font-weight: 600;
    color: #111111;
    min-width: 84px;
    text-align: right;
}

/* ---------- Parties (full-width) ---------- */
.invoice-parties {
    display: block;
    margin-bottom: 14px;
    flex-shrink: 0;
}

.invoice-party-block {
    border: 1px solid #d4d4d4;
    border-left: 4px solid #075985;
    border-radius: 4px;
    padding: 10px 14px;
    background: #fafafa;
    width: 100%;
    box-sizing: border-box;
}

.invoice-party-label {
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: #777777;
    margin-bottom: 6px;
}

.invoice-party-name {
    font-size: 13px;
    font-weight: 700;
    color: #111111;
    margin-bottom: 6px;
}

.invoice-party-details {
    display: flex;
    flex-wrap: wrap;
    gap: 4px 24px;
}

.invoice-party-line {
    display: inline-flex;
    align-items: baseline;
    gap: 6px;
    font-size: 10px;
    color: #444444;
    line-height: 1.4;
}

.invoice-party-line-block {
    flex-basis: 100%;
}

.invoice-party-line-key {
    font-size: 8.5px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #999999;
}

/* ---------- Continuation banner ---------- */
.invoice-continuation {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 10px;
    background: #f2f2f2;
    border: 1px solid #d4d4d4;
    border-radius: 3px;
    margin-bottom: 12px;
    flex-shrink: 0;
}

.invoice-continuation-title {
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.04em;
    color: #111111;
}

.invoice-continuation-sub {
    font-size: 9.5px;
    color: #666666;
    margin-top: 1px;
}

.invoice-continuation-right {
    font-size: 9.5px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #666666;
}

/* ============================================================================
   TABLE — head is fixed; body fills the rest of the wrapper.
   ============================================================================ */

.invoice-table-wrap {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    border: 1px solid #d4d4d4;
    border-radius: 3px;
    overflow: hidden;
    margin-bottom: 12px;
}

.invoice-table-body {
    flex: 1;
    min-height: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.invoice-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 11px;
    table-layout: fixed;
}

.invoice-table thead {
    background: #f2f2f2;
    flex-shrink: 0;
}

.invoice-table th {
    text-align: left;
    padding: 8px 8px;
    border-bottom: 1px solid #d4d4d4;
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #555555;
}

.invoice-table td {
    padding: 8px 8px;
    border-bottom: 1px solid #e5e7eb;
    color: #222222;
    vertical-align: top;
}

/* Every row gets a clean horizontal separator */
.invoice-table tbody tr td {
    border-bottom: 1px solid #e5e7eb;
}

.invoice-table tbody tr:last-child td {
    border-bottom: 1px solid #d4d4d4;
}

.invoice-table tbody tr:nth-child(even) {
    background: #fafafa;
}

.invoice-table .col-num {
    width: 8%;
    text-align: center;
}

.invoice-table .col-type {
    width: 16%;
}

.invoice-table .col-ref {
    width: 56%;
}

.invoice-table .col-amt {
    width: 20%;
    text-align: right;
}

.invoice-table .col-disc-type {
    width: 18%;
}

.invoice-table .capitalize {
    text-transform: capitalize;
}

.invoice-table-empty {
    text-align: center;
    color: #999999;
    font-style: italic;
    padding: 14px 8px;
}

.invoice-table-compact {
    font-size: 10.5px;
}

.invoice-table-compact th,
.invoice-table-compact td {
    padding: 6px 8px;
}

.desc-cell {
    white-space: pre-wrap;
    word-break: break-word;
}

/* ---------- Discounts section ---------- */
.invoice-discounts {
    margin-bottom: 12px;
    flex-shrink: 0;
}

.invoice-discounts .invoice-table {
    border: 1px solid #d4d4d4;
    border-radius: 3px;
    overflow: hidden;
}

.invoice-discounts .invoice-table tbody tr td {
    border-bottom: 1px solid #e5e7eb;
}

.invoice-discounts .invoice-table tbody tr:last-child td {
    border-bottom: none;
}

.invoice-discounts .invoice-table tbody tr:nth-child(even) {
    background: #fafafa;
}

.invoice-section-title {
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: #555555;
    margin: 0 0 6px;
}

/* ---------- Bottom row ---------- */
.invoice-bottom-row {
    margin-top: 14px;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    gap: 24px;
    flex-shrink: 0;
    page-break-inside: avoid;
    break-inside: avoid;
}

.invoice-bottom-left {
    flex: 1;
    min-width: 0;
    padding-bottom: 8px;
}

.invoice-bottom-right {
    flex-shrink: 0;
    width: 70mm;
}

.invoice-signature-slot {
    width: 60mm;
    text-align: center;
}

.invoice-signature-line {
    border-top: 1px solid #111111;
    margin-bottom: 4px;
}

.invoice-signature-label {
    font-size: 9.5px;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #111111;
}

.invoice-signature-sub {
    font-size: 9px;
    color: #777777;
    margin-top: 2px;
}

/* ---------- Totals box ---------- */
.invoice-totals-box {
    border: 1px solid #d4d4d4;
    border-radius: 3px;
    overflow: hidden;
}

.invoice-totals-row {
    display: flex;
    justify-content: space-between;
    padding: 7px 11px;
    border-bottom: 1px solid #eeeeee;
    font-size: 11px;
}

.invoice-totals-row:last-child {
    border-bottom: none;
}

.invoice-totals-label {
    color: #666666;
}

.invoice-totals-value {
    font-weight: 600;
    color: #111111;
}

.invoice-totals-value-red {
    color: #b91c1c;
}

.invoice-totals-row-grand {
    border-top: 2px solid #111111;
    background: #fafafa;
    padding: 10px 11px;
}

.invoice-totals-label-lg {
    font-size: 11px;
    font-weight: 800;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: #111111;
}

.invoice-totals-value-lg {
    font-size: 14px;
    font-weight: 800;
    color: #111111;
}

/* ---------- Notes ---------- */
.invoice-notes-section {
    margin-top: 12px;
    padding-top: 8px;
    border-top: 1px solid #e5e7eb;
    flex-shrink: 0;
}

.invoice-notes {
    font-size: 10.5px;
    line-height: 1.5;
    color: #444444;
    white-space: pre-wrap;
    margin: 0;
}

/* ---------- Footer ---------- */
.invoice-footer {
    margin: 12px -14mm 0 -14mm;
    padding: 10px 14mm;
    background: #075985;
    color: white;
    text-align: center;
    flex-shrink: 0;
}

.invoice-footer-main {
    font-size: 8.5px;
    font-weight: 600;
    margin: 0;
}

.invoice-footer-sub {
    font-size: 8.5px;
    opacity: 0.95;
    margin: 2px 0 0;
}

/* ============================================================================
   Print isolation
   ============================================================================ */

body>.invoice-print-root {
    position: absolute !important;
    left: -100000px !important;
    top: 0 !important;
    width: 210mm;
    pointer-events: none;
    opacity: 0;
}

@media print {
    @page {
        size: A4 portrait;
        margin: 0;
    }

    html,
    body {
        background: #ffffff !important;
        margin: 0 !important;
        padding: 0 !important;
        width: 210mm !important;
        overflow: visible !important;
    }

    body * {
        visibility: hidden !important;
    }

    .invoice-print-root,
    .invoice-print-root * {
        visibility: visible !important;
    }

    body>.invoice-print-root {
        position: static !important;
        left: auto !important;
        top: auto !important;
        width: auto !important;
        pointer-events: auto !important;
        opacity: 1 !important;
    }

    .invoice-print-root {
        display: block !important;
        width: 210mm !important;
        margin: 0 !important;
    }

    .invoice-page {
        box-shadow: none !important;
        margin: 0 !important;
        page-break-after: always;
        break-after: page;
    }

    .invoice-page:last-child {
        page-break-after: auto;
        break-after: auto;
    }
}
</style>