// src/types/customerAdditionalBill.ts

export interface CustomerAdditionalBill {
    id: number
    user_id: number
    customer_id: number
    amount: number
    description: string
    total_paid: number
    total_paid_through: string | null
    remaining: number
    created_at: string
    updated_at: string
}

export interface CreateCustomerAdditionalBillPayload {
    customer_id: number
    amount: number
    description: string
}

export interface UpdateCustomerAdditionalBillPayload {
    amount?: number
    description?: string
}

export interface CreateBillPaymentPayload {
    amount: number
    payment_date?: string
    payment_method?: 'cash' | 'bank_transfer' | 'check' | 'mobile_banking'
    reference_number?: string | null
    notes?: string | null
}

export type CustomerAdditionalBillSortField =
    | 'customer_id'
    | 'amount'
    | 'total_paid'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'