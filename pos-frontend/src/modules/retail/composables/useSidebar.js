import Swal from 'sweetalert2';
import { onMounted, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import api from '../../../api.js';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const getImageUrl = (path) => {
	if (!path || path === 'null' || path === 'undefined' || path === '') return '';
	if (path.startsWith('http://') || path.startsWith('https://')) return path;
	return `${API_BASE_URL}${path}`;
};

export function useSidebar() {
	const router = useRouter();
	const route = useRoute();
	const sidebarOpen = ref(false);

	const openGroups = ref({
		stock: route.path.includes('master-produk') || route.path.includes('penerimaan-barang') || route.path.includes('stock-opname') || route.path.includes('retur-barang'),
		admin: route.path.startsWith('/retail/dashboard') || route.path.startsWith('/retail/karyawan') || route.path.includes('report'),
	});

	const toggleGroup = (group) => {
		openGroups.value[group] = !openGroups.value[group];
	};

	const user = ref({
		name: localStorage.getItem('name') || 'User',
		role: localStorage.getItem('role') || 'staff',
		storeName: localStorage.getItem('storeName') || 'POS UMKM',
		storeLogo: getImageUrl(localStorage.getItem('storeLogo')),
		foto_url: getImageUrl(localStorage.getItem('foto_url')),
	});

	// Handler fungsi bernama agar bisa di-remove tanpa menyebabkan Memory Leak
	const syncUserData = () => {
		user.value.name = localStorage.getItem('name') || 'User';
		user.value.foto_url = getImageUrl(localStorage.getItem('foto_url'));
	};

	const syncStoreData = () => {
		user.value.storeName = localStorage.getItem('storeName') || 'POS UMKM';
		user.value.storeLogo = getImageUrl(localStorage.getItem('storeLogo'));
	};

	onMounted(() => {
		syncUserData();
		syncStoreData();
		user.value.role = localStorage.getItem('role') || 'staff';

		window.addEventListener('storage', syncUserData);
		window.addEventListener('profile-updated', syncUserData);
		window.addEventListener('store-updated', syncStoreData);
	});

	onUnmounted(() => {
		window.removeEventListener('storage', syncUserData);
		window.removeEventListener('profile-updated', syncUserData);
		window.removeEventListener('store-updated', syncStoreData);
	});

	const logout = () => {
		Swal.fire({
			title: 'Sesi Operasional',
			text: 'Pilih tindakan untuk mengakhiri sesi Anda saat ini.',
			icon: 'question',
			showCancelButton: true,
			showDenyButton: true,
			confirmButtonText: 'Ganti Cabang',
			denyButtonText: 'Keluar Total',
			cancelButtonText: 'Batal',
			buttonsStyling: false,
			customClass: {
				popup: 'rounded-[32px] p-6 border border-slate-100 shadow-2xl',
				title: 'text-2xl font-black text-slate-800 tracking-tight',
				htmlContainer: 'text-sm font-bold text-slate-500 mt-2 mb-6',
				actions: 'flex flex-col sm:flex-row gap-3 w-full',
				confirmButton: 'w-full sm:flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-indigo-200 active:scale-95 order-1 sm:order-1',
				denyButton: 'w-full sm:flex-1 bg-rose-500 hover:bg-rose-600 text-white font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all shadow-xl shadow-rose-200 active:scale-95 order-2 sm:order-2',
				cancelButton: 'w-full sm:flex-1 bg-slate-100 hover:bg-slate-200 text-slate-500 hover:text-slate-700 font-black text-[10px] md:text-xs uppercase tracking-widest py-4 rounded-2xl transition-all active:scale-95 order-3 sm:order-3',
			},
		}).then(async (result) => {
			if (result.isConfirmed) {
				try {
					const resSetting = await api.get('/retail/store/settings');
					const freshData = resSetting.data.data;
					let tempStores = localStorage.getItem('temp_stores');

					if (tempStores) {
						let storesArr = JSON.parse(tempStores);
						storesArr = storesArr.map((store) => {
							if (store.id === freshData.id) {
								return { ...store, nama_toko: freshData.nama_toko, logo_url: freshData.logo_url, kota: freshData.kota, subscription_plan: freshData.subscription_plan };
							}
							return store;
						});
						localStorage.setItem('temp_stores', JSON.stringify(storesArr));
					}
				} catch (e) {
					console.error('Sinkronisasi list cabang ditunda:', e.message);
				}

				// Wipe clean status toko aktif sebelum transisi agar tidak nyangkut
				localStorage.removeItem('store_id');
				localStorage.removeItem('storeName');
				localStorage.removeItem('storeLogo');
				localStorage.removeItem('subscriptionPlan');

				router.push('/select-store');
			} else if (result.isDenied) {
				localStorage.clear();
				window.location.href = '/login'; // Gunakan Hard Reload untuk menjamin interceptor axios bersih total
			}
		});
	};

	return { route, sidebarOpen, openGroups, user, toggleGroup, logout };
}
