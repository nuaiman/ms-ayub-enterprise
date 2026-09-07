// src/types/brokerVehicleBill.ts

export interface BrokerVehicleBill {
    id: number
    vehicle_id: number
    transport_id: number
    vehicle_number: string
    broker_id: number | null
    broker_name: string
    joma_cost: number
    vehicle_cost: number
    bill_amount: number
    paid_amount: number  // From broker_total_paid
    status: 'unpaid' | 'paid' | 'cancelled'
    payment_date: string | null
    notes: string | null
    created_at: string
    updated_at: string
}

export type BrokerVehicleBillSortField = 'vehicle_number' | 'broker_name' | 'transport_id' | 'bill_amount' | 'status' | 'created_at'
export type SortDirection = 'asc' | 'desc'