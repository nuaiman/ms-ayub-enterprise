// src/stores/godownStoreBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GodownStoreBillStore, GodownStoreBillSortField, SortDirection } from '@/types/godownStoreBill'
import { useStoresStore } from './stores'
import { useLotsStore } from './lots'
import { useGodownsStore } from './godowns'
import { useCustomersStore } from './customers'
import { getOriginalStock } from '@/utils/storeReconstruction'

export const useGodownStoreBillsStore = defineStore('godownStoreBills', () => {
    const storesStore = useStoresStore()
    const lotsStore = useLotsStore()
    const godownsStore = useGodownsStore()
    const customersStore = useCustomersStore()

    const searchQuery = ref('')
    const monthFilter = ref('')
    const sortField = ref<GodownStoreBillSortField>('lot_name')
    const sortDirection = ref<SortDirection>('asc')

    const storeBillData = computed<GodownStoreBillStore[]>(() => {
        const stores = storesStore.stores.filter(s =>
            s.is_active &&
            s.godown_cut > 0 &&
            (s.quantity > 0 || s.weight > 0)
        )

        return stores.map(store => {
            const lot = lotsStore.getLotById(store.lot_id)
            const godown = godownsStore.getGodownById(store.godown_id)
            const customer = lot?.customer_id ? customersStore.getCustomerById(lot.customer_id) : undefined

            const original = getOriginalStock(store)
            let monthlyBill = 0
            if (store.store_bill_type === 'quantity') {
                monthlyBill = original.quantity * store.godown_cut
            } else {
                monthlyBill = original.weight * store.godown_cut
            }

            const totalBilled = store.last_paid_amount || 0
            const outstanding = monthlyBill - totalBilled

            return {
                ...store,
                lot_name: lot ? lotsStore.getLotDisplayName(lot) : `Lot #${store.lot_id}`,
                item_name: lot ? lotsStore.getLotDisplayName(lot) : 'Unknown',
                customer_name: customer ? (customer.company_name || customer.contact_person || `Customer #${customer.id}`) : 'N/A',
                godown_name: godown ? godown.name : `Godown #${store.godown_id}`,
                monthly_bill: monthlyBill,
                total_billed: totalBilled,
                outstanding: outstanding > 0 ? outstanding : 0
            }
        })
    })

    const filteredStoreBills = computed(() => {
        let result = [...storeBillData.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(s =>
                s.lot_name?.toLowerCase().includes(query) ||
                s.item_name?.toLowerCase().includes(query) ||
                s.customer_name?.toLowerCase().includes(query) ||
                s.godown_name?.toLowerCase().includes(query)
            )
        }

        if (monthFilter.value) {
            const [year, month] = monthFilter.value.split('-').map(Number)

            if (year && month && !isNaN(year) && !isNaN(month)) {
                const billMonth = new Date(year, month - 1)

                result = result.filter(s => {
                    const start = new Date(s.billing_start)

                    if (s.billing_end) {
                        const end = new Date(s.billing_end)
                        return start <= billMonth && billMonth <= end
                    }
                    return start <= billMonth
                })
            }
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'lot_name':
                    comparison = (a.lot_name || '').localeCompare(b.lot_name || '')
                    break
                case 'godown_name':
                    comparison = (a.godown_name || '').localeCompare(b.godown_name || '')
                    break
                case 'customer_name':
                    comparison = (a.customer_name || '').localeCompare(b.customer_name || '')
                    break
                case 'monthly_bill':
                    comparison = (a.monthly_bill || 0) - (b.monthly_bill || 0)
                    break
                case 'total_billed':
                    comparison = (a.total_billed || 0) - (b.total_billed || 0)
                    break
                case 'outstanding':
                    comparison = (a.outstanding || 0) - (b.outstanding || 0)
                    break
                case 'last_paid_through':
                    comparison = (a.last_paid_through || '').localeCompare(b.last_paid_through || '')
                    break
                default:
                    comparison = 0
            }
            return sortDirection.value === 'desc' ? -comparison : comparison
        })

        return result
    })

    const totalOutstanding = computed(() => {
        return filteredStoreBills.value.reduce((sum, s) => sum + (s.outstanding || 0), 0)
    })

    const totalMonthlyBill = computed(() => {
        return filteredStoreBills.value.reduce((sum, s) => sum + (s.monthly_bill || 0), 0)
    })

    const totalPaid = computed(() => {
        return filteredStoreBills.value.reduce((sum, s) => sum + (s.total_billed || 0), 0)
    })

    const availableMonths = computed(() => {
        const months = new Set<string>()
        storeBillData.value.forEach(s => {
            const start = new Date(s.billing_start)
            const month = `${start.getFullYear()}-${String(start.getMonth() + 1).padStart(2, '0')}`
            months.add(month)
            if (s.billing_end) {
                const end = new Date(s.billing_end)
                const endMonth = `${end.getFullYear()}-${String(end.getMonth() + 1).padStart(2, '0')}`
                months.add(endMonth)
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

    const setSort = (field: GodownStoreBillSortField) => {
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

    const getStatusBadgeClass = (store: GodownStoreBillStore): string => {
        if (!store.is_active) return 'border-(--color-red) text-(--color-red)'
        if (store.outstanding && store.outstanding > 0) return 'border-(--color-yellow) text-(--color-yellow)'
        if (store.last_paid_amount > 0) return 'border-(--color-green) text-(--color-green)'
        return 'border-(--color-blue) text-(--color-blue)'
    }

    const getStatusLabel = (store: GodownStoreBillStore): string => {
        if (!store.is_active) return 'Inactive'
        if (store.outstanding && store.outstanding > 0) return 'Outstanding'
        if (store.last_paid_amount > 0) return 'Paid'
        return 'Active'
    }

    return {
        searchQuery,
        monthFilter,
        sortField,
        sortDirection,

        storeBillData,
        filteredStoreBills,
        totalOutstanding,
        totalMonthlyBill,
        totalPaid,
        availableMonths,

        setSearchQuery,
        clearSearch,
        setSort,
        formatMonthYear,
        getStatusBadgeClass,
        getStatusLabel,
    }
})