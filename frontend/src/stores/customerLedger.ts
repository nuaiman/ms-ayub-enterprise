// src/stores/customerLedger.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useCustomersStore } from './customers'
import { useLotsStore } from './lots'
import { useStoresStore } from './stores'
import { useDeliveriesStore } from './deliveries'
import { useDeliveryItemsStore } from './deliveryItems'
import { useDamagesStore } from './damages'
import { useTransportsStore } from './transports'
import { useVehiclesStore } from './vehicles'
import { useGodownsStore } from './godowns'

export type LedgerEventType =
    | 'customer_created'
    | 'lot_created'
    | 'store_created'
    | 'delivery_created'
    | 'delivery_item_added'
    | 'damage_recorded'
    | 'transport_created'
    | 'vehicle_added'
    | 'storage_payment'
    | 'unload_payment'
    | 'delivery_payment'
    | 'transport_payment'

export interface LedgerEvent {
    id: string
    type: LedgerEventType
    date: string
    title: string
    description: string
    reference?: string
    amount?: number
    amountLabel?: string
    amountKind?: 'debit' | 'credit' | 'neutral'
    icon: string
    color: string
    meta?: Record<string, string | number | null | undefined>
}

const CUSTOMER_VISIBLE_TYPES: LedgerEventType[] = [
    'customer_created',
    'delivery_created',
    'delivery_item_added',
    'transport_created',
    'storage_payment',
    'unload_payment',
    'delivery_payment',
    'transport_payment',
]

const CUSTOMER_HIDDEN_META: Partial<Record<LedgerEventType, string[]>> = {
    delivery_item_added: ['loading_rate', 'vehicle_number', 'driver_number'],
    transport_created: ['office_commission_amount', 'customer_total_paid'],
    storage_payment: ['unload_rate'],
    unload_payment: ['unload_rate'],
}

export const useCustomerLedgerStore = defineStore('customerLedger', () => {
    const customersStore = useCustomersStore()
    const lotsStore = useLotsStore()
    const storesStore = useStoresStore()
    const deliveriesStore = useDeliveriesStore()
    const deliveryItemsStore = useDeliveryItemsStore()
    const damagesStore = useDamagesStore()
    const transportsStore = useTransportsStore()
    const vehiclesStore = useVehiclesStore()
    const godownsStore = useGodownsStore()

    const selectedCustomerId = ref<number | null>(null)
    const typeFilter = ref<LedgerEventType[]>([])
    const dateFrom = ref<string>('')
    const dateTo = ref<string>('')
    const searchQuery = ref('')
    const customerView = ref(false)

    const buildEvents = (customerId: number): LedgerEvent[] => {
        const events: LedgerEvent[] = []

        const customer = customersStore.getCustomerById(customerId)
        if (!customer) return events

        // --- Customer created ---
        events.push({
            id: `customer_${customer.id}`,
            type: 'customer_created',
            date: customer.created_at,
            title: 'Customer created',
            description: customer.company_name || customer.contact_person || `Customer #${customer.id}`,
            reference: customer.phone,
            icon: '👤',
            color: 'blue',
            meta: {
                phone: customer.phone,
                email: customer.email,
                address: customer.address,
                notes: customer.notes,
            },
        })

        // --- Lots of this customer ---
        const customerLots = lotsStore.lots.filter(l => l.customer_id === customerId)
        const customerLotIds = new Set(customerLots.map(l => l.id))

        for (const lot of customerLots) {
            events.push({
                id: `lot_${lot.id}`,
                type: 'lot_created',
                date: lot.created_at,
                title: `Lot #${lot.lot_number} created`,
                description: lotsStore.getLotDisplayName(lot),
                reference: `Lot #${lot.id}`,
                icon: '🏷️',
                color: 'indigo',
                meta: {
                    product_name: lot.product_name,
                    category: lot.category,
                    customer_charge_type: lot.customer_charge_type,
                    majhi_bill_type: lot.majhi_bill_type,
                    customer_storage_rate: lot.customer_storage_rate,
                    unload_rate: lot.unload_rate,
                    majhi_cut: lot.majhi_cut,
                },
            })
        }

        // --- Stores of those lots ---
        const customerStores = storesStore.stores.filter(s => customerLotIds.has(s.lot_id))

        for (const store of customerStores) {
            const lot = lotsStore.getLotById(store.lot_id)
            const godown = godownsStore.getGodownById(store.godown_id)
            events.push({
                id: `store_${store.id}`,
                type: 'store_created',
                date: store.created_at,
                title: `Store created`,
                description: `${lot ? `Lot #${lot.lot_number}` : `Lot #${store.lot_id}`} @ ${godown ? godown.name : `Godown #${store.godown_id}`}`,
                reference: `Store #${store.id}`,
                icon: '🏢',
                color: 'teal',
                meta: {
                    quantity: `${store.quantity} ${store.quantity_unit}`,
                    weight: `${store.weight} ${store.weight_unit}`,
                    store_bill_type: store.store_bill_type,
                    godown_cut: store.godown_cut,
                    is_active: store.is_active ? 'Active' : 'Inactive',
                },
            })
        }

        // --- Deliveries for this customer ---
        const customerDeliveries = deliveriesStore.deliveries.filter(d => d.customer_id === customerId)
        const customerDeliveryIds = new Set(customerDeliveries.map(d => d.id))

        for (const delivery of customerDeliveries) {
            events.push({
                id: `delivery_${delivery.id}`,
                type: 'delivery_created',
                date: delivery.delivery_date,
                title: `Delivery #${delivery.id} created`,
                description: delivery.from_location
                    ? `${delivery.from_location}${delivery.to_location ? ` → ${delivery.to_location}` : ''}`
                    : `Delivery #${delivery.id}`,
                reference: `Delivery #${delivery.id}`,
                icon: '🚚',
                color: 'green',
                meta: {
                    receiver_name: delivery.receiver_name,
                    receiver_phone: delivery.receiver_phone,
                    from_location: delivery.from_location,
                    to_location: delivery.to_location,
                },
            })
        }

        // --- Delivery items (via the customer's deliveries) ---
        const customerDeliveryItems = deliveryItemsStore.deliveryItems.filter(di =>
            customerDeliveryIds.has(di.delivery_id)
        )

        for (const di of customerDeliveryItems) {
            const lot = lotsStore.getLotById(di.lot_id)
            const lotName = lot ? lotsStore.getLotDisplayName(lot) : `Lot #${di.lot_id}`

            const billAmount =
                di.customer_charge_type === 'quantity'
                    ? (di.loading_rate || 0) * (di.quantity || 0)
                    : (di.loading_rate || 0) * (di.weight || 0)

            events.push({
                id: `delivery_item_${di.id}`,
                type: 'delivery_item_added',
                date: di.created_at,
                title: `Delivery item added`,
                description: `${lotName} from ${lot ? `Lot #${lot.lot_number}` : `Lot #${di.lot_id}`}`,
                reference: `Delivery #${di.delivery_id}`,
                amount: billAmount,
                amountLabel: 'Billed',
                amountKind: 'debit',
                icon: '📦',
                color: 'lime',
                meta: {
                    quantity: `${di.quantity} ${di.quantity_unit}`,
                    weight: `${di.weight} ${di.weight_unit}`,
                    loading_rate: di.loading_rate,
                    vehicle_number: di.vehicle_number,
                    driver_number: di.driver_number,
                    customer_paid_unload_amount: di.customer_paid_unload_amount,
                },
            })

            if ((di.customer_paid_unload_amount || 0) > 0) {
                events.push({
                    id: `delivery_payment_${di.id}`,
                    type: 'delivery_payment',
                    date: di.updated_at,
                    title: `Delivery bill payment`,
                    description: `Payment for delivery item #${di.id} (${lotName})`,
                    reference: `Delivery #${di.delivery_id}`,
                    amount: di.customer_paid_unload_amount || 0,
                    amountLabel: 'Paid',
                    amountKind: 'credit',
                    icon: '💵',
                    color: 'emerald',
                    meta: {
                        bill_amount: billAmount,
                        paid_amount: di.customer_paid_unload_amount,
                    },
                })
            }
        }

        // --- Damages for this customer's stores ---
        const customerStoreIds = new Set(customerStores.map(s => s.id))
        const customerDamages = damagesStore.damages.filter(d => customerStoreIds.has(d.store_id))

        for (const damage of customerDamages) {
            events.push({
                id: `damage_${damage.id}`,
                type: 'damage_recorded',
                date: damage.damage_date,
                title: `Damage recorded`,
                description: damage.reason,
                reference: `Damage #${damage.id}`,
                amount: damage.amount,
                amountLabel: 'Damage',
                amountKind: 'neutral',
                icon: '🔴',
                color: 'red',
                meta: {
                    quantity: `${damage.quantity} ${damage.quantity_unit}`,
                    weight: `${damage.weight} ${damage.weight_unit}`,
                    notes: damage.notes,
                },
            })
        }

        // --- Storage payments (from lots) ---
        for (const lot of customerLots) {
            const lotName = lotsStore.getLotDisplayName(lot)

            if ((lot.customer_last_paid_amount || 0) > 0) {
                events.push({
                    id: `storage_payment_${lot.id}`,
                    type: 'storage_payment',
                    date: lot.customer_last_paid_through || lot.updated_at,
                    title: `Storage payment`,
                    description: `Storage payment for ${lotName} (Lot #${lot.lot_number})`,
                    reference: `Lot #${lot.id}`,
                    amount: lot.customer_last_paid_amount || 0,
                    amountLabel: 'Paid',
                    amountKind: 'credit',
                    icon: '💰',
                    color: 'emerald',
                    meta: {
                        paid_through: lot.customer_last_paid_through,
                        total_paid: lot.customer_last_paid_amount,
                    },
                })
            }

            if ((lot.customer_paid_unload_amount || 0) > 0) {
                events.push({
                    id: `unload_payment_${lot.id}`,
                    type: 'unload_payment',
                    date: lot.updated_at,
                    title: `Unload payment`,
                    description: `Unload payment for ${lotName} (Lot #${lot.lot_number})`,
                    reference: `Lot #${lot.id}`,
                    amount: lot.customer_paid_unload_amount || 0,
                    amountLabel: 'Paid',
                    amountKind: 'credit',
                    icon: '💵',
                    color: 'emerald',
                    meta: {
                        unload_rate: lot.unload_rate,
                        total_paid: lot.customer_paid_unload_amount,
                    },
                })
            }
        }

        // --- Transports for this customer ---
        const customerTransports = transportsStore.transports.filter(t => t.customer_id === customerId)

        for (const transport of customerTransports) {
            events.push({
                id: `transport_${transport.id}`,
                type: 'transport_created',
                date: transport.transport_date,
                title: `Transport #${transport.id} created`,
                description: transport.from_location
                    ? `${transport.from_location}${transport.to_location ? ` → ${transport.to_location}` : ''}`
                    : `Transport #${transport.id}`,
                reference: `Transport #${transport.id}`,
                amount: transport.office_commission_amount,
                amountLabel: 'Commission',
                amountKind: 'neutral',
                icon: '🚛',
                color: 'orange',
                meta: {
                    vehicle_quantity: transport.vehicle_quantity,
                    delivery_type: transport.delivery_type,
                    customer_total_paid: transport.customer_total_paid,
                    office_commission_amount: transport.office_commission_amount,
                },
            })

            if ((transport.customer_total_paid || 0) > 0) {
                events.push({
                    id: `transport_payment_${transport.id}`,
                    type: 'transport_payment',
                    date: transport.updated_at,
                    title: `Transport payment`,
                    description: `Payment for Transport #${transport.id}`,
                    reference: `Transport #${transport.id}`,
                    amount: transport.customer_total_paid || 0,
                    amountLabel: 'Paid',
                    amountKind: 'credit',
                    icon: '💵',
                    color: 'emerald',
                    meta: {
                        total_paid: transport.customer_total_paid,
                    },
                })
            }
        }

        // --- Vehicles of the customer's transports ---
        const customerTransportIds = new Set(customerTransports.map(t => t.id))
        const customerVehicles = vehiclesStore.vehicles.filter(v => customerTransportIds.has(v.transport_id))

        for (const vehicle of customerVehicles) {
            events.push({
                id: `vehicle_${vehicle.id}`,
                type: 'vehicle_added',
                date: vehicle.created_at,
                title: `Vehicle added`,
                description: `${vehicle.vehicle_number}`,
                reference: `Transport #${vehicle.transport_id}`,
                amount: vehicle.customer_charge,
                amountLabel: 'Customer Charge',
                amountKind: 'debit',
                icon: '🚗',
                color: 'amber',
                meta: {
                    driver_name: vehicle.driver_name,
                    driver_phone: vehicle.driver_phone,
                    joma_cost: vehicle.joma_cost,
                    vehicle_cost: vehicle.vehicle_cost,
                    other_cost: vehicle.other_cost,
                    labour_cost: vehicle.labour_cost,
                    demarage_amount: vehicle.demarage_amount,
                    demarage_reason: vehicle.demarage_reason,
                    customer_charge: vehicle.customer_charge,
                },
            })
        }

        return events
    }

    const sanitizeEvent = (event: LedgerEvent): LedgerEvent | null => {
        if (!CUSTOMER_VISIBLE_TYPES.includes(event.type)) {
            return null
        }

        const hiddenKeys = CUSTOMER_HIDDEN_META[event.type] || []

        let meta: Record<string, string | number | null | undefined> | undefined
        if (event.meta) {
            const cleaned: Record<string, string | number | null | undefined> = {}
            for (const [k, v] of Object.entries(event.meta)) {
                if (!hiddenKeys.includes(k)) cleaned[k] = v
            }
            meta = Object.keys(cleaned).length > 0 ? cleaned : undefined
        }

        let amount = event.amount
        let amountLabel = event.amountLabel
        let amountKind = event.amountKind
        if (event.type === 'transport_created') {
            amount = undefined
            amountLabel = undefined
            amountKind = undefined
        }

        return {
            ...event,
            amount,
            amountLabel,
            amountKind,
            meta,
        }
    }

    const allEvents = computed<LedgerEvent[]>(() => {
        if (!selectedCustomerId.value) return []
        return buildEvents(selectedCustomerId.value)
    })

    const visibleEvents = computed<LedgerEvent[]>(() => {
        if (!customerView.value) return allEvents.value
        const out: LedgerEvent[] = []
        for (const e of allEvents.value) {
            const s = sanitizeEvent(e)
            if (s) out.push(s)
        }
        return out
    })

    const filteredEvents = computed<LedgerEvent[]>(() => {
        let result = [...visibleEvents.value]

        if (typeFilter.value.length > 0) {
            const allowed = new Set(typeFilter.value)
            result = result.filter(e => allowed.has(e.type))
        }

        if (dateFrom.value) {
            const from = new Date(dateFrom.value + 'T00:00:00').getTime()
            result = result.filter(e => new Date(e.date).getTime() >= from)
        }
        if (dateTo.value) {
            const to = new Date(dateTo.value + 'T23:59:59').getTime()
            result = result.filter(e => new Date(e.date).getTime() <= to)
        }

        if (searchQuery.value) {
            const q = searchQuery.value.toLowerCase()
            result = result.filter(e =>
                e.title.toLowerCase().includes(q) ||
                e.description.toLowerCase().includes(q) ||
                (e.reference && e.reference.toLowerCase().includes(q))
            )
        }

        result.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())

        return result
    })

    const summary = computed(() => {
        let totalBilled = 0
        let totalPaid = 0

        for (const e of visibleEvents.value) {
            if (e.amountKind === 'debit' && typeof e.amount === 'number') {
                totalBilled += e.amount
            }
            if (e.amountKind === 'credit' && typeof e.amount === 'number') {
                totalPaid += e.amount
            }
        }

        return {
            totalBilled,
            totalPaid,
            outstanding: totalBilled - totalPaid,
            eventCount: visibleEvents.value.length,
        }
    })

    const setCustomer = (id: number | null) => {
        selectedCustomerId.value = id
    }

    const setTypeFilter = (types: LedgerEventType[]) => {
        typeFilter.value = types
    }

    const toggleType = (type: LedgerEventType) => {
        const idx = typeFilter.value.indexOf(type)
        if (idx === -1) {
            typeFilter.value.push(type)
        } else {
            typeFilter.value.splice(idx, 1)
        }
    }

    const clearTypeFilter = () => {
        typeFilter.value = []
    }

    const setDateRange = (from: string, to: string) => {
        dateFrom.value = from
        dateTo.value = to
    }

    const clearDateRange = () => {
        dateFrom.value = ''
        dateTo.value = ''
    }

    const setSearchQuery = (q: string) => {
        searchQuery.value = q
    }

    const clearSearch = () => {
        searchQuery.value = ''
    }

    const setCustomerView = (on: boolean) => {
        customerView.value = on
    }

    const resetFilters = () => {
        typeFilter.value = []
        dateFrom.value = ''
        dateTo.value = ''
        searchQuery.value = ''
    }

    const typeLabel = (type: LedgerEventType): string => {
        const labels: Record<LedgerEventType, string> = {
            customer_created: 'Customer',
            lot_created: 'Lot',
            store_created: 'Store',
            delivery_created: 'Delivery',
            delivery_item_added: 'Delivery Item',
            damage_recorded: 'Damage',
            transport_created: 'Transport',
            vehicle_added: 'Vehicle',
            storage_payment: 'Storage Payment',
            unload_payment: 'Unload Payment',
            delivery_payment: 'Delivery Payment',
            transport_payment: 'Transport Payment',
        }
        return labels[type] || type
    }

    const formatEventDate = (dateStr: string): string => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        })
    }

    return {
        selectedCustomerId,
        typeFilter,
        dateFrom,
        dateTo,
        searchQuery,
        customerView,

        allEvents,
        visibleEvents,
        filteredEvents,
        summary,

        setCustomer,
        setTypeFilter,
        toggleType,
        clearTypeFilter,
        setDateRange,
        clearDateRange,
        setSearchQuery,
        clearSearch,
        setCustomerView,
        resetFilters,

        typeLabel,
        formatEventDate,
    }
})