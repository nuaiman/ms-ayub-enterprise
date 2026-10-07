// src/types/godown.ts

export interface Godown {
  id: number
  name: string
  phone: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export interface CreateGodownPayload {
  name: string
  phone?: string | null
  notes?: string | null
}

export interface UpdateGodownPayload {
  name?: string
  phone?: string | null
  notes?: string | null
}

export type GodownSortField = 'name' | 'phone' | 'created_at'
export type SortDirection = 'asc' | 'desc'