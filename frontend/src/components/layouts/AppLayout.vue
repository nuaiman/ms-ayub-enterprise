<!-- src/components/layouts/AppLayout.vue -->
<template>
    <div class="min-h-screen bg-(--color-bg) text-(--color-text-primary) flex selection:bg-(--color-blue)/20">
        <AppSidebar :is-open="sidebarOpen" :is-collapsed="isCollapsed" :menu-groups="menuGroups" @close="closeSidebar"
            @hover-enter="handleSidebarHoverEnter" @hover-leave="handleSidebarHoverLeave" @logout="handleLogout"
            @change-password="handleOpenChangePassword" @reset-all-passwords="handleOpenResetAll"
            @download-backup="handleDownloadBackup" />

        <div v-if="sidebarOpen" class="fixed inset-0 z-20 bg-black/30 backdrop-blur-sm lg:hidden transition-opacity"
            @click="closeSidebar" />

        <div class="flex-1 flex flex-col min-w-0">
            <AppHeader :menu-groups="menuGroups" @toggle-sidebar="toggleSidebar" @logout="handleLogout"
                @change-password="handleOpenChangePassword" @reset-all-passwords="handleOpenResetAll"
                @download-backup="handleDownloadBackup" />

            <main class="flex-1 px-3 sm:px-4 lg:px-6 py-4 sm:py-5">
                <RouterView />
            </main>
        </div>

        <!-- Change Password Dialog -->
        <BaseDialog v-model="showChangePassword" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-blue)/10 text-(--color-blue) flex items-center justify-center shrink-0">
                        <KeyRound class="w-5 h-5" :stroke-width="2" />
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Change Password</h2>
                        <p class="text-xs text-(--color-text-secondary)">Update your account password</p>
                    </div>
                </div>

                <div class="space-y-3">
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            Current Password
                        </label>
                        <input v-model="changePasswordForm.current" type="password" placeholder="Enter current password"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
                    </div>
                    <div>
                        <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                            New Password
                        </label>
                        <input v-model="changePasswordForm.new" type="password" placeholder="Enter new password"
                            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent"
                            @keyup.enter="submitChangePassword" />
                    </div>
                </div>
            </div>

            <template #actions>
                <button type="button" @click="showChangePassword = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>
                <button type="button" @click="submitChangePassword"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-opacity">
                    Update Password
                </button>
            </template>
        </BaseDialog>

        <!-- Reset All Passwords Dialog -->
        <BaseDialog v-model="showResetAll" max-width="sm">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-yellow)/10 text-(--color-yellow) flex items-center justify-center shrink-0">
                        <Lock class="w-5 h-5" :stroke-width="2" />
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Reset All Passwords</h2>
                        <p class="text-xs text-(--color-text-secondary)">
                            This will reset passwords for all users except root
                        </p>
                    </div>
                </div>

                <div>
                    <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
                        New Password
                    </label>
                    <input v-model="resetAllForm.new" type="password" placeholder="Enter new password for all users"
                        class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent"
                        @keyup.enter="submitResetAll" />
                </div>
            </div>

            <template #actions>
                <button type="button" @click="showResetAll = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>
                <button type="button" @click="submitResetAll"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-yellow) text-white rounded-lg hover:opacity-90 transition-opacity">
                    Reset All
                </button>
            </template>
        </BaseDialog>

        <!-- Logout Dialog -->
        <BaseDialog v-model="showLogoutDialog" max-width="sm" @close="showLogoutDialog = false">
            <div class="space-y-4">
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-(--color-red)/10 text-(--color-red) flex items-center justify-center shrink-0">
                        <AlertTriangle class="w-5 h-5" :stroke-width="2" />
                    </div>
                    <div>
                        <h2 class="text-lg font-bold text-(--color-text-primary)">Confirm Logout</h2>
                        <p class="text-xs text-(--color-text-secondary)">Are you sure you want to log out?</p>
                    </div>
                </div>
                <p class="text-sm text-(--color-text-secondary)">
                    You'll need to sign in again to access your workspace.
                </p>
            </div>

            <template #actions>
                <button type="button" @click="showLogoutDialog = false"
                    class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">
                    Cancel
                </button>
                <button type="button" @click="confirmLogout"
                    class="px-4 py-2 text-sm font-semibold bg-(--color-red) text-white rounded-lg hover:opacity-90 transition-opacity">
                    Logout
                </button>
            </template>
        </BaseDialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { KeyRound, Lock, AlertTriangle } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { useSettingsStore } from '@/stores/settings'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import { push } from 'notivue'
import type { MenuGroup } from './types'

const router = useRouter()
const auth = useAuthStore()
const usersStore = useUsersStore()
const settingsStore = useSettingsStore()

const sidebarOpen = ref(false)
const showLogoutDialog = ref(false)
const showChangePassword = ref(false)
const showResetAll = ref(false)

const isMobile = ref(window.innerWidth < 1024)
const isHovered = ref(false)
let hoverLeaveTimer: ReturnType<typeof setTimeout> | null = null

const isCollapsed = computed(() => {
    if (isMobile.value) return false
    return !isHovered.value
})

const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
}

const closeSidebar = () => {
    sidebarOpen.value = false
}

const handleSidebarHoverEnter = () => {
    if (hoverLeaveTimer) {
        clearTimeout(hoverLeaveTimer)
        hoverLeaveTimer = null
    }
    if (!isMobile.value) {
        isHovered.value = true
    }
}

const handleSidebarHoverLeave = () => {
    if (isMobile.value) return
    if (hoverLeaveTimer) {
        clearTimeout(hoverLeaveTimer)
    }
    hoverLeaveTimer = setTimeout(() => {
        hoverLeaveTimer = null
        const el = document.querySelector('aside')
        if (el && el.matches(':hover')) {
            return
        }
        isHovered.value = false
    }, 120)
}

const updateMobileStatus = () => {
    isMobile.value = window.innerWidth < 1024

    if (isMobile.value) {
        isHovered.value = false
    }
}

onMounted(() => {
    window.addEventListener('resize', updateMobileStatus)
})

onUnmounted(() => {
    window.removeEventListener('resize', updateMobileStatus)
    if (hoverLeaveTimer) {
        clearTimeout(hoverLeaveTimer)
        hoverLeaveTimer = null
    }
})

// ============= ACTIONS =============
const handleLogout = () => {
    showLogoutDialog.value = true
}

const confirmLogout = async () => {
    await auth.logout()
    router.push('/auth')
}

const changePasswordForm = ref({ current: '', new: '' })

const handleOpenChangePassword = () => {
    changePasswordForm.value = { current: '', new: '' }
    showChangePassword.value = true
}

const submitChangePassword = async () => {
    if (!changePasswordForm.value.current || !changePasswordForm.value.new) {
        push.error('Please fill in all fields')
        return
    }
    if (changePasswordForm.value.new.length < 8) {
        push.error('Password must be at least 8 characters')
        return
    }
    const success = await auth.changePassword(
        changePasswordForm.value.current,
        changePasswordForm.value.new
    )
    if (success) {
        changePasswordForm.value = { current: '', new: '' }
        showChangePassword.value = false
    }
}

const resetAllForm = ref({ new: '' })

const handleOpenResetAll = () => {
    resetAllForm.value = { new: '' }
    showResetAll.value = true
}

const submitResetAll = async () => {
    if (!resetAllForm.value.new) {
        push.error('Please enter a new password')
        return
    }
    if (resetAllForm.value.new.length < 8) {
        push.error('Password must be at least 8 characters')
        return
    }
    const success = await usersStore.resetAllPasswords(resetAllForm.value.new)
    if (success) {
        resetAllForm.value = { new: '' }
        showResetAll.value = false
    }
}

const handleDownloadBackup = () => {
    settingsStore.downloadBackup()
}

// ============= MENU GROUPS =============
const menuGroups: MenuGroup[] = [
    {
        id: 'parties',
        label: 'Parties',
        icon: 'parties',
        items: [
            { path: '/godowns', label: 'Godowns', icon: 'godowns' },
            { path: '/majhis', label: 'Majhis', icon: 'majhis' },
            { path: '/brokers', label: 'Brokers', icon: 'brokers' },
        ],
    },
    {
        id: 'warehouse',
        label: 'Warehouse',
        icon: 'warehouse',
        items: [
            { path: '/lots', label: 'Lots', icon: 'lots' },
            { path: '/deliveries', label: 'Deliveries', icon: 'deliveries' },
            { path: '/store-adjustments', label: 'Store Adjustments', icon: 'store-adjustments' },
            { path: '/store-transfers', label: 'Store Transfers', icon: 'store-transfers' },
            { path: '/lot-transfers', label: 'Lot Transfers', icon: 'lot-transfers' },
            { path: '/damages', label: 'Damages', icon: 'damages' },
        ],
    },
    {
        id: 'bills',
        label: 'Bills',
        icon: 'bills',
        items: [
            { path: '/invoices', label: 'Invoices', icon: 'invoices' },
            { path: '/godown-bills', label: 'Godown Bills', icon: 'godown-bills' },
            { path: '/majhi-bills', label: 'Majhi Bills', icon: 'majhi-bills' },
            { path: '/customer-store-bills', label: 'Customer Store Bills', icon: 'customer-store-bills' },
            { path: '/customer-delivery-bills', label: 'Customer Delivery Bills', icon: 'customer-delivery-bills' },
            { path: '/customer-additional-bills', label: 'Customer Additional Bills', icon: 'customer-additional-bills' },
        ],
    },
    {
        id: 'office',
        label: 'Office',
        icon: 'office',
        items: [
            { path: '/users', label: 'Users', icon: 'users' },
            { path: '/salaries', label: 'Salaries', icon: 'salaries' },
            { path: '/incomes', label: 'Incomes', icon: 'incomes' },
            { path: '/expenses', label: 'Expenses', icon: 'expenses' },
            { path: '/logs', label: 'Logs', icon: 'logs' },
        ],
    },
]
</script>