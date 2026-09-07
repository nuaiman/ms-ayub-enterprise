// src/types/godownStoreBill.ts

// This is the store data enriched with computed fields for billing
export interface GodownStoreBillStore {
    id: number
    lot_id: number
    godown_id: number
    store_bill_type: 'weight' | 'quantity'
    godown_cut: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    is_active: boolean
    billing_start: string
    billing_end: string | null
    last_paid_through: string | null
    last_paid_amount: number
    notes: string | null
    // Computed fields (added by the store)
    lot_name?: string
    item_name?: string
    customer_name?: string
    godown_name?: string
    monthly_bill?: number
    total_billed?: number
    outstanding?: number
}

export type GodownStoreBillSortField =
    | 'lot_name'
    | 'godown_name'
    | 'customer_name'
    | 'monthly_bill'
    | 'total_billed'
    | 'outstanding'
    | 'last_paid_through'

export type SortDirection = 'asc' | 'desc'