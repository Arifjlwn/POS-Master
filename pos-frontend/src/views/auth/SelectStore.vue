<script setup>
import Swal from 'sweetalert2';
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';
import PricingPlan from '../../modules/retail/components/PricingPlan.vue';

const router = useRouter();
const stores = ref([]);
const userName = ref('');
const isLoading = ref(false);

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

// ==============================================================
// PENGATURAN NAMA BRAND APLIKASI
// ==============================================================
const BRAND_NAME = 'ARZURA';

// STATE UTAMA UNTUK MODAL
const showPricingModal = ref(false);

// 🚀 TIMPA HOOK ONMOUNTED DI SELECTSTORE.VUE LU PAKAI INI RIF!
onMounted(async () => {
	const name = localStorage.getItem('temp_name');
	userName.value = name || 'Owner';

	const tempStores = localStorage.getItem('temp_stores');
	stores.value = tempStores ? JSON.parse(tempStores) : [];

	try {
		const res = await api.get('/me');
		if (res.data && res.data.stores) {
			stores.value = res.data.stores;
			localStorage.setItem('temp_stores', JSON.stringify(res.data.stores));
		}
	} catch (err) {
		console.error('Gagal menyinkronkan data ruko realtime dari database:', err);
	}
});

// Urutkan toko berdasarkan ID bawaan
const sortedStores = computed(() => {
	return [...stores.value].sort((a, b) => a.id - b.id);
});

// Menentukan kondisi apakah user sedang melakukan "Ekspansi" (jika sudah punya toko aktif)
const isExpansionMode = computed(() => {
	// 🚀 FILTER CERDAS: Hanya dianggap ekspansi kalau sudah punya toko yang statusnya 'active' !
	return stores.value.some((store) => store.subscription_status === 'active' || store.subscription_status === 'ACTIVE');
});

const selectBranch = async (store) => {
	// 🚀 1. INTERCEPTOR TOKO PENDING (Penyelamat Transaksi Tertunda !)
	if (store.subscription_status === 'pending' || store.subscription_status === 'PENDING') {
		Swal.fire({
			title: 'Pembayaran Tertunda',
			text: `Gerai "${store.nama_toko}" belum diaktifkan. Ingin melanjutkan proses pembayaran gateway sekarang?`,
			icon: 'info',
			showCancelButton: true,
			confirmButtonColor: '#4f46e5',
			cancelButtonColor: '#64748b',
			confirmButtonText: 'Ya, Bayar Sekarang',
			cancelButtonText: 'Nanti Saja',
			customClass: { popup: 'rounded-[32px]' },
		}).then((result) => {
			if (result.isConfirmed) {
				// Simpan data cadangan ke localStorage agar halaman setup-toko tahu paket apa yang mau dibayar
				localStorage.setItem('pendingIndustry', store.industry || 'retail');
				localStorage.setItem('pendingPlan', store.subscription_plan || 'basic');

				// Lempar balik ke setup toko dengan membawa parameter ID Toko lama yang pending tadi !
				router.push(`/setup-toko?is_expansion=true&resume_store_id=${store.id}`);
			}
		});
		return; // Blokir eksekusi agar dia tidak menembak login kasir !
	}

	// 🚀 2. SELEKSI AKSES KASIR NORMAL (Jika status sudah 'active')
	isLoading.value = true;
	try {
		const res = await api.post('/auth/select-store', { store_id: store.id });

		localStorage.setItem('token', res.data.token);
		localStorage.setItem('store_id', res.data.store_id);
		localStorage.setItem('storeName', res.data.store_name || 'POS UMKM');
		localStorage.setItem('storeLogo', res.data.store_logo || res.data.logo_url || '');
		localStorage.setItem('subscriptionPlan', res.data.subscription_plan || 'basic');

		let finalRole = 'owner';
		if (res.data.role) {
			finalRole = res.data.role.toLowerCase();
		} else if (res.data.user && res.data.user.role) {
			finalRole = res.data.user.role.toLowerCase();
		} else {
			const savedRole = localStorage.getItem('role');
			if (savedRole) finalRole = savedRole.toLowerCase();
		}
		localStorage.setItem('role', finalRole);

		localStorage.setItem('name', res.data.name || res.data.user?.name || 'Owner');
		localStorage.setItem('foto_url', res.data.foto_url || res.data.user?.foto_url || '');

		router.push('/retail/pos/riwayat');
	} catch (error) {
		console.error('Gagal masuk sesi toko:', error);
		Swal.fire({
			icon: 'error',
			title: 'Akses Gagal',
			text: error.response?.data?.error || 'Gagal masuk ke toko ini.',
			confirmButtonColor: '#ef4444',
			customClass: { popup: 'rounded-[32px]' },
		});
	} finally {
		isLoading.value = false;
	}
};

const handlePilihPaket = (payload) => {
	localStorage.setItem('pendingIndustry', payload.industry);
	localStorage.setItem('pendingPlan', payload.plan);
	showPricingModal.value = false;
	router.push('/setup-toko?is_expansion=true');
};

const handleLogout = () => {
	Swal.fire({
		title: 'Keluar Aplikasi?',
		text: 'Anda harus memasukkan ulang kredensial untuk mengakses data kasir kembali.',
		icon: 'warning',
		showCancelButton: true,
		confirmButtonColor: '#4f46e5',
		cancelButtonColor: '#f43f5e',
		confirmButtonText: 'Ya, Keluar Sesi',
		cancelButtonText: 'Batal',
		customClass: {
			popup: 'rounded-[32px] p-6 md:p-8',
			confirmButton: 'rounded-xl px-5 py-3 text-xs font-black uppercase tracking-wider',
			cancelButton: 'rounded-xl px-5 py-3 text-xs font-black uppercase tracking-wider',
		},
	}).then((result) => {
		if (result.isConfirmed) {
			localStorage.clear();
			router.push('/login');
		}
	});
};

const getPlanStyle = (plan) => {
	const p = plan ? plan.toLowerCase() : 'basic';
	if (p === 'premium') return 'bg-amber-50 text-amber-700 border-amber-200 ring-amber-500 shadow-amber-100/50';
	if (p === 'pro') return 'bg-indigo-50 text-indigo-700 border-indigo-200 ring-indigo-500 shadow-indigo-100/50';
	return 'bg-sky-50 text-sky-700 border-sky-200 ring-sky-500 shadow-sky-100/50';
};

const cleanLogoUrl = (url) => {
	if (!url) return '';
	if (url.startsWith('http://') || url.startsWith('https://')) {
		return url;
	}
	let cleanPath = url.replace(/http:\/\/localhost:8080/g, '');
	return `${API_BASE_URL}${cleanPath}`;
};

const getStoreLabel = (store, index) => {
	// 🚀 JIKA PENDING: Berikan informasi tegas mendeteksi status belum diaktivasi !
	if (store.subscription_status === 'pending' || store.subscription_status === 'PENDING') {
		return 'MENUNGGU PEMBAYARAN';
	}
	if (index === 0) return 'Toko Utama';
	const urutanAngka = ['Kedua', 'Ketiga', 'Keempat', 'Kelima', 'Keenam', 'Ketujuh', 'Kedelapan', 'Kesembilan'];
	if (index - 1 < urutanAngka.length) {
		return `Cabang ${urutanAngka[index - 1]}`;
	}
	return `Cabang Ke-${index + 1}`;
};
</script>

<template>
	<div class="min-h-screen bg-[#FAFCFF] flex flex-col items-center py-6 md:py-12 px-4 md:px-8 font-sans relative overflow-x-hidden select-none">
		<div class="absolute top-[-10%] left-[-10%] w-[30rem] md:w-[50rem] h-[30rem] md:h-[50rem] bg-gradient-to-br from-indigo-200/30 to-purple-200/20 rounded-full filter blur-[100px] md:blur-[150px] pointer-events-none"></div>
		<div class="absolute bottom-[-10%] right-[-10%] w-[30rem] md:w-[50rem] h-[30rem] md:h-[50rem] bg-gradient-to-tr from-blue-200/30 to-sky-200/20 rounded-full filter blur-[100px] md:blur-[150px] pointer-events-none"></div>

		<div class="w-full max-w-6xl flex justify-end mb-6 md:mb-0 md:absolute md:top-6 md:right-12 z-20">
			<button @click="handleLogout" class="flex items-center gap-2.5 px-4 md:px-5 py-2.5 md:py-3 bg-white/80 backdrop-blur-md border border-rose-100 hover:border-rose-200 text-rose-600 hover:text-white hover:bg-rose-600 rounded-2xl font-black text-[10px] md:text-[11px] uppercase tracking-widest shadow-sm hover:shadow-xl hover:-translate-y-0.5 transition-all duration-300 w-full sm:w-auto justify-center">
				<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
				</svg>
				Keluar Sesi
			</button>
		</div>

		<div class="w-full max-w-6xl z-10 flex flex-col items-center animate-fade-in-up md:mt-8">
			<div class="text-center mb-10 md:mb-16 w-full px-2">
				<div class="inline-flex items-center justify-center px-4 md:px-5 py-2 md:py-2.5 bg-white rounded-2xl shadow-sm border border-slate-100/80 mb-4 md:mb-6 hover:scale-105 transition-transform duration-300">
					<div class="font-black text-xl md:text-2xl text-slate-900 tracking-tighter leading-none flex items-center gap-1.5 uppercase">
						{{ BRAND_NAME }}
						<span class="text-indigo-600 bg-indigo-50 px-2 md:px-2.5 py-0.5 md:py-1 rounded-xl border border-indigo-100/70 text-[10px] md:text-xs tracking-normal normal-case">POS</span>
					</div>
				</div>
				<h1 class="text-3xl md:text-6xl font-black text-slate-900 tracking-tight mb-4 leading-tight">
					Selamat Datang,
					<span class="bg-gradient-to-r from-indigo-600 to-blue-600 bg-clip-text text-transparent">{{ userName ? userName.split(' ')[0] : 'Owner' }}</span>
					!
				</h1>
				<p class="text-slate-400 font-bold text-[9px] md:text-[11px] uppercase tracking-widest bg-white border border-slate-100/60 px-4 md:px-6 py-2 md:py-2.5 rounded-full inline-block shadow-sm">
					{{ stores.length > 0 ? 'Silakan pilih gerai toko untuk memulai transaksi kasir' : 'Inisialisasi Pembuatan Toko Baru' }}
				</p>
			</div>

			<div v-if="stores.length > 0" class="w-full">
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 md:gap-8 w-full relative z-10 items-stretch">
					<div v-for="(store, idx) in sortedStores" :key="store.id" @click="selectBranch(store)" class="group bg-white rounded-[32px] md:rounded-[36px] p-6 md:p-7 border border-slate-100/80 shadow-sm hover:shadow-[0_25px_60px_-15px_rgba(99,102,241,0.15)] hover:-translate-y-2 transition-all duration-500 cubic-bezier-card cursor-pointer flex flex-col relative overflow-hidden justify-between">
						<div class="absolute top-0 left-0 w-full h-1.5 transition-all duration-500" :class="store.subscription_status === 'pending' ? 'bg-rose-500' : getPlanStyle(store.subscription_plan).includes('ring-amber-500') ? 'bg-amber-500' : getPlanStyle(store.subscription_plan).includes('ring-indigo-500') ? 'bg-indigo-500' : 'bg-sky-500'"></div>

						<div>
							<div class="flex items-start justify-between mb-6 md:mb-8">
								<div class="w-14 h-14 md:w-16 md:h-16 rounded-[20px] md:rounded-[22px] bg-slate-50 border border-slate-100 flex items-center justify-center overflow-hidden shrink-0 group-hover:scale-105 shadow-inner transition-all duration-500">
									<img v-if="store.logo_url" :src="cleanLogoUrl(store.logo_url)" class="w-full h-full object-cover" />
									<svg v-else class="w-6 h-6 md:w-7 h-7 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
									</svg>
								</div>

								<span :class="store.subscription_status === 'pending' ? 'bg-rose-50 text-rose-600 border-rose-200 font-extrabold' : 'bg-slate-50 border-slate-200/50'" class="text-[9px] md:text-[10px] font-black uppercase tracking-widest px-3 py-1 md:py-1.5 rounded-xl transition-all duration-300">
									{{ getStoreLabel(store, idx) }}
								</span>
							</div>

							<div class="mb-6 md:mb-8">
								<h3 class="text-xl md:text-2xl font-black text-slate-800 mb-1.5 md:mb-2 line-clamp-1 uppercase tracking-tight group-hover:text-indigo-600 transition-colors duration-300">
									{{ store.nama_toko }}
								</h3>
								<p class="text-[11px] md:text-xs font-bold text-slate-400 uppercase tracking-wider line-clamp-1 flex items-center gap-1.5">
									<svg class="w-3.5 h-3.5 text-indigo-500 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
										<path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
									</svg>
									{{ store.kota || 'Lokasi Belum Diatur' }}
								</p>
							</div>
						</div>

						<div class="mt-auto flex items-center justify-between border-t border-slate-50 pt-4 md:pt-5">
							<div class="flex flex-col">
								<span class="text-[8px] font-black text-slate-300 uppercase tracking-widest leading-none">Lisensi</span>
								<span class="text-[9px] md:text-[10px] font-black uppercase tracking-wider mt-1 px-2 py-0.5 border rounded-md" :class="getPlanStyle(store.subscription_plan)">
									{{ store.subscription_plan || 'Basic' }}
								</span>
							</div>

							<div :class="store.subscription_status === 'pending' ? 'text-rose-600 group-hover:text-rose-700' : 'text-slate-700 group-hover:text-indigo-600'" class="text-[10px] md:text-xs font-black uppercase tracking-widest flex items-center gap-1 transition-colors duration-300">
								{{ store.subscription_status === 'pending' ? 'Selesaikan Aktivasi' : 'Kelola Toko' }}
								<svg class="w-3.5 h-3.5 md:w-4 h-4 group-hover:translate-x-1.5 transition-transform duration-300 ease-out" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
									<path stroke-linecap="round" stroke-linejoin="round" d="M14 5l7 7m0 0l-7 7m7-7H3" />
								</svg>
							</div>
						</div>
					</div>

					<div @click="showPricingModal = true" class="bg-white rounded-[32px] md:rounded-[36px] p-6 md:p-7 border-2 border-dashed border-slate-200 hover:border-indigo-400 hover:bg-indigo-50/5 hover:-translate-y-1.5 transition-all duration-500 cubic-bezier-card cursor-pointer flex flex-col items-center justify-center text-center min-h-[200px] md:min-h-[230px] group shadow-sm">
						<div class="w-12 h-12 md:w-14 md:h-14 rounded-full bg-slate-50 border border-slate-200 group-hover:bg-indigo-100/50 group-hover:scale-105 transition-all duration-300 flex items-center justify-center mb-3 md:mb-4 shadow-sm">
							<svg class="w-5 h-5 md:w-6 h-6 text-slate-400 group-hover:text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
							</svg>
						</div>
						<h3 class="text-xs md:text-sm font-black text-slate-700 group-hover:text-indigo-800 uppercase tracking-widest mb-1 transition-colors">Buka Cabang Baru</h3>
						<p class="text-[11px] md:text-xs font-bold text-slate-400 px-4 md:px-6 leading-relaxed">Ekspansi bisnis Anda dengan mendaftarkan lokasi toko atau gerai baru.</p>
					</div>
				</div>
			</div>

			<div v-else class="w-full max-w-xl text-center py-12 md:py-16 px-6 md:px-8 bg-white rounded-[36px] md:rounded-[44px] border border-slate-100 shadow-xl">
				<div class="w-20 h-20 md:w-24 md:h-24 bg-indigo-50 border border-indigo-100 text-indigo-600 rounded-[28px] md:rounded-[32px] flex items-center justify-center mx-auto mb-6 md:mb-8 shadow-sm">
					<svg class="w-10 h-10 md:w-12 md:h-12 animate-pulse" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
					</svg>
				</div>
				<h2 class="text-2xl md:text-3xl font-black text-slate-800 uppercase tracking-tight mb-3">Belum Ada Toko Terdaftar</h2>
				<p class="text-slate-400 font-bold text-[11px] md:text-sm leading-relaxed mb-6 md:mb-8 px-4">Akun Anda telah aktif secara global. Silakan tentukan kluster industri untuk mengonfigurasi database toko pertama Anda.</p>
				<button @click="showPricingModal = true" class="w-full py-3.5 md:py-4 bg-indigo-600 hover:bg-slate-900 text-white rounded-2xl font-black text-xs uppercase tracking-widest shadow-xl shadow-indigo-100 transition-all duration-300">Buat Toko Pertama Anda</button>
			</div>
		</div>

		<Transition name="modal-fade">
			<div v-if="showPricingModal" class="fixed inset-0 z-[100] flex items-center justify-center p-3 sm:p-4 md:p-6 overflow-y-auto">
				<div @click="stores.length > 0 ? (showPricingModal = false) : null" class="fixed inset-0 bg-slate-950/50 backdrop-blur-md transition-opacity duration-300"></div>

				<div class="bg-white w-full max-w-6xl rounded-[32px] md:rounded-[44px] shadow-2xl p-5 sm:p-6 md:p-10 relative border border-slate-100 my-auto max-h-[92vh] overflow-y-auto custom-scrollbar flex flex-col items-center z-10 transform transition-all duration-500 cubic-bezier-modal">
					<PricingPlan :is-expansion="isExpansionMode" :show-close="stores.length > 0" @close="showPricingModal = false" @select-plan="handlePilihPaket" />
				</div>
			</div>
		</Transition>

		<div v-if="isLoading" class="fixed inset-0 z-[150] bg-slate-950/40 backdrop-blur-sm flex items-center justify-center p-4">
			<div class="bg-white p-6 md:p-8 rounded-[24px] md:rounded-[28px] shadow-2xl border border-slate-100 flex flex-col items-center max-w-xs w-full text-center">
				<div class="w-11 h-11 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
				<div class="text-[10px] font-black text-slate-600 uppercase tracking-widest animate-pulse">Memuat Toko</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
/* CSS Tetap Mewah Konsisten */
.cubic-bezier-card {
	transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}
.cubic-bezier-modal {
	transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}
.animate-fade-in-up {
	animation: fadeInUp 0.6s cubic-bezier(0.25, 1, 0.5, 1) both;
}
@keyframes fadeInUp {
	from {
		opacity: 0;
		transform: translateY(16px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}
.modal-fade-enter-from {
	opacity: 0;
}
.modal-fade-enter-from .cubic-bezier-modal {
	opacity: 0;
	transform: scale(0.96) translateY(12px);
}
.modal-fade-enter-active,
.modal-fade-leave-active {
	transition: all 0.35s ease-out;
}
.modal-fade-leave-to {
	opacity: 0;
}
.modal-fade-leave-to .cubic-bezier-modal {
	opacity: 0;
	transform: scale(0.97) translateY(8px);
}
.custom-scrollbar::-webkit-scrollbar {
	width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #e2e8f0;
	border-radius: 10px;
}
</style>
