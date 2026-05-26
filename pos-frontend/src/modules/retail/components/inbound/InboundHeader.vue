<script setup>
defineProps({ currentUser: Object });
</script>

<template>
    <div class="bg-gradient-to-br from-slate-900 to-indigo-900 rounded-[40px] p-8 md:p-10 mb-10 text-white shadow-2xl flex flex-col md:flex-row items-center justify-between relative overflow-hidden border border-white/10">
        <div class="z-10">
            <h1 class="text-3xl font-black tracking-tighter mb-2 uppercase italic leading-none">Inbound <span class="text-blue-400">Logistics</span></h1>
            <p class="text-blue-200 font-bold text-[10px] uppercase tracking-[0.2em] flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
                Penerimaan Stok Barang
            </p>
        </div>
        <div class="z-10 mt-4 md:mt-0 px-4 py-2 bg-white/10 rounded-2xl border border-white/20 backdrop-blur-md">
            <span class="text-[9px] font-black uppercase tracking-widest text-blue-300">Mode Akses: {{ currentUser.role }}</span>
        </div>
    </div>
</template>