<!-- src/components/features/customers/CustomerForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Customer Information -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Customer Information
            </h3>
            <div class="space-y-4">
                <!-- Company Name -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Company Name
                    </label>
                    <input v-model="form.company_name" type="text" placeholder="Enter company name"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Contact Person -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Contact Person
                    </label>
                    <input v-model="form.contact_person" type="text" placeholder="Enter contact person name"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Phone -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Phone <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.phone" type="tel" placeholder="Enter phone number" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Email -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Email
                    </label>
                    <input v-model="form.email" type="email" placeholder="Enter email address"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>

                <!-- Address -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Address
                    </label>
                    <textarea v-model="form.address" rows="3" placeholder="Enter full address"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
                </div>

                <!-- Notes -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Notes
                    </label>
                    <textarea v-model="form.notes" rows="2" placeholder="Enter any notes about this customer"
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Customer' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Customer } from '@/types/customer'
import { useCustomersStore } from '@/stores/customers'
import { push } from 'notivue'

const props = defineProps<{
    customer?: Customer | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'customer-created': []
    'customer-updated': []
    'cancel': []
}>()

const customersStore = useCustomersStore()
const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit' || !!props.customer)

const form = ref({
    company_name: '',
    contact_person: '',
    phone: '',
    email: '',
    address: '',
    notes: '',
})

const initializeForm = () => {
    if (props.customer) {
        form.value = {
            company_name: props.customer.company_name || '',
            contact_person: props.customer.contact_person || '',
            phone: props.customer.phone || '',
            email: props.customer.email || '',
            address: props.customer.address || '',
            notes: props.customer.notes || '',
        }
    } else {
        form.value = {
            company_name: '',
            contact_person: '',
            phone: '',
            email: '',
            address: '',
            notes: '',
        }
    }
}

watch(() => props.customer, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.customer) {
        initializeForm()
    } else {
        form.value = {
            company_name: '',
            contact_person: '',
            phone: '',
            email: '',
            address: '',
            notes: '',
        }
    }
}

const submit = async () => {
    // Validate: at least company_name or contact_person
    if (!form.value.company_name.trim() && !form.value.contact_person.trim()) {
        push.error('Either company name or contact person is required')
        return
    }

    if (!form.value.phone.trim()) {
        push.error('Phone number is required')
        return
    }

    if (form.value.email && !isValidEmail(form.value.email)) {
        push.error('Please enter a valid email address')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.customer) {
            const success = await customersStore.updateCustomer(props.customer.id, {
                company_name: form.value.company_name.trim() || null,
                contact_person: form.value.contact_person.trim() || null,
                phone: form.value.phone.trim(),
                email: form.value.email.trim() || null,
                address: form.value.address.trim() || null,
                notes: form.value.notes.trim() || null,
            })

            if (success) {
                push.success('Customer updated successfully!')
                emit('customer-updated')
            }
        } else {
            const newCustomer = await customersStore.createCustomer({
                company_name: form.value.company_name.trim() || null,
                contact_person: form.value.contact_person.trim() || null,
                phone: form.value.phone.trim(),
                email: form.value.email.trim() || null,
                address: form.value.address.trim() || null,
                notes: form.value.notes.trim() || null,
            })

            if (newCustomer) {
                push.success('Customer created successfully!')
                resetForm()
                emit('customer-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update customer' : 'Failed to create customer')
    } finally {
        submitting.value = false
    }
}

const isValidEmail = (email: string): boolean => {
    const pattern = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/
    return pattern.test(email)
}
</script>