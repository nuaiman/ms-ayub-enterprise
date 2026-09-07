// src/types/transport.ts

export type DeliveryType = 'local' | 'district'

export interface Transport {
  id: number
  user_id: number
  customer_id: number | null
  from_location: string
  to_location: string | null
  vehicle_quantity: number
  delivery_type: DeliveryType | null
  notes: string | null
  transport_date: string
  office_commission_amount: number
  image_url: string | null
  customer_total_paid: number  // NEW - not required for creation
  created_at: string
  updated_at: string
}

export interface CreateTransportPayload {
  customer_id?: number | null
  from_location: string
  to_location?: string | null
  vehicle_quantity: number
  delivery_type?: DeliveryType | null
  notes?: string | null
  transport_date?: string
  office_commission_amount: number
  // customer_total_paid is NOT required - defaults to 0
}

export interface UpdateTransportPayload {
  customer_id?: number | null
  from_location?: string
  to_location?: string | null
  vehicle_quantity?: number
  delivery_type?: DeliveryType | null
  notes?: string | null
  transport_date?: string
  office_commission_amount?: number
  // customer_total_paid is NOT included - use separate endpoint
}

export interface UpdateCustomerPaymentPayload {
  customer_total_paid: number  // NEW - for dedicated endpoint
}

export type TransportSortField = 'customer_id' | 'from_location' | 'to_location' | 'vehicle_quantity' | 'delivery_type' | 'transport_date' | 'office_commission_amount' | 'created_at'
export type SortDirection = 'asc' | 'desc'