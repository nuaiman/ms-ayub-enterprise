// src/types/store.ts

export interface Store {
  id: number
  user_id: number
  lot_id: number
  godown_id: number
  weight: number
  quantity: number
  start_date: string
  is_active: boolean
  image_url: string | null
  created_at: string
  updated_at: string
}

export interface CreateStorePayload {
  lot_id: number
  godown_id: number
  weight: number
  quantity: number
  start_date?: string
  is_active?: boolean
}

export interface UpdateStorePayload {
  weight?: number
  quantity?: number
  start_date?: string
  is_active?: boolean
}

export type StoreSortField =
  | 'lot_id'
  | 'godown_id'
  | 'weight'
  | 'quantity'
  | 'start_date'
  | 'is_active'
  | 'created_at'

export type SortDirection = 'asc' | 'desc'