// src/types/majhiLoadingBill.ts

export interface MajhiLoadingBill {
    id: number
    delivery_item_id: number
    delivery_id: number
    delivery_date: string
    majhi_id: number | null
    majhi_name: string
    item_name: string
    lot_id: number
    store_id: number
    majhi_bill_type: 'weight' | 'quantity' | 'job'
    majhi_cut: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    bill_amount: number
    paid_amount: number  // From majhi_total_paid
    status: 'unpaid' | 'paid' | 'cancelled'
    payment_date: string | null
    notes: string | null
    created_at: string
    updated_at: string
}

export type MajhiLoadingBillSortField = 'majhi_name' | 'item_name' | 'bill_amount' | 'status' | 'delivery_date' | 'created_at'
export type SortDirection = 'asc' | 'desc'