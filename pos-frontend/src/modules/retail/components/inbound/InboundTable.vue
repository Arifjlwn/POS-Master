<script setup>
defineProps({
	cartLPB: Array,
	isOwner: Boolean,
	hitungTotalStok: Function,
	hitungModalPerPcs: Function,
});

const emit = defineEmits(['remove']);

const formatNumber = (val) => {
	if (val === null || val === undefined || val === '') return '';
	return Number(val).toLocaleString('id-ID');
};

const handleInputHarga = (item, event) => {
	let raw = event.target.value.replace(/\D/g, '');
	let num = raw ? parseInt(raw, 10) : 0;
	item.harga_beli_input = num;
	event.target.value = raw ? num.toLocaleString('id-ID') : '';
};

// 🚀 TAMBAHAN BARU: Pelindung Anti-Minus & Desimal untuk QTY
const handleInputQty = (item, field, event) => {
	let raw = event.target.value.replace(/\D/g, ''); // Basmi semua selain angka murni
	let num = raw ? parseInt(raw, 10) : 0;
	item[field] = num;
	event.target.value = num; // Balikin ke UI dalam bentuk angka bersih
};
</script>

<template>
	<table class="w-full text-left whitespace-nowrap border-collapse">
		<thead>
			<tr class="bg-slate-50/50 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] border-b border-slate-100">
				<th class="p-6 w-[30%]">Nama Barang & Kemasan</th>
				<th class="p-6 text-center">Qty Masuk (Berdasarkan Kemasan)</th>
				<th class="p-6 text-center bg-blue-50/30 text-blue-600">
					Total Masuk
					<br />
					(Satuan Dasar)
				</th>
				<th class="p-6 text-right">
					Harga Total Tagihan
					<br />
					(Untuk Keseluruhan Qty)
				</th>
				<th v-if="isOwner" class="p-6 text-right">
					HPP Baru
					<br />
					(Per Satuan Dasar)
				</th>
				<th class="p-6 text-center">Aksi</th>
			</tr>
		</thead>
		<tbody class="divide-y divide-slate-50">
			<tr v-if="cartLPB.length === 0">
				<td :colspan="isOwner ? 6 : 5" class="p-16 text-center text-slate-400 font-black text-[10px] uppercase tracking-[0.3em] opacity-50">Belum ada barang di scan</td>
			</tr>
			<tr v-for="(item, index) in cartLPB" :key="item.product_id" :class="isOwner && hitungModalPerPcs(item) > item.harga_jual_saat_ini ? 'bg-red-50/50' : 'hover:bg-slate-50/50'">
				<td class="p-6">
					<div class="font-black text-slate-800 text-sm uppercase tracking-tight">{{ item.nama_produk }}</div>

					<div class="mt-1 flex flex-wrap gap-1">
						<span v-if="item.is_nested" class="px-2 py-0.5 bg-indigo-50 border border-indigo-100 text-indigo-600 rounded text-[9px] font-black uppercase tracking-widest">3 Lapis ({{ item.satuan_besar || 'BESAR' }} ➝ {{ item.satuan_tengah || 'TENGAH' }} ➝ {{ item.satuan_dasar }})</span>
						<span v-else-if="item.has_satuan_besar" class="px-2 py-0.5 bg-emerald-50 border border-emerald-100 text-emerald-600 rounded text-[9px] font-black uppercase tracking-widest">2 Lapis ({{ item.satuan_besar || 'BESAR' }} ➝ {{ item.satuan_dasar }})</span>
						<span v-else class="px-2 py-0.5 bg-slate-100 border border-slate-200 text-slate-500 rounded text-[9px] font-black uppercase tracking-widest">Eceran ({{ item.satuan_dasar }})</span>
					</div>

					<div v-if="isOwner && hitungModalPerPcs(item) > item.harga_jual_saat_ini" class="mt-2 text-[8px] font-black text-red-600 animate-pulse uppercase bg-red-100 px-2 py-1 rounded inline-block">⚠️ HPP melebihi Harga Jual!</div>
				</td>

				<td class="p-6">
					<div class="flex flex-wrap items-center justify-center gap-2">
						<div v-if="item.has_satuan_besar" class="flex flex-col items-center gap-1">
							<input :value="item.qty_besar" @input="handleInputQty(item, 'qty_besar', $event)" type="text" inputmode="numeric" class="w-16 p-2 bg-white border-2 border-slate-200 rounded-xl text-center font-black text-slate-700 outline-none focus:border-indigo-500 shadow-inner" />
							<span class="text-[8px] font-black text-slate-400 uppercase">{{ item.satuan_besar || 'KARTON' }}</span>
						</div>

						<div v-if="item.is_nested" class="flex flex-col items-center gap-1">
							<span class="text-slate-300 font-bold mb-5">+</span>
						</div>
						<div v-if="item.is_nested" class="flex flex-col items-center gap-1">
							<input :value="item.qty_tengah" @input="handleInputQty(item, 'qty_tengah', $event)" type="text" inputmode="numeric" class="w-16 p-2 bg-white border-2 border-slate-200 rounded-xl text-center font-black text-emerald-600 outline-none focus:border-emerald-500 shadow-inner" />
							<span class="text-[8px] font-black text-emerald-400 uppercase">{{ item.satuan_tengah || 'BUNGKUS' }}</span>
						</div>

						<div class="flex flex-col items-center gap-1">
							<span v-if="item.has_satuan_besar" class="text-slate-300 font-bold mb-5">+</span>
						</div>
						<div class="flex flex-col items-center gap-1">
							<input :value="item.qty_dasar" @input="handleInputQty(item, 'qty_dasar', $event)" type="text" inputmode="numeric" class="w-16 p-2 bg-white border-2 border-slate-200 rounded-xl text-center font-black text-blue-600 outline-none focus:border-blue-500 shadow-inner" />
							<span class="text-[8px] font-black text-blue-400 uppercase">{{ item.satuan_dasar || 'PCS' }}</span>
						</div>
					</div>
				</td>

				<td class="p-6 text-center bg-blue-50/20">
					<div class="flex flex-col items-center">
						<span class="text-lg font-black text-blue-700 leading-none">{{ formatNumber(hitungTotalStok(item)) }}</span>
						<span class="text-[8px] font-black text-blue-400 uppercase mt-1 tracking-widest">{{ item.satuan_dasar }}</span>
					</div>
				</td>

				<td class="p-6 min-w-[200px]">
					<div class="relative w-full">
						<span class="absolute inset-y-0 left-0 pl-4 flex items-center text-xs font-black text-slate-400 italic">Rp</span>
						<input :value="formatNumber(item.harga_beli_input)" @input="handleInputHarga(item, $event)" type="text" inputmode="numeric" class="w-full pl-11 pr-4 py-3 bg-white border-2 border-slate-200 rounded-xl text-right font-black text-slate-800 outline-none focus:border-amber-500 transition-all text-sm shadow-inner" />
					</div>
					<div class="text-right mt-1.5 text-[8px] font-black text-indigo-400 uppercase italic">*Total uang keluar untuk {{ formatNumber(hitungTotalStok(item)) }} {{ item.satuan_dasar }}</div>
				</td>

				<td v-if="isOwner" class="p-6 text-right">
					<div class="text-base font-black tracking-tight text-emerald-600" :class="{ 'text-red-600': hitungModalPerPcs(item) > item.harga_jual_saat_ini }">Rp {{ formatNumber(hitungModalPerPcs(item)) }}</div>
					<div class="text-[8px] font-black text-slate-400 uppercase mt-1">/ {{ item.satuan_dasar }}</div>
				</td>

				<td class="p-6 text-center">
					<button @click="$emit('remove', index)" class="p-3 bg-slate-50 hover:bg-red-50 rounded-2xl transition-all group border border-transparent hover:border-red-100 shadow-sm hover:shadow-md hover:-translate-y-0.5">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-300 group-hover:text-red-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
					</button>
				</td>
			</tr>
		</tbody>
	</table>
</template>
