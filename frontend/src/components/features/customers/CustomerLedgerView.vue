<!-- src/components/features/customers/CustomerLedgerView.vue -->
<template>
    <div class="ledger-print-root">
        <div v-for="(page, pageIdx) in pages" :key="pageIdx" class="ledger-page">
            <!-- ============================================================ -->
            <!-- HEADER (every page)                                            -->
            <!-- ============================================================ -->
            <header class="ledger-header">
                <div class="ledger-header-left">
                    <div class="ledger-logo">
                        <img src="@/assets/logo.png" alt="Logo" />
                    </div>
                    <div class="ledger-header-text">
                        <h1 class="ledger-company-name">M/S. AYUB ENTERPRISE</h1>
                        <p class="ledger-company-line">Ali Hossain Chairman Building (1st Floor), 958/27, Strand Road,
                            Mazirghat, Chattogram.</p>
                        <p class="ledger-company-line">Cell : 01813-397288, 01705-727492, 01793-287729, 01864-106406</p>
                        <p class="ledger-company-line">E-mail : mdayubenterprise@gmail.com</p>
                    </div>
                </div>
                <div class="ledger-header-right">
                    <div class="ledger-title">STATEMENT OF ACCOUNT</div>
                    <div class="ledger-meta-line">
                        <span class="ledger-meta-key">Generated</span>
                        <span class="ledger-meta-val">{{ generatedDate }}</span>
                    </div>
                    <div class="ledger-meta-line">
                        <span class="ledger-meta-key">Page</span>
                        <span class="ledger-meta-val">{{ pageIdx + 1 }} of {{ pages.length }}</span>
                    </div>
                </div>
            </header>

            <!-- ============================================================ -->
            <!-- BILL TO CARD (page 1, full width)                              -->
            <!-- ============================================================ -->
            <section v-if="pageIdx === 0" class="ledger-parties">
                <div class="ledger-party-block">
                    <div class="ledger-party-top">
                        <div class="ledger-party-left">
                            <div class="ledger-party-label">Statement For</div>
                            <div class="ledger-party-name">{{ customerName }}</div>
                            <div class="ledger-party-details">
                                <div v-if="customerPhone" class="ledger-party-line">
                                    <span class="ledger-party-line-key">Phone</span>
                                    <span>{{ customerPhone }}</span>
                                </div>
                                <div v-if="customerAddress" class="ledger-party-line ledger-party-line-block">
                                    <span class="ledger-party-line-key">Address</span>
                                    <span>{{ customerAddress }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="ledger-party-right">
                            <div class="ledger-party-label">Statement Period</div>
                            <div class="ledger-party-period">{{ periodLabel }}</div>
                            <div class="ledger-party-line ledger-party-line-right">
                                <span class="ledger-party-line-key">Events</span>
                                <span>{{ eventCount }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="ledger-party-summary">
                        <div class="ledger-party-summary-label">Summary</div>
                        <div class="ledger-party-summary-grid">
                            <div class="ledger-party-summary-item">
                                <span class="ledger-party-summary-key">Billed</span>
                                <span class="ledger-party-summary-val">BDT {{ formatMoney(summary.totalBilled) }}</span>
                            </div>
                            <div class="ledger-party-summary-item">
                                <span class="ledger-party-summary-key">Paid</span>
                                <span class="ledger-party-summary-val">BDT {{ formatMoney(summary.totalPaid) }}</span>
                            </div>
                            <div class="ledger-party-summary-item ledger-party-summary-item-strong">
                                <span class="ledger-party-summary-key">Outstanding</span>
                                <span class="ledger-party-summary-val">BDT {{ formatMoney(summary.outstanding) }}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Continuation banner (page 2+) -->
            <div v-else class="ledger-continuation">
                <div class="ledger-continuation-left">
                    <div class="ledger-continuation-title">
                        Statement of Account — {{ customerName }}
                    </div>
                    <div class="ledger-continuation-sub">
                        {{ periodLabel }}
                    </div>
                </div>
                <div class="ledger-continuation-right">
                    Page {{ pageIdx + 1 }} of {{ pages.length }}
                </div>
            </div>

            <!-- ============================================================ -->
            <!-- ITEMS TABLE                                                    -->
            <!-- ============================================================ -->
            <section class="ledger-table-wrap">
                <table class="ledger-table">
                    <thead>
                        <tr>
                            <th class="col-date">Date</th>
                            <th class="col-type">Type</th>
                            <th class="col-ref">Reference</th>
                            <th class="col-desc">Description</th>
                            <th class="col-amt text-right">Debit</th>
                            <th class="col-amt text-right">Credit</th>
                            <th class="col-amt text-right">Balance</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="page.rows.length === 0">
                            <td colspan="7" class="ledger-table-empty">
                                No events in the selected period.
                            </td>
                        </tr>
                        <tr v-for="row in page.rows" :key="row.event.id">
                            <td class="col-date">{{ shortDate(row.event.date) }}</td>
                            <td class="col-type">{{ typeLabel(row.event.type) }}</td>
                            <td class="col-ref">{{ row.event.reference || '—' }}</td>
                            <td class="col-desc">
                                <div class="desc-title">{{ row.event.title }}</div>
                                <div v-if="row.event.description" class="desc-sub">{{ row.event.description }}</div>
                            </td>
                            <td class="col-amt text-right">
                                <span v-if="row.event.amountKind === 'debit'">
                                    {{ formatMoney(row.event.amount || 0) }}
                                </span>
                                <span v-else class="dash">—</span>
                            </td>
                            <td class="col-amt text-right">
                                <span v-if="row.event.amountKind === 'credit'">
                                    {{ formatMoney(row.event.amount || 0) }}
                                </span>
                                <span v-else class="dash">—</span>
                            </td>
                            <td class="col-amt text-right balance-cell">
                                {{ formatMoney(row.balance) }}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </section>

            <!-- ============================================================ -->
            <!-- SIGNATURE (left) + SUMMARY (right) — last page only            -->
            <!-- ============================================================ -->
            <section v-if="page.isLast" class="ledger-bottom-row">
                <div class="ledger-bottom-left">
                    <div class="ledger-signature-slot">
                        <div class="ledger-signature-line"></div>
                        <div class="ledger-signature-label">Authorized Signature</div>
                    </div>
                </div>

                <div class="ledger-bottom-right">
                    <div class="ledger-totals-box">
                        <div class="ledger-totals-row">
                            <span class="ledger-totals-label">Total Billed</span>
                            <span class="ledger-totals-value">
                                {{ formatMoney(summary.totalBilled) }}
                            </span>
                        </div>
                        <div class="ledger-totals-row">
                            <span class="ledger-totals-label">Total Paid</span>
                            <span class="ledger-totals-value">
                                {{ formatMoney(summary.totalPaid) }}
                            </span>
                        </div>
                        <div class="ledger-totals-row ledger-totals-row-grand">
                            <span class="ledger-totals-label-lg">Outstanding</span>
                            <span class="ledger-totals-value-lg">
                                {{ formatMoney(summary.outstanding) }}
                            </span>
                        </div>
                    </div>
                </div>
            </section>

            <!-- ============================================================ -->
            <!-- FOOTER (every page)                                            -->
            <!-- ============================================================ -->
            <footer class="ledger-footer">
                <p class="ledger-footer-main">M/S. AYUB ENTERPRISE · Commission Agent — Godown &amp; Logistics</p>
                <p class="ledger-footer-sub">Generated on {{ generatedDate }} · This is a computer-generated statement.
                </p>
            </footer>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { LedgerEvent, LedgerEventType } from '@/stores/customerLedger'

const props = defineProps<{
    customerName: string
    customerPhone?: string | null
    customerAddress?: string | null
    events: LedgerEvent[]
    summary: {
        totalBilled: number
        totalPaid: number
        outstanding: number
        eventCount: number
    }
    dateFrom?: string
    dateTo?: string
}>()

// ---------------------------------------------------------------------------
// Labels
// ---------------------------------------------------------------------------

const typeLabelMap: Record<LedgerEventType, string> = {
    customer_created: 'Customer',
    lot_created: 'Lot',
    store_created: 'Store',
    delivery_created: 'Delivery',
    delivery_item_added: 'Delivery Item',
    damage_recorded: 'Damage',
    lot_transferred_out: 'Lot Out',
    lot_transferred_in: 'Lot In',
    customer_store_bill: 'Store Bill',
    customer_delivery_bill: 'Delivery Bill',
    customer_additional_bill: 'Additional Bill',
    invoice_created: 'Invoice',
    payment_received: 'Payment',
}

const typeLabel = (t: LedgerEventType) => typeLabelMap[t] || t

const eventCount = computed(() => props.events.length)

const generatedDate = computed(() =>
    new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
)

const periodLabel = computed(() => {
    if (props.dateFrom && props.dateTo) {
        return `${shortDate(props.dateFrom)} — ${shortDate(props.dateTo)}`
    }
    if (props.dateFrom) return `From ${shortDate(props.dateFrom)}`
    if (props.dateTo) return `Up to ${shortDate(props.dateTo)}`
    return 'All time'
})

// ---------------------------------------------------------------------------
// Formatting
// ---------------------------------------------------------------------------

const formatMoney = (val: number): string => {
    return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    }).format(val || 0)
}

const shortDate = (dateStr: string): string => {
    const d = new Date(dateStr)
    return d.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: '2-digit' })
}

// ---------------------------------------------------------------------------
// Rows + running balance
// ---------------------------------------------------------------------------

interface Row {
    event: LedgerEvent
    balance: number
}

interface Page {
    rows: Row[]
    isLast: boolean
}

const allRows = computed<Row[]>(() => {
    const sorted = [...props.events].sort(
        (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
    )

    let running = 0
    const rows: Row[] = []
    for (const e of sorted) {
        if (e.amountKind === 'debit' && typeof e.amount === 'number') {
            running += e.amount
        } else if (e.amountKind === 'credit' && typeof e.amount === 'number') {
            running -= e.amount
        }
        rows.push({ event: e, balance: running })
    }
    return rows
})

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

const ROWS_FIRST_PAGE = 15
const ROWS_CONTINUATION = 24
const ROWS_LAST_PAGE = 16

const pages = computed<Page[]>(() => {
    const rows = allRows.value

    if (rows.length === 0) {
        return [{ rows: [], isLast: true }]
    }

    const result: Page[] = []
    let cursor = 0

    if (rows.length <= ROWS_LAST_PAGE) {
        return [{ rows, isLast: true }]
    }

    const first = rows.slice(cursor, cursor + ROWS_FIRST_PAGE)
    cursor += first.length
    result.push({ rows: first, isLast: false })

    while (cursor < rows.length) {
        const remaining = rows.length - cursor

        if (remaining <= ROWS_LAST_PAGE) {
            const slice = rows.slice(cursor, cursor + remaining)
            cursor += slice.length
            result.push({ rows: slice, isLast: true })
            break
        }

        const slice = rows.slice(cursor, cursor + ROWS_CONTINUATION)
        cursor += slice.length
        const isLast = cursor >= rows.length
        result.push({ rows: slice, isLast })
    }

    return result
})
</script>

<style>
/* ============================================================================
   PRINT ROOT
   ============================================================================ */

.ledger-print-root {
    width: 210mm;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 0;
}

/* ---------- A4 page ---------- */
.ledger-page {
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

.ledger-page:last-child {
    page-break-after: auto;
    break-after: auto;
    margin-bottom: 0;
}

/* ---------- Header ---------- */
.ledger-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 20px;
    padding-bottom: 10px;
    border-bottom: 2px solid #111111;
    margin-bottom: 12px;
    flex-shrink: 0;
}

.ledger-header-left {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    min-width: 0;
}

.ledger-logo {
    width: 46px;
    height: 46px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
}

.ledger-logo img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    display: block;
}

.ledger-header-text {
    min-width: 0;
}

.ledger-company-name {
    font-size: 15px;
    font-weight: 800;
    letter-spacing: 0.04em;
    margin: 0;
    color: #111111;
    text-transform: uppercase;
}

.ledger-company-line {
    font-size: 10px;
    color: #555555;
    margin: 1px 0 0;
    line-height: 1.35;
}

.ledger-header-right {
    text-align: right;
    flex-shrink: 0;
    min-width: 160px;
}

.ledger-title {
    font-size: 18px;
    font-weight: 900;
    letter-spacing: 0.12em;
    color: #111111;
    line-height: 1;
    margin-bottom: 6px;
    text-transform: uppercase;
}

.ledger-meta-line {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    font-size: 10px;
    line-height: 1.5;
}

.ledger-meta-key {
    color: #777777;
}

.ledger-meta-val {
    font-weight: 600;
    color: #111111;
    min-width: 84px;
    text-align: right;
}

/* ---------- Bill To card (full width, page 1) ---------- */
.ledger-parties {
    display: block;
    margin-bottom: 14px;
    flex-shrink: 0;
}

.ledger-party-block {
    border: 1px solid #d4d4d4;
    border-left: 4px solid #075985;
    border-radius: 4px;
    padding: 10px 14px;
    background: #fafafa;
    width: 100%;
    box-sizing: border-box;
}

.ledger-party-top {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 24px;
}

.ledger-party-left {
    flex: 1;
    min-width: 0;
}

.ledger-party-right {
    flex-shrink: 0;
    text-align: right;
    min-width: 60mm;
}

.ledger-party-label {
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: #777777;
    margin-bottom: 6px;
}

.ledger-party-name {
    font-size: 13px;
    font-weight: 700;
    color: #111111;
    margin-bottom: 6px;
}

.ledger-party-period {
    font-size: 12px;
    font-weight: 700;
    color: #111111;
    margin-bottom: 4px;
}

.ledger-party-details {
    display: flex;
    flex-wrap: wrap;
    gap: 4px 24px;
}

.ledger-party-line {
    display: inline-flex;
    align-items: baseline;
    gap: 6px;
    font-size: 10px;
    color: #444444;
    line-height: 1.4;
}

.ledger-party-line-right {
    justify-content: flex-end;
}

.ledger-party-line-block {
    flex-basis: 100%;
}

.ledger-party-line-key {
    font-size: 8.5px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #999999;
}

/* Summary sub-block inside the Bill To card */
.ledger-party-summary {
    margin-top: 10px;
    padding-top: 8px;
    border-top: 1px solid #e2e8f0;
}

.ledger-party-summary-label {
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: #777777;
    margin-bottom: 6px;
}

.ledger-party-summary-grid {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 4px 24px;
}

.ledger-party-summary-item {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 8px;
    font-size: 10px;
    color: #444444;
    line-height: 1.45;
}

.ledger-party-summary-item-strong {
    padding-left: 12px;
    border-left: 2px solid #111111;
    margin-left: -2px;
}

.ledger-party-summary-key {
    color: #666666;
}

.ledger-party-summary-val {
    color: #111111;
    font-weight: 600;
    font-variant-numeric: tabular-nums;
}

.ledger-party-summary-item-strong .ledger-party-summary-key {
    color: #111111;
    font-weight: 700;
}

.ledger-party-summary-item-strong .ledger-party-summary-val {
    color: #111111;
    font-weight: 800;
}

/* ---------- Continuation banner ---------- */
.ledger-continuation {
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

.ledger-continuation-title {
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.04em;
    color: #111111;
}

.ledger-continuation-sub {
    font-size: 9.5px;
    color: #666666;
    margin-top: 1px;
}

.ledger-continuation-right {
    font-size: 9.5px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #666666;
}

/* ============================================================================
   TABLE
   ============================================================================ */

.ledger-table-wrap {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    border: 1px solid #d4d4d4;
    border-radius: 3px;
    overflow: hidden;
    margin-bottom: 12px;
}

.ledger-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 11px;
    table-layout: fixed;
}

.ledger-table thead {
    background: #f2f2f2;
    flex-shrink: 0;
}

.ledger-table th {
    text-align: left;
    padding: 8px 8px;
    border-bottom: 1px solid #d4d4d4;
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #555555;
}

.ledger-table td {
    padding: 8px 8px;
    border-bottom: 1px solid #e5e7eb;
    color: #222222;
    vertical-align: top;
}

.ledger-table tbody tr td {
    border-bottom: 1px solid #e5e7eb;
}

.ledger-table tbody tr:last-child td {
    border-bottom: 1px solid #d4d4d4;
}

.ledger-table tbody tr:nth-child(even) {
    background: #fafafa;
}

.ledger-table .text-right {
    text-align: right;
}

.ledger-table .col-date {
    width: 11%;
}

.ledger-table .col-type {
    width: 13%;
}

.ledger-table .col-ref {
    width: 14%;
}

.ledger-table .col-desc {
    width: 32%;
}

.ledger-table .col-amt {
    width: 10%;
    text-align: right;
    font-variant-numeric: tabular-nums;
}

.desc-title {
    font-weight: 600;
    color: #111111;
    line-height: 1.3;
}

.desc-sub {
    font-size: 9px;
    color: #666666;
    margin-top: 2px;
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.35;
}

.balance-cell {
    font-weight: 700;
    color: #111111;
}

.dash {
    color: #cbd5e1;
}

.ledger-table-empty {
    text-align: center;
    color: #999999;
    font-style: italic;
    padding: 14px 8px;
}

/* ---------- Bottom row (signature left + totals right) ---------- */
.ledger-bottom-row {
    margin-top: 14px;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    gap: 24px;
    flex-shrink: 0;
    page-break-inside: avoid;
    break-inside: avoid;
}

.ledger-bottom-left {
    flex: 1;
    min-width: 0;
    padding-bottom: 8px;
}

.ledger-bottom-right {
    flex-shrink: 0;
    width: 70mm;
}

.ledger-signature-slot {
    width: 60mm;
    text-align: center;
    padding-top: 12mm;
}

.ledger-signature-line {
    border-top: 1px solid #111111;
    margin-bottom: 4px;
}

.ledger-signature-label {
    font-size: 9.5px;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #111111;
}

/* ---------- Totals box ---------- */
.ledger-totals-box {
    border: 1px solid #d4d4d4;
    border-radius: 3px;
    overflow: hidden;
}

.ledger-totals-row {
    display: flex;
    justify-content: space-between;
    padding: 7px 11px;
    border-bottom: 1px solid #eeeeee;
    font-size: 11px;
}

.ledger-totals-row:last-child {
    border-bottom: none;
}

.ledger-totals-label {
    color: #666666;
}

.ledger-totals-value {
    font-weight: 600;
    color: #111111;
    font-variant-numeric: tabular-nums;
}

.ledger-totals-row-grand {
    border-top: 2px solid #111111;
    background: #fafafa;
    padding: 10px 11px;
}

.ledger-totals-label-lg {
    font-size: 11px;
    font-weight: 800;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: #111111;
}

.ledger-totals-value-lg {
    font-size: 14px;
    font-weight: 800;
    color: #111111;
    font-variant-numeric: tabular-nums;
}

/* ---------- Footer ---------- */
.ledger-footer {
    margin: 12px -14mm 0 -14mm;
    padding: 10px 14mm;
    background: #075985;
    color: #ffffff;
    text-align: center;
    flex-shrink: 0;
}

.ledger-footer-main {
    font-size: 8.5px;
    font-weight: 600;
    margin: 0;
}

.ledger-footer-sub {
    font-size: 8.5px;
    opacity: 0.95;
    margin: 2px 0 0;
}

/* ============================================================================
   PRINT ISOLATION
   ============================================================================ */

body>.ledger-print-root {
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
        height: auto !important;
        min-height: 0 !important;
        overflow: visible !important;
    }

    /* LEADING-BLANK FIX */
    body>*:not(.ledger-print-root) {
        display: none !important;
    }

    #app,
    body>#app {
        display: none !important;
    }

    .ledger-print-root,
    .ledger-print-root * {
        visibility: visible !important;
    }

    body>.ledger-print-root,
    body>#app>.ledger-print-root {
        position: static !important;
        left: auto !important;
        top: auto !important;
        width: 210mm !important;
        pointer-events: auto !important;
        opacity: 1 !important;
        margin: 0 !important;
        padding: 0 !important;
        display: block !important;
    }

    .ledger-print-root {
        display: block !important;
        width: 210mm !important;
        margin: 0 !important;
        padding: 0 !important;
    }

    /* TRAILING-BLANK FIX */
    .ledger-page {
        box-shadow: none !important;
        margin: 0 !important;
        margin-bottom: 0 !important;
        page-break-after: always;
        break-after: page;
        height: 297mm !important;
        max-height: 297mm !important;
        overflow: hidden !important;
    }

    .ledger-page:last-child {
        page-break-after: auto !important;
        break-after: auto !important;
        page-break-inside: avoid;
        break-inside: avoid;
        margin: 0 !important;
        margin-bottom: 0 !important;
    }
}
</style>