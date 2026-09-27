// src/types/invoice.ts

export type InvoiceParty = 'customer' | 'broker' | 'godown' | 'majhi'

export type InvoiceType =
    | 'customer'
    | 'broker_vehicle'
    | 'godown_store'
    | 'majhi'

export type InvoiceStatus = 'draft' | 'finalized' | 'printed' | 'cancelled'

export interface InvoiceItem {
    id: string
    source_type:
    | 'storage_bill'
    | 'lot_bill'
    | 'delivery_bill'
    | 'transport_bill'
    | 'additional_charge'
    | 'broker_vehicle_bill'
    | 'godown_store_bill'
    | 'majhi_lot_bill'
    | 'majhi_loading_bill'
    source_id: number
    item: string
    date: string
    description: string
    quantity: number
    rate: number
    amount: number          // Outstanding amount
    total_amount: number    // Full bill amount
    paid_amount: number     // Amount already paid
}

export interface Invoice {
    id: string
    number: string
    party_type: InvoiceParty
    party_id: number
    party_name: string
    type: InvoiceType
    date: string
    items: InvoiceItem[]
    total: number
    received: number
    notes: string
    status: InvoiceStatus
    created_at: string
    updated_at: string
}

export interface AvailableItem {
    id: number
    source_type: InvoiceItem['source_type']
    item: string
    date: string
    description: string
    quantity: number
    rate: number
    amount: number
    total_amount: number
    paid_amount: number
    status: 'unpaid' | 'partial' | 'paid'
    selected: boolean
}

export const generateInvoiceNumber = (party: InvoiceParty, sequence: number): string => {
    const year = new Date().getFullYear()
    const prefixes: Record<InvoiceParty, string> = {
        customer: 'CUS',
        broker: 'BRK',
        godown: 'GDN',
        majhi: 'MJH',
    }
    const seq = String(sequence).padStart(4, '0')
    return `${prefixes[party]}-${year}-${seq}`
}

export const partyTypeLabel = (party: InvoiceParty): string => {
    const map: Record<InvoiceParty, string> = {
        customer: 'Customer',
        broker: 'Broker',
        godown: 'Godown',
        majhi: 'Majhi',
    }
    return map[party]
}

export const invoiceTitle = (party: InvoiceParty): string => {
    const map: Record<InvoiceParty, string> = {
        customer: 'Customer Bill',
        broker: 'Broker Vehicle Bill',
        godown: 'Godown Store Bill',
        majhi: 'Majhi Bill',
    }
    return map[party]
}