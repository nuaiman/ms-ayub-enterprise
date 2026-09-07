import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/axios'
import type {
  Salary,
  CreateSalaryPayload,
  UpdateSalaryPayload,
  MarkSalaryAsPaidPayload,
  SalaryFilters,
  SalaryStatus,
} from '@/types/salary'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'
import { useUsersStore } from './users'

export const useSalariesStore = defineStore('salaries', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()

  const salaries = ref<Salary[]>([])
  const searchQuery = ref('')

  const filteredSalaries = computed(() => {
    let result = [...salaries.value]

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      const usersStore = useUsersStore()
      result = result.filter((s) => {
        const employeeName = usersStore.getUserName(s.employee_id).toLowerCase()
        return (
          employeeName.includes(query) ||
          s.month_year.includes(query) ||
          s.status.toLowerCase().includes(query) ||
          (s.notes && s.notes.toLowerCase().includes(query))
        )
      })
    }

    return result
  })

  const statusCounts = computed(() => {
    const counts: Record<SalaryStatus, number> = {
      draft: 0,
      paid: 0,
      cancelled: 0,
    }
    salaries.value.forEach((s) => {
      if (counts[s.status] !== undefined) {
        counts[s.status]++
      }
    })
    return counts
  })

  const totalSalaries = computed(() => salaries.value.length)
  const draftCount = computed(() => salaries.value.filter(s => s.status === 'draft').length)

  const getEmployeeBasicSalary = (employeeId: number): number => {
    const usersStore = useUsersStore()
    const user = usersStore.users.find((u) => u.id === employeeId)
    return user?.monthly_salary || 0
  }

  const calculateTotalSalary = (salary: Salary): number => {
    const basicSalary = getEmployeeBasicSalary(salary.employee_id)
    return basicSalary + salary.bonus - salary.deductions
  }

  const fetchSalaries = async (filters?: SalaryFilters) => {
    displayLoader()
    try {
      const params: Record<string, string> = {}
      if (filters?.employee_id) params.employee_id = filters.employee_id
      if (filters?.month_year) params.month_year = filters.month_year
      if (filters?.status) params.status = filters.status

      const res = await api.get<ApiResponse<Salary[]>>('/salaries', { params })
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      salaries.value = res.data.data
      return salaries.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch salaries')
      return []
    } finally {
      destroyLoader()
    }
  }

  const fetchSalariesByEmployee = async (employeeId: number) => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Salary[]>>(`/salaries/employee/${employeeId}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      salaries.value = res.data.data
      return salaries.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch employee salaries')
      return []
    } finally {
      destroyLoader()
    }
  }

  const fetchSalariesByMonth = async (monthYear: string) => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Salary[]>>(`/salaries/month/${monthYear}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      salaries.value = res.data.data
      return salaries.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch monthly salaries')
      return []
    } finally {
      destroyLoader()
    }
  }

  const createSalary = async (payload: CreateSalaryPayload): Promise<Salary | null> => {
    displayLoader()
    try {
      if (!payload.employee_id) {
        push.error('Employee is required')
        return null
      }
      if (!payload.month_year) {
        push.error('Month is required')
        return null
      }
      if (payload.bonus < 0) {
        push.error('Bonus cannot be negative')
        return null
      }
      if (payload.deductions < 0) {
        push.error('Deductions cannot be negative')
        return null
      }

      const res = await api.post<ApiResponse<Salary>>('/salaries', payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      salaries.value.push(res.data.data)
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create salary')
      return null
    } finally {
      destroyLoader()
    }
  }

  const updateSalary = async (id: number, payload: UpdateSalaryPayload): Promise<Salary | null> => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Salary>>(`/salaries/${id}`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      const index = salaries.value.findIndex((s) => s.id === id)
      if (index !== -1) {
        salaries.value[index] = res.data.data
      }
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update salary')
      return null
    } finally {
      destroyLoader()
    }
  }

  const markSalaryAsPaid = async (id: number, payload: MarkSalaryAsPaidPayload): Promise<Salary | null> => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Salary>>(`/salaries/${id}/pay`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      const index = salaries.value.findIndex((s) => s.id === id)
      if (index !== -1) {
        salaries.value[index] = res.data.data
      }
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to mark salary as paid')
      return null
    } finally {
      destroyLoader()
    }
  }

  const cancelSalary = async (id: number): Promise<Salary | null> => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Salary>>(`/salaries/${id}/cancel`)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      const index = salaries.value.findIndex((s) => s.id === id)
      if (index !== -1) {
        salaries.value[index] = res.data.data
      }
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to cancel salary')
      return null
    } finally {
      destroyLoader()
    }
  }

  const deleteSalary = async (id: number): Promise<boolean> => {
    displayLoader()
    try {
      const res = await api.delete<ApiResponse<null>>(`/salaries/${id}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      salaries.value = salaries.value.filter((s) => s.id !== id)
      push.success(res.data.message)
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete salary')
      return false
    } finally {
      destroyLoader()
    }
  }

  const getSalaryById = (id: number): Salary | undefined => {
    return salaries.value.find((s) => s.id === id)
  }

  return {
    salaries,
    searchQuery,
    filteredSalaries,
    statusCounts,
    totalSalaries,
    draftCount,
    fetchSalaries,
    fetchSalariesByEmployee,
    fetchSalariesByMonth,
    createSalary,
    updateSalary,
    markSalaryAsPaid,
    cancelSalary,
    deleteSalary,
    getSalaryById,
    getEmployeeBasicSalary,
    calculateTotalSalary,
  }
})