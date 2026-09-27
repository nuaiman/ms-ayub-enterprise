// src/stores/brokerVehicleBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useVehiclesStore } from './vehicles'
import { useBrokersStore } from './brokers'
import type {
    BrokerVehicleBill,
    BrokerVehicleBillSortField,
    SortDirection,
} from '@/types/brokerVehicleBill'
import { push } from 'notivue'

export const useBrokerVehicleBillsStore = defineStore('brokerVehicleBills', () => {
    const vehiclesStore = useVehiclesStore()
    const brokersStore = useBrokersStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<BrokerVehicleBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    const bills = computed<BrokerVehicleBill[]>(() => {
        const result: BrokerVehicleBill[] = []

        for (const vehicle of vehiclesStore.vehicles) {
            const billAmount = (vehicle.joma_cost || 0) + (vehicle.vehicle_cost || 0)
            if (billAmount === 0) continue

            const brokerName = vehicle.broker_id
                ? brokersStore.getBrokerName(vehicle.broker_id)
                : 'No Broker'

            const paidAmount = vehicle.total_paid_to_broker || 0
            const status: 'unpaid' | 'paid' = paidAmount >= billAmount ? 'paid' : 'unpaid'

            result.push({
                id: vehicle.id,
                vehicle_number: vehicle.vehicle_number,
                transport_id: vehicle.transport_id,
                broker_id: vehicle.broker_id,
                broker_name: brokerName,
                joma_cost: vehicle.joma_cost || 0,
                vehicle_cost: vehicle.vehicle_cost || 0,
                bill_amount: billAmount,
                paid_amount: paidAmount,
                status,
                payment_date: status === 'paid' ? new Date().toISOString() : null,
                notes: vehicle.notes || null,
                created_at: vehicle.created_at,
                updated_at: vehicle.updated_at,
            })
        }

        return result
    })

    const filteredBills = computed(() => {
        let result = [...bills.value]

        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(bill =>
                bill.vehicle_number.toLowerCase().includes(query) ||
                bill.broker_name.toLowerCase().includes(query) ||
                String(bill.transport_id).includes(query) ||
                bill.status.toLowerCase().includes(query)
            )
        }

        if (statusFilter.value) {
            result = result.filter(bill => bill.status === statusFilter.value)
        }

        result.sort((a, b) => {
            let comparison = 0
            switch (sortField.value) {
                case 'vehicle_number':
                    comparison = a.vehicle_number.localeCompare(b.vehicle_number)
                    break
                case 'broker_name':
                    comparison = a.broker_name.localeCompare(b.broker_name)
                    break
                case 'transport_id':
                    comparison = a.transport_id - b.transport_id
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
            .reduce((sum, b) => sum + (b.bill_amount - b.paid_amount), 0)
    )

    const markBillAsPaid = async (
        vehicleId: number,
        payload: { amount: number; payment_date?: string; notes?: string | null }
    ): Promise<boolean> => {
        try {
            const vehicle = vehiclesStore.getVehicleById(vehicleId)
            if (!vehicle) {
                push.error('Vehicle not found')
                return false
            }

            const bill = bills.value.find(b => b.id === vehicleId)
            if (!bill) {
                push.error('Bill not found')
                return false
            }

            if (bill.status === 'paid') {
                push.info('Bill is already fully paid')
                return true
            }

            const newTotalPaid = (vehicle.total_paid_to_broker || 0) + payload.amount

            if (newTotalPaid > bill.bill_amount) {
                const remaining = bill.bill_amount - bill.paid_amount
                push.error(`Payment amount exceeds remaining balance of ${remaining.toFixed(2)}`)
                return false
            }

            const result = await vehiclesStore.updateVehicleBrokerPayment(vehicleId, {
                total_paid_to_broker: newTotalPaid,
            })

            if (result) {
                await vehiclesStore.fetchVehicles()
                push.success('Payment recorded successfully')
                return true
            }

            return false
        } catch (error) {
            console.error('Error marking broker vehicle bill as paid:', error)
            push.error('Failed to record payment')
            return false
        }
    }

    const cancelBill = async (vehicleId: number): Promise<boolean> => {
        try {
            const vehicle = vehiclesStore.getVehicleById(vehicleId)
            if (!vehicle) {
                push.error('Vehicle not found')
                return false
            }

            const result = await vehiclesStore.updateVehicle(vehicleId, {
                joma_cost: 0,
                vehicle_cost: 0,
            })

            if (result) {
                await vehiclesStore.fetchVehicles()
                push.success('Bill cancelled successfully')
                return true
            }

            return false
        } catch (error) {
            console.error('Error cancelling broker vehicle bill:', error)
            push.error('Failed to cancel bill')
            return false
        }
    }

    const setSort = (field: BrokerVehicleBillSortField) => {
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

    const getBillById = (id: number): BrokerVehicleBill | undefined => {
        return bills.value.find(b => b.id === id)
    }

    const getBillsByBrokerId = (brokerId: number): BrokerVehicleBill[] => {
        return bills.value.filter(b => b.broker_id === brokerId)
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
        getBillsByBrokerId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
    }
})