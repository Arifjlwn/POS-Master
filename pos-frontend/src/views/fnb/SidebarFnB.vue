<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';

const route = useRoute();
const router = useRouter();

const isMobileMenuOpen = ref(false);

// State Data Profil
const userName = ref('Kasir');
const userRole = ref('Staff FnB');
const userInitial = ref('KR');

// 🚀 STATE BARU BUAT NAMA TOKO
const storeName = ref('Resto POS'); 
const storeInitial = ref('RP'); // Inisial toko kalau kepanjangan

onMounted(() => {
    // Tarik Profil User
    const name = localStorage.getItem('name') || 'Kasir Resto';
    userName.value = name;
    userRole.value = localStorage.getItem('role') || 'F&B Outlet';
    userInitial.value = name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase();

    // 🚀 TARIK NAMA TOKO DARI LOCAL STORAGE (Hasil dari Setup / Login)
    const savedStore = localStorage.getItem('storeName');
    if (savedStore) {
        storeName.value = savedStore;
        // Bikin inisial dari nama toko (misal: Ayam Bakar Badas -> AB)
        storeInitial.value = savedStore.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase();
    }
});

const handleLogout = () => {
    Swal.fire({
        title: 'Akhiri Sesi?',
        text: 'Pastikan antrean dapur aman ya bosku!',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#ef4444',
        confirmButtonText: 'Ya, Keluar'
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.removeItem('token');
            router.push('/login');
        }
    });
};
</script>

<template>
    <div class="flex flex-col lg:flex-row h-screen bg-[#F8FAFC] w-full overflow-hidden font-sans relative">

        <header class="lg:hidden h-16 bg-white border-b border-slate-200 flex items-center justify-between px-4 shrink-0 z-40 shadow-sm">
            <button @click="isMobileMenuOpen = true" class="w-10 h-10 flex items-center justify-center text-indigo-600 hover:bg-indigo-50 rounded-xl transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="4" x2="20" y1="12" y2="12"/><line x1="4" x2="20" y1="6" y2="6"/><line x1="4" x2="20" y1="18" y2="18"/></svg>
            </button>
            
            <div class="text-center overflow-hidden px-4">
                <h1 class="font-black text-indigo-900 text-sm uppercase tracking-widest truncate">{{ storeName }}</h1>
            </div>

            <div class="w-8 h-8 rounded-lg bg-slate-900 text-white flex items-center justify-center font-black text-[10px] uppercase shadow-sm shrink-0">
                {{ userInitial }}
            </div>
        </header>

        <div v-if="isMobileMenuOpen" @click="isMobileMenuOpen = false" class="lg:hidden fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-[55] transition-opacity"></div>

        <aside :class="isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'" 
               class="fixed lg:relative top-0 left-0 h-full w-[260px] bg-[#0B1121] border-r border-white/5 flex flex-col justify-between shrink-0 transition-transform duration-300 ease-out z-[60] shadow-2xl lg:shadow-none">
            
            <div>
                <div class="h-24 hidden lg:flex items-center justify-start px-7 border-b border-white/5 bg-white/[0.02] relative overflow-hidden">
                    <div class="absolute -top-10 -left-10 w-32 h-32 bg-indigo-500/20 rounded-full blur-3xl pointer-events-none"></div>
                    <div class="flex items-center gap-3 relative z-10 w-full">
                        
                        <div class="w-12 h-12 bg-gradient-to-br from-indigo-500 to-blue-600 rounded-[14px] flex items-center justify-center shadow-lg shadow-indigo-500/40 shrink-0 font-black text-white text-lg border border-white/10">
                            {{ storeInitial }}
                        </div>

                        <div class="overflow-hidden flex-1 pr-2">
                            <h2 class="font-black text-white tracking-widest uppercase text-sm leading-tight truncate" :title="storeName">{{ storeName }}</h2>
                            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1 block">POS FnB App</span>
                        </div>
                    </div>
                </div>

                <div class="lg:hidden h-16 flex items-center justify-between px-5 border-b border-white/5">
                    <span class="font-black text-white tracking-widest uppercase text-sm">Menu Utama</span>
                    <button @click="isMobileMenuOpen = false" class="w-8 h-8 flex items-center justify-center text-slate-400 hover:text-white bg-white/5 rounded-lg">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>

                <nav class="p-4 space-y-2 mt-2">
                    <p class="hidden lg:block text-[9px] font-black text-slate-500 uppercase tracking-[0.2em] px-3 mb-3">Menu Utama</p>

                    <router-link to="/fnb/kasir" @click="isMobileMenuOpen = false" class="relative flex items-center gap-4 px-4 py-3.5 rounded-xl transition-all duration-300 group overflow-hidden" :class="route.path.includes('/fnb/kasir') ? 'bg-indigo-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0 z-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2"/><path d="M3 11h18"/><path d="M12 11v8"/></svg>
                        <span class="font-black text-xs uppercase tracking-widest z-10">Kasir POS</span>
                    </router-link>

                    <router-link to="/fnb/dapur" @click="isMobileMenuOpen = false" class="relative flex items-center gap-4 px-4 py-3.5 rounded-xl transition-all duration-300 group overflow-hidden" :class="route.path.includes('/fnb/dapur') ? 'bg-orange-500 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0 z-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 13.87A4 4 0 0 1 7.41 6a5.11 5.11 0 0 1 1.05-1.54 5 5 0 0 1 7.08 0A5.11 5.11 0 0 1 16.59 6 4 4 0 0 1 18 13.87V21H6Z"/><line x1="6" y1="17" x2="18" y2="17"/></svg>
                        <span class="font-black text-xs uppercase tracking-widest z-10">Monitor Dapur</span>
                    </router-link>

                    <router-link to="/fnb/master-menu" @click="isMobileMenuOpen = false" class="relative flex items-center gap-4 px-4 py-3.5 rounded-xl transition-all duration-300 group overflow-hidden" :class="route.path.includes('/fnb/master-menu') ? 'bg-blue-600 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0 z-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M12 2v20M17 5H9.5a4.5 4.5 0 0 0 0 9H15a4.5 4.5 0 0 1 0 9H6.5"/></svg>
                        <span class="font-black text-xs uppercase tracking-widest z-10">Master Menu</span>
                    </router-link>

                    <router-link to="/fnb/laporan" @click="isMobileMenuOpen = false" class="relative flex items-center gap-4 px-4 py-3.5 rounded-xl transition-all duration-300 group overflow-hidden" :class="route.path.includes('/fnb/laporan') ? 'bg-emerald-500 text-white' : 'text-slate-400 hover:bg-slate-800 hover:text-white'">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0 z-10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
                        <span class="font-black text-xs uppercase tracking-widest z-10">Laporan Resto</span>
                    </router-link>
                </nav>
            </div>

            <div class="p-4 border-t border-white/5 space-y-3 bg-white/[0.01]">
                <div class="flex items-center gap-3 p-3 rounded-xl bg-slate-800/50 border border-white/5">
                    <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shrink-0 font-black text-white text-xs">
                        {{ userInitial }}
                    </div>
                    <div class="overflow-hidden">
                        <h4 class="text-xs font-black text-white uppercase tracking-wider truncate">{{ userName }}</h4>
                        <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5 truncate">{{ userRole }}</p>
                    </div>
                </div>
                <button @click="handleLogout" class="w-full flex items-center gap-4 px-4 py-3.5 text-rose-400 hover:bg-rose-500 hover:text-white rounded-xl transition-all group border border-transparent hover:border-rose-500/50">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                    <span class="font-black text-xs uppercase tracking-widest">Keluar Sesi</span>
                </button>
            </div>
        </aside>

        <main class="flex-1 overflow-hidden relative">
            <slot />
        </main>

    </div>
</template>