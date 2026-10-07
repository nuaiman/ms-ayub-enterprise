// src/stores/lotTransfers.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    LotTransfer,
    CreateLotTransferPayload,
    CreateLotTransferResponse,
    LotTransferSortField,
    SortDirection,
} from '@/types/lotTransfer'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useLotTransfersStore = defineStore('lotTransfers', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    const transfers = ref<LotTransfer[]>([])
    const searchQuery = ref('')
    const sortField = ref<LotTransferSortField>('transferred_at')
    const sortDirection = ref<SortDirection>('desc')

    const filteredTransfers = computed(() => {
        let result = [...transfers.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(t =>
                String(t.lot_id).includes(query) ||
                String(t.from_customer_id).includes(query) ||
                String(t.to_customer_id).includes(query) ||
                String(t.user_id).includes(query) ||
                (t.notes && t.notes.toLowerCase().includes(query))
            )
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'transferred_at':
                    comparison =
                        new Date(a.transferred_at).getTime() -
                        new Date(b.transferred_at).getTime()
                    break
                case 'lot_id':
                    comparison = a.lot_id - b.lot_id
                    break
                case 'user_id':
                    comparison = a.user_id - b.user_id
                    break
                case 'from_customer_id':
                    comparison = a.from_customer_id - b.from_customer_id
                    break
                case 'to_customer_id':
                    comparison = a.to_customer_id - b.to_customer_id
                    break
                default:
                    comparison = 0
            }
            return sortDirection.value === 'desc' ? -comparison : comparison
        })

        return result
    })

    const totalTransfers = computed(() => transfers.value.length)

    // =========================================================================
    // FETCH
    // =========================================================================

    const fetchAllTransfers = async () => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<LotTransfer[]>>(
                '/lot-transfers'
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            transfers.value = res.data.data
            return transfers.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to fetch lot transfers'
            )
            return []
        } finally {
            destroyLoader()
        }
    }

    const fetchTransfersByLot = async (lotId: number) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<LotTransfer[]>>(
                `/lots/${lotId}/transfers`
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            transfers.value = res.data.data
            return transfers.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to fetch lot transfers'
            )
            return []
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // CREATE
    // =========================================================================

    const createTransfer = async (
        lotId: number,
        payload: CreateLotTransferPayload
    ): Promise<CreateLotTransferResponse | null> => {
        displayLoader()
        try {
            if (!payload.to_customer_id) {
                push.error('Target customer is required')
                return null
            }

            const res = await api.post<ApiResponse<CreateLotTransferResponse>>(
                `/lots/${lotId}/transfer`,
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }

            transfers.value.unshift(res.data.data.transfer)

            const billsCreated = res.data.data.bills_created
            if (billsCreated > 0) {
                push.success(
                    `Lot transferred. ${billsCreated} customer store bill${billsCreated === 1 ? '' : 's'} created.`
                )
            } else {
                push.success(res.data.message || 'Lot transferred')
            }
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to transfer lot'
            )
            return null
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // DELETE
    // =========================================================================

    const deleteTransfer = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<unknown>>(
                `/lot-transfers/${id}`
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            transfers.value = transfers.value.filter(t => t.id !== id)
            push.success(res.data.message || 'Transfer deleted')
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to delete transfer'
            )
            return false
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // SORT / SEARCH
    // =========================================================================

    const setSort = (field: LotTransferSortField) => {
        if (sortField.value === field) {
            sortDirection.value =
                sortDirection.value === 'asc' ? 'desc' : 'asc'
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

    const getTransferById = (id: number): LotTransfer | undefined =>
        transfers.value.find(t => t.id === id)

    return {
        // State
        transfers,
        searchQuery,
        sortField,
        sortDirection,

        // Computed
        filteredTransfers,
        totalTransfers,

        // Fetch
        fetchAllTransfers,
        fetchTransfersByLot,

        // CRUD
        createTransfer,
        deleteTransfer,

        // Sort / Search
        setSort,
        setSearchQuery,
        clearSearch,

        // Utilities
        getTransferById,
    }
})