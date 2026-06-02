<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const router = useRouter();
const route = useRoute(); // 🚀 TAMBAHIN INI BUAT NANGKEP URL PARAMETER

// State Form
const form = ref({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
    tempat_lahir: '',
    tanggal_lahir: '',
    no_hp: ''
});

const showPassword = ref(false);
const isLoading = ref(false);

// 🚀 STATE BUAT NANGKEP TITIPAN PAKET DARI LANDING PAGE
const selectedIndustry = ref('');
const selectedPlan = ref('');

onMounted(() => {
    // Tangkap dari URL dan simpan ke localStorage
    if (route.query.industry) {
        selectedIndustry.value = route.query.industry;
        localStorage.setItem('pendingIndustry', route.query.industry);
    }
    if (route.query.plan) {
        selectedPlan.value = route.query.plan;
        localStorage.setItem('pendingPlan', route.query.plan);
    }
});

// 🚀 FUNGSI KEMBALI KE LANDING PAGE
const goHome = () => {
    router.push('/');
};

// --- VALIDASI NO HP ---
const formatNoHpRegister = () => {
    let val = String(form.value.no_hp);
    if (val.startsWith('0')) {
        val = val.substring(1);
    }
    if (val.startsWith('62')) {
        val = val.substring(2);
    }
    form.value.no_hp = val;
};

// --- VALIDASI PASSWORD ---
const hasUppercase = computed(() => /[A-Z]/.test(form.value.password));
const hasNumber = computed(() => /[0-9]/.test(form.value.password));
const isPasswordValid = computed(() => form.value.password.length >= 6 && hasUppercase.value && hasNumber.value);
const isMatch = computed(() => form.value.password === form.value.confirmPassword && form.value.confirmPassword !== '');

const handleRegister = async () => {
    if (!isPasswordValid.value) {
        return Swal.fire({ icon: 'error', title: 'Password Lemah', text: 'Password wajib minimal 6 karakter, mengandung huruf besar dan angka.' });
    }
    if (!isMatch.value) {
        return Swal.fire({ icon: 'error', title: 'Password Tidak Cocok', text: 'Konfirmasi password harus sama.' });
    }

    isLoading.value = true;
    try {
        const response = await api.post('/register', {
            name: form.value.name,
            email: form.value.email,
            password: form.value.password,
            tempat_lahir: form.value.tempat_lahir,
            tanggal_lahir: form.value.tanggal_lahir,
            no_hp: `62${form.value.no_hp}`,
            role: 'owner'
        });

        Swal.fire({ 
            icon: 'success', 
            title: 'Pendaftaran Berhasil!', 
            text: 'Akun Anda telah tercatat. Silakan pilih metode untuk menerima kode verifikasi OTP.', 
            confirmButtonColor: '#4f46e5' 
        });

        router.push({ 
            path: '/select-verify', 
            query: { email: form.value.email, phone: form.value.no_hp } 
        });
    } catch (error) {
        const msg = error.response?.data?.error || 'Terjadi kesalahan pendaftaran.';
        Swal.fire({ icon: 'error', title: 'Pendaftaran Gagal', text: msg, confirmButtonColor: '#ef4444' });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] font-sans p-6 relative overflow-hidden">
        
        <div class="absolute -top-24 -right-24 w-96 h-96 bg-blue-100/50 rounded-full blur-3xl"></div>
        <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-indigo-100/50 rounded-full blur-3xl"></div>
        
        <!-- 🚀 TAMBAHAN: Tombol Floating Back ke Beranda -->
        <button 
            @click="goHome"
            class="absolute top-6 left-6 md:top-10 md:left-10 flex items-center gap-2 px-5 py-3 bg-white/60 hover:bg-white backdrop-blur-md border border-slate-200/50 text-slate-500 hover:text-indigo-600 rounded-2xl font-black text-[10px] uppercase tracking-widest transition-all hover:shadow-lg hover:-translate-x-1 group z-50"
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 transition-transform group-hover:-translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Kembali ke Beranda
        </button>

        <div class="w-full max-w-2xl relative mt-12 md:mt-0">
            <div class="bg-white rounded-[40px] p-8 md:p-12 shadow-2xl relative border border-white">
                
                <div class="text-center mb-10">
                    <h2 class="text-4xl font-black text-slate-900 tracking-tighter">NEXA<span class="text-indigo-600">POS</span></h2>
                    <p class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.3em] mt-2">Daftarkan akun owner untuk mulai bisnis</p>

                    <div v-if="selectedPlan" class="mt-5 inline-flex items-center gap-2 bg-indigo-50 border border-indigo-100 px-4 py-2 rounded-full shadow-sm animate-[fadeIn_0.5s_ease-out]">
                        <span class="w-2 h-2 rounded-full bg-indigo-500 animate-pulse"></span>
                        <span class="text-[9px] font-black text-indigo-700 uppercase tracking-widest">
                            Paket Terpilih: {{ selectedIndustry }} - {{ selectedPlan }}
                        </span>
                    </div>
                </div>

                <form @submit.prevent="handleRegister" class="space-y-5">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
                        
                        <div class="space-y-1">
                            <label class="label-style">Nama Lengkap</label>
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg></div>
                                <input v-model="form.name" type="text" required class="input-field" placeholder="Nama Owner">
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Alamat Email</label>
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg></div>
                                <input v-model="form.email" type="email" required class="input-field" placeholder="email@bisnis.com">
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Tempat Lahir</label>
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg></div>
                                <input v-model="form.tempat_lahir" type="text" required class="input-field" placeholder="Kota Kelahiran">
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Tanggal Lahir</label>
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg></div>
                                <input v-model="form.tanggal_lahir" type="date" required class="input-field">
                            </div>
                        </div>

                        <div class="space-y-1 md:col-span-2">
                            <label class="label-style">Nomor WhatsApp / HP</label>
                            <div class="flex items-center bg-white border-2 border-slate-200 rounded-2xl focus-within:border-indigo-500 focus-within:ring-4 focus-within:ring-indigo-500/10 transition-all shadow-sm overflow-hidden group">
                                <div class="pl-4 pr-3 py-4 bg-slate-50 border-r border-slate-200 flex items-center justify-center gap-2 select-none shrink-0">
                                    <span class="text-slate-500 font-black text-sm leading-none">+62</span>
                                </div>
                                <input v-model="form.no_hp" @input="formatNoHpRegister" type="number" required class="w-full px-4 py-4 bg-transparent outline-none font-black text-slate-800 placeholder:text-slate-300 placeholder:font-bold border-none ring-0 focus:ring-0" placeholder="81234567890">
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Password</label>
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg></div>
                                <input v-model="form.password" :type="showPassword ? 'text' : 'password'" required class="input-field" placeholder="••••••••">
                            </div>
                            <div class="flex gap-2 mt-2 px-1">
                                <div class="flex items-center gap-1 text-[9px] font-bold" :class="form.password.length >= 6 ? 'text-emerald-500' : 'text-slate-300'">
                                    <span v-if="form.password.length >= 6">✓</span> 6 Karakter
                                </div>
                                <div class="flex items-center gap-1 text-[9px] font-bold" :class="hasUppercase ? 'text-emerald-500' : 'text-slate-300'">
                                    <span v-if="hasUppercase">✓</span> Huruf Besar
                                </div>
                                <div class="flex items-center gap-1 text-[9px] font-bold" :class="hasNumber ? 'text-emerald-500' : 'text-slate-300'">
                                    <span v-if="hasNumber">✓</span> Angka
                                </div>
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Konfirmasi Password</label>
                            <div class="relative group">
                                <div class="icon-container" :class="isMatch ? 'text-emerald-500' : ''">
                                    <svg v-if="!isMatch" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
                                </div>
                                <input v-model="form.confirmPassword" :type="showPassword ? 'text' : 'password'" required class="input-field" :class="isMatch ? 'border-emerald-100 bg-emerald-50/30 text-emerald-900' : ''" placeholder="Ulangi Password">
                            </div>
                            <p v-if="isMatch" class="text-[9px] text-emerald-600 font-bold mt-1 ml-1 uppercase tracking-widest animate-[fadeIn_0.3s_ease-out]">Password Cocok ✨</p>
                        </div>
                    </div>

                    <button type="submit" :disabled="isLoading || !isPasswordValid || !isMatch" class="btn-primary mt-6">
                        <span v-if="!isLoading" class="flex items-center gap-2 justify-center">
                            DAFTARKAN AKUN OWNER
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                        </span>
                        <span v-else class="animate-pulse flex items-center justify-center gap-2">
                            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            MENYIMPAN DATA...
                        </span>
                    </button>
                </form>

                <div class="mt-8 text-center">
                    <p class="text-xs font-bold text-slate-400 uppercase tracking-tight">
                        Sudah punya akun? <router-link to="/login" class="text-indigo-600 hover:underline">Masuk Sekarang</router-link>
                    </p>
                </div>

            </div>
        </div>
    </div>
</template>

<style scoped>
.label-style { @apply text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1 block; }
.icon-container { @apply absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-indigo-600 transition-colors; }
.input-field { @apply w-full pl-12 pr-4 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300; }
.btn-primary { @apply w-full py-5 rounded-[20px] bg-indigo-600 text-white font-black text-sm uppercase tracking-widest shadow-xl shadow-indigo-200 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed; }

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-5px); }
    to { opacity: 1; transform: translateY(0); }
}
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
    -webkit-appearance: none; 
    margin: 0; 
}
input[type=number] { -moz-appearance: textfield; }
</style>