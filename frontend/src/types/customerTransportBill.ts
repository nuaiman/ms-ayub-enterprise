// src/types/customerTransportBill.ts

import type { CustomerChargeUnit } from './transport'

export interface CustomerTransportBillVehicle {
    vehicle_id: number
    vehicle_number: string
    broker_name: string | null
    joma_cost: number
    vehicle_cost: number
}

export interface CustomerTransportBill {
    id: number                 // transport id
    customer_id: number
    customer_name: string
    from_location: string
    to_location: string | null
    vehicle_quantity: number
    total_vehicles: number
    customer_charge_unit: CustomerChargeUnit
    customer_total_unit: number
    customer_charge_per_unit: number
    bill_amount: number
    paid_amount: number
    payment_date: string | null
    status: 'unpaid' | 'paid' | 'cancelled'
    notes: string | null
    vehicles: CustomerTransportBillVehicle[]
    created_at: string
    updated_at: string
}

export type CustomerTransportBillSortField =
    | 'id'
    | 'customer_name'
    | 'from_location'
    | 'to_location'
    | 'vehicle_quantity'
    | 'bill_amount'
    | 'status'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'