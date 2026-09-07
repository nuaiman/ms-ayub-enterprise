// src/types/majhi.ts

export interface Majhi {
  id: number
  name: string
  phone: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

export interface CreateMajhiPayload {
  name: string
  phone?: string | null
  notes?: string | null
}

export interface UpdateMajhiPayload {
  name?: string
  phone?: string | null
  notes?: string | null
}

export type MajhiSortField = 'name' | 'phone' | 'created_at'
export type SortDirection = 'asc' | 'desc'