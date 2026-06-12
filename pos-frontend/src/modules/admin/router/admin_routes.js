// 🚀 Sub-Router Lokal Khusus Domain Modul Admin Core SaaS

// Import Master Layout lu (Lazy load biar enteng )
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
			requiresGuest: true, // Opsional: Biar yg udah login mental balik ke dashboard
		},
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
			{
				path: '',
				redirect: '/admin/dashboard',
			},

			{
				path: 'dashboard',
				name: 'MissionControl',
				component: () => import('../views/MissionControl.vue'),
			},

			{
				path: 'analytics',
				name: 'Analytics',
				component: () => import('../views/analytics.vue'),
			},

			{
				path: 'tenant-hub',
				name: 'AdminTenantHub',
				component: () => import('../views/TenantManagementHub.vue'),
			},

			{
				path: 'audit',
				name: 'AuditLogs',
				component: () => import('../views/AuditLogs.vue'),
			},
		],
	},
];

export default adminRoutes;
