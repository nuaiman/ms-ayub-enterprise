<!-- src/components/features/salaries/SalaryForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Employee Selection -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Employee Information
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Employee <span class="text-(--color-red)">*</span>
                    </label>
                    <select v-model="form.employee_id" required :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed">
                        <option value="">Select an employee</option>
                        <option v-for="user in employeeOptions" :key="user.id" :value="user.id">
                            {{ user.name }} ({{ user.role }})
                        </option>
                    </select>
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Employee cannot be changed
                    </p>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Month <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.month_year" type="month" required :disabled="isEditMode"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
                    <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Month cannot be changed</p>
                </div>
            </div>
        </div>

        <!-- Salary Details -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Salary Details
            </h3>

            <!-- Employee Basic Salary (read-only) -->
            <div class="p-4 rounded-lg bg-(--color-muted-bg)/30 border border-(--color-border) mb-4">
                <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-(--color-text-primary)">Employee Basic Salary</span>
                    <span class="text-lg font-semibold text-(--color-text-primary)">{{
                        formatCurrency(employeeBasicSalary) }}</span>
                </div>
                <p class="text-xs text-(--color-text-secondary) mt-1">This is the employee's monthly salary from their
                    profile</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Bonus</label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.bonus" type="number" step="0.01" min="0" placeholder="0.00"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Deductions</label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
                        <input v-model.number="form.deductions" type="number" step="0.01" min="0" placeholder="0.00"
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>
            </div>
        </div>

        <!-- Total Calculation -->
        <div class="p-4 rounded-lg bg-(--color-muted-bg)/50 border border-(--color-border)">
            <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-(--color-text-primary)">Total Salary</span>
                <span class="text-xl font-bold text-(--color-blue)">{{ formatCurrency(totalAmount) }}</span>
            </div>
            <div class="flex items-center justify-between mt-1 text-xs text-(--color-text-secondary)">
                <span>Basic Salary + Bonus - Deductions</span>
                <span>{{ formatCurrency(employeeBasicSalary) }} + {{ formatCurrency(form.bonus || 0) }} - {{
                    formatCurrency(form.deductions || 0) }}</span>
            </div>
        </div>

        <!-- Notes -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Additional Information
            </h3>
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Notes</label>
                <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this salary"
                    class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none"></textarea>
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Salary' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Salary } from '@/types/salary'
import { useSalariesStore } from '@/stores/salaries'
import { useUsersStore } from '@/stores/users'
import { formatCurrency } from '@/utils/currency'
import { push } from 'notivue'

const props = defineProps<{
    salary?: Salary | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'salary-created': []
    'salary-updated': []
    'cancel': []
}>()

const salariesStore = useSalariesStore()
const usersStore = useUsersStore()

const isEditMode = computed(() => props.mode === 'edit' || !!props.salary)

const employeeOptions = computed(() => {
    return usersStore.users
        .filter(u => u.is_active)
        .map(u => ({
            id: u.id,
            name: u.name,
            role: u.role
        }))
})

const employeeBasicSalary = computed(() => {
    if (!form.value.employee_id) return 0
    const employee = usersStore.users.find(u => u.id === form.value.employee_id)
    return employee?.monthly_salary || 0
})

const form = ref({
    employee_id: null as number | null,
    month_year: '',
    bonus: 0,
    deductions: 0,
    notes: '',
})

const submitting = ref(false)

const totalAmount = computed(() => {
    return (employeeBasicSalary.value || 0) + (form.value.bonus || 0) - (form.value.deductions || 0)
})

const initializeForm = () => {
    if (props.salary) {
        form.value = {
            employee_id: props.salary.employee_id,
            month_year: props.salary.month_year,
            bonus: props.salary.bonus || 0,
            deductions: props.salary.deductions || 0,
            notes: props.salary.notes || '',
        }
    } else {
        const now = new Date()
        const monthYear = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
        form.value = {
            employee_id: null,
            month_year: monthYear,
            bonus: 0,
            deductions: 0,
            notes: '',
        }
    }
}

watch(() => props.salary, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.salary) {
        initializeForm()
    } else {
        const now = new Date()
        const monthYear = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
        form.value = {
            employee_id: null,
            month_year: monthYear,
            bonus: 0,
            deductions: 0,
            notes: '',
        }
    }
}

const submit = async () => {
    if (!form.value.employee_id) {
        push.error('Please select an employee')
        return
    }

    if (!form.value.month_year) {
        push.error('Please select a month')
        return
    }

    if (form.value.bonus < 0) {
        push.error('Bonus cannot be negative')
        return
    }

    if (form.value.deductions < 0) {
        push.error('Deductions cannot be negative')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.salary) {
            const success = await salariesStore.updateSalary(props.salary.id, {
                bonus: form.value.bonus,
                deductions: form.value.deductions,
                notes: form.value.notes || null,
            })

            if (success) {
                push.success('Salary updated successfully!')
                emit('salary-updated')
            }
        } else {
            const newSalary = await salariesStore.createSalary({
                employee_id: form.value.employee_id,
                month_year: form.value.month_year,
                bonus: form.value.bonus,
                deductions: form.value.deductions,
                notes: form.value.notes || null,
            })

            if (newSalary) {
                push.success('Salary created successfully!')
                resetForm()
                emit('salary-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update salary' : 'Failed to create salary')
    } finally {
        submitting.value = false
    }
}
</script>