<!-- src/components/features/transports/TransportForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Transport Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Transport Information
            </h3>

            <TransportFields v-model:customer-id="form.customer_id" v-model:from-location="form.from_location"
                v-model:to-location="form.to_location" v-model:vehicle-quantity="form.vehicle_quantity"
                v-model:transport-type="form.transport_type" v-model:transport-date="form.transport_date"
                v-model:office-commission-amount="form.office_commission_amount" v-model:notes="form.notes"
                :customer-options="customerOptions" :disabled="submitting || !!justCreatedId" :required="true" />
        </div>

        <!-- Customer Billing -->
        <div v-if="!justCreatedId" class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Billing
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Charge Unit <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.customer_charge_unit" :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="vehicle">Vehicle</option>
                        <option value="weight">Weight</option>
                        <option value="quantity">Quantity</option>
                    </select>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Total Unit <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model.number="form.customer_total_unit" type="number" step="0.01" min="0"
                        placeholder="0.00" :disabled="submitting"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Charge Per Unit <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.customer_charge_per_unit" type="number" step="0.01" min="0"
                            placeholder="0.00" :disabled="submitting"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    </div>
                </div>
            </div>

            <!-- Live total charge preview -->
            <div class="mt-4 p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-(--color-text-primary)">Customer Total Charge</span>
                    <span class="text-lg font-bold text-(--color-blue)">
                        {{ formatCurrency(customerTotalCharge) }}
                    </span>
                </div>
                <p class="text-xs text-(--color-text-secondary) mt-1">
                    {{ formatCurrency(form.customer_charge_per_unit) }} ৳— {{ form.customer_total_unit }}
                </p>
            </div>
        </div>

        <!-- Customer payment summary (edit mode) -->
        <div v-if="!justCreatedId && isEditMode && transport" class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Payment
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Paid</p>
                    <p class="text-lg font-bold text-(--color-green)">
                        {{ formatCurrency(transport.customer_total_paid) }}
                    </p>
                </div>
                <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border)">
                    <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Outstanding</p>
                    <p class="text-lg font-bold"
                        :class="customerOutstanding > 0 ? 'text-(--color-red)' : 'text-(--color-green)'">
                        {{ formatCurrency(customerOutstanding) }}
                    </p>
                </div>
            </div>
            <p class="text-xs text-(--color-text-secondary) mt-2">
                Record customer payments from the Customer Transport Bills page.
            </p>
        </div>

        <!-- Vehicles Section (hidden when retrying image) -->
        <div v-if="!justCreatedId" class="border-t border-(--color-border) pt-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider">
                    Vehicles
                </h3>
                <div class="flex items-center gap-3">
                    <span class="text-xs text-(--color-text-secondary) bg-(--color-muted-bg) px-2.5 py-1 rounded-full">
                        {{ totalVehicleCount }} / {{ form.vehicle_quantity }} vehicle(s)
                    </span>
                </div>
            </div>

            <!-- Existing vehicles (edit mode) -->
            <div v-if="isEditMode && existingVehicles.length > 0" class="mb-4">
                <h4 class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider mb-2">
                    Existing Vehicles
                </h4>
                <div class="space-y-2">
                    <div v-for="vehicle in existingVehicles" :key="vehicle.id"
                        class="flex items-center justify-between p-3 rounded-lg bg-(--color-muted-bg)/20 border border-(--color-border)">
                        <div class="flex items-center gap-3">
                            <span class="text-sm font-medium text-(--color-text-primary)">
                                {{ vehicle.vehicle_number }}
                            </span>
                            <span class="text-xs text-(--color-text-secondary)">
                                Broker: {{ getBrokerName(vehicle.broker_id) }}
                            </span>
                        </div>
                        <div class="flex items-center gap-3">
                            <span class="text-xs text-(--color-text-secondary)">
                                Joma: {{ formatCurrency(vehicle.joma_cost) }}
                            </span>
                            <span class="text-xs text-(--color-text-secondary)">
                                Vehicle: {{ formatCurrency(vehicle.vehicle_cost) }}
                            </span>
                            <button type="button" @click="removeExistingVehicle(vehicle.id)"
                                class="text-xs text-(--color-red) hover:bg-(--color-red)/10 px-2 py-1 rounded transition-colors">
                                Remove
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- New vehicle inline sub-forms -->
            <div v-if="form.vehicles.length > 0" class="space-y-4">
                <div class="flex items-center justify-between">
                    <h4 class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">
                        {{ isEditMode ? 'Add New Vehicles' : 'Vehicle Details' }}
                    </h4>
                    <span class="text-xs text-(--color-text-secondary)">
                        {{ form.vehicles.length }} to add
                    </span>
                </div>

                <div class="space-y-4">
                    <div v-for="(vehicle, index) in form.vehicles" :key="index"
                        class="relative p-4 rounded-lg border border-(--color-border) bg-(--color-muted-bg)/10">
                        <button type="button" @click="removeNewVehicle(index)"
                            class="absolute top-2 right-2 p-1 rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-colors">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>

                        <h4 class="text-sm font-medium text-(--color-text-primary) flex items-center gap-2 mb-3">
                            <span
                                class="w-6 h-6 rounded-full bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center text-xs font-bold">
                                {{ index + 1 }}
                            </span>
                            Vehicle #{{ index + 1 }}
                        </h4>

                        <VehicleFields v-model:vehicle-number="vehicle.vehicle_number"
                            v-model:broker-id="vehicle.broker_id" v-model:joma-cost="vehicle.joma_cost"
                            v-model:vehicle-cost="vehicle.vehicle_cost" v-model:other-cost="vehicle.other_cost"
                            v-model:labour-cost="vehicle.labour_cost" v-model:demarage-cost="vehicle.demarage_cost"
                            v-model:notes="vehicle.notes" :broker-options="brokerOptions" :disabled="submitting"
                            :required="false" />
                    </div>
                </div>
            </div>

            <div v-else
                class="text-center py-6 text-sm text-(--color-text-secondary) border border-dashed border-(--color-border) rounded-lg">
                No new vehicles to add.
            </div>

            <div class="mt-4">
                <button type="button" @click="addNewVehicle" :disabled="!canAddMoreVehicles"
                    class="w-full py-3 text-sm font-medium rounded-lg border-2 border-dashed border-(--color-border) text-(--color-text-secondary) hover:border-(--color-blue) hover:text-(--color-blue) hover:bg-(--color-blue)/5 transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-40 disabled:cursor-not-allowed">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                    Add Vehicle
                </button>
                <p v-if="!canAddMoreVehicles" class="text-xs text-(--color-text-secondary) mt-2 text-center">
                    Vehicle quantity limit reached
                </p>
            </div>
        </div>

        <!-- Image -->
        <div class="border-t border-(--color-border) pt-6">
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Transport Image <span class="normal-case text-xs font-normal">(Optional)</span>
            </h3>
            <div class="flex flex-col sm:flex-row items-start gap-4">
                <div class="shrink-0">
                    <div v-if="imagePreview"
                        class="relative w-24 h-24 rounded-lg overflow-hidden border border-(--color-border)">
                        <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                        <button type="button" @click="removeImage"
                            class="absolute top-1 right-1 w-5 h-5 bg-(--color-red) text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">৳—</button>
                    </div>
                    <div v-else
                        class="w-24 h-24 rounded-lg border-2 border-dashed border-(--color-border) flex items-center justify-center bg-(--color-muted-bg)">
                        <svg class="w-10 h-10 text-(--color-text-secondary)" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                    </div>
                </div>

                <div class="flex-1 space-y-3">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Upload
                            Image</label>
                        <div class="flex flex-wrap gap-2">
                            <label
                                class="px-4 py-2 text-sm font-medium rounded-lg cursor-pointer bg-(--color-muted-bg) border border-(--color-border) hover:bg-(--color-muted-bg)/70 transition-all duration-200 inline-flex items-center gap-2"
                                :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                                </svg>
                                {{ uploading ? 'Uploading...' : 'Choose File' }}
                                <input type="file" accept="image/*" class="hidden" @change="handleFileSelect"
                                    :disabled="uploading" />
                            </label>
                            <button v-if="imagePreview" type="button" @click="removeImage"
                                class="px-4 py-2 text-sm font-medium rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-all duration-200">
                                Remove
                            </button>
                        </div>
                    </div>

                    <div v-if="uploading" class="w-full bg-(--color-muted-bg) rounded-full h-1.5 overflow-hidden">
                        <div class="bg-(--color-blue) h-full rounded-full transition-all duration-300"
                            :style="{ width: uploadProgress + '%' }"></div>
                    </div>

                    <p v-if="uploadError" class="text-xs text-(--color-red)">{{ uploadError }}</p>
                    <p v-else-if="justCreatedId" class="text-xs text-(--color-yellow)">
                        Transport was created but the image upload failed. Retry below or remove the image.
                    </p>
                    <p v-else-if="isEditMode && props.transport?.image_url"
                        class="text-xs text-(--color-text-secondary)">
                        Current image will be replaced
                    </p>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')" :disabled="submitting || uploading"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
                Cancel
            </button>
            <button type="submit" :disabled="submitting || uploading || (!justCreatedId && !canSubmit)"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting || uploading" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ justCreatedId ? 'Retrying image...' : (isEditMode ? 'Saving...' : 'Creating...') }}
                </span>
                <span v-else>
                    {{ justCreatedId ? 'Retry Image Upload' : isEditMode ? 'Save Changes' : 'Create Transport' }}
                </span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { Transport, CustomerChargeUnit } from '@/types/transport'
import { useTransportsStore } from '@/stores/transports'
import { useCustomersStore } from '@/stores/customers'
import { useBrokersStore } from '@/stores/brokers'
import { useVehiclesStore } from '@/stores/vehicles'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { formatDateForBackend } from '@/utils/date'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'
import TransportFields from './TransportFields.vue'
import VehicleFields from '@/components/features/vehicles/VehicleFields.vue'

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
    transport?: Transport | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'transport-created': []
    'transport-updated': []
    'cancel': []
}>()

const transportsStore = useTransportsStore()
const customersStore = useCustomersStore()
const brokersStore = useBrokersStore()
const vehiclesStore = useVehiclesStore()

const submitting = ref(false)
const isEditMode = computed(() => props.mode === 'edit' || !!props.transport)

const customerOptions = computed(() => customersStore.customers)
const brokerOptions = computed(() => brokersStore.brokers)

const justCreatedId = ref<number | null>(null)

// Image state
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)
const imageToDelete = ref(false)

const clearImageState = () => {
    imageFile.value = null
    imagePreview.value = null
    uploadError.value = null
    uploadProgress.value = 0
    imageToDelete.value = false
    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null
    if (fileInput) fileInput.value = ''
}

const handleFileSelect = (event: Event) => {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (!file) return

    if (!file.type.startsWith('image/')) {
        uploadError.value = 'Please select an image file'
        return
    }
    if (file.size > 5 * 1024 * 1024) {
        uploadError.value = 'Image size should be less than 5MB'
        return
    }

    uploadError.value = null
    imageFile.value = file
    imageToDelete.value = false

    const reader = new FileReader()
    reader.onload = (e) => {
        imagePreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
}

const removeImage = () => {
    imageFile.value = null
    imagePreview.value = null
    uploadError.value = null
    uploadProgress.value = 0
    imageToDelete.value = true
    const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null
    if (fileInput) fileInput.value = ''
}

const today = new Date().toISOString().slice(0, 10)

const form = ref({
    customer_id: null as number | null,
    from_location: '',
    to_location: '',
    vehicle_quantity: 1,
    transport_type: null as 'local' | 'district' | null,
    transport_date: today,
    office_commission_amount: 0,
    customer_charge_unit: 'vehicle' as CustomerChargeUnit,
    customer_total_unit: 0,
    customer_charge_per_unit: 0,
    notes: '',
    vehicles: [] as VehicleRow[],
})

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

const existingVehicles = computed(() => {
    if (!props.transport) return []
    return vehiclesStore.getVehiclesByTransportId(props.transport.id)
})

const totalVehicleCount = computed(() => existingVehicles.value.length + form.value.vehicles.length)

const canAddMoreVehicles = computed(() => {
    return totalVehicleCount.value < form.value.vehicle_quantity
})

const customerTotalCharge = computed(() => {
    return (form.value.customer_charge_per_unit || 0) * (form.value.customer_total_unit || 0)
})

const customerOutstanding = computed(() => {
    if (!props.transport) return 0
    return Math.max(0, (props.transport.customer_total_charge || 0) - (props.transport.customer_total_paid || 0))
})

const canSubmit = computed(() => {
    if (!form.value.from_location.trim()) return false
    const allHaveNumbers = form.value.vehicles.every(v => v.vehicle_number.trim() !== '')
    if (!allHaveNumbers) return false
    const allHaveBrokers = form.value.vehicles.every(v => v.broker_id !== null)
    if (!allHaveBrokers) return false
    return true
})

const addNewVehicle = () => {
    if (!canAddMoreVehicles.value) {
        push.warning('Vehicle quantity limit reached')
        return
    }
    form.value.vehicles.push(createEmptyVehicle())
}

const removeNewVehicle = (index: number) => {
    form.value.vehicles.splice(index, 1)
}

const removeExistingVehicle = async (vehicleId: number) => {
    if (!confirm('Are you sure you want to remove this vehicle from the transport?')) return
    const success = await vehiclesStore.deleteVehicle(vehicleId)
    if (success) {
        push.success('Vehicle removed from transport')
        await vehiclesStore.fetchVehicles()
    }
}

const getBrokerName = (id: number): string => brokersStore.getBrokerName(id)

const initializeForm = () => {
    if (props.transport) {
        const date = props.transport.transport_date
            ? new Date(props.transport.transport_date).toISOString().slice(0, 10)
            : today

        form.value = {
            customer_id: props.transport.customer_id || null,
            from_location: props.transport.from_location || '',
            to_location: props.transport.to_location || '',
            vehicle_quantity: props.transport.vehicle_quantity || 1,
            transport_type: props.transport.transport_type || null,
            transport_date: date,
            office_commission_amount: props.transport.office_commission_amount || 0,
            customer_charge_unit: props.transport.customer_charge_unit || 'vehicle',
            customer_total_unit: props.transport.customer_total_unit || 0,
            customer_charge_per_unit: props.transport.customer_charge_per_unit || 0,
            notes: props.transport.notes || '',
            vehicles: [],
        }

        justCreatedId.value = null
        clearImageState()
        if (isEditMode.value && props.transport.image_url) {
            imagePreview.value = getImageUrl(props.transport.image_url) || null
        }
    } else {
        form.value = {
            customer_id: null,
            from_location: '',
            to_location: '',
            vehicle_quantity: 1,
            transport_type: null,
            transport_date: today,
            office_commission_amount: 0,
            customer_charge_unit: 'vehicle',
            customer_total_unit: 0,
            customer_charge_per_unit: 0,
            notes: '',
            vehicles: [createEmptyVehicle()],
        }
        justCreatedId.value = null
        clearImageState()
    }
}

watch(() => props.transport, initializeForm, { immediate: true })

const resetForm = () => {
    initializeForm()
}

const submit = async () => {
    // Retry image path
    if (justCreatedId.value !== null) {
        if (!imageFile.value) {
            resetForm()
            emit('transport-created')
            return
        }
        uploading.value = true
        uploadError.value = null
        uploadProgress.value = 0
        try {
            const url = await uploadImage('transports', justCreatedId.value, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                return
            }
            push.success('Transport and image created successfully!')
            resetForm()
            emit('transport-created')
        } finally {
            uploading.value = false
        }
        return
    }

    // Validate
    if (!form.value.from_location.trim()) {
        push.error('From location is required')
        return
    }
    if (!form.value.customer_id) {
        push.error('Customer is required')
        return
    }
    if (!form.value.transport_date) {
        push.error('Transport date is required')
        return
    }
    if (form.value.vehicle_quantity < 0) {
        push.error('Vehicle quantity cannot be negative')
        return
    }

    const invalidNewVehicles = form.value.vehicles.filter(v => !v.vehicle_number.trim())
    if (invalidNewVehicles.length > 0) {
        push.error(`Please enter vehicle numbers for all ${form.value.vehicles.length} new vehicles`)
        return
    }
    const missingBrokers = form.value.vehicles.filter(v => v.broker_id === null)
    if (missingBrokers.length > 0) {
        push.error(`Please select a broker for all ${form.value.vehicles.length} new vehicles`)
        return
    }

    submitting.value = true
    const transportDate = formatDateForBackend(form.value.transport_date) || ''

    try {
        let transportId: number

        if (isEditMode.value && props.transport) {
            const success = await transportsStore.updateTransport(props.transport.id, {
                customer_id: form.value.customer_id,
                from_location: form.value.from_location.trim(),
                to_location: form.value.to_location?.trim() || null,
                vehicle_quantity: form.value.vehicle_quantity,
                transport_type: form.value.transport_type,
                transport_date: transportDate,
                office_commission_amount: form.value.office_commission_amount,
                customer_charge_unit: form.value.customer_charge_unit,
                customer_total_unit: form.value.customer_total_unit,
                customer_charge_per_unit: form.value.customer_charge_per_unit,
                notes: form.value.notes.trim() || null,
            })

            if (!success) {
                submitting.value = false
                return
            }

            transportId = props.transport.id
            push.success('Transport updated successfully!')
        } else {
            const newTransport = await transportsStore.createTransport({
                customer_id: form.value.customer_id,
                from_location: form.value.from_location.trim(),
                to_location: form.value.to_location?.trim() || null,
                vehicle_quantity: form.value.vehicle_quantity,
                transport_type: form.value.transport_type,
                transport_date: transportDate,
                office_commission_amount: form.value.office_commission_amount,
                customer_charge_unit: form.value.customer_charge_unit,
                customer_total_unit: form.value.customer_total_unit,
                customer_charge_per_unit: form.value.customer_charge_per_unit,
                notes: form.value.notes.trim() || null,
            })

            if (!newTransport) {
                submitting.value = false
                return
            }

            transportId = newTransport.id
            push.success('Transport created successfully!')
        }

        // Create new vehicles
        let createdCount = 0
        let failedCount = 0

        for (const vehicle of form.value.vehicles) {
            if (!vehicle.vehicle_number.trim() || vehicle.broker_id === null) {
                failedCount++
                continue
            }

            const result = await vehiclesStore.createVehicle({
                transport_id: transportId,
                vehicle_number: vehicle.vehicle_number.trim(),
                broker_id: vehicle.broker_id,
                joma_cost: vehicle.joma_cost,
                vehicle_cost: vehicle.vehicle_cost,
                other_cost: vehicle.other_cost,
                labour_cost: vehicle.labour_cost,
                demarage_cost: vehicle.demarage_cost,
                notes: vehicle.notes?.trim() || null,
            })

            if (result) createdCount++
            else failedCount++
        }

        if (createdCount > 0 && failedCount === 0) {
            push.success(`${createdCount} vehicle(s) added successfully!`)
        } else if (createdCount > 0 && failedCount > 0) {
            push.warning(`${createdCount} vehicle(s) added, ${failedCount} failed`)
        }

        // Image handling
        if (imageFile.value) {
            uploading.value = true
            uploadProgress.value = 0
            uploadError.value = null
            const url = await uploadImage('transports', transportId, imageFile.value, (p) => {
                uploadProgress.value = p
            })
            uploading.value = false
            if (!url) {
                uploadError.value = 'Failed to upload image. Please try again.'
                if (!isEditMode.value) {
                    justCreatedId.value = transportId
                }
                push.warning('Transport saved but image upload failed. Retry below.')
                await Promise.all([transportsStore.fetchTransports(), vehiclesStore.fetchVehicles()])
                return
            }
        } else if (isEditMode.value && props.transport && imageToDelete.value && props.transport.image_url) {
            await deleteImage('transports', props.transport.id, props.transport.image_url)
        }

        if (isEditMode.value) {
            emit('transport-updated')
        } else {
            emit('transport-created')
        }

        resetForm()
        await Promise.all([transportsStore.fetchTransports(), vehiclesStore.fetchVehicles()])
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update transport' : 'Failed to create transport')
    } finally {
        submitting.value = false
    }
}

onMounted(async () => {
    await Promise.all([
        customersStore.fetchCustomers(),
        brokersStore.fetchBrokers(),
        transportsStore.fetchTransports(),
        vehiclesStore.fetchVehicles(),
    ])
})
</script>