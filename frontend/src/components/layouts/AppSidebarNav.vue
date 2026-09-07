<!-- src/components/layouts/AppSidebarNav.vue -->
<template>
    <nav class="flex-1 overflow-y-auto py-3 px-2.5 space-y-0.5">
        <!-- Dashboard -->
        <AppSidebarNavItem to="/dashboard" icon="📊" label="Dashboard" :active="$route.path === '/dashboard'"
            :collapsed="collapsed" @click="closeSidebar" />

        <!-- Customers -->
        <AppSidebarNavItem to="/customers" icon="👥" label="Customers" :active="$route.path === '/customers'"
            :collapsed="collapsed" @click="closeSidebar" />

        <!-- Navigation Groups -->
        <AppSidebarNavGroup v-for="group in menuGroups" :key="group.id" :group="group" :collapsed="collapsed"
            @close="closeSidebar" />

        <!-- Support -->
        <AppSidebarNavItem to="/support" icon="❓" label="Support" :active="$route.path === '/support'"
            :collapsed="collapsed" @click="closeSidebar" />
    </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import AppSidebarNavItem from './AppSidebarNavItem.vue'
import AppSidebarNavGroup from './AppSidebarNavGroup.vue'

defineProps<{
    collapsed: boolean
    menuGroups: any[]
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

const $route = useRoute()

const closeSidebar = () => {
    emit('close')
}
</script>