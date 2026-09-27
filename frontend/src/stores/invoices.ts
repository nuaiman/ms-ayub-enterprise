// src/stores/invoices.ts

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
    Invoice,
    InvoiceParty,
    AvailableItem,
} from '@/types/invoice'
import { generateInvoiceNumber } from '@/types/invoice'
import { useCustomerStorageBillsStore } from './customerStorageBills'
import { useCustomerLotBillsStore } from './customerLotBills'
import { useCustomerDeliveryBillsStore } from './customerDeliveryBills'
import { useCustomerTransportBillsStore } from './customerTransportBills'
import { useCustomerAdditionalBillsStore } from './customerAdditionalBills'
import { useBrokerVehicleBillsStore } from './brokerVehicleBills'
import { useGodownStoreBillsStore } from './godownStoreBills'
import { useMajhiLotBillsStore } from './majhiLotBills'
import { useMajhiLoadingBillsStore } from './majhiLoadingBills'
import { push } from 'notivue'

const STORAGE_KEY = 'invoices_data_v2'

export const useInvoiceStore = defineStore('invoice', () => {
    const invoices = ref<Invoice[]>([])
    const currentInvoice = ref<Invoice | null>(null)

    const loadFromLocalStorage = () => {
        try {
            const data = localStorage.getItem(STORAGE_KEY)
            if (data) {
                invoices.value = JSON.parse(data)
            }
        } catch (error) {
            console.error('Failed to load invoices:', error)
        }
    }

    const saveToLocalStorage = () => {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(invoices.value))
        } catch (error) {
            console.error('Failed to save invoices:', error)
        }
    }

    // =========================================================================
    // AVAILABLE ITEMS PER PARTY
    // =========================================================================

    const itemsForCustomer = (customerId: number): AvailableItem[] => {
        const items: AvailableItem[] = []

        const storageStore = useCustomerStorageBillsStore()
        for (const bill of storageStore.customerBillData) {
            if (bill.customer_id !== customerId || bill.outstanding <= 0) continue
            const status = bill.total_paid > 0 ? 'partial' : 'unpaid'
            items.push({
                id: bill.lot_id,
                source_type: 'storage_bill',
                item: 'Storage Bill',
                date: bill.billing_start ? new Date(bill.billing_start).toLocaleDateString() : '',
                description: `${bill.item_name || `Lot #${bill.lot_id}`} - ${bill.months_billed} month(s)`,
                quantity: 1,
                rate: bill.monthly_bill || 0,
                amount: bill.outstanding,
                total_amount: bill.total_billed,
                paid_amount: bill.total_paid,
                status,
                selected: false,
            })
        }

        const lotStore = useCustomerLotBillsStore()
        for (const bill of lotStore.bills) {
            if (bill.customer_id !== customerId) continue
            const outstanding = bill.bill_amount - bill.paid_amount
            if (outstanding <= 0) continue
            const status = bill.paid_amount > 0 ? 'partial' : 'unpaid'
            items.push({
                id: bill.lot_id,
                source_type: 'lot_bill',
                item: 'Lot Unload Bill',
                date: bill.created_at ? new Date(bill.created_at).toLocaleDateString() : '',
                description: `${bill.item_name || `Lot #${bill.lot_id}`}`,
                quantity: 1,
                rate: bill.bill_amount || 0,
                amount: outstanding,
                total_amount: bill.bill_amount,
                paid_amount: bill.paid_amount,
                status,
                selected: false,
            })
        }

        const deliveryStore = useCustomerDeliveryBillsStore()
        for (const bill of deliveryStore.bills) {
            if (bill.customer_id !== customerId) continue
            const outstanding = bill.bill_amount - bill.paid_amount
            if (outstanding <= 0) continue
            const status = bill.paid_amount > 0 ? 'partial' : 'unpaid'
            items.push({
                id: bill.delivery_item_id,
                source_type: 'delivery_bill',
                item: 'Delivery Bill',
                date: bill.delivery_date ? new Date(bill.delivery_date).toLocaleDateString() : '',
                description: `${bill.item_name || `Delivery #${bill.delivery_id}`}`,
                quantity: 1,
                rate: bill.bill_amount || 0,
                amount: outstanding,
                total_amount: bill.bill_amount,
                paid_amount: bill.paid_amount,
                status,
                selected: false,
            })
        }

        const transportStore = useCustomerTransportBillsStore()
        for (const bill of transportStore.bills) {
            if (bill.customer_id !== customerId) continue
            const outstanding = bill.bill_amount - bill.paid_amount
            if (outstanding <= 0) continue
            const status = bill.paid_amount > 0 ? 'partial' : 'unpaid'
            items.push({
                id: bill.id,
                source_type: 'transport_bill',
                item: 'Transport Bill',
                date: bill.created_at ? new Date(bill.created_at).toLocaleDateString() : '',
                description: `${bill.from_location} → ${bill.to_location || 'N/A'} (${bill.total_vehicles} vehicles)`,
                quantity: bill.total_vehicles || 1,
                rate: bill.customer_charge_per_unit,
                amount: outstanding,
                total_amount: bill.bill_amount,
                paid_amount: bill.paid_amount,
                status,
                selected: false,
            })
        }

        const additionalStore = useCustomerAdditionalBillsStore()
        for (const bill of additionalStore.bills) {
            if (bill.customer_id !== customerId || bill.outstanding <= 0) continue
            const status = bill.paid_amount > 0 ? 'partial' : 'unpaid'
            items.push({
                id: bill.id,
                source_type: 'additional_charge',
                item: 'Additional Charge',
                date: bill.created_at ? new Date(bill.created_at).toLocaleDateString() : '',
                description: `${bill.description} (${bill.entity_type} #${bill.entity_id})`,
                quantity: 1,
                rate: bill.amount || 0,
                amount: bill.outstanding,
                total_amount: bill.amount,
                paid_amount: bill.paid_amount,
                status,
                selected: false,
            })
        }

        return items
    }

    const itemsForBroker = (brokerId: number): AvailableItem[] => {
        const items: AvailableItem[] = []
        const store = useBrokerVehicleBillsStore()

        for (const bill of store.bills) {
            if (bill.broker_id !== brokerId) continue
            if (bill.status === 'paid') continue
            const outstanding = bill.bill_amount - bill.paid_amount
            if (outstanding <= 0) continue

            items.push({
                id: bill.id,
                source_type: 'broker_vehicle_bill',
                item: 'Broker Vehicle Bill',
                date: bill.created_at ? new Date(bill.created_at).toLocaleDateString() : '',
                description: `Vehicle ${bill.vehicle_number} (Transport #${bill.transport_id})`,
                quantity: 1,
                rate: bill.bill_amount || 0,
                amount: outstanding,
                total_amount: bill.bill_amount,
                paid_amount: bill.paid_amount,
                status: bill.paid_amount > 0 ? 'partial' : 'unpaid',
                selected: false,
            })
        }

        return items
    }

    const itemsForGodown = (godownId: number): AvailableItem[] => {
        const items: AvailableItem[] = []
        const store = useGodownStoreBillsStore()

        for (const bill of store.storeBillData) {
            if (bill.godown_id !== godownId) continue
            const outstanding = bill.outstanding || 0
            if (outstanding <= 0) continue

            items.push({
                id: bill.id,
                source_type: 'godown_store_bill',
                item: 'Godown Store Bill',
                date: bill.billing_start ? new Date(bill.billing_start).toLocaleDateString() : '',
                description: `${bill.lot_name || `Lot #${bill.lot_id}`} @ ${bill.godown_name || `Godown #${bill.godown_id}`}`,
                quantity: 1,
                rate: bill.godown_cut || 0,
                amount: outstanding,
                total_amount: (bill.last_paid_amount || 0) + outstanding,
                paid_amount: bill.last_paid_amount || 0,
                status: (bill.last_paid_amount || 0) > 0 ? 'partial' : 'unpaid',
                selected: false,
            })
        }

        return items
    }

    const itemsForMajhi = (majhiId: number): AvailableItem[] => {
        const items: AvailableItem[] = []

        const lotStore = useMajhiLotBillsStore()
        for (const bill of lotStore.bills) {
            if (bill.majhi_id !== majhiId) continue
            const outstanding = bill.bill_amount - bill.paid_amount
            if (outstanding <= 0) continue
            items.push({
                id: bill.lot_id,
                source_type: 'majhi_lot_bill',
                item: 'Majhi Lot Bill',
                date: bill.created_at ? new Date(bill.created_at).toLocaleDateString() : '',
                description: `${bill.item_name || `Lot #${bill.lot_id}`} (${bill.majhi_bill_type})`,
                quantity: 1,
                rate: bill.majhi_cut || 0,
                amount: outstanding,
                total_amount: bill.bill_amount,
                paid_amount: bill.paid_amount,
                status: bill.paid_amount > 0 ? 'partial' : 'unpaid',
                selected: false,
            })
        }

        const loadingStore = useMajhiLoadingBillsStore()
        for (const bill of loadingStore.bills) {
            if (bill.majhi_id !== majhiId) continue
            const outstanding = bill.bill_amount - bill.paid_amount
            if (outstanding <= 0) continue
            items.push({
                id: bill.delivery_item_id,
                source_type: 'majhi_loading_bill',
                item: 'Majhi Loading Bill',
                date: bill.delivery_date ? new Date(bill.delivery_date).toLocaleDateString() : '',
                description: `${bill.item_name || `Delivery #${bill.delivery_id}`} (${bill.majhi_bill_type})`,
                quantity: 1,
                rate: bill.majhi_cut || 0,
                amount: outstanding,
                total_amount: bill.bill_amount,
                paid_amount: bill.paid_amount,
                status: bill.paid_amount > 0 ? 'partial' : 'unpaid',
                selected: false,
            })
        }

        return items
    }

    const getAvailableItemsForParty = (party: InvoiceParty, partyId: number): AvailableItem[] => {
        switch (party) {
            case 'customer': return itemsForCustomer(partyId)
            case 'broker': return itemsForBroker(partyId)
            case 'godown': return itemsForGodown(partyId)
            case 'majhi': return itemsForMajhi(partyId)
            default: return []
        }
    }

    // =========================================================================
    // INVOICE CREATION
    // =========================================================================

    const createInvoice = (
        party: InvoiceParty,
        partyId: number,
        partyName: string,
    ): Invoice => {
        const type: Invoice['type'] = party === 'customer' ? 'customer'
            : party === 'broker' ? 'broker_vehicle'
                : party === 'godown' ? 'godown_store'
                    : 'majhi'

        const sequence = invoices.value.filter(inv => inv.party_type === party).length + 1
        const number = generateInvoiceNumber(party, sequence)

        const newInvoice: Invoice = {
            id: `inv_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            number,
            party_type: party,
            party_id: partyId,
            party_name: partyName,
            type,
            date: new Date().toISOString().slice(0, 10),
            items: [],
            total: 0,
            received: 0,
            notes: '',
            status: 'draft',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        }

        currentInvoice.value = newInvoice
        return newInvoice
    }

    const buildInvoiceFromSelectedItems = (selectedItems: AvailableItem[]): Invoice => {
        if (!currentInvoice.value) {
            throw new Error('No current invoice')
        }

        const invoice = currentInvoice.value

        invoice.items = selectedItems.map(item => ({
            id: `item_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            source_type: item.source_type,
            source_id: item.id,
            item: item.item,
            date: item.date,
            description: item.description,
            quantity: item.quantity,
            rate: item.rate,
            amount: item.amount,
            total_amount: item.total_amount,
            paid_amount: item.paid_amount,
        }))

        invoice.total = invoice.items.reduce((sum, item) => sum + item.amount, 0)
        invoice.received = invoice.items.reduce((sum, item) => sum + item.paid_amount, 0)
        invoice.updated_at = new Date().toISOString()

        saveToLocalStorage()
        return invoice
    }

    const updateInvoiceField = <K extends keyof Invoice>(field: K, value: Invoice[K]) => {
        if (currentInvoice.value) {
            currentInvoice.value[field] = value
            currentInvoice.value.updated_at = new Date().toISOString()
            saveToLocalStorage()
        }
    }

    const saveCurrentInvoice = () => {
        if (!currentInvoice.value) return
        const existingIndex = invoices.value.findIndex(inv => inv.id === currentInvoice.value!.id)
        if (existingIndex !== -1) {
            invoices.value[existingIndex] = currentInvoice.value
        } else {
            invoices.value.push(currentInvoice.value)
        }
        saveToLocalStorage()
        push.success('Invoice saved successfully!')
    }

    const getInvoiceById = (id: string): Invoice | undefined => {
        return invoices.value.find(inv => inv.id === id)
    }

    const deleteInvoice = (id: string): boolean => {
        const index = invoices.value.findIndex(inv => inv.id === id)
        if (index !== -1) {
            invoices.value.splice(index, 1)
            saveToLocalStorage()
            push.success('Invoice deleted successfully!')
            return true
        }
        push.error('Invoice not found')
        return false
    }

    const clearCurrentInvoice = () => {
        currentInvoice.value = null
    }

    loadFromLocalStorage()

    return {
        invoices,
        currentInvoice,
        getAvailableItemsForParty,
        createInvoice,
        buildInvoiceFromSelectedItems,
        updateInvoiceField,
        saveCurrentInvoice,
        getInvoiceById,
        deleteInvoice,
        clearCurrentInvoice,
        saveToLocalStorage,
        loadFromLocalStorage,
    }
})