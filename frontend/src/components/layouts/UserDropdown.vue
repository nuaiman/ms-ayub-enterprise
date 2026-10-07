<!-- src/components/layouts/UserDropdown.vue -->
<template>
    <div ref="dropdownRef" class="relative" :class="variant === 'sidebar' ? 'w-full' : ''">
        <!-- Trigger -->
        <button @click="toggleDropdown" :class="[
            'flex items-center rounded-xl transition-colors duration-150',
            variant === 'header'
                ? 'w-8 h-8 justify-center bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 text-xs font-bold hover:bg-(--color-blue)/20'
                : collapsed
                    ? 'w-10 h-10 mx-auto justify-center bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 text-xs font-bold'
                    : 'w-full p-2 gap-2.5 bg-(--color-muted-bg)/30 hover:bg-(--color-muted-bg)/50',
        ]">
            <template v-if="variant === 'header'">
                {{ initials }}
            </template>

            <template v-else>
                <div
                    class="w-8 h-8 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-[11px] font-bold shrink-0">
                    {{ initials }}
                </div>
                <div v-if="!collapsed" class="flex-1 min-w-0 text-left">
                    <p class="text-[12.5px] font-semibold text-(--color-text-primary) truncate leading-tight">
                        {{ auth.user?.name }}
                    </p>
                    <p class="text-[10px] text-(--color-text-secondary)/60 truncate capitalize mt-0.5">
                        {{ auth.user?.role }}
                    </p>
                </div>
                <ChevronDown v-if="!collapsed"
                    class="w-3.5 h-3.5 text-(--color-text-secondary)/40 shrink-0 transition-transform duration-200"
                    :class="{ 'rotate-180': open }" :stroke-width="2.5" />
            </template>
        </button>

        <!-- Dropdown -->
        <div v-if="open"
            class="absolute bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50"
            :class="variant === 'header'
                ? 'right-0 top-10 w-56'
                : collapsed
                    ? 'bottom-full left-1/2 -translate-x-1/2 mb-2 w-56'
                    : 'bottom-full left-0 right-0 mb-2'">
            <!-- User info header -->
            <div class="px-4 py-3 border-b border-(--color-border)/40">
                <div class="flex items-center gap-3">
                    <div
                        class="w-9 h-9 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-sm font-bold shrink-0">
                        {{ initials }}
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
                            :class="roleBadgeClass">
                            {{ auth.user?.role }}
                        </span>
                    </div>
                </div>
            </div>

            <!-- Menu items -->
            <div class="py-1">
                <button @click="handleThemeToggle"
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                    <Sun v-if="!themeStore.isDark" class="w-4 h-4 shrink-0" :stroke-width="2" />
                    <Moon v-else class="w-4 h-4 shrink-0" :stroke-width="2" />
                    {{ themeStore.isDark ? 'Light Mode' : 'Dark Mode' }}
                </button>

                <div class="border-t border-(--color-border)/30 my-1"></div>

                <button @click="handleAction('changePassword')"
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                    <KeyRound class="w-4 h-4 shrink-0" :stroke-width="2" />
                    Change Password
                </button>

                <template v-if="auth.user?.role === 'admin'">
                    <div class="border-t border-(--color-border)/30 my-1"></div>

                    <button @click="handleAction('downloadBackup')"
                        class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <Download class="w-4 h-4 shrink-0" :stroke-width="2" />
                        Download Backup
                    </button>

                    <button @click="handleAction('resetAllPasswords')"
                        class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
                        <Lock class="w-4 h-4 shrink-0" :stroke-width="2" />
                        Reset All Passwords
                    </button>
                </template>

                <div class="border-t border-(--color-border)/30 my-1"></div>

                <button @click="handleAction('logout')"
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-(--color-red) hover:bg-(--color-red)/10 transition-colors">
                    <LogOut class="w-4 h-4 shrink-0" :stroke-width="2" />
                    Logout
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
    ChevronDown,
    KeyRound,
    Download,
    Lock,
    LogOut,
    Sun,
    Moon,
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useClickOutside } from '@/composables/useClickOutside'

type UserAction = 'logout' | 'changePassword' | 'resetAllPasswords' | 'downloadBackup'

withDefaults(defineProps<{
    variant?: 'header' | 'sidebar'
    collapsed?: boolean
}>(), {
    variant: 'header',
    collapsed: false,
})

const emit = defineEmits<{
    (e: UserAction): void
}>()

const auth = useAuthStore()
const themeStore = useThemeStore()

const dropdownRef = ref<HTMLElement | null>(null)
const open = ref(false)

useClickOutside(dropdownRef, () => {
    open.value = false
})

const initials = computed(() => {
    if (!auth.user?.name) return '?'
    return auth.user.name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
})

const roleBadgeClass = computed(() => {
    const role = auth.user?.role || ''
    switch (role) {
        case 'admin':
            return 'bg-(--color-blue)/10 text-(--color-blue)'
        case 'manager':
            return 'bg-(--color-yellow)/10 text-(--color-yellow)'
        default:
            return 'bg-(--color-muted-bg) text-(--color-text-secondary)'
    }
})

const toggleDropdown = () => {
    open.value = !open.value
}

const handleThemeToggle = () => {
    themeStore.toggleTheme()
}

const handleAction = (action: UserAction) => {
    open.value = false
    emit(action)
}
</script>