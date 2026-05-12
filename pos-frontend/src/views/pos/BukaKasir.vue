<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import Swal from 'sweetalert2';

const role = localStorage.getItem('role');
const router = useRouter();
const name = localStorage.getItem('name'); 
const stationNumber = ref('01');
const modalAwal = ref(0);
const loading = ref(false);

onMounted(async () => {
    try {
        const res = await axios.get('http://localhost:8080/api/pos/check-session', {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        });
        
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
        Swal.fire('Eits!', 'Modal awal tidak boleh kosong.', 'warning');
        return;
    }

    loading.value = true;
    try {
        const token = localStorage.getItem('token');
        
        const res = await axios.post('http://localhost:8080/api/pos/open-session', {
            station_number: stationNumber.value,
            modal_awal: parseFloat(modalAwal.value)
        }, {
            headers: { Authorization: `Bearer ${token}` }
        });

        Swal.fire({
            icon: 'success',
            title: 'KASIR DIBUKA',
            text: `Selamat bertugas!`,
            timer: 1500,
            showConfirmButton: false
        }).then(() => {
            router.push('/pos/kasir'); 
        });

    } catch (error) {
        const msg = error.response?.data?.error || 'Gagal membuka kasir';
        
        // 🚀 LOGIKA: Jika ditolak karena belum absen
        if (msg.toLowerCase().includes('absen')) {
    Swal.fire({
        title: 'Belum Absen!',
        text: 'Mas Arif wajib absen wajah dulu.',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Absen Sekarang 📸',
    }).then((result) => {
        if (result.isConfirmed) {
            // Ganti router.push dengan ini:
            window.location.href = '/absensi'; 
        }
    });
        } else if (error.response?.status === 401) {
            router.push('/login');
        } else {
            Swal.fire('Gagal', msg, 'error');
        }
    } finally {
        loading.value = false;
    }
};
</script>

<template>
    <div class="min-h-screen bg-slate-900 flex items-center justify-center p-4">
        <div class="bg-white rounded-[45px] p-10 w-full max-w-md shadow-2xl transition-all border-[8px] border-slate-800">
            <div class="text-center mb-8">
                <div class="bg-blue-100 w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-4">
                    <span class="text-4xl">🔐</span>
                </div>
                <h1 class="text-3xl font-black text-slate-800 tracking-tight tracking-widest">INITIALIZATION</h1>
                <p class="text-slate-400 font-bold text-xs uppercase tracking-[0.2em]">ARZU STORE - POS SYSTEM</p>
            </div>

            <div class="space-y-6">
                <div class="bg-slate-50 p-4 rounded-3xl border-2 border-slate-100 flex items-center justify-between">
                    <div>
                        <label class="text-[10px] font-black text-slate-400 uppercase ml-1">Operator</label>
                        <div class="text-lg font-black text-slate-700 flex items-center gap-2">
                            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                            {{ name }}
                        </div>
                    </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                    <div v-for="n in ['01', '02', '03']" :key="n" 
                        @click="stationNumber = n"
                        :class="stationNumber === n ? 'border-blue-600 bg-blue-50 text-blue-600' : 'border-slate-100 text-slate-400'"
                        class="cursor-pointer border-2 p-4 rounded-2xl text-center font-black transition-all hover:border-blue-300">
                        STATION {{ n }}
                    </div>
                </div>

                <div>
                    <label class="text-[10px] font-black text-slate-400 uppercase ml-1 tracking-widest">Modal Awal (Drawer)</label>
                    <div class="relative mt-2">
                        <span class="absolute left-6 top-1/2 -translate-y-1/2 font-black text-blue-600 text-2xl">Rp</span>
                        <input 
                            type="text" 
                            placeholder="0"
                            :value="modalAwal === 0 ? '' : modalAwal.toLocaleString('id-ID')"
                            @input="handleInputModal"
                            class="w-full bg-blue-50/50 border-2 border-blue-100 p-6 pl-20 rounded-[30px] font-black text-3xl text-blue-600 focus:border-blue-500 outline-none transition-all"
                        >
                    </div>
                </div>

                <div class="pt-2">
                    <button @click="handleBukaKasir" :disabled="loading"
                        class="w-full bg-blue-600 hover:bg-blue-700 text-white p-6 rounded-[30px] font-black text-lg shadow-xl shadow-blue-200 transition-all active:scale-95 disabled:opacity-50">
                        {{ loading ? 'PROSES DATA...' : 'BUKA KASIR SEKARANG 🚀' }}
                    </button>
                </div>

                <div class="text-center">
                    <button 
                        @click="router.push('/dashboard')" 
                        v-if="role === 'owner'" 
                        class="text-slate-400 hover:text-slate-600 text-xs font-bold uppercase tracking-widest transition-colors">
                        Kembali ke Dashboard
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>