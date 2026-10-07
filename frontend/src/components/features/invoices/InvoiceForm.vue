<!-- src/components/features/invoices/InvoiceForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- ENTITY -->
        <div v-if="!isEditMode">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Entity
            </h3>

            <EntityPicker v-model:entity-type="form.entity_type" v-model:entity-id="form.entity_id"
                :disabled="submitting" />
        </div>

        <!-- LOCKED ENTITY (edit mode) -->
        <div v-else class="p-3 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
            <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Entity</p>
            <p class="text-sm text-(--color-text-primary) mt-0.5">
                <span class="capitalize">{{ form.entity_type }}</span> — {{ lockedEntityLabel }}
            </p>
        </div>

        <!-- UNBILLED BILLS (create only) -->
        <div v-if="!isEditMode && form.entity_type && form.entity_id" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-3">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Select Bills
                </h3>
                <div class="flex items-center gap-3">
                    <span class="text-xs text-(--color-text-secondary)">
                        {{ selectedCount }} / {{ unbilled.length }} selected
                    </span>
                    <button v-if="unbilled.length > 0" type="button" @click="toggleSelectAll"
                        class="text-xs font-medium text-(--color-blue) hover:opacity-80 transition-opacity">
                        {{ allSelected ? 'Clear all' : 'Select all' }}
                    </button>
                </div>
            </div>

            <div v-if="loadingUnbilled" class="text-sm text-(--color-text-secondary) py-4">
                Loading unpaid bills...
            </div>

            <div v-else-if="unbilled.length === 0"
                class="text-center py-6 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No unpaid bills found for this entity.
            </div>

            <div v-else class="rounded-lg border border-(--color-border) overflow-hidden">
                <div
                    class="grid grid-cols-12 gap-2 px-3 py-2 bg-(--color-muted-bg)/40 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider items-center">
                    <div class="col-span-1 flex items-center justify-center">
                        <input type="checkbox" :checked="allSelected" :indeterminate.prop="someSelected && !allSelected"
                            @change="toggleSelectAll"
                            class="w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 cursor-pointer" />
                    </div>
                    <div class="col-span-5">Bill</div>
                    <div class="col-span-2 text-right">Total</div>
                    <div class="col-span-2 text-right">Paid</div>
                    <div class="col-span-2 text-right">Remaining</div>
                </div>
                <label v-for="b in unbilled" :key="b.bill_type + '-' + b.bill_id"
                    class="grid grid-cols-12 gap-2 px-3 py-2 border-t border-(--color-border) text-sm cursor-pointer transition-colors"
                    :class="isSelected(b) ? 'bg-(--color-blue)/5' : 'hover:bg-(--color-muted-bg)/30'">
                    <div class="col-span-1 flex items-center justify-center">
                        <input type="checkbox" :checked="isSelected(b)" @change="toggleBill(b)"
                            class="w-4 h-4 rounded border-(--color-border) text-(--color-blue) focus:ring-2 focus:ring-(--color-blue)/20 focus:ring-offset-0 cursor-pointer" />
                    </div>
                    <div class="col-span-5 text-(--color-text-primary) truncate">{{ b.label }}</div>
                    <div class="col-span-2 text-right text-(--color-text-secondary)">{{ formatCurrency(b.total) }}</div>
                    <div class="col-span-2 text-right text-(--color-text-secondary)">{{ formatCurrency(b.paid) }}</div>
                    <div class="col-span-2 text-right font-semibold text-(--color-text-primary)">{{
                        formatCurrency(b.remaining) }}</div>
                </label>
            </div>
        </div>

        <!-- DISCOUNTS -->
        <div class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-3">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Discounts
                </h3>
                <button type="button" @click="addDiscount" :disabled="submitting"
                    class="text-xs font-medium text-(--color-blue) hover:opacity-80 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed">
                    + Add discount
                </button>
            </div>

            <div v-if="form.discounts.length === 0"
                class="text-center py-4 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No discounts applied.
            </div>

            <div v-else class="space-y-3">
                <div v-for="(d, idx) in form.discounts" :key="idx"
                    class="grid grid-cols-12 gap-2 items-end p-3 rounded-lg bg-(--color-muted-bg)/10 border border-(--color-border)">
                    <!-- Type -->
                    <div class="col-span-3">
                        <label class="text-xs font-medium text-(--color-text-secondary) block mb-1">Type</label>
                        <select v-model="d.type" :disabled="submitting"
                            class="w-full px-2 py-1.5 rounded-md bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) text-sm focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                            <option value="flat">Flat</option>
                            <option value="percent">Percent</option>
                        </select>
                    </div>

                    <!-- Value -->
                    <div class="col-span-3">
                        <label class="text-xs font-medium text-(--color-text-secondary) block mb-1">Value</label>
                        <input v-model.number="d.value" type="number" step="0.01" min="0" :disabled="submitting"
                            class="w-full px-2 py-1.5 rounded-md bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) text-sm focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>

                    <!-- Reason -->
                    <div class="col-span-5">
                        <label class="text-xs font-medium text-(--color-text-secondary) block mb-1">Reason</label>
                        <input v-model="d.reason" type="text" placeholder="Optional" :disabled="submitting"
                            class="w-full px-2 py-1.5 rounded-md bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) text-sm placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>

                    <!-- Remove -->
                    <div class="col-span-1 flex justify-end">
                        <button type="button" @click="removeDiscount(idx)" :disabled="submitting"
                            class="w-8 h-8 flex items-center justify-center rounded-md text-(--color-red) hover:bg-(--color-red)/10 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- TOTALS PREVIEW -->
        <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border) space-y-2">
            <div class="flex items-center justify-between text-sm">
                <span class="text-(--color-text-secondary)">Subtotal</span>
                <span class="font-semibold text-(--color-text-primary)">{{ formatCurrency(previewSubtotal) }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
                <span class="text-(--color-text-secondary)">Discount</span>
                <span class="font-semibold text-(--color-red)">− {{ formatCurrency(previewDiscount) }}</span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-(--color-border)">
                <span class="text-sm font-semibold text-(--color-text-primary)">Total</span>
                <span class="text-lg font-bold text-(--color-blue)">{{ formatCurrency(previewTotal) }}</span>
            </div>
        </div>

        <!-- NOTES -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Notes
            </label>
            <textarea v-model="form.notes" rows="3" placeholder="Optional notes" :disabled="submitting"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>

        <!-- ACTIONS -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitDisabled"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ isEditMode ? 'Saving...' : 'Creating...' }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Invoice' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type {
    Invoice,
    InvoiceDetail,
    UnbilledBill,
    EntityType,
    DiscountType,
    BillRef,
} from '@/types/invoice'
import { useInvoicesStore } from '@/stores/invoices'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useBrokersStore } from '@/stores/brokers'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'
import EntityPicker from './EntityPicker.vue'

interface DiscountForm {
    type: DiscountType
    value: number
    reason: string
}

const props = defineProps<{
    invoice?: InvoiceDetail | Invoice | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'invoice-created': [invoice: InvoiceDetail]
    'invoice-updated': [invoice: InvoiceDetail]
    'cancel': []
}>()

const invoicesStore = useInvoicesStore()
const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const brokersStore = useBrokersStore()

const submitting = ref(false)
const loadingUnbilled = ref(false)
const unbilled = ref<UnbilledBill[]>([])
const selectedBills = ref<Set<string>>(new Set())

const isEditMode = computed(() => props.mode === 'edit' || !!props.invoice)

const form = ref({
    entity_type: null as EntityType | null,
    entity_id: null as number | null,
    discounts: [] as DiscountForm[],
    notes: '',
})

const lockedEntityLabel = computed(() => {
    if (!props.invoice) return ''
    switch (props.invoice.entity_type) {
        case 'customer': {
            const c = customersStore.getCustomerById(props.invoice.entity_id)
            return c ? (c.company_name || c.contact_person || `Customer #${c.id}`) : `Customer #${props.invoice.entity_id}`
        }
        case 'majhi':
            return majhisStore.getMajhiName(props.invoice.entity_id)
        case 'godown':
            return godownsStore.getGodownName(props.invoice.entity_id)
        case 'broker':
            return brokersStore.getBrokerName(props.invoice.entity_id)
        default:
            return `#${props.invoice.entity_id}`
    }
})

const billKey = (b: UnbilledBill): string => `${b.bill_type}-${b.bill_id}`

const isSelected = (b: UnbilledBill): boolean => selectedBills.value.has(billKey(b))

const selectedCount = computed(() => selectedBills.value.size)

const allSelected = computed(() =>
    unbilled.value.length > 0 && selectedBills.value.size === unbilled.value.length
)

const someSelected = computed(() =>
    selectedBills.value.size > 0 && !allSelected.value
)

const toggleBill = (b: UnbilledBill) => {
    const key = billKey(b)
    const next = new Set(selectedBills.value)
    if (next.has(key)) {
        next.delete(key)
    } else {
        next.add(key)
    }
    selectedBills.value = next
}

const toggleSelectAll = () => {
    if (allSelected.value) {
        selectedBills.value = new Set()
    } else {
        selectedBills.value = new Set(unbilled.value.map(billKey))
    }
}

const previewSubtotal = computed(() => {
    if (isEditMode.value) {
        return props.invoice?.subtotal ?? 0
    }
    let total = 0
    for (const b of unbilled.value) {
        if (selectedBills.value.has(billKey(b))) {
            total += b.remaining
        }
    }
    return total
})

const previewDiscount = computed(() => {
    const sub = previewSubtotal.value
    let total = 0
    for (const d of form.value.discounts) {
        if (d.type === 'flat') {
            total += d.value || 0
        } else {
            total += (sub * (d.value || 0)) / 100
        }
    }
    if (total > sub) total = sub
    return total
})

const previewTotal = computed(() => {
    const t = previewSubtotal.value - previewDiscount.value
    return t < 0 ? 0 : t
})

const submitDisabled = computed(() => {
    if (submitting.value) return true
    if (isEditMode.value) return false
    if (!form.value.entity_type || !form.value.entity_id) return true
    if (selectedBills.value.size === 0) return true
    return false
})

// =============================================================================
// UNBILLED PREVIEW
// =============================================================================

const loadUnbilled = async () => {
    if (isEditMode.value) return
    if (!form.value.entity_type || !form.value.entity_id) {
        unbilled.value = []
        selectedBills.value = new Set()
        return
    }
    loadingUnbilled.value = true
    try {
        unbilled.value = await invoicesStore.fetchUnbilled(form.value.entity_type, form.value.entity_id)
        // Auto-select all bills by default — user can uncheck.
        selectedBills.value = new Set(unbilled.value.map(billKey))
    } finally {
        loadingUnbilled.value = false
    }
}

watch(() => [form.value.entity_type, form.value.entity_id], loadUnbilled)

// =============================================================================
// DISCOUNTS
// =============================================================================

const addDiscount = () => {
    form.value.discounts.push({ type: 'flat', value: 0, reason: '' })
}

const removeDiscount = (idx: number) => {
    form.value.discounts.splice(idx, 1)
}

// =============================================================================
// INIT
// =============================================================================

const initialize = () => {
    if (props.invoice) {
        form.value = {
            entity_type: props.invoice.entity_type,
            entity_id: props.invoice.entity_id,
            discounts: ((props.invoice as InvoiceDetail).discounts ?? []).map((d) => ({
                type: d.type,
                value: d.value,
                reason: d.reason ?? '',
            })),
            notes: props.invoice.notes ?? '',
        }
    } else {
        form.value = {
            entity_type: null,
            entity_id: null,
            discounts: [],
            notes: '',
        }
        unbilled.value = []
        selectedBills.value = new Set()
    }
}

watch(() => props.invoice, initialize, { immediate: true })

// =============================================================================
// SUBMIT
// =============================================================================

const submit = async () => {
    if (!isEditMode.value) {
        if (!form.value.entity_type) {
            push.error('Please select an entity type')
            return
        }
        if (!form.value.entity_id) {
            push.error('Please select an entity')
            return
        }
        if (selectedBills.value.size === 0) {
            push.error('Please select at least one bill')
            return
        }
    }

    // Validate discounts
    for (const d of form.value.discounts) {
        if (d.value < 0) {
            push.error('Discount value cannot be negative')
            return
        }
        if (d.type === 'percent' && d.value > 100) {
            push.error('Percent discount cannot exceed 100')
            return
        }
    }

    submitting.value = true

    try {
        const discounts = form.value.discounts.map((d) => ({
            type: d.type,
            value: d.value,
            reason: d.reason.trim() || null,
        }))

        if (isEditMode.value && props.invoice) {
            const result = await invoicesStore.updateInvoice(props.invoice.id, {
                discounts,
                notes: form.value.notes.trim() || null,
            })
            if (result) {
                emit('invoice-updated', result)
            }
        } else {
            const bills: BillRef[] = unbilled.value
                .filter((b) => selectedBills.value.has(billKey(b)))
                .map((b) => ({ bill_type: b.bill_type, bill_id: b.bill_id }))

            const result = await invoicesStore.createInvoice({
                entity_type: form.value.entity_type!,
                entity_id: form.value.entity_id!,
                bills,
                discounts,
                notes: form.value.notes.trim() || null,
            })
            if (result) {
                initialize()
                emit('invoice-created', result)
            }
        }
    } finally {
        submitting.value = false
    }
}

// =============================================================================
// PRELOAD
// =============================================================================

onMounted(async () => {
    if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
    if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
    if (godownsStore.godowns.length === 0) await godownsStore.fetchGodowns()
    if (brokersStore.brokers.length === 0) await brokersStore.fetchBrokers()
})
</script>