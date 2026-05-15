import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/auth/Login.vue';

const routes = [
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
        path: '/select-verify',
        name: 'SelectVerify',
        component: () => import('../views/auth/SelectionVerify.vue')
    },
    { 
        path: '/verify',
        name: 'Verify',
        component: () => import('../views/auth/VerifyEmail.vue')
    },
    { 
        path: '/login',
        name: 'Login',
        component: Login
    },
    { 
        path: '/dashboard', 
        component: () => import('../views/dashboard/Dashboard.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/pos/kasir', 
        component: () => import('../views/pos/Kasir.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/pos/buka-kasir', 
        component: () => import('../views/pos/BukaKasir.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/absensi', 
        component: () => import('../views/sdm/Absensi.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/schedule', 
        component: () => import('../views/sdm/Schedule.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/stock-opname', 
        component: () => import('../views/inventori/StockOpname.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/stock-opname/report',
        name: 'SrockOpnameReport', 
        component: () => import('../views/inventori/StockOpnameReport.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/riwayat', 
        component: () => import('../views/pos/Riwayat.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/produk', 
        component: () => import('../views/produk/MasterProduk.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/penerimaan-barang', 
        component: () => import('../views/produk/PenerimaanBarang.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/karyawan', 
        component: () => import('../views/sdm/ManageKaryawan.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/setup', 
        component: () => import('../views/setting/SetupToko.vue'),
        meta: { requiresAuth: true } 
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

// SATPAM FRONTEND (Route Guard)
router.beforeEach((to, from, next) => {
    const role = localStorage.getItem('role');
    const token = localStorage.getItem('token');
    
    // 1. Cek Login (Wajib ada token untuk semua rute kecuali login)
    if (to.meta.requiresAuth && !token) {
        return next('/login');
    }

    // 2. Logic khusus area POS
    if (to.path.startsWith('/pos')) {
        // Jika tujuannya ke mesin kasir (/pos/kasir), 
        // tapi dia belum inisialisasi modal, kita jangan kasih lewat!
        // Tapi pengecekan inisial modal ini paling akurat dilakukan di onMounted Kasir.vue
        // Untuk di router, kita kasih lewat dulu aja
        return next();
    }

    // 3. Proteksi Dashboard (Hanya Owner)
    if (to.path === '/dashboard' && role !== 'owner') {
        return next('/absensi');
    }

    next();
});

export default router;