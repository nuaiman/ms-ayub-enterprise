// src/types/majhiLotBill.ts

export interface MajhiLotBill {
    id: number
    lot_id: number
    item_name: string
    majhi_id: number | null
    majhi_name: string
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

export type MajhiLotBillSortField = 'majhi_name' | 'item_name' | 'bill_amount' | 'status' | 'created_at'
export type SortDirection = 'asc' | 'desc'