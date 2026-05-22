export default [
  {
    path: '/retail/dashboard',
    name: 'RetailDashboard',
    // 🚀 Arahkan ke sub-folder views/dashboard/ tempat file DashboardRetail lu berada
    component: () => import('../views/dashboard/DashboardRetail.vue'), 
    meta: { requiresAuth: true, role: 'owner' }
  },
  {
    path: '/retail/stock-opname',
    name: 'RetailStockOpname',
    component: () => import('../views/inventori/StockOpname.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/pos',
    name: 'RetailPOS',
    component: () => import('../views/pos/PosRetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/retail/pos/buka-kasir',
    name: 'BukaKasir',
    component: () => import('../views/pos/BukaKasir.vue'), // arahkan tepat ke file lego anyar lu beb
    meta: { requiresAuth: true }
  },
]