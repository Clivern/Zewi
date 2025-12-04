import axios from 'axios'

// Get API base URL from window (injected at runtime) or environment variable, fallback to empty string
const API_BASE_URL = (typeof window !== 'undefined' && window.API_BASE_URL !== undefined)
	? window.API_BASE_URL
	: (import.meta.env.VITE_API_BASE_URL || '')

// Create axios instance with base configuration
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

/**
 * Get the current state value
 * @returns {Promise<{state: string|null}>}
 */
export async function getState() {
  try {
    const response = await apiClient.get('/api/v1/state')
    return response.data
  } catch (error) {
    if (error.response) {
      // Server responded with error status
      throw new Error(error.response.data?.error || `Server error: ${error.response.status}`)
    } else if (error.request) {
      // Request made but no response received
      throw new Error('No response from server. Please check your connection.')
    } else {
      // Something else happened
      throw new Error(error.message || 'Failed to fetch state')
    }
  }
}

/**
 * Update the state value
 * @param {string} value - The new state value
 * @returns {Promise<{message: string, state: string}>}
 */
export async function updateState(value) {
  try {
    const response = await apiClient.put('/api/v1/state', { value })
    return response.data
  } catch (error) {
    if (error.response) {
      // Server responded with error status
      throw new Error(error.response.data?.error || `Server error: ${error.response.status}`)
    } else if (error.request) {
      // Request made but no response received
      throw new Error('No response from server. Please check your connection.')
    } else {
      // Something else happened
      throw new Error(error.message || 'Failed to update state')
    }
  }
}
