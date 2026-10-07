<template>
  <div class="min-h-screen bg-(--color-bg) flex items-center justify-center p-6">
    <div class="text-center">
      <div
        class="w-24 h-24 mx-auto rounded-full bg-(--color-muted-bg) border border-(--color-border) flex items-center justify-center mb-6">
        <svg class="w-12 h-12 text-(--color-text-secondary)" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
            d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h1 class="text-6xl font-bold text-(--color-text-primary) mb-2">404</h1>
      <p class="text-lg font-medium text-(--color-text-primary) mb-2">Page not found</p>
      <p class="text-sm text-(--color-text-secondary) mb-8">The page you're looking for doesn't exist.</p>
      <router-link to="/customers"
        class="inline-flex px-6 py-2.5 bg-(--color-blue) hover:opacity-90 text-white rounded-lg font-semibold transition-opacity">
        Go to Dashboard
      </router-link>
    </div>
  </div>
</template>