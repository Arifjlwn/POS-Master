import { createRouter, createWebHistory } from 'vue-router'

// Import rute modular yang sudah kita pisah-pisah
// import fnbRoutes from '@/modules/fnb/router/index.js'
import retailRoutes from '../modules/retail/router/retail_routes.js'
// import laundryRoutes from '@/modules/jasalayanan/laundry/router/laundry_routes.js'

const baseRoutes = [
  {
    path: '/',
    name: 'LandingPage',
    component: () => import('../views/public/LandingPage.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/auth/Register.vue')
  },
  {
    path: '/select-verify',
    name: 'SelectVerify',
    component: () => import('../views/auth/SelectionVerify.vue')
  },
  {
    path: '/verify-otp',
    name: 'VerifyOTP',
    component: () => import('../views/auth/VerifyOTP.vue')
  },
  {
    path: '/setup-toko',
    name: 'SetupToko',
    component: () => import('../views/onboarding/SetupToko.vue'),
    meta: { requiresAuth: true }
  },
  // {
  //   path: '/:pathMatch(.*)*',
  //   name: 'NotFound',
  //   component: () => import('../views/errors/NotFound.vue')
  // }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...baseRoutes,
    // ...fnbRoutes,
    ...retailRoutes,
    // ...laundryRoutes
  ]
})

// Satpam Guard JWT bawaan lu tetep ditaruh di sini
router.beforeEach((to, from) => {
  const token = localStorage.getItem('token') // sesuaikan dengan cara simpan lu
  
  if (to.meta.requiresAuth && !token) {
    return('/login')
  } else {
    return
  }
})

export default router