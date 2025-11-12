<template>
  <div class="min-h-screen bg-theme-bg">
    <NavBar />

    <main class="w-full py-8 px-6 lg:px-8">
      <div class="page-header">
        <h1 class="page-title">Profile</h1>
        <p class="page-subtitle">View your account information</p>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
        <p class="text-sm text-red-800">{{ error }}</p>
      </div>

      <!-- Success Message -->
      <div v-if="successMessage" class="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
        <p class="text-sm text-green-800">{{ successMessage }}</p>
      </div>

      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="text-theme-textLight">Loading profile...</div>
      </div>

      <div v-else>
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <!-- Main Content -->
          <div class="lg:col-span-2 space-y-6">
            <!-- Account Information -->
            <div class="card">
              <h3 class="text-lg font-semibold text-theme-text mb-6">Account Information</h3>
              <div class="form-group">
                <div>
                  <label class="form-label">Name</label>
                  <input
                    v-model="profile.name"
                    type="text"
                    class="input-field"
                    placeholder="Enter your name"
                    :disabled="savingProfile"
                  />
                </div>
                <div>
                  <label class="form-label">Email Address</label>
                  <input
                    v-model="profile.email"
                    type="email"
                    class="input-field"
                    placeholder="Enter your email"
                    :disabled="savingProfile"
                  />
                </div>
                <div class="flex justify-end pt-2">
                  <button
                    @click="saveProfile"
                    class="btn-primary"
                    :disabled="savingProfile"
                  >
                    <span v-if="savingProfile">Saving...</span>
                    <span v-else>Save Changes</span>
                  </button>
                </div>
              </div>
            </div>

            <!-- Change Password -->
            <div class="card">
              <h3 class="text-lg font-semibold text-theme-text mb-6">Change Password</h3>
              <div class="form-group">
                <div>
                  <label class="form-label">Current Password</label>
                  <input
                    v-model="passwordForm.currentPassword"
                    type="password"
                    class="input-field"
                    placeholder="Enter current password"
                    :disabled="savingPassword"
                  />
                </div>
                <div>
                  <label class="form-label">New Password</label>
                  <input
                    v-model="passwordForm.newPassword"
                    type="password"
                    class="input-field"
                    placeholder="Enter new password"
                    :disabled="savingPassword"
                  />
                  <p class="text-xs text-theme-textLight mt-1">
                    Password must be at least 8 characters long and contain uppercase, lowercase, number, and special character.
                  </p>
                </div>
                <div>
                  <label class="form-label">Confirm New Password</label>
                  <input
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    class="input-field"
                    placeholder="Confirm new password"
                    :disabled="savingPassword"
                  />
                </div>
                <div class="flex justify-end pt-2">
                  <button
                    @click="savePassword"
                    class="btn-primary"
                    :disabled="savingPassword || !isPasswordFormValid"
                  >
                    <span v-if="savingPassword">Updating...</span>
                    <span v-else>Update Password</span>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Sidebar - API Key -->
          <div class="space-y-6">
            <div class="card">
              <h3 class="text-lg font-semibold text-theme-text mb-6">API Key</h3>
              <div class="form-group">
                <div>
                  <label class="form-label">API Key</label>
                  <div class="flex flex-col gap-2">
                    <div class="flex items-start gap-2">
                      <textarea
                        :value="displayedAPIKey"
                        class="input-field font-mono text-xs resize-none overflow-hidden min-h-[2.5rem]"
                        readonly
                        :disabled="loadingAPIKey"
                        rows="1"
                        style="word-break: break-all; white-space: pre-wrap;"
                        ref="apiKeyTextarea"
                      />
                    </div>
                    <div class="flex items-center gap-2">
                      <button
                        @click="toggleAPIKeyVisibility"
                        type="button"
                        class="px-3 py-2 text-sm text-theme-text hover:text-theme-textLight border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed flex-1"
                        :disabled="loadingAPIKey || !apiKey"
                      >
                        {{ showAPIKey ? 'Hide' : 'Show' }}
                      </button>
                      <button
                        @click="copyAPIKey"
                        type="button"
                        class="px-3 py-2 text-sm text-theme-text hover:text-theme-textLight border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed flex-1"
                        :disabled="loadingAPIKey || !apiKey"
                      >
                        {{ copiedAPIKey ? 'Copied!' : 'Copy' }}
                      </button>
                    </div>
                  </div>
                  <p class="text-xs text-theme-textLight mt-2">
                    Use this API key to authenticate API requests. Keep it secure and never share it publicly.
                  </p>
                </div>
                <div class="flex justify-end pt-2">
                  <button
                    @click="openRotateModal"
                    class="btn-secondary w-full"
                    :disabled="rotatingAPIKey || loadingAPIKey"
                  >
                    Rotate API Key
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Rotate API Key Confirmation Modal -->
      <div v-if="showRotateModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeRotateModal">
        <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
          <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeRotateModal"></div>
          <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
            <div class="bg-white px-6 pt-6 pb-4">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-theme-text">Rotate API Key</h3>
                <button
                  type="button"
                  @click="closeRotateModal"
                  class="text-theme-textLight hover:text-theme-text"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <p class="text-sm text-theme-textLight">
                Are you sure you want to rotate your API key? The current key will be invalidated and you will need to update any applications using it.
              </p>
            </div>
            <div class="bg-gray-50 px-6 py-4 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeRotateModal"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                @click="rotateAPIKey"
                class="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 focus:outline-none transition-all duration-150"
                :disabled="rotatingAPIKey"
              >
                <span v-if="!rotatingAPIKey">Rotate</span>
                <span v-else class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Rotating...
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '../components/NavBar.vue'
import { authAPI } from '../api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const savingProfile = ref(false)
const savingPassword = ref(false)
const error = ref(null)
const successMessage = ref(null)

const profile = reactive({
  id: null,
  name: '',
  email: '',
  role: '',
  isActive: false,
  lastLoginAt: null,
  createdAt: null,
  updatedAt: null
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const apiKey = ref(null)
const showAPIKey = ref(false)
const copiedAPIKey = ref(false)
const loadingAPIKey = ref(false)
const rotatingAPIKey = ref(false)
const apiKeyTextarea = ref(null)
const showRotateModal = ref(false)

// Computed property for displayed API key value
const displayedAPIKey = computed(() => {
  if (!apiKey.value) return '••••••••••••••••••••••••••••••••'
  if (showAPIKey.value) return apiKey.value
  return '•'.repeat(apiKey.value.length)
})

// Load profile data
const loadProfile = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await authAPI.getProfile()
    // Handle both nested { user: {...} } and direct user object structures
    const userData = response.data?.user || response.data
    if (userData) {
      profile.id = userData.id
      profile.name = userData.name || ''
      profile.email = userData.email || ''
      profile.role = userData.role || ''
      profile.isActive = userData.isActive !== false
      profile.lastLoginAt = userData.lastLoginAt || null
      profile.createdAt = userData.createdAt || null
      profile.updatedAt = userData.updatedAt || null
    }
  } catch (err) {
    console.error('Load profile error:', err)
    error.value = err.response?.data?.errorMessage || 'Failed to load profile. Please try again.'
  } finally {
    loading.value = false
  }
}

// Check if password form is valid
const isPasswordFormValid = computed(() => {
  return passwordForm.currentPassword.trim() !== '' &&
         passwordForm.newPassword.trim() !== '' &&
         passwordForm.confirmPassword.trim() !== '' &&
         passwordForm.newPassword === passwordForm.confirmPassword &&
         passwordForm.newPassword.length >= 8
})

// Save profile
const saveProfile = async () => {
  if (savingProfile.value) return
  savingProfile.value = true
  error.value = null
  successMessage.value = null
  try {
    const updateData = {
      name: profile.name.trim(),
      email: profile.email.trim(),
    }
    const response = await authAPI.updateProfile(updateData)
    if (response.data?.successMessage || response.data?.user) {
      successMessage.value = response.data.successMessage || 'Profile updated successfully'
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
      // Reload profile to get latest data
      await loadProfile()
      // Refresh auth store to get updated user data
      await authStore.checkAuth()
    }
  } catch (err) {
    console.error('Update profile error:', err)
    error.value = err.response?.data?.errorMessage || 'Failed to update profile. Please try again.'
  } finally {
    savingProfile.value = false
  }
}

// Save password
const savePassword = async () => {
  if (savingPassword.value || !isPasswordFormValid.value) return

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    error.value = 'New password and confirm password do not match'
    return
  }

  savingPassword.value = true
  error.value = null
  successMessage.value = null

  try {
    const updateData = {
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword
    }
    const response = await authAPI.updatePassword(updateData)
    if (response.data?.message) {
      successMessage.value = response.data.message || 'Password updated successfully'
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
      // Reset password form
      passwordForm.currentPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
    }
  } catch (err) {
    console.error('Update password error:', err)
    error.value = err.response?.data?.errorMessage || 'Failed to update password. Please try again.'
  } finally {
    savingPassword.value = false
  }
}

// Load API key
const loadAPIKey = async () => {
  loadingAPIKey.value = true
  error.value = null
  try {
    const response = await authAPI.getAPIKey()
    if (response.data?.apiKey) {
      apiKey.value = response.data.apiKey
      adjustTextareaHeight()
    }
  } catch (err) {
    console.error('Load API key error:', err)
    // Don't show error for API key loading, just log it
  } finally {
    loadingAPIKey.value = false
  }
}

// Adjust textarea height
const adjustTextareaHeight = () => {
  nextTick(() => {
    if (apiKeyTextarea.value) {
      apiKeyTextarea.value.style.height = 'auto'
      apiKeyTextarea.value.style.height = apiKeyTextarea.value.scrollHeight + 'px'
    }
  })
}

// Toggle API key visibility
const toggleAPIKeyVisibility = () => {
  showAPIKey.value = !showAPIKey.value
  adjustTextareaHeight()
}

// Watch for API key changes to adjust textarea height
watch([apiKey, showAPIKey], () => {
  adjustTextareaHeight()
})

// Copy API key to clipboard
const copyAPIKey = async () => {
  if (!apiKey.value) return

  try {
    await navigator.clipboard.writeText(apiKey.value)
    copiedAPIKey.value = true
    setTimeout(() => {
      copiedAPIKey.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy API key:', err)
    error.value = 'Failed to copy API key to clipboard'
  }
}

// Open rotate API key modal
const openRotateModal = () => {
  showRotateModal.value = true
}

// Close rotate API key modal
const closeRotateModal = () => {
  showRotateModal.value = false
}

// Rotate API key
const rotateAPIKey = async () => {
  if (rotatingAPIKey.value) return

  rotatingAPIKey.value = true
  error.value = null
  successMessage.value = null

  try {
    const response = await authAPI.rotateAPIKey()
    if (response.data?.apiKey) {
      apiKey.value = response.data.apiKey
      showAPIKey.value = true // Show the new key
      adjustTextareaHeight()
      successMessage.value = response.data.message || 'API key rotated successfully'
      setTimeout(() => {
        successMessage.value = null
      }, 5000)
      closeRotateModal()
    }
  } catch (err) {
    console.error('Rotate API key error:', err)
    error.value = err.response?.data?.errorMessage || 'Failed to rotate API key. Please try again.'
  } finally {
    rotatingAPIKey.value = false
  }
}

// Check authentication on mount
onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  await Promise.all([loadProfile(), loadAPIKey()])
})
</script>
