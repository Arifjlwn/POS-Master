<script setup>
// Props nerima data dari Induk, gak perlu import composable di sini
defineProps({ 
    searchQuery: String 
});

// Emit buat ngasih perintah ke Induk (Parent)
defineEmits(['update:searchQuery', 'search', 'scan']);
</script>

<template>
    <div class="flex items-center gap-3">
        <div class="relative flex-1">
            <div class="absolute inset-y-0 left-0 pl-5 flex items-center pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-slate-400 group-focus-within:text-blue-500 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
            </div>
            
            <input 
                :value="searchQuery" 
                @input="$emit('update:searchQuery', $event.target.value); $emit('search')"
                type="text" 
                placeholder="Scan Barcode / Cari Produk..." 
                class="w-full pl-14 pr-4 py-5 bg-slate-50 border-2 border-transparent rounded-[24px] outline-none font-bold text-sm focus:border-blue-500 focus:bg-white transition-all shadow-inner"
            >
        </div>
        
        <button @click="$emit('scan')" class="shrink-0 p-5 bg-blue-50 hover:bg-blue-600 text-blue-600 hover:text-white rounded-[24px] transition-all border border-blue-100 shadow-sm" title="Scan pakai Kamera">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
        </button>
    </div>
</template>