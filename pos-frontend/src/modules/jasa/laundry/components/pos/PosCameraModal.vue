<script setup>
defineProps({ isCameraOpen: Boolean, cameraTarget: String });
defineEmits(['jepret', 'close', 'video-bind']);
</script>

<template>
	<div v-if="isCameraOpen" class="fixed inset-0 z-[99999] bg-black/90 backdrop-blur-md flex flex-col items-center justify-center p-4">
		<div class="w-full max-w-md bg-slate-900 rounded-3xl overflow-hidden border border-slate-800 shadow-2xl relative flex flex-col">
			<div class="p-4 border-b border-slate-800 text-center text-white font-black text-xs uppercase tracking-widest shrink-0">Kamera Rana Dokumentasi Sesi</div>
			<div class="relative bg-black flex-1 min-h-[240px]">
				<video
					:ref="
						(el) => {
							if (el) $emit('video-bind', el);
						}
					"
					autoplay
					playsinline
					class="w-full h-64 object-cover"></video>
			</div>
			<div class="p-4 bg-slate-900 border-t border-slate-800 flex justify-center gap-3 shrink-0">
				<button type="button" @click="$emit('close')" class="px-6 py-3 bg-slate-800 hover:bg-slate-700 text-slate-400 font-black text-xs uppercase tracking-widest rounded-xl transition-all">Batal</button>
				<button type="button" @click="$emit('jepret')" class="px-8 py-3 bg-white text-slate-900 font-black text-xs uppercase tracking-widest rounded-xl shadow-lg transition-all active:scale-95 flex items-center gap-2">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
						<circle cx="12" cy="12" r="3" />
						<path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
					</svg>
					<span>Tangkap Foto</span>
				</button>
			</div>
		</div>
	</div>
</template>
