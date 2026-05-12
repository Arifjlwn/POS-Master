<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// State (Variabel penampung inputan)
const email = ref('');
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
        identifier: email.value,
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

    // 2. Logika Redirect Cerdas (Setup Toko vs Dashboard)
    if (data.has_setup_store === false) {
      alert("Login Berhasil! Silakan lengkapi data toko Anda terlebih dahulu.");
      router.push('/setup'); // Arahkan ke halaman setup toko
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
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold text-center mb-6 text-gray-800">Login POS UMKM</h2>
      
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 text-red-700 rounded text-sm text-center">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700">Email</label>
          <input 
            v-model="email" 
            type="email" 
            required 
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="admin@toko.com"
          >
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700">Password</label>
          <input 
            v-model="password" 
            type="password" 
            required 
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="••••••••"
          >
        </div>

        <button 
          type="submit" 
          :disabled="isLoading"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {{ isLoading ? 'Memproses...' : 'Masuk' }}
        </button>
      </form>
    </div>
  </div>
</template>