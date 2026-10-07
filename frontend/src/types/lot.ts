// src/types/lot.ts

export interface Lot {
  id: number
  user_id: number
  customer_id: number
  lot_number: string
  product_name: string
  weight_unit: string
  quantity_unit: string
  created_at: string
  updated_at: string
}

export interface CreateLotPayload {
  customer_id: number
  lot_number: string
  product_name: string
  weight_unit?: string
  quantity_unit?: string
}

export interface UpdateLotPayload {
  customer_id?: number
  lot_number?: string
  product_name?: string
  weight_unit?: string
  quantity_unit?: string
}

export type LotSortField = 'customer_id' | 'lot_number' | 'product_name' | 'created_at'
export type SortDirection = 'asc' | 'desc'