<script setup>
defineProps({
    mrr: Number,
    arr: Number,
    totalTenants: Number,
    newTenants: Number,
    formatRupiah: Function
});
</script>

<template>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
        <div class="bg-gradient-to-br from-[#131B2E] to-[#1a243d] border border-slate-800 p-6 rounded-3xl shadow-xl relative overflow-hidden group hover:border-indigo-500/40 transition-all">
            <div class="absolute top-0 right-0 p-4 opacity-10 text-indigo-400 group-hover:opacity-20 transition-all">
                <svg class="w-24 h-24" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            </div>
            <span class="text-[10px] text-slate-500 font-black uppercase tracking-widest block">MRR (Monthly Recurring Revenue)</span>
            <h2 class="text-2xl sm:text-3xl font-black text-white mt-3 font-mono tracking-tight">{{ formatRupiah(mrr) }}</h2>
            <p class="text-[10px] text-emerald-400 font-bold uppercase tracking-wider mt-2 flex items-center gap-1">
                <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3"><path stroke-linecap="round" stroke-linejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" /></svg>
                Bulan Berjalan Aktif
            </p>
        </div>

        <div class="bg-gradient-to-br from-[#131B2E] to-[#1a243d] border border-slate-800 p-6 rounded-3xl shadow-xl relative overflow-hidden group hover:border-purple-500/40 transition-all">
            <div class="absolute top-0 right-0 p-4 opacity-10 text-purple-400 group-hover:opacity-20 transition-all">
                <svg class="w-24 h-24" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75m0-.75A2.25 2.25 0 016 2.25h13.5A2.25 2.25 0 0121 4.5v12a2.25 2.25 0 01-2.25 2.25M3.75 4.5h.007m-.008 11.25h.007m0 0h-.007m0 0l2.25 13.5" /></svg>
            </div>
            <span class="text-[10px] text-slate-500 font-black uppercase tracking-widest block">ARR (Annual Recurring Revenue)</span>
            <h2 class="text-2xl sm:text-3xl font-black text-purple-400 mt-3 font-mono tracking-tight">{{ formatRupiah(arr) }}</h2>
            <p class="text-[10px] text-slate-500 font-bold uppercase tracking-wider mt-2">Proyeksi Pendapatan Tahunan Platform</p>
        </div>

        <div class="bg-[#131B2E] border border-slate-800 p-6 rounded-3xl shadow-xl flex flex-col justify-between hover:border-blue-500/40 transition-all md:col-span-2 lg:col-span-1">
            <div class="flex items-center justify-between">
                <div>
                    <span class="text-[10px] text-slate-500 font-black uppercase tracking-widest block">Total Akumulasi Tenant</span>
                    <h2 class="text-3xl font-black text-blue-400 mt-1 font-mono">{{ totalTenants }} <span class="text-xs text-slate-500 font-bold">Gerai</span></h2>
                </div>
                <div class="px-3 py-1 bg-emerald-500/10 border border-emerald-500/20 rounded-xl text-emerald-400 text-[10px] font-black uppercase tracking-wider">
                    +{{ newTenants }} Bulan Ini
                </div>
            </div>
            <div class="w-full bg-slate-800 h-1.5 rounded-full mt-4 overflow-hidden">
                <div class="bg-blue-500 h-full rounded-full transition-all" :style="{ width: totalTenants > 0 ? '100%' : '0%' }"></div>
            </div>
        </div>
    </div>
</template>