// src/types/deliveryItem.ts

export interface DeliveryItem {
  id: number
  delivery_id: number
  store_id: number
  majhi_id: number | null
  lot_id: number
  vehicle_number: string | null
  driver_number: string | null
  quantity: number
  quantity_unit: string
  weight: number
  weight_unit: string
  loading_rate: number
  majhi_cut: number
  notes: string | null
  customer_charge_type: 'weight' | 'quantity'
  customer_paid_unload_amount: number
  majhi_bill_type: 'weight' | 'quantity' | 'job'
  majhi_total_paid: number
  created_at: string
  updated_at: string
}

export interface CreateDeliveryItemPayload {
  delivery_id: number
  store_id: number
  majhi_id?: number | null
  lot_id: number
  vehicle_number?: string | null
  driver_number?: string | null
  quantity: number
  quantity_unit?: string
  weight: number
  weight_unit?: string
  loading_rate: number
  majhi_cut: number
  notes?: string | null
  customer_charge_type?: 'weight' | 'quantity'
  customer_paid_unload_amount?: number
  majhi_bill_type?: 'weight' | 'quantity' | 'job'
  majhi_total_paid?: number
}

export interface UpdateDeliveryItemPayload {
  store_id?: number
  majhi_id?: number | null
  vehicle_number?: string | null
  driver_number?: string | null
  quantity?: number
  quantity_unit?: string
  weight?: number
  weight_unit?: string
  loading_rate?: number
  majhi_cut?: number
  notes?: string | null
  customer_charge_type?: 'weight' | 'quantity'
  customer_paid_unload_amount?: number
  majhi_bill_type?: 'weight' | 'quantity' | 'job'
  majhi_total_paid?: number
}

export type DeliveryItemSortField = 'delivery_id' | 'store_id' | 'lot_id' | 'quantity' | 'weight' | 'created_at'
export type SortDirection = 'asc' | 'desc'