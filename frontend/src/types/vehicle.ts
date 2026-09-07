// src/types/vehicle.ts

export interface Vehicle {
  id: number
  user_id: number
  transport_id: number
  vehicle_number: string
  broker_id: number | null
  driver_name: string | null
  driver_phone: string | null
  joma_cost: number
  vehicle_cost: number
  customer_charge: number
  other_cost: number
  labour_cost: number
  demarage_amount: number
  demarage_reason: string | null
  broker_total_paid: number  // NEW
  created_at: string
  updated_at: string
}

export interface CreateVehiclePayload {
  transport_id: number
  vehicle_number: string
  broker_id?: number | null
  driver_name?: string | null
  driver_phone?: string | null
  joma_cost: number
  vehicle_cost: number
  customer_charge: number
  other_cost: number
  labour_cost: number
  demarage_amount: number
  demarage_reason?: string | null
  // broker_total_paid is NOT needed for creation - defaults to 0
}

export interface UpdateVehiclePayload {
  vehicle_number?: string
  broker_id?: number | null
  driver_name?: string | null
  driver_phone?: string | null
  joma_cost?: number
  vehicle_cost?: number
  customer_charge?: number
  other_cost?: number
  labour_cost?: number
  demarage_amount?: number
  demarage_reason?: string | null
  // broker_total_paid is NOT included here - use separate endpoint
}

export interface UpdateBrokerPaymentPayload {
  broker_total_paid: number  // NEW - for dedicated endpoint
}

export type VehicleSortField = 'transport_id' | 'vehicle_number' | 'broker_id' | 'driver_name' | 'joma_cost' | 'vehicle_cost' | 'customer_charge' | 'created_at'
export type SortDirection = 'asc' | 'desc'