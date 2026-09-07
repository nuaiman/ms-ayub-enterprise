import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/axios'
import type { Customer, CreateCustomerPayload, UpdateCustomerPayload } from '@/types/customer'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useCustomersStore = defineStore('customers', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()

  const customers = ref<Customer[]>([])
  const searchQuery = ref('')

  const fetchCustomers = async () => {
    displayLoader()
    try {
      const res = await api.get<ApiResponse<Customer[]>>('/customers')
      if (!res.data.success) {
        push.error(res.data.message)
        return []
      }
      customers.value = res.data.data
      return customers.value
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to fetch customers')
      return []
    } finally {
      destroyLoader()
    }
  }

  const createCustomer = async (payload: CreateCustomerPayload): Promise<Customer | null> => {
    displayLoader()
    try {
      if (!payload.company_name && !payload.contact_person) {
        push.error('Either company name or contact person is required')
        return null
      }

      const res = await api.post<ApiResponse<Customer>>('/customers', payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      customers.value.push(res.data.data)
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to create customer')
      return null
    } finally {
      destroyLoader()
    }
  }

  const updateCustomer = async (id: number, payload: UpdateCustomerPayload): Promise<Customer | null> => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<Customer>>(`/customers/${id}`, payload)
      if (!res.data.success) {
        push.error(res.data.message)
        return null
      }
      const index = customers.value.findIndex(c => c.id === id)
      if (index !== -1) {
        customers.value[index] = res.data.data
      }
      push.success(res.data.message)
      return res.data.data
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to update customer')
      return null
    } finally {
      destroyLoader()
    }
  }

  const deleteCustomer = async (id: number): Promise<boolean> => {
    displayLoader()
    try {
      const res = await api.delete<ApiResponse<null>>(`/customers/${id}`)
      if (!res.data.success) {
        push.error(res.data.message)
        return false
      }
      customers.value = customers.value.filter(c => c.id !== id)
      push.success(res.data.message)
      return true
    } catch (error) {
      const err = error as AxiosError<ApiResponse<null>>
      push.error(err.response?.data?.message || 'Failed to delete customer')
      return false
    } finally {
      destroyLoader()
    }
  }

  const getCustomerName = (id: number): string => {
    const customer = customers.value.find(c => c.id === id)
    if (!customer) return `Customer #${id}`
    return customer.company_name || customer.contact_person || `Customer #${id}`
  }

  const getCustomerById = (id: number): Customer | undefined => {
    return customers.value.find(c => c.id === id)
  }

  return {
    customers,
    searchQuery,
    fetchCustomers,
    createCustomer,
    updateCustomer,
    deleteCustomer,
    getCustomerName,
    getCustomerById,
  }
})