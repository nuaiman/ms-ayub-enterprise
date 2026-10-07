// src/types/income.ts

export interface Income {
    id: number
    user_id: number
    title: string
    amount: number
    income_date: string
    notes: string | null
    image_url: string | null
    created_at: string
    updated_at: string
}

export interface CreateIncomePayload {
    title: string
    amount: number
    income_date?: string
    notes?: string | null
}

export interface UpdateIncomePayload {
    title?: string
    amount?: number
    income_date?: string
    notes?: string | null
}

export type IncomeSortField = 'title' | 'amount' | 'income_date' | 'user_id' | 'created_at'
export type SortDirection = 'asc' | 'desc'