<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Swal from 'sweetalert2';

const router = useRouter();
const route = useRoute();
const sidebarOpen = ref(false);

// State untuk buka-tutup grup menu
const openGroups = ref({
    stock: route.path.includes('barang') || route.path.includes('opname') || route.path.includes('produk') || route.path.includes('returns'),
    admin: route.path.startsWith('/dashboard') || route.path.startsWith('/karyawan') || route.path.startsWith('/setup')
});

const toggleGroup = (group) => {
    openGroups.value[group] = !openGroups.value[group];
};

const user = ref({
    name: localStorage.getItem('name') || 'User',
    role: localStorage.getItem('role') || 'staff',
    storeName: localStorage.getItem('storeName') || 'POS UMKM'
});

onMounted(() => {
    user.value.name = localStorage.getItem('name') || 'User';
    user.value.role = localStorage.getItem('role') || 'staff';
});

const logout = () => {
    Swal.fire({
        title: 'Mau keluar, Bos?',
        text: "Pastikan semua kerjaan hari ini sudah tersimpan ya!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#4f46e5',
        cancelButtonColor: '#94a3b8',
        confirmButtonText: 'Ya, Logout',
        cancelButtonText: 'Batal',
        customClass: {
            popup: 'rounded-[32px]',
            confirmButton: 'rounded-[16px] font-black px-6 py-3',
            cancelButton: 'rounded-[16px] font-black px-6 py-3'
        }
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.clear();
            router.push('/login');
        }
    });
};
</script>

<template>
    <div class="min-h-screen bg-[#F8FAFC] flex flex-col relative overflow-hidden font-sans selection:bg-indigo-100 selection:text-indigo-600">
        
        <header class="bg-white/80 backdrop-blur-md border-b border-slate-100 flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4 sticky top-0 z-40 shadow-sm">
            <div class="flex items-center gap-4 sm:gap-5">
                <button @click="sidebarOpen = true" class="group flex flex-col gap-1.5 p-2 rounded-xl hover:bg-indigo-50 transition-all active:scale-95">
                    <div class="w-6 h-0.5 bg-indigo-600 rounded-full group-hover:w-4 transition-all"></div>
                    <div class="w-6 h-0.5 bg-indigo-600 rounded-full"></div>
                    <div class="w-4 h-0.5 bg-indigo-600 rounded-full group-hover:w-6 transition-all"></div>
                </button>
                <div class="flex flex-col">
                    <div class="font-black text-lg sm:text-xl text-slate-900 tracking-tighter leading-none">
                        POS<span class="text-indigo-600">UMKM</span>
                    </div>
                    <span class="text-[8px] sm:text-[9px] font-black text-slate-400 uppercase tracking-widest truncate max-w-[120px] sm:max-w-xs">{{ user.storeName }}</span>
                </div>
            </div>
            
            <div class="flex items-center gap-4">
                <div class="flex items-center gap-3 sm:pl-4 sm:border-l border-slate-100">
                    <div class="text-right hidden sm:block">
                        <div class="text-xs font-black text-slate-800 uppercase leading-none">{{ user.name.split(' ')[0] }}</div>
                        <span class="text-[9px] font-bold text-indigo-600 uppercase tracking-tighter bg-indigo-50 px-1.5 py-0.5 rounded-md mt-1 inline-block border border-indigo-100">{{ user.role }}</span>
                    </div>
                    <div class="w-9 h-9 sm:w-10 sm:h-10 rounded-[14px] bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center text-white text-xs sm:text-sm font-black shadow-md border-2 border-white ring-2 ring-slate-100">
                        {{ user.name.substring(0, 2).toUpperCase() }}
                    </div>
                </div>
            </div>
        </header>

        <Transition name="fade">
            <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-slate-950/40 backdrop-blur-sm z-40 transition-all"></div>
        </Transition>

        <aside :class="sidebarOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full'" class="fixed inset-y-0 left-0 w-[280px] sm:w-[320px] bg-white transform transition-transform duration-500 ease-[cubic-bezier(0.34,1.56,0.64,1)] z-50 flex flex-col border-r border-slate-100">
            
            <div class="p-6 sm:p-8 flex items-center justify-between bg-slate-50/50 border-b border-slate-100">
                <div class="flex flex-col">
                    <div class="font-black text-xl sm:text-2xl text-slate-900 tracking-tighter leading-none">
                        POS<span class="text-indigo-600">UMKM</span>
                    </div>
                    <span class="text-[9px] sm:text-[10px] font-bold text-slate-400 uppercase tracking-[0.3em] mt-1">Management System</span>
                </div>
                <button @click="sidebarOpen = false" class="w-9 h-9 sm:w-10 sm:h-10 flex items-center justify-center rounded-xl bg-white border border-slate-200 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all shadow-sm active:scale-95">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <nav class="flex-1 px-4 sm:px-6 py-6 space-y-8 overflow-y-auto custom-scrollbar">
                
                <div>
                    <div class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] px-2 mb-4 flex items-center gap-2">
                        <span class="w-1.5 h-1.5 bg-indigo-500 rounded-full shadow-[0_0_8px_rgba(99,102,241,0.6)]"></span> Operasional
                    </div>
                    <div class="space-y-1.5">
                        <router-link to="/pos/kasir" @click="sidebarOpen = false" class="nav-link group" :class="{ 'active': route.path === '/pos/kasir' }">
                            <svg class="icon group-hover:-rotate-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                            <span>Mesin Kasir (POS)</span>
                        </router-link>

                        <router-link to="/riwayat" @click="sidebarOpen = false" class="nav-link group" :class="{ 'active': route.path.startsWith('/riwayat') }">
                            <svg class="icon group-hover:scale-110" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" /></svg>
                            <span>Riwayat Transaksi</span>
                        </router-link>

                        <router-link to="/absensi" @click="sidebarOpen = false" class="nav-link group" :class="{ 'active': route.path === '/absensi' }">
                            <svg class="icon group-hover:translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                            <span>Absensi</span>
                        </router-link>

                        <router-link to="/schedule" @click="sidebarOpen = false" class="nav-link group" :class="{ 'active': route.path === '/schedule' }">
                            <svg class="icon group-hover:rotate-45" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                            <span>Jadwal Shift (TSM)</span>
                        </router-link>
                    </div>
                </div>

                <div>
                    <button @click="toggleGroup('stock')" class="group-btn hover:text-emerald-600 group">
                        <span class="flex items-center gap-2"><span class="w-1.5 h-1.5 bg-emerald-500 rounded-full shadow-[0_0_8px_rgba(16,185,129,0.6)]"></span> Inventori & Stok</span>
                        <svg :class="openGroups.stock ? 'rotate-180 text-emerald-500' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    
                    <div v-show="openGroups.stock" class="mt-2 space-y-1 ml-4 border-l-2 border-slate-100 pl-2">
                        
                        <router-link v-if="user.role === 'owner'" to="/produk" @click="sidebarOpen = false" class="sub-link hover:text-emerald-600 hover:bg-emerald-50 hover:border-emerald-500" :class="{ 'active-sub !text-emerald-700 !bg-emerald-50/80 !border-emerald-500': route.path.startsWith('/produk') }">
                            Master Data Produk
                        </router-link>

                        <router-link to="/penerimaan-barang" @click="sidebarOpen = false" class="sub-link hover:text-emerald-600 hover:bg-emerald-50 hover:border-emerald-500" :class="{ 'active-sub !text-emerald-700 !bg-emerald-50/80 !border-emerald-500': route.path === '/penerimaan-barang' }">
                            Terima Barang (LPB)
                        </router-link>
                        <router-link to="/stock-opname" @click="sidebarOpen = false" class="sub-link hover:text-emerald-600 hover:bg-emerald-50 hover:border-emerald-500" :class="{ 'active-sub !text-emerald-700 !bg-emerald-50/80 !border-emerald-500': route.path === '/stock-opname' }">
                            Stock Opname
                        </router-link>
                        <router-link to="/retur-barang" @click="sidebarOpen = false" class="sub-link hover:text-rose-600 hover:bg-rose-50 hover:border-rose-500" :class="{ 'active-sub !text-rose-700 !bg-rose-50/80 !border-rose-500': route.path.startsWith('/returns') }">
                            Waste & Retur
                        </router-link>
                    </div>
                </div>

                <div v-if="user.role !== 'staff'">
                    <button @click="toggleGroup('admin')" class="group-btn hover:text-amber-600 group">
                        <span class="flex items-center gap-2"><span class="w-1.5 h-1.5 bg-amber-500 rounded-full shadow-[0_0_8px_rgba(245,158,11,0.6)]"></span> Administrasi Toko</span>
                        <svg :class="openGroups.admin ? 'rotate-180 text-amber-500' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    
                    <div v-show="openGroups.admin" class="mt-2 space-y-1 ml-4 border-l-2 border-slate-100 pl-2">
                        
                        <router-link to="/stock-opname/report" @click="sidebarOpen = false" class="sub-link hover:text-amber-600 hover:bg-amber-50 hover:border-amber-500" :class="{ 'active-sub !text-amber-700 !bg-amber-50/80 !border-amber-500': route.path === '/stock-opname/report' }">
                            Laporan Hasil SO
                        </router-link>

                        <router-link to="/retur-barang/report" @click="sidebarOpen = false" class="sub-link hover:text-amber-600 hover:bg-amber-50 hover:border-amber-500" :class="{ 'active-sub !text-amber-700 !bg-amber-50/80 !border-amber-500': route.path === '/retur-barang/report' }">
                            Laporan Retur Barang
                        </router-link>

                        <template v-if="user.role === 'owner'">
                            <router-link to="/dashboard" @click="sidebarOpen = false" class="sub-link hover:text-amber-600 hover:bg-amber-50 hover:border-amber-500" :class="{ 'active-sub !text-amber-700 !bg-amber-50/80 !border-amber-500': route.path === '/dashboard' }">
                                Dashboard Analitik
                            </router-link>
                            <router-link to="/karyawan" @click="sidebarOpen = false" class="sub-link hover:text-amber-600 hover:bg-amber-50 hover:border-amber-500" :class="{ 'active-sub !text-amber-700 !bg-amber-50/80 !border-amber-500': route.path.startsWith('/karyawan') }">
                                Manajemen Karyawan
                            </router-link>
                            <router-link to="/setup" @click="sidebarOpen = false" class="sub-link hover:text-amber-600 hover:bg-amber-50 hover:border-amber-500" :class="{ 'active-sub !text-amber-700 !bg-amber-50/80 !border-amber-500': route.path.startsWith('/setup') }">
                                Pengaturan Sistem
                            </router-link>
                        </template>
                    </div>
                </div>

                <div v-if="user.role === 'staff'" class="px-4 py-4 bg-slate-50 border-2 border-dashed border-slate-200 rounded-2xl flex items-center justify-center gap-2 mt-4 text-slate-400">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                    <span class="text-[9px] font-black uppercase tracking-[0.2em]">Akses Manajemen Dibatasi</span>
                </div>

            </nav>

            <div class="p-4 sm:p-6 bg-slate-50/80 border-t border-slate-100">
                <div class="flex items-center justify-between p-3 sm:p-4 bg-white rounded-2xl shadow-sm border border-slate-100 hover:border-slate-200 transition-colors">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-[14px] bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center text-white font-black text-sm shadow-md">
                            {{ user.name.substring(0, 2).toUpperCase() }}
                        </div>
                        <div class="flex flex-col max-w-[100px] sm:max-w-[120px]">
                            <div class="text-xs font-black text-slate-800 uppercase leading-none truncate">{{ user.name }}</div>
                            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1 truncate">{{ user.role }}</span>
                        </div>
                    </div>
                    <button @click="logout" class="w-9 h-9 sm:w-10 sm:h-10 flex items-center justify-center rounded-xl bg-slate-50 text-slate-400 hover:text-white hover:bg-rose-500 transition-all active:scale-95 group" title="Logout">
                        <svg class="w-4 h-4 sm:w-5 sm:h-5 group-hover:translate-x-0.5 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3 3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                    </button>
                </div>
            </div>
        </aside>

        <main class="flex-1 w-full max-w-full overflow-y-auto bg-[#F8FAFC] h-[calc(100vh-64px)] relative scroll-smooth custom-scrollbar">
            <slot /> 
        </main>

    </div>
</template>

<style scoped>
/* Navigation Link Base */
.nav-link {
    @apply flex items-center gap-3.5 px-4 py-3.5 rounded-[16px] text-xs font-black tracking-tight transition-all duration-300 border-2 border-transparent;
    @apply text-slate-500 hover:bg-white hover:border-slate-100 hover:shadow-sm;
}

/* Icon Base */
.icon {
    @apply w-[18px] h-[18px] transition-transform duration-300;
}

/* Active State Operasional */
.nav-link.active {
    @apply bg-slate-900 text-white shadow-xl shadow-slate-200 border-slate-900 translate-x-1;
}
.nav-link.active .icon {
    @apply text-indigo-400;
}

/* Group Button Base */
.group-btn {
    @apply w-full flex items-center justify-between px-2 py-2 text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] transition-colors;
}

/* Sub-links Base */
.sub-link {
    @apply flex items-center gap-4 px-4 py-2.5 rounded-xl text-[11px] font-bold text-slate-500 transition-all duration-200 border-l-[3px] border-transparent;
}

/* Base style buat active sublink, dioverride di template via inline class :class */
.active-sub {
    @apply shadow-sm;
}

/* Custom Transitions Overlay */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* Custom Scrollbar */
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px;}
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>