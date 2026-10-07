// src/stores/invoices.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    Invoice,
    InvoiceDetail,
    UnbilledBill,
    CreateInvoicePayload,
    UpdateInvoicePayload,
    EntityType,
    InvoiceSortField,
    SortDirection,
} from '@/types/invoice'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useInvoicesStore = defineStore('invoices', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    // ============= STATE =============
    const invoices = ref<Invoice[]>([])
    const searchQuery = ref('')
    const sortField = ref<InvoiceSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    // ============= COMPUTED =============
    const filteredInvoices = computed(() => {
        let result = [...invoices.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(
                (i) =>
                    i.entity_type.toLowerCase().includes(query) ||
                    String(i.entity_id).includes(query) ||
                    String(i.total).includes(query) ||
                    String(i.subtotal).includes(query) ||
                    (i.notes && i.notes.toLowerCase().includes(query))
            )
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'entity_id':
                    comparison = a.entity_id - b.entity_id
                    break
                case 'subtotal':
                    comparison = a.subtotal - b.subtotal
                    break
                case 'discount_amount':
                    comparison = a.discount_amount - b.discount_amount
                    break
                case 'total':
                    comparison = a.total - b.total
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

    const totalInvoices = computed(() => invoices.value.length)
    const totalSubtotal = computed(() => invoices.value.reduce((sum, i) => sum + i.subtotal, 0))
    const totalDiscount = computed(() => invoices.value.reduce((sum, i) => sum + i.discount_amount, 0))
    const totalAmount = computed(() => invoices.value.reduce((sum, i) => sum + i.total, 0))

    // ============= FETCH =============
    const fetchInvoices = async (params?: {
        entity_type?: EntityType
        entity_id?: number
    }) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<Invoice[]>>('/invoices', { params })
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            invoices.value = res.data.data
            return invoices.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch invoices')
            return []
        } finally {
            destroyLoader()
        }
    }

    const fetchInvoiceDetail = async (id: number): Promise<InvoiceDetail | null> => {
        try {
            const res = await api.get<ApiResponse<InvoiceDetail>>(`/invoices/${id}`)
            if (!res.data.success) return null
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch invoice')
            return null
        }
    }

    const fetchUnbilled = async (
        entityType: EntityType,
        entityID: number
    ): Promise<UnbilledBill[]> => {
        try {
            const res = await api.get<ApiResponse<UnbilledBill[]>>('/invoices/unbilled', {
                params: { entity_type: entityType, entity_id: entityID },
            })
            if (!res.data.success) return []
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch unbilled bills')
            return []
        }
    }

    // ============= CREATE =============
    const createInvoice = async (payload: CreateInvoicePayload): Promise<InvoiceDetail | null> => {
        displayLoader()
        try {
            if (!payload.entity_type) {
                push.error('Entity type is required')
                return null
            }
            if (!payload.entity_id) {
                push.error('Entity is required')
                return null
            }
            if (!payload.bills || payload.bills.length === 0) {
                push.error('Select at least one bill')
                return null
            }

            const res = await api.post<ApiResponse<InvoiceDetail>>('/invoices', payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }

            invoices.value.unshift(res.data.data)
            push.success(res.data.message || 'Invoice created')
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to create invoice')
            return null
        } finally {
            destroyLoader()
        }
    }

    // ============= UPDATE =============
    const updateInvoice = async (
        id: number,
        payload: UpdateInvoicePayload
    ): Promise<InvoiceDetail | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<InvoiceDetail>>(`/invoices/${id}`, payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }

            const index = invoices.value.findIndex((i) => i.id === id)
            if (index !== -1) invoices.value[index] = res.data.data

            push.success(res.data.message || 'Invoice updated')
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to update invoice')
            return null
        } finally {
            destroyLoader()
        }
    }

    // ============= DELETE =============
    const deleteInvoice = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<null>>(`/invoices/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            invoices.value = invoices.value.filter((i) => i.id !== id)
            push.success(res.data.message || 'Invoice deleted')
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete invoice')
            return false
        } finally {
            destroyLoader()
        }
    }

    // ============= SORT / SEARCH =============
    const setSort = (field: InvoiceSortField) => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
        } else {
            sortField.value = field
            sortDirection.value = 'desc'
        }
    }

    const setSearchQuery = (query: string) => {
        searchQuery.value = query
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    // ============= UTILITIES =============
    const getInvoiceById = (id: number): Invoice | undefined =>
        invoices.value.find((i) => i.id === id)

    return {
        // State
        invoices,
        searchQuery,
        sortField,
        sortDirection,

        // Computed
        filteredInvoices,
        totalInvoices,
        totalSubtotal,
        totalDiscount,
        totalAmount,

        // Fetch
        fetchInvoices,
        fetchInvoiceDetail,
        fetchUnbilled,

        // CRUD
        createInvoice,
        updateInvoice,
        deleteInvoice,

        // Sort / Search
        setSort,
        setSearchQuery,
        clearSearch,

        // Utilities
        getInvoiceById,
    }
})