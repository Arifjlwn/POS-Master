<script setup>
defineProps({ 
    searchProductQuery: String, 
    isDropdownOpen: Boolean, 
    filteredProducts: Array, 
    selectedProduct: Object, 
    form: Object, 
    alasanOptions: Array 
});

const emit = defineEmits(['update:searchProductQuery', 'update:form', 'start-scanner', 'select-product', 'clear-product', 'add-to-cart']);

// 🚀 SECURITY: Fungsi untuk memotong angka minus dan huruf 'e' dari inputan kasir
const sanitizeNumber = (value) => {
    let sanitized = String(value).replace(/[^0-9]/g, '');
    return sanitized === '' ? 0 : Number(sanitized);
};
</script>

<template>
    <div class="bg-white rounded-[32px] p-6 md:p-8 shadow-sm border border-slate-100 xl:sticky xl:top-24">
        <h2 class="text-lg font-black text-slate-800 uppercase tracking-tight mb-6 flex items-center gap-2 border-b border-slate-100 pb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-rose-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M11 2a2 2 0 0 0-2 2v5H4a2 2 0 0 0-2 2v2c0 1.1.9 2 2 2h5v5c0 1.1.9 2 2 2h2a2 2 0 0 0 2-2v-5h5a2 2 0 0 0 2-2v-2a2 2 0 0 0-2-2h-5V4a2 2 0 0 0-2-2h-2z"/></svg>
            Pilih Barang
        </h2>

        <form @submit.prevent="$emit('add-to-cart')" class="flex flex-col gap-5">
            <div class="relative">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Pilih Barang Bermasalah</label>
                
                <div v-if="selectedProduct" class="flex items-center justify-between p-4 bg-rose-50 border-2 border-rose-200 rounded-2xl">
                    <div class="flex flex-col">
                        <span class="text-rose-900 font-black text-sm uppercase leading-tight truncate">{{ selectedProduct.nama_produk }}</span>
                        <span class="text-[10px] font-bold text-rose-500 mt-0.5">Stok Tersedia: {{ selectedProduct.stok }} {{ selectedProduct.satuan_dasar }}</span>
                    </div>
                    <button type="button" @click="$emit('clear-product')" class="w-8 h-8 rounded-xl bg-white text-rose-400 hover:text-rose-600 hover:bg-rose-100 flex items-center justify-center transition-all shadow-sm">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>

                <div v-else class="relative">
                    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M4 5v14M8 5v14M12 5v14M16 5v14M20 5v14"/><line x1="2" y1="12" x2="22" y2="12" stroke="red" stroke-width="1.5"/></svg>
                    </div>
                    <input 
                        :value="searchProductQuery" 
                        @input="$emit('update:searchProductQuery', $event.target.value)" 
                        type="text" 
                        class="w-full pl-12 pr-14 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-bold text-slate-800 text-sm transition-all placeholder:text-slate-300" 
                        placeholder="Ketik nama atau SKU..."
                    >
                    <button type="button" @click="$emit('start-scanner')" class="absolute inset-y-0 right-2 my-auto w-10 h-10 bg-slate-200 hover:bg-rose-100 text-slate-500 hover:text-rose-600 rounded-xl flex items-center justify-center transition-all">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M23 19a2 2 0 0 1-2-2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>
                    </button>
                    
                    <div v-if="isDropdownOpen && filteredProducts.length > 0" class="absolute z-50 w-full mt-2 bg-white border border-slate-100 rounded-2xl shadow-xl overflow-hidden max-h-60 overflow-y-auto">
                        <div 
                            v-for="prod in filteredProducts" 
                            :key="prod.id" 
                            @click="$emit('select-product', prod)"
                            class="p-4 border-b border-slate-50 hover:bg-rose-50 cursor-pointer transition-colors flex justify-between items-center group"
                        >
                            <div>
                                <div class="font-black text-slate-800 text-xs uppercase group-hover:text-rose-700">{{ prod.nama_produk }}</div>
                                <div class="text-[9px] font-bold text-slate-400 mt-1">{{ prod.sku }}</div>
                            </div>
                            <span class="text-[10px] font-black bg-slate-100 text-slate-500 px-2 py-1 rounded-lg group-hover:bg-rose-100 group-hover:text-rose-600">Sisa: {{ prod.stok }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="selectedProduct" class="p-4 bg-slate-50 rounded-2xl border-2 border-slate-100">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-3 block text-center">Jumlah Barang Yang Diretur</label>
                
                <div class="flex flex-wrap items-center justify-center gap-2">
                    <div v-if="selectedProduct.satuan_besar" class="flex flex-col items-center gap-1 w-20">
                        <input v-model.number="form.qty_besar" @input="form.qty_besar = sanitizeNumber($event.target.value)" type="text" inputmode="numeric" class="w-full p-2 bg-white border-2 border-slate-200 rounded-xl text-center font-black text-rose-600 outline-none focus:border-rose-500 shadow-inner">
                        <span class="text-[8px] font-black text-slate-400 uppercase">{{ selectedProduct.satuan_besar }}</span>
                    </div>

                    <span v-if="selectedProduct.is_nested_uom" class="font-bold text-slate-300 mb-4">+</span>
                    <div v-if="selectedProduct.is_nested_uom" class="flex flex-col items-center gap-1 w-20">
                        <input v-model.number="form.qty_tengah" @input="form.qty_tengah = sanitizeNumber($event.target.value)" type="text" inputmode="numeric" class="w-full p-2 bg-white border-2 border-slate-200 rounded-xl text-center font-black text-rose-600 outline-none focus:border-rose-500 shadow-inner">
                        <span class="text-[8px] font-black text-slate-400 uppercase">{{ selectedProduct.satuan_tengah }}</span>
                    </div>

                    <span v-if="selectedProduct.satuan_besar" class="font-bold text-slate-300 mb-4">+</span>
                    <div class="flex flex-col items-center gap-1 w-20">
                        <input v-model.number="form.qty_dasar" @input="form.qty_dasar = sanitizeNumber($event.target.value)" type="text" inputmode="numeric" class="w-full p-2 bg-white border-2 border-slate-200 rounded-xl text-center font-black text-rose-600 outline-none focus:border-rose-500 shadow-inner">
                        <span class="text-[8px] font-black text-slate-400 uppercase">{{ selectedProduct.satuan_dasar }}</span>
                    </div>
                </div>
            </div>

            <div v-else>
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Qty Retur</label>
                <input type="text" disabled class="w-full px-4 py-3.5 bg-slate-100 border-2 border-slate-100 rounded-2xl text-slate-400 font-black cursor-not-allowed " placeholder="Pilih barang dulu...">
            </div>

            <div>
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Alasan</label>
                <div class="relative">
                    <select v-model="form.alasan" required class="w-full px-4 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-bold text-slate-800 text-xs appearance-none cursor-pointer transition-all disabled:opacity-50 disabled:cursor-not-allowed" :disabled="!selectedProduct">
                        <option value="" disabled selected hidden>Pilih...</option>
                        <option v-for="opt in alasanOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
                    </select>
                    <div class="absolute inset-y-0 right-3 flex items-center pointer-events-none">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                    </div>
                </div>
            </div>

            <div>
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 mb-2 block">Keterangan Opsional</label>
                <textarea v-model="form.catatan" rows="1" class="w-full px-4 py-3.5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-rose-500 focus:bg-white outline-none font-medium text-slate-700 text-sm transition-all resize-none placeholder:text-slate-300 disabled:opacity-50 disabled:cursor-not-allowed" placeholder="Contoh: Jatuh pecah saat diturunin dari truk..." :disabled="!selectedProduct"></textarea>
            </div>

            <button type="submit" :disabled="!selectedProduct" class="w-full bg-slate-100 hover:bg-rose-50 border-2 border-slate-200 hover:border-rose-200 text-slate-600 hover:text-rose-600 py-4 rounded-[20px] font-black text-xs uppercase tracking-[0.2em] transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed mt-2 flex items-center justify-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
                TAMBAH KE KERANJANG
            </button>
        </form>
    </div>
</template>