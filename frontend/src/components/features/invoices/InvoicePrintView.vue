<!-- src/components/features/invoices/InvoicePrintView.vue -->
<template>
    <div class="invoice-page-inner">
        <div class="invoice-content">
            <!-- Company Header -->
            <header class="company-header">
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

            <!-- Meta -->
            <div class="meta-row">
                <div class="ref-box">
                    <span class="ref-label">Ref. Invoice #</span>
                    <span class="ref-value">{{ invoice?.number || '—' }}</span>
                </div>
                <div class="title-center">
                    <h2>{{ invoice?.type === 'godown' ? 'Godown Bill' : 'Transport Bill' }}</h2>
                    <div class="title-underline"></div>
                </div>
                <div class="date-box">
                    <span class="label">Date</span>
                    <span class="date-value">{{ invoice?.date || '—' }}</span>
                </div>
            </div>

            <!-- Bill To -->
            <div class="bill-grid">
                <div class="bill-box">
                    <div class="bill-label">Bill To</div>
                    <div class="bill-value">{{ invoice?.customer_name || '—' }}</div>
                </div>
                <div class="bill-box">
                    <div class="bill-label">Invoice Type</div>
                    <div class="bill-value capitalize">{{ invoice?.type || '—' }}</div>
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
                        <tr v-for="(item, idx) in invoice?.items || []" :key="item.id" :class="{ 'border-t': idx > 0 }">
                            <td>{{ item.item || '—' }}</td>
                            <td>{{ item.date || '—' }}</td>
                            <td>{{ item.description || '—' }}</td>
                            <td style="text-align:right;">{{ item.quantity || 0 }}</td>
                            <td style="text-align:right;">{{ formatMoney(item.rate) }}</td>
                            <td style="text-align:right; font-weight:600;">{{ formatMoney(item.amount) }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Bottom -->
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

            <!-- Notes -->
            <div v-if="invoice?.notes" class="notes-print">
                <div class="label">Notes</div>
                <div class="text">{{ invoice.notes }}</div>
            </div>
        </div>

        <!-- Footer -->
        <footer class="invoice-footer">
            <p style="font-weight:500;">Ali Hossain Chairman Building (1st Floor), 958/27, Strand Road, Mazirghat,
                Chattogram.</p>
            <p style="opacity:0.95;">Cell : 01813-397288, 01705-727492, 01793-287709, 01864-106406 • E-mail :
                mdayubenterprise@gmail.com</p>
        </footer>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Invoice } from '@/types/invoice'

const props = defineProps<{
    invoice: Invoice | null
}>()

const formatMoney = (val: number): string => {
    return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
    }).format(val || 0)
}
</script>

<style scoped>
.invoice-page-inner {
    width: 210mm;
    min-height: 297mm;
    background: white;
    color: #1f2937;
    display: flex;
    flex-direction: column;
}

.invoice-content {
    display: flex;
    flex: 1;
    flex-direction: column;
    padding: 11mm 14mm 7mm 14mm;
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

.table-wrap {
    margin-top: 1.25rem;
    flex: 1;
    display: flex;
    flex-direction: column;
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
}

.invoice-table td:last-child {
    border-right: none;
}

.invoice-table .border-t {
    border-top: 1px solid #d1d5db;
}

.invoice-table tbody tr:last-child td {
    border-bottom: 1px solid #1f2937;
}

.bottom-grid {
    margin-top: 1.25rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
    align-items: end;
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
    margin-top: 2rem;
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
    margin-top: 1rem;
    border-top: 1px solid #e5e7eb;
    padding-top: 0.5rem;
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
}

@media print {
    .invoice-page-inner {
        width: 210mm;
        min-height: 297mm;
        box-shadow: none;
        margin: 0;
    }
}
</style>