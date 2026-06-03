export default [
    // === DASHBOARD (PREMIUM - LEVEL 3) ===
    {
        path: "/retail/dashboard",
        name: "RetailDashboard",
        component: () => import("../views/dashboard/DashboardRetail.vue"),
        meta: { requiresAuth: true, role: "owner", minPlanLevel: 3 }, // 🔒 OWNER ONLY
    },

    // === Stock Opname (PREMIUM - LEVEL 3) ===
    // ✅ DIBUKA UNTUK KARYAWAN/KASIR
    {
        path: "/retail/stock-opname",
        name: "StockOpname",
        component: () => import("../views/inventori/StockOpname.vue"),
        meta: { requiresAuth: true, minPlanLevel: 3 },
    },
    {
        path: "/retail/stock-opname/report",
        name: "StockOpnameReport",
        component: () => import("../views/inventori/StockOpnameReport.vue"),
        meta: { requiresAuth: true, minPlanLevel: 3 },
    },

    // === POS (BASIC - LEVEL 1) ===
    // ✅ BEBAS AKSES BUAT KASIR
    {
        path: "/retail/pos",
        name: "RetailPOS",
        component: () => import("../views/pos/PosRetail.vue"),
        meta: { requiresAuth: true, minPlanLevel: 1 },
    },
    {
        path: "/retail/pos/buka-kasir",
        name: "BukaKasir",
        component: () => import("../views/pos/BukaKasir.vue"),
        meta: { requiresAuth: true, minPlanLevel: 1 },
    },
    {
        path: "/retail/pos/riwayat",
        name: "Riwayat",
        component: () => import("../views/pos/Riwayat.vue"),
        meta: { requiresAuth: true, minPlanLevel: 1 },
    },

    // === MASTER PRODUK (BASIC - LEVEL 1) ===
    {
        path: "/retail/produk/master-produk",
        name: "MasterProduk",
        component: () => import("../views/produk/MasterProduk.vue"),
        meta: { requiresAuth: true, role: "owner", minPlanLevel: 1 }, // 🔒 OWNER ONLY
    },

    // ✅ DIBUKA UNTUK KARYAWAN/KASIR (Terima Barang)
    {
        path: "/retail/produk/penerimaan-barang",
        name: "PenerimaanBarang",
        component: () => import("../views/produk/PenerimaanBarang.vue"),
        meta: { requiresAuth: true, minPlanLevel: 1 },
    },

    // === RETUR (PREMIUM - LEVEL 3) ===
    // ✅ DIBUKA UNTUK KARYAWAN/KASIR
    {
        path: "/retail/produk/retur-barang",
        name: "ReturBarang",
        component: () => import("../views/produk/ReturBarang.vue"),
        meta: { requiresAuth: true, minPlanLevel: 3 },
    },
    {
        path: "/retail/produk/retur-barang/report",
        name: "ReportReturBarang",
        component: () => import("../views/produk/ReturReport.vue"),
        meta: { requiresAuth: true, minPlanLevel: 3 },
    },

    // === SDM (PRO - LEVEL 2) ===
    // ✅ BEBAS AKSES BUAT KASIR (Absen & Jadwal)
    {
        path: "/retail/sdm/absensi",
        name: "Absensi",
        component: () => import("../views/sdm/Absensi.vue"),
        meta: { requiresAuth: true, minPlanLevel: 2 },
    },
    {
        path: "/retail/sdm/schedule",
        name: "Schedule",
        component: () => import("../views/sdm/Schedule.vue"),
        meta: { requiresAuth: true, minPlanLevel: 2 },
    },

    // 🔒 OWNER ONLY (Tambah Karyawan / Gaji)
    {
        path: "/retail/sdm/karyawan",
        name: "Karyawan",
        component: () => import("../views/sdm/ManageEmployee.vue"),
        meta: { requiresAuth: true, role: "owner", minPlanLevel: 2 },
    },

    // === PENGATURAN TOKO (BASIC - LEVEL 1) ===
    {
        path: "/retail/settings",
        name: "StoreSettings",
        component: () => import("../views/settings/StoreSetting.vue"),
        meta: { requiresAuth: true, role: "owner", minPlanLevel: 1 }, // 🔒 OWNER ONLY
    },
    {
        path: "/retail/account",
        name: "AccountSettings",
        component: () => import("../views/settings/AccountSetting.vue"),
        meta: { requiresAuth: true, minPlanLevel: 1 },
    },

    // === FITUR DEWA (PREMIUM - LEVEL 3) ===
    {
        path: "/retail/settings/whatsapp",
        name: "WhatsappSettings",
        component: () => import("../views/settings/WhatsappSetting.vue"),
        meta: { requiresAuth: true, role: "owner", minPlanLevel: 3 }, // 🔒 OWNER ONLY
    },
];
