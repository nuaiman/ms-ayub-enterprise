// src/types/storeTransfer.ts

export interface StoreTransfer {
    id: number
    user_id: number
    store_id: number
    from_godown_id: number
    to_godown_id: number
    weight_at_transfer: number
    quantity_at_transfer: number
    notes: string | null
    transferred_at: string
    created_at: string
    updated_at: string
}

export interface CreateStoreTransferPayload {
    to_godown_id: number
    notes?: string | null
}

export interface CreateStoreTransferResponse {
    transfer: StoreTransfer
    store: {
        id: number
        godown_id: number
        weight: number
        quantity: number
        is_active: boolean
        [key: string]: unknown
    }
}

export type StoreTransferSortField =
    | 'transferred_at'
    | 'store_id'
    | 'user_id'
    | 'from_godown_id'
    | 'to_godown_id'

export type SortDirection = 'asc' | 'desc'