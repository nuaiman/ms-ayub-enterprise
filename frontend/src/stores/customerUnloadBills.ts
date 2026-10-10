// src/stores/customerUnloadBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    CustomerUnloadBill,
    CreateCustomerUnloadBillPayload,
    UpdateCustomerUnloadBillPayload,
    CreateBillPaymentPayload,
    CustomerUnloadBillSortField,
    SortDirection,
} from '@/types/customerUnloadBill'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useCustomerUnloadBillsStore = defineStore('customerUnloadBills', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    const bills = ref<CustomerUnloadBill[]>([])
    const searchQuery = ref('')
    const sortField = ref<CustomerUnloadBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    const filteredBills = computed(() => {
        let result = [...bills.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(b =>
                b.bill_type.toLowerCase().includes(query) ||
                String(b.rate).includes(query) ||
                String(b.total_paid).includes(query)
            )
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'customer_id':
                    comparison = a.customer_id - b.customer_id
                    break
                case 'store_id':
                    comparison = a.store_id - b.store_id
                    break
                case 'bill_type':
                    comparison = a.bill_type.localeCompare(b.bill_type)
                    break
                case 'rate':
                    comparison = a.rate - b.rate
                    break
                case 'total_paid':
                    comparison = a.total_paid - b.total_paid
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
    const totalBilled = computed(() => bills.value.reduce((sum, b) => sum + b.total_amount, 0))
    const totalPaid = computed(() => bills.value.reduce((sum, b) => sum + b.total_paid, 0))

    // =========================================================================
    // FETCH
    // =========================================================================

    const fetchCustomerUnloadBills = async (params?: {
        customer_id?: number
        store_id?: number
    }) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<CustomerUnloadBill[]>>(
                '/customer-unload-bills',
                { params }
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            bills.value = res.data.data
            return bills.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch customer unload bills')
            return []
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // CREATE
    // =========================================================================

    const createCustomerUnloadBill = async (
        payload: CreateCustomerUnloadBillPayload
    ): Promise<CustomerUnloadBill | null> => {
        displayLoader()
        try {
            if (!payload.customer_id) {
                push.error('Customer is required')
                return null
            }
            if (!payload.store_id) {
                push.error('Store is required')
                return null
            }

            const res = await api.post<ApiResponse<CustomerUnloadBill>>(
                '/customer-unload-bills',
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            bills.value.unshift(res.data.data)
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to create customer unload bill')
            return null
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // UPDATE
    // =========================================================================

    const updateCustomerUnloadBill = async (
        id: number,
        payload: UpdateCustomerUnloadBillPayload
    ): Promise<CustomerUnloadBill | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<CustomerUnloadBill>>(
                `/customer-unload-bills/${id}`,
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            const index = bills.value.findIndex(b => b.id === id)
            if (index !== -1) bills.value[index] = res.data.data
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to update customer unload bill')
            return null
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // PAYMENTS
    // =========================================================================

    const createCustomerUnloadBillPayment = async (
        billId: number,
        payload: CreateBillPaymentPayload
    ): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.post<ApiResponse<unknown>>(
                `/customer-unload-bills/${billId}/payments`,
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            push.success(res.data.message || 'Payment recorded')
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to record payment')
            return false
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // DELETE
    // =========================================================================

    const deleteCustomerUnloadBill = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<null>>(`/customer-unload-bills/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            bills.value = bills.value.filter(b => b.id !== id)
            push.success(res.data.message)
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete customer unload bill')
            return false
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // SORT / SEARCH
    // =========================================================================

    const setSort = (field: CustomerUnloadBillSortField) => {
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

    const getBillById = (id: number): CustomerUnloadBill | undefined =>
        bills.value.find(b => b.id === id)

    return {
        bills,
        searchQuery,
        sortField,
        sortDirection,

        filteredBills,
        totalBills,
        totalBilled,
        totalPaid,

        fetchCustomerUnloadBills,
        createCustomerUnloadBill,
        updateCustomerUnloadBill,
        createCustomerUnloadBillPayment,
        deleteCustomerUnloadBill,

        setSort,
        setSearchQuery,
        clearSearch,

        getBillById,
    }
})