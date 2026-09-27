<!-- src/components/features/invoices/InvoicePartySelect.vue -->
<template>
    <div class="space-y-6">
        <div class="p-6 rounded-xl bg-(--color-muted-bg)/30 border border-(--color-border)">
            <div class="flex items-start gap-3 mb-6">
                <div
                    class="w-10 h-10 rounded-xl bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center shrink-0">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                </div>
                <div>
                    <h2 class="text-base font-semibold text-(--color-text-primary)">Select Party</h2>
                    <p class="text-sm text-(--color-text-secondary)">Choose who this invoice is for</p>
                </div>
            </div>

            <!-- Party Type -->
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-2">
                    Party Type <span class="text-(--color-red)">*</span>
                </label>
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
                    <button v-for="p in partyTypes" :key="p.value" type="button" @click="handlePartyTypeChange(p.value)"
                        class="px-3 py-3 rounded-xl border text-sm font-medium transition-all duration-200 flex flex-col items-center gap-2"
                        :class="localParty === p.value
                            ? 'border-(--color-blue) bg-(--color-blue)/10 text-(--color-blue)'
                            : 'border-(--color-border) text-(--color-text-secondary) hover:bg-(--color-muted-bg)'">
                        <span v-html="p.icon" class="w-5 h-5 flex items-center justify-center"></span>
                        {{ p.label }}
                    </button>
                </div>
            </div>

            <!-- Party Picker -->
            <div class="mt-6">
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    {{ partyTypeLabel(localParty) }} <span class="text-(--color-red)">*</span>
                </label>
                <select v-model="localPartyId"
                    class="w-full px-4 py-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-2 focus:ring-(--color-blue)/20 focus:border-(--color-blue) transition-all duration-200 appearance-none">
                    <option :value="null">Select a {{ partyTypeLabel(localParty).toLowerCase() }}...</option>
                    <option v-for="option in partyOptions" :key="option.id" :value="option.id">
                        {{ option.label }}
                    </option>
                </select>
            </div>

            <!-- Summary -->
            <div v-if="localPartyId && selectedPartyOption" class="mt-4 pt-4 border-t border-(--color-border)">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                    <div>
                        <p class="text-xs text-(--color-text-secondary)">{{ partyTypeLabel(localParty) }}</p>
                        <p class="font-medium text-(--color-text-primary)">{{ selectedPartyOption.label }}</p>
                    </div>
                    <div v-if="selectedPartyOption.sub">
                        <p class="text-xs text-(--color-text-secondary)">Details</p>
                        <p class="font-medium text-(--color-text-primary)">{{ selectedPartyOption.sub }}</p>
                    </div>
                </div>
                <p class="text-xs text-(--color-text-secondary) mt-3">
                    All unpaid and partially paid bills for this {{ partyTypeLabel(localParty).toLowerCase() }} will be
                    available for selection
                </p>
            </div>
        </div>

        <div class="flex justify-end">
            <button @click="handleNext" :disabled="!localPartyId"
                class="px-6 py-2.5 bg-(--color-blue) text-white rounded-xl text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100 flex items-center gap-2">
                Next Step
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useCustomersStore } from '@/stores/customers'
import { useBrokersStore } from '@/stores/brokers'
import { useGodownsStore } from '@/stores/godowns'
import { useMajhisStore } from '@/stores/majhis'
import type { InvoiceParty } from '@/types/invoice'
import { partyTypeLabel } from '@/types/invoice'

const props = defineProps<{
    party: InvoiceParty
    partyId: number | null
}>()

const emit = defineEmits<{
    'update:party': [value: InvoiceParty]
    'update:partyId': [value: number | null]
    'next': []
}>()

const customersStore = useCustomersStore()
const brokersStore = useBrokersStore()
const godownsStore = useGodownsStore()
const majhisStore = useMajhisStore()

const localParty = ref<InvoiceParty>(props.party)
const localPartyId = ref<number | null>(props.partyId)

const partyTypes: { value: InvoiceParty; label: string; icon: string }[] = [
    {
        value: 'customer',
        label: 'Customer',
        icon: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
    },
    {
        value: 'broker',
        label: 'Broker',
        icon: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>`,
    },
    {
        value: 'godown',
        label: 'Godown',
        icon: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>`,
    },
    {
        value: 'majhi',
        label: 'Majhi',
        icon: `<svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>`,
    },
]

interface PartyOption {
    id: number
    label: string
    sub?: string
}

const partyOptions = computed<PartyOption[]>(() => {
    switch (localParty.value) {
        case 'customer':
            return customersStore.customers.map(c => ({
                id: c.id,
                label: c.company_name || c.contact_person || `Customer #${c.id}`,
                sub: c.phone,
            }))
        case 'broker':
            return brokersStore.brokers.map(b => ({
                id: b.id,
                label: b.name,
                sub: b.phone || undefined,
            }))
        case 'godown':
            return godownsStore.godowns.map(g => ({
                id: g.id,
                label: g.name,
                sub: g.phone || undefined,
            }))
        case 'majhi':
            return majhisStore.majhis.map(m => ({
                id: m.id,
                label: m.name,
                sub: m.phone || undefined,
            }))
        default:
            return []
    }
})

const selectedPartyOption = computed(() => {
    if (!localPartyId.value) return null
    return partyOptions.value.find(o => o.id === localPartyId.value) || null
})

watch(localParty, () => {
    localPartyId.value = null
    emit('update:party', localParty.value)
    emit('update:partyId', null)
})

watch(localPartyId, (val) => {
    emit('update:partyId', val)
})

const handlePartyTypeChange = (value: InvoiceParty) => {
    localParty.value = value
}

const handleNext = () => {
    if (localPartyId.value) {
        emit('next')
    }
}
</script>