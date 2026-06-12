<script setup>
defineProps({
    planCounts: Object
});
</script>

<template>
    <div class="grid grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between hover:border-indigo-500/40 transition-all">
            <span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Masa Trial</span>
            <h2 class="text-xl sm:text-2xl font-black text-indigo-400 mt-2">
                {{ planCounts.trial }} <span class="text-xs text-slate-500 font-bold">Ruko</span>
            </h2>
        </div>
        <div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between hover:border-amber-500/40 transition-all">
            <span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Paket Basic</span>
            <h2 class="text-xl sm:text-2xl font-black text-amber-400 mt-2">
                {{ planCounts.basic }} <span class="text-xs text-slate-500 font-bold">Ruko</span>
            </h2>
        </div>
        <div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between hover:border-emerald-500/40 transition-all">
            <span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Paket Pro</span>
            <h2 class="text-xl sm:text-2xl font-black text-emerald-400 mt-2">
                {{ planCounts.pro }} <span class="text-xs text-slate-500 font-bold">Ruko</span>
            </h2>
        </div>
        <div class="bg-[#131B2E] border border-slate-800 p-4 rounded-2xl shadow-xl flex flex-col justify-between hover:border-purple-500/40 transition-all">
            <span class="text-[10px] text-slate-500 font-black uppercase tracking-wider">Paket Premium</span>
            <h2 class="text-xl sm:text-2xl font-black text-purple-400 mt-2">
                {{ planCounts.premium }} <span class="text-xs text-slate-500 font-bold">Ruko</span>
            </h2>
        </div>
    </div>
</template>