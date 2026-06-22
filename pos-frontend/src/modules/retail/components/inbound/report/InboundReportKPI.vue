<script setup>
defineProps({
	stats: { type: Object, required: true },
	formatRupiah: { type: Function, required: true },
});
</script>

<template>
	<div class="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-8">
		<div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex items-center gap-5">
			<div class="w-14 h-14 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center shrink-0">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
			</div>
			<div>
				<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Nilai LPB</p>
				<h3 class="text-xl font-black text-slate-800 tracking-tight">{{ formatRupiah(stats.totalModal) }}</h3>
			</div>
		</div>

		<div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex items-center gap-5">
			<div class="w-14 h-14 bg-sky-50 text-sky-600 rounded-2xl flex items-center justify-center shrink-0">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
			</div>
			<div>
				<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Barang Masuk</p>
				<h3 class="text-xl font-black text-slate-800 tracking-tight">
					{{ stats.totalBarang }}
					<span class="text-xs font-bold text-slate-500">Unit</span>
				</h3>
			</div>
		</div>

		<div class="bg-white p-6 rounded-[24px] border border-slate-200 shadow-sm flex items-center gap-5">
			<div class="w-14 h-14 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center shrink-0">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
			</div>
			<div>
				<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Faktur</p>
				<h3 class="text-xl font-black text-slate-800 tracking-tight">
					{{ stats.totalFaktur }}
					<span class="text-xs font-bold text-slate-500">Dokumen</span>
				</h3>
			</div>
		</div>
	</div>
</template>
