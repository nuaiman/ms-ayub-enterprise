// src/types/damage.ts

export interface Damage {
    id: number
    user_id: number
    store_id: number
    quantity: number
    quantity_unit: string
    weight: number
    weight_unit: string
    damage_date: string
    reason: string
    amount: number
    notes: string | null
    image_url: string | null
    created_at: string
    updated_at: string
}

export interface CreateDamagePayload {
    store_id: number
    quantity?: number
    weight?: number
    damage_date?: string
    reason: string
    amount?: number
    notes?: string | null
    image_url?: string | null
}

export interface UpdateDamagePayload {
    quantity?: number
    weight?: number
    damage_date?: string
    reason?: string
    amount?: number
    notes?: string | null
    image_url?: string | null
}

export type DamageSortField =
    | 'store_id'
    | 'user_id'
    | 'damage_date'
    | 'amount'
    | 'quantity'
    | 'weight'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'