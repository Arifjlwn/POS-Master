<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import Swal from 'sweetalert2';
import api from '../../api.js';

const router = useRouter();

const form = ref({
    nama_toko: '',
    tipe_bisnis: 'Retail (Toko Kelontong/Minimarket)',
    alamat_toko: '',
    telepon: '',
    fitur_opsional: ['kasir'],
});

const isLoading = ref(false);
const businessTypes = ['Retail (Toko Kelontong/Minimarket)', 'F&B (Resto/Cafe)', 'Jasa (Laundry/Barbershop)', 'Lainnya'];

const submit = async () => {
    isLoading.value = true;
    try {
        const response = await api.post('/setup', {
            nama_toko: form.value.nama_toko,
            tipe_bisnis: form.value.tipe_bisnis,
            alamat_toko: form.value.alamat_toko,
            telepon: form.value.telepon
        });
        
        // UPDATE TOKEN SECARA SENYAP
        if (response.data && response.data.token) {
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('storeName', form.value.nama_toko);
        }
        
        await Swal.fire({
            icon: 'success',
            title: 'Toko Berhasil Dibuka!',
            text: `Selamat datang, Bos ${localStorage.getItem('name')}! Mari mulai bisnismu.`,
            confirmButtonColor: '#2563eb',
            customClass: { popup: 'rounded-[32px]' }
        });
        
        router.push('/dashboard'); 
    } catch (error) {
        Swal.fire({
            icon: 'error',
            title: 'Gagal Setup Toko',
            text: error.response?.data?.error || 'Terjadi kesalahan sistem.',
            confirmButtonColor: '#ef4444'
        });
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen bg-[#F8FAFC] flex flex-col justify-center py-12 sm:px-6 lg:px-8 font-sans relative overflow-hidden">
        
        <div class="absolute -top-24 -left-24 w-96 h-96 bg-blue-100/40 rounded-full blur-3xl"></div>
        <div class="absolute -bottom-24 -right-24 w-96 h-96 bg-indigo-100/40 rounded-full blur-3xl"></div>

        <div class="sm:mx-auto sm:w-full sm:max-w-xl text-center relative z-10">
            <div class="w-20 h-20 bg-blue-600 rounded-[28px] flex items-center justify-center mx-auto shadow-2xl shadow-blue-200 mb-6 transform -rotate-6 transition-transform hover:rotate-0 duration-500">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9h18v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V9Z"/><path d="m3 9 2.45-4.9A2 2 0 0 1 7.24 3h9.52a2 2 0 0 1 1.8 1.1L21 9"/><path d="M12 3v6"/></svg>
            </div>
            <h2 class="text-4xl font-black text-slate-900 tracking-tighter">Konfigurasi <span class="text-blue-600">Toko</span></h2>
            <p class="mt-3 text-slate-400 font-bold text-xs uppercase tracking-[0.2em]">Langkah Terakhir Menuju Kesuksesan</p>
        </div>

        <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-xl px-4 relative z-10">
            <div class="bg-white/80 backdrop-blur-xl py-10 px-8 shadow-2xl rounded-[40px] border border-white">
                <form class="space-y-8" @submit.prevent="submit">

                    <div class="space-y-6">
                        <div class="flex items-center gap-3 border-b border-slate-100 pb-3">
                            <div class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 font-black text-sm">1</div>
                            <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Identitas Bisnis</h3>
                        </div>

                        <div class="grid grid-cols-1 gap-5">
                            <div class="space-y-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Nama Brand / Toko</label>
                                <input v-model="form.nama_toko" type="text" required class="input-modern" placeholder="Contoh: Toko Berkah Jaya">
                            </div>

                            <div class="space-y-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Kategori Bisnis</label>
                                <select v-model="form.tipe_bisnis" class="input-modern bg-white cursor-pointer appearance-none">
                                    <option v-for="type in businessTypes" :key="type" :value="type">{{ type }}</option>
                                </select>
                            </div>

                            <div class="space-y-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Alamat Operasional</label>
                                <textarea v-model="form.alamat_toko" rows="2" required class="input-modern resize-none" placeholder="Alamat lengkap lokasi toko..."></textarea>
                            </div>

                            <div class="space-y-2">
                                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">No. WhatsApp Bisnis</label>
                                <div class="relative">
                                    <span class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 font-bold text-sm">+62</span>
                                    <input v-model="form.telepon" type="text" class="input-modern pl-14" placeholder="812345678xx">
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="space-y-6">
                        <div class="flex items-center gap-3 border-b border-slate-100 pb-3">
                            <div class="w-8 h-8 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-600 font-black text-sm">2</div>
                            <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">Fitur SaaS Aktif</h3>
                        </div>

                        <div class="grid grid-cols-1 gap-3">
                            <label class="feature-card border-blue-200 bg-blue-50/50 cursor-not-allowed">
                                <input type="checkbox" value="kasir" v-model="form.fitur_opsional" disabled class="w-5 h-5 text-blue-600 rounded-lg border-gray-300">
                                <div class="ml-4">
                                    <span class="font-black text-blue-900 block text-sm uppercase">Cloud POS System</span>
                                    <span class="text-blue-600 text-[10px] font-bold uppercase tracking-widest italic">Modul Inti (Wajib)</span>
                                </div>
                            </label>

                            <label class="feature-card border-slate-100 hover:bg-slate-50 cursor-pointer transition-all">
                                <input type="checkbox" value="absensi" v-model="form.fitur_opsional" class="w-5 h-5 text-blue-600 rounded-lg border-gray-300">
                                <div class="ml-4">
                                    <span class="font-black text-slate-800 block text-sm uppercase">Biometric Attendance</span>
                                    <span class="text-slate-400 text-[10px] font-medium uppercase tracking-tight">Presensi Wajah & Kehadiran</span>
                                </div>
                            </label>
                        </div>
                    </div>

                    <button type="submit" :disabled="isLoading" class="btn-submit">
                        <template v-if="!isLoading">
                            Luncurkan Bisnis Sekarang 🚀
                        </template>
                        <template v-else>
                            <svg class="animate-spin h-5 w-5 text-white mr-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            MENYIAPKAN INFRASTRUKTUR...
                        </template>
                    </button>

                </form>
            </div>
            <p class="mt-8 text-center text-[10px] font-black text-slate-300 uppercase tracking-[0.4em]">Integrated Business Intelligence &copy; 2026</p>
        </div>
    </div>
</template>

<style scoped>
.input-modern {
    @apply block w-full px-5 py-4 bg-slate-50 border-2 border-slate-50 rounded-2xl focus:bg-white focus:border-blue-600 focus:ring-4 focus:ring-blue-500/5 outline-none font-bold text-slate-800 transition-all placeholder:text-slate-300 placeholder:font-medium;
}

.feature-card {
    @apply flex items-center p-5 border-2 rounded-[24px] transition-all;
}

.btn-submit {
    @apply w-full flex items-center justify-center py-5 px-6 rounded-[24px] shadow-2xl shadow-blue-200 text-sm font-black text-white bg-blue-600 hover:bg-slate-900 transition-all active:scale-95 disabled:opacity-50 uppercase tracking-widest;
}

/* Custom Scrollbar */
textarea::-webkit-scrollbar { width: 4px; }
textarea::-webkit-scrollbar-thumb { background: #E2E8F0; border-radius: 10px; }
</style>