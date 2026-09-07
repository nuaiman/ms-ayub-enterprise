// src/stores/invoices.ts

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
    Invoice,
    InvoiceType,
    AvailableItem,
} from '@/types/invoice'
import { generateInvoiceNumber } from '@/types/invoice'
import { useCustomerStorageBillsStore } from './customerStorageBills'
import { useCustomerLotBillsStore } from './customerLotBills'
import { useCustomerDeliveryBillsStore } from './customerDeliveryBills'
import { useCustomerTransportBillsStore } from './customerTransportBills'
import { push } from 'notivue'

const STORAGE_KEY = 'invoices_data'

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

    // Get available items (unpaid/partial bills) for a customer
    const getAvailableItemsForCustomer = (customerId: number, type: InvoiceType): AvailableItem[] => {
        const items: AvailableItem[] = []

        if (type === 'godown') {
            // Storage Bills
            const storageStore = useCustomerStorageBillsStore()
            const storageBills = storageStore.customerBillData.filter(
                bill => bill.customer_id === customerId && bill.outstanding > 0
            )
            for (const bill of storageBills) {
                const status = bill.outstanding > 0 && bill.total_paid > 0 ? 'partial' : 'unpaid'
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
                    status: status,
                    selected: false,
                })
            }

            // Lot Unload Bills
            const lotStore = useCustomerLotBillsStore()
            const lotBills = lotStore.bills.filter(
                bill => bill.customer_id === customerId && (bill.bill_amount - bill.paid_amount) > 0
            )
            for (const bill of lotBills) {
                const outstanding = bill.bill_amount - bill.paid_amount
                const status = outstanding > 0 && bill.paid_amount > 0 ? 'partial' : 'unpaid'
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
                    status: status,
                    selected: false,
                })
            }

            // Delivery Bills
            const deliveryStore = useCustomerDeliveryBillsStore()
            const deliveryBills = deliveryStore.bills.filter(
                bill => bill.customer_id === customerId && (bill.bill_amount - bill.paid_amount) > 0
            )
            for (const bill of deliveryBills) {
                const outstanding = bill.bill_amount - bill.paid_amount
                const status = outstanding > 0 && bill.paid_amount > 0 ? 'partial' : 'unpaid'
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
                    status: status,
                    selected: false,
                })
            }
        }

        if (type === 'transport') {
            // Transport Bills
            const transportStore = useCustomerTransportBillsStore()
            const transportBills = transportStore.bills.filter(
                bill => bill.customer_id === customerId && (bill.bill_amount - bill.paid_amount) > 0
            )
            for (const bill of transportBills) {
                const outstanding = bill.bill_amount - bill.paid_amount
                const status = outstanding > 0 && bill.paid_amount > 0 ? 'partial' : 'unpaid'
                items.push({
                    id: bill.transport_id,
                    source_type: 'transport_bill',
                    item: 'Transport Bill',
                    date: bill.created_at ? new Date(bill.created_at).toLocaleDateString() : '',
                    description: `${bill.from_location} → ${bill.to_location || 'N/A'} (${bill.total_vehicles} vehicles)`,
                    quantity: bill.total_vehicles || 1,
                    rate: bill.bill_amount / (bill.total_vehicles || 1),
                    amount: outstanding,
                    total_amount: bill.bill_amount,
                    paid_amount: bill.paid_amount,
                    status: status,
                    selected: false,
                })
            }
        }

        return items
    }

    const createInvoice = (type: InvoiceType, customerId: number, customerName: string): Invoice => {
        const sequence = invoices.value.filter(inv => inv.type === type).length + 1
        const number = generateInvoiceNumber(type, sequence)

        const newInvoice: Invoice = {
            id: `inv_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            number,
            type,
            customer_id: customerId,
            customer_name: customerName,
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

        // Build items
        invoice.items = selectedItems.map(item => ({
            id: `item_${Date.now()}_${Math.random().toString(36).substr(2, 6)}`,
            source_type: item.source_type,
            source_id: item.id,
            item: item.item,
            date: item.date,
            description: item.description,
            quantity: item.quantity,
            rate: item.rate,
            amount: item.amount,        // Outstanding amount
            total_amount: item.total_amount,
            paid_amount: item.paid_amount,
        }))

        // Calculate totals
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
        getAvailableItemsForCustomer,
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