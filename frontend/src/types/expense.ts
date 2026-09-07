// src/types/expense.ts

export interface Expense {
  id: number
  user_id: number
  title: string
  amount: number
  expense_date: string
  notes: string | null
  image_url: string | null
  created_at: string
  updated_at: string
}

export interface CreateExpensePayload {
  title: string
  amount: number
  expense_date?: string
  notes?: string | null
}

export interface UpdateExpensePayload {
  title?: string
  amount?: number
  expense_date?: string
  notes?: string | null
}

// Added 'user_id' to the sort fields
export type ExpenseSortField = 'title' | 'amount' | 'expense_date' | 'user_id' | 'created_at'
export type SortDirection = 'asc' | 'desc'