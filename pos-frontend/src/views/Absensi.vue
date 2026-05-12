<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue';
import Sidebar from '../components/Sidebar.vue';
import api from '../api.js';

// --- STATE WAKTU REALTIME ---
const currentTime = ref('');
let timer;

const updateTime = () => {
    const now = new Date();
    currentTime.value = now.toLocaleTimeString('id-ID', {
        hour: '2-digit', minute: '2-digit', second: '2-digit'
    }).replace(/:/g, '.');
};

// --- STATE DATA ---
const karyawan = ref([]);
const riwayat = ref([]);
const isLoading = ref(true);

// Ambil data User ID yang lagi login dari JWT/LocalStorage
// Fungsi parse JWT sederhana buat ngambil user_id
const getPayloadFromToken = () => {
    const token = localStorage.getItem('token');
    if (!token) return { user_id: 0, role: '' };
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return { user_id: payload.user_id, role: payload.role };
    } catch (e) {
        return { user_id: 0, role: '' };
    }
};
const currentUser = ref(getPayloadFromToken());

// Filter Tanggal (Default hari ini)
const today = new Date().toISOString().split('T')[0];
const tanggalDipilih = ref(today);

// --- STATE KAMERA SELFIE ---
const showCameraModal = ref(false);
const videoRef = ref(null);
const stream = ref(null);
const absenTarget = ref({ id: null, nama: '', jenis: '' });
const isSubmitting = ref(false);

// --- FUNGSI TARIK DATA DARI GOLANG ---
const fetchData = async () => {
    isLoading.value = true;
    try {
        // 1. Ambil daftar karyawan (untuk tombol absen)
        const resKaryawan = await api.get('/employees');
        karyawan.value = resKaryawan.data.data || [];

        // 2. Ambil riwayat absen berdasarkan tanggal
        const resRiwayat = await api.get('/attendance', {
            params: { tanggal: tanggalDipilih.value }
        });
        riwayat.value = resRiwayat.data.data || [];
    } catch (error) {
        console.error("Gagal tarik data absensi:", error);
    } finally {
        isLoading.value = false;
    }
};

onMounted(() => {
    updateTime();
    timer = setInterval(updateTime, 1000);
    fetchData();
});

onUnmounted(() => {
    clearInterval(timer);
    stopCamera(); // Pastikan kamera mati saat pindah halaman
});

watch(tanggalDipilih, () => {
    fetchData();
});

// --- FUNGSI KAMERA & ABSEN ---
const mulaiAbsen = async (id, nama, jenis) => {
    absenTarget.value = { id, nama, jenis };
    showCameraModal.value = true;

    // Nyalakan Kamera
    try {
        stream.value = await navigator.mediaDevices.getUserMedia({ video: { facingMode: "user" } });
        await nextTick(); // Tunggu modal render
        if (videoRef.value) {
            videoRef.value.srcObject = stream.value;
        }
    } catch (err) {
        alert("Kamera tidak terdeteksi atau diblokir. Tidak bisa melakukan absensi wajah!");
        showCameraModal.value = false;
    }
};

const stopCamera = () => {
    if (stream.value) {
        stream.value.getTracks().forEach(track => track.stop());
    }
    showCameraModal.value = false;
};

const jepretDanKirim = async () => {
    isSubmitting.value = true;

    // 1. Tangkap gambar dari video
    const canvas = document.createElement('canvas');
    canvas.width = videoRef.value.videoWidth;
    canvas.height = videoRef.value.videoHeight;
    canvas.getContext('2d').drawImage(videoRef.value, 0, 0);
    
    // 2. Ubah jadi Base64
    const photoBase64 = canvas.toDataURL('image/jpeg', 0.8);

    // 3. Matikan kamera
    stopCamera();

    // 4. Kirim ke Golang
    try {
        await api.post('/attendance', {
            user_id: absenTarget.value.id,
            jenis: absenTarget.value.jenis, // 'Masuk' atau 'Pulang'
            foto: photoBase64 // Kirim foto
        });

        alert(`✅ Absen ${absenTarget.value.jenis} Berhasil! Selamat bekerja!`);
        fetchData(); // Refresh log absen
    } catch (error) {
        alert(error.response?.data?.error || "Gagal melakukan absensi");
    } finally {
        isSubmitting.value = false;
    }
};

// --- FUNGSI EXPORT ---
const downloadLaporan = async () => {
    const token = localStorage.getItem('token');
    const [tahun, bulan] = tanggalDipilih.value.split('-');
    
    try {
        const response = await fetch(`http://localhost:8080/api/attendance/export?bulan=${bulan}&tahun=${tahun}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        
        if (!response.ok) throw new Error("Gagal ekspor");

        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `Laporan_Absensi_${bulan}_${tahun}.csv`;
        document.body.appendChild(a);
        a.click();
        a.remove();
    } catch (error) {
        alert("Gagal mengunduh file Laporan.");
    }
};
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-8 max-w-7xl mx-auto font-sans">

            <div class="bg-blue-800 rounded-3xl p-6 md:p-8 mb-8 text-white shadow-xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-blue-900">
                <div class="absolute -right-10 -top-20 opacity-10 text-[200px] font-black italic pointer-events-none">⏱️</div>
                <div class="relative z-10 text-center md:text-left mb-6 md:mb-0">
                    <h1 class="text-3xl font-black tracking-tight mb-2">Mesin Absensi Wajah</h1>
                    <p class="text-blue-200 font-medium text-sm">Validasi kehadiran dengan foto *real-time* webcam/HP.</p>
                </div>
                <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner">
                    <div class="text-sm font-bold text-blue-200 uppercase tracking-widest mb-1">Waktu Saat Ini</div>
                    <div class="font-mono text-4xl md:text-5xl font-black text-yellow-300 tracking-wider drop-shadow-md">
                        {{ currentTime || '00.00.00' }}
                    </div>
                </div>
            </div>

            <h2 class="text-lg font-black text-gray-800 mb-4 flex items-center gap-2">
                👥 Daftar Karyawan Toko
            </h2>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-10">
                <div v-for="user in karyawan" :key="user.id"
                     class="bg-white rounded-3xl p-5 shadow-sm border transition-all duration-300"
                     :class="user.id === currentUser.user_id ? 'border-blue-300 ring-2 ring-blue-100 shadow-md transform hover:-translate-y-1' : 'border-gray-100 opacity-80'">

                    <div class="flex items-center gap-4 mb-5 border-b border-gray-50 pb-4">
                        <div class="w-14 h-14 rounded-full flex items-center justify-center shrink-0 border-2 border-white shadow-md"
                             :class="user.id === currentUser.user_id ? 'bg-gradient-to-br from-blue-500 to-blue-600 text-white' : 'bg-gray-100 text-gray-400'">
                            <span class="font-black text-lg uppercase">{{ user.name.substring(0, 2) }}</span>
                        </div>
                        <div>
                            <h3 class="font-black text-lg leading-tight" :class="user.id === currentUser.user_id ? 'text-blue-900' : 'text-gray-600'">
                                {{ user.name }}
                            </h3>
                            <div class="text-[10px] font-bold px-2 py-0.5 mt-1 rounded-md inline-block uppercase tracking-widest"
                                 :class="user.role === 'owner' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'">
                                {{ user.role }} • {{ user.nik || 'OWNER' }}
                            </div>
                        </div>
                    </div>

                    <div class="flex gap-3">
                        <button
                            @click="mulaiAbsen(user.id, user.name, 'Masuk')"
                            :disabled="user.id !== currentUser.user_id"
                            :class="user.id === currentUser.user_id ? 'bg-green-50 text-green-700 hover:bg-green-600 hover:text-white border-green-200 active:scale-95' : 'bg-gray-50 text-gray-400 border-gray-100 cursor-not-allowed'"
                            class="flex-1 py-3 rounded-xl font-black text-sm transition-all border flex items-center justify-center gap-1.5">
                            <span class="text-lg" :class="{'opacity-50 grayscale': user.id !== currentUser.user_id}">📸</span> Masuk
                        </button>

                        <button
                            @click="mulaiAbsen(user.id, user.name, 'Pulang')"
                            :disabled="user.id !== currentUser.user_id"
                            :class="user.id === currentUser.user_id ? 'bg-orange-50 text-orange-700 hover:bg-orange-600 hover:text-white border-orange-200 active:scale-95' : 'bg-gray-50 text-gray-400 border-gray-100 cursor-not-allowed'"
                            class="flex-1 py-3 rounded-xl font-black text-sm transition-all border flex items-center justify-center gap-1.5">
                            <span class="text-lg" :class="{'opacity-50 grayscale': user.id !== currentUser.user_id}">📸</span> Pulang
                        </button>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
                <div class="p-4 md:p-6 border-b border-gray-100 flex flex-col md:flex-row md:justify-between md:items-center bg-gray-50 gap-4">
                    <h3 class="font-black text-gray-800 text-lg flex items-center gap-2">
                        📋 Log Kehadiran
                    </h3>

                    <div class="flex items-center gap-3">
                        <input type="date" v-model="tanggalDipilih" class="px-4 py-2 bg-white border border-gray-200 rounded-xl text-sm font-bold text-gray-700 shadow-sm focus:ring-blue-500 focus:border-blue-500 cursor-pointer outline-none">

                        <button v-if="currentUser.role === 'owner'" @click="downloadLaporan" class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-xl font-bold text-sm shadow-md transition-all flex items-center gap-2 active:scale-95">
                            <span>📊</span> Download CSV
                        </button>
                    </div>
                </div>

                <div class="overflow-x-auto">
                    <table class="w-full text-left whitespace-nowrap">
                        <thead class="bg-white border-b border-gray-100">
                            <tr>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest">Karyawan</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Foto Masuk</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Jam Masuk</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Foto Pulang</th>
                                <th class="px-6 py-4 text-xs font-black text-gray-400 uppercase tracking-widest text-center">Jam Pulang</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-50">
                            <tr v-if="isLoading">
                                <td colspan="5" class="px-6 py-12 text-center text-gray-400 font-medium">Memuat data absensi...</td>
                            </tr>
                            <tr v-else-if="riwayat.length === 0">
                                <td colspan="5" class="px-6 py-12 text-center text-gray-400 font-medium italic">Belum ada data absensi terekam pada tanggal ini.</td>
                            </tr>
                            <tr v-for="log in riwayat" :key="log.id" class="hover:bg-gray-50 transition-colors">
                                <td class="px-6 py-4 font-bold text-gray-800">{{ log.User?.name }}</td>
                                
                                <td class="px-6 py-4 text-center">
                                    <div v-if="log.foto_masuk" class="w-10 h-10 rounded-lg mx-auto bg-gray-200 overflow-hidden border border-gray-300">
                                        <img :src="log.foto_masuk" class="w-full h-full object-cover">
                                    </div>
                                    <span v-else class="text-gray-300">-</span>
                                </td>
                                <td class="px-6 py-4 text-center">
                                    <span v-if="log.jam_masuk" class="bg-green-100 text-green-700 font-bold px-3 py-1 rounded-lg text-sm">{{ log.jam_masuk }}</span>
                                    <span v-else class="text-gray-300">-</span>
                                </td>

                                <td class="px-6 py-4 text-center">
                                    <div v-if="log.foto_pulang" class="w-10 h-10 rounded-lg mx-auto bg-gray-200 overflow-hidden border border-gray-300">
                                        <img :src="log.foto_pulang" class="w-full h-full object-cover">
                                    </div>
                                    <span v-else class="text-gray-300">-</span>
                                </td>
                                <td class="px-6 py-4 text-center">
                                    <span v-if="log.jam_pulang" class="bg-orange-100 text-orange-700 font-bold px-3 py-1 rounded-lg text-sm">{{ log.jam_pulang }}</span>
                                    <span v-else class="text-gray-300">-</span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div v-if="showCameraModal" class="fixed inset-0 bg-black/90 z-[100] flex items-center justify-center p-4 backdrop-blur-sm">
            <div class="bg-white rounded-3xl overflow-hidden shadow-2xl w-full max-w-md flex flex-col">
                <div class="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                    <div>
                        <h3 class="font-black text-gray-800">📸 Selfie Absen {{ absenTarget.jenis }}</h3>
                        <p class="text-xs text-gray-500 font-bold mt-0.5">{{ absenTarget.nama }}</p>
                    </div>
                    <button @click="stopCamera" class="w-8 h-8 rounded-full bg-gray-200 text-gray-600 font-black hover:bg-red-100 hover:text-red-600 transition-colors">✕</button>
                </div>
                
                <div class="relative bg-gray-900 w-full aspect-[3/4] flex items-center justify-center">
                    <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover transform scale-x-[-1]"></video>
                    <div class="absolute inset-0 border-[10px] border-white/20 pointer-events-none rounded-2xl m-4"></div>
                    <div class="absolute bottom-4 left-0 right-0 text-center pointer-events-none">
                        <span class="bg-black/50 text-white text-xs font-bold px-3 py-1 rounded-full backdrop-blur-md">Posisikan wajah Anda di tengah</span>
                    </div>
                </div>

                <div class="p-4 bg-white">
                    <button @click="jepretDanKirim" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-4 rounded-2xl font-black text-lg shadow-lg transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2">
                        {{ isSubmitting ? 'Mengirim Data...' : '◎ Jepret & Absen Sekarang' }}
                    </button>
                </div>
            </div>
        </div>

    </Sidebar>
</template>