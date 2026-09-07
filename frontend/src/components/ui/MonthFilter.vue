<!-- src/components/ui/MonthFilter.vue -->
<template>
    <div class="relative" ref="dropdownRef">
        <!-- Trigger Button -->
        <button @click="toggleDropdown"
            class="px-3 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent flex items-center gap-2 min-w-35 justify-between">
            <span>{{ selectedLabel }}</span>
            <svg class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform"
                :class="{ 'rotate-180': isOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7" />
            </svg>
        </button>

        <!-- Dropdown Menu -->
        <div v-if="isOpen"
            class="absolute left-0 top-full mt-1 w-56 max-h-80 overflow-y-auto bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg z-50 py-1">
            <!-- All Months Option -->
            <div @click="selectMonth('')"
                class="px-3 py-2 text-sm hover:bg-(--color-muted-bg) cursor-pointer transition-colors"
                :class="modelValue === '' ? 'bg-(--color-blue)/10 text-(--color-blue)' : 'text-(--color-text-primary)'">
                All Months
            </div>

            <!-- Divider -->
            <div class="border-t border-(--color-border)/40 my-1"></div>

            <!-- Year Groups -->
            <div v-for="group in groupedMonths" :key="group.year">
                <!-- Year Header (click to toggle) -->
                <div @click="toggleYear(group.year)"
                    class="flex items-center justify-between px-3 py-1.5 text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider hover:bg-(--color-muted-bg)/50 cursor-pointer transition-colors">
                    <span>{{ group.year }}</span>
                    <svg class="w-3 h-3 transition-transform" :class="{ 'rotate-180': isYearExpanded(group.year) }"
                        fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7" />
                    </svg>
                </div>

                <!-- Months under this year -->
                <div v-if="isYearExpanded(group.year)" class="ml-2">
                    <div v-for="month in group.months" :key="month.value" @click="selectMonth(month.value)"
                        class="px-3 py-1.5 text-sm hover:bg-(--color-muted-bg) cursor-pointer transition-colors rounded-l-lg"
                        :class="modelValue === month.value ? 'bg-(--color-blue)/10 text-(--color-blue)' : 'text-(--color-text-primary)'">
                        {{ month.label }}
                    </div>
                </div>
            </div>

            <!-- Empty State -->
            <div v-if="groupedMonths.length === 0" class="px-3 py-4 text-sm text-(--color-text-secondary) text-center">
                No data available
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useClickOutside } from '@/composables/useClickOutside'

const props = defineProps<{
    modelValue: string
    items: string[]  // Array of month strings in format "YYYY-MM"
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
}>()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

// Track expanded years
const expandedYears = ref<Set<string>>(new Set())

// Group months by year with deduplication
const groupedMonths = computed(() => {
    const yearMap = new Map<string, Set<string>>()

    // Use Set to deduplicate months within each year
    props.items.forEach(month => {
        const [year, monthNum] = month.split('-')
        if (!year || !monthNum) return

        if (!yearMap.has(year)) {
            yearMap.set(year, new Set())
        }
        yearMap.get(year)!.add(month)
    })

    // Convert to array and sort
    const result = Array.from(yearMap.entries())
        .map(([year, monthSet]) => ({
            year,
            months: Array.from(monthSet)
                .sort((a, b) => b.localeCompare(a)) // Descending (newest first)
                .map(m => ({
                    value: m,
                    label: formatMonthLabel(m)
                }))
        }))
        .sort((a, b) => b.year.localeCompare(a.year)) // Descending years

    // Auto-expand the year that contains the selected month
    if (props.modelValue) {
        const [selectedYear] = props.modelValue.split('-')
        if (selectedYear) {
            expandedYears.value.add(selectedYear)
        }
    }

    // If no selected month, expand the first year with data
    if (!props.modelValue && result.length > 0 && expandedYears.value.size === 0) {
        expandedYears.value.add(result[0]!.year)
    }

    return result
})

// Selected label for the trigger button
const selectedLabel = computed(() => {
    if (!props.modelValue) return 'All Months'
    return formatMonthLabel(props.modelValue)
})

// Format "YYYY-MM" to "Month Year"
const formatMonthLabel = (monthYear: string): string => {
    const [year, month] = monthYear.split('-')
    if (!year || !month) return monthYear
    const date = new Date(parseInt(year), parseInt(month) - 1)
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
}

// Toggle dropdown
const toggleDropdown = () => {
    isOpen.value = !isOpen.value
}

// Select a month
const selectMonth = (value: string) => {
    emit('update:modelValue', value)
    isOpen.value = false
}

// Toggle year expansion
const toggleYear = (year: string) => {
    if (expandedYears.value.has(year)) {
        expandedYears.value.delete(year)
    } else {
        expandedYears.value.add(year)
    }
}

// Check if year is expanded
const isYearExpanded = (year: string): boolean => {
    return expandedYears.value.has(year)
}

// Click outside to close
useClickOutside(dropdownRef, () => {
    isOpen.value = false
})

// Close on escape key
const handleEscape = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && isOpen.value) {
        isOpen.value = false
    }
}

onMounted(() => {
    document.addEventListener('keydown', handleEscape)
})

onUnmounted(() => {
    document.removeEventListener('keydown', handleEscape)
})
</script>