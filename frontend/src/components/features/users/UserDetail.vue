<!-- src/components/features/users/UserDetail.vue -->
<template>
  <div v-if="user" class="space-y-6">
    <!-- Header -->
    <div class="flex items-start gap-4">
      <!-- Avatar -->
      <div class="shrink-0">
        <div v-if="user.image_url" class="w-20 h-20 rounded-full overflow-hidden border-2 border-(--color-border)">
          <img :src="getImageUrl(user.image_url)" :alt="user.name" class="w-full h-full object-cover" />
        </div>
        <div v-else
          class="w-20 h-20 rounded-full bg-(--color-muted-bg) border-2 border-(--color-border) flex items-center justify-center">
          <span class="text-2xl font-semibold text-(--color-text-secondary)">{{ getInitials(user.name) }}</span>
        </div>
      </div>

      <!-- Info -->
      <div class="flex-1 min-w-0">
        <h2 class="text-2xl font-bold text-(--color-text-primary)">{{ user.name }}</h2>
        <div class="flex items-center gap-2 flex-wrap mt-1">
          <span class="text-sm text-(--color-text-secondary)">@{{ user.username }}</span>
          <span class="w-1 h-1 rounded-full bg-(--color-text-secondary)"></span>
          <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize"
            :class="usersStore.getRoleBadgeClass(user.role)">
            <span class="w-1 h-1 rounded-full mr-1" :class="usersStore.getRoleDotClass(user.role)"></span>
            {{ user.role }}
          </span>
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border"
            :class="user.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-yellow) text-(--color-yellow)'">
            <span class="w-1.5 h-1.5 rounded-full"
              :class="user.is_active ? 'bg-(--color-green)' : 'bg-(--color-yellow)'"></span>
            {{ user.is_active ? 'Active' : 'Inactive' }}
          </span>
        </div>
      </div>
    </div>

    <!-- Meta -->
    <div class="flex flex-wrap items-center gap-4 pb-4 border-b border-(--color-border)">
      <span class="text-xs text-(--color-text-secondary)">ID: {{ user.id }}</span>
      <span class="w-px h-4 bg-(--color-border)"></span>
      <span class="text-xs text-(--color-text-secondary)">Created: {{ usersStore.formatDate(user.created_at) }}</span>
      <span class="w-px h-4 bg-(--color-border)"></span>
      <span class="text-xs text-(--color-text-secondary)">Updated: {{ usersStore.formatDate(user.updated_at) }}</span>
    </div>

    <!-- Details -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="space-y-4">
        <!-- Email -->
        <div>
          <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Email</p>
          <p class="text-sm text-(--color-text-primary)">{{ user.email || '—' }}</p>
        </div>

        <!-- Phone -->
        <div>
          <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Phone</p>
          <p class="text-sm text-(--color-text-primary)">{{ user.phone || '—' }}</p>
        </div>

        <!-- Address -->
        <div>
          <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Address</p>
          <p class="text-sm text-(--color-text-primary)">{{ user.address || '—' }}</p>
        </div>
      </div>

      <div class="space-y-4">
        <!-- Salary -->
        <div>
          <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">Monthly Salary</p>
          <p class="text-sm font-semibold text-(--color-text-primary)">{{ formatCurrency(user.monthly_salary) }}</p>
        </div>

        <!-- ID Type -->
        <div>
          <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">ID Type</p>
          <p class="text-sm text-(--color-text-primary) capitalize">{{ user.id_type || '—' }}</p>
        </div>

        <!-- ID Number -->
        <div>
          <p class="text-xs font-medium text-(--color-text-secondary) uppercase tracking-wider">ID Number</p>
          <p class="text-sm text-(--color-text-primary)">{{ user.id_number || '—' }}</p>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex flex-wrap items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
      <button v-if="canEdit" @click="emit('edit', user)"
        class="px-4 py-2 text-sm font-medium rounded-lg bg-(--color-blue) text-white hover:opacity-90 transition-all duration-200">
        Edit User
      </button>
      <button @click="emit('close')"
        class="px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200">
        Close
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { User } from '@/types/auth'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { getImageUrl } from '@/utils/image'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
  user: User | null
}>()

const emit = defineEmits<{
  'close': []
  'edit': [user: User]
  'updated': []
}>()

const auth = useAuthStore()
const usersStore = useUsersStore()

const canEdit = computed(() => {
  const currentUser = auth.user
  if (!currentUser || !props.user) return false
  if (currentUser.role === 'admin') return true
  if (currentUser.role === 'manager' && props.user.role !== 'admin') return true
  return currentUser.id === props.user.id
})

const getInitials = (name: string): string => {
  return name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
}
</script>