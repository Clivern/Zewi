<template>
  <nav class="bg-white border-b border-theme-border relative">
    <div class="w-full px-6 lg:px-8">
      <div class="flex justify-between h-14">
        <div class="flex items-center">
          <div class="flex-shrink-0 flex items-center">
            <router-link to="/admin/dashboard">
              <img src="/logo.png" alt="Zewi Logo" class="h-8 w-auto">
            </router-link>
          </div>
          <!-- Desktop Navigation -->
          <div class="hidden md:ml-8 md:flex md:space-x-1">
            <router-link
              to="/admin/dashboard"
              :class="[
                'inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                $route.path === '/admin/dashboard'
                  ? 'bg-theme-hover text-theme-text'
                  : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
              ]"
            >
              Dashboard
            </router-link>
            <router-link
              to="/admin/buckets"
              :class="[
                'inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                $route.path.startsWith('/admin/buckets')
                  ? 'bg-theme-hover text-theme-text'
                  : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
              ]"
            >
              Buckets
            </router-link>
            <router-link
              to="/admin/users"
              :class="[
                'inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                $route.path === '/admin/users'
                  ? 'bg-theme-hover text-theme-text'
                  : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
              ]"
            >
              Users
            </router-link>
            <router-link
              to="/admin/settings"
              :class="[
                'inline-flex items-center px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
                $route.path === '/admin/settings'
                  ? 'bg-theme-hover text-theme-text'
                  : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
              ]"
            >
              Settings
            </router-link>
          </div>
        </div>
        <div class="flex items-center">
          <div v-if="isAuthenticated" class="flex items-center space-x-3">
            <!-- User Dropdown Menu -->
            <div class="relative hidden sm:block" ref="dropdownRef">
              <button
                @click="toggleUserDropdown"
                class="flex items-center space-x-2 p-1 rounded-full hover:bg-theme-hover focus:outline-none focus:ring-2 focus:ring-theme-text focus:ring-offset-2 transition-colors"
                aria-label="User menu"
                aria-expanded="false"
              >
                <div class="w-8 h-8 rounded-full overflow-hidden bg-theme-text flex items-center justify-center">
                  <img
                    v-if="userAvatar"
                    :src="userAvatar"
                    :alt="userEmail"
                    class="w-full h-full object-cover"
                  />
                  <span v-else class="text-white text-xs font-medium">
                    {{ userInitials }}
                  </span>
                </div>
                <svg
                  class="w-4 h-4 text-theme-textLight transition-transform"
                  :class="{ 'rotate-180': userDropdownOpen }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <!-- Dropdown Menu -->
              <transition
                enter-active-class="transition ease-out duration-100"
                enter-from-class="opacity-0 scale-95"
                enter-to-class="opacity-100 scale-100"
                leave-active-class="transition ease-in duration-75"
                leave-from-class="opacity-100 scale-100"
                leave-to-class="opacity-0 scale-95"
              >
                <div
                  v-if="userDropdownOpen"
                  class="absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-50"
                >
                  <div class="py-1">
                    <div class="px-4 py-3 border-b border-theme-border">
                      <p class="text-sm font-medium text-theme-text">{{ userName }}</p>
                      <p class="text-sm text-theme-textLight truncate">{{ userEmail }}</p>
                    </div>
                    <router-link
                      to="/admin/profile"
                      @click="closeUserDropdown"
                      class="block px-4 py-2 text-sm text-theme-textLight hover:bg-theme-hover hover:text-theme-text transition-colors"
                    >
                      Profile
                    </router-link>
                    <div class="border-t border-theme-border my-1"></div>
                    <button
                      @click="handleLogout"
                      class="block w-full text-left px-4 py-2 text-sm text-theme-textLight hover:bg-theme-hover hover:text-theme-text transition-colors"
                    >
                      Logout
                    </button>
                  </div>
                </div>
              </transition>
            </div>

            <!-- Mobile Menu Button -->
            <button
              @click="toggleMobileMenu"
              class="md:hidden p-2 rounded-md text-theme-text hover:bg-theme-hover focus:outline-none"
              aria-label="Toggle menu"
            >
              <svg
                v-if="!mobileMenuOpen"
                class="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
              <svg
                v-else
                class="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Menu -->
    <transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <div
        v-if="mobileMenuOpen && isAuthenticated"
        class="md:hidden border-t border-theme-border bg-white"
      >
        <div class="px-6 py-4 space-y-1">
          <router-link
            to="/admin/dashboard"
            @click="closeMobileMenu"
            class="block px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="
              $route.path === '/admin/dashboard'
                ? 'bg-theme-hover text-theme-text'
                : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
            "
          >
            Dashboard
          </router-link>
          <router-link
            to="/admin/buckets"
            @click="closeMobileMenu"
            class="block px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="
              $route.path.startsWith('/admin/buckets')
                ? 'bg-theme-hover text-theme-text'
                : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
            "
          >
            Buckets
          </router-link>
          <router-link
            to="/admin/users"
            @click="closeMobileMenu"
            class="block px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="
              $route.path === '/admin/users'
                ? 'bg-theme-hover text-theme-text'
                : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
            "
          >
            Users
          </router-link>
          <router-link
            to="/admin/settings"
            @click="closeMobileMenu"
            class="block px-3 py-2 rounded-md text-sm font-medium transition-colors"
            :class="
              $route.path === '/admin/settings'
                ? 'bg-theme-hover text-theme-text'
                : 'text-theme-textLight hover:bg-theme-hover hover:text-theme-text'
            "
          >
            Settings
          </router-link>
          <div class="border-t border-theme-border mt-2 pt-2">
            <div class="flex items-center space-x-3 px-3 py-2 mb-2">
              <div class="w-10 h-10 rounded-full overflow-hidden bg-theme-text flex items-center justify-center flex-shrink-0">
                <img
                  v-if="userAvatar"
                  :src="userAvatar"
                  :alt="userEmail"
                  class="w-full h-full object-cover"
                />
                <span v-else class="text-white text-sm font-medium">
                  {{ userInitials }}
                </span>
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-theme-text truncate">{{ userName }}</p>
                <p class="text-xs text-theme-textLight truncate">{{ userEmail }}</p>
                <p class="text-xs text-theme-textLight truncate">{{ userRole }}</p>
              </div>
            </div>
            <router-link
              to="/admin/profile"
              @click="closeMobileMenu"
              class="block px-3 py-2 rounded-md text-sm font-medium text-theme-textLight hover:bg-theme-hover hover:text-theme-text transition-colors"
            >
              Profile
            </router-link>
            <button
              @click="handleMobileLogout"
              class="w-full text-left px-3 py-2 rounded-md text-sm font-medium text-theme-textLight hover:bg-theme-hover hover:text-theme-text transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </transition>
  </nav>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const mobileMenuOpen = ref(false)
const userDropdownOpen = ref(false)
const dropdownRef = ref(null)

// Get user data from auth store
const user = computed(() => authStore.currentUser)

const isAuthenticated = computed(() => authStore.isAuthenticated)

// Get user name (fallback to email if name not available)
const userName = computed(() => {
  if (user.value?.name && user.value.name.trim()) {
    return user.value.name
  }
  return user.value?.email || 'User'
})

// Get user email
const userEmail = computed(() => {
  return user.value?.email || 'User'
})

// Get user role (capitalize first letter)
const userRole = computed(() => {
  const role = user.value?.role || ''
  return role.charAt(0).toUpperCase() + role.slice(1)
})

// Get user avatar URL
const userAvatar = computed(() => {
  return user.value?.avatar || null
})

// Get user initials for avatar fallback
const userInitials = computed(() => {
  if (!user.value) return 'U'
  // Prefer name initials if available
  if (user.value.name && user.value.name.trim()) {
    const nameParts = user.value.name.trim().split(/\s+/)
    if (nameParts.length >= 2) {
      // First letter of first name + first letter of last name
      return (nameParts[0].charAt(0) + nameParts[nameParts.length - 1].charAt(0)).toUpperCase()
    }
    return nameParts[0].charAt(0).toUpperCase()
  }
  // Fallback to email first letter
  if (user.value.email) {
    return user.value.email.charAt(0).toUpperCase()
  }
  return 'U'
})

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}

const toggleUserDropdown = () => {
  userDropdownOpen.value = !userDropdownOpen.value
}

const closeUserDropdown = () => {
  userDropdownOpen.value = false
}

const handleLogout = () => {
  closeUserDropdown()
  authStore.logout()
  router.push('/login')
}

const handleMobileLogout = () => {
  closeMobileMenu()
  handleLogout()
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    userDropdownOpen.value = false
  }
}

// Close mobile menu and dropdown when route changes
watch(() => route.path, () => {
  mobileMenuOpen.value = false
  userDropdownOpen.value = false
})

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
