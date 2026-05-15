<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Swal from 'sweetalert2';

const router = useRouter();
const route = useRoute();
const sidebarOpen = ref(false);

// State untuk buka-tutup grup menu
const openGroups = ref({
    stock: route.path.includes('barang') || route.path.includes('opname'),
    admin: route.path.startsWith('/dashboard') || route.path === '/karyawan' || route.path === '/produk' || route.path === '/setup'
});

const toggleGroup = (group) => {
    openGroups.value[group] = !openGroups.value[group];
};

const user = ref({
    name: localStorage.getItem('name') || 'User',
    role: localStorage.getItem('role') || 'staff',
    storeName: localStorage.getItem('storeName') || 'Indo UMKM'
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
        confirmButtonColor: '#2563eb',
        cancelButtonColor: '#64748b',
        confirmButtonText: 'Ya, Logout',
        cancelButtonText: 'Batal',
        customClass: {
            popup: 'rounded-[32px]',
            confirmButton: 'rounded-xl font-black px-6',
            cancelButton: 'rounded-xl font-black px-6'
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
    <div class="min-h-screen bg-[#F8FAFC] flex flex-col relative overflow-hidden font-sans">
        
        <header class="bg-white/80 backdrop-blur-md border-b border-gray-100 flex items-center justify-between px-6 py-4 sticky top-0 z-40">
            <div class="flex items-center gap-5">
                <button @click="sidebarOpen = true" class="group flex flex-col gap-1.5 p-2 rounded-xl hover:bg-blue-50 transition-all">
                    <div class="w-6 h-0.5 bg-blue-600 rounded-full group-hover:w-4 transition-all"></div>
                    <div class="w-6 h-0.5 bg-blue-600 rounded-full"></div>
                    <div class="w-4 h-0.5 bg-blue-600 rounded-full group-hover:w-6 transition-all"></div>
                </button>
                <div class="flex flex-col">
                    <div class="font-black text-xl text-slate-900 tracking-tighter leading-none">
                        POS<span class="text-blue-600">UMKM</span>
                    </div>
                    <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">{{ user.storeName }}</span>
                </div>
            </div>
            
            <div class="flex items-center gap-4">
                <div class="hidden sm:flex items-center gap-3 pl-4 border-l border-gray-100">
                    <div class="text-right">
                        <div class="text-xs font-black text-slate-800 uppercase leading-none">{{ user.name.split(' ')[0] }}</div>
                        <span class="text-[9px] font-bold text-blue-600 uppercase tracking-tighter bg-blue-50 px-1.5 py-0.5 rounded-md mt-1 inline-block">{{ user.role }}</span>
                    </div>
                    <div class="w-9 h-9 rounded-xl bg-slate-900 flex items-center justify-center text-white text-xs font-black ring-4 ring-slate-50">
                        {{ user.name.substring(0, 2).toUpperCase() }}
                    </div>
                </div>
            </div>
        </header>

        <Transition name="fade">
            <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-slate-900/60 backdrop-blur-[2px] z-40 transition-all"></div>
        </Transition>

        <aside :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full'" class="fixed inset-y-0 left-0 w-80 bg-white transform transition-all duration-500 ease-[cubic-bezier(0.34,1.56,0.64,1)] z-50 flex flex-col border-r border-gray-100">
            
            <div class="p-8 flex items-center justify-between">
                <div class="flex flex-col">
                    <div class="font-black text-2xl text-slate-900 tracking-tighter leading-none">
                        POS<span class="text-blue-600">UMKM</span>
                    </div>
                    <span class="text-[10px] font-bold text-slate-400 uppercase tracking-[0.3em] mt-1">Management System</span>
                </div>
                <button @click="sidebarOpen = false" class="w-10 h-10 flex items-center justify-center rounded-2xl bg-slate-50 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <nav class="flex-1 px-6 space-y-8 overflow-y-auto custom-scrollbar pb-10">
                
                <div>
                    <div class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] px-2 mb-4 flex items-center gap-2">
                        <span class="w-1.5 h-1.5 bg-blue-500 rounded-full"></span> Menu Operasional
                    </div>
                    <div class="space-y-1">
                        <router-link to="/pos/kasir" @click="sidebarOpen = false" class="nav-link" :class="{ 'active': route.path === '/pos/kasir' }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                            <span>Mesin Kasir (POS)</span>
                        </router-link>

                        <router-link to="/riwayat" @click="sidebarOpen = false" class="nav-link" :class="{ 'active': route.path.startsWith('/riwayat') }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" /></svg>
                            <span>Jurnal Transaksi</span>
                        </router-link>

                        <router-link to="/absensi" @click="sidebarOpen = false" class="nav-link" :class="{ 'active': route.path === '/absensi' }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" /></svg>
                            <span>Presensi</span>
                        </router-link>

                        <router-link to="/schedule" @click="sidebarOpen = false" class="nav-link" :class="{ 'active': route.path === '/schedule' }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                            <span>Setting Shift (TSM)</span>
                        </router-link>
                    </div>
                </div>

                <div>
                    <button @click="toggleGroup('stock')" class="group-btn">
                        <span class="flex items-center gap-2"><span class="w-1.5 h-1.5 bg-amber-500 rounded-full"></span> Inventori & Stok</span>
                        <svg :class="openGroups.stock ? 'rotate-180' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    <div v-show="openGroups.stock" class="mt-2 space-y-1 ml-4 border-l-2 border-slate-100 pl-2">
                        <router-link to="/penerimaan-barang" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path === '/penerimaan-barang' }">Terima Barang (LPB)</router-link>
                        <router-link to="/stock-opname" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path === '/stock-opname' }">Stock Opname</router-link>
                    </div>
                </div>

                <div>
                    <button @click="toggleGroup('admin')" class="group-btn">
                        <span class="flex items-center gap-2"><span class="w-1.5 h-1.5 bg-indigo-500 rounded-full"></span> Administrasi Toko</span>
                        <svg :class="openGroups.admin ? 'rotate-180' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    <div v-show="openGroups.admin" class="mt-2 space-y-1 ml-4 border-l-2 border-slate-100 pl-2">
                        <template v-if="user.role === 'manager' || user.role === 'supervisor'">
                            <router-link to="/stock-opname/report" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path === '/stock-opname/report' }">
                                Laporan Hasil SO
                            </router-link>
                        </template>

                        <template v-if="user.role === 'owner'">
                            <router-link to="/dashboard" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path === '/dashboard' }">
                                Dashboard Analitik
                            </router-link>
                            <router-link to="/stock-opname/report" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path === '/stock-opname/report' }">
                                Laporan Hasil SO
                            </router-link>
                            <router-link to="/produk" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path.startsWith('/produk') }">
                                Master Data Produk
                            </router-link>
                            <router-link to="/karyawan" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path.startsWith('/karyawan') }">
                                Manajemen Karyawan
                            </router-link>
                            <router-link to="/setup" @click="sidebarOpen = false" class="sub-link" :class="{ 'active-sub': route.path.startsWith('/setup') }">
                                Pengaturan Sistem
                            </router-link>
                        </template>
                        
                        <div v-if="user.role === 'staff'" class="px-4 py-2 text-[10px] font-bold text-slate-400 italic">
                            Akses Administrasi Terbatas 🔒
                        </div>
                    </div>
                </div>
            </nav>

            <div class="p-6 bg-slate-50 border-t border-gray-100">
                <div class="flex items-center justify-between p-3 bg-white rounded-2xl shadow-sm border border-slate-100">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-xl bg-blue-600 flex items-center justify-center text-white font-black text-sm">
                            {{ user.name.substring(0, 2).toUpperCase() }}
                        </div>
                        <div class="flex flex-col">
                            <div class="text-xs font-black text-slate-800 uppercase leading-none">{{ user.name.split(' ')[0] }}</div>
                            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1">{{ user.role }}</span>
                        </div>
                    </div>
                    <button @click="logout" class="w-9 h-9 flex items-center justify-center rounded-xl text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3 3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                    </button>
                </div>
            </div>
        </aside>

        <main class="flex-1 w-full max-w-full overflow-y-auto bg-[#F8FAFC] h-[calc(100vh-73px)] relative scroll-smooth">
            <div class="p-4 sm:p-8">
                <slot /> 
            </div>
        </main>

    </div>
</template>

<style scoped>
/* Navigation Link Base */
.nav-link {
    @apply flex items-center gap-4 px-4 py-3.5 rounded-[20px] text-[13px] font-black tracking-tight transition-all duration-300 border border-transparent;
    @apply text-slate-500 hover:bg-white hover:border-slate-100 hover:shadow-sm hover:translate-x-1;
}

/* Icon Base */
.icon {
    @apply w-5 h-5 transition-transform duration-300;
}

/* Active State */
.nav-link.active {
    @apply bg-slate-900 text-white shadow-xl shadow-slate-200 translate-x-1;
}
.nav-link.active .icon {
    @apply text-blue-400;
}

/* Group Button Base */
.group-btn {
    @apply w-full flex items-center justify-between px-2 py-2 text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] hover:text-blue-600 transition-colors;
}

/* Sub-links Base */
.sub-link {
    @apply flex items-center gap-4 px-4 py-2.5 rounded-xl text-[12px] font-bold text-slate-500 transition-all duration-200 border-l-2 border-transparent;
    @apply hover:text-blue-600 hover:bg-blue-50/50 hover:pl-6;
}

.sub-link.active-sub {
    @apply text-blue-600 bg-blue-50 border-blue-600 pl-6;
}

/* Custom Transitions */
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

/* Custom Scrollbar */
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #E2E8F0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #CBD5E1; }
</style>