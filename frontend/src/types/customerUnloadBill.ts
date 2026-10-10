// src/types/customerUnloadBill.ts

export type CustomerUnloadBillType = 'weight' | 'quantity'

export interface CustomerUnloadBill {
    id: number
    user_id: number
    customer_id: number
    store_id: number
    bill_type: CustomerUnloadBillType
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

export interface CreateCustomerUnloadBillPayload {
    customer_id: number
    store_id: number
    bill_type: CustomerUnloadBillType
    rate: number
}

export interface UpdateCustomerUnloadBillPayload {
    bill_type?: CustomerUnloadBillType
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

export type CustomerUnloadBillSortField =
    | 'customer_id'
    | 'store_id'
    | 'bill_type'
    | 'rate'
    | 'total_paid'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'