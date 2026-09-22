// src/stores/customerStorageBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useLotsStore } from './lots'
import { useStoresStore } from './stores'
import { useCustomersStore } from './customers'
import { useDeliveryItemsStore } from './deliveryItems'
import { useDamagesStore } from './damages'

export interface CustomerStorageBill {
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
    const customersStore = useCustomersStore()
    const deliveryItemsStore = useDeliveryItemsStore()
    const damagesStore = useDamagesStore()

    const searchQuery = ref('')
    const monthFilter = ref('')
    const sortField = ref<'customer_name' | 'item_name' | 'monthly_bill' | 'outstanding' | 'total_billed'>('customer_name')
    const sortDirection = ref<'asc' | 'desc'>('asc')

    const getMonthsBetween = (start: Date, end: Date): number => {
        const months = (end.getFullYear() - start.getFullYear()) * 12 + (end.getMonth() - start.getMonth())
        return Math.max(1, months + 1)
    }

    const startOfMonth = (d: Date): Date => new Date(d.getFullYear(), d.getMonth(), 1)

    const getOutflowSince = (storeId: number, since: Date): { quantity: number; weight: number } => {
        let quantity = 0
        let weight = 0

        const deliveryItems = deliveryItemsStore.deliveryItems.filter(di => di.store_id === storeId)
        for (const di of deliveryItems) {
            const created = new Date(di.created_at)
            if (created >= since) {
                quantity += di.quantity || 0
                weight += di.weight || 0
            }
        }

        const damages = damagesStore.damages.filter(d => d.store_id === storeId)
        for (const d of damages) {
            const when = d.damage_date ? new Date(d.damage_date) : new Date(d.created_at)
            if (when >= since) {
                quantity += d.quantity || 0
                weight += d.weight || 0
            }
        }

        return { quantity, weight }
    }

    const getOpeningBalance = (
        currentQuantity: number,
        currentWeight: number,
        storeId: number,
        monthStart: Date,
    ): { quantity: number; weight: number } => {
        const outflow = getOutflowSince(storeId, monthStart)
        return {
            quantity: Math.max(0, currentQuantity + outflow.quantity),
            weight: Math.max(0, currentWeight + outflow.weight),
        }
    }

    const customerBillData = computed<CustomerStorageBill[]>(() => {
        const allLots = lotsStore.lots.filter(l => l.customer_storage_rate > 0)
        const now = new Date()
        const currentMonthStart = startOfMonth(now)

        const result: CustomerStorageBill[] = []

        for (const lot of allLots) {
            const stores = storesStore.getStoresByLotId(lot.id)
            if (stores.length === 0) continue

            let quantityUnit = 'units'
            let weightUnit = 'kg'
            let billingStart: Date | null = null
            let billingEnd: Date | null = null

            for (const store of stores) {
                if (store.quantity_unit) quantityUnit = store.quantity_unit
                if (store.weight_unit) weightUnit = store.weight_unit

                if (store.billing_start) {
                    const s = new Date(store.billing_start)
                    if (!billingStart || s < billingStart) billingStart = s
                }
                if (store.billing_end) {
                    const e = new Date(store.billing_end)
                    if (!billingEnd || e > billingEnd) billingEnd = e
                }
            }

            if (!billingStart) continue

            const billAsOf = billingEnd && billingEnd < currentMonthStart
                ? startOfMonth(billingEnd)
                : currentMonthStart

            let openingQuantity = 0
            let openingWeight = 0
            for (const store of stores) {
                const opening = getOpeningBalance(
                    store.quantity,
                    store.weight,
                    store.id,
                    billAsOf,
                )
                openingQuantity += opening.quantity
                openingWeight += opening.weight
            }

            let monthlyBill = 0
            if (lot.customer_charge_type === 'quantity') {
                monthlyBill = openingQuantity * lot.customer_storage_rate
            } else {
                monthlyBill = openingWeight * lot.customer_storage_rate
            }

            let totalBilled = 0
            const cursor = startOfMonth(billingStart)
            const lastMonth = billAsOf

            while (cursor <= lastMonth) {
                let monthOpeningQuantity = 0
                let monthOpeningWeight = 0
                for (const store of stores) {
                    const opening = getOpeningBalance(
                        store.quantity,
                        store.weight,
                        store.id,
                        cursor,
                    )
                    monthOpeningQuantity += opening.quantity
                    monthOpeningWeight += opening.weight
                }

                if (lot.customer_charge_type === 'quantity') {
                    totalBilled += monthOpeningQuantity * lot.customer_storage_rate
                } else {
                    totalBilled += monthOpeningWeight * lot.customer_storage_rate
                }

                cursor.setMonth(cursor.getMonth() + 1)
            }

            const monthsBilled = getMonthsBetween(billingStart, billAsOf)

            const customerId = lot.customer_id || null
            const customerName = customerId ? customersStore.getCustomerName(customerId) : 'No Customer'

            const totalPaid = lot.customer_last_paid_amount || 0
            const outstanding = totalBilled - totalPaid

            result.push({
                id: lot.id,
                lot_id: lot.id,
                item_name: lotsStore.getLotDisplayName(lot),
                customer_id: customerId,
                customer_name: customerName,
                customer_charge_type: lot.customer_charge_type,
                customer_storage_rate: lot.customer_storage_rate,
                quantity: openingQuantity,
                quantity_unit: quantityUnit,
                weight: openingWeight,
                weight_unit: weightUnit,
                monthly_bill: monthlyBill,
                billing_start: billingStart.toISOString(),
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

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(b =>
                b.customer_name.toLowerCase().includes(query) ||
                b.item_name.toLowerCase().includes(query) ||
                String(b.lot_id).includes(query)
            )
        }

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

    const recordCustomerPayment = async (lotId: number, amount: number, paidThrough: string | null): Promise<boolean> => {
        const lot = lotsStore.getLotById(lotId)
        if (!lot) {
            return false
        }

        const bill = customerBillData.value.find(b => b.lot_id === lotId)
        if (!bill) {
            return false
        }

        const newTotalPaid = (lot.customer_last_paid_amount || 0) + amount

        if (newTotalPaid > bill.total_billed) {
            return false
        }

        const result = await lotsStore.updateLot(lotId, {
            customer_last_paid_amount: newTotalPaid,
            customer_last_paid_through: paidThrough,
        })

        if (result) {
            await lotsStore.fetchLots()
            return true
        }

        return false
    }

    return {
        searchQuery,
        monthFilter,
        sortField,
        sortDirection,

        customerBillData,
        filteredCustomerBills,
        totalOutstanding,
        totalMonthlyBill,
        totalBilled,
        totalPaid,
        availableMonths,

        setSearchQuery,
        clearSearch,
        setSort,
        formatMonthYear,
        recordCustomerPayment,
    }
})