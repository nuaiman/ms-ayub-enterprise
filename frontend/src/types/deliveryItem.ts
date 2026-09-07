// src/types/deliveryItem.ts

export interface DeliveryItem {
  id: number
  delivery_id: number
  store_id: number
  majhi_id: number | null
  item_id: number
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
  // NEW BILLING FIELDS (snapshot at time of delivery)
  customer_charge_type: 'weight' | 'quantity'
  customer_paid_unload_amount: number  // Hidden - not displayed in forms
  majhi_bill_type: 'weight' | 'quantity' | 'job'
  majhi_total_paid: number  // Hidden - not displayed in forms
  created_at: string
  updated_at: string
}

export interface CreateDeliveryItemPayload {
  delivery_id: number
  store_id: number
  majhi_id?: number | null
  item_id: number
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
  // NEW - optional, defaults from lot if not provided
  customer_charge_type?: 'weight' | 'quantity'
  customer_paid_unload_amount?: number  // Hidden - not displayed
  majhi_bill_type?: 'weight' | 'quantity' | 'job'
  majhi_total_paid?: number  // Hidden - not displayed
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
  // NEW - optional
  customer_charge_type?: 'weight' | 'quantity'
  customer_paid_unload_amount?: number  // Hidden - not displayed
  majhi_bill_type?: 'weight' | 'quantity' | 'job'
  majhi_total_paid?: number  // Hidden - not displayed
}

export type DeliveryItemSortField = 'delivery_id' | 'store_id' | 'item_id' | 'lot_id' | 'quantity' | 'weight' | 'created_at'
export type SortDirection = 'asc' | 'desc'