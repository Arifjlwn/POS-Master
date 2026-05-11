<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api.js';

const router = useRouter();

// State Form Murni Vue
const form = ref({
    nama_toko: '',
    tipe_bisnis: 'Retail', // Default Retail
    alamat_toko: '',
    telepon: '',
    qris_image: null,
    fitur_opsional: ['kasir'],
});

const imagePreview = ref(null);
const isLoading = ref(false);

// List Tipe Bisnis untuk dipilih
const businessTypes = ['Retail (Toko Kelontong/Minimarket)', 'F&B (Resto/Cafe)', 'Jasa (Laundry/Barbershop)', 'Lainnya'];

// Fungsi untuk menangani upload gambar dan preview
const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
        form.value.qris_image = file;
        imagePreview.value = URL.createObjectURL(file);
    }
};

const submit = async () => {
    isLoading.value = true;
    
    // Gunakan FormData karena ada file (QRIS)
    const formData = new FormData();
    formData.append('nama_toko', form.value.nama_toko);
    formData.append('tipe_bisnis', form.value.tipe_bisnis);
    formData.append('alamat_toko', form.value.alamat_toko);
    formData.append('telepon', form.value.telepon);
    formData.append('fitur_opsional', JSON.stringify(form.value.fitur_opsional));
    if (form.value.qris_image) {
        formData.append('qris_image', form.value.qris_image);
    }

    try {
        await api.post('/setup-toko', formData);
        alert('Toko Berhasil Dibuat! Selamat datang Bos Arif.');
        router.push('/dashboard'); // Langsung gas ke Dashboard
    } catch (error) {
        alert('Gagal setup toko. Cek koneksi ke server Golang.');
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen bg-gray-100 flex flex-col justify-center py-10 sm:px-6 lg:px-8 font-sans">
        <div class="sm:mx-auto sm:w-full sm:max-w-xl text-center">
            <div class="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg mb-4">
                <span class="text-3xl">🚀</span>
            </div>
            <h2 class="text-3xl font-black text-gray-900">Mulai Bisnismu</h2>
            <p class="mt-2 text-sm text-gray-600 font-medium">Sesuaikan aplikasi dengan tipe bisnismu.</p>
        </div>

        <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-xl px-4">
            <div class="bg-white py-8 px-6 shadow-xl rounded-3xl border border-gray-100">
                <form class="space-y-6" @submit.prevent="submit">

                    <div>
                        <h3 class="text-lg font-black text-gray-800 border-b border-gray-100 pb-2 mb-4">Informasi Bisnis</h3>
                        <div class="space-y-4">
                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Nama Toko/Bisnis <span class="text-red-500">*</span></label>
                                <input v-model="form.nama_toko" type="text" required class="block w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none font-bold transition-all">
                            </div>

                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Tipe Bisnis <span class="text-red-500">*</span></label>
                                <select v-model="form.tipe_bisnis" class="block w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none font-bold bg-white cursor-pointer transition-all">
                                    <option v-for="type in businessTypes" :key="type" :value="type">{{ type }}</option>
                                </select>
                            </div>

                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1">Alamat Lengkap <span class="text-red-500">*</span></label>
                                <textarea v-model="form.alamat_toko" rows="2" required class="block w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none font-medium transition-all"></textarea>
                            </div>

                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1">WhatsApp / Telepon Bisnis</label>
                                <input v-model="form.telepon" type="text" class="block w-full px-4 py-3 border border-gray-200 rounded-2xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none font-bold transition-all" placeholder="08xxxx">
                            </div>
                        </div>
                    </div>

                    <div>
                        <h3 class="text-lg font-black text-gray-800 border-b border-gray-100 pb-2 mb-4 mt-6">Pilih Modul Aktif</h3>
                        <div class="space-y-3">
                            <label class="flex items-start p-4 border-2 border-blue-100 bg-blue-50 rounded-2xl cursor-not-allowed transition-all">
                                <div class="flex items-center h-5">
                                    <input type="checkbox" value="kasir" v-model="form.fitur_opsional" disabled class="w-5 h-5 text-blue-600 rounded-lg border-gray-300">
                                </div>
                                <div class="ml-3">
                                    <span class="font-black text-blue-900 block text-sm">Sistem POS Kasir</span>
                                    <span class="text-blue-700 text-[11px] font-bold uppercase tracking-tight">Modul Wajib Terpasang</span>
                                </div>
                            </label>

                            <label class="flex items-start p-4 border border-gray-200 rounded-2xl cursor-pointer hover:bg-gray-50 transition-all">
                                <div class="flex items-center h-5">
                                    <input type="checkbox" value="absensi" v-model="form.fitur_opsional" class="w-5 h-5 text-blue-600 rounded-lg border-gray-300">
                                </div>
                                <div class="ml-3">
                                    <span class="font-black text-gray-800 block text-sm">Absensi Karyawan</span>
                                    <span class="text-gray-500 text-[11px] font-medium">Rekam jam masuk & pulang karyawan secara digital.</span>
                                </div>
                            </label>
                        </div>
                    </div>

                    <div class="pt-4">
                        <button type="submit" :disabled="isLoading" class="w-full flex justify-center py-4 px-4 rounded-2xl shadow-lg text-base font-black text-white bg-blue-700 hover:bg-blue-800 focus:outline-none transition-all disabled:opacity-50 transform hover:-translate-y-1 active:scale-95">
                            {{ isLoading ? 'Mempersiapkan Sistem...' : 'Buka Toko Sekarang 🚀' }}
                        </button>
                    </div>

                </form>
            </div>
        </div>
    </div>
</template>