<!-- src/views/SupportView.vue -->
<template>
  <div class="max-w-6xl mx-auto">
    <!-- Header -->
    <div class="text-center mb-12">
      <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-accent/10 border border-accent/20 mb-5">
        <span class="relative flex h-2 w-2">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-accent"></span>
        </span>
        <span class="text-xs font-medium text-accent">We're here to help</span>
      </div>
      <h1 class="text-4xl md:text-6xl font-bold text-primary mb-4 tracking-tight">
        Need <span class="text-accent">Support</span>?
      </h1>
      <p class="text-secondary text-lg max-w-2xl mx-auto">Get in touch with the developer or find answers to your
        questions</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Developer Card -->
      <div
        class="group relative overflow-hidden rounded-3xl bg-surface border border-default shadow-sm hover:shadow-xl transition-all duration-500 hover:-translate-y-1">
        <div
          class="absolute top-0 right-0 w-80 h-80 bg-accent/5 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2 group-hover:bg-accent/10 transition-all duration-700">
        </div>

        <div class="relative p-6 md:p-8">
          <div class="flex items-start gap-5 mb-6">
            <div class="relative">
              <div
                class="w-20 h-20 rounded-2xl bg-linear-to-br from-accent to-accent/70 flex items-center justify-center shadow-lg shadow-accent/25">
                <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                </svg>
              </div>
              <div
                class="absolute -bottom-1 -right-1 w-6 h-6 rounded-full bg-success-text border-2 border-surface flex items-center justify-center">
                <svg class="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
              </div>
            </div>
            <div>
              <h2 class="text-2xl font-bold text-primary">Developer</h2>
              <p class="text-sm text-muted">Software Developer</p>
              <div class="flex items-center gap-1 mt-1">
                <span class="inline-flex w-1.5 h-1.5 rounded-full bg-success-text animate-pulse"></span>
                <span class="text-xs text-muted">Available for support</span>
              </div>
            </div>
          </div>

          <div class="space-y-3.5">
            <!-- Email -->
            <div
              class="flex items-center gap-4 p-4 rounded-2xl bg-surface-alt/50 border border-default/50 hover:border-accent/30 hover:bg-surface-alt/70 transition-all duration-300 group/email">
              <div
                class="w-11 h-11 rounded-xl bg-linear-to-br from-green-500/15 to-green-500/5 flex items-center justify-center shrink-0 border border-green-500/10">
                <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-[10px] font-semibold text-muted uppercase tracking-wider">Email</p>
                <a :href="'mailto:' + developer.email"
                  class="text-sm text-primary hover:text-accent transition-colors truncate block font-medium">
                  {{ developer.email }}
                </a>
              </div>
              <button @click="copyEmail"
                class="p-2.5 hover:bg-surface rounded-xl transition-all duration-200 group/btn shrink-0">
                <svg class="w-4 h-4 text-muted group-hover/btn:text-accent transition-colors" fill="none"
                  stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              </button>
            </div>

            <!-- YouTube -->
            <div
              class="flex items-center gap-4 p-4 rounded-2xl bg-surface-alt/50 border border-default/50 hover:border-accent/30 hover:bg-surface-alt/70 transition-all duration-300">
              <div
                class="w-11 h-11 rounded-xl bg-linear-to-br from-red-500/15 to-red-500/5 flex items-center justify-center shrink-0 border border-red-500/10">
                <svg class="w-5 h-5 text-red-600" fill="currentColor" viewBox="0 0 24 24">
                  <path
                    d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" />
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-[10px] font-semibold text-muted uppercase tracking-wider">YouTube</p>
                <a :href="developer.youtube" target="_blank" rel="noopener noreferrer"
                  class="text-sm text-primary hover:text-accent transition-colors truncate block font-medium">
                  {{ developer.youtube.replace('https://www.youtube.com/@', '@') }}
                </a>
              </div>
              <div class="w-4 h-4 text-muted">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
              </div>
            </div>

            <!-- Status Badge -->
            <div class="flex items-center gap-3 p-3 rounded-xl bg-info-bg/30 border border-info-border/30">
              <svg class="w-4 h-4 text-info-text shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="text-xs text-info-text">Usually responds within 24 hours</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Support Card -->
      <div
        class="group relative overflow-hidden rounded-3xl bg-surface border border-default shadow-sm hover:shadow-xl transition-all duration-500 hover:-translate-y-1">
        <div
          class="absolute bottom-0 left-0 w-80 h-80 bg-accent/5 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2 group-hover:bg-accent/10 transition-all duration-700">
        </div>

        <div class="relative p-6 md:p-8">
          <div class="flex items-start gap-5 mb-6">
            <div class="relative">
              <div
                class="w-20 h-20 rounded-2xl bg-linear-to-br from-warning-text to-warning-text/70 flex items-center justify-center shadow-lg shadow-warning-text/25">
                <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M18.364 5.636L9.172 14.828a4 4 0 01-5.657 0L2.879 14.12m15.485 0l-5.657 5.657a4 4 0 01-5.657 0L5.636 18.364m-1.414 1.414L2.879 18.364" />
                </svg>
              </div>
            </div>
            <div>
              <h2 class="text-2xl font-bold text-primary">Support</h2>
              <p class="text-sm text-muted">How can we help?</p>
              <div class="flex items-center gap-1 mt-1">
                <span class="inline-flex w-1.5 h-1.5 rounded-full bg-accent"></span>
                <span class="text-xs text-muted">We're here for you</span>
              </div>
            </div>
          </div>

          <div class="space-y-3.5">
            <div
              class="p-4 rounded-2xl bg-surface-alt/50 border border-default/50 hover:border-accent/30 transition-all duration-300">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-success-bg/50 flex items-center justify-center shrink-0">
                  <svg class="w-4 h-4 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div>
                  <p class="text-[10px] font-semibold text-muted uppercase tracking-wider">Response Time</p>
                  <p class="text-sm text-primary font-medium flex items-center gap-2">
                    <span class="w-1.5 h-1.5 rounded-full bg-success-text animate-pulse"></span>
                    Usually within 24 hours
                  </p>
                </div>
              </div>
            </div>

            <div
              class="p-4 rounded-2xl bg-surface-alt/50 border border-default/50 hover:border-accent/30 transition-all duration-300">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-info-bg/50 flex items-center justify-center shrink-0">
                  <svg class="w-4 h-4 text-info-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <div>
                  <p class="text-[10px] font-semibold text-muted uppercase tracking-wider">Best Way to Reach</p>
                  <p class="text-sm text-primary font-medium">📧 Email us with your query</p>
                </div>
              </div>
            </div>

            <div
              class="p-4 rounded-2xl bg-surface-alt/50 border border-default/50 hover:border-accent/30 transition-all duration-300">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-warning-bg/50 flex items-center justify-center shrink-0">
                  <svg class="w-4 h-4 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                </div>
                <div>
                  <p class="text-[10px] font-semibold text-muted uppercase tracking-wider">Support Hours</p>
                  <p class="text-sm text-primary font-medium">Mon-Fri, 9AM - 6PM (GMT+6)</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div
      class="mt-8 relative overflow-hidden rounded-3xl bg-surface border border-default shadow-sm hover:shadow-xl transition-all duration-500">
      <div
        class="absolute top-0 left-1/2 w-96 h-96 bg-accent/5 rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2">
      </div>

      <div class="relative p-6 md:p-8">
        <div class="flex items-center gap-4 mb-6">
          <div class="w-12 h-12 rounded-2xl bg-accent/10 flex items-center justify-center shrink-0">
            <svg class="w-6 h-6 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-primary">Quick Actions</h3>
            <p class="text-sm text-muted">Get help instantly with these quick actions</p>
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <button @click="copyEmail"
            class="group relative flex items-center gap-4 p-5 rounded-2xl border border-default/50 bg-surface-alt/30 hover:border-accent/50 hover:bg-accent/5 transition-all duration-300 overflow-hidden">
            <div
              class="absolute inset-0 bg-linear-to-r from-accent/0 via-accent/0 to-accent/0 group-hover:from-accent/5 group-hover:via-accent/0 group-hover:to-accent/5 transition-all duration-500">
            </div>
            <div
              class="w-14 h-14 rounded-2xl bg-accent/10 flex items-center justify-center shrink-0 group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
              </svg>
            </div>
            <div class="text-left flex-1">
              <p class="text-sm font-semibold text-primary group-hover:text-accent transition-colors">Copy Email</p>
              <p class="text-xs text-muted">Copy developer email to clipboard</p>
            </div>
            <svg
              class="w-5 h-5 text-muted group-hover:text-accent group-hover:translate-x-1 transition-all duration-300"
              fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <a :href="developer.youtube" target="_blank" rel="noopener noreferrer"
            class="group relative flex items-center gap-4 p-5 rounded-2xl border border-default/50 bg-surface-alt/30 hover:border-accent/50 hover:bg-accent/5 transition-all duration-300 overflow-hidden">
            <div
              class="absolute inset-0 bg-linear-to-r from-accent/0 via-accent/0 to-accent/0 group-hover:from-accent/5 group-hover:via-accent/0 group-hover:to-accent/5 transition-all duration-500">
            </div>
            <div
              class="w-14 h-14 rounded-2xl bg-red-500/10 flex items-center justify-center shrink-0 group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6 text-red-600" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" />
              </svg>
            </div>
            <div class="text-left flex-1">
              <p class="text-sm font-semibold text-primary group-hover:text-accent transition-colors">Visit YouTube</p>
              <p class="text-xs text-muted">Watch tutorials and updates</p>
            </div>
            <svg
              class="w-5 h-5 text-muted group-hover:text-accent group-hover:translate-x-1 transition-all duration-300"
              fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </a>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="mt-10 text-center">
      <p class="text-sm text-muted">
        Made with <span class="text-accent">❤️</span> by <span class="text-accent font-medium">Boiddutik</span>
      </p>
      <p class="text-xs text-muted mt-1">© {{ new Date().getFullYear() }} MS Ayub Enterprise. All rights reserved.</p>
    </div>
  </div>

  <!-- Toast notification -->
  <Transition enter-active-class="transition ease-out duration-300" enter-from-class="opacity-0 translate-y-4 scale-95"
    enter-to-class="opacity-100 translate-y-0 scale-100" leave-active-class="transition ease-in duration-200"
    leave-from-class="opacity-100 translate-y-0 scale-100" leave-to-class="opacity-0 translate-y-4 scale-95">
    <div v-if="showToast"
      class="fixed bottom-6 right-6 bg-surface border border-success-border rounded-2xl shadow-2xl p-4 max-w-sm z-50">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-success-bg flex items-center justify-center shrink-0">
          <svg class="w-5 h-5 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <div>
          <p class="text-sm font-semibold text-primary">Copied!</p>
          <p class="text-xs text-muted">Email copied to clipboard</p>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useClipboardStore } from '@/stores/clipboard'

const clipboardStore = useClipboardStore()

const developer = {
  name: 'Boiddutik',
  email: 'boiddutik@gmail.com',
  youtube: 'https://www.youtube.com/@boiddutik'
}

const showToast = ref(false)
let toastTimeout: ReturnType<typeof setTimeout> | null = null

const copyEmail = async () => {
  const success = await clipboardStore.copyToClipboard(developer.email)

  if (success) {
    showToast.value = true
    if (toastTimeout) clearTimeout(toastTimeout)
    toastTimeout = setTimeout(() => {
      showToast.value = false
    }, 2500)
  }
}
</script>

<style scoped>
.animate-slide-up {
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.animate-ping {
  animation: ping 2s cubic-bezier(0, 0, 0.2, 1) infinite;
}

@keyframes ping {

  75%,
  100% {
    transform: scale(2);
    opacity: 0;
  }
}
</style>