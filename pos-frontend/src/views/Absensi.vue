<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue';
import Sidebar from '../components/Sidebar.vue';
import api from '../api.js';
import Swal from 'sweetalert2';
import * as faceapi from 'face-api.js';

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
const filterMode = ref('harian');
const bulanDipilih = ref(new Date().toISOString().slice(0, 7)); // YYYY-MM
const isAiLoading = ref(true);

const getPayloadFromToken = () => {
    const token = localStorage.getItem('token');
    const name = localStorage.getItem('name') || 'User';
    const role = localStorage.getItem('role') || 'kasir';
    
    if (!token) return { user_id: 0, role: '', name: '' };
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return { user_id: payload.user_id, role: role.toLowerCase(), name: name };
    } catch (e) {
        return { user_id: 0, role: '', name: '' };
    }
};
const currentUser = ref(getPayloadFromToken());

const today = new Date().toISOString().split('T')[0];
const tanggalDipilih = ref(today);

// --- 🚀 LOGIC DETEKSI STATUS ABSEN REAKTIF (FIXED) ---
const hasAbsenMasuk = computed(() => {
    return riwayat.value.some(log => {
        // Ambil 10 karakter pertama saja (YYYY-MM-DD) dari string tanggal API
        const tglClean = log.tanggal ? log.tanggal.substring(0, 10) : "";
        
        return Number(log.user_id) === Number(currentUser.value.user_id) && 
            tglClean === today && 
            log.jam_masuk != null
    });
});

const hasAbsenPulang = computed(() => {
    return riwayat.value.some(log => {
        const tglClean = log.tanggal ? log.tanggal.substring(0, 10) : "";
        
        return Number(log.user_id) === Number(currentUser.value.user_id) && 
               tglClean === today && 
               log.jam_pulang != null && log.jam_pulang !== ""
    });
});

// --- STATE KAMERA ---
const showCameraModal = ref(false);
const videoRef = ref(null);
const stream = ref(null);
const absenTarget = ref({ id: null, nama: '', jenis: '' });
const isSubmitting = ref(false);

// --- FUNGSI TARIK DATA ---
const fetchData = async () => {
    isLoading.value = true;
    try {
        if (currentUser.value.role === 'owner') {
            const resKaryawan = await api.get('/employees');
            karyawan.value = resKaryawan.data.data || [];
        } else {
            karyawan.value = [{ 
                id: currentUser.value.user_id, 
                name: currentUser.value.name, 
                role: currentUser.value.role, 
                nik: 'KARYAWAN' 
            }];
        }

        const params = {};
        if (filterMode.value === 'harian') {
            params.tanggal = tanggalDipilih.value;
        } else {
            const [tahun, bulan] = bulanDipilih.value.split('-');
            params.bulan = bulan;
            params.tahun = tahun;
        }

        const resRiwayat = await api.get('/attendance', { params });
        // 🚀 Gunakan spread operator agar Vue mendeteksi perubahan array secara total
        riwayat.value = [...(resRiwayat.data.data || [])];
    } catch (error) {
        console.error("Gagal tarik data:", error);
    } finally {
        isLoading.value = false;
    }
};

const lihatFoto = (url, nama, tipe, jam) => {
    if (!url) return;
    Swal.fire({
        title: `Foto ${tipe}`,
        html: `<div class="text-left font-bold text-sm">👤 Karyawan: ${nama}<br>⏰ Jam: ${jam}</div>`,
        imageUrl: url,
        imageAlt: 'Foto Absensi',
        confirmButtonText: 'Tutup',
        confirmButtonColor: '#2563eb',
        customClass: { image: 'rounded-2xl border-4 border-gray-100 shadow-lg' }
    });
};

// LOGIC AI ABSENSI
onMounted(async () => {
    updateTime();
    timer = setInterval(updateTime, 1000);
    
    // 🚀 LOAD MODELS AI SAAT HALAMAN DIBUKA
    try {
        await Promise.all([
            faceapi.nets.tinyFaceDetector.loadFromUri('/models'),
            faceapi.nets.faceLandmark68Net.loadFromUri('/models'),
            faceapi.nets.faceRecognitionNet.loadFromUri('/models')
        ]);
        isAiLoading.value = false;
        console.log("AI Models Loaded! 🤖");
    } catch (e) {
        console.error("Gagal load AI models", e);
    }
    
    fetchData();
});

onUnmounted(() => {
    clearInterval(timer);
    stopCamera();
});

watch([tanggalDipilih, bulanDipilih, filterMode], () => fetchData());

const mulaiAbsen = async (id, nama, jenis) => {
    absenTarget.value = { id, nama, jenis };
    showCameraModal.value = true;
    try {
        stream.value = await navigator.mediaDevices.getUserMedia({ 
            video: { facingMode: "user", width: 1280, height: 720 } 
        });
        await nextTick();
        if (videoRef.value) videoRef.value.srcObject = stream.value;
    } catch (err) {
        Swal.fire('Kamera Error', 'Nyalakan izin kamera ya!', 'error');
        showCameraModal.value = false;
    }
};

const stopCamera = () => {
    if (stream.value) stream.value.getTracks().forEach(track => track.stop());
    showCameraModal.value = false;
};

const jepretDanKirim = async () => {
    if (isAiLoading.value) return Swal.fire('Tunggu', 'AI sedang pemanasan...', 'warning');
    
    isSubmitting.value = true;
    
    try {
        // 🚀 1. Ambil data profil diri sendiri (Pasti diizinkan Backend)
        const resMe = await api.get('/me'); 
        const targetEmp = resMe.data; // Mengambil data user yang sedang login

        // Pastikan field foto_url ada di response /api/me
        if (!targetEmp.foto_url) {
            throw new Error("Foto Master Anda belum ada di sistem. Minta Owner untuk upload foto Anda.");
        }

        // 🚀 2. Load Foto Master dari Server (Gunakan URL Lengkap)
        const masterUrl = 'http://localhost:8080' + targetEmp.foto_url;
        const imgMaster = await faceapi.fetchImage(masterUrl);
        
        const masterDetections = await faceapi.detectSingleFace(imgMaster, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!masterDetections) throw new Error("AI tidak menemukan wajah di Foto Master. Coba upload foto baru yang lebih jelas.");

        const faceMatcher = new faceapi.FaceMatcher(masterDetections);

        // 🚀 3. Deteksi Wajah di Kamera (Video)
        const queryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!queryDetections) throw new Error("Wajah tidak terdeteksi di kamera! Pastikan wajah terlihat jelas dan cukup cahaya.");

        // Bandingkan wajah kamera dengan wajah master
        const bestMatch = faceMatcher.findBestMatch(queryDetections.descriptor);
        
        // Tolerance: 0.5 (Makin kecil makin ketat)
        if (bestMatch.distance > 0.5) {
            throw new Error("Wajah tidak cocok! Verifikasi gagal.");
        }

        // 🚀 4. Lolos Verifikasi, Kecilkan Foto untuk dikirim ke Backend
        const canvas = document.createElement('canvas');
        canvas.width = 320; 
        canvas.height = 400;
        const ctx = canvas.getContext('2d');
        ctx.drawImage(videoRef.value, 0, 0, canvas.width, canvas.height);

        // Quality 0.4 biar gak kena 413 Payload Too Large lagi
        const photoBase64 = canvas.toDataURL('image/jpeg', 0.4); 
        
        stopCamera(); // Matikan kamera setelah sukses

        // Kirim ke database
        await api.post('/attendance', {
            user_id: targetEmp.user_id, // Sesuaikan dengan key di response /api/me
            jenis: absenTarget.value.jenis,
            foto: photoBase64
        });

        Swal.fire({
            icon: 'success',
            title: `Absen ${absenTarget.value.jenis} Berhasil!`,
            text: 'Selamat bekerja!',
            timer: 2000,
            showConfirmButton: false
        });

        fetchData(); // Refresh tabel riwayat absen

    } catch (error) {
        console.error("DEBUG ABSEN:", error);
        Swal.fire('Gagal', error.message || "Terjadi kesalahan", 'error');
    } finally {
        isSubmitting.value = false;
    }
};
const downloadLaporan = async () => {
    let bulan, tahun;
    filterMode.value === 'harian' ? [tahun, bulan] = tanggalDipilih.value.split('-').slice(0, 2) : [tahun, bulan] = bulanDipilih.value.split('-');
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8080/api/attendance/export?bulan=${bulan}&tahun=${tahun}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url; a.download = `ABSENSI_${bulan}_${tahun}.csv`; a.click();
    } catch (error) { Swal.fire('Gagal!', 'Gagal download laporan.', 'error'); }
};
</script>

<template>
    <Sidebar>
        <div class="p-6 md:p-8 max-w-7xl mx-auto font-sans">
            <div class="bg-blue-800 rounded-3xl p-6 md:p-8 mb-8 text-white shadow-xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-blue-900">
                <div class="absolute -right-10 -top-20 opacity-10 text-[200px] font-black italic pointer-events-none">⏱️</div>
                <div class="relative z-10 text-center md:text-left mb-6 md:mb-0">
                    <h1 class="text-3xl font-black tracking-tight mb-2 uppercase">Presensi Wajah</h1>
                    <p class="text-blue-200 font-medium text-sm italic">Capture Selfie + Timestamp Realtime</p>
                </div>
                <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner">
                    <div class="text-[10px] font-black text-blue-200 uppercase tracking-[0.2em] mb-1">Status Shift</div>
                    <div class="font-mono text-2xl font-black text-yellow-300">
                        {{ hasAbsenMasuk ? (hasAbsenPulang ? '🏁 SELESAI' : '🟢 AKTIF') : '⚪ BELUM' }}
                    </div>
                </div>
            </div>

            <h2 class="text-lg font-black text-gray-800 mb-4 flex items-center gap-2 uppercase tracking-tight">👤 Panel Absensi Anda</h2>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-10">
                <div v-for="user in karyawan" :key="user.id"
                     class="bg-white rounded-3xl p-5 shadow-sm border-2 transition-all duration-300"
                     :class="user.id == currentUser.user_id ? 'border-blue-500 shadow-blue-100 ring-4 ring-blue-50' : 'border-gray-100 opacity-60'">

                    <div class="flex items-center gap-4 mb-5 pb-4 border-b border-gray-50">
                        <div class="w-14 h-14 rounded-full flex items-center justify-center shrink-0 border-2 border-white shadow-md bg-blue-600 text-white font-black text-lg uppercase">{{ user.name.substring(0, 2) }}</div>
                        <div>
                            <h3 class="font-black text-lg leading-tight text-gray-800 uppercase">{{ user.name }}</h3>
                            <div class="text-[10px] font-black px-2 py-0.5 mt-1 rounded bg-blue-100 text-blue-700 uppercase tracking-widest">{{ user.role }} • {{ user.nik || 'ID:'+user.id }}</div>
                        </div>
                    </div>

                    <div class="flex gap-3">
                        <button v-if="!hasAbsenMasuk" @click="mulaiAbsen(user.id, user.name, 'Masuk')"
                            :disabled="user.id != currentUser.user_id"
                            class="flex-1 py-3.5 rounded-2xl font-black text-sm transition-all border-2 flex items-center justify-center gap-2 bg-green-50 text-green-700 border-green-200 hover:bg-green-600 hover:text-white disabled:opacity-30">
                            📸 MASUK
                        </button>
                        
                        <button v-if="hasAbsenMasuk && !hasAbsenPulang" @click="mulaiAbsen(user.id, user.name, 'Pulang')"
                            :disabled="user.id != currentUser.user_id"
                            class="flex-1 py-3.5 rounded-2xl font-black text-sm transition-all border-2 flex items-center justify-center gap-2 bg-orange-50 text-orange-700 border-orange-200 hover:bg-orange-600 hover:text-white disabled:opacity-30">
                            🏁 PULANG
                        </button>

                        <div v-if="hasAbsenPulang" class="w-full py-3.5 text-center bg-gray-100 text-gray-400 font-black rounded-2xl text-xs uppercase tracking-widest border-2 border-dashed border-gray-200">
                            TUGAS SELESAI HARI INI
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
                <div class="p-5 md:p-6 border-b border-gray-100 flex flex-col md:flex-row md:justify-between md:items-center bg-gray-50 gap-4">
                    <h3 class="font-black text-gray-800 text-lg flex items-center gap-2 uppercase tracking-tight">📋 Log Riwayat</h3>
                    <div class="flex flex-wrap items-center gap-3">
                        <select v-model="filterMode" class="px-4 py-2 bg-white border-2 border-gray-200 rounded-xl text-sm font-black text-gray-700 outline-none">
                            <option value="harian">HARIAN</option><option value="bulanan">BULANAN</option>
                        </select>
                        <input v-if="filterMode === 'harian'" type="date" v-model="tanggalDipilih" class="px-4 py-2 border-2 border-gray-200 rounded-xl text-sm font-black text-gray-700">
                        <input v-else type="month" v-model="bulanDipilih" class="px-4 py-2 border-2 border-gray-200 rounded-xl text-sm font-black text-gray-700">
                        <button v-if="currentUser.role === 'owner'" @click="downloadLaporan" class="bg-emerald-600 text-white px-5 py-2 rounded-xl font-black text-xs uppercase shadow-md">Ekspor CSV</button>
                    </div>
                </div>
                <div class="overflow-x-auto">
                    <table class="w-full text-left whitespace-nowrap">
                        <thead class="bg-white border-b border-gray-100">
                            <tr class="text-[10px] font-black text-gray-400 uppercase tracking-[0.2em]">
                                <th class="px-6 py-5">Karyawan</th>
                                <th class="px-6 py-5 text-center">Foto Masuk</th>
                                <th class="px-6 py-5 text-center">Jam Masuk</th>
                                <th class="px-6 py-5 text-center">Foto Pulang</th>
                                <th class="px-6 py-5 text-center">Jam Pulang</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-gray-50">
                            <tr v-if="isLoading"><td colspan="5" class="px-6 py-12 text-center text-gray-400 font-bold uppercase animate-pulse">Sinkronisasi...</td></tr>
                            <tr v-else-if="riwayat.length === 0"><td colspan="5" class="px-6 py-12 text-center text-gray-400 font-bold italic">Data Kosong.</td></tr>
                            <tr v-for="log in riwayat" :key="log.id" class="hover:bg-blue-50/30 transition-colors group">
                                <td class="px-6 py-5">
                                    <div class="font-black text-gray-800 uppercase text-sm">{{ log.User?.name }}</div>
                                    <div class="text-[9px] text-gray-400 font-bold tracking-widest">{{ log.User?.nik || 'ID: '+log.user_id }}</div>
                                </td>
                                <td class="px-6 py-5 text-center">
                                    <div v-if="log.foto_masuk" @click="lihatFoto(log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-14 h-14 rounded-xl mx-auto border-4 border-white shadow-sm overflow-hidden cursor-zoom-in"><img :src="log.foto_masuk" class="w-full h-full object-cover"></div>
                                    <span v-else class="text-gray-200">❌</span>
                                </td>
                                <td class="px-6 py-5 text-center"><span v-if="log.jam_masuk" class="bg-green-100 text-green-700 font-black px-3 py-1.5 rounded-lg text-xs">{{ log.jam_masuk }}</span></td>
                                <td class="px-6 py-5 text-center">
                                    <div v-if="log.foto_pulang" @click="lihatFoto(log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-14 h-14 rounded-xl mx-auto border-4 border-white shadow-sm overflow-hidden cursor-zoom-in"><img :src="log.foto_pulang" class="w-full h-full object-cover"></div>
                                    <span v-else class="text-gray-200">❌</span>
                                </td>
                                <td class="px-6 py-5 text-center"><span v-if="log.jam_pulang" class="bg-orange-100 text-orange-700 font-black px-3 py-1.5 rounded-lg text-xs">{{ log.jam_pulang }}</span></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div v-if="showCameraModal" class="fixed inset-0 bg-gray-900/95 z-[100] flex items-center justify-center p-4 backdrop-blur-md">
            <div class="bg-white rounded-[40px] overflow-hidden shadow-2xl w-full max-w-md border-8 border-white">
                <div class="p-6 border-b border-gray-50 flex justify-between items-center bg-white">
                    <div>
                        <h3 class="font-black text-gray-800 uppercase tracking-tighter text-xl">Selfie Absen</h3>
                        <p class="text-xs text-blue-600 font-black uppercase mt-1">{{ absenTarget.jenis }} • {{ absenTarget.nama }}</p>
                    </div>
                    <button @click="stopCamera" class="w-10 h-10 rounded-full bg-gray-100 text-gray-400 hover:text-red-600 transition-all font-black">✕</button>
                </div>
                <div class="relative bg-black w-full aspect-[3/4] flex items-center justify-center overflow-hidden">
                    <div v-if="isAiLoading" class="absolute inset-0 bg-slate-900/90 backdrop-blur-md flex flex-col items-center justify-center z-50 text-white">
        <div class="w-10 h-10 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <p class="font-black text-[10px] uppercase tracking-[0.2em]">AI Engine Pemanasan...</p>
    </div>
                    <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover transform scale-x-[-1]"></video>
                    <div class="absolute bottom-6 left-6 text-yellow-400 font-mono text-xs drop-shadow-[0_2px_2px_rgba(0,0,0,0.8)]">
                        <div class="font-black">{{ absenTarget.nama.toUpperCase() }}</div>
                        <div class="text-lg font-black tracking-widest">{{ currentTime }}</div>
                        <div class="font-bold opacity-80">{{ new Date().toLocaleDateString('id-ID') }}</div>
                    </div>
                </div>
                <div class="p-6 bg-white">
                    <button @click="jepretDanKirim" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-blue-700 text-white py-5 rounded-[24px] font-black text-xl shadow-2xl transition-all disabled:opacity-50">
                        {{ isSubmitting ? 'MENGIRIM...' : '◎ JEPRET & ABSEN' }}
                    </button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }
</style>