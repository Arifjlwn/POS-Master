import Swal from 'sweetalert2';
import { createRouter, createWebHistory } from 'vue-router';

// import fnbRoutes from '@/modules/fnb/router/index.js'
import retailRoutes from '../modules/retail/router/retail_routes.js';
// import laundryRoutes from '@/modules/jasalayanan/laundry/router/laundry_routes.js'

const baseRoutes = [
    {
        path: '/',
        name: 'LandingPage',
        component: () => import('../views/public/LandingPage.vue')
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/auth/Register.vue')
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/auth/Login.vue')
    },
    {
        path: '/forgot-password',
        name: 'ForgotPassword',
        component: () => import('../views/auth/ForgotPassword.vue')
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
        path: '/reset-password',
        name: 'ResetPassword',
        component: () => import('../views/auth/ResetPassword.vue')
    },
    {
        path: '/setup-toko',
        name: 'SetupToko',
        component: () => import('../views/onboarding/SetupToko.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/select-store',
        name: 'SelectStore',
        component: () => import('../views/auth/SelectStore.vue'),
        meta: { requiresAuth: true }
    }
];

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        ...baseRoutes,
        // ...fnbRoutes,
        ...retailRoutes
        // ...laundryRoutes
    ]
});

// 🛡️ SATPAM GLOBAL VUE ROUTER
router.beforeEach((to, from) => {
    const token = localStorage.getItem('token');
    const userRole = localStorage.getItem('role') || 'staff';

    // 🚀 1. LOGIKA ANTI-MUNDUR (GUEST GUARD)
    const guestRoutes = ['/login', '/register', '/select-verify', '/verify-otp'];

    // Kalau dia mau balik ke landing page, TAPI niatnya buat EKSPANSI, izinkan!
    if (to.path === '/' && to.query.action === 'expansion') {
        // Biarkan lewat (jangan return apa-apa di block ini)
    }
    // Kalau dia mau ke landing page atau route guest lain tapi BUKAN ekspansi, lempar!
    else if (token && (guestRoutes.includes(to.path) || to.path === '/')) {
        // Kalau dia owner lempar ke select store aja (biar lebih rapi)
        if (userRole === 'owner') return '/select-store';
        // Kalau kasir, lempar ke POS
        return '/retail/pos';
    }

    // 🚀 2. CEK AUTENTIKASI (HARUS LOGIN)
    if (to.meta.requiresAuth && !token) {
        return '/login';
    }

    // 🚀 3. CEK OTORISASI JABATAN (KASIR VS OWNER)
    if (to.meta.role === 'owner' && userRole !== 'owner') {
        Swal.fire({
            icon: 'error',
            title: 'Akses Ilegal!',
            text: 'Waduh, halaman ini cuma boleh dibuka sama Owner!',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[32px]' }
        });
        return '/retail/pos/riwayat';
    }

    // 🚀 4. LOGIKA SATPAM SAAS (SISTEM KASTA LEVEL)
    if (to.meta.minPlanLevel) {
        const subPlan = localStorage.getItem('subscriptionPlan') || 'basic';

        const getPlanLevel = (plan) => {
            const p = plan.toLowerCase();
            if (p === 'premium' || p === 'trial') return 3;
            if (p === 'pro') return 2;
            return 1;
        };

        const currentLevel = getPlanLevel(subPlan);

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

export default router;
