<script setup>
import { ref } from 'vue';

const props = defineProps({
	// Menyediakan default data agar jika dipanggil tanpa props di manapun, langsung muncul otomatis
	plansData: {
		type: Object,
		default: () => ({
			retail: [
				{ id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Uji coba gratis untuk melihat kesesuaian sistem dengan toko Anda.', features: ['Semua Fitur Unlocked'] },
				{ id: 'basic', name: 'Retail Basic', price: 'Rp 49k', duration: '/Bulan', desc: 'Solusi tepat untuk pengelolaan satu toko skala kecil.', features: ['POS Kasir Penjualan', 'Manajemen Stok & LPB', 'Struk Thermal Bluetooth', 'Riwayat Log Struk'] },
				{ id: 'pro', name: 'Retail Pro', price: 'Rp 149k', duration: '/Bulan', desc: 'Cocok untuk toko ritel yang mulai mengelola staf operasional.', features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Pengaturan Shift', 'Laporan Ekspor (CSV/Excel)'] },
				{ id: 'premium', name: 'Retail Premium', price: 'Rp 299k', duration: '/Bulan', desc: 'Kendali penuh untuk bisnis skala besar & audit inventaris.', features: ['Semua Fitur Pro', 'Audit Stock Opname (SO)', 'Dashboard Analitik Laba', 'Notifikasi WhatsApp System'] },
			],
			fnb: [
				{ id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Uji coba gratis modul resto untuk kelancaran pesanan dapur.', features: ['Semua Fitur Unlocked'] },
				{ id: 'basic', name: 'F&B Basic', price: 'Rp 59k', duration: '/Bulan', desc: 'Sistem operasional efisien untuk kedai atau coffee shop.', features: ['Manajemen Layout Meja', 'Cetak Tiket Dapur (Kitchen)', 'Pajak & Service Charge', 'Struk Thermal Bluetooth'] },
				{ id: 'pro', name: 'F&B Pro', price: 'Rp 169k', duration: '/Bulan', desc: 'Ideal untuk restoran yang butuh kontrol manajemen staf teratur.', features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Shift Kerja', 'Split Bill & Gabung Meja'] },
				{ id: 'premium', name: 'F&B Premium', price: 'Rp 349k', duration: '/Bulan', desc: 'Skalabilitas bisnis franchise dengan laporan analitik terpusat.', features: ['Semua Fitur Pro', 'Resep Bahan Baku (BOM)', 'Self-Order QR Menu', 'Notifikasi WhatsApp System'] },
			],
			jasa: [
				{ id: 'trial', name: 'Starter Trial', price: 'Rp 0', duration: '14 Hari', desc: 'Coba modul layanan untuk bengkel, salon, atau laundry.', features: ['Semua Fitur Unlocked'] },
				{ id: 'basic', name: 'Service Basic', price: 'Rp 49k', duration: '/Bulan', desc: 'Sistem tracking pesanan yang rapi untuk bisnis jasa kecil.', features: ['Tracking Status Pesanan', 'Cetak Nota / Resi Barcode', 'Manajemen Layanan & Tarif', 'Laporan Pendapatan'] },
				{ id: 'pro', name: 'Service Pro', price: 'Rp 159k', duration: '/Bulan', desc: 'Sistem otomatisasi performa operasional tim staf.', features: ['Semua Fitur Basic', 'Manajemen Staf & HR', 'Absensi & Shift Kerja', 'Laporan Kinerja Bulanan'] },
				{ id: 'premium', name: 'Service Premium', price: 'Rp 329k', duration: '/Bulan', desc: 'Manajemen booking tingkat lanjut dengan pengingat otomatis.', features: ['Semua Fitur Pro', 'Dashboard Analitik', 'Sistem Booking Reservasi', 'Notifikasi WhatsApp System'] },
			],
		}),
	},
	// Mengetahui context apakah sedang dalam mode ekspansi atau registrasi biasa
	isExpansion: {
		type: Boolean,
		default: false,
	},
	// Mengontrol visibilitas tombol silang close
	showClose: {
		type: Boolean,
		default: true,
	},
});

const emit = defineEmits(['select-plan', 'close']);

// State Internal Tab
const activePricingTab = ref('retail');

const industriesList = [
	{ id: 'retail', title: 'Retail', isReady: true },
	{ id: 'fnb', title: 'F&B', isReady: false },
	{ id: 'jasa', title: 'Jasa', isReady: false },
];

const selectIndustryTab = (ind) => {
	if (!ind.isReady) return;
	activePricingTab.value = ind.id;
};

// Pengaturan warna teks label judul yang konsisten & kontras tinggi
const getPlanStyle = (planId) => {
	const p = planId ? planId.toLowerCase() : 'basic';
	if (p === 'premium') return 'text-amber-400 font-extrabold tracking-[0.25em]';
	if (p === 'pro') return 'text-indigo-600';
	if (p === 'basic') return 'text-blue-500';
	return 'text-slate-400';
};

const handleSelect = (planId) => {
	emit('select-plan', { industry: activePricingTab.value, plan: planId });
};
</script>

<template>
	<div class="w-full relative px-1 sm:px-4">
		<button v-if="props.showClose" @click="emit('close')" type="button" class="absolute -top-4 right-0 sm:right-4 w-10 h-10 bg-white hover:bg-rose-500 text-slate-500 hover:text-white rounded-full flex items-center justify-center transition-all duration-300 hover:rotate-90 focus:outline-none focus:ring-2 focus:ring-rose-500/20 shadow-sm z-50 border border-slate-200" title="Tutup">
			<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
				<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
			</svg>
		</button>

		<div class="text-center mb-12">
			<h2 class="text-3xl md:text-5xl font-black text-slate-900 tracking-tighter mb-4 uppercase">
				Investasi
				<span class="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600">Transparan</span>
			</h2>
			<p class="text-slate-400 font-bold text-xs sm:text-sm uppercase tracking-[0.2em] max-w-xl mx-auto leading-relaxed">Pilih modul industri Anda, dan temukan skalabilitas murni penunjang bisnis.</p>
		</div>

		<div class="w-full flex justify-center mb-16">
			<div class="bg-slate-100 p-1.5 rounded-2xl border border-slate-200/60 shadow-inner flex w-full max-w-md">
				<button v-for="ind in industriesList" :key="ind.id" @click="selectIndustryTab(ind)" :class="[activePricingTab === ind.id ? 'bg-slate-900 text-white shadow-lg shadow-slate-900/20' : 'text-slate-500 hover:text-slate-800', !ind.isReady ? 'opacity-40 cursor-not-allowed' : 'hover:bg-slate-600']" class="flex-1 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all text-center duration-300 focus:outline-none">
					{{ ind.title }}
				</button>
			</div>
		</div>

		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 sm:gap-8 items-stretch select-none">
			<div
				v-for="plan in props.plansData[activePricingTab]"
				:key="plan.id"
				class="p-6 sm:p-8 rounded-[36px] border transition-all duration-500 cubic-bezier-card flex flex-col relative group bg-white"
				:class="[
					plan.id === 'trial' ? 'border-slate-200/60 shadow-sm hover:shadow-xl hover:border-slate-300 hover:-translate-y-1' : '',
					plan.id === 'basic' ? 'border-slate-200 shadow-sm hover:border-blue-400 hover:shadow-xl hover:shadow-blue-500/5 hover:-translate-y-1' : '',
					plan.id === 'pro' ? 'border-2 border-indigo-600 shadow-xl shadow-indigo-100/70 lg:scale-[1.03] hover:scale-[1.05] lg:hover:scale-[1.06] z-10 hover:shadow-[0_30px_60px_-15px_rgba(79,70,229,0.25)]' : '',
					plan.id === 'premium' ? '!bg-slate-950 border-slate-900 shadow-xl text-white hover:shadow-[0_30px_60px_-15px_rgba(0,0,0,0.4)] hover:-translate-y-1' : '',
				]">
				<div v-if="plan.id === 'pro'" class="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest shadow-md whitespace-nowrap">REKOMENDASI UMKM</div>
				<div v-if="plan.id === 'premium'" class="absolute -top-3.5 left-1/2 -translate-x-1/2 bg-gradient-to-r from-amber-500 to-amber-600 text-slate-950 px-4 py-1.5 rounded-full text-[9px] font-black uppercase tracking-widest shadow-md whitespace-nowrap">ENTERPRISE CORE</div>

				<div class="mb-6 mt-1">
					<h3 class="font-black text-[11px] uppercase tracking-[0.2em] mb-3 group-hover:scale-105 origin-left transition-transform duration-300" :class="getPlanStyle(plan.id)">
						{{ plan.name }}
					</h3>
					<div class="flex items-baseline gap-0.5">
						<span class="text-3xl sm:text-4xl font-black tracking-tighter transition-colors duration-300" :class="[plan.id === 'premium' ? 'text-white' : 'text-slate-900', plan.id === 'pro' ? 'text-indigo-600' : '', plan.id === 'basic' ? 'group-hover:text-blue-600' : '']">{{ plan.price }}</span>
						<span class="font-bold text-[9px] uppercase tracking-widest ml-1 text-slate-400">{{ plan.duration }}</span>
					</div>
					<p class="text-xs font-medium mt-4 h-12 leading-relaxed transition-colors duration-300" :class="plan.id === 'premium' ? 'text-slate-400 group-hover:text-slate-300' : 'text-slate-500 group-hover:text-slate-600'">{{ plan.desc }}</p>
				</div>

				<ul class="space-y-3.5 mb-8 flex-1 border-t pt-5 transition-colors duration-300" :class="plan.id === 'premium' ? 'border-slate-800/80' : 'border-slate-100'">
					<li v-for="feat in plan.features" :key="feat" class="flex items-start gap-2.5 text-xs font-bold leading-tight group-hover:translate-x-0.5 transition-transform duration-300" :class="plan.id === 'premium' ? 'text-slate-300 group-hover:text-white' : 'text-slate-700 group-hover:text-slate-900'">
						<svg class="w-4 h-4 shrink-0 mt-0.5 transition-transform duration-300 group-hover:scale-110" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3.5" :class="[plan.id === 'premium' ? 'text-amber-400' : '', plan.id === 'pro' ? 'text-indigo-500' : '', plan.id === 'basic' ? 'text-blue-500' : '', plan.id === 'trial' ? 'text-slate-400' : '']">
							<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
						</svg>
						{{ feat }}
					</li>
				</ul>

				<button
					@click="handleSelect(plan.id)"
					:disabled="props.isExpansion && plan.id === 'trial'"
					:class="[
						plan.id === 'trial' ? 'bg-slate-50 hover:bg-slate-200 text-slate-700 border border-slate-200 focus:ring-slate-200' : '',
						plan.id === 'basic' ? 'bg-slate-900 text-white hover:bg-blue-600 shadow-md hover:shadow-blue-600/20 focus:ring-blue-500' : '',
						plan.id === 'pro' ? 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:from-slate-900 hover:to-slate-900 shadow-lg shadow-indigo-600/10 hover:shadow-none focus:ring-indigo-500' : '',
						plan.id === 'premium' ? 'bg-gradient-to-r from-amber-400 to-amber-500 text-slate-950 hover:from-white hover:to-white shadow-lg shadow-amber-500/5 hover:shadow-none focus:ring-amber-400' : '',
						props.isExpansion && plan.id === 'trial' ? '!opacity-30 !cursor-not-allowed hover:!bg-slate-50 text-slate-400' : '',
					]"
					class="block w-full text-center py-3.5 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all duration-300 transform active:scale-95 focus:outline-none focus:ring-2 focus:ring-offset-2">
					{{ props.isExpansion ? (plan.id === 'trial' ? 'TIDAK BERLAKU' : 'Pilih ' + plan.name) : plan.id === 'trial' ? 'Mulai Eksplorasi' : 'Luncurkan ' + plan.name }}
				</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
/* Kurva elastis kustom untuk animasi elevation card */
.cubic-bezier-card {
	transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}
</style>
