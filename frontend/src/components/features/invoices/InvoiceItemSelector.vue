<!-- src/components/features/invoices/InvoiceItemSelector.vue -->
<template>
    <div class="space-y-6">
        <!-- Header -->
        <div
            class="flex flex-col sm:flex-row items-start sm:items-center justify-between p-4 rounded-xl bg-(--color-muted-bg)/30 border border-(--color-border) gap-3">
            <div>
                <p class="text-sm font-medium text-(--color-text-primary)">
                    {{ partyTypeLabel(party) }}: <span class="font-semibold">{{ partyName }}</span>
                </p>
                <p class="text-xs text-(--color-text-secondary) mt-0.5">
                    Select unpaid or partially paid bills
                </p>
            </div>
            <div class="flex items-center gap-4">
                <span class="text-sm font-medium text-(--color-blue)">
                    {{ selectedCount }} selected
                </span>
                <button @click="toggleAll"
                    class="text-xs font-medium text-(--color-blue) hover:text-(--color-blue)/70 transition-colors"
                    :disabled="availableItems.length === 0">
                    {{ allSelected ? 'Deselect All' : 'Select All' }}
                </button>
            </div>
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
                This {{ partyTypeLabel(party).toLowerCase() }} has no unpaid or partially paid bills
            </p>
            <button @click="emit('back')"
                class="mt-4 text-sm text-(--color-blue) hover:underline inline-flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
                Go back to party selection
            </button>
        </div>

        <!-- Items List -->
        <div v-else class="space-y-4 max-h-125 overflow-y-auto pr-1">
            <div v-for="(group, groupName) in groupedItems" :key="groupName"
                class="rounded-xl border border-(--color-border) bg-(--color-surface) overflow-hidden">
                <div
                    class="flex items-center justify-between p-4 bg-(--color-muted-bg)/20 border-b border-(--color-border)">
                    <div class="flex items-center gap-3">
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

                <div class="divide-y divide-(--color-border)">
                    <div v-for="item in group.items" :key="`${item.source_type}_${item.id}`"
                        class="flex items-start gap-3 p-4 hover:bg-(--color-muted-bg)/20 transition-colors cursor-pointer group/item"
                        @click="handleItemClick(item, $event)">
                        <div @click.stop>
                            <input type="checkbox" :checked="item.selected" @change="toggleItem(item)"
                                class="mt-1 w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 transition-all duration-200 shrink-0" />
                        </div>

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
                    {{ selectedCount }} items selected
                </p>
                <div class="flex flex-wrap items-center gap-4 text-sm mt-1">
                    <span>Total Due: <span class="font-semibold text-(--color-red)">৳ {{ formatAmount(selectedTotal)
                            }}</span></span>
                    <span>Already Paid: <span class="font-semibold text-(--color-green)">৳ {{
                        formatAmount(selectedPaid) }}</span></span>
                </div>
            </div>
            <div class="flex items-center gap-3 w-full sm:w-auto">
                <button @click="emit('back')"
                    class="flex-1 sm:flex-none px-5 py-2.5 text-sm font-medium rounded-xl border border-(--color-border) hover:bg-(--color-muted-bg) transition-all duration-200">
                    Back
                </button>
                <button @click="handleNext" :disabled="selectedCount === 0"
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
import type { AvailableItem, InvoiceParty } from '@/types/invoice'
import { partyTypeLabel } from '@/types/invoice'
import { push } from 'notivue'

const props = defineProps<{
    party: InvoiceParty
    partyId: number | null
    partyName: string
    selectedItems: AvailableItem[]
}>()

const emit = defineEmits<{
    'update:selectedItems': [value: AvailableItem[]]
    'next': []
    'back': []
}>()

const invoiceStore = useInvoiceStore()
const availableItems = ref<AvailableItem[]>([])

const selectedCount = computed(() => availableItems.value.filter(i => i.selected).length)
const selectedTotal = computed(() => availableItems.value.filter(i => i.selected).reduce((s, i) => s + i.amount, 0))
const selectedPaid = computed(() => availableItems.value.filter(i => i.selected).reduce((s, i) => s + i.paid_amount, 0))
const allSelected = computed(() => availableItems.value.length > 0 && availableItems.value.every(i => i.selected))

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
        const g = groups[key]
        if (g) {
            g.allSelected = g.items.every(i => i.selected)
            g.selectedCount = g.items.filter(i => i.selected).length
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
        additional_charge: 'Additional Charges',
        broker_vehicle_bill: 'Broker Vehicle Bills',
        godown_store_bill: 'Godown Store Bills',
        majhi_lot_bill: 'Majhi Lot Bills',
        majhi_loading_bill: 'Majhi Loading Bills',
    }
    return labels[key] || key
}

const formatAmount = (v: number): string => v.toFixed(2)

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

const handleItemClick = (item: AvailableItem, event: MouseEvent) => {
    const target = event.target as HTMLElement
    if (target.tagName === 'INPUT' && (target as HTMLInputElement).type === 'checkbox') return
    toggleItem(item)
}

const toggleItem = (item: AvailableItem) => {
    item.selected = !item.selected
    emit('update:selectedItems', availableItems.value)
}

const toggleGroup = (groupName: string, event: Event) => {
    const checked = (event.target as HTMLInputElement).checked
    const group = groupedItems.value[groupName]
    if (!group) return
    for (const item of group.items) item.selected = checked
    emit('update:selectedItems', availableItems.value)
}

const toggleAll = () => {
    const newState = !allSelected.value
    for (const item of availableItems.value) item.selected = newState
    emit('update:selectedItems', availableItems.value)
}

const handleNext = () => {
    if (selectedCount.value === 0) {
        push.error('Please select at least one item')
        return
    }
    emit('next')
}

watch(() => [props.party, props.partyId], () => {
    if (props.partyId) {
        availableItems.value = invoiceStore.getAvailableItemsForParty(props.party, props.partyId)
        for (const item of availableItems.value) item.selected = false
    } else {
        availableItems.value = []
    }
}, { immediate: true })

watch(() => props.selectedItems, (newItems) => {
    if (newItems.length > 0 && availableItems.value.length > 0) {
        for (const local of availableItems.value) {
            const match = newItems.find(i => i.source_type === local.source_type && i.id === local.id)
            if (match) local.selected = match.selected
        }
    }
}, { deep: true })
</script>