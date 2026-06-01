import { createRouter, createWebHistory } from 'vue-router'
import Swal from 'sweetalert2' // 🚀 WAJIB IMPORT INI BUAT POP-UP SATPAM

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

// 🛡️ SATPAM GLOBAL VUE ROUTER
router.beforeEach((to, from) => {
  const token = localStorage.getItem('token');
  
  // 🚀 1. LOGIKA ANTI-MUNDUR (GUEST GUARD)
  // Daftar halaman yang HARAM dibuka kalau udah login
  const guestRoutes = ['/login', '/register', '/select-verify', '/verify-otp', '/'];
  
  if (token && guestRoutes.includes(to.path)) {
    // Kalau udah punya token tapi nekat buka halaman login, tendang ke POS!
    return '/retail/pos'; 
  }

  // 🚀 2. CEK AUTENTIKASI (HARUS LOGIN)
  if (to.meta.requiresAuth && !token) {
    return '/login';
  }

  // 🚀 3. LOGIKA SATPAM SAAS (SISTEM KASTA LEVEL)
  if (to.meta.minPlanLevel) {
    const subPlan = localStorage.getItem('subscriptionPlan') || 'basic';

    const getPlanLevel = (plan) => {
      const p = plan.toLowerCase();
      if (p === 'premium' || p === 'trial') return 3;
      if (p === 'pro') return 2;
      return 1; // Basic
    };

    const currentLevel = getPlanLevel(subPlan);

    // Kalau kasta user kurang dari syarat halaman
    if (currentLevel < to.meta.minPlanLevel) {
      Swal.fire({
        icon: 'error',
        title: 'Akses Ditolak 🛑',
        text: 'Halaman ini eksklusif untuk paket langganan di atas Anda. Silakan hubungi Owner untuk melakukan upgrade.',
        confirmButtonColor: '#ef4444',
        customClass: { popup: 'rounded-[32px]' }
      });
      
      return '/retail/pos/riwayat'; 
    }
  }

  return true;
});

export default router