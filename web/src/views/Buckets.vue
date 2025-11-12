<template>
  <div class="min-h-screen bg-notion-bg">
    <NavBar />

    <!-- Main Content -->
    <main class="w-full py-8 px-6 lg:px-8">
      <!-- Page Header -->
      <div class="mb-8 flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-semibold text-notion-text">Buckets</h1>
          <p class="text-notion-textLight mt-2">Manage your storage buckets</p>
        </div>
        <button
          @click="openCreateModal"
          class="btn-primary"
        >
          Create Bucket
        </button>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="mb-6 rounded-md border border-green-200 bg-green-50 p-4">
        <p class="text-sm text-green-800">{{ successMessage }}</p>
      </div>
      <div v-if="errorMessage" class="mb-6 rounded-md border border-red-200 bg-red-50 p-4">
        <p class="text-sm text-red-800">{{ errorMessage }}</p>
      </div>

      <!-- Buckets Table -->
      <div class="bg-white rounded-lg border border-notion-border overflow-hidden">
        <div v-if="loading" class="p-8 text-center">
          <svg class="animate-spin h-6 w-6 mx-auto text-notion-text" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-notion-textLight mt-2">Loading buckets...</p>
        </div>

        <div v-else-if="buckets.length === 0" class="p-8 text-center">
          <p class="text-notion-textLight">No buckets found</p>
        </div>

        <table v-else class="min-w-full divide-y divide-notion-border">
          <thead class="bg-notion-hover">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Region
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Created
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-notion-border">
            <tr v-for="bucket in buckets" :key="bucket.name" class="hover:bg-notion-hover">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-notion-textLight mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                  <router-link
                    :to="`/admin/buckets/${encodeURIComponent(bucket.name)}`"
                    class="text-sm font-medium text-blue-600 hover:text-blue-900"
                  >
                    {{ bucket.name }}
                  </router-link>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ bucket.region || 'N/A' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ formatDate(bucket.createdAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <router-link
                  :to="`/admin/buckets/${encodeURIComponent(bucket.name)}`"
                  class="text-blue-600 hover:text-blue-900 mr-4"
                >
                  View
                </router-link>
                <button
                  @click="openDeleteModal(bucket)"
                  class="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="total > 0" class="bg-white px-6 py-4 border-t border-notion-border flex items-center justify-between">
          <div class="text-sm text-notion-textLight">
            Showing {{ offset + 1 }} to {{ Math.min(offset + limit, total) }} of {{ total }} buckets
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="goToPage(offset - limit)"
              :disabled="offset === 0"
              class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            <button
              @click="goToPage(offset + limit)"
              :disabled="offset + limit >= total"
              class="btn-secondary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Create Bucket Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeCreateModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeCreateModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleCreate">
            <div class="bg-white px-6 pt-6 pb-4">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-notion-text">Create Bucket</h3>
                <button
                  type="button"
                  @click="closeCreateModal"
                  class="text-notion-textLight hover:text-notion-text"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="space-y-4">
                <div>
                  <label for="create-name" class="block text-sm font-medium text-notion-text mb-2">
                    Bucket Name *
                  </label>
                  <input
                    id="create-name"
                    v-model="createForm.name"
                    type="text"
                    required
                    class="input-field"
                    placeholder="my-bucket"
                    pattern="[a-z0-9.-]+"
                    title="Bucket name must be lowercase letters, numbers, dots, or hyphens"
                  />
                  <p class="text-xs text-notion-textLight mt-1">Bucket names must be lowercase and can contain letters, numbers, dots, and hyphens</p>
                </div>
                <div>
                  <label for="create-region" class="block text-sm font-medium text-notion-text mb-2">
                    Region
                  </label>
                  <input
                    id="create-region"
                    v-model="createForm.region"
                    type="text"
                    class="input-field"
                    placeholder="us-east-1"
                  />
                </div>
              </div>
            </div>
            <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeCreateModal"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="createLoading"
              >
                <span v-if="!createLoading">Create</span>
                <span v-else class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Creating...
                </span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeDeleteModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeDeleteModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-6 pt-6 pb-4">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-semibold text-notion-text">Delete Bucket</h3>
              <button
                type="button"
                @click="closeDeleteModal"
                class="text-notion-textLight hover:text-notion-text"
              >
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <p class="text-sm text-notion-textLight">
              Are you sure you want to delete bucket <strong class="text-notion-text">{{ bucketToDelete?.name }}</strong>?
              This action cannot be undone and will delete all objects in the bucket.
            </p>
          </div>
          <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
            <button
              type="button"
              @click="closeDeleteModal"
              class="btn-secondary"
            >
              Cancel
            </button>
            <button
              @click="handleDelete"
              class="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 focus:outline-none transition-all duration-150"
              :disabled="deleteLoading"
            >
              <span v-if="!deleteLoading">Delete</span>
              <span v-else class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                Deleting...
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { bucketAPI } from '@/api'
import NavBar from '@/components/NavBar.vue'
import { formatDate } from '@/utils/helpers'

const router = useRouter()
const authStore = useAuthStore()

// State
const loading = ref(false)
const buckets = ref([])
const total = ref(0)
const limit = ref(50)
const offset = ref(0)
const successMessage = ref(null)
const errorMessage = ref(null)

// Modal states
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const createLoading = ref(false)
const deleteLoading = ref(false)
const bucketToDelete = ref(null)

// Forms
const createForm = reactive({
  name: '',
  region: ''
})

const loadBuckets = async () => {
  loading.value = true
  errorMessage.value = null

  try {
    const response = await bucketAPI.getBuckets({
      limit: limit.value,
      offset: offset.value
    })

    if (response.data) {
      buckets.value = response.data.buckets || []
      total.value = response.data._meta?.total || 0
    }
  } catch (err) {
    console.error('Failed to load buckets:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to load buckets'
  } finally {
    loading.value = false
  }
}

const goToPage = (newOffset) => {
  if (newOffset < 0) return
  if (newOffset >= total.value) return
  offset.value = newOffset
  loadBuckets()
}

const openCreateModal = () => {
  createForm.name = ''
  createForm.region = ''
  showCreateModal.value = true
}

const closeCreateModal = () => {
  showCreateModal.value = false
  createForm.name = ''
  createForm.region = ''
}

const handleCreate = async () => {
  createLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const response = await bucketAPI.createBucket({
      name: createForm.name.toLowerCase(),
      region: createForm.region || undefined
    })

    if (response.data) {
      successMessage.value = 'Bucket created successfully'
      closeCreateModal()
      loadBuckets()
      setTimeout(() => {
        successMessage.value = null
      }, 3000)
    }
  } catch (err) {
    console.error('Failed to create bucket:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else {
      errorMessage.value = 'Failed to create bucket. Please try again.'
    }
  } finally {
    createLoading.value = false
  }
}

const openDeleteModal = (bucket) => {
  bucketToDelete.value = bucket
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  bucketToDelete.value = null
}

const handleDelete = async () => {
  if (!bucketToDelete.value) return

  deleteLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    await bucketAPI.deleteBucket(bucketToDelete.value.id)
    successMessage.value = 'Bucket deleted successfully'
    closeDeleteModal()
    loadBuckets()
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err) {
    console.error('Failed to delete bucket:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else {
      errorMessage.value = 'Failed to delete bucket. Please try again.'
    }
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  loadBuckets()
})
</script>

