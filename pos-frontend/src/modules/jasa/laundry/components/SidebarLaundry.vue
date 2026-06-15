<script setup>
import Swal from 'sweetalert2';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../../../api.js'; 
import { useSidebar } from '../composables/useSidebar.js';

const router = useRouter();
const { route, sidebarOpen, openGroups, user, toggleGroup, logout } = useSidebar();

const subPlan = ref(localStorage.getItem('subscriptionPlan') || 'basic');

const getPlanLevel = (plan) => {
    const p = plan ? plan.toLowerCase() : 'basic';
    if (p === 'premium' || p === 'trial') return 3;
    if (p === 'pro') return 2;
    return 1;
};
const planLevel = computed(() => getPlanLevel(subPlan.value));

// 🛡️ LOCKING GATEWAY SYSTEM: Alert paywall pembatasan fitur premium ruko laundry bray bray!
const triggerUpgrade = (fiturName, minLevel) => {
    sidebarOpen.value = false;
    let targetPlan = minLevel === 2 ? 'PRO' : 'PREMIUM';

    Swal.fire({
        icon: 'warning',
        title: 'Fitur Terkunci',
        html: `Fitur <b>${fiturName}</b> eksklusif untuk paket <b>${targetPlan}</b>.<br><br>Tingkatkan paket Anda sekarang untuk membuka potensi maksimal bisnis.`,
        confirmButtonText: 'Upgrade Sekarang',
        showCancelButton: true,
        cancelButtonText: 'Nanti Dulu',
        confirmButtonColor: '#4f46e5',
        cancelButtonColor: '#94a3b8',
        customClass: { popup: 'rounded-[32px]' },
    }).then((result) => {
        if (result.isConfirmed) router.push('/laundry/akun');
    });
};

const updatePlanRealtime = () => {
    subPlan.value = localStorage.getItem('subscriptionPlan') || 'basic';
};

// 🚀 NAVIGATION SWITCH MULTI-BRANCH: Balik ke gerbang pilih cabang instan tanpa logout bray!
const goToSelectStore = () => {
    sidebarOpen.value = false;
    router.push('/select-store');
};

onMounted(async () => {
    window.addEventListener('store-updated', updatePlanRealtime);
    window.addEventListener('storage', updatePlanRealtime);

    const role = localStorage.getItem('role') || user.role;

    try {
        const res = await api.get('/laundry/setting');
        const data = res.data;
        const status = data.subscription_status || 'active';

        if (data.nama_toko) {
            user.storeName = data.nama_toko;
            localStorage.setItem('storeName', data.nama_toko);
        }

        if (data.qris_image) {
            user.storeLogo = data.qris_image;
            localStorage.setItem('storeLogo', data.qris_image);
        }

        if (localStorage.getItem('subscriptionPlan')) {
            subPlan.value = localStorage.getItem('subscriptionPlan');
        }

        let isDead = false;
        if (status !== 'active' && status !== 'ACTIVE') {
            isDead = true;
        }

        if (isDead) {
            if (role === 'owner' && route.path !== '/laundry/akun') {
                router.push('/laundry/akun');
            } else if (role !== 'owner' && route.path !== '/laundry/pos') {
                router.push('/laundry/pos');
            }
        }
    } catch (e) {
        if (e.response && (e.response.status === 402 || e.response.status === 403)) {
            if (role === 'owner' && route.path !== '/laundry/akun') {
                router.push('/laundry/akun');
            } else if (role !== 'owner' && route.path !== '/laundry/pos') {
                router.push('/laundry/pos');
            }
        }
    }
});

onUnmounted(() => {
    window.removeEventListener('store-updated', updatePlanRealtime);
    window.removeEventListener('storage', updatePlanRealtime);
});
</script>

<template>
    <div class="h-screen w-screen flex flex-col bg-[#F8FAFC] font-sans overflow-hidden relative selection:bg-indigo-100 selection:text-indigo-600">
        <!-- 💻 MASTER TOP BANNER HEADER -->
        <header class="bg-white/80 backdrop-blur-md border-b border-slate-100 flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4 sticky top-0 z-40 shadow-sm shrink-0 no-print">
            <div class="flex items-center gap-4 sm:gap-5">
                <button @click="sidebarOpen = true" class="group flex flex-col gap-1.5 p-2 rounded-xl hover:bg-indigo-50 transition-all active:scale-95">
                    <div class="w-6 h-0.5 bg-indigo-600 rounded-full group-hover:w-4 transition-all"></div>
                    <div class="w-6 h-0.5 bg-indigo-600 rounded-full"></div>
                    <div class="w-6 h-0.5 bg-indigo-600 rounded-full group-hover:w-6 transition-all"></div>
                </button>

                <!-- 🚀 FIX: Logo/Nama Toko di Header Atas Kalau Diklik Bisa Auto-Pindah Cabang bray! -->
                <div class="flex items-center gap-3 cursor-pointer hover:opacity-80 active:scale-98 transition-all" title="Ganti / Pilih Cabang Toko Lain" @click="goToSelectStore">
                    <img v-if="user.storeLogo && user.storeLogo !== 'null' && user.storeLogo !== ''" :src="user.storeLogo" class="h-11 sm:h-14 max-w-[200px] object-contain origin-left transition-all" alt="Logo Toko" />
                    <div v-else class="font-black text-lg sm:text-xl text-slate-900 tracking-tighter leading-none">
                        <span class="text-indigo-600">{{ user.storeName || 'ARZU LAUNDRY' }}</span>
                    </div>
                    <span class="text-[8px] px-2 py-0.5 rounded-md font-black uppercase tracking-widest text-white shadow-sm shrink-0 self-center" :class="planLevel === 3 ? 'bg-amber-500' : planLevel === 2 ? 'bg-indigo-500' : 'bg-slate-400'">{{ subPlan }}</span>
                </div>
            </div>

            <div class="flex items-center gap-4">
                <div class="flex items-center gap-3 sm:pl-4 sm:border-l border-slate-100">
                    <div class="text-right hidden sm:block">
                        <div class="text-xs font-black text-slate-800 uppercase leading-none">{{ user.name ? user.name.split(' ')[0] : 'OWNER' }}</div>
                        <span class="text-[9px] font-bold text-indigo-600 uppercase tracking-tighter bg-indigo-50 px-1.5 py-0.5 rounded-md mt-1 inline-block border border-indigo-100">{{ user.role }}</span>
                    </div>
                    <div class="w-9 h-9 sm:w-10 sm:h-10 rounded-[14px] bg-gradient-to-br from-slate-800 to-slate-900 flex items-center justify-center text-white text-xs sm:text-sm font-black shadow-md border-2 border-white ring-2 ring-slate-100 overflow-hidden">
                        <img v-if="user.foto_url && user.foto_url !== 'null' && user.foto_url !== ''" :src="user.foto_url" class="w-full h-full object-cover" />
                        <span v-else>{{ user.name ? user.name.substring(0, 2).toUpperCase() : 'LD' }}</span>
                    </div>
                </div>
            </div>
        </header>

        <!-- Slide Backdrop Overlay bray bray bray -->
        <Transition name="fade">
            <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-slate-950/40 backdrop-blur-sm z-40 transition-all"></div>
        </Transition>

        <!-- 🌌 ASIDE DRAWER PANEL (DARK INDIGO MEWAH - #1e1b4b) -->
        <aside :class="sidebarOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full'" class="no-print fixed inset-y-0 left-0 w-[280px] sm:w-[320px] bg-[#1e1b4b] text-white transform transition-transform duration-500 ease-[cubic-bezier(0.34,1.56,0.64,1)] z-50 flex flex-col border-r border-slate-800 h-full">
            
            <!-- 🚀 FIX: Logo/Nama Di Dalam Drawer Kalau Diklik Juga Bisa Auto-Pindah Cabang bray! -->
            <div class="p-6 sm:p-8 flex items-center justify-between bg-black/10 border-b border-white/10 shrink-0 cursor-pointer hover:bg-black/20 transition-colors" title="Ganti / Pilih Cabang Toko Lain" @click="goToSelectStore">
                <div class="flex flex-col">
                    <img v-if="user.storeLogo && user.storeLogo !== 'null' && user.storeLogo !== ''" :src="user.storeLogo" class="h-16 sm:h-20 max-w-[240px] object-contain mb-2 origin-left" alt="Logo Toko" />
                    <div v-else class="font-black text-xl sm:text-2xl text-white tracking-tighter leading-none">{{ user.storeName || 'ARZU LAUNDRY' }}</div>
                    <span class="text-[9px] sm:text-[10px] font-black text-indigo-400 uppercase tracking-[0.3em] mt-1.5">LAUNDRY POS SYSTEM</span>
                </div>
            </div>

            <!-- Area Scrollable Nav Link Item bray -->
            <nav class="flex-1 px-4 sm:px-6 py-6 space-y-8 overflow-y-auto custom-scrollbar min-h-0">
                <!-- SEGMENT OPERASIONAL -->
                <div>
                    <div class="text-[9px] font-black text-indigo-400 uppercase tracking-[0.2em] px-2 mb-4 flex items-center gap-2">
                        <span class="w-1.5 h-1.5 bg-indigo-400 rounded-full shadow-[0_0_8px_rgba(129,140,248,0.6)]"></span>
                        Operasional Ruko
                    </div>
                    <div class="space-y-1.5">
                        <router-link to="/laundry/pos" @click="sidebarOpen = false" class="nav-link" :class="{ active: route.path === '/laundry/pos' }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
                            <span>Kasir Laundry</span>
                        </router-link>
                        <router-link to="/laundry/pos/riwayat" @click="sidebarOpen = false" class="nav-link" :class="{ active: route.path.startsWith('/laundry/pos/riwayat') }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5H7a2 2 0 00-2-2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" /></svg>
                            <span>Riwayat Transaksi</span>
                        </router-link>
                        <router-link to="/laundry/status" @click="sidebarOpen = false" class="nav-link" :class="{ active: route.path === '/laundry/status' }">
                            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z M12 6v6l4 2" /></svg>
                            <span>Status Antrean Rak</span>
                        </router-link>
                    </div>
                </div>

                <!-- SEGMENT KATALOG -->
                <div>
                    <button @click="toggleGroup('stock')" class="group-btn hover:text-emerald-400">
                        <span class="flex items-center gap-2">
                            <span class="w-1.5 h-1.5 bg-emerald-500 rounded-full shadow-[0_0_8px_rgba(16,185,129,0.6)]"></span>
                            Katalog Jasa & Paket
                        </span>
                        <svg :class="openGroups.stock ? 'rotate-180 text-emerald-400' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    <div v-show="openGroups.stock" class="mt-2 space-y-1 ml-4 border-l-2 border-white/10 pl-2">
                        <router-link v-if="user.role === 'owner'" to="/laundry/master-layanan" @click="sidebarOpen = false" class="sub-link hover:text-emerald-400 hover:bg-white/5 hover:border-emerald-500" :class="{ 'active-sub !text-emerald-400 !bg-white/10 !border-emerald-500': route.path.startsWith('/laundry/master-layanan') }">Master Katalog Harga</router-link>
                    </div>
                </div>

                <!-- SEGMENT ADM / FINANCE -->
                <div v-if="user.role !== 'staff'">
                    <button @click="toggleGroup('admin')" class="group-btn hover:text-amber-400">
                        <span class="flex items-center gap-2">
                            <span class="w-1.5 h-1.5 bg-amber-500 rounded-full shadow-[0_0_8px_rgba(245,158,11,0.6)]"></span>
                            Administrasi Finansial
                        </span>
                        <svg :class="openGroups.admin ? 'rotate-180 text-amber-400' : ''" class="w-3 h-3 transition-transform" fill="none" stroke="currentColor" stroke-width="3" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
                    </button>
                    <div v-show="openGroups.admin" class="mt-2 space-y-1 ml-4 border-l-2 border-white/10 pl-2">
                        <template v-if="user.role === 'owner'">
                            <a v-if="planLevel < 3" href="#" @click.prevent="triggerUpgrade('Laporan Omset Lanjutan', 3)" class="sub-link-locked">
                                <span class="flex-1">Laporan Omset & Grafik</span>
                                <span class="text-amber-500">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2" /><path d="M7 11V7a5 5 0 0 1 10 0v4" /></svg>
                                </span>
                            </a>
                            <router-link v-else to="/laundry/laporan" @click="sidebarOpen = false" class="sub-link hover:text-amber-400 hover:bg-white/5 hover:border-emerald-500" :class="{ 'active-sub !text-amber-400 !bg-white/10 !border-emerald-500': route.path === '/laundry/laporan' }">Laporan Omset & Grafik</router-link>
                            
                            <a v-if="planLevel < 2" href="#" @click.prevent="triggerUpgrade('Manajemen Karyawan Jasa', 2)" class="sub-link-locked">
                                <span class="flex-1">Manajemen Karyawan</span>
                                <span class="text-amber-500">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2" /><path d="M7 11V7a5 5 0 0 1 10 0v4" /></svg>
                                </span>
                            </a>
                            <router-link v-else to="/laundry/karyawan" @click="sidebarOpen = false" class="sub-link hover:text-amber-400 hover:bg-white/5 hover:border-emerald-500" :class="{ 'active-sub !text-amber-400 !bg-white/10 !border-emerald-500': route.path === '/laundry/karyawan' }">Manajemen Karyawan</router-link>
                        </template>
                    </div>
                </div>

                <!-- SEGMENT SYSTEM CLOUD -->
                <div>
                    <div class="text-[9px] font-black text-indigo-400 uppercase tracking-[0.2em] px-2 mb-4 mt-6 flex items-center gap-2">
                        <span class="w-1.5 h-1.5 bg-blue-400 rounded-full shadow-[0_0_8px_rgba(59,130,246,0.6)]"></span>
                        Sistem & Akun Cloud
                    </div>
                    <div class="space-y-1.5">
                        <template v-if="user.role === 'owner'">
                            <a v-if="planLevel < 3" href="#" @click.prevent="triggerUpgrade('Integrasi WhatsApp Bot Nota', 3)" class="nav-link-locked">
                                <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 21l1.65-3.8a9 9 0 1 1 3.4 2.9L3 21" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 10a.5.5 0 0 0 1 0V9a.5.5 0 0 0-1 0v1a5 5 0 0 0 5 5h1a.5.5 0 0 0 0-1h-1a.5.5 0 0 0 0 1" /></svg>
                                <span class="flex-1">Integrasi WA</span>
                                <span class="text-amber-500">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2" /><path d="M7 11V7a5 5 0 0 1 10 0v4" /></svg>
                                </span>
                            </a>
                            <router-link v-else to="/laundry/whatsapp" @click="sidebarOpen = false" class="nav-link group" :class="{ active: route.path === '/laundry/whatsapp' }">
                                <svg class="icon group-hover:rotate-12 group-hover:text-emerald-400 transition-all" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 21l1.65-3.8a9 9 0 1 1 3.4 2.9L3 21" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 10a.5.5 0 0 0 1 0V9a.5.5 0 0 0-1 0v1a5 5 0 0 0 5 5h1a.5.5 0 0 0 0-1h-1a.5.5 0 0 0 0 1" /></svg>
                                <span>Integrasi WA</span>
                            </router-link>
                            
                            <router-link to="/laundry/setting" @click="sidebarOpen = false" class="nav-link group" :class="{ active: route.path === '/laundry/setting' }">
                                <svg class="icon group-hover:rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                                <span>Pengaturan Toko</span>
                            </router-link>
                        </template>
                        <router-link to="/laundry/akun" @click="sidebarOpen = false" class="nav-link group" :class="{ active: route.path === '/laundry/akun' }">
                            <svg class="icon group-hover:scale-110" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
                            <span>Akun Saya</span>
                        </router-link>
                    </div>
                </div>
            </nav>

            <!-- Bagian Profil Bawah & Tombol Logout bray -->
            <div class="p-4 sm:p-6 bg-black/10 border-t border-white/10 shrink-0">
                <div class="flex items-center justify-between p-3 sm:p-4 bg-white/5 rounded-2xl border border-white/5 hover:border-white/10 transition-colors">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-[14px] bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center text-white font-black text-sm shadow-md overflow-hidden">
                            <img v-if="user.foto_url && user.foto_url !== 'null' && user.foto_url !== ''" :src="user.foto_url" class="w-full h-full object-cover" />
                            <span v-else>{{ user.name ? user.name.substring(0, 2).toUpperCase() : 'LD' }}</span>
                        </div>
                        <div class="flex flex-col max-w-[100px] sm:max-w-[120px]">
                            <div class="text-xs font-black text-white uppercase leading-none truncate">{{ user.name || 'OWNER' }}</div>
                            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-1 truncate">{{ user.role }}</span>
                        </div>
                    </div>
                    <button @click="logout" class="w-9 h-9 sm:w-10 sm:h-10 flex items-center justify-center rounded-xl bg-white/5 text-slate-400 hover:text-white hover:bg-rose-600 transition-all active:scale-95 group" title="Logout">
                        <svg class="w-4 h-4 sm:w-5 sm:h-5 group-hover:translate-x-0.5 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                    </button>
                </div>
            </div>
        </aside>

        <!-- Area Dynamic Slot Render Konten Utama -->
        <main class="flex-1 w-full overflow-y-auto bg-[#F8FAFC] relative scroll-smooth custom-scrollbar print:overflow-visible print:bg-white">
            <slot />
        </main>
    </div>
</template>

<style scoped>
.nav-link {
    @apply flex items-center gap-3.5 px-4 py-3.5 rounded-[16px] text-xs font-black tracking-tight transition-all duration-300 border-2 border-transparent text-slate-400 hover:bg-white/5 hover:text-white;
}
.icon {
    @apply w-[18px] h-[18px] transition-transform duration-300;
}
.nav-link.active {
    @apply bg-indigo-600 text-white shadow-xl shadow-indigo-900/50 border-indigo-500 translate-x-1;
}
.nav-link.active .icon {
    @apply text-white;
}
.nav-link-locked {
    @apply flex items-center gap-3.5 px-4 py-3.5 rounded-[16px] text-xs font-black tracking-tight border-2 border-transparent text-slate-500 bg-black/10 hover:bg-black/20 cursor-not-allowed opacity-60;
}
.sub-link-locked {
    @apply flex items-center gap-4 px-4 py-2.5 rounded-xl text-[11px] font-bold text-slate-500 border-l-[3px] border-transparent hover:bg-black/10 cursor-not-allowed opacity-60;
}
.group-btn {
    @apply w-full flex items-center justify-between px-2 py-2 text-[9px] font-black text-slate-400 uppercase tracking-[0.2em] transition-colors;
}
.sub-link {
    @apply flex items-center gap-4 px-4 py-2.5 rounded-xl text-[11px] font-bold text-slate-400 transition-all duration-200 border-l-[3px] border-transparent;
}
.active-sub {
    @apply border-indigo-500 text-indigo-400 font-black;
}
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
input[type='number']::-webkit-inner-spin-button, input[type='number']::-webkit-outer-spin-button { -webkit-appearance: none; margin: 0; }
input[type='number'] { -moz-appearance: textfield; }
</style>