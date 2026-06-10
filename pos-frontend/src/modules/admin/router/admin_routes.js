// 🚀 Sub-Router Lokal Khusus Domain Modul Admin Core SaaS

// Import Master Layout lu (Lazy load biar enteng bray)
const AdminLayout = () => import('../components/layout/AdminLayout.vue');

const adminRoutes = [
    // ==========================================
    // 1. RUTE GUEST (TANPA LAYOUT)
    // ==========================================
    {
        path: '/admin/login',
        name: 'AdminLogin',
        component: () => import('../views/AdminLogin.vue'),
        meta: { 
            requiresGuest: true // Opsional: Biar yg udah login mental balik ke dashboard
        }
    },

    // ==========================================
    // 2. RUTE TERPROTEKSI (DIBUNGKUS ADMIN LAYOUT)
    // ==========================================
    {
        path: '/admin',
        component: AdminLayout,
        // 🔒 Gembok ganda: Wajib login (token ada) & Wajib punya pangkat super_admin di DB!
        meta: {
            requiresAuth: true,
            requiresAdmin: true,
        },
        children: [
            // Redirect otomatis dari /admin ke /admin/dashboard
            {
                path: '',
                redirect: '/admin/dashboard'
            },
            
            // 🚀 Halaman Telemetri Mission Control (URL: /admin/dashboard)
            {
                path: 'dashboard',
                name: 'MissionControl',
                component: () => import('../views/MissionControl.vue'),
            },

            // 💡 Nanti kalau lu mau nambahin halaman Tenant Management, Billing, dll, selipin di bawah sini bray!
            // {
            //     path: 'tenants',
            //     name: 'TenantManagement',
            //     component: () => import('../views/TenantManagement.vue'),
            // },
        ]
    }
];

export default adminRoutes;