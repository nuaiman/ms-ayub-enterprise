<template>
  <div class="min-h-screen bg-app flex items-center justify-center p-6">
    <div class="text-center">
      <div
        class="w-24 h-24 mx-auto rounded-full bg-surface-alt border border-default flex items-center justify-center mb-6">
        <svg class="w-12 h-12 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
            d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h1 class="text-6xl font-bold text-primary mb-2">404</h1>
      <p class="text-lg font-medium text-primary mb-2">Page not found</p>
      <p class="text-sm text-muted mb-8">The page you're looking for doesn't exist.</p>
      <router-link to="/dashboard"
        class="inline-flex px-6 py-2.5 bg-accent hover:bg-accent-hover text-accent-foreground rounded-lg font-semibold transition-colors">
        Go to Dashboard
      </router-link>
    </div>
  </div>
</template>