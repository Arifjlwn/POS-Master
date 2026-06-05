<script setup>
import { computed, ref } from 'vue';
const props = defineProps({ show: Boolean, isEditing: Boolean, isSubmitting: Boolean, form: Object, categories: Array, imagePreview: String, stokDalamKarton: Number, eceranTambahan: Number });
const emit = defineEmits(['close', 'submit', 'start-scanner', 'file-change', 'update:stokDalamKarton', 'update:eceranTambahan']);
const imageInput = ref(null);
const showCategoryDropdown = ref(false);
const triggerImageUpload = () => imageInput.value.click();
const filteredCategories = computed(() => {
	if (!props.form.category) return props.categories;
	const q = props.form.category.toLowerCase();
	return props.categories.filter((c) => c.toLowerCase().includes(q));
});
const selectCategory = (cat) => {
	props.form.category = cat;
	showCategoryDropdown.value = false;
};
const formatNumber = (val) => (val === null || val === undefined || val === '' ? '' : Number(val).toLocaleString('id-ID'));
const handleInputForm = (key, event) => {
	let raw = event.target.value.replace(/\D/g, ''),
		num = raw ? parseInt(raw, 10) : 0;
	props.form[key] = num;
	event.target.value = raw ? num.toLocaleString('id-ID') : '';
};
const handleInputEmit = (emitName, event) => {
	let raw = event.target.value.replace(/\D/g, ''),
		num = raw ? parseInt(raw, 10) : 0;
	emit(emitName, num);
	event.target.value = raw ? num.toLocaleString('id-ID') : '';
};
const enableNestedUom = () => {
	props.form.is_nested_uom = true;
	if (props.form.satuan_dasar === 'ML') {
		if (!props.form.satuan_tengah) props.form.satuan_tengah = 'LITER';
		if (!props.form.isi_tengah_ke_dasar) props.form.isi_tengah_ke_dasar = 1000;
	} else if (props.form.satuan_dasar === 'GRAM') {
		if (!props.form.satuan_tengah) props.form.satuan_tengah = 'KG';
		if (!props.form.isi_tengah_ke_dasar) props.form.isi_tengah_ke_dasar = 1000;
	}
};
</script>

<template>
	<div v-if="show" class="fixed inset-0 bg-slate-900/80 flex items-center justify-center z-[100] p-4 backdrop-blur-sm no-print">
		<div class="bg-white rounded-[32px] shadow-2xl w-full max-w-xl flex flex-col max-h-[90vh] overflow-hidden border border-slate-100">
			<div class="p-6 border-b border-slate-50 bg-slate-50/50 flex justify-between items-center shrink-0">
				<h2 class="text-xl font-black text-slate-800 uppercase italic">{{ isEditing ? 'Edit Data Produk' : 'Registrasi Produk Baru' }}</h2>
				<button @click="emit('close')" class="p-2 rounded-xl bg-white text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all border border-slate-100 shadow-sm">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
						<line x1="18" y1="6" x2="6" y2="18" />
						<line x1="6" y1="6" x2="18" y2="18" />
					</svg>
				</button>
			</div>
			<div class="p-6 md:p-8 overflow-y-auto custom-scrollbar">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-5">
					<div class="md:col-span-2 flex items-center gap-5 p-4 rounded-[24px] bg-slate-50 border border-slate-100 mb-2">
						<div @click="triggerImageUpload" class="w-16 h-16 rounded-[18px] border-2 border-dashed border-slate-300 flex items-center justify-center bg-white cursor-pointer overflow-hidden shadow-inner shrink-0 hover:border-blue-400 transition-colors">
							<img v-if="imagePreview" :src="imagePreview" class="w-full h-full object-cover" />
							<svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-slate-300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z" />
								<circle cx="12" cy="13" r="3" />
							</svg>
						</div>
						<div class="flex-1">
							<input v-model="form.name" type="text" placeholder="NAMA BARANG..." class="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-black text-sm uppercase mb-2 text-slate-800 transition-colors" />
							<div class="flex gap-2">
								<input v-model="form.sku" type="text" placeholder="BARCODE / SKU..." class="w-full px-4 py-2 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-xs uppercase text-slate-800 transition-colors" />
								<button @click.prevent="emit('start-scanner')" class="px-4 bg-blue-50 text-blue-600 hover:bg-blue-600 hover:text-white rounded-xl border border-blue-100 transition-colors shadow-sm">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
										<path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
										<path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
									</svg>
								</button>
							</div>
						</div>
						<input type="file" ref="imageInput" @change="(e) => emit('file-change', e)" accept="image/*" class="hidden" />
					</div>
					<div class="md:col-span-2 relative">
						<label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5 block">Kategori</label>
						<div v-if="showCategoryDropdown" @click="showCategoryDropdown = false" class="fixed inset-0 z-40"></div>
						<div class="relative z-50">
							<input v-model="form.category" @focus="showCategoryDropdown = true" placeholder="Pilih / Ketik Kategori Baru..." class="w-full px-4 py-3.5 rounded-xl border border-slate-200 focus:border-blue-600 outline-none font-bold text-sm bg-white uppercase transition-all text-slate-800" />
							<div @click="showCategoryDropdown = !showCategoryDropdown" class="absolute inset-y-0 right-0 pr-4 flex items-center cursor-pointer hover:scale-110 transition-transform">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 transition-transform duration-200 hover:text-blue-500" :class="showCategoryDropdown ? 'rotate-180' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
							</div>
						</div>
						<transition enter-active-class="transition duration-200 ease-out" enter-from-class="transform scale-95 opacity-0" enter-to-class="transform scale-100 opacity-100" leave-active-class="transition duration-100 ease-in" leave-from-class="transform scale-100 opacity-100" leave-to-class="transform scale-95 opacity-0">
							<div v-if="showCategoryDropdown" class="absolute z-50 w-full mt-2 bg-white border border-slate-200 rounded-2xl shadow-xl max-h-52 overflow-y-auto custom-scrollbar overflow-hidden flex flex-col">
								<div v-for="cat in filteredCategories" :key="cat" @click="selectCategory(cat)" class="px-4 py-3 hover:bg-blue-50 cursor-pointer text-xs md:text-sm font-bold text-slate-700 transition-colors uppercase border-b border-slate-50 last:border-0">{{ cat }}</div>
								<div v-if="form.category && !categories.some((c) => c.toLowerCase() === form.category.toLowerCase())" @click="showCategoryDropdown = false" class="px-4 py-3 bg-indigo-50 hover:bg-indigo-100 cursor-pointer text-xs md:text-sm font-black text-indigo-700 transition-colors uppercase flex items-center gap-2 sticky bottom-0 border-t border-indigo-100">
									<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" /></svg>
									<span class="truncate">Buat Baru: "{{ form.category }}"</span>
								</div>
								<div v-if="filteredCategories.length === 0 && !form.category" class="p-4 text-center text-xs font-bold text-slate-400 italic">Belum ada kategori tersedia.</div>
							</div>
						</transition>
					</div>
					<div class="md:col-span-2 p-5 bg-slate-900 rounded-[28px] text-white shadow-xl mt-2 mb-2 relative">
						<div class="flex items-center justify-between mb-4">
							<div class="flex items-center gap-2">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-blue-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
									<path d="m7.5 4.27 9 5.15" />
									<path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z" />
									<path d="M12 22V12" />
								</svg>
								<h4 class="font-black text-[10px] uppercase tracking-[0.2em]">Konversi & Satuan Jual</h4>
							</div>
							<span v-if="isEditing" class="text-[8px] font-black text-rose-400 uppercase tracking-widest bg-rose-900/30 px-2 py-1 rounded border border-rose-900/50 flex items-center gap-1 z-20">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
									<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
									<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
								</svg>
								Struktur Terkunci
							</span>
						</div>
						<div :class="isEditing ? 'opacity-50 pointer-events-none' : ''">
							<div class="grid grid-cols-2 gap-4 mb-4">
								<div>
									<label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Satuan Dasar Jual Terkecil</label>
									<select v-model="form.satuan_dasar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl outline-none font-black text-xs uppercase cursor-pointer text-white">
										<option value="PCS">PCS</option>
										<option value="POUCH">POUCH</option>
										<option value="KG">KG</option>
										<option value="GRAM">GRAM</option>
										<option value="LITER">LITER</option>
										<option value="ML">ML</option>
										<option value="PACK">PACK</option>
										<option value="BOX">BOX</option>
										<option value="BOTOL">BOTOL</option>
										<option value="BATANG">BATANG</option>
										<option value="BUNGKUS">BUNGKUS</option>
									</select>
								</div>
								<div>
									<label class="text-[8px] font-black text-slate-500 uppercase block mb-1">Beli Dalam Karton/Kemasan Besar?</label>
									<div @click="form.has_satuan_besar = !form.has_satuan_besar" class="w-full p-3.5 bg-slate-800 border border-slate-700 rounded-xl font-black text-[10px] uppercase cursor-pointer flex items-center justify-between transition-colors hover:border-blue-500">
										{{ form.has_satuan_besar ? 'YA (AKTIF)' : 'TIDAK (HANYA PCS)' }}
										<div :class="form.has_satuan_besar ? 'bg-blue-500' : 'bg-slate-600'" class="w-2 h-2 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.5)] transition-colors"></div>
									</div>
								</div>
							</div>
							<div v-if="form.has_satuan_besar" class="col-span-2 flex flex-col gap-4 pt-4 border-t border-slate-800 mt-2">
								<div>
									<label class="text-[8px] font-black text-blue-400 uppercase block mb-1">Sebutannya Apa ? ( Kemasan Paling Besar )</label>
									<input v-model="form.satuan_besar" type="text" placeholder="KARTON / SLOP / KARUNG / JERIGEN" class="w-full md:w-1/2 p-3 bg-slate-800 border border-blue-900 focus:border-blue-500 rounded-xl outline-none font-black text-xs uppercase text-white transition-all" />
									<button v-if="!form.is_nested_uom" @click="enableNestedUom" class="mt-3 text-[9px] font-black text-blue-400 hover:text-blue-300 uppercase tracking-widest text-left flex items-start sm:items-center gap-1.5 w-full sm:w-max bg-blue-900/30 px-3 py-2.5 rounded-xl border border-blue-900/50 transition-all hover:bg-blue-900/50 active:scale-95 shadow-sm whitespace-normal">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 flex-shrink-0 mt-0.5 sm:mt-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" /></svg>
										<span>Tambah Kemasan Tengah ( Cth: Slop -> Bungkus -> Batang )</span>
									</button>
								</div>
								<div v-if="form.is_nested_uom" class="animate-[fadeInUp_0.3s_ease-out]">
									<div class="flex items-center justify-between md:w-1/2 mb-1">
										<label class="text-[8px] font-black text-sky-400 uppercase block">Sebutannya Apa ? (Kemasan Tengah)</label>
										<button @click="form.is_nested_uom = false" class="text-[8px] font-black text-rose-400 hover:text-rose-300 uppercase tracking-widest flex items-center gap-1 bg-rose-900/20 px-2 py-1 rounded-lg">
											<svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
											Batal
										</button>
									</div>
									<input v-model="form.satuan_tengah" type="text" placeholder="Cth: BUNGKUS / RENTENG / LITER / KG" class="w-full md:w-1/2 p-3 bg-slate-800 border border-sky-900 focus:border-sky-500 rounded-xl outline-none font-black text-xs uppercase text-white transition-all" />
								</div>
								<div class="p-4 bg-slate-800/50 border border-slate-700 rounded-2xl flex flex-col gap-4 mt-2 shadow-inner">
									<div class="flex items-center gap-2 mb-1">
										<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" /></svg>
										<span class="text-[9px] font-black text-amber-400 uppercase tracking-widest">Kalkulator Isi Otomatis</span>
									</div>
									<div v-if="form.is_nested_uom" class="flex flex-col gap-4">
										<div class="grid grid-cols-1 md:grid-cols-2 gap-4 transition-all duration-300">
											<div>
												<label class="text-[8px] font-black text-slate-400 uppercase block mb-1">1 {{ form.satuan_besar || 'BESAR' }} isi berapa {{ form.satuan_tengah || 'TENGAH' }}?</label>
												<input :value="formatNumber(form.isi_besar_ke_tengah)" @input="handleInputForm('isi_besar_ke_tengah', $event)" type="text" inputmode="numeric" placeholder="Contoh: 10" class="w-full p-3 bg-slate-900 border border-slate-700 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-center text-white transition-all" />
											</div>
											<div>
												<label class="text-[8px] font-black text-slate-400 uppercase block mb-1">1 {{ form.satuan_tengah || 'TENGAH' }} isi berapa {{ form.satuan_dasar }}?</label>
												<input :value="formatNumber(form.isi_tengah_ke_dasar)" @input="handleInputForm('isi_tengah_ke_dasar', $event)" type="text" inputmode="numeric" placeholder="Contoh: 16" class="w-full p-3 bg-slate-900 border border-slate-700 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-center text-white transition-all" />
											</div>
										</div>
										<div class="text-[9px] font-black text-emerald-400 uppercase tracking-widest italic text-center mt-1 bg-emerald-900/20 py-2.5 rounded-lg border border-emerald-900/50">
											<span v-if="form.isi_per_besar > 0">Sistem Menyimpan: 1 {{ form.satuan_besar || 'BESAR' }} = {{ formatNumber(form.isi_per_besar) }} {{ form.satuan_dasar }}</span>
											<span v-else>Isi angka di atas untuk hitung otomatis...</span>
										</div>
									</div>
									<div v-else-if="form.satuan_dasar === 'GRAM' || form.satuan_dasar === 'ML'" class="flex flex-col gap-2">
										<label class="text-[8px] font-black text-slate-400 uppercase block">1 {{ form.satuan_besar || (form.satuan_dasar === 'GRAM' ? 'KARUNG' : 'JERIGEN') }} isinya berapa {{ form.satuan_dasar === 'GRAM' ? 'KG' : 'LITER' }}?</label>
										<div class="flex items-center gap-3">
											<input :value="formatNumber(form.input_kg)" @input="handleInputForm('input_kg', $event)" type="text" inputmode="numeric" :placeholder="form.satuan_dasar === 'GRAM' ? 'Contoh: 20' : 'Contoh: 18'" class="w-full md:w-1/3 p-3 bg-slate-900 border border-slate-700 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-white transition-all text-center" />
											<span class="text-xs font-black text-slate-400">{{ form.satuan_dasar === 'GRAM' ? 'KG' : 'LITER' }}</span>
										</div>
									</div>
									<div v-else class="flex flex-col gap-2">
										<label class="text-[8px] font-black text-slate-400 uppercase block">1 {{ form.satuan_besar || 'BESAR' }} isi berapa {{ form.satuan_dasar }}?</label>
										<div class="flex gap-2 items-center">
											<input :value="formatNumber(form.isi_per_besar)" @input="handleInputForm('isi_per_besar', $event)" type="text" inputmode="numeric" placeholder="Contoh: 24" class="w-full md:w-1/3 p-3 bg-slate-900 border border-slate-700 focus:border-blue-500 rounded-xl outline-none font-black text-xs text-center text-white transition-all" />
											<span class="text-xs font-black text-slate-400">{{ form.satuan_dasar }}</span>
										</div>
									</div>
								</div>
							</div>
						</div>
						<div v-if="form.has_satuan_besar" class="grid grid-cols-2 gap-4 pt-4 border-t border-slate-800 mt-4">
							<div :class="isEditing ? 'opacity-50 pointer-events-none' : ''">
								<label class="text-[8px] font-black text-amber-400 uppercase block mb-1">Harga Beli 1 {{ form.satuan_besar || 'KEMASAN' }}</label>
								<input :value="formatNumber(form.harga_beli_besar)" @input="handleInputForm('harga_beli_besar', $event)" type="text" inputmode="numeric" placeholder="Rp" class="w-full p-3 bg-amber-900/20 border border-amber-900 focus:border-amber-500 rounded-xl outline-none font-black text-xs text-amber-400 transition-all" />
							</div>
							<div>
								<label class="text-[8px] font-black text-emerald-400 uppercase block mb-1">Harga Jual 1 {{ form.satuan_besar || 'KEMASAN' }}</label>
								<input :value="formatNumber(form.harga_jual_besar)" @input="handleInputForm('harga_jual_besar', $event)" type="text" inputmode="numeric" placeholder="Rp" class="w-full p-3 bg-emerald-900/20 border border-emerald-900 focus:border-emerald-500 rounded-xl outline-none font-black text-xs text-emerald-400 transition-all" />
							</div>
						</div>
					</div>
					<div class="p-5 border rounded-[28px] transition-all duration-300 flex flex-col justify-between h-full" :class="form.has_satuan_besar ? 'bg-slate-100/80 border-transparent' : 'bg-white border-slate-200'">
						<div>
							<label class="text-[9px] font-black uppercase tracking-widest mb-3 block text-center" :class="form.has_satuan_besar ? 'text-indigo-500' : 'text-slate-400'">Harga Modal Dasar ({{ form.satuan_dasar }})</label>
						</div>
						<div class="relative w-full my-auto">
							<span class="absolute inset-y-0 left-0 pl-5 flex items-center text-sm font-black" :class="form.has_satuan_besar ? 'text-indigo-400' : 'text-slate-400'">Rp</span>
							<input :value="formatNumber(form.cost_price)" disabled type="text" :class="form.has_satuan_besar ? 'text-indigo-600 bg-slate-200/40 cursor-not-allowed border-transparent shadow-none' : 'text-slate-800 bg-white border-slate-200 focus:border-blue-600 outline-none shadow-inner'" class="w-full pl-12 pr-4 py-5 rounded-2xl text-3xl font-black text-center border-2 transition-all" />
						</div>
						<div class="mt-4 flex items-center justify-center">
							<span v-if="form.has_satuan_besar" class="text-[8px] font-black text-indigo-500 uppercase tracking-widest italic">* Terkunci dari kalkulator grosir</span>
							<span v-else class="text-[8px] font-black text-transparent select-none uppercase tracking-widest italic">* Filler</span>
						</div>
					</div>
					<div class="p-5 bg-blue-50 border border-blue-100 rounded-[28px] shadow-sm flex flex-col justify-between h-full">
						<label class="text-[9px] font-black text-blue-500 uppercase tracking-widest mb-3 block text-center">Patokan Harga Jual Eceran</label>
						<div class="flex flex-col gap-2 my-auto">
							<div class="relative">
								<span class="absolute inset-y-0 left-0 pl-5 flex items-center text-sm font-black text-blue-400">Rp</span>
								<input :value="formatNumber(form.harga_eceran_tampil)" @input="handleInputForm('harga_eceran_tampil', $event)" type="text" inputmode="numeric" placeholder="Harga..." class="w-full pl-12 pr-4 py-4 rounded-2xl bg-white border border-blue-200 focus:border-blue-600 outline-none font-black text-xl text-blue-700 shadow-inner transition-all text-center" />
							</div>
							<div class="text-center font-black text-[9px] text-blue-400 uppercase italic py-1">UNTUK SETIAP</div>
							<div class="relative flex items-center">
								<input :value="formatNumber(form.qty_eceran_tampil)" @input="handleInputForm('qty_eceran_tampil', $event)" type="text" inputmode="numeric" class="w-full pl-4 pr-16 py-4 rounded-2xl bg-white border border-blue-200 focus:border-blue-600 outline-none font-black text-xl text-center text-blue-700 shadow-inner transition-all" />
								<span class="absolute inset-y-0 right-0 pr-5 flex items-center font-black text-[10px] text-blue-600 uppercase">{{ form.satuan_dasar }}</span>
							</div>
						</div>
						<div class="mt-4 pt-4 border-t border-blue-200/50 flex flex-col gap-2 px-1">
							<p class="text-[9px] font-black text-slate-500 uppercase tracking-widest italic text-center">*Disimpan di sistem: Rp {{ formatNumber(form.price) }} / {{ form.satuan_dasar }}</p>
							<p class="text-[9px] font-black uppercase tracking-widest bg-emerald-100/50 px-2 py-2 rounded-lg text-center" :class="form.price - form.cost_price < 0 ? 'text-rose-600 bg-rose-100/50' : 'text-emerald-600 bg-emerald-100/50'">Profit: Rp {{ formatNumber((form.price - form.cost_price) * form.qty_eceran_tampil) }} / {{ formatNumber(form.qty_eceran_tampil) }} {{ form.satuan_dasar }}</p>
						</div>
					</div>
					<div class="md:col-span-2 grid grid-cols-1 gap-4 bg-slate-50 p-5 rounded-[24px] border border-slate-200 relative overflow-hidden">
						<div v-if="isEditing" class="absolute inset-0 bg-slate-100/80 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center p-4">
							<div class="bg-white px-6 py-4 rounded-2xl shadow-sm border border-rose-200 flex flex-col items-center text-center max-w-sm">
								<svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-rose-500 mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
									<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
									<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
								</svg>
								<span class="font-black text-[10px] uppercase text-rose-600 tracking-widest">Kolom Stok Dikunci!</span>
								<span class="text-[9px] font-bold text-slate-500 mt-1">Gunakan menu "Stok Masuk / Keluar" untuk mengubah kuantitas stok.</span>
							</div>
						</div>
						<div v-if="form.has_satuan_besar" class="grid grid-cols-1 sm:grid-cols-3 gap-4 items-center" :class="isEditing ? 'opacity-30 pointer-events-none' : ''">
							<div>
								<label class="text-[9px] font-black text-indigo-600 uppercase tracking-widest mb-2 block">Jumlah {{ form.satuan_besar }}</label>
								<input :value="formatNumber(stokDalamKarton)" @input="handleInputEmit('update:stokDalamKarton', $event)" type="text" inputmode="numeric" placeholder="0" class="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 focus:border-indigo-600 outline-none font-black text-lg text-indigo-600 transition-all shadow-sm" />
							</div>
							<div>
								<label class="text-[9px] font-black text-amber-600 uppercase tracking-widest mb-2 block">+ Lebih Eceran ({{ form.satuan_dasar }})</label>
								<input :value="formatNumber(eceranTambahan)" @input="handleInputEmit('update:eceranTambahan', $event)" type="text" inputmode="numeric" placeholder="0" class="w-full px-4 py-3.5 rounded-xl bg-white border border-slate-200 focus:border-amber-600 outline-none font-black text-lg text-amber-600 transition-all shadow-sm" />
							</div>
							<div>
								<label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block">Total Stok Akhir</label>
								<input :value="formatNumber(form.stock)" type="text" disabled class="w-full px-4 py-3.5 rounded-xl bg-slate-100 border border-transparent font-black text-xl text-slate-500 text-center" />
							</div>
						</div>
						<div v-else :class="isEditing ? 'opacity-30 pointer-events-none' : ''">
							<label class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-2 block text-center">Stok Awal Fisik (dalam {{ form.satuan_dasar }})</label>
							<input :value="formatNumber(form.stock)" @input="handleInputForm('stock', $event)" type="text" inputmode="numeric" placeholder="0" class="w-full px-4 py-4 rounded-2xl bg-white border-2 border-slate-200 focus:border-blue-600 outline-none font-black text-center text-2xl text-slate-800 transition-all" />
						</div>
					</div>
				</div>
			</div>
			<div class="p-6 bg-slate-50 border-t border-slate-100 shrink-0">
				<button @click="emit('submit')" :disabled="isSubmitting" class="w-full py-5 font-black text-xs uppercase tracking-[0.2em] bg-blue-600 text-white rounded-[24px] shadow-xl shadow-blue-200 hover:bg-slate-900 transition-all active:scale-95 flex items-center justify-center gap-3 disabled:opacity-50">
					<template v-if="isSubmitting">
						<div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
						Menyimpan Data...
					</template>
					<template v-else>
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
							<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
							<polyline points="17 21 17 13 7 13 7 21" />
							<polyline points="7 3 7 8 15 8" />
						</svg>
						Simpan Perubahan Produk
					</template>
				</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
@keyframes fadeInUp {
	from {
		opacity: 0;
		transform: translateY(5px);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}
.custom-scrollbar::-webkit-scrollbar {
	height: 6px;
	width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
	background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
	background: #cbd5e1;
	border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
	background: #94a3b8;
}
</style>
