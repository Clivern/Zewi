import { createRouter, createWebHistory } from 'vue-router'
import Landing from '@/views/Landing.vue'
import State from '@/views/State.vue'

const routes = [  
  {
    path: '/',
    name: 'Landing',
    component: Landing
  },
  {
    path: '/state',
    name: 'State',
    component: State
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'Landing',
    component: Landing
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
