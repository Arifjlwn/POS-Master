<script setup>
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';

const route = useRoute();
const router = useRouter();
const isSidebarOpen = ref(false);

// 🔥 DINAMIS: Nama toko bakal update otomatis kalau kamu ganti di menu Setting
const storeName = ref(localStorage.getItem('store_name') || 'ARZU');
const shortName = computed(() => storeName.value.split(' ')[0].substring(0, 10));

const menuItems = [
    { name: 'Kasir Laundry', path: '/laundry/pos', icon: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z M4 8h16 M8 12h4 M8 16h8' },
    { name: 'Status Cucian', path: '/laundry/status', icon: 'M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z M12 6v6l4 2' },
    { name: 'Master Layanan', path: '/laundry/master-layanan', icon: 'M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z M3.27 6.96L12 12.01l8.73-5.05 M12 22.08V12' },
    { name: 'Laporan', path: '/laundry/laporan', icon: 'M18 20V10 M12 20V4 M6 20v-6 M3 20h18' },
    // 🚀 MENU BARU: PENGATURAN TOKO
    { name: 'Setting Toko', path: '/laundry/setting-toko', icon: 'M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6z M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z' }
];

const handleLogout = () => {
    Swal.fire({
        title: 'Keluar Kasir?', 
        text: 'Sesi kamu akan diakhiri.',
        icon: 'question', 
        showCancelButton: true,
        confirmButtonColor: '#e11d48', 
        cancelButtonColor: '#94a3b8',
        confirmButtonText: 'Ya, Keluar'
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.clear();
            router.push('/login');
        }
    });
};
</script>

<template>
    <div class="flex h-screen bg-[#F8FAFC] overflow-hidden">
        
        <div v-if="isSidebarOpen" @click="isSidebarOpen = false" class="fixed inset-0 bg-slate-900/60 z-[60] lg:hidden backdrop-blur-sm transition-opacity"></div>

        <aside :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'" 
               class="fixed inset-y-0 left-0 z-[70] w-[260px] bg-[#1e1b4b] text-white flex flex-col transition-transform duration-300 ease-in-out lg:static lg:shrink-0 shadow-2xl lg:shadow-none border-r border-slate-800">
            
            <div class="p-7 flex items-center gap-4 shrink-0 border-b border-white/10 bg-black/10">
                <div class="w-12 h-12 bg-gradient-to-br from-indigo-500 to-indigo-700 rounded-2xl flex items-center justify-center shadow-lg shadow-indigo-500/30 border border-indigo-400/20">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22v-5"/><path d="M9 7V2"/><path d="M15 7V2"/><rect width="16" height="15" x="4" y="7" rx="2" ry="2"/></svg>
                </div>
                <div class="overflow-hidden">
                    <h2 class="font-black text-base tracking-[0.2em] uppercase leading-none truncate">{{ shortName }}</h2>
                    <p class="text-[9px] font-black text-indigo-400 tracking-[0.3em] mt-1.5 uppercase truncate">LAUNDRY POS</p>
                </div>
            </div>

            <nav class="flex-1 overflow-y-auto p-4 space-y-2">
                <router-link v-for="menu in menuItems" :key="menu.name" :to="menu.path" @click="isSidebarOpen = false"
                    :class="route.path.includes(menu.path) ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-900/50 ring-1 ring-white/10' : 'text-slate-400 hover:bg-white/10 hover:text-white'" 
                    class="flex items-center gap-4 px-5 py-4 rounded-[16px] transition-all group active:scale-95 relative overflow-hidden">
                    
                    <div v-if="route.path.includes(menu.path)" class="absolute left-0 top-0 bottom-0 w-1.5 bg-white rounded-r-full"></div>
                    
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 transition-transform group-hover:scale-110 relative z-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path :d="menu.icon"/></svg>
                    <span class="font-black text-[11px] uppercase tracking-[0.15em] relative z-10">{{ menu.name }}</span>
                </router-link>
            </nav>

            <div class="p-4 mt-auto border-t border-white/5">
                <button @click="handleLogout" class="w-full flex items-center gap-4 px-5 py-4 text-rose-400 bg-rose-500/5 hover:bg-rose-500/15 hover:text-rose-300 rounded-[16px] transition-all group active:scale-95 border border-transparent hover:border-rose-500/20">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                    <span class="font-black text-[11px] uppercase tracking-[0.15em]">Keluar Sistem</span>
                </button>
            </div>
        </aside>

        <main class="flex-1 flex flex-col min-w-0 h-full relative overflow-hidden">
            
            <header class="lg:hidden bg-white/80 backdrop-blur-md border-b border-slate-200/60 px-5 py-4 flex items-center justify-between shrink-0 z-40 sticky top-0">
                <div class="flex items-center gap-3">
                    <button @click="isSidebarOpen = true" class="p-2.5 bg-slate-50 hover:bg-slate-100 text-slate-700 rounded-xl active:scale-90 transition-all border border-slate-200/50">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
                    </button>
                    <span class="font-black text-[11px] text-slate-700 uppercase tracking-widest">{{ shortName }} LAUNDRY</span>
                </div>
                <div class="w-9 h-9 rounded-full bg-indigo-50 border border-indigo-100 flex items-center justify-center shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4.5 h-4.5 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/></svg>
                </div>
            </header>
            
            <div class="flex-1 h-full overflow-hidden bg-[#F8FAFC]">
                <slot />
            </div>
        </main>
    </div>
</template>