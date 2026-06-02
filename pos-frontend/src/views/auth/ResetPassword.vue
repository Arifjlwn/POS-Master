<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const route = useRoute();
const router = useRouter();

// Nangkep email dan token/otp dari halaman verifikasi sebelumnya via URL query
const email = ref('');
const token = ref('');

const password = ref('');
const confirmPassword = ref('');
const showPassword = ref(false);
const isLoading = ref(false);

onMounted(() => {
    email.value = route.query.email || '';
    token.value = route.query.token || route.query.otp || '';
    
    // Keamanan Berlapis: Kalo ga ada email atau token, tendang balik ke login
    if (!email.value || !token.value) {
        Swal.fire({
            icon: 'error',
            title: 'Sesi Tidak Valid',
            text: 'Permintaan reset password tidak sah. Silakan ulangi proses.',
            confirmButtonColor: '#4f46e5'
        });
        router.push('/login');
    }
});

// --- VALIDASI ENTITY PASSWORD (SAMA KAYA REGISTER BIAR KONSISTEN) ---
const hasUppercase = computed(() => /[A-Z]/.test(password.value));
const hasNumber = computed(() => /[0-9]/.test(password.value));
const isPasswordValid = computed(() => password.value.length >= 6 && hasUppercase.value && hasNumber.value);
const isMatch = computed(() => password.value === confirmPassword.value && confirmPassword.value !== '');

const handleResetPassword = async () => {
    if (!isPasswordValid.value || !isMatch.value) return;

    isLoading.value = true;
    try {
        // 🚀 TEMBAK API GOLANG UNTUK UPDATE DATABASE AUTOMATICALLY
        await api.post('/reset-password', {
            email: email.value,
            token: token.value, // OTP / Token verifikasi
            password: password.value // Password baru yang udah lolos enkripsi/hashing di Go nanti
        });

        await Swal.fire({
            icon: 'success',
            title: 'Password Diperbarui!',
            text: 'Password baru Anda berhasil disimpan di sistem. Silakan login kembali.',
            confirmButtonColor: '#2563eb',
            customClass: { popup: 'rounded-[32px]' }
        });

        // Bersihkan token sisa biar ga bisa di-back
        router.replace('/login');
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Gagal Mereset',
            text: error.response?.data?.error || 'Token kadaluarsa atau terjadi gangguan server.',
            confirmButtonColor: '#ef4444',
            customClass: { popup: 'rounded-[32px]' }
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

    <div class="w-full max-w-md relative group z-10">
      <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-[40px] blur opacity-20 group-hover:opacity-30 transition duration-1000"></div>
      
      <div class="bg-white rounded-[40px] p-8 md:p-12 shadow-2xl relative border border-white flex flex-col">
        
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 bg-indigo-600 rounded-3xl shadow-xl shadow-indigo-200 mb-6">
             <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
          </div>
          <h2 class="text-3xl font-black text-slate-900 tracking-tighter uppercase">Kredensial Baru</h2>
          <p class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.2em] mt-2">Buat password baru akun managemen anda</p>
        </div>

        <form @submit.prevent="handleResetPassword" class="space-y-5">
          
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Password Baru</label>
            <div class="relative group">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-indigo-600 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              </div>
              <input 
                v-model="password" 
                :type="showPassword ? 'text' : 'password'" 
                required 
                class="w-full pl-12 pr-12 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-indigo-600 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300"
                placeholder="••••••••"
              >
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 hover:text-indigo-600 transition-all"
              >
                <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/><circle cx="12" cy="12" r="3"/></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9.88 9.88 1.45 1.45"/><path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"/><path d="M6.61 6.61A13.52 13.52 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"/><line x1="2" x2="22" y1="2" y2="22"/><path d="M13.41 13.41a2 2 0 0 1-2.82-2.82"/></svg>
              </button>
            </div>
            
            <div class="flex gap-2 mt-2 px-1">
                <div class="flex items-center gap-1 text-[9px] font-bold" :class="password.length >= 6 ? 'text-emerald-500' : 'text-slate-300'">
                    <span v-if="password.length >= 6">✓</span> 6 Karakter
                </div>
                <div class="flex items-center gap-1 text-[9px] font-bold" :class="hasUppercase ? 'text-emerald-500' : 'text-slate-300'">
                    <span v-if="hasUppercase">✓</span> Huruf Besar
                </div>
                <div class="flex items-center gap-1 text-[9px] font-bold" :class="hasNumber ? 'text-emerald-500' : 'text-slate-300'">
                    <span v-if="hasNumber">✓</span> Angka
                </div>
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Konfirmasi Password Baru</label>
            <div class="relative group">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-indigo-600 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              </div>
              <input 
                v-model="confirmPassword" 
                :type="showPassword ? 'text' : 'password'" 
                required 
                class="w-full pl-12 pr-4 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-indigo-600 font-bold text-slate-800 outline-none transition-all placeholder:text-slate-300"
                placeholder="••••••••"
                :class="isMatch ? 'border-emerald-100 bg-emerald-50/20 text-emerald-900' : ''"
              >
            </div>
            <p v-if="isMatch" class="text-[9px] text-emerald-600 font-bold mt-1 ml-1 uppercase tracking-widest">Password Cocok ✨</p>
          </div>

          <button 
            type="submit" 
            :disabled="isLoading || !isPasswordValid || !isMatch"
            class="w-full py-4 rounded-2xl bg-indigo-600 text-white font-black text-sm uppercase tracking-widest shadow-xl shadow-indigo-200 hover:bg-slate-900 hover:shadow-slate-200 transition-all active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center gap-3 mt-6"
          >
            <template v-if="!isLoading">
              SIMPAN PASSWORD BARU
              <svg xmlns="http://www.w3.org/2000/xl" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
            </template>
            <template v-else>
              <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              MENYIMPAN DI DATABASE...
            </template>
          </button>
        </form>

        <button @click="router.push('/login')" class="mt-6 text-[10px] font-black text-slate-300 uppercase tracking-widest hover:text-slate-500 transition-colors">Batal & Login</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.transition-all {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>