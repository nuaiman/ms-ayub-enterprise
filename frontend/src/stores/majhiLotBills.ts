// src/stores/majhiLotBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useLotsStore } from './lots'
import { useItemsStore } from './items'
import { useMajhisStore } from './majhis'
import { useStoresStore } from './stores'
import type { MajhiLotBill, MajhiLotBillSortField, SortDirection } from '@/types/majhiLotBill'
import { push } from 'notivue'

export const useMajhiLotBillsStore = defineStore('majhiLotBills', () => {
    const lotsStore = useLotsStore()
    const itemsStore = useItemsStore()
    const majhisStore = useMajhisStore()
    const storesStore = useStoresStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<MajhiLotBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    // ============= HELPERS =============

    // Calculate bill amount for a lot based on majhi_bill_type
    const calculateBillAmount = (lot: any, stores: any[]): number => {
        let totalQuantity = 0
        let totalWeight = 0

        // Sum up quantities and weights from all stores
        for (const store of stores) {
            totalQuantity += store.quantity
            totalWeight += store.weight
        }

        // Calculate based on majhi_bill_type
        switch (lot.majhi_bill_type) {
            case 'quantity':
                return lot.majhi_cut * totalQuantity
            case 'weight':
                return lot.majhi_cut * totalWeight
            case 'job':
                return lot.majhi_cut // Fixed amount for job
            default:
                return 0
        }
    }

    // ============= COMPUTED =============

    // Generate bills for all lots with majhi_id and majhi_cut > 0
    const bills = computed<MajhiLotBill[]>(() => {
        // Get all active lots with majhi_id and majhi_cut > 0
        const activeLots = lotsStore.lots.filter(l =>
            l.is_active &&
            l.majhi_id !== null &&
            l.majhi_cut > 0
        )

        const result: MajhiLotBill[] = []

        for (const lot of activeLots) {
            // Get all stores for this lot
            const stores = storesStore.getStoresByLotId(lot.id)

            // Skip if no stores
            if (stores.length === 0) continue

            // Get item and majhi info
            const item = itemsStore.getItemById(lot.item_id)
            const majhiName = lot.majhi_id ? majhisStore.getMajhiName(lot.majhi_id) : 'No Majhi'

            // Calculate bill amount
            const billAmount = calculateBillAmount(lot, stores)

            // Get paid amount from lot
            const paidAmount = lot.majhi_total_paid || 0

            // Determine status
            let status: 'unpaid' | 'paid' | 'cancelled' = 'unpaid'
            if (paidAmount >= billAmount) {
                status = 'paid'
            }

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
                majhi_id: lot.majhi_id,
                majhi_name: majhiName,
                majhi_bill_type: lot.majhi_bill_type,
                majhi_cut: lot.majhi_cut,
                quantity: totalQuantity,
                quantity_unit: quantityUnit,
                weight: totalWeight,
                weight_unit: weightUnit,
                bill_amount: billAmount,
                paid_amount: paidAmount,
                status: status,
                payment_date: null, // We don't track payment date separately for majhi
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
                bill.majhi_name.toLowerCase().includes(query) ||
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
                case 'majhi_name':
                    comparison = a.majhi_name.localeCompare(b.majhi_name)
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

    // Mark a bill as paid (updates the lot's majhi_total_paid)
    const markBillAsPaid = async (lotId: number, payload: {
        amount: number  // Payment amount
        notes?: string | null
    }): Promise<boolean> => {
        try {
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
            const newTotalPaid = (lot.majhi_total_paid || 0) + payload.amount

            // Validate: Cannot pay more than bill amount
            if (newTotalPaid > bill.bill_amount) {
                const remaining = bill.bill_amount - bill.paid_amount
                push.error(`Payment amount exceeds remaining balance of ${remaining.toFixed(2)}`)
                return false
            }

            const result = await lotsStore.updateLotMajhiPayment(lotId, newTotalPaid)

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
        } catch (error) {
            console.error('Error marking majhi bill as paid:', error)
            push.error('Failed to record payment')
            return false
        }
    }

    // Cancel a bill (sets majhi_cut to 0 and removes majhi_id)
    const cancelBill = async (lotId: number): Promise<boolean> => {
        try {
            const lot = lotsStore.getLotById(lotId)
            if (!lot) {
                push.error('Lot not found')
                return false
            }

            const result = await lotsStore.updateLot(lotId, {
                majhi_cut: 0,
                majhi_id: null,
                notes: `Majhi bill cancelled at ${new Date().toISOString()}${lot.notes ? ` - ${lot.notes}` : ''}`,
            })

            if (result) {
                await lotsStore.fetchLots()
                push.success('Bill cancelled successfully')
                return true
            }

            return false
        } catch (error) {
            console.error('Error cancelling majhi bill:', error)
            push.error('Failed to cancel bill')
            return false
        }
    }

    // ============= SORT =============
    const setSort = (field: MajhiLotBillSortField) => {
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

    const getBillById = (id: number): MajhiLotBill | undefined => {
        return bills.value.find(b => b.id === id)
    }

    const getBillsByLotId = (lotId: number): MajhiLotBill | undefined => {
        return bills.value.find(b => b.lot_id === lotId)
    }

    const getBillsByMajhiId = (majhiId: number): MajhiLotBill[] => {
        return bills.value.filter(b => b.majhi_id === majhiId)
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

    const getBillTypeLabel = (billType: string): string => {
        switch (billType) {
            case 'quantity':
                return 'Per Quantity'
            case 'weight':
                return 'Per Weight'
            case 'job':
                return 'Job (Fixed)'
            default:
                return billType
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
        getBillsByMajhiId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        getBillTypeLabel,
        formatBillDate,
    }
})