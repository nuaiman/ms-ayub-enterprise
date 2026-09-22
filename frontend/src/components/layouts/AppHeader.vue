<!-- src/components/layouts/AppHeader.vue -->
<template>
    <header
        class="sticky top-0 z-10 bg-(--color-surface)/80 backdrop-blur-xl border-b border-(--color-border)/40 h-14 flex items-center px-4 sm:px-6">
        <!-- Mobile Menu Toggle -->
        <button @click="$emit('toggleSidebar')"
            class="lg:hidden p-2 -ml-2 rounded-lg text-(--color-text-secondary) hover:bg-(--color-muted-bg) hover:text-(--color-text-primary) transition-colors"
            aria-label="Toggle sidebar">
            <Menu class="w-5 h-5" :stroke-width="2" />
        </button>

        <!-- Breadcrumbs -->
        <div class="flex-1 min-w-0 ml-2 lg:ml-0">
            <nav aria-label="Breadcrumb" class="flex items-center gap-1 text-xs text-(--color-text-secondary)/70">
                <div v-for="(crumb, idx) in breadcrumbs" :key="idx" class="flex items-center gap-1">
                    <router-link v-if="crumb.path && idx < breadcrumbs.length - 1" :to="crumb.path"
                        class="hover:text-(--color-blue) transition-colors truncate max-w-24">
                        {{ crumb.label }}
                    </router-link>
                    <span v-else class="text-(--color-text-primary) font-medium truncate max-w-32">
                        {{ crumb.label }}
                    </span>
                    <ChevronRight v-if="idx < breadcrumbs.length - 1" class="w-3 h-3 text-(--color-border) shrink-0"
                        :stroke-width="2" />
                </div>
            </nav>
        </div>

        <!-- Right Actions - User Menu -->
        <div ref="dropdownRef" class="flex items-center gap-1.5 ml-auto shrink-0">
            <div class="relative">
                <button @click="toggleDropdown"
                    class="w-8 h-8 rounded-lg bg-(--color-blue)/15 text-(--color-blue) border border-(--color-blue)/15 flex items-center justify-center text-xs font-bold shrink-0 hover:bg-(--color-blue)/20 transition-colors">
                    {{ userInitials }}
                </button>

                <!-- Dropdown -->
                <div v-if="dropdownOpen"
                    class="absolute right-0 top-10 w-56 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50">
                    <!-- User Info -->
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

                    <!-- Menu Items -->
                    <UserMenuItems @change-password="handleChangePassword"
                        @reset-all-passwords="handleResetAllPasswords" @download-backup="handleDownloadBackup"
                        @logout="handleLogout" />
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { Menu, ChevronRight } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useClickOutside } from '@/composables/useClickOutside'
import UserMenuItems from './UserMenuItems.vue'

const emit = defineEmits<{
    (e: 'toggleSidebar'): void
    (e: 'logout'): void
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
}>()

const route = useRoute()
const auth = useAuthStore()

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

const handleLogout = () => {
    dropdownOpen.value = false
    emit('logout')
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

const pageTitle = computed(() => {
    const titleMap: Record<string, string> = {
        '/dashboard': 'Dashboard',
        '/customers': 'Customers',
        '/brokers': 'Brokers',
        '/majhis': 'Majhis',
        '/godowns': 'Godowns',
        '/lots': 'Lots',
        '/stores': 'Stores',
        '/damages': 'Damages',
        '/deliveries': 'Deliveries',
        '/delivery-items': 'Delivery Items',
        '/transports': 'Transports',
        '/vehicles': 'Vehicles',
        '/godown-store-bills': 'Godown Store Bills',
        '/customer-storage-bills': 'Customer Storage Bills',
        '/customer-lot-bills': 'Customer Unload Bills',
        '/majhi-lot-bills': 'Majhi Lot Bills',
        '/customer-delivery-bills': 'Customer Delivery Bills',
        '/majhi-loading-bills': 'Majhi Loading Bills',
        '/broker-vehicle-bills': 'Broker Vehicle Bills',
        '/customer-transport-bills': 'Customer Transport Bills',
        '/invoices': 'Invoices',
        '/users': 'Users',
        '/expenses': 'Expenses',
        '/logs': 'Logs',
        '/support': 'Support',
    }
    for (const [path, title] of Object.entries(titleMap)) {
        if (route.path.startsWith(path)) return title
    }
    return 'MS Ayub Enterprise'
})

interface Breadcrumb {
    label: string
    path?: string
}

const breadcrumbs = computed<Breadcrumb[]>(() => {
    if (route.path === '/dashboard') {
        return [{ label: 'Dashboard' }]
    }
    if (route.path === '/customers') {
        return [{ label: 'Customers' }]
    }
    if (route.path === '/support') {
        return [{ label: 'Support' }]
    }

    const menuGroups = [
        { id: 'warehouse', label: 'Warehouse', items: [] },
        { id: 'transport', label: 'Transport', items: [] },
        { id: 'payments', label: 'Payments', items: [] },
        { id: 'office', label: 'Office', items: [] },
    ]

    const items: Breadcrumb[] = []
    const activeGroup = menuGroups.find(group =>
        group.items.some((item: any) => route.path.startsWith(item.path))
    )
    if (activeGroup) {
        items.push({ label: activeGroup.label })
    }
    items.push({ label: pageTitle.value })
    return items
})
</script>