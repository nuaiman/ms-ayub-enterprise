// src/types/vehicle.ts

export interface Vehicle {
  id: number
  user_id: number
  transport_id: number
  vehicle_number: string
  broker_id: number
  joma_cost: number
  vehicle_cost: number
  total_paid_to_broker: number
  other_cost: number
  labour_cost: number
  demarage_cost: number
  notes: string | null
  created_at: string
  updated_at: string
}

export interface CreateVehiclePayload {
  transport_id: number
  vehicle_number: string
  broker_id: number
  joma_cost: number
  vehicle_cost: number
  other_cost: number
  labour_cost: number
  demarage_cost: number
  notes?: string | null
}

export interface UpdateVehiclePayload {
  vehicle_number?: string
  broker_id?: number
  joma_cost?: number
  vehicle_cost?: number
  other_cost?: number
  labour_cost?: number
  demarage_cost?: number
  notes?: string | null
}

export interface UpdateBrokerPaymentPayload {
  total_paid_to_broker: number
}

export type VehicleSortField =
  | 'transport_id'
  | 'vehicle_number'
  | 'broker_id'
  | 'joma_cost'
  | 'vehicle_cost'
  | 'created_at'

export type SortDirection = 'asc' | 'desc'