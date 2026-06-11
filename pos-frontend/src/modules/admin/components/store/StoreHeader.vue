<script setup>
defineProps({
	isLoading: Boolean,
});
defineEmits(['refresh']);
</script>

<template>
	<div class="flex justify-between items-center mb-8">
		<div>
			<h1 class="text-2xl font-black tracking-tight uppercase">Manajemen Mitra Toko</h1>
			<p class="text-slate-400 text-xs mt-1">Total ruko retail terdaftar di platform SaaS ARZURA POS.</p>
		</div>
		<button @click="$emit('refresh')" :disabled="isLoading" class="px-4 py-2.5 bg-slate-800 hover:bg-slate-700 disabled:opacity-50 rounded-xl text-xs font-bold transition-all flex items-center gap-2">
			<svg :class="{ 'animate-spin': isLoading }" class="w-4 h-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
			</svg>
			Refresh Data
		</button>
	</div>
</template>
