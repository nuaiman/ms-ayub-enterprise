<!-- src/components/features/invoices/EntityPicker.vue -->
<template>
    <div class="space-y-4">
        <!-- Entity Type -->
        <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                Entity Type <span class="text-(--color-red)">*</span>
            </label>
            <div
                class="grid grid-cols-2 sm:grid-cols-4 gap-2 p-1 rounded-lg bg-(--color-muted-bg) border border-(--color-border)">
                <button v-for="opt in entityTypeOptions" :key="opt.value" type="button" :disabled="disabled"
                    @click="selectType(opt.value)"
                    class="px-3 py-2 rounded-md text-sm font-medium capitalize transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="entityType === opt.value
                        ? 'bg-(--color-surface) text-(--color-text-primary) shadow-sm'
                        : 'text-(--color-text-secondary) hover:text-(--color-text-primary)'">
                    {{ opt.label }}
                </button>
            </div>
        </div>

        <!-- Entity Selection -->
        <div v-if="entityType">
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                {{ entityTypeLabel }} <span class="text-(--color-red)">*</span>
            </label>

            <div v-if="loadingEntities"
                class="px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-sm text-(--color-text-secondary)">
                Loading {{ entityTypeLabel.toLowerCase() }}s...
            </div>

            <select v-else :value="entityId ?? ''" @change="onSelect" :disabled="disabled || entityOptions.length === 0"
                class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                <option value="">Select {{ entityTypeLabel.toLowerCase() }}</option>
                <option v-for="opt in entityOptions" :key="opt.id" :value="opt.id">
                    {{ opt.label }}
                </option>
            </select>

            <p v-if="!loadingEntities && entityOptions.length === 0" class="text-xs text-(--color-text-secondary) mt-1">
                No {{ entityTypeLabel.toLowerCase() }}s available.
            </p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { EntityType } from '@/types/invoice'
import { useCustomersStore } from '@/stores/customers'
import { useMajhisStore } from '@/stores/majhis'
import { useGodownsStore } from '@/stores/godowns'
import { useBrokersStore } from '@/stores/brokers'

interface EntityOption {
    id: number
    label: string
}

const props = withDefaults(defineProps<{
    entityType: EntityType | null
    entityId: number | null
    disabled?: boolean
    lockEntityType?: boolean
}>(), {
    disabled: false,
    lockEntityType: false,
})

const emit = defineEmits<{
    (e: 'update:entityType', value: EntityType | null): void
    (e: 'update:entityId', value: number | null): void
}>()

const customersStore = useCustomersStore()
const majhisStore = useMajhisStore()
const godownsStore = useGodownsStore()
const brokersStore = useBrokersStore()

const loadingEntities = ref(false)

const entityTypeOptions: { value: EntityType; label: string }[] = [
    { value: 'customer', label: 'Customer' },
    { value: 'majhi', label: 'Majhi' },
    { value: 'godown', label: 'Godown' },
    { value: 'broker', label: 'Broker' },
]

const entityTypeLabel = computed(() => {
    switch (props.entityType) {
        case 'customer': return 'Customer'
        case 'majhi': return 'Majhi'
        case 'godown': return 'Godown'
        case 'broker': return 'Broker'
        default: return ''
    }
})

const entityOptions = computed<EntityOption[]>(() => {
    switch (props.entityType) {
        case 'customer':
            return customersStore.customers.map((c) => ({
                id: c.id,
                label: c.company_name || c.contact_person || `Customer #${c.id}`,
            }))
        case 'majhi':
            return majhisStore.majhis.map((m) => ({ id: m.id, label: m.name }))
        case 'godown':
            return godownsStore.godowns.map((g) => ({ id: g.id, label: g.name }))
        case 'broker':
            return brokersStore.brokers.map((b) => ({ id: b.id, label: b.name }))
        default:
            return []
    }
})

const selectType = (t: EntityType) => {
    if (props.disabled || props.lockEntityType) return
    if (props.entityType === t) return
    emit('update:entityType', t)
    emit('update:entityId', null)
}

const onSelect = (e: Event) => {
    const raw = (e.target as HTMLSelectElement).value
    emit('update:entityId', raw ? parseInt(raw) : null)
}

const ensureLoaded = async (t: EntityType) => {
    loadingEntities.value = true
    try {
        switch (t) {
            case 'customer':
                if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
                break
            case 'majhi':
                if (majhisStore.majhis.length === 0) await majhisStore.fetchMajhis()
                break
            case 'godown':
                if (godownsStore.godowns.length === 0) await godownsStore.fetchGodowns()
                break
            case 'broker':
                if (brokersStore.brokers.length === 0) await brokersStore.fetchBrokers()
                break
        }
    } finally {
        loadingEntities.value = false
    }
}

watch(() => props.entityType, (t) => {
    if (t) ensureLoaded(t)
})

onMounted(() => {
    if (props.entityType) ensureLoaded(props.entityType)
})
</script>