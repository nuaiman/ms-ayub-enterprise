// src/stores/deliveries.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    Delivery,
    DeliveryItem,
    CreateDeliveryPayload,
    UpdateDeliveryPayload,
    CreateDeliveryItemPayload,
    UpdateDeliveryItemPayload,
    DeliveryDetailResponse,
    DeliverySortField,
    SortDirection,
} from '@/types/delivery'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useDeliveriesStore = defineStore('deliveries', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    // ============= STATE =============
    const deliveries = ref<Delivery[]>([])
    const searchQuery = ref('')
    const sortField = ref<DeliverySortField>('delivery_date')
    const sortDirection = ref<SortDirection>('desc')

    // Per-delivery item cache, keyed by delivery id.
    const itemsByDelivery = ref<Record<number, DeliveryItem[]>>({})
    const itemsLoading = ref<Record<number, boolean>>({})

    // ============= COMPUTED =============
    const filteredDeliveries = computed(() => {
        let result = [...deliveries.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(d =>
                String(d.id).includes(query) ||
                (d.receiver_name && d.receiver_name.toLowerCase().includes(query)) ||
                (d.receiver_phone && d.receiver_phone.toLowerCase().includes(query)) ||
                (d.from_location && d.from_location.toLowerCase().includes(query)) ||
                (d.to_location && d.to_location.toLowerCase().includes(query)) ||
                (d.notes && d.notes.toLowerCase().includes(query))
            )
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'customer_id':
                    comparison = (a.customer_id ?? 0) - (b.customer_id ?? 0)
                    break
                case 'delivery_date':
                    comparison = new Date(a.delivery_date).getTime() - new Date(b.delivery_date).getTime()
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

    const totalDeliveries = computed(() => deliveries.value.length)

    const totalItemCount = computed(() =>
        Object.values(itemsByDelivery.value).reduce((sum, arr) => sum + arr.length, 0)
    )

    // ============= FETCH =============

    const fetchDeliveries = async (params?: { customer_id?: number }) => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<Delivery[]>>('/deliveries', { params })
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            deliveries.value = res.data.data
            // Cache is now stale for anything no longer present.
            const validIds = new Set(deliveries.value.map(d => d.id))
            for (const id of Object.keys(itemsByDelivery.value)) {
                if (!validIds.has(Number(id))) delete itemsByDelivery.value[Number(id)]
            }
            return deliveries.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch deliveries')
            return []
        } finally {
            destroyLoader()
        }
    }

    const fetchDeliveryDetail = async (id: number): Promise<DeliveryDetailResponse | null> => {
        try {
            const res = await api.get<ApiResponse<DeliveryDetailResponse>>(`/deliveries/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch delivery')
            return null
        }
    }

    // Low-level fetch (does NOT touch cache).
    const fetchDeliveryItems = async (deliveryId: number): Promise<DeliveryItem[]> => {
        try {
            const res = await api.get<ApiResponse<DeliveryItem[]>>(`/deliveries/${deliveryId}/items`)
            if (!res.data.success) return []
            return res.data.data
        } catch {
            return []
        }
    }

    // Cache-aware loader. If already cached, resolves immediately.
    const loadItemsFor = async (deliveryId: number, force = false): Promise<DeliveryItem[]> => {
        if (!force && itemsByDelivery.value[deliveryId] !== undefined) {
            return itemsByDelivery.value[deliveryId] ?? []
        }
        itemsLoading.value[deliveryId] = true
        try {
            const items = await fetchDeliveryItems(deliveryId)
            itemsByDelivery.value[deliveryId] = items
            return items
        } finally {
            itemsLoading.value[deliveryId] = false
        }
    }

    const getItemsFor = (deliveryId: number): DeliveryItem[] =>
        itemsByDelivery.value[deliveryId] ?? []

    const invalidateItems = (deliveryId: number) => {
        delete itemsByDelivery.value[deliveryId]
    }

    // ============= CREATE =============

    const createDelivery = async (payload: CreateDeliveryPayload): Promise<Delivery | null> => {
        displayLoader()
        try {
            const res = await api.post<ApiResponse<Delivery>>('/deliveries', payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            deliveries.value.unshift(res.data.data)
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to create delivery')
            return null
        } finally {
            destroyLoader()
        }
    }

    /**
     * Creates a single delivery item. On failure the full backend
     * response is logged to the console so the caller can actually
     * debug *why* it was rejected.
     */
    const createDeliveryItem = async (
        deliveryId: number,
        payload: CreateDeliveryItemPayload
    ): Promise<DeliveryItem | null> => {
        try {
            const res = await api.post<ApiResponse<DeliveryItem>>(
                `/deliveries/${deliveryId}/items`,
                payload
            )

            if (!res.data.success) {
                // Backend said success:false — surface the message and log everything.
                console.error('[createDeliveryItem] API success=false', {
                    deliveryId,
                    payload,
                    response: res.data,
                })
                push.error(res.data.message || 'Failed to create delivery item')
                return null
            }

            // Append to cache if we already have it.
            const existing = itemsByDelivery.value[deliveryId]
            if (existing !== undefined) {
                existing.push(res.data.data)
            }

            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            console.error('[createDeliveryItem] Request threw', {
                deliveryId,
                payload,
                status: err.response?.status,
                responseData: err.response?.data,
                error: err,
            })
            push.error(
                err.response?.data?.message ||
                `Failed to create delivery item (HTTP ${err.response?.status ?? 'network error'})`
            )
            return null
        }
    }

    // ============= UPDATE =============

    const updateDelivery = async (
        id: number,
        payload: UpdateDeliveryPayload
    ): Promise<Delivery | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<Delivery>>(`/deliveries/${id}`, payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            const index = deliveries.value.findIndex(d => d.id === id)
            if (index !== -1) deliveries.value[index] = res.data.data
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to update delivery')
            return null
        } finally {
            destroyLoader()
        }
    }

    const updateDeliveryItem = async (
        itemId: number,
        payload: UpdateDeliveryItemPayload
    ): Promise<DeliveryItem | null> => {
        try {
            const res = await api.patch<ApiResponse<DeliveryItem>>(
                `/delivery-items/${itemId}`,
                payload
            )
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            // Update in whatever cache list contains this item.
            for (const key of Object.keys(itemsByDelivery.value)) {
                const list = itemsByDelivery.value[Number(key)]
                if (!list) continue
                const idx = list.findIndex(i => i.id === itemId)
                if (idx !== -1) {
                    list[idx] = res.data.data
                    break
                }
            }
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to update delivery item')
            return null
        }
    }

    // Refresh the master record after an image change so `image_url` is current.
    const refreshDelivery = async (id: number): Promise<Delivery | null> => {
        try {
            const res = await api.get<ApiResponse<DeliveryDetailResponse>>(`/deliveries/${id}`)
            if (!res.data.success) return null
            const index = deliveries.value.findIndex(d => d.id === id)
            if (index !== -1) deliveries.value[index] = res.data.data.delivery
            return res.data.data.delivery
        } catch {
            return null
        }
    }

    // ============= DELETE =============

    const deleteDelivery = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<null>>(`/deliveries/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            deliveries.value = deliveries.value.filter(d => d.id !== id)
            delete itemsByDelivery.value[id]
            push.success(res.data.message)
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete delivery')
            return false
        } finally {
            destroyLoader()
        }
    }

    const deleteDeliveryItem = async (itemId: number): Promise<boolean> => {
        try {
            const res = await api.delete<ApiResponse<null>>(`/delivery-items/${itemId}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            for (const key of Object.keys(itemsByDelivery.value)) {
                const list = itemsByDelivery.value[Number(key)]
                if (!list) continue
                const idx = list.findIndex(i => i.id === itemId)
                if (idx !== -1) {
                    list.splice(idx, 1)
                    break
                }
            }
            push.success(res.data.message)
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete delivery item')
            return false
        }
    }

    // ============= SORT / SEARCH =============

    const setSort = (field: DeliverySortField) => {
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

    const getDeliveryById = (id: number): Delivery | undefined =>
        deliveries.value.find(d => d.id === id)

    return {
        // State
        deliveries,
        searchQuery,
        sortField,
        sortDirection,
        itemsByDelivery,
        itemsLoading,

        // Computed
        filteredDeliveries,
        totalDeliveries,
        totalItemCount,

        // Fetch
        fetchDeliveries,
        fetchDeliveryDetail,
        fetchDeliveryItems,
        loadItemsFor,
        getItemsFor,
        invalidateItems,
        refreshDelivery,

        // CRUD
        createDelivery,
        createDeliveryItem,
        updateDelivery,
        updateDeliveryItem,
        deleteDelivery,
        deleteDeliveryItem,

        // Sort / Search
        setSort,
        setSearchQuery,
        clearSearch,

        // Utilities
        getDeliveryById,
    }
})