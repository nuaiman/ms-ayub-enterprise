// src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api, { setAccessToken } from '@/utils/axios'
import type { User } from '@/types/auth'
import type { ApiResponse } from '@/types/api'
import { push } from 'notivue'
import { useGlobalLoader } from 'vue-global-loader'
import type { AxiosError } from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const { displayLoader, destroyLoader } = useGlobalLoader()

  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('access_token'))
  const isLoading = ref(false)

  const isAuthenticated = computed(() => !!user.value && !!token.value)

  const setSession = (newToken: string, newUser: User) => {
    console.log('[AUTH] Setting session:', { newToken, newUser })
    token.value = newToken
    user.value = newUser
    setAccessToken(newToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
    localStorage.setItem('access_token', newToken)
  }

  const updateToken = (newToken: string) => {
    token.value = newToken
    setAccessToken(newToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  const clearSession = () => {
    token.value = null
    user.value = null
    setAccessToken(null)
    delete api.defaults.headers.common['Authorization']
    localStorage.removeItem('access_token')
    document.cookie = 'refresh_token=; Max-Age=0; path=/'
  }

  const login = async (username: string, password: string): Promise<boolean> => {
    displayLoader()
    try {
      console.log('[AUTH] Login attempt:', username)

      const res = await api.post<ApiResponse<{ access_token: string }>>('/users/login', {
        username,
        password,
      })

      console.log('[AUTH] Login response:', res.data)

      if (!res.data.success) {
        push.error(res.data.message || 'Login failed')
        return false
      }

      const newToken = res.data.data.access_token
      console.log('[AUTH] Token received:', newToken)

      // Set token before making the next request
      setAccessToken(newToken)
      api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`

      console.log('[AUTH] Fetching current user...')
      const me = await api.get<ApiResponse<User>>('/users/current-user')

      console.log('[AUTH] Current user response:', me.data)

      if (!me.data.success) {
        push.error('Failed to get user info')
        return false
      }

      const userData = me.data.data
      console.log('[AUTH] User data:', userData)

      setSession(newToken, userData)
      push.success('Welcome back!')
      return true
    } catch (error) {
      console.error('[AUTH] Login error:', error)
      const err = error as AxiosError<{ message: string }>
      const message = err.response?.data?.message || 'Login failed. Please check your credentials.'
      push.error(message)
      clearSession()
      return false
    } finally {
      destroyLoader()
    }
  }

  const logout = async () => {
    displayLoader()
    try {
      await api.delete('/users/logout')
    } catch {
      // Ignore errors on logout
    } finally {
      clearSession()
      push.info('Logged out successfully')
      destroyLoader()
    }
  }

  const initAuth = async (): Promise<boolean> => {
    isLoading.value = true
    try {
      const storedToken = localStorage.getItem('access_token')
      console.log('[AUTH] Init auth, stored token:', storedToken ? 'exists' : 'none')

      if (!storedToken) {
        console.log('[AUTH] No stored token found')
        return false
      }

      setAccessToken(storedToken)
      api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`

      console.log('[AUTH] Validating token with /users/current-user...')
      const me = await api.get<ApiResponse<User>>('/users/current-user')
      console.log('[AUTH] Init auth user response:', me.data)

      if (me.data.success) {
        console.log('[AUTH] Token valid, restoring session for user:', me.data.data.username)
        setSession(storedToken, me.data.data)
        return true
      }

      console.log('[AUTH] Token invalid or expired, clearing session')
      clearSession()
      return false
    } catch (error) {
      console.error('[AUTH] Init auth error:', error)
      // Only clear if it's an auth error (401), not network errors
      const err = error as AxiosError
      if (err.response?.status === 401) {
        clearSession()
      }
      return false
    } finally {
      isLoading.value = false
    }
  }

  const changePassword = async (currentPassword: string, newPassword: string): Promise<boolean> => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<null>>('/users/change-password', {
        current_password: currentPassword,
        new_password: newPassword,
      })

      if (!res.data.success) {
        push.error(res.data.message || 'Password change failed')
        return false
      }

      push.success(res.data.message || 'Password changed successfully')
      return true
    } catch (error) {
      const err = error as AxiosError<{ message: string }>
      push.error(err.response?.data?.message || 'Password change failed')
      return false
    } finally {
      destroyLoader()
    }
  }

  const updateProfile = async (payload: {
    name: string
    email: string | null
    phone: string | null
    address: string | null
    id_type: string | null
    id_number: string | null
  }): Promise<boolean> => {
    displayLoader()
    try {
      const res = await api.patch<ApiResponse<User>>(`/users/${user.value?.id}/profile`, payload)

      if (!res.data.success) {
        push.error(res.data.message || 'Profile update failed')
        return false
      }

      user.value = res.data.data
      push.success(res.data.message || 'Profile updated successfully')
      return true
    } catch (error) {
      const err = error as AxiosError<{ message: string }>
      push.error(err.response?.data?.message || 'Profile update failed')
      return false
    } finally {
      destroyLoader()
    }
  }

  const refreshSession = async (): Promise<boolean> => {
    displayLoader()
    try {
      const res = await api.post<ApiResponse<{ access_token: string }>>('/users/refresh-session')

      if (!res.data.success) {
        clearSession()
        return false
      }

      const newToken = res.data.data.access_token
      updateToken(newToken)

      const me = await api.get<ApiResponse<User>>('/users/current-user')

      if (me.data.success) {
        user.value = me.data.data
        return true
      }

      clearSession()
      return false
    } catch {
      clearSession()
      return false
    } finally {
      destroyLoader()
    }
  }

  return {
    // State
    user,
    token,
    isLoading,

    // Computed
    isAuthenticated,

    // Actions
    login,
    logout,
    initAuth,
    changePassword,
    updateProfile,
    refreshSession,
    updateToken,
    clearSession,
  }
})