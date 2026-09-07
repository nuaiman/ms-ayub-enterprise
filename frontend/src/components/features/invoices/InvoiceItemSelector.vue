<!-- src/components/features/invoices/InvoiceItemSelector.vue -->
<template>
    <div class="space-y-6">
        <!-- Header -->
        <div
            class="flex flex-col sm:flex-row items-start sm:items-center justify-between p-4 rounded-xl bg-(--color-muted-bg)/30 border border-(--color-border) gap-3">
            <div>
                <p class="text-sm font-medium text-(--color-text-primary)">
                    Customer: <span class="font-semibold">{{ customerName }}</span>
                </p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5">
                    Select unpaid or partially paid bills (max 12 items)
                </p>
            </div>
            <div class="flex items-center gap-4">
                <span class="text-sm font-medium"
                    :class="selectedCount > 12 ? 'text-(--color-red)' : 'text-(--color-blue)'">
                    {{ selectedCount }}/12 selected
                </span>
                <button @click="toggleAll"
                    class="text-xs font-medium text-(--color-blue) hover:text-(--color-blue)/70 transition-colors"
                    :disabled="availableItems.length === 0">
                    {{ allSelected ? 'Deselect All' : 'Select All' }}
                </button>
            </div>
        </div>

        <!-- Warning when at limit -->
        <div v-if="selectedCount >= 12"
            class="p-3 rounded-xl bg-(--color-yellow)/10 border border-(--color-yellow)/20 text-(--color-yellow) text-sm flex items-center gap-2">
            <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span>Maximum 12 items per invoice. Unselect some items to add more.</span>
        </div>

        <!-- No Items -->
        <div v-if="availableItems.length === 0" class="text-center py-16">
            <div class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
                <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
            </div>
            <p class="text-sm font-medium text-(--color-text-primary) mt-4">No outstanding bills</p>
            <p class="text-xs text-(--color-text-secondary) mt-1">
                This customer has no unpaid or partially paid {{ isGodown ? 'storage, lot, or delivery' : 'transport' }}
                bills
            </p>
            <button @click="emit('back')"
                class="mt-4 text-sm text-(--color-blue) hover:underline inline-flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                Go back to customer selection
            </button>
        </div>

        <!-- Items List -->
        <div v-else class="space-y-4 max-h-125 overflow-y-auto pr-1">
            <div v-for="(group, groupName) in groupedItems" :key="groupName"
                class="rounded-xl border border-(--color-border) bg-(--color-surface) overflow-hidden">
                <!-- Group Header -->
                <div
                    class="flex items-center justify-between p-4 bg-(--color-muted-bg)/20 border-b border-(--color-border)">
                    <div class="flex items-center gap-3">
                        <!-- Group checkbox - convenience to select/deselect all -->
                        <input type="checkbox" :checked="group.allSelected"
                            @change="toggleGroup(groupName as string, $event)"
                            class="w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 transition-all duration-200" />
                        <div>
                            <h3 class="text-sm font-semibold text-(--color-text-primary)">{{ groupLabel(groupName as
                                string) }}</h3>
                            <p class="text-xs text-(--color-text-secondary)">{{ group.items.length }} items outstanding
                            </p>
                        </div>
                    </div>
                    <span class="text-xs font-medium text-(--color-blue) bg-(--color-blue)/10 px-3 py-1 rounded-lg">
                        {{ group.selectedCount }} selected
                    </span>
                </div>

                <!-- Individual Items -->
                <div class="divide-y divide-(--color-border)">
                    <div v-for="item in group.items" :key="`${item.source_type}_${item.id}`"
                        class="flex items-start gap-3 p-4 hover:bg-(--color-muted-bg)/20 transition-colors cursor-pointer group/item"
                        :class="{ 'opacity-50 cursor-not-allowed': !item.selected && selectedCount >= 12 }"
                        @click="handleItemClick(item, $event)">
                        <!-- Individual checkbox - clicking this ONLY toggles the checkbox -->
                        <div @click.stop>
                            <input type="checkbox" :checked="item.selected" @change="toggleItem(item)"
                                :disabled="!item.selected && selectedCount >= 12"
                                class="mt-1 w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 transition-all duration-200 shrink-0 disabled:opacity-50 disabled:cursor-not-allowed" />
                        </div>

                        <!-- Item Details - clicking this also toggles the item -->
                        <div class="flex-1 min-w-0">
                            <div class="flex items-center gap-3">
                                <span class="text-sm font-medium text-(--color-text-primary)">{{ item.item }}</span>
                                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium"
                                    :class="getStatusClass(item)">
                                    {{ getStatusLabel(item) }}
                                </span>
                            </div>
                            <div
                                class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-(--color-text-secondary) mt-0.5">
                                <span>{{ item.date }}</span>
                                <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)/40"></span>
                                <span class="truncate">{{ item.description }}</span>
                            </div>
                            <div class="flex items-center gap-4 mt-1 text-xs">
                                <span>Total: <span class="font-medium text-(--color-text-primary)">৳ {{
                                    formatAmount(item.total_amount) }}</span></span>
                                <span v-if="item.paid_amount > 0">
                                    Paid: <span class="font-medium text-(--color-green)">৳ {{
                                        formatAmount(item.paid_amount) }}</span>
                                </span>
                                <span>
                                    Due: <span class="font-medium text-(--color-red)">৳ {{ formatAmount(item.amount)
                                    }}</span>
                                </span>
                            </div>
                        </div>

                        <!-- Amount -->
                        <div class="text-right shrink-0">
                            <p class="text-sm font-semibold text-(--color-red)">৳ {{ formatAmount(item.amount) }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Summary -->
        <div v-if="availableItems.length > 0"
            class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-4 rounded-xl border border-(--color-border) bg-(--color-muted-bg)/20">
            <div>
                <p class="text-sm font-medium text-(--color-text-primary)">
                    {{ selectedCount }} of 12 items selected
                </p>
                <div class="flex flex-wrap items-center gap-4 text-sm mt-1">
                    <span>Total Due: <span class="font-semibold text-(--color-red)">৳ {{ formatAmount(selectedTotal)
                    }}</span></span>
                    <span>Already Paid: <span class="font-semibold text-(--color-green)">৳ {{
                        formatAmount(selectedPaid) }}</span></span>
                </div>
                <div class="w-full mt-2 h-1.5 rounded-full bg-(--color-muted-bg) overflow-hidden">
                    <div class="h-full rounded-full transition-all duration-300"
                        :class="selectedCount > 12 ? 'bg-(--color-red)' : 'bg-(--color-blue)'"
                        :style="{ width: `${Math.min((selectedCount / 12) * 100, 100)}%` }"></div>
                </div>
            </div>
            <div class="flex items-center gap-3 w-full sm:w-auto">
                <button @click="emit('back')"
                    class="flex-1 sm:flex-none px-5 py-2.5 text-sm font-medium rounded-xl border border-(--color-border) hover:bg-(--color-muted-bg) transition-all duration-200">
                    Back
                </button>
                <button @click="handleNext" :disabled="selectedCount === 0 || selectedCount > 12"
                    class="flex-1 sm:flex-none px-6 py-2.5 bg-(--color-blue) text-white rounded-xl text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100 flex items-center justify-center gap-2">
                    Review Invoice
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useInvoiceStore } from '@/stores/invoices'
import { useCustomersStore } from '@/stores/customers'
import type { AvailableItem, InvoiceType } from '@/types/invoice'
import { push } from 'notivue'

const props = defineProps<{
    type: InvoiceType
    customerId: number | null
    selectedItems: AvailableItem[]
}>()

const emit = defineEmits<{
    'update:selectedItems': [value: AvailableItem[]]
    'next': []
    'back': []
}>()

const MAX_ITEMS = 12
const invoiceStore = useInvoiceStore()
const customersStore = useCustomersStore()

// Local copy of available items with selection state
const availableItems = ref<AvailableItem[]>([])

const isGodown = computed(() => props.type === 'godown')

const customerName = computed(() => {
    if (!props.customerId) return 'Unknown'
    const customer = customersStore.getCustomerById(props.customerId)
    return customer?.company_name || customer?.contact_person || `Customer #${props.customerId}`
})

const selectedCount = computed(() => {
    return availableItems.value.filter(item => item.selected).length
})

const selectedTotal = computed(() => {
    return availableItems.value
        .filter(item => item.selected)
        .reduce((sum, item) => sum + item.amount, 0)
})

const selectedPaid = computed(() => {
    return availableItems.value
        .filter(item => item.selected)
        .reduce((sum, item) => sum + item.paid_amount, 0)
})

const allSelected = computed(() => {
    return availableItems.value.length > 0 && availableItems.value.every(item => item.selected)
})

const groupedItems = computed(() => {
    const groups: Record<string, { items: AvailableItem[], allSelected: boolean, selectedCount: number }> = {}

    for (const item of availableItems.value) {
        let group = groups[item.source_type]
        if (!group) {
            group = { items: [], allSelected: false, selectedCount: 0 }
            groups[item.source_type] = group
        }
        group.items.push(item)
    }

    for (const key of Object.keys(groups)) {
        const group = groups[key]
        if (group) {
            group.allSelected = group.items.every(item => item.selected)
            group.selectedCount = group.items.filter(item => item.selected).length
        }
    }

    return groups
})

const groupLabel = (key: string): string => {
    const labels: Record<string, string> = {
        storage_bill: 'Storage Bills',
        lot_bill: 'Lot Unload Bills',
        delivery_bill: 'Delivery Bills',
        transport_bill: 'Transport Bills',
    }
    return labels[key] || key
}

const formatAmount = (value: number): string => {
    return value.toFixed(2)
}

const getStatusClass = (item: AvailableItem): string => {
    if (item.status === 'paid') return 'bg-(--color-green)/10 text-(--color-green)'
    if (item.status === 'partial') return 'bg-(--color-yellow)/10 text-(--color-yellow)'
    return 'bg-(--color-red)/10 text-(--color-red)'
}

const getStatusLabel = (item: AvailableItem): string => {
    if (item.status === 'paid') return 'Fully Paid'
    if (item.status === 'partial') return 'Partially Paid'
    return 'Unpaid'
}

// ============================================================
// FIX: Handle both row click AND checkbox click with proper type checking
// ============================================================
const handleItemClick = (item: AvailableItem, event: MouseEvent) => {
    // Get the target element
    const target = event.target as HTMLElement

    // Check if the click was on a checkbox input
    if (target.tagName === 'INPUT' && (target as HTMLInputElement).type === 'checkbox') {
        // Checkbox click - let the @change handler handle it
        return
    }
    // Row click - toggle the item
    toggleItem(item)
}

// Toggle individual item
const toggleItem = (item: AvailableItem) => {
    // If trying to select and already at max, prevent it
    if (!item.selected && selectedCount.value >= MAX_ITEMS) {
        push.warning(`Maximum ${MAX_ITEMS} items allowed per invoice`)
        return
    }

    // Toggle ONLY this item
    item.selected = !item.selected

    // Emit the updated items list to parent
    emit('update:selectedItems', availableItems.value)
}

// Toggle all items in a group
const toggleGroup = (groupName: string, event: Event) => {
    const checked = (event.target as HTMLInputElement).checked
    const group = groupedItems.value[groupName]
    if (!group) return

    const unselectedItems = group.items.filter(item => !item.selected)
    const itemsToSelect = checked ? unselectedItems.length : 0
    const currentSelected = selectedCount.value

    if (checked && currentSelected + itemsToSelect > MAX_ITEMS) {
        const availableSlots = MAX_ITEMS - currentSelected
        let selected = 0
        for (const item of group.items) {
            if (selected < availableSlots && !item.selected) {
                item.selected = true
                selected++
            }
        }
        if (selected < itemsToSelect) {
            push.warning(`Only ${selected} of ${itemsToSelect} items selected (max ${MAX_ITEMS})`)
        }
    } else {
        for (const item of group.items) {
            item.selected = checked
        }
    }

    emit('update:selectedItems', availableItems.value)
}

// Toggle all items
const toggleAll = () => {
    const newState = !allSelected.value

    if (newState && availableItems.value.length > MAX_ITEMS) {
        let count = 0
        for (const item of availableItems.value) {
            if (count < MAX_ITEMS) {
                item.selected = true
                count++
            } else {
                item.selected = false
            }
        }
        push.warning(`Selected ${count} of ${availableItems.value.length} items (max ${MAX_ITEMS})`)
    } else {
        for (const item of availableItems.value) {
            item.selected = newState
        }
    }
    emit('update:selectedItems', availableItems.value)
}

const handleNext = () => {
    if (selectedCount.value === 0) {
        push.error('Please select at least one item')
        return
    }
    if (selectedCount.value > MAX_ITEMS) {
        push.error(`Maximum ${MAX_ITEMS} items allowed`)
        return
    }
    emit('next')
}

// Load available items when customer changes
watch(() => props.customerId, (customerId) => {
    if (customerId) {
        availableItems.value = invoiceStore.getAvailableItemsForCustomer(customerId, props.type)
        for (const item of availableItems.value) {
            item.selected = false
        }
    } else {
        availableItems.value = []
    }
}, { immediate: true })

// Sync with parent selectedItems
watch(() => props.selectedItems, (newItems) => {
    if (newItems.length > 0 && availableItems.value.length > 0) {
        for (const localItem of availableItems.value) {
            const match = newItems.find(item =>
                item.source_type === localItem.source_type &&
                item.id === localItem.id
            )
            if (match) {
                localItem.selected = match.selected
            }
        }
    }
}, { deep: true })
</script>