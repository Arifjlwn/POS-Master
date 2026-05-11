<script setup>
import { Head, useForm } from '@inertiajs/vue3';
import { ref } from 'vue';

const form = useForm({
    nama_toko: '',
    alamat_toko: '',
    telepon: '',
    qris_image: null,
    // Default: Fitur POS Kasir wajib nyala, Absensi mati
    fitur_opsional: ['kasir'],
});

const imagePreview = ref(null);

// Fungsi untuk menangani upload gambar dan preview
const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
        form.qris_image = file;
        imagePreview.value = URL.createObjectURL(file);
    }
};

const submit = () => {
    // Karena ada file gambar, Inertia otomatis mengubah ini jadi FormData
    form.post('/setup-toko');
};
</script>

<template>
    <Head title="Setup Toko" />

    <div class="min-h-screen bg-gray-100 flex flex-col justify-center py-10 sm:px-6 lg:px-8 font-sans">
        <div class="sm:mx-auto sm:w-full sm:max-w-xl text-center">
            <div class="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg mb-4">
                <span class="text-3xl">🚀</span>
            </div>
            <h2 class="text-3xl font-black text-gray-900">Mulai Bisnismu</h2>
            <p class="mt-2 text-sm text-gray-600">Sesuaikan aplikasi dengan kebutuhan operasional tokomu.</p>
        </div>

        <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-xl">
            <div class="bg-white py-8 px-6 shadow-xl sm:rounded-2xl sm:px-10 border border-gray-100">
                <form class="space-y-6" @submit.prevent="submit">

                    <div>
                        <h3 class="text-lg font-bold text-gray-800 border-b pb-2 mb-4">Informasi Dasar</h3>
                        <div class="space-y-4">
                            <div>
                                <label class="block text-sm font-bold text-gray-700 mb-1">Nama Toko <span class="text-red-500">*</span></label>
                                <input v-model="form.nama_toko" type="text" required class="block w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-blue-500 focus:border-blue-500 text-sm">
                            </div>
                            <div>
                                <label class="block text-sm font-bold text-gray-700 mb-1">Alamat Lengkap <span class="text-red-500">*</span></label>
                                <textarea v-model="form.alamat_toko" rows="2" required class="block w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-blue-500 focus:border-blue-500 text-sm"></textarea>
                            </div>
                            <div>
                                <label class="block text-sm font-bold text-gray-700 mb-1">WhatsApp / Telepon</label>
                                <input v-model="form.telepon" type="text" class="block w-full px-4 py-2.5 border border-gray-300 rounded-xl focus:ring-blue-500 focus:border-blue-500 text-sm">
                            </div>
                        </div>
                    </div>

                    <div>
                        <h3 class="text-lg font-bold text-gray-800 border-b pb-2 mb-4 mt-6">Metode Pembayaran (QRIS)</h3>
                        <div class="mt-1 flex items-center gap-4">
                            <div v-if="imagePreview" class="h-24 w-24 rounded-xl border border-gray-200 overflow-hidden shrink-0">
                                <img :src="imagePreview" class="h-full w-full object-cover">
                            </div>
                            <div class="flex-1">
                                <label class="block text-sm font-bold text-gray-700 mb-1">Upload QR Code Toko (Opsional)</label>
                                <input type="file" @change="handleImageChange" accept="image/*" class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 cursor-pointer">
                                <p class="text-xs text-gray-500 mt-1">Format: JPG, PNG maksimal 2MB. Akan dicetak di struk.</p>
                            </div>
                        </div>
                    </div>

                    <div>
                        <h3 class="text-lg font-bold text-gray-800 border-b pb-2 mb-4 mt-6">Modul Aplikasi Aktif</h3>
                        <div class="space-y-3">
                            <label class="flex items-start p-3 border border-blue-200 bg-blue-50 rounded-xl cursor-not-allowed">
                                <div class="flex items-center h-5">
                                    <input type="checkbox" value="kasir" v-model="form.fitur_opsional" disabled class="w-5 h-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500">
                                </div>
                                <div class="ml-3 text-sm">
                                    <span class="font-bold text-blue-900 block">Sistem POS Kasir (Wajib)</span>
                                    <span class="text-blue-700 text-xs">Modul utama untuk transaksi, stok, dan master produk.</span>
                                </div>
                            </label>

                            <label class="flex items-start p-3 border border-gray-200 rounded-xl cursor-pointer hover:bg-gray-50 transition-colors">
                                <div class="flex items-center h-5">
                                    <input type="checkbox" value="absensi" v-model="form.fitur_opsional" class="w-5 h-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500">
                                </div>
                                <div class="ml-3 text-sm">
                                    <span class="font-bold text-gray-900 block">Absensi Karyawan</span>
                                    <span class="text-gray-500 text-xs">Aktifkan menu absen masuk/pulang untuk karyawan toko.</span>
                                </div>
                            </label>
                        </div>
                    </div>

                    <div class="pt-4">
                        <button type="submit" :disabled="form.processing" class="w-full flex justify-center py-3.5 px-4 rounded-xl shadow-md text-sm font-black text-white bg-blue-800 hover:bg-blue-900 focus:outline-none transition-all disabled:opacity-50">
                            {{ form.processing ? 'Memproses Sistem...' : 'Buat Toko & Buka Sistem 🚀' }}
                        </button>
                    </div>

                </form>
            </div>
        </div>
    </div>
</template>
