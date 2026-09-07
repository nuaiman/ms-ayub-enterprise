<!-- src/components/features/expenses/ExpenseForm.vue -->
<template>
    <form @submit.prevent="submit" class="space-y-6">
        <!-- Expense Details -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Expense Information
            </h3>
            <div class="space-y-4">
                <!-- Title with Autocomplete -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Title <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative" ref="dropdownRef">
                        <!-- Input -->
                        <input ref="inputRef" v-model="searchTerm" type="text" placeholder="Enter expense title"
                            required @focus="openDropdown" @input="onInput" @keydown.down="selectNext"
                            @keydown.up="selectPrevious" @keydown.enter="handleEnter" @keydown.escape="closeDropdown"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />

                        <!-- Dropdown Icon -->
                        <button type="button" @click="toggleDropdown"
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 9l-7 7-7-7" />
                            </svg>
                        </button>

                        <!-- Dropdown List -->
                        <Transition enter-active-class="transition ease-out duration-200"
                            enter-from-class="opacity-0 -translate-y-1" enter-to-class="opacity-100 translate-y-0"
                            leave-active-class="transition ease-in duration-150"
                            leave-from-class="opacity-100 translate-y-0" leave-to-class="opacity-0 -translate-y-1">
                            <div v-if="isDropdownOpen && filteredSuggestions.length > 0"
                                class="absolute left-0 right-0 top-full mt-1 max-h-48 overflow-y-auto bg-(--color-surface) border border-(--color-border) rounded-lg shadow-lg z-50 py-1">
                                <button v-for="(suggestion, index) in filteredSuggestions" :key="suggestion"
                                    type="button" @click="selectSuggestion(suggestion)"
                                    @mouseenter="highlightedIndex = index"
                                    class="w-full px-3 py-2 text-sm text-left transition-colors" :class="[
                                        highlightedIndex === index
                                            ? 'bg-(--color-blue)/10 text-(--color-blue)'
                                            : 'text-(--color-text-primary) hover:bg-(--color-muted-bg)'
                                    ]">
                                    <span class="flex items-center gap-2">
                                        <span>{{ suggestion }}</span>
                                        <span class="text-xs text-(--color-text-secondary) ml-auto">
                                            {{ getUsageCount(suggestion) }}×
                                        </span>
                                    </span>
                                </button>
                            </div>
                        </Transition>
                    </div>
                    <p class="text-xs text-(--color-text-secondary) mt-1">
                        {{ filteredSuggestions.length > 0 ? `${filteredSuggestions.length} suggestion(s) available` :
                        'Type to see suggestions' }}
                    </p>
                </div>

                <!-- Amount -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Amount <span class="text-(--color-red)">*</span>
                    </label>
                    <div class="relative">
                        <span
                            class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">à§³</span>
                        <input v-model.number="form.amount" type="number" step="0.01" min="0" placeholder="0.00"
                            required
                            class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                </div>

                <!-- Date -->
                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        Expense Date <span class="text-(--color-red)">*</span>
                    </label>
                    <input v-model="form.expense_date" type="date" required
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                </div>
            </div>
        </div>

        <!-- Notes -->
        <div>
            <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
                Additional Information
            </h3>
            <div>
                <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Notes</label>
                <textarea v-model="form.notes" rows="3" placeholder="Enter any notes about this expense"
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
                <span v-else>{{ isEditMode ? 'Save Changes' : 'Create Expense' }}</span>
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick, onUnmounted } from 'vue'
import type { Expense } from '@/types/expense'
import { useExpensesStore } from '@/stores/expenses'
import { push } from 'notivue'

const props = defineProps<{
    expense?: Expense | null
    mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
    'expense-created': []
    'expense-updated': []
    'cancel': []
}>()

const expensesStore = useExpensesStore()
const submitting = ref(false)

const isEditMode = computed(() => props.mode === 'edit' || !!props.expense)

// ============================================================
// AUTOCOMPLETE STATE
// ============================================================
const dropdownRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)
const searchTerm = ref('')
const isDropdownOpen = ref(false)
const highlightedIndex = ref(-1)

// Get unique expense titles from existing expenses
const uniqueTitles = computed(() => {
    const titles = expensesStore.expenses
        .map(e => e.title)
        .filter(title => title && title.trim())
        .filter((title, index, self) => self.indexOf(title) === index) // Unique
        .sort((a, b) => a.localeCompare(b))
    return titles
})

// Filter suggestions based on search term
const filteredSuggestions = computed(() => {
    if (!searchTerm.value.trim()) {
        return uniqueTitles.value.slice(0, 10)
    }
    const term = searchTerm.value.toLowerCase()
    return uniqueTitles.value
        .filter(title => title.toLowerCase().includes(term))
        .slice(0, 10)
})

// Get usage count for a title
const getUsageCount = (title: string): number => {
    if (!title) return 0
    return expensesStore.expenses.filter(e => e.title === title).length
}

// ============================================================
// AUTOCOMPLETE METHODS
// ============================================================
const openDropdown = () => {
    if (!isEditMode.value && filteredSuggestions.value.length > 0) {
        isDropdownOpen.value = true
        highlightedIndex.value = -1
    }
}

const closeDropdown = () => {
    isDropdownOpen.value = false
    highlightedIndex.value = -1
}

const toggleDropdown = () => {
    if (isDropdownOpen.value) {
        closeDropdown()
    } else {
        openDropdown()
    }
}

const selectSuggestion = (suggestion: string) => {
    if (!suggestion) return
    searchTerm.value = suggestion
    form.value.title = suggestion
    closeDropdown()
    // Focus the next field (amount)
    nextTick(() => {
        const amountInput = document.querySelector('input[type="number"]') as HTMLInputElement
        if (amountInput) amountInput.focus()
    })
}

const onInput = () => {
    form.value.title = searchTerm.value
    if (searchTerm.value.trim()) {
        isDropdownOpen.value = true
        highlightedIndex.value = -1
    } else {
        closeDropdown()
    }
}

const selectNext = () => {
    if (highlightedIndex.value < filteredSuggestions.value.length - 1) {
        highlightedIndex.value++
    }
}

const selectPrevious = () => {
    if (highlightedIndex.value > 0) {
        highlightedIndex.value--
    }
}

const handleEnter = () => {
    if (highlightedIndex.value >= 0 && highlightedIndex.value < filteredSuggestions.value.length) {
        const suggestion = filteredSuggestions.value[highlightedIndex.value]
        if (suggestion) {
            selectSuggestion(suggestion)
        }
    } else {
        closeDropdown()
    }
}

// ============================================================
// FORM STATE
// ============================================================
const form = ref({
    title: '',
    amount: null as number | null,
    expense_date: '',
    notes: '',
})

const initializeForm = () => {
    if (props.expense) {
        const date = new Date(props.expense.expense_date)
        form.value = {
            title: props.expense.title || '',
            amount: props.expense.amount || null,
            expense_date: date.toISOString().slice(0, 10),
            notes: props.expense.notes || '',
        }
        searchTerm.value = props.expense.title || ''
    } else {
        const now = new Date()
        form.value = {
            title: '',
            amount: null,
            expense_date: now.toISOString().slice(0, 10),
            notes: '',
        }
        searchTerm.value = ''
    }
}

watch(() => props.expense, initializeForm, { immediate: true })

const resetForm = () => {
    if (isEditMode.value && props.expense) {
        initializeForm()
    } else {
        const now = new Date()
        form.value = {
            title: '',
            amount: null,
            expense_date: now.toISOString().slice(0, 10),
            notes: '',
        }
        searchTerm.value = ''
    }
}

// ============================================================
// SUBMIT
// ============================================================
const submit = async () => {
    const title = form.value.title.trim()
    if (!title) {
        push.error('Title is required')
        return
    }

    if (!form.value.amount || form.value.amount <= 0) {
        push.error('Amount must be greater than 0')
        return
    }

    if (!form.value.expense_date) {
        push.error('Expense date is required')
        return
    }

    submitting.value = true

    try {
        if (isEditMode.value && props.expense) {
            const success = await expensesStore.updateExpense(props.expense.id, {
                title: title,
                amount: form.value.amount,
                expense_date: form.value.expense_date,
                notes: form.value.notes.trim() || null,
            })

            if (success) {
                push.success('Expense updated successfully!')
                emit('expense-updated')
            }
        } else {
            const newExpense = await expensesStore.createExpense({
                title: title,
                amount: form.value.amount,
                expense_date: form.value.expense_date,
                notes: form.value.notes.trim() || null,
            })

            if (newExpense) {
                push.success('Expense created successfully!')
                resetForm()
                emit('expense-created')
            }
        }
    } catch (error) {
        console.error('Error:', error)
        push.error(isEditMode.value ? 'Failed to update expense' : 'Failed to create expense')
    } finally {
        submitting.value = false
    }
}

// Close dropdown on click outside
const handleClickOutside = (event: MouseEvent) => {
    if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
        closeDropdown()
    }
}

onMounted(() => {
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>