// src/types/delivery.ts

export interface Delivery {
    id: number
    user_id: number
    customer_id: number | null
    delivery_date: string
    receiver_name: string | null
    receiver_phone: string | null
    from_location: string | null
    to_location: string | null
    notes: string | null
    image_url: string | null
    created_at: string
    updated_at: string
}

export interface DeliveryItem {
    id: number
    delivery_id: number
    store_id: number
    majhi_id: number
    vehicle_number: string | null
    driver_number: string | null
    quantity: number
    weight: number
    created_at: string
    updated_at: string
}

export interface CreateDeliveryPayload {
    customer_id?: number | null
    delivery_date?: string
    receiver_name?: string | null
    receiver_phone?: string | null
    from_location?: string | null
    to_location?: string | null
    notes?: string | null
}

export interface UpdateDeliveryPayload {
    customer_id?: number | null
    delivery_date?: string
    receiver_name?: string | null
    receiver_phone?: string | null
    from_location?: string | null
    to_location?: string | null
    notes?: string | null
}

export interface CreateDeliveryItemPayload {
    store_id: number
    majhi_id: number
    vehicle_number?: string | null
    driver_number?: string | null
    quantity: number
    weight: number
}

export interface UpdateDeliveryItemPayload {
    store_id?: number
    majhi_id?: number
    vehicle_number?: string | null
    driver_number?: string | null
    quantity?: number
    weight?: number
}

export interface DeliveryDetailResponse {
    delivery: Delivery
    items: DeliveryItem[]
}

export type DeliverySortField =
    | 'customer_id'
    | 'delivery_date'
    | 'created_at'

export type SortDirection = 'asc' | 'desc'