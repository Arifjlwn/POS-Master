<script setup>
defineProps({
	activeTab: { type: String, required: true },
	isEditing: { type: Boolean, required: true },
	formJasa: { type: Object, required: true },
	formParfum: { type: Object, required: true },
});

const emit = defineEmits(['save', 'cancel']);
</script>

<template>
	<div class="bg-white p-5 sm:p-6 md:p-8 rounded-[24px] border border-slate-200 shadow-xl shadow-slate-200/50 mb-8 animate-[fadeInDown_0.3s_ease-out]">
		<div class="flex items-center gap-3 mb-6 border-b border-slate-100 pb-4">
			<div class="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
					<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
					<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
				</svg>
			</div>
			<h3 class="text-sm font-black text-slate-800 uppercase tracking-widest">
				{{ activeTab === 'jasa' ? (isEditing ? 'Edit Layanan Jasa' : 'Input Jasa Baru') : 'Input Parfum Baru' }}
			</h3>
		</div>

		<div v-if="activeTab === 'jasa'" class="grid grid-cols-1 sm:grid-cols-2 gap-5 md:gap-6">
			<div>
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Layanan / Paket</label>
				<input v-model="formJasa.nama_produk" type="text" placeholder="Contoh: Cuci Karpet Tebal" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300" />
			</div>
			<div>
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Harga Jual</label>
				<input v-model="formJasa.harga_jual" type="number" placeholder="0" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-black text-sm text-slate-800 transition-all placeholder:text-slate-300" />
			</div>
			<div>
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Tipe Satuan</label>
				<select v-model="formJasa.satuan_dasar" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 cursor-pointer transition-all">
					<option value="KG">Kiloan (KG)</option>
					<option value="PCS">Satuan (PCS / Helai)</option>
					<option value="M2">Meter Persegi (M2)</option>
				</select>
			</div>
			<div>
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Estimasi Waktu Pengerjaan</label>
				<select v-model="formJasa.estimasi" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 cursor-pointer transition-all">
					<option value="3 Jam">Express (3 Jam)</option>
					<option value="6 Jam">Express (6 Jam)</option>
					<option value="1 Hari">Kilat (1 Hari)</option>
					<option value="2 Hari">Reguler (2 Hari)</option>
					<option value="3 Hari">Standar (3 Hari)</option>
					<option value="7 Hari">Khusus (7 Hari)</option>
				</select>
			</div>
		</div>

		<div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-5 md:gap-6">
			<div>
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Nama Varian Parfum</label>
				<input v-model="formParfum.nama" type="text" placeholder="Contoh: Aroma Sakura Premium, Downy Mistik..." class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-sm text-slate-800 transition-all placeholder:text-slate-300" />
			</div>
			<div>
				<label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Biaya Tambahan (Charge)</label>
				<input v-model="formParfum.harga" type="number" placeholder="Set 0 jika gratis / bawaan" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-200 rounded-xl outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-black text-sm text-slate-800 transition-all placeholder:text-slate-300" />
			</div>
		</div>

		<div class="mt-8 flex justify-end gap-3">
			<button type="button" @click="emit('cancel')" class="px-6 py-4 bg-slate-100 hover:bg-slate-200 text-slate-600 rounded-xl font-black text-xs uppercase tracking-widest transition-all">Batal</button>
			<button type="button" @click="emit('save')" class="w-full sm:w-auto bg-slate-900 hover:bg-slate-800 text-white px-8 py-4 rounded-xl font-black text-xs uppercase tracking-[0.15em] transition-all active:scale-95 shadow-xl shadow-slate-900/20 flex items-center justify-center gap-2">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
					<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
					<polyline points="17 21 17 13 7 13 7 21" />
					<polyline points="7 3 7 8 15 8" />
				</svg>
				{{ isEditing ? 'Simpan Perubahan' : 'Simpan Data Katalog' }}
			</button>
		</div>
	</div>
</template>
