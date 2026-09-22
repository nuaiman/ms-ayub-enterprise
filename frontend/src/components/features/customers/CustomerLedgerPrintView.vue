<!-- src/components/features/customers/CustomerLedgerPrintView.vue -->
<template>
    <div class="ledger-print-root">
        <div v-for="(page, pageIdx) in pages" :key="pageIdx" class="ledger-page" :ref="el => setPageRef(el, pageIdx)">
            <!-- ============ HEADER ============ -->
            <header class="page-header">
                <div v-if="pageIdx === 0" class="company-header">
                    <div class="logo-area">
                        <div class="logo-box">
                            <img src="@/assets/logo.png" alt="logo" />
                        </div>
                        <div>
                            <div class="company-name">M/S. AYUB ENTERPRISE</div>
                            <div class="company-divider"></div>
                        </div>
                    </div>
                    <div>
                        <div class="business-type">Commission Agent</div>
                        <div class="subtitle">Godown &amp; Logistics</div>
                    </div>
                </div>

                <div v-if="pageIdx === 0" class="meta-row">
                    <div class="ref-box">
                        <span class="ref-label">Statement For</span>
                        <span class="ref-value">{{ customerName }}</span>
                    </div>
                    <div class="title-center">
                        <h2>Statement of Account</h2>
                        <div class="title-underline"></div>
                    </div>
                    <div class="date-box">
                        <span class="label">Generated</span>
                        <span class="date-value">{{ generatedDate }}</span>
                    </div>
                </div>

                <div v-if="pageIdx === 0" class="bill-grid">
                    <div class="bill-box">
                        <div class="bill-label">Customer</div>
                        <div class="bill-value">
                            <div class="customer-name">{{ customerName }}</div>
                            <div v-if="customerPhone" class="customer-sub">{{ customerPhone }}</div>
                            <div v-if="customerAddress" class="customer-sub">{{ customerAddress }}</div>
                        </div>
                    </div>
                    <div class="bill-box">
                        <div class="bill-label">Period</div>
                        <div class="bill-value">
                            <div class="customer-name">{{ periodLabel }}</div>
                            <div class="customer-sub">{{ eventCount }} event(s)</div>
                        </div>
                    </div>
                </div>

                <div v-if="pageIdx === 0" class="summary-grid">
                    <div class="summary-cell">
                        <div class="summary-label">Total Billed</div>
                        <div class="summary-value">BDT. {{ formatMoney(summary.totalBilled) }}</div>
                    </div>
                    <div class="summary-cell">
                        <div class="summary-label">Total Paid</div>
                        <div class="summary-value paid">BDT. {{ formatMoney(summary.totalPaid) }}</div>
                    </div>
                    <div class="summary-cell">
                        <div class="summary-label">Outstanding</div>
                        <div class="summary-value outstanding">BDT. {{ formatMoney(summary.outstanding) }}</div>
                    </div>
                </div>

                <div v-else class="continuation-header">
                    <div class="continuation-left">
                        <div class="continuation-title">Statement of Account</div>
                        <div class="continuation-sub">{{ customerName }}</div>
                    </div>
                    <div class="continuation-right">
                        Page {{ pageIdx + 1 }} of {{ pages.length }}
                    </div>
                </div>
            </header>

            <!-- ============ CONTENT ============ -->
            <main class="page-content">
                <div class="timeline-header">
                    <div class="th-date">Date</div>
                    <div class="th-event">
                        <span class="th-type">Type</span>
                        <span class="th-desc">Event</span>
                        <span class="th-amount">Amount</span>
                    </div>
                </div>

                <div v-if="page.length === 0" class="empty-state">
                    No events in the selected period.
                </div>

                <div v-else class="timeline">
                    <div v-for="event in page" :key="event.id" class="timeline-row">
                        <div class="timeline-date">
                            {{ shortDate(event.date) }}
                        </div>

                        <div class="timeline-body">
                            <div class="row-head">
                                <div class="row-head-left">
                                    <span class="type-chip">{{ typeLabel(event.type) }}</span>
                                    <span v-if="event.reference" class="ref-inline">{{ event.reference }}</span>
                                </div>
                                <div v-if="typeof event.amount === 'number' && event.amount !== 0" class="row-amount"
                                    :class="event.amountKind">
                                    {{ amountSign(event.amountKind) }}{{ formatMoney(event.amount) }}
                                    <span v-if="event.amountLabel" class="amount-label">{{ event.amountLabel }}</span>
                                </div>
                            </div>

                            <div class="row-title">{{ event.title }}</div>
                            <div class="row-desc">{{ event.description }}</div>

                            <div v-if="event.meta && Object.keys(event.meta).length > 0" class="row-meta">
                                <span v-for="(value, key) in event.meta" :key="key" class="meta-item">
                                    <span class="meta-key">{{ humanKey(String(key)) }}:</span>
                                    <span class="meta-val">{{ value ?? 'â€”' }}</span>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            </main>

            <!-- ============ FOOTER ============ -->
            <footer class="page-footer">
                <p class="footer-line-1">Ali Hossain Chairman Building (1st Floor), 958/27, Strand Road, Mazirghat,
                    Chattogram.</p>
                <p class="footer-line-2">Cell : 01813-397288, 01705-727492, 01793-287709, 01864-106406 â€¢ E-mail :
                    mdayubenterprise@gmail.com</p>
            </footer>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, type ComponentPublicInstance } from 'vue'
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
    customerView?: boolean
}>()

const pageRefs = ref<HTMLElement[]>([])
const setPageRef = (el: Element | ComponentPublicInstance | null, idx: number) => {
    if (el instanceof HTMLElement) {
        pageRefs.value[idx] = el
    }
}
defineExpose({ pageRefs })

const typeLabelMap: Record<LedgerEventType, string> = {
    customer_created: 'Customer',
    lot_created: 'Lot',
    store_created: 'Store',
    delivery_created: 'Delivery',
    delivery_item_added: 'Delivery Item',
    damage_recorded: 'Damage',
    transport_created: 'Transport',
    vehicle_added: 'Vehicle',
    storage_payment: 'Storage Payment',
    unload_payment: 'Unload Payment',
    delivery_payment: 'Delivery Payment',
    transport_payment: 'Transport Payment',
}

const typeLabel = (t: LedgerEventType) => typeLabelMap[t] || t

const eventCount = computed(() => props.events.length)

const generatedDate = computed(() => {
    return new Date().toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    })
})

const periodLabel = computed(() => {
    if (props.dateFrom && props.dateTo) {
        return `${shortDate(props.dateFrom)} â†’ ${shortDate(props.dateTo)}`
    }
    if (props.dateFrom) return `From ${shortDate(props.dateFrom)}`
    if (props.dateTo) return `Up to ${shortDate(props.dateTo)}`
    return 'All time'
})

const PAGE_HEIGHT_MM = 297
const PAGE_VERTICAL_PADDING_MM = 22
const USABLE_MM = PAGE_HEIGHT_MM - PAGE_VERTICAL_PADDING_MM

const PAGE1_HEADER_MM = 110
const PAGE_OTHER_HEADER_MM = 22
const TIMELINE_HEADER_MM = 8
const PAGE_FOOTER_MM = 14

const EVENT_BASE_MM = 16
const EVENT_META_MM = 6

const estimateEventHeight = (event: LedgerEvent): number => {
    let h = EVENT_BASE_MM
    const metaCount = event.meta ? Object.keys(event.meta).length : 0
    if (metaCount > 0) {
        const lines = Math.ceil(metaCount / 3)
        h += EVENT_META_MM * lines
    }
    return h
}

const pages = computed<LedgerEvent[][]>(() => {
    const events = props.events

    if (events.length === 0) {
        return [[]]
    }

    const result: LedgerEvent[][] = []
    let currentPage: LedgerEvent[] = []
    let isFirstPage = true

    const availableFor = (first: boolean): number => {
        const header = first ? PAGE1_HEADER_MM : PAGE_OTHER_HEADER_MM
        return USABLE_MM - header - TIMELINE_HEADER_MM - PAGE_FOOTER_MM
    }

    let available = availableFor(true)

    for (const event of events) {
        const h = estimateEventHeight(event)

        if (h > available && currentPage.length > 0) {
            result.push(currentPage)
            currentPage = []
            isFirstPage = false
            available = availableFor(false)
        }

        currentPage.push(event)
        available -= h
    }

    if (currentPage.length > 0) {
        result.push(currentPage)
    }

    return result
})

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

const amountSign = (kind?: 'debit' | 'credit' | 'neutral'): string => {
    if (kind === 'debit') return '+'
    if (kind === 'credit') return 'âˆ’'
    return ''
}

const humanKey = (key: string): string => {
    return key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}
</script>

<style scoped>
.ledger-print-root {
    display: flex;
    flex-direction: column;
    gap: 8mm;
    background: transparent;
}

.ledger-page {
    width: 210mm;
    height: 297mm;
    background: white;
    color: #1f2937;
    display: flex;
    flex-direction: column;
    padding: 11mm 14mm;
    box-sizing: border-box;
    overflow: hidden;
    position: relative;
}

.page-header {
    flex-shrink: 0;
}

.page-content {
    flex: 1;
    min-height: 0;
    margin-top: 4mm;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.page-footer {
    flex-shrink: 0;
    margin-top: auto;
    background: #075985;
    color: white;
    padding: 8px 12px;
    font-size: 8.5px;
    line-height: 1.5;
    text-align: center;
    border-radius: 2px;
}

.footer-line-1 {
    font-weight: 500;
}

.footer-line-2 {
    opacity: 0.95;
}

.company-header {
    border-bottom: 2px solid #1f2937;
    padding-bottom: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.logo-area {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.logo-box {
    display: flex;
    height: 18mm;
    width: 22mm;
    align-items: center;
    justify-content: center;
}

.logo-box img {
    max-height: 100%;
    max-width: 100%;
    object-fit: contain;
}

.company-name {
    font-family: 'Cinzel', 'Times New Roman', serif;
    font-size: 22px;
    font-weight: 700;
    line-height: 1;
    letter-spacing: 0.04em;
    color: #172554;
}

.company-divider {
    margin-top: 4px;
    height: 2px;
    width: 80px;
    background: #0288d1;
}

.business-type {
    font-family: 'Playfair Display', Georgia, serif;
    font-size: 15px;
    font-weight: 700;
    font-style: italic;
    color: #b91c1c;
    text-align: right;
}

.subtitle {
    font-size: 9px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #6b7280;
    text-align: right;
    margin-top: 2px;
}

.meta-row {
    margin-top: 1.25rem;
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
}

.ref-box {
    display: flex;
    border: 1px solid #374151;
    font-size: 10px;
}

.ref-label {
    border-right: 1px solid #374151;
    background: #f3f4f6;
    padding: 4px 12px;
    font-weight: 600;
    text-transform: uppercase;
}

.ref-value {
    padding: 4px 16px;
    font-weight: 700;
    color: #1f2937;
    min-width: 40mm;
    text-align: center;
}

.title-center {
    text-align: center;
}

.title-center h2 {
    font-family: 'Cinzel', 'Times New Roman', serif;
    font-size: 18px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: #111827;
}

.title-underline {
    margin: 4px auto 0;
    height: 2px;
    width: 48px;
    background: #0288d1;
}

.date-box {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 10px;
}

.date-box .label {
    font-weight: 600;
    text-transform: uppercase;
}

.date-value {
    border: 1px solid #374151;
    background: white;
    padding: 4px 16px;
    font-weight: 500;
    min-width: 30mm;
    text-align: center;
}

.bill-grid {
    margin-top: 1.25rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

.bill-box {
    border: 1px solid #374151;
}

.bill-box .bill-label {
    border-bottom: 1px solid #374151;
    background: #f3f4f6;
    padding: 4px 12px;
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
}

.bill-box .bill-value {
    padding: 10px 12px;
    min-height: 15mm;
    font-size: 11px;
}

.customer-name {
    font-weight: 700;
    font-size: 12px;
    color: #111827;
}

.customer-sub {
    font-size: 10px;
    color: #4b5563;
    margin-top: 2px;
}

.summary-grid {
    margin-top: 1rem;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 0;
    border: 1px solid #1f2937;
}

.summary-cell {
    padding: 8px 12px;
    border-right: 1px solid #1f2937;
}

.summary-cell:last-child {
    border-right: none;
}

.summary-label {
    font-size: 8.5px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #6b7280;
    font-weight: 700;
}

.summary-value {
    font-size: 13px;
    font-weight: 700;
    color: #111827;
    margin-top: 2px;
}

.summary-value.paid {
    color: #15803d;
}

.summary-value.outstanding {
    color: #b91c1c;
}

.continuation-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #374151;
    padding-bottom: 6px;
}

.continuation-title {
    font-family: 'Cinzel', 'Times New Roman', serif;
    font-size: 13px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #111827;
}

.continuation-sub {
    font-size: 10px;
    color: #4b5563;
    margin-top: 1px;
}

.continuation-right {
    font-size: 10px;
    color: #6b7280;
    font-weight: 600;
}

.timeline-header {
    display: grid;
    grid-template-columns: 32mm 1fr;
    background: #f3f4f6;
    border: 1px solid #374151;
    border-bottom: none;
    font-size: 8.5px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #374151;
}

.timeline-header .th-date {
    padding: 4px 10px;
    border-right: 1px solid #374151;
}

.timeline-header .th-event {
    padding: 4px 10px;
    display: flex;
    justify-content: space-between;
    gap: 8px;
}

.timeline-header .th-type {
    flex: 0 0 auto;
}

.timeline-header .th-desc {
    flex: 1;
    margin-left: 6px;
}

.timeline-header .th-amount {
    flex: 0 0 auto;
}

.empty-state {
    padding: 20px 16px;
    font-size: 10px;
    color: #6b7280;
    text-align: center;
    border: 1px solid #374151;
    border-top: none;
}

.timeline {
    display: flex;
    flex-direction: column;
    border: 1px solid #374151;
    border-top: none;
}

.timeline-row {
    display: grid;
    grid-template-columns: 32mm 1fr;
    border-bottom: 1px solid #e5e7eb;
    page-break-inside: avoid;
    break-inside: avoid;
}

.timeline-row:last-child {
    border-bottom: none;
}

.timeline-date {
    padding: 6px 10px;
    font-size: 9.5px;
    font-weight: 600;
    color: #374151;
    background: #f9fafb;
    border-right: 1px solid #e5e7eb;
}

.timeline-body {
    padding: 6px 10px;
}

.row-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
}

.row-head-left {
    display: flex;
    align-items: center;
    gap: 6px;
    min-width: 0;
}

.type-chip {
    font-size: 8px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    padding: 2px 6px;
    background: #e5e7eb;
    color: #374151;
    border-radius: 3px;
}

.ref-inline {
    font-family: 'Courier New', monospace;
    font-size: 9px;
    color: #6b7280;
}

.row-amount {
    font-size: 10.5px;
    font-weight: 700;
    white-space: nowrap;
}

.row-amount.debit {
    color: #1e40af;
}

.row-amount.credit {
    color: #15803d;
}

.row-amount.neutral {
    color: #6b7280;
}

.amount-label {
    font-size: 8px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    margin-left: 4px;
    color: #9ca3af;
    font-weight: 600;
}

.row-title {
    font-size: 10.5px;
    font-weight: 600;
    color: #111827;
    margin-top: 3px;
}

.row-desc {
    font-size: 9.5px;
    color: #4b5563;
    margin-top: 1px;
}

.row-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 3px 12px;
    margin-top: 4px;
    padding-top: 4px;
    border-top: 1px dashed #e5e7eb;
}

.meta-item {
    font-size: 9px;
    color: #374151;
}

.meta-key {
    color: #6b7280;
    font-weight: 600;
    margin-right: 3px;
}

.meta-val {
    color: #111827;
}

@media print {
    .ledger-print-root {
        gap: 0;
    }

    .ledger-page {
        box-shadow: none;
        margin: 0;
        page-break-after: always;
    }

    .ledger-page:last-child {
        page-break-after: auto;
    }
}
</style>