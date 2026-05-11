<script setup>
import { Head, Link, useForm } from '@inertiajs/vue3';
import InputError from '@/Components/InputError.vue';

const form = useForm({
    name: '',
    email: '',
    password: '',
    password_confirmation: '',
});

const submit = () => {
    form.post(route('register'), {
        onFinish: () => form.reset('password', 'password_confirmation'),
    });
};
</script>

<template>
    <Head title="Daftar Akun" />

    <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8 font-sans">
        <div class="sm:mx-auto sm:w-full sm:max-w-md text-center">
            <div class="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg mb-4 transform -rotate-3 hover:rotate-0 transition-transform duration-300">
                <span class="text-3xl text-white font-black italic">POS</span>
            </div>
            <h2 class="text-3xl font-black text-gray-900 tracking-tight">
                Buat Akun Baru
            </h2>
            <p class="mt-2 text-sm text-gray-500 font-medium">
                Satu akun untuk kelola seluruh operasional tokomu.
            </p>
        </div>

        <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
            <div class="bg-white py-8 px-6 shadow-2xl sm:rounded-3xl sm:px-10 border border-gray-100">
                <form class="space-y-5" @submit.prevent="submit">

                    <div>
                        <label class="block text-sm font-bold text-gray-700 mb-1">Nama Lengkap</label>
                        <input v-model="form.name" type="text" required autofocus
                            class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-medium"
                            placeholder="Contoh: Arif Juliawan">
                        <InputError class="mt-1" :message="form.errors.name" />
                    </div>

                    <div>
                        <label class="block text-sm font-bold text-gray-700 mb-1">Alamat Email</label>
                        <input v-model="form.email" type="email" required
                            class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-medium"
                            placeholder="nama@toko.com">
                        <InputError class="mt-1" :message="form.errors.email" />
                    </div>

                    <div>
                        <label class="block text-sm font-bold text-gray-700 mb-1">Kata Sandi</label>
                        <input v-model="form.password" type="password" required
                            class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-medium"
                            placeholder="Minimal 8 karakter">
                        <InputError class="mt-1" :message="form.errors.password" />
                    </div>

                    <div>
                        <label class="block text-sm font-bold text-gray-700 mb-1">Konfirmasi Sandi</label>
                        <input v-model="form.password_confirmation" type="password" required
                            class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-xl shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-medium"
                            placeholder="Ulangi kata sandi">
                        <InputError class="mt-1" :message="form.errors.password_confirmation" />
                    </div>

                    <div class="pt-2">
                        <button type="submit" :disabled="form.processing"
                            class="w-full flex justify-center py-3.5 px-4 border border-transparent rounded-xl shadow-md text-sm font-black text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all disabled:opacity-50 hover:-translate-y-0.5">
                            {{ form.processing ? 'Menyiapkan Akun...' : 'Daftar Sekarang 🚀' }}
                        </button>
                    </div>

                    <div class="mt-6 text-center text-sm font-medium text-gray-600 border-t pt-4 border-gray-100">
                        Sudah punya akun?
                        <Link :href="route('login')" class="text-blue-600 hover:text-blue-800 font-bold hover:underline transition-all">
                            Masuk di sini
                        </Link>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>
