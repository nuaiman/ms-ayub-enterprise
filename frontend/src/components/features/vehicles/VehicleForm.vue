<!-- src/components/features/vehicles/VehicleForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Transport Selection -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Select Transport
            </h3>

            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                    Transport <span class="text-(--color-red)">*</span>
                </label>
                <select v-model="form.transport_id" @change="onTransportChange" required :disabled="submitting"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                    <option :value="null">Select a transport</option>
                    <option v-for="transport in transportOptions" :key="transport.id" :value="transport.id">
                        #{{ transport.id }} - {{ getTransportDisplayName(transport) }}
                    </option>
                </select>
            </div>

            <!-- Vehicle Count Info -->
            <div v-if="selectedTransport"
                class="mt-3 p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <div class="flex items-center justify-between text-sm">
                    <span class="text-(--color-text-secondary)">Vehicles for this transport:</span>
                    <span class="font-medium text-(--color-text-primary)">
                        {{ currentVehicleCount }} / {{ selectedTransport.vehicle_quantity }}
                    </span>
                </div>
                <div class="w-full h-2 mt-2 rounded-full bg-(--color-muted-bg) overflow-hidden">
                    <div class="h-full rounded-full transition-all duration-500" :class="[
                        currentVehicleCount >= selectedTransport.vehicle_quantity ? 'bg-(--color-red)' : 'bg-(--color-blue)',
                        currentVehicleCount === 0 ? 'w-0' : ''
                    ]" :style="{ width: progressPercentage + '%' }">
                    </div>
                </div>
                <p v-if="currentVehicleCount >= selectedTransport.vehicle_quantity"
                    class="text-xs text-(--color-red) mt-1.5">
                    All vehicles already added for this transport
                </p>
                <p v-else-if="currentVehicleCount > 0" class="text-xs text-(--color-text-secondary) mt-1.5">
                    {{ selectedTransport.vehicle_quantity - currentVehicleCount }} vehicle(s) remaining
                </p>
            </div>
        </div>

        <!-- Vehicle Fields -->
        <div v-if="selectedTransport && remainingVehicles > 0" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Vehicle Details
                </h3>
                <span
                    class="text-xs font-medium text-(--color-text-secondary) bg-(--color-muted-bg) px-2.5 py-1 rounded-full">
                    {{ totalVehiclesToAdd }} vehicle(s) to add
                </span>
            </div>

            <!-- Main Vehicle -->
            <div class="p-4 rounded-lg border border-(--color-blue)/30 bg-(--color-blue)/5">
                <h4 class="text-sm font-medium text-(--color-text-primary) flex items-center gap-2 mb-3">
                    <span
                        class="w-6 h-6 rounded-full bg-(--color-blue) text-white flex items-center justify-center text-xs font-bold">1</span>
                    Main Vehicle
                </h4>

                <VehicleFields v-model:vehicle-number="form.main_vehicle.vehicle_number"
                    v-model:broker-id="form.main_vehicle.broker_id" v-model:joma-cost="form.main_vehicle.joma_cost"
                    v-model:vehicle-cost="form.main_vehicle.vehicle_cost"
                    v-model:other-cost="form.main_vehicle.other_cost"
                    v-model:labour-cost="form.main_vehicle.labour_cost"
                    v-model:demarage-cost="form.main_vehicle.demarage_cost" v-model:notes="form.main_vehicle.notes"
                    :broker-options="brokerOptions" :disabled="submitting" :required="true" />
            </div>

            <!-- Additional Vehicles -->
            <div v-if="remainingVehicles > 1" class="mt-4 space-y-4">
                <div v-for="(vehicle, index) in form.additional_vehicles" :key="index"
                    class="p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10">

                    <div class="flex items-center justify-between mb-3">
                        <h4 class="text-sm font-medium text-(--color-text-primary) flex items-center gap-2">
                            <span
                                class="w-6 h-6 rounded-full bg-(--color-muted-bg) text-(--color-text-secondary) flex items-center justify-center text-xs font-bold">
                                {{ index + 2 }}
                            </span>
                            Additional Vehicle #{{ index + 2 }}
                        </h4>
                        <button type="button" @click="removeAdditionalVehicle(index)"
                            :disabled="form.additional_vehicles.length <= 1"
                            class="text-xs text-(--color-red) hover:bg-(--color-red)/10 px-3 py-1 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed">
                            Remove
                        </button>
                    </div>

                    <VehicleFields v-model:vehicle-number="vehicle.vehicle_number" v-model:broker-id="vehicle.broker_id"
                        v-model:joma-cost="vehicle.joma_cost" v-model:vehicle-cost="vehicle.vehicle_cost"
                        v-model:other-cost="vehicle.other_cost" v-model:labour-cost="vehicle.labour_cost"
                        v-model:demarage-cost="vehicle.demarage_cost" v-model:notes="vehicle.notes"
                        :broker-options="brokerOptions" :disabled="submitting" :required="false" />
                </div>

                <div v-if="remainingVehicles - 1 > form.additional_vehicles.length" class="mt-4">
                    <button type="button" @click="addAdditionalVehicle"
                        class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 5v14M5 12h14" />
                        </svg>
                        Add Additional Vehicle
                    </button>
                    <p class="text-xs text-(--color-text-secondary) mt-2">
                        {{ form.additional_vehicles.length }} of {{ remainingVehicles - 1 }} additional vehicles added
                    </p>
                </div>
            </div>
        </div>

        <!-- No vehicles remaining message -->
        <div v-else-if="selectedTransport && remainingVehicles === 0"
            class="text-center py-8 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
            All vehicles already added for this transport
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || !canSubmit"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ isEditMode ? 'Saving...' : `Adding ${totalVehiclesToAdd} vehicle(s)...` }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : `Add ${totalVehiclesToAdd} Vehicle(s)` }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Vehicle } from '@/types/vehicle'
import type { Transport } from '@/types/transport'
import { useVehiclesStore } from '@/stores/vehicles'
import { useTransportsStore } from '@/stores/transports'
import { useBrokersStore } from '@/stores/brokers'
import { push } from 'notivue'
import VehicleFields from './VehicleFields.vue'

interface VehicleRow {
    vehicle_number: string
    broker_id: number | null
    joma_cost: number
    vehicle_cost: number
    other_cost: number
    labour_cost: number
    demarage_cost: number
    notes: string
}

const props = defineProps<{
    vehicle?: Vehicle | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'vehicle-created': []
    'vehicle-updated': []
    'cancel': []
}>()

const vehiclesStore = useVehiclesStore()
const transportsStore = useTransportsStore()
const brokersStore = useBrokersStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.vehicle)

const transportOptions = computed(() => transportsStore.transports)
const brokerOptions = computed(() => brokersStore.brokers)

const createEmptyVehicle = (): VehicleRow => ({
    vehicle_number: '',
    broker_id: null,
    joma_cost: 0,
    vehicle_cost: 0,
    other_cost: 0,
    labour_cost: 0,
    demarage_cost: 0,
    notes: '',
})

const form = ref({
    transport_id: null as number | null,
    main_vehicle: createEmptyVehicle(),
    additional_vehicles: [] as VehicleRow[],
})

const selectedTransport = computed(() => {
    if (!form.value.transport_id) return null
    return transportsStore.getTransportById(form.value.transport_id) || null
})

const currentVehicleCount = computed(() => {
    if (!form.value.transport_id) return 0
    return vehiclesStore.getVehiclesByTransportId(form.value.transport_id).length
})

const remainingVehicles = computed(() => {
    if (!selectedTransport.value) return 0
    return Math.max(0, selectedTransport.value.vehicle_quantity - currentVehicleCount.value)
})

const progressPercentage = computed(() => {
    if (!selectedTransport.value || selectedTransport.value.vehicle_quantity === 0) return 0
    return Math.min((currentVehicleCount.value / selectedTransport.value.vehicle_quantity) * 100, 100)
})

const totalVehiclesToAdd = computed(() => 1 + form.value.additional_vehicles.length)

const canSubmit = computed(() => {
    if (isEditMode.value) {
        return !!form.value.transport_id && !!form.value.main_vehicle.vehicle_number.trim() && !!form.value.main_vehicle.broker_id
    }
    if (!form.value.transport_id) return false
    if (!form.value.main_vehicle.vehicle_number.trim()) return false
    if (!form.value.main_vehicle.broker_id) return false
    const allValid = form.value.additional_vehicles.every(v => v.vehicle_number.trim() !== '' && v.broker_id !== null)
    if (!allValid) return false
    if (form.value.additional_vehicles.length > remainingVehicles.value - 1) return false
    return true
})

const getTransportDisplayName = (transport: Transport): string => {
    const fromTo = `${transport.from_location}${transport.to_location ? ` ৳ ${transport.to_location}` : ''}`
    return `${fromTo} (${transport.vehicle_quantity} vehicles)`
}

const initializeAdditionalVehicles = () => {
    if (!selectedTransport.value) return
    const remaining = remainingVehicles.value - 1
    const currentAdditional = form.value.additional_vehicles.length

    if (remaining <= 0) {
        form.value.additional_vehicles = []
        return
    }

    if (currentAdditional < remaining) {
        const toAdd = remaining - currentAdditional
        for (let i = 0; i < toAdd; i++) {
            form.value.additional_vehicles.push(createEmptyVehicle())
        }
    } else if (currentAdditional > remaining) {
        form.value.additional_vehicles = form.value.additional_vehicles.slice(0, remaining)
    }
}

const addAdditionalVehicle = () => {
    const remaining = remainingVehicles.value - 1 - form.value.additional_vehicles.length
    if (remaining > 0) {
        form.value.additional_vehicles.push(createEmptyVehicle())
    } else {
        push.warning('Cannot add more vehicles than the transport allows')
    }
}

const removeAdditionalVehicle = (index: number) => {
    if (form.value.additional_vehicles.length > 1) {
        form.value.additional_vehicles.splice(index, 1)
    }
}

const onTransportChange = () => {
    form.value.main_vehicle = createEmptyVehicle()
    form.value.additional_vehicles = []
    if (selectedTransport.value) {
        initializeAdditionalVehicles()
    }
}

const initializeForm = () => {
    if (props.vehicle && isEditMode.value) {
        form.value = {
            transport_id: props.vehicle.transport_id,
            main_vehicle: {
                vehicle_number: props.vehicle.vehicle_number || '',
                broker_id: props.vehicle.broker_id || null,
                joma_cost: props.vehicle.joma_cost || 0,
                vehicle_cost: props.vehicle.vehicle_cost || 0,
                other_cost: props.vehicle.other_cost || 0,
                labour_cost: props.vehicle.labour_cost || 0,
                demarage_cost: props.vehicle.demarage_cost || 0,
                notes: props.vehicle.notes || '',
            },
            additional_vehicles: [],
        }
    } else {
        form.value = {
            transport_id: null,
            main_vehicle: createEmptyVehicle(),
            additional_vehicles: [],
        }
    }
}

watch(() => props.vehicle, initializeForm, { immediate: true })

watch(() => remainingVehicles.value, (newRemaining) => {
    if (!isEditMode.value && form.value.transport_id) {
        if (newRemaining <= 1) {
            form.value.additional_vehicles = []
        } else {
            initializeAdditionalVehicles()
        }
    }
})

const resetForm = () => {
    if (isEditMode.value && props.vehicle) {
        initializeForm()
    } else {
        form.value = {
            transport_id: null,
            main_vehicle: createEmptyVehicle(),
            additional_vehicles: [],
        }
    }
}

const submit = async () => {
    if (!form.value.transport_id) {
        push.error('Please select a transport')
        return
    }
    if (!form.value.main_vehicle.vehicle_number.trim()) {
        push.error('Main vehicle number is required')
        return
    }
    if (!form.value.main_vehicle.broker_id) {
        push.error('Broker is required for the main vehicle')
        return
    }

    const invalidVehicles = form.value.additional_vehicles.filter(v => !v.vehicle_number.trim() || !v.broker_id)
    if (invalidVehicles.length > 0) {
        push.error(`Please fill in vehicle numbers and brokers for all ${form.value.additional_vehicles.length} additional vehicles`)
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.vehicle) {
            const success = await vehiclesStore.updateVehicle(props.vehicle.id, {
                vehicle_number: form.value.main_vehicle.vehicle_number.trim(),
                broker_id: form.value.main_vehicle.broker_id,
                joma_cost: form.value.main_vehicle.joma_cost,
                vehicle_cost: form.value.main_vehicle.vehicle_cost,
                other_cost: form.value.main_vehicle.other_cost,
                labour_cost: form.value.main_vehicle.labour_cost,
                demarage_cost: form.value.main_vehicle.demarage_cost,
                notes: form.value.main_vehicle.notes.trim() || null,
            })

            if (success) {
                push.success('Vehicle updated successfully!')
                emit('vehicle-updated')
            }
        } else {
            const newVehicle = await vehiclesStore.createVehicle({
                transport_id: form.value.transport_id,
                vehicle_number: form.value.main_vehicle.vehicle_number.trim(),
                broker_id: form.value.main_vehicle.broker_id,
                joma_cost: form.value.main_vehicle.joma_cost,
                vehicle_cost: form.value.main_vehicle.vehicle_cost,
                other_cost: form.value.main_vehicle.other_cost,
                labour_cost: form.value.main_vehicle.labour_cost,
                demarage_cost: form.value.main_vehicle.demarage_cost,
                notes: form.value.main_vehicle.notes.trim() || null,
            })

            if (!newVehicle) {
                submitting.value = false
                return
            }

            let createdCount = 0
            let failedCount = 0

            for (const vehicle of form.value.additional_vehicles) {
                if (!vehicle.vehicle_number.trim() || !vehicle.broker_id) {
                    failedCount++
                    continue
                }

                const result = await vehiclesStore.createVehicle({
                    transport_id: form.value.transport_id,
                    vehicle_number: vehicle.vehicle_number.trim(),
                    broker_id: vehicle.broker_id,
                    joma_cost: vehicle.joma_cost,
                    vehicle_cost: vehicle.vehicle_cost,
                    other_cost: vehicle.other_cost,
                    labour_cost: vehicle.labour_cost,
                    demarage_cost: vehicle.demarage_cost,
                    notes: vehicle.notes.trim() || null,
                })

                if (result) createdCount++
                else failedCount++
            }

            if (createdCount > 0 && failedCount === 0) {
                push.success(`${createdCount + 1} vehicle(s) added successfully!`)
            } else if (createdCount > 0 && failedCount > 0) {
                push.warning(`${createdCount} vehicle(s) added, ${failedCount} failed`)
            } else {
                push.success('Main vehicle created successfully!')
            }

            resetForm()
            emit('vehicle-created')
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update vehicle' : 'Failed to create vehicle')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await Promise.all([
        transportsStore.fetchTransports(),
        vehiclesStore.fetchVehicles(),
        brokersStore.fetchBrokers(),
    ])
})
</script>