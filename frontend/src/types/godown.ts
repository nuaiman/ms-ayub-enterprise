// src/types/godown.ts

export interface Godown {
  id: number
  name: string
  phone: string | null
  notes: string | null
  is_active: boolean
  monthly_rent: number
  created_at: string
  updated_at: string
}

export interface CreateGodownPayload {
  name: string
  phone?: string | null
  notes?: string | null
  is_active?: boolean
  monthly_rent?: number
}

export interface UpdateGodownPayload {
  name?: string
  phone?: string | null
  notes?: string | null
  is_active?: boolean
  monthly_rent?: number
}

export type GodownSortField = 'name' | 'phone' | 'is_active' | 'monthly_rent' | 'created_at'
export type SortDirection = 'asc' | 'desc'