import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuth } from '@/services/auth'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Admin from '@/views/Admin.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from, next) => {
  const { verify, isAuthenticated } = useAuth()

  if (to.meta.requiresAuth) {
    if (!isAuthenticated()) {
      next('/login')
      return
    }

    const isValid = await verify()
    if (!isValid) {
      next('/login')
      return
    }
  }

  next()
})

export default router
