<!-- src/components/features/invoices/InvoicePreview.vue -->
<template>
    <div class="space-y-6">
        <!-- Invoice Details -->
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
                    <p class="text-sm text-(--color-text-secondary)">Review and download your invoice</p>
                </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Invoice Number</label>
                    <input :value="localInvoice?.number || ''"
                        @input="updateField('number', ($event.target as HTMLInputElement).value)" type="text"
                        class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Invoice Date</label>
                    <input :value="localInvoice?.date || ''"
                        @input="updateField('date', ($event.target as HTMLInputElement).value)" type="date"
                        class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Bill To</label>
                    <input :value="localInvoice?.customer_name || ''"
                        @input="updateField('customer_name', ($event.target as HTMLInputElement).value)" type="text"
                        class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200" />
                </div>
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Notes</label>
                    <textarea :value="localInvoice?.notes || ''"
                        @input="updateField('notes', ($event.target as HTMLTextAreaElement).value)" rows="2"
                        placeholder="Additional notes..."
                        class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200 resize-none"></textarea>
                </div>
            </div>

            <!-- Financial Summary -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-4 pt-4 border-t border-(--color-border)">
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30">
                    <p class="text-xs text-(--color-text-secondary)">Total Items</p>
                    <p class="text-lg font-bold text-(--color-text-primary)">{{ localInvoice?.items.length || 0 }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30">
                    <p class="text-xs text-(--color-text-secondary)">Total Outstanding</p>
                    <p class="text-lg font-bold text-(--color-red)">৳{{ formatAmount(localInvoice?.total || 0) }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30">
                    <p class="text-xs text-(--color-text-secondary)">Already Paid</p>
                    <p class="text-lg font-bold text-(--color-green)">৳{{ formatAmount(localInvoice?.received || 0)
                    }}</p>
                </div>
                <div class="p-3 rounded-lg bg-(--color-muted-bg)/30">
                    <p class="text-xs text-(--color-text-secondary)">Total Invoice</p>
                    <p class="text-lg font-bold text-(--color-blue)">৳{{ formatAmount(localInvoice?.total || 0) }}
                    </p>
                </div>
            </div>
        </div>

        <!-- Invoice Preview -->
        <div class="rounded-xl border border-(--color-border) bg-(--color-surface) overflow-hidden">
            <div
                class="flex items-center justify-between p-4 border-b border-(--color-border) bg-(--color-muted-bg)/10">
                <div>
                    <h2 class="text-base font-semibold text-(--color-text-primary)">Invoice Preview</h2>
                    <p class="text-sm text-(--color-text-secondary)">A4 · Ready to download</p>
                </div>
                <span
                    class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium border border-(--color-green) text-(--color-green)">
                    <span class="w-1.5 h-1.5 rounded-full bg-(--color-green)"></span>
                    Ready
                </span>
            </div>

            <div class="p-4 bg-(--color-muted-bg)/10 overflow-auto">
                <div ref="invoicePrintRef" class="invoice-page">
                    <InvoicePrintView :invoice="localInvoice" />
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 border-t border-(--color-border)">
            <button @click="emit('back')"
                class="w-full sm:w-auto px-5 py-2.5 text-sm font-medium rounded-xl border border-(--color-border) hover:bg-(--color-muted-bg) transition-all duration-200 flex items-center justify-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                Back to Items
            </button>
            <div class="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
                <button @click="handleDownloadPDF"
                    class="w-full sm:w-auto px-6 py-2.5 bg-(--color-blue) text-white rounded-xl text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 flex items-center justify-center gap-2">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    Download PDF
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import type { Invoice, InvoiceType } from '@/types/invoice'
import InvoicePrintView from './InvoicePrintView.vue'
import { push } from 'notivue'
import jsPDF from 'jspdf'
import html2canvas from 'html2canvas'

const props = defineProps<{
    type: InvoiceType
    invoice: Invoice | null
}>()

const emit = defineEmits<{
    'update:invoice': [value: Partial<Invoice>]
    'back': []
    'download': []
}>()

const localInvoice = ref<Invoice | null>(null)
const invoicePrintRef = ref<HTMLElement | null>(null)
const isDownloading = ref(false)

const formatAmount = (value: number): string => {
    return value.toFixed(2)
}

const updateField = <K extends keyof Invoice>(field: K, value: Invoice[K]) => {
    if (localInvoice.value) {
        localInvoice.value[field] = value
        emit('update:invoice', { [field]: value })
    }
}

const handleDownloadPDF = async () => {
    if (isDownloading.value) return

    isDownloading.value = true
    push.info('Generating PDF...')

    try {
        await new Promise(resolve => setTimeout(resolve, 200))

        const element = invoicePrintRef.value
        if (!element) {
            push.error('Failed to generate PDF')
            isDownloading.value = false
            return
        }

        const canvas = await html2canvas(element, {
            scale: 2,
            useCORS: true,
            logging: false,
            backgroundColor: '#ffffff',
            width: 794,
            height: 1123,
        })

        const imgData = canvas.toDataURL('image/png')
        const pdf = new jsPDF({
            orientation: 'portrait',
            unit: 'px',
            format: [794, 1123],
        })

        pdf.addImage(imgData, 'PNG', 0, 0, 794, 1123)

        const filename = `invoice-${localInvoice.value?.number || 'draft'}-${new Date().toISOString().slice(0, 10)}.pdf`
        pdf.save(filename)

        push.success('PDF downloaded successfully!')
        emit('download')
    } catch (error) {
        console.error('PDF generation error:', error)
        push.error('Failed to generate PDF')
    } finally {
        isDownloading.value = false
    }
}

const handleGlobalDownload = () => {
    handleDownloadPDF()
}

watch(() => props.invoice, (newInvoice) => {
    localInvoice.value = newInvoice ? { ...newInvoice } : null
}, { immediate: true })

onMounted(() => {
    window.addEventListener('download-invoice', handleGlobalDownload)
})

onUnmounted(() => {
    window.removeEventListener('download-invoice', handleGlobalDownload)
})
</script>

<style scoped>
.invoice-page {
    width: 210mm;
    min-height: 297mm;
    max-width: 210mm;
    background: white;
    color: #1f2937;
    margin: 0 auto;
    overflow: hidden;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}
</style>