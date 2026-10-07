<!-- src/components/layouts/AppHeader.vue -->
<template>
    <header
        class="sticky top-0 z-10 bg-(--color-surface)/80 backdrop-blur-xl border-b border-(--color-border)/40 h-14 flex items-center px-4 sm:px-6">
        <button @click="$emit('toggleSidebar')"
            class="lg:hidden p-2 -ml-2 rounded-lg text-(--color-text-secondary) hover:bg-(--color-muted-bg) hover:text-(--color-text-primary) transition-colors"
            aria-label="Toggle sidebar">
            <Menu class="w-5 h-5" :stroke-width="2" />
        </button>

        <div class="flex-1 min-w-0 ml-2 lg:ml-0">
            <nav aria-label="Breadcrumb" class="flex items-center gap-1 text-xs text-(--color-text-secondary)/70">
                <div v-for="(crumb, idx) in breadcrumbs" :key="idx" class="flex items-center gap-1 min-w-0">
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

        <div class="ml-auto shrink-0">
            <UserDropdown variant="header" @logout="$emit('logout')" @change-password="$emit('changePassword')"
                @reset-all-passwords="$emit('resetAllPasswords')" @download-backup="$emit('downloadBackup')" />
        </div>
    </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Menu, ChevronRight } from 'lucide-vue-next'
import UserDropdown from './UserDropdown.vue'
import type { MenuGroup } from './types'

const props = defineProps<{
    menuGroups: MenuGroup[]
}>()

defineEmits<{
    (e: 'toggleSidebar'): void
    (e: 'logout'): void
    (e: 'changePassword'): void
    (e: 'resetAllPasswords'): void
    (e: 'downloadBackup'): void
}>()

const route = useRoute()

const titleMap: Record<string, string> = {
    '/customers': 'Customers',
    '/brokers': 'Brokers',
    '/majhis': 'Majhis',
    '/godowns': 'Godowns',
    '/lots': 'Lots',
    '/stores': 'Stores',
    '/users': 'Users',
    '/salaries': 'Salaries',
    '/expenses': 'Expenses',
    '/logs': 'Logs',
    '/support': 'Support',
    '/majhi-bills': 'Majhi Bills',
    '/godown-bills': 'Godown Bills',
    '/customer-store-bills': 'Customer Store Bills',
}

const pageTitle = computed(() => {
    let best = ''
    for (const path of Object.keys(titleMap)) {
        if (route.path.startsWith(path) && path.length > best.length) {
            best = path
        }
    }
    return best ? (titleMap[best] ?? 'MS Ayub Enterprise') : 'MS Ayub Enterprise'
})

interface Breadcrumb {
    label: string
    path?: string
}

const breadcrumbs = computed<Breadcrumb[]>(() => {
    if (route.path === '/customers') return [{ label: 'Customers' }]
    if (route.path === '/support') return [{ label: 'Support' }]

    const items: Breadcrumb[] = []
    const activeGroup = props.menuGroups.find(g =>
        g.items.some(item => route.path.startsWith(item.path))
    )
    if (activeGroup) items.push({ label: activeGroup.label })
    items.push({ label: pageTitle.value })
    return items
})
</script>