// src/types/customerStorageBill.ts

export interface CustomerStorageBill {
    id: number
    lot_id: number
    item_name: string
    customer_id: number | null
    customer_name: string
    customer_charge_type: 'weight' | 'quantity'
    customer_storage_rate: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    monthly_bill: number
    billing_start: string | null
    billing_end: string | null
    months_billed: number
    total_billed: number
    total_paid: number
    outstanding: number
    last_paid_through: string | null
    notes: string | null
}

export type CustomerStorageBillSortField = 'customer_name' | 'item_name' | 'monthly_bill' | 'total_billed' | 'outstanding'
export type SortDirection = 'asc' | 'desc'