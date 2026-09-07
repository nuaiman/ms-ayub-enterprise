// src/types/customer.ts

export interface Customer {
  id: number
  company_name: string | null
  contact_person: string | null
  phone: string
  email: string | null
  address: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export interface CreateCustomerPayload {
  company_name?: string | null
  contact_person?: string | null
  phone: string
  email?: string | null
  address?: string | null
  notes?: string | null
}

export interface UpdateCustomerPayload {
  company_name?: string | null
  contact_person?: string | null
  phone?: string
  email?: string | null
  address?: string | null
  notes?: string | null
}

export type CustomerSortField = 'company_name' | 'contact_person' | 'phone' | 'email' | 'address' | 'created_at' // <-- Added address to sort
export type SortDirection = 'asc' | 'desc'