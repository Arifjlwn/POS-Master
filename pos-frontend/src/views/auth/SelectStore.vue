<script setup>
import Swal from 'sweetalert2';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';

const router = useRouter();
const stores = ref([]);
const userName = ref('');
const isLoading = ref(false);

onMounted(() => {
    const tempStores = localStorage.getItem('temp_stores');
    const name = localStorage.getItem('temp_name');

    if (!tempStores) {
        router.push('/login');
        return;
    }

    stores.value = JSON.parse(tempStores);
    userName.value = name || 'Owner';
});

const selectBranch = async (storeId) => {
    isLoading.value = true;
    try {
        const res = await api.post('/auth/select-store', {
            store_id: storeId
        });

        // 1. TIMPA TOKEN LAMA DENGAN TOKEN FINAL
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('store_id', res.data.store_id);
        localStorage.setItem('store_name', res.data.store_name);
        localStorage.setItem('subscriptionPlan', res.data.subscription_plan);

        // 🚀 INI DIA YANG KETINGGALAN: SET ROLE LU JADI OWNER LAGI!
        // Ambil role dari respon Golang (yang mana udah kita set 'owner' dari token)
        const userRole = res.data.role ? res.data.role.toLowerCase() : 'owner';
        localStorage.setItem('role', userRole);

        // TANGKAP NAMA DAN FOTO DARI GOLANG!
        localStorage.setItem('name', res.data.name || 'Owner');
        localStorage.setItem('foto_url', res.data.foto_url || '');

        // 2. BERSIHKAN MEMORY SEMENTARA
        localStorage.removeItem('temp_stores');
        localStorage.removeItem('temp_name');

        // 3. MASUK KE DASHBOARD UTAMA
        router.push('/retail/dashboard');
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Akses Ditolak',
            text: error.response?.data?.error || 'Gagal masuk ke cabang ini.',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[32px]' }
        });
    } finally {
        isLoading.value = false;
    }
};

const getPlanStyle = (plan) => {
    const p = plan ? plan.toLowerCase() : 'basic';
    if (p === 'premium') return 'bg-amber-50 text-amber-700 border-amber-200 ring-amber-500';
    if (p === 'pro') return 'bg-indigo-50 text-indigo-700 border-indigo-200 ring-indigo-500';
    return 'bg-sky-50 text-sky-700 border-sky-200 ring-sky-500';
};
</script>

<template>
    <div class="min-h-screen bg-slate-50 flex flex-col items-center justify-center p-4 md:p-8 font-sans relative overflow-hidden">
        <div
            class="absolute top-[-10%] left-[-10%] w-96 h-96 bg-indigo-400 rounded-full mix-blend-multiply filter blur-[128px] opacity-20 animate-blob"></div>
        <div
            class="absolute bottom-[-10%] right-[-10%] w-96 h-96 bg-sky-400 rounded-full mix-blend-multiply filter blur-[128px] opacity-20 animate-blob animation-delay-2000"></div>

        <div class="w-full max-w-6xl z-10">
            <div class="text-center mb-12">
                <div class="inline-flex items-center justify-center p-3 bg-white rounded-2xl shadow-sm border border-slate-100 mb-6">
                    <div class="font-black text-2xl text-slate-900 tracking-tighter leading-none">
                        NEXA
                        <span class="text-indigo-600">POS</span>
                    </div>
                </div>
                <h1 class="text-3xl md:text-5xl font-black text-slate-800 tracking-tight mb-3">Selamat Datang, {{ userName.split(' ')[0] }}</h1>
                <p class="text-slate-500 font-bold text-sm uppercase tracking-widest">Silakan pilih cabang operasional yang ingin Anda akses</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <div
                    v-for="store in stores"
                    :key="store.id"
                    @click="selectBranch(store.id)"
                    class="group bg-white rounded-[32px] p-6 border border-slate-100 shadow-sm hover:shadow-2xl hover:-translate-y-2 transition-all duration-300 cursor-pointer relative overflow-hidden flex flex-col">
                    <div
                        class="absolute top-0 left-0 w-full h-1.5 transition-all duration-300"
                        :class="getPlanStyle(store.subscription_plan).split(' ')[3].replace('ring-', 'bg-')"></div>

                    <div class="flex items-start justify-between mb-6">
                        <div
                            class="w-14 h-14 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-center overflow-hidden shrink-0 group-hover:scale-110 transition-transform duration-300">
                            <img
                                v-if="store.logo_url"
                                :src="store.logo_url.startsWith('http') ? store.logo_url : 'http://localhost:8080' + store.logo_url"
                                class="w-full h-full object-cover" />
                            <svg v-else class="w-6 h-6 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                            </svg>
                        </div>

                        <span
                            class="text-[9px] font-black uppercase tracking-widest px-3 py-1.5 rounded-lg border shadow-sm"
                            :class="getPlanStyle(store.subscription_plan)">
                            {{ store.subscription_plan || 'BASIC' }}
                        </span>
                    </div>

                    <div class="flex-1">
                        <h3 class="text-xl font-black text-slate-800 mb-1 line-clamp-1">
                            {{ store.nama_toko }}
                        </h3>
                        <p class="text-xs font-bold text-slate-400 uppercase tracking-wider line-clamp-1 flex items-center gap-1">
                            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                            </svg>
                            {{ store.kota || 'Lokasi Belum Diatur' }}
                        </p>
                    </div>

                    <div class="mt-8 flex items-center justify-between border-t border-slate-50 pt-4">
                        <div class="text-[10px] font-black text-slate-400 uppercase tracking-widest">ID: {{ store.id }}</div>
                        <div
                            class="text-xs font-black text-indigo-600 uppercase tracking-widest group-hover:translate-x-1 transition-transform flex items-center gap-1">
                            Masuk
                            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                                <path stroke-linecap="round" stroke-linejoin="round" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                            </svg>
                        </div>
                    </div>
                </div>

                <div
                    @click="router.push('/?action=expansion#pricing')"
                    class="bg-transparent rounded-[32px] p-6 border-2 border-dashed border-slate-300 hover:border-indigo-400 hover:bg-indigo-50/50 transition-all duration-300 cursor-pointer flex flex-col items-center justify-center text-center min-h-[240px] group">
                    <div
                        class="w-14 h-14 rounded-full bg-slate-200 group-hover:bg-indigo-200 flex items-center justify-center mb-4 transition-colors">
                        <svg
                            class="w-6 h-6 text-slate-500 group-hover:text-indigo-600"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            stroke-width="3">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
                        </svg>
                    </div>
                    <h3 class="text-base font-black text-slate-600 group-hover:text-indigo-800 uppercase tracking-wider mb-1">Buka Cabang Baru</h3>
                    <p class="text-xs font-bold text-slate-400 px-4">Ekspansi bisnis Anda dengan menambah infrastruktur baru.</p>
                </div>
            </div>
        </div>

        <div v-if="isLoading" class="fixed inset-0 z-50 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center">
            <div class="bg-white p-6 rounded-2xl shadow-xl flex flex-col items-center">
                <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
                <div class="text-xs font-black text-slate-600 uppercase tracking-widest animate-pulse">Menyiapkan Akses...</div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.animate-blob {
    animation: blob 7s infinite;
}
.animation-delay-2000 {
    animation-delay: 2s;
}
@keyframes blob {
    0% {
        transform: translate(0px, 0px) scale(1);
    }
    33% {
        transform: translate(30px, -50px) scale(1.1);
    }
    66% {
        transform: translate(-20px, 20px) scale(0.9);
    }
    100% {
        transform: translate(0px, 0px) scale(1);
    }
}
</style>
