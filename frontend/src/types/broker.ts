// src/types/broker.ts

export interface Broker {
  id: number
  name: string
  phone: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export interface CreateBrokerPayload {
  name: string
  phone?: string | null
  notes?: string | null
}

export interface UpdateBrokerPayload {
  name?: string
  phone?: string | null
  notes?: string | null
}

export type BrokerSortField = 'name' | 'phone' | 'created_at'
export type SortDirection = 'asc' | 'desc'