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
        // Width based on state
        isMobile ? 'w-70' : (isCollapsed ? 'lg:w-17' : 'lg:w-60'),
        'w-70'  // Mobile: full width
    ]">
        <!-- Logo -->
        <!-- ✅ On mobile, always show full logo (collapsed=false) -->
        <AppSidebarLogo :collapsed="isMobile ? false : isCollapsed" />

        <!-- Navigation -->
        <!-- ✅ On mobile, always show full nav (collapsed=false) -->
        <AppSidebarNav :collapsed="isMobile ? false : isCollapsed" :menu-groups="menuGroups" @close="closeSidebar" />

        <!-- User Card - Hidden on mobile (shown in header) -->
        <div class="hidden lg:block">
            <!-- ✅ On desktop only, use isCollapsed -->
            <AppSidebarUser :collapsed="isCollapsed" @logout="handleLogout" @change-password="handleChangePassword"
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
    isCollapsed: boolean
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