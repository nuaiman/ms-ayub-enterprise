<!-- src/components/layouts/AppSidebar.vue -->
<template>
    <aside
        class="fixed inset-y-0 left-0 z-30 flex flex-col bg-(--color-surface)/95 backdrop-blur-sm border-r border-(--color-border)/50 w-72 shrink-0 transition-transform duration-200 ease-out"
        :class="isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0 lg:static'">
        <AppSidebarLogo />
        <AppSidebarNav :menu-groups="menuGroups" @close="closeSidebar" />
        <AppSidebarUser @logout="handleLogout" @change-password="handleChangePassword"
            @reset-all-passwords="handleResetAllPasswords" @download-backup="handleDownloadBackup" />
    </aside>
</template>

<script setup lang="ts">
import AppSidebarLogo from './AppSidebarLogo.vue'
import AppSidebarNav from './AppSidebarNav.vue'
import AppSidebarUser from './AppSidebarUser.vue'

defineProps<{
    isOpen: boolean
    menuGroups: Array<{
        id: string
        label: string
        icon: string
        items: Array<{ path: string; label: string; icon: string }>
    }>
}>()

const emit = defineEmits<{
    (e: 'close'): void
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