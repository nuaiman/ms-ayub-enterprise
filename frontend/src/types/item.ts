// src/types/item.ts

export interface Item {
  id: number
  user_id: number
  customer_id: number | null
  product_name: string | null
  category: string | null
  is_active: boolean
  notes: string | null
  image_url: string | null
  created_at: string
  updated_at: string
}

export interface CreateItemPayload {
  customer_id?: number | null
  product_name?: string | null
  category?: string | null
  notes?: string | null
}

export interface UpdateItemPayload {
  customer_id?: number | null
  product_name?: string | null
  category?: string | null
  notes?: string | null
}

export type ItemSortField = 'product_name' | 'category' | 'is_active' | 'created_at'
export type SortDirection = 'asc' | 'desc'