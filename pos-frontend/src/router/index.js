import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/auth/Login.vue';

const routes = [
    // ==========================================
    // MENU UTAMA WEBSITE (GLOBAL)
    // ==========================================
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
        path: '/setup', 
        component: () => import('../views/setting/SetupToko.vue'),
        meta: { requiresAuth: true } 
    },

    // ==========================================
    // --- TIPE BISNIS: RETAIL ---
    // ==========================================
    { 
        path: '/retail/dashboard', 
        component: () => import('../views/retail/dashboard/Dashboard.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/pos', 
        component: () => import('../views/retail/pos/Kasir.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/pos/buka-kasir', 
        component: () => import('../views/retail/pos/BukaKasir.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/absensi', 
        component: () => import('../views/retail/sdm/Absensi.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/schedule', 
        component: () => import('../views/retail/sdm/Schedule.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/stock-opname', 
        component: () => import('../views/retail/inventori/StockOpname.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/stock-opname/report',
        name: 'StockOpnameReport', 
        component: () => import('../views/retail/inventori/StockOpnameReport.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/riwayat', 
        component: () => import('../views/retail/pos/Riwayat.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/master-produk', 
        component: () => import('../views/retail/produk/MasterProduk.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/penerimaan-barang', 
        component: () => import('../views/retail/produk/PenerimaanBarang.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/retur-barang', 
        component: () => import('../views/retail/produk/ReturBarang.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/retur-barang/report', 
        component: () => import('../views/retail/produk/ReturReport.vue'),
        meta: { requiresAuth: true } 
    },
    { 
        path: '/retail/karyawan', 
        component: () => import('../views/retail/sdm/ManageKaryawan.vue'),
        meta: { requiresAuth: true } 
    },

    // ==========================================
    // --- TIPE BISNIS: LAYANAN & JASA ---
    // ==========================================

    // ==========================================
    //            --- LAUNDRY ---
    // ==========================================
    {
        path: '/laundry/pos',
        name: 'PosLaundry',
        component: () => import('../views/layanan-jasa/laundry/PosLaundry.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/laundry/master-layanan',
        name: 'MasterLayanan',
        component: () => import('../views/layanan-jasa/laundry/MasterLayanan.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/laundry/status',
        name: 'StatusCucian',
        component: () => import('../views/layanan-jasa/laundry/StatusCucian.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/laundry/laporan',
        name: 'Laporan',
        component: () => import('../views/layanan-jasa/laundry/LaporanLaundry.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/laundry/setting',
        name: 'Setting',
        component: () => import('../views/layanan-jasa/laundry/SettingToko.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/laundry/karyawan',
        name: 'ManajemenKasirLaundry',
        component: () => import('../views/layanan-jasa/laundry/ManajemenKasir.vue'),
        meta: { requiresAuth: true }
    },

    // ==========================================
    // --- TIPE BISNIS: Food & Beverage ---
    // ==========================================
    {
        path: '/fnb/kasir',
        name: 'PosFnB',
        component: () => import('../views/fnb/KasirFnB.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/fnb/dapur',
        name: 'DapurFnB',
        component: () => import('../views/fnb/DapurFnB.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/fnb/laporan',
        name: 'LaporanFnB',
        component: () => import('../views/fnb/LaporanFnB.vue'),
        meta: { requiresAuth: true }
    },

];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

// ==========================================
// SATPAM FRONTEND (Route Guard)
// ==========================================
router.beforeEach((to, from, next) => {
    const role = localStorage.getItem('role');
    const token = localStorage.getItem('token');
    
    // 1. Cek Login (Wajib ada token untuk semua rute kecuali public)
    if (to.meta.requiresAuth && !token) {
        return next('/login');
    }

    // 2. Logic khusus area POS (Pake .includes biar bisa baca /retail/pos atau /laundry/pos)
    if (to.path.includes('/pos')) {
        return next();
    }

    // 3. Proteksi Dashboard (Hanya Owner) - Pake .includes biar jalan di semua tipe bisnis
    if (to.path.includes('/dashboard') && role !== 'owner') {
        // Kalau karyawan nyasar ke dashboard, lempar sesuai tipe bisnisnya!
        if (to.path.includes('/retail')) {
            return next('/retail/absensi');
        } else if (to.path.includes('/laundry')) {
            return next('/laundry/pos'); // Contoh lemparan kalau karyawan laundry
        } else {
            return next('/setup'); // Fallback
        }
    }

    next();
});

export default router;