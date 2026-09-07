<!-- src/components/layouts/AppSidebarNavGroup.vue -->
<template>
    <div class="space-y-0.5 relative">
        <!-- Group Header -->
        <button @click="handleToggle"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)"
            :class="[
                isActive(group.id) ? 'text-(--color-blue)' : '',
                collapsed ? 'justify-center px-0' : ''
            ]">
            <span
                class="w-7 h-7 rounded-lg bg-(--color-muted-bg)/50 flex items-center justify-center text-sm shrink-0 border border-(--color-border)/30"
                :class="{ 'ring-2 ring-(--color-blue)/25 bg-(--color-blue)/10': collapsed && isActive(group.id) }">
                {{ group.icon }}
            </span>
            <!-- Show text when NOT collapsed -->
            <span v-if="!collapsed" class="flex-1 text-left truncate text-[13px]">{{ group.label }}</span>
            <!-- Show chevron when NOT collapsed -->
            <svg v-if="!collapsed" class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0"
                :class="{ 'rotate-180': isExpanded }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7" />
            </svg>
        </button>

        <!-- Children - only show when expanded AND not collapsed -->
        <div v-if="isExpanded && !collapsed" class="ml-4 pl-3 space-y-0.5 border-l-2 border-(--color-border)/40">
            <AppSidebarNavItem v-for="item in group.items" :key="item.path" :to="item.path" :icon="item.icon"
                :label="item.label" :active="$route.path === item.path" :collapsed="false" @click="closeSidebar" />
        </div>

        <!-- Collapsed Flyout (desktop only) -->
        <div v-if="collapsed"
            class="hidden lg:block absolute left-full top-0 ml-3 w-48 p-1.5 rounded-xl bg-(--color-surface) border border-(--color-border) shadow-xl z-50 space-y-0.5">
            <div
                class="px-3 py-1.5 border-b border-(--color-border)/40 flex items-center gap-2 text-[11px] font-semibold text-(--color-text-primary)">
                <span>{{ group.icon }}</span>
                <span>{{ group.label }}</span>
            </div>
            <AppSidebarNavItem v-for="item in group.items" :key="item.path" :to="item.path" :icon="item.icon"
                :label="item.label" :active="$route.path === item.path" :collapsed="false" @click="closeSidebar" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebarNavItem from './AppSidebarNavItem.vue'

const props = defineProps<{
    group: {
        id: string
        label: string
        icon: string
        items: Array<{ path: string; label: string; icon: string }>
    }
    collapsed: boolean
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

const $route = useRoute()

const STORAGE_KEY = 'sidebar_expanded_groups'

const getDefaultExpanded = (): string[] => {
    try {
        const saved = localStorage.getItem(STORAGE_KEY)
        if (saved) {
            const parsed = JSON.parse(saved)
            if (Array.isArray(parsed)) return parsed
        }
    } catch { }
    // Auto-expand Bills and Invoices groups by default
    return ['bills', 'invoices']
}

const expandedGroups = ref<string[]>(getDefaultExpanded())

const isExpanded = computed(() => expandedGroups.value.includes(props.group.id))

const saveExpandedState = () => {
    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(expandedGroups.value))
    } catch { }
}

const handleToggle = () => {
    if (props.collapsed) return
    const index = expandedGroups.value.indexOf(props.group.id)
    if (index === -1) {
        expandedGroups.value.push(props.group.id)
    } else {
        expandedGroups.value.splice(index, 1)
    }
    saveExpandedState()
}

const isActive = (groupId: string) => {
    const groupPaths: Record<string, string[]> = {
        parties: ['/brokers', '/majhis', '/godowns', '/rents'],
        warehouse: ['/items', '/lots', '/stores', '/damages', '/deliveries', '/delivery-items'],
        transport: ['/transports', '/vehicles'],
        bills: [
            '/godown-store-bills',
            '/customer-storage-bills',
            '/customer-lot-bills',
            '/majhi-lot-bills',
            '/customer-delivery-bills',
            '/majhi-loading-bills',
            '/broker-vehicle-bills',
            '/customer-transport-bills',
        ],
        invoices: ['/invoices'],
        office: ['/users', '/salaries', '/expenses', '/logs'],
    }
    return groupPaths[groupId]?.some(path => $route.path.startsWith(path)) || false
}

const closeSidebar = () => {
    emit('close')
}

// Auto-expand if active
watch(() => $route.path, (newPath) => {
    const groupMap: Record<string, string[]> = {
        parties: ['/brokers', '/majhis', '/godowns'],
        warehouse: ['/items', '/lots', '/stores', '/damages', '/deliveries', '/delivery-items'],
        transport: ['/transports', '/vehicles'],
        bills: [
            '/godown-store-bills',
            '/customer-storage-bills',
            '/customer-lot-bills',
            '/majhi-lot-bills',
            '/customer-delivery-bills',
            '/majhi-loading-bills',
            '/broker-vehicle-bills',
            '/customer-transport-bills',
        ],
        invoices: ['/invoices'],
        office: ['/users', '/salaries', '/expenses', '/logs'],
    }

    for (const [id, paths] of Object.entries(groupMap)) {
        if (paths.some(path => newPath.startsWith(path))) {
            if (!expandedGroups.value.includes(id)) {
                expandedGroups.value.push(id)
                saveExpandedState()
            }
        }
    }
}, { immediate: true })
</script>