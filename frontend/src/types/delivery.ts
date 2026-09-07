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

export type DeliverySortField = 'customer_id' | 'delivery_date' | 'receiver_name' | 'from_location' | 'to_location' | 'created_at'
export type SortDirection = 'asc' | 'desc'