// src/types/rent.ts

export type RentStatus = 'draft' | 'paid' | 'cancelled'

export type PaymentMethod = 'cash' | 'bank_transfer' | 'check' | 'mobile_banking'

export interface Rent {
    id: number
    user_id: number
    godown_id: number
    month_year: string
    amount: number
    status: RentStatus
    payment_date: string | null
    payment_method: PaymentMethod | null
    reference_number: string | null
    notes: string | null
    created_at: string
    updated_at: string
}

export interface CreateRentPayload {
    godown_id: number
    month_year: string
    amount: number
    notes?: string | null
}

export interface UpdateRentPayload {
    amount?: number
    notes?: string | null
}

export interface MarkRentAsPaidPayload {
    payment_date?: string
    payment_method?: PaymentMethod
    reference_number?: string | null
}

export interface RentFilters {
    godown_id?: string
    month_year?: string
    status?: RentStatus
}

export type RentSortField = 'godown_id' | 'month_year' | 'amount' | 'status' | 'created_at'
export type SortDirection = 'asc' | 'desc'