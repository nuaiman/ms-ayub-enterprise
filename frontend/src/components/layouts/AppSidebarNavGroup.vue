<!-- src/components/layouts/AppSidebarNavGroup.vue -->
<template>
    <div class="space-y-0.5 relative">
        <!-- Group Header -->
        <button @click="handleToggle"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150"
            :class="[
                isActive(group.id)
                    ? 'text-(--color-blue)'
                    : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)',
                collapsed ? 'justify-center px-0' : ''
            ]">
            <span
                class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
                :class="[
                    collapsed && isActive(group.id)
                        ? 'ring-2 ring-(--color-blue)/25 bg-(--color-blue)/10 border-(--color-blue)/20'
                        : 'bg-(--color-muted-bg)/50 border-(--color-border)/30'
                ]">
                <component :is="groupIconComponent" class="w-3.5 h-3.5" :stroke-width="2" />
            </span>
            <span v-if="!collapsed" class="flex-1 text-left truncate text-[13px]">
                {{ group.label }}
            </span>
            <ChevronDown v-if="!collapsed"
                class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform duration-200"
                :class="{ 'rotate-180': isExpanded }" :stroke-width="2.5" />
        </button>

        <!-- Children - only show when expanded AND not collapsed -->
        <div v-if="isExpanded && !collapsed" class="ml-4 pl-3 space-y-0.5 border-l-2 border-(--color-border)/40">
            <AppSidebarNavItem v-for="item in group.items" :key="item.path" :to="item.path" :icon="item.icon"
                :label="item.label" :active="$route.path === item.path" :collapsed="false" @click="closeSidebar" />
        </div>

        <!-- Collapsed Flyout (desktop only) -->
        <div v-if="collapsed"
            class="hidden lg:block absolute left-full top-0 ml-3 w-48 p-1.5 rounded-xl bg-(--color-surface) border border-(--color-border) shadow-xl z-50 space-y-0.5 opacity-0 pointer-events-none group-hover:opacity-100 group-hover:pointer-events-auto transition-opacity">
            <div
                class="px-3 py-1.5 border-b border-(--color-border)/40 flex items-center gap-2 text-[11px] font-semibold text-(--color-text-primary)">
                <component :is="groupIconComponent" class="w-3 h-3" :stroke-width="2.5" />
                <span>{{ group.label }}</span>
            </div>
            <AppSidebarNavItem v-for="item in group.items" :key="item.path" :to="item.path" :icon="item.icon"
                :label="item.label" :active="$route.path === item.path" :collapsed="false" @click="closeSidebar" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Component } from 'vue'
import { useRoute } from 'vue-router'
import {
    ChevronDown,
    LayoutDashboard,
    Users,
    Handshake,
    HardHat,
    Warehouse,
    Package,
    ClipboardList,
    AlertCircle,
    Truck,
    PackageCheck,
    Car,
    FileText,
    Receipt,
    UserCog,
    Wallet,
    CreditCard,
    ScrollText,
    HelpCircle,
    Briefcase,
    Boxes,
    type LucideIcon,
} from 'lucide-vue-next'
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
            if (Array.isArray(parsed) && parsed.length > 0) return parsed
        }
    } catch { }
    return []
}

const expandedGroups = ref<string[]>(getDefaultExpanded())

const isExpanded = computed(() => expandedGroups.value.includes(props.group.id))

const saveExpandedState = () => {
    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(expandedGroups.value))
    } catch {
        // Ignore storage errors
    }
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
        warehouse: ['/lots', '/stores', '/damages', '/deliveries', '/delivery-items'],
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
    return groupPaths[groupId]?.some(path => $route.path.startsWith(path)) ?? false
}

const closeSidebar = () => {
    emit('close')
}

// Icon mapping registry for groups
const groupIconRegistry: Record<string, Component> = {
    parties: Handshake,
    warehouse: Boxes,
    transport: Truck,
    bills: Receipt,
    invoices: FileText,
    office: Briefcase,
}

const groupIconComponent = computed<Component>(() => {
    return groupIconRegistry[props.group.id] ?? LayoutDashboard
})

// Auto-expand if active
watch(
    () => $route.path,
    (newPath) => {
        const groupMap: Record<string, string[]> = {
            parties: ['/brokers', '/majhis', '/godowns'],
            warehouse: ['/lots', '/stores', '/damages', '/deliveries', '/delivery-items'],
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
    },
    { immediate: true }
)
</script>