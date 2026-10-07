// src/types/customerStoreBill.ts

export type CustomerStoreBillType = 'weight' | 'quantity'

export interface CustomerStoreBill {
    id: number
    user_id: number
    customer_id: number
    store_id: number
    month_year: string
    bill_type: CustomerStoreBillType
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

export interface CreateCustomerStoreBillPayload {
    customer_id: number
    store_id: number
    month_year?: string   // optional — backend derives from store.start_date if omitted
    bill_type: CustomerStoreBillType
    rate: number
}

export interface UpdateCustomerStoreBillPayload {
    bill_type?: CustomerStoreBillType
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

export type CustomerStoreBillSortField =
    | 'store_id'
    | 'customer_id'
    | 'bill_type'
    | 'rate'
    | 'total_paid'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'