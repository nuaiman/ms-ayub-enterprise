// src/types/invoice.ts

export type EntityType = 'customer' | 'majhi' | 'godown' | 'broker'

export type DiscountType = 'flat' | 'percent'

export type UnbilledBillType =
    | 'customer_store'
    | 'customer_delivery'
    | 'customer_additional'
    | 'majhi'
    | 'godown'

export interface Discount {
    id: number
    invoice_id: number
    type: DiscountType
    value: number
    computed_amount: number
    reason: string | null
    created_at: string
    updated_at: string
}

export interface InvoiceItem {
    id: number
    invoice_id: number
    bill_type: UnbilledBillType
    bill_id: number
    amount: number
    created_at: string
}

export interface Invoice {
    id: number
    user_id: number
    entity_type: EntityType
    entity_id: number
    subtotal: number
    discount_amount: number
    total: number
    notes: string | null
    created_at: string
    updated_at: string
}

export interface InvoiceDetail extends Invoice {
    items: InvoiceItem[]
    discounts: Discount[]
}

export interface UnbilledBill {
    bill_type: UnbilledBillType
    bill_id: number
    label: string
    total: number
    paid: number
    remaining: number
    date: string
}

export interface DiscountInput {
    type: DiscountType
    value: number
    reason?: string | null
}

export interface BillRef {
    bill_type: UnbilledBillType
    bill_id: number
}

export interface CreateInvoicePayload {
    entity_type: EntityType
    entity_id: number
    bills: BillRef[]
    discounts?: DiscountInput[]
    notes?: string | null
}

export interface UpdateInvoicePayload {
    discounts?: DiscountInput[]
    notes?: string | null
}

export type InvoiceSortField =
    | 'entity_id'
    | 'subtotal'
    | 'discount_amount'
    | 'total'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'