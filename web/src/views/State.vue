<template>
  <div class="min-h-screen bg-[#F7F6F3] flex flex-col items-center justify-center px-4">
    <div class="w-full max-w-2xl">
      <div class="card">
        <div class="page-header">
          <h1 class="page-title">State Management</h1>
          <p class="page-subtitle">View and update the current state value</p>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="text-center py-8">
          <p class="text-[#6B6966]">Loading state...</p>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
          <p class="text-red-800 text-sm">{{ error }}</p>
        </div>

        <!-- Success Message -->
        <div v-if="success" class="mb-6 p-4 bg-green-50 border border-green-200 rounded-md">
          <p class="text-green-800 text-sm">{{ success }}</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="form-group" v-if="!loading">
          <div>
            <label for="state" class="form-label">Current State Value</label>
            <input
              id="state"
              v-model="stateValue"
              type="text"
              class="input-field"
              placeholder="Enter state value"
              :disabled="updating"
            />
            <p v-if="currentState !== null && currentState !== stateValue" class="text-xs text-[#6B6966] mt-2">
              Current value: <span class="font-mono">{{ currentState || '(empty)' }}</span>
            </p>
          </div>

          <div class="flex gap-3">
            <button
              type="submit"
              class="btn-primary flex-1"
              :disabled="updating || stateValue === currentState"
            >
              <span v-if="updating">Updating...</span>
              <span v-else>Update State</span>
            </button>
            <button
              type="button"
              class="btn-secondary"
              @click="fetchState"
              :disabled="updating || loading"
            >
              Refresh
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getState, updateState } from '@/api/index.js'

const stateValue = ref('')
const currentState = ref(null)
const loading = ref(false)
const updating = ref(false)
const error = ref('')
const success = ref('')

const fetchState = async () => {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await getState()
    currentState.value = response.state
    stateValue.value = response.state || ''
  } catch (err) {
    error.value = err.message || 'Failed to fetch state'
    console.error('Error fetching state:', err)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  updating.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await updateState(stateValue.value)
    currentState.value = response.state
    success.value = 'State updated successfully'
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      success.value = ''
    }, 3000)
  } catch (err) {
    error.value = err.message || 'Failed to update state'
    console.error('Error updating state:', err)
  } finally {
    updating.value = false
  }
}

onMounted(() => {
  fetchState()
})
</script>
