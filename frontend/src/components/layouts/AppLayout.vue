<!-- src/components/layouts/AppLayout.vue -->
<template>
    <div class="min-h-screen bg-(--color-bg) text-(--color-text-primary) flex selection:bg-(--color-blue)/20">
        <!-- Sidebar -->
        <AppSidebar :is-open="sidebarOpen" :is-collapsed="isCollapsed" :is-mobile="isMobile" :menu-groups="menuGroups"
            @toggle="toggleSidebar" @close="closeSidebar" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave"
            @logout="handleLogout" @change-password="handleOpenChangePassword" @reset-all-passwords="handleOpenResetAll"
            @download-backup="handleDownloadBackup" />

        <!-- Mobile Overlay -->
        <div v-if="sidebarOpen" class="fixed inset-0 z-20 bg-black/30 backdrop-blur-sm lg:hidden transition-opacity"
            @click="closeSidebar" />

        <!-- Main Content -->
        <div class="flex-1 flex flex-col min-w-0">
            <!-- Header -->
            <AppHeader @toggle-sidebar="toggleSidebar" @logout="handleLogout"
                @change-password="handleOpenChangePassword" @reset-all-passwords="handleOpenResetAll"
                @download-backup="handleDownloadBackup" />

            <!-- Page Content -->
            <main class="flex-1 px-3 sm:px-4 lg:px-6 py-4 sm:py-5">
                <slot />
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
                    You'll need to sign in again to access your dashboard.
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
import { useRouter } from 'vue-router'
import { KeyRound, Lock, AlertTriangle } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { useSettingsStore } from '@/stores/settings'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import { push } from 'notivue'

const router = useRouter()
const auth = useAuthStore()
const usersStore = useUsersStore()
const settingsStore = useSettingsStore()

// ============= STATE =============
const sidebarOpen = ref(false)
const showLogoutDialog = ref(false)
const showChangePassword = ref(false)
const showResetAll = ref(false)

// ============= RESPONSIVE =============
const isMobile = ref(window.innerWidth < 1024)
const isHovered = ref(false)

// Collapsed only applies on desktop
const isCollapsed = computed(() => {
    if (isMobile.value) return false
    return !isHovered.value
})

// ============= SIDEBAR HANDLERS =============
const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
}

const closeSidebar = () => {
    sidebarOpen.value = false
}

const handleMouseEnter = () => {
    if (!isMobile.value) {
        isHovered.value = true
    }
}

const handleMouseLeave = () => {
    if (!isMobile.value) {
        isHovered.value = false
    }
}

// ============= RESPONSIVE =============
const updateMobileStatus = () => {
    const wasMobile = isMobile.value
    isMobile.value = window.innerWidth < 1024

    if (wasMobile && !isMobile.value && sidebarOpen.value) {
        sidebarOpen.value = false
    }
}

onMounted(() => {
    window.addEventListener('resize', updateMobileStatus)
})

onUnmounted(() => {
    window.removeEventListener('resize', updateMobileStatus)
})

// ============= LOGOUT =============
const handleLogout = () => {
    showLogoutDialog.value = true
}

const confirmLogout = async () => {
    await auth.logout()
    router.push('/auth')
}

// ============= CHANGE PASSWORD =============
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

// ============= RESET ALL PASSWORDS =============
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

// ============= BACKUP =============
const handleDownloadBackup = () => {
    settingsStore.downloadBackup()
}

// ============= MENU GROUPS =============
const menuGroups = [
    {
        id: 'parties',
        label: 'Parties',
        icon: 'parties',
        items: [
            { path: '/brokers', label: 'Brokers', icon: 'brokers' },
            { path: '/majhis', label: 'Majhis', icon: 'majhis' },
            { path: '/godowns', label: 'Godowns', icon: 'godowns' },
        ],
    },
    {
        id: 'warehouse',
        label: 'Warehouse',
        icon: 'warehouse',
        items: [
            { path: '/lots', label: 'Lots', icon: 'lots' },
            { path: '/stores', label: 'Stores', icon: 'stores' },
            { path: '/damages', label: 'Damages', icon: 'damages' },
            { path: '/deliveries', label: 'Deliveries', icon: 'deliveries' },
            { path: '/delivery-items', label: 'Delivery Items', icon: 'delivery-items' },
        ],
    },
    {
        id: 'transport',
        label: 'Transport',
        icon: 'transport',
        items: [
            { path: '/transports', label: 'Transports', icon: 'transports' },
            { path: '/vehicles', label: 'Vehicles', icon: 'vehicles' },
        ],
    },
    {
        id: 'bills',
        label: 'Bills',
        icon: 'bills',
        items: [
            { path: '/godown-store-bills', label: 'Godown Store Bills', icon: 'godown-store-bills' },
            { path: '/customer-storage-bills', label: 'Customer Storage Bills', icon: 'customer-storage-bills' },
            { path: '/customer-lot-bills', label: 'Customer Unload Bills', icon: 'customer-lot-bills' },
            { path: '/majhi-lot-bills', label: 'Majhi Lot Bills', icon: 'majhi-lot-bills' },
            { path: '/customer-delivery-bills', label: 'Customer Delivery Bills', icon: 'customer-delivery-bills' },
            { path: '/majhi-loading-bills', label: 'Majhi Loading Bills', icon: 'majhi-loading-bills' },
            { path: '/broker-vehicle-bills', label: 'Broker Vehicle Bills', icon: 'broker-vehicle-bills' },
            { path: '/customer-transport-bills', label: 'Customer Transport Bills', icon: 'customer-transport-bills' },
        ],
    },
    {
        id: 'invoices',
        label: 'Invoices',
        icon: 'invoices',
        items: [
            { path: '/invoices', label: 'Create Invoice', icon: 'invoices' },
        ],
    },
    {
        id: 'office',
        label: 'Office',
        icon: 'office',
        items: [
            { path: '/users', label: 'Users', icon: 'users' },
            { path: '/salaries', label: 'Salaries', icon: 'salaries' },
            { path: '/expenses', label: 'Expenses', icon: 'expenses' },
            { path: '/logs', label: 'Logs', icon: 'logs' },
        ],
    },
]
</script>