<script setup>
import { ref } from 'vue';
import { Link, usePage, router } from '@inertiajs/vue3';

// Tarik data global dari Inertia
const page = usePage();

// State sidebarOpen sekarang default-nya false untuk semua perangkat
const sidebarOpen = ref(false);

// Fungsi Logout
const logout = () => {
    router.post(route('logout'));
};
</script>

<template>
    <div class="min-h-screen bg-gray-50 flex flex-col relative overflow-hidden font-sans">

        <header class="bg-white shadow-sm border-b border-gray-200 flex items-center justify-between px-4 py-3 sticky top-0 z-40">
            <div class="flex items-center gap-4">
                <button @click="sidebarOpen = true" class="text-gray-500 hover:text-blue-600 focus:outline-none transition-colors p-2 rounded-xl hover:bg-blue-50 active:scale-95">
                    <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                </button>

                <div class="font-black text-2xl text-blue-600 tracking-tighter hidden sm:block drop-shadow-sm">
                    POS<span class="text-gray-800">UMKM</span>
                </div>
            </div>

            <div class="flex items-center gap-3">
                <div class="hidden sm:flex items-center gap-2.5 px-3 py-1.5 bg-gray-50 rounded-full border border-gray-200 shadow-sm cursor-pointer hover:bg-gray-100 transition-colors">
                    <div class="w-7 h-7 rounded-full bg-blue-600 flex items-center justify-center text-white text-xs font-black shadow-inner uppercase">
                        {{ $page.props.auth.user.name.substring(0, 2) }}
                    </div>
                    <div class="flex flex-col pr-2">
                        <span class="text-xs font-black text-gray-800 leading-none">{{ $page.props.auth.user.name }}</span>
                        <span class="text-[10px] font-bold text-gray-500 capitalize">{{ $page.props.auth.user.role }}</span>
                    </div>
                </div>
            </div>
        </header>

        <div
            v-if="sidebarOpen"
            @click="sidebarOpen = false"
            class="fixed inset-0 bg-gray-900/40 z-40 backdrop-blur-sm transition-opacity"
        ></div>

        <div
            :class="sidebarOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full'"
            class="fixed inset-y-0 left-0 w-72 bg-white border-r border-gray-100 transform transition-all duration-300 ease-[cubic-bezier(0.4,0,0.2,1)] z-50 flex flex-col"
        >
            <div class="flex items-center justify-between h-[61px] px-6 border-b border-gray-100 bg-white shrink-0">
                <div class="font-black text-2xl text-blue-600 tracking-tighter drop-shadow-sm">
                    POS<span class="text-gray-800">UMKM</span>
                </div>
                <button @click="sidebarOpen = false" class="p-2 -mr-2 text-gray-400 hover:text-red-500 rounded-xl hover:bg-red-50 transition-colors active:scale-95">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <nav class="flex-1 px-4 py-6 space-y-2.5 overflow-y-auto custom-scrollbar">

                <div class="text-xs font-black text-gray-400 uppercase tracking-widest px-4 mb-4">Menu Utama</div>

                <Link
                    href="/pos"
                    @click="sidebarOpen = false"
                    class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all"
                    :class="$page.url === '/pos' ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50 hover:text-blue-600'"
                >
                    <span class="text-xl">🛒</span>
                    POS Kasir
                </Link>

                <Link
                    v-if="$page.props.auth.store?.fitur_opsional?.includes('absensi')"
                    href="/absensi"
                    @click="sidebarOpen = false"
                    class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all"
                    :class="$page.url.startsWith('/absensi') ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50 hover:text-blue-600'"
                >
                    <span class="text-xl">⏱️</span>
                    Absensi Karyawan
                </Link>

                <div v-if="$page.props.auth.user.role === 'owner'" class="pt-4 mt-4 border-t border-gray-100">
                    <div class="text-xs font-black text-gray-400 uppercase tracking-widest px-4 mb-4">Administrasi Toko</div>

                    <Link
                        href="/dashboard"
                        @click="sidebarOpen = false"
                        class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all"
                        :class="$page.url.startsWith('/dashboard') ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50 hover:text-blue-600'"
                    >
                        <span class="text-xl">📊</span>
                        Dashboard
                    </Link>

                    <Link
                        href="/products"
                        @click="sidebarOpen = false"
                        class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all"
                        :class="$page.url.startsWith('/products') ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50 hover:text-blue-600'"
                    >
                        <span class="text-xl">📦</span>
                        Master Produk
                    </Link>

                    <Link
                        href="/karyawan"
                        @click="sidebarOpen = false"
                        class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all"
                        :class="$page.url.startsWith('/karyawan') ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50 hover:text-blue-600'"
                    >
                        <span class="text-xl">👥</span>
                        Manajemen Karyawan
                    </Link>

                    <Link
                        href="/pengaturan"
                        @click="sidebarOpen = false"
                        class="flex items-center gap-4 px-4 py-3.5 rounded-2xl text-sm font-bold transition-all"
                        :class="$page.url.startsWith('/pengaturan') ? 'bg-blue-600 text-white shadow-md shadow-blue-200' : 'text-gray-600 hover:bg-gray-50 hover:text-blue-600'"
                    >
                        <span class="text-xl">⚙️</span>
                        Pengaturan Toko
                    </Link>
                </div>
            </nav>

            <div class="p-5 border-t border-gray-100 bg-gray-50 shrink-0">

                <div class="mb-4 pb-4 border-b border-gray-200">
                    <div class="text-[10px] text-gray-400 font-black uppercase tracking-widest mb-1.5">Lokasi Kerja</div>
                    <div class="flex items-center gap-2 text-sm font-black text-blue-900">
                        <span>🏪</span>
                        {{ $page.props.auth.store?.nama_toko || 'Belum Ada Toko' }}
                    </div>
                </div>

                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-full bg-blue-100 border-2 border-white shadow-sm overflow-hidden flex items-center justify-center">
                            <span class="text-blue-700 font-black text-sm uppercase">{{ $page.props.auth.user.name.substring(0, 2) }}</span>
                        </div>
                        <div>
                            <div class="text-sm font-black text-gray-800 leading-tight">{{ $page.props.auth.user.name.split(' ')[0] }}</div>
                            <div class="text-[10px] font-bold text-gray-500 capitalize uppercase tracking-wider">{{ $page.props.auth.user.role }}</div>
                        </div>
                    </div>

                    <button @click="logout" title="Keluar" class="w-10 h-10 rounded-full bg-white border border-gray-200 text-gray-500 hover:text-red-600 hover:border-red-200 hover:bg-red-50 flex items-center justify-center transition-all shadow-sm active:scale-95">
                        <svg class="w-5 h-5 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                        </svg>
                    </button>
                </div>

            </div>
        </div>

        <main class="flex-1 w-full max-w-full overflow-y-auto bg-gray-50 h-[calc(100vh-61px)]">
            <transition name="fade">
                <div v-if="$page.props.flash.message" class="m-6 p-4 bg-green-500 text-black rounded-2xl shadow-lg flex items-center justify-between z-50">
                    <div class="flex items-center gap-3">
                        <span class="text-xl">✅</span>
                        <span class="font-bold tracking-wide">{{ $page.props.flash.message }}</span>
                    </div>
                    <button @click="$page.props.flash.message = null" class="text-white/70 hover:text-white font-black text-lg active:scale-95">✕</button>
                </div>
            </transition>
            <slot />
        </main>

    </div>
</template>

<style scoped>
/* Custom Scrollbar biar menunya halus */
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }

/* Transisi untuk Flash Message */
.fade-enter-active, .fade-leave-active { transition: opacity 0.4s ease, transform 0.4s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(-10px); }
</style>
