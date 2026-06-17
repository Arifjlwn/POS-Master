<script setup>
defineProps({
	showModal: Boolean,
	totalTagihan: Number,
	qrisStoreUrl: String,
	buktiTransferData: String,
	formatRupiah: Function,
});
defineEmits(['close', 'open-camera', 'remove-photo', 'confirm']);
</script>

<template>
	<Teleport to="body">
		<div v-if="showModal" class="fixed inset-0 z-[9999] flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 animate-[fadeIn_0.2s_ease-out]">
			<div class="bg-white rounded-3xl shadow-2xl w-full max-w-sm md:max-w-2xl overflow-hidden flex flex-col border border-slate-200">
				<div class="bg-slate-900 p-4 text-center text-white shrink-0">
					<h3 class="font-black text-xs uppercase tracking-widest">SCAN QRIS PEMBAYARAN TOKO</h3>
					<p class="text-xs font-bold text-slate-400 mt-1">Total Tagihan: {{ formatRupiah(totalTagihan) }}</p>
				</div>
				<div class="p-5 grid grid-cols-1 md:grid-cols-2 gap-5 max-h-[70vh] overflow-y-auto custom-scrollbar items-center bg-slate-50/30 flex-1">
					<div class="bg-white rounded-2xl border border-slate-200 p-4 flex items-center justify-center shadow-sm">
						<img :src="qrisStoreUrl || 'https://upload.wikimedia.org/wikipedia/commons/d/d0/QR_code_for_mobile_English_Wikipedia.svg'" alt="QRIS Toko" class="w-full max-w-[170px] aspect-square object-contain mix-blend-multiply" />
					</div>
					<div class="w-full border border-dashed border-slate-300 rounded-2xl p-4 bg-white flex flex-col items-center justify-center relative min-h-[160px] shadow-sm">
						<button v-if="!buktiTransferData" @click="$emit('open-camera')" class="flex flex-col items-center text-slate-500 hover:text-slate-800 transition-colors my-2">
							<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
								<rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
								<circle cx="8.5" cy="8.5" r="1.5" />
								<polyline points="21 15 16 10 5 21" />
							</svg>
							<span class="text-[9px] font-black uppercase tracking-widest text-center">Foto Bukti Transfer HP Pelanggan</span>
						</button>
						<div v-else class="w-full h-full relative z-20 rounded-xl overflow-hidden border border-slate-200">
							<img :src="buktiTransferData" class="w-full h-36 object-cover" />
							<button @click="$emit('remove-photo')" class="absolute top-2 right-2 bg-rose-600 text-white w-5 h-5 rounded-full flex items-center justify-center shadow-lg">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
									<line x1="18" y1="6" x2="6" y2="18" />
									<line x1="6" y1="6" x2="18" y2="18" />
								</svg>
							</button>
						</div>
					</div>
				</div>
				<div class="p-4 bg-white border-t border-slate-200 flex gap-2 shrink-0">
					<button @click="$emit('close')" class="w-1/3 py-2.5 text-slate-500 font-black text-xs uppercase tracking-widest hover:bg-slate-100 rounded-xl transition-all">Batal</button>
					<button @click="$emit('confirm')" :disabled="!buktiTransferData" class="w-2/3 py-2.5 bg-slate-900 text-white font-black text-xs uppercase tracking-widest rounded-xl disabled:opacity-40 flex justify-center items-center gap-2 transition-all">Sahkan Bayar</button>
				</div>
			</div>
		</div>
	</Teleport>
</template>
