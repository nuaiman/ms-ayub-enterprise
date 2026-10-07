<!-- src/components/layouts/AppSidebar.vue -->
<template>
    <aside ref="asideEl" @mouseenter="onMouseEnter" @mouseleave="onMouseLeave" :class="[
        'fixed inset-y-0 left-0 z-30 flex flex-col shrink-0',
        'bg-(--color-surface)/95 backdrop-blur-sm',
        'border-r border-(--color-border)/50',
        'lg:static lg:translate-x-0',
        isOpen ? 'translate-x-0' : '-translate-x-full',
        isExpanded ? 'w-60' : 'w-16',
        'transition-[width,transform] duration-200 ease-out',
    ]">
        <!-- Logo -->
        <div class="h-14 flex items-center border-b border-(--color-border)/40 shrink-0"
            :class="isExpanded ? 'px-4' : 'justify-center px-2'">
            <div class="flex items-center gap-3 overflow-hidden w-full" :class="isExpanded ? '' : 'justify-center'">
                <div
                    class="w-8 h-8 shrink-0 overflow-hidden rounded-lg bg-(--color-muted-bg) border border-(--color-border)/50 p-0.5 flex items-center justify-center">
                    <img src="@/assets/logo.png" alt="MS Ayub Enterprise"
                        class="w-full h-full block object-cover rounded-md" />
                </div>
                <div v-if="isExpanded" class="flex flex-col whitespace-nowrap">
                    <span class="text-sm font-semibold tracking-tight text-(--color-text-primary) leading-tight">
                        MS Ayub
                    </span>
                    <span class="text-[10px] font-medium text-(--color-text-secondary)/60 tracking-[0.15em] uppercase">
                        Enterprise
                    </span>
                </div>
            </div>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 overflow-y-auto py-3 space-y-0.5" :class="isExpanded ? 'px-2.5' : 'px-2'">
            <!-- Customers (top-level) -->
            <RouterLink to="/customers"
                class="group relative flex items-center rounded-lg text-sm font-medium w-full transition-colors duration-150"
                :class="[
                    route.path === '/customers'
                        ? 'bg-(--color-blue)/10 text-(--color-blue)'
                        : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)',
                    isExpanded ? 'gap-3 px-3 py-2' : 'justify-center p-1',
                ]" @click="onNavClick">
                <span
                    class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
                    :class="route.path === '/customers'
                        ? 'bg-(--color-blue)/15 border-(--color-blue)/20'
                        : 'bg-(--color-muted-bg)/50 border-(--color-border)/30 group-hover:border-(--color-border)/50'">
                    <Users class="w-3.5 h-3.5" :stroke-width="2" />
                </span>
                <span v-if="isExpanded" class="truncate text-[13px]">Customers</span>
                <div v-if="!isExpanded"
                    class="hidden lg:block absolute left-full top-1/2 -translate-y-1/2 ml-3 px-2.5 py-1 rounded-lg bg-(--color-surface) border border-(--color-border) text-[11px] font-medium text-(--color-text-primary) shadow-lg whitespace-nowrap z-50 pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity">
                    Customers
                </div>
            </RouterLink>

            <!-- Groups -->
            <div v-for="group in menuGroups" :key="group.id" class="relative w-full group">
                <button @click="toggleGroup(group.id)"
                    class="w-full flex items-center rounded-lg text-sm font-medium transition-colors duration-150"
                    :class="[
                        isGroupActive(group)
                            ? 'text-(--color-blue)'
                            : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)',
                        isExpanded ? 'gap-3 px-3 py-2' : 'justify-center p-1',
                    ]">
                    <span
                        class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
                        :class="isGroupActive(group)
                            ? 'bg-(--color-blue)/10 border-(--color-blue)/20'
                            : 'bg-(--color-muted-bg)/50 border-(--color-border)/30'">
                        <component :is="groupIcon(group.id)" class="w-3.5 h-3.5" :stroke-width="2" />
                    </span>
                    <template v-if="isExpanded">
                        <span class="flex-1 text-left truncate text-[13px]">{{ group.label }}</span>
                        <ChevronDown
                            class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform duration-200"
                            :class="{ 'rotate-180': isGroupExpanded(group.id) }" :stroke-width="2.5" />
                    </template>
                </button>

                <!-- Children (expanded) -->
                <div v-if="isExpanded && isGroupExpanded(group.id)"
                    class="ml-4 pl-3 mt-0.5 space-y-0.5 border-l-2 border-(--color-border)/40">
                    <RouterLink v-for="item in group.items" :key="item.path" :to="item.path"
                        class="group/item flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium w-full transition-colors duration-150"
                        :class="route.path === item.path
                            ? 'bg-(--color-blue)/10 text-(--color-blue)'
                            : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)'"
                        @click="onNavClick">
                        <span
                            class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
                            :class="route.path === item.path
                                ? 'bg-(--color-blue)/15 border-(--color-blue)/20'
                                : 'bg-(--color-muted-bg)/50 border-(--color-border)/30 group-hover/item:border-(--color-border)/50'">
                            <component :is="itemIcon(item.icon)" class="w-3.5 h-3.5" :stroke-width="2" />
                        </span>
                        <span class="truncate text-[13px]">{{ item.label }}</span>
                    </RouterLink>
                </div>

                <!-- Flyout (collapsed) -->
                <div v-if="!isExpanded"
                    class="hidden lg:block absolute left-full top-0 ml-3 w-48 p-1.5 rounded-xl bg-(--color-surface) border border-(--color-border) shadow-xl z-50 space-y-0.5 opacity-0 pointer-events-none group-hover:opacity-100 group-hover:pointer-events-auto transition-opacity duration-150">
                    <div
                        class="px-3 py-1.5 border-b border-(--color-border)/40 flex items-center gap-2 text-[11px] font-semibold text-(--color-text-primary)">
                        <component :is="groupIcon(group.id)" class="w-3 h-3" :stroke-width="2.5" />
                        <span>{{ group.label }}</span>
                    </div>
                    <RouterLink v-for="item in group.items" :key="item.path" :to="item.path"
                        class="flex items-center gap-2 px-3 py-2 rounded-lg text-[13px] font-medium transition-colors"
                        :class="route.path === item.path
                            ? 'bg-(--color-blue)/10 text-(--color-blue)'
                            : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)'"
                        @click="onNavClick">
                        <component :is="itemIcon(item.icon)" class="w-3.5 h-3.5 shrink-0" :stroke-width="2" />
                        <span class="truncate">{{ item.label }}</span>
                    </RouterLink>
                </div>
            </div>

            <!-- Support (top-level) -->
            <RouterLink to="/support"
                class="group relative flex items-center rounded-lg text-sm font-medium w-full transition-colors duration-150"
                :class="[
                    route.path === '/support'
                        ? 'bg-(--color-blue)/10 text-(--color-blue)'
                        : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)',
                    isExpanded ? 'gap-3 px-3 py-2' : 'justify-center p-1',
                ]" @click="onNavClick">
                <span
                    class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
                    :class="route.path === '/support'
                        ? 'bg-(--color-blue)/15 border-(--color-blue)/20'
                        : 'bg-(--color-muted-bg)/50 border-(--color-border)/30 group-hover:border-(--color-border)/50'">
                    <HelpCircle class="w-3.5 h-3.5" :stroke-width="2" />
                </span>
                <span v-if="isExpanded" class="truncate text-[13px]">Support</span>
                <div v-if="!isExpanded"
                    class="hidden lg:block absolute left-full top-1/2 -translate-y-1/2 ml-3 px-2.5 py-1 rounded-lg bg-(--color-surface) border border-(--color-border) text-[11px] font-medium text-(--color-text-primary) shadow-lg whitespace-nowrap z-50 pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity">
                    Support
                </div>
            </RouterLink>
        </nav>

        <!-- User Card -->
        <div class="p-2.5 border-t border-(--color-border)/40 shrink-0">
            <UserDropdown variant="sidebar" :collapsed="!isExpanded" @logout="emit('logout')"
                @change-password="emit('changePassword')" @reset-all-passwords="emit('resetAllPasswords')"
                @download-backup="emit('downloadBackup')" />
        </div>
    </aside>
</template>

<script setup lang="ts">
import { reactive, watch, onMounted, onUnmounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import {
    Users,
    HelpCircle,
    ChevronDown,
    Handshake,
    Boxes,
    Briefcase,
    LayoutDashboard,
    Receipt,
    ArrowRightLeft,
    TrendingUp,
    Truck,
    TriangleAlert
} from 'lucide-vue-next'
import type { Component } from 'vue'
import UserDropdown from './UserDropdown.vue'
import type { MenuGroup } from './types'

const props = defineProps<{
    isOpen: boolean
    menuGroups: MenuGroup[]
}>()

const emit = defineEmits<{
    (e: 'close'): void
    (e: 'logout'): void
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
}>()

const route = useRoute()

// =============================================================================
// HOVER-TO-EXPAND (self-contained)
// =============================================================================

const asideEl = ref<HTMLElement | null>(null)
const isMobile = ref(window.innerWidth < 1024)
const isHovered = ref(false)

const isExpanded = computed(() => {
    if (isMobile.value) return true
    return isHovered.value
})

let lastX = -1
let lastY = -1

const onDocMouseMove = (e: MouseEvent) => {
    lastX = e.clientX
    lastY = e.clientY
}

const pointerIsInside = (): boolean => {
    const el = asideEl.value
    if (!el) return false
    if (lastX < 0 || lastY < 0) return el.matches(':hover')
    const r = el.getBoundingClientRect()
    return lastX >= r.left && lastX <= r.right && lastY >= r.top && lastY <= r.bottom
}

let leaveTimer: ReturnType<typeof setTimeout> | null = null

const onMouseEnter = () => {
    if (leaveTimer) {
        clearTimeout(leaveTimer)
        leaveTimer = null
    }
    if (!isMobile.value) {
        isHovered.value = true
    }
}

const onMouseLeave = () => {
    if (isMobile.value) return
    if (leaveTimer) clearTimeout(leaveTimer)
    leaveTimer = setTimeout(() => {
        leaveTimer = null
        if (pointerIsInside()) return
        if (asideEl.value && asideEl.value.matches(':hover')) return
        isHovered.value = false
    }, 120)
}

// =============================================================================
// NAV CLICK / MOBILE
// =============================================================================

const onNavClick = () => {
    if (isMobile.value) {
        emit('close')
    }
}

const updateMobile = () => {
    const wasMobile = isMobile.value
    isMobile.value = window.innerWidth < 1024

    if (!wasMobile && isMobile.value) {
        isHovered.value = false
        if (leaveTimer) {
            clearTimeout(leaveTimer)
            leaveTimer = null
        }
    }
}

onMounted(() => {
    document.addEventListener('mousemove', onDocMouseMove, { passive: true })
    window.addEventListener('resize', updateMobile)
})

onUnmounted(() => {
    document.removeEventListener('mousemove', onDocMouseMove)
    window.removeEventListener('resize', updateMobile)
    if (leaveTimer) {
        clearTimeout(leaveTimer)
        leaveTimer = null
    }
})

// =============================================================================
// GROUP EXPANSION (persisted)
// =============================================================================

const STORAGE_KEY = 'sidebar_expanded_groups'

const loadExpanded = (): Record<string, boolean> => {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        if (!raw) return {}
        const parsed = JSON.parse(raw)
        if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
            return parsed as Record<string, boolean>
        }
    } catch {
        // ignore
    }
    return {}
}

const expandedGroups = reactive<Record<string, boolean>>(loadExpanded())

const saveExpanded = () => {
    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(expandedGroups))
    } catch {
        // ignore
    }
}

const isGroupExpanded = (id: string) => expandedGroups[id] === true

const toggleGroup = (id: string) => {
    if (!isExpanded.value) return
    expandedGroups[id] = !expandedGroups[id]
    saveExpanded()
}

const isGroupActive = (group: MenuGroup) =>
    group.items.some(item => route.path.startsWith(item.path))

watch(
    () => route.path,
    () => {
        for (const group of props.menuGroups) {
            if (isGroupActive(group) && expandedGroups[group.id] !== true) {
                expandedGroups[group.id] = true
                saveExpanded()
            }
        }
    },
    { immediate: true }
)

// =============================================================================
// ICONS
// =============================================================================

const groupIcons: Record<string, Component> = {
    parties: Handshake,
    warehouse: Boxes,
    office: Briefcase,
    bills: Receipt,
}

const itemIcons: Record<string, Component> = {
    users: Users,
    brokers: Handshake,
    majhis: Users,
    godowns: Boxes,
    lots: Boxes,
    stores: Boxes,
    salaries: Briefcase,
    expenses: Briefcase,
    incomes: TrendingUp,
    logs: Briefcase,
    'majhi-bills': Receipt,
    'godown-bills': Receipt,
    'customer-store-bills': Receipt,
    'customer-delivery-bills': Receipt,
    'customer-additional-bills': Receipt,
    'invoices': Receipt,
    'store-adjustments': ArrowRightLeft,
    'store-transfers': ArrowRightLeft,
    'lot-transfers': ArrowRightLeft,
    'deliveries': Truck,
    'damages': TriangleAlert,
}

const groupIcon = (id: string): Component => groupIcons[id] ?? LayoutDashboard
const itemIcon = (icon: string): Component => itemIcons[icon] ?? LayoutDashboard
</script>