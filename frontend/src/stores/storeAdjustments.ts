// src/stores/storeAdjustments.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    StoreAdjustment,
    CreateStoreAdjustmentPayload,
    CreateStoreAdjustmentResponse,
    StoreAdjustmentSortField,
    SortDirection,
} from '@/types/storeAdjustment'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useStoreAdjustmentsStore = defineStore('storeAdjustments', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    const adjustments = ref<StoreAdjustment[]>([])
    const searchQuery = ref('')
    const sortField = ref<StoreAdjustmentSortField>('adjusted_at')
    const sortDirection = ref<SortDirection>('desc')

    const filteredAdjustments = computed(() => {
        let result = [...adjustments.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(a =>
                String(a.store_id).includes(query) ||
                String(a.user_id).includes(query) ||
                a.adjustment_type.toLowerCase().includes(query) ||
                (a.reason && a.reason.toLowerCase().includes(query)) ||
                (a.notes && a.notes.toLowerCase().includes(query))
            )
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'adjusted_at':
                    comparison =
                        new Date(a.adjusted_at).getTime() -
                        new Date(b.adjusted_at).getTime()
                    break
                case 'store_id':
                    comparison = a.store_id - b.store_id
                    break
                case 'user_id':
                    comparison = a.user_id - b.user_id
                    break
                case 'weight_delta':
                    comparison = a.weight_delta - b.weight_delta
                    break
                case 'quantity_delta':
                    comparison = a.quantity_delta - b.quantity_delta
                    break
                default:
                    comparison = 0
            }
            return sortDirection.value === 'desc' ? -comparison : comparison
        })

        return result
    })

    const totalAdjustments = computed(() => adjustments.value.length)

    // =========================================================================
    // FETCH
    // =========================================================================

    const fetchAllAdjustments = async () => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<StoreAdjustment[]>>(
                '/store-adjustments'
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            adjustments.value = res.data.data
            return adjustments.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to fetch store adjustments'
            )
            return []
        } finally {
            destroyLoader()
        }
    }

    const fetchAdjustmentsByStore = async (storeId: number) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<StoreAdjustment[]>>(
                `/stores/${storeId}/adjustments`
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            adjustments.value = res.data.data
            return adjustments.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to fetch store adjustments'
            )
            return []
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // CREATE
    // =========================================================================

    const createAdjustment = async (
        storeId: number,
        payload: CreateStoreAdjustmentPayload
    ): Promise<CreateStoreAdjustmentResponse | null> => {
        displayLoader()
        try {
            if (!payload.input_weight && !payload.input_quantity) {
                push.error(
                    'At least one of weight or quantity must be non-zero'
                )
                return null
            }

            const res = await api.post<ApiResponse<CreateStoreAdjustmentResponse>>(
                `/stores/${storeId}/adjustments`,
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }

            adjustments.value.unshift(res.data.data.adjustment)
            push.success(res.data.message || 'Adjustment recorded')
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to create adjustment'
            )
            return null
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // DELETE
    // =========================================================================

    const deleteAdjustment = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<unknown>>(
                `/store-adjustments/${id}`
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            adjustments.value = adjustments.value.filter(a => a.id !== id)
            push.success(res.data.message || 'Adjustment deleted')
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(
                err.response?.data?.message || 'Failed to delete adjustment'
            )
            return false
        } finally {
            destroyLoader()
        }
    }

    // =========================================================================
    // SORT / SEARCH
    // =========================================================================

    const setSort = (field: StoreAdjustmentSortField) => {
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

    const getAdjustmentById = (id: number): StoreAdjustment | undefined =>
        adjustments.value.find(a => a.id === id)

    return {
        // State
        adjustments,
        searchQuery,
        sortField,
        sortDirection,

        // Computed
        filteredAdjustments,
        totalAdjustments,

        // Fetch
        fetchAllAdjustments,
        fetchAdjustmentsByStore,

        // CRUD
        createAdjustment,
        deleteAdjustment,

        // Sort / Search
        setSort,
        setSearchQuery,
        clearSearch,

        // Utilities
        getAdjustmentById,
    }
})