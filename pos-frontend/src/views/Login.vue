<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router'; // <-- Router murni Vue
import api from '../api'; // <-- Panggil kurir

const identifier = ref('');
const password = ref('');
const errorMessage = ref('');
const isLoading = ref(false);
const router = useRouter();

const handleLogin = async () => {
    isLoading.value = true;
    errorMessage.value = '';

    try {
        const response = await api.post('/login', {
            identifier: identifier.value,
            password: password.value
        });
        
        // Simpan tiket di brankas
        localStorage.setItem('token', response.data.token);
        
        // Meluncur ke Dashboard tanpa kedip!
        router.push('/dashboard');
    } catch (error) {
        if (error.response && error.response.status === 401) {
            errorMessage.value = 'Email/NIK atau Password salah!';
        } else {
            errorMessage.value = 'Gagal terhubung ke server.';
        }
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-100">
        <div class="bg-white p-8 rounded-xl shadow-lg w-full max-w-md">
            <div class="text-center mb-8">
                <h1 class="text-3xl font-black text-gray-800">
                    <span class="text-red-600 italic">Indo</span><span class="text-blue-600 italic">UMKM</span>
                </h1>
                <p class="text-gray-500 mt-2 font-medium">Sistem POS Enterprise</p>
            </div>

            <form @submit.prevent="handleLogin" class="space-y-6">
                <div v-if="errorMessage" class="bg-red-50 text-red-600 p-3 rounded-lg text-sm font-bold text-center border border-red-200">
                    {{ errorMessage }}
                </div>
                <div>
                    <label class="block text-sm font-bold text-gray-700 mb-1">Email / NIK Kasir</label>
                    <input type="text" v-model="identifier" required class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
                </div>
                <div>
                    <label class="block text-sm font-bold text-gray-700 mb-1">Password</label>
                    <input type="password" v-model="password" required class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all">
                </div>
                <button type="submit" :disabled="isLoading" class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded-lg transition-all shadow-md disabled:opacity-50">
                    {{ isLoading ? 'Memproses...' : 'MASUK' }}
                </button>
            </form>
        </div>
    </div>
</template>