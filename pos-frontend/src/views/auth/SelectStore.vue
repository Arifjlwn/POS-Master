<script setup>
import Swal from 'sweetalert2';
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';

const router = useRouter();
const stores = ref([]);
const userName = ref('');
const isLoading = ref(false);

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

// ==============================================================
// STATE & DATA REAKTIF MODAL PRICING
// ==============================================================
const showPricingModal = ref(false);
const activePricingTab = ref('retail');

const industries = [
	{
		id: 'retail',
		title: 'Retail & Distribusi',
		desc: 'Supermarket, Minimarket, Butik, Toko Kelontong',
	},
	{
		id: 'fnb',
		title: 'Food & Beverage',
		desc: 'Cafe, Restoran, Kedai Kopi, Rumah Makan',
	},
	{
		id: 'jasa',
		title: 'Layanan & Jasa',
		desc: 'Laundry, Barbershop, Salon, Bengkel',
	},
];

// SINKRONISASI ENTERPRISE: Menyelaraskan manifest teks fitur dengan gembok middleware backend Go
const pricingPlans = {
	retail: [
		{
			id: 'trial',
			name: 'Starter Trial',
			price: 'Rp 0',
			duration: '14 Hari',
			desc: 'Uji coba gratis untuk melihat kesesuaian sistem dengan toko Anda.', features: ['Semua Fitur Unlocked']
		},
		{
			id: 'basic',
			name: 'Retail Basic',
			price: '49k',
			duration: '/Bulan',
			desc: 'Solusi tepat untuk pengelolaan satu toko skala kecil.',
			features: ['POS Kasir', 'Manajemen Stok & LPB', 'Struk Thermal Bluetooth', 'Riwayat Struk'],
		},
		{
			id: 'pro',
			name: 'Retail Pro',
			price: '149k',
			duration: '/Bulan',
			desc: 'Cocok untuk toko ritel yang mulai mengelola staf operasional.',
			features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Pengaturan Shift', 'Laporan Ekspor (CSV/Excel)'],
		},
		{
			id: 'premium',
			name: 'Retail Premium',
			price: '299k',
			duration: '/Bulan',
			desc: 'Kendali penuh untuk bisnis skala besar & audit inventaris.',
			features: ['Semua Fitur Pro', 'Audit Stock Opname', 'Dashboard Analitik', 'Notifikasi WhatsApp System'],
		},
	],
	fnb: [
		{
			id: 'trial',
			name: 'Starter Trial',
			price: 'Rp 0',
			duration: '14 Hari',
			desc: 'Uji coba gratis modul resto untuk kelancaran pesanan dapur.',
			features: ['Semua Fitur Unlocked'],
		},
		{
			id: 'basic',
			name: 'F&B Basic',
			price: '59k',
			duration: '/Bulan',
			desc: 'Sistem operasional praktis untuk kedai atau coffee shop.',
			features: ['Manajemen Layout Meja', 'Cetak Tiket Dapur (Kitchen)', 'Pajak & Service Charge', 'Struk Thermal Bluetooth'],
		},
		{
			id: 'pro',
			name: 'F&B Pro',
			price: '169k',
			duration: '/Bulan',
			desc: 'Ideal untuk restoran yang butuh kontrol manajemen staf teratur.',
			features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Shift Kerja', 'Split Bill & Gabung Meja'],
		},
		{
			id: 'premium',
			name: 'F&B Premium',
			price: '349k',
			duration: '/Bulan',
			desc: 'Skalabilitas bisnis franchise dengan laporan analitik terpusat.',
			features: ['Semua Fitur Pro', 'Resep Bahan Baku (BOM)', 'Self-Order QR Menu', 'Notifikasi WhatsApp System'],
		},
	],
	jasa: [
		{
			id: 'trial',
			name: 'Starter Trial',
			price: 'Rp 0',
			duration: '14 Hari',
			desc: 'Coba modul layanan untuk bengkel, salon, atau laundry.',
			features: ['Semua Fitur Unlocked'],
		},
		{
			id: 'basic',
			name: 'Service Basic',
			price: '49k',
			duration: '/Bulan',
			desc: 'Sistem tracking pesanan yang rapi untuk bisnis jasa kecil.',
			features: ['Tracking Status Pesanan', 'Cetak Nota / Resi Barcode', 'Manajemen Layanan & Tarif', 'Laporan Pendapatan'],
		},
		{
			id: 'pro',
			name: 'Service Pro',
			price: '159k',
			duration: '/Bulan',
			desc: 'Sistem otomatisasi performa operasional tim staf.',
			features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Shift Kerja', 'Laporan Kinerja Bulanan'],
		},
		{
			id: 'premium',
			name: 'Service Premium',
			price: '329k',
			duration: '/Bulan',
			desc: 'Manajemen booking tingkat lanjut dengan pengingat otomatis.',
			features: ['Semua Fitur Pro', 'Dashboard Analitik', 'Sistem Booking Reservasi', 'Notifikasi WhatsApp System'],
		},
	],
};

onMounted(() => {
	const tempStores = localStorage.getItem('temp_stores');
	const name = localStorage.getItem('temp_name');

	if (!tempStores && !name) {
		router.push('/login');
		return;
	}

	stores.value = tempStores ? JSON.parse(tempStores) : [];
	userName.value = name || 'Owner';
});

const sortedStores = computed(() => {
	return [...stores.value].sort((a, b) => a.id - b.id);
});

const filteredPlans = computed(() => {
	const currentPlans = pricingPlans[activePricingTab.value] || [];
	if (stores.value.length > 0) {
		return currentPlans.filter((plan) => plan.id !== 'trial');
	}
	return currentPlans;
});

const selectBranch = async (storeId) => {
	isLoading.value = true;
	try {
		const res = await api.post('/auth/select-store', { store_id: storeId });

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

const handlePilihPaketEkspansi = (industry, planId) => {
	localStorage.setItem('pendingIndustry', industry);
	localStorage.setItem('pendingPlan', planId);
	showPricingModal.value = false;
	router.push('/setup-toko?is_expansion=true');
};

const getPlanStyle = (plan) => {
	const p = plan ? plan.toLowerCase() : 'basic';
	if (p === 'premium') return 'bg-amber-50 text-amber-700 border-amber-200 ring-amber-500 text-amber-500';
	if (p === 'pro') return 'bg-indigo-50 text-indigo-700 border-indigo-200 ring-indigo-500 text-indigo-500';
	return 'bg-sky-50 text-sky-700 border-sky-200 ring-sky-500 text-sky-500';
};

const cleanLogoUrl = (url) => {
	if (!url) return '';
	if (url.startsWith('http://') || url.startsWith('https://')) {
		return url;
	}
	let cleanPath = url.replace(/http:\/\/localhost:8080/g, '');
	return `${API_BASE_URL}${cleanPath}`;
};

const getStoreLabel = (index) => {
	if (index === 0) return 'Toko Pertama (Utama)';
	const urutanAngka = ['Kedua', 'Ketiga', 'Keempat', 'Kelima', 'Keenam', 'Ketujuh', 'Kedelapan', 'Kesembilan'];
	if (index - 1 < urutanAngka.length) {
		return `Cabang ${urutanAngka[index - 1]}`;
	}
	return `Cabang Ke-${index + 1}`;
};
</script>

<template>
	<div class="min-h-screen bg-[#F8FAFC] flex flex-col items-center py-12 px-4 md:px-8 font-sans relative overflow-x-hidden select-none">
		<div class="absolute top-[-15%] left-[-15%] w-[45rem] h-[45rem] bg-gradient-to-br from-indigo-300/20 to-purple-300/10 rounded-full filter blur-[140px] pointer-events-none"></div>
		<div class="absolute bottom-[-15%] right-[-15%] w-[45rem] h-[45rem] bg-gradient-to-tr from-blue-300/20 to-sky-300/10 rounded-full filter blur-[140px] pointer-events-none"></div>

		<div class="w-full max-w-6xl z-10 flex flex-col items-center animate-fade-in-up">
			<div class="text-center mb-16">
				<div class="inline-flex items-center justify-center px-4 py-2 bg-white rounded-2xl shadow-xs border border-slate-100 mb-6 hover:scale-105 transition-transform duration-300">
					<div class="font-black text-2xl text-slate-900 tracking-tighter leading-none flex items-center gap-1.5">
						NEXA
						<span class="text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-lg border border-indigo-100">POS</span>
					</div>
				</div>
				<h1 class="text-3xl md:text-5xl font-black text-slate-800 tracking-tight mb-4">Selamat Datang, {{ userName.split(' ')[0] }}!</h1>
				<p class="text-slate-400 font-bold text-xs md:text-sm uppercase tracking-wider bg-white border border-slate-100 px-5 py-2 rounded-full inline-block shadow-3xs">
					{{ stores.length > 0 ? 'Silakan pilih gerai toko untuk memulai transaksi kasir' : 'Inisialisasi Pembuatan Toko Baru' }}
				</p>
			</div>

			<div v-if="stores.length > 0" class="w-full">
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full relative z-10 items-stretch">
					<div v-for="(store, idx) in sortedStores" :key="store.id" @click="selectBranch(store.id)" class="group bg-white rounded-[32px] p-6 border border-slate-100 shadow-2xs hover:shadow-2xl hover:-translate-y-2 transition-all duration-500 cubic-bezier-card cursor-pointer flex flex-col relative overflow-hidden justify-between">
						<div class="absolute top-0 left-0 w-full h-1.5 transition-all duration-500" :class="getPlanStyle(store.subscription_plan).split(' ')[3].replace('ring-', 'bg-')"></div>

						<div>
							<div class="flex items-start justify-between mb-6">
								<div class="w-14 h-14 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-center overflow-hidden shrink-0 group-hover:scale-105 transition-all duration-500">
									<img v-if="store.logo_url" :src="cleanLogoUrl(store.logo_url)" class="w-full h-full object-cover" />
									<svg v-else class="w-6 h-6 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
									</svg>
								</div>
								<span class="text-[10px] font-black uppercase tracking-widest bg-slate-50 border border-slate-200/60 px-3 py-1.5 rounded-xl transition-all duration-300 group-hover:bg-indigo-50 group-hover:text-indigo-600 group-hover:border-indigo-100">
									{{ getStoreLabel(idx) }}
								</span>
							</div>

							<div class="mb-6">
								<h3 class="text-xl font-black text-slate-800 mb-1.5 line-clamp-1 uppercase tracking-tight group-hover:text-indigo-600 transition-colors duration-300">
									{{ store.nama_toko }}
								</h3>
								<p class="text-xs font-bold text-slate-400 uppercase tracking-wider line-clamp-1 flex items-center gap-1">
									<svg class="w-3.5 h-3.5 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
										<path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
									</svg>
									{{ store.kota || 'Lokasi Belum Diatur' }}
								</p>
							</div>
						</div>

						<div class="mt-auto flex items-center justify-between border-t border-slate-50/80 pt-4">
							<div class="flex flex-col">
								<span class="text-[8px] font-black text-slate-300 uppercase tracking-widest leading-none">Status</span>
								<span class="text-[10px] font-black text-emerald-500 uppercase tracking-wider mt-0.5 flex items-center gap-1">
									<span class="w-2 h-2 bg-emerald-500 rounded-full relative flex">
										<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
										<span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
									</span>
									Kasir Ready
								</span>
							</div>
							<div class="text-xs font-black text-slate-700 group-hover:text-indigo-600 uppercase tracking-widest flex items-center gap-0.5 transition-colors duration-300">
								Buka Sesi Kasir
								<svg class="w-4 h-4 group-hover:translate-x-1.5 transition-transform duration-300 ease-out" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
									<path stroke-linecap="round" stroke-linejoin="round" d="M14 5l7 7m0 0l-7 7m7-7H3" />
								</svg>
							</div>
						</div>
					</div>

					<div @click="showPricingModal = true" class="bg-white rounded-[32px] p-6 border-2 border-dashed border-slate-200 hover:border-indigo-400 hover:bg-indigo-50/10 hover:-translate-y-1.5 transition-all duration-500 cubic-bezier-card cursor-pointer flex flex-col items-center justify-center text-center min-h-[210px] group shadow-3xs">
						<div class="w-12 h-12 rounded-full bg-slate-50 border border-slate-200/60 group-hover:bg-indigo-100 group-hover:scale-105 transition-all duration-300 flex items-center justify-center mb-3">
							<svg class="w-5 h-5 text-slate-400 group-hover:text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
								<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
							</svg>
						</div>
						<h3 class="text-sm font-black text-slate-600 group-hover:text-indigo-800 uppercase tracking-wider mb-0.5 transition-colors">Buka Cabang Baru</h3>
						<p class="text-[11px] font-bold text-slate-400 px-4 leading-relaxed">Daftarkan lokasi outlet atau cabang toko baru Anda ke dalam sistem.</p>
					</div>
				</div>
			</div>

			<div v-else class="w-full max-w-xl text-center py-12 px-6 bg-white rounded-[40px] border border-slate-100 shadow-xl">
				<div class="w-20 h-20 bg-indigo-50 border border-indigo-100 text-indigo-600 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-3xs">
					<svg class="w-10 h-10 animate-bounce" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
					</svg>
				</div>
				<h2 class="text-2xl font-black text-slate-800 uppercase tracking-tight mb-2">Belum Ada Toko Terdaftar</h2>
				<p class="text-slate-400 font-bold text-xs md:text-sm leading-relaxed mb-8 px-4">Akun Anda telah aktif secara global. Silakan pilih jenis bidang usaha untuk mengaktifkan sistem kasir pada toko pertama Anda.</p>
				<button @click="showPricingModal = true" class="w-full py-4 bg-indigo-600 hover:bg-slate-900 text-white rounded-2xl font-black text-xs uppercase tracking-widest shadow-xl transition-all duration-300 transform active:scale-95">Buat Toko Pertama Anda</button>
			</div>
		</div>

		<!-- REVISI TRANSPARAN MODAL: Menggunakan TransitionGroup dan Unique Key anti-kaku  -->
		<Transition name="modal-fade">
			<div v-if="showPricingModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-6 overflow-y-auto">
				<div @click="stores.length > 0 ? (showPricingModal = false) : null" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm transition-opacity duration-300"></div>

				<div class="bg-white w-full max-w-6xl rounded-[40px] shadow-2xl p-6 md:p-10 relative border border-slate-100 my-auto max-h-[92vh] overflow-y-auto custom-scrollbar flex flex-col items-center z-10 transform transition-all duration-500 cubic-bezier-modal">
					<button v-if="stores.length > 0" @click="showPricingModal = false" class="absolute top-6 right-6 w-10 h-10 bg-slate-50 hover:bg-rose-50 text-slate-400 hover:text-rose-500 rounded-full flex items-center justify-center transition-all duration-300 hover:rotate-90 focus:outline-none">
						<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>

					<div class="text-center mb-10">
						<h2 class="text-3xl md:text-4xl font-black text-slate-800 tracking-tight">
							INVESTASI
							<span class="text-indigo-600">TRANSPARAN</span>
						</h2>
						<p class="text-slate-400 font-bold text-xs md:text-sm mt-2 max-w-2xl uppercase tracking-wider">PILIH MODUL INDUSTRI ANDA, DAN TEMUKAN SKALABILITAS YANG DIRANCANG KHUSUS UNTUK BISNIS ANDA.</p>
					</div>

					<div class="bg-white/80 backdrop-blur-md p-1.5 rounded-[24px] border border-slate-200/80 flex w-full max-w-xl shadow-xl shadow-slate-100/50 mb-12 shrink-0">
						<button v-for="ind in industries" :key="ind.id" @click="activePricingTab = ind.id" :class="activePricingTab === ind.id ? 'bg-slate-900 text-white shadow-md' : 'text-slate-500 hover:text-slate-800 hover:bg-slate-50/60'" class="flex-1 px-2 py-3.5 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all text-center duration-500 focus:outline-none">
							{{ ind.title }}
						</button>
					</div>

					<!-- 🚀 KUNCI SAKTI: Mengubah pembungkus menjadi TransitionGroup untuk handle animasi list  -->
					<Transition name="tab-swap" mode="out-in">
						<div :key="activePricingTab" class="flex flex-wrap justify-center gap-6 w-full items-stretch min-h-[400px]">
							<div
								v-for="plan in filteredPlans"
								:key="plan.id"
								class="p-8 rounded-[32px] border flex flex-col relative group w-full sm:w-[calc(50%-12px)] lg:w-[calc(25%-18px)] min-w-[245px] transition-all duration-500"
								:class="[plan.id === 'trial' ? 'bg-white border-slate-200 shadow-sm border-2 border-dashed' : '', plan.id === 'basic' ? 'bg-white border-slate-200 shadow-sm hover:border-sky-300 hover:shadow-xl' : '', plan.id === 'pro' ? 'bg-white border-2 border-indigo-600 shadow-2xl lg:scale-105 z-10' : '', plan.id === 'premium' ? 'bg-slate-900 border-slate-800 shadow-2xl z-0' : '']">
								<div v-if="plan.id === 'pro'" class="absolute -top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-indigo-600 to-blue-600 text-white px-5 py-2 rounded-full text-[9px] font-black uppercase tracking-widest shadow-xl whitespace-nowrap animate-pulse">REKOMENDASI UMKM</div>
								<div v-if="plan.id === 'premium'" class="absolute -top-4 left-1/2 -translate-x-1/2 bg-gradient-to-r from-amber-500 to-amber-600 text-slate-950 px-5 py-2 rounded-full text-[9px] font-black uppercase tracking-widest shadow-xl whitespace-nowrap">FITUR ENTERPRISE</div>

								<div class="mb-8 mt-2 flex-1">
									<h3 class="font-black text-[11px] uppercase tracking-[0.2em] mb-4" :class="[plan.id === 'trial' ? 'text-indigo-500' : '', plan.id === 'basic' ? 'text-sky-500' : '', plan.id === 'pro' ? 'text-indigo-600' : '', plan.id === 'premium' ? 'text-amber-400' : '']">
										{{ plan.name }}
									</h3>
									<div class="flex items-baseline gap-1">
										<span class="text-4xl lg:text-5xl font-black tracking-tighter" :class="plan.id === 'premium' ? 'text-white' : 'text-slate-900'">{{ plan.price }}</span>
										<span class="font-bold text-[10px] uppercase tracking-widest mb-1 ml-1" :class="plan.id === 'premium' ? 'text-slate-400' : 'text-slate-400'">{{ plan.duration }}</span>
									</div>
									<p class="text-xs font-bold mt-5 h-12 leading-relaxed" :class="plan.id === 'premium' ? 'text-slate-400' : 'text-slate-500'">
										{{ plan.desc }}
									</p>

									<ul class="space-y-4 mb-2 mt-6 border-t pt-6" :class="plan.id === 'premium' ? 'border-slate-800' : 'border-slate-100'">
										<li v-for="feat in plan.features" :key="feat" class="flex items-start gap-3 text-xs font-bold leading-tight" :class="plan.id === 'premium' ? 'text-slate-300' : 'text-slate-700'">
											<svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3" :class="[plan.id === 'premium' ? 'text-amber-400' : '', plan.id === 'pro' ? 'text-indigo-500' : '', plan.id === 'basic' ? 'text-sky-500' : '', plan.id === 'trial' ? 'text-indigo-500' : '']">
												<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
											</svg>
											{{ feat }}
										</li>
									</ul>
								</div>

								<button
									@click="handlePilihPaketEkspansi(activePricingTab, plan.id)"
									:class="[
										plan.id === 'trial' ? 'bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white border border-indigo-100' : '',
										plan.id === 'basic' ? 'bg-slate-900 text-white hover:bg-sky-600 shadow-lg' : '',
										plan.id === 'pro' ? 'bg-gradient-to-r from-indigo-600 to-blue-600 text-white hover:from-slate-900 hover:to-slate-900 shadow-xl' : '',
										plan.id === 'premium' ? 'bg-gradient-to-r from-amber-400 to-amber-500 text-slate-950 font-extrabold hover:from-white hover:to-white shadow-xl' : '',
									]"
									class="block w-full text-center py-4 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all duration-300 transform active:scale-95 group-hover:scale-[1.01]">
									PILIH {{ plan.name }}
								</button>
							</div>
						</div>
					</Transition>
				</div>
			</div>
		</Transition>

		<div v-if="isLoading" class="fixed inset-0 z-[150] bg-slate-900/50 backdrop-blur-xs flex items-center justify-center">
			<div class="bg-white p-7 rounded-[24px] shadow-2xl border border-slate-100 flex flex-col items-center">
				<div class="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
				<div class="text-xs font-black text-slate-600 uppercase tracking-widest animate-pulse">Memuat Konfigurasi Toko...</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
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

.animate-fade-in-cards {
	animation: fadeInCards 0.35s ease-out both;
}
@keyframes fadeInCards {
	from {
		opacity: 0;
		transform: translateY(6px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.modal-fade-enter-from .bg-slate-900\/60 {
	opacity: 0;
}
.modal-fade-enter-from .cubic-bezier-modal {
	opacity: 0;
	transform: scale(0.95) translateY(10px);
}

.modal-fade-enter-active,
.modal-fade-leave-active {
	transition: all 0.4s ease-out;
}

.modal-fade-leave-to .bg-slate-900\/60 {
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
	background: #cbd5e1;
	border-radius: 10px;
}

/* 🚀 CSS PRESETS: Efek cross-fade premium yang sangat bersahabat dengan RAM komputer kantor  */
.tab-swap-enter-from {
	opacity: 0;
	transform: translateY(4px);
}
.tab-swap-leave-to {
	opacity: 0;
	transform: translateY(-4px);
}
.tab-swap-enter-active,
.tab-swap-leave-active {
	transition:
		opacity 0.2s ease,
		transform 0.2s ease;
	will-change: transform, opacity;
}
</style>
