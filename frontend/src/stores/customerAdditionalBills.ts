// src/stores/customerAdditionalBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    AdditionalCharge,
    CustomerAdditionalBill,
    CustomerAdditionalBillSortField,
    SortDirection,
} from '@/types/customerAdditionalBill'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'
import { useCustomersStore } from './customers'

export const useCustomerAdditionalBillsStore = defineStore('customerAdditionalBills', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    const customersStore = useCustomersStore()

    // ============= STATE =============
    const charges = ref<AdditionalCharge[]>([])
    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | ''>('')
    const sortField = ref<CustomerAdditionalBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    // ============= COMPUTED =============

    const bills = computed<CustomerAdditionalBill[]>(() => {
        return charges.value.map(charge => {
            const customerName = customersStore.getCustomerName(charge.customer_id)
            const paid = charge.customer_total_paid || 0
            const outstanding = Math.max(0, charge.amount - paid)
            const status: 'unpaid' | 'paid' = paid >= charge.amount ? 'paid' : 'unpaid'

            return {
                id: charge.id,
                customer_id: charge.customer_id,
                customer_name: customerName,
                amount: charge.amount,
                description: charge.description,
                paid_amount: paid,
                outstanding,
                status,
                payment_date: charge.customer_total_paid_through,
                created_at: charge.created_at,
                updated_at: charge.updated_at,
            }
        })
    })

    const filteredBills = computed(() => {
        let result = [...bills.value]

        if (searchQuery.value) {
            const q = searchQuery.value.toLowerCase()
            result = result.filter(bill =>
                bill.customer_name.toLowerCase().includes(q) ||
                bill.description.toLowerCase().includes(q) ||
                bill.status.toLowerCase().includes(q)
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
                case 'amount':
                    comparison = a.amount - b.amount
                    break
                case 'description':
                    comparison = a.description.localeCompare(b.description)
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
    const totalAmount = computed(() => bills.value.reduce((sum, b) => sum + b.amount, 0))
    const totalUnpaidAmount = computed(() =>
        bills.value
            .filter(b => b.status === 'unpaid')
            .reduce((sum, b) => sum + b.outstanding, 0)
    )

    // ============= ACTIONS =============

    const fetchCharges = async () => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<AdditionalCharge[]>>('/customer-additional-charges')
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            charges.value = res.data.data
            return charges.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch additional charges')
            return []
        } finally {
            destroyLoader()
        }
    }

    const createCharge = async (payload: {
        customer_id: number
        amount: number
        description: string
    }): Promise<AdditionalCharge | null> => {
        displayLoader()
        try {
            const res = await api.post<ApiResponse<AdditionalCharge>>('/customer-additional-charges', payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            charges.value.push(res.data.data)
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to create additional charge')
            return null
        } finally {
            destroyLoader()
        }
    }

    const updateCharge = async (
        id: number,
        payload: { amount?: number; description?: string }
    ): Promise<AdditionalCharge | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<AdditionalCharge>>(`/customer-additional-charges/${id}`, payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            const index = charges.value.findIndex(c => c.id === id)
            if (index !== -1) {
                charges.value[index] = res.data.data
            }
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to update additional charge')
            return null
        } finally {
            destroyLoader()
        }
    }

    const markBillAsPaid = async (
        id: number,
        payload: { amount: number; payment_date?: string; notes?: string | null }
    ): Promise<boolean> => {
        try {
            const charge = charges.value.find(c => c.id === id)
            if (!charge) {
                push.error('Additional charge not found')
                return false
            }

            const bill = bills.value.find(b => b.id === id)
            if (!bill) {
                push.error('Bill not found')
                return false
            }

            if (bill.status === 'paid') {
                push.info('Bill is already fully paid')
                return true
            }

            const newTotalPaid = (charge.customer_total_paid || 0) + payload.amount

            if (newTotalPaid > bill.amount) {
                const remaining = bill.amount - bill.paid_amount
                push.error(`Payment amount exceeds remaining balance of ${remaining.toFixed(2)}`)
                return false
            }

            displayLoader()
            try {
                const res = await api.patch<ApiResponse<AdditionalCharge>>(
                    `/customer-additional-charges/${id}/customer-payment`,
                    {
                        customer_total_paid: newTotalPaid,
                        customer_total_paid_through: payload.payment_date || null,
                    }
                )

                if (!res.data.success) {
                    push.error(res.data.message)
                    return false
                }

                const index = charges.value.findIndex(c => c.id === id)
                if (index !== -1) {
                    charges.value[index] = res.data.data
                }

                push.success('Payment recorded successfully')
                return true
            } finally {
                destroyLoader()
            }
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to record payment')
            return false
        }
    }

    const deleteCharge = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<null>>(`/customer-additional-charges/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            charges.value = charges.value.filter(c => c.id !== id)
            push.success(res.data.message)
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete additional charge')
            return false
        } finally {
            destroyLoader()
        }
    }

    // ============= SORT / SEARCH =============

    const setSort = (field: CustomerAdditionalBillSortField) => {
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

    const setStatusFilter = (status: 'unpaid' | 'paid' | '') => {
        statusFilter.value = status
    }

    // ============= UTILITIES =============

    const getStatusBadgeClass = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'border-(--color-yellow) text-(--color-yellow)'
            case 'paid':
                return 'border-(--color-green) text-(--color-green)'
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
            minute: '2-digit',
        })
    }

    const getChargesByCustomerId = (customerId: number): CustomerAdditionalBill[] => {
        return bills.value.filter(b => b.customer_id === customerId)
    }

    return {
        // State
        charges,
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
        fetchCharges,
        createCharge,
        updateCharge,
        markBillAsPaid,
        deleteCharge,

        // Sort / Search
        setSort,
        setSearchQuery,
        clearSearch,
        setStatusFilter,

        // Utilities
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        formatBillDate,
        getChargesByCustomerId,
    }
})