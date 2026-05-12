<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// State (Ganti nama variabel biar lebih general)
const identifier = ref(''); 
const password = ref('');
const errorMessage = ref('');
const isLoading = ref(false);

// Fungsi untuk hit API Login
const handleLogin = async () => {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        identifier: identifier.value, // 🚀 Dikirim ke Golang sebagai identifier
        password: password.value
      })
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Gagal login, silakan coba lagi.');
    }

    // 1. Simpan Karcis VIP (Token) ke brankas browser (Local Storage)
    localStorage.setItem('token', data.token);
    localStorage.setItem('role', data.role);
    localStorage.setItem('name', data.name);

    // 2. Logika Redirect Cerdas (Setup Toko vs Dashboard)
    if (data.has_setup_store === false) {
      alert("Login Berhasil! Silakan lengkapi data toko Anda terlebih dahulu.");
      router.push('/setup-toko'); // Arahkan ke halaman setup toko
    } else {
      router.push('/dashboard'); // Arahkan ke dashboard
    }

  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md border-t-8 border-blue-600">
      <div class="text-center mb-8">
        <h2 class="text-3xl font-black text-gray-900">POS UMKM</h2>
        <p class="text-gray-500 font-medium mt-1">Silakan masuk ke akun Anda</p>
      </div>
      
      <div v-if="errorMessage" class="mb-5 p-4 bg-red-50 border border-red-200 text-red-600 rounded-xl text-sm font-bold text-center">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-5">
        <div>
          <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Email Owner / NIK Kasir</label>
          <input 
            v-model="identifier" 
            type="text" 
            required 
            class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-bold text-gray-800 outline-none transition-all"
            placeholder="admin@toko.com atau 20260001"
          >
        </div>

        <div>
          <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Password</label>
          <input 
            v-model="password" 
            type="password" 
            required 
            class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-bold text-gray-800 outline-none transition-all"
            placeholder="••••••••"
          >
        </div>

        <button 
          type="submit" 
          :disabled="isLoading"
          class="w-full flex justify-center py-3.5 px-4 border border-transparent rounded-xl shadow-lg shadow-blue-100 text-sm font-black text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 transition-all active:scale-95 mt-2"
        >
          {{ isLoading ? 'Memeriksa Data...' : 'Masuk ke Sistem 🚀' }}
        </button>
      </form>
    </div>
  </div>
</template>