<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import Swal from 'sweetalert2'; // 🚀 Import SweetAlert2
import api from '../api.js';

const router = useRouter();

const identifier = ref(''); 
const password = ref('');
const showPassword = ref(false);
const errorMessage = ref('');
const isLoading = ref(false);

const handleLogin = async () => {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    // 🚀 1. Tembak API menggunakan instance Axios kita yang sudah nge-bind .env
    const response = await api.post('/login', {
      identifier: identifier.value,
      password: password.value
    });

    // 🚀 2. Di Axios, datanya langsung nangkring di response.data tanpa .json()
    const data = response.data; 
    console.log("Data dari API:", data);

    // 3. Simpan ke Local Storage
    localStorage.setItem('token', data.token);
    localStorage.setItem('role', data.role.toLowerCase()); // Simpan kecil semua biar aman dicek
    localStorage.setItem('name', data.name || '');
    localStorage.setItem('storeName', data.store_name || 'Toko Belum Di-Setup');

    // 🚀 4. SWEETALERT LOGIN BERHASIL
    const Toast = Swal.mixin({
      toast: true,
      position: 'top-end',
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
    });

    Toast.fire({
      icon: 'success',
      title: `Halo, ${data.name}!`,
      text: 'Berhasil masuk ke sistem.'
    });

    // 🚀 5. LOGIKA REDIRECT CERDAS BERDASARKAN ROLE
    if (data.has_setup_store === false) {
      Swal.fire({
        icon: 'info',
        title: 'Setup Toko',
        text: 'Lengkapi data toko Anda untuk memulai.',
        confirmButtonColor: '#2563eb'
      });
      router.push('/setup-toko');
    } else {
      if (data.role.toLowerCase() === 'owner') {
        router.push('/dashboard');
      } else {
        router.push('/absensi'); 
      }
    }

  } catch (error) {
    // 🚀 6. TANGKAP ERROR KASUS AXIOS (Membaca response error dari Go Backend Mas)
    const msg = error.response?.data?.error || error.message || 'Gagal login, silakan coba lagi.';
    errorMessage.value = msg;

    Swal.fire({
      icon: 'error',
      title: 'Login Gagal',
      text: msg,
      confirmButtonColor: '#ef4444'
    });
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100 font-sans">
    <div class="bg-white p-8 rounded-3xl shadow-2xl w-full max-w-md border-t-[12px] border-blue-600 relative overflow-hidden">
      <div class="absolute -right-10 -top-10 w-32 h-32 bg-blue-50 rounded-full"></div>
      
      <div class="text-center mb-8 relative">
        <h2 class="text-4xl font-black text-gray-900 tracking-tighter">POS<span class="text-blue-600">UMKM</span></h2>
        <p class="text-gray-400 font-bold text-sm mt-1 uppercase tracking-widest">Sistem Manajemen Retail</p>
      </div>
      
      <div v-if="errorMessage" class="mb-5 p-4 bg-red-50 border-l-4 border-red-500 text-red-700 rounded-xl text-xs font-black uppercase transition-all animate-shake">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6 relative">
        <div>
          <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 ml-1">Identitas Masuk</label>
          <div class="relative">
            <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">👤</span>
            <input 
              v-model="identifier" 
              type="text" 
              required 
              class="w-full pl-11 pr-4 py-4 bg-gray-50 border-2 border-gray-100 rounded-2xl focus:ring-4 focus:ring-blue-500/10 focus:border-blue-500 font-bold text-gray-800 outline-none transition-all placeholder:text-gray-300"
              placeholder="Email / NIK Karyawan"
            >
          </div>
        </div>

        <div>
          <label class="block text-[10px] font-black text-gray-400 uppercase tracking-[0.2em] mb-2 ml-1">Kata Sandi</label>
          <div class="relative">
            <span class="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">🔒</span>
            <input 
              v-model="password" 
              :type="showPassword ? 'text' : 'password'" 
              required 
              class="w-full pl-11 pr-4 py-4 bg-gray-50 border-2 border-gray-100 rounded-2xl focus:ring-4 focus:ring-blue-500/10 focus:border-blue-500 font-bold text-gray-800 outline-none transition-all placeholder:text-gray-300"
              placeholder="••••••••"
            >
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-blue-600 transition-colors focus:outline-none"
            >
              <span v-if="showPassword">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </span>
      <span v-else>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.542-7a10.05 10.05 0 014.13-4.13m4.13-4.13A10.05 10.05 0 0112 5c4.478 0 8.268 2.943 9.542 7a10.05 10.05 0 01-4.13 4.13m-4.13-4.13L3 3m3.24 3.24l10.52 10.52" />
        </svg>
      </span>
            </button>
          </div>
        </div>

        <button 
          type="submit" 
          :disabled="isLoading"
          class="w-full flex justify-center py-4 px-4 border border-transparent rounded-2xl shadow-xl shadow-blue-200 text-sm font-black text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 transition-all active:scale-95 mt-4"
        >
          <span v-if="!isLoading" class="flex items-center gap-2">MASUK KE SISTEM 🚀</span>
          <span v-else class="flex items-center gap-2">
            <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            MEMERIKSA...
          </span>
        </button>
      </form>
      
      <div class="mt-8 text-center text-[10px] font-black text-gray-300 uppercase tracking-widest">
        &copy; 2026 - ARIF JULIAWAN
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}
.animate-shake {
  animation: shake 0.2s ease-in-out 0s 2;
}
</style>