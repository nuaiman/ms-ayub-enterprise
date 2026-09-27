<!-- src/components/features/invoices/InvoiceGenerator.vue -->
<template>
    <div class="invoice-generator">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
            <div class="flex items-center gap-4">
                <button @click="goBack"
                    class="w-10 h-10 flex items-center justify-center border border-(--color-border) rounded-xl hover:bg-(--color-muted-bg) transition-all duration-200 group">
                    <svg class="w-5 h-5 text-(--color-text-secondary) group-hover:text-(--color-text-primary) transition-colors"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                    </svg>
                </button>
                <div>
                    <h1 class="text-xl font-bold text-(--color-text-primary)">Create Invoice</h1>
                    <p class="text-sm text-(--color-text-secondary) mt-0.5">Select party, bills, and generate</p>
                </div>
            </div>

            <div class="flex items-center gap-2">
                <button @click="saveDraft"
                    class="h-10 px-4 flex items-center gap-2 border border-(--color-border) rounded-xl text-sm font-medium text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-all duration-200">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                    </svg>
                    <span class="hidden sm:inline">Save Draft</span>
                    <span class="sm:hidden">Draft</span>
                </button>
                <button @click="downloadPDF"
                    class="h-10 px-5 flex items-center gap-2 bg-(--color-blue) text-white rounded-xl text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span class="hidden sm:inline">Download PDF</span>
                    <span class="sm:hidden">PDF</span>
                </button>
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading" class="flex items-center justify-center py-16">
            <div class="text-center space-y-4">
                <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                <p class="text-sm text-(--color-text-secondary)">Loading invoice data...</p>
            </div>
        </div>

        <!-- Steps -->
        <div v-else>
            <div class="flex items-center gap-3 mb-8">
                <div v-for="(stepInfo, index) in steps" :key="index" class="flex items-center gap-3"
                    :class="index < steps.length - 1 ? 'flex-1' : ''">
                    <div class="flex items-center gap-2">
                        <span
                            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold transition-all duration-300"
                            :class="step >= index + 1
                                ? 'bg-(--color-blue) text-white shadow-lg shadow-(--color-blue)/20'
                                : 'bg-(--color-muted-bg) text-(--color-text-secondary)'">
                            {{ index + 1 }}
                        </span>
                        <span class="text-sm font-medium transition-colors duration-300 hidden sm:block"
                            :class="step >= index + 1 ? 'text-(--color-text-primary)' : 'text-(--color-text-secondary)'">
                            {{ stepInfo.label }}
                        </span>
                    </div>
                    <div v-if="index < steps.length - 1"
                        class="flex-1 h-0.5 bg-(--color-border) transition-colors duration-300"
                        :class="step > index + 1 ? 'bg-(--color-blue)' : ''"></div>
                </div>
            </div>

            <div class="bg-(--color-surface) rounded-2xl border border-(--color-border) p-6 shadow-sm">
                <!-- Step 1: Party Selection -->
                <div v-if="step === 1" class="space-y-6">
                    <InvoicePartySelect v-model:party="party" v-model:party-id="partyId" @next="goToStep2" />
                </div>

                <!-- Step 2: Item Selection -->
                <div v-if="step === 2" class="space-y-6">
                    <InvoiceItemSelector :party="party" :party-id="partyId" :party-name="partyName"
                        v-model:selected-items="selectedItems" @next="goToStep3" @back="goToStep1" />
                </div>

                <!-- Step 3: Preview -->
                <div v-if="step === 3" class="space-y-6">
                    <InvoicePreview :invoice="currentInvoice" @update:invoice="updateInvoice" @back="goToStep2" />
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useInvoiceStore } from '@/stores/invoices'
import { useCustomersStore } from '@/stores/customers'
import { useBrokersStore } from '@/stores/brokers'
import { useGodownsStore } from '@/stores/godowns'
import { useMajhisStore } from '@/stores/majhis'
import { useCustomerStorageBillsStore } from '@/stores/customerStorageBills'
import { useCustomerLotBillsStore } from '@/stores/customerLotBills'
import { useCustomerDeliveryBillsStore } from '@/stores/customerDeliveryBills'
import { useCustomerTransportBillsStore } from '@/stores/customerTransportBills'
import { useCustomerAdditionalBillsStore } from '@/stores/customerAdditionalBills'
import { useBrokerVehicleBillsStore } from '@/stores/brokerVehicleBills'
import { useGodownStoreBillsStore } from '@/stores/godownStoreBills'
import { useMajhiLotBillsStore } from '@/stores/majhiLotBills'
import { useMajhiLoadingBillsStore } from '@/stores/majhiLoadingBills'
import { useLotsStore } from '@/stores/lots'
import { useStoresStore } from '@/stores/stores'
import { useDeliveriesStore } from '@/stores/deliveries'
import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useTransportsStore } from '@/stores/transports'
import { useVehiclesStore } from '@/stores/vehicles'
import { useDamagesStore } from '@/stores/damages'
import type { Invoice, AvailableItem, InvoiceParty } from '@/types/invoice'
import { push } from 'notivue'
import InvoicePartySelect from '@/components/features/invoices/InvoicePartySelect.vue'
import InvoiceItemSelector from '@/components/features/invoices/InvoiceItemSelector.vue'
import InvoicePreview from '@/components/features/invoices/InvoicePreview.vue'

const router = useRouter()
const invoiceStore = useInvoiceStore()
const customersStore = useCustomersStore()
const brokersStore = useBrokersStore()
const godownsStore = useGodownsStore()
const majhisStore = useMajhisStore()
const lotsStore = useLotsStore()
const storesStore = useStoresStore()
const deliveriesStore = useDeliveriesStore()
const deliveryItemsStore = useDeliveryItemsStore()
const transportsStore = useTransportsStore()
const vehiclesStore = useVehiclesStore()
const damagesStore = useDamagesStore()

const step = ref(1)
const party = ref<InvoiceParty>('customer')
const partyId = ref<number | null>(null)
const selectedItems = ref<AvailableItem[]>([])
const isLoading = ref(true)

const steps = [
    { label: 'Select Party' },
    { label: 'Choose Bills' },
    { label: 'Review & Download' },
]

const currentInvoice = computed(() => invoiceStore.currentInvoice)

const partyName = computed(() => {
    if (!partyId.value) return ''
    switch (party.value) {
        case 'customer': return customersStore.getCustomerName(partyId.value)
        case 'broker': return brokersStore.getBrokerName(partyId.value)
        case 'godown': return godownsStore.getGodownName(partyId.value)
        case 'majhi': return majhisStore.getMajhiName(partyId.value)
        default: return ''
    }
})

const loadInvoiceData = async () => {
    isLoading.value = true
    try {
        await Promise.all([
            customersStore.fetchCustomers(),
            brokersStore.fetchBrokers(),
            godownsStore.fetchGodowns(),
            majhisStore.fetchMajhis(),
            lotsStore.fetchLots(),
            storesStore.fetchStores(),
            deliveriesStore.fetchDeliveries(),
            deliveryItemsStore.fetchDeliveryItems(),
            transportsStore.fetchTransports(),
            vehiclesStore.fetchVehicles(),
            damagesStore.fetchDamages(),
        ])
    } catch (error) {
        console.error('[INVOICE] Error loading data:', error)
        push.error('Failed to load invoice data')
    } finally {
        isLoading.value = false
    }
}

onMounted(async () => {
    await loadInvoiceData()
})

const goToStep2 = () => {
    if (!partyId.value) {
        push.error('Please select a party')
        return
    }
    step.value = 2
}

const goToStep3 = () => {
    const selected = selectedItems.value.filter(item => item.selected)
    if (selected.length === 0) {
        push.error('Please select at least one item')
        return
    }

    if (!partyId.value) {
        push.error('Party not selected')
        return
    }

    let invoice = currentInvoice.value
    if (!invoice || invoice.party_id !== partyId.value || invoice.party_type !== party.value) {
        invoice = invoiceStore.createInvoice(party.value, partyId.value, partyName.value)
    }

    invoiceStore.buildInvoiceFromSelectedItems(selected)
    step.value = 3
}

const goToStep1 = () => { step.value = 1 }
const goToStep2Back = () => { step.value = 2 }

const updateInvoice = (updates: Partial<Invoice>) => {
    if (currentInvoice.value) {
        Object.keys(updates).forEach((key) => {
            const field = key as keyof Invoice
            const value = updates[field]
            if (value !== undefined) {
                invoiceStore.updateInvoiceField(field, value)
            }
        })
    }
}

const saveDraft = () => {
    if (currentInvoice.value && currentInvoice.value.items.length > 0) {
        invoiceStore.saveCurrentInvoice()
    } else {
        push.warning('No items to save')
    }
}

const downloadPDF = () => {
    const downloadEvent = new CustomEvent('download-invoice')
    window.dispatchEvent(downloadEvent)
}

const goBack = () => {
    router.back()
}
</script>