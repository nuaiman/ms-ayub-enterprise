// src/stores/customerTransportBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useTransportsStore } from './transports'
import { useVehiclesStore } from './vehicles'
import { useCustomersStore } from './customers'
import { useBrokersStore } from './brokers'
import type {
    CustomerTransportBill,
    CustomerTransportBillVehicle,
    CustomerTransportBillSortField,
} from '@/types/customerTransportBill'
import type { SortDirection } from '@/types/transport'
import { push } from 'notivue'

export const useCustomerTransportBillsStore = defineStore('customerTransportBills', () => {
    const transportsStore = useTransportsStore()
    const vehiclesStore = useVehiclesStore()
    const customersStore = useCustomersStore()
    const brokersStore = useBrokersStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<CustomerTransportBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    const getVehicleBreakdown = (transportId: number): CustomerTransportBillVehicle[] => {
        const vehicles = vehiclesStore.getVehiclesByTransportId(transportId)
        return vehicles.map(v => ({
            vehicle_id: v.id,
            vehicle_number: v.vehicle_number,
            broker_name: v.broker_id ? brokersStore.getBrokerName(v.broker_id) : null,
            joma_cost: v.joma_cost || 0,
            vehicle_cost: v.vehicle_cost || 0,
        }))
    }

    const bills = computed<CustomerTransportBill[]>(() => {
        const transports = transportsStore.transports
        const result: CustomerTransportBill[] = []

        for (const transport of transports) {
            const vehicles = vehiclesStore.getVehiclesByTransportId(transport.id)
            if (vehicles.length === 0) continue

            const billAmount = transport.customer_total_charge || 0
            if (billAmount === 0) continue

            const customerName = customersStore.getCustomerName(transport.customer_id)

            const paidAmount = transport.customer_total_paid || 0
            const status = paidAmount >= billAmount ? 'paid' : 'unpaid'

            result.push({
                id: transport.id,
                customer_id: transport.customer_id,
                customer_name: customerName,
                from_location: transport.from_location,
                to_location: transport.to_location,
                vehicle_quantity: transport.vehicle_quantity,
                total_vehicles: vehicles.length,
                customer_charge_unit: transport.customer_charge_unit,
                customer_total_unit: transport.customer_total_unit,
                customer_charge_per_unit: transport.customer_charge_per_unit,
                bill_amount: billAmount,
                paid_amount: paidAmount,
                payment_date: transport.customer_total_paid_through,
                status,
                notes: transport.notes || null,
                vehicles: getVehicleBreakdown(transport.id),
                created_at: transport.created_at,
                updated_at: transport.updated_at,
            })
        }

        return result
    })

    const filteredBills = computed(() => {
        let result = [...bills.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(bill =>
                bill.customer_name.toLowerCase().includes(query) ||
                bill.from_location.toLowerCase().includes(query) ||
                (bill.to_location && bill.to_location.toLowerCase().includes(query)) ||
                String(bill.id).includes(query) ||
                bill.status.toLowerCase().includes(query)
            )
        }

        if (statusFilter.value) {
            result = result.filter(bill => bill.status === statusFilter.value)
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'id':
                    comparison = a.id - b.id
                    break
                case 'customer_name':
                    comparison = a.customer_name.localeCompare(b.customer_name)
                    break
                case 'from_location':
                    comparison = a.from_location.localeCompare(b.from_location)
                    break
                case 'to_location':
                    comparison = (a.to_location || '').localeCompare(b.to_location || '')
                    break
                case 'vehicle_quantity':
                    comparison = a.vehicle_quantity - b.vehicle_quantity
                    break
                case 'bill_amount':
                    comparison = a.bill_amount - b.bill_amount
                    break
                case 'status':
                    comparison = a.status.localeCompare(b.status)
                    break
                case 'created_at':
                    comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
                    break
                default:
                    comparison = 0
            }
            return sortDirection.value === 'desc' ? -comparison : comparison
        })

        return result
    })

    const totalBills = computed(() => bills.value.length)
    const totalUnpaid = computed(() => bills.value.filter(b => b.status === 'unpaid').length)
    const totalPaid = computed(() => bills.value.filter(b => b.status === 'paid').length)
    const totalAmount = computed(() => bills.value.reduce((sum, b) => sum + b.bill_amount, 0))
    const totalUnpaidAmount = computed(() =>
        bills.value
            .filter(b => b.status === 'unpaid')
            .reduce((sum, b) => sum + b.bill_amount, 0)
    )

    const markBillAsPaid = async (
        transportId: number,
        payload: { amount: number; payment_date?: string; notes?: string | null }
    ): Promise<boolean> => {
        try {
            const transport = transportsStore.getTransportById(transportId)
            if (!transport) {
                push.error('Transport not found')
                return false
            }

            const bill = bills.value.find(b => b.id === transportId)
            if (!bill) {
                push.error('Bill not found')
                return false
            }

            if (bill.status === 'paid') {
                push.info('Bill is already fully paid')
                return true
            }

            const newTotalPaid = (transport.customer_total_paid || 0) + payload.amount

            if (newTotalPaid > bill.bill_amount) {
                const remaining = bill.bill_amount - bill.paid_amount
                push.error(`Payment amount exceeds remaining balance of ${remaining.toFixed(2)}`)
                return false
            }

            const result = await transportsStore.updateCustomerPayment(transportId, {
                customer_total_paid: newTotalPaid,
                customer_total_paid_through: payload.payment_date || null,
            })

            if (result) {
                await transportsStore.fetchTransports()
                push.success('Payment recorded successfully')
                return true
            }

            return false
        } catch (error) {
            console.error('Error marking customer transport bill as paid:', error)
            push.error('Failed to record payment')
            return false
        }
    }

    const cancelBill = async (transportId: number): Promise<boolean> => {
        try {
            const transport = transportsStore.getTransportById(transportId)
            if (!transport) {
                push.error('Transport not found')
                return false
            }

            const result = await transportsStore.updateTransport(transportId, {
                customer_charge_per_unit: 0,
            })

            if (result) {
                await transportsStore.fetchTransports()
                push.success('Bill cancelled successfully')
                return true
            }

            return false
        } catch (error) {
            console.error('Error cancelling customer transport bill:', error)
            push.error('Failed to cancel bill')
            return false
        }
    }

    const setSort = (field: CustomerTransportBillSortField) => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
        } else {
            sortField.value = field
            sortDirection.value = 'asc'
        }
    }

    const setSearchQuery = (query: string) => {
        searchQuery.value = query
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    const setStatusFilter = (status: 'unpaid' | 'paid' | 'cancelled' | '') => {
        statusFilter.value = status
    }

    const getBillById = (id: number): CustomerTransportBill | undefined => {
        return bills.value.find(b => b.id === id)
    }

    const getBillsByCustomerId = (customerId: number): CustomerTransportBill[] => {
        return bills.value.filter(b => b.customer_id === customerId)
    }

    const getStatusBadgeClass = (status: string): string => {
        switch (status) {
            case 'unpaid': return 'border-(--color-yellow) text-(--color-yellow)'
            case 'paid': return 'border-(--color-green) text-(--color-green)'
            case 'cancelled': return 'border-(--color-red) text-(--color-red)'
            default: return 'border-(--color-border) text-(--color-text-secondary)'
        }
    }

    const getStatusDotClass = (status: string): string => {
        switch (status) {
            case 'unpaid': return 'bg-(--color-yellow)'
            case 'paid': return 'bg-(--color-green)'
            case 'cancelled': return 'bg-(--color-red)'
            default: return 'bg-(--color-text-secondary)'
        }
    }

    const getStatusLabel = (status: string): string => {
        switch (status) {
            case 'unpaid': return 'Unpaid'
            case 'paid': return 'Paid'
            case 'cancelled': return 'Cancelled'
            default: return status
        }
    }

    const formatBillDate = (dateStr: string): string => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            month: 'short', day: 'numeric', year: 'numeric',
            hour: '2-digit', minute: '2-digit'
        })
    }

    return {
        searchQuery,
        statusFilter,
        sortField,
        sortDirection,

        bills,
        filteredBills,
        totalBills,
        totalUnpaid,
        totalPaid,
        totalAmount,
        totalUnpaidAmount,

        markBillAsPaid,
        cancelBill,

        setSort,
        setSearchQuery,
        clearSearch,
        setStatusFilter,

        getBillById,
        getBillsByCustomerId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        formatBillDate,
    }
})