<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';

const props = defineProps({
	plansData: {
		type: Object,
		default: () => ({
			retail: [
				{ id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Uji coba gratis untuk melihat kesesuaian sistem dengan toko Anda.', features: ['Semua Fitur Unlocked', 'Batas 1 Terminal Perangkat', 'Masa Trial 14 Hari'] },
				{ id: 'basic', name: 'Retail Basic', price: 'Rp 49k', duration: '/Bulan', desc: 'Solusi tepat untuk pengelolaan satu toko skala kecil.', features: ['POS Kasir Penjualan', 'Manajemen Stok & LPB', 'Struk Thermal Bluetooth', 'Riwayat Log Struk'] },
				{ id: 'pro', name: 'Retail Pro', price: 'Rp 149k', duration: '/Bulan', desc: 'Cocok untuk toko ritel yang mulai mengelola staf operasional.', features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Pengaturan Shift', 'Laporan Ekspor (CSV/Excel)'] },
				{ id: 'premium', name: 'Retail Premium', price: 'Rp 299k', duration: '/Bulan', desc: 'Kendali penuh untuk bisnis skala besar & audit inventaris.', features: ['Semua Fitur Pro', 'Audit Stock Opname (SO)', 'Dashboard Analitik Laba', 'Notifikasi WhatsApp System'] },
			],
			fnb: [
				{ id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Uji coba gratis modul resto untuk kelancaran pesanan dapur.', features: ['Semua Fitur Unlocked', 'Masa Trial 14 Hari'] },
				{ id: 'basic', name: 'F&B Basic', price: 'Rp 59k', duration: '/Bulan', desc: 'Sistem operasional efisien untuk kedai atau coffee shop.', features: ['Manajemen Layout Meja', 'Cetak Tiket Dapur (Kitchen)', 'Pajak & Service Charge', 'Struk Thermal Bluetooth'] },
				{ id: 'pro', name: 'F&B Pro', price: 'Rp 169k', duration: '/Bulan', desc: 'Ideal untuk restoran yang butuh kontrol manajemen staf teratur.', features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Shift Kerja', 'Split Bill & Gabung Meja'] },
				{ id: 'premium', name: 'F&B Premium', price: 'Rp 349k', duration: '/Bulan', desc: 'Skalabilitas bisnis franchise dengan laporan analitik terpusat.', features: ['Semua Fitur Pro', 'Resep Bahan Baku (BOM)', 'Self-Order QR Menu', 'Notifikasi WhatsApp System'] },
			],
			jasa: [
				{ id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Uji coba penuh modul manajemen bisnis jasa terpadu untuk ruko Anda.', features: ['Semua Fitur Unlocked', 'Batas 1 Terminal Perangkat', 'Masa Trial 14 Hari'] },
				{ id: 'basic', name: 'Jasa Basic', price: 'Rp 49k', duration: '/Bulan', desc: 'Sistem pencatatan transaksi dan kasir digital instan ruko jasa rintisan.', features: ['POS Kasir Jasa', 'Katalog Tarif Kustom', 'Cetak Nota & Struk Thermal', 'Riwayat Transaksi Masuk'] },
				{ id: 'pro', name: 'Jasa Pro', price: 'Rp 149k', duration: '/Bulan', desc: 'Solusi optimal bagi outlet jasa yang ingin mengontrol performa tim staf.', features: ['Semua Fitur Jasa Basic', 'Sistem Tracker Status Kerja', 'Manajemen Akun Tim Staf', 'Absensi & Jadwal Shift Kerja'] },
				{ id: 'premium', name: 'Jasa Premium', price: 'Rp 299k', duration: '/Bulan', desc: 'Otomatisasi bisnis jasa tingkat lanjut dengan kecerdasan analitik laba rugi.', features: ['Semua Fitur Jasa Pro', 'Alokasi Nomor Rak / Slot', 'Integrasi WhatsApp Gateway', 'Dashboard Analitik Omset', 'Layanan Support 24/7'] },
			],
		}),
	},
	isExpansion: { type: Boolean, default: false },
	showClose: { type: Boolean, default: true },
	lockIndustry: { type: String, default: '' },
});

const emit = defineEmits(['select-plan', 'close']);

const activePricingTab = ref('retail');

const industriesList = [
	{ id: 'retail', title: 'Retail', isReady: true },
	{ id: 'fnb', title: 'F&B (soon)', isReady: false },
	{ id: 'jasa', title: 'Jasa', isReady: true },
];

const normalizeIndustryId = (id) => (id === 'laundry' || id === 'jasa' ? 'jasa' : id);

const filteredIndustriesList = computed(() => {
	if (!props.lockIndustry) return industriesList;
	const targetId = normalizeIndustryId(props.lockIndustry);
	return industriesList.filter((ind) => ind.id === targetId);
});

onMounted(() => {
	if (props.lockIndustry) {
		activePricingTab.value = normalizeIndustryId(props.lockIndustry);
	} else {
		const savedTab = localStorage.getItem('pendingIndustry');
		if (savedTab) activePricingTab.value = normalizeIndustryId(savedTab);
	}
	window.addEventListener('switch-pricing-tab', handleTabSwitchEvent);
});

onUnmounted(() => {
	window.removeEventListener('switch-pricing-tab', handleTabSwitchEvent);
});

const handleTabSwitchEvent = (e) => {
	activePricingTab.value = normalizeIndustryId(e.detail);
};

const selectIndustryTab = (ind) => {
	if (!ind.isReady) return;
	activePricingTab.value = ind.id;
};

const getPlanStyle = (planId) => {
	const p = planId ? planId.toLowerCase() : 'basic';
	if (p === 'premium') return 'text-amber-400 font-extrabold tracking-[0.2em]';
	if (p === 'pro') return 'text-indigo-500';
	if (p === 'basic') return 'text-blue-500';
	return 'text-slate-400';
};

const handleSelect = (planId) => {
	emit('select-plan', { industry: activePricingTab.value, plan: planId });
};
</script>

<template>
	<div class="w-full relative px-4 sm:px-6 py-6 md:py-10 bg-slate-50/40 rounded-3xl">
		<button v-if="props.showClose" @click="emit('close')" type="button" class="absolute top-4 right-4 w-10 h-10 bg-white hover:bg-rose-500 text-slate-400 hover:text-white rounded-xl flex items-center justify-center transition-all duration-300 hover:rotate-90 focus:outline-none border border-slate-200/80 shadow-sm z-50" title="Tutup">
			<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
				<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
			</svg>
		</button>

		<div class="text-center mb-10 md:mb-14">
			<h2 class="text-3xl sm:text-4xl md:text-5xl font-black text-slate-900 tracking-tight uppercase leading-none">
				Investasi
				<span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-blue-600 to-purple-600">Transparan</span>
			</h2>
			<p class="text-slate-400 font-bold text-[11px] sm:text-xs uppercase tracking-[0.25em] max-w-xl mx-auto leading-relaxed mt-4 px-2">Pilih modul industri Anda, dan temukan skalabilitas murni penunjang bisnis.</p>
		</div>

		<div class="w-full flex justify-center mb-12 px-4">
			<div class="bg-slate-200/60 p-1 rounded-2xl border border-slate-300/40 flex w-full max-w-sm">
				<button v-for="ind in filteredIndustriesList" :key="ind.id" @click="selectIndustryTab(ind)" :disabled="!ind.isReady" :class="[activePricingTab === ind.id ? 'bg-white text-slate-900 shadow-sm font-black' : 'text-slate-400 hover:text-slate-600 font-bold', !ind.isReady ? 'opacity-40 cursor-not-allowed' : '']" class="flex-1 py-2.5 rounded-xl text-[10px] sm:text-xs uppercase tracking-widest transition-all duration-300 focus:outline-none">
					{{ ind.title }}
				</button>
			</div>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6 lg:gap-8 items-stretch select-none">
			<div
				v-for="plan in props.plansData[activePricingTab]"
				:key="plan.id"
				class="group p-6 sm:p-8 rounded-[32px] border flex flex-col relative bg-white transition-all duration-500 cubic-bezier-card justify-between"
				:class="[
					plan.id === 'trial' ? 'border-slate-200 shadow-sm hover:shadow-xl hover:border-slate-300 hover:-translate-y-2' : '',
					plan.id === 'basic' ? 'border-slate-200 shadow-sm hover:border-blue-400 hover:shadow-xl hover:shadow-blue-500/5 hover:-translate-y-2' : '',
					plan.id === 'pro' ? 'border-2 border-indigo-600 shadow-lg shadow-indigo-100 z-10 hover:shadow-[0_25px_50px_-12px_rgba(79,70,229,0.22)] hover:-translate-y-2.5 scale-100 md:scale-[1.01]' : '',
					plan.id === 'premium' ? '!bg-slate-950 border-slate-900 shadow-xl text-white hover:shadow-[0_30px_60px_-12px_rgba(0,0,0,0.55)] hover:-translate-y-2 hover:border-amber-500/40' : '',
				]">
				<div v-if="plan.id === 'pro'" class="absolute -top-3 left-1/2 -translate-x-1/2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-4 py-1 rounded-full text-[9px] font-black uppercase tracking-widest shadow-md shadow-indigo-500/20 flex items-center gap-1.5 whitespace-nowrap z-20">
					<svg class="w-3 h-3 text-amber-300 fill-current" viewBox="0 0 24 24">
						<path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z" />
					</svg>
					REKOMENDASI UMKM
				</div>

				<div v-if="plan.id === 'premium'" class="absolute -top-3 left-1/2 -translate-x-1/2 bg-gradient-to-r from-amber-400 to-amber-500 text-slate-950 px-4 py-1 rounded-full text-[9px] font-black uppercase tracking-widest shadow-md shadow-amber-500/20 flex items-center gap-1.5 whitespace-nowrap z-20">
					<svg class="w-3 h-3 text-slate-950 fill-current" viewBox="0 0 24 24">
						<path d="M2 22h20V2zM12 2l4 8h6l-5 4 2 8-7-5-7 5 2-8-5-4h6z" />
					</svg>
					ENTERPRISE CORE
				</div>

				<div class="mb-6 mt-2">
					<h3 class="font-black text-[11px] sm:text-xs uppercase tracking-[0.25em] mb-4 transition-transform duration-300 group-hover:translate-x-1" :class="getPlanStyle(plan.id)">
						{{ plan.name }}
					</h3>
					<div class="flex items-baseline gap-1">
						<span class="text-4xl sm:text-5xl font-black tracking-tight transition-colors duration-300" :class="[plan.id === 'premium' ? 'text-white' : 'text-slate-900', plan.id === 'pro' ? 'text-indigo-600' : '']">
							{{ plan.price }}
						</span>
						<span class="font-bold text-[10px] uppercase tracking-widest ml-1 text-slate-400">
							{{ plan.duration }}
						</span>
					</div>
					<p class="text-xs font-medium mt-4 min-h-[48px] leading-relaxed transition-colors duration-300" :class="plan.id === 'premium' ? 'text-slate-400' : 'text-slate-500'">
						{{ plan.desc }}
					</p>
				</div>

				<ul class="space-y-4 mb-8 flex-1 border-t pt-6 transition-colors duration-300" :class="plan.id === 'premium' ? 'border-slate-800/80' : 'border-slate-100'">
					<li v-for="feat in plan.features" :key="feat" class="flex items-start gap-3 text-xs font-medium leading-snug transition-all duration-300" :class="plan.id === 'premium' ? 'text-slate-300 group-hover:text-white' : 'text-slate-600 group-hover:text-slate-900'">
						<svg class="w-4 h-4 shrink-0 mt-0.5 transition-transform duration-300 group-hover:scale-125" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3" :class="[plan.id === 'premium' ? 'text-amber-400' : '', plan.id === 'pro' ? 'text-indigo-500' : '', plan.id === 'basic' ? 'text-blue-500' : '', plan.id === 'trial' ? 'text-slate-400' : '']">
							<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
						</svg>
						<span class="truncate-feature">{{ feat }}</span>
					</li>
				</ul>

				<button
					@click="handleSelect(plan.id)"
					:disabled="props.isExpansion && plan.id === 'trial'"
					:class="[
						plan.id === 'trial' ? 'bg-slate-50 hover:bg-slate-100 text-slate-700 border border-slate-200 hover:border-slate-300' : '',
						plan.id === 'basic' ? 'bg-slate-900 text-white hover:bg-blue-600 shadow-md hover:shadow-blue-600/20' : '',
						plan.id === 'pro' ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:from-indigo-700 hover:to-purple-700 shadow-lg shadow-indigo-600/10' : '',
						plan.id === 'premium' ? 'bg-gradient-to-r from-amber-400 to-amber-500 text-slate-950 hover:opacity-95 shadow-lg shadow-amber-500/10' : '',
						props.isExpansion && plan.id === 'trial' ? '!opacity-20 !cursor-not-allowed hover:!bg-slate-50 text-slate-400' : '',
					]"
					class="block w-full text-center py-3.5 rounded-2xl font-black text-xs uppercase tracking-widest transition-all duration-300 transform active:scale-98 focus:outline-none shadow-sm shrink-0">
					{{ props.isExpansion ? (plan.id === 'trial' ? 'TIDAK BERLAKU' : 'Pilih Paket') : plan.id === 'trial' ? 'Mulai Eksplorasi' : 'Luncurkan Paket' }}
				</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
.cubic-bezier-card {
	transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}

.truncate-feature {
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
	text-overflow: ellipsis;
}
</style>
