<!-- src/components/features/transports/TransportRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Customer - 3 columns -->
        <div class="col-span-3 min-w-0">
            <div class="flex items-center gap-2.5">
                <div class="shrink-0">
                    <div v-if="transport.image_url"
                        class="w-9 h-9 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="getImageUrl(transport.image_url)" :alt="customerLabel"
                            class="w-full h-full object-cover" />
                    </div>
                    <div v-else
                        class="w-9 h-9 rounded-lg bg-(--color-blue)/10 border border-(--color-border) flex items-center justify-center">
                        <svg class="w-4 h-4 text-(--color-blue)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M8 18L12 22M12 22L16 18M12 22V10M21 14L12 10L3 14M21 14L12 18M21 14V18M3 14V18M3 14L12 18M3 14L12 10M3 14V10M21 10L12 6M3 10L12 6M21 10L12 14M3 10L12 14" />
                        </svg>
                    </div>
                </div>
                <div class="min-w-0">
                    <div class="font-medium text-(--color-text-primary) truncate text-sm">
                        {{ customerLabel }}
                    </div>
                </div>
            </div>
        </div>

        <!-- Transport Type - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block capitalize">
                {{ transport.transport_type || '৳' }}
            </span>
        </div>

        <!-- From - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ transport.from_location }}
            </span>
        </div>

        <!-- To - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ transport.to_location || '৳' }}
            </span>
            <span class="text-xs text-(--color-text-secondary)/70 truncate block mt-0.5">
                {{ formatDate(transport.transport_date) }}
            </span>
        </div>

        <!-- Vehicle Count - 1 column -->
        <div class="col-span-1">
            <span class="text-sm text-(--color-text-secondary)">
                {{ transport.vehicle_quantity }}
            </span>
        </div>

        <!-- Commission - 1 column -->
        <div class="col-span-1">
            <span class="text-sm font-semibold text-(--color-text-primary)">
                {{ formatCurrency(transport.office_commission_amount) }}
            </span>
        </div>

        <!-- Actions - 1 column -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-52 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
                    <button @click="handleView"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                        </svg>
                        View Details
                    </button>

                    <button @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>

                    <button @click="handleOpenDemarage" :disabled="transportVehicles.length === 0"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-yellow) hover:bg-(--color-muted-bg) transition-colors disabled:opacity-40 disabled:cursor-not-allowed">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        Add Demarage
                    </button>

                    <button @click="handleDelete"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-red) hover:bg-(--color-muted-bg) transition-colors border-t border-(--color-border) mt-1 pt-1">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        Delete
                    </button>
                </div>
            </Transition>

            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>

        <!-- Demarage Dialog -->
        <BaseDialog v-model="demarageDialogOpen" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-yellow)/10 text-(--color-yellow) flex items-center justify-center shrink-0">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Add Demarage</h2>
                        <p class="text-xs text-(--color-text-secondary)">
                            Transport #{{ transport.id }}
                        </p>
                    </div>
                </div>

                <div class="space-y-3">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Vehicle <span class="text-(--color-red)">*</span>
                        </label>
                        <select v-model="demarageForm.vehicleId"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
                            <option :value="null">Select a vehicle</option>
                            <option v-for="v in transportVehicles" :key="v.id" :value="v.id">
                                {{ v.vehicle_number }}
                            </option>
                        </select>
                    </div>

                    <!-- Existing value notice -->
                    <div v-if="selectedDemarageVehicle && selectedDemarageVehicle.demarage_cost > 0"
                        class="p-3 rounded-lg bg-(--color-yellow)/10 border border-(--color-yellow)/20 text-xs text-(--color-yellow)">
                        Existing demarage on {{ selectedDemarageVehicle.vehicle_number }}:
                        {{ formatCurrency(selectedDemarageVehicle.demarage_cost) }}.
                        Saving will overwrite it.
                    </div>

                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Demarage Amount <span class="text-(--color-red)">*</span>
                        </label>
                        <div class="relative">
                            <span
                                class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                            <input v-model.number="demarageForm.amount" type="number" step="0.01" min="0"
                                placeholder="0.00" :disabled="!demarageForm.vehicleId"
                                class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                        </div>
                    </div>
                </div>
            </div>

            <template #actions>
                <button @click="demarageDialogOpen = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>
                <button @click="confirmDemarage" :disabled="!demarageForm.vehicleId || savingDemarage"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-yellow) text-white rounded-lg hover:opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
                    {{ savingDemarage ? 'Saving...' : 'Save Demarage' }}
                </button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Transport } from '@/types/transport'
import { useCustomersStore } from '@/stores/customers'
import { useVehiclesStore } from '@/stores/vehicles'
import { formatCurrency } from '@/utils/currency'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'
import BaseDialog from '@/components/ui/BaseDialog.vue'

const props = defineProps<{
    transport: Transport
}>()

const emit = defineEmits<{
    'view': [transport: Transport]
    'edit': [transport: Transport]
    'delete': [transport: Transport]
    'updated': []
}>()

const customersStore = useCustomersStore()
const vehiclesStore = useVehiclesStore()
const isOpen = ref(false)

const customerLabel = computed(() => {
    if (!props.transport.customer_id) return '৳'
    return customersStore.getCustomerName(props.transport.customer_id)
})

const transportVehicles = computed(() => vehiclesStore.getVehiclesByTransportId(props.transport.id))

const formatDate = (dateStr: string): string => {
    return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
    })
}

const toggleMenu = () => { isOpen.value = !isOpen.value }
const closeMenu = () => { isOpen.value = false }
const handleView = () => { closeMenu(); emit('view', props.transport) }
const handleEdit = () => { closeMenu(); emit('edit', props.transport) }
const handleDelete = () => { closeMenu(); emit('delete', props.transport) }

// ============= Demarage =============

const demarageDialogOpen = ref(false)
const savingDemarage = ref(false)
const demarageForm = ref({
    vehicleId: null as number | null,
    amount: 0,
})

const selectedDemarageVehicle = computed(() => {
    if (!demarageForm.value.vehicleId) return null
    return transportVehicles.value.find(v => v.id === demarageForm.value.vehicleId) || null
})

const handleOpenDemarage = () => {
    closeMenu()
    demarageForm.value = {
        vehicleId: null,
        amount: 0,
    }
    demarageDialogOpen.value = true
}

watch(() => demarageForm.value.vehicleId, (newId) => {
    if (!newId) {
        demarageForm.value.amount = 0
        return
    }
    const v = transportVehicles.value.find(x => x.id === newId)
    if (v) {
        demarageForm.value.amount = v.demarage_cost || 0
    }
})

const confirmDemarage = async () => {
    if (!demarageForm.value.vehicleId) {
        push.error('Please select a vehicle')
        return
    }
    if (demarageForm.value.amount < 0) {
        push.error('Demarage amount cannot be negative')
        return
    }

    savingDemarage.value = true
    try {
        const result = await vehiclesStore.updateVehicle(demarageForm.value.vehicleId, {
            demarage_cost: demarageForm.value.amount,
        })

        if (result) {
            demarageDialogOpen.value = false
            emit('updated')
        }
    } finally {
        savingDemarage.value = false
    }
}
</script>