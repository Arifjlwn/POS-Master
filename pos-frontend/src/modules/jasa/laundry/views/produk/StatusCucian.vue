<script setup>
import SidebarLaundry from '../../components/SidebarLaundry.vue';
import { useKanbanLaundry } from '../../composables/useKanbanLaundry.js';

// 🚀 CLEAN ARCHITECTURE: Tarik semua otak sirkuit dari Composable!
const { riwayat, isLoading, searchQuery, orderAntri, orderProses, orderSelesai, formatDate, updateStatusKanban, prosesPengambilan } = useKanbanLaundry();
</script>

<template>
	<SidebarLaundry>
		<div class="flex-1 flex flex-col h-full bg-[#F8FAFC] overflow-hidden relative font-sans">
			<div class="p-5 md:p-6 shrink-0 bg-white border-b border-slate-200 flex flex-col sm:flex-row justify-between sm:items-center z-10 shadow-sm relative gap-4">
				<div class="flex items-center gap-4">
					<div class="w-12 h-12 bg-indigo-50 border border-indigo-100 rounded-2xl flex items-center justify-center shrink-0 shadow-inner">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2" />
						</svg>
					</div>
					<div>
						<h1 class="text-xl font-black tracking-tight uppercase text-slate-900 leading-tight">Papan Operasional</h1>
						<p class="text-[10px] font-bold text-slate-400 mt-1 uppercase tracking-widest flex items-center gap-2">
							Alur Sirkuit Cucian
							<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
						</p>
					</div>
				</div>

				<div class="relative w-full sm:w-64">
					<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
							<path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
						</svg>
					</div>
					<input v-model="searchQuery" type="text" placeholder="Cari Resi / Pelanggan..." class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 outline-none font-bold text-xs text-slate-800 transition-all shadow-inner placeholder:text-slate-400" />
				</div>
			</div>

			<div class="flex-1 overflow-y-auto custom-scrollbar p-5 md:p-6 pb-20">
				<div v-if="isLoading && riwayat.length === 0" class="flex flex-col items-center justify-center py-20 text-slate-400">
					<div class="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
					<p class="font-black text-[10px] uppercase tracking-widest animate-pulse">Memuat Topologi Data...</p>
				</div>

				<div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
					<div class="bg-white p-4 rounded-[28px] border border-slate-200 flex flex-col h-[70vh] shadow-sm">
						<div class="flex items-center justify-between mb-5 px-2">
							<h3 class="font-black text-[11px] uppercase tracking-widest text-slate-800 flex items-center gap-2">
								<div class="w-2.5 h-2.5 rounded-full bg-rose-400 ring-4 ring-rose-50"></div>
								Antrean Masuk
							</h3>
							<span class="bg-slate-100 text-slate-600 text-[10px] font-black px-2.5 py-1 rounded-lg border border-slate-200">{{ orderAntri.length }}</span>
						</div>
						<div class="flex-1 overflow-y-auto custom-scrollbar pr-1 space-y-3">
							<div v-for="order in orderAntri" :key="order.id" class="bg-slate-50 p-4 rounded-2xl border border-slate-200/60 hover:border-indigo-300 transition-colors group">
								<div class="flex justify-between items-start mb-2">
									<span class="text-[9px] font-black uppercase tracking-widest text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded border border-indigo-100">{{ order.invoice }}</span>
								</div>
								<p class="font-black text-sm text-slate-800 uppercase mb-3">{{ order.pelanggan }}</p>

								<div class="bg-white p-2.5 rounded-xl border border-slate-100 mb-3 space-y-2 shadow-sm">
									<div class="flex justify-between items-center gap-2">
										<span class="text-[8px] font-black uppercase tracking-widest text-slate-400">Layanan:</span>
										<span class="text-[9px] font-bold text-slate-700 truncate max-w-[120px]">{{ order.layanan || 'Paket Laundry' }} ({{ order.berat_kg || order.BeratKg || 0 }} {{ order.satuan_dasar || 'KG' }})</span>
									</div>
									<div class="flex justify-between items-center gap-2">
										<span class="text-[8px] font-black uppercase tracking-widest text-slate-400">Rak:</span>
										<span v-if="order.nomor_rak && order.nomor_rak !== '-'" class="text-[8px] font-black text-indigo-600 bg-indigo-50 px-1.5 py-0.5 rounded uppercase tracking-widest border border-indigo-100">
											{{ order.nomor_rak }}
										</span>
										<span v-else class="text-[8px] font-bold text-slate-400 italic">Belum Ada</span>
									</div>
								</div>

								<div class="flex items-center gap-1.5 text-[9px] font-bold text-rose-500 mb-3 bg-white w-fit px-2.5 py-1 rounded-lg border border-rose-100 shadow-sm">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<circle cx="12" cy="12" r="10" />
										<path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6l4 2" />
									</svg>
									Target: {{ formatDate(order.estimasi_waktu) }}
								</div>

								<div class="flex justify-between items-center mt-3 pt-3 border-t border-slate-200">
									<span :class="order.status_bayar === 'BELUM_LUNAS' ? 'bg-rose-100 text-rose-600' : 'bg-emerald-100 text-emerald-700'" class="text-[8px] font-black px-2 py-1 rounded uppercase tracking-widest border border-white/50">
										{{ order.status_bayar === 'BELUM_LUNAS' ? 'Piutang' : 'Lunas' }}
									</span>
									<button @click="updateStatusKanban(order, 'PROSES')" class="bg-slate-900 hover:bg-indigo-600 text-white text-[9px] font-black px-3.5 py-2 rounded-xl uppercase tracking-widest transition-all shadow-md flex items-center gap-1 active:scale-95">
										Mulai Cuci
										<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>
									</button>
								</div>
							</div>
						</div>
					</div>

					<div class="bg-white p-4 rounded-[28px] border border-slate-200 flex flex-col h-[70vh] shadow-sm">
						<div class="flex items-center justify-between mb-5 px-2">
							<h3 class="font-black text-[11px] uppercase tracking-widest text-slate-800 flex items-center gap-2">
								<div class="w-2.5 h-2.5 rounded-full bg-sky-400 ring-4 ring-sky-50 animate-pulse"></div>
								Sedang Dicuci
							</h3>
							<span class="bg-slate-100 text-slate-600 text-[10px] font-black px-2.5 py-1 rounded-lg border border-slate-200">{{ orderProses.length }}</span>
						</div>
						<div class="flex-1 overflow-y-auto custom-scrollbar pr-1 space-y-3">
							<div v-for="order in orderProses" :key="order.id" class="bg-sky-50/30 p-4 rounded-2xl border border-sky-100 hover:border-sky-300 transition-colors group">
								<div class="flex justify-between items-start mb-2">
									<span class="text-[9px] font-black uppercase tracking-widest text-sky-600 bg-white px-2 py-0.5 rounded border border-sky-100">{{ order.invoice }}</span>
								</div>
								<p class="font-black text-sm text-slate-800 uppercase mb-3">{{ order.pelanggan }}</p>

								<div class="bg-white p-2.5 rounded-xl border border-slate-100 mb-3 space-y-2 shadow-sm">
									<div class="flex justify-between items-center gap-2">
										<span class="text-[8px] font-black uppercase tracking-widest text-slate-400">Layanan:</span>
										<span class="text-[9px] font-bold text-slate-700 truncate max-w-[120px]">{{ order.layanan || 'Paket Laundry' }} ({{ order.berat_kg || order.BeratKg || 0 }} {{ order.satuan_dasar || 'KG' }})</span>
									</div>
									<div class="flex justify-between items-center gap-2">
										<span class="text-[8px] font-black uppercase tracking-widest text-slate-400">Rak:</span>
										<span v-if="order.nomor_rak && order.nomor_rak !== '-'" class="text-[8px] font-black text-indigo-600 bg-indigo-50 px-1.5 py-0.5 rounded uppercase tracking-widest border border-indigo-100">
											{{ order.nomor_rak }}
										</span>
										<span v-else class="text-[8px] font-bold text-slate-400 italic">Belum Ada</span>
									</div>
								</div>

								<div class="flex items-center gap-1.5 text-[9px] font-bold text-sky-600 mb-3 bg-white w-fit px-2.5 py-1 rounded-lg border border-sky-100 shadow-sm">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<circle cx="12" cy="12" r="10" />
										<path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6l4 2" />
									</svg>
									Target: {{ formatDate(order.estimasi_waktu) }}
								</div>

								<div class="flex justify-between items-center mt-3 pt-3 border-t border-sky-100/50">
									<button @click="updateStatusKanban(order, 'ANTRI')" class="text-[9px] font-black text-slate-400 hover:text-rose-500 uppercase transition-colors tracking-widest active:scale-95">Batalkan</button>
									<button @click="updateStatusKanban(order, 'SELESAI')" class="bg-sky-600 hover:bg-emerald-500 text-white text-[9px] font-black px-3.5 py-2 rounded-xl uppercase tracking-widest transition-all shadow-md flex items-center gap-1 active:scale-95">
										Selesai
										<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
									</button>
								</div>
							</div>
						</div>
					</div>

					<div class="bg-white p-4 rounded-[28px] border border-slate-200 flex flex-col h-[70vh] shadow-sm">
						<div class="flex items-center justify-between mb-5 px-2">
							<h3 class="font-black text-[11px] uppercase tracking-widest text-slate-800 flex items-center gap-2">
								<div class="w-2.5 h-2.5 rounded-full bg-emerald-400 ring-4 ring-emerald-50"></div>
								Siap Diambil
							</h3>
							<span class="bg-slate-100 text-slate-600 text-[10px] font-black px-2.5 py-1 rounded-lg border border-slate-200">{{ orderSelesai.length }}</span>
						</div>
						<div class="flex-1 overflow-y-auto custom-scrollbar pr-1 space-y-3">
							<div v-for="order in orderSelesai" :key="order.id" class="bg-emerald-50/30 p-4 rounded-2xl border border-emerald-100 hover:border-emerald-300 transition-colors group relative overflow-hidden">
								<div class="flex justify-between items-start mb-2 relative z-10">
									<span class="text-[9px] font-black uppercase tracking-widest text-emerald-700 bg-white px-2 py-0.5 rounded border border-emerald-200">{{ order.invoice }}</span>
								</div>
								<p class="font-black text-sm text-slate-800 uppercase mb-3 relative z-10">{{ order.pelanggan }}</p>

								<div class="bg-white p-2.5 rounded-xl border border-emerald-100/50 mb-3 space-y-2 shadow-sm relative z-10">
									<div class="flex justify-between items-center gap-2">
										<span class="text-[8px] font-black uppercase tracking-widest text-slate-400">Layanan:</span>
										<span class="text-[9px] font-bold text-slate-700 truncate max-w-[120px]">{{ order.layanan || 'Paket Laundry' }} ({{ order.berat_kg || order.BeratKg || 0 }} {{ order.satuan_dasar || 'KG' }})</span>
									</div>
									<div class="flex justify-between items-center gap-2">
										<span class="text-[8px] font-black uppercase tracking-widest text-slate-400">Rak:</span>
										<span v-if="order.nomor_rak && order.nomor_rak !== '-'" class="text-[8px] font-black text-indigo-600 bg-indigo-50 px-1.5 py-0.5 rounded uppercase tracking-widest border border-indigo-100 shadow-sm">
											{{ order.nomor_rak }}
										</span>
										<span v-else class="text-[8px] font-bold text-slate-400 italic">Belum Ada</span>
									</div>
								</div>

								<div class="flex justify-end items-center mt-3 pt-3 border-t border-emerald-100/50 relative z-10">
									<button v-if="order.status_bayar === 'BELUM_LUNAS'" @click="prosesPengambilan(order)" class="w-full bg-rose-500 hover:bg-rose-600 text-white text-[9px] font-black px-3 py-3 rounded-xl uppercase tracking-[0.2em] transition-all shadow-lg shadow-rose-200 animate-pulse flex items-center justify-center gap-2">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
										Bayar & Serahkan
									</button>
									<button v-else @click="prosesPengambilan(order)" class="w-full bg-slate-800 hover:bg-slate-900 text-white text-[9px] font-black px-3 py-3 rounded-xl uppercase tracking-widest transition-all shadow-md flex items-center justify-center gap-2 active:scale-95">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
										Serahkan (Lunas)
									</button>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</SidebarLaundry>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
	width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
	background: #94a3b8;
}
@keyframes fadeIn {
	from {
		opacity: 0;
		transform: translateY(10px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}
.animate-fade-in {
	animation: fadeIn 0.3s ease-out;
}
</style>
