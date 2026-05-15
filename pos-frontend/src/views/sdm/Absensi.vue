<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue';
import Sidebar from '../../components/Sidebar.vue';
import api from '../../api.js';
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
const urutanTanggalTerbaru = ref(true);
const isLoading = ref(true);
const filterMode = ref('harian');
const bulanDipilih = ref(new Date().toISOString().slice(0, 7));
const isAiLoading = ref(true);

const toggleSortTanggal = () => {
    urutanTanggalTerbaru.value = !urutanTanggalTerbaru.value;
    if (urutanTanggalTerbaru.value) {
        riwayat.value.sort((a, b) => b.tanggal.localeCompare(a.tanggal));
    } else {
        riwayat.value.sort((a, b) => a.tanggal.localeCompare(b.tanggal));
    }
};

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

const today = new Date().toLocaleDateString('en-CA');
const tanggalDipilih = ref(today);

const cekAbsenMasukKaryawan = (userId) => {
    return riwayat.value.some(log => {
        if (!log.tanggal) return false;
        const tglClean = log.tanggal.substring(0, 10);
        return Number(log.user_id) === Number(userId) && 
               tglClean === today && 
               log.jam_masuk != null;
    });
};

const cekAbsenPulangKaryawan = (userId) => {
    return riwayat.value.some(log => {
        if (!log.tanggal) return false;
        const tglClean = log.tanggal.substring(0, 10);
        return Number(log.user_id) === Number(userId) && 
               tglClean === today && 
               log.jam_pulang != null && log.jam_pulang !== "";
    });
};

// --- STATE KAMERA ---
const showCameraModal = ref(false);
const videoRef = ref(null);
const stream = ref(null);
const absenTarget = ref({ id: null, nama: '', jenis: '' });
const isSubmitting = ref(false);

const fetchData = async () => {
    isLoading.value = true;
    try {
        const resKaryawan = await api.get('/employees');
        const allEmp = resKaryawan.data.data || [];
        const staffSaja = allEmp.filter(e => e.role !== 'owner');

        if (currentUser.value.role === 'owner') {
            if (staffSaja.length === 0) {
                karyawan.value = allEmp; 
            } else {
                karyawan.value = staffSaja; 
            }
        } else {
            karyawan.value = staffSaja;
        }

        const [tahun, bulan] = bulanDipilih.value.split('-');
        const tanggalAwalBulan = `${bulanDipilih.value}-01`;
        const tanggalAkhirBulan = `${bulanDipilih.value}-${new Date(tahun, bulan, 0).getDate()}`;

        const params = {};
        if (filterMode.value === 'harian') {
            params.tanggal = tanggalDipilih.value;
        } else {
            params.bulan = bulan;
            params.tahun = tahun;
        }

        const [resRiwayat, resSched] = await Promise.all([
            api.get('/attendance', { params }),
            api.get('/schedules', { params: { 
                start_date: filterMode.value === 'harian' ? tanggalDipilih.value : tanggalAwalBulan, 
                end_date: filterMode.value === 'harian' ? tanggalDipilih.value : tanggalAkhirBulan 
            } })
        ]);

        const dataAbsenReal = resRiwayat.data.data || [];
        const dataJadwalReal = resSched.data.data || [];
        const matriksGabungan = [];

        if (filterMode.value === 'harian') {
            karyawan.value.forEach(emp => {
                const empKey = emp.id || emp.user_id;
                const jadwalHariIni = dataJadwalReal.find(s => Number(s.user_id) === Number(empKey) && s.tanggal.substring(0, 10) === tanggalDipilih.value);
                const absenHariIni = dataAbsenReal.find(a => Number(a.user_id) === Number(empKey) && a.tanggal.substring(0, 10) === tanggalDipilih.value);

                if (jadwalHariIni && jadwalHariIni.shift_type === 'OFF') return;

                matriksGabungan.push({
                    id: absenHariIni?.id || `temp-harian-${empKey}`,
                    user_id: empKey,
                    tanggal: tanggalDipilih.value,
                    User: emp,
                    shift: jadwalHariIni ? jadwalHariIni.shift_type.replace(' (Approved)','').replace(' (Pending)','') : 'Belum Set',
                    foto_masuk: absenHariIni?.foto_masuk || null,
                    jam_masuk: absenHariIni?.jam_masuk || null,
                    foto_pulang: absenHariIni?.foto_pulang || null,
                    jam_pulang: absenHariIni?.jam_pulang || null,
                    status: absenHariIni?.status || (tanggalDipilih.value < today ? 'Mangkir' : 'Belum Absen')
                });
            });
        } else {
            const jumlahHari = new Date(tahun, bulan, 0).getDate();
            karyawan.value.forEach(emp => {
                const empKey = emp.id || emp.user_id;

                for (let hari = 1; hari <= jumlahHari; hari++) {
                    const tglLoopStr = `${bulanDipilih.value}-${String(hari).padStart(2, '0')}`;
                    const absenMatch = dataAbsenReal.find(a => Number(a.user_id) === Number(empKey) && a.tanggal.substring(0, 10) === tglLoopStr);
                    const jadwalMatch = dataJadwalReal.find(s => Number(s.user_id) === Number(empKey) && s.tanggal.substring(0, 10) === tglLoopStr);

                    const shiftClean = jadwalMatch ? jadwalMatch.shift_type.replace(' (Approved)','').replace(' (Pending)','') : 'Belum Set';
                    let statusDinamis = 'Belum Absen';
                    
                    if (absenMatch) {
                        statusDinamis = absenMatch.status;
                    } else if (shiftClean === 'OFF') {
                        statusDinamis = 'Libur (OFF)';
                    } else if (emp.role === 'owner') {
                        statusDinamis = 'Owner';
                    } else if (tglLoopStr < today) {
                        statusDinamis = 'Mangkir';
                    }

                    if (tglLoopStr <= today || absenMatch) {
                        matriksGabungan.push({
                            id: absenMatch?.id || `temp-bulan-${empKey}-${tglLoopStr}`,
                            user_id: empKey,
                            tanggal: tglLoopStr,
                            User: emp,
                            shift: shiftClean,
                            foto_masuk: absenMatch?.foto_masuk || null,
                            jam_masuk: absenMatch?.jam_masuk || null,
                            foto_pulang: absenMatch?.foto_pulang || null,
                            jam_pulang: absenMatch?.jam_pulang || null,
                            status: statusDinamis
                        });
                    }
                }
            });
        }

        riwayat.value = [...matriksGabungan].sort((a, b) => b.tanggal.localeCompare(a.tanggal));
    } catch (error) {
        console.error("Gagal sinkronisasi data log:", error);
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

onMounted(async () => {
    updateTime();
    timer = setInterval(updateTime, 1000);
    
    try {
        await Promise.all([
            faceapi.nets.tinyFaceDetector.loadFromUri('/models'),
            faceapi.nets.faceLandmark68Net.loadFromUri('/models'),
            faceapi.nets.faceRecognitionNet.loadFromUri('/models')
        ]);
        isAiLoading.value = false;
        console.log("AI Models Loaded!");
    } catch (e) {
        console.error("Gagal load AI models. Pastikan folder /models ada di public folder Vue.", e);
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
    if (isAiLoading.value) return Swal.fire('Tunggu', 'Sistem AI sedang memuat model...', 'warning');
    isSubmitting.value = true;
    
    try {
        const resMe = await api.get('/me'); 
        const targetEmp = resMe.data;

        if (!targetEmp.biometric_url) {
            throw new Error("Data Biometrik Wajah belum terdaftar. Minta HRD untuk mendaftarkan wajah Anda.");
        }

        // 🚀 FIX: Sanitasi URL Foto agar bebas error CORS / Path nyasar
        const baseUrl = import.meta.env.VITE_API_BASE_URL.replace(/\/$/, ''); // Buang slash di akhir base url
        const photoPath = targetEmp.biometric_url.startsWith('/') ? targetEmp.biometric_url : `/${targetEmp.biometric_url}`; // Pastikan ada slash di awal path
        const masterUrl = `${baseUrl}${photoPath}`;
        
        console.log("Mencoba load foto master Biometrik dari:", masterUrl);
        const imgMaster = await faceapi.fetchImage(masterUrl);
        
        const masterDetections = await faceapi.detectSingleFace(imgMaster, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!masterDetections) throw new Error("Wajah di foto sistem tidak terbaca oleh AI. Harap ganti foto profil.");

        const faceMatcher = new faceapi.FaceMatcher(masterDetections);
        const queryDetections = await faceapi.detectSingleFace(videoRef.value, new faceapi.TinyFaceDetectorOptions())
            .withFaceLandmarks().withFaceDescriptor();

        if (!queryDetections) throw new Error("Wajah tidak terdeteksi di kamera! Pastikan wajah di tengah dan terang.");

        const bestMatch = faceMatcher.findBestMatch(queryDetections.descriptor);
        
        if (bestMatch.distance > 0.5) {
            throw new Error("Wajah tidak cocok dengan database! (Verifikasi Ditolak)");
        }

        const canvas = document.createElement('canvas');
        canvas.width = 320; 
        canvas.height = 400;
        const ctx = canvas.getContext('2d');
        ctx.drawImage(videoRef.value, 0, 0, canvas.width, canvas.height);
        const photoBase64 = canvas.toDataURL('image/jpeg', 0.4); 
        
        stopCamera();

        await api.post('/attendance', {
            user_id: targetEmp.user_id, 
            jenis: absenTarget.value.jenis,
            foto: photoBase64
        });

        Swal.fire({
            icon: 'success',
            title: `Absen ${absenTarget.value.jenis} Berhasil!`,
            text: 'Data telah tercatat di server.',
            timer: 2000,
            showConfirmButton: false
        });

        fetchData(); 

    } catch (error) {
        console.error("DEBUG ABSEN ERROR:", error);
        // Error handling yang lebih spesifik
        let errorMsg = error.message || "Terjadi kesalahan";
        if (errorMsg.includes("Failed to fetch")) {
            errorMsg = "Gagal memuat foto server. Cek koneksi atau setting CORS backend.";
        }
        Swal.fire('Gagal Verifikasi', errorMsg, 'error');
    } finally {
        isSubmitting.value = false;
    }
};

const downloadLaporan = async () => {
    let bulan, tahun;
    filterMode.value === 'harian' ? [tahun, bulan] = tanggalDipilih.value.split('-').slice(0, 2) : [tahun, bulan] = bulanDipilih.value.split('-');
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/attendance/export?bulan=${bulan}&tahun=${tahun}`, {
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
        <div class="p-4 md:p-8 max-w-7xl mx-auto font-sans">
            
            <div class="bg-blue-800 rounded-3xl p-6 md:p-8 mb-8 text-white shadow-xl flex flex-col md:flex-row items-center justify-between overflow-hidden relative border border-blue-900">
                <div class="absolute -right-5 -top-10 opacity-10 text-blue-200 pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-64 h-64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                </div>
                
                <div class="relative z-10 text-center md:text-left mb-6 md:mb-0">
                    <h1 class="text-3xl font-black tracking-tight mb-2 uppercase">Presensi Karyawan</h1>
                    <p class="text-blue-200 font-medium text-sm italic flex items-center justify-center md:justify-start gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 9a2 2 0 0 1-2 2H6l-4 4V4c0-1.1.9-2 2-2h8a2 2 0 0 1 2 2v5Z"/><path d="M18 9h2a2 2 0 0 1 2 2v11l-4-4h-6a2 2 0 0 1-2-2v-1"/></svg>
                        Absen Tepat Waktu Ya Tim!
                    </p>
                </div>
                
                <div class="relative z-10 bg-white/10 backdrop-blur-md border border-white/20 rounded-2xl p-4 md:px-8 text-center shadow-inner">
                    <div class="text-[10px] font-black text-blue-200 uppercase tracking-[0.2em] mb-1">Status Shift Anda</div>
                    <div class="font-mono text-2xl font-black flex items-center justify-center gap-2"
                         :class="cekAbsenMasukKaryawan(currentUser.user_id) ? (cekAbsenPulangKaryawan(currentUser.user_id) ? 'text-slate-300' : 'text-emerald-400') : 'text-amber-300'">
                        <template v-if="cekAbsenMasukKaryawan(currentUser.user_id)">
                            <template v-if="cekAbsenPulangKaryawan(currentUser.user_id)">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                                SELESAI
                            </template>
                            <template v-else>
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 animate-pulse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                AKTIF
                            </template>
                        </template>
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                            BELUM
                        </template>
                    </div>
                </div>
            </div>

            <h2 class="text-lg font-black text-slate-800 mb-4 flex items-center gap-2 uppercase tracking-tight">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                Panel Absensi
            </h2>

            <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 mb-10">
                <div v-for="user in karyawan" :key="user.id"
                     class="bg-white rounded-3xl p-5 shadow-sm border-2 transition-all duration-300 relative overflow-hidden group"
                     :class="user.id == currentUser.user_id ? 'border-blue-500 shadow-blue-100 ring-4 ring-blue-50' : 'border-slate-100 opacity-60'">
                    
                    <div class="flex items-center gap-4 mb-5 pb-4 border-b border-slate-50 relative z-10">
                        <div class="w-14 h-14 rounded-[16px] flex items-center justify-center shrink-0 border-2 border-white shadow-md bg-slate-900 text-white font-black text-lg uppercase">{{ user.name.substring(0, 2) }}</div>
                        <div>
                            <h3 class="font-black text-lg leading-tight text-slate-800 uppercase">{{ user.name }}</h3>
                            <div class="text-[9px] font-black px-2 py-1 mt-1 rounded bg-slate-100 text-slate-500 uppercase tracking-widest inline-block">{{ user.role }} • {{ user.nik || 'ID:'+user.id }}</div>
                        </div>
                    </div>

                    <div class="flex gap-3 relative z-10">
                        <button v-if="!cekAbsenMasukKaryawan(user.id)" @click="mulaiAbsen(user.id, user.name, 'Masuk')"
                            :disabled="user.id != currentUser.user_id"
                            class="flex-1 py-3.5 rounded-[18px] font-black text-[11px] uppercase tracking-widest transition-all border-2 flex items-center justify-center gap-2 bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-600 hover:text-white disabled:opacity-30">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/></svg>
                            Absen Masuk
                        </button>
                        
                        <button v-if="cekAbsenMasukKaryawan(user.id) && !cekAbsenPulangKaryawan(user.id)" @click="mulaiAbsen(user.id, user.name, 'Pulang')"
                            :disabled="user.id != currentUser.user_id"
                            class="flex-1 py-3.5 rounded-[18px] font-black text-[11px] uppercase tracking-widest transition-all border-2 flex items-center justify-center gap-2 bg-amber-50 text-amber-700 border-amber-200 hover:bg-amber-500 hover:text-white disabled:opacity-30">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                            Absen Pulang
                        </button>

                        <div v-if="cekAbsenPulangKaryawan(user.id)" class="w-full py-3.5 text-center bg-slate-50 text-slate-400 font-black rounded-[18px] text-[10px] uppercase tracking-widest border border-slate-200 flex items-center justify-center gap-2">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                            Selesai Hari Ini
                        </div>
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-3xl shadow-xl border border-slate-100 overflow-hidden">
                <div class="p-5 md:p-6 border-b border-slate-100 flex flex-col md:flex-row md:justify-between md:items-center bg-slate-50/50 gap-4">
                    <h3 class="font-black text-slate-800 text-lg flex items-center gap-2 uppercase tracking-tight">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
                        Log Riwayat
                    </h3>
                    <div class="flex flex-wrap items-center gap-3">
                        <select v-model="filterMode" class="px-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase tracking-widest shadow-sm">
                            <option value="harian">Harian</option><option value="bulanan">Bulanan</option>
                        </select>
                        <input v-if="filterMode === 'harian'" type="date" v-model="tanggalDipilih" class="px-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase shadow-sm">
                        <input v-else type="month" v-model="bulanDipilih" class="px-4 py-2 bg-white border border-slate-200 rounded-xl text-xs font-black text-slate-700 outline-none uppercase shadow-sm">
                        <button v-if="currentUser.role === 'owner'" @click="downloadLaporan" class="bg-slate-900 text-white px-5 py-2.5 rounded-xl font-black text-[10px] uppercase tracking-widest shadow-md flex items-center gap-2 hover:bg-blue-600 transition-colors">
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="3" y2="15"/></svg>
                            Ekspor CSV
                        </button>
                    </div>
                </div>
                
                <div class="overflow-x-auto custom-scrollbar">
                    <table class="w-full text-left whitespace-nowrap">
                        <thead class="bg-white border-b border-slate-100">
                            <tr class="text-[9px] font-black text-slate-400 uppercase tracking-[0.2em]">
                                <th class="px-6 py-5">Karyawan</th>
                                <th v-if="filterMode === 'bulanan'" @click="toggleSortTanggal" class="px-6 py-5 text-center cursor-pointer select-none hover:bg-slate-50 transition-colors group">
                                    <div class="flex items-center justify-center gap-2">
                                        Tanggal
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-blue-600 transition-transform" :class="!urutanTanggalTerbaru ? 'rotate-180' : ''" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m6 9 6 6 6-6"/></svg>
                                    </div>
                                </th>
                                <th class="px-6 py-5 text-center">Shift</th>
                                <th class="px-6 py-5 text-center">Foto Masuk</th>
                                <th class="px-6 py-5 text-center">Jam Masuk</th>
                                <th class="px-6 py-5 text-center">Foto Pulang</th>
                                <th class="px-6 py-5 text-center">Jam Pulang</th>
                                <th class="px-6 py-5 text-center">Status</th>
                            </tr>
                        </thead>
                        <tbody class="divide-y divide-slate-50">
                            <tr v-if="isLoading">
                                <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest animate-pulse">
                                    Sinkronisasi Data...
                                </td>
                            </tr>
                            <tr v-else-if="riwayat.length === 0">
                                <td :colspan="filterMode === 'bulanan' ? 8 : 7" class="px-6 py-12 text-center text-slate-400 font-black text-xs uppercase tracking-widest opacity-50">
                                    Data Riwayat Kosong
                                </td>
                            </tr>
                            <tr v-else v-for="log in riwayat" :key="log.id" class="hover:bg-slate-50/50 transition-colors group">
                                <td class="px-6 py-5">
                                    <div class="font-black text-slate-800 uppercase text-xs">{{ log.User?.name }}</div>
                                    <div class="text-[9px] text-slate-400 font-bold tracking-widest">{{ log.User?.nik || 'ID: '+log.user_id }}</div>
                                </td>
                                
                                <td v-if="filterMode === 'bulanan'" class="px-6 py-5 text-center font-mono text-[11px] font-black text-slate-600">
                                    {{ log.tanggal ? log.tanggal.substring(0, 10) : '-' }}
                                </td>

                                <td class="px-6 py-5 text-center font-black text-[10px] text-blue-600 uppercase tracking-widest">
                                    <span class="bg-blue-50 border border-blue-100 px-3 py-1 rounded-lg">{{ log.shift }}</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <div v-if="log.foto_masuk" @click="lihatFoto(log.foto_masuk, log.User?.name, 'Masuk', log.jam_masuk)" class="w-12 h-12 rounded-xl mx-auto border-2 border-slate-200 shadow-sm overflow-hidden cursor-zoom-in hover:border-blue-500 transition-colors">
                                        <img :src="log.foto_masuk" class="w-full h-full object-cover">
                                    </div>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <span v-if="log.jam_masuk" class="bg-emerald-50 text-emerald-700 border border-emerald-100 font-black px-3 py-1.5 rounded-lg text-[10px]">
                                        {{ log.jam_masuk }}
                                    </span>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <div v-if="log.foto_pulang" @click="lihatFoto(log.foto_pulang, log.User?.name, 'Pulang', log.jam_pulang)" class="w-12 h-12 rounded-xl mx-auto border-2 border-slate-200 shadow-sm overflow-hidden cursor-zoom-in hover:border-blue-500 transition-colors">
                                        <img :src="log.foto_pulang" class="w-full h-full object-cover">
                                    </div>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>

                                <td class="px-6 py-5 text-center">
                                    <span v-if="log.jam_pulang" class="bg-amber-50 text-amber-700 border border-amber-100 font-black px-3 py-1.5 rounded-lg text-[10px]">
                                        {{ log.jam_pulang }}
                                    </span>
                                    <span v-else class="text-slate-300 font-black text-xs">-</span>
                                </td>
                                
                                <td class="px-6 py-5 text-center">
                                    <span v-if="log.status === 'Hadir'" class="inline-flex items-center gap-1.5 bg-emerald-50 text-emerald-700 border border-emerald-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                                        Hadir
                                    </span>
                                    <span v-else-if="log.status === 'Lupa Absen Pulang'" class="inline-flex items-center gap-1.5 bg-amber-50 text-amber-700 border border-amber-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
                                        Lupa Pulang
                                    </span>
                                    <span v-else-if="log.status === 'Mangkir'" class="inline-flex items-center gap-1.5 bg-red-50 text-red-700 border border-red-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                                        Mangkir
                                    </span>
                                    <span v-else-if="log.status === 'Libur (OFF)'" class="inline-flex items-center gap-1.5 bg-slate-50 text-slate-500 border border-slate-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="m9 16 2 2 4-4"/></svg>
                                        OFF
                                    </span>
                                    <span v-else-if="log.status === 'Belum Absen'" class="inline-flex items-center gap-1.5 bg-blue-50 text-blue-700 border border-blue-200 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                                        Belum Absen
                                    </span>
                                    <span v-else class="bg-slate-100 text-slate-600 font-black px-3 py-1.5 rounded-full text-[9px] uppercase tracking-widest shadow-sm">
                                        {{ log.status }}
                                    </span>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div v-if="showCameraModal" class="fixed inset-0 bg-slate-950/95 z-[100] flex items-center justify-center p-4 backdrop-blur-xl">
            <div class="bg-white rounded-[40px] overflow-hidden shadow-2xl w-full max-w-sm border-[10px] border-slate-900/5">
                <div class="p-6 border-b border-slate-50 flex justify-between items-center bg-white">
                    <div>
                        <h3 class="font-black text-slate-800 uppercase tracking-tighter text-xl italic">Verification</h3>
                        <p class="text-[10px] text-blue-600 font-black uppercase tracking-widest mt-0.5">{{ absenTarget.jenis }} • {{ absenTarget.nama }}</p>
                    </div>
                    <button @click="stopCamera" class="w-10 h-10 rounded-[14px] bg-slate-100 text-slate-400 hover:text-red-500 hover:bg-red-50 transition-all flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                    </button>
                </div>
                <div class="relative bg-slate-900 w-full aspect-[3/4] flex items-center justify-center overflow-hidden">
                    <div v-if="isAiLoading" class="absolute inset-0 bg-slate-950/90 backdrop-blur-md flex flex-col items-center justify-center z-50 text-white">
                        <div class="w-12 h-12 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4"></div>
                        <p class="font-black text-[10px] uppercase tracking-[0.2em] animate-pulse">Initializing AI Engine...</p>
                    </div>
                    <video ref="videoRef" autoplay playsinline class="w-full h-full object-cover transform scale-x-[-1]"></video>
                    
                    <div class="absolute inset-0 border-[16px] border-black/20 pointer-events-none"></div>
                    <div class="absolute inset-x-12 inset-y-16 border-2 border-white/30 rounded-[40px] pointer-events-none">
                        <div class="absolute top-0 left-0 w-8 h-8 border-t-4 border-l-4 border-white/80 rounded-tl-[38px]"></div>
                        <div class="absolute top-0 right-0 w-8 h-8 border-t-4 border-r-4 border-white/80 rounded-tr-[38px]"></div>
                        <div class="absolute bottom-0 left-0 w-8 h-8 border-b-4 border-l-4 border-white/80 rounded-bl-[38px]"></div>
                        <div class="absolute bottom-0 right-0 w-8 h-8 border-b-4 border-r-4 border-white/80 rounded-br-[38px]"></div>
                    </div>

                    <div class="absolute bottom-6 left-6 text-yellow-400 font-mono text-xs drop-shadow-md">
                        <div class="font-black tracking-widest">{{ absenTarget.nama.toUpperCase() }}</div>
                        <div class="text-xl font-black tracking-widest mt-1">{{ currentTime }}</div>
                        <div class="font-bold opacity-80 mt-0.5">{{ new Date().toLocaleDateString('id-ID') }}</div>
                    </div>
                </div>
                <div class="p-6 bg-white">
                    <button @click="jepretDanKirim" :disabled="isSubmitting" class="w-full bg-blue-600 hover:bg-slate-900 text-white py-5 rounded-[24px] font-black text-sm uppercase tracking-[0.2em] shadow-xl shadow-blue-200 transition-all disabled:opacity-50 flex items-center justify-center gap-3 active:scale-95">
                        <template v-if="isSubmitting">
                            <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            Verifying...
                        </template>
                        <template v-else>
                            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
                            Capture & Verify
                        </template>
                    </button>
                </div>
            </div>
        </div>
    </Sidebar>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 5px; height: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>