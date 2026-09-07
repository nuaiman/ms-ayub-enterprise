// src/types/dashboard.ts

export interface DashboardMetrics {
    // Received from customers
    totalReceived: number
    monthlyReceived: number
    collectionRate: number

    // Due from customers
    totalDue: number
    overdue: number
    outstandingCustomers: number

    // Payment obligations (what we owe)
    majhiUnpaidAmount: number
    brokerUnpaidAmount: number
    rentUnpaidAmount: number

    // Outstanding by customer
    outstandingByCustomer: OutstandingCustomer[]

    // Outstanding by type
    outstandingByType: OutstandingByType[]
}

export interface OutstandingCustomer {
    name: string
    amount: number
}

export interface OutstandingByType {
    name: string
    amount: number
    percentage: number
    color: string
}

export interface CustomerRevenue {
    name: string
    initials: string
    phone: string
    storage: number
    unload: number
    delivery: number
    transport: number
    total: number
    paid: number
    due: number
    paymentPercentage: number
}

export interface MonthOption {
    value: string
    label: string
}

export interface DashboardData {
    metrics: DashboardMetrics
    customerRevenue: CustomerRevenue[]
    availableMonths: MonthOption[]
}