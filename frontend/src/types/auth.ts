// src/types/auth.ts

export type Role = 'admin' | 'manager' | 'accounts' | 'staff'

export type IDType =
  | 'nid'
  | 'passport'
  | 'driving_license'
  | 'birth_certificate'
  | 'trade_license'
  | 'other'

export interface User {
  id: number
  name: string
  username: string
  email: string | null
  phone: string | null
  address: string | null
  id_type: IDType | null
  id_number: string | null
  image_url: string | null
  role: Role
  is_active: boolean
  monthly_salary: number
  created_at: string
  updated_at: string
}

export interface CreateUserPayload {
  name: string
  username: string
  password: string
  role: Role
  email?: string | null
  phone?: string | null
  address?: string | null
  id_type?: string | null
  id_number?: string | null
  monthly_salary?: number
  image_url?: string | null
}

export interface UpdateProfilePayload {
  name: string
  email: string | null
  phone: string | null
  address: string | null
  id_type: IDType | null
  id_number: string | null
}

export interface UpdateUserProfilePayload {
  name: string
  email: string | null
  phone: string | null
  address: string | null
  id_type: IDType | null
  id_number: string | null
}

export interface UpdateUserSalaryPayload {
  monthly_salary: number
}


export type UserSortField = 'name' | 'username' | 'role' | 'created_at'

export type SortDirection = 'asc' | 'desc'

export type UserSortKey = 'newest' | 'oldest' | 'name_asc' | 'name_desc'