// src/stores/customerLotBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useLotsStore } from './lots'
import { useItemsStore } from './items'
import { useCustomersStore } from './customers'
import { useStoresStore } from './stores'
import type { CustomerLotBill, CustomerLotBillSortField, SortDirection } from '@/types/customerLotBill'
import { push } from 'notivue'

export const useCustomerLotBillsStore = defineStore('customerLotBills', () => {
    const lotsStore = useLotsStore()
    const itemsStore = useItemsStore()
    const customersStore = useCustomersStore()
    const storesStore = useStoresStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<CustomerLotBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    // ============= HELPERS =============

    // Calculate bill amount for a lot (unload_rate × quantity/weight)
    const calculateBillAmount = (lot: any, stores: any[]): number => {
        let totalQuantity = 0
        let totalWeight = 0

        // Sum up quantities and weights from all stores
        for (const store of stores) {
            totalQuantity += store.quantity
            totalWeight += store.weight
        }

        // Calculate based on customer_charge_type
        if (lot.customer_charge_type === 'quantity') {
            return lot.unload_rate * totalQuantity
        } else {
            return lot.unload_rate * totalWeight
        }
    }

    // ============= COMPUTED =============

    // Generate bills for all lots with unload_rate > 0
    const bills = computed<CustomerLotBill[]>(() => {
        // Get all lots with unload_rate > 0 (even if inactive)
        const allLots = lotsStore.lots.filter(l =>
            l.unload_rate > 0
        )

        const result: CustomerLotBill[] = []

        for (const lot of allLots) {
            // Get all stores for this lot
            const stores = storesStore.getStoresByLotId(lot.id)

            // Skip if no stores
            if (stores.length === 0) continue

            // Get item and customer info
            const item = itemsStore.getItemById(lot.item_id)
            const customerId = item?.customer_id || null
            const customerName = customerId ? customersStore.getCustomerName(customerId) : 'No Customer'

            // Calculate bill amount
            const billAmount = calculateBillAmount(lot, stores)

            // Skip if bill amount is 0
            if (billAmount === 0) continue

            // Get paid amount from lot
            const paidAmount = lot.customer_paid_unload_amount || 0

            // Determine status - 'paid' only if fully paid, otherwise 'unpaid'
            const status = paidAmount >= billAmount ? 'paid' : 'unpaid'

            // Use lot's customer_last_paid_through as payment date if paid
            const paymentDate = status === 'paid' ? lot.customer_last_paid_through : null

            // Sum up quantities and weights
            let totalQuantity = 0
            let totalWeight = 0
            let quantityUnit = 'units'
            let weightUnit = 'kg'

            for (const store of stores) {
                totalQuantity += store.quantity
                totalWeight += store.weight
                if (store.quantity_unit) quantityUnit = store.quantity_unit
                if (store.weight_unit) weightUnit = store.weight_unit
            }

            result.push({
                id: lot.id,
                lot_id: lot.id,
                item_name: item ? itemsStore.getItemDisplayName(item) : `Item #${lot.item_id}`,
                customer_id: customerId,
                customer_name: customerName,
                customer_charge_type: lot.customer_charge_type,
                unload_rate: lot.unload_rate,
                quantity: totalQuantity,
                quantity_unit: quantityUnit,
                weight: totalWeight,
                weight_unit: weightUnit,
                bill_amount: billAmount,
                paid_amount: paidAmount,
                status: status,
                payment_date: paymentDate,
                notes: lot.notes || null,
                created_at: lot.created_at,
                updated_at: lot.updated_at,
            })
        }

        return result
    })

    const filteredBills = computed(() => {
        let result = [...bills.value]

        // Filter by search query
        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(bill =>
                bill.customer_name.toLowerCase().includes(query) ||
                bill.item_name.toLowerCase().includes(query) ||
                String(bill.lot_id).includes(query) ||
                bill.status.toLowerCase().includes(query) ||
                (bill.notes && bill.notes.toLowerCase().includes(query))
            )
        }

        // Filter by status
        if (statusFilter.value) {
            result = result.filter(bill => bill.status === statusFilter.value)
        }

        // Sort
        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'customer_name':
                    comparison = a.customer_name.localeCompare(b.customer_name)
                    break
                case 'item_name':
                    comparison = a.item_name.localeCompare(b.item_name)
                    break
                case 'bill_amount':
                    comparison = a.bill_amount - b.bill_amount
                    break
                case 'status':
                    comparison = a.status.localeCompare(b.status)
                    break
                case 'created_at':
                    comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
                    break
                default:
                    comparison = 0
            }
            return sortDirection.value === 'desc' ? -comparison : comparison
        })

        return result
    })

    const totalBills = computed(() => bills.value.length)
    const totalUnpaid = computed(() => bills.value.filter(b => b.status === 'unpaid').length)
    const totalPaid = computed(() => bills.value.filter(b => b.status === 'paid').length)
    const totalAmount = computed(() => bills.value.reduce((sum, b) => sum + b.bill_amount, 0))
    const totalUnpaidAmount = computed(() => {
        return bills.value
            .filter(b => b.status === 'unpaid')
            .reduce((sum, b) => sum + b.bill_amount, 0)
    })

    // ============= ACTIONS =============

    // Mark a bill as paid (updates the lot's customer_paid_unload_amount)
    const markBillAsPaid = async (lotId: number, payload: {
        amount: number  // Payment amount
        payment_date?: string
        notes?: string | null
    }): Promise<boolean> => {
        const lot = lotsStore.getLotById(lotId)
        if (!lot) {
            push.error('Lot not found')
            return false
        }

        const bill = bills.value.find(b => b.lot_id === lotId)
        if (!bill) {
            push.error('Bill not found')
            return false
        }

        if (bill.status === 'paid') {
            push.info('Bill is already fully paid')
            return true
        }

        // Calculate new total paid (add payment amount to existing paid)
        const newTotalPaid = (lot.customer_paid_unload_amount || 0) + payload.amount

        // Validate: Cannot pay more than bill amount
        if (newTotalPaid > bill.bill_amount) {
            const remaining = bill.bill_amount - bill.paid_amount
            push.error(`Payment amount exceeds remaining balance of ${remaining.toFixed(2)}`)
            return false
        }

        const result = await lotsStore.updateLotCustomerUnloadPayment(
            lotId,
            newTotalPaid,
            payload.payment_date || null
        )

        if (result && payload.notes) {
            await lotsStore.updateLot(lotId, {
                notes: payload.notes
            })
        }

        if (result) {
            await lotsStore.fetchLots()
            push.success('Payment recorded successfully')
            return true
        }

        return false
    }

    // Cancel a bill (sets unload_rate to 0)
    const cancelBill = async (lotId: number): Promise<boolean> => {
        const lot = lotsStore.getLotById(lotId)
        if (!lot) {
            push.error('Lot not found')
            return false
        }

        const result = await lotsStore.updateLot(lotId, {
            unload_rate: 0,
            notes: `Bill cancelled at ${new Date().toISOString()}${lot.notes ? ` - ${lot.notes}` : ''}`,
        })

        if (result) {
            await lotsStore.fetchLots()
            push.success('Bill cancelled successfully')
            return true
        }

        return false
    }

    // ============= SORT =============
    const setSort = (field: CustomerLotBillSortField) => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
        } else {
            sortField.value = field
            sortDirection.value = 'asc'
        }
    }

    // ============= SEARCH =============
    const setSearchQuery = (query: string) => {
        searchQuery.value = query
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    const setStatusFilter = (status: 'unpaid' | 'paid' | 'cancelled' | '') => {
        statusFilter.value = status
    }

    // ============= UTILITIES =============

    const getBillById = (id: number): CustomerLotBill | undefined => {
        return bills.value.find(b => b.id === id)
    }

    const getBillsByLotId = (lotId: number): CustomerLotBill | undefined => {
        return bills.value.find(b => b.lot_id === lotId)
    }

    const getBillsByCustomerId = (customerId: number): CustomerLotBill[] => {
        return bills.value.filter(b => b.customer_id === customerId)
    }

    const getStatusBadgeClass = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'border-(--color-yellow) text-(--color-yellow)'
            case 'paid':
                return 'border-(--color-green) text-(--color-green)'
            case 'cancelled':
                return 'border-(--color-red) text-(--color-red)'
            default:
                return 'border-(--color-border) text-(--color-text-secondary)'
        }
    }

    const getStatusDotClass = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'bg-(--color-yellow)'
            case 'paid':
                return 'bg-(--color-green)'
            case 'cancelled':
                return 'bg-(--color-red)'
            default:
                return 'bg-(--color-text-secondary)'
        }
    }

    const getStatusLabel = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'Unpaid'
            case 'paid':
                return 'Paid'
            case 'cancelled':
                return 'Cancelled'
            default:
                return status
        }
    }

    const formatBillDate = (dateStr: string): string => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        })
    }

    return {
        // State
        searchQuery,
        statusFilter,
        sortField,
        sortDirection,

        // Computed
        bills,
        filteredBills,
        totalBills,
        totalUnpaid,
        totalPaid,
        totalAmount,
        totalUnpaidAmount,

        // Actions
        markBillAsPaid,
        cancelBill,

        // Sort
        setSort,

        // Search
        setSearchQuery,
        clearSearch,
        setStatusFilter,

        // Utilities
        getBillById,
        getBillsByLotId,
        getBillsByCustomerId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        formatBillDate,
    }
})