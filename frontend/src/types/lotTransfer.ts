// src/types/lotTransfer.ts

export type LotBillType = 'weight' | 'quantity'

export interface LotTransfer {
    id: number
    user_id: number
    lot_id: number
    from_customer_id: number
    to_customer_id: number
    notes: string | null
    transferred_at: string
    created_at: string
    updated_at: string
}

export interface CreateLotTransferPayload {
    to_customer_id: number
    notes?: string | null
    // Optional billing block — if bill_type is set and rate > 0, backend
    // creates a customer_store_bill for each active store under the lot,
    // for the current month, under the new customer.
    bill_type?: LotBillType
    rate?: number
}

export interface CreateLotTransferResponse {
    transfer: LotTransfer
    lot: {
        id: number
        customer_id: number
        lot_number: string
        product_name: string
        [key: string]: unknown
    }
    bills_created: number
}

export type LotTransferSortField =
    | 'transferred_at'
    | 'lot_id'
    | 'user_id'
    | 'from_customer_id'
    | 'to_customer_id'

export type SortDirection = 'asc' | 'desc'