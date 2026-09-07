<!-- src/components/features/users/UserRow.vue -->
<template>
  <div 
    class="grid grid-cols-12 items-center py-3 px-3 border-b border-(--color-border) transition-all duration-200 hover:bg-(--color-muted-bg)/30 cursor-pointer"
    @click="handleClick"
  >
    <!-- Name -->
    <div class="col-span-3">
      <div class="flex items-center gap-3">
        <div class="shrink-0">
          <div v-if="user.image_url" class="w-9 h-9 rounded-full overflow-hidden border border-(--color-border)">
            <img :src="getImageUrl(user.image_url)" :alt="user.name" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-9 h-9 rounded-full bg-(--color-muted-bg) border border-(--color-border) flex items-center justify-center">
            <span class="text-xs font-semibold text-(--color-text-secondary)">{{ getInitials(user.name) }}</span>
          </div>
        </div>
        <div class="min-w-0">
          <div class="font-medium text-(--color-text-primary) truncate text-sm">{{ user.name }}</div>
          <div class="text-xs text-(--color-text-secondary) truncate">@{{ user.username }}</div>
        </div>
      </div>
    </div>

    <!-- Role -->
    <div class="col-span-2">
      <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize" :class="usersStore.getRoleBadgeClass(user.role)">
        <span class="w-1 h-1 rounded-full mr-1" :class="usersStore.getRoleDotClass(user.role)"></span>
        {{ user.role }}
      </span>
    </div>

    <!-- Salary -->
    <div class="col-span-1">
      <span class="text-xs text-(--color-text-secondary)">{{ formatCurrency(user.monthly_salary) }}</span>
    </div>

    <!-- Status -->
    <div class="col-span-2">
      <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium border" :class="user.is_active ? 'border-(--color-green) text-(--color-green)' : 'border-(--color-yellow) text-(--color-yellow)'">
        <span class="w-1.5 h-1.5 rounded-full" :class="user.is_active ? 'bg-(--color-green)' : 'bg-(--color-yellow)'"></span>
        {{ user.is_active ? 'Active' : 'Inactive' }}
      </span>
    </div>

    <!-- Created -->
    <div class="col-span-2">
      <div class="text-xs text-(--color-text-secondary) leading-tight">{{ usersStore.formatDate(user.created_at) }}</div>
      <div class="text-[10px] text-(--color-text-secondary)/60 leading-tight">{{ usersStore.formatTime(user.created_at) }}</div>
    </div>

    <!-- Actions -->
    <div class="col-span-2 flex justify-end relative" @click.stop>
      <button 
        @click="toggleMenu"
        class="w-7 h-7 flex items-center justify-center border border-(--color-border) rounded-md hover:bg-(--color-muted-bg) transition-all duration-200"
      >
        <svg class="w-3.5 h-3.5 text-(--color-text-secondary)" fill="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="5" r="1.5" />
          <circle cx="12" cy="12" r="1.5" />
          <circle cx="12" cy="19" r="1.5" />
        </svg>
      </button>

      <!-- Dropdown -->
      <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 scale-100 translate-y-0" leave-to-class="opacity-0 scale-95 translate-y-1">
        <div v-if="isOpen" class="absolute right-0 top-9 w-48 bg-(--color-surface) border border-(--color-border) rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <button @click="handleView" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            View Details
          </button>

          <button v-if="canEdit" @click="handleEdit" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit Details
          </button>

          <button v-if="canChangePassword" @click="handleChangePassword" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
            Change Password
          </button>

          <button v-if="canChangeRole" @click="handleChangeRole" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            Change Role
          </button>

          <button v-if="canToggleActive" @click="handleToggleActive" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-(--color-muted-bg) transition-colors" :class="user.is_active ? 'text-(--color-yellow)' : 'text-(--color-green)'">
            <svg v-if="user.is_active" class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
            <svg v-else class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ user.is_active ? 'Deactivate' : 'Activate' }}
          </button>
        </div>
      </Transition>

      <!-- Backdrop -->
      <div v-if="isOpen" class="fixed inset-0 z-40" @click="closeMenu"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { User } from '@/types/auth'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { getImageUrl } from '@/utils/image'
import { formatCurrency } from '@/utils/currency'

const props = defineProps<{
  user: User
}>()

const emit = defineEmits<{
  'view': [user: User]
  'edit': [user: User]
  'change-password': [user: User]
  'change-role': [user: User]
  'toggle-active': [user: User]
  'updated': []
}>()

const auth = useAuthStore()
const usersStore = useUsersStore()

const isOpen = ref(false)

const currentUser = computed(() => auth.user)

const canEdit = computed(() => {
  if (!currentUser.value) return false
  if (currentUser.value.role === 'admin') return true
  if (currentUser.value.role === 'manager' && props.user.role !== 'admin') return true
  return currentUser.value.id === props.user.id
})

const canChangePassword = computed(() => {
  if (!currentUser.value) return false
  if (currentUser.value.role === 'admin') return true
  if (currentUser.value.role === 'manager' && props.user.role !== 'admin') return true
  return currentUser.value.id === props.user.id
})

const canChangeRole = computed(() => {
  if (!currentUser.value) return false
  if (props.user.role === 'admin') return false
  return currentUser.value.role === 'admin'
})

const canToggleActive = computed(() => {
  if (!currentUser.value) return false
  if (props.user.id === 1) return false
  if (props.user.role === 'admin') return false
  return currentUser.value.role === 'admin' || currentUser.value.role === 'manager'
})

const getInitials = (name: string): string => {
  return name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
}

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const closeMenu = () => {
  isOpen.value = false
}

const handleClick = () => {
  emit('view', props.user)
}

const handleView = () => {
  closeMenu()
  emit('view', props.user)
}

const handleEdit = () => {
  closeMenu()
  emit('edit', props.user)
}

const handleChangePassword = () => {
  closeMenu()
  emit('change-password', props.user)
}

const handleChangeRole = () => {
  closeMenu()
  emit('change-role', props.user)
}

const handleToggleActive = async () => {
  closeMenu()
  const success = await usersStore.toggleActive(props.user.id)
  if (success) {
    emit('updated')
  }
}
</script>