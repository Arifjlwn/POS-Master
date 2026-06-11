<script setup>
import { computed } from 'vue';

const emit = defineEmits(['toggle-sidebar']);

// 🚀 TERIMA PROPS REALTIME: Navbar sekarang bisa mantau kesehatan server dari dashboard pusat bray!
const props = defineProps({
	dbStatus: {
		type: String,
		default: 'Online & Stabil',
	},
	latency: {
		type: String,
		default: '0ms',
	},
});

// Ambil data nama real dari session login admin
const founderName = computed(() => {
	return localStorage.getItem('name') || 'Super Admin';
});

// JURUS DINAMIS AVATAR: Otomatis generate inisial
const avatarInitials = computed(() => {
	const name = founderName.value.trim();
	if (!name) return 'SA';

	const words = name.split(/\s+/);
	if (words.length === 1) {
		return words[0].substring(0, 2).toUpperCase();
	}
	return (words[0][0] + words[1][0]).toUpperCase();
});

// 🚀 JURUS DINAMIS INDIKATOR SERVER: Menganalisis kesehatan peladen secara riil!
const isSystemHealthy = computed(() => {
	if (props.dbStatus !== 'Online & Stabil') return 'down';

	// Potong string "ms" buat ambil angka latensinya doang bray
	const msValue = parseInt(props.latency.replace('ms', '')) || 0;
	if (msValue > 300) return 'overload';

	return 'healthy';
});
</script>

<template>
	<header class="h-20 bg-[#131B2E] border-b border-slate-800 flex items-center justify-between px-4 lg:px-8 z-30 sticky top-0 backdrop-blur-md bg-opacity-90 transition-all">
		<div class="flex items-center gap-4">
			<button @click="emit('toggle-sidebar')" class="p-2 rounded-xl text-slate-400 hover:bg-slate-800 hover:text-white transition-colors">
				<svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
				</svg>
			</button>
			<h2 class="text-white font-black tracking-tight text-lg hidden sm:block">Mission Control</h2>
		</div>

		<div class="flex items-center gap-4">
			<div v-if="isSystemHealthy === 'healthy'" class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
				<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
				<span class="text-emerald-400 text-[10px] font-black uppercase tracking-widest">Sistem Stabil ({{ props.latency }})</span>
			</div>

			<div v-else-if="isSystemHealthy === 'overload'" class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-amber-500/10 border border-amber-500/20 rounded-full">
				<span class="w-2 h-2 rounded-full bg-amber-500 animate-ping"></span>
				<span class="text-amber-400 text-[10px] font-black uppercase tracking-widest">Server Overload ({{ props.latency }})</span>
			</div>

			<div v-else class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-rose-500/10 border border-rose-500/20 rounded-full">
				<span class="w-2 h-2 rounded-full bg-rose-500 animate-bounce"></span>
				<span class="text-rose-400 text-[10px] font-black uppercase tracking-widest">Sistem Down !</span>
			</div>

			<div class="flex items-center gap-3 pl-4 border-l border-slate-700">
				<div class="text-right hidden sm:block">
					<p class="text-white text-sm font-bold">{{ founderName }}</p>
					<p class="text-indigo-400 text-[10px] font-black uppercase tracking-widest">Founder / Root</p>
				</div>
				<div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-600 to-purple-600 p-[2px]">
					<div class="w-full h-full bg-[#131B2E] rounded-[10px] flex items-center justify-center">
						<span class="text-white font-black text-sm select-none">{{ avatarInitials }}</span>
					</div>
				</div>
			</div>
		</div>
	</header>
</template>
