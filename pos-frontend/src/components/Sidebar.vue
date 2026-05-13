<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Swal from 'sweetalert2';

const router = useRouter();
const route = useRoute();

const sidebarOpen = ref(false);

// State untuk buka-tutup grup menu
const openGroups = ref({
    stock: true, // Default terbuka
    admin: false
});

const toggleGroup = (group) => {
    openGroups.value[group] = !openGroups.value[group];
};

const user = ref({
    name: localStorage.getItem('name') || 'User',
    role: localStorage.getItem('role') || 'kasir',
    storeName: localStorage.getItem('storeName') || 'Indo UMKM'
});

onMounted(() => {
    user.value.name = localStorage.getItem('name') || 'User';
    user.value.role = localStorage.getItem('role') || 'kasir';
});

const logout = () => {
    Swal.fire({
        title: 'Mau keluar, Bos?',
        text: "Pastikan semua kerjaan sudah beres ya!",
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#2563eb',
        cancelButtonColor: '#64748b',
        confirmButtonText: 'Ya, Logout Sekarang',
        cancelButtonText: 'Batal',
        reverseButtons: true
    }).then((result) => {
        if (result.isConfirmed) {
            localStorage.clear();
            router.push('/login');
        }
    });
};
</script>

<template>
    <div class="min-h-screen bg-gray-50 flex flex-col relative overflow-hidden font-sans">
        
        <header class="bg-white shadow-sm border-b border-gray-200 flex items-center justify-between px-4 py-3 sticky top-0 z-40">
            <div class="flex items-center gap-4">
                <button @click="sidebarOpen = true" class="text-gray-500 hover:text-blue-600 p-2 rounded-xl hover:bg-blue-50 active:scale-95 transition-all">
                    <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                </button>
                <div class="font-black text-2xl text-blue-600 tracking-tighter hidden sm:block">
                    POS<span class="text-gray-800">UMKM</span>
                </div>
            </div>
            <div class="flex items-center gap-3">
                <div class="hidden sm:flex items-center gap-2.5 px-3 py-1.5 bg-gray-50 rounded-full border border-gray-200">
                    <div class="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center text-white text-[10px] font-black uppercase">
                        {{ user.name.substring(0, 2) }}
                    </div>
                    <div class="flex flex-col pr-2 text-left">
                        <span class="text-xs font-black text-gray-800 leading-none uppercase">{{ user.name.split(' ')[0] }}</span>
                        <span class="text-[10px] font-bold text-gray-500 uppercase">{{ user.role }}</span>
                    </div>
                </div>
            </div>
        </header>

        <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-gray-900/40 z-40 backdrop-blur-sm transition-opacity"></div>

        <div :class="sidebarOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full'" class="fixed inset-y-0 left-0 w-72 bg-white border-r border-gray-100 transform transition-all duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] z-50 flex flex-col">
            
            <div class="flex items-center justify-between h-[61px] px-6 border-b border-gray-100 bg-white">
                <div class="font-black text-2xl text-blue-600 tracking-tighter">POS<span class="text-gray-800">UMKM</span></div>
                <button @click="sidebarOpen = false" class="p-2 text-gray-400 hover:text-red-500 rounded-xl hover:bg-red-50 transition-colors">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <nav class="flex-1 px-4 py-6 space-y-1.5 overflow-y-auto custom-scrollbar">
                
                <div class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] px-4 mb-2">Main Menu</div>
                <router-link to="/pos/kasir" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-3 rounded-2xl text-sm font-bold transition-all" :class="route.path === '/pos/kasir' ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50'">
                    <span class="text-lg">🛒</span> POS Kasir
                </router-link>
                <router-link to="/riwayat" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-3 rounded-2xl text-sm font-bold transition-all" :class="route.path.startsWith('/riwayat') ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50'">
                    <span class="text-lg">📜</span> Riwayat
                </router-link>
                <router-link to="/absensi" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-3 rounded-2xl text-sm font-bold transition-all" :class="route.path === '/absensi' ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50'">
                    <span class="text-lg">📸</span> Absensi
                </router-link>

                <div class="pt-4">
                    <button @click="toggleGroup('stock')" class="w-full flex items-center justify-between px-4 py-2 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] hover:text-blue-600 transition-colors">
                        <span>📦 Stock & Inventory</span>
                        <svg :class="openGroups.stock ? 'rotate-180' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    
                    <div v-show="openGroups.stock" class="mt-1 space-y-1 ml-2 border-l-2 border-gray-50">
                        <router-link to="/penerimaan-barang" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path === '/penerimaan-barang' ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                            <span class="text-base">🚚</span> Terima Barang
                        </router-link>
                        <router-link to="/stock-opname" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path === '/stock-opname' ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                            <span class="text-base">🔍</span> Stock Opname
                        </router-link>
                        
                    </div>
                </div>

                <div class="pt-4">
                    <button @click="toggleGroup('admin')" class="w-full flex items-center justify-between px-4 py-2 text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] hover:text-blue-600 transition-colors">
                        <span>⚙️ Administrasi</span>
                        <svg :class="openGroups.admin ? 'rotate-180' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 9l-7 7-7-7" /></svg>
                    </button>

                    <div v-show="openGroups.admin" class="mt-1 space-y-1 ml-2 border-l-2 border-gray-50">
                        <template v-if="user.role === 'owner'">
                            <router-link to="/dashboard" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path.startsWith('/dashboard') ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                                <span class="text-base">📊</span> Dashboard
                            </router-link>
                            <router-link to="/stock-opname/report" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path === '/stock-opname/report' ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                <span class="text-base">📉</span> Laporan Selisih SO
            </router-link>
                            <router-link to="/produk" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path.startsWith('/produk') ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                            <span class="text-base">📁</span> Master Produk
                            </router-link>
                            <router-link to="/karyawan" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path.startsWith('/karyawan') ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                                <span class="text-base">👥</span> Karyawan
                            </router-link>
                            <router-link to="/setup" @click="sidebarOpen = false" class="flex items-center gap-4 px-4 py-2.5 rounded-xl text-sm font-bold transition-all" :class="route.path.startsWith('/setup') ? 'text-blue-600 bg-blue-50' : 'text-gray-500 hover:text-blue-600'">
                                <span class="text-base">⚙️</span> Toko
                            </router-link>
                        </template>
                        <div v-else class="px-4 py-2 text-[10px] font-bold text-gray-400 italic">Menu ini hanya untuk Owner 🔒</div>
                    </div>
                </div>

            </nav>

            <div class="p-5 border-t border-gray-100 bg-gray-50 shrink-0">
                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-full bg-blue-100 border-2 border-white shadow-sm flex items-center justify-center">
                            <span class="text-blue-700 font-black text-sm uppercase">{{ user.name.substring(0, 2) }}</span>
                        </div>
                        <div class="text-left">
                            <div class="text-sm font-black text-gray-800 leading-tight uppercase">{{ user.name.split(' ')[0] }}</div>
                            <div class="text-[10px] font-bold text-gray-500 uppercase tracking-widest">{{ user.role }}</div>
                        </div>
                    </div>
                    <button @click="logout" class="w-10 h-10 rounded-full bg-white border border-gray-200 text-gray-500 hover:text-red-600 hover:bg-red-50 flex items-center justify-center transition-all shadow-sm active:scale-95">
                        <svg class="w-5 h-5 ml-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                    </button>
                </div>
            </div>
        </div>

        <main class="flex-1 w-full max-w-full overflow-y-auto bg-gray-50 h-[calc(100vh-61px)] relative">
            <slot /> 
        </main>

    </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }
</style>