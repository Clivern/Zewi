<template>
  <div class="min-h-screen bg-theme-bg flex items-center justify-center px-4">
    <div class="max-w-2xl w-full text-center">
      <!-- 500 Content -->
      <div class="mb-8">
        <h1 class="text-9xl font-bold text-theme-text mb-4">500</h1>
        <h2 class="text-3xl font-semibold text-theme-text mb-4">Internal Server Error</h2>
        <p class="text-lg text-theme-textLight mb-4">
          We're sorry, but something went wrong on our end.
        </p>
      </div>

      <!-- Actions -->
      <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
        <router-link to="/" class="btn-primary">
          <svg class="w-5 h-5 inline-block mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Home
        </router-link>
      </div>

      <!-- Support Info -->
      <div class="mt-12 pt-8 border-t border-theme-border">
        <p class="text-sm text-theme-textLight mb-2">
          If this problem persists, please contact our support team.
        </p>
        <a href="mailto:support@ziee.sh" class="text-sm text-theme-text hover:text-theme-textLight transition-colors">
          support@ziee.sh
        </a>
      </div>
    </div>
  </div>
</template>

