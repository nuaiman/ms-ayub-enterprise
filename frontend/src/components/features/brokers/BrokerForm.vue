<!-- src/components/features/brokers/BrokerForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Broker Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Broker Information
            </h3>
            <div class="space-y-4">
                <!-- Name -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Name <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.name" type="text" placeholder="Enter broker name" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Phone -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Phone
                    </label>
                    <input v-model="form.phone" type="tel" placeholder="Enter phone number"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Notes -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this broker"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
            <button type="button" @click="emit('cancel')"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
                Cancel
            </button>
            <button type="submit" :disabled="submitting"
                class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
                <span v-if="submitting" class="inline-flex items-center justify-center gap-2">
                    <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                        <path class="opacity-75" fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ isEditMode ? 'Saving...' : 'Creating...' }}
                </span>
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Broker' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Broker } from '@/types/broker'
import { useBrokersStore } from '@/stores/brokers'
import { push } from 'notivue'

const props = defineProps<{
    broker?: Broker | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'broker-created': []
    'broker-updated': []
    'cancel': []
}>()

const brokersStore = useBrokersStore()
const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit' || !!props.broker)

const form = ref({
    name: '',
    phone: '',
    notes: '',
})

const initializeForm = () => {
    if (props.broker) {
        form.value = {
            name: props.broker.name || '',
            phone: props.broker.phone || '',
            notes: props.broker.notes || '',
        }
    } else {
        form.value = {
            name: '',
            phone: '',
            notes: '',
        }
    }
}

watch(() => props.broker, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.broker) {
        initializeForm()
    } else {
        form.value = {
            name: '',
            phone: '',
            notes: '',
        }
    }
}

const submit = async () => {
    if (!form.value.name.trim()) {
        push.error('Name is required')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.broker) {
            const success = await brokersStore.updateBroker(props.broker.id, {
                name: form.value.name.trim(),
                phone: form.value.phone.trim() || null,
                notes: form.value.notes.trim() || null,
            })

            if (success) {
                push.success('Broker updated successfully!')
                emit('broker-updated')
            }
        } else {
            const newBroker = await brokersStore.createBroker({
                name: form.value.name.trim(),
                phone: form.value.phone.trim() || null,
                notes: form.value.notes.trim() || null,
            })

            if (newBroker) {
                push.success('Broker created successfully!')
                resetForm()
                emit('broker-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update broker' : 'Failed to create broker')
    } finally {
        submitting.value = false
    }
}
</script>