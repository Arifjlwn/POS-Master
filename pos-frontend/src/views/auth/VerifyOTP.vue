<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const route = useRoute();
const router = useRouter();

// Ambil data dari query parameter halaman sebelumnya
const email = route.query.email || '';
const phone = route.query.phone || '';

const otp = ref('');
const isLoading = ref(false);

// --- LOGIKA TIMER & RESEND OTP ---
const INITIAL_TIME = 180; // 3 Menit (180 Detik)
const timeLeft = ref(INITIAL_TIME);
let timerInterval = null;

const startTimer = () => {
    clearInterval(timerInterval);
    timeLeft.value = INITIAL_TIME;
    timerInterval = setInterval(() => {
        if (timeLeft.value > 0) {
            timeLeft.value--;
        } else {
            clearInterval(timerInterval);
        }
    }, 1000);
};

// Format tampilan detik ke MM:SS (Contoh: 02:45)
const formattedTime = computed(() => {
    const minutes = Math.floor(timeLeft.value / 60);
    const seconds = timeLeft.value % 60;
    return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
});

const isTimerActive = computed(() => timeLeft.value > 0);

// --- HIT API VERIFIKASI ---
const handleVerify = async () => {
    if (otp.value.length < 6) return;
    
    isLoading.value = true;
    try {
        await api.post('/verify-otp', {
            email: email,
            otp: otp.value
        });
        
        await Swal.fire({
            icon: 'success',
            title: 'Verifikasi Sukses!',
            text: 'Akun Owner Anda sudah aktif, mari konfigurasi toko Anda.',
            confirmButtonColor: '#2563eb'
        });
        
        router.push('/login');
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Verifikasi Gagal',
            text: error.response?.data?.error || 'Kode OTP salah atau sudah kadaluarsa.',
            confirmButtonColor: '#ef4444'
        });
    } finally {
        isLoading.value = false;
    }
};

// --- HIT API RESEND OTP ---
const handleResendOTP = async () => {
    if (isTimerActive.value) return;

    isLoading.value = true;
    try {
        // Panggil endpoint register ulang atau endpoint khusus resend jika ada
        await api.post('/register', {
            // Karena backend menggunakan upsert/create, hit register dengan email yang sama akan trigger OTP baru
            email: email,
            resend: true // flag tambahan opsional jika backend butuh
        });

        Swal.fire({
            icon: 'success',
            title: 'Kode Baru Dikirim!',
            text: 'Silakan cek kembali kotak masuk email atau folder spam Anda.',
            confirmButtonColor: '#2563eb'
        });

        otp.value = ''; // Reset inputan kotak
        startTimer();   // Jalankan ulang timer mundur
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Gagal Kirim Ulang',
            text: error.response?.data?.error || 'Terjadi kesalahan sistem, coba lagi nanti.',
            confirmButtonColor: '#ef4444'
        });
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => {
    startTimer();
});

onUnmounted(() => {
    clearInterval(timerInterval);
});
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-[#F8FAFC] p-6 relative overflow-hidden font-sans">
        
        <div class="absolute -top-24 -right-24 w-96 h-96 bg-blue-100/40 rounded-full blur-3xl opacity-70"></div>
        <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-indigo-100/40 rounded-full blur-3xl opacity-70"></div>

        <div class="w-full max-w-md relative z-10">
            <div class="bg-white rounded-[40px] p-10 shadow-2xl border border-white text-center">
                
                <div class="w-20 h-20 bg-indigo-50 text-indigo-600 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-sm border border-indigo-100/30">
                    <svg v-if="email && !phone" xmlns="http://www.w3.org/2000/svg" class="w-9 h-9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-9 h-9" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                </div>

                <h2 class="text-3xl font-black text-slate-900 tracking-tighter uppercase">Masukkan OTP</h2>
                
                <p class="text-slate-400 font-bold text-[10px] uppercase tracking-widest mt-2 mb-8 leading-relaxed">
                    Kode rahasia 6-digit telah dikirim ke <br/>
                    <span class="text-indigo-600 font-black text-xs normal-case tracking-normal">{{ email || phone }}</span>
                </p>

                <form @submit.prevent="handleVerify" class="space-y-6">
                    <div class="relative">
                        <input 
                            v-model="otp" 
                            type="text" 
                            maxlength="6" 
                            pattern="[0-9]*" 
                            inputmode="numeric"
                            class="w-full text-center text-4xl font-black tracking-[0.4em] pl-[0.4em] py-5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-600 outline-none transition-all placeholder:text-slate-200" 
                            placeholder="000000"
                            required
                        >
                    </div>

                    <div class="flex items-center justify-center gap-2 px-4 py-3 bg-slate-50 border border-slate-100 rounded-xl text-xs">
                        <span class="font-bold text-slate-400 uppercase tracking-wider">Sesi Berlaku:</span>
                        <span :class="['font-black tracking-wider', timeLeft < 30 ? 'text-red-500 animate-pulse' : 'text-indigo-600']">
                            {{ formattedTime }}
                        </span>
                    </div>

                    <button 
                        type="submit"
                        :disabled="otp.length < 6 || isLoading || !isTimerActive" 
                        class="w-full py-5 rounded-2xl bg-indigo-600 text-white font-black text-xs uppercase tracking-widest shadow-xl shadow-indigo-200 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-30 disabled:pointer-events-none"
                    >
                        {{ isLoading ? 'Memproses...' : 'Aktifkan Akun Owner 🚀' }}
                    </button>
                </form>

                <div class="mt-8 pt-6 border-t border-slate-100">
                    <p class="text-xs font-bold text-slate-400 uppercase tracking-tight">
                        Tidak menerima kode OTP? <br class="sm:hidden"/>
                        <button 
                            @click="handleResendOTP" 
                            :disabled="isTimerActive || isLoading"
                            :class="['ml-1 font-black transition-colors outline-none focus:underline', isTimerActive ? 'text-slate-300 cursor-not-allowed' : 'text-indigo-600 hover:text-slate-900']"
                        >
                            Kirim Ulang Kode
                        </button>
                    </p>
                </div>

            </div>
        </div>
    </div>
</template>

<style scoped>
/* Menghilangkan panah spinner bawaan input type number di beberapa browser */
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
input[type=number] {
    -moz-appearance: textfield;
}
</style>