// src/types/godownBill.ts

export type GodownBillType = 'weight' | 'quantity' | 'fixed'

export interface GodownBill {
    id: number
    user_id: number
    godown_id: number
    store_id: number
    month_year: string
    bill_type: GodownBillType
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

export interface CreateGodownBillPayload {
    godown_id: number
    store_id: number
    month_year?: string   // optional — backend derives from store.start_date if omitted
    bill_type: GodownBillType
    rate: number
}

export interface UpdateGodownBillPayload {
    bill_type?: GodownBillType
    rate?: number
    weight_at_billing?: number
    quantity_at_billing?: number
}

export interface CreateBillPaymentPayload {
    amount: number
    payment_date?: string
    payment_method?: 'cash' | 'bank_transfer' | 'check' | 'mobile_banking'
    reference_number?: string | null
    notes?: string | null
}

export type GodownBillSortField =
    | 'store_id'
    | 'godown_id'
    | 'bill_type'
    | 'rate'
    | 'total_paid'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'