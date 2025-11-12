<template>
  <div class="min-h-screen bg-notion-bg">
    <NavBar />

    <!-- Main Content -->
    <main class="w-full py-8 px-6 lg:px-8">
      <!-- Page Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-notion-text">Settings</h1>
        <p class="text-sm text-notion-textLight mt-2">Manage application settings and preferences</p>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="mb-6 rounded-md border border-green-200 bg-green-50 p-4">
        <p class="text-sm text-green-800">{{ successMessage }}</p>
      </div>
      <div v-if="errorMessage" class="mb-6 rounded-md border border-red-200 bg-red-50 p-4">
        <p class="text-sm text-red-800">{{ errorMessage }}</p>
      </div>

      <!-- Settings Form -->
      <form @submit.prevent="handleSave" class="space-y-6">
        <!-- General Settings Box -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-4">General</h2>
          <div class="border-b border-notion-border mb-4"></div>
          <div class="space-y-4">
            <div>
              <label for="applicationName" class="block text-sm font-medium text-notion-text mb-2">
                Application Name
              </label>
              <input
                id="applicationName"
                v-model="form.applicationName"
                type="text"
                required
                class="input-field max-w-md"
                placeholder="Zewi"
                :disabled="loading"
              >
            </div>
            <div>
              <label for="applicationEmail" class="block text-sm font-medium text-notion-text mb-2">
                Application Email
              </label>
              <input
                id="applicationEmail"
                v-model="form.applicationEmail"
                type="email"
                required
                class="input-field max-w-md"
                placeholder="admin@zewi.com"
                :disabled="loading"
              >
            </div>
            <div>
              <label for="applicationURL" class="block text-sm font-medium text-notion-text mb-2">
                Application URL
              </label>
              <input
                id="applicationURL"
                v-model="form.applicationURL"
                type="url"
                required
                class="input-field max-w-md"
                placeholder="http://zewi.com"
                :disabled="loading"
              >
            </div>
          </div>
        </div>

        <!-- SMTP Configuration Box -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-4">SMTP Configuration</h2>
          <div class="border-b border-notion-border mb-4"></div>
          <div class="space-y-4">
            <div>
              <label for="smtpServer" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Server
              </label>
              <input
                id="smtpServer"
                v-model="form.smtpServer"
                type="text"
                class="input-field max-w-md"
                placeholder="smtp.zewi.com"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpPort" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Port
              </label>
              <input
                id="smtpPort"
                v-model="form.smtpPort"
                type="text"
                class="input-field max-w-xs"
                placeholder="587"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpFromEmail" class="block text-sm font-medium text-notion-text mb-2">
                SMTP From Email
              </label>
              <input
                id="smtpFromEmail"
                v-model="form.smtpFromEmail"
                type="email"
                class="input-field max-w-md"
                placeholder="noreply@zewi.com"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpUsername" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Username
              </label>
              <input
                id="smtpUsername"
                v-model="form.smtpUsername"
                type="text"
                class="input-field max-w-md"
                placeholder="smtpuser"
                :disabled="loading"
              >
            </div>

            <div>
              <label for="smtpPassword" class="block text-sm font-medium text-notion-text mb-2">
                SMTP Password
              </label>
              <input
                id="smtpPassword"
                v-model="form.smtpPassword"
                type="password"
                class="input-field max-w-md"
                placeholder="Enter SMTP password"
                :disabled="loading"
              >
            </div>

            <div class="flex items-start">
              <input
                id="smtpUseTLS"
                v-model="form.smtpUseTLS"
                type="checkbox"
                class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
                :disabled="loading"
              >
              <label for="smtpUseTLS" class="ml-3 block">
                <span class="text-sm font-medium text-notion-text">Use TLS</span>
                <p class="text-xs text-notion-textLight mt-1">Enable TLS for SMTP connection</p>
              </label>
            </div>
          </div>
        </div>

        <!-- Maintenance Mode Box -->
        <div class="bg-white rounded-lg border border-notion-border p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-notion-text mb-4">Maintenance</h2>
          <div class="border-b border-notion-border mb-4"></div>
          <div class="flex items-start">
            <input
              id="maintenanceMode"
              v-model="form.maintenanceMode"
              type="checkbox"
              class="h-4 w-4 mt-1 rounded text-notion-text focus:ring-notion-text border-notion-border"
              :disabled="loading"
            >
            <label for="maintenanceMode" class="ml-3 block">
              <span class="text-sm font-medium text-notion-text">Enable Maintenance Mode</span>
              <p class="text-xs text-notion-textLight mt-1">Put the application in maintenance mode</p>
            </label>
          </div>
        </div>

        <!-- Submit Button -->
        <div class="flex justify-end gap-3">
          <button
            type="submit"
            class="btn-primary"
            :disabled="loading"
          >
            <span v-if="!loading">Save Settings</span>
            <span v-else class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              Saving...
            </span>
          </button>
        </div>
      </form>
    </main>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { settingsAPI } from '@/api'
import NavBar from '@/components/NavBar.vue'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const successMessage = ref(null)
const errorMessage = ref(null)

const form = reactive({
  applicationURL: '',
  applicationEmail: '',
  applicationName: '',
  maintenanceMode: false,
  smtpServer: '',
  smtpPort: '',
  smtpFromEmail: '',
  smtpUsername: '',
  smtpPassword: '',
  smtpUseTLS: true
})


const loadSettings = async () => {
  loading.value = true
  errorMessage.value = null

  try {
    const response = await settingsAPI.getSettings()
    const settings = response.data?.settings

    if (settings) {
      form.applicationURL = settings.ApplicationURL || ''
      form.applicationEmail = settings.ApplicationEmail || ''
      form.applicationName = settings.ApplicationName || ''
      form.maintenanceMode = settings.MaintenanceMode || false
      form.smtpServer = settings.SMTPServer || ''
      form.smtpPort = settings.SMTPPort || ''
      form.smtpFromEmail = settings.SMTPFromEmail || ''
      form.smtpUsername = settings.SMTPUsername || ''
      form.smtpPassword = settings.SMTPPassword || ''
      form.smtpUseTLS = settings.SMTPUseTLS || false
    }
  } catch (err) {
    console.error('Failed to load settings:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// Watch for message changes and scroll to top
watch([successMessage, errorMessage], () => {
  if (successMessage.value || errorMessage.value) {
    nextTick(() => {
      scrollToTop()
    })
  }
})

const handleSave = async () => {
  loading.value = true
  successMessage.value = null
  errorMessage.value = null

  try {
    const payload = {
      applicationURL: form.applicationURL,
      applicationEmail: form.applicationEmail,
      applicationName: form.applicationName,
      maintenanceMode: form.maintenanceMode,
      smtpUseTLS: form.smtpUseTLS
    }

    // Only include SMTP fields if they have values
    if (form.smtpServer){ payload.smtpServer = form.smtpServer }
    if (form.smtpPort){ payload.smtpPort = form.smtpPort }
    if (form.smtpFromEmail){ payload.smtpFromEmail = form.smtpFromEmail }
    if (form.smtpUsername){ payload.smtpUsername = form.smtpUsername }
    if (form.smtpPassword){ payload.smtpPassword = form.smtpPassword }

    const response = await settingsAPI.updateSettings(payload)

    if (response.data?.successMessage) {
      successMessage.value = response.data.successMessage
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
    }
  } catch (err) {
    console.error('Failed to save settings:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to save settings. Please try again.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

