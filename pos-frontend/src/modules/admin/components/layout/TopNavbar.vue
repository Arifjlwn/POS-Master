<script setup>
const emit = defineEmits(['toggle-sidebar']);
const founderName = localStorage.getItem('name') || 'Arif Juliawan';
</script>

<template>
    <header class="h-20 bg-[#131B2E] border-b border-slate-800 flex items-center justify-between px-4 lg:px-8 z-30 sticky top-0 backdrop-blur-md bg-opacity-90 transition-all">
        <div class="flex items-center gap-4">
            <button @click="emit('toggle-sidebar')" class="p-2 rounded-xl text-slate-400 hover:bg-slate-800 hover:text-white transition-colors">
                <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
                </svg>
            </button>
            <h2 class="text-white font-black tracking-tight text-lg hidden sm:block">Mission Control</h2>
        </div>

        <div class="flex items-center gap-4">
            <div class="hidden md:flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/20 rounded-full">
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                <span class="text-emerald-400 text-[10px] font-black uppercase tracking-widest">Sistem Stabil</span>
            </div>

            <div class="flex items-center gap-3 pl-4 border-l border-slate-700">
                <div class="text-right hidden sm:block">
                    <p class="text-white text-sm font-bold">{{ founderName }}</p>
                    <p class="text-indigo-400 text-[10px] font-black uppercase tracking-widest">Founder / Root</p>
                </div>
                <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-600 to-purple-600 p-[2px]">
                    <div class="w-full h-full bg-[#131B2E] rounded-[10px] flex items-center justify-center">
                        <span class="text-white font-black text-sm">AJ</span>
                    </div>
                </div>
            </div>
        </div>
    </header>
</template>