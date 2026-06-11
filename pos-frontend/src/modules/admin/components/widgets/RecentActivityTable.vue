<script setup>
// 🚀 Menggunakan prop agar data bersifat dinamis dan dapat dikendalikan dari halaman utama (reactive)
defineProps({
	activities: {
		type: Array,
		required: true,
	},
});
</script>

<template>
	<div class="bg-[#131B2E] border border-slate-800 rounded-[24px] sm:rounded-[28px] p-5 sm:p-6 shadow-xl flex flex-col justify-between h-full w-full overflow-hidden">
		<div>
			<div class="flex items-center justify-between mb-4 border-b border-slate-800 pb-4">
				<div class="min-w-0 flex-1 pr-2">
					<h3 class="text-xs sm:text-sm font-black text-white uppercase tracking-wide truncate">Live Audit Logs Activity</h3>
					<p class="text-[9px] sm:text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-0.5 truncate">Realtime Global Event Stream Detector</p>
				</div>
				<span class="w-2 h-2 rounded-full bg-indigo-500 animate-ping shrink-0"></span>
			</div>

			<div class="space-y-3 max-h-[260px] overflow-y-auto pr-1 custom-scrollbar">
				<template v-if="activities && activities.length > 0">
					<div v-for="log in activities" :key="log.id" class="flex items-start gap-2.5 p-2.5 sm:p-3 bg-[#0B0F19]/60 border border-slate-800/40 rounded-xl transition-all hover:bg-[#0B0F19]">
						<span class="text-[9px] sm:text-[10px] font-black text-indigo-400 bg-indigo-500/10 px-1.5 py-0.5 rounded-md mt-0.5 shrink-0">
							{{ log.time }}
						</span>

						<p class="text-[11px] sm:text-xs text-slate-300 font-medium leading-relaxed flex-1 break-words overflow-hidden">
							<span class="text-white font-bold shrink-0">{{ log.user }}</span>
							— {{ log.action }}
						</p>
					</div>
				</template>

				<div v-else class="flex flex-col items-center justify-center py-8 text-center opacity-60">
					<svg class="w-8 h-8 text-slate-600 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
					</svg>
					<p class="text-[10px] sm:text-xs text-slate-400 font-bold uppercase tracking-widest">Belum Ada Aktivitas Terekam</p>
				</div>
			</div>
		</div>
	</div>
</template>
