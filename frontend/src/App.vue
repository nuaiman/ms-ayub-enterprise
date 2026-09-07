<!-- src/App.vue -->
<template>
    <div class="min-h-screen bg-app text-primary transition-theme">
        <!-- Show loading spinner while auth is initializing -->
        <div v-if="auth.isLoading" class="fixed inset-0 flex items-center justify-center bg-app z-50">
            <div class="flex flex-col items-center gap-4">
                <div class="w-12 h-12 border-4 border-accent border-t-transparent rounded-full animate-spin"></div>
                <p class="text-sm text-secondary">Loading...</p>
            </div>
        </div>
        <RouterView v-else />

        <!-- Notivue Notifications -->
        <Notivue v-slot="item">
            <Notification :item="item" />
        </Notivue>
    </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Notivue, Notification } from 'notivue'

const auth = useAuthStore()
</script>