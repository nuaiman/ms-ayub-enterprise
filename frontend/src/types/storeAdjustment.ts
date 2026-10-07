// src/types/storeAdjustment.ts

export type StoreAdjustmentType = 'delta' | 'absolute'

export interface StoreAdjustment {
    id: number
    user_id: number
    store_id: number
    adjustment_type: StoreAdjustmentType
    input_weight: number
    input_quantity: number
    weight_delta: number
    quantity_delta: number
    reason: string | null
    notes: string | null
    adjusted_at: string
    created_at: string
    updated_at: string
}

export interface CreateStoreAdjustmentPayload {
    adjustment_type: StoreAdjustmentType
    input_weight: number
    input_quantity: number
    reason?: string | null
    notes?: string | null
}

export interface CreateStoreAdjustmentResponse {
    adjustment: StoreAdjustment
    store: {
        id: number
        weight: number
        quantity: number
        is_active: boolean
        // ...other store fields are present but we only care about these here
        [key: string]: unknown
    }
}

export type StoreAdjustmentSortField =
    | 'adjusted_at'
    | 'store_id'
    | 'user_id'
    | 'weight_delta'
    | 'quantity_delta'

export type SortDirection = 'asc' | 'desc'