// src/stores/customerLotBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useLotsStore } from './lots'
import { useCustomersStore } from './customers'
import { useStoresStore } from './stores'
import { getOriginalStockTotals } from '@/utils/storeReconstruction'
import type { CustomerLotBill, CustomerLotBillSortField, SortDirection } from '@/types/customerLotBill'
import { push } from 'notivue'

export const useCustomerLotBillsStore = defineStore('customerLotBills', () => {
    const lotsStore = useLotsStore()
    const customersStore = useCustomersStore()
    const storesStore = useStoresStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<CustomerLotBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    const calculateBillAmount = (lot: any, stores: any[]): number => {
        const originals = getOriginalStockTotals(stores)

        if (lot.customer_charge_type === 'quantity') {
            return lot.unload_rate * originals.quantity
        } else {
            return lot.unload_rate * originals.weight
        }
    }

    const bills = computed<CustomerLotBill[]>(() => {
        const allLots = lotsStore.lots.filter(l => l.unload_rate > 0)

        const result: CustomerLotBill[] = []

        for (const lot of allLots) {
            const stores = storesStore.getStoresByLotId(lot.id)
            if (stores.length === 0) continue

            const customerId = lot.customer_id || null
            const customerName = customerId ? customersStore.getCustomerName(customerId) : 'No Customer'

            const billAmount = calculateBillAmount(lot, stores)
            if (billAmount === 0) continue

            const paidAmount = lot.customer_paid_unload_amount || 0
            const status = paidAmount >= billAmount ? 'paid' : 'unpaid'
            const paymentDate = status === 'paid' ? lot.customer_last_paid_through : null

            const originals = getOriginalStockTotals(stores)

            let quantityUnit = 'units'
            let weightUnit = 'kg'
            for (const store of stores) {
                if (store.quantity_unit) quantityUnit = store.quantity_unit
                if (store.weight_unit) weightUnit = store.weight_unit
            }

            result.push({
                id: lot.id,
                lot_id: lot.id,
                item_name: lotsStore.getLotDisplayName(lot),
                customer_id: customerId,
                customer_name: customerName,
                customer_charge_type: lot.customer_charge_type,
                unload_rate: lot.unload_rate,
                quantity: originals.quantity,
                quantity_unit: quantityUnit,
                weight: originals.weight,
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

        if (statusFilter.value) {
            result = result.filter(bill => bill.status === statusFilter.value)
        }

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

    const markBillAsPaid = async (lotId: number, payload: {
        amount: number
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

        const newTotalPaid = (lot.customer_paid_unload_amount || 0) + payload.amount

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

    const setSort = (field: CustomerLotBillSortField) => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
        } else {
            sortField.value = field
            sortDirection.value = 'asc'
        }
    }

    const setSearchQuery = (query: string) => {
        searchQuery.value = query
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    const setStatusFilter = (status: 'unpaid' | 'paid' | 'cancelled' | '') => {
        statusFilter.value = status
    }

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
        searchQuery,
        statusFilter,
        sortField,
        sortDirection,

        bills,
        filteredBills,
        totalBills,
        totalUnpaid,
        totalPaid,
        totalAmount,
        totalUnpaidAmount,

        markBillAsPaid,
        cancelBill,

        setSort,
        setSearchQuery,
        clearSearch,
        setStatusFilter,

        getBillById,
        getBillsByLotId,
        getBillsByCustomerId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        formatBillDate,
    }
})