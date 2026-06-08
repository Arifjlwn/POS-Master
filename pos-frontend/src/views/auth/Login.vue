<script setup>
import Swal from 'sweetalert2';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';

const router = useRouter();

const identifier = ref('');
const password = ref('');
const showPassword = ref(false);
const errorMessage = ref('');
const isLoading = ref(false);

const goHome = () => {
    router.push('/');
};

const goToForgotPassword = () => {
    router.push('/forgot-password');
};

const handleLogin = async () => {
    isLoading.value = true;
    errorMessage.value = '';

    try {
        const response = await api.post('/login', {
            identifier: identifier.value,
            password: password.value,
        });

        const data = response.data;

        // 1. LOGIKA MULTI-OUTLET (JIKA BACKEND MINTA PILIH CABANG)
        if (data.require_select) {
            // Simpan token sementara (hanya valid untuk endpoint /select-store)
            localStorage.setItem('token', data.token);

            // Simpan data toko dan nama ke local storage sementara
            localStorage.setItem('temp_stores', JSON.stringify(data.stores || []));
            localStorage.setItem('temp_name', data.name || 'Owner');

            // Lempar ke halaman Pilih Cabang
            router.push('/select-store');
            return; // Hentikan eksekusi script di sini
        }

        // 2. LOGIKA SINGLE-OUTLET (LANGSUNG MASUK TANPA PILIH CABANG)
        // Simpan data Auth Final ke Local Storage
        localStorage.setItem('token', data.token);
        localStorage.setItem('role', data.role.toLowerCase());
        localStorage.setItem('name', data.name || '');
        localStorage.setItem('storeName', data.store_name || 'Toko Belum Di-Setup');
        localStorage.setItem('storeLogo', data.store_logo || '');
        localStorage.setItem('foto_url', data.foto_url || '');
        localStorage.setItem('qrisImage', data.qris_image || '');
        localStorage.setItem('qrisName', data.qris_name || '');
        localStorage.setItem('subscriptionPlan', data.subscription_plan || 'basic');
        localStorage.setItem('fiturAktif', data.fitur_aktif || '[]');

        if (data.business_type) {
            localStorage.setItem('businessType', data.business_type);
        }

        const Toast = Swal.mixin({
            toast: true,
            position: 'top-end',
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            customClass: { popup: 'rounded-2xl font-sans' },
        });

        Toast.fire({
            icon: 'success',
            title: `Halo, ${data.name}!`,
            text: 'Selamat datang di sistem manajemen.',
        });

        if (data.has_setup_store === false) {
            router.push('/setup-toko');
        } else {
            const roleUser = data.role.toLowerCase();
            let tipeBisnis = (data.business_type || data.tipe_bisnis || '').toLowerCase();

            if (!tipeBisnis || tipeBisnis === '') {
                if (identifier.value.toLowerCase().includes('laundry')) {
                    tipeBisnis = 'laundry';
                } else {
                    tipeBisnis = 'retail';
                }
            }

            if (tipeBisnis.includes('laundry')) {
                router.push(roleUser === 'owner' ? '/laundry/laporan' : '/laundry/pos');
            } else if (tipeBisnis.includes('bengkel')) {
                router.push(roleUser === 'owner' ? '/bengkel/dashboard' : '/bengkel/pos');
            } else if (tipeBisnis.includes('salon') || tipeBisnis.includes('barbershop')) {
                router.push(roleUser === 'owner' ? '/salon/dashboard' : '/salon/pos');
            } else if (tipeBisnis.includes('cuci')) {
                router.push(roleUser === 'owner' ? '/cuci-kendaraan/dashboard' : '/cuci-kendaraan/pos');
            } else if (tipeBisnis.includes('f&b') || tipeBisnis.includes('food')) {
                router.push(roleUser === 'owner' ? '/fnb/laporan' : '/fnb/kasir');
            } else {
                router.push('/retail/pos/riwayat');
            }
        }
    } catch (error) {
        const msg = error.response?.data?.error || error.message || 'Gagal login, silakan coba lagi.';
        errorMessage.value = msg;

        Swal.fire({
            icon: 'error',
            title: 'Akses Ditolak',
            text: msg,
            confirmButtonColor: '#2563eb',
            customClass: { popup: 'rounded-[32px]' },
        });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] font-sans p-4 relative overflow-hidden">
        <div class="absolute -top-24 -left-24 w-96 h-96 bg-blue-100/50 rounded-full blur-3xl"></div>
        <div class="absolute -bottom-24 -right-24 w-96 h-96 bg-indigo-100/50 rounded-full blur-3xl"></div>

        <div class="w-full max-w-md relative group mt-10 md:mt-0 z-10">
            <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-[40px] blur opacity-20 group-hover:opacity-30 transition duration-1000"></div>

            <div class="bg-white rounded-[40px] p-8 md:p-12 shadow-2xl relative border border-white flex flex-col">
                <div class="text-center mb-10 relative">
                    
                    <!-- 🚀 FIX: Mengganti Ikon SVG dengan File Logo Asli 512x512 -->
                    <div class="inline-flex items-center justify-center w-20 h-20 bg-white rounded-3xl shadow-xl shadow-blue-100/60 mb-6 transform -rotate-6 group-hover:rotate-0 transition-transform duration-500 border border-slate-50 overflow-hidden p-2">
                        <img src="/logo.png" alt="ARZURA POS Logo" class="w-full h-full object-contain" />
                    </div>

                    <h2 class="text-3xl font-black text-slate-900 tracking-tighter">
                        ARZURA
                        <span class="text-blue-600">POS</span>
                    </h2>
                    <p class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.3em] mt-2">Enterprise Retail Management</p>
                </div>

                <form @submit.prevent="handleLogin" class="space-y-5">
                    <div class="space-y-2">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Email Owner / No. WA Kasir</label>
                        <div class="relative group">
                            <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-600 transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
                                    <circle cx="12" cy="7" r="4" />
                                </svg>
                            </div>
                            <input v-model="identifier" type="text" required class="w-full pl-12 pr-4 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-blue-600 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300 placeholder:font-medium" placeholder="contoh@email.com / 08123456..." />
                        </div>
                    </div>

                    <div class="space-y-2">
                        <div class="flex justify-between items-center px-1">
                            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Password</label>
                            <a href="#" @click.prevent="goToForgotPassword" class="text-[10px] font-black text-blue-600 uppercase hover:underline">Lupa Password?</a>
                        </div>
                        <div class="relative group">
                            <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-600 transition-colors">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                                </svg>
                            </div>
                            <input v-model="password" :type="showPassword ? 'text' : 'password'" required class="w-full pl-12 pr-12 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-blue-600 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300" placeholder="••••••••" />
                            <button type="button" @click="showPassword = !showPassword" class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 hover:text-blue-600 transition-all">
                                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
                                    <circle cx="12" cy="12" r="3" />
                                </svg>
                                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M9.88 9.88 1.45 1.45" />
                                    <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68" />
                                    <path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
                                    <line x1="2" x2="22" y1="2" y2="22" />
                                    <path d="M13.41 13.41a2 2 0 0 1-2.82-2.82" />
                                </svg>
                            </button>
                        </div>
                    </div>

                    <button type="submit" :disabled="isLoading" class="w-full py-4 rounded-2xl bg-blue-600 text-white font-black text-sm uppercase tracking-widest shadow-xl shadow-blue-200 hover:bg-slate-900 hover:shadow-slate-200 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-3 mt-4">
                        <template v-if="!isLoading">
                            MASUK
                            <svg xmlns="http://www.w3.org/2000/xl" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M5 12h14" />
                                <path d="m12 5 7 7-7 7" />
                            </svg>
                        </template>
                        <template v-else>
                            <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            SINKRONISASI...
                        </template>
                    </button>
                </form>

                <button @click="goHome" class="mt-6 w-full py-3 rounded-xl bg-slate-50 text-slate-500 font-bold text-[10px] uppercase tracking-widest hover:bg-slate-100 hover:text-blue-600 transition-all flex items-center justify-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                    </svg>
                    Kembali ke Halaman Utama
                </button>

                <div class="mt-8 text-center">
                    <p class="text-[9px] font-black text-slate-300 uppercase tracking-[0.4em]">Integrated Business Intelligence &copy; 2026</p>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Transisi Smooth */
.transition-all {
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Animasi Error shake yang lebih halus */
.animate-shake {
    animation: shake 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
}

@keyframes shake {
    10%,
    90% {
        transform: translate3d(-1px, 0, 0);
    }
    20%,
    80% {
        transform: translate3d(2px, 0, 0);
    }
    30%,
    50%,
    70% {
        transform: translate3d(-4px, 0, 0);
    }
    40%,
    60% {
        transform: translate3d(4px, 0, 0);
    }
}
</style>