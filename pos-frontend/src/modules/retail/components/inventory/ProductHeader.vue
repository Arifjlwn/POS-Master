<script setup>
import { ref } from 'vue';

defineProps({ searchQuery: String, selectedCategory: String, categories: Array });
const emit = defineEmits(['update:searchQuery', 'update:selectedCategory', 'export', 'trigger-import', 'handle-import', 'add-new']);

const importInput = ref(null);
const isCategoryOpen = ref(false);

const onImportClick = () => {
	emit('trigger-import');
	importInput.value.click();
};
const selectCategory = (cat) => {
	emit('update:selectedCategory', cat);
	isCategoryOpen.value = false;
};
</script>

<template>
	<div class="mb-6 no-print">
		<div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[32px] p-6 md:p-10 mb-8 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden">
			<div class="z-10 text-center md:text-left">
				<h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-1 uppercase italic">
					Master
					<span class="text-blue-400">Inventory</span>
				</h1>
				<p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em]">Katalog Produk & Harga Gudang</p>
			</div>
			<div class="z-10 mt-6 md:mt-0 flex flex-wrap justify-center gap-3">
				<button @click="emit('export')" class="bg-emerald-500/20 hover:bg-emerald-500 text-emerald-400 hover:text-white px-5 py-3 rounded-[16px] font-black transition-all flex items-center gap-2 text-[10px] uppercase tracking-widest border border-emerald-500/50">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
					Ekspor
				</button>
				<button @click="onImportClick" class="bg-amber-500/20 hover:bg-amber-500 text-amber-400 hover:text-white px-5 py-3 rounded-[16px] font-black transition-all flex items-center gap-2 text-[10px] uppercase tracking-widest border border-amber-500/50">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
					Impor
				</button>
				<input type="file" ref="importInput" class="hidden" accept=".csv" @change="(e) => emit('handle-import', e)" />
				<button @click="emit('add-new')" class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-3 rounded-[16px] font-black transition-all shadow-lg flex items-center gap-2 text-[10px] uppercase tracking-widest active:scale-95 border border-blue-400">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" /></svg>
					Tambah
				</button>
			</div>
		</div>

		<div class="flex flex-col sm:flex-row gap-4">
			<div class="relative flex-1 group">
				<div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
				</div>
				<input :value="searchQuery" @input="emit('update:searchQuery', $event.target.value)" type="text" placeholder="Cari nama barang atau barcode..." class="block w-full pl-14 pr-4 py-4 bg-white rounded-2xl border-2 border-slate-100 shadow-sm focus:border-blue-600 outline-none font-bold text-sm transition-all text-slate-700" />
			</div>

			<div class="w-full sm:w-64 shrink-0 relative">
				<div v-if="isCategoryOpen" @click="isCategoryOpen = false" class="fixed inset-0 z-40"></div>
				<div @click="isCategoryOpen = !isCategoryOpen" class="relative block w-full pl-4 pr-10 py-4 bg-white rounded-2xl border-2 shadow-sm text-sm font-bold cursor-pointer outline-none transition-all uppercase flex items-center justify-between" :class="isCategoryOpen ? 'z-50 border-blue-600 text-blue-700' : 'z-10 border-slate-100 text-slate-700 hover:border-blue-300'">
					<span class="truncate">{{ selectedCategory || 'SEMUA KATEGORI' }}</span>
					<div class="absolute inset-y-0 right-0 pr-4 flex items-center pointer-events-none transition-transform duration-300" :class="isCategoryOpen ? 'rotate-180 text-blue-600' : 'text-slate-400'">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
					</div>
				</div>

				<transition enter-active-class="transition duration-200 ease-out" enter-from-class="transform scale-95 opacity-0 -translate-y-2" enter-to-class="transform scale-100 opacity-100 translate-y-0" leave-active-class="transition duration-150 ease-in" leave-from-class="transform scale-100 opacity-100 translate-y-0" leave-to-class="transform scale-95 opacity-0 -translate-y-2">
					<div v-if="isCategoryOpen" class="absolute z-50 w-full mt-2 bg-white border border-slate-100 rounded-2xl shadow-xl max-h-60 overflow-y-auto custom-scrollbar flex flex-col py-2">
						<div @click="selectCategory('')" class="px-5 py-3.5 hover:bg-blue-50 cursor-pointer text-xs font-black transition-colors uppercase border-b border-slate-50 last:border-0" :class="selectedCategory === '' ? 'text-blue-600 bg-blue-50/50' : 'text-slate-500 hover:text-blue-600'">SEMUA KATEGORI</div>
						<div v-for="cat in categories" :key="cat" @click="selectCategory(cat)" class="px-5 py-3.5 hover:bg-blue-50 cursor-pointer text-xs font-black transition-colors uppercase border-b border-slate-50 last:border-0" :class="selectedCategory === cat ? 'text-blue-600 bg-blue-50/50' : 'text-slate-500 hover:text-blue-600'">{{ cat }}</div>
					</div>
				</transition>
			</div>
		</div>
	</div>
</template>
