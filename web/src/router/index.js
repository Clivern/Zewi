import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Login from '@/views/Login.vue'
import Dashboard from '@/views/Dashboard.vue'
import Setup from '@/views/Setup.vue'
import Settings from '@/views/Settings.vue'
import Users from '@/views/Users.vue'
import Buckets from '@/views/Buckets.vue'
import BucketFiles from '@/views/BucketFiles.vue'
import Profile from '@/views/Profile.vue'
import NotFound from '@/views/NotFound.vue'
import ServerError from '@/views/ServerError.vue'

const routes = [
  {
    path: '/setup',
    name: 'Setup',
    component: Setup,
    meta: { requiresGuest: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/admin/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/users',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/settings',
    name: 'Settings',
    component: Settings,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/buckets',
    name: 'Buckets',
    component: Buckets,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/buckets/:bucketName',
    name: 'BucketFiles',
    component: BucketFiles,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  },
  {
    path: '/500',
    name: 'ServerError',
    component: ServerError
  },
  {
    path: '/404',
    name: 'NotFound',
    component: NotFound
  },
  {
    path: '/',
    redirect: '/admin/dashboard'
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard for authentication and setup
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Wait for auth initialization to complete before making decisions
  // This prevents redirects during the initial auth check on page refresh
  while (authStore.isInitializing) {
    await new Promise(resolve => setTimeout(resolve, 50))
  }

  const isAuthenticated = authStore.isAuthenticated

  // Check if setup is completed (you can implement this check via API)
  // For now, we'll allow access to setup page if not authenticated

  // Allow access to error pages without authentication
  if (to.name === 'NotFound' || to.name === 'ServerError') {
    next()
    return
  }

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next('/admin/dashboard')
  } else {
    next()
  }
})

// Global error handler
router.onError((error) => {
  console.error('Router error:', error)
  router.push('/500')
})

export default router
