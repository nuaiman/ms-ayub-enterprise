// src/stores/storeTransfers.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    StoreTransfer,
    CreateStoreTransferPayload,
    CreateStoreTransferResponse,
    StoreTransferSortField,
    SortDirection,
} from '@/types/storeTransfer'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useStoreTransfersStore = defineStore('storeTransfers', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    const transfers = ref<StoreTransfer[]>([])
    const searchQuery = ref('')
    const sortField = ref<StoreTransferSortField>('transferred_at')
    const sortDirection = ref<SortDirection>('desc')

    const filteredTransfers = computed(() => {
        let result = [...transfers.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(t =>
                String(t.store_id).includes(query) ||
                String(t.from_godown_id).includes(query) ||
                String(t.to_godown_id).includes(query) ||
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
                case 'store_id':
                    comparison = a.store_id - b.store_id
                    break
                case 'user_id':
                    comparison = a.user_id - b.user_id
                    break
                case 'from_godown_id':
                    comparison = a.from_godown_id - b.from_godown_id
                    break
                case 'to_godown_id':
                    comparison = a.to_godown_id - b.to_godown_id
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
            const res = await api.get<ApiResponse<StoreTransfer[]>>(
                '/store-transfers'
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
                err.response?.data?.message || 'Failed to fetch store transfers'
            )
            return []
        } finally {
            destroyLoader()
        }
    }

    const fetchTransfersByStore = async (storeId: number) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<StoreTransfer[]>>(
                `/stores/${storeId}/transfers`
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
                err.response?.data?.message || 'Failed to fetch store transfers'
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
        storeId: number,
        payload: CreateStoreTransferPayload
    ): Promise<CreateStoreTransferResponse | null> => {
        displayLoader()
        try {
            if (!payload.to_godown_id) {
                push.error('Target godown is required')
                return null
            }

            const res = await api.post<ApiResponse<CreateStoreTransferResponse>>(
                `/stores/${storeId}/transfer`,
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }

            transfers.value.unshift(res.data.data.transfer)
            push.success(res.data.message || 'Store transferred')
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to transfer store'
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
                `/store-transfers/${id}`
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

    const setSort = (field: StoreTransferSortField) => {
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

    const getTransferById = (id: number): StoreTransfer | undefined =>
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
        fetchTransfersByStore,

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