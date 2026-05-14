<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const router = useRouter();

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

// --- 🚀 VALIDASI PASSWORD ---
const hasUppercase = computed(() => /[A-Z]/.test(form.value.password));
const hasNumber = computed(() => /[0-9]/.test(form.value.password));
const isPasswordValid = computed(() => form.value.password.length >= 6 && hasUppercase.value && hasNumber.value);
const isMatch = computed(() => form.value.password === form.value.confirmPassword && form.value.confirmPassword !== '');

const handleRegister = async () => {
    // Validasi Akhir sebelum kirim
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
            no_hp: form.value.no_hp,
            role: 'owner'
        });

        Swal.fire({
            icon: 'success',
            title: 'Daftar Berhasil!',
            text: 'Silakan masukkan kode OTP dari email Anda.',
            confirmButtonColor: '#2563eb'
        });
        router.push({
            path: '/select-verify',
            query: { 
                email: form.value.email,
                phone: form.value.no_hp
            }
        });

    } catch (error) {
        const msg = error.response?.data?.error || 'Terjadi kesalahan pendaftaran.';
        Swal.fire({ 
          icon: 'error',
          title: 'Pendaftaran Gagal', 
          text: msg,
          confirmButtonColor: '#ef4444'
        });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] font-sans p-6 relative overflow-hidden">
        
        <div class="absolute -top-24 -right-24 w-96 h-96 bg-blue-100/50 rounded-full blur-3xl"></div>
        <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-indigo-100/50 rounded-full blur-3xl"></div>

        <div class="w-full max-w-2xl relative">
            <div class="bg-white rounded-[40px] p-8 md:p-12 shadow-2xl relative border border-white">
                
                <div class="text-center mb-10">
                    <h2 class="text-4xl font-black text-slate-900 tracking-tighter">POS<span class="text-blue-600">UMKM</span></h2>
                    <p class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.3em] mt-2">Daftarkan akun owner untuk mulai bisnis</p>
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
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg></div>
                                <input v-model="form.no_hp" type="text" required class="input-field" placeholder="08xxxxxxxxxx">
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Password</label>
                            <div class="relative group">
                                <div class="icon-container"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg></div>
                                <input v-model="form.password" :type="showPassword ? 'text' : 'password'" required class="input-field" placeholder="••••••••">
                            </div>
                            <div class="flex gap-2 mt-2 px-1">
                                <div class="flex items-center gap-1 text-[9px] font-bold" :class="form.password.length >= 6 ? 'text-green-500' : 'text-slate-300'">
                                    <span v-if="form.password.length >= 6">✓</span> 6 Karakter
                                </div>
                                <div class="flex items-center gap-1 text-[9px] font-bold" :class="hasUppercase ? 'text-green-500' : 'text-slate-300'">
                                    <span v-if="hasUppercase">✓</span> Huruf Besar
                                </div>
                                <div class="flex items-center gap-1 text-[9px] font-bold" :class="hasNumber ? 'text-green-500' : 'text-slate-300'">
                                    <span v-if="hasNumber">✓</span> Angka
                                </div>
                            </div>
                        </div>

                        <div class="space-y-1">
                            <label class="label-style">Konfirmasi Password</label>
                            <div class="relative group">
                                <div class="icon-container" :class="isMatch ? 'text-green-500' : ''">
                                    <svg v-if="!isMatch" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
                                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
                                </div>
                                <input v-model="form.confirmPassword" :type="showPassword ? 'text' : 'password'" required class="input-field" :class="isMatch ? 'border-green-100 bg-green-50/30' : ''" placeholder="Ulangi Password">
                            </div>
                            <p v-if="isMatch" class="text-[9px] text-green-600 font-bold mt-1 ml-1 uppercase tracking-widest">Password Sudah Cocok ✨</p>
                        </div>
                    </div>

                    <button type="submit" :disabled="isLoading || !isPasswordValid || !isMatch" class="btn-primary mt-6">
                        <span v-if="!isLoading" class="flex items-center gap-2 justify-center">
                            DAFTARKAN AKUN OWNER <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
                        </span>
                        <span v-else class="animate-pulse">MENYIMPAN DATA...</span>
                    </button>
                </form>

                <div class="mt-8 text-center">
                    <p class="text-xs font-bold text-slate-400 uppercase tracking-tight">
                        Sudah punya akun? 
                        <router-link to="/login" class="text-blue-600 hover:underline">Masuk Sekarang</router-link>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.label-style {
    @apply text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1;
}
.icon-container {
    @apply absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-blue-600 transition-colors;
}
.input-field {
    @apply w-full pl-12 pr-4 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-blue-600 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300;
}
.btn-primary {
    @apply w-full py-5 rounded-2xl bg-blue-600 text-white font-black text-sm uppercase tracking-widest shadow-xl shadow-blue-200 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-30;
}
</style>