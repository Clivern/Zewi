<template>
  <div class="min-h-screen bg-notion-bg">
    <NavBar />

    <!-- Main Content -->
    <main class="w-full py-8 px-6 lg:px-8">
      <!-- Breadcrumb Navigation -->
      <div class="mb-4">
        <nav class="flex items-center space-x-2 text-sm">
          <router-link
            to="/admin/buckets"
            class="text-blue-600 hover:text-blue-900"
          >
            Buckets
          </router-link>
          <svg class="h-4 w-4 text-notion-textLight" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
          <span class="text-notion-text">{{ bucketName }}</span>
          <template v-if="currentPath">
            <svg class="h-4 w-4 text-notion-textLight" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
            <span class="text-notion-text">{{ currentPath }}</span>
          </template>
        </nav>
      </div>

      <!-- Page Header -->
      <div class="mb-8 flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-semibold text-notion-text">{{ bucketName }}</h1>
          <p class="text-notion-textLight mt-2">Browse directories and files</p>
        </div>
        <div class="flex items-center space-x-3">
          <button
            @click="openUploadModal"
            class="btn-primary"
          >
            Upload File
          </button>
          <button
            @click="openCreateFolderModal"
            class="btn-secondary"
          >
            Create Folder
          </button>
        </div>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="mb-6 rounded-md border border-green-200 bg-green-50 p-4">
        <p class="text-sm text-green-800">{{ successMessage }}</p>
      </div>
      <div v-if="errorMessage" class="mb-6 rounded-md border border-red-200 bg-red-50 p-4">
        <p class="text-sm text-red-800">{{ errorMessage }}</p>
      </div>

      <!-- Files and Directories Table -->
      <div class="bg-white rounded-lg border border-notion-border overflow-hidden">
        <div v-if="loading" class="p-8 text-center">
          <svg class="animate-spin h-6 w-6 mx-auto text-notion-text" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-notion-textLight mt-2">Loading...</p>
        </div>

        <div v-else-if="items.length === 0" class="p-8 text-center">
          <p class="text-notion-textLight">This directory is empty</p>
        </div>

        <table v-else class="min-w-full divide-y divide-notion-border">
          <thead class="bg-notion-hover">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Size
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Modified
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-notion-textLight uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-notion-border">
            <!-- Directories first -->
            <tr v-for="item in directories" :key="item.name" class="hover:bg-notion-hover">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-yellow-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                  </svg>
                  <button
                    @click="navigateToPath(item.name)"
                    class="text-sm font-medium text-blue-600 hover:text-blue-900"
                  >
                    {{ item.name }}
                  </button>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                —
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ formatDate(item.modifiedAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="openDeleteModal(item, 'directory')"
                  class="text-red-600 hover:text-red-900"
                >
                  Delete
                </button>
              </td>
            </tr>
            <!-- Files -->
            <tr v-for="item in files" :key="item.name" class="hover:bg-notion-hover">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-blue-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <span class="text-sm font-medium text-notion-text">{{ item.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ formatSize(item.size) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-notion-textLight">
                {{ formatDate(item.modifiedAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="downloadFile(item.name)"
                  class="text-blue-600 hover:text-blue-900 mr-4"
                >
                  Download
                </button>
                <button
                  @click="openDeleteModal(item, 'file')"
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
            Showing {{ offset + 1 }} to {{ Math.min(offset + limit, total) }} of {{ total }} items
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

    <!-- Upload File Modal -->
    <div v-if="showUploadModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeUploadModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeUploadModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleUpload">
            <div class="bg-white px-6 pt-6 pb-4">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-notion-text">Upload File</h3>
                <button
                  type="button"
                  @click="closeUploadModal"
                  class="text-notion-textLight hover:text-notion-text"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="space-y-4">
                <div>
                  <label for="upload-file" class="block text-sm font-medium text-notion-text mb-2">
                    Select File *
                  </label>
                  <input
                    id="upload-file"
                    ref="fileInput"
                    type="file"
                    required
                    class="input-field"
                    @change="handleFileSelect"
                  />
                </div>
                <div v-if="uploadForm.fileName">
                  <p class="text-sm text-notion-textLight">Selected: {{ uploadForm.fileName }}</p>
                </div>
              </div>
            </div>
            <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeUploadModal"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="uploadLoading || !uploadForm.file"
              >
                <span v-if="!uploadLoading">Upload</span>
                <span v-else class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                  Uploading...
                </span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Create Folder Modal -->
    <div v-if="showCreateFolderModal" class="fixed inset-0 z-50 overflow-y-auto" @click.self="closeCreateFolderModal">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75" @click="closeCreateFolderModal"></div>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleCreateFolder">
            <div class="bg-white px-6 pt-6 pb-4">
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-notion-text">Create Folder</h3>
                <button
                  type="button"
                  @click="closeCreateFolderModal"
                  class="text-notion-textLight hover:text-notion-text"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="space-y-4">
                <div>
                  <label for="folder-name" class="block text-sm font-medium text-notion-text mb-2">
                    Folder Name *
                  </label>
                  <input
                    id="folder-name"
                    v-model="folderForm.name"
                    type="text"
                    required
                    class="input-field"
                    placeholder="my-folder"
                  />
                </div>
              </div>
            </div>
            <div class="bg-notion-hover px-6 py-4 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeCreateFolderModal"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="folderLoading"
              >
                <span v-if="!folderLoading">Create</span>
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
              <h3 class="text-lg font-semibold text-notion-text">Delete {{ itemToDeleteType === 'directory' ? 'Folder' : 'File' }}</h3>
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
              Are you sure you want to delete {{ itemToDeleteType === 'directory' ? 'folder' : 'file' }} <strong class="text-notion-text">{{ itemToDelete?.name }}</strong>?
              {{ itemToDeleteType === 'directory' ? 'This will delete all contents inside the folder.' : '' }}
              This action cannot be undone.
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { bucketAPI } from '@/api'
import NavBar from '@/components/NavBar.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Get bucket name from route params
const bucketName = computed(() => decodeURIComponent(route.params.bucketName || ''))
const currentPath = computed(() => {
  const path = route.query.path || ''
  return path ? decodeURIComponent(path) : ''
})

// State
const loading = ref(false)
const items = ref([])
const total = ref(0)
const limit = ref(50)
const offset = ref(0)
const successMessage = ref(null)
const errorMessage = ref(null)

// Modal states
const showUploadModal = ref(false)
const showCreateFolderModal = ref(false)
const showDeleteModal = ref(false)
const uploadLoading = ref(false)
const folderLoading = ref(false)
const deleteLoading = ref(false)
const itemToDelete = ref(null)
const itemToDeleteType = ref('file')
const fileInput = ref(null)

// Forms
const uploadForm = reactive({
  file: null,
  fileName: ''
})

const folderForm = reactive({
  name: ''
})

// Computed properties
const directories = computed(() => {
  return items.value.filter(item => item.type === 'directory' || item.isDirectory)
})

const files = computed(() => {
  return items.value.filter(item => item.type === 'file' || !item.isDirectory)
})


const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime()) || date.getFullYear() < 1970) {
      return 'N/A'
    }

    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMs / 3600000)
    const diffDays = Math.floor(diffMs / 86400000)

    if (diffMins < 1) return 'Just now'
    if (diffMins < 60) return `${diffMins} ${diffMins === 1 ? 'minute' : 'minutes'} ago`
    if (diffHours < 24) return `${diffHours} ${diffHours === 1 ? 'hour' : 'hours'} ago`
    if (diffDays < 7) return `${diffDays} ${diffDays === 1 ? 'day' : 'days'} ago`

    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
    const month = months[date.getMonth()]
    const day = date.getDate()
    const year = date.getFullYear()
    const isCurrentYear = year === now.getFullYear()

    if (isCurrentYear) {
      return `${month} ${day}`
    } else {
      return `${month} ${day}, ${year}`
    }
  } catch {
    return 'N/A'
  }
}

const formatSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const loadItems = async () => {
  loading.value = true
  errorMessage.value = null

  try {
    const params = {
      limit: limit.value,
      offset: offset.value
    }
    if (currentPath.value) {
      params.path = currentPath.value
    }

    const response = await bucketAPI.listBucketItems(bucketName.value, params)

    if (response.data) {
      items.value = response.data.items || []
      total.value = response.data._meta?.total || 0
    }
  } catch (err) {
    console.error('Failed to load items:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to load items'
  } finally {
    loading.value = false
  }
}

const goToPage = (newOffset) => {
  if (newOffset < 0) return
  if (newOffset >= total.value) return
  offset.value = newOffset
  loadItems()
}

const navigateToPath = (pathName) => {
  const newPath = currentPath.value
    ? `${currentPath.value}/${pathName}`
    : pathName
  router.push({
    name: 'BucketFiles',
    params: { bucketName: bucketName.value },
    query: { path: encodeURIComponent(newPath) }
  })
}

const openUploadModal = () => {
  uploadForm.file = null
  uploadForm.fileName = ''
  showUploadModal.value = true
}

const closeUploadModal = () => {
  showUploadModal.value = false
  uploadForm.file = null
  uploadForm.fileName = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    uploadForm.file = file
    uploadForm.fileName = file.name
  }
}

const handleUpload = async () => {
  if (!uploadForm.file) return

  uploadLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const formData = new FormData()
    formData.append('file', uploadForm.file)

    const path = currentPath.value ? `${currentPath.value}/${uploadForm.fileName}` : uploadForm.fileName

    await bucketAPI.uploadFile(bucketName.value, path, formData)

    successMessage.value = 'File uploaded successfully'
    closeUploadModal()
    loadItems()
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err) {
    console.error('Failed to upload file:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else {
      errorMessage.value = 'Failed to upload file. Please try again.'
    }
  } finally {
    uploadLoading.value = false
  }
}

const openCreateFolderModal = () => {
  folderForm.name = ''
  showCreateFolderModal.value = true
}

const closeCreateFolderModal = () => {
  showCreateFolderModal.value = false
  folderForm.name = ''
}

const handleCreateFolder = async () => {
  folderLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const path = currentPath.value
      ? `${currentPath.value}/${folderForm.name}`
      : folderForm.name

    await bucketAPI.createFolder(bucketName.value, path)

    successMessage.value = 'Folder created successfully'
    closeCreateFolderModal()
    loadItems()
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err) {
    console.error('Failed to create folder:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else {
      errorMessage.value = 'Failed to create folder. Please try again.'
    }
  } finally {
    folderLoading.value = false
  }
}

const openDeleteModal = (item, type) => {
  itemToDelete.value = item
  itemToDeleteType.value = type
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  itemToDelete.value = null
  itemToDeleteType.value = 'file'
}

const handleDelete = async () => {
  if (!itemToDelete.value) return

  deleteLoading.value = true
  errorMessage.value = null
  successMessage.value = null

  try {
    const path = currentPath.value
      ? `${currentPath.value}/${itemToDelete.value.name}`
      : itemToDelete.value.name

    await bucketAPI.deleteItem(bucketName.value, path, itemToDeleteType.value)

    successMessage.value = `${itemToDeleteType.value === 'directory' ? 'Folder' : 'File'} deleted successfully`
    closeDeleteModal()
    loadItems()
    setTimeout(() => {
      successMessage.value = null
    }, 3000)
  } catch (err) {
    console.error('Failed to delete item:', err)
    if (err.response?.data?.errorMessage) {
      errorMessage.value = err.response.data.errorMessage
    } else {
      errorMessage.value = `Failed to delete ${itemToDeleteType.value}. Please try again.`
    }
  } finally {
    deleteLoading.value = false
  }
}

const downloadFile = async (fileName) => {
  try {
    const path = currentPath.value
      ? `${currentPath.value}/${fileName}`
      : fileName

    const response = await bucketAPI.downloadFile(bucketName.value, path)

    // Create a blob and download it
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Failed to download file:', err)
    errorMessage.value = err.response?.data?.errorMessage || 'Failed to download file'
  }
}

// Watch for route changes to reload items
watch(() => route.query.path, () => {
  loadItems()
})

onMounted(() => {
  if (!bucketName.value) {
    router.push('/admin/buckets')
    return
  }
  loadItems()
})
</script>

