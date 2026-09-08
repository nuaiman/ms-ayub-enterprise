<!-- src/components/layouts/AppSidebar.vue -->
<template>
    <aside @mouseenter="$emit('mouseenter')" @mouseleave="$emit('mouseleave')" :class="[
        'fixed inset-y-0 left-0 z-30 flex flex-col justify-between shrink-0',
        'bg-(--color-surface)/95 backdrop-blur-sm',
        'border-r border-(--color-border)/50',
        // Mobile: slide in/out
        isOpen ? 'translate-x-0' : '-translate-x-full',
        // Desktop: always visible
        'lg:translate-x-0 lg:static',
        // Fixed width - no collapse
        'w-60 lg:w-60'
    ]">
        <!-- Logo - always full -->
        <AppSidebarLogo :collapsed="false" />

        <!-- Navigation - always full -->
        <AppSidebarNav :collapsed="false" :menu-groups="menuGroups" @close="closeSidebar" />

        <!-- User Card - always visible -->
        <div>
            <AppSidebarUser :collapsed="false" @logout="handleLogout" @change-password="handleChangePassword"
                @reset-all-passwords="handleResetAllPasswords" @download-backup="handleDownloadBackup" />
        </div>
    </aside>
</template>

<script setup lang="ts">
import AppSidebarLogo from './AppSidebarLogo.vue'
import AppSidebarNav from './AppSidebarNav.vue'
import AppSidebarUser from './AppSidebarUser.vue'

defineProps<{
    isOpen: boolean
    isCollapsed: boolean // Keep prop for compatibility but don't use
    isMobile: boolean
    menuGroups: any[]
}>()

const emit = defineEmits<{
    (e: 'toggle'): void
    (e: 'close'): void
    (e: 'mouseenter'): void
    (e: 'mouseleave'): void
    (e: 'logout'): void
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
}>()

const closeSidebar = () => emit('close')
const handleLogout = () => emit('logout')
const handleChangePassword = () => emit('changePassword')
const handleResetAllPasswords = () => emit('resetAllPasswords')
const handleDownloadBackup = () => emit('downloadBackup')
</script>