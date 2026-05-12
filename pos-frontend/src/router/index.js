import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';

const routes = [
    { path: '/', redirect: '/login' },
    { 
        path: '/login',
        name: 'Login',
        component: Login
    },
    { 
        path: '/dashboard', 
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/kasir', 
        component: () => import('../views/Kasir.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/produk', 
        component: () => import('../views/MasterProduk.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/setup', 
        component: () => import('../views/SetupToko.vue'),
        meta: { requiresAuth: true } 
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

// SATPAM FRONTEND (Route Guard)
router.beforeEach((to, from, next) => {
    const isAuthenticated = localStorage.getItem('token');

    // Kalau mau ke halaman yang dikunci tapi belum punya token, tendang ke login!
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login');
    } 
    // Kalau udah login tapi iseng buka halaman login lagi, arahkan ke dashboard!
    else if (to.path === '/login' && isAuthenticated) {
        next('/dashboard');
    } 
    else {
        next(); // Silakan lewat
    }
});

export default router;