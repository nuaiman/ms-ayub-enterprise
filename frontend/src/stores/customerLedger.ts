// src/stores/customerLedger.ts

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useCustomersStore } from './customers'
import { useLotsStore } from './lots'
import { useStoresStore } from './stores'
import { useGodownsStore } from './godowns'
import { useMajhisStore } from './majhis'
import { useDeliveriesStore } from './deliveries'
import { useDamagesStore } from './damages'
import { useCustomerStoreBillsStore } from './customerStoreBills'
import { useCustomerDeliveryBillsStore } from './customerDeliveryBills'
import { useCustomerAdditionalBillsStore } from './customerAdditionalBills'
import { useInvoicesStore } from './invoices'
import { useLotTransfersStore } from './lotTransfers'
import api from '@/utils/axios'
import type { ApiResponse } from '@/types/api'

export type LedgerEventType =
    | 'customer_created'
    | 'lot_created'
    | 'store_created'
    | 'delivery_created'
    | 'delivery_item_added'
    | 'damage_recorded'
    | 'lot_transferred_out'
    | 'lot_transferred_in'
    | 'customer_store_bill'
    | 'customer_delivery_bill'
    | 'customer_additional_bill'
    | 'invoice_created'
    | 'payment_received'

export type LedgerAmountKind = 'debit' | 'credit' | 'neutral'

export interface LedgerEvent {
    id: string
    type: LedgerEventType
    date: string
    title: string
    description: string
    reference?: string
    amount?: number
    amountLabel?: string
    amountKind?: LedgerAmountKind
    icon: string
    meta?: Record<string, string | number | null | undefined>
}

interface LedgerPayment {
    id: number
    bill_type: string
    bill_id: number
    amount: number
    payment_date: string
    payment_method?: string | null
    reference_number?: string | null
    notes?: string | null
}

const CUSTOMER_VISIBLE_TYPES: LedgerEventType[] = [
    'customer_created',
    'lot_created',
    'store_created',
    'delivery_created',
    'delivery_item_added',
    'damage_recorded',
    'lot_transferred_out',
    'lot_transferred_in',
    'customer_store_bill',
    'customer_delivery_bill',
    'customer_additional_bill',
    'invoice_created',
    'payment_received',
]

const CUSTOMER_HIDDEN_META_KEYS: string[] = ['user_id', 'internal_notes', 'rate']

const ICONS: Record<LedgerEventType, string> = {
    customer_created: '👤',
    lot_created: '🏷️',
    store_created: '🏢',
    delivery_created: '🚚',
    delivery_item_added: '📦',
    damage_recorded: '🔴',
    lot_transferred_out: '➡️',
    lot_transferred_in: '⬅️',
    customer_store_bill: '🧾',
    customer_delivery_bill: '🧾',
    customer_additional_bill: '🧾',
    invoice_created: '📄',
    payment_received: '💵',
}

export const useCustomerLedgerStore = defineStore('customerLedger', () => {
    const customersStore = useCustomersStore()
    const lotsStore = useLotsStore()
    const storesStore = useStoresStore()
    const godownsStore = useGodownsStore()
    const majhisStore = useMajhisStore()
    const deliveriesStore = useDeliveriesStore()
    const damagesStore = useDamagesStore()
    const customerStoreBillsStore = useCustomerStoreBillsStore()
    const customerDeliveryBillsStore = useCustomerDeliveryBillsStore()
    const customerAdditionalBillsStore = useCustomerAdditionalBillsStore()
    const invoicesStore = useInvoicesStore()
    const lotTransfersStore = useLotTransfersStore()

    const selectedCustomerId = ref<number | null>(null)
    const typeFilter = ref<LedgerEventType[]>([])
    const dateFrom = ref<string>('')
    const dateTo = ref<string>('')
    const searchQuery = ref('')
    const customerView = ref(false)
    const isLoading = ref(false)

    const allEvents = ref<LedgerEvent[]>([])

    // -------------------------------------------------------------------------
    // MANUAL SELECTION STATE
    // -------------------------------------------------------------------------
    // When selectedEventIds is non-empty, print/export uses only the selected
    // events. When it's empty, everything visible (post-filter) is used.
    const selectedEventIds = ref<Set<string>>(new Set())

    const paymentsCache = ref<Record<string, LedgerPayment[]>>({})

    // =========================================================================
    // HELPERS
    // =========================================================================

    const lotLabel = (lotId: number): string => {
        const lot = lotsStore.getLotById(lotId)
        if (!lot) return `Lot #${lotId}`
        return `${lot.product_name} · Lot ${lot.lot_number}`
    }

    const lotLabelFromStore = (storeId: number): string => {
        const store = storesStore.getStoreById(storeId)
        if (!store) return `Store #${storeId}`
        return lotLabel(store.lot_id)
    }

    const godownName = (godownId: number): string => godownsStore.getGodownName(godownId)
    const majhiName = (majhiId: number): string => majhisStore.getMajhiName(majhiId)

    const fetchPaymentsForBill = async (
        billType: string,
        billId: number
    ): Promise<LedgerPayment[]> => {
        const key = `${billType}:${billId}`
        const cached = paymentsCache.value[key]
        if (cached) return cached

        const routeMap: Record<string, string> = {
            customer_store: `/customer-store-bills/${billId}/payments`,
            customer_delivery: `/customer-delivery-bills/${billId}/payments`,
            customer_additional: `/customer-additional-bills/${billId}/payments`,
        }
        const route = routeMap[billType]
        if (!route) {
            paymentsCache.value[key] = []
            return []
        }

        try {
            const res = await api.get<
                ApiResponse<{
                    payments: Array<{
                        id: number
                        bill_type: string
                        bill_id: number
                        amount: number
                        payment_date: string
                        payment_method?: string | null
                        reference_number?: string | null
                        notes?: string | null
                    }>
                }>
            >(route)

            if (!res.data.success || !res.data.data?.payments) {
                paymentsCache.value[key] = []
                return []
            }

            const mapped: LedgerPayment[] = res.data.data.payments.map((p) => ({
                id: p.id,
                bill_type: billType,
                bill_id: billId,
                amount: p.amount,
                payment_date: p.payment_date,
                payment_method: p.payment_method ?? null,
                reference_number: p.reference_number ?? null,
                notes: p.notes ?? null,
            }))

            paymentsCache.value[key] = mapped
            return mapped
        } catch {
            paymentsCache.value[key] = []
            return []
        }
    }

    const paymentToEvent = (p: LedgerPayment, billLabel: string): LedgerEvent => ({
        id: `payment_${p.id}`,
        type: 'payment_received',
        date: p.payment_date,
        title: 'Payment received',
        description: billLabel,
        reference: p.reference_number || `Payment #${p.id}`,
        amount: p.amount,
        amountLabel: 'Paid',
        amountKind: 'credit',
        icon: ICONS.payment_received,
        meta: {
            method: p.payment_method,
            reference_number: p.reference_number,
            notes: p.notes,
        },
    })

    // =========================================================================
    // BUILD EVENTS
    // =========================================================================

    const buildEvents = async (customerId: number): Promise<LedgerEvent[]> => {
        const events: LedgerEvent[] = []

        const customer = customersStore.getCustomerById(customerId)
        if (!customer) return events

        // 1. Customer created
        events.push({
            id: `customer_${customer.id}`,
            type: 'customer_created',
            date: customer.created_at,
            title: 'Customer created',
            description: customer.company_name || customer.contact_person || `Customer #${customer.id}`,
            reference: customer.phone,
            amountKind: 'neutral',
            icon: ICONS.customer_created,
            meta: {
                phone: customer.phone,
                email: customer.email,
                address: customer.address,
                notes: customer.notes,
            },
        })

        // 2. Lots
        const customerLots = lotsStore.lots.filter((l) => l.customer_id === customerId)
        const customerLotIds = new Set(customerLots.map((l) => l.id))

        for (const lot of customerLots) {
            events.push({
                id: `lot_${lot.id}`,
                type: 'lot_created',
                date: lot.created_at,
                title: `Lot ${lot.lot_number} created`,
                description: lot.product_name,
                reference: `Lot #${lot.id}`,
                amountKind: 'neutral',
                icon: ICONS.lot_created,
                meta: {
                    product_name: lot.product_name,
                    weight_unit: lot.weight_unit,
                    quantity_unit: lot.quantity_unit,
                },
            })
        }

        // 3. Stores
        const customerStores = storesStore.stores.filter((s) => customerLotIds.has(s.lot_id))
        const customerStoreIds = new Set(customerStores.map((s) => s.id))

        for (const store of customerStores) {
            const lot = lotsStore.getLotById(store.lot_id)
            events.push({
                id: `store_${store.id}`,
                type: 'store_created',
                date: store.created_at,
                title: `Store #${store.id} created`,
                description: `${lotLabel(store.lot_id)} @ ${godownName(store.godown_id)}`,
                reference: `Store #${store.id}`,
                amountKind: 'neutral',
                icon: ICONS.store_created,
                meta: {
                    godown: godownName(store.godown_id),
                    quantity: `${store.quantity} ${lot?.quantity_unit ?? ''}`.trim(),
                    weight: `${store.weight} ${lot?.weight_unit ?? ''}`.trim(),
                    start_date: store.start_date,
                    is_active: store.is_active ? 'Active' : 'Inactive',
                },
            })
        }

        // 4. Deliveries
        const customerDeliveries = deliveriesStore.deliveries.filter((d) => d.customer_id === customerId)

        for (const delivery of customerDeliveries) {
            events.push({
                id: `delivery_${delivery.id}`,
                type: 'delivery_created',
                date: delivery.delivery_date,
                title: `Delivery #${delivery.id} created`,
                description:
                    delivery.from_location || delivery.to_location
                        ? `${delivery.from_location || '—'} → ${delivery.to_location || '—'}`
                        : `Delivery #${delivery.id}`,
                reference: `Delivery #${delivery.id}`,
                amountKind: 'neutral',
                icon: ICONS.delivery_created,
                meta: {
                    receiver_name: delivery.receiver_name,
                    receiver_phone: delivery.receiver_phone,
                    from_location: delivery.from_location,
                    to_location: delivery.to_location,
                },
            })
        }

        // 5. Delivery items
        const itemsByDelivery = deliveriesStore.itemsByDelivery
        for (const delivery of customerDeliveries) {
            const items = itemsByDelivery[delivery.id] ?? []
            for (const item of items) {
                const store = storesStore.getStoreById(item.store_id)
                const lot = store ? lotsStore.getLotById(store.lot_id) : null
                events.push({
                    id: `delivery_item_${item.id}`,
                    type: 'delivery_item_added',
                    date: item.created_at,
                    title: 'Delivery item added',
                    description: lot
                        ? `${lot.product_name} · Lot ${lot.lot_number}`
                        : `Store #${item.store_id}`,
                    reference: `Delivery #${delivery.id}`,
                    amountKind: 'neutral',
                    icon: ICONS.delivery_item_added,
                    meta: {
                        store: `Store #${item.store_id}`,
                        majhi: majhiName(item.majhi_id),
                        quantity: `${item.quantity} ${lot?.quantity_unit ?? ''}`.trim(),
                        weight: `${item.weight} ${lot?.weight_unit ?? ''}`.trim(),
                        vehicle_number: item.vehicle_number,
                        driver_number: item.driver_number,
                    },
                })
            }
        }

        // 6. Damages
        const customerDamages = damagesStore.damages.filter((d) => customerStoreIds.has(d.store_id))
        for (const damage of customerDamages) {
            events.push({
                id: `damage_${damage.id}`,
                type: 'damage_recorded',
                date: damage.damage_date,
                title: 'Damage recorded',
                description: damage.reason,
                reference: `Damage #${damage.id}`,
                amount: damage.amount,
                amountLabel: 'Damage',
                amountKind: 'debit',
                icon: ICONS.damage_recorded,
                meta: {
                    store: `Store #${damage.store_id}`,
                    quantity: `${damage.quantity} ${damage.quantity_unit}`.trim(),
                    weight: `${damage.weight} ${damage.weight_unit}`.trim(),
                    notes: damage.notes,
                },
            })
        }

        // 7. Lot transfers
        for (const t of lotTransfersStore.transfers) {
            if (t.from_customer_id === customerId) {
                events.push({
                    id: `lot_transfer_out_${t.id}`,
                    type: 'lot_transferred_out',
                    date: t.transferred_at,
                    title: 'Lot transferred out',
                    description: `${lotLabel(t.lot_id)} → ${customersStore.getCustomerName(t.to_customer_id)}`,
                    reference: `Transfer #${t.id}`,
                    amountKind: 'neutral',
                    icon: ICONS.lot_transferred_out,
                    meta: {
                        to_customer: customersStore.getCustomerName(t.to_customer_id),
                        notes: t.notes,
                    },
                })
            }
            if (t.to_customer_id === customerId) {
                events.push({
                    id: `lot_transfer_in_${t.id}`,
                    type: 'lot_transferred_in',
                    date: t.transferred_at,
                    title: 'Lot transferred in',
                    description: `${lotLabel(t.lot_id)} ← ${customersStore.getCustomerName(t.from_customer_id)}`,
                    reference: `Transfer #${t.id}`,
                    amountKind: 'neutral',
                    icon: ICONS.lot_transferred_in,
                    meta: {
                        from_customer: customersStore.getCustomerName(t.from_customer_id),
                        notes: t.notes,
                    },
                })
            }
        }

        // 8. Customer store bills + payments
        const csBills = customerStoreBillsStore.bills.filter((b) => b.customer_id === customerId)
        for (const b of csBills) {
            events.push({
                id: `customer_store_bill_${b.id}`,
                type: 'customer_store_bill',
                date: b.created_at,
                title: 'Customer store bill',
                description: `${lotLabelFromStore(b.store_id)} · ${b.month_year}`,
                reference: `Bill #${b.id}`,
                amount: b.total_amount,
                amountLabel: 'Billed',
                amountKind: 'debit',
                icon: ICONS.customer_store_bill,
                meta: {
                    month: b.month_year,
                    store: `Store #${b.store_id}`,
                    bill_type: b.bill_type,
                    quantity: `${b.quantity_at_billing} ${b.quantity_unit_at_billing ?? ''}`.trim(),
                    weight: `${b.weight_at_billing} ${b.weight_unit_at_billing ?? ''}`.trim(),
                    rate: b.rate,
                },
            })

            const payments = await fetchPaymentsForBill('customer_store', b.id)
            for (const p of payments) {
                events.push(paymentToEvent(p, `Customer store bill #${b.id}`))
            }
        }

        // 9. Customer delivery bills + payments
        const cdBills = customerDeliveryBillsStore.bills.filter((b) => b.customer_id === customerId)
        for (const b of cdBills) {
            events.push({
                id: `customer_delivery_bill_${b.id}`,
                type: 'customer_delivery_bill',
                date: b.created_at,
                title: 'Customer delivery bill',
                description: `Delivery item #${b.delivery_item_id}`,
                reference: `Bill #${b.id}`,
                amount: b.total_amount,
                amountLabel: 'Billed',
                amountKind: 'debit',
                icon: ICONS.customer_delivery_bill,
                meta: {
                    delivery_item: `Item #${b.delivery_item_id}`,
                    bill_type: b.bill_type,
                    quantity: `${b.quantity_at_billing} ${b.quantity_unit_at_billing ?? ''}`.trim(),
                    weight: `${b.weight_at_billing} ${b.weight_unit_at_billing ?? ''}`.trim(),
                    rate: b.rate,
                },
            })

            const payments = await fetchPaymentsForBill('customer_delivery', b.id)
            for (const p of payments) {
                events.push(paymentToEvent(p, `Customer delivery bill #${b.id}`))
            }
        }

        // 10. Customer additional bills + payments
        const caBills = customerAdditionalBillsStore.bills.filter((b) => b.customer_id === customerId)
        for (const b of caBills) {
            events.push({
                id: `customer_additional_bill_${b.id}`,
                type: 'customer_additional_bill',
                date: b.created_at,
                title: 'Customer additional bill',
                description: b.description,
                reference: `Bill #${b.id}`,
                amount: b.amount,
                amountLabel: 'Billed',
                amountKind: 'debit',
                icon: ICONS.customer_additional_bill,
                meta: {
                    description: b.description,
                },
            })

            const payments = await fetchPaymentsForBill('customer_additional', b.id)
            for (const p of payments) {
                events.push(paymentToEvent(p, `Customer additional bill #${b.id}`))
            }
        }

        // 11. Invoices
        const invoices = invoicesStore.invoices.filter(
            (i) => i.entity_type === 'customer' && i.entity_id === customerId
        )
        for (const inv of invoices) {
            events.push({
                id: `invoice_${inv.id}`,
                type: 'invoice_created',
                date: inv.created_at,
                title: `Invoice INV-${inv.id} created`,
                description: inv.notes || 'Invoice issued',
                reference: `INV-${inv.id}`,
                amount: inv.total,
                amountLabel: 'Invoiced',
                amountKind: 'debit',
                icon: ICONS.invoice_created,
                meta: {
                    subtotal: inv.subtotal,
                    discount: inv.discount_amount,
                    total: inv.total,
                    notes: inv.notes,
                },
            })
        }

        return events
    }

    // =========================================================================
    // SANITIZE (customer view)
    // =========================================================================

    const sanitizeEvent = (event: LedgerEvent): LedgerEvent | null => {
        if (!CUSTOMER_VISIBLE_TYPES.includes(event.type)) return null

        let meta = event.meta
        if (meta) {
            const cleaned: Record<string, string | number | null | undefined> = {}
            for (const [k, v] of Object.entries(meta)) {
                if (!CUSTOMER_HIDDEN_META_KEYS.includes(k)) cleaned[k] = v
            }
            meta = Object.keys(cleaned).length > 0 ? cleaned : undefined
        }

        return { ...event, meta }
    }

    // =========================================================================
    // REACTIVE SELECTORS
    // =========================================================================

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
            result = result.filter((e) => allowed.has(e.type))
        }

        if (dateFrom.value) {
            const from = new Date(dateFrom.value + 'T00:00:00').getTime()
            result = result.filter((e) => new Date(e.date).getTime() >= from)
        }
        if (dateTo.value) {
            const to = new Date(dateTo.value + 'T23:59:59').getTime()
            result = result.filter((e) => new Date(e.date).getTime() <= to)
        }

        if (searchQuery.value) {
            const q = searchQuery.value.toLowerCase()
            result = result.filter(
                (e) =>
                    e.title.toLowerCase().includes(q) ||
                    e.description.toLowerCase().includes(q) ||
                    (e.reference && e.reference.toLowerCase().includes(q))
            )
        }

        result.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())

        return result
    })

    // -------------------------------------------------------------------------
    // SELECTION-AWARE DERIVATIONS
    // -------------------------------------------------------------------------

    const hasSelection = computed(() => selectedEventIds.value.size > 0)

    /**
     * Returns the set of events that should actually be considered "in scope"
     * for summary + print/export. If the user made a manual selection, we honor
     * that (intersected with the current filters so a hidden item can't sneak
     * back in via the selection set). Otherwise we fall back to filteredEvents.
     */
    const activeEvents = computed<LedgerEvent[]>(() => {
        const filtered = filteredEvents.value
        if (selectedEventIds.value.size === 0) return filtered

        const selected = selectedEventIds.value
        return filtered.filter((e) => selected.has(e.id))
    })

    const selectedCount = computed(() => {
        // Count only items in the current filtered view
        const filtered = filteredEvents.value
        let n = 0
        for (const e of filtered) {
            if (selectedEventIds.value.has(e.id)) n++
        }
        return n
    })

    const allVisibleSelected = computed(() => {
        const filtered = filteredEvents.value
        if (filtered.length === 0) return false
        for (const e of filtered) {
            if (!selectedEventIds.value.has(e.id)) return false
        }
        return true
    })

    const someVisibleSelected = computed(() => {
        const filtered = filteredEvents.value
        if (filtered.length === 0) return false
        let any = false
        for (const e of filtered) {
            if (selectedEventIds.value.has(e.id)) {
                any = true
                break
            }
        }
        return any && !allVisibleSelected.value
    })

    // -------------------------------------------------------------------------
    // SUMMARY (now derived from activeEvents)
    // -------------------------------------------------------------------------

    const summary = computed(() => {
        let totalBilled = 0
        let totalPaid = 0

        for (const e of activeEvents.value) {
            if (e.amountKind === 'debit' && typeof e.amount === 'number') totalBilled += e.amount
            if (e.amountKind === 'credit' && typeof e.amount === 'number') totalPaid += e.amount
        }

        return {
            totalBilled,
            totalPaid,
            outstanding: totalBilled - totalPaid,
            eventCount: activeEvents.value.length,
        }
    })

    // =========================================================================
    // ACTIONS
    // =========================================================================

    const setCustomer = async (id: number | null) => {
        if (selectedCustomerId.value === id) return
        selectedCustomerId.value = id
        paymentsCache.value = {}
        selectedEventIds.value = new Set()
        if (!id) {
            allEvents.value = []
            return
        }
        isLoading.value = true
        try {
            allEvents.value = await buildEvents(id)
        } finally {
            isLoading.value = false
        }
    }

    const refresh = async () => {
        if (!selectedCustomerId.value) return
        paymentsCache.value = {}
        isLoading.value = true
        try {
            allEvents.value = await buildEvents(selectedCustomerId.value)
            // Drop any selection entries that no longer exist
            const valid = new Set(allEvents.value.map((e) => e.id))
            const next = new Set<string>()
            for (const id of selectedEventIds.value) {
                if (valid.has(id)) next.add(id)
            }
            selectedEventIds.value = next
        } finally {
            isLoading.value = false
        }
    }

    const toggleType = (type: LedgerEventType) => {
        const idx = typeFilter.value.indexOf(type)
        if (idx === -1) typeFilter.value.push(type)
        else typeFilter.value.splice(idx, 1)
    }

    const setDateRange = (from: string, to: string) => {
        dateFrom.value = from
        dateTo.value = to
    }

    const setSearchQuery = (q: string) => {
        searchQuery.value = q
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

    // -------------------------------------------------------------------------
    // SELECTION ACTIONS
    // -------------------------------------------------------------------------

    const isEventSelected = (id: string): boolean => selectedEventIds.value.has(id)

    const toggleEvent = (id: string) => {
        const next = new Set(selectedEventIds.value)
        if (next.has(id)) next.delete(id)
        else next.add(id)
        selectedEventIds.value = next
    }

    const selectAllVisible = () => {
        const next = new Set(selectedEventIds.value)
        for (const e of filteredEvents.value) {
            next.add(e.id)
        }
        selectedEventIds.value = next
    }

    const clearSelection = () => {
        selectedEventIds.value = new Set()
    }

    // =========================================================================
    // LABELS
    // =========================================================================

    const typeLabel = (type: LedgerEventType): string => {
        const labels: Record<LedgerEventType, string> = {
            customer_created: 'Customer',
            lot_created: 'Lot',
            store_created: 'Store',
            delivery_created: 'Delivery',
            delivery_item_added: 'Delivery Item',
            damage_recorded: 'Damage',
            lot_transferred_out: 'Lot Out',
            lot_transferred_in: 'Lot In',
            customer_store_bill: 'Store Bill',
            customer_delivery_bill: 'Delivery Bill',
            customer_additional_bill: 'Additional Bill',
            invoice_created: 'Invoice',
            payment_received: 'Payment',
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
        // Filters / view state
        selectedCustomerId,
        typeFilter,
        dateFrom,
        dateTo,
        searchQuery,
        customerView,
        isLoading,

        // Data
        allEvents,
        visibleEvents,
        filteredEvents,
        activeEvents,
        summary,

        // Selection
        selectedEventIds,
        hasSelection,
        selectedCount,
        allVisibleSelected,
        someVisibleSelected,
        isEventSelected,
        toggleEvent,
        selectAllVisible,
        clearSelection,

        // Actions
        setCustomer,
        refresh,
        toggleType,
        setDateRange,
        setSearchQuery,
        setCustomerView,
        resetFilters,

        // Labels
        typeLabel,
        formatEventDate,
    }
})