<script setup>
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

// Nanti fungsi logout disesuaikan sama yang udah ada di auth kamu ya beb
const handleLogout = () => {
    localStorage.removeItem('token');
    router.push('/login');
};
</script>

<template>
    <div class="flex h-screen bg-slate-50 w-full overflow-hidden font-sans">
        
        <aside class="w-20 md:w-64 bg-slate-900 border-r border-slate-800 flex flex-col justify-between shrink-0 transition-all duration-300 z-50 shadow-2xl">
            <div>
                <div class="h-20 flex items-center justify-center md:justify-start md:px-8 border-b border-slate-800/50 bg-slate-900/50">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 bg-indigo-500 rounded-xl flex items-center justify-center shadow-lg shadow-indigo-500/30">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/><path d="M7 2v20"/><path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"/></svg>
                        </div>
                        <div class="hidden md:block">
                            <h2 class="font-black text-white tracking-widest uppercase text-sm leading-none">POS Resto</h2>
                            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Sistem Kasir F&B</span>
                        </div>
                    </div>
                </div>

                <nav class="p-3 md:p-4 space-y-2 mt-4">
                    <p class="hidden md:block text-[10px] font-black text-slate-500 uppercase tracking-widest px-4 mb-2">Menu Utama</p>

                    <router-link to="/fnb/kasir" class="group flex items-center gap-4 px-3 md:px-4 py-3 md:py-3.5 rounded-xl transition-all" :class="route.path.includes('/fnb/kasir') ? 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2"/><path d="M3 11h18"/><path d="M12 11v8"/></svg>
                        <span class="hidden md:block font-bold text-xs uppercase tracking-widest">Kasir POS</span>
                    </router-link>

                    <router-link to="/fnb/dapur" class="group flex items-center gap-4 px-3 md:px-4 py-3 md:py-3.5 rounded-xl transition-all" :class="route.path.includes('/fnb/dapur') ? 'bg-orange-500/10 text-orange-400 border border-orange-500/20' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 13.87A4 4 0 0 1 7.41 6a5.11 5.11 0 0 1 1.05-1.54 5 5 0 0 1 7.08 0A5.11 5.11 0 0 1 16.59 6 4 4 0 0 1 18 13.87V21H6Z"/><line x1="6" y1="17" x2="18" y2="17"/></svg>
                        <span class="hidden md:block font-bold text-xs uppercase tracking-widest">Monitor Dapur</span>
                    </router-link>

                    <router-link to="/fnb/laporan" class="group flex items-center gap-4 px-3 md:px-4 py-3 md:py-3.5 rounded-xl transition-all" :class="route.path.includes('/fnb/laporan') ? 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
                        <span class="hidden md:block font-bold text-xs uppercase tracking-widest">Laporan Resto</span>
                    </router-link>
                </nav>
            </div>

            <div class="p-3 md:p-4 border-t border-slate-800">
                <button @click="handleLogout" class="w-full flex items-center justify-center md:justify-start gap-4 px-3 md:px-4 py-3 text-rose-400 hover:bg-rose-500/10 hover:text-rose-300 rounded-xl transition-all font-bold text-xs uppercase tracking-widest">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                    <span class="hidden md:block">Keluar</span>
                </button>
            </div>
        </aside>

        <main class="flex-1 overflow-hidden bg-[#F8FAFC]">
            <slot />
        </main>
    </div>
</template>