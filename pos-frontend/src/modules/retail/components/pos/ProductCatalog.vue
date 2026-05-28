<script setup>

defineProps({
    searchQuery: String,
    filteredProducts: Array,
    heldOrders: Array,
    getImageUrl: Function
});

const emit = defineEmits(['update:searchQuery', 'barcode-scan', 'start-scanner', 'show-held', 'add-to-cart']);
</script>

<template>
    <div class="flex-1 flex flex-col min-h-0 w-full lg:w-8/12 xl:w-9/12 transition-all duration-300 no-print">
        <div class="flex gap-2 md:gap-3 shrink-0 mb-3 md:mb-4 items-stretch h-12 md:h-14">
            <div class="relative flex-1 group h-full">
                <div class="absolute inset-y-0 left-0 pl-4 md:pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input
                    type="text"
                    :value="searchQuery"
                    @input="emit('update:searchQuery', $event.target.value)" 
                    @keydown.enter.prevent="emit('barcode-scan')"
                    placeholder="Cari atau Scan..."
                    class="w-full h-full pl-12 md:pl-14 pr-4 rounded-[16px] md:rounded-[20px] border-2 border-slate-200 focus:border-indigo-600 shadow-sm text-slate-800 font-bold bg-white text-xs md:text-sm transition-all outline-none"
                >
            </div>
            
            <button v-if="heldOrders?.length > 0" @click="emit('show-held')" class="lg:hidden shrink-0 bg-amber-100 hover:bg-amber-500 text-amber-600 hover:text-white px-4 md:px-5 rounded-[16px] md:rounded-[20px] transition-all border-2 border-amber-200 flex items-center justify-center shadow-sm relative animate-pulse">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                <span class="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-rose-500 text-[10px] font-black text-white shadow-md">{{ heldOrders.length }}</span>
            </button>
            
            <button @click="emit('start-scanner')" class="shrink-0 bg-indigo-100 hover:bg-indigo-600 text-indigo-600 hover:text-white px-4 md:px-5 rounded-[16px] md:rounded-[20px] transition-all border-2 border-indigo-200 flex items-center justify-center shadow-sm h-full">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 md:w-7 md:h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </button>
        </div>

        <div v-if="filteredProducts.length === 0" class="flex-1 flex flex-col items-center justify-center bg-white/50 rounded-[24px] md:rounded-[32px] border-2 border-dashed border-slate-300">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 md:w-24 md:h-24 text-slate-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" /></svg>
            <p class="text-slate-400 font-black text-sm md:text-lg uppercase tracking-widest text-center">Produk Tidak Ditemukan</p>
        </div>

        <div v-else class="flex-1 overflow-y-auto custom-scrollbar pr-2 pb-4">
            <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 xl:grid-cols-5 gap-3 md:gap-4">
                <div v-for="product in filteredProducts" :key="product.id" @click="emit('add-to-cart', product)"
                    class="bg-white rounded-[16px] md:rounded-[24px] shadow-sm hover:shadow-xl hover:ring-2 hover:ring-indigo-500 transition-all duration-200 overflow-hidden cursor-pointer border border-slate-100 group flex flex-col transform hover:-translate-y-1">
                    
                    <div class="relative pt-2 px-2 md:pt-3 md:px-3">
                        <div class="bg-slate-50 rounded-xl md:rounded-2xl overflow-hidden aspect-square flex items-center justify-center border border-slate-100">
                            <img 
                                :src="getImageUrl(product.image) || 'https://placehold.co/150x150?text=No+Foto'" 
                                :alt="product.name"
                                class="w-full h-full object-contain mix-blend-multiply p-3 md:p-4 group-hover:scale-110 transition-transform duration-300"
                            >
                        </div>
                        <div v-if="product.stock <= 0" class="absolute top-3 right-3 md:top-5 md:right-5 text-[8px] md:text-[9px] font-black px-1.5 md:px-2 py-0.5 md:py-1 rounded-md shadow-sm bg-rose-500 text-white animate-pulse uppercase tracking-widest">
                            HABIS
                        </div>
                    </div>

                    <div class="p-2 md:p-4 flex flex-col flex-1 text-center justify-between gap-1 md:gap-2">
                        <h2 class="font-bold text-slate-700 text-[10px] md:text-[11px] line-clamp-2 leading-tight uppercase" :title="product.name">{{ product.name }}</h2>
                        <p class="text-indigo-700 font-black text-xs md:text-sm">Rp {{ product.price.toLocaleString('id-ID') }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>