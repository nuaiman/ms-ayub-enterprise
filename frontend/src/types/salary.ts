// src/types/salary.ts

export type SalaryStatus = 'draft' | 'paid' | 'cancelled'

export type PaymentMethod = 'cash' | 'bank_transfer' | 'check' | 'mobile_banking'

export interface Salary {
  id: number
  user_id: number          // Created by
  employee_id: number      // Employee being paid
  month_year: string       // Format: YYYY-MM
  bonus: number
  deductions: number
  status: SalaryStatus
  payment_date: string | null
  payment_method: PaymentMethod | null
  reference_number: string | null
  notes: string | null
  created_at: string
  updated_at: string
  // These come from the API response (calculated on the fly)
  basic_salary?: number
  total_salary?: number
}

export interface CreateSalaryPayload {
  employee_id: number
  month_year: string
  bonus: number
  deductions: number
  notes?: string | null
}

export interface UpdateSalaryPayload {
  bonus?: number
  deductions?: number
  notes?: string | null
}

export interface MarkSalaryAsPaidPayload {
  payment_date?: string
  payment_method?: PaymentMethod
  reference_number?: string | null
}

export interface SalaryFilters {
  employee_id?: string
  month_year?: string
  status?: SalaryStatus
}

export type SalarySortField = 'employee_id' | 'month_year' | 'status' | 'created_at'
export type SortDirection = 'asc' | 'desc'