<template>
    <div class="login-container">
        <h2>Login POS UMKM</h2>

        <form @submit.prevent="handleLogin">
            <div class="form-group">
                <label>Email / NIK Kasir</label>
                <input
                    type="text"
                    v-model="identifier"
                    placeholder="Masukkan Email atau NIK"
                    required
                />
            </div>

            <div class="form-group">
                <label>Password</label>
                <input
                    type="password"
                    v-model="password"
                    placeholder="Masukkan Password"
                    required
                />
            </div>

            <button type="submit" :disabled="isLoading">
                {{ isLoading ? 'Memeriksa...' : 'Masuk' }}
            </button>

            <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
        </form>
    </div>
</template>

<script setup>
<<<<<<< HEAD
import { ref } from 'vue';
import api from '../../api.js'; // Panggil kurir Axios (Posisinya 2 folder di luar)
=======
import { Head, Link, useForm } from '@inertiajs/vue3';
import InputError from '@/Components/InputError.vue';
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970

// State (Wadah penampung inputan)
const identifier = ref('');
const password = ref('');
const errorMessage = ref('');
const isLoading = ref(false);

// Fungsi saat tombol Masuk diklik
const handleLogin = async () => {
    isLoading.value = true;
    errorMessage.value = '';

    try {
        // Tembak API Golang
        const response = await api.post('/login', {
            identifier: identifier.value,
            password: password.value
        });

        // Kalau sukses, simpan tiket JWT ke brankas browser
        const token = response.data.token;
        localStorage.setItem('token', token);

        // Pakai cara bawaan browser/Inertia untuk pindah halaman
        window.location.href = '/dashboard';

    } catch (error) {
        // Tangkap pesan error dari Golang
        if (error.response && error.response.data.error) {
            errorMessage.value = error.response.data.error;
        } else {
            errorMessage.value = "Gagal terhubung ke server.";
        }
    } finally {
        isLoading.value = false;
    }
};
</script>

<<<<<<< HEAD
<style scoped>
.login-container {
    max-width: 400px;
    margin: 100px auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 8px;
    font-family: sans-serif;
}
.form-group {
    margin-bottom: 15px;
}
.form-group label {
    display: block;
    margin-bottom: 5px;
}
.form-group input {
    width: 100%;
    padding: 8px;
    box-sizing: border-box;
}
button {
    width: 100%;
    padding: 10px;
    background-color: #007BFF;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}
button:disabled {
    background-color: #aaa;
}
.error-text {
    color: red;
    margin-top: 10px;
    text-align: center;
}
</style>
=======
<template>
    <Head title="Masuk" />

    <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8 font-sans">

        <div class="sm:mx-auto sm:w-full sm:max-w-md text-center">
            <div class="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg mb-4 transform -rotate-3 hover:rotate-0 transition-transform duration-300">
                <span class="text-3xl text-white font-black italic">POS</span>
            </div>
            <h2 class="text-3xl font-black text-gray-900 tracking-tight">
                Selamat Datang
            </h2>
            <p class="mt-2 text-sm text-gray-500 font-medium">
                Masuk untuk mengelola operasional tokomu.
            </p>
        </div>

        <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
            <div class="bg-white py-8 px-6 shadow-2xl sm:rounded-3xl sm:px-10 border border-gray-100">

                <div v-if="status" class="mb-5 p-4 rounded-xl bg-green-50 border border-green-200 text-sm font-bold text-green-700 text-center shadow-sm">
                    {{ status }}
                </div>

                <form class="space-y-5" @submit.prevent="submit">

                    <div>
                        <label class="block text-sm font-bold text-gray-700 mb-1">Email Owner / NIK Karyawan</label>
                        <input v-model="form.email" type="text" required autofocus autocomplete="username"
                            class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-medium"
                            placeholder="nama@toko.com atau 20260001">
                        <InputError class="mt-1" :message="form.errors.email" />
                    </div>

                    <div>
                        <label class="block text-sm font-bold text-gray-700 mb-1">Kata Sandi</label>
                        <input v-model="form.password" type="password" required autocomplete="current-password"
                            class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-medium"
                            placeholder="Masukkan kata sandi">
                        <InputError class="mt-1" :message="form.errors.password" />
                    </div>

                    <div class="flex items-center justify-between pt-1">
                        <div class="flex items-center">
                            <input v-model="form.remember" id="remember" type="checkbox"
                                class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded cursor-pointer">
                            <label for="remember" class="ml-2 block text-sm text-gray-700 font-bold cursor-pointer">
                                Ingat saya
                            </label>
                        </div>

                        <div class="text-sm" v-if="canResetPassword">
                            <Link :href="route('password.request')" class="font-bold text-blue-600 hover:text-blue-800 hover:underline transition-all">
                                Lupa sandi?
                            </Link>
                        </div>
                    </div>

                    <div class="pt-2">
                        <button type="submit" :disabled="form.processing"
                            class="w-full flex justify-center py-3.5 px-4 border border-transparent rounded-xl shadow-md text-sm font-black text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all disabled:opacity-50 hover:-translate-y-0.5">
                            {{ form.processing ? 'Memeriksa Data...' : 'Masuk 🚀' }}
                        </button>
                    </div>

                    <div class="mt-6 text-center text-sm font-medium text-gray-600 border-t pt-5 border-gray-100">
                        Belum punya akun toko?
                        <Link :href="route('register')" class="text-blue-600 hover:text-blue-800 font-bold hover:underline transition-all">
                            Daftar di sini
                        </Link>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>
>>>>>>> fe5108ac44404475b1c66a1cb99fadaf9190e970
