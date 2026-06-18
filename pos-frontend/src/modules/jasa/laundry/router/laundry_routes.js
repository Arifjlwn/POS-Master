// 🧺 SIRKUIT NAVIGASI LAUNDRY DENGAN PENGAWALAN KASTA LISENSI SAAS (ANTI-BYPASS)
export default [
	// =====================================================================
	// 💳 FITUR KASIR UTAMA & STATUS (BASIC - LEVEL 1)
	// =====================================================================
	// ✅ Boleh diakses Kasir/Staff ruko laundry bray
	{
		path: '/laundry/pos',
		name: 'LaundryPOS',
		component: () => import('../views/pos/PosLaundry.vue'),
		meta: { requiresAuth: true, minPlanLevel: 1 },
	},
	{
		path: '/laundry/pos/riwayat',
		name: 'LaundryRiwayat',
		component: () => import('../views/pos/RiwayatLaundry.vue'),
		meta: { requiresAuth: true, minPlanLevel: 1 },
	},
	{
		path: '/laundry/status',
		name: 'LaundryStatus',
		component: () => import('../views/StatusCucian.vue'),
		meta: { requiresAuth: true, minPlanLevel: 1 },
	},

	// =====================================================================
	// 🧺 MASTER KATALOG PRODUK (BASIC - LEVEL 1)
	// =====================================================================
	// 🔒 Khusus tingkat Pemilik (Owner) ruko
	{
		path: '/laundry/master-layanan',
		name: 'LaundryMasterLayanan',
		component: () => import('../views/produk/MasterLayanan.vue'),
		meta: { requiresAuth: true, role: 'owner', minPlanLevel: 1 },
	},

	// =====================================================================
	// 👥 TIM KARYAWAN & ABSENSI (PRO - LEVEL 2)
	// =====================================================================
	// 🔒 Hanya ruko paket PRO ke atas yang bisa manage karyawan
	{
		path: '/laundry/karyawan',
		name: 'LaundryKaryawan',
		component: () => import('../views/sdm/ManajemenKaryawan.vue'),
		meta: { requiresAuth: true, role: 'owner', minPlanLevel: 2 },
	},

	// =====================================================================
	// 📈 DASHBOARD & BI (PREMIUM - LEVEL 3)
	// =====================================================================
	// 🔒 Laporan finansial mendalam masuk kasta premium bray!
	{
		path: '/laundry/laporan',
		name: 'LaundryLaporan',
		component: () => import('../views/LaporanLaundry.vue'),
		meta: { requiresAuth: true, role: 'owner', minPlanLevel: 3 },
	},

	// =====================================================================
	// ⚙️ SEKTOR REGULASI SETTING (UNIVERSAL ROUTING CONNECTED)
	// =====================================================================
	// Info Dasar & Akun cukup di LEVEL 1
	{
		path: '/laundry/setting',
		name: 'LaundryStoreSetting',
		component: () => import('../views/settings/StoreSetting.vue'),
		meta: { requiresAuth: true, role: 'owner', minPlanLevel: 1 },
	},
	{
		path: '/laundry/akun',
		name: 'LaundryAccountSetting',
		component: () => import('../views/settings/AccountSetting.vue'),
		meta: { requiresAuth: true, minPlanLevel: 1 },
	},

	// 🚀 NEW: MANAJEMEN DENAH RAK BAJU (BASIC - LEVEL 1)
	// Owner butuh akses ini buat ngatur kapasitas lemari/rak ruko mereka
	{
		path: '/laundry/rak',
		name: 'LaundryManageRack',
		component: () => import('../views/settings/ManageRackView.vue'),
		meta: { requiresAuth: true, role: 'owner', minPlanLevel: 1 },
	},

	// 🔒 INTEGRASI WHATSAPP GATEWAY (PREMIUM - LEVEL 3)
	// Fitur kirim nota struk otomatis via Fonnte dikunci saklek di level tertinggi bray!
	{
		path: '/laundry/whatsapp',
		name: 'LaundryWhatsappSetting',
		component: () => import('../views/settings/WhatsappSetting.vue'),
		meta: { requiresAuth: true, role: 'owner', minPlanLevel: 3 },
	},
];
