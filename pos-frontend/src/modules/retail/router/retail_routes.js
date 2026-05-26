export default [
  // === DASHBOARD ===
  {
    path: '/retail/dashboard',
    name: 'RetailDashboard',
    // 🚀 Arahkan ke sub-folder views/dashboard/ tempat file DashboardRetail lu berada
    component: () => import('../views/dashboard/DashboardRetail.vue'), 
    meta: { requiresAuth: true, role: 'owner' }
  },

  // === INVENTORI ===
  {
    path: '/retail/inventori/stock-opname',
    name: 'StockOpname',
    component: () => import('../views/inventori/StockOpname.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/inventori/report',
    name: 'StockOpnameReport',
    component: () => import('../views/inventori/StockOpnameReport.vue'),
    meta: { requiresAuth: true }
  },

  // === POS ===
  {
    path: '/retail/pos',
    name: 'RetailPOS',
    component: () => import('../views/pos/PosRetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/pos/buka-kasir',
    name: 'BukaKasir',
    component: () => import('../views/pos/BukaKasir.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/pos/riwayat',
    name: 'Riwayat',
    component: () => import('../views/pos/Riwayat.vue'),
    meta: { requiresAuth: true }
  },

  // === MASTER PRODUK ===
  {
    path: '/retail/produk/master-produk',
    name: 'MasterProduk',
    component: () => import('../views/produk/MasterProduk.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/produk/penerimaan-barang',
    name: 'PenerimaanBarang',
    component: () => import('../views/produk/PenerimaanBarang.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/produk/retur-barang',
    name: 'ReturBarang',
    component: () => import('../views/produk/ReturBarang.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/produk/retur-barang/report',
    name: 'ReportReturBarang',
    component: () => import('../views/produk/ReturReport.vue'),
    meta: { requiresAuth: true }
  },

  // === SDM ===
  {
    path: '/retail/sdm/absensi',
    name: 'Absensi',
    component: () => import('../views/sdm/Absensi.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/sdm/schedule',
    name: 'Schedule',
    component: () => import('../views/sdm/Schedule.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/sdm/karyawan',
    name: 'Karyawan',
    component: () => import('../views/sdm/ManageKaryawan.vue'),
    meta: { requiresAuth: true }
  },
]