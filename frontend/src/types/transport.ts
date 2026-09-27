// src/types/transport.ts

export type TransportType = 'local' | 'district'
export type CustomerChargeUnit = 'vehicle' | 'weight' | 'quantity'

export interface Transport {
  id: number
  user_id: number
  customer_id: number
  from_location: string
  to_location: string | null
  vehicle_quantity: number
  transport_date: string
  transport_type: TransportType | null
  image_url: string | null
  notes: string | null
  office_commission_amount: number

  // Customer billing
  customer_charge_unit: CustomerChargeUnit
  customer_total_unit: number
  customer_charge_per_unit: number
  customer_total_charge: number
  customer_total_paid: number
  customer_total_paid_through: string | null

  created_at: string
  updated_at: string
}

export interface CreateTransportPayload {
  customer_id: number
  from_location: string
  to_location?: string | null
  vehicle_quantity: number
  transport_date?: string
  transport_type?: TransportType | null
  notes?: string | null
  office_commission_amount: number

  customer_charge_unit: CustomerChargeUnit
  customer_total_unit: number
  customer_charge_per_unit: number
}

export interface UpdateTransportPayload {
  customer_id?: number
  from_location?: string
  to_location?: string | null
  vehicle_quantity?: number
  transport_date?: string
  transport_type?: TransportType | null
  notes?: string | null
  office_commission_amount?: number

  customer_charge_unit?: CustomerChargeUnit
  customer_total_unit?: number
  customer_charge_per_unit?: number
}

export interface UpdateTransportCustomerPaymentPayload {
  customer_total_paid: number
  customer_total_paid_through?: string | null
}

export type TransportSortField =
  | 'customer_id'
  | 'from_location'
  | 'to_location'
  | 'vehicle_quantity'
  | 'transport_type'
  | 'transport_date'
  | 'office_commission_amount'
  | 'created_at'

export type SortDirection = 'asc' | 'desc'