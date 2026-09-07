// src/stores/customerStorageBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useLotsStore } from './lots'
import { useStoresStore } from './stores'
import { useItemsStore } from './items'
import { useCustomersStore } from './customers'

interface CustomerStorageBill {
    id: number
    lot_id: number
    item_name: string
    customer_id: number | null
    customer_name: string
    customer_charge_type: 'weight' | 'quantity'
    customer_storage_rate: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    monthly_bill: number
    billing_start: string | null
    billing_end: string | null
    months_billed: number
    total_billed: number
    total_paid: number
    outstanding: number
    last_paid_through: string | null
    notes: string | null
}

export const useCustomerStorageBillsStore = defineStore('customerStorageBills', () => {
    const lotsStore = useLotsStore()
    const storesStore = useStoresStore()
    const itemsStore = useItemsStore()
    const customersStore = useCustomersStore()

    const searchQuery = ref('')
    const monthFilter = ref('')
    const sortField = ref<'customer_name' | 'item_name' | 'monthly_bill' | 'outstanding' | 'total_billed'>('customer_name')
    const sortDirection = ref<'asc' | 'desc'>('asc')

    // ============= HELPERS =============

    // Calculate months between two dates
    const getMonthsBetween = (start: Date, end: Date): number => {
        const months = (end.getFullYear() - start.getFullYear()) * 12 + (end.getMonth() - start.getMonth())
        return Math.max(1, months)
    }

    // Get current month start date
    const getCurrentMonthStart = (): Date => {
        const now = new Date()
        return new Date(now.getFullYear(), now.getMonth(), 1)
    }

    // ============= COMPUTED =============

    const customerBillData = computed<CustomerStorageBill[]>(() => {
        // Get all lots with customer_storage_rate > 0 (even if inactive)
        const allLots = lotsStore.lots.filter(l =>
            l.customer_storage_rate > 0
        )

        const currentMonthStart = getCurrentMonthStart()
        const result: CustomerStorageBill[] = []

        for (const lot of allLots) {
            // Get all stores for this lot
            const stores = storesStore.getStoresByLotId(lot.id)

            // Skip if no stores
            if (stores.length === 0) continue

            // Sum up quantities and weights from all stores
            let totalQuantity = 0
            let totalWeight = 0
            let quantityUnit = 'units'
            let weightUnit = 'kg'

            // Track billing period across all stores
            let billingStart: Date | null = null
            let billingEnd: Date | null = null

            for (const store of stores) {
                totalQuantity += store.quantity
                totalWeight += store.weight
                if (store.quantity_unit) quantityUnit = store.quantity_unit
                if (store.weight_unit) weightUnit = store.weight_unit

                // Track earliest billing_start
                if (store.billing_start) {
                    const start = new Date(store.billing_start)
                    if (!billingStart || start < billingStart) {
                        billingStart = start
                    }
                }

                // Track latest billing_end (null means active)
                if (store.billing_end) {
                    const end = new Date(store.billing_end)
                    if (!billingEnd || end > billingEnd) {
                        billingEnd = end
                    }
                }
            }

            // If no billing_start from stores, skip this lot
            if (!billingStart) {
                continue
            }

            // Calculate monthly bill
            let monthlyBill = 0
            if (lot.customer_charge_type === 'quantity') {
                monthlyBill = totalQuantity * lot.customer_storage_rate
            } else {
                monthlyBill = totalWeight * lot.customer_storage_rate
            }

            // Calculate duration in months
            const endDate = billingEnd || currentMonthStart
            const monthsBilled = getMonthsBetween(billingStart, endDate)

            // Total billed = monthly_bill × months_billed
            const totalBilled = monthlyBill * monthsBilled

            // Get item and customer info
            const item = itemsStore.getItemById(lot.item_id)
            const customerId = item?.customer_id || null
            const customerName = customerId ? customersStore.getCustomerName(customerId) : 'No Customer'

            // Total paid from lot
            const totalPaid = lot.customer_last_paid_amount || 0
            const outstanding = totalBilled - totalPaid

            // Always include the bill, even if fully paid (shows as paid status)
            result.push({
                id: lot.id,
                lot_id: lot.id,
                item_name: item ? itemsStore.getItemDisplayName(item) : `Item #${lot.item_id}`,
                customer_id: customerId,
                customer_name: customerName,
                customer_charge_type: lot.customer_charge_type,
                customer_storage_rate: lot.customer_storage_rate,
                quantity: totalQuantity,
                quantity_unit: quantityUnit,
                weight: totalWeight,
                weight_unit: weightUnit,
                monthly_bill: monthlyBill,
                billing_start: billingStart ? billingStart.toISOString() : null,
                billing_end: billingEnd ? billingEnd.toISOString() : null,
                months_billed: monthsBilled,
                total_billed: totalBilled,
                total_paid: totalPaid,
                outstanding: outstanding > 0 ? outstanding : 0,
                last_paid_through: lot.customer_last_paid_through || null,
                notes: lot.notes || null,
            })
        }

        return result
    })

    const filteredCustomerBills = computed(() => {
        let result = [...customerBillData.value]

        // Filter by search
        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(b =>
                b.customer_name.toLowerCase().includes(query) ||
                b.item_name.toLowerCase().includes(query) ||
                String(b.lot_id).includes(query)
            )
        }

        // Filter by month (billing_start month)
        if (monthFilter.value) {
            const [year, month] = monthFilter.value.split('-').map(Number)
            if (year && month && !isNaN(year) && !isNaN(month)) {
                const filterMonth = new Date(year, month - 1)
                result = result.filter(b => {
                    if (!b.billing_start) return false
                    const start = new Date(b.billing_start)
                    const startMonth = new Date(start.getFullYear(), start.getMonth(), 1)
                    const end = b.billing_end ? new Date(b.billing_end) : new Date()
                    const endMonth = new Date(end.getFullYear(), end.getMonth(), 1)
                    return startMonth <= filterMonth && filterMonth <= endMonth
                })
            }
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
                case 'monthly_bill':
                    comparison = a.monthly_bill - b.monthly_bill
                    break
                case 'total_billed':
                    comparison = a.total_billed - b.total_billed
                    break
                case 'outstanding':
                    comparison = a.outstanding - b.outstanding
                    break
                default:
                    comparison = 0
            }
            return sortDirection.value === 'desc' ? -comparison : comparison
        })

        return result
    })

    const totalOutstanding = computed(() => {
        return filteredCustomerBills.value.reduce((sum, b) => sum + b.outstanding, 0)
    })

    const totalMonthlyBill = computed(() => {
        return filteredCustomerBills.value.reduce((sum, b) => sum + b.monthly_bill, 0)
    })

    const totalBilled = computed(() => {
        return filteredCustomerBills.value.reduce((sum, b) => sum + b.total_billed, 0)
    })

    const totalPaid = computed(() => {
        return filteredCustomerBills.value.reduce((sum, b) => sum + b.total_paid, 0)
    })

    // ============= AVAILABLE MONTHS =============
    const availableMonths = computed(() => {
        const months = new Set<string>()
        customerBillData.value.forEach(b => {
            if (b.billing_start) {
                const start = new Date(b.billing_start)
                const month = `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}`
                months.add(month)
            }
            if (b.billing_end) {
                const end = new Date(b.billing_end)
                const month = `${end.getFullYear()}-${String(end.getMonth() + 1).padStart(2, '0')}`
                months.add(month)
            }
        })
        return Array.from(months).sort((a, b) => b.localeCompare(a))
    })

    // ============= ACTIONS =============

    const setSearchQuery = (query: string) => {
        searchQuery.value = query
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    const setSort = (field: 'customer_name' | 'item_name' | 'monthly_bill' | 'total_billed' | 'outstanding') => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
        } else {
            sortField.value = field
            sortDirection.value = 'asc'
        }
    }

    const formatMonthYear = (monthYear: string): string => {
        const [year, month] = monthYear.split('-')
        if (!year || !month) return monthYear
        const date = new Date(parseInt(year), parseInt(month) - 1)
        return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    }

    // Record payment for a customer (updates lot's customer_last_paid_amount and customer_last_paid_through)
    const recordCustomerPayment = async (lotId: number, amount: number, paidThrough: string | null): Promise<boolean> => {
        const lot = lotsStore.getLotById(lotId)
        if (!lot) return false

        const newTotalPaid = (lot.customer_last_paid_amount || 0) + amount

        const result = await lotsStore.updateLot(lotId, {
            customer_last_paid_amount: newTotalPaid,
            customer_last_paid_through: paidThrough,
        })

        return result !== null
    }

    return {
        // State
        searchQuery,
        monthFilter,
        sortField,
        sortDirection,

        // Computed
        customerBillData,
        filteredCustomerBills,
        totalOutstanding,
        totalMonthlyBill,
        totalBilled,
        totalPaid,
        availableMonths,

        // Actions
        setSearchQuery,
        clearSearch,
        setSort,
        formatMonthYear,
        recordCustomerPayment,
    }
})