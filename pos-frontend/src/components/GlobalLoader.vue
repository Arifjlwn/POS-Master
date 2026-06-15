<script setup>
defineProps({
	isActive: {
		type: Boolean,
		default: false,
	},
});
</script>

<template>
	<transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0 backdrop-blur-none" enter-to-class="opacity-100 backdrop-blur-sm" leave-active-class="transition duration-300 ease-in" leave-from-class="opacity-100 backdrop-blur-sm" leave-to-class="opacity-0 backdrop-blur-none">
		<div v-if="isActive" class="fixed inset-0 z-[9999] flex flex-col items-center justify-center bg-slate-950/40 backdrop-blur-md no-print pointer-events-auto select-none">
			<div class="relative flex flex-col items-center justify-center p-8 rounded-[32px] bg-white/70 border border-white/50 shadow-2xl shadow-slate-900/10 max-w-xs w-full text-center animate-[scaleUp_0.3s_cubic-bezier(0.34,1.56,0.64,1)]">
				<div class="relative w-16 h-16 mb-4 flex items-center justify-center">
					<div class="absolute inset-0 border-4 border-slate-100 border-t-indigo-600 rounded-full animate-spin"></div>
					<div class="absolute w-12 h-12 border-4 border-transparent border-t-violet-500 border-b-violet-500 rounded-full animate-[spin_1.2s_linear_infinite_reverse]"></div>
					<div class="absolute w-8 h-8 border-4 border-slate-100 border-l-blue-500 rounded-full animate-[spin_0.8s_linear_infinite]"></div>
				</div>

				<h3 class="text-xs font-black text-slate-800 uppercase tracking-[0.25em] pl-[0.25em]">ARZURA POS</h3>
				<div class="mt-1.5 flex items-center gap-1.5">
					<span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest animate-pulse">Sinkronisasi Server...</span>
				</div>

				<div class="absolute -inset-1 rounded-[36px] bg-gradient-to-r from-blue-500/10 to-indigo-500/10 blur-xl -z-10 animate-pulse"></div>
			</div>
		</div>
	</transition>
</template>

<style scoped>
@keyframes scaleUp {
	from {
		opacity: 0;
		transform: scale(0.92);
	}
	to {
		opacity: 1;
		transform: scale(1);
	}
}

/* Custom cubic-bezier transition untuk transisi blur kasta tertinggi bray */
.backdrop-blur-sm {
	backdrop-filter: blur(8px) !important;
	-webkit-backdrop-filter: blur(8px) !important;
}
</style>
