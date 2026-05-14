<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../api.js';
import Swal from 'sweetalert2';

const role = localStorage.getItem('role');
const router = useRouter();
const name = localStorage.getItem('name'); 
const storeName = localStorage.getItem('storeName');
const stationNumber = ref('01');
const modalAwal = ref(0);
const loading = ref(false);

onMounted(async () => {
    try {
        const res = await api.get('/pos/check-session');
        if (res.data.has_session) {
            router.push('/pos/kasir'); 
        }
    } catch (error) {
        console.error("Gagal cek session", error);
    }
});

const handleInputModal = (e) => {
    const val = e.target.value.replace(/\D/g, '');
    modalAwal.value = val ? parseInt(val, 10) : 0;
};

const handleBukaKasir = async () => {
    if (modalAwal.value <= 0) {
        Swal.fire({
            icon: 'warning',
            title: 'Modal Kosong?',
            text: 'Masukkan modal awal untuk uang kembalian di laci kasir.',
            confirmButtonColor: '#2563eb'
        });
        return;
    }

    loading.value = true;
    try {
        await api.post('/pos/open-session', {
            station_number: stationNumber.value,
            modal_awal: parseFloat(modalAwal.value)
        });

        Swal.fire({
            icon: 'success',
            title: 'SESSION OPENED',
            text: `Kasir Station ${stationNumber.value} berhasil dibuka.`,
            timer: 1500,
            showConfirmButton: false,
            customClass: { popup: 'rounded-[32px]' }
        }).then(() => {
            router.push('/pos/kasir'); 
        });

    } catch (error) {
        const msg = error.response?.data?.error || 'Gagal membuka kasir';
        
        if (msg.toLowerCase().includes('absen')) {
            Swal.fire({
                title: 'Otorisasi Gagal',
                text: 'Sistem mendeteksi Anda belum melakukan absensi hari ini.',
                icon: 'error',
                showCancelButton: true,
                confirmButtonText: 'Absen Sekarang 📸',
                confirmButtonColor: '#2563eb',
                customClass: { popup: 'rounded-[32px]' }
            }).then((result) => {
                if (result.isConfirmed) {
                    window.location.href = '/absensi'; 
                }
            });
        } else {
            Swal.fire('Error', msg, 'error');
        }
    } finally {
        loading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen bg-slate-950 flex items-center justify-center p-6 relative overflow-hidden">
        
        <div class="absolute top-0 left-0 w-full h-full opacity-10 pointer-events-none">
            <div class="absolute -top-24 -left-24 w-96 h-96 bg-blue-600 rounded-full blur-[120px]"></div>
            <div class="absolute -bottom-24 -right-24 w-96 h-96 bg-indigo-600 rounded-full blur-[120px]"></div>
        </div>

        <div class="w-full max-w-lg relative">
            <div class="bg-white rounded-[48px] p-8 md:p-12 shadow-2xl border-[12px] border-slate-900/5 relative overflow-hidden">
                
                <div class="text-center mb-10">
                    <div class="w-20 h-20 bg-slate-900 rounded-[28px] flex items-center justify-center mx-auto mb-6 shadow-xl shadow-blue-500/20 transform -rotate-3 hover:rotate-0 transition-all duration-500">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-blue-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20 20H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2z"/><path d="M2 12h20"/><path d="M7 16h.01"/><path d="M11 16h.01"/><path d="M15 16h.01"/><path d="M7 8h.01"/><path d="M11 8h.01"/><path d="M15 8h.01"/></svg>
                    </div>
                    <h1 class="text-3xl font-black text-slate-900 tracking-tighter uppercase">Point of <span class="text-blue-600">Sale</span></h1>
                    <p class="text-slate-400 font-bold text-[10px] uppercase tracking-[0.4em] mt-1">{{ storeName }}</p>
                </div>

                <div class="space-y-8">
                    <div class="bg-slate-50 p-5 rounded-[32px] border border-slate-100 flex items-center justify-between">
                        <div class="flex items-center gap-4">
                            <div class="w-12 h-12 rounded-2xl bg-white border border-slate-200 flex items-center justify-center text-lg shadow-sm">👤</div>
                            <div>
                                <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block">Logged Operator</label>
                                <div class="text-sm font-black text-slate-800 uppercase flex items-center gap-2">
                                    <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
                                    {{ name }}
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="space-y-3">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-2 italic">Select Device Station</label>
                        <div class="grid grid-cols-3 gap-3">
                            <button v-for="n in ['01', '02', '03']" :key="n" 
                                @click="stationNumber = n"
                                :class="stationNumber === n ? 'bg-slate-900 text-white shadow-xl shadow-slate-200 scale-105' : 'bg-slate-50 text-slate-400 grayscale border-transparent'"
                                class="flex flex-col items-center gap-2 py-4 rounded-[28px] border-2 transition-all duration-300">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="12" x="3" y="4" rx="2" ry="2"/><line x1="2" x2="22" y1="20" y2="20"/><line x1="12" x2="12" y1="16" y2="20"/></svg>
                                <span class="text-[10px] font-black tracking-tighter">POS {{ n }}</span>
                            </button>
                        </div>
                    </div>

                    <div class="space-y-3">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-2">Floating Capital (Modal Awal)</label>
                        <div class="relative group">
                            <div class="absolute left-6 top-1/2 -translate-y-1/2 flex flex-col items-center">
                                <span class="text-[10px] font-black text-blue-400 uppercase leading-none">IDR</span>
                                <span class="text-xl font-black text-blue-600">Rp</span>
                            </div>
                            <input 
                                type="text" 
                                placeholder="0"
                                :value="modalAwal === 0 ? '' : modalAwal.toLocaleString('id-ID')"
                                @input="handleInputModal"
                                class="w-full bg-blue-50/30 border-2 border-blue-100 p-8 pl-24 rounded-[36px] font-black text-4xl text-slate-900 focus:border-blue-600 focus:bg-white outline-none transition-all placeholder:text-slate-200"
                            >
                        </div>
                    </div>

                    <div class="pt-4">
                        <button @click="handleBukaKasir" :disabled="loading"
                            class="w-full bg-blue-600 hover:bg-slate-900 text-white p-6 rounded-[32px] font-black text-sm uppercase tracking-[0.2em] shadow-2xl shadow-blue-200 hover:shadow-slate-200 transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-4">
                            <span v-if="!loading">Initialize Session</span>
                            <span v-else class="animate-pulse">Accessing Server...</span>
                            <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                        </button>
                        
                        <div class="mt-8 text-center" v-if="role === 'owner'">
                            <button @click="router.push('/dashboard')" 
                                class="text-[10px] font-black text-slate-300 hover:text-blue-600 uppercase tracking-[0.3em] transition-colors flex items-center justify-center gap-2 mx-auto">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m15 18-6-6 6-6"/></svg>
                                Back to Dashboard
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Transisi Smooth */
.transition-all {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>