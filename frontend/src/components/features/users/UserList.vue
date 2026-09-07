<!-- src/components/features/users/UserList.vue -->
<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-(--color-text-primary)">Users</h2>
        <span class="text-sm text-(--color-text-secondary) bg-(--color-muted-bg) px-2 py-0.5 rounded-md">
          {{ usersStore.users.length }}
        </span>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <!-- Search -->
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <input :value="usersStore.searchQuery" @input="handleSearch" type="text" placeholder="Search users..."
            class="w-full sm:w-56 pl-9 pr-8 py-2 rounded-lg text-sm bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--color-text-secondary)" fill="none"
            stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <button v-if="usersStore.searchQuery" @click="usersStore.clearSearch" type="button"
            class="absolute right-2.5 top-1/2 -translate-y-1/2 text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Copy Button -->
        <button @click="handleCopyToClipboard"
          class="h-9 w-9 flex items-center justify-center border border-(--color-border) rounded-lg text-(--color-text-secondary) hover:bg-(--color-muted-bg) transition-colors relative shrink-0"
          title="Copy table to clipboard">
          <svg v-if="!clipboardStore.copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <svg v-else class="w-4 h-4 text-(--color-green)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </button>

        <!-- Create -->
        <button v-if="canCreate" @click="createDialogOpen = true"
          class="h-9 px-4 flex items-center gap-2 bg-(--color-blue) text-white rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap shrink-0">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
          </svg>
          <span class="hidden sm:inline">Create</span>
        </button>
      </div>
    </div>

    <!-- Table - overflow-x-auto for horizontal scroll on small screens -->
    <div class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-3xl">
        <!-- Header Row -->
        <div
          class="grid grid-cols-12 items-center py-3 px-3 border-b border-(--color-border) text-xs font-semibold text-(--color-text-secondary) uppercase tracking-wider bg-(--color-muted-bg)/30 rounded-t-lg">
          <div class="col-span-3">Name</div>
          <div class="col-span-2">Role</div>
          <div class="col-span-1">Salary</div>
          <div class="col-span-2">Status</div>
          <div class="col-span-2">Created</div>
          <div class="col-span-2 text-right">Actions</div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="text-center space-y-4">
            <svg class="animate-spin w-10 h-10 text-(--color-blue) mx-auto" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            <p class="text-sm text-(--color-text-secondary)">Loading users...</p>
          </div>
        </div>

        <!-- Empty -->
        <div v-else-if="usersStore.filteredUsers.length === 0" class="flex items-center justify-center py-12">
          <div class="text-center space-y-3">
            <div class="w-16 h-16 mx-auto rounded-full bg-(--color-muted-bg) flex items-center justify-center">
              <svg class="w-8 h-8 text-(--color-text-secondary)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-(--color-text-primary)">No users found</p>
            <p class="text-xs text-(--color-text-secondary)">{{ usersStore.searchQuery ? 'Try adjusting your search' :
              'Create a new user to get started' }}</p>
          </div>
        </div>

        <!-- Rows -->
        <div v-else>
          <UserRow v-for="user in usersStore.filteredUsers" :key="user.id" :user="user" @view="openDetailDialog"
            @edit="handleEditUser" @change-password="handleChangePassword" @change-role="handleChangeRole"
            @toggle-active="handleToggleActive" @updated="fetchUsers" />
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div v-if="!loading && usersStore.filteredUsers.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-(--color-border) shrink-0">
      <p class="text-xs text-(--color-text-secondary)">Showing {{ usersStore.filteredUsers.length }} of {{
        usersStore.users.length }} users</p>
    </div>

    <!-- Dialogs -->
    <!-- Create Dialog -->
    <BaseDialog v-model="createDialogOpen" max-width="3xl">
      <div class="mb-6">
        <h2 class="text-xl font-bold text-(--color-text-primary)">Create New User</h2>
        <p class="text-sm text-(--color-text-secondary) mt-1">Add a new team member to the system</p>
      </div>
      <UserForm mode="create" @user-created="handleUserCreated" @cancel="createDialogOpen = false" />
    </BaseDialog>

    <!-- Detail Dialog -->
    <BaseDialog v-model="detailDialogOpen" max-width="3xl">
      <UserDetail v-if="selectedUser" :user="selectedUser" @close="closeDetailDialog" @edit="handleEditUserFromDetail"
        @updated="fetchUsers" />
    </BaseDialog>

    <!-- Edit Dialog -->
    <BaseDialog v-model="editDialogOpen" max-width="3xl">
      <div class="mb-6">
        <h2 class="text-xl font-bold text-(--color-text-primary)">Edit User</h2>
        <p class="text-sm text-(--color-text-secondary) mt-1">Update user information</p>
      </div>
      <UserForm v-if="selectedUser" mode="edit" :user="selectedUser" @user-updated="handleUserUpdated"
        @cancel="editDialogOpen = false" />
    </BaseDialog>

    <!-- Change Password Dialog -->
    <BaseDialog v-model="passwordDialogOpen" max-width="sm">
      <div class="mb-4">
        <h2 class="text-lg font-bold text-(--color-text-primary)">Change Password</h2>
        <p class="text-sm text-(--color-text-secondary) mt-0.5">For <span
            class="font-medium text-(--color-text-primary)">{{ selectedUser?.name }}</span></p>
      </div>
      <div class="space-y-4">
        <input v-model="newPassword" type="password" placeholder="Enter new password"
          class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent"
          @keyup.enter="submitPasswordChange" />
      </div>
      <template #actions>
        <button @click="passwordDialogOpen = false"
          class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
        <button @click="submitPasswordChange"
          class="px-4 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-colors">Update</button>
      </template>
    </BaseDialog>

    <!-- Change Role Dialog -->
    <BaseDialog v-model="roleDialogOpen" max-width="sm">
      <div class="mb-4">
        <h2 class="text-lg font-bold text-(--color-text-primary)">Change Role</h2>
        <p class="text-sm text-(--color-text-secondary) mt-0.5">For <span
            class="font-medium text-(--color-text-primary)">{{
              selectedUser?.name }}</span></p>
      </div>
      <div class="space-y-4">
        <select v-model="newRole"
          class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
          <option value="staff">Staff</option>
          <option value="accounts">Accounts</option>
          <option value="manager">Manager</option>
        </select>
      </div>
      <template #actions>
        <button @click="roleDialogOpen = false"
          class="px-4 py-2 text-sm rounded-lg hover:bg-(--color-muted-bg) transition-colors">Cancel</button>
        <button @click="submitRoleChange"
          class="px-4 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-colors">Update</button>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUsersStore } from '@/stores/users'
import { useAuthStore } from '@/stores/auth'
import { useClipboardStore } from '@/stores/clipboard'
import type { User, Role } from '@/types/auth'
import UserRow from './UserRow.vue'
import UserForm from './UserForm.vue'
import UserDetail from './UserDetail.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const usersStore = useUsersStore()
const auth = useAuthStore()
const clipboardStore = useClipboardStore()

const loading = ref(true)
const createDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const editDialogOpen = ref(false)
const passwordDialogOpen = ref(false)
const roleDialogOpen = ref(false)
const selectedUser = ref<User | null>(null)
const newPassword = ref('')
const newRole = ref<Role>('staff')

const canCreate = computed(() => {
  const role = auth.user?.role
  return role === 'admin' || role === 'manager'
})

const fetchUsers = async () => {
  loading.value = true
  try {
    await usersStore.fetchUsers()
  } finally {
    loading.value = false
  }
}

const handleSearch = (e: Event) => {
  const target = e.target as HTMLInputElement
  usersStore.setSearchQuery(target.value)
}

const handleCopyToClipboard = async () => {
  const data = usersStore.generateClipboardData()
  if (data) {
    await clipboardStore.copyToClipboard(data)
  }
}

const openDetailDialog = (user: User) => {
  selectedUser.value = user
  detailDialogOpen.value = true
}

const closeDetailDialog = () => {
  detailDialogOpen.value = false
  setTimeout(() => {
    selectedUser.value = null
  }, 300)
}

const handleEditUser = (user: User) => {
  selectedUser.value = user
  editDialogOpen.value = true
}

const handleEditUserFromDetail = (user: User) => {
  detailDialogOpen.value = false
  setTimeout(() => {
    selectedUser.value = user
    editDialogOpen.value = true
  }, 300)
}

const handleChangePassword = (user: User) => {
  selectedUser.value = user
  newPassword.value = ''
  passwordDialogOpen.value = true
}

const handleChangeRole = (user: User) => {
  selectedUser.value = user
  newRole.value = user.role
  roleDialogOpen.value = true
}

const handleToggleActive = async (user: User) => {
  await usersStore.toggleActive(user.id)
}

const submitPasswordChange = async () => {
  if (!selectedUser.value || !newPassword.value) return
  if (newPassword.value.length < 8) {
    push.error('Password must be at least 8 characters')
    return
  }
  const success = await usersStore.changeUserPassword(selectedUser.value.id, newPassword.value)
  if (success) {
    passwordDialogOpen.value = false
    await fetchUsers()
  }
}

const submitRoleChange = async () => {
  if (!selectedUser.value) return
  const success = await usersStore.changeRole(selectedUser.value.id, newRole.value)
  if (success) {
    roleDialogOpen.value = false
    await fetchUsers()
  }
}

const handleUserCreated = async () => {
  createDialogOpen.value = false
  await fetchUsers()
}

const handleUserUpdated = async () => {
  editDialogOpen.value = false
  detailDialogOpen.value = false
  await fetchUsers()
}

onMounted(() => {
  fetchUsers()
})
</script>