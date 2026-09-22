// src/utils/storeReconstruction.ts

import { useDeliveryItemsStore } from '@/stores/deliveryItems'
import { useDamagesStore } from '@/stores/damages'
import type { Store } from '@/types/store'

export interface OriginalStock {
    quantity: number
    weight: number
}

/**
 * Reconstructs what a store originally held at creation time.
 *
 * The backend decrements store.quantity / store.weight when a delivery_item
 * is inserted (via SQLite trigger) and the damage handler decrements them
 * directly. This walks that same outflow history and adds it back to the
 * store's current values, yielding the original received amounts.
 *
 * If a delivery_item or damage row is later deleted, the backend already
 * restores the value, so the reconstruction stays correct.
 *
 * Assumption: stores are not manually edited via the Store form after
 * creation. If they are, the reconstruction will over-state the original.
 */
export const getOriginalStock = (store: Store): OriginalStock => {
    const deliveryItemsStore = useDeliveryItemsStore()
    const damagesStore = useDamagesStore()

    let quantity = store.quantity || 0
    let weight = store.weight || 0

    for (const di of deliveryItemsStore.deliveryItems) {
        if (di.store_id === store.id) {
            quantity += di.quantity || 0
            weight += di.weight || 0
        }
    }

    for (const d of damagesStore.damages) {
        if (d.store_id === store.id) {
            quantity += d.quantity || 0
            weight += d.weight || 0
        }
    }

    return { quantity, weight }
}

/**
 * Sum of original quantity/weight across an array of stores.
 */
export const getOriginalStockTotals = (stores: Store[]): OriginalStock => {
    let quantity = 0
    let weight = 0
    for (const store of stores) {
        const original = getOriginalStock(store)
        quantity += original.quantity
        weight += original.weight
    }
    return { quantity, weight }
}