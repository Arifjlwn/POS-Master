<script setup>
import { onMounted, ref } from 'vue';
import { getCategoryIcon } from '../../../../utils/categoryHelper.js';

defineProps({
	searchQuery: String,
	filteredProducts: Array,
	heldOrders: Array,
	getImageUrl: Function,
	categories: Array,
	selectedCategory: String,

	isLoading: {
		type: Boolean,
		default: false,
	},
});

const emit = defineEmits(['update:searchQuery', 'update:selectedCategory', 'barcode-scan', 'start-scanner', 'show-held', 'add-to-cart']);

const searchInputRef = ref(null);

onMounted(() => {
	if (searchInputRef.value) {
		searchInputRef.value.focus();
	}
});

const forceFocusScanner = () => {
	if (searchInputRef.value) {
		searchInputRef.value.focus();
	}
};
</script>

<template>
	<div @click="forceFocusScanner" class="flex-1 flex flex-col md:flex-row min-h-0 w-full lg:w-8/12 xl:w-9/12 gap-4 transition-all duration-300 no-print">
		<div class="flex md:flex-col overflow-x-auto md:overflow-x-visible md:overflow-y-auto gap-2 pb-2 md:pb-0 shrink-0 custom-scrollbar md:w-44 xl:w-52 scroll-smooth" @click.stop>
			<button type="button" @click="emit('update:selectedCategory', '')" :class="!selectedCategory ? 'bg-indigo-600 text-white shadow-md border-indigo-600' : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'" class="px-4 py-2.5 md:py-3.5 rounded-xl font-black text-[10px] uppercase tracking-widest border transition-all shrink-0 flex items-center gap-2 active:scale-95 shadow-sm md:w-full md:justify-start">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
					<path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
				</svg>
				<span>Semua Produk</span>
			</button>

			<button v-for="cat in categories" :key="cat" type="button" @click="emit('update:selectedCategory', cat)" :class="selectedCategory === cat ? 'bg-indigo-600 text-white shadow-md border-indigo-600' : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50'" class="px-4 py-2.5 md:py-3.5 rounded-xl font-black text-[10px] uppercase tracking-widest border transition-all shrink-0 flex items-center gap-2 active:scale-95 shadow-sm md:w-full md:justify-start truncate" :title="cat">
				<span v-html="getCategoryIcon(cat)" class="flex items-center justify-center shrink-0"></span>
				<span class="truncate">{{ cat }}</span>
			</button>
		</div>

		<div class="flex-1 flex flex-col min-h-0">
			<div class="flex gap-2 md:gap-3 shrink-0 mb-3 md:mb-4 items-stretch h-12 md:h-14">
				<div class="relative flex-1 group h-full">
					<div class="absolute inset-y-0 left-0 pl-4 md:pl-5 flex items-center pointer-events-none">
						<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
					</div>
					<input ref="searchInputRef" type="text" :value="searchQuery" @input="emit('update:searchQuery', $event.target.value)" @keydown.enter.prevent="emit('barcode-scan')" placeholder="Cari Nama Produk / Scan Barcode SKU..." class="w-full h-full pl-12 md:pl-14 pr-4 rounded-[16px] md:rounded-[20px] border-2 border-slate-200 focus:border-indigo-600 shadow-sm text-slate-800 font-bold bg-white text-xs md:text-sm transition-all outline-none" />
				</div>

				<button v-if="heldOrders?.length > 0" @click.stop="emit('show-held')" class="lg:hidden shrink-0 bg-amber-100 hover:bg-amber-500 text-amber-600 hover:text-white px-4 md:px-5 rounded-[16px] md:rounded-[20px] transition-all border-2 border-amber-200 flex items-center justify-center shadow-sm relative animate-pulse">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
					<span class="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-rose-500 text-[10px] font-black text-white shadow-md">{{ heldOrders.length }}</span>
				</button>

				<button @click.stop="emit('start-scanner')" class="shrink-0 bg-indigo-100 hover:bg-indigo-600 text-indigo-600 hover:text-white px-4 md:px-5 rounded-[16px] md:rounded-[20px] transition-all border-2 border-indigo-200 flex items-center justify-center shadow-sm h-full">
					<svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-7 md:h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
						<path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
					</svg>
				</button>
			</div>

			<div v-if="isLoading" class="flex-1 overflow-y-auto custom-scrollbar pr-2 pb-4">
				<div class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 md:gap-4">
					<div v-for="n in 8" :key="'shimmer_' + n" class="bg-white rounded-[16px] md:rounded-[24px] border border-slate-100 p-3 flex flex-col gap-3">
						<div class="bg-slate-100 rounded-xl md:rounded-2xl aspect-square shimmer-wave"></div>
						<div class="h-3.5 bg-slate-100 rounded-md w-3/4 mx-auto shimmer-wave"></div>
						<div class="h-4 bg-slate-100 rounded-md w-1/2 mx-auto shimmer-wave"></div>
					</div>
				</div>
			</div>

			<div v-else-if="filteredProducts.length === 0" class="flex-1 flex flex-col items-center justify-center bg-white/50 rounded-[24px] md:rounded-[32px] border-2 border-dashed border-slate-300">
				<svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 md:w-24 md:h-24 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" /></svg>
				<p class="text-slate-400 font-black text-sm md:text-lg uppercase tracking-widest text-center">Produk Tidak Terdaftar</p>
			</div>

			<div v-else class="flex-1 overflow-y-auto custom-scrollbar pr-2 pb-4">
				<div class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 md:gap-4">
					<div
						v-for="product in filteredProducts"
						:key="product.id"
						@click.stop="!product.stock || Number(product.stock) <= 0 || product.stock === '0' ? null : emit('add-to-cart', product)"
						:class="['bg-white rounded-[16px] md:rounded-[24px] shadow-sm transition-all duration-200 overflow-hidden border border-slate-100 group flex flex-col transform', !product.stock || Number(product.stock) <= 0 || product.stock === '0' ? 'opacity-50 cursor-not-allowed pointer-events-none select-none' : 'hover:shadow-xl hover:ring-2 hover:ring-indigo-500 hover:-translate-y-1 cursor-pointer']">
						<div class="relative pt-2 px-2 md:pt-3 md:px-3">
							<div class="bg-slate-100 rounded-xl md:rounded-2xl overflow-hidden aspect-square flex items-center justify-center border border-slate-100">
								<img :src="getImageUrl(product.image) || 'https://placehold.co/150x150?text=No+Foto'" :alt="product.name" class="w-full h-full object-contain mix-blend-multiply p-3 md:p-4 group-hover:scale-110 transition-transform duration-300" />
							</div>
							<div v-if="!product.stock || Number(product.stock) <= 0 || product.stock === '0'" class="absolute top-3 right-3 md:top-5 md:right-5 text-[8px] md:text-[9px] font-black px-1.5 md:px-2 py-0.5 md:py-1 rounded-md shadow-sm bg-red-700 text-white animate-pulse uppercase tracking-widest z-10">HABIS</div>
						</div>

						<div class="p-2 md:p-4 flex flex-col flex-1 text-center justify-between gap-1 md:gap-2">
							<h2 class="font-bold text-slate-700 text-[10px] md:text-[11px] line-clamp-2 leading-tight uppercase" :title="product.name">{{ product.name }}</h2>
							<div class="flex flex-col gap-0.5">
								<p class="text-indigo-700 font-black text-xs md:text-sm">Rp {{ product.price.toLocaleString('id-ID') }}</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
/* 🚀 THE ULTIMATE SHIMMER INFRASTRUCTURE ANCHOR */
.shimmer-wave {
	position: relative;
	overflow: hidden;
	background-color: #f1f5f9 !important; /* Slate 100 */
}

.shimmer-wave::after {
	position: absolute;
	top: 0;
	right: 0;
	bottom: 0;
	left: 0;
	transform: translateX(-100%);
	background-image: linear-gradient(90deg, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0.6) 20%, rgba(255, 255, 255, 0.9) 60%, rgba(255, 255, 255, 0) 100/100);
	animation: shimmerSwipe 1.6s infinite;
	content: '';
}

@keyframes shimmerSwipe {
	100% {
		transform: translateX(100%);
	}
}

.custom-scrollbar::-webkit-scrollbar {
	height: 4px;
	width: 4px;
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
