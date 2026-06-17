<script setup>
defineProps({
	showModal: Boolean,
	perfumes: Array,
	formatRupiah: Function,
});
defineEmits(['close', 'toggle-status']);
</script>

<template>
	<Teleport to="body">
		<div v-if="showModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]">
			<div class="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden flex flex-col border border-slate-200">
				<div class="bg-slate-900 p-5 text-center text-white flex justify-between items-center shrink-0">
					<h3 class="font-black text-xs uppercase tracking-widest">Saklar Ketersediaan Parfum</h3>
					<button @click="$emit('close')" class="text-white hover:opacity-70">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
							<line x1="18" y1="6" x2="6" y2="18" />
							<line x1="6" y1="6" x2="18" y2="18" />
						</svg>
					</button>
				</div>
				<div class="p-4 max-h-[60vh] overflow-y-auto custom-scrollbar space-y-2 bg-slate-50/50 flex-1">
					<div v-if="perfumes.length === 0" class="text-center py-6 text-xs font-bold text-slate-400 uppercase">Belum ada varian parfum di Master Layanan.</div>
					<div v-for="perfume in perfumes" :key="perfume.id" class="p-3 border border-slate-200 rounded-xl flex items-center justify-between bg-white shadow-sm">
						<div>
							<h4 class="font-black text-slate-800 text-xs uppercase">🌸 {{ perfume.nama }}</h4>
							<p class="text-[9px] font-bold text-slate-400 mt-0.5">Charge: {{ perfume.harga > 0 ? formatRupiah(perfume.harga) : 'Gratis' }}</p>
						</div>
						<button @click="$emit('toggle-status', perfume)" :class="perfume.status === 'Tersedia' ? 'bg-emerald-500' : 'bg-slate-200'" class="w-12 h-6 rounded-full p-0.5 transition-all duration-300 flex items-center relative shadow-inner">
							<div :class="perfume.status === 'Tersedia' ? 'translate-x-6' : 'translate-x-0'" class="w-5 h-5 bg-white rounded-full shadow-md transition-transform duration-300"></div>
						</button>
					</div>
				</div>
				<div class="p-4 bg-white border-t border-slate-200 text-center shrink-0">
					<button @click="$emit('close')" class="bg-slate-900 text-white font-black text-[10px] tracking-widest px-6 py-2.5 rounded-xl uppercase">Selesai Monitor</button>
				</div>
			</div>
		</div>
	</Teleport>
</template>
