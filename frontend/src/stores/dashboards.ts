// src/stores/dashboards.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

import { useCustomersStore } from './customers'
import { useCustomerStorageBillsStore } from './customerStorageBills'
import { useCustomerLotBillsStore } from './customerLotBills'
import { useCustomerDeliveryBillsStore } from './customerDeliveryBills'
import { useCustomerTransportBillsStore } from './customerTransportBills'
import { useMajhiLotBillsStore } from './majhiLotBills'
import { useMajhiLoadingBillsStore } from './majhiLoadingBills'
import { useBrokerVehicleBillsStore } from './brokerVehicleBills'
import { useRentsStore } from './rents'
import { useLotsStore } from './lots'
import { useStoresStore } from './stores'
import { useGodownsStore } from './godowns'
import { useSalariesStore } from './salaries'
import { useExpensesStore } from './expenses'
import { useDeliveriesStore } from './deliveries'
import { useDeliveryItemsStore } from './deliveryItems'
import { useTransportsStore } from './transports'
import { useVehiclesStore } from './vehicles'
import { useMajhisStore } from './majhis'
import { useBrokersStore } from './brokers'

import type {
    DashboardMetrics,
    OutstandingCustomer,
    OutstandingByType,
    CustomerRevenue,
    MonthOption,
    DashboardData,
} from '@/types/dashboard'

export const useDashboardStore = defineStore('dashboard', () => {
    const selectedMonth = ref<string>(getCurrentMonth())
    const isLoading = ref(false)
    const isLoaded = ref(false)

    function getCurrentMonth(): string {
        const now = new Date()
        return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
    }

    function getMonthLabel(monthYear: string): string {
        const [year, month] = monthYear.split('-')
        if (!year || !month) return monthYear
        const date = new Date(parseInt(year), parseInt(month) - 1)
        return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    }

    function getInitials(name: string): string {
        if (!name) return '??'
        return name
            .split(' ')
            .map((word) => word[0])
            .join('')
            .toUpperCase()
            .slice(0, 2)
    }

    async function loadDashboardData() {
        if (isLoading.value) return
        if (isLoaded.value) return

        isLoading.value = true

        try {
            const customersStore = useCustomersStore()
            const lotsStore = useLotsStore()
            const storesStore = useStoresStore()
            const godownsStore = useGodownsStore()
            const salariesStore = useSalariesStore()
            const expensesStore = useExpensesStore()
            const deliveriesStore = useDeliveriesStore()
            const deliveryItemsStore = useDeliveryItemsStore()
            const transportsStore = useTransportsStore()
            const vehiclesStore = useVehiclesStore()
            const majhisStore = useMajhisStore()
            const brokersStore = useBrokersStore()
            const rentsStore = useRentsStore()

            await Promise.all([
                customersStore.fetchCustomers(),
                lotsStore.fetchLots(),
                storesStore.fetchStores(),
                godownsStore.fetchGodowns(),
                salariesStore.fetchSalaries(),
                expensesStore.fetchExpenses(),
                deliveriesStore.fetchDeliveries(),
                deliveryItemsStore.fetchDeliveryItems(),
                transportsStore.fetchTransports(),
                vehiclesStore.fetchVehicles(),
                majhisStore.fetchMajhis(),
                brokersStore.fetchBrokers(),
                rentsStore.fetchRents ? rentsStore.fetchRents() : Promise.resolve(),
                rentsStore.fetchCurrentMonthRents ? rentsStore.fetchCurrentMonthRents() : Promise.resolve(),
            ])

            isLoaded.value = true
        } catch (error) {
            console.error('Failed to load dashboard data:', error)
        } finally {
            isLoading.value = false
        }
    }

    function resetLoadedState() {
        isLoaded.value = false
    }

    const metrics = computed<DashboardMetrics>(() => {
        const storageStore = useCustomerStorageBillsStore()
        const lotStore = useCustomerLotBillsStore()
        const deliveryStore = useCustomerDeliveryBillsStore()
        const transportStore = useCustomerTransportBillsStore()
        const majhiLotStore = useMajhiLotBillsStore()
        const majhiLoadingStore = useMajhiLoadingBillsStore()
        const brokerStore = useBrokerVehicleBillsStore()
        const rentStore = useRentsStore()

        const storagePaid = storageStore.customerBillData?.reduce((sum, b) => sum + (b.total_paid || 0), 0) || 0
        const lotPaid = lotStore.bills?.reduce((sum, b) => sum + (b.paid_amount || 0), 0) || 0
        const deliveryPaid = deliveryStore.bills?.reduce((sum, b) => sum + (b.paid_amount || 0), 0) || 0
        const transportPaid = transportStore.bills?.reduce((sum, b) => sum + (b.paid_amount || 0), 0) || 0

        const totalReceived = storagePaid + lotPaid + deliveryPaid + transportPaid

        const currentMonth = new Date().getMonth()
        const currentYear = new Date().getFullYear()

        const monthlyStoragePaid = storageStore.customerBillData
            ?.filter((b) => b.billing_start && new Date(b.billing_start).getMonth() === currentMonth && new Date(b.billing_start).getFullYear() === currentYear)
            .reduce((sum, b) => sum + (b.total_paid || 0), 0) || 0

        const monthlyLotPaid = lotStore.bills
            ?.filter((b) => new Date(b.created_at).getMonth() === currentMonth && new Date(b.created_at).getFullYear() === currentYear)
            .reduce((sum, b) => sum + (b.paid_amount || 0), 0) || 0

        const monthlyDeliveryPaid = deliveryStore.bills
            ?.filter((b) => new Date(b.created_at).getMonth() === currentMonth && new Date(b.created_at).getFullYear() === currentYear)
            .reduce((sum, b) => sum + (b.paid_amount || 0), 0) || 0

        const monthlyTransportPaid = transportStore.bills
            ?.filter((b) => new Date(b.created_at).getMonth() === currentMonth && new Date(b.created_at).getFullYear() === currentYear)
            .reduce((sum, b) => sum + (b.paid_amount || 0), 0) || 0

        const monthlyReceived = monthlyStoragePaid + monthlyLotPaid + monthlyDeliveryPaid + monthlyTransportPaid

        const storageDue = storageStore.customerBillData?.reduce((sum, b) => sum + (b.outstanding || 0), 0) || 0
        const lotDue = lotStore.bills?.reduce((sum, b) => sum + ((b.bill_amount || 0) - (b.paid_amount || 0)), 0) || 0
        const deliveryDue = deliveryStore.bills?.reduce((sum, b) => sum + ((b.bill_amount || 0) - (b.paid_amount || 0)), 0) || 0
        const transportDue = transportStore.bills?.reduce((sum, b) => sum + ((b.bill_amount || 0) - (b.paid_amount || 0)), 0) || 0

        const totalDue = storageDue + lotDue + deliveryDue + transportDue

        const thirtyDaysAgo = new Date()
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)

        const overdueStorage = storageStore.customerBillData
            ?.filter((b) => b.billing_start && new Date(b.billing_start) < thirtyDaysAgo)
            .reduce((sum, b) => sum + (b.outstanding || 0), 0) || 0

        const overdueLot = lotStore.bills
            ?.filter((b) => new Date(b.created_at) < thirtyDaysAgo)
            .reduce((sum, b) => sum + ((b.bill_amount || 0) - (b.paid_amount || 0)), 0) || 0

        const overdueDelivery = deliveryStore.bills
            ?.filter((b) => new Date(b.created_at) < thirtyDaysAgo)
            .reduce((sum, b) => sum + ((b.bill_amount || 0) - (b.paid_amount || 0)), 0) || 0

        const overdueTransport = transportStore.bills
            ?.filter((b) => new Date(b.created_at) < thirtyDaysAgo)
            .reduce((sum, b) => sum + ((b.bill_amount || 0) - (b.paid_amount || 0)), 0) || 0

        const overdue = overdueStorage + overdueLot + overdueDelivery + overdueTransport

        const customersWithDue = new Set<number>()
        storageStore.customerBillData
            ?.filter((b) => (b.outstanding || 0) > 0 && b.customer_id)
            .forEach((b) => b.customer_id && customersWithDue.add(b.customer_id))
        lotStore.bills
            ?.filter((b) => ((b.bill_amount || 0) - (b.paid_amount || 0)) > 0 && b.customer_id)
            .forEach((b) => b.customer_id && customersWithDue.add(b.customer_id))
        deliveryStore.bills
            ?.filter((b) => ((b.bill_amount || 0) - (b.paid_amount || 0)) > 0 && b.customer_id)
            .forEach((b) => b.customer_id && customersWithDue.add(b.customer_id))
        transportStore.bills
            ?.filter((b) => ((b.bill_amount || 0) - (b.paid_amount || 0)) > 0 && b.customer_id)
            .forEach((b) => b.customer_id && customersWithDue.add(b.customer_id))

        const outstandingCustomers = customersWithDue.size

        const totalBilled = (storageStore.customerBillData?.reduce((sum, b) => sum + (b.total_billed || 0), 0) || 0) +
            (lotStore.bills?.reduce((sum, b) => sum + (b.bill_amount || 0), 0) || 0) +
            (deliveryStore.bills?.reduce((sum, b) => sum + (b.bill_amount || 0), 0) || 0) +
            (transportStore.bills?.reduce((sum, b) => sum + (b.bill_amount || 0), 0) || 0)

        const collectionRate = totalBilled > 0 ? Math.round((totalReceived / totalBilled) * 100) : 0

        const majhiUnpaidAmount = (majhiLotStore.totalUnpaidAmount || 0) + (majhiLoadingStore.totalUnpaidAmount || 0)

        const brokerUnpaidAmount = brokerStore.totalUnpaidAmount || 0

        const rentUnpaidAmount = rentStore.totalOutstanding || 0

        const customerMap = new Map<string, number>()
        const customersStore = useCustomersStore()

        storageStore.customerBillData
            ?.filter((b) => (b.outstanding || 0) > 0 && b.customer_id)
            .forEach((b) => {
                if (b.customer_id) {
                    const name = customersStore.getCustomerName(b.customer_id)
                    customerMap.set(name, (customerMap.get(name) || 0) + (b.outstanding || 0))
                }
            })

        lotStore.bills
            ?.filter((b) => ((b.bill_amount || 0) - (b.paid_amount || 0)) > 0 && b.customer_id)
            .forEach((b) => {
                if (b.customer_id) {
                    const name = customersStore.getCustomerName(b.customer_id)
                    customerMap.set(name, (customerMap.get(name) || 0) + ((b.bill_amount || 0) - (b.paid_amount || 0)))
                }
            })

        deliveryStore.bills
            ?.filter((b) => ((b.bill_amount || 0) - (b.paid_amount || 0)) > 0 && b.customer_id)
            .forEach((b) => {
                if (b.customer_id) {
                    const name = customersStore.getCustomerName(b.customer_id)
                    customerMap.set(name, (customerMap.get(name) || 0) + ((b.bill_amount || 0) - (b.paid_amount || 0)))
                }
            })

        transportStore.bills
            ?.filter((b) => ((b.bill_amount || 0) - (b.paid_amount || 0)) > 0 && b.customer_id)
            .forEach((b) => {
                if (b.customer_id) {
                    const name = customersStore.getCustomerName(b.customer_id)
                    customerMap.set(name, (customerMap.get(name) || 0) + ((b.bill_amount || 0) - (b.paid_amount || 0)))
                }
            })

        const outstandingByCustomer: OutstandingCustomer[] = Array.from(customerMap.entries())
            .map(([name, amount]) => ({ name, amount }))
            .sort((a, b) => b.amount - a.amount)
            .slice(0, 4)

        const totalOutstanding = storageDue + lotDue + deliveryDue + transportDue

        const outstandingByType: OutstandingByType[] = [
            {
                name: 'Storage',
                amount: storageDue,
                percentage: totalOutstanding > 0 ? Math.round((storageDue / totalOutstanding) * 100) : 0,
                color: '#3b82f6',
            },
            {
                name: 'Transport',
                amount: transportDue,
                percentage: totalOutstanding > 0 ? Math.round((transportDue / totalOutstanding) * 100) : 0,
                color: '#eab308',
            },
            {
                name: 'Unload',
                amount: lotDue,
                percentage: totalOutstanding > 0 ? Math.round((lotDue / totalOutstanding) * 100) : 0,
                color: '#8b5cf6',
            },
            {
                name: 'Delivery',
                amount: deliveryDue,
                percentage: totalOutstanding > 0 ? Math.round((deliveryDue / totalOutstanding) * 100) : 0,
                color: '#22c55e',
            },
        ].filter((item) => item.amount > 0)

        return {
            totalReceived,
            monthlyReceived,
            collectionRate,
            totalDue,
            overdue,
            outstandingCustomers,
            majhiUnpaidAmount,
            brokerUnpaidAmount,
            rentUnpaidAmount,
            outstandingByCustomer,
            outstandingByType,
        }
    })

    const customerRevenue = computed<CustomerRevenue[]>(() => {
        const customersStore = useCustomersStore()
        const storageStore = useCustomerStorageBillsStore()
        const lotStore = useCustomerLotBillsStore()
        const deliveryStore = useCustomerDeliveryBillsStore()
        const transportStore = useCustomerTransportBillsStore()

        const allCustomers = customersStore.customers || []

        const result: CustomerRevenue[] = []

        for (const customer of allCustomers) {
            const storageBills = storageStore.customerBillData?.filter(
                (b) => b.customer_id === customer.id
            ) || []
            const storageTotal = storageBills.reduce((sum, b) => sum + (b.total_billed || 0), 0)
            const storagePaid = storageBills.reduce((sum, b) => sum + (b.total_paid || 0), 0)

            const lotBills = lotStore.bills?.filter((b) => b.customer_id === customer.id) || []
            const lotTotal = lotBills.reduce((sum, b) => sum + (b.bill_amount || 0), 0)
            const lotPaid = lotBills.reduce((sum, b) => sum + (b.paid_amount || 0), 0)

            const deliveryBills = deliveryStore.bills?.filter((b) => b.customer_id === customer.id) || []
            const deliveryTotal = deliveryBills.reduce((sum, b) => sum + (b.bill_amount || 0), 0)
            const deliveryPaid = deliveryBills.reduce((sum, b) => sum + (b.paid_amount || 0), 0)

            const transportBills = transportStore.bills?.filter((b) => b.customer_id === customer.id) || []
            const transportTotal = transportBills.reduce((sum, b) => sum + (b.bill_amount || 0), 0)
            const transportPaid = transportBills.reduce((sum, b) => sum + (b.paid_amount || 0), 0)

            const total = storageTotal + lotTotal + deliveryTotal + transportTotal
            const paid = storagePaid + lotPaid + deliveryPaid + transportPaid
            const due = total - paid
            const paymentPercentage = total > 0 ? Math.round((paid / total) * 100) : 0

            if (total > 0) {
                const name = customer.company_name || customer.contact_person || `Customer #${customer.id}`
                const initials = getInitials(customer.company_name || customer.contact_person || '')

                result.push({
                    name,
                    initials,
                    phone: customer.phone || '',
                    storage: storageTotal,
                    unload: lotTotal,
                    delivery: deliveryTotal,
                    transport: transportTotal,
                    total,
                    paid,
                    due,
                    paymentPercentage,
                })
            }
        }

        return result.sort((a, b) => b.total - a.total)
    })

    const availableMonths = computed<MonthOption[]>(() => {
        const months = new Set<string>()
        const storageStore = useCustomerStorageBillsStore()

        storageStore.customerBillData?.forEach((b) => {
            if (b.billing_start) {
                const date = new Date(b.billing_start)
                const month = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
                months.add(month)
            }
        })

        return Array.from(months)
            .sort((a, b) => b.localeCompare(a))
            .map((month) => ({
                value: month,
                label: getMonthLabel(month),
            }))
    })

    function setSelectedMonth(month: string) {
        selectedMonth.value = month
    }

    function getDashboardData(): DashboardData {
        return {
            metrics: metrics.value,
            customerRevenue: customerRevenue.value,
            availableMonths: availableMonths.value,
        }
    }

    return {
        selectedMonth,
        isLoading,
        isLoaded,

        metrics,
        customerRevenue,
        availableMonths,

        loadDashboardData,
        resetLoadedState,
        setSelectedMonth,
        getDashboardData,
    }
})