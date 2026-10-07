// src/stores/customerAdditionalBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    CustomerAdditionalBill,
    CreateCustomerAdditionalBillPayload,
    UpdateCustomerAdditionalBillPayload,
    CreateBillPaymentPayload,
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

    const bills = ref<CustomerAdditionalBill[]>([])
    const searchQuery = ref('')
    const sortField = ref<CustomerAdditionalBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    const filteredBills = computed(() => {
        let result = [...bills.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            const customersStore = useCustomersStore()
            result = result.filter(b =>
                b.description.toLowerCase().includes(query) ||
                String(b.amount).includes(query) ||
                String(b.total_paid).includes(query) ||
                customersStore.getCustomerName(b.customer_id).toLowerCase().includes(query)
            )
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'customer_id':
                    comparison = a.customer_id - b.customer_id
                    break
                case 'amount':
                    comparison = a.amount - b.amount
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
    const totalBilled = computed(() => bills.value.reduce((sum, b) => sum + b.amount, 0))
    const totalPaid = computed(() => bills.value.reduce((sum, b) => sum + b.total_paid, 0))

    // =========================================================================
    // FETCH
    // =========================================================================

    const fetchCustomerAdditionalBills = async (params?: { customer_id?: number }) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<CustomerAdditionalBill[]>>(
                '/customer-additional-bills',
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
            push.error(err.response?.data?.message || 'Failed to fetch customer additional bills')
            return []
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // CREATE
    // =========================================================================

    const createCustomerAdditionalBill = async (
        payload: CreateCustomerAdditionalBillPayload
    ): Promise<CustomerAdditionalBill | null> => {
        displayLoader()
        try {
            if (!payload.customer_id) {
                push.error('Customer is required')
                return null
            }
            if (!payload.amount || payload.amount <= 0) {
                push.error('Amount must be greater than 0')
                return null
            }
            if (!payload.description || !payload.description.trim()) {
                push.error('Description is required')
                return null
            }

            const res = await api.post<ApiResponse<CustomerAdditionalBill>>(
                '/customer-additional-bills',
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
            push.error(err.response?.data?.message || 'Failed to create customer additional bill')
            return null
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // UPDATE
    // =========================================================================

    const updateCustomerAdditionalBill = async (
        id: number,
        payload: UpdateCustomerAdditionalBillPayload
    ): Promise<CustomerAdditionalBill | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<CustomerAdditionalBill>>(
                `/customer-additional-bills/${id}`,
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
            push.error(err.response?.data?.message || 'Failed to update customer additional bill')
            return null
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // PAYMENTS
    // =========================================================================

    const createCustomerAdditionalBillPayment = async (
        billId: number,
        payload: CreateBillPaymentPayload
    ): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.post<ApiResponse<unknown>>(
                `/customer-additional-bills/${billId}/payments`,
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

    const deleteCustomerAdditionalBill = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<null>>(`/customer-additional-bills/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            bills.value = bills.value.filter(b => b.id !== id)
            push.success(res.data.message)
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete customer additional bill')
            return false
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // SORT / SEARCH
    // =========================================================================

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

    const getBillById = (id: number): CustomerAdditionalBill | undefined =>
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

        fetchCustomerAdditionalBills,
        createCustomerAdditionalBill,
        updateCustomerAdditionalBill,
        createCustomerAdditionalBillPayment,
        deleteCustomerAdditionalBill,

        setSort,
        setSearchQuery,
        clearSearch,

        getBillById,
    }
})