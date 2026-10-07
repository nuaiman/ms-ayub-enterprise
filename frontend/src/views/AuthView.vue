<template>
  <main class="min-h-screen w-full flex items-center justify-center p-4 sm:p-6
             bg-(--color-bg) text-(--color-text-primary)">
    <div class="w-full max-w-md space-y-6">

      <!-- Main Card Container -->
      <div class="p-7 sm:p-9 rounded-3xl
                  bg-(--color-surface)
                  border border-(--color-border)
                  shadow-2xl shadow-black/5 space-y-7">

        <!-- Brand Header -->
        <header class="flex items-center gap-4">
          <div class="w-13 h-13 shrink-0 overflow-hidden rounded-2xl
                      bg-(--color-muted-bg)
                      border border-(--color-border) p-0.5">
            <img src="@/assets/logo.png" alt="MS Ayub Enterprise"
              class="w-full h-full block object-cover rounded-[14px]" />
          </div>

          <div class="min-w-0">
            <h1 class="text-lg font-semibold tracking-tight text-(--color-text-primary)">
              MS Ayub Enterprise
            </h1>

            <p class="text-xs text-(--color-text-secondary)">
              Sign in to access your workspace
            </p>
          </div>
        </header>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">

          <!-- Username -->
          <div class="space-y-1.5">
            <label for="username" class="block text-[11px] font-semibold uppercase tracking-wider
                         text-(--color-text-secondary)">
              Username
            </label>

            <div class="relative flex items-center">
              <span class="absolute left-3.5 text-(--color-icon-secondary) pointer-events-none">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </span>

              <input id="username" v-model="username" type="text" autocomplete="username"
                placeholder="Enter your username" :disabled="loading" class="w-full h-11 pl-10 pr-4 rounded-xl text-sm
                       bg-(--color-muted-bg)
                       border border-(--color-border)
                       text-(--color-text-primary)
                       placeholder:text-(--color-text-secondary)
                       placeholder:opacity-60
                       outline-none
                       transition
                       hover:border-(--color-chevron)
                       focus:bg-(--color-surface)
                       focus:border-(--color-blue)
                       focus:ring-2
                       focus:ring-[color-mix(in_srgb,var(--color-blue)_15%,transparent)]
                       disabled:opacity-60
                       disabled:cursor-not-allowed" />
            </div>
          </div>

          <!-- Password -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <label for="password" class="block text-[11px] font-semibold uppercase tracking-wider
                           text-(--color-text-secondary)">
                Password
              </label>

              <button type="button" class="text-xs font-medium
                           text-(--color-blue)
                           transition-opacity
                           hover:opacity-75" @click="showForgotPassword = !showForgotPassword">
                Forgot password?
              </button>
            </div>

            <div class="relative flex items-center">
              <span class="absolute left-3.5 text-(--color-icon-secondary) pointer-events-none">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </span>

              <input id="password" v-model="password" :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password" placeholder="Enter your password" :disabled="loading" class="w-full h-11 pl-10 pr-10 rounded-xl text-sm
                       bg-(--color-muted-bg)
                       border border-(--color-border)
                       text-(--color-text-primary)
                       placeholder:text-(--color-text-secondary)
                       placeholder:opacity-60
                       outline-none
                       transition
                       hover:border-(--color-chevron)
                       focus:bg-(--color-surface)
                       focus:border-(--color-blue)
                       focus:ring-2
                       focus:ring-[color-mix(in_srgb,var(--color-blue)_15%,transparent)]
                       disabled:opacity-60
                       disabled:cursor-not-allowed" />

              <button type="button" :disabled="loading" :aria-label="showPassword ? 'Hide password' : 'Show password'"
                class="absolute right-2.5 flex items-center justify-center
                       w-7 h-7 rounded-lg
                       text-(--color-icon-secondary)
                       bg-transparent
                       transition-colors
                       hover:text-(--color-icon-primary)
                       hover:bg-(--color-surface-alt)
                       active:scale-95
                       disabled:opacity-50
                       disabled:cursor-not-allowed" @click="showPassword = !showPassword">

                <svg v-if="!showPassword" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"
                  aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  <circle cx="12" cy="12" r="3" stroke-width="1.8" />
                </svg>

                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M3 3l18 18" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8"
                    d="M10.584 10.587A2 2 0 0013.414 13.4" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8"
                    d="M9.88 4.24A10.72 10.72 0 0112 4c4.477 0 8.268 2.943 9.542 7a10.77 10.77 0 01-4.126 5.407" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8"
                    d="M6.228 6.228C4.489 7.535 3.207 9.421 2.458 12 3.732 16.057 7.523 19 12 19c1.61 0 3.13-.37 4.48-1.03" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Forgot Password Notice -->
          <Transition enter-active-class="transition-all duration-150 ease-out"
            leave-active-class="transition-all duration-150 ease-in" enter-from-class="opacity-0 -translate-y-1"
            leave-to-class="opacity-0 -translate-y-1">
            <div v-if="showForgotPassword" class="rounded-xl px-3.5 py-2.5 text-xs leading-relaxed
                       bg-(--color-card-secondary)
                       border border-(--color-border)
                       text-(--color-text-secondary)">
              Password resets are managed by your administrator. Please contact technical support.
            </div>
          </Transition>

          <!-- Login Error -->
          <Transition enter-active-class="transition-all duration-150 ease-out"
            leave-active-class="transition-all duration-150 ease-in" enter-from-class="opacity-0 -translate-y-1"
            leave-to-class="opacity-0 -translate-y-1">
            <div v-if="error" role="alert" class="rounded-xl px-3.5 py-2.5 text-xs font-medium
                       bg-[color-mix(in_srgb,var(--color-red)_8%,var(--color-card-secondary))]
                       border border-[color-mix(in_srgb,var(--color-red)_25%,var(--color-border))]
                       text-(--color-red)">
              {{ error }}
            </div>
          </Transition>

          <!-- Submit -->
          <button type="submit" :disabled="loading || !username.trim() || !password" class="w-full h-11 rounded-xl px-5 mt-2
                     text-sm font-semibold
                     flex items-center justify-center gap-2
                     bg-(--color-button-primary)
                     text-(--color-on-button-primary)
                     border border-(--color-button-primary)
                     transition-opacity
                     hover:opacity-92
                     active:scale-[0.985]
                     disabled:opacity-45
                     disabled:cursor-not-allowed">
            <svg v-if="loading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2.5" class="opacity-25" />
              <path fill="currentColor" d="M21 12a9 9 0 00-9-9v2a7 7 0 017 7h2z" />
            </svg>

            <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
          </button>
        </form>

        <!-- Support Box Inside Card Footer -->
        <div class="pt-4 border-t border-(--color-border) text-center text-xs text-(--color-text-secondary)">
          Having trouble logging in?
          <router-link to="/support" class="ml-1 font-medium text-(--color-blue) transition-opacity hover:opacity-75">
            Contact Support
          </router-link>
        </div>

      </div>

      <!-- Footer -->
      <footer class="text-center text-xs text-(--color-text-secondary) space-y-1">
        <p>&copy; {{ currentYear }} MS Ayub Enterprise. All rights reserved.</p>
        <p>
          Powered by
          <span class="font-medium text-(--color-text-primary)">Boiddutik</span>
        </p>
      </footer>

    </div>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const showForgotPassword = ref(false)
const loading = ref(false)
const error = ref('')

const currentYear = new Date().getFullYear()

const handleLogin = async () => {
  error.value = ''
  const trimmedUsername = username.value.trim()

  if (!trimmedUsername || !password.value) {
    error.value = 'Please enter both username and password.'
    return
  }

  loading.value = true

  try {
    const success = await auth.login(trimmedUsername, password.value)

    if (success) {
      await router.push('/customers')
      return
    }

    error.value = 'Invalid username or password.'
  } catch (err) {
    error.value = 'Unable to sign in. Please try again.'
    console.error('[AuthView] Login error:', err)
  } finally {
    loading.value = false
  }
}
</script>
