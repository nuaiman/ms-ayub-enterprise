// src/types/lot.ts

export type CustomerChargeType = 'weight' | 'quantity'
export type MajhiBillType = 'weight' | 'quantity' | 'job'

export interface Lot {
  id: number
  user_id: number
  customer_id: number | null
  product_name: string | null
  category: string | null
  lot_number: number
  customer_charge_type: CustomerChargeType
  majhi_bill_type: MajhiBillType
  customer_storage_rate: number
  unload_rate: number
  majhi_id: number | null
  majhi_cut: number
  is_active: boolean
  notes: string | null
  image_url: string | null
  customer_last_paid_through: string | null
  customer_last_paid_amount: number
  customer_paid_unload_amount: number
  majhi_total_paid: number
  created_at: string
  updated_at: string
}

export interface CreateLotPayload {
  customer_id?: number | null
  product_name?: string | null
  category?: string | null
  lot_number: number
  customer_charge_type: CustomerChargeType
  majhi_bill_type: MajhiBillType
  customer_storage_rate: number
  unload_rate: number
  majhi_id?: number | null
  majhi_cut: number
  notes?: string | null
  customer_last_paid_through?: string | null
  customer_last_paid_amount?: number
  customer_paid_unload_amount?: number
  majhi_total_paid?: number
}

export interface UpdateLotPayload {
  customer_id?: number | null
  product_name?: string | null
  category?: string | null
  customer_charge_type?: CustomerChargeType
  majhi_bill_type?: MajhiBillType
  customer_storage_rate?: number
  unload_rate?: number
  majhi_id?: number | null
  majhi_cut?: number
  notes?: string | null
  customer_last_paid_through?: string | null
  customer_last_paid_amount?: number
  customer_paid_unload_amount?: number
  majhi_total_paid?: number
}

export type LotSortField = 'customer_id' | 'lot_number' | 'customer_charge_type' | 'is_active' | 'created_at'
export type SortDirection = 'asc' | 'desc'