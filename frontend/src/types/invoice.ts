// src/types/invoice.ts

export type InvoiceType = 'godown' | 'transport'
export type InvoiceStatus = 'draft' | 'finalized' | 'printed' | 'cancelled'

export interface InvoiceItem {
    id: string
    source_type: 'storage_bill' | 'lot_bill' | 'delivery_bill' | 'transport_bill'
    source_id: number
    item: string
    date: string
    description: string
    quantity: number
    rate: number
    amount: number          // Outstanding amount (what's still due)
    total_amount: number    // Full bill amount
    paid_amount: number     // Amount already paid
}

export interface Invoice {
    id: string
    number: string
    type: InvoiceType
    customer_id: number | null
    customer_name: string
    date: string
    items: InvoiceItem[]
    total: number           // Sum of all outstanding amounts
    received: number        // Sum of all paid amounts from selected items
    notes: string
    status: InvoiceStatus
    created_at: string
    updated_at: string
}

export interface AvailableItem {
    id: number
    source_type: 'storage_bill' | 'lot_bill' | 'delivery_bill' | 'transport_bill'
    item: string
    date: string
    description: string
    quantity: number
    rate: number
    amount: number           // Outstanding amount (bill - paid)
    total_amount: number     // Full bill amount
    paid_amount: number      // Amount already paid
    status: 'unpaid' | 'partial' | 'paid'
    selected: boolean
}

export const generateInvoiceNumber = (type: InvoiceType, sequence: number): string => {
    const year = new Date().getFullYear()
    const prefix = type === 'godown' ? 'GOD' : 'TRN'
    const seq = String(sequence).padStart(4, '0')
    return `${prefix}-${year}-${seq}`
}