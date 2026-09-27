<!-- src/components/layouts/AppSidebarNav.vue -->
<template>
    <nav class="flex-1 overflow-y-auto py-3 px-2.5 space-y-0.5">
        <AppSidebarNavItem to="/dashboard" icon="dashboard" label="Dashboard" :active="$route.path === '/dashboard'"
            @click="closeSidebar" />

        <AppSidebarNavItem to="/customers" icon="users" label="Customers" :active="$route.path === '/customers'"
            @click="closeSidebar" />

        <AppSidebarNavGroup v-for="group in menuGroups" :key="group.id" :group="group" @close="closeSidebar" />

        <AppSidebarNavItem to="/support" icon="support" label="Support" :active="$route.path === '/support'"
            @click="closeSidebar" />
    </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import AppSidebarNavItem from './AppSidebarNavItem.vue'
import AppSidebarNavGroup from './AppSidebarNavGroup.vue'

defineProps<{
    menuGroups: Array<{
        id: string
        label: string
        icon: string
        items: Array<{ path: string; label: string; icon: string }>
    }>
}>()

const emit = defineEmits<{
    (e: 'close'): void
}>()

const $route = useRoute()
const closeSidebar = () => emit('close')
</script>