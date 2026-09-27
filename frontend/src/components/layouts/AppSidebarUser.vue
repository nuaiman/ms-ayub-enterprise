<!-- src/components/layouts/AppSidebarUser.vue -->
<template>
    <div ref="dropdownRef" class="p-2.5 border-t border-(--color-border)/40 shrink-0 relative">
        <!-- User Card - Click to toggle dropdown -->
        <button @click="toggleDropdown"
            class="w-full flex items-center p-2 gap-2.5 rounded-xl bg-(--color-muted-bg)/30 hover:bg-(--color-muted-bg)/50 transition-colors duration-150">
            <div
                class="w-8 h-8 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-[11px] font-bold shrink-0">
                {{ userInitials }}
            </div>
            <div class="flex-1 min-w-0 text-left">
                <p class="text-[12.5px] font-semibold text-(--color-text-primary) truncate leading-tight">
                    {{ auth.user?.name }}
                </p>
                <p class="text-[10px] text-(--color-text-secondary)/60 truncate capitalize mt-0.5">
                    {{ auth.user?.role }}
                </p>
            </div>
            <ChevronDown class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform duration-200"
                :class="{ 'rotate-180': dropdownOpen }" :stroke-width="2.5" />
        </button>

        <!-- Dropdown Menu -->
        <div v-if="dropdownOpen"
            class="absolute bottom-full left-0 right-0 mb-2 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50">
            <div class="px-4 py-3 border-b border-(--color-border)/40">
                <div class="flex items-center gap-3">
                    <div
                        class="w-9 h-9 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-sm font-bold shrink-0">
                        {{ userInitials }}
                    </div>
                    <div class="min-w-0">
                        <div class="text-sm font-semibold text-(--color-text-primary) truncate">
                            {{ auth.user?.name }}
                        </div>
                        <div class="text-xs text-(--color-text-secondary) truncate">
                            @{{ auth.user?.username }}
                        </div>
                        <span
                            class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-medium capitalize mt-0.5"
                            :class="getRoleBadgeClass(auth.user?.role || '')">
                            {{ auth.user?.role }}
                        </span>
                    </div>
                </div>
            </div>

            <UserMenuItems @change-password="handleChangePassword" @reset-all-passwords="handleResetAllPasswords"
                @download-backup="handleDownloadBackup" @logout="handleLogout" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ChevronDown } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import UserMenuItems from './UserMenuItems.vue'

const emit = defineEmits<{
    (e: 'logout'): void
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
}>()

const auth = useAuthStore()

const dropdownRef = ref<HTMLElement | null>(null)
const dropdownOpen = ref(false)

const userInitials = computed(() => {
    if (!auth.user?.name) return '?'
    return auth.user.name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
})

const getRoleBadgeClass = (role: string): string => {
    switch (role) {
        case 'admin':
            return 'bg-(--color-blue)/10 text-(--color-blue)'
        case 'manager':
            return 'bg-(--color-yellow)/10 text-(--color-yellow)'
        default:
            return 'bg-(--color-muted-bg) text-(--color-text-secondary)'
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