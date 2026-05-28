<script setup>

defineProps({
    products: Array,
    isLoading: Boolean,
    currentPage: Number,
    totalPages: Number,
    visiblePages: Array,
    getImageUrl: Function
});

const emit = defineEmits(['edit', 'delete', 'change-page']);
</script>

<template>
    <div class="w-full bg-white rounded-[32px] shadow-sm border border-slate-100 overflow-hidden flex flex-col no-print">
        <div class="overflow-x-auto custom-scrollbar">
            <table class="w-full text-left border-collapse whitespace-nowrap min-w-[900px]">
                <thead>
                    <tr class="bg-slate-50 border-b border-slate-100">
                        <th class="p-5 font-black uppercase tracking-widest text-[10px] text-slate-400">Produk & Detail</th>
                        <th class="p-5 font-black uppercase tracking-widest text-[10px] text-slate-400 text-center">Kategori</th>
                        <th class="p-5 font-black uppercase tracking-widest text-[10px] text-right text-slate-400">Modal Dasar</th>
                        <th class="p-5 font-black uppercase tracking-widest text-[10px] text-right text-slate-400">Harga Jual</th>
                        <th class="p-5 font-black uppercase tracking-widest text-[10px] text-center text-slate-400">Total Stok</th>
                        <th class="p-5 font-black uppercase tracking-widest text-[10px] text-center text-slate-400">Aksi</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-slate-50">
                    <tr v-if="isLoading">
                        <td colspan="6" class="p-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">Memuat Katalog Produk...</td>
                    </tr>
                    <tr v-else-if="products.length === 0">
                        <td colspan="6" class="p-16 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">Produk Tidak Ditemukan</td>
                    </tr>
                    <tr v-else v-for="product in products" :key="product.id" class="hover:bg-blue-50/30 transition-colors group">
                        
                        <!-- 1. PRODUK & DETAIL -->
                        <td class="p-5 flex items-center gap-4">
                            <div class="w-12 h-12 md:w-14 md:h-14 rounded-[16px] border-2 border-white shadow-sm bg-slate-50 flex items-center justify-center text-slate-300 overflow-hidden shrink-0 group-hover:border-blue-200 transition-colors">
                                <img v-if="product.gambar" :src="getImageUrl(product.gambar)" class="w-full h-full object-cover">
                                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                            </div>
                            <div>
                                <div class="font-black text-slate-800 text-xs md:text-sm uppercase">{{ product.nama_produk }}</div>
                                <div class="flex flex-wrap items-center gap-2 mt-1.5">
                                    <span class="text-[8px] md:text-[9px] font-black bg-slate-100 text-slate-500 px-2 py-0.5 rounded border border-slate-200 uppercase tracking-widest">{{ product.sku || 'NO-SKU' }}</span>
                                    <span v-if="product.satuan_besar" class="text-[8px] md:text-[9px] font-black bg-purple-50 text-purple-600 px-2 py-0.5 rounded border border-purple-100 uppercase tracking-widest" title="Isi Konversi">
                                        1 {{ product.satuan_besar }} = {{ product.isi_per_besar }} {{ product.satuan_dasar }}
                                    </span>
                                </div>
                            </div>
                        </td>

                        <!-- 2. KATEGORI -->
                        <td class="p-5 text-center">
                            <span class="bg-blue-50 text-blue-600 border border-blue-100 px-3 py-1.5 rounded-lg font-black text-[9px] uppercase tracking-widest">{{ product.kategori || 'General' }}</span>
                        </td>

                        <!-- 3. HARGA MODAL -->
                        <td class="p-5 text-right">
                            <div class="text-slate-500 font-black text-xs">Rp {{ (product.harga_modal || 0).toLocaleString('id-ID') }}</div>
                            <div class="text-[8px] font-bold text-slate-400 uppercase mt-0.5 tracking-widest">/ {{ product.satuan_dasar || 'PCS' }}</div>
                        </td>

                        <!-- 4. HARGA JUAL (ECERAN & GROSIR) -->
                        <td class="p-5 text-right">
                            <!-- Eceran -->
                            <div class="font-black text-slate-800 text-sm">Rp {{ (product.harga_jual || 0).toLocaleString('id-ID') }}</div>
                            <div class="text-[8px] font-bold text-slate-500 uppercase tracking-widest">/ {{ product.satuan_dasar || 'PCS' }}</div>
                            
                            <!-- Grosir (Muncul kalau ada) -->
                            <div v-if="product.satuan_besar" class="mt-2 pt-2 border-t border-slate-200 border-dashed">
                                <div class="font-black text-emerald-600 text-xs">Rp {{ (product.harga_jual_besar || 0).toLocaleString('id-ID') }}</div>
                                <div class="text-[8px] font-bold text-emerald-500 uppercase tracking-widest">/ {{ product.satuan_besar }}</div>
                            </div>
                        </td>

                        <!-- 5. STOK CERDAS (TRANSLATE KE BAHASA MANUSIA) -->
                        <td class="p-5 text-center">
                            <span class="px-3 py-1.5 text-[10px] rounded-lg font-black tracking-widest shadow-sm inline-flex items-center gap-1" :class="(product.stok || 0) > 10 ? 'bg-blue-50 text-blue-600 border border-blue-100' : 'bg-red-50 text-red-600 border border-red-100'">
                                {{ (product.stok || 0).toLocaleString('id-ID') }} <span class="uppercase text-[8px] opacity-70">{{ product.satuan_dasar || 'PCS' }}</span>
                            </span>
                            
                            <!-- RUMUS TRANSLATE KONVERSI -->
                            <div v-if="product.satuan_besar && product.isi_per_besar > 0" class="mt-2 flex flex-col items-center gap-0.5">
                                <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest italic">Setara Dengan:</span>
                                <div class="bg-slate-50 border border-slate-200 px-2 py-1 rounded text-[9px] font-black tracking-widest flex items-center gap-1">
                                    <span class="text-indigo-600">📦 {{ Math.floor((product.stok || 0) / product.isi_per_besar).toLocaleString('id-ID') }} {{ product.satuan_besar }}</span>
                                    <span v-if="(product.stok || 0) % product.isi_per_besar > 0" class="text-amber-500">+ {{ ((product.stok || 0) % product.isi_per_besar).toLocaleString('id-ID') }} {{ product.satuan_dasar }}</span>
                                </div>
                            </div>
                        </td>

                        <!-- 6. AKSI -->
                        <td class="p-5 text-center">
                            <div class="flex justify-center gap-2">
                                <button @click="emit('edit', product)" class="p-2.5 bg-slate-50 border border-slate-200 text-slate-400 rounded-xl hover:bg-blue-600 hover:text-white hover:border-blue-600 transition-colors" title="Edit">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2-2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                                </button>
                                <button @click="emit('delete', product.id)" class="p-2.5 bg-slate-50 border border-slate-200 text-slate-400 rounded-xl hover:bg-red-500 hover:text-white hover:border-red-500 transition-colors" title="Hapus">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h14" /></svg>
                                </button>
                            </div>
                        </td>

                    </tr>
                </tbody>
            </table>
        </div>

        <div class="p-4 bg-slate-50 border-t border-slate-100 flex flex-col md:flex-row justify-between items-center gap-4 shrink-0">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest text-center md:text-left">
                Halaman <span class="text-blue-600">{{ currentPage }}</span> dari {{ totalPages }} 
                <span v-if="products.length > 0" class="hidden sm:inline">| Menampilkan {{ products.length }} Item</span>
            </span>
            
            <div class="flex flex-wrap justify-center gap-1.5 md:gap-2">
                <button @click="emit('change-page', 1)" :disabled="currentPage === 1 || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" /></svg>
                </button>
                <button @click="emit('change-page', currentPage - 1)" :disabled="currentPage === 1 || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" /></svg>
                </button>
                
                <button v-for="page in visiblePages" :key="page" @click="emit('change-page', page)" :disabled="isLoading" 
                    :class="currentPage === page ? 'bg-blue-600 text-white border-blue-600 shadow-md shadow-blue-200' : 'bg-white text-slate-600 border-slate-200 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50'"
                    class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center rounded-lg text-[10px] md:text-xs font-black transition-all shadow-sm">
                    {{ page }}
                </button>
                
                <button @click="emit('change-page', currentPage + 1)" :disabled="currentPage === totalPages || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>
                </button>
                <button @click="emit('change-page', totalPages)" :disabled="currentPage === totalPages || isLoading" class="w-8 h-8 md:w-9 md:h-9 flex items-center justify-center bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-blue-600 hover:border-blue-300 hover:bg-blue-50 disabled:opacity-40 shadow-sm transition-all">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M13 5l7 7-7 7M5 5l7 7-7 7" /></svg>
                </button>
            </div>
        </div>
    </div>
</template>