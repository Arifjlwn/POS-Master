<script setup>
defineProps({ 
    activeTab: String,
    isOwner: Boolean, // 🚀 Terima prop role
    isKlaimEligible: Boolean, // 🚀 TERIMA PROP INI
    daysLeftKlaim: Number     // 🚀 TERIMA PROP INI JUGA
});
defineEmits(['update:activeTab']);
</script>

<template>
    <div>
        <div class="bg-gradient-to-br from-slate-900 via-slate-800 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-6 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
            <div class="absolute -left-10 -bottom-10 opacity-10">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-64 h-64" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" /></svg>
            </div>
            
            <div class="z-10 text-center md:text-left">
                <h1 class="text-3xl md:text-4xl font-black tracking-tighter mb-2 uppercase italic leading-none">Stock <span class="text-indigo-400">Opname</span></h1>
                <p class="text-slate-300 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center justify-center md:justify-start gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
                    Penyesuaian Stok Fisik vs Sistem
                </p>
            </div>

            <div class="z-10 mt-6 md:mt-0 flex items-center gap-3 bg-amber-500/20 px-5 py-3 rounded-2xl border border-amber-500/30 backdrop-blur-md">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
                <div class="flex flex-col">
                    <span class="text-[9px] font-black text-amber-200 uppercase tracking-widest">Metode Audit</span>
                    <span class="text-sm font-black text-amber-400 uppercase tracking-widest">
                        {{ isOwner ? 'GUIDED COUNTING' : 'BLIND COUNTING' }}
                    </span>
                </div>
            </div>
        </div>

        <div class="flex justify-end mb-8">
            <div class="bg-white p-1.5 rounded-2xl shadow-sm border border-slate-100 flex gap-2">
                <button @click="$emit('update:activeTab', 'SO')" 
                        :class="activeTab === 'SO' ? 'bg-indigo-600 text-white shadow-md' : 'text-slate-500 hover:bg-slate-50'"
                        class="px-6 py-2.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all">
                    Audit Reguler
                </button>
                
                <button v-if="isKlaimEligible"
        @click="$emit('update:activeTab', 'KLAIM')" 
        :class="activeTab === 'KLAIM' 
            ? 'bg-amber-500 text-white shadow-lg shadow-amber-200 border-2 border-amber-500' 
            : 'text-slate-600 hover:text-amber-600 hover:bg-amber-50 border-2 border-transparent hover:border-amber-200'"
        class="px-5 py-2.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all flex items-center gap-2.5">
    
    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
    </svg>
    
    Klaim Barang
    
    <span :class="activeTab === 'KLAIM' ? 'bg-amber-600 text-white shadow-inner' : 'bg-amber-100 text-amber-600'" 
            class="px-2.5 py-1 rounded-md text-[9px] font-black tracking-widest ml-1">
        SISA {{ daysLeftKlaim }} HARI
    </span>
</button>

<button v-else
        disabled
        class="px-6 py-2.5 rounded-xl font-black text-xs uppercase tracking-widest transition-all flex items-center gap-2 bg-slate-100 text-slate-400 cursor-not-allowed border border-slate-200"
        title="Klaim sudah dilakukan atau masa klaim 7 hari sudah lewat">
    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
    {{ isKlaimEligible ? 'Sudah Diklaim' : 'Klaim Terkunci' }}
</button>
            </div>
        </div>
    </div>
</template>