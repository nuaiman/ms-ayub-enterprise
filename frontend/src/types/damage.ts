// src/types/damage.ts

export interface Damage {
  id: number
  store_id: number
  user_id: number
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
  quantity: number
  quantity_unit?: string
  weight: number
  weight_unit?: string
  damage_date?: string
  reason: string
  amount: number
  notes?: string | null
}

export interface UpdateDamagePayload {
  quantity?: number
  quantity_unit?: string
  weight?: number
  weight_unit?: string
  damage_date?: string
  reason?: string
  amount?: number
  notes?: string | null
}

export type DamageSortField = 'store_id' | 'user_id' | 'quantity' | 'weight' | 'damage_date' | 'amount' | 'created_at'
export type SortDirection = 'asc' | 'desc'