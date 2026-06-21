<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import PricingPlan from '../../components/PricingPlan.vue';

const router = useRouter();
const route = useRoute();
const isMenuOpen = ref(false);
const isScrolled = ref(false);

// Mengamati apakah user masuk lewat action ekspansi toko
const isExpansionMode = computed(() => route.query.action === 'expansion');

onMounted(() => {
	window.addEventListener('scroll', () => {
		isScrolled.value = window.scrollY > 20;
	});

	if (route.query.action === 'expansion') {
		setTimeout(() => {
			scrollToSection('pricing');
		}, 500);
	}
});

const features = [
	{
		title: 'Sistem Kasir Pintar',
		desc: 'Transaksi instan dengan dukungan barcode, manajemen cepat, dan pencetakan struk thermal otomatis.',
		icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>`,
	},
	{
		title: 'Kontrol Multi-Gudang',
		desc: 'Visibilitas inventaris real-time di seluruh cabang. Dilengkapi fitur Stock Opname (SO) dan manajemen retur barang.',
		icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg>`,
	},
	{
		title: 'Manajemen SDM & Shift',
		desc: 'Lacak kehadiran tim operasional toko, atur rotasi jadwal shift kerja, dan amankan logs aktivitas operator.',
		icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`,
	},
	{
		title: 'Ekosistem Terintegrasi',
		desc: 'Dari notifikasi WhatsApp operasional, integrasi pembayaran digital QRIS, hingga laporan analitik profit komprehensif.',
		icon: `<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg>`,
	},
];

const industries = [
	{
		id: 'retail',
		title: 'Retail & Distribusi',
		desc: 'Supermarket, Butik, Minimarket, Toko Kelontong, Elektronik',
		icon: 'M16 11V7a4 4 0 0 0-8 0v4M5 9h14l1 12H4L5 9z',
		isReady: true,
	},
	{
		id: 'fnb',
		title: 'Food & Beverage (Cooming Soon)',
		desc: 'Cafe, Restoran, Kedai Kopi, Waralaba / Franchise',
		icon: 'M18 8h1a4 4 0 0 1 0 8h-1M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z M6 1v3 M10 1v3 M14 1v3',
		isReady: false,
	},
	{
		id: 'jasa',
		title: 'Layanan Jasa (Cooming Soon)',
		desc: 'Laundry, Barbershop, Bengkel, Steam Kendaraan',
		icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05',
		isReady: false,
	},
];

// 🚀 UPGRADE: Fungsi scroll pintar yang merubah tab pricing otomatis!
const scrollToSection = (id, selectedIndustry = null) => {
	isMenuOpen.value = false;

	if (selectedIndustry) {
		localStorage.setItem('activeLandingTab', selectedIndustry);
		window.dispatchEvent(new CustomEvent('switch-pricing-tab', { detail: selectedIndustry }));
	}

	const el = document.getElementById(id);
	if (el) el.scrollIntoView({ behavior: 'smooth' });
};

// Menangani event emit tombol pilih paket dari dalam komponen PricingPlan
const handlePilihPaket = (payload) => {
	const { industry, plan } = payload;

	if (isExpansionMode.value && plan === 'trial') return;

	localStorage.setItem('pendingIndustry', industry);
	localStorage.setItem('pendingPlan', plan);

	if (isExpansionMode.value) {
		router.push('/setup-toko?is_expansion=true');
	} else {
		router.push('/register');
	}
};
</script>

<template>
	<div class="min-h-screen bg-[#F8FAFC] font-sans text-slate-900 selection:bg-indigo-600 selection:text-white overflow-x-hidden antialiased">
		<nav :class="['fixed top-0 w-full z-50 transition-all duration-300', isScrolled ? 'bg-white/90 backdrop-blur-xl border-b border-slate-200/60 py-3.5 shadow-lg shadow-slate-100/40' : 'bg-transparent py-5 sm:py-6']">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex items-center justify-between">
				<div class="flex items-center cursor-pointer select-none" @click="scrollToSection('hero')">
					<img src="/favicon.svg" alt="ARZURA POS LOGO" class="h-8 md:h-9 w-auto object-contain" />
					<span class="ml-2.5 font-black text-lg tracking-tighter text-slate-900 uppercase">
						ARZURA
						<span class="text-indigo-600 italic">POS</span>
					</span>
				</div>

				<div class="hidden md:flex items-center gap-8 lg:gap-10">
					<button @click="scrollToSection('industri')" class="text-[11px] font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Solusi</button>
					<button @click="scrollToSection('features')" class="text-[11px] font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Teknologi</button>
					<button @click="scrollToSection('pricing')" class="text-[11px] font-black uppercase tracking-widest text-slate-500 hover:text-indigo-600 transition-colors">Harga</button>

					<template v-if="!isExpansionMode">
						<router-link to="/login" class="text-[11px] font-black uppercase tracking-widest text-slate-800 hover:text-indigo-600 transition-colors ml-2 border-l pl-5 border-slate-200">LOGIN</router-link>
						<button @click="scrollToSection('pricing')" class="bg-slate-900 text-white text-[10px] font-black uppercase tracking-[0.2em] px-6 py-3 rounded-xl hover:bg-indigo-600 shadow-xl hover:shadow-indigo-200/50 transition-all active:scale-95 duration-300">Mulai Trial</button>
					</template>
					<template v-else>
						<router-link to="/select-store" class="text-[11px] font-black uppercase tracking-widest text-rose-600 hover:text-rose-700 transition-colors ml-2 border-l pl-5 border-slate-200">Batal Ekspansi</router-link>
					</template>
				</div>

				<button @click="isMenuOpen = !isMenuOpen" class="md:hidden p-2 text-slate-600 hover:text-indigo-600 transition-colors focus:outline-none rounded-xl active:bg-slate-100">
					<svg v-if="!isMenuOpen" xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" /></svg>
					<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
				</button>
			</div>

			<transition name="fade">
				<div v-show="isMenuOpen" class="md:hidden absolute top-full left-0 w-full bg-white border-b border-slate-200 shadow-2xl flex flex-col py-4 px-4 gap-1">
					<button @click="scrollToSection('industri')" class="text-left text-xs font-black uppercase tracking-widest text-slate-600 py-3 px-3 rounded-xl hover:bg-slate-50 transition-colors">Solusi Industri</button>
					<button @click="scrollToSection('features')" class="text-left text-xs font-black uppercase tracking-widest text-slate-600 py-3 px-3 rounded-xl hover:bg-slate-50 transition-colors">Teknologi Utama</button>
					<button @click="scrollToSection('pricing')" class="text-left text-xs font-black uppercase tracking-widest text-slate-600 py-3 px-3 rounded-xl hover:bg-slate-50 transition-colors">Daftar Harga</button>
					<router-link to="/login" class="text-left text-xs font-black uppercase tracking-widest text-slate-600 py-3 px-3 rounded-xl hover:bg-slate-50 transition-colors">Login</router-link>
				</div>
			</transition>
		</nav>

		<section id="hero" class="relative pt-32 sm:pt-40 md:pt-44 pb-20 md:pb-28 overflow-hidden">
			<div class="absolute -top-40 -left-40 w-[45rem] h-[45rem] bg-indigo-300/20 rounded-full blur-3xl opacity-70 pointer-events-none"></div>
			<div class="absolute top-20 -right-20 w-[35rem] h-[35rem] bg-violet-300/20 rounded-full blur-3xl opacity-70 pointer-events-none"></div>

			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10 flex flex-col lg:flex-row items-center gap-12 lg:gap-8">
				<div class="flex-1 text-center lg:text-left">
					<div class="inline-flex items-center gap-2 px-4 py-2 bg-white border border-slate-200/80 rounded-full mb-6 shadow-sm">
						<span class="w-2 h-2 rounded-full bg-gradient-to-r from-indigo-600 to-purple-500 animate-pulse"></span>
						<span class="text-[9px] font-black text-slate-500 uppercase tracking-[0.2em]">Infrastruktur SaaS POS Premium Cloud</span>
					</div>

					<h1 class="text-4xl sm:text-5xl md:text-6xl font-black text-slate-900 tracking-tighter mb-6 leading-[1.08] uppercase">
						Orkestrasi Bisnis
						<br class="hidden lg:block" />
						Terpadu Dalam
						<br />
						<span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-purple-600 to-blue-500 italic">Satu Ekosistem.</span>
					</h1>

					<p class="max-w-xl mx-auto lg:mx-0 text-slate-500 font-bold text-sm sm:text-base mb-10 leading-relaxed">Tinggalkan sistem konvensional terpisah. Kendalikan transaksi kasir toko, manajemen rantai pasok inventaris harian, hingga analitik keuntungan bersih riil secara terpusat dari satu dashboard.</p>

					<div class="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4">
						<button @click="scrollToSection('pricing')" class="w-full sm:w-auto bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-8 py-4 rounded-xl font-black text-xs uppercase tracking-widest shadow-xl shadow-indigo-200 hover:from-slate-900 hover:to-slate-900 transition-all duration-300 hover:-translate-y-0.5 active:scale-95 flex items-center justify-center gap-2">
							Lihat Paket Sistem
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14M12 5l7 7-7 7" /></svg>
						</button>
						<span class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-2 sm:mt-0">Proses Akses Setup &lt; 2 Menit</span>
					</div>
				</div>

				<div class="flex-1 w-full max-w-md lg:max-w-none relative mt-4 lg:mt-0">
					<svg class="absolute -top-8 -left-8 w-24 h-24 text-slate-200/60" fill="currentColor" viewBox="0 0 100 100">
						<pattern id="grid-dots" x="0" y="0" width="16" height="16" patternUnits="userSpaceOnUse"><circle cx="2" cy="2" r="1.5"></circle></pattern>
						<rect width="100" height="100" fill="url(#grid-dots)"></rect>
					</svg>

					<div class="relative bg-white p-6 sm:p-8 rounded-[36px] shadow-2xl shadow-indigo-950/5 border border-slate-100 z-10 transform hover:-translate-y-1.5 transition-transform duration-500">
						<div class="flex items-center justify-between mb-6 pb-5 border-b border-slate-100">
							<div class="flex items-center gap-3">
								<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white shadow-lg shadow-indigo-100">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" /></svg>
								</div>
								<div>
									<div class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-0.5">Live Operational Analytics</div>
									<div class="text-2xl font-black text-slate-800 tracking-tighter">Rp 24.580.000</div>
								</div>
							</div>
							<span class="bg-emerald-50 text-emerald-600 px-2.5 py-1 rounded-lg text-[9px] font-black uppercase tracking-widest flex items-center gap-0.5 border border-emerald-100">
								<svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18" /></svg>
								12.5%
							</span>
						</div>

						<div class="space-y-3.5">
							<div class="flex justify-between items-center text-[10px] font-black text-slate-400 uppercase tracking-wider mb-1">
								<span>Aktivitas Alur Data</span>
								<span>Status Gateway</span>
							</div>
							<div class="flex items-center justify-between p-3.5 bg-slate-50/80 rounded-xl border border-slate-100">
								<div class="flex items-center gap-2.5">
									<div class="w-2 h-2 rounded-full bg-emerald-500"></div>
									<span class="text-xs font-black text-slate-700">Audit Inventaris Multi-Gudang</span>
								</div>
								<span class="text-[9px] font-black text-slate-400 uppercase tracking-wider">Selesai</span>
							</div>
							<div class="flex items-center justify-between p-3.5 bg-slate-50/80 rounded-xl border border-slate-100">
								<div class="flex items-center gap-2.5">
									<div class="w-2 h-2 rounded-full bg-indigo-500 animate-pulse"></div>
									<span class="text-xs font-black text-slate-700">Backup Secure Cloud Storage</span>
								</div>
								<span class="text-[9px] font-black text-indigo-500 uppercase tracking-wider">Aktif</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</section>

		<section id="industri" class="py-24 bg-slate-900 text-white relative overflow-hidden">
			<div class="absolute -bottom-48 -left-48 w-80 h-80 bg-indigo-500/10 rounded-full blur-3xl pointer-events-none"></div>
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
				<div class="text-center mb-16">
					<h2 class="text-3xl md:text-5xl font-black tracking-tighter mb-4 uppercase">
						Arsitektur Modular
						<br />
						<span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-purple-400 to-blue-400">Sesuai Klaster Kebutuhan</span>
					</h2>
					<p class="text-slate-400 font-bold text-xs sm:text-sm uppercase tracking-[0.2em] max-w-xl mx-auto leading-relaxed">Sistem cerdas otomatis terkonfigurasi membelah fitur sesuai model operasional bisnis Anda.</p>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8">
					<div v-for="ind in industries" :key="ind.title" @click="ind.isReady ? scrollToSection('pricing', ind.id) : null" :class="[!ind.isReady ? 'opacity-40 cursor-not-allowed border-slate-800 bg-slate-950/20 select-none' : 'bg-slate-800/40 cursor-pointer border-slate-800/80 hover:border-indigo-500 hover:bg-slate-800/90']" class="p-6 sm:p-8 rounded-[28px] border transition-all duration-300 flex flex-col items-start gap-6 group shadow-xl relative">
						<div v-if="!ind.isReady" class="absolute top-4 right-4 bg-amber-500/20 border border-amber-500/30 text-amber-400 px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest">SOON / DEVELOPING</div>

						<div class="w-14 h-14 rounded-xl bg-slate-800 text-indigo-400 flex items-center justify-center border border-slate-700" :class="{ 'group-hover:bg-gradient-to-tr group-hover:from-indigo-600 group-hover:to-purple-500 group-hover:text-white group-hover:scale-105 group-hover:border-transparent transition-all duration-300': ind.isReady }">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path :d="ind.icon" /></svg>
						</div>
						<div>
							<h4 class="font-black text-lg uppercase tracking-wider text-white mb-2" :class="{ 'group-hover:text-indigo-300 transition-colors': ind.isReady }">{{ ind.title }}</h4>
							<p class="text-[11px] font-medium text-slate-400 leading-relaxed">{{ ind.desc }}</p>
						</div>
						<div class="mt-auto pt-5 w-full flex items-center justify-between text-indigo-400 text-[10px] font-black uppercase tracking-widest border-t border-slate-800/50" :class="{ 'group-hover:text-white transition-colors': ind.isReady }">
							<span>{{ ind.isReady ? 'Buka Paket Integrasi' : 'Klaster Dikunci' }}</span>
							<svg v-if="ind.isReady" class="w-4 h-4 transform group-hover:translate-x-1.5 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M17 8l4 4m0 0l-4 4m4-4H3" /></svg>
							<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-slate-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
								<rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
								<path d="M7 11V7a5 5 0 0 1 10 0v4" />
							</svg>
						</div>
					</div>
				</div>
			</div>
		</section>

		<section id="features" class="py-24 bg-white relative border-t border-slate-100">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="text-center mb-16">
					<h2 class="text-3xl md:text-5xl font-black text-slate-900 tracking-tighter mb-4 uppercase">
						Teknologi Inti
						<span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">Sistem POS</span>
					</h2>
					<p class="text-slate-400 font-bold text-xs sm:text-sm uppercase tracking-[0.2em]">Pilar utama penahan laju efisiensi manajemen operasional harian Anda.</p>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 sm:gap-8">
					<div v-for="feat in features" :key="feat.title" class="bg-[#F8FAFC]/80 p-6 sm:p-8 rounded-[28px] border border-slate-200/50 hover:shadow-xl hover:shadow-indigo-950/5 hover:border-indigo-100 transition-all duration-300 group hover:-translate-y-1 flex flex-col">
						<div class="w-12 h-12 bg-white border border-slate-200 text-indigo-600 rounded-xl flex items-center justify-center mb-5 group-hover:bg-gradient-to-tr group-hover:from-indigo-600 group-hover:to-purple-500 group-hover:text-white transition-all duration-300 shadow-sm" v-html="feat.icon"></div>
						<h3 class="text-base font-black text-slate-800 mb-2.5 tracking-tight uppercase">{{ feat.title }}</h3>
						<p class="text-slate-500 font-medium text-xs leading-relaxed">{{ feat.desc }}</p>
					</div>
				</div>
			</div>
		</section>

		<section id="pricing" class="py-24 bg-[#F8FAFC] border-t border-slate-200/50 relative overflow-hidden">
			<div class="absolute top-1/4 left-1/2 -translate-x-1/2 w-[55rem] h-[55rem] bg-indigo-50/40 rounded-full blur-3xl pointer-events-none"></div>

			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
				<PricingPlan :is-expansion="isExpansionMode" :show-close="false" @select-plan="handlePilihPaket" />
			</div>
		</section>

		<section class="py-20 bg-[#F8FAFC]">
			<div class="max-w-5xl mx-auto px-4 sm:px-6">
				<div class="bg-gradient-to-br from-indigo-950 via-slate-950 to-black rounded-[40px] p-10 md:p-16 text-center relative overflow-hidden shadow-2xl">
					<div class="absolute top-0 right-0 w-72 h-72 bg-indigo-500/10 rounded-full blur-3xl pointer-events-none"></div>
					<div class="absolute bottom-0 left-0 w-72 h-72 bg-purple-500/10 rounded-full blur-3xl pointer-events-none"></div>

					<div class="relative z-10">
						<h2 class="text-3xl md:text-5xl font-black text-white tracking-tighter mb-4 leading-tight uppercase">
							Siap Elevasi Bisnis
							<br />
							Hari Ini?
						</h2>

						<button @click="scrollToSection('pricing')" class="inline-flex items-center gap-2 bg-white text-slate-900 px-8 py-4 rounded-full font-black text-xs uppercase tracking-widest hover:bg-indigo-50 shadow-xl transition-all duration-300 active:scale-95 mt-4">
							Mulai Integrasi Sistem
							<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M5 12h14M12 5l7 7-7 7" /></svg>
						</button>

						<div class="mt-10 flex flex-row items-center justify-center gap-6 text-indigo-200/70 text-[9px] font-black uppercase tracking-widest">
							<span class="flex items-center gap-1">
								<svg class="w-3.5 h-3.5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
								Akses Instan
							</span>
							<span class="opacity-30">•</span>
							<span class="flex items-center gap-1">
								<svg class="w-3.5 h-3.5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
								Transparansi Finansial
							</span>
						</div>
					</div>
				</div>
			</div>
		</section>

		<footer class="bg-white border-t border-slate-200/80 pt-16 pb-12 font-sans">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="grid grid-cols-2 md:grid-cols-12 gap-8 pb-12 border-b border-slate-100 text-left">
					<div class="col-span-2 md:col-span-4 flex flex-col gap-4">
						<div class="flex items-center">
							<img src="/favicon.svg" alt="ARZURA POS LOGO" class="h-6 w-auto object-contain" />
							<div class="ml-2 font-black text-2xl tracking-tighter text-slate-900 uppercase">
								ARZURA
								<span class="text-indigo-600">POS</span>
							</div>
						</div>
						<p class="text-xs text-slate-500 font-medium max-w-sm leading-relaxed">Infrastruktur sistem POS & manajemen ERP modern multi-cabang terintegrasi arsitektur cloud terpusat. Solusi andalan UMKM mempercepat kemajuan bisnis operasional.</p>
					</div>

					<div class="col-span-1 md:col-span-2 flex flex-col gap-3">
						<h4 class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Ekosistem</h4>
						<ul class="flex flex-col gap-2.5 text-xs font-bold text-slate-600">
							<li><button @click="scrollToSection('features')" class="hover:text-indigo-600 transition-colors">Kasir Pintar POS</button></li>
							<li><button @click="scrollToSection('features')" class="hover:text-indigo-600 transition-colors">Multi-Gudang Cloud</button></li>
							<li><button @click="scrollToSection('features')" class="hover:text-indigo-600 transition-colors">Manajemen Shift HR</button></li>
							<li><button @click="scrollToSection('pricing')" class="hover:text-indigo-600 transition-colors">Integrasi WhatsApp</button></li>
						</ul>
					</div>

					<div class="col-span-1 md:col-span-2 flex flex-col gap-3">
						<h4 class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Klaster Solusi</h4>
						<ul class="flex flex-col gap-2.5 text-xs font-bold text-slate-600">
							<li><button @click="scrollToSection('pricing')" class="hover:text-indigo-600 transition-colors">Ritel & Minimarket</button></li>
							<li><button class="text-slate-400 cursor-not-allowed line-through" disabled>Layanan Jasa (Soon)</button></li>
							<li><button class="text-slate-400 cursor-not-allowed line-through" disabled>F&B Resto (Soon)</button></li>
							<li><button @click="scrollToSection('hero')" class="hover:text-indigo-600 transition-colors">Multi-Store Expansion</button></li>
						</ul>
					</div>

					<div class="col-span-2 md:col-span-4 flex flex-col gap-3">
						<h4 class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Keamanan Server</h4>
						<div class="flex flex-col gap-2 text-xs font-medium text-slate-500 leading-relaxed">
							<p class="font-bold text-slate-700">ARZURA SaaS Cloud Platform Architecture</p>
							<p>Seluruh transmisi data transaksi dienkripsi penuh menggunakan protocol SSL berlapis di atas arsitektur infrastruktur terdistribusi global.</p>
							<div class="flex items-center gap-3 mt-1 pt-3 border-t border-slate-100">
								<span class="bg-indigo-50 text-indigo-600 font-black px-2.5 py-1 rounded text-[9px] uppercase tracking-widest border border-indigo-100 animate-pulse">GATEWAY ENGINE: ONLINE</span>
							</div>
						</div>
					</div>
				</div>

				<div class="pt-8 flex flex-col md:flex-row items-center justify-between gap-4 text-center md:text-left">
					<div class="text-[10px] font-black text-slate-400 uppercase tracking-[1.5px] md:order-2 flex flex-wrap justify-center gap-x-6 gap-y-2">
						<a href="#" class="hover:text-indigo-600 transition-colors">Kebijakan Privasi</a>
						<a href="#" class="hover:text-indigo-600 transition-colors">Syarat Ketentuan</a>
						<a href="#" class="hover:text-indigo-600 transition-colors">Pusat Bantuan</a>
					</div>
					<div class="text-[10px] font-black text-slate-500 uppercase tracking-[1.5px]">
						Hak Cipta &copy; 2026 Developed by
						<span class="text-slate-800 font-black">Arif Juliawan</span>
						• All Rights Reserved.
					</div>
				</div>
			</div>
		</footer>
	</div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
	transition:
		opacity 0.2s ease-out,
		transform 0.2s ease-out;
}
.fade-enter-from,
.fade-leave-to {
	opacity: 0;
	transform: translateY(-8px);
}
</style>
