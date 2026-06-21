import Swal from 'sweetalert2';
import { createRouter, createWebHistory } from 'vue-router';
import { useLoading } from '../composables/useLoading.js';

import adminRoutes from '../modules/admin/router/admin_routes.js';
import retailRoutes from '../modules/retail/router/retail_routes.js';
import laundryRoutes from '../modules/jasa/laundry/router/laundry_routes.js';

const baseRoutes = [
    {
        path: '/',
        name: 'LandingPage',
        component: () => import('../views/public/LandingPage.vue'),
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/auth/Register.vue'),
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/auth/Login.vue'),
    },
    {
        path: '/forgot-password',
        name: 'ForgotPassword',
        component: () => import('../views/auth/ForgotPassword.vue'),
    },
    {
        path: '/select-verify',
        name: 'SelectVerify',
        component: () => import('../views/auth/SelectionVerify.vue'),
    },
    {
        path: '/verify-otp',
        name: 'VerifyOTP',
        component: () => import('../views/auth/VerifyOTP.vue'),
    },
    {
        path: '/reset-password',
        name: 'ResetPassword',
        component: () => import('../views/auth/ResetPassword.vue'),
    },
    {
        path: '/setup-toko',
        name: 'SetupToko',
        component: () => import('../views/onboarding/SetupToko.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/select-store',
        name: 'SelectStore',
        component: () => import('../views/auth/SelectStore.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/admin/login',
        name: 'LoginAdmin',
        component: () => import('../modules/admin/views/AdminLogin.vue'),
    },
];

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [...baseRoutes, ...adminRoutes, ...retailRoutes, ...laundryRoutes],
});

// 🛡️ SATPAM GLOBAL VUE ROUTER (FOUNDER/ENTERPRISE EDITION)
router.beforeEach((to, from) => {
    const { startLoading } = useLoading();
    startLoading();

    const token = localStorage.getItem('token');
    const userRole = localStorage.getItem('role') || 'staff';
    
    // 🚀 DYNAMIC DETECTOR: Lacak kasta industri user dari local storage bray!
    const rawIndustry = localStorage.getItem('pendingIndustry') || localStorage.getItem('user_industry') || 'retail';
    const userIndustry = rawIndustry.toLowerCase();

    // Helper penentu home rute dinamis kasir staff bray bray bray
    const getStaffHomeRoute = () => {
        if (userIndustry === 'jasa' || userIndustry === 'laundry') return '/laundry/pos/riwayat';
        if (userIndustry === 'fnb') return '/fnb/pos';
        return '/retail/pos';
    };

    const guestRoutes = ['/login', '/register', '/select-verify', '/verify-otp', '/admin/login'];

    // 🚀 1. LOGIKA ANTI-MUNDUR (GUEST GUARD - DYNAMIC REDIRECT)
    if (to.path === '/' && to.query.action === 'expansion') {
        // Biarkan lewat bray
    } else if (token && (guestRoutes.includes(to.path) || to.path === '/')) {
        if (userRole === 'super_admin') return '/admin/dashboard';
        if (userRole === 'owner') return '/select-store';
        return getStaffHomeRoute(); // ◄ FIX: Pengalihan dinamis anti-hardcode retail!
    }

    // 🚀 2. CEK AUTENTIKASI
    if (to.meta.requiresAuth && !token) {
        if (to.path.startsWith('/admin')) return '/admin/login';
        return '/login';
    }

	// 🚀 NEW GUARD: GEBREG CELAH ROOT ADMIN!
    if (to.meta.requiresAdmin && userRole !== 'super_admin') {
        Swal.fire({
            icon: 'error',
            title: 'Pelanggaran Otoritas! 🛑',
            text: 'Deteksi Akses Ilegal! Perangkat Anda tidak memiliki izin valid Root Admin platform.',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[32px]' },
        });
        return '/login';
    }

    // 🚀 3. CEK OTORISASI JABATAN
    if (to.meta.role === 'owner' && userRole !== 'owner') {
        Swal.fire({
            icon: 'error',
            title: 'Akses Ilegal!',
            text: 'Waduh, halaman ini cuma boleh dibuka sama Owner!',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[32px]' },
        });
        return getStaffHomeRoute(); // ◄ FIX: Mentalin staff liar balik ke jalurnya masing-masing
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
                customClass: { popup: 'rounded-[32px]' },
            });

            return getStaffHomeRoute(); // ◄ FIX: Balikin ruko ruko trial ke gerbang aslinya
        }
    }

    return true;
});

router.afterEach(() => {
    const { stopLoading } = useLoading();
    stopLoading();
});

export default router;