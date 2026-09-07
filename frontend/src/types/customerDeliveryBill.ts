// src/types/customerDeliveryBill.ts

export interface CustomerDeliveryBill {
    id: number
    delivery_item_id: number
    delivery_id: number
    delivery_date: string
    customer_id: number | null
    customer_name: string
    item_name: string
    lot_id: number
    store_id: number
    customer_charge_type: 'weight' | 'quantity'
    loading_rate: number
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

export type CustomerDeliveryBillSortField = 'customer_name' | 'item_name' | 'bill_amount' | 'status' | 'delivery_date' | 'created_at'
export type SortDirection = 'asc' | 'desc'