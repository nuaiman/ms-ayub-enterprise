<!-- src/views/UsersView.vue -->
<template>
    <div class="flex flex-col h-full min-h-[calc(100vh-120px)]">
        <!-- Stats -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 shrink-0">
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Total Users</p>
                <p class="text-2xl font-bold text-(--color-text-primary) mt-1">{{ usersStore.users.length }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Active</p>
                <p class="text-2xl font-bold text-(--color-green) mt-1">{{ activeUsers }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Admins</p>
                <p class="text-2xl font-bold text-(--color-blue) mt-1">{{ adminCount }}</p>
            </div>
            <div class="rounded-xl p-4 bg-(--color-surface) border border-(--color-border)">
                <p class="text-xs text-(--color-text-secondary) uppercase tracking-wider">Managers</p>
                <p class="text-2xl font-bold text-(--color-yellow) mt-1">{{ managerCount }}</p>
            </div>
        </div>

        <!-- User List -->
        <div class="flex-1 min-h-0 mt-6">
            <UserList />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useUsersStore } from '@/stores/users'
import UserList from '@/components/features/users/UserList.vue'

const usersStore = useUsersStore()

const activeUsers = computed(() => {
    return usersStore.users.filter(u => u.is_active).length
})

const adminCount = computed(() => {
    return usersStore.users.filter(u => u.role === 'admin').length
})

const managerCount = computed(() => {
    return usersStore.users.filter(u => u.role === 'manager').length
})

onMounted(() => {
    usersStore.fetchUsers()
})
</script>