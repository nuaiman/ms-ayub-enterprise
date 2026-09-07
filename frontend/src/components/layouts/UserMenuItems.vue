<!-- src/components/layouts/UserMenuItems.vue -->
<template>
    <div class="py-1">
        <!-- Theme Toggle -->
        <button @click="handleThemeToggle"
            class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <svg v-if="!themeStore.isDark" class="w-4 h-4 shrink-0" fill="none" stroke="currentColor"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <svg v-else class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
            {{ themeStore.isDark ? 'Light Mode' : 'Dark Mode' }}
        </button>

        <div class="border-t border-(--color-border)/30 my-1"></div>

        <!-- Change Password -->
        <button @click="handleChangePassword"
            class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
            Change Password
        </button>

        <!-- Admin Only -->
        <template v-if="auth.user?.role === 'admin'">
            <div class="border-t border-(--color-border)/30 my-1"></div>

            <!-- Download Backup -->
            <button @click="handleDownloadBackup"
                class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Download Backup
            </button>

            <!-- Reset All Passwords -->
            <button @click="handleResetAllPasswords"
                class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                Reset All Passwords
            </button>
        </template>

        <div class="border-t border-(--color-border)/30 my-1"></div>

        <!-- Logout -->
        <button @click="handleLogout"
            class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-red) hover:bg-(--color-red)/10 transition-colors">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Logout
        </button>
    </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

const emit = defineEmits<{
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
    (e: 'logout'): void
}>()

const auth = useAuthStore()
const themeStore = useThemeStore()

const handleThemeToggle = () => {
    themeStore.toggleTheme()
}

const handleChangePassword = () => {
    emit('changePassword')
}

const handleResetAllPasswords = () => {
    emit('resetAllPasswords')
}

const handleDownloadBackup = () => {
    emit('downloadBackup')
}

const handleLogout = () => {
    emit('logout')
}
</script>