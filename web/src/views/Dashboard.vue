<template>
  <div class="min-h-screen bg-notion-bg">
    <NavBar />

    <!-- Main Content -->
    <main class="w-full py-8 px-6 lg:px-8">
      <!-- Page Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-notion-text">Dashboard</h1>
        <p class="text-sm text-notion-textLight mt-2">Monitor your S3-compatible object storage</p>
      </div>

      <!-- Content Area -->
      <div class="bg-white rounded-lg border border-notion-border p-6">
        <p class="text-notion-textLight py-8">Dashboard content will be added here</p>
      </div>
    </main>
  </div>
</template>

<script setup>
import NavBar from '@/components/NavBar.vue'
</script>

