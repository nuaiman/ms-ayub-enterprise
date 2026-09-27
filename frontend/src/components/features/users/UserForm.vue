<!-- src/components/features/users/UserForm.vue -->
<template>
  <form @submit.prevent="submit" class="space-y-6">
    <!-- Personal Information -->
    <div>
      <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
        Personal Information
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="md:col-span-2">
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
            Full Name <span class="text-(--color-red)">*</span>
          </label>
          <input v-model="form.name" type="text" placeholder="John Doe" required :disabled="!!justCreatedId"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Email</label>
          <input v-model="form.email" type="email" placeholder="john@example.com" :disabled="!!justCreatedId"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Phone</label>
          <input v-model="form.phone" type="tel" placeholder="+1 234 567 890" :disabled="!!justCreatedId"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>

        <div class="md:col-span-2">
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Address</label>
          <textarea v-model="form.address" rows="2" placeholder="Enter full address" :disabled="!!justCreatedId"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent resize-none disabled:opacity-50 disabled:cursor-not-allowed"></textarea>
        </div>
      </div>
    </div>

    <!-- Account Settings (hidden when retrying image) -->
    <div v-if="!justCreatedId">
      <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
        Account Settings
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
            Username <span class="text-(--color-red)">*</span>
          </label>
          <input v-model="form.username" type="text" placeholder="johndoe" required :disabled="isEditMode"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
          <p v-if="isEditMode" class="text-xs text-(--color-text-secondary) mt-1">Username cannot be changed</p>
        </div>

        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
            Password <span v-if="!isEditMode" class="text-(--color-red)">*</span>
            <span v-else class="text-xs font-normal text-(--color-text-secondary)">(Leave blank to keep
              current)</span>
          </label>
          <div class="relative">
            <input v-model="form.password" :type="showPassword ? 'text' : 'password'"
              :placeholder="isEditMode ? 'Enter new password' : 'Enter password'" :required="!isEditMode"
              class="w-full px-3 py-2 pr-10 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
            <button type="button" @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-(--color-text-secondary) hover:text-(--color-text-primary) transition-colors">
              <svg v-if="!showPassword" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
        </div>

        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">
            Role <span class="text-(--color-red)">*</span>
          </label>
          <select v-model="form.role" required
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
            <option value="staff">Staff</option>
            <option value="accounts">Accounts</option>
            <option v-if="canCreateManager" value="manager">Manager</option>
          </select>
        </div>

        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Monthly Salary</label>
          <div class="relative">
            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-(--color-text-secondary)">৳</span>
            <input v-model.number="form.monthly_salary" type="number" step="0.01" min="0" placeholder="0.00"
              class="w-full pl-7 pr-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent" />
          </div>
        </div>
      </div>
    </div>

    <!-- Identification (hidden when retrying image) -->
    <div v-if="!justCreatedId">
      <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
        Identification <span class="normal-case text-xs font-normal">(Optional)</span>
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">ID Type</label>
          <select v-model="form.id_type"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent">
            <option value="">Select ID Type</option>
            <option value="nid">NID</option>
            <option value="passport">Passport</option>
            <option value="driving_license">Driving License</option>
            <option value="birth_certificate">Birth Certificate</option>
            <option value="trade_license">Trade License</option>
            <option value="other">Other</option>
          </select>
        </div>

        <div>
          <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">ID Number</label>
          <input v-model="form.id_number" type="text" placeholder="Enter ID number" :disabled="!form.id_type"
            class="w-full px-3 py-2 rounded-lg bg-(--color-muted-bg) border border-(--color-border) text-(--color-text-primary) placeholder:text-(--color-text-secondary) focus:outline-none focus:ring-1 focus:ring-(--color-blue) focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
        </div>
      </div>
    </div>

    <!-- Profile Image -->
    <div class="border-t border-(--color-border) pt-6">
      <h3 class="text-sm font-semibold text-(--color-text-secondary) uppercase tracking-wider mb-4">
        Profile Image <span class="normal-case text-xs font-normal">(Optional)</span>
      </h3>
      <div class="flex flex-col sm:flex-row items-start gap-4">
        <div class="shrink-0">
          <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-(--color-border)">
            <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
            <button type="button" @click="removeImage"
              class="absolute top-1 right-1 w-5 h-5 bg-(--color-red) text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">৳—</button>
          </div>
          <div v-else
            class="w-24 h-24 rounded-lg border-2 border-dashed border-(--color-border) flex items-center justify-center bg-(--color-muted-bg)">
            <svg class="w-10 h-10 text-(--color-text-secondary)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
        </div>

        <div class="flex-1 space-y-3">
          <div>
            <label class="text-sm font-medium text-(--color-text-primary) block mb-1.5">Upload
              Image</label>
            <div class="flex flex-wrap gap-2">
              <label
                class="px-4 py-2 text-sm font-medium rounded-lg cursor-pointer bg-(--color-muted-bg) border border-(--color-border) hover:bg-(--color-muted-bg)/70 transition-all duration-200 inline-flex items-center gap-2"
                :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
                {{ uploading ? 'Uploading...' : 'Choose File' }}
                <input type="file" accept="image/*" class="hidden" @change="handleFileSelect" :disabled="uploading" />
              </label>
              <button v-if="imagePreview" type="button" @click="removeImage"
                class="px-4 py-2 text-sm font-medium rounded-lg text-(--color-red) hover:bg-(--color-red)/10 transition-all duration-200">
                Remove
              </button>
            </div>
          </div>

          <div v-if="uploading" class="w-full bg-(--color-muted-bg) rounded-full h-1.5 overflow-hidden">
            <div class="bg-(--color-blue) h-full rounded-full transition-all duration-300"
              :style="{ width: uploadProgress + '%' }"></div>
          </div>

          <p v-if="uploadError" class="text-xs text-(--color-red)">{{ uploadError }}</p>
          <p v-else-if="justCreatedId" class="text-xs text-(--color-yellow)">
            User was created but the image upload failed. Retry below or remove the image.
          </p>
          <p v-else-if="isEditMode && props.user?.image_url" class="text-xs text-(--color-text-secondary)">
            Current image will be replaced
          </p>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex flex-col sm:flex-row items-center justify-end gap-3 pt-4 border-t border-(--color-border)">
      <button type="button" @click="emit('cancel')" :disabled="submitting || uploading"
        class="w-full sm:w-auto px-4 py-2 text-sm font-medium rounded-lg hover:bg-(--color-muted-bg) transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed">
        Cancel
      </button>
      <button type="submit" :disabled="submitting || uploading"
        class="w-full sm:w-auto px-6 py-2 text-sm font-semibold bg-(--color-blue) text-white rounded-lg hover:opacity-90 transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
        <span v-if="submitting || uploading" class="inline-flex items-center justify-center gap-2">
          <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          {{ justCreatedId ? 'Retrying image...' : (isEditMode ? 'Saving...' : 'Creating...') }}
        </span>
        <span v-else>{{ justCreatedId ? 'Retry Image Upload' : (isEditMode ? 'Save Changes' : 'Create User') }}</span>
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { IDType, Role, User } from '@/types/auth'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const props = defineProps<{
  user?: User | null
  mode?: 'create' | 'edit'
}>()

const emit = defineEmits<{
  'user-created': []
  'user-updated': []
  'cancel': []
}>()

const auth = useAuthStore()
const usersStore = useUsersStore()

const isEditMode = computed(() => props.mode === 'edit' || !!props.user)
const canCreateManager = computed(() => auth.user?.role === 'admin')

const form = ref({
  name: '',
  username: '',
  password: '',
  email: '',
  phone: '',
  address: '',
  role: 'staff' as Role,
  id_type: '' as IDType | '',
  id_number: '',
  monthly_salary: null as number | null,
})

const showPassword = ref(false)
const submitting = ref(false)

const justCreatedId = ref<number | null>(null)

// IMAGE STATE + HELPERS
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)
const imageToDelete = ref(false)

const clearImageState = () => {
  imageFile.value = null
  imagePreview.value = null
  uploadError.value = null
  uploadProgress.value = 0
  imageToDelete.value = false
  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null
  if (fileInput) fileInput.value = ''
}

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Please select an image file'
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = 'Image size should be less than 5MB'
    return
  }

  uploadError.value = null
  imageFile.value = file
  imageToDelete.value = false

  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

const removeImage = () => {
  imageFile.value = null
  imagePreview.value = null
  uploadError.value = null
  uploadProgress.value = 0
  imageToDelete.value = true
  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null
  if (fileInput) fileInput.value = ''
}

const initializeForm = () => {
  if (props.user) {
    form.value = {
      name: props.user.name || '',
      username: props.user.username || '',
      password: '',
      email: props.user.email || '',
      phone: props.user.phone || '',
      address: props.user.address || '',
      role: props.user.role || 'staff',
      id_type: props.user.id_type || '',
      id_number: props.user.id_number || '',
      monthly_salary: props.user.monthly_salary || null,
    }
    justCreatedId.value = null
    clearImageState()
    if (props.user.image_url) {
      imagePreview.value = getImageUrl(props.user.image_url) || null
    }
  } else {
    justCreatedId.value = null
    clearImageState()
  }
}

watch(() => props.user, initializeForm, { immediate: true })

const resetForm = () => {
  if (isEditMode.value && props.user) {
    initializeForm()
  } else {
    form.value = {
      name: '',
      username: '',
      password: '',
      email: '',
      phone: '',
      address: '',
      role: 'staff',
      id_type: '',
      id_number: '',
      monthly_salary: null,
    }
    showPassword.value = false
    justCreatedId.value = null
    clearImageState()
  }
}

const submit = async () => {
  // RETRY IMAGE PATH
  if (justCreatedId.value !== null) {
    if (!imageFile.value) {
      resetForm()
      emit('user-created')
      return
    }
    uploading.value = true
    uploadError.value = null
    uploadProgress.value = 0
    try {
      const url = await uploadImage('users', justCreatedId.value, imageFile.value, (p) => {
        uploadProgress.value = p
      })
      if (!url) {
        uploadError.value = 'Failed to upload image. Please try again.'
        return
      }
      push.success('User and image created successfully!')
      resetForm()
      emit('user-created')
    } finally {
      uploading.value = false
    }
    return
  }

  // NORMAL VALIDATION
  if (!form.value.name) {
    push.error('Name is required')
    return
  }
  if (!form.value.username) {
    push.error('Username is required')
    return
  }
  if (!isEditMode.value && !form.value.password) {
    push.error('Password is required')
    return
  }
  if (form.value.password && form.value.password.length < 8) {
    push.error('Password must be at least 8 characters')
    return
  }

  submitting.value = true

  try {
    let userId: number

    if (isEditMode.value && props.user) {
      const profileSuccess = await usersStore.updateUserProfile(props.user.id, {
        name: form.value.name,
        email: form.value.email || null,
        phone: form.value.phone || null,
        address: form.value.address || null,
        id_type: form.value.id_type as IDType || null,
        id_number: form.value.id_number || null,
      })

      if (!profileSuccess) {
        submitting.value = false
        return
      }

      if (form.value.monthly_salary !== props.user.monthly_salary) {
        await usersStore.updateUserSalary(props.user.id, {
          monthly_salary: form.value.monthly_salary || 0
        })
      }

      if (form.value.password) {
        await usersStore.changeUserPassword(props.user.id, form.value.password)
      }

      userId = props.user.id
      push.success('User updated successfully!')
    } else {
      const newUser = await usersStore.createUser({
        name: form.value.name,
        username: form.value.username,
        password: form.value.password,
        role: form.value.role,
        email: form.value.email || null,
        phone: form.value.phone || null,
        address: form.value.address || null,
        id_type: form.value.id_type || null,
        id_number: form.value.id_number || null,
        monthly_salary: form.value.monthly_salary || 0,
        image_url: null,
      })

      if (!newUser) {
        submitting.value = false
        return
      }

      userId = newUser.id
      push.success('User created successfully!')
    }

    // IMAGE HANDLING
    if (imageFile.value) {
      uploading.value = true
      uploadProgress.value = 0
      uploadError.value = null
      const url = await uploadImage('users', userId, imageFile.value, (p) => {
        uploadProgress.value = p
      })
      uploading.value = false
      if (!url) {
        uploadError.value = 'Failed to upload image. Please try again.'
        if (!isEditMode.value) {
          justCreatedId.value = userId
        }
        push.warning('User saved but image upload failed. Retry below.')
        await usersStore.fetchUsers()
        return
      }
    } else if (isEditMode.value && props.user && imageToDelete.value && props.user.image_url) {
      await deleteImage('users', props.user.id, props.user.image_url)
    }

    if (isEditMode.value) {
      emit('user-updated')
    } else {
      resetForm()
      emit('user-created')
    }
  } catch (error) {
    console.error('Error:', error)
    push.error(isEditMode.value ? 'Failed to update user' : 'Failed to create user')
  } finally {
    submitting.value = false
  }
}
</script>