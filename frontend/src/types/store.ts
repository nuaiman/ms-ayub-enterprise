// src/types/store.ts

export type StoreBillType = 'weight' | 'quantity'

export interface Store {
  id: number
  lot_id: number
  godown_id: number
  store_bill_type: StoreBillType
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
  created_at: string
  updated_at: string
}

export interface CreateStorePayload {
  lot_id: number
  godown_id: number
  store_bill_type?: StoreBillType
  godown_cut: number
  quantity: number
  quantity_unit?: string
  weight: number
  weight_unit?: string
  billing_start?: string
  billing_end?: string | null
  last_paid_through?: string | null
  last_paid_amount?: number
  notes?: string | null
}

export interface UpdateStorePayload {
  store_bill_type?: StoreBillType
  godown_cut?: number
  quantity?: number
  quantity_unit?: string
  weight?: number
  weight_unit?: string
  billing_start?: string
  billing_end?: string | null
  last_paid_through?: string | null
  last_paid_amount?: number
  notes?: string | null
}

export type StoreSortField = 'lot_id' | 'godown_id' | 'quantity' | 'weight' | 'is_active' | 'created_at'
export type SortDirection = 'asc' | 'desc'