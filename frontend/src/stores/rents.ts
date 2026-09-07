import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
    Rent,
    CreateRentPayload,
    UpdateRentPayload,
    MarkRentAsPaidPayload,
    RentFilters,
    RentStatus,
} from '@/types/rent'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'
import { useGodownsStore } from './godowns'
import { useUsersStore } from './users'

export const useRentsStore = defineStore('rents', () => {
    const { displayLoader, destroyLoader } = useGlobalLoader()

    const rents = ref<Rent[]>([])
    const searchQuery = ref('')

    const filteredRents = computed(() => {
        let result = [...rents.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            const godownsStore = useGodownsStore()
            const usersStore = useUsersStore()
            result = result.filter(
                (r) =>
                    String(r.amount).includes(query) ||
                    r.month_year.includes(query) ||
                    r.status.toLowerCase().includes(query) ||
                    (r.notes && r.notes.toLowerCase().includes(query)) ||
                    godownsStore.getGodownName(r.godown_id).toLowerCase().includes(query) ||
                    usersStore.getUserName(r.user_id).toLowerCase().includes(query)
            )
        }

        return result
    })

    const statusCounts = computed(() => {
        const counts: Record<RentStatus, number> = {
            draft: 0,
            paid: 0,
            cancelled: 0,
        }
        rents.value.forEach((r) => {
            if (counts[r.status] !== undefined) {
                counts[r.status]++
            }
        })
        return counts
    })

    const totalRents = computed(() => rents.value.length)
    const draftCount = computed(() => rents.value.filter(r => r.status === 'draft').length)
    const totalOutstanding = computed(() => {
        return rents.value
            .filter((r) => r.status === 'draft')
            .reduce((sum, r) => sum + r.amount, 0)
    })

    const fetchRents = async (filters?: RentFilters) => {
        displayLoader()
        try {
            const params: Record<string, string> = {}
            if (filters?.godown_id) params.godown_id = filters.godown_id
            if (filters?.month_year) params.month_year = filters.month_year
            if (filters?.status) params.status = filters.status

            const res = await api.get<ApiResponse<Rent[]>>('/rents', { params })
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            rents.value = res.data.data
            return rents.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch rents')
            return []
        } finally {
            destroyLoader()
        }
    }

    const fetchCurrentMonthRents = async () => {
        displayLoader()
        try {
            const res = await api.get<ApiResponse<Rent[]>>('/rents/current-month')
            if (!res.data.success) {
                push.error(res.data.message)
                return []
            }
            rents.value = res.data.data
            return rents.value
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to fetch current month rents')
            return []
        } finally {
            destroyLoader()
        }
    }

    const createRent = async (payload: CreateRentPayload): Promise<Rent | null> => {
        displayLoader()
        try {
            if (!payload.godown_id) {
                push.error('Godown is required')
                return null
            }
            if (!payload.month_year) {
                push.error('Month is required')
                return null
            }
            if (payload.amount < 0) {
                push.error('Amount cannot be negative')
                return null
            }

            const res = await api.post<ApiResponse<Rent>>('/rents', payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            rents.value.push(res.data.data)
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to create rent')
            return null
        } finally {
            destroyLoader()
        }
    }

    const updateRent = async (id: number, payload: UpdateRentPayload): Promise<Rent | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<Rent>>(`/rents/${id}`, payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            const index = rents.value.findIndex((r) => r.id === id)
            if (index !== -1) {
                rents.value[index] = res.data.data
            }
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to update rent')
            return null
        } finally {
            destroyLoader()
        }
    }

    const markRentAsPaid = async (id: number, payload: MarkRentAsPaidPayload): Promise<Rent | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<Rent>>(`/rents/${id}/pay`, payload)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            const index = rents.value.findIndex((r) => r.id === id)
            if (index !== -1) {
                rents.value[index] = res.data.data
            }
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to mark rent as paid')
            return null
        } finally {
            destroyLoader()
        }
    }

    const cancelRent = async (id: number): Promise<Rent | null> => {
        displayLoader()
        try {
            const res = await api.patch<ApiResponse<Rent>>(`/rents/${id}/cancel`)
            if (!res.data.success) {
                push.error(res.data.message)
                return null
            }
            const index = rents.value.findIndex((r) => r.id === id)
            if (index !== -1) {
                rents.value[index] = res.data.data
            }
            push.success(res.data.message)
            return res.data.data
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to cancel rent')
            return null
        } finally {
            destroyLoader()
        }
    }

    const deleteRent = async (id: number): Promise<boolean> => {
        displayLoader()
        try {
            const res = await api.delete<ApiResponse<null>>(`/rents/${id}`)
            if (!res.data.success) {
                push.error(res.data.message)
                return false
            }
            rents.value = rents.value.filter((r) => r.id !== id)
            push.success(res.data.message)
            return true
        } catch (error) {
            const err = error as AxiosError<ApiResponse<null>>
            push.error(err.response?.data?.message || 'Failed to delete rent')
            return false
        } finally {
            destroyLoader()
        }
    }

    const getRentById = (id: number): Rent | undefined => {
        return rents.value.find((r) => r.id === id)
    }

    return {
        rents,
        searchQuery,
        filteredRents,
        statusCounts,
        totalRents,
        draftCount,
        totalOutstanding,
        fetchRents,
        fetchCurrentMonthRents,
        createRent,
        updateRent,
        markRentAsPaid,
        cancelRent,
        deleteRent,
        getRentById,
    }
})