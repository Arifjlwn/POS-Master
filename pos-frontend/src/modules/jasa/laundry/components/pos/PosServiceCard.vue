<script setup>
defineProps({
	service: { type: Object, required: true },
	formatRupiah: { type: Function, required: true },
});
defineEmits(['click']);
</script>

<template>
	<div @click="$emit('click')" class="bg-white p-4 md:p-5 rounded-[24px] border border-slate-200 hover:border-slate-800 hover:shadow-xl hover:shadow-slate-100 transition-all active:scale-95 flex flex-col justify-between h-32 md:h-36 relative overflow-hidden group">
		<div class="z-10">
			<span class="text-[9px] font-black text-slate-700 bg-slate-100 px-2.5 py-1 rounded-lg uppercase tracking-wider mb-2 inline-block border border-slate-200">EST: {{ service.estimasi_durasi || '1' }} {{ service.estimasi_satuan || 'Hari' }}</span>
			<h3 class="text-xs md:text-sm font-black text-slate-800 uppercase leading-tight line-clamp-2 pr-2">
				{{ service.nama_produk }}
			</h3>
		</div>
		<p class="text-sm md:text-base font-black text-slate-900 z-10">
			{{ formatRupiah(service.harga_jual) }}
			<span class="text-[10px] text-slate-400 font-bold uppercase">/ {{ service.satuan_dasar || 'KG' }}</span>
		</p>
	</div>
</template>
