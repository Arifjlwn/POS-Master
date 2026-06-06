<script setup>
import { ref, onMounted } from 'vue';

const props = defineProps({ 
    searchQuery: String,
    disabled: Boolean,
    products: Array 
});

defineEmits(['update:searchQuery', 'search', 'scan', 'add']);

const inputRef = ref(null);

onMounted(() => {
    // Sekarang variabel 'props' sudah legal dan bisa dibaca dengan aman !
    if (inputRef.value && !props.disabled) {
        inputRef.value.focus();
    }
});
</script>

<template>
    <div class="bg-white p-4 md:p-6 rounded-[32px] shadow-sm border border-slate-100 relative select-none">
        <div class="flex items-center gap-3">
            <div class="relative flex-1 group">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input 
                    ref="inputRef"
                    :value="searchQuery" 
                    @input="$emit('update:searchQuery', $event.target.value); $emit('search')" 
                    :disabled="disabled" 
                    :tabindex="disabled ? -1 : 0"
                    type="text" 
                    placeholder="Scan Barcode SKU / Ketik Nama Barang Ritel..." 
                    :class="disabled ? 'bg-slate-100 text-slate-400 cursor-not-allowed border-transparent shadow-inner' : 'bg-slate-50 border-2 border-transparent focus:border-indigo-500 focus:bg-white text-slate-800'"
                    class="w-full pl-14 pr-10 py-4 md:py-5 rounded-[24px] outline-none font-black text-xs md:text-sm transition-all placeholder:text-slate-400 shadow-sm"
                >
                
                <button v-if="searchQuery && searchQuery.length > 0" @click="$emit('update:searchQuery', ''); products.value = []" class="absolute inset-y-0 right-4 flex items-center text-slate-400 hover:text-rose-500 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 md:w-5 md:h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>
            
            <button @click="$emit('scan')" :disabled="disabled" class="shrink-0 p-4 md:p-5 bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white rounded-[24px] transition-all border border-indigo-100 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-indigo-50 disabled:hover:text-indigo-600 active:scale-95" title="Scan Barcode">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 md:w-6 md:h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </button>
        </div>
        
        <div v-if="searchQuery && searchQuery.length >= 2" class="absolute left-4 right-4 md:left-6 md:right-6 mt-3 bg-white border border-slate-200 rounded-[24px] shadow-2xl z-[100] overflow-hidden max-h-60 overflow-y-auto custom-scrollbar animate-fade-in">
            
            <div v-if="products.length">
                <div v-for="p in products" :key="p.id" @click="$emit('add', p)" class="p-4 md:p-5 hover:bg-indigo-50 cursor-pointer border-b border-slate-100 last:border-0 flex justify-between items-center group transition-all">
                    <div class="flex items-center gap-3 md:gap-4 min-w-0">
                        <div class="w-9 h-9 md:w-10 md:h-10 bg-slate-50 border border-slate-100 rounded-xl flex items-center justify-center text-lg md:text-xl shadow-inner group-hover:bg-indigo-100 group-hover:border-indigo-200 transition-all shrink-0">📦</div>
                        <div class="flex flex-col min-w-0 text-left">
                            <span class="font-black text-slate-800 uppercase text-xs md:text-sm group-hover:text-indigo-600 transition-colors truncate" :title="p.nama_produk">{{ p.nama_produk }}</span>
                            <span class="text-[9px] md:text-[10px] text-slate-400 font-bold tracking-wider mt-0.5 uppercase">SKU: {{ p.sku || p.SKU || '-' }}</span>
                        </div>
                    </div>
                    <button class="text-[9px] md:text-[10px] bg-slate-900 group-hover:bg-indigo-600 text-white px-3.5 py-1.5 md:px-4 md:py-2 rounded-xl font-black uppercase tracking-widest shadow-sm transition-colors shrink-0 active:scale-95">Pilih</button>
                </div>
            </div>

            <div v-else class="p-6 text-center flex flex-col items-center justify-center gap-2">
                <div class="text-2xl">❌</div>
                <div class="font-black text-slate-400 uppercase tracking-widest text-xs md:text-sm">Produk Tidak Terdaftar Dalam Master Data</div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.15s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }

.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
</style>