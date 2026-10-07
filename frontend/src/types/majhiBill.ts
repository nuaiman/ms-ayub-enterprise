// src/types/majhiBill.ts

export type MajhiBillType = 'weight' | 'quantity' | 'job'

export interface MajhiBill {
    id: number
    user_id: number
    majhi_id: number
    store_id: number | null
    delivery_item_id: number | null
    bill_type: MajhiBillType
    rate: number
    weight_at_billing: number
    quantity_at_billing: number
    weight_unit_at_billing: string | null
    quantity_unit_at_billing: string | null
    total_amount: number
    total_paid: number
    total_paid_through: string | null
    remaining: number
    created_at: string
    updated_at: string
}

export interface CreateMajhiBillPayload {
    majhi_id: number
    store_id?: number | null
    delivery_item_id?: number | null
    bill_type: MajhiBillType
    rate: number
}

export interface UpdateMajhiBillPayload {
    bill_type?: MajhiBillType
    rate?: number
    weight_at_billing?: number
    quantity_at_billing?: number
}

export interface CreateMajhiBillPaymentPayload {
    amount: number
    payment_date?: string
    payment_method?: 'cash' | 'bank_transfer' | 'check' | 'mobile_banking'
    reference_number?: string | null
    notes?: string | null
}

export type MajhiBillSortField =
    | 'bill_type'
    | 'rate'
    | 'total_paid'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'