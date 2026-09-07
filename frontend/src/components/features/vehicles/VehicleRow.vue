<!-- src/components/features/vehicles/VehicleRow.vue -->
<template>
    <div class="grid grid-cols-12 items-center w-full py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
        @click="handleView">
        <!-- Vehicle Number - 4 columns -->
        <div class="col-span-4 min-w-0">
            <div class="font-medium text-(--color-text-primary) truncate text-sm">
                {{ vehicle.vehicle_number }}
            </div>
        </div>

        <!-- Transport - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                #{{ vehicle.transport_id }}
            </span>
        </div>

        <!-- Broker - 2 columns -->
        <div class="col-span-2 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ getBrokerName(vehicle.broker_id) }}
            </span>
        </div>

        <!-- Driver - 3 columns -->
        <div class="col-span-3 min-w-0">
            <span class="text-sm text-(--color-text-secondary) truncate block">
                {{ vehicle.driver_name || '—' }}
            </span>
            <span v-if="vehicle.driver_phone" class="text-xs text-(--color-text-secondary) truncate block">
                {{ vehicle.driver_phone }}
            </span>
        </div>

        <!-- Actions - 1 column, right aligned -->
        <div class="col-span-1 flex items-center justify-end relative" @click.stop>
            <button @click="toggleMenu"
                class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200">
                <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="5" r="1.5" />
                    <circle cx="12" cy="12" r="1.5" />
                    <circle cx="12" cy="19" r="1.5" />
                </svg>
            </button>

            <!-- Dropdown -->
            <Transition enter-active-class="transition ease-out duration-200"
                enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition ease-in duration-150"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 translate-y-1">
                <div v-if="isOpen"
                    class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
                    <!-- View -->
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

                    <!-- Edit -->
                    <button @click="handleEdit"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                    </button>

                    <!-- Delete -->
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

            <!-- Backdrop -->
            <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Vehicle } from '@/types/vehicle'
import { useBrokersStore } from '@/stores/brokers'

const props = defineProps<{
    vehicle: Vehicle
}>()

const emit = defineEmits<{
    'view': [vehicle: Vehicle]
    'edit': [vehicle: Vehicle]
    'delete': [vehicle: Vehicle]
    'updated': []
}>()

const brokersStore = useBrokersStore()
const isOpen = ref(false)

const getBrokerName = (id: number | null): string => {
    if (!id) return '—'
    return brokersStore.getBrokerName(id)
}

const toggleMenu = () => {
    isOpen.value = !isOpen.value
}

const closeMenu = () => {
    isOpen.value = false
}

const handleView = () => {
    closeMenu()
    emit('view', props.vehicle)
}

const handleEdit = () => {
    closeMenu()
    emit('edit', props.vehicle)
}

const handleDelete = () => {
    closeMenu()
    emit('delete', props.vehicle)
}
</script>