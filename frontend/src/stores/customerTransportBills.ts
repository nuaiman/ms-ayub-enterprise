// src/stores/customerTransportBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useTransportsStore } from './transports'
import { useVehiclesStore } from './vehicles'
import { useCustomersStore } from './customers'
import { useBrokersStore } from './brokers'
import type { CustomerTransportBill, CustomerTransportBillVehicle, CustomerTransportBillSortField, SortDirection } from '@/types/customerTransportBill'

export const useCustomerTransportBillsStore = defineStore('customerTransportBills', () => {
    const transportsStore = useTransportsStore()
    const vehiclesStore = useVehiclesStore()
    const customersStore = useCustomersStore()
    const brokersStore = useBrokersStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<CustomerTransportBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    // ============= HELPERS =============

    // Get vehicle breakdown for a transport
    const getVehicleBreakdown = (transportId: number): CustomerTransportBillVehicle[] => {
        const vehicles = vehiclesStore.getVehiclesByTransportId(transportId)
        return vehicles.map(v => ({
            vehicle_id: v.id,
            vehicle_number: v.vehicle_number,
            customer_charge: v.customer_charge || 0,
            joma_cost: v.joma_cost || 0,
            vehicle_cost: v.vehicle_cost || 0,
            broker_name: v.broker_id ? brokersStore.getBrokerName(v.broker_id) : null,
        }))
    }

    // ============= COMPUTED =============

    // Generate bills for all transports with vehicles that have customer_charge > 0
    const bills = computed<CustomerTransportBill[]>(() => {
        // Get all transports
        const transports = transportsStore.transports

        const result: CustomerTransportBill[] = []

        for (const transport of transports) {
            // Get vehicles for this transport
            const vehicles = vehiclesStore.getVehiclesByTransportId(transport.id)

            // Skip if no vehicles
            if (vehicles.length === 0) continue

            // Calculate total customer charge
            const totalCustomerCharge = vehicles.reduce((sum, v) => sum + (v.customer_charge || 0), 0)

            // Skip if no customer charge
            if (totalCustomerCharge === 0) continue

            // Get customer name
            const customerName = transport.customer_id
                ? customersStore.getCustomerName(transport.customer_id)
                : 'No Customer'

            // Get paid amount from transport
            const paidAmount = transport.customer_total_paid || 0

            // Determine status - 'paid' only if fully paid, otherwise 'unpaid'
            const status = paidAmount >= totalCustomerCharge ? 'paid' : 'unpaid'

            // Get vehicle breakdown
            const vehicleBreakdown = getVehicleBreakdown(transport.id)

            result.push({
                id: transport.id,
                transport_id: transport.id,
                customer_id: transport.customer_id,
                customer_name: customerName,
                from_location: transport.from_location,
                to_location: transport.to_location,
                vehicle_quantity: transport.vehicle_quantity,
                total_vehicles: vehicles.length,
                bill_amount: totalCustomerCharge,
                paid_amount: paidAmount,
                status: status,
                payment_date: status === 'paid' ? new Date().toISOString() : null,
                notes: transport.notes || null,
                vehicles: vehicleBreakdown,
                created_at: transport.created_at,
                updated_at: transport.updated_at,
            })
        }

        return result
    })

    const filteredBills = computed(() => {
        let result = [...bills.value]

        // Filter by search query
        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(bill =>
                bill.customer_name.toLowerCase().includes(query) ||
                bill.from_location.toLowerCase().includes(query) ||
                (bill.to_location && bill.to_location.toLowerCase().includes(query)) ||
                String(bill.transport_id).includes(query) ||
                bill.status.toLowerCase().includes(query)
            )
        }

        // Filter by status
        if (statusFilter.value) {
            result = result.filter(bill => bill.status === statusFilter.value)
        }

        // Sort
        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'transport_id':
                    comparison = a.transport_id - b.transport_id
                    break
                case 'customer_name':
                    comparison = a.customer_name.localeCompare(b.customer_name)
                    break
                case 'from_location':
                    comparison = a.from_location.localeCompare(b.from_location)
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
    const totalUnpaidAmount = computed(() => {
        return bills.value
            .filter(b => b.status === 'unpaid')
            .reduce((sum, b) => sum + b.bill_amount, 0)
    })

    // ============= ACTIONS =============

    // Mark a bill as paid (updates the transport's customer_total_paid)
    const markBillAsPaid = async (transportId: number, _payload: {
        payment_date?: string
        notes?: string | null
    }): Promise<boolean> => {
        try {
            // Get the transport
            const transport = transportsStore.getTransportById(transportId)
            if (!transport) {
                return false
            }

            // Get the bill
            const bill = bills.value.find(b => b.transport_id === transportId)
            if (!bill) {
                return false
            }

            if (bill.status === 'paid') {
                return true
            }

            // Calculate new total paid (full amount)
            const newTotalPaid = bill.bill_amount

            // Update the transport's customer_total_paid
            const result = await transportsStore.updateCustomerPayment(transportId, {
                customer_total_paid: newTotalPaid
            })

            return result !== null
        } catch (error) {
            console.error('Error marking customer transport bill as paid:', error)
            return false
        }
    }

    // Cancel a bill (sets customer_charge to 0 on all vehicles in the transport)
    const cancelBill = async (transportId: number): Promise<boolean> => {
        try {
            const transport = transportsStore.getTransportById(transportId)
            if (!transport) {
                return false
            }

            // Get all vehicles for this transport
            const vehicles = vehiclesStore.getVehiclesByTransportId(transportId)

            // Set customer_charge to 0 on each vehicle
            let successCount = 0
            for (const vehicle of vehicles) {
                const result = await vehiclesStore.updateVehicle(vehicle.id, {
                    customer_charge: 0,
                })
                if (result) {
                    successCount++
                }
            }

            if (successCount === vehicles.length) {
                return true
            } else {
                console.warn(`Cancelled ${successCount}/${vehicles.length} vehicles`)
                return successCount > 0
            }
        } catch (error) {
            console.error('Error cancelling customer transport bill:', error)
            return false
        }
    }

    // ============= SORT =============
    const setSort = (field: CustomerTransportBillSortField) => {
        if (sortField.value === field) {
            sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
        } else {
            sortField.value = field
            sortDirection.value = 'asc'
        }
    }

    // ============= SEARCH =============
    const setSearchQuery = (query: string) => {
        searchQuery.value = query
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    const setStatusFilter = (status: 'unpaid' | 'paid' | 'cancelled' | '') => {
        statusFilter.value = status
    }

    // ============= UTILITIES =============

    const getBillById = (id: number): CustomerTransportBill | undefined => {
        return bills.value.find(b => b.id === id)
    }

    const getBillsByTransportId = (transportId: number): CustomerTransportBill | undefined => {
        return bills.value.find(b => b.transport_id === transportId)
    }

    const getBillsByCustomerId = (customerId: number): CustomerTransportBill[] => {
        return bills.value.filter(b => b.customer_id === customerId)
    }

    const getStatusBadgeClass = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'border-(--color-yellow) text-(--color-yellow)'
            case 'paid':
                return 'border-(--color-green) text-(--color-green)'
            case 'cancelled':
                return 'border-(--color-red) text-(--color-red)'
            default:
                return 'border-(--color-border) text-(--color-text-secondary)'
        }
    }

    const getStatusDotClass = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'bg-(--color-yellow)'
            case 'paid':
                return 'bg-(--color-green)'
            case 'cancelled':
                return 'bg-(--color-red)'
            default:
                return 'bg-(--color-text-secondary)'
        }
    }

    const getStatusLabel = (status: string): string => {
        switch (status) {
            case 'unpaid':
                return 'Unpaid'
            case 'paid':
                return 'Paid'
            case 'cancelled':
                return 'Cancelled'
            default:
                return status
        }
    }

    const formatBillDate = (dateStr: string): string => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        })
    }

    return {
        // State
        searchQuery,
        statusFilter,
        sortField,
        sortDirection,

        // Computed
        bills,
        filteredBills,
        totalBills,
        totalUnpaid,
        totalPaid,
        totalAmount,
        totalUnpaidAmount,

        // Actions
        markBillAsPaid,
        cancelBill,

        // Sort
        setSort,

        // Search
        setSearchQuery,
        clearSearch,
        setStatusFilter,

        // Utilities
        getBillById,
        getBillsByTransportId,
        getBillsByCustomerId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        formatBillDate,
    }
})