<script setup>
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const route = useRoute();
const router = useRouter();
const email = route.query.email; // Ambil email dari URL
const otp = ref('');
const isLoading = ref(false);

const handleVerify = async () => {
    isLoading.value = true;
    try {
        await api.post('/verify-otp', {
            email: email,
            otp: otp.value
        });
        
        await Swal.fire({
            icon: 'success',
            title: 'Verifikasi Sukses!',
            text: 'Akun Anda sudah aktif, silakan login.',
            confirmButtonColor: '#2563eb'
        });
        
        router.push('/login');
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'OTP Gagal',
            text: error.response?.data?.error || 'Kode salah atau kadaluarsa.'
        });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] p-6">
        <div class="w-full max-w-md bg-white rounded-[40px] p-10 shadow-2xl border border-white text-center">
            <div class="w-16 h-16 bg-blue-50 text-blue-600 rounded-3xl flex items-center justify-center mx-auto mb-6">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/><rect width="20" height="16" x="2" y="4" rx="2"/></svg>
            </div>
            <h2 class="text-3xl font-black text-slate-900 tracking-tighter">CEK EMAIL BOS!</h2>
            <p class="text-slate-400 font-bold text-[10px] uppercase tracking-widest mt-2 mb-8">Kode verifikasi dikirim ke <span class="text-blue-600">{{ email }}</span></p>

            <form @submit.prevent="handleVerify" class="space-y-6">
                <input v-model="otp" type="text" maxlength="6" class="w-full text-center text-4xl font-black tracking-[0.5em] py-5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:border-blue-600 outline-none transition-all placeholder:text-slate-200" placeholder="000000">
                <button :disabled="otp.length < 6 || isLoading" class="w-full py-5 rounded-2xl bg-blue-600 text-white font-black uppercase tracking-widest shadow-xl shadow-blue-200 transition-all active:scale-95 disabled:opacity-50">
                    {{ isLoading ? 'VERIFIKASI...' : 'AKTIFKAN AKUN 🚀' }}
                </button>
            </form>
        </div>
    </div>
</template>