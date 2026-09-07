// src/types/customerTransportBill.ts

export interface CustomerTransportBillVehicle {
    vehicle_id: number
    vehicle_number: string
    customer_charge: number
    joma_cost: number
    vehicle_cost: number
    broker_name: string | null
}

export interface CustomerTransportBill {
    id: number
    transport_id: number
    customer_id: number | null
    customer_name: string
    from_location: string
    to_location: string | null
    vehicle_quantity: number
    total_vehicles: number
    bill_amount: number
    paid_amount: number  // From customer_total_paid
    status: 'unpaid' | 'paid' | 'cancelled'
    payment_date: string | null
    notes: string | null
    vehicles: CustomerTransportBillVehicle[]  // Breakdown of vehicles
    created_at: string
    updated_at: string
}

export type CustomerTransportBillSortField =
    | 'transport_id'
    | 'customer_name'
    | 'from_location'
    | 'to_location'      // Added
    | 'vehicle_quantity'
    | 'bill_amount'
    | 'status'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'