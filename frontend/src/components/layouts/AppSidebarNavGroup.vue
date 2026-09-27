<!-- src/components/layouts/AppSidebarNavGroup.vue -->
<template>
    <div class="space-y-0.5">
        <!-- Group Header -->
        <button @click="toggle"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150"
            :class="isActive
                ? 'text-(--color-blue)'
                : 'text-(--color-text-secondary) hover:bg-(--color-muted-bg)/60 hover:text-(--color-text-primary)'">
            <span
                class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border transition-colors duration-150"
                :class="isActive
                    ? 'bg-(--color-blue)/10 border-(--color-blue)/20'
                    : 'bg-(--color-muted-bg)/50 border-(--color-border)/30'">
                <component :is="groupIcon" class="w-3.5 h-3.5" :stroke-width="2" />
            </span>
            <span class="flex-1 text-left text-[13px] leading-tight">{{ group.label }}</span>
            <ChevronDown class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform duration-200"
                :class="{ 'rotate-180': isExpanded }" :stroke-width="2.5" />
        </button>

        <!-- Children -->
        <div v-if="isExpanded" class="ml-4 pl-3 space-y-0.5 border-l-2 border-(--color-border)/40">
            <AppSidebarNavItem v-for="item in group.items" :key="item.path" :to="item.path" :icon="item.icon"
                :label="item.label" :active="$route.path === item.path" @click="closeSidebar" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Component } from 'vue'
import { useRoute } from 'vue-router'
import {
    ChevronDown,
    Handshake,
    Boxes,
    Truck,
    Receipt,
    FileText,
    Briefcase,
    LayoutDashboard,
} from 'lucide-vue-next'
import AppSidebarNavItem from './AppSidebarNavItem.vue'

const props = defineProps<{
    group: {
        id: string
        label: string
        icon: string
        items: Array<{ path: string; label: string; icon: string }>
    }
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

const $route = useRoute()

// In-memory only. Starts collapsed. User toggles manually. No persistence.
const isExpanded = ref(false)

const isActive = computed(() => {
    return props.group.items.some(item => $route.path.startsWith(item.path))
})

const groupIconRegistry: Record<string, Component> = {
    parties: Handshake,
    warehouse: Boxes,
    transport: Truck,
    bills: Receipt,
    invoices: FileText,
    office: Briefcase,
}

const groupIcon = computed<Component>(() => {
    return groupIconRegistry[props.group.id] ?? LayoutDashboard
})

const toggle = () => {
    isExpanded.value = !isExpanded.value
}

const closeSidebar = () => emit('close')
</script>