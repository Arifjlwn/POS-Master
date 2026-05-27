<script setup>
defineProps({ 
    searchQuery: String, 
    products: Array 
});
defineEmits(['update:searchQuery', 'search', 'scan', 'add']);
</script>

<template>
    <div class="bg-white p-6 rounded-[32px] shadow-sm border border-slate-100 relative">
        <div class="flex items-center gap-3">
            <div class="relative flex-1 group">
                <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                </div>
                <input 
                    :value="searchQuery" 
                    @input="$emit('update:searchQuery', $event.target.value); $emit('search')" 
                    type="text" 
                    placeholder="Scan Barcode / Ketik Nama Barang..." 
                    class="w-full pl-14 pr-4 py-5 bg-slate-50 border-2 border-transparent rounded-[24px] outline-none font-bold text-sm focus:border-indigo-500 focus:bg-white transition-all text-slate-800 placeholder:text-slate-400"
                >
            </div>
            
            <button @click="$emit('scan')" class="shrink-0 p-5 bg-indigo-50 hover:bg-indigo-600 text-indigo-600 hover:text-white rounded-[24px] transition-all border border-indigo-100 shadow-sm" title="Scan Barcode">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
            </button>
        </div>
        
        <div v-if="products.length && searchQuery.length > 0" class="absolute left-6 right-6 mt-3 bg-white border border-slate-100 rounded-[28px] shadow-2xl z-[100] overflow-hidden">
            <div v-for="p in products" :key="p.id" @click="$emit('add', p)" class="p-5 hover:bg-indigo-50 cursor-pointer border-b last:border-0 flex justify-between items-center group transition-all">
                <div class="flex items-center gap-4">
                    <div class="w-10 h-10 bg-slate-100 rounded-xl flex items-center justify-center text-xl shadow-inner group-hover:bg-indigo-100 transition-colors">📦</div>
                    <div class="font-black text-slate-800 uppercase text-sm group-hover:text-indigo-600 transition-colors">{{ p.nama_produk }}</div>
                </div>
                <span class="text-[9px] bg-slate-900 text-white px-4 py-2 rounded-xl font-black uppercase tracking-widest shadow-sm">Pilih</span>
            </div>
        </div>
    </div>
</template>