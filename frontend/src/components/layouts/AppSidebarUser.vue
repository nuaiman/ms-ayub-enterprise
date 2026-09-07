<!-- src/components/layouts/AppSidebarUser.vue -->
<template>
    <div ref="dropdownRef" class="p-2.5 border-t border-(--color-border)/40 shrink-0 relative">
        <!-- User Card - Click to toggle dropdown -->
        <button @click="toggleDropdown" class="w-full flex items-center rounded-xl transition-colors" :class="[
            collapsed ? 'lg:justify-center lg:p-1' : 'lg:p-2 lg:gap-2.5 lg:bg-(--color-muted-bg)/30 lg:hover:bg-(--color-muted-bg)/50'
        ]">
            <div
                class="w-8 h-8 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-[11px] font-bold shrink-0">
                {{ userInitials }}
            </div>
            <!-- ✅ Show text when NOT collapsed -->
            <div v-if="!collapsed" class="flex-1 min-w-0 text-left">
                <p class="text-[12.5px] font-semibold text-(--color-text-primary) truncate leading-tight">{{
                    auth.user?.name }}</p>
                <p class="text-[10px] text-(--color-text-secondary)/60 truncate capitalize mt-0.5">{{ auth.user?.role }}
                </p>
            </div>
            <!-- ✅ Show chevron when NOT collapsed -->
            <svg v-if="!collapsed" class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform"
                :class="{ 'rotate-180': dropdownOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7" />
            </svg>
        </button>

        <!-- Dropdown Menu -->
        <div v-if="dropdownOpen"
            class="absolute bottom-full left-0 right-0 mb-2 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50"
            :class="collapsed ? 'lg:w-56 lg:left-auto lg:right-0' : ''">
            <!-- User Info Header -->
            <div class="px-4 py-3 border-b border-(--color-border)/40">
                <div class="flex items-center gap-3">
                    <div
                        class="w-9 h-9 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-sm font-bold shrink-0">
                        {{ userInitials }}
                    </div>
                    <div class="min-w-0">
                        <div class="text-sm font-semibold text-(--color-text-primary) truncate">{{ auth.user?.name }}
                        </div>
                        <div class="text-xs text-(--color-text-secondary) truncate">@{{ auth.user?.username }}</div>
                        <span
                            class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-medium capitalize mt-0.5"
                            :class="getRoleBadgeClass(auth.user?.role || '')">
                            {{ auth.user?.role }}
                        </span>
                    </div>
                </div>
            </div>

            <!-- Menu Items - Shared Component -->
            <UserMenuItems @change-password="handleChangePassword" @reset-all-passwords="handleResetAllPasswords"
                @download-backup="handleDownloadBackup" @logout="handleLogout" />
        </div>

        <!-- Collapsed User Popover Card (desktop only) -->
        <div v-if="collapsed && !dropdownOpen"
            class="hidden lg:block absolute left-full bottom-0 ml-3 w-44 p-2.5 rounded-xl bg-(--color-surface) border border-(--color-border) shadow-xl opacity-0 pointer-events-none group-hover/user:opacity-100 group-hover/user:pointer-events-auto z-50 space-y-2">
            <div>
                <p class="text-[12px] font-semibold text-(--color-text-primary) truncate">{{ auth.user?.name }}</p>
                <p class="text-[10px] text-(--color-text-secondary)/60 capitalize">{{ auth.user?.role }}</p>
            </div>
            <button @click="handleLogout"
                class="w-full flex items-center gap-2 px-2.5 py-1.5 rounded-lg text-[11px] font-medium text-(--color-red) bg-(--color-red)/10 hover:bg-(--color-red)/20 transition-colors">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                Logout
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useClickOutside } from '@/composables/useClickOutside'
import UserMenuItems from './UserMenuItems.vue'

const props = defineProps<{
    collapsed: boolean
}>()

const emit = defineEmits<{
    (e: 'logout'): void
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
}>()

const auth = useAuthStore()
const themeStore = useThemeStore()

const dropdownRef = ref<HTMLElement | null>(null)
const dropdownOpen = ref(false)

useClickOutside(dropdownRef, () => {
    dropdownOpen.value = false
})

const handleResize = () => {
    dropdownOpen.value = false
}

onMounted(() => {
    window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
})

const userInitials = computed(() => {
    if (!auth.user?.name) return '?'
    return auth.user.name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

const getRoleBadgeClass = (role: string): string => {
    switch (role) {
        case 'admin': return 'bg-(--color-blue)/10 text-(--color-blue)'
        case 'manager': return 'bg-(--color-yellow)/10 text-(--color-yellow)'
        default: return 'bg-(--color-muted-bg) text-(--color-text-secondary)'
    }
}

const toggleDropdown = () => {
    dropdownOpen.value = !dropdownOpen.value
}

const handleChangePassword = () => {
    dropdownOpen.value = false
    emit('changePassword')
}

const handleResetAllPasswords = () => {
    dropdownOpen.value = false
    emit('resetAllPasswords')
}

const handleDownloadBackup = () => {
    dropdownOpen.value = false
    emit('downloadBackup')
}

const handleLogout = () => {
    dropdownOpen.value = false
    emit('logout')
}
</script>