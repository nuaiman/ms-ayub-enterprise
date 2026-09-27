// src/types/brokerVehicleBill.ts

export interface BrokerVehicleBill {
    id: number                 // vehicle id
    vehicle_number: string
    transport_id: number
    broker_id: number
    broker_name: string
    joma_cost: number
    vehicle_cost: number
    bill_amount: number        // joma_cost + vehicle_cost
    paid_amount: number        // from vehicle.total_paid_to_broker
    status: 'unpaid' | 'paid' | 'cancelled'
    payment_date: string | null
    notes: string | null
    created_at: string
    updated_at: string
}

export type BrokerVehicleBillSortField =
    | 'vehicle_number'
    | 'broker_name'
    | 'transport_id'
    | 'bill_amount'
    | 'status'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'