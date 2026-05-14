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
        component: () => import('../views/Absensi.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/schedule', 
        component: () => import('../views/Schedule.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/stock-opname', 
        component: () => import('../views/StockOpname.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/stock-opname/report',
        name: 'SrockOpnameReport', 
        component: () => import('../views/StockOpnameReport.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/riwayat', 
        component: () => import('../views/Riwayat.vue'),
        meta: { requiresAuth: true } // Kunci pintu!
    },
    { 
        path: '/produk', 
        component: () => import('../views/MasterProduk.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/penerimaan-barang', 
        component: () => import('../views/PenerimaanBarang.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/karyawan', 
        component: () => import('../views/ManageKaryawan.vue'),
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
        return next('/pos/buka-kasir');
    }

    next();
});

export default router;