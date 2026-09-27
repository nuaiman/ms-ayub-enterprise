<!-- src/components/features/invoices/InvoicePrintView.vue -->
<template>
    <div class="invoice-print-root">
        <div v-for="(page, pageIdx) in pages" :key="pageIdx" class="invoice-page-inner"
            :ref="el => setPageRef(el, pageIdx)">
            <div class="invoice-content">
                <!-- Company Header (page 1 only) -->
                <header v-if="pageIdx === 0" class="company-header">
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
                </header>

                <!-- Meta (page 1 only) -->
                <div v-if="pageIdx === 0" class="meta-row">
                    <div class="ref-box">
                        <span class="ref-label">Ref. Invoice #</span>
                        <span class="ref-value">{{ invoice?.number || '—' }}</span>
                    </div>
                    <div class="title-center">
                        <h2>{{ title }}</h2>
                        <div class="title-underline"></div>
                    </div>
                    <div class="date-box">
                        <span class="label">Date</span>
                        <span class="date-value">{{ invoice?.date || '—' }}</span>
                    </div>
                </div>

                <!-- Bill to grid (page 1 only) -->
                <div v-if="pageIdx === 0" class="bill-grid">
                    <div class="bill-box">
                        <div class="bill-label">Bill To</div>
                        <div class="bill-value">{{ invoice?.party_name || '—' }}</div>
                    </div>
                    <div class="bill-box">
                        <div class="bill-label">Party Type</div>
                        <div class="bill-value capitalize">{{ partyLabel }}</div>
                    </div>
                </div>

                <!-- Continuation header (subsequent pages) -->
                <div v-else class="continuation-header">
                    <div class="continuation-left">
                        <div class="continuation-title">{{ title }}</div>
                        <div class="continuation-sub">
                            {{ invoice?.party_name || '—' }}
                            <span v-if="invoice?.number" class="continuation-ref"> · Ref: {{ invoice.number }}</span>
                        </div>
                    </div>
                    <div class="continuation-right">
                        Page {{ pageIdx + 1 }} of {{ pages.length }}
                    </div>
                </div>

                <!-- Items Table -->
                <div class="table-wrap">
                    <table class="invoice-table">
                        <thead>
                            <tr>
                                <th style="width:18%;">Item</th>
                                <th style="width:22%;">Date</th>
                                <th style="width:30%;">Description</th>
                                <th style="width:10%;">Qty</th>
                                <th style="width:10%;">Rate</th>
                                <th style="width:10%;">Amount</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(item, idx) in page" :key="item.id" :class="{ 'border-t': idx > 0 }">
                                <td>{{ item.item || '—' }}</td>
                                <td>{{ item.date || '—' }}</td>
                                <td>{{ item.description || '—' }}</td>
                                <td style="text-align:right;">{{ item.quantity || 0 }}</td>
                                <td style="text-align:right;">{{ formatMoney(item.rate) }}</td>
                                <td style="text-align:right; font-weight:600;">{{ formatMoney(item.amount) }}</td>
                            </tr>
                            <tr v-if="page.length === 0">
                                <td colspan="6" style="text-align:center; color:#6b7280; padding:12px;">
                                    No items on this page.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- Totals + signature (every page) -->
                <div class="bottom-grid">
                    <div class="signature-area">
                        <div style="height:12mm;"></div>
                        <div class="signature-line">Receiver's Signature</div>
                    </div>
                    <div>
                        <div class="totals-box">
                            <div class="totals-row">
                                <span class="label">Total Outstanding</span>
                                <span class="value">BDT. {{ formatMoney(invoice?.total || 0) }}</span>
                            </div>
                            <div class="totals-row bg-light">
                                <span class="label">Already Paid</span>
                                <span class="value" style="color: #15803d;">BDT. {{ formatMoney(invoice?.received || 0)
                                    }}</span>
                            </div>
                            <div class="totals-row">
                                <span class="label">Total Invoice</span>
                                <span class="value">BDT. {{ formatMoney(invoice?.total || 0) }}</span>
                            </div>
                        </div>
                        <div class="company-signature">
                            <div class="signature-for">For Ayub Enterprise</div>
                            <div class="signature-line-bottom">Authorized Signature</div>
                        </div>
                    </div>
                </div>

                <!-- Notes (page 1 only) -->
                <div v-if="pageIdx === 0 && invoice?.notes" class="notes-print">
                    <div class="label">Notes</div>
                    <div class="text">{{ invoice.notes }}</div>
                </div>
            </div>

            <!-- Footer (every page) -->
            <footer class="invoice-footer">
                <p style="font-weight:500;">Ali Hossain Chairman Building (1st Floor), 958/27, Strand Road, Mazirghat,
                    Chattogram.</p>
                <p style="opacity:0.95;">Cell : 01813-397288, 01705-727492, 01793-287709, 01864-106406 · E-mail :
                    mdayubenterprise@gmail.com</p>
            </footer>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, type ComponentPublicInstance } from 'vue'
import type { Invoice, InvoiceItem } from '@/types/invoice'
import { invoiceTitle, partyTypeLabel } from '@/types/invoice'

const props = defineProps<{
    invoice: Invoice | null
}>()

const title = computed(() => props.invoice ? invoiceTitle(props.invoice.party_type) : 'Invoice')
const partyLabel = computed(() => props.invoice ? partyTypeLabel(props.invoice.party_type) : '—')

const pageRefs = ref<HTMLElement[]>([])
const setPageRef = (el: Element | ComponentPublicInstance | null, idx: number) => {
    if (el instanceof HTMLElement) {
        pageRefs.value[idx] = el
    }
}
defineExpose({ pageRefs })

const PAGE_HEIGHT_MM = 297
const PAGE_VERTICAL_PADDING_MM = 18
const USABLE_MM = PAGE_HEIGHT_MM - PAGE_VERTICAL_PADDING_MM

const PAGE1_HEADER_MM = 90
const PAGE_OTHER_HEADER_MM = 14
const TABLE_HEADER_MM = 10
const BOTTOM_BLOCK_MM = 60
const NOTES_BLOCK_MM = 20
const FOOTER_MM = 16

const ROW_BASE_MM = 8
const ROW_DESC_EXTRA_MM = 4

const estimateRowHeight = (item: InvoiceItem): number => {
    let h = ROW_BASE_MM
    const descLen = (item.description || '').length
    if (descLen > 40) h += ROW_DESC_EXTRA_MM
    if (descLen > 90) h += ROW_DESC_EXTRA_MM
    return h
}

const pages = computed<InvoiceItem[][]>(() => {
    const items = props.invoice?.items || []
    if (items.length === 0) return [[]]

    const result: InvoiceItem[][] = []
    let currentPage: InvoiceItem[] = []
    let isFirst = true

    const availableFor = (first: boolean): number => {
        const header = first ? PAGE1_HEADER_MM : PAGE_OTHER_HEADER_MM
        const notes = first && props.invoice?.notes ? NOTES_BLOCK_MM : 0
        return USABLE_MM - header - TABLE_HEADER_MM - BOTTOM_BLOCK_MM - FOOTER_MM - notes
    }

    let available = availableFor(true)

    for (const item of items) {
        const h = estimateRowHeight(item)
        if (h > available && currentPage.length > 0) {
            result.push(currentPage)
            currentPage = []
            isFirst = false
            available = availableFor(false)
        }
        currentPage.push(item)
        available -= h
    }

    if (currentPage.length > 0) result.push(currentPage)
    return result
})

const formatMoney = (val: number): string => {
    return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    }).format(val || 0)
}
</script>

<style scoped>
.invoice-print-root {
    display: flex;
    flex-direction: column;
    gap: 8mm;
    background: transparent;
}

.invoice-page-inner {
    width: 210mm;
    height: 297mm;
    background: white;
    color: #1f2937;
    display: flex;
    flex-direction: column;
    padding: 11mm 14mm 7mm 14mm;
    box-sizing: border-box;
    overflow: hidden;
}

.invoice-content {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-height: 0;
}

.company-header {
    border-bottom: 2px solid #1f2937;
    padding-bottom: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-shrink: 0;
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
    flex-shrink: 0;
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
    min-width: 18mm;
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
    flex-shrink: 0;
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
    padding: 12px 12px;
    min-height: 15mm;
    font-size: 11px;
    font-weight: 600;
    display: flex;
    align-items: center;
}

.continuation-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #374151;
    padding-bottom: 6px;
    flex-shrink: 0;
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

.continuation-ref {
    color: #6b7280;
}

.continuation-right {
    font-size: 10px;
    color: #6b7280;
    font-weight: 600;
}

.table-wrap {
    margin-top: 1rem;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    border: 1px solid #1f2937;
}

.invoice-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 10px;
}

.invoice-table th {
    border-bottom: 2px solid #1f2937;
    background: #f3f4f6;
    padding: 8px 4px;
    text-align: center;
    font-weight: 700;
    text-transform: uppercase;
}

.invoice-table td {
    padding: 8px 4px;
    vertical-align: top;
    border-right: 1px solid #1f2937;
    text-align: center;
    word-break: break-word;
}

.invoice-table td:last-child {
    border-right: none;
}

.invoice-table .border-t {
    border-top: 1px solid #d1d5db;
}

/* Closes the table under the last row */
.invoice-table tbody tr:last-child td {
    border-bottom: 1px solid #1f2937;
}

.bottom-grid {
    margin-top: 1rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    align-items: end;
    flex-shrink: 0;
    page-break-inside: avoid;
    break-inside: avoid;
}

.signature-area {
    padding-bottom: 0.5rem;
}

.signature-line {
    margin-top: 12mm;
    border-top: 1px solid #1f2937;
    padding-top: 4px;
    text-align: center;
    font-size: 9px;
    font-weight: 500;
    width: 48mm;
}

.totals-box {
    border: 1px solid #1f2937;
    font-size: 10px;
}

.totals-row {
    display: flex;
    justify-content: space-between;
    padding: 8px 12px;
    border-bottom: 1px solid #374151;
}

.totals-row:last-child {
    border-bottom: none;
}

.totals-row.bg-light {
    background: #f9fafb;
}

.totals-row .label {
    font-weight: 600;
}

.totals-row .value {
    font-weight: 700;
}

.company-signature {
    margin-left: auto;
    margin-top: 1.5rem;
    width: 55mm;
    text-align: center;
}

.signature-for {
    font-size: 11px;
    font-weight: 700;
    color: #1f2937;
}

.signature-line-bottom {
    margin-top: 8mm;
    border-top: 1px solid #1f2937;
    padding-top: 4px;
    font-size: 9px;
    font-weight: 500;
    color: #4b5563;
}

.notes-print {
    margin-top: 0.75rem;
    border-top: 1px solid #e5e7eb;
    padding-top: 0.5rem;
    flex-shrink: 0;
}

.notes-print .label {
    font-size: 8px;
    font-weight: 700;
    text-transform: uppercase;
    color: #9ca3af;
}

.notes-print .text {
    font-size: 9px;
    line-height: 1.5;
    color: #4b5563;
}

.invoice-footer {
    background: #075985;
    color: white;
    padding: 10px 14mm;
    font-size: 8.5px;
    line-height: 1.6;
    text-align: center;
    flex-shrink: 0;
    margin: 0 -14mm -7mm -14mm;
}

@media print {
    .invoice-print-root {
        gap: 0;
    }

    .invoice-page-inner {
        box-shadow: none;
        margin: 0;
        page-break-after: always;
    }

    .invoice-page-inner:last-child {
        page-break-after: auto;
    }
}
</style>