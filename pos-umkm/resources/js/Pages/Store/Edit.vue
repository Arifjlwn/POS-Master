<script setup>
import { Head, useForm } from '@inertiajs/vue3';
import AppLayout from '@/Layouts/AppLayout.vue';

const props = defineProps({
    store: Object
});

// Isi form otomatis dengan data toko saat ini
const form = useForm({
    nama_toko: props.store.nama_toko,
    alamat_toko: props.store.alamat_toko,
    telepon: props.store.telepon,
    // Pastikan array tidak null
    fitur_opsional: props.store.fitur_opsional || ['kasir'],
});

const submit = () => {
    form.post(route('store.update'), {
        preserveScroll: true
    });
};
</script>

<template>
    <Head title="Pengaturan Toko" />

    <AppLayout>
        <div class="p-6 md:p-8 max-w-4xl mx-auto">
            <div class="mb-8">
                <h1 class="text-3xl font-black text-gray-900 leading-tight">Pengaturan Toko</h1>
                <p class="text-gray-500 font-medium">Ubah profil bisnis dan kelola modul aplikasi yang aktif.</p>
            </div>

            <transition name="fade">
                <div v-if="$page.props.flash.message" class="mb-6 p-4 bg-green-500 text-black rounded-2xl shadow-lg shadow-green-200 flex items-center justify-between">
                    <div class="flex items-center gap-3">
                        <span class="text-xl">✅</span>
                        <span class="font-bold">{{ $page.props.flash.message }}</span>
                    </div>
                    <button @click="$page.props.flash.message = null" class="text-white/70 hover:text-white font-black">✕</button>
                </div>
            </transition>

            <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
                <form @submit.prevent="submit" class="p-6 md:p-8 space-y-8">

                    <div>
                        <h3 class="text-lg font-black text-gray-800 border-b pb-3 mb-5 flex items-center gap-2">
                            <span>🏪</span> Profil Bisnis
                        </h3>
                        <div class="space-y-4">
                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Nama Toko</label>
                                <input v-model="form.nama_toko" type="text" required class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-blue-500 focus:border-blue-500 font-bold text-sm">
                            </div>
                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Alamat Lengkap</label>
                                <textarea v-model="form.alamat_toko" rows="2" required class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-blue-500 focus:border-blue-500 font-bold text-sm"></textarea>
                            </div>
                            <div>
                                <label class="block text-xs font-black text-gray-400 uppercase tracking-widest mb-1.5">Telepon / WhatsApp</label>
                                <input v-model="form.telepon" type="text" class="w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-xl focus:ring-blue-500 focus:border-blue-500 font-bold text-sm">
                            </div>
                        </div>
                    </div>

                    <div>
                        <h3 class="text-lg font-black text-gray-800 border-b pb-3 mb-5 flex items-center gap-2 mt-8">
                            <span>🧩</span> Modul Aplikasi
                        </h3>
                        <div class="space-y-4">

                            <label class="flex items-center justify-between p-4 border border-blue-200 bg-blue-50/50 rounded-2xl cursor-not-allowed">
                                <div class="flex flex-col">
                                    <span class="font-black text-blue-900">Sistem POS Kasir</span>
                                    <span class="text-blue-600 text-xs font-medium">Modul wajib untuk transaksi dan produk.</span>
                                </div>
                                <input type="checkbox" checked disabled class="w-6 h-6 text-blue-600 rounded-md border-gray-300 focus:ring-blue-500 opacity-50">
                            </label>

                            <label class="flex items-center justify-between p-4 border border-gray-200 hover:border-blue-300 rounded-2xl cursor-pointer transition-colors" :class="{'bg-blue-50/30 border-blue-300': form.fitur_opsional.includes('absensi')}">
                                <div class="flex flex-col">
                                    <span class="font-black text-gray-900">Absensi & Manajemen Karyawan</span>
                                    <span class="text-gray-500 text-xs font-medium">Nyalakan jika toko dikelola oleh lebih dari 1 orang.</span>
                                </div>
                                <input type="checkbox" value="absensi" v-model="form.fitur_opsional" class="w-6 h-6 text-blue-600 rounded-md border-gray-300 focus:ring-blue-500">
                            </label>

                        </div>
                    </div>

                    <div class="pt-6 border-t border-gray-100 flex justify-end">
                        <button type="submit" :disabled="form.processing" class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3.5 rounded-xl font-black text-sm shadow-lg shadow-blue-200 transition-all active:scale-95">
                            {{ form.processing ? 'Menyimpan...' : 'Simpan Pengaturan 🚀' }}
                        </button>
                    </div>

                </form>
            </div>
        </div>
    </AppLayout>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.5s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
