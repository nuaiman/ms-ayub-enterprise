// src/types/customerLotBill.ts

export interface CustomerLotBill {
    id: number
    lot_id: number
    item_name: string
    customer_id: number | null
    customer_name: string
    customer_charge_type: 'weight' | 'quantity'
    unload_rate: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    bill_amount: number
    paid_amount: number  // From customer_paid_unload_amount
    status: 'unpaid' | 'paid' | 'cancelled'
    payment_date: string | null
    notes: string | null
    created_at: string
    updated_at: string
}

export type CustomerLotBillSortField = 'customer_name' | 'item_name' | 'bill_amount' | 'status' | 'created_at'
export type SortDirection = 'asc' | 'desc'