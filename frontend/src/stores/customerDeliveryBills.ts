// src/stores/customerDeliveryBills.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useDeliveryItemsStore } from './deliveryItems'
import { useDeliveriesStore } from './deliveries'
import { useItemsStore } from './items'
import { useCustomersStore } from './customers'
import { useLotsStore } from './lots'
import type { CustomerDeliveryBill, CustomerDeliveryBillSortField, SortDirection } from '@/types/customerDeliveryBill'

export const useCustomerDeliveryBillsStore = defineStore('customerDeliveryBills', () => {
    const deliveryItemsStore = useDeliveryItemsStore()
    const deliveriesStore = useDeliveriesStore()
    const itemsStore = useItemsStore()
    const customersStore = useCustomersStore()
    const lotsStore = useLotsStore()

    const searchQuery = ref('')
    const statusFilter = ref<'unpaid' | 'paid' | 'cancelled' | ''>('')
    const sortField = ref<CustomerDeliveryBillSortField>('created_at')
    const sortDirection = ref<SortDirection>('desc')

    // ============= HELPERS =============

    // Calculate bill amount for a delivery item (loading_rate × quantity/weight)
    const calculateBillAmount = (item: any): number => {
        if (item.customer_charge_type === 'quantity') {
            return item.loading_rate * item.quantity
        } else {
            return item.loading_rate * item.weight
        }
    }

    // ============= COMPUTED =============

    // Generate bills for all delivery items with loading_rate > 0
    const bills = computed<CustomerDeliveryBill[]>(() => {
        // Get all delivery items with loading_rate > 0 (even if previously paid)
        const items = deliveryItemsStore.deliveryItems.filter(item => item.loading_rate > 0)

        const result: CustomerDeliveryBill[] = []

        for (const item of items) {
            // Get delivery info
            const delivery = deliveriesStore.getDeliveryById(item.delivery_id)
            if (!delivery) continue

            // Get customer info
            const customerId = delivery.customer_id || null
            const customerName = customerId ? customersStore.getCustomerName(customerId) : 'No Customer'

            // Get item and lot info
            const lot = lotsStore.getLotById(item.lot_id)
            const lotItem = lot ? itemsStore.getItemById(lot.item_id) : null

            // Calculate bill amount
            const billAmount = calculateBillAmount(item)

            // Skip if bill amount is 0
            if (billAmount === 0) continue

            // Get paid amount from delivery item
            const paidAmount = item.customer_paid_unload_amount || 0

            // Determine status - 'paid' only if fully paid, otherwise 'unpaid'
            const status = paidAmount >= billAmount ? 'paid' : 'unpaid'

            result.push({
                id: item.id,
                delivery_item_id: item.id,
                delivery_id: item.delivery_id,
                delivery_date: delivery.delivery_date,
                customer_id: customerId,
                customer_name: customerName,
                item_name: lotItem ? itemsStore.getItemDisplayName(lotItem) : `Item #${item.item_id}`,
                lot_id: item.lot_id,
                store_id: item.store_id,
                customer_charge_type: item.customer_charge_type,
                loading_rate: item.loading_rate,
                quantity: item.quantity,
                quantity_unit: item.quantity_unit,
                weight: item.weight,
                weight_unit: item.weight_unit,
                bill_amount: billAmount,
                paid_amount: paidAmount,
                status: status,
                payment_date: status === 'paid' ? new Date().toISOString() : null,
                notes: item.notes || null,
                created_at: item.created_at,
                updated_at: item.updated_at,
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
                bill.item_name.toLowerCase().includes(query) ||
                String(bill.delivery_id).includes(query) ||
                String(bill.lot_id).includes(query) ||
                bill.status.toLowerCase().includes(query) ||
                (bill.notes && bill.notes.toLowerCase().includes(query))
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
                case 'customer_name':
                    comparison = a.customer_name.localeCompare(b.customer_name)
                    break
                case 'item_name':
                    comparison = a.item_name.localeCompare(b.item_name)
                    break
                case 'bill_amount':
                    comparison = a.bill_amount - b.bill_amount
                    break
                case 'status':
                    comparison = a.status.localeCompare(b.status)
                    break
                case 'delivery_date':
                    comparison = new Date(a.delivery_date).getTime() - new Date(b.delivery_date).getTime()
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

    // Mark a bill as paid (updates the delivery item's customer_paid_unload_amount)
    const markBillAsPaid = async (deliveryItemId: number, payload: {
        payment_date?: string
        notes?: string | null
    }): Promise<boolean> => {
        try {
            // Get the delivery item
            const item = deliveryItemsStore.getDeliveryItemById(deliveryItemId)
            if (!item) {
                return false
            }

            // Get the bill
            const bill = bills.value.find(b => b.delivery_item_id === deliveryItemId)
            if (!bill) {
                return false
            }

            if (bill.status === 'paid') {
                return true
            }

            // Calculate new total paid (full amount)
            const newTotalPaid = bill.bill_amount

            // Update the delivery item's customer_paid_unload_amount
            const result = await deliveryItemsStore.updateDeliveryItemCustomerUnloadPayment(
                deliveryItemId,
                newTotalPaid
            )

            if (result) {
                // Also update the notes if provided
                if (payload.notes) {
                    await deliveryItemsStore.updateDeliveryItem(deliveryItemId, {
                        notes: payload.notes
                    })
                }
                return true
            }

            return false
        } catch (error) {
            console.error('Error marking customer delivery bill as paid:', error)
            return false
        }
    }

    // Cancel a bill (sets loading_rate to 0)
    const cancelBill = async (deliveryItemId: number): Promise<boolean> => {
        try {
            const item = deliveryItemsStore.getDeliveryItemById(deliveryItemId)
            if (!item) {
                return false
            }

            const result = await deliveryItemsStore.updateDeliveryItem(deliveryItemId, {
                loading_rate: 0,
                notes: `Customer delivery bill cancelled at ${new Date().toISOString()}${item.notes ? ` - ${item.notes}` : ''}`,
            })

            return result !== null
        } catch (error) {
            console.error('Error cancelling customer delivery bill:', error)
            return false
        }
    }

    // ============= SORT =============
    const setSort = (field: CustomerDeliveryBillSortField) => {
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

    const getBillById = (id: number): CustomerDeliveryBill | undefined => {
        return bills.value.find(b => b.id === id)
    }

    const getBillsByDeliveryId = (deliveryId: number): CustomerDeliveryBill[] => {
        return bills.value.filter(b => b.delivery_id === deliveryId)
    }

    const getBillsByCustomerId = (customerId: number): CustomerDeliveryBill[] => {
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

    const getChargeTypeLabel = (chargeType: 'weight' | 'quantity'): string => {
        return chargeType === 'weight' ? 'Weight' : 'Quantity'
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

    const formatDeliveryDate = (dateStr: string): string => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
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
        getBillsByDeliveryId,
        getBillsByCustomerId,
        getStatusBadgeClass,
        getStatusDotClass,
        getStatusLabel,
        getChargeTypeLabel,
        formatBillDate,
        formatDeliveryDate,
    }
})