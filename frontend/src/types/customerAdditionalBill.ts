// src/types/customerAdditionalBill.ts

export interface AdditionalCharge {
    id: number
    user_id: number
    customer_id: number
    amount: number
    description: string
    customer_total_paid: number
    customer_total_paid_through: string | null
    created_at: string
    updated_at: string
}

export interface CreateAdditionalChargePayload {
    customer_id: number
    amount: number
    description: string
}

export interface UpdateAdditionalChargePayload {
    amount?: number
    description?: string
}

export interface UpdateAdditionalChargePaymentPayload {
    customer_total_paid: number
    customer_total_paid_through?: string | null
}

// View model for the bills page — enriched with customer name.
export interface CustomerAdditionalBill {
    id: number
    customer_id: number
    customer_name: string
    amount: number
    description: string
    paid_amount: number
    outstanding: number
    status: 'unpaid' | 'paid'
    payment_date: string | null
    created_at: string
    updated_at: string
}

export type CustomerAdditionalBillSortField =
    | 'customer_name'
    | 'amount'
    | 'description'
    | 'status'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'