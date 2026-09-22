<!-- src/components/layouts/UserMenuItems.vue -->
<template>
    <div class="py-1">
        <!-- Theme Toggle -->
        <button @click="handleThemeToggle"
            class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <Sun v-if="!themeStore.isDark" class="w-4 h-4 shrink-0" :stroke-width="2" />
            <Moon v-else class="w-4 h-4 shrink-0" :stroke-width="2" />
            {{ themeStore.isDark ? 'Light Mode' : 'Dark Mode' }}
        </button>

        <div class="border-t border-(--color-border)/30 my-1"></div>

        <!-- Change Password -->
        <button @click="handleChangePassword"
            class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <KeyRound class="w-4 h-4 shrink-0" :stroke-width="2" />
            Change Password
        </button>

        <!-- Admin Only -->
        <template v-if="auth.user?.role === 'admin'">
            <div class="border-t border-(--color-border)/30 my-1"></div>

            <!-- Download Backup -->
            <button @click="handleDownloadBackup"
                class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                <Download class="w-4 h-4 shrink-0" :stroke-width="2" />
                Download Backup
            </button>

            <!-- Reset All Passwords -->
            <button @click="handleResetAllPasswords"
                class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                <Lock class="w-4 h-4 shrink-0" :stroke-width="2" />
                Reset All Passwords
            </button>
        </template>

        <div class="border-t border-(--color-border)/30 my-1"></div>

        <!-- Logout -->
        <button @click="handleLogout"
            class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-red) hover:bg-(--color-red)/10 transition-colors">
            <LogOut class="w-4 h-4 shrink-0" :stroke-width="2" />
            Logout
        </button>
    </div>
</template>

<script setup lang="ts">
import { Sun, Moon, KeyRound, Download, Lock, LogOut } from 'lucide-vue-next'
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